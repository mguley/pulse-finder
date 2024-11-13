package validators

import (
	"interfaces/api/utils/validators"
	"interfaces/api/vacancy/dto"
	"strings"
	"time"
)

// RequestValidator is responsible for validating Vacancy request DTOs.
type RequestValidator struct {
	*validators.Validator // Embeds the general Validator to leverage its validation functions.
}

// NewRequestValidator creates and returns a new instance of RequestValidator.
// It retrieves a Validator instance from the pool for efficient memory usage.
func NewRequestValidator() *RequestValidator {
	return &RequestValidator{validators.GetValidator()}
}

// Validate performs the validation logic on the provided vacancy request DTO.
func (v *RequestValidator) Validate(r *dto.Request) bool {
	// Check if each required field is provided and non-empty
	v.checkRequiredField(r.Title, "title")
	v.checkRequiredField(r.Company, "company")
	v.checkRequiredField(r.Description, "description")
	v.checkRequiredField(r.Location, "location")

	// Validate PostedAt field
	v.validatePostedAt(r.PostedAt)

	return v.Valid()
}

// ValidateForUpdate performs validation on the vacancy request DTO for updating an existing vacancy.
func (v *RequestValidator) ValidateForUpdate(r *dto.Request) bool {
	if r.Title != nil {
		v.checkRequiredField(r.Title, "title")
	}
	if r.Company != nil {
		v.checkRequiredField(r.Company, "company")
	}
	if r.Description != nil {
		v.checkRequiredField(r.Description, "description")
	}
	if r.Location != nil {
		v.checkRequiredField(r.Location, "location")
	}
	if r.PostedAt != nil {
		v.validatePostedAt(r.PostedAt)
	}

	return v.Valid()
}

// checkRequiredField checks if a required field is provided and non-empty.
func (v *RequestValidator) checkRequiredField(field *string, fieldName string) {
	if field == nil || strings.TrimSpace(*field) == "" {
		v.AddError(fieldName, fieldName+" must be provided and cannot be empty or whitespace")
	}
}

// validatePostedAt checks if the PostedAt field is valid.
func (v *RequestValidator) validatePostedAt(postedAt *string) {
	if postedAt == nil || strings.TrimSpace(*postedAt) == "" {
		v.AddError("posted_at", "posted_at must be provided and cannot be empty or whitespace")
		return
	}

	_, err := time.Parse(time.DateOnly, strings.TrimSpace(*postedAt))
	if err != nil {
		v.AddError("posted_at", "posted_at must be a valid format. Example: 2006-01-02")
	}
}
