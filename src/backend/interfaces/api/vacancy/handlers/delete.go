package handlers

import (
	"application/vacancy"
	"interfaces/api/utils"
	"net/http"
)

// DeleteVacancyHandler handles HTTP requests for deleting a job vacancy by its unique identifier.
type DeleteVacancyHandler struct {
	*utils.Handler   // HTTP handler utility.
	*utils.Errors    // Error handler for standardized error responses.
	*vacancy.Service // Vacancy service for business logic.
}

// NewDeleteVacancyHandler creates and returns a new instance of DeleteVacancyHandler.
func NewDeleteVacancyHandler(
	handler *utils.Handler,
	errors *utils.Errors,
	service *vacancy.Service,
) *DeleteVacancyHandler {
	return &DeleteVacancyHandler{
		Handler: handler,
		Errors:  errors,
		Service: service,
	}
}

// Execute processes the HTTP request to delete a job vacancy by its ID.
func (h *DeleteVacancyHandler) Execute(w http.ResponseWriter, r *http.Request) {
	id, err := h.ExtractId(r)
	if err != nil {
		h.NotFoundResponse(w, r)
		return
	}

	if err = h.Service.DeleteVacancy(r.Context(), id); err != nil {
		h.NotFoundResponse(w, r)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
