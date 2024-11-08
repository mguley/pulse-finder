package entity

// HealthCheck represents the status of the application's health.
type HealthCheck struct {
	status      string
	environment string
	version     string
}

// GetStatus returns the health status as a string.
func (e *HealthCheck) GetStatus() string {
	return e.status
}

// SetStatus sets the health check status.
func (e *HealthCheck) SetStatus(status string) *HealthCheck {
	e.status = status
	return e
}

// GetEnvironment returns the environment in which the application is running.
func (e *HealthCheck) GetEnvironment() string {
	return e.environment
}

// SetEnvironment sets the environment for the application (e.g., "production", "development")
func (e *HealthCheck) SetEnvironment(environment string) *HealthCheck {
	e.environment = environment
	return e
}

// GetVersion returns the current version of the application.
func (e *HealthCheck) GetVersion() string {
	return e.version
}

// SetVersion sets the version for the application.
func (e *HealthCheck) SetVersion(version string) *HealthCheck {
	e.version = version
	return e
}
