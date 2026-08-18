package storage

import (
	"testing"

	"competitive-analysis-agent/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestNewPostgresDB_InvalidHost(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Host:     "invalid-host-that-does-not-exist",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Name:     "test",
		SSLMode:  "disable",
	}

	db, err := NewPostgresDB(cfg)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestNewPostgresDB_InvalidPort(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     99999,
		User:     "postgres",
		Password: "postgres",
		Name:     "test",
		SSLMode:  "disable",
	}

	db, err := NewPostgresDB(cfg)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestNewPostgresDB_EmptyPassword(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "",
		Name:     "test",
		SSLMode:  "disable",
	}

	db, err := NewPostgresDB(cfg)
	if err != nil {
		assert.Nil(t, db)
	}
}

func TestPostgresDB_DB(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Host:     "invalid",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Name:     "test",
		SSLMode:  "disable",
	}

	db, err := NewPostgresDB(cfg)
	if err != nil {
		assert.Nil(t, db)
		return
	}

	assert.NotNil(t, db)
	sqlDB := db.DB()
	assert.NotNil(t, sqlDB)
}

func TestPostgresDB_Close(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Host:     "invalid",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Name:     "test",
		SSLMode:  "disable",
	}

	db, err := NewPostgresDB(cfg)
	if err != nil {
		assert.Nil(t, db)
		return
	}

	err = db.Close()
	assert.NoError(t, err)
}

func TestDatabaseConfig_Struct(t *testing.T) {
	cfg := config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		Name:     "testdb",
		SSLMode:  "disable",
	}

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 5432, cfg.Port)
	assert.Equal(t, "postgres", cfg.User)
	assert.Equal(t, "password", cfg.Password)
	assert.Equal(t, "testdb", cfg.Name)
	assert.Equal(t, "disable", cfg.SSLMode)
}

func TestRedisConfig_Struct(t *testing.T) {
	cfg := config.RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "password",
		DB:       0,
	}

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 6379, cfg.Port)
	assert.Equal(t, "password", cfg.Password)
	assert.Equal(t, 0, cfg.DB)
}

func TestAppConfig_Struct(t *testing.T) {
	cfg := config.AppConfig{
		Host: "0.0.0.0",
		Port: 8080,
		Mode: "development",
	}

	assert.Equal(t, "0.0.0.0", cfg.Host)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "development", cfg.Mode)
}

func TestLLMConfig_Struct(t *testing.T) {
	cfg := config.LLMConfig{
		Provider: "qwen",
		APIKey:   "test-key",
		BaseURL:  "https://api.example.com",
	}

	assert.Equal(t, "qwen", cfg.Provider)
	assert.Equal(t, "test-key", cfg.APIKey)
	assert.Equal(t, "https://api.example.com", cfg.BaseURL)
}

func TestAsynqConfig_Struct(t *testing.T) {
	cfg := config.AsynqConfig{
		RedisHost:     "localhost",
		RedisPassword: "password",
		RedisDB:       0,
	}

	assert.Equal(t, "localhost", cfg.RedisHost)
	assert.Equal(t, "password", cfg.RedisPassword)
	assert.Equal(t, 0, cfg.RedisDB)
}
