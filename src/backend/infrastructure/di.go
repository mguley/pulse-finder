package infrastructure

import (
	"application/config"
	"application/dependency"
	appEvent "application/event"
	"infrastructure/event"
	"log"
)

// Container provides a lazily initialized set of dependencies for the infrastructure layer.
type Container struct {
	EventDispatcher dependency.LazyDependency[appEvent.Dispatcher]
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
	return c
}
