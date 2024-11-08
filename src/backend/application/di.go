package application

import (
	"application/config"
	"application/dependency"
)

// Container is a struct that holds all the dependencies for the application.
// It acts as a central registry for services, ensuring that dependencies are managed in a lazy loaded manner.
type Container struct {
	Config dependency.LazyDependency[*config.Configuration]
}

// NewContainer creates and returns a new instance of Container.
func NewContainer() *Container {
	c := &Container{}

	c.Config = dependency.LazyDependency[*config.Configuration]{
		InitFunc: config.LoadConfig,
	}

	return c
}
