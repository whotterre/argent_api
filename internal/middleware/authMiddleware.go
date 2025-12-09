package middleware

import (
	"net/http"
	"strings"
	"whotterre/argent/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequireAuth(authService services.AuthService, apiKeyService services.APIKeyService, requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		apiKeyHeader := c.GetHeader("x-api-key")

		var userID uuid.UUID
		var err error

		if authHeader != "" {
			// JWT auth
			if !strings.HasPrefix(authHeader, "Bearer ") {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
				c.Abort()
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			userID, err = authService.GetUserIDFromJWT(tokenString)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
				c.Abort()
				return
			}
		} else if apiKeyHeader != "" {
			// API key auth
			if requiredPermission == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "API key not allowed for this endpoint"})
				c.Abort()
				return
			}

			apiKey, err := apiKeyService.ValidateAPIKey(apiKeyHeader, uuid.Nil, requiredPermission)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key or insufficient permissions"})
				c.Abort()
				return
			}

			userID = apiKey.UserID
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header or x-api-key required"})
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("user_id", userID)
		c.Next()
	}
}
