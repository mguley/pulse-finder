package handler

import (
	"application/vacancy"
	"context"
	"domain/vacancy/entity"
	"infrastructure/grpc/vacancy/validators"
	vacancyv1 "infrastructure/proto/vacancy/gen"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// CreateVacancy handles the gRPC request to create a new job vacancy.
// It validates the request, processes the data, and saves the vacancy to the database.
func (s *VacancyService) CreateVacancy(
	ctx context.Context,
	req *vacancyv1.CreateVacancyRequest,
) (*vacancyv1.CreateVacancyResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "context canceled")
	default:
		return s.handleCreateVacancy(ctx, req)
	}
}

// handleCreateVacancy processes the CreateVacancyRequest by validating, mapping and saving it.
func (s *VacancyService) handleCreateVacancy(
	ctx context.Context,
	req *vacancyv1.CreateVacancyRequest,
) (*vacancyv1.CreateVacancyResponse, error) {
	// Validate the request
	if err := s.validator.ValidateCreateVacancyRequest(req); err != nil {
		return nil, err
	}

	// Map the request to a domain entity
	v := s.mapToEntity(req)
	defer v.Release()

	// Save vacancy
	if err := s.service.CreateVacancy(ctx, v); err != nil {
		return nil, status.Errorf(codes.Internal, "create vacancy: %v", err)
	}

	return s.sendResponse(v), nil
}

// mapToEntity converts a gRPC CreateVacancyRequest into a domain Vacancy entity.
func (s *VacancyService) mapToEntity(req *vacancyv1.CreateVacancyRequest) *entity.Vacancy {
	v := entity.GetVacancy()

	v.SetTitle(req.Title).SetCompany(req.Company).SetDescription(req.Description).SetLocation(req.Location)
	if postedAt, err := time.Parse(s.dateFormat, req.PostedAt); err == nil {
		v.SetPostedAt(postedAt)
	}

	return v
}

// sendResponse converts a domain Vacancy entity into a gRPC CreateVacancyResponse.
func (s *VacancyService) sendResponse(e *entity.Vacancy) *vacancyv1.CreateVacancyResponse {
	return &vacancyv1.CreateVacancyResponse{
		Id:          e.GetId(),
		Title:       e.GetTitle(),
		Company:     e.GetCompany(),
		Description: e.GetDescription(),
		PostedAt:    e.GetPostedAt().Format(s.dateFormat),
		Location:    e.GetLocation(),
	}
}
