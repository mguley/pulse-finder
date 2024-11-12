package handlers

import (
	"application/vacancy"
	"domain/vacancy/entity"
	"interfaces/api/utils"
	"interfaces/api/vacancy/dto"
	"net/http"
)

// GetVacancyHandler handles HTTP requests for retrieving a job vacancy by its unique identifier.
type GetVacancyHandler struct {
	*utils.Handler   // HTTP handler utility.
	*utils.Errors    // Error handler for standardized error responses.
	*vacancy.Service // Vacancy service for business logic.
}

// NewGetVacancyHandler creates and returns a new instance of GetVacancyHandler.
func NewGetVacancyHandler(
	handler *utils.Handler,
	errors *utils.Errors,
	service *vacancy.Service,
) *GetVacancyHandler {
	return &GetVacancyHandler{
		Handler: handler,
		Errors:  errors,
		Service: service,
	}
}

// Execute processes the HTTP request to retrieve a job vacancy by its ID.
func (h *GetVacancyHandler) Execute(w http.ResponseWriter, r *http.Request) {
	id, err := h.ExtractId(r)
	if err != nil {
		h.NotFoundResponse(w, r)
		return
	}

	// Retrieve vacancy
	v, err := h.Service.GetVacancy(r.Context(), id)
	if err != nil {
		h.NotFoundResponse(w, r)
		return
	}

	// Send success response
	h.sendSuccessResponse(w, r, v)
}

// sendSuccessResponse sends a success response containing the retrieved Vacancy's data.
func (h *GetVacancyHandler) sendSuccessResponse(w http.ResponseWriter, r *http.Request, e *entity.Vacancy) {
	response := dto.GetResponse().FromEntity(e)
	defer response.Release()

	if err := h.WriteJson(w, http.StatusOK, response, nil); err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}
