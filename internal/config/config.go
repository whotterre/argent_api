package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port               string `mapstructure:"PORT"`
	DatabaseURL        string `mapstructure:"DATABASE_URL"`
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL  string `mapstructure:"GOOGLE_REDIRECT_URL"`
	JWTSecret          string `mapstructure:"JWT_SECRET"`
	PaystackSecret     string `mapstructure:"PAYSTACK_SECRET"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("../../")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// Enable automatic environment variable reading
	viper.AutomaticEnv()

	// Bind environment variables explicitly to ensure they're read
	// even when .env file doesn't exist
	viper.BindEnv("PORT")
	viper.BindEnv("DATABASE_URL")
	viper.BindEnv("GOOGLE_CLIENT_ID")
	viper.BindEnv("GOOGLE_CLIENT_SECRET")
	viper.BindEnv("GOOGLE_REDIRECT_URL")
	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("PAYSTACK_SECRET")

	// Try to read from .env file
	err = viper.ReadInConfig()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
		// Clear error to allow loading from environment variables
		err = nil
	}

	err = viper.Unmarshal(&config)
	return
}
