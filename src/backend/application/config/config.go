package config

import (
	"os"
	"strconv"
)

// Configuration holds the main application configuration settings.
type Configuration struct {
	Port int    // Application server port
	Env  string // Environment (e.g., "development", "production")
}

// LoadConfig loads the configuration settings from environment variables, falling back to default
// values if variables are not set.
func LoadConfig() *Configuration {
	config := &Configuration{
		Port: getEnvAsInt("PORT", 4005),
		Env:  getEnv("ENV", "development"),
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
