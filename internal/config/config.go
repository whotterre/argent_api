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

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	err = viper.Unmarshal(&config)
	return
}
