package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Config holds the server configuration settings.
type Config struct {
	TLSEnabled bool   // Whether TLS is enabled
	CertFile   string // Path to the TLS certificate file
	KeyFile    string // Path to the TLS key file
	Port       string // Port the server listens on
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

// NewGRPCServer initializes a gRPC server with the provided options.
func NewGRPCServer(opts ...Option) (*grpc.Server, *Config, error) {
	config := &Config{
		TLSEnabled: false,
	}

	// Apply options to configure the server
	for _, opt := range opts {
		opt(config)
	}

	// Create gRPC server with or without TLS
	var grpcServer *grpc.Server
	if config.TLSEnabled {
		cred, err := credentials.NewServerTLSFromFile(config.CertFile, config.KeyFile)
		if err != nil {
			return nil, nil, err
		}
		grpcServer = grpc.NewServer(grpc.Creds(cred))
	} else {
		grpcServer = grpc.NewServer()
	}

	return grpcServer, config, nil
}
