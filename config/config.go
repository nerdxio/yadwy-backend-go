package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort uint16
	DB         DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfig() Config {
	cfg := Config{
		ServerPort: 3000,
		DB: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "user",
			Password: "pass",
			DBName:   "yadwy_db",
			SSLMode:  "disable",
		},
	}

	if serverPort, exists := os.LookupEnv("SERVER_PORT"); exists {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = uint16(port)
		}
	}

	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		cfg.DB.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		if port, err := strconv.Atoi(dbPort); err == nil {
			cfg.DB.Port = port
		}
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		cfg.DB.User = dbUser
	}
	if dbPass := os.Getenv("DB_PASSWORD"); dbPass != "" {
		cfg.DB.Password = dbPass
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		cfg.DB.DBName = dbName
	}
	if sslMode := os.Getenv("DB_SSL_MODE"); sslMode != "" {
		cfg.DB.SSLMode = sslMode
	}

	return cfg
}
