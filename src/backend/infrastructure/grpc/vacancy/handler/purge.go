package handler

import (
	"context"
	vacancyv1 "infrastructure/proto/vacancy/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PurgeVacancies handles the gRPC request to remove all job vacancies from the database.
// It processes the request and invokes the appropriate service logic.
func (s *VacancyService) PurgeVacancies(
	ctx context.Context,
	req *vacancyv1.PurgeVacanciesRequest,
) (*vacancyv1.PurgeVacanciesResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "context canceled")
	default:
		return s.handlePurgeVacancies(ctx)
	}
}

// handlePurgeVacancies performs the actual logic of purging all job vacancies.
func (s *VacancyService) handlePurgeVacancies(ctx context.Context) (*vacancyv1.PurgeVacanciesResponse, error) {
	if err := s.service.PurgeVacancies(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "purge vacancies: %v", err)
	}
	return &vacancyv1.PurgeVacanciesResponse{
		Message: codes.OK.String(),
	}, nil
}
