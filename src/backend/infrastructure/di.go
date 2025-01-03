package infrastructure

import (
	"application/auth"
	"application/config"
	"application/dependency"
	appEvent "application/event"
	"application/vacancy"
	"domain/vacancy/repository"
	"infrastructure/database"
	"infrastructure/event"
	authHandler "infrastructure/grpc/auth/handler"
	authServer "infrastructure/grpc/auth/server"
	vacancyHandler "infrastructure/grpc/vacancy/handler"
	vacancyServer "infrastructure/grpc/vacancy/server"
	infraVacancy "infrastructure/vacancy"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Container provides a lazily initialized set of dependencies for the infrastructure layer.
type Container struct {
	EventDispatcher      dependency.LazyDependency[appEvent.Dispatcher]
	JwtAuthService       dependency.LazyDependency[*auth.Service]
	DB                   dependency.LazyDependency[*pgxpool.Pool]
	VacancyRepository    dependency.LazyDependency[repository.VacancyRepository]
	VacancyService       dependency.LazyDependency[*vacancy.Service]
	AuthServiceServer    dependency.LazyDependency[*authHandler.Service]
	AuthServer           dependency.LazyDependency[*authServer.AuthServer]
	VacancyServiceServer dependency.LazyDependency[*vacancyHandler.VacancyService]
	VacancyServer        dependency.LazyDependency[*vacancyServer.VacancyServer]
}

// NewContainer initializes and returns a new Container with lazy dependencies for the infrastructure layer.
func NewContainer(cfg *config.Configuration) *Container {
	c := &Container{
		EventDispatcher: dependency.LazyDependency[appEvent.Dispatcher]{
			InitFunc: func() appEvent.Dispatcher {
				d, err := event.NewNatsEventDispatcher(cfg.Nats.URL)
				if err != nil {
					log.Fatalf("Failed to initialize NATS event dispatcher: %v", err)
				}
				return d
			},
		},
	}
	c.DB = dependency.LazyDependency[*pgxpool.Pool]{
		InitFunc: func() *pgxpool.Pool {
			db, err := database.NewPostgresDB(cfg.DB.DSN)
			if err != nil {
				log.Fatalf("Failed to initialize database: %v", err)
			}
			return db
		},
	}
	c.VacancyRepository = dependency.LazyDependency[repository.VacancyRepository]{
		InitFunc: func() repository.VacancyRepository {
			return infraVacancy.NewPgxVacancyRepository(c.DB.Get())
		},
	}
	c.JwtAuthService = dependency.LazyDependency[*auth.Service]{
		InitFunc: func() *auth.Service { return auth.NewService(cfg) },
	}
	c.VacancyService = dependency.LazyDependency[*vacancy.Service]{
		InitFunc: func() *vacancy.Service {
			return vacancy.NewService(c.VacancyRepository.Get(), c.EventDispatcher.Get())
		},
	}

	// gRPC services
	c.AuthServiceServer = dependency.LazyDependency[*authHandler.Service]{
		InitFunc: func() *authHandler.Service {
			return authHandler.NewService(c.JwtAuthService.Get())
		},
	}
	c.AuthServer = dependency.LazyDependency[*authServer.AuthServer]{
		InitFunc: func() *authServer.AuthServer {
			var env, port, certFile, keyFile = cfg.Env, cfg.GRPC.AuthServerPort, cfg.TLSConfig.Certificate, cfg.TLSConfig.Key
			instance, err := authServer.NewAuthServer(env, port, certFile, keyFile)
			if err != nil {
				log.Fatalf("Failed to initialize gRPC Auth server: %v", err)
			}
			return instance
		},
	}
	c.VacancyServiceServer = dependency.LazyDependency[*vacancyHandler.VacancyService]{
		InitFunc: func() *vacancyHandler.VacancyService {
			return vacancyHandler.NewVacancyService(c.VacancyService.Get())
		},
	}
	c.VacancyServer = dependency.LazyDependency[*vacancyServer.VacancyServer]{
		InitFunc: func() *vacancyServer.VacancyServer {
			var env, port, certFile, keyFile = cfg.Env, cfg.GRPC.VacancyServerPort, cfg.TLSConfig.Certificate,
				cfg.TLSConfig.Key
			instance, err := vacancyServer.NewVacancyServer(env, port, certFile, keyFile)
			if err != nil {
				log.Fatalf("Failed to initialize gRPC Vacancy server: %v", err)
			}
			return instance
		},
	}

	return c
}
