package supervisor

import (
	"context"
	"fmt"
	"strings"
)

// Scheduler decomposes user queries into subtasks and routes them.
type Scheduler struct {
	agents map[TaskType]Agent
}

// NewScheduler creates a new Scheduler.
func NewScheduler() *Scheduler {
	return &Scheduler{
		agents: make(map[TaskType]Agent),
	}
}

// RegisterAgent registers an agent for a specific task type.
func (s *Scheduler) RegisterAgent(taskType TaskType, agent Agent) {
	s.agents[taskType] = agent
}

// DecomposeQuery breaks down a user query into subtasks.
func (s *Scheduler) DecomposeQuery(ctx context.Context, query string) ([]*Task, error) {
	var tasks []*Task

	lowerQuery := strings.ToLower(query)

	// Analyze query to determine required tasks
	if containsAny(lowerQuery, []string{"trend", "market", "demand", "趋势", "市场"}) {
		tasks = append(tasks, &Task{
			Type:        TaskTypeMarketTrend,
			Description: "Analyze market trends and demand signals",
			Query:       query,
			Priority:    1,
		})
	}

	if containsAny(lowerQuery, []string{"competitor", "competitors", "竞品", "竞争"}) {
		tasks = append(tasks, &Task{
			Type:        TaskTypeCompetitor,
			Description: "Analyze competitor pricing and positioning",
			Query:       query,
			Priority:    2,
		})
	}

	if containsAny(lowerQuery, []string{"price", "pricing", "定价", "价格"}) {
		tasks = append(tasks, &Task{
			Type:        TaskTypePricing,
			Description: "Develop pricing strategy",
			Query:       query,
			Priority:    3,
		})
	}

	if containsAny(lowerQuery, []string{"supply", "chain", "供应链", "物流"}) {
		tasks = append(tasks, &Task{
			Type:        TaskTypeSupplyChain,
			Description: "Analyze supply chain and logistics",
			Query:       query,
			Priority:    4,
		})
	}

	// If no specific tasks identified, default to market trend
	if len(tasks) == 0 {
		tasks = append(tasks, &Task{
			Type:        TaskTypeMarketTrend,
			Description: "General market analysis",
			Query:       query,
			Priority:    1,
		})
	}

	return tasks, nil
}

// ScheduleAndExecute decomposes query and executes all tasks via registered agents.
func (s *Scheduler) ScheduleAndExecute(ctx context.Context, query string) ([]*Task, error) {
	tasks, err := s.DecomposeQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to decompose query: %w", err)
	}

	for _, task := range tasks {
		agent, ok := s.agents[task.Type]
		if !ok {
			task.Result = fmt.Errorf("no agent registered for task type: %s", task.Type)
			continue
		}

		result, err := agent.Execute(ctx, task)
		if err != nil {
			task.Result = fmt.Errorf("agent %s failed: %w", agent.Name(), err)
			continue
		}
		task.Result = result
	}

	return tasks, nil
}

func containsAny(s string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.Contains(s, kw) {
			return true
		}
	}
	return false
}
