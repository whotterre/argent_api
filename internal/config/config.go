package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DatabaseURL        string
	BaseURL            string
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
	config.BaseURL = os.Getenv("BASE_URL")
	config.GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	config.GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	config.GoogleRedirectURL = os.Getenv("GOOGLE_REDIRECT_URL")
	config.JWTSecret = os.Getenv("JWT_SECRET")
	config.PaystackSecret = os.Getenv("PAYSTACK_SECRET")

	// Debug log
	log.Printf("Config loaded: PORT=%s, DATABASE_URL=%s, BASE_URL=%s", config.Port, config.DatabaseURL, config.BaseURL)

	return config, nil
}
