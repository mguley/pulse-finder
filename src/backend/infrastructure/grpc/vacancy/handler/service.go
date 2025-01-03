package handler

import (
	"application/vacancy"
	"context"
	vacancyv1 "infrastructure/proto/vacancy/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VacancyService struct {
	vacancyv1.UnimplementedVacancyServiceServer
	service *vacancy.Service
}

func NewVacancyService(service *vacancy.Service) *VacancyService {
	return &VacancyService{service: service}
}

func (s *VacancyService) CreateVacancy(
	ctx context.Context,
	req *vacancyv1.CreateVacancyRequest,
) (*vacancyv1.CreateVacancyResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "context canceled")
	default:
		return s.process(req)
	}
}

func (s *VacancyService) process(req *vacancyv1.CreateVacancyRequest) (*vacancyv1.CreateVacancyResponse, error) {
	// todo
	return nil, nil
}
