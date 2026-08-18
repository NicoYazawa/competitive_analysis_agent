package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"competitive-analysis-agent/internal/agents"
	"competitive-analysis-agent/internal/api"
	"competitive-analysis-agent/internal/api/handlers"
	"competitive-analysis-agent/internal/config"
	"competitive-analysis-agent/internal/llm"
	"competitive-analysis-agent/internal/storage"
	"competitive-analysis-agent/internal/storage/repository"
	"competitive-analysis-agent/internal/supervisor"
)

func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Determine config path
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	configPath := fmt.Sprintf("configs/%s.yaml", env)

	cfg, err := config.Load(configPath)
	if err != nil {
		logger.Error("Failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Initialize PostgreSQL
	db, err := storage.NewPostgresDB(&cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to PostgreSQL", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("Connected to PostgreSQL", slog.String("host", cfg.Database.Host), slog.Int("port", cfg.Database.Port))

	// Initialize Redis
	cache, err := storage.NewRedisCache(&cfg.Redis)
	if err != nil {
		logger.Error("Failed to connect to Redis", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer cache.Close()
	logger.Info("Connected to Redis", slog.String("host", cfg.Redis.Host), slog.Int("port", cfg.Redis.Port))

	// Initialize repositories
	compRepo := repository.NewCompetitorRepository(db)
	priceRepo := repository.NewPriceHistoryRepository(db)

	// Initialize LLM client
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

	// Initialize agents
	marketAgent := agents.NewMarketAgent(llmClient)
	competitorAgent := agents.NewCompetitorAgent(llmClient)
	supplyChainAgent := agents.NewSupplyChainAgent(llmClient)

	// Initialize supervisor and register agents
	scheduler := supervisor.NewScheduler()
	scheduler.RegisterAgent(supervisor.TaskTypeMarketTrend, marketAgent)
	scheduler.RegisterAgent(supervisor.TaskTypeCompetitor, competitorAgent)
	scheduler.RegisterAgent(supervisor.TaskTypeSupplyChain, supplyChainAgent)

	aggregator := supervisor.NewAggregator()

	// Initialize handlers
	competitorHandler := handlers.NewCompetitorHandler(compRepo, priceRepo)
	trendHandler := handlers.NewTrendHandler(nil, scheduler, aggregator) // marketAgent=nil until fully wired
	pricingHandler := handlers.NewPricingHandler()
	productHandler := handlers.NewProductHandler()

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	logger.Info("Starting API server", slog.String("addr", addr))
	logger.Info("Health check: http://" + addr + "/api/health")

	// Graceful shutdown
	srv := &http.Server{Addr: addr, Handler: api.NewRouter(competitorHandler, trendHandler, pricingHandler, productHandler)}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	sig := <-sigChan
	logger.Info("Received signal, shutting down...", slog.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Shutdown error", slog.String("error", err.Error()))
	}

	logger.Info("Server stopped")
}
