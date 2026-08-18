package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCompetitorRepository_GetByID_Success tests GetByID with sqlmock
func TestCompetitorRepository_GetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewCompetitorRepository(newSqlmockDB(db))
	ctx := context.Background()
	compID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "name", "platform", "platform_product_id", "current_price",
		"currency", "rating", "review_count", "seller_rating", "seller_review_count",
		"source_url", "created_at", "updated_at",
	}).AddRow(
		compID, "Test Product", "amazon", "B001", 99.99,
		"USD", 4.5, 1000, 4.8, 500,
		"https://amazon.com/dp/B001", now, now,
	)

	mock.ExpectQuery("SELECT (.+) FROM competitors WHERE id = \\$1").
		WithArgs(compID).
		WillReturnRows(rows)

	comp, err := repo.GetByID(ctx, compID)
	require.NoError(t, err)
	require.NotNil(t, comp)
	assert.Equal(t, compID, comp.ID)
	assert.Equal(t, "Test Product", comp.Name)
	assert.Equal(t, "amazon", comp.Platform)
	assert.Equal(t, 99.99, comp.CurrentPrice)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestCompetitorRepository_GetByID_NotFound tests GetByID returning no rows
func TestCompetitorRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewCompetitorRepository(newSqlmockDB(db))
	ctx := context.Background()
	compID := uuid.New()

	mock.ExpectQuery("SELECT (.+) FROM competitors WHERE id = \\$1").
		WithArgs(compID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "platform", "platform_product_id", "current_price",
			"currency", "rating", "review_count", "seller_rating", "seller_review_count",
			"source_url", "created_at", "updated_at",
		}))

	comp, err := repo.GetByID(ctx, compID)
	require.NoError(t, err)
	assert.Nil(t, comp)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestCompetitorRepository_GetByID_Error tests GetByID with DB error
func TestCompetitorRepository_GetByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewCompetitorRepository(newSqlmockDB(db))
	ctx := context.Background()
	compID := uuid.New()

	mock.ExpectQuery("SELECT (.+) FROM competitors WHERE id = \\$1").
		WithArgs(compID).
		WillReturnError(context.DeadlineExceeded)

	comp, err := repo.GetByID(ctx, compID)
	require.Error(t, err)
	assert.Nil(t, comp)
	assert.Contains(t, err.Error(), "query competitor error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestCompetitorRepository_GetByPlatformAndProductID_Success tests GetByPlatformAndProductID
func TestCompetitorRepository_GetByPlatformAndProductID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewCompetitorRepository(newSqlmockDB(db))
	ctx := context.Background()
	compID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "name", "platform", "platform_product_id", "current_price",
		"currency", "rating", "review_count", "seller_rating", "seller_review_count",
		"source_url", "created_at", "updated_at",
	}).AddRow(
		compID, "Test Product", "amazon", "B001", 99.99,
		"USD", 4.5, 1000, 4.8, 500,
		"https://amazon.com/dp/B001", now, now,
	)

	mock.ExpectQuery("SELECT (.+) FROM competitors WHERE platform = \\$1 AND platform_product_id = \\$2").
		WithArgs("amazon", "B001").
		WillReturnRows(rows)

	comp, err := repo.GetByPlatformAndProductID(ctx, "amazon", "B001")
	require.NoError(t, err)
	require.NotNil(t, comp)
	assert.Equal(t, "Test Product", comp.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestCompetitorRepository_GetByPlatformAndProductID_NotFound tests GetByPlatformAndProductID returning no rows
func TestCompetitorRepository_GetByPlatformAndProductID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewCompetitorRepository(newSqlmockDB(db))
	ctx := context.Background()

	mock.ExpectQuery("SELECT (.+) FROM competitors WHERE platform = \\$1 AND platform_product_id = \\$2").
		WithArgs("amazon", "B001").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "platform", "platform_product_id", "current_price",
			"currency", "rating", "review_count", "seller_rating", "seller_review_count",
			"source_url", "created_at", "updated_at",
		}))

	comp, err := repo.GetByPlatformAndProductID(ctx, "amazon", "B001")
	require.NoError(t, err)
	assert.Nil(t, comp)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestCompetitorRepository_CountByPlatform_Success tests CountByPlatform
