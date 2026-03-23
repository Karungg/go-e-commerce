package main

import (
	"fmt"
	"go-e-commerce/internal/config"
	delivery "go-e-commerce/internal/delivery/http"
	"go-e-commerce/internal/repository"
	"go-e-commerce/internal/security"
	"go-e-commerce/internal/usecase"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	db := config.SetupDatabase(cfg)

	jwtAuth := security.NewJWTAuth(cfg.JWTSecret, cfg.JWTExpirationHours)

	userRepo := repository.NewUserRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	sellerRepo := repository.NewSellerRepository(db)

	authUsecase := usecase.NewAuthUseCase(db, userRepo, customerRepo, sellerRepo, jwtAuth)

	router := gin.Default()

	api := router.Group("/api")
	delivery.NewAuthController(api, authUsecase)

	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Server running on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}