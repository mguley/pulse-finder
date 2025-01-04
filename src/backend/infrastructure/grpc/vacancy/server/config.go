package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Config holds the server configuration settings.
type Config struct {
	TLSEnabled   bool                // Whether TLS is enabled
	CertFile     string              // Path to the TLS certificate file
	KeyFile      string              // Path to the TLS key file
	Port         string              // Port the server listens on
	Interceptors []grpc.ServerOption // Interceptors and other gRPC server options
}

// Option defines a functional option for configuring the server.
type Option func(*Config)

// WithTLS enables TLS for the gRPC server and sets the certificate and key files.
func WithTLS(certFile, keyFile string) Option {
	return func(c *Config) {
		c.TLSEnabled = true
		c.CertFile = certFile
		c.KeyFile = keyFile
	}
}

// WithPort sets the server's listening port.
func WithPort(port string) Option {
	return func(c *Config) {
		c.Port = port
	}
}

// WithInterceptors adds interceptors and other options for the gRPC server.
func WithInterceptors(interceptors ...grpc.ServerOption) Option {
	return func(c *Config) {
		c.Interceptors = append(c.Interceptors, interceptors...)
	}
}

// NewGRPCServer initializes a gRPC server with the provided options.
func NewGRPCServer(opts ...Option) (*grpc.Server, *Config, error) {
	config := &Config{
		TLSEnabled:   false,
		Interceptors: []grpc.ServerOption{}, // Default to no interceptors
	}

	// Apply options to configure the server
	for _, opt := range opts {
		opt(config)
	}

	// gRPC server options
	var serverOpts []grpc.ServerOption

	// Add TLS credentials if enabled
	if config.TLSEnabled {
		cred, err := credentials.NewServerTLSFromFile(config.CertFile, config.KeyFile)
		if err != nil {
			return nil, nil, err
		}
		serverOpts = append(serverOpts, grpc.Creds(cred))
	}

	// Add interceptors if present
	serverOpts = append(serverOpts, config.Interceptors...)

	grpcServer := grpc.NewServer(serverOpts...)
	return grpcServer, config, nil
}
