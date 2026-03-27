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
	authCtrl "go-e-commerce/internal/delivery/http/auth"
	categoryCtrl "go-e-commerce/internal/delivery/http/category"
	productCtrl "go-e-commerce/internal/delivery/http/product"
	"go-e-commerce/internal/delivery/http/route"
	authRepo "go-e-commerce/internal/repository/auth"
	categoryRepo "go-e-commerce/internal/repository/category"
	productRepo "go-e-commerce/internal/repository/product"
	"go-e-commerce/internal/repository"
	"go-e-commerce/internal/security"
	authUC "go-e-commerce/internal/usecase/auth"
	categoryUC "go-e-commerce/internal/usecase/category"
	productUC "go-e-commerce/internal/usecase/product"

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

	userRepo := authRepo.NewUserRepository(db)
	customerRepo := authRepo.NewCustomerRepository(db)
	sellerRepo := authRepo.NewSellerRepository(db)
	catRepo := categoryRepo.NewCategoryRepository(db)
	prodRepo := productRepo.NewProductRepository(db)
	txManager := repository.NewTransactionManager(db)

	authUsecase := authUC.NewAuthUseCase(txManager, logger, userRepo, customerRepo, sellerRepo, jwtAuth)
	categoryUsecase := categoryUC.NewCategoryUseCase(catRepo)
	productUsecase := productUC.NewProductUseCase(prodRepo)

	router := gin.Default()
	api := router.Group("/api")
	
	authController := authCtrl.NewAuthController(authUsecase)
	categoryController := categoryCtrl.NewCategoryController(categoryUsecase)
	productController := productCtrl.NewProductController(productUsecase)
	route.SetupRoutes(api, authController, categoryController, productController, jwtAuth)

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