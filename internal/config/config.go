package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Auth     AuthConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URL string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port     string
	Platform string
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret string
	PolkaKey  string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file (optional in production)
	_ = godotenv.Load()

	cfg := &Config{
		Database: DatabaseConfig{
			URL: os.Getenv("DB_URL"),
		},
		Server: ServerConfig{
			Port:     getEnvOrDefault("PORT", "8080"),
			Platform: getEnvOrDefault("PLATFORM", "prod"),
		},
		Auth: AuthConfig{
			JWTSecret: os.Getenv("JWT_SECRET"),
			PolkaKey:  os.Getenv("POLKA_KEY"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("DB_URL is required")
	}
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.Auth.PolkaKey == "" {
		return fmt.Errorf("POLKA_KEY is required")
	}
	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
