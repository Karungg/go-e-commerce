package main

import (
	"go-e-commerce/internal/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	db := config.SetupDatabase(cfg)
	_ = db
}