func TestCompetitorRepository_CountByPlatform_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewCompetitorRepository(newSqlmockDB(db))
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(42)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM competitors WHERE platform = \\$1").
		WithArgs("amazon").
		WillReturnRows(rows)

	count, err := repo.CountByPlatform(ctx, "amazon")
	require.NoError(t, err)
	assert.Equal(t, 42, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestCompetitorRepository_CountByPlatform_Error tests CountByPlatform with error
func TestCompetitorRepository_CountByPlatform_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewCompetitorRepository(newSqlmockDB(db))
	ctx := context.Background()

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM competitors WHERE platform = \\$1").
		WithArgs("amazon").
		WillReturnError(context.DeadlineExceeded)

	count, err := repo.CountByPlatform(ctx, "amazon")
	require.Error(t, err)
	assert.Equal(t, 0, count)
	assert.Contains(t, err.Error(), "count competitors error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestCompetitorRepository_ListByPlatform_Success tests ListByPlatform
func TestCompetitorRepository_ListByPlatform_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewCompetitorRepository(newSqlmockDB(db))
	ctx := context.Background()
	id1, id2 := uuid.New(), uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "name", "platform", "platform_product_id", "current_price",
		"currency", "rating", "review_count", "seller_rating", "seller_review_count",
		"source_url", "created_at", "updated_at",
	}).AddRow(id1, "Product1", "amazon", "P1", 99.99, "USD", 4.5, 100, 4.8, 50, "url1", now, now).
		AddRow(id2, "Product2", "amazon", "P2", 149.99, "USD", 4.3, 200, 4.6, 100, "url2", now, now)

	mock.ExpectQuery("SELECT (.+) FROM competitors WHERE platform = \\$1 ORDER BY updated_at DESC LIMIT \\$2 OFFSET \\$3").
		WithArgs("amazon", 10, 0).
		WillReturnRows(rows)

	competitors, err := repo.ListByPlatform(ctx, "amazon", 10, 0)
	require.NoError(t, err)
	assert.Len(t, competitors, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestCompetitorRepository_ListByPlatform_Empty tests ListByPlatform with no results
func TestCompetitorRepository_ListByPlatform_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewCompetitorRepository(newSqlmockDB(db))
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{
		"id", "name", "platform", "platform_product_id", "current_price",
		"currency", "rating", "review_count", "seller_rating", "seller_review_count",
		"source_url", "created_at", "updated_at",
	})
	mock.ExpectQuery("SELECT (.+) FROM competitors WHERE platform = \\$1 ORDER BY updated_at DESC LIMIT \\$2 OFFSET \\$3").
		WithArgs("amazon", 10, 0).
		WillReturnRows(rows)

	competitors, err := repo.ListByPlatform(ctx, "amazon", 10, 0)
	require.NoError(t, err)
	assert.Len(t, competitors, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestCompetitorRepository_ListAll_Success tests ListAll
func TestCompetitorRepository_ListAll_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewCompetitorRepository(newSqlmockDB(db))
	ctx := context.Background()
	id1 := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "name", "platform", "platform_product_id", "current_price",
		"currency", "rating", "review_count", "seller_rating", "seller_review_count",
		"source_url", "created_at", "updated_at",
	}).AddRow(id1, "Product1", "amazon", "P1", 99.99, "USD", 4.5, 100, 4.8, 50, "url1", now, now)

	mock.ExpectQuery("SELECT (.+) FROM competitors ORDER BY updated_at DESC LIMIT \\$1 OFFSET \\$2").
		WithArgs(10, 0).
		WillReturnRows(rows)

	competitors, err := repo.ListAll(ctx, 10, 0)
	require.NoError(t, err)
	assert.Len(t, competitors, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestCompetitorRepository_ListAll_Error tests ListAll with error
func TestCompetitorRepository_ListAll_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewCompetitorRepository(newSqlmockDB(db))
	ctx := context.Background()

	mock.ExpectQuery("SELECT (.+) FROM competitors").
		WillReturnError(context.DeadlineExceeded)

	competitors, err := repo.ListAll(ctx, 10, 0)
	require.Error(t, err)
	assert.Nil(t, competitors)
	assert.Contains(t, err.Error(), "query competitors error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestCompetitorRepository_ListByPlatform_ScanError tests ListByPlatform with scan error
func TestCompetitorRepository_ListByPlatform_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewCompetitorRepository(newSqlmockDB(db))
	ctx := context.Background()

	// Fewer columns than expected - will cause scan error
	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(uuid.New(), "Test")
	mock.ExpectQuery("SELECT (.+) FROM competitors WHERE platform = \\$1").
		WithArgs("amazon").
		WillReturnRows(rows)

	competitors, err := repo.ListByPlatform(ctx, "amazon", 10, 0)
	require.Error(t, err)
	assert.Nil(t, competitors)
}

// TestCompetitorRepository_GetByID_ScanError tests GetByID with scan error
func TestCompetitorRepository_GetByID_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewCompetitorRepository(newSqlmockDB(db))
	ctx := context.Background()
	compID := uuid.New()

	// Fewer columns than expected - will cause scan error
	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(compID, "Test")
	mock.ExpectQuery("SELECT (.+) FROM competitors WHERE id = \\$1").
		WithArgs(compID).
		WillReturnRows(rows)

	comp, err := repo.GetByID(ctx, compID)
	require.Error(t, err)
	assert.Nil(t, comp)
}
