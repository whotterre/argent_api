package routes

import (
	"whotterre/argent/internal/config"
	"whotterre/argent/internal/handlers"
	"whotterre/argent/internal/middleware"
	"whotterre/argent/internal/repositories"
	"whotterre/argent/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRoutes(app *gin.Engine, cfg config.Config, db *gorm.DB) {
	// Auth modules
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg)
	authHandler := handlers.NewAuthHandler(authService, cfg)

	// API Key modules
	auth := app.Group("/auth")
	// Google Auth routes
	auth.GET("/google", authHandler.HandleGoogleLogin)
	auth.GET("/google/callback", authHandler.HandleGoogleCallback)
	// API Key routes
	apiKeyRepo := repositories.NewAPIKeyRepository(db)
	apiKeyService := services.NewAPIKeyService(apiKeyRepo)
	apiKeyHandler := handlers.NewAPIKeyHandler(apiKeyService)

	apiKey := app.Group("/keys")
	apiKey.Use(middleware.RequireAuth(authService, apiKeyService, ""))
	apiKey.POST("/create", apiKeyHandler.CreateAPIKey)
	apiKey.POST("/rollover", apiKeyHandler.RolloverAPIKey)

	// Wallet modules
	walletRepo := repositories.NewWalletRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	walletService := services.NewWalletService(walletRepo, transactionRepo, userRepo, cfg.PaystackSecret, db, cfg)
	walletHandler := handlers.NewWalletHandler(walletService)

	wallet := app.Group("/wallet")
	wallet.Use(middleware.RequireAuth(authService, apiKeyService, "read"))
	wallet.POST("/deposit", walletHandler.Deposit)
	wallet.GET("/balance", walletHandler.GetBalance)
	wallet.POST("/transfer", walletHandler.Transfer)
	wallet.GET("/transactions", walletHandler.GetTransactions)
	wallet.GET("/deposit/:reference/status", walletHandler.GetDepositStatus)

	// Public wallet endpoints (no auth required)
	app.POST("/wallet/paystack/webhook", walletHandler.Webhook)
	app.GET("/wallet/deposit/callback", walletHandler.DepositCallback)
	// Swagger docs
	app.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

