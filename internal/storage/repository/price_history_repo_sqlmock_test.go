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

// TestPriceHistoryRepository_GetByID_Success tests GetByID with sqlmock
func TestPriceHistoryRepository_GetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPriceHistoryRepository(newSqlmockDB(db))
	ctx := context.Background()
	id := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "competitor_id", "price", "currency", "recorded_at",
	}).AddRow(id, uuid.New(), 99.99, "USD", now)

	mock.ExpectQuery("SELECT (.+) FROM price_history WHERE id = \\$1").
		WithArgs(id).
		WillReturnRows(rows)

	p, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, p)
	assert.Equal(t, id, p.ID)
	assert.Equal(t, 99.99, p.Price)
	assert.Equal(t, "USD", p.Currency)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestPriceHistoryRepository_GetByID_NotFound tests GetByID returning no rows
func TestPriceHistoryRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPriceHistoryRepository(newSqlmockDB(db))
	ctx := context.Background()
	id := uuid.New()

	mock.ExpectQuery("SELECT (.+) FROM price_history WHERE id = \\$1").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "competitor_id", "price", "currency", "recorded_at",
		}))

	p, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	assert.Nil(t, p)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestPriceHistoryRepository_GetByID_Error tests GetByID with DB error
func TestPriceHistoryRepository_GetByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPriceHistoryRepository(newSqlmockDB(db))
	ctx := context.Background()
	id := uuid.New()

	mock.ExpectQuery("SELECT (.+) FROM price_history WHERE id = \\$1").
		WithArgs(id).
		WillReturnError(context.DeadlineExceeded)

	p, err := repo.GetByID(ctx, id)
	require.Error(t, err)
	assert.Nil(t, p)
	assert.Contains(t, err.Error(), "query price history error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestPriceHistoryRepository_GetLatest_Success tests GetLatest with sqlmock
func TestPriceHistoryRepository_GetLatest_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPriceHistoryRepository(newSqlmockDB(db))
	ctx := context.Background()
	compID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "competitor_id", "price", "currency", "recorded_at",
	}).AddRow(uuid.New(), compID, 149.99, "USD", now)

	mock.ExpectQuery("SELECT (.+) FROM price_history WHERE competitor_id = \\$1").
		WithArgs(compID).
		WillReturnRows(rows)

	p, err := repo.GetLatest(ctx, compID)
	require.NoError(t, err)
	require.NotNil(t, p)
	assert.Equal(t, compID, p.CompetitorID)
	assert.Equal(t, 149.99, p.Price)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestPriceHistoryRepository_GetLatest_NotFound tests GetLatest returning no rows
func TestPriceHistoryRepository_GetLatest_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPriceHistoryRepository(newSqlmockDB(db))
	ctx := context.Background()
	compID := uuid.New()

	mock.ExpectQuery("SELECT (.+) FROM price_history WHERE competitor_id = \\$1").
		WithArgs(compID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "competitor_id", "price", "currency", "recorded_at",
		}))

	p, err := repo.GetLatest(ctx, compID)
	require.NoError(t, err)
	assert.Nil(t, p)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestPriceHistoryRepository_GetLatest_Error tests GetLatest with DB error
