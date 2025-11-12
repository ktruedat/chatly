package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ktruedat/chatly/internal/application"
	"github.com/ktruedat/chatly/internal/application/config"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, proceeding with environment variables")
	}

	configPath := getEnv("CONFIG_PATH", "config.yaml")
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	app := application.New(cfg)
	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Run application
	log.Printf("Chatly Backend starting...")
	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
