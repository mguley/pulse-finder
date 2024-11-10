package config

import (
	"os"
	"strconv"
)

// Configuration holds the main application configuration settings.
type Configuration struct {
	Port int            // Application server port
	Env  string         // Environment (e.g., "development", "production")
	Jwt  JWTConfig      // Jwt configuration for authentication
	DB   DatabaseConfig // Database configuration
}

// JWTConfig holds configuration for JWT-based authentication.
type JWTConfig struct {
	Secret string // Secret key for signing JWT tokens
}

// DatabaseConfig holds settings for database connection.
type DatabaseConfig struct {
	DSN string // Data source name for database connection
}

// LoadConfig loads the configuration settings from environment variables, falling back to default
// values if variables are not set.
func LoadConfig() *Configuration {
	config := &Configuration{
		Port: getEnvAsInt("PORT", 4005),
		Env:  getEnv("ENV", "development"),
		Jwt: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
		},
		DB: DatabaseConfig{
			DSN: getEnv("DB_DSN", ""),
		},
	}

	return config
}

// getEnv fetches the value of an environment variable or returns a fallback.
func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

// getEnvAsInt fetches the value of an environment variable as an integer or returns a fallback.
func getEnvAsInt(key string, fallback int) int {
	v := getEnv(key, "")
	if value, err := strconv.Atoi(v); err == nil {
		return value
	}
	return fallback
}
