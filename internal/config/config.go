package config

import (
	"github.com/spf13/viper"

	"yadwy-backend/internal/database"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database database.Config
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port int
	Host string
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.sslmode", "disable")

	// Read environment variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
