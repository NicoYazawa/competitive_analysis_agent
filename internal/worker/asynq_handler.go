package worker

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/hibiken/asynq"
)

// TaskHandler Asynq任务处理器包装器
type TaskHandler struct {
	logger   *slog.Logger
	inner    *Handler
}

// NewTaskHandler 创建TaskHandler
func NewTaskHandler(logger *slog.Logger, inner *Handler) *TaskHandler {
	return &TaskHandler{
		logger: logger,
		inner:  inner,
	}
}

// ProcessTask 实现asynq.Handler接口
func (h *TaskHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var payload TaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		h.logger.Error("Failed to parse task payload",
			slog.String("error", err.Error()),
			slog.String("type", task.Type()))
		return err
	}

	h.logger.Info("Processing task",
		slog.String("task_id", payload.TaskID),
		slog.String("type", string(payload.Type)))

	return h.inner.ProcessTask(ctx, task)
}
