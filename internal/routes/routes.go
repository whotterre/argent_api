package routes

import (
	"net/http"
	"whotterre/argent/internal/config"
	"whotterre/argent/internal/handlers"
	"whotterre/argent/internal/repositories"
	"whotterre/argent/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(app *gin.Engine, cfg config.Config, db *gorm.DB) {
	// Auth modules
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg)
	authHandler := handlers.NewAuthHandler(authService)


	auth := app.Group("/auth")
	// Google Auth routes
	auth.GET("/google", authHandler.HandleGoogleLogin)
	auth.GET("/google/callback", dummy)
	
}

func dummy(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{"hello": "Hello"})
}