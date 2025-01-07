package handler

import (
	"context"
	"domain/vacancy/entity"
	vacancyv1 "infrastructure/proto/vacancy/gen"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateVacancy handles the gRPC request to create a new job vacancy.
// It validates the request, processes the data, and saves the vacancy to the database.
func (s *VacancyService) CreateVacancy(
	ctx context.Context,
	req *vacancyv1.CreateVacancyRequest,
) (*vacancyv1.CreateVacancyResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "request canceled")
	default:
		return s.handleCreateVacancy(ctx, req)
	}
}

// handleCreateVacancy processes the CreateVacancyRequest by validating, mapping, and saving it.
func (s *VacancyService) handleCreateVacancy(
	ctx context.Context,
	req *vacancyv1.CreateVacancyRequest,
) (*vacancyv1.CreateVacancyResponse, error) {
	// Validate the request
	if err := s.validator.ValidateCreateVacancyRequest(req); err != nil {
		return nil, err
	}

	v := s.toEntity(req)
	defer v.Release()

	// Save vacancy
	if err := s.service.CreateVacancy(ctx, v); err != nil {
		return nil, status.Errorf(codes.Internal, "create vacancy: %v", err)
	}

	return s.toResponse(v), nil
}

// toEntity converts a gRPC CreateVacancyRequest into a domain Vacancy entity.
func (s *VacancyService) toEntity(req *vacancyv1.CreateVacancyRequest) *entity.Vacancy {
	v := entity.GetVacancy()

	v.SetTitle(req.Title).
		SetCompany(req.Company).
		SetDescription(req.Description).
		SetLocation(req.Location)

	if postedAt, err := time.Parse(s.dateFormat, req.PostedAt); err == nil {
		v.SetPostedAt(postedAt)
	}
	return v
}

// toResponse converts a domain Vacancy entity into a gRPC CreateVacancyResponse.
func (s *VacancyService) toResponse(e *entity.Vacancy) *vacancyv1.CreateVacancyResponse {
	return &vacancyv1.CreateVacancyResponse{
		Id:          e.GetId(),
		Title:       e.GetTitle(),
		Company:     e.GetCompany(),
		Description: e.GetDescription(),
		PostedAt:    e.GetPostedAt().Format(s.dateFormat),
		Location:    e.GetLocation(),
	}
}
