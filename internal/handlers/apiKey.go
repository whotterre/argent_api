package handlers

import (
	"errors"
	"log"
	"net/http"
	"slices"

	"whotterre/argent/internal/customErrors"
	"whotterre/argent/internal/dto"
	"whotterre/argent/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type APIKeyHandler struct {
	apiKeyService services.APIKeyService
}

func NewAPIKeyHandler(apiKeyService services.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{
		apiKeyService: apiKeyService,
	}
}

func (h *APIKeyHandler) CreateAPIKey(c *gin.Context) {
	var req dto.CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Failed to read request body because ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
	}

	validExpiryInputs := []string{"1H", "1D", "1M", "1Y"}
	// Check if the expiry input isn't there
	if !slices.Contains(validExpiryInputs, req.Expiry) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid expiry format. Supported formats are '1H', '1D', '1M', '1Y'",
		})
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	response, err := h.apiKeyService.CreateAPIKey(req, userID)
	if err != nil {
		if errors.Is(err, customErrors.ErrorActiveAPIKeysExceeded) {
			c.JSON(http.StatusPreconditionFailed, gin.H{
				"error": "You can only have five active API keys",
			})
			return
		}
		if errors.Is(err, customErrors.ErrInvalidPermission) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid permission provided",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create API key",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"api_key":    response.APIKey,
		"expires_at": response.ExpiresAt,
	})
}

func (h *APIKeyHandler) RolloverAPIKey(c *gin.Context){
	var req dto.RolloverAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Failed to parse request body because", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Something went wrong while rolling over the API key on our end",
		})
	}

	response, err := h.apiKeyService.RolloverAPIKey(&req)
	if err != nil {
		log.Println("Failed to roll over API key because", err.Error())
		return 
	}
	
	c.JSON(http.StatusOK, gin.H{
		"api_key":    response.APIKey,
		"expires_at": response.ExpiresAt,
	})
}