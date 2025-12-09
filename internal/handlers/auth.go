package handlers

import (
	"log"
	"net/http"
	"whotterre/argent/internal/config"
	"whotterre/argent/internal/dto"
	"whotterre/argent/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
	cfg         config.Config
}

func NewAuthHandler(authService services.AuthService, cfg config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
	}
}

func (h *AuthHandler) HandleGoogleLogin(c *gin.Context) {
	log.Println("Initiating Google login....")
	// Generate CSRF token
	stateToken := h.authService.GenerateStateToken()
	// Store in cookie
	c.SetCookie("oauth_state", stateToken, 300, "/", "", false, true)

	authURL := h.authService.GetAuthCodeURL(stateToken)

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *AuthHandler) HandleGoogleCallback(c *gin.Context) {
	savedState, err := c.Cookie("oauth_state")
	if err != nil {
		log.Printf("OAuth state cookie not found: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Authentication session expired. Please try again.",
		})
		return
	}

	queryState := c.Query("state")
	if queryState != savedState {
		log.Printf("State mismatch: expected %s, got %s", savedState, queryState)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid authentication request. Possible CSRF attack.",
		})
		return
	}

	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	if errorParam := c.Query("error"); errorParam != "" {
		log.Printf("OAuth error from Google: %s", errorParam)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Google authentication failed: " + errorParam,
		})
		return
	}

	code := c.Query("code")
	if code == "" {
		log.Println("No authorization code received")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No authorization code received from Google.",
		})
		return
	}

	token, err := h.authService.ExchangeCode(code)
	if err != nil {
		log.Printf("Failed to exchange code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to complete authentication with Google.",
		})
		return
	}

	userInfo, err := h.authService.GetGoogleUserInfo(token.AccessToken)
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve user information from Google.",
		})
		return
	}

	if userInfo.Email == "" || userInfo.ID == "" {
		log.Printf("Incomplete user info received: %+v", userInfo)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Incomplete user information received from Google.",
		})
		return
	}

	newUser := dto.CreateNewUserRequest{
		GoogleID:  userInfo.ID,
		Email:     userInfo.Email,
		FirstName: userInfo.FamilyName,
		LastName:  userInfo.GivenName,
	}

	user, err := h.authService.FindOrCreateUser(&newUser)
	if err != nil {
		return
	}

	jwtToken, err := h.authService.GenerateJWT(user, h.cfg.JWTSecret)
	if err != nil {
		log.Printf("Failed to generate JWT: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to generate authentication token.",
		})
		return
	}

	log.Printf("User authenticated successfully: %s (%s)", user.Email, user.ID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Authentication successful",
		"token":   jwtToken,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"wallet":     user.Wallet,
		},
	})
}
