package handlers

import (
	"application/vacancy"
	"domain/vacancy/entity"
	"fmt"
	"interfaces/api/utils"
	"interfaces/api/vacancy/dto"
	"interfaces/api/vacancy/validators"
	"net/http"
)

// UpdateVacancyHandler handles HTTP requests for updating an existing vacancy by its ID.
type UpdateVacancyHandler struct {
	*utils.Handler               // HTTP handler utility.
	*utils.Errors                // Error handler for standardized error responses.
	*vacancy.Service             // Vacancy service for business logic.
	*validators.RequestValidator // Vacancy request validator.
}

// NewUpdateVacancyHandler creates and returns a new instance of UpdateVacancyHandler.
func NewUpdateVacancyHandler(
	handler *utils.Handler,
	errors *utils.Errors,
	service *vacancy.Service,
	validator *validators.RequestValidator,
) *UpdateVacancyHandler {
	return &UpdateVacancyHandler{
		Handler:          handler,
		Errors:           errors,
		Service:          service,
		RequestValidator: validator,
	}
}

// Execute processes the HTTP request to update an existing vacancy.
func (h *UpdateVacancyHandler) Execute(w http.ResponseWriter, r *http.Request) {
	// Parse and validate request
	request, err := h.parseAndValidateRequest(w, r)
	if err != nil {
		return
	}
	defer request.Release()

	// Retrieve vacancy by ID
	e, err := h.Service.GetVacancy(r.Context(), *request.ID)
	if err != nil {
		h.NotFoundResponse(w, r)
		return
	}
	request.ToEntity(e)

	// Save vacancy
	if err = h.Service.UpdateVacancy(r.Context(), e); err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	// Send success response
	h.sendSuccessResponse(w, r, e)
}

// parseAndValidateRequest parses the incoming request and validates it for updating a vacancy.
// Returns the validated request or an error if parsing/validation fails.
func (h *UpdateVacancyHandler) parseAndValidateRequest(w http.ResponseWriter, r *http.Request) (*dto.Request, error) {
	id, err := h.ExtractId(r)
	if err != nil {
		h.NotFoundResponse(w, r)
		return nil, err
	}

	request := dto.GetRequest()
	if err = h.ReadJson(w, r, &request); err != nil {
		h.FailedValidationResponse(w, r, map[string]string{"error": err.Error()})
		h.RequestValidator.ClearErrors()
		return nil, err
	}

	// Assign the extracted ID to the request
	request.ID = &id
	if !h.RequestValidator.ValidateForUpdate(request) {
		h.FailedValidationResponse(w, r, h.RequestValidator.Errors)
		h.RequestValidator.ClearErrors()
		return nil, fmt.Errorf("validation failed")
	}

	return request, nil
}

// sendSuccessResponse sends a success response with the updated vacancy data.
func (h *UpdateVacancyHandler) sendSuccessResponse(w http.ResponseWriter, r *http.Request, e *entity.Vacancy) {
	response := dto.GetResponse().FromEntity(e)
	defer response.Release()

	w.Header().Add("Location", fmt.Sprintf("/v1/vacancies/%d", e.GetId()))
	if err := h.WriteJson(w, http.StatusOK, response, nil); err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}
