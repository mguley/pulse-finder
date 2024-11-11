package vacancy

import (
	"application/event"
	"context"
	"domain/vacancy/entity"
	"domain/vacancy/events"
	"domain/vacancy/repository"
)

// Service provides application services for managing job vacancies.
type Service struct {
	repository repository.VacancyRepository
	dispatcher event.Dispatcher
}

// NewService initializes a new Service.
func NewService(r repository.VacancyRepository, d event.Dispatcher) *Service {
	return &Service{repository: r, dispatcher: d}
}

// CreateVacancy saves a new job vacancy to the database and dispatches a VacancyCreatedEvent.
// Returns an error if saving the vacancy or dispatching the event fails.
func (s *Service) CreateVacancy(ctx context.Context, v *entity.Vacancy) error {
	if err := s.repository.Save(ctx, v); err != nil {
		return err
	}
	e := events.NewVacancyCreatedEvent(v.GetId())
	if err := s.dispatcher.Dispatch(e); err != nil {
		return err
	}
	return nil
}

// UpdateVacancy updates an existing job vacancy in the database and dispatches a VacancyUpdatedEvent.
// Returns an error if updating the vacancy or dispatching the event fails.
func (s *Service) UpdateVacancy(ctx context.Context, v *entity.Vacancy) error {
	if err := s.repository.Update(ctx, v); err != nil {
		return err
	}
	e := events.NewVacancyUpdatedEvent(v.GetId())
	if err := s.dispatcher.Dispatch(e); err != nil {
		return err
	}
	return nil
}

// DeleteVacancy deletes an existing job vacancy from the database and dispatches a VacancyDeletedEvent.
// Returns an error if deleting the vacancy or dispatching the event fails.
func (s *Service) DeleteVacancy(ctx context.Context, id int64) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return err
	}
	e := events.NewVacancyDeletedEvent(id)
	if err := s.dispatcher.Dispatch(e); err != nil {
		return err
	}
	return nil
}

// GetVacancy retrieves a job vacancy by its unique ID from the database.
// Returns the vacancy if found, or an error if retrieval fails.
func (s *Service) GetVacancy(ctx context.Context, id int64) (*entity.Vacancy, error) {
	v, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// ListVacancies retrieves a list of all job vacancies from the database.
// Returns a slice of job vacancies or an error if retrieval fails.
func (s *Service) ListVacancies(ctx context.Context) ([]*entity.Vacancy, error) {
	list, err := s.repository.GetList(ctx)
	if err != nil {
		return nil, err
	}
	return list, nil
}