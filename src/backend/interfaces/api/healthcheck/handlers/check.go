package handlers

import (
	"application/healthcheck"
	"interfaces/api/healthcheck/dto"
	"interfaces/api/utils"
	"net/http"
)

// HealthCheckHandler handles HTTP requests for health checks.
type HealthCheckHandler struct {
	*utils.Handler       // HTTP handler utility
	*utils.Errors        // Error handling utility
	*healthcheck.Service // Health check service
}

// NewHealthCheckHandler creates a new HealthCheckHandler instance.
func NewHealthCheckHandler(
	handler *utils.Handler,
	errors *utils.Errors,
	service *healthcheck.Service,
) *HealthCheckHandler {
	return &HealthCheckHandler{
		Handler: handler,
		Errors:  errors,
		Service: service,
	}
}

// Execute processes a health check request and writes the JSON response.
func (h *HealthCheckHandler) Execute(w http.ResponseWriter, r *http.Request) {
	response := dto.GetResponse()
	defer response.Release()

	h.Service.Process(response)
	if err := h.WriteJson(w, http.StatusOK, response, nil); err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}
