package services

import (
	"crypto/rand"
	"encoding/base64"
	"whotterre/argent/internal/config"
	"whotterre/argent/internal/repositories"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService interface {
	GenerateStateToken() string
	GetAuthCodeURL(state string) string
}

type authService struct {
	authRepo    repositories.UserRepository
	oauthConfig *oauth2.Config
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
