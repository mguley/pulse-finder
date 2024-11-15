package handlers

import (
	"application/vacancy"
	"domain/vacancy/entity"
	"fmt"
	"interfaces/api/utils"
	"interfaces/api/vacancy/dto"
	"interfaces/api/vacancy/dto/list"
	"interfaces/api/vacancy/validators"
	"net/http"
)

// ListVacancyHandler handles the HTTP requests for listing job vacancies.
type ListVacancyHandler struct {
	*utils.Handler               // HTTP handler utility.
	*utils.Errors                // Error handler for standardized error responses.
	*vacancy.Service             // Vacancy service for business logic.
	*validators.RequestValidator // Vacancy request validator.
}

// NewListVacancyHandler creates and returns a new instance of ListVacancyHandler.
func NewListVacancyHandler(
	handler *utils.Handler,
	errors *utils.Errors,
	service *vacancy.Service,
	validator *validators.RequestValidator,
) *ListVacancyHandler {
	return &ListVacancyHandler{
		Handler:          handler,
		Errors:           errors,
		Service:          service,
		RequestValidator: validator,
	}
}

// Execute processes the HTTP request to list job vacancies.
func (h *ListVacancyHandler) Execute(w http.ResponseWriter, r *http.Request) {
	// Parse and validate the request
	rq, err := h.parseAndValidateRequest(w, r)
	if err != nil {
		return
	}
	defer rq.Release()

	// Fetch filtered vacancies
	items, err := h.Service.ListFilteredVacancies(r.Context(), *rq.Title, *rq.Company, *rq.Filters.Page,
		*rq.Filters.PageSize, *rq.Filters.SortField, *rq.Filters.SortOrder)
	if err != nil {
		h.NotFoundResponse(w, r)
		return
	}

	// Send success response
	h.sendSuccessResponse(w, r, items)
}

// parseAndValidateRequest reads, parses, and validates the incoming query parameters from the request URL.
// Returns a validated request DTO or an error if the validation fails.
func (h *ListVacancyHandler) parseAndValidateRequest(w http.ResponseWriter, r *http.Request) (*list.Request, error) {
	request := list.GetRequest()
	q := r.URL.Query()

	// Extract query parameters
	title := h.GetQueryString(q, "title", "")
	company := h.GetQueryString(q, "company", "")
	sortField := h.GetQueryString(q, "sort_field", "id")
	sortOrder := h.GetQueryString(q, "sort_order", "desc")
	page := h.GetQueryInt(q, "page", 1)
	pageSize := h.GetQueryInt(q, "page_size", 10)

	// Validate
	if !h.RequestValidator.ValidateFilters(page, pageSize, sortField) {
		h.FailedValidationResponse(w, r, h.RequestValidator.Errors)
		h.RequestValidator.ClearErrors()
		return nil, fmt.Errorf("validation failed")
	}

	// Populate the request DTO with the validated parameters
	request.Title = &title
	request.Company = &company
	request.Filters.Page = &page
	request.Filters.PageSize = &pageSize
	request.Filters.SortField = &sortField
	request.Filters.SortOrder = &sortOrder

	return request, nil
}

// sendSuccessResponse sends a success response with the list of vacancies.
func (h *ListVacancyHandler) sendSuccessResponse(w http.ResponseWriter, r *http.Request, data []*entity.Vacancy) {
	response := dto.GetResponse()
	defer response.Release()

	items := response.ToList(data)
	if err := h.WriteJson(w, http.StatusOK, items, nil); err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}
