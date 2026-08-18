package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"competitive-analysis-agent/internal/config"

	_ "github.com/lib/pq"
)

// PostgresDB PostgreSQL 连接管理
type PostgresDB struct {
	db  *sql.DB
	cfg *config.DatabaseConfig
}

// NewPostgresDB 创建 PostgreSQL 连接
func NewPostgresDB(cfg *config.DatabaseConfig) (*PostgresDB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 配置连接池
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{db: db, cfg: cfg}, nil
}

// DB 获取数据库连接
func (p *PostgresDB) DB() *sql.DB {
	return p.db
}

// Close 关闭连接
func (p *PostgresDB) Close() error {
	return p.db.Close()
}

// Ping 检查连接
func (p *PostgresDB) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

// Exec 执行 SQL
func (p *PostgresDB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return p.db.ExecContext(ctx, query, args...)
}

// Query 查询
func (p *PostgresDB) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return p.db.QueryContext(ctx, query, args...)
}

// QueryRow 查询单行
func (p *PostgresDB) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return p.db.QueryRowContext(ctx, query, args...)
}

// BeginTx 开始事务
func (p *PostgresDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return p.db.BeginTx(ctx, opts)
}

// WithTransaction 执行事务
func (p *PostgresDB) WithTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := p.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction error: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback error: %v, original error: %w", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit error: %w", err)
	}

	return nil
}

// HealthCheck 健康检查
func (p *PostgresDB) HealthCheck(ctx context.Context) error {
	return p.db.PingContext(ctx)
}
