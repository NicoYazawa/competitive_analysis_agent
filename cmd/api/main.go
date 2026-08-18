package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"competitive-analysis-agent/internal/api"
	"competitive-analysis-agent/internal/config"
)

func main() {
	// Determine config path
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	configPath := fmt.Sprintf("configs/%s.yaml", env)

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)

	log.Printf("Starting server on %s", addr)
	log.Printf("Health check: http://%s/api/health", addr)

	if err := http.ListenAndServe(addr, api.NewRouter()); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