func TestPriceHistoryRepository_GetLatest_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPriceHistoryRepository(newSqlmockDB(db))
	ctx := context.Background()
	compID := uuid.New()

	mock.ExpectQuery("SELECT (.+) FROM price_history WHERE competitor_id = \\$1").
		WithArgs(compID).
		WillReturnError(context.DeadlineExceeded)

	p, err := repo.GetLatest(ctx, compID)
	require.Error(t, err)
	assert.Nil(t, p)
	assert.Contains(t, err.Error(), "query price history error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestPriceHistoryRepository_Count_Success tests Count with sqlmock
func TestPriceHistoryRepository_Count_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPriceHistoryRepository(newSqlmockDB(db))
	ctx := context.Background()
	compID := uuid.New()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(42)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM price_history WHERE competitor_id = \\$1").
		WithArgs(compID).
		WillReturnRows(rows)

	count, err := repo.Count(ctx, compID)
	require.NoError(t, err)
	assert.Equal(t, 42, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestPriceHistoryRepository_Count_Empty tests Count returning 0
func TestPriceHistoryRepository_Count_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPriceHistoryRepository(newSqlmockDB(db))
	ctx := context.Background()
	compID := uuid.New()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM price_history WHERE competitor_id = \\$1").
		WithArgs(compID).
		WillReturnRows(rows)

	count, err := repo.Count(ctx, compID)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestPriceHistoryRepository_Count_Error tests Count with DB error
func TestPriceHistoryRepository_Count_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPriceHistoryRepository(newSqlmockDB(db))
	ctx := context.Background()
	compID := uuid.New()

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM price_history WHERE competitor_id = \\$1").
		WithArgs(compID).
		WillReturnError(context.DeadlineExceeded)

	count, err := repo.Count(ctx, compID)
	require.Error(t, err)
	assert.Equal(t, 0, count)
	assert.Contains(t, err.Error(), "count price history error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestPriceHistoryRepository_GetAveragePrice_Success tests GetAveragePrice with sqlmock
func TestPriceHistoryRepository_GetAveragePrice_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPriceHistoryRepository(newSqlmockDB(db))
	ctx := context.Background()
	compID := uuid.New()

	rows := sqlmock.NewRows([]string{"avg"}).AddRow(123.45)
	mock.ExpectQuery("SELECT AVG\\(price\\) FROM price_history WHERE competitor_id = \\$1").
		WithArgs(compID, 7).
		WillReturnRows(rows)

	avg, err := repo.GetAveragePrice(ctx, compID, 7)
	require.NoError(t, err)
	assert.Equal(t, 123.45, avg)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestPriceHistoryRepository_GetAveragePrice_Null tests GetAveragePrice with NULL result
func TestPriceHistoryRepository_GetAveragePrice_Null(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPriceHistoryRepository(newSqlmockDB(db))
	ctx := context.Background()
	compID := uuid.New()

	rows := sqlmock.NewRows([]string{"avg"}).AddRow(nil)
	mock.ExpectQuery("SELECT AVG\\(price\\) FROM price_history WHERE competitor_id = \\$1").
		WithArgs(compID, 30).
		WillReturnRows(rows)

	avg, err := repo.GetAveragePrice(ctx, compID, 30)
	require.NoError(t, err)
	assert.Equal(t, 0.0, avg)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestPriceHistoryRepository_GetAveragePrice_Error tests GetAveragePrice with DB error
func TestPriceHistoryRepository_GetAveragePrice_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPriceHistoryRepository(newSqlmockDB(db))
	ctx := context.Background()
	compID := uuid.New()

	mock.ExpectQuery("SELECT AVG\\(price\\) FROM price_history WHERE competitor_id = \\$1").
		WithArgs(compID, 7).
		WillReturnError(context.DeadlineExceeded)

	avg, err := repo.GetAveragePrice(ctx, compID, 7)
	require.Error(t, err)
	assert.Equal(t, 0.0, avg)
	assert.Contains(t, err.Error(), "query average price error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestPriceHistoryRepository_GetMinMaxPrice_Success tests GetMinMaxPrice with sqlmock
func TestPriceHistoryRepository_GetMinMaxPrice_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPriceHistoryRepository(newSqlmockDB(db))
	ctx := context.Background()
	compID := uuid.New()

	rows := sqlmock.NewRows([]string{"min", "max"}).AddRow(50.00, 199.99)
	mock.ExpectQuery("SELECT MIN\\(price\\), MAX\\(price\\) FROM price_history WHERE competitor_id = \\$1").
		WithArgs(compID, 30).
		WillReturnRows(rows)

	min, max, err := repo.GetMinMaxPrice(ctx, compID, 30)
	require.NoError(t, err)
	assert.Equal(t, 50.00, min)
	assert.Equal(t, 199.99, max)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestPriceHistoryRepository_GetMinMaxPrice_Error tests GetMinMaxPrice with DB error
func TestPriceHistoryRepository_GetMinMaxPrice_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPriceHistoryRepository(newSqlmockDB(db))
	ctx := context.Background()
	compID := uuid.New()

	mock.ExpectQuery("SELECT MIN\\(price\\), MAX\\(price\\) FROM price_history WHERE competitor_id = \\$1").
		WithArgs(compID, 30).
		WillReturnError(context.DeadlineExceeded)

	min, max, err := repo.GetMinMaxPrice(ctx, compID, 30)
	require.Error(t, err)
	assert.Equal(t, 0.0, min)
	assert.Equal(t, 0.0, max)
	assert.Contains(t, err.Error(), "query min max price error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestPriceHistoryRepository_GetByID_ScanError tests GetByID with scan error
func TestPriceHistoryRepository_GetByID_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPriceHistoryRepository(newSqlmockDB(db))
	ctx := context.Background()
	id := uuid.New()

	// Fewer columns than expected - will cause scan error
	rows := sqlmock.NewRows([]string{"id", "competitor_id"}).AddRow(id, uuid.New())
	mock.ExpectQuery("SELECT (.+) FROM price_history WHERE id = \\$1").
		WithArgs(id).
		WillReturnRows(rows)

	p, err := repo.GetByID(ctx, id)
	require.Error(t, err)
	assert.Nil(t, p)
}
