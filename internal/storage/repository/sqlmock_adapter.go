package repository

import (
	"context"
	"database/sql"
)

// sqlmockDB wraps a *sql.DB to implement the Queryer interface with context support.
// This is needed because go-sqlmock's DB doesn't have context-aware methods.
type sqlmockDB struct {
	db *sql.DB
}

func (s *sqlmockDB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s *sqlmockDB) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func (s *sqlmockDB) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

// newSqlmockDB creates a Queryer-compatible wrapper around a *sql.DB from go-sqlmock.
func newSqlmockDB(db *sql.DB) Queryer {
	return &sqlmockDB{db: db}
}
