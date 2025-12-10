package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"whotterre/argent/internal/config"
	"whotterre/argent/internal/dto"
	"whotterre/argent/internal/models"
	"whotterre/argent/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService interface {
	GenerateStateToken() string
	GetAuthCodeURL(state string) string
	ExchangeCode(code string) (*oauth2.Token, error)
	GenerateJWT(user *models.User, jwtSecret string) (string, error)
	GetGoogleUserInfo(accessToken string) (*models.GoogleUserInfo, error)
	FindOrCreateUser(newUser *dto.CreateNewUserRequest) (*models.User, error)
	ParseJWT(tokenString string) (*jwt.MapClaims, error)
	GetUserIDFromJWT(tokenString string) (uuid.UUID, error)
}

type authService struct {
	authRepo    repositories.UserRepository
	oauthConfig *oauth2.Config
	jwtSecret   string
}

func NewAuthService(authRepo repositories.UserRepository, cfg config.Config) AuthService {
	redirectURL := cfg.BaseURL + "/auth/google/callback"
	log.Printf("OAuth Redirect URL: %s", redirectURL)
	return &authService{
		authRepo: authRepo,
		oauthConfig: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		jwtSecret: cfg.JWTSecret,
	}
}

// Generates a 32 character CSRF token string
func (s *authService) GenerateStateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (s *authService) GetAuthCodeURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *authService) ExchangeCode(code string) (*oauth2.Token, error) {
	return s.oauthConfig.Exchange(context.Background(), code)
}

func (s *authService) GetGoogleUserInfo(accessToken string) (*models.GoogleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo models.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func (s *authService) GenerateJWT(user *models.User, jwtSecret string) (string, error) {
	// jwt payload
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (s *authService) FindOrCreateUser(newUser *dto.CreateNewUserRequest) (*models.User, error) {
	user, err := s.authRepo.FindOrCreateUser(newUser)
	if err != nil {
		log.Println("Failed to find or create user because", err.Error())
		return nil, err
	}
	return user, nil
}

func (s *authService) ParseJWT(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (s *authService) GetUserIDFromJWT(tokenString string) (uuid.UUID, error) {
	claims, err := s.ParseJWT(tokenString)
	if err != nil {
		return uuid.Nil, err
	}

	userIDStr, ok := (*claims)["user_id"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("user_id claim not found or not a string")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user_id format: %v", err)
	}

	return userID, nil
}
