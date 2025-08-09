package main

import (
	"log"

	"github.com/andreyxaxa/order_svc/config"
	"github.com/andreyxaxa/order_svc/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	// Config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error :%s", err)
	}

	// Run
	app.Run(cfg)
}
