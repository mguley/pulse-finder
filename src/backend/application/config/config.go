package config

import (
	"os"
	"strconv"
)

// Configuration holds the main application configuration settings.
type Configuration struct {
	Port      int            // Application server port.
	Env       string         // Environment (e.g., "development", "production").
	Jwt       JWTConfig      // Jwt configuration for authentication.
	DB        DatabaseConfig // Database configuration for connecting to the data source.
	Nats      NatsConfig     // NATS configuration.
	GRPC      GrpcConfig     // Configuration for gRPC server settings.
	TLSConfig TLSConfig      // Configuration for TLS settings.
}

// GrpcConfig holds settings for gRPC servers.
type GrpcConfig struct {
	AuthServerPort    string // Port for the Auth gRPC server.
	VacancyServerPort string // Port for the Vacancy gRPC server.
}

// TLSConfig holds settings for TLS.
type TLSConfig struct {
	Certificate string // Path to the TLS certificate file.
	Key         string // Path to the TLS key file.
}

// NatsConfig holds configuration settings for connecting to a NATS server.
type NatsConfig struct {
	URL string // URL of the NATS server.
}

// JWTConfig holds configuration settings for JWT-based authentication.
type JWTConfig struct {
	Secret string // Secret key for signing JWT tokens.
}

// DatabaseConfig holds settings for database connection.
type DatabaseConfig struct {
	DSN string // Data source name for database connection.
}

// LoadConfig loads the configuration settings from environment variables, falling back to default
// values if variables are not set.
func LoadConfig() *Configuration {
	config := &Configuration{
		Port: getEnvAsInt("PORT", 4005),
		Env:  getEnv("ENV", "dev"),
		Jwt: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
		},
		DB: DatabaseConfig{
			DSN: getEnv("DB_DSN", ""),
		},
		Nats: NatsConfig{
			URL: getEnv("NATS_URL", ""),
		},
		GRPC: GrpcConfig{
			AuthServerPort:    getEnv("GRPC_AUTH_SERVER_PORT", ""),
			VacancyServerPort: getEnv("GRPC_VACANCY_SERVER_PORT", ""),
		},
		TLSConfig: TLSConfig{
			Certificate: getEnv("TLS_CERTIFICATE", ""),
			Key:         getEnv("TLS_KEY", ""),
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
