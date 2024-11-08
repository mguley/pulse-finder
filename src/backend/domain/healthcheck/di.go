package healthcheck

import (
	"application/config"
	"application/dependency"
	"application/healthcheck"
	"interfaces/api/healthcheck/handlers"
	"interfaces/api/utils"
)

// Container provides a lazily initialized set of dependencies for the health check domain.
type Container struct {
	HealthCheckService dependency.LazyDependency[*healthcheck.Service]
	HealthCheckHandler dependency.LazyDependency[*handlers.HealthCheckHandler]
}

// NewContainer initializes and returns a new Container with lazy dependencies for the health check domain.
func NewContainer(cfg *config.Configuration, h *utils.Handler, e *utils.Errors) *Container {
	c := &Container{
		HealthCheckService: dependency.LazyDependency[*healthcheck.Service]{
			InitFunc: func() *healthcheck.Service {
				return healthcheck.NewService(cfg)
			},
		},
	}
	c.HealthCheckHandler = dependency.LazyDependency[*handlers.HealthCheckHandler]{
		InitFunc: func() *handlers.HealthCheckHandler {
			return handlers.NewHealthCheckHandler(h, e, c.HealthCheckService.Get())
		},
	}
	return c
}
