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

// CreateAPIKey godoc
// @Summary Create a new API key
// @Description Generate a new API key with specified permissions and expiry
// @Tags api-keys
// @Accept json
// @Produce json
// @Param request body dto.CreateAPIKeyRequest true "API key creation request"
// @Success 200 {object} dto.CreateAPIKeyResponse "API key creation response"
// @Failure 400 {object} map[string]string "error"
// @Failure 412 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Security BearerAuth
// @Router /keys/create [post]
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

// RolloverAPIKey godoc
// @Summary Rollover an existing API key
// @Description Generate a new API key to replace an existing one
// @Tags api-keys
// @Accept json
// @Produce json
// @Param request body dto.RolloverAPIKeyRequest true "API key rollover request"
// @Success 200 {object} dto.RolloverAPIKeyResponse "API key rollover response"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Security BearerAuth
// @Router /keys/rollover [post]
func (h *APIKeyHandler) RolloverAPIKey(c *gin.Context) {
	var req dto.RolloverAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Failed to parse request body because", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Something went wrong while rolling over the API key on our end",
		})
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	response, err := h.apiKeyService.RolloverAPIKey(&req, userID)
	if err != nil {
		log.Println("Failed to roll over API key because", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"api_key":    response.APIKey,
		"expires_at": response.ExpiresAt,
	})
}
