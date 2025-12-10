package config

import (
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
	viper.AutomaticEnv()

	err = viper.Unmarshal(&config)
	return
}
