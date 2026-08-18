package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// Create temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.yaml")

	configContent := `
database:
  host: "localhost"
  port: 5432
  user: "testuser"
  password: "testpass"
  name: "testdb"
  sslmode: "disable"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

app:
  host: "0.0.0.0"
  port: 8080
  mode: "test"

llm:
  provider: "qwen"
  api_key: "test-key"
  base_url: "https://test.example.com"

asynq:
  redis_host: "localhost:6379"
  redis_password: ""
  redis_db: 1
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	cfg, err := Load(configPath)
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// Database assertions
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, "testuser", cfg.Database.User)
	assert.Equal(t, "testpass", cfg.Database.Password)
	assert.Equal(t, "testdb", cfg.Database.Name)
	assert.Equal(t, "disable", cfg.Database.SSLMode)

	// Redis assertions
	assert.Equal(t, "localhost", cfg.Redis.Host)
	assert.Equal(t, 6379, cfg.Redis.Port)
	assert.Equal(t, 0, cfg.Redis.DB)

	// App assertions
	assert.Equal(t, "0.0.0.0", cfg.App.Host)
	assert.Equal(t, 8080, cfg.App.Port)
	assert.Equal(t, "test", cfg.App.Mode)

	// LLM assertions
	assert.Equal(t, "qwen", cfg.LLM.Provider)
	assert.Equal(t, "test-key", cfg.LLM.APIKey)
	assert.Equal(t, "https://test.example.com", cfg.LLM.BaseURL)

	// Asynq assertions
	assert.Equal(t, "localhost:6379", cfg.Asynq.RedisHost)
	assert.Equal(t, 1, cfg.Asynq.RedisDB)
}

func TestLoadWithEnvVar(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.yaml")

	os.Setenv("TEST_DB_PASSWORD", "secret-password")
	defer os.Unsetenv("TEST_DB_PASSWORD")

	configContent := `
database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "${TEST_DB_PASSWORD}"
  name: "testdb"
  sslmode: "disable"
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	cfg, err := Load(configPath)
	require.NoError(t, err)
	assert.Equal(t, "secret-password", cfg.Database.Password)
}

func TestLoadFileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	assert.Error(t, err)
}

func TestLoadInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")

	err := os.WriteFile(configPath, []byte("invalid: yaml: content:"), 0644)
	require.NoError(t, err)

	_, err = Load(configPath)
	assert.Error(t, err)
}

func TestMustLoad(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.yaml")

	configContent := `
app:
  host: "localhost"
  port: 8080
  mode: "test"
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	cfg := MustLoad(configPath)
	assert.NotNil(t, cfg)
	assert.Equal(t, 8080, cfg.App.Port)
}

func TestMustLoadPanic(t *testing.T) {
	assert.Panics(t, func() {
		MustLoad("/nonexistent/path/config.yaml")
	})
}
