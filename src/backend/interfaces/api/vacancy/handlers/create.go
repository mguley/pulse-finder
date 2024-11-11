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

// CreateVacancyHandler handles the HTTP requests for creating a new job vacancy.
type CreateVacancyHandler struct {
	*utils.Handler               // HTTP handler utility.
	*utils.Errors                // Error handler for standardized error responses.
	*vacancy.Service             // Vacancy service for business logic.
	*validators.RequestValidator // Vacancy request validator.
}

// NewCreateVacancyHandler creates and returns a new instance of CreateVacancyHandler.
func NewCreateVacancyHandler(
	handler *utils.Handler,
	errors *utils.Errors,
	service *vacancy.Service,
	validator *validators.RequestValidator,
) *CreateVacancyHandler {
	return &CreateVacancyHandler{
		Handler:          handler,
		Errors:           errors,
		Service:          service,
		RequestValidator: validator,
	}
}

// Execute processes the HTTP request to create a new job vacancy.
func (h *CreateVacancyHandler) Execute(w http.ResponseWriter, r *http.Request) {
	// Parse and validate request
	request, err := h.parseAndValidateRequest(w, r)
	if err != nil {
		return
	}
	defer request.Release()

	// Map request DTO to Vacancy entity
	e := entity.GetVacancy()
	defer e.Release()
	request.ToEntity(e)

	// Save vacancy
	if err = h.Service.CreateVacancy(r.Context(), e); err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	// Send success response
	h.sendSuccessResponse(w, r, e)
}

// parseAndValidateRequest reads, parses, and validates the incoming JSON request body.
// Returns the validated request or an error if validation fails.
func (h *CreateVacancyHandler) parseAndValidateRequest(w http.ResponseWriter, r *http.Request) (*dto.Request, error) {
	request := dto.GetRequest()
	if err := h.ReadJson(w, r, &request); err != nil {
		h.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return nil, err
	}

	if !h.RequestValidator.Validate(request) {
		h.FailedValidationResponse(w, r, h.RequestValidator.Errors)
		h.RequestValidator.ClearErrors()
		return nil, fmt.Errorf("validation failed")
	}

	return request, nil
}

// sendSuccessResponse sends a success response with the created Vacancy's data.
func (h *CreateVacancyHandler) sendSuccessResponse(w http.ResponseWriter, r *http.Request, e *entity.Vacancy) {
	response := dto.GetResponse().FromEntity(e)
	defer response.Release()

	w.Header().Add("Location", fmt.Sprintf("/v1/vacancies/%d", e.GetId()))
	if err := h.WriteJson(w, http.StatusCreated, response, nil); err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}
