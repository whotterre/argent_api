package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(app *gin.Engine) {
	// Google Auth routes
	app.GET("/auth/google", dummy)
	app.GET("/auth/google/callback", dummy)
	
}

func dummy(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{"hello": "Hello"})
}