package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DatabaseURL        string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	JWTSecret          string
	PaystackSecret     string
}

func LoadConfig() (config Config, err error) {
	// Load .env file if it exists (ignore error if not)
	_ = godotenv.Load()

	// Read from environment variables
	config.Port = os.Getenv("PORT")
	config.DatabaseURL = os.Getenv("DATABASE_URL")
	config.GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	config.GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	config.GoogleRedirectURL = os.Getenv("GOOGLE_REDIRECT_URL")
	config.JWTSecret = os.Getenv("JWT_SECRET")
	config.PaystackSecret = os.Getenv("PAYSTACK_SECRET")

	// Log if DATABASE_URL is empty (for debugging)
	if config.DatabaseURL == "" {
		log.Println("Warning: DATABASE_URL is not set")
	}

	return config, nil
}
