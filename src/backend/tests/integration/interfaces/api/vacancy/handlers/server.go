package handlers

import (
	"context"
	"log"
	"net/http/httptest"
	"testing"
	"tests"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

// TestServer contains the components for a test HTTP server and database.
type TestServer struct {
	Container *tests.TestContainer
	Server    *httptest.Server
	Router    *httprouter.Router
	DB        *pgxpool.Pool
}

// SetupTestServer initializes the test container and server with customizable routes.
func SetupTestServer(
	t *testing.T, configureRoutes func(router *httprouter.Router, container *tests.TestContainer)) *TestServer {
	container := tests.NewTestContainer()

	router := httprouter.New()
	configureRoutes(router, container)
	server := httptest.NewServer(router)

	t.Cleanup(func() {
		server.Close()
		Teardown(container.DB.Get())
	})

	return &TestServer{
		Container: container,
		Server:    server,
		Router:    router,
		DB:        container.DB.Get(),
	}
}

// Teardown cleans up the database by truncating tables.
func Teardown(db *pgxpool.Pool) {
	ctx := context.Background()
	_, err := db.Exec(ctx, "TRUNCATE TABLE job_vacancies RESTART IDENTITY CASCADE;")
	if err != nil {
		log.Fatalf("failed to truncate job_vacancies: %v", err)
	}
}
