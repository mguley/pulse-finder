package tests

import (
	"application/config"
	"application/dependency"
	"domain/vacancy/repository"
	"infrastructure/database"
	"infrastructure/vacancy"
	"interfaces/api/utils"
	"log"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TestContainer manages the dependencies for integration tests.
// It provides lazy initialization for core and domain specific dependencies.
type TestContainer struct {
	Config            dependency.LazyDependency[*config.Configuration]
	DB                dependency.LazyDependency[*pgxpool.Pool]
	Handler           dependency.LazyDependency[*utils.Handler]
	Errors            dependency.LazyDependency[*utils.Errors]
	VacancyRepository dependency.LazyDependency[repository.VacancyRepository]
}

// NewTestContainer creates a new instance of TestContainer.
func NewTestContainer() *TestContainer {
	c := &TestContainer{}

	initCoreDependencies(c)
	initVacancyDomainDependencies(c)

	return c
}

// initVacancyDomainDependencies initializes dependencies related to the vacancy domain.
func initVacancyDomainDependencies(c *TestContainer) {
	c.VacancyRepository = dependency.LazyDependency[repository.VacancyRepository]{
		InitFunc: func() repository.VacancyRepository {
			return vacancy.NewPgxVacancyRepository(c.DB.Get())
		},
	}
}

// initCoreDependencies initializes core application dependencies.
func initCoreDependencies(c *TestContainer) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	c.Config = dependency.LazyDependency[*config.Configuration]{
		InitFunc: config.LoadConfig,
	}
	c.DB = dependency.LazyDependency[*pgxpool.Pool]{
		InitFunc: func() *pgxpool.Pool {
			db, err := database.NewPostgresDB(c.Config.Get().DB.DSN)
			if err != nil {
				log.Fatalf("Failed to initialize database: %v", err)
			}
			return db
		},
	}
	c.Handler = dependency.LazyDependency[*utils.Handler]{
		InitFunc: utils.NewHandler,
	}
	c.Errors = dependency.LazyDependency[*utils.Errors]{
		InitFunc: func() *utils.Errors {
			return utils.NewErrors(logger, c.Handler.Get())
		},
	}
}
