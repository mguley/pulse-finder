package validators

import (
	"interfaces/api/utils/validators"
	"interfaces/api/vacancy/dto"
	"strings"
	"time"
)

// RequestValidator is responsible for validating Vacancy request DTOs.
type RequestValidator struct {
	*validators.Validator          // Embeds the general Validator to leverage its validation functions.
	list                  []string // A list of fields allowed to be validated.
}

// NewRequestValidator creates and returns a new instance of RequestValidator.
// It retrieves a Validator instance from the pool for efficient memory usage.
func NewRequestValidator() *RequestValidator {
	return &RequestValidator{
		Validator: validators.GetValidator(),
		list:      []string{"id", "title", "company"},
	}
}

// Validate performs validation on the provided vacancy request DTO for creating a new vacancy.
func (v *RequestValidator) Validate(r *dto.Request) bool {
	return v.performValidation(r, true)
}

// ValidateForUpdate performs validation on the provided vacancy request DTO for updating an existing vacancy.
// It checks only the non-nil fields, allowing partial updates.
func (v *RequestValidator) ValidateForUpdate(r *dto.Request) bool {
	return v.performValidation(r, false)
}

// ValidateFilters validates the pagination and sort filter fields.
func (v *RequestValidator) ValidateFilters(page, size int, sort string) bool {
	v.Check(page > 0, "page", "page must be greater than zero")
	v.Check(page <= 1_000, "page", "page must be less than or equal to 1_000")

	v.Check(size > 0, "size", "size must be greater than zero")
	v.Check(size <= 750, "size", "size must be less than or equal to 750")

	if sort != "" {
		v.Check(v.PermittedValue(sort, v.list...), "sort", "sort contains an invalid value")
	}
	return v.Valid()
}

// performValidation performs the common validation logic for both creation and update scenarios.
func (v *RequestValidator) performValidation(r *dto.Request, checkRequired bool) bool {
	v.validateField(r.Title, "title", checkRequired)
	v.validateField(r.Company, "company", checkRequired)
	v.validateField(r.Description, "description", checkRequired)
	v.validateField(r.Location, "location", checkRequired)
	v.validatePostedAt(r.PostedAt, checkRequired)

	return v.Valid()
}

// validateField checks if a field is provided and non-empty based on the given flag.
func (v *RequestValidator) validateField(field *string, fieldName string, checkRequired bool) {
	if checkRequired || field != nil {
		if field == nil || strings.TrimSpace(*field) == "" {
			v.AddError(fieldName, fieldName+" must be provided and cannot be empty or whitespace")
		}
	}
}

// validatePostedAt checks if the PostedAt field is valid.
func (v *RequestValidator) validatePostedAt(postedAt *string, checkRequired bool) {
	if checkRequired || postedAt != nil {
		if postedAt == nil || strings.TrimSpace(*postedAt) == "" {
			v.AddError("posted_at", "posted_at must be provided and cannot be empty or whitespace")
			return
		}

		_, err := time.Parse(time.DateOnly, strings.TrimSpace(*postedAt))
		if err != nil {
			v.AddError("posted_at", "posted_at must be in the format YYYY-MM-DD. Example: 2006-01-02")
		}
	}
}
