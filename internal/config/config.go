// Package config provides application configuration loader
package config

import (
	"os"
	"time"
)

// Config holds all application configurations
type Config struct {
	Env    string
	Server ServerConfig
	DB     PostgresConfig
	Auth   AuthConfig
}

// ServerConfig holds HTTP server settings
type ServerConfig struct {
	Host string
	Port string
}

// PostgresConfig holds PostgreSQL connection settings
type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

// AuthConfig holds JWT authentication settings
type AuthConfig struct {
	SigningKey      string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// LoadConfig reads environment variables and returns Config
func LoadConfig() (*Config, error) {
	accessTTL, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_TTL"))
	if err != nil {
		accessTTL = time.Hour * 1
	}

	refreshTTL, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_TTL"))
	if err != nil {
		refreshTTL = time.Hour * 24 * 30
	}
	cfg := &Config{
		Env: os.Getenv("ENV_LOG"),
		Server: ServerConfig{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		DB: PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_EXTERNAL_PORT"),
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName:   os.Getenv("POSTGRES_DB"),
		},
		Auth: AuthConfig{
			SigningKey:      os.Getenv("SIGNING_KEY"),
			AccessTokenTTL:  accessTTL,
			RefreshTokenTTL: refreshTTL,
		},
	}

	return cfg, nil
}
