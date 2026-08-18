package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	App      AppConfig      `yaml:"app"`
	LLM      LLMConfig      `yaml:"llm"`
	Asynq    AsynqConfig    `yaml:"asynq"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"sslmode"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type AppConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type LLMConfig struct {
	Provider string `yaml:"provider"`
	APIKey   string `yaml:"api_key"`
	BaseURL  string `yaml:"base_url"`
}

type AsynqConfig struct {
	RedisHost     string `yaml:"redis_host"`
	RedisPassword string `yaml:"redis_password"`
	RedisDB       int    `yaml:"redis_db"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Expand environment variables
	content := os.ExpandEnv(string(data))

	var cfg Config
	if err := yaml.Unmarshal([]byte(content), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func MustLoad(path string) *Config {
	cfg, err := Load(path)
	if err != nil {
		panic(err)
	}
	return cfg
}
