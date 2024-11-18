package tests

import (
	"application/config"
	"application/dependency"
	appEvent "application/event"
	"application/vacancy"
	"domain/vacancy/repository"
	"infrastructure/database"
	"infrastructure/event"
	infraVacancy "infrastructure/vacancy"
	"interfaces/api/utils"
	"interfaces/api/vacancy/handlers"
	apiValidators "interfaces/api/vacancy/validators"
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
	EventDispatcher   dependency.LazyDependency[appEvent.Dispatcher]
	VacancyValidator  dependency.LazyDependency[*apiValidators.RequestValidator]
	VacancyService    dependency.LazyDependency[*vacancy.Service]
	CreateHandler     dependency.LazyDependency[*handlers.CreateVacancyHandler]
	DeleteHandler     dependency.LazyDependency[*handlers.DeleteVacancyHandler]
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
			return infraVacancy.NewPgxVacancyRepository(c.DB.Get())
		},
	}
	c.EventDispatcher = dependency.LazyDependency[appEvent.Dispatcher]{
		InitFunc: func() appEvent.Dispatcher {
			d, err := event.NewNatsEventDispatcher(c.Config.Get().Nats.URL)
			if err != nil {
				log.Fatalf("Failed to initialize NATS event dispatcher: %v", err)
			}
			return d
		},
	}
	c.VacancyValidator = dependency.LazyDependency[*apiValidators.RequestValidator]{
		InitFunc: apiValidators.NewRequestValidator,
	}
	c.VacancyService = dependency.LazyDependency[*vacancy.Service]{
		InitFunc: func() *vacancy.Service {
			return vacancy.NewService(c.VacancyRepository.Get(), c.EventDispatcher.Get())
		},
	}
	c.CreateHandler = dependency.LazyDependency[*handlers.CreateVacancyHandler]{
		InitFunc: func() *handlers.CreateVacancyHandler {
			return handlers.NewCreateVacancyHandler(
				c.Handler.Get(), c.Errors.Get(), c.VacancyService.Get(), c.VacancyValidator.Get())
		},
	}
	c.DeleteHandler = dependency.LazyDependency[*handlers.DeleteVacancyHandler]{
		InitFunc: func() *handlers.DeleteVacancyHandler {
			return handlers.NewDeleteVacancyHandler(c.Handler.Get(), c.Errors.Get(), c.VacancyService.Get())
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
