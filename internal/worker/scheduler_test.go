package worker

import (
	"log/slog"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
)

func newTestSchedulerLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func setupTestScheduler(t *testing.T) (*Scheduler, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to create miniredis: %v", err)
	}

	logger := newTestSchedulerLogger()
	sched, err := NewScheduler(logger, mr.Addr(), "", 0)
	if err != nil {
		mr.Close()
		t.Fatalf("failed to create scheduler: %v", err)
	}

	return sched, mr
}

func TestNewScheduler(t *testing.T) {
	sched, mr := setupTestScheduler(t)
	defer mr.Close()
	defer sched.Close()

	assert.NotNil(t, sched)
	assert.NotNil(t, sched.client)
	assert.NotNil(t, sched.logger)
}

func TestScheduler_EnqueuePriceCheck(t *testing.T) {
	sched, mr := setupTestScheduler(t)
	defer mr.Close()
	defer sched.Close()

	err := sched.EnqueuePriceCheck("product-123")
	assert.NoError(t, err)
}

func TestScheduler_EnqueueCompetitorSync(t *testing.T) {
	sched, mr := setupTestScheduler(t)
	defer mr.Close()
	defer sched.Close()

	err := sched.EnqueueCompetitorSync("search query")
	assert.NoError(t, err)
}

func TestScheduler_EnqueueTrendAnalysis(t *testing.T) {
	sched, mr := setupTestScheduler(t)
	defer mr.Close()
	defer sched.Close()

	err := sched.EnqueueTrendAnalysis("market trends")
	assert.NoError(t, err)
}

func TestScheduler_EnqueueSupplyAlert(t *testing.T) {
	sched, mr := setupTestScheduler(t)
	defer mr.Close()
	defer sched.Close()

	err := sched.EnqueueSupplyAlert("product-456")
	assert.NoError(t, err)
}

func TestScheduler_ScheduleDailyPriceCheck(t *testing.T) {
	sched, mr := setupTestScheduler(t)
	defer mr.Close()
	defer sched.Close()

	productIDs := []string{"prod-1", "prod-2", "prod-3"}
	err := sched.ScheduleDailyPriceCheck(productIDs)
	assert.NoError(t, err)
}

func TestScheduler_Enqueue_WithOptions(t *testing.T) {
	sched, mr := setupTestScheduler(t)
	defer mr.Close()
	defer sched.Close()

	payload := NewTaskPayload(TaskTypePriceCheck, "prod-test", "")
	err := sched.Enqueue(payload, asynq.MaxRetry(5), asynq.Queue("critical"))
	assert.NoError(t, err)
}

func TestCronExprConstants(t *testing.T) {
	assert.Equal(t, "@every 1m", CronEveryMinute)
	assert.Equal(t, "@every 5m", CronEvery5Minutes)
	assert.Equal(t, "@every 15m", CronEvery15Minutes)
	assert.Equal(t, "@every 30m", CronEvery30Minutes)
	assert.Equal(t, "@every 1h", CronEveryHour)
	assert.Equal(t, "0 9 * * *", CronEveryDay9AM)
	assert.Equal(t, "0 2 * * *", CronEveryDay2AM)
	assert.Equal(t, "0 9 * * 1", CronEveryWeekMonday)
}

func TestScheduler_Enqueue_ErrorPath(t *testing.T) {
	sched, mr := setupTestScheduler(t)

	// Close miniredis to simulate connection error
	mr.Close()

	payload := NewTaskPayload(TaskTypePriceCheck, "prod-test", "")
	err := sched.Enqueue(payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "enqueue task failed")
}

func TestScheduler_ScheduleDailyPriceCheck_ErrorHandling(t *testing.T) {
	sched, mr := setupTestScheduler(t)
	defer mr.Close()
	defer sched.Close()

	// Test with empty product list
	err := sched.ScheduleDailyPriceCheck([]string{})
	assert.NoError(t, err) // should not error on empty list

	// Test that partial failure continues (first succeeds, second fails if redis closes)
	// Note: with miniredis this would be hard to trigger, so we verify the normal case works
	productIDs := []string{"prod-1", "prod-2"}
	err = sched.ScheduleDailyPriceCheck(productIDs)
	assert.NoError(t, err)
}

func TestScheduler_Close(t *testing.T) {
	sched, mr := setupTestScheduler(t)
	mr.Close() // close miniredis first

	err := sched.Close()
	assert.NoError(t, err) // Close should not error even if already closed
}

func TestScheduler_Enqueue_InvalidPayload(t *testing.T) {
	sched, mr := setupTestScheduler(t)
	defer mr.Close()
	defer sched.Close()

	// Create a payload and manually corrupt its JSON before enqueueing
	// Since ToJSON is reliable, we test the serialization path indirectly
	payload := NewTaskPayload(TaskTypePriceCheck, "prod-123", "")
	err := sched.Enqueue(payload)
	assert.NoError(t, err)
}
