package handler

import (
	"application/auth"
	"application/config"
	"application/dependency"
	appEvent "application/event"
	"application/vacancy"
	"domain/vacancy/repository"
	"infrastructure/database"
	"infrastructure/event"
	vacancyHandler "infrastructure/grpc/vacancy/handler"
	"infrastructure/grpc/vacancy/validators"
	infraVacancy "infrastructure/vacancy"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TestContainer holds dependencies for the integration tests.
type TestContainer struct {
	Config               dependency.LazyDependency[*config.Configuration]
	JwtService           dependency.LazyDependency[*auth.Service]
	EventDispatcher      dependency.LazyDependency[appEvent.Dispatcher]
	DB                   dependency.LazyDependency[*pgxpool.Pool]
	VacancyRepository    dependency.LazyDependency[repository.VacancyRepository]
	VacancyService       dependency.LazyDependency[*vacancy.Service]
	Validator            dependency.LazyDependency[validators.Validator]
	VacancyServiceServer dependency.LazyDependency[*vacancyHandler.VacancyService]
}

// NewTestContainer initializes a new test container.
func NewTestContainer() *TestContainer {
	c := &TestContainer{}

	c.Config = dependency.LazyDependency[*config.Configuration]{
		InitFunc: config.LoadConfig,
	}
	c.JwtService = dependency.LazyDependency[*auth.Service]{
		InitFunc: func() *auth.Service { return auth.NewService(c.Config.Get()) },
	}
	c.EventDispatcher = dependency.LazyDependency[appEvent.Dispatcher]{
		InitFunc: func() appEvent.Dispatcher {
			instance, err := event.NewNatsEventDispatcher(c.Config.Get().Nats.URL)
			if err != nil {
				log.Fatalf("fail to create nats event dispatcher: %v", err)
			}
			return instance
		},
	}
	c.DB = dependency.LazyDependency[*pgxpool.Pool]{
		InitFunc: func() *pgxpool.Pool {
			instance, err := database.NewPostgresDB(c.Config.Get().DB.DSN)
			if err != nil {
				log.Fatalf("fail to create postgres db: %v", err)
			}
			return instance
		},
	}
	c.VacancyRepository = dependency.LazyDependency[repository.VacancyRepository]{
		InitFunc: func() repository.VacancyRepository {
			return infraVacancy.NewPgxVacancyRepository(c.DB.Get())
		},
	}
	c.VacancyService = dependency.LazyDependency[*vacancy.Service]{
		InitFunc: func() *vacancy.Service {
			return vacancy.NewService(c.VacancyRepository.Get(), c.EventDispatcher.Get())
		},
	}
	c.Validator = dependency.LazyDependency[validators.Validator]{
		InitFunc: func() validators.Validator {
			return validators.NewVacancyValidator()
		},
	}
	c.VacancyServiceServer = dependency.LazyDependency[*vacancyHandler.VacancyService]{
		InitFunc: func() *vacancyHandler.VacancyService {
			return vacancyHandler.NewVacancyService(c.VacancyService.Get(), c.Validator.Get())
		},
	}

	return c
}
