package handler

import (
	"application/auth"
	"application/config"
	"application/dependency"
	authHandler "infrastructure/grpc/auth/handler"
)

// TestContainer holds dependencies for the integration tests.
type TestContainer struct {
	Config            dependency.LazyDependency[*config.Configuration]
	JwtAuthService    dependency.LazyDependency[*auth.Service]
	AuthServiceServer dependency.LazyDependency[*authHandler.Service]
}

// NewTestContainer initializes a new test container.
func NewTestContainer() *TestContainer {
	c := &TestContainer{}

	c.Config = dependency.LazyDependency[*config.Configuration]{
		InitFunc: config.LoadConfig,
	}
	c.JwtAuthService = dependency.LazyDependency[*auth.Service]{
		InitFunc: func() *auth.Service { return auth.NewService(c.Config.Get()) },
	}
	c.AuthServiceServer = dependency.LazyDependency[*authHandler.Service]{
		InitFunc: func() *authHandler.Service {
			return authHandler.NewService(c.JwtAuthService.Get())
		},
	}

	return c
}
