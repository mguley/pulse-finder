package handler

import (
	"application/auth"
	"context"
	"infrastructure/grpc/vacancy/interceptors"
	vacancyv1 "infrastructure/proto/vacancy/gen"
	"net"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SetupTestContainer initializes the TestContainer.
func SetupTestContainer(t *testing.T) (vacancyv1.VacancyServiceClient, *auth.Service) {
	container := NewTestContainer()

	// Create a listener for the in-process gRPC server
	listener, err := net.Listen("tcp", ":0") // Use a random available port
	require.NoError(t, err, "Failed to create listener")

	// Initialize the gRPC server with the JWT interceptor
	jwtService := container.JwtService.Get()
	server := grpc.NewServer(grpc.UnaryInterceptor(interceptors.JwtVacancyInterceptor(jwtService)))

	// Register the VacancyService
	vacancyService := container.VacancyServiceServer.Get()
	vacancyv1.RegisterVacancyServiceServer(server, vacancyService)

	// Start the server
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

	// Cleans up the database by truncating tables.
	t.Cleanup(func() {
		teardown(container.DB.Get(), t)
	})

	return vacancyv1.NewVacancyServiceClient(conn), jwtService
}

// teardown cleans up the database by truncating tables.
func teardown(db *pgxpool.Pool, t *testing.T) {
	ctx := context.Background()
	_, err := db.Exec(ctx, "TRUNCATE TABLE job_vacancies RESTART IDENTITY CASCADE;")
	require.NoError(t, err, "Failed to truncate job_vacancies")
}
