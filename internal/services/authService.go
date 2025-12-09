package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"whotterre/argent/internal/config"
	"whotterre/argent/internal/dto"
	"whotterre/argent/internal/models"
	"whotterre/argent/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
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
}

type authService struct {
	authRepo    repositories.UserRepository
	oauthConfig *oauth2.Config
	jwtSecret   string
}

func NewAuthService(authRepo repositories.UserRepository, cfg config.Config) AuthService {
	return &authService{
		authRepo: authRepo,
		oauthConfig: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
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
