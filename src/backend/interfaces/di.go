package interfaces

import (
	appAuth "application/auth"
	"application/config"
	"application/dependency"
	"interfaces/api/utils"
	"interfaces/middleware/auth"
)

// Container provides a lazily initialized set of dependencies for the interfaces layer.
type Container struct {
	JwtAuthService    dependency.LazyDependency[*appAuth.Service]
	JwtAuthMiddleware dependency.LazyDependency[*auth.JwtAuthMiddleware]
}

// NewContainer initializes and returns a new Container with lazy dependencies for the interfaces layer.
func NewContainer(cfg *config.Configuration, e *utils.Errors) *Container {
	c := &Container{
		JwtAuthService: dependency.LazyDependency[*appAuth.Service]{
			InitFunc: func() *appAuth.Service { return appAuth.NewService(cfg) },
		},
	}
	c.JwtAuthMiddleware = dependency.LazyDependency[*auth.JwtAuthMiddleware]{
		InitFunc: func() *auth.JwtAuthMiddleware {
			return auth.NewJwtAuthMiddleware(c.JwtAuthService.Get(), e)
		},
	}

	return c
}
