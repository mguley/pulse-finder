package application

import (
	"application/config"
	"application/dependency"
	diAuth "domain/auth"
	diHealthcheck "domain/healthcheck"
	diVacancy "domain/vacancy"
	"github.com/jackc/pgx/v5/pgxpool"
	diInfrastructure "infrastructure"
	"infrastructure/database"
	"interfaces/api/utils"
	"log"
	"log/slog"
	"os"
)

// Container is a struct that holds all the dependencies for the application.
// It acts as a central registry for services, ensuring that dependencies are managed in a lazy loaded manner.
type Container struct {
	Config                  dependency.LazyDependency[*config.Configuration]
	DB                      dependency.LazyDependency[*pgxpool.Pool]
	Handler                 dependency.LazyDependency[*utils.Handler]
	Errors                  dependency.LazyDependency[*utils.Errors]
	InfrastructureContainer dependency.LazyDependency[*diInfrastructure.Container]
	HealthCheckContainer    dependency.LazyDependency[*diHealthcheck.Container]
	JwtAuthContainer        dependency.LazyDependency[*diAuth.Container]
	VacancyContainer        dependency.LazyDependency[*diVacancy.Container]
}

// NewContainer creates and returns a new instance of Container.
// Each dependency is configured to initialize only when first accessed.
func NewContainer() *Container {
	container := &Container{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create container with base dependencies
	container.Config = dependency.LazyDependency[*config.Configuration]{
		InitFunc: config.LoadConfig,
	}
	container.Handler = dependency.LazyDependency[*utils.Handler]{
		InitFunc: utils.NewHandler,
	}
	container.Errors = dependency.LazyDependency[*utils.Errors]{
		InitFunc: func() *utils.Errors {
			return utils.NewErrors(logger, container.Handler.Get())
		},
	}

	// Database
	container.DB = dependency.LazyDependency[*pgxpool.Pool]{
		InitFunc: func() *pgxpool.Pool {
			db, err := database.NewPostgresDB(container.Config.Get().DB.DSN)
			if err != nil {
				log.Fatalf("Failed to initialize database: %v", err)
			}
			return db
		},
	}

	// Domain/layer containers
	container.InfrastructureContainer = dependency.LazyDependency[*diInfrastructure.Container]{
		InitFunc: func() *diInfrastructure.Container {
			return diInfrastructure.NewContainer(container.Config.Get())
		},
	}
	container.HealthCheckContainer = dependency.LazyDependency[*diHealthcheck.Container]{
		InitFunc: func() *diHealthcheck.Container {
			return diHealthcheck.NewContainer(
				container.Config.Get(),
				container.Handler.Get(),
				container.Errors.Get())
		},
	}
	container.JwtAuthContainer = dependency.LazyDependency[*diAuth.Container]{
		InitFunc: func() *diAuth.Container {
			return diAuth.NewContainer(
				container.Config.Get(),
				container.Handler.Get(),
				container.Errors.Get())
		},
	}
	container.VacancyContainer = dependency.LazyDependency[*diVacancy.Container]{
		InitFunc: func() *diVacancy.Container {
			return diVacancy.NewContainer(
				container.DB.Get(),
				container.InfrastructureContainer.Get().EventDispatcher.Get(),
				container.Handler.Get(),
				container.Errors.Get())
		},
	}

	return container
}
