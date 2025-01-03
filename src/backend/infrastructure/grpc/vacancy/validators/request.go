package validators

import (
	"fmt"
	vacancyv1 "infrastructure/proto/vacancy/gen"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Validator defines the interface for validating gRPC requests.
type Validator interface {
	ValidateCreateVacancyRequest(req *vacancyv1.CreateVacancyRequest) error
}

// VacancyValidator implements validation rules for gRPC vacancy requests.
type VacancyValidator struct {
	dateFormat string
}

// NewVacancyValidator creates a new instance of VacancyValidator.
func NewVacancyValidator() *VacancyValidator {
	return &VacancyValidator{dateFormat: "2006-01-02"}
}

// ValidateCreateVacancyRequest validates the CreateVacancyRequest fields.
func (v *VacancyValidator) ValidateCreateVacancyRequest(req *vacancyv1.CreateVacancyRequest) error {
	var validationErrors []error

	if err := validateStringField(req.Title, "title"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validateStringField(req.Company, "company"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validateStringField(req.Description, "description"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validateDateField(req.PostedAt, "posted_at", v.dateFormat); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validateStringField(req.Location, "location"); err != nil {
		validationErrors = append(validationErrors, err)
	}

	return combineErrors(validationErrors)
}

// validateStringField checks if a string field is provided and non-empty.
func validateStringField(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return status.Errorf(codes.InvalidArgument,
			"%s must be provided and cannot be empty or whitespace", fieldName)
	}
	return nil
}

// validateDateField checks if a date field is in the expected format.
func validateDateField(value, fieldName, format string) error {
	if strings.TrimSpace(value) == "" {
		return status.Errorf(codes.InvalidArgument,
			"%s must be provided and cannot be empty or whitespace", fieldName)
	}
	if _, err := time.Parse(format, value); err != nil {
		return status.Errorf(codes.InvalidArgument,
			"%s must be in the format %s. Example: %s",
			fieldName, format, time.Now().Format(format))
	}
	return nil
}

// combineErrors combines multiple errors into a single gRPC error.
func combineErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	var sb strings.Builder
	for i, err := range errs {
		sb.WriteString(fmt.Sprintf("Error %d: %v\n", i+1, err.Error()))
	}
	return status.Error(codes.InvalidArgument, sb.String())
}
