package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-e-commerce/internal/config"
	delivery "go-e-commerce/internal/delivery/http"
	"go-e-commerce/internal/delivery/http/route"
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
	config.RunMigrations(db)

	jwtAuth := security.NewJWTAuth(cfg.JWTSecret, cfg.JWTExpirationHours)

	userRepo := repository.NewUserRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	sellerRepo := repository.NewSellerRepository(db)
	txManager := repository.NewTransactionManager(db)

	authUsecase := usecase.NewAuthUseCase(txManager, logger, userRepo, customerRepo, sellerRepo, jwtAuth)

	router := gin.Default()
	api := router.Group("/api")
	
	authController := delivery.NewAuthController(authUsecase)
	route.SetupRoutes(api, authController, jwtAuth)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	go func() {
		logger.Info("Server running", slog.String("port", serverAddr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	stop()
	logger.Info("Shutting down gracefully, press Ctrl+C again to force")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("Server exiting")
}