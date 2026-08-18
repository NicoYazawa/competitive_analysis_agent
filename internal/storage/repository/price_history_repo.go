package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"competitive-analysis-agent/internal/domain/entity"
	"competitive-analysis-agent/internal/storage"

	"github.com/google/uuid"
)

// PriceHistoryRepository 价格历史仓储
type PriceHistoryRepository struct {
	db *storage.PostgresDB
}

// NewPriceHistoryRepository 创建价格历史仓储
func NewPriceHistoryRepository(db *storage.PostgresDB) *PriceHistoryRepository {
	return &PriceHistoryRepository{db: db}
}

// Create 创建价格记录
func (r *PriceHistoryRepository) Create(ctx context.Context, p *entity.PriceHistory) error {
	if err := p.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	if p.RecordedAt.IsZero() {
		p.RecordedAt = time.Now()
	}

	query := `
		INSERT INTO price_history (id, competitor_id, price, currency, recorded_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		p.ID, p.CompetitorID, p.Price, p.Currency, p.RecordedAt,
	)

	if err != nil {
		return fmt.Errorf("insert price history error: %w", err)
	}

	return nil
}

// GetByID 根据 ID 获取
func (r *PriceHistoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.PriceHistory, error) {
	query := `
		SELECT id, competitor_id, price, currency, recorded_at
		FROM price_history
		WHERE id = $1
	`

	var p entity.PriceHistory
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.CompetitorID, &p.Price, &p.Currency, &p.RecordedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query price history error: %w", err)
	}

	return &p, nil
}

// GetByCompetitorID 获取竞品的所有价格历史
func (r *PriceHistoryRepository) GetByCompetitorID(ctx context.Context, competitorID uuid.UUID, limit int) ([]*entity.PriceHistory, error) {
	query := `
		SELECT id, competitor_id, price, currency, recorded_at
		FROM price_history
		WHERE competitor_id = $1
		ORDER BY recorded_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, competitorID, limit)
	if err != nil {
		return nil, fmt.Errorf("query price history error: %w", err)
	}
	defer rows.Close()

	var history []*entity.PriceHistory
	for rows.Next() {
		var p entity.PriceHistory
		err := rows.Scan(
			&p.ID, &p.CompetitorID, &p.Price, &p.Currency, &p.RecordedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan price history error: %w", err)
		}
		history = append(history, &p)
	}

	return history, nil
}

// GetLatest 获取最新价格
func (r *PriceHistoryRepository) GetLatest(ctx context.Context, competitorID uuid.UUID) (*entity.PriceHistory, error) {
	query := `
		SELECT id, competitor_id, price, currency, recorded_at
		FROM price_history
		WHERE competitor_id = $1
		ORDER BY recorded_at DESC
		LIMIT 1
	`

	var p entity.PriceHistory
	err := r.db.QueryRow(ctx, query, competitorID).Scan(
		&p.ID, &p.CompetitorID, &p.Price, &p.Currency, &p.RecordedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query price history error: %w", err)
	}

	return &p, nil
}

// GetLatestPrices 获取多个竞品的最新价格
func (r *PriceHistoryRepository) GetLatestPrices(ctx context.Context, competitorIDs []uuid.UUID) (map[uuid.UUID]*entity.PriceHistory, error) {
	if len(competitorIDs) == 0 {
		return make(map[uuid.UUID]*entity.PriceHistory), nil
	}

	query := `
		SELECT DISTINCT ON (competitor_id) id, competitor_id, price, currency, recorded_at
		FROM price_history
		WHERE competitor_id = ANY($1)
		ORDER BY competitor_id, recorded_at DESC
	`

	rows, err := r.db.Query(ctx, query, competitorIDs)
	if err != nil {
		return nil, fmt.Errorf("query price history error: %w", err)
	}
	defer rows.Close()

	result := make(map[uuid.UUID]*entity.PriceHistory)
	for rows.Next() {
		var p entity.PriceHistory
		err := rows.Scan(
			&p.ID, &p.CompetitorID, &p.Price, &p.Currency, &p.RecordedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan price history error: %w", err)
		}
		result[p.CompetitorID] = &p
	}

	return result, nil
}

// GetPriceRange 获取价格区间
func (r *PriceHistoryRepository) GetPriceRange(ctx context.Context, competitorID uuid.UUID, startTime, endTime time.Time) ([]*entity.PriceHistory, error) {
	query := `
		SELECT id, competitor_id, price, currency, recorded_at
		FROM price_history
		WHERE competitor_id = $1 AND recorded_at >= $2 AND recorded_at <= $3
		ORDER BY recorded_at ASC
	`

	rows, err := r.db.Query(ctx, query, competitorID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("query price history error: %w", err)
	}
	defer rows.Close()

	var history []*entity.PriceHistory
	for rows.Next() {
		var p entity.PriceHistory
		err := rows.Scan(
			&p.ID, &p.CompetitorID, &p.Price, &p.Currency, &p.RecordedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan price history error: %w", err)
		}
		history = append(history, &p)
	}

	return history, nil
}

