package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"competitive-analysis-agent/internal/config"
	"competitive-analysis-agent/internal/llm"
	"competitive-analysis-agent/internal/scraper"
	"competitive-analysis-agent/internal/scraper/platforms"
	"competitive-analysis-agent/internal/storage"
	"competitive-analysis-agent/internal/storage/repository"
	"competitive-analysis-agent/internal/worker"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

func main() {
	// 命令行参数
	configPath := flag.String("config", "configs/development.yaml", "path to config file")
	flag.Parse()

	// 初始化日志
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Error("Failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// 初始化 Redis 客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Asynq.RedisHost,
		Password: cfg.Asynq.RedisPassword,
		DB:       cfg.Asynq.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Error("Failed to connect to Redis", slog.String("error", err.Error()))
		os.Exit(1)
	}
	cancel()
	logger.Info("Connected to Redis", slog.String("addr", cfg.Asynq.RedisHost))

	// 初始化数据库（用于存储爬取结果）
	db, err := storage.NewPostgresDB(&cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to PostgreSQL", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("Connected to PostgreSQL", slog.String("host", cfg.Database.Host), slog.Int("port", cfg.Database.Port))

	// 初始化 LLM 客户端
	llmClient := llm.NewOpenAICompatProvider(
		cfg.LLM.APIKey,
		cfg.LLM.BaseURL,
		"qwen-max",
		cfg.LLM.Provider,
	)
	if cfg.LLM.Provider == "qwen" {
		llmClient.SetAPIPath("/v1/chat/completions")
	}
	logger.Info("LLM client initialized", slog.String("provider", cfg.LLM.Provider))

	// 初始化 Repository
	compRepo := repository.NewCompetitorRepository(db)

	// 初始化 DataCleaner
	dataCleaner := scraper.NewDataCleaner()

	// 初始化多平台爬虫
	var multiScraper *platforms.MultiPlatformScraper
	{
		var amazon *platforms.AmazonScraper
		var aliexpress *platforms.AliExpressScraper
		var ebay *platforms.EbayScraper
		var temu *platforms.TemuScraper
		var shopify *platforms.ShopifyScraper

		proxyPool := scraper.NewProxyPool(nil) // 生产环境配置真实代理
		proxyList := proxyPool.GetAll()

		// 初始化各平台爬虫（仅当浏览器可用时）
		if s, err := platforms.NewAmazonScraper(proxyList); err == nil {
			amazon = s
			defer s.Close()
		} else {
			logger.Warn("Amazon scraper unavailable", slog.String("error", err.Error()))
		}

		if s, err := platforms.NewAliExpressScraper(proxyList); err == nil {
			aliexpress = s
			defer s.Close()
		} else {
			logger.Warn("AliExpress scraper unavailable", slog.String("error", err.Error()))
		}

		if s, err := platforms.NewEbayScraper(proxyList); err == nil {
			ebay = s
			defer s.Close()
		} else {
			logger.Warn("eBay scraper unavailable", slog.String("error", err.Error()))
		}

		if s, err := platforms.NewTemuScraper(proxyList); err == nil {
			temu = s
			defer s.Close()
		} else {
			logger.Warn("Temu scraper unavailable", slog.String("error", err.Error()))
		}

		// Shopify 使用 Colly（无需浏览器），从配置文件读取 stores
		shopifyStores := cfg.ShopifyStores
		if len(shopifyStores) == 0 {
			shopifyStores = []string{} // 使用空切片，DiscoverStores 可后续添加
		}
		shopify = platforms.NewShopifyScraper(shopifyStores)

		multiScraper = platforms.NewMultiPlatformScraper(
			amazon, aliexpress, ebay, temu, shopify, dataCleaner,
		)
		logger.Info("Multi-platform scraper initialized",
			slog.Bool("amazon", amazon != nil),
			slog.Bool("aliexpress", aliexpress != nil),
			slog.Bool("ebay", ebay != nil),
			slog.Bool("temu", temu != nil),
			slog.Bool("shopify", shopify != nil),
		)
	}

	// 初始化 Worker 依赖适配器
	multiAdapter := worker.NewMultiPlatformAdapter(multiScraper, dataCleaner, compRepo)
	llmAdapter := worker.NewLLMAnalyzerAdapter(llmClient)

	// 初始化 Handler
	handler := worker.NewHandler(logger, multiAdapter, llmAdapter)

	// 初始化 Asynq Server
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.Asynq.RedisHost,
			Password: cfg.Asynq.RedisPassword,
			DB:       cfg.Asynq.RedisDB,
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				logger.Error("Task failed",
					slog.String("type", task.Type()),
					slog.String("error", err.Error()),
				)
			}),
		},
	)

	// 注册处理器
	mux := asynq.NewServeMux()
	mux.Handle(string(worker.TaskTypePriceCheck), worker.NewTaskHandler(logger, handler))
	mux.Handle(string(worker.TaskTypeCompetitorSync), worker.NewTaskHandler(logger, handler))
	mux.Handle(string(worker.TaskTypeTrendAnalysis), worker.NewTaskHandler(logger, handler))
	mux.Handle(string(worker.TaskTypeSupplyAlert), worker.NewTaskHandler(logger, handler))

	// 启动 Worker
	logger.Info("Starting Asynq worker...")

	errChan := make(chan error, 1)
	go func() {
		if err := srv.Start(mux); err != nil {
			errChan <- err
		}
	}()

	// 等待信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		logger.Info("Received signal, shutting down...", slog.String("signal", sig.String()))
	case err := <-errChan:
		logger.Error("Worker error", slog.String("error", err.Error()))
	}

	// 优雅关闭
	multiScraper.Close()
	srv.Shutdown()
	redisClient.Close()
	logger.Info("Worker stopped")
}
