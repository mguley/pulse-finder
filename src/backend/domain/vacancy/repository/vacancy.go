package repository

import (
	"context"
	"domain/vacancy/entity"
)

// VacancyRepository defines the interface for interacting with job vacancy data.
type VacancyRepository interface {
	// Save persists a new job vacancy into the data source.
	// Returns an error if the operation fails.
	Save(ctx context.Context, vacancy *entity.Vacancy) error

	// Get retrieves a job vacancy by its unique ID.
	// Returns a pointer to the Vacancy entity and an error if retrieval fails.
	Get(ctx context.Context, id int64) (*entity.Vacancy, error)

	// Update modifies an existing job vacancy in the data source.
	// Returns an error if the vacancy does not exist or if the operation fails.
	Update(ctx context.Context, vacancy *entity.Vacancy) error

	// Delete removes a job vacancy from the data source by its unique ID.
	// Returns an error if the vacancy does not exist or if the deletion operation fails.
	Delete(ctx context.Context, id int64) error

	// GetList retrieves a list of all job vacancies from the data source.
	// Returns a slice of Vacancy pointers and an error if the operation fails.
	GetList(ctx context.Context) ([]*entity.Vacancy, error)

	// GetFilteredList retrieves a list of job vacancies from the data source based on filter criteria.
	// Returns a slice of Vacancy pointers and an error if the operation fails.
	GetFilteredList(
		ctx context.Context,
		title, company string,
		page, pageSize int,
		sortField, sortOrder string) ([]*entity.Vacancy, error)

	// Purge removes all job vacancies from the data source.
	// Returns an error if the operation fails.
	Purge(ctx context.Context) error
}