// Delete 删除价格记录
func (r *PriceHistoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM price_history WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete price history error: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("price history not found")
	}

	return nil
}

// DeleteByCompetitorID 删除竞品的所有价格历史
func (r *PriceHistoryRepository) DeleteByCompetitorID(ctx context.Context, competitorID uuid.UUID) error {
	query := `DELETE FROM price_history WHERE competitor_id = $1`

	_, err := r.db.Exec(ctx, query, competitorID)
	if err != nil {
		return fmt.Errorf("delete price history error: %w", err)
	}

	return nil
}

// Count 获取记录数量
func (r *PriceHistoryRepository) Count(ctx context.Context, competitorID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM price_history WHERE competitor_id = $1`

	var count int
	err := r.db.QueryRow(ctx, query, competitorID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count price history error: %w", err)
	}

	return count, nil
}

// DetectPriceChange 检测价格变化
func (r *PriceHistoryRepository) DetectPriceChange(ctx context.Context, competitorID uuid.UUID, threshold float64) (*PriceChange, error) {
	// 获取最近两条记录
	query := `
		SELECT id, competitor_id, price, currency, recorded_at
		FROM price_history
		WHERE competitor_id = $1
		ORDER BY recorded_at DESC
		LIMIT 2
	`

	rows, err := r.db.Query(ctx, query, competitorID)
	if err != nil {
		return nil, fmt.Errorf("query price history error: %w", err)
	}
	defer rows.Close()

	var latest, previous *entity.PriceHistory
	count := 0
	for rows.Next() {
		var p entity.PriceHistory
		err := rows.Scan(
			&p.ID, &p.CompetitorID, &p.Price, &p.Currency, &p.RecordedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan price history error: %w", err)
		}
		count++
		if count == 1 {
			latest = &p
		} else if count == 2 {
			previous = &p
		}
	}

	if latest == nil {
		return nil, nil
	}

	if previous == nil {
		return &PriceChange{
			CompetitorID: competitorID,
			OldPrice:     0,
			NewPrice:     latest.Price,
			Change:       0,
			ChangePercent: 0,
			Direction:    "new",
		}, nil
	}

	change := latest.Price - previous.Price
	changePercent := 0.0
	if previous.Price > 0 {
		changePercent = (change / previous.Price) * 100
	}

	direction := "unchanged"
	if change > threshold {
		direction = "increased"
	} else if change < -threshold {
		direction = "decreased"
	}

	return &PriceChange{
		CompetitorID:    competitorID,
		OldPrice:        previous.Price,
		NewPrice:        latest.Price,
		Change:          change,
		ChangePercent:   changePercent,
		Direction:       direction,
		OldRecordedAt:   previous.RecordedAt,
		NewRecordedAt:   latest.RecordedAt,
	}, nil
}

// PriceChange 价格变化信息
type PriceChange struct {
	CompetitorID    uuid.UUID
	OldPrice        float64
	NewPrice        float64
	Change          float64
	ChangePercent   float64
	Direction       string // "increased", "decreased", "unchanged", "new"
	OldRecordedAt   time.Time
	NewRecordedAt   time.Time
}

// GetAveragePrice 获取平均价格
func (r *PriceHistoryRepository) GetAveragePrice(ctx context.Context, competitorID uuid.UUID, days int) (float64, error) {
	query := `
		SELECT AVG(price)
		FROM price_history
		WHERE competitor_id = $1 AND recorded_at >= NOW() - INTERVAL '1 day' * $2
	`

	var avg sql.NullFloat64
	err := r.db.QueryRow(ctx, query, competitorID, days).Scan(&avg)
	if err != nil {
		return 0, fmt.Errorf("query average price error: %w", err)
	}

	if !avg.Valid {
		return 0, nil
	}

	return avg.Float64, nil
}

// GetMinMaxPrice 获取最低最高价
func (r *PriceHistoryRepository) GetMinMaxPrice(ctx context.Context, competitorID uuid.UUID, days int) (min, max float64, err error) {
	query := `
		SELECT MIN(price), MAX(price)
		FROM price_history
		WHERE competitor_id = $1 AND recorded_at >= NOW() - INTERVAL '1 day' * $2
	`

	err = r.db.QueryRow(ctx, query, competitorID, days).Scan(&min, &max)
	if err != nil {
		return 0, 0, fmt.Errorf("query min max price error: %w", err)
	}

	return min, max, nil
}
