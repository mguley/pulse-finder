package handler

import (
	"context"
	"fmt"
	vacancyv1 "infrastructure/proto/vacancy/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteVacancy handles the gRPC request to delete an existing job vacancy by its ID.
func (s *VacancyService) DeleteVacancy(
	ctx context.Context,
	req *vacancyv1.DeleteVacancyRequest,
) (*vacancyv1.DeleteVacancyResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "request canceled")
	default:
		return s.handleDeleteVacancy(ctx, req)
	}
}

// handleDeleteVacancy processes the DeleteVacancyRequest by validating and deleting the vacancy.
func (s *VacancyService) handleDeleteVacancy(
	ctx context.Context,
	req *vacancyv1.DeleteVacancyRequest,
) (*vacancyv1.DeleteVacancyResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid vacancy ID: %d", req.GetId())
	}

	if err := s.service.DeleteVacancy(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return &vacancyv1.DeleteVacancyResponse{
		Message: fmt.Sprintf("Vacancy with ID %d deleted", req.GetId()),
	}, nil
}
