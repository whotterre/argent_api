package main

import (
	"whotterre/argent/internal/config"
	"whotterre/argent/internal/initializers"
	"whotterre/argent/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	cfg, err := config.LoadConfig()
	if err != nil {
	   return 
	}

	// Connect to database
	initializers.ConnectToDB(cfg.DatabaseURL)
	db := initializers.DB
	routes.SetupRoutes(app, cfg, db)
	
	port := ":" + cfg.Port
	app.Run(port)
}