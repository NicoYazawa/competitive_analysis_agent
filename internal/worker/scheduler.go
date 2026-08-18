package worker

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/hibiken/asynq"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	client *asynq.Client
	logger *slog.Logger
}

// NewScheduler 创建调度器
func NewScheduler(logger *slog.Logger, redisAddr, redisPassword string, redisDB int) (*Scheduler, error) {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	return &Scheduler{
		client: client,
		logger: logger,
	}, nil
}

// Close 关闭调度器
func (s *Scheduler) Close() error {
	return s.client.Close()
}

// Enqueue 入队任务
func (s *Scheduler) Enqueue(payload *TaskPayload, opts ...asynq.Option) error {
	payloadJSON, err := payload.ToJSON()
	if err != nil {
		return fmt.Errorf("serialize payload failed: %w", err)
	}

	task := asynq.NewTask(string(payload.Type), []byte(payloadJSON))
	info, err := s.client.Enqueue(task, opts...)
	if err != nil {
		return fmt.Errorf("enqueue task failed: %w", err)
	}

	s.logger.Info("Task enqueued",
		slog.String("task_id", payload.TaskID),
		slog.String("type", string(payload.Type)),
		slog.String("queue", info.Queue),
		slog.Int("retry", info.Retried))

	return nil
}

// EnqueuePriceCheck 入队价格检查任务
func (s *Scheduler) EnqueuePriceCheck(productID string) error {
	payload := NewTaskPayload(TaskTypePriceCheck, productID, "")

	opts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.Timeout(10 * time.Minute),
		asynq.Queue("critical"),
	}

	return s.Enqueue(payload, opts...)
}

// EnqueueCompetitorSync 入队竞品同步任务
func (s *Scheduler) EnqueueCompetitorSync(query string) error {
	payload := NewTaskPayload(TaskTypeCompetitorSync, "", query)

	opts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.Timeout(15 * time.Minute),
		asynq.Queue("default"),
	}

	return s.Enqueue(payload, opts...)
}

// EnqueueTrendAnalysis 入队趋势分析任务
func (s *Scheduler) EnqueueTrendAnalysis(query string) error {
	payload := NewTaskPayload(TaskTypeTrendAnalysis, "", query)

	opts := []asynq.Option{
		asynq.MaxRetry(2),
		asynq.Timeout(5 * time.Minute),
		asynq.Queue("low"),
	}

	return s.Enqueue(payload, opts...)
}

// EnqueueSupplyAlert 入队供应链预警任务
func (s *Scheduler) EnqueueSupplyAlert(productID string) error {
	payload := NewTaskPayload(TaskTypeSupplyAlert, productID, "")

	opts := []asynq.Option{
		asynq.MaxRetry(5),
		asynq.Timeout(2 * time.Minute),
		asynq.Queue("critical"),
	}

	return s.Enqueue(payload, opts...)
}

// ScheduleDailyPriceCheck 调度每日价格检查
func (s *Scheduler) ScheduleDailyPriceCheck(productIDs []string) error {
	for _, productID := range productIDs {
		if err := s.EnqueuePriceCheck(productID); err != nil {
			s.logger.Error("Failed to schedule daily price check",
				slog.String("product_id", productID),
				slog.String("error", err.Error()))
			continue
		}
	}

	s.logger.Info("Daily price check scheduled",
		slog.Int("count", len(productIDs)))

	return nil
}

// CronExpr cron表达式帮助
type CronExpr struct {
	EveryMinute     string
	Every5Minutes   string
	Every15Minutes  string
	Every30Minutes  string
	EveryHour       string
	EveryDay9AM     string
	EveryDay2AM     string
	EveryWeekMonday string
}

// CronExpr constants
var CronEveryMinute = "@every 1m"
var CronEvery5Minutes = "@every 5m"
var CronEvery15Minutes = "@every 15m"
var CronEvery30Minutes = "@every 30m"
var CronEveryHour = "@every 1h"
var CronEveryDay9AM = "0 9 * * *"
var CronEveryDay2AM = "0 2 * * *"
var CronEveryWeekMonday = "0 9 * * 1"
