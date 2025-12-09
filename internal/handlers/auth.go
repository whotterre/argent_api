package handlers

import (
	"net/http"
	"whotterre/argent/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) HandleGoogleLogin(c *gin.Context){
	// Generate CSRF token
	stateToken := h.authService.GenerateStateToken()
	// Store in cookie
	c.SetCookie("oauth_State", stateToken, 300, "/", "", false, true)

	authURL := h.authService.GetAuthCodeURL(stateToken)

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}