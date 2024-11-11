package vacancy

import (
	"context"
	"domain/vacancy/entity"
	"domain/vacancy/repository"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PgxVacancyRepository struct {
	repository.VacancyRepository
	db *pgxpool.Pool
}

func NewPgxVacancyRepository(db *pgxpool.Pool) *PgxVacancyRepository {
	return &PgxVacancyRepository{db: db}
}

// Save inserts a new item into the database and retrieves the generated ID and version.
func (r *PgxVacancyRepository) Save(ctx context.Context, v *entity.Vacancy) error {
	baseQuery := `
		INSERT INTO job_vacancies (title, company, description, posted_at, location)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, version
	`
	args := []any{v.GetTitle(), v.GetCompany(), v.GetDescription(), v.GetPostedAt(), v.GetLocation()}

	return r.withTransaction(ctx, func(tx pgx.Tx) error {
		var id int64
		var version int32

		// Execute the insert query and retrieve the generated id and version.
		if err := tx.QueryRow(ctx, baseQuery, args...).Scan(&id, &version); err != nil {
			return fmt.Errorf("failed to save vacancy: %w", err)
		}
		v.SetId(id).SetVersion(version)
		return nil
	})
}

// Get retrieves an item from the database by its ID.
func (r *PgxVacancyRepository) Get(ctx context.Context, id int64) (*entity.Vacancy, error) {
	baseQuery := `SELECT * FROM job_vacancies WHERE id = $1`
	row := r.db.QueryRow(ctx, baseQuery, id)
	v := &entity.Vacancy{}

	var itemId int64
	var title, company, description, location string
	var postedAt time.Time
	var version int32

	// Scan the row into vacancy fields.
	if err := row.Scan(&itemId, &title, &company, &description, &postedAt, &location, &version); err != nil {
		return nil, fmt.Errorf("failed to fetch vacancy: %w", err)
	}

	v.SetId(itemId).SetTitle(title).SetCompany(company).SetDescription(description).SetPostedAt(postedAt).
		SetLocation(location).SetVersion(version)
	return v, nil
}

// Update modifies an existing item in the database with new data, using optimistic concurrency control.
func (r *PgxVacancyRepository) Update(ctx context.Context, v *entity.Vacancy) error {
	baseQuery := `
		UPDATE job_vacancies
		SET title = $1, company = $2, description = $3, posted_at = $4, location = $5, version = version + 1
		WHERE id = $6 AND version = $7
		RETURNING version
	`
	args := []any{v.GetTitle(), v.GetCompany(), v.GetDescription(), v.GetPostedAt(), v.GetLocation(), v.GetId(), v.GetVersion()}

	return r.withTransaction(ctx, func(tx pgx.Tx) error {
		var version int32
		// Execute the update query and retrieve the new version.
		if err := tx.QueryRow(ctx, baseQuery, args...).Scan(&version); err != nil {
			return fmt.Errorf("failed to update vacancy: %w", err)
		}
		v.SetVersion(version)
		return nil
	})
}

// Delete removes an item from the database by its ID.
func (r *PgxVacancyRepository) Delete(ctx context.Context, id int64) error {
	baseQuery := `DELETE FROM job_vacancies WHERE id = $1`

	return r.withTransaction(ctx, func(tx pgx.Tx) error {
		commandTag, err := tx.Exec(ctx, baseQuery, id)
		if err != nil {
			return fmt.Errorf("failed to delete vacancy: %w", err)
		}
		if commandTag.RowsAffected() == 0 {
			return fmt.Errorf("vacancy with id %d does not exist", id)
		}
		return nil
	})
}

// GetList retrieves a list of items from the database.
func (r *PgxVacancyRepository) GetList(ctx context.Context) ([]*entity.Vacancy, error) {
	baseQuery := `SELECT * FROM job_vacancies`
	rows, err := r.db.Query(ctx, baseQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vacancies: %w", err)
	}
	defer rows.Close()

	var list []*entity.Vacancy
	for rows.Next() {
		var v entity.Vacancy
		var itemId int64
		var title, company, description, location string
		var postedAt time.Time
		var version int32

		// Scan the current row into vacancy fields.
		if err = rows.Scan(&itemId, &title, &company, &description, &postedAt, &location, &version); err != nil {
			return nil, fmt.Errorf("failed to scan vacancy: %w", err)
		}

		v.SetId(itemId).SetTitle(title).SetCompany(company).SetDescription(description).SetPostedAt(postedAt).
			SetLocation(location).SetVersion(version)
		list = append(list, &v)
	}

	// Check for row iteration errors.
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return list, nil
}

// withTransaction manages database transactions, allowing rollback on errors and commit on success.
func (r *PgxVacancyRepository) withTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	// Start a transaction.
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Execute the function within the transaction.
	if err = fn(tx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("transaction rollback failed: %w, original error: %v", rbErr, err)
		}
		return err
	}

	// Commit the transaction on success.
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
