package handler

import (
	authv1 "infrastructure/proto/auth/gen"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SetupTestContainer initializes the TestContainer.
func SetupTestContainer(t *testing.T) authv1.AuthServiceClient {
	container := NewTestContainer()

	// Create a listener for the in-process gRPC server
	listener, err := net.Listen("tcp", ":0") // Use a random available port
	require.NoError(t, err, "Failed to create listener")

	// Initialize the gRPC server and register the AuthServiceServer
	server := grpc.NewServer()
	authService := container.AuthServiceServer.Get()
	authv1.RegisterAuthServiceServer(server, authService)

	// Start the server in a goroutine
	go func() {
		err = server.Serve(listener)
		require.NoError(t, err, "Failed to start gRPC server")
	}()

	// Ensure server stops after tests
	t.Cleanup(func() {
		server.GracefulStop()
	})

	// Set up a gRPC client to interact with the server
	target := listener.Addr().String()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient(target, opts...)
	require.NoError(t, err, "Failed to connect to gRPC server")

	// Ensure client connection is closed after tests
	t.Cleanup(func() {
		err = conn.Close()
		require.NoError(t, err, "Failed to close gRPC connection")
	})

	return authv1.NewAuthServiceClient(conn)
}
