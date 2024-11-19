package vacancy

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"testing"
	"tests"
)

// TestContainer wraps the tests.TestContainer and adds convenience for accessing dependencies.
type TestContainer struct {
	Container *tests.TestContainer
	DB        *pgxpool.Pool
}

// SetupTestDatabase initializes a new test database and returns a TestContainer for managing dependencies.
func SetupTestDatabase(t *testing.T) *TestContainer {
	container := tests.NewTestContainer()

	t.Cleanup(func() {
		Teardown(container.DB.Get())
	})

	return &TestContainer{
		Container: container,
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
