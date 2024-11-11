package vacancy

import (
	"context"
	"domain/vacancy/entity"
	"domain/vacancy/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxVacancyRepository struct {
	repository.VacancyRepository
	db *pgxpool.Pool
}

func NewPgxVacancyRepository(db *pgxpool.Pool) *PgxVacancyRepository {
	return &PgxVacancyRepository{db: db}
}

func (r *PgxVacancyRepository) Save(ctx context.Context, v *entity.Vacancy) error {
	// todo
	return nil
}

func (r *PgxVacancyRepository) Get(ctx context.Context, id int64) (*entity.Vacancy, error) {
	// todo
	return nil, nil
}

func (r *PgxVacancyRepository) Update(ctx context.Context, v *entity.Vacancy) error {
	// todo
	return nil
}

func (r *PgxVacancyRepository) Delete(ctx context.Context, id int64) error {
	// todo
	return nil
}

func (r *PgxVacancyRepository) GetList(ctx context.Context) ([]*entity.Vacancy, error) {
	// todo
	return nil, nil
}
