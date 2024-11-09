package application

import (
	"application/config"
	"application/dependency"
	diAuth "domain/auth"
	diHealthcheck "domain/healthcheck"
	"interfaces/api/utils"
	"log/slog"
	"os"
)

// Container is a struct that holds all the dependencies for the application.
// It acts as a central registry for services, ensuring that dependencies are managed in a lazy loaded manner.
type Container struct {
	Config               dependency.LazyDependency[*config.Configuration]
	Handler              dependency.LazyDependency[*utils.Handler]
	Errors               dependency.LazyDependency[*utils.Errors]
	HealthCheckContainer dependency.LazyDependency[*diHealthcheck.Container]
	JwtAuthContainer     dependency.LazyDependency[*diAuth.Container]
}

// NewContainer creates and returns a new instance of Container.
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

	// Domain containers
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

	return container
}
