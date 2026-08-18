package worker

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// TaskType 任务类型
type TaskType string

const (
	TaskTypePriceCheck    TaskType = "price_check"     // 每日价格检查
	TaskTypeCompetitorSync TaskType = "competitor_sync" // 竞品数据同步
	TaskTypeTrendAnalysis TaskType = "trend_analysis"  // 趋势分析
	TaskTypeSupplyAlert   TaskType = "supply_alert"    // 供应链预警
)

// TaskPayload 任务载荷
type TaskPayload struct {
	TaskID    string    `json:"task_id"`
	Type      TaskType  `json:"type"`
	ProductID string    `json:"product_id,omitempty"`
	Query     string    `json:"query,omitempty"`
	Schedule  string    `json:"schedule,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// NewTaskPayload 创建任务载荷
func NewTaskPayload(taskType TaskType, productID, query string) *TaskPayload {
	return &TaskPayload{
		TaskID:    uuid.New().String(),
		Type:      taskType,
		ProductID: productID,
		Query:     query,
		CreatedAt: time.Now(),
	}
}

// ToJSON 序列化为JSON
func (p *TaskPayload) ToJSON() (string, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ParseTaskPayload 解析任务载荷
func ParseTaskPayload(data string) (*TaskPayload, error) {
	var payload TaskPayload
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

// RetryPolicy 重试策略
type RetryPolicy struct {
	MaxRetry  int           // 最大重试次数
	Initial   time.Duration // 初始间隔
	MaxDelay  time.Duration // 最大间隔
	Multiplier float64      // 间隔倍数
}

// DefaultRetryPolicy 默认重试策略
var DefaultRetryPolicy = RetryPolicy{
	MaxRetry:   3,
	Initial:    10 * time.Second,
	MaxDelay:   10 * time.Minute,
	Multiplier: 2.0,
}
