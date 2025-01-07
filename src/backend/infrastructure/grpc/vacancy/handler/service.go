package handler

import (
	"application/vacancy"
	"infrastructure/grpc/vacancy/validators"
	vacancyv1 "infrastructure/proto/vacancy/gen"
)

// VacancyService implements the gRPC VacancyServiceServer interface.
// It provides methods to manage job vacancies via gRPC.
type VacancyService struct {
	vacancyv1.UnimplementedVacancyServiceServer
	service    *vacancy.Service     // Application service for managing vacancies.
	validator  validators.Validator // Validator for validating incoming gRPC requests.
	dateFormat string               // Date format used for parsing and formatting dates.
}

// NewVacancyService initializes and returns a new instance of VacancyService.
func NewVacancyService(service *vacancy.Service, validator validators.Validator) *VacancyService {
	return &VacancyService{
		service:    service,
		validator:  validator,
		dateFormat: "2006-01-02",
	}
}
