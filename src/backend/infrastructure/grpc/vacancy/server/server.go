package server

import (
	"errors"
	"fmt"
	vacancyv1 "infrastructure/proto/vacancy/gen"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

// VacancyServer is a high-level wrapper for the gRPC VacancyService server.
// It manages server initialization, registration and lifecycle.
type VacancyServer struct {
	grpcServer *grpc.Server // The gRPC server to serve RPC requests.
	listener   net.Listener // The network listener for the server.
	port       string       // The port the server listens on.
	env        string       // The environment (e.g., "prod" or "dev").
}

// NewVacancyServer creates a new instance of VacancyServer based on the provided configuration.
func NewVacancyServer(env, port, certFile, keyFile string) (*VacancyServer, error) {
	var (
		grpcServer   *grpc.Server
		serverConfig *Config
		listener     net.Listener
		err          error
	)

	switch env {
	case "prod":
		// Enable TLS in production
		grpcServer, serverConfig, err = NewGRPCServer(
			WithTLS(certFile, keyFile),
			WithPort(port))
	case "dev":
		grpcServer, serverConfig, err = NewGRPCServer(
			WithPort(port))
	default:
		return nil, errors.New("unsupported environment; must be \"prod\" or \"dev\"")
	}

	if err != nil {
		return nil, fmt.Errorf("create gRPC server: %w", err)
	}

	// Start listening
	listener, err = net.Listen("tcp", fmt.Sprintf(":%s", serverConfig.Port))
	if err != nil {
		return nil, fmt.Errorf("create listener: %w", err)
	}

	return &VacancyServer{
		grpcServer: grpcServer,
		listener:   listener,
		port:       port,
		env:        env,
	}, nil
}

// RegisterService registers the VacancyService implementation with the gRPC server.
func (s *VacancyServer) RegisterService(service vacancyv1.VacancyServiceServer) {
	vacancyv1.RegisterVacancyServiceServer(s.grpcServer, service)
}

// Start starts the gRPC server and begins listening for incoming requests.
func (s *VacancyServer) Start() {
	log.Printf("Starting the Vacancy gRPC server on %s (env: %s)...", s.listener.Addr(), s.env)
	go func() {
		if err := s.grpcServer.Serve(s.listener); err != nil {
			log.Printf("gRPC server failed to serve: %v", err)
		}
	}()
}

// WaitForShutdown gracefully shuts down the server upon receiving termination signals.
func (s *VacancyServer) WaitForShutdown() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	sig := <-signalChan
	log.Printf("Received signal: %s. Initiating shutdown...", sig)

	s.grpcServer.GracefulStop()
	log.Println("Vacancy gRPC server stopped gracefully.")
}
