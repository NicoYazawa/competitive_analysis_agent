package repository

import (
	"context"
	"database/sql"
)

// Queryer interface for database query operations
type Queryer interface {
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row
}
