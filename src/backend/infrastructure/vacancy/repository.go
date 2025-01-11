package vacancy

import (
	"context"
	"domain/vacancy/entity"
	"fmt"
	"infrastructure/persistence/criteria"
	"infrastructure/persistence/query"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PgxVacancyRepository implements the VacancyRepository interface using pgx.
type PgxVacancyRepository struct {
	db *pgxpool.Pool // Connection pool for database interactions.
}

// NewPgxVacancyRepository initializes a new instance of PgxVacancyRepository with a database connection pool.
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
	args := []any{v.GetTitle(), v.GetCompany(), v.GetDescription(), v.GetPostedAt(), v.GetLocation(), v.GetId(),
		v.GetVersion()}

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

// GetFilteredList retrieves a list of items based on the specified filter criteria.
func (r *PgxVacancyRepository) GetFilteredList(
	ctx context.Context,
	title, company string,
	page, pageSize int,
	sortField, sortOrder string,
) ([]*entity.Vacancy, error) {
	baseQuery := `SELECT * FROM job_vacancies`
	qb := query.GetBuilder(baseQuery)
	defer qb.Release()

	criteriaBuilder := criteria.GetSearchCriteriaBuilder()
	defer criteriaBuilder.Release()

	// Add filters based on provided parameters
	if title != "" {
		criteriaBuilder.AddFilter("title", "ILIKE", fmt.Sprintf("%%%s%%", title))
	}
	if company != "" {
		criteriaBuilder.AddFilter("company", "ILIKE", fmt.Sprintf("%%%s%%", company))
	}

	// Set logical operator for combining filters, default to "AND"
	criteriaBuilder.SetLogicalOperator("AND")
	searchCriteria := criteriaBuilder.Build()

	// Apply search criteria to the QueryBuilder
	qb.ApplySearchCriteria(searchCriteria)

	// Set sorting and pagination
	if sortField != "" {
		if sortOrder == "" {
			sortOrder = "ASC"
		}
		qb.SetOrderBy(sortField, sortOrder)
	} else {
		qb.SetOrderBy("title", "ASC") // Default sorting by title
	}
	qb.SetPagination(page, pageSize)

	// Build the final query and arguments
	q, args := qb.Build(searchCriteria)

	// Execute the query
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vacancies: %w", err)
	}
	defer rows.Close()

	// Process the results
	var list []*entity.Vacancy
	for rows.Next() {
		var v entity.Vacancy
		var id int64
		var title, company, description, location string
		var postedAt time.Time
		var version int32

		// Scan the row into vacancy fields
		if err = rows.Scan(&id, &title, &company, &description, &postedAt, &location, &version); err != nil {
			return nil, fmt.Errorf("failed to scan vacancy: %w", err)
		}

		// Map scanned values
		v.SetId(id).SetTitle(title).SetCompany(company).SetDescription(description).SetPostedAt(postedAt).
			SetLocation(location).SetVersion(version)
		list = append(list, &v)
	}

	// Check for row iteration errors.
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return list, nil
}

// Purge removes all job vacancies from the database by truncating the table.
func (r *PgxVacancyRepository) Purge(ctx context.Context) error {
	baseQuery := `TRUNCATE TABLE job_vacancies RESTART IDENTITY CASCADE`

	return r.withTransaction(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, baseQuery)
		if err != nil {
			return fmt.Errorf("failed to purge vacancies: %w", err)
		}
		return nil
	})
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
