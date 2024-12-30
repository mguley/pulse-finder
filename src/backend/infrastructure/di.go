package infrastructure

import (
	"application/auth"
	"application/config"
	"application/dependency"
	appEvent "application/event"
	"infrastructure/event"
	authHandler "infrastructure/grpc/auth/handler"
	"infrastructure/grpc/auth/server"
	"log"
)

// Container provides a lazily initialized set of dependencies for the infrastructure layer.
type Container struct {
	EventDispatcher   dependency.LazyDependency[appEvent.Dispatcher]
	JwtAuthService    dependency.LazyDependency[*auth.Service]
	AuthServiceServer dependency.LazyDependency[*authHandler.Service]
	AuthServer        dependency.LazyDependency[*server.AuthServer]
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

	// gRPC services
	c.JwtAuthService = dependency.LazyDependency[*auth.Service]{
		InitFunc: func() *auth.Service { return auth.NewService(cfg) },
	}
	c.AuthServiceServer = dependency.LazyDependency[*authHandler.Service]{
		InitFunc: func() *authHandler.Service {
			return authHandler.NewService(c.JwtAuthService.Get())
		},
	}
	c.AuthServer = dependency.LazyDependency[*server.AuthServer]{
		InitFunc: func() *server.AuthServer {
			var env, port, certFile, keyFile = cfg.Env, cfg.GRPC.AuthServerPort, cfg.TLSConfig.Certificate, cfg.TLSConfig.Key
			authServer, err := server.NewAuthServer(env, port, certFile, keyFile)
			if err != nil {
				log.Fatalf("Failed to initialize gRPC Auth server: %v", err)
			}
			return authServer
		},
	}

	return c
}
