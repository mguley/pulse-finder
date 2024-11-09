package auth

import (
	"application/auth"
	"application/config"
	"application/dependency"
	"interfaces/api/auth/handlers"
	"interfaces/api/utils"
)

// Container provides a lazily initialized set of dependencies for the auth domain.
type Container struct {
	JwtAuthService dependency.LazyDependency[*auth.Service]
	JwtAuthHandler dependency.LazyDependency[*handlers.JwtTokenHandler]
}

// NewContainer initializes and returns a new Container with lazy dependencies for the auth domain.
func NewContainer(cfg *config.Configuration, h *utils.Handler, e *utils.Errors) *Container {
	c := &Container{
		JwtAuthService: dependency.LazyDependency[*auth.Service]{
			InitFunc: func() *auth.Service { return auth.NewService(cfg) },
		},
	}
	c.JwtAuthHandler = dependency.LazyDependency[*handlers.JwtTokenHandler]{
		InitFunc: func() *handlers.JwtTokenHandler {
			return handlers.NewJwtTokenHandler(h, e, c.JwtAuthService.Get())
		},
	}

	return c
}
