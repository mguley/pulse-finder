package validators

import (
	"interfaces/api/utils/validators"
	"interfaces/api/vacancy/dto"
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
	// todo add validation rules

	return v.Valid()
}
