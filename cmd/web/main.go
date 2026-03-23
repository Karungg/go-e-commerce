package main

import (
	"fmt"
	"log/slog"
	"os"

	"go-e-commerce/internal/config"
	delivery "go-e-commerce/internal/delivery/http"
	"go-e-commerce/internal/repository"
	"go-e-commerce/internal/security"
	"go-e-commerce/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db := config.SetupDatabase(cfg)

	jwtAuth := security.NewJWTAuth(cfg.JWTSecret, cfg.JWTExpirationHours)

	userRepo := repository.NewUserRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	sellerRepo := repository.NewSellerRepository(db)

	authUsecase := usecase.NewAuthUseCase(db, logger, userRepo, customerRepo, sellerRepo, jwtAuth)

	router := gin.Default()
	api := router.Group("/api")
	delivery.NewAuthController(api, authUsecase)

	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
	logger.Info("Server running", slog.String("port", serverAddr))
	
	if err := router.Run(serverAddr); err != nil {
		logger.Error("Server failed to start", slog.Any("error", err))
		os.Exit(1)
	}
}