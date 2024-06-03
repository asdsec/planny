package configs

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// Config stores the configuration for the application
type Config struct {
	Environment          string        `mapstructure:"ENVIRONMENT"`
	DatabaseURL          string        `mapstructure:"DATABASE_URL"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

// Load loads the configuration from the environment variables
func Load() (config Config, err error) {
	viper.AddConfigPath("configs")
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
		return
	}

	err = viper.Unmarshal(&config)
	return
}
