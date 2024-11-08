package healthcheck

import (
	"application/config"
	"interfaces/api/healthcheck/dto"
	"runtime/debug"
	"time"
)

// Service provides health check details.
type Service struct {
	config *config.Configuration
}

// NewService creates a new instance of Service.
func NewService(c *config.Configuration) *Service {
	return &Service{config: c}
}

// Process populates the response with health check details.
func (s *Service) Process(r *dto.Response) {
	timestamp := time.Now().UTC().Format(time.DateTime)
	status := "available"

	// Populate the response with health check information
	r.Status = &status
	r.SystemInfo.Environment = &s.config.Env
	r.SystemInfo.Version = s.getRevision()
	r.Timestamp = &timestamp
}

// getRevision retrieves the VCS revision, if available.
func (s *Service) getRevision() *string {
	revision := "unknown"
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				revision = setting.Value
				break
			}
		}
	}

	return &revision
}
