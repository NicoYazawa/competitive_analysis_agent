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

// CompetitorRepository 竞品仓储
type CompetitorRepository struct {
	db *storage.PostgresDB
}

// NewCompetitorRepository 创建竞品仓储
func NewCompetitorRepository(db *storage.PostgresDB) *CompetitorRepository {
	return &CompetitorRepository{db: db}
}

// Create 创建竞品
func (r *CompetitorRepository) Create(ctx context.Context, c *entity.Competitor) error {
	if err := c.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	query := `
		INSERT INTO competitors (id, name, platform, platform_product_id, current_price,
			currency, rating, review_count, seller_rating, seller_review_count, source_url,
			created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.Exec(ctx, query,
		c.ID, c.Name, c.Platform, c.PlatformProductID, c.CurrentPrice,
		c.Currency, c.Rating, c.ReviewCount, c.SellerRating, c.SellerReviewCount,
		c.SourceURL, c.CreatedAt, c.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("insert competitor error: %w", err)
	}

	return nil
}

// GetByID 根据 ID 获取
func (r *CompetitorRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Competitor, error) {
	query := `
		SELECT id, name, platform, platform_product_id, current_price, currency,
			rating, review_count, seller_rating, seller_review_count, source_url,
			created_at, updated_at
		FROM competitors
		WHERE id = $1
	`

	var c entity.Competitor
	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Platform, &c.PlatformProductID, &c.CurrentPrice,
		&c.Currency, &c.Rating, &c.ReviewCount, &c.SellerRating, &c.SellerReviewCount,
		&c.SourceURL, &c.CreatedAt, &c.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query competitor error: %w", err)
	}

	return &c, nil
}

// GetByPlatformAndProductID 根据平台和产品 ID 获取
func (r *CompetitorRepository) GetByPlatformAndProductID(ctx context.Context, platform, productID string) (*entity.Competitor, error) {
	query := `
		SELECT id, name, platform, platform_product_id, current_price, currency,
			rating, review_count, seller_rating, seller_review_count, source_url,
			created_at, updated_at
		FROM competitors
		WHERE platform = $1 AND platform_product_id = $2
	`

	var c entity.Competitor
	err := r.db.QueryRow(ctx, query, platform, productID).Scan(
		&c.ID, &c.Name, &c.Platform, &c.PlatformProductID, &c.CurrentPrice,
		&c.Currency, &c.Rating, &c.ReviewCount, &c.SellerRating, &c.SellerReviewCount,
		&c.SourceURL, &c.CreatedAt, &c.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query competitor error: %w", err)
	}

	return &c, nil
}

// Update 更新竞品
func (r *CompetitorRepository) Update(ctx context.Context, c *entity.Competitor) error {
	if err := c.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	c.UpdatedAt = time.Now()

	query := `
		UPDATE competitors
		SET name = $2, platform = $3, platform_product_id = $4, current_price = $5,
			currency = $6, rating = $7, review_count = $8, seller_rating = $9,
			seller_review_count = $10, source_url = $11, updated_at = $12
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query,
		c.ID, c.Name, c.Platform, c.PlatformProductID, c.CurrentPrice,
		c.Currency, c.Rating, c.ReviewCount, c.SellerRating, c.SellerReviewCount,
		c.SourceURL, c.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("update competitor error: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("competitor not found")
	}

	return nil
}

// Delete 删除竞品
func (r *CompetitorRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM competitors WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete competitor error: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("competitor not found")
	}

	return nil
}

// ListByPlatform 列出平台下的所有竞品
func (r *CompetitorRepository) ListByPlatform(ctx context.Context, platform string, limit, offset int) ([]*entity.Competitor, error) {
	query := `
		SELECT id, name, platform, platform_product_id, current_price, currency,
			rating, review_count, seller_rating, seller_review_count, source_url,
			created_at, updated_at
		FROM competitors
		WHERE platform = $1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, platform, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query competitors error: %w", err)
	}
	defer rows.Close()

	var competitors []*entity.Competitor
	for rows.Next() {
		var c entity.Competitor
		err := rows.Scan(
			&c.ID, &c.Name, &c.Platform, &c.PlatformProductID, &c.CurrentPrice,
			&c.Currency, &c.Rating, &c.ReviewCount, &c.SellerRating, &c.SellerReviewCount,
			&c.SourceURL, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan competitor error: %w", err)
		}
		competitors = append(competitors, &c)
	}

	return competitors, nil
}

// ListAll 列出所有竞品
func (r *CompetitorRepository) ListAll(ctx context.Context, limit, offset int) ([]*entity.Competitor, error) {
	query := `
		SELECT id, name, platform, platform_product_id, current_price, currency,
			rating, review_count, seller_rating, seller_review_count, source_url,
			created_at, updated_at
		FROM competitors
		ORDER BY updated_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query competitors error: %w", err)
	}
	defer rows.Close()

	var competitors []*entity.Competitor
	for rows.Next() {
		var c entity.Competitor
		err := rows.Scan(
			&c.ID, &c.Name, &c.Platform, &c.PlatformProductID, &c.CurrentPrice,
			&c.Currency, &c.Rating, &c.ReviewCount, &c.SellerRating, &c.SellerReviewCount,
			&c.SourceURL, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan competitor error: %w", err)
		}
		competitors = append(competitors, &c)
	}

	return competitors, nil
}

// CountByPlatform 统计平台下竞品数量
func (r *CompetitorRepository) CountByPlatform(ctx context.Context, platform string) (int, error) {
	query := `SELECT COUNT(*) FROM competitors WHERE platform = $1`

	var count int
	err := r.db.QueryRow(ctx, query, platform).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count competitors error: %w", err)
	}

	return count, nil
}

// UpdatePrice 更新竞品价格
func (r *CompetitorRepository) UpdatePrice(ctx context.Context, id uuid.UUID, price float64, currency string) error {
	query := `
		UPDATE competitors
		SET current_price = $2, currency = $3, updated_at = $4
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id, price, currency, time.Now())
	if err != nil {
		return fmt.Errorf("update price error: %w", err)
	}

	return nil
}

// Upsert 插入或更新
func (r *CompetitorRepository) Upsert(ctx context.Context, c *entity.Competitor) error {
	if err := c.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	c.UpdatedAt = time.Now()
	c.CreatedAt = time.Now()

	query := `
		INSERT INTO competitors (id, name, platform, platform_product_id, current_price,
			currency, rating, review_count, seller_rating, seller_review_count, source_url,
			created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (platform, platform_product_id) DO UPDATE SET
			name = EXCLUDED.name,
			current_price = EXCLUDED.current_price,
			currency = EXCLUDED.currency,
			rating = EXCLUDED.rating,
			review_count = EXCLUDED.review_count,
			seller_rating = EXCLUDED.seller_rating,
			seller_review_count = EXCLUDED.seller_review_count,
			source_url = EXCLUDED.source_url,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.Exec(ctx, query,
		c.ID, c.Name, c.Platform, c.PlatformProductID, c.CurrentPrice,
		c.Currency, c.Rating, c.ReviewCount, c.SellerRating, c.SellerReviewCount,
		c.SourceURL, c.CreatedAt, c.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("upsert competitor error: %w", err)
	}

	return nil
}
