package vacancy

import (
	"context"
	"domain/vacancy/entity"
	"log"
	"testing"
	"tests"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func teardown(db *pgxpool.Pool) {
	ctx := context.Background()
	_, err := db.Exec(ctx, "TRUNCATE TABLE job_vacancies RESTART IDENTITY CASCADE;")
	if err != nil {
		log.Fatalf("failed to truncate job_vacancies: %v", err)
	}
}

func TestPgxVacancyRepository_Save(t *testing.T) {
	c := tests.NewTestContainer()
	db := c.DB.Get()
	defer teardown(db)

	r := c.VacancyRepository.Get()
	ctx := context.Background()

	v := &entity.Vacancy{}
	v.SetTitle("Integration Test Engineer").
		SetCompany("Tech Corp").
		SetDescription("Responsible for integration testing").
		SetPostedAt(time.Now()).
		SetLocation("Remote")

	err := r.Save(ctx, v)
	require.NoError(t, err)
	assert.NotZero(t, v.GetId(), "Expected non-zero ID after save")
	assert.Equal(t, int32(1), v.GetVersion(), "Expected version to be 1 after save")
}

func TestPgxVacancyRepository_Get(t *testing.T) {
	c := tests.NewTestContainer()
	db := c.DB.Get()
	defer teardown(db)

	r := c.VacancyRepository.Get()
	ctx := context.Background()

	// Insert a test item into the database
	v := &entity.Vacancy{}
	v.SetTitle("Software Engineer").
		SetCompany("Tech Innovations").
		SetDescription("Develop cutting-edge software solutions").
		SetPostedAt(time.Now()).
		SetLocation("New York")

	err := r.Save(ctx, v)
	require.NoError(t, err)

	// Fetch the saved item by ID
	result, err := r.Get(ctx, v.GetId())
	require.NoError(t, err)
	assert.Equal(t, v.GetId(), result.GetId(), "Expected the fetched ID to match the inserted ID")
	assert.Equal(t, v.GetTitle(), result.GetTitle(), "Expected the title to match")
	assert.Equal(t, v.GetCompany(), result.GetCompany(), "Expected the company to match")
	assert.Equal(t, v.GetDescription(), result.GetDescription(), "Expected the description to match")
	assert.Equal(t, v.GetLocation(), result.GetLocation(), "Expected the location to match")
}

func TestPgxVacancyRepository_Update(t *testing.T) {
	c := tests.NewTestContainer()
	db := c.DB.Get()
	defer teardown(db)

	r := c.VacancyRepository.Get()
	ctx := context.Background()

	// Insert a test item into the database
	v := &entity.Vacancy{}
	v.SetTitle("Software Engineer").
		SetCompany("Tech Innovations").
		SetDescription("Develop cutting-edge software solutions").
		SetPostedAt(time.Now()).
		SetLocation("New York")

	err := r.Save(ctx, v)
	require.NoError(t, err)

	// Update the vacancy's details
	v.SetTitle("Senior Software Engineer").
		SetCompany("Innovative Tech").
		SetDescription("Lead the development of innovative solutions").
		SetLocation("San Francisco")

	err = r.Update(ctx, v)
	require.NoError(t, err)

	// Fetch the updated vacancy
	updated, err := r.Get(ctx, v.GetId())
	require.NoError(t, err)
	assert.Equal(t, v.GetTitle(), updated.GetTitle(), "Expected the title to match the updated value")
	assert.Equal(t, v.GetCompany(), updated.GetCompany(), "Expected the company to match the updated value")
	assert.Equal(t, v.GetDescription(), updated.GetDescription(), "Expected the description to match the updated value")
	assert.Equal(t, v.GetLocation(), updated.GetLocation(), "Expected the location to match the updated value")
	assert.Equal(t, v.GetVersion(), updated.GetVersion(), "Expected the version to be incremented")
}

func TestPgxVacancyRepository_Update_VersionMismatch(t *testing.T) {
	c := tests.NewTestContainer()
	db := c.DB.Get()
	defer teardown(db)

	r := c.VacancyRepository.Get()
	ctx := context.Background()

	// Insert a test vacancy into the database
	v := &entity.Vacancy{}
	v.SetTitle("Software Engineer").
		SetCompany("Tech Innovations").
		SetDescription("Develop cutting-edge software solutions").
		SetPostedAt(time.Now()).
		SetLocation("New York")

	err := r.Save(ctx, v)
	require.NoError(t, err)

	// Simulate a version mismatch by manually incrementing the version
	v.SetVersion(v.GetVersion() + 1)

	// Attempt to update with the mismatched version
	v.SetTitle("Senior Software Engineer")
	err = r.Update(ctx, v)

	// Assert that an error occurred
	assert.Error(t, err, "Expected an error due to version mismatch")
	assert.Contains(t, err.Error(), "failed to update vacancy", "Expected the error message to indicate update failure")

	// Fetch the vacancy to verify it wasn't updated
	original, err := r.Get(ctx, v.GetId())
	require.NoError(t, err)
	assert.NotEqual(t, "Senior Software Engineer", original.GetTitle(), "Expected the title to remain unchanged")
	assert.Equal(t, int32(1), original.GetVersion(), "Expected the version to remain unchanged")
}

func TestPgxVacancyRepository_Delete(t *testing.T) {
	c := tests.NewTestContainer()
	db := c.DB.Get()
	defer teardown(db)

	r := c.VacancyRepository.Get()
	ctx := context.Background()

	// Insert a test item into the database
	v := &entity.Vacancy{}
	v.SetTitle("Software Engineer").
		SetCompany("Tech Innovations").
		SetDescription("Develop cutting-edge software solutions").
		SetPostedAt(time.Now()).
		SetLocation("New York")

	err := r.Save(ctx, v)
	require.NoError(t, err)

	// Delete the vacancy by ID
	err = r.Delete(ctx, v.GetId())
	require.NoError(t, err)

	// Verify the vacancy no longer exists
	_, err = r.Get(ctx, v.GetId())
	assert.Error(t, err, "Expected an error since the vacancy should no longer exist")
	assert.Contains(t, err.Error(), "failed to fetch vacancy", "Expected the error message to indicate fetch failure")
}

func TestPgxVacancyRepository_Delete_NonExistentID(t *testing.T) {
	c := tests.NewTestContainer()
	db := c.DB.Get()
	defer teardown(db)

	r := c.VacancyRepository.Get()
	ctx := context.Background()

	// Attempt to delete a non-existent vacancy
	err := r.Delete(ctx, 9999) // Assuming this ID does not exist
	assert.Error(t, err, "Expected an error for non-existent ID")
	assert.Contains(t, err.Error(), "does not exist", "Expected the error message to indicate non-existent ID")
}

func TestPgxVacancyRepository_GetList(t *testing.T) {
	c := tests.NewTestContainer()
	db := c.DB.Get()
	defer teardown(db)

	r := c.VacancyRepository.Get()
	ctx := context.Background()

	// Insert multiple vacancies into the database
	v1 := &entity.Vacancy{}
	v1.SetTitle("Software Engineer").
		SetCompany("Tech Innovations").
		SetDescription("Develop cutting-edge software solutions").
		SetPostedAt(time.Now()).
		SetLocation("New York")
	err := r.Save(ctx, v1)
	require.NoError(t, err)

	v2 := &entity.Vacancy{}
	v2.SetTitle("Data Scientist").
		SetCompany("AI Corp").
		SetDescription("Analyze data and build AI models").
		SetPostedAt(time.Now()).
		SetLocation("San Francisco")
	err = r.Save(ctx, v2)
	require.NoError(t, err)

	// Fetch the list of vacancies
	list, err := r.GetList(ctx)
	require.NoError(t, err)
	assert.Len(t, list, 2, "Expected exactly 2 vacancies in the list")

	// Verify the contents of the list
	assert.Equal(t, v1.GetTitle(), list[0].GetTitle(), "Expected the first vacancy title to match")
	assert.Equal(t, v2.GetTitle(), list[1].GetTitle(), "Expected the second vacancy title to match")
}

func TestPgxVacancyRepository_GetList_EmptyDatabase(t *testing.T) {
	c := tests.NewTestContainer()
	db := c.DB.Get()
	defer teardown(db)

	r := c.VacancyRepository.Get()
	ctx := context.Background()

	// Fetch the list of vacancies from an empty table
	list, err := r.GetList(ctx)
	require.NoError(t, err)
	assert.Empty(t, list, "Expected the vacancy list to be empty")
}

func TestPgxVacancyRepository_GetFilteredList(t *testing.T) {
	c := tests.NewTestContainer()
	db := c.DB.Get()
	defer teardown(db)

	r := c.VacancyRepository.Get()
	ctx := context.Background()

	// Insert test vacancies
	v1 := &entity.Vacancy{}
	v1.SetTitle("Software Engineer").
		SetCompany("Tech Innovations").
		SetDescription("Develop cutting-edge software solutions").
		SetPostedAt(time.Now()).
		SetLocation("New York")
	err := r.Save(ctx, v1)
	require.NoError(t, err)

	v2 := &entity.Vacancy{}
	v2.SetTitle("Data Scientist").
		SetCompany("AI Corp").
		SetDescription("Analyze data and build AI models").
		SetPostedAt(time.Now()).
		SetLocation("San Francisco")
	err = r.Save(ctx, v2)
	require.NoError(t, err)

	// Fetch filtered vacancies
	list, err := r.GetFilteredList(ctx, "Software", "Tech", 1, 10, "title", "ASC")
	require.NoError(t, err)
	assert.Len(t, list, 1, "Expected exactly 1 vacancy in the filtered list")

	// Verify the filtered result
	assert.Equal(t, "Software Engineer", list[0].GetTitle(), "Expected the title to match")
	assert.Equal(t, "Tech Innovations", list[0].GetCompany(), "Expected the company to match")
}

func TestPgxVacancyRepository_GetFilteredList_NoFilters(t *testing.T) {
	c := tests.NewTestContainer()
	db := c.DB.Get()
	defer teardown(db)

	r := c.VacancyRepository.Get()
	ctx := context.Background()

	// Insert test vacancies
	v1 := &entity.Vacancy{}
	v1.SetTitle("Software Engineer").
		SetCompany("Tech Innovations").
		SetDescription("Develop cutting-edge software solutions").
		SetPostedAt(time.Now()).
		SetLocation("New York")
	err := r.Save(ctx, v1)
	require.NoError(t, err)

	v2 := &entity.Vacancy{}
	v2.SetTitle("Data Scientist").
		SetCompany("AI Corp").
		SetDescription("Analyze data and build AI models").
		SetPostedAt(time.Now()).
		SetLocation("San Francisco")
	err = r.Save(ctx, v2)
	require.NoError(t, err)

	// Fetch all vacancies without filters
	list, err := r.GetFilteredList(ctx, "", "", 1, 10, "", "")
	require.NoError(t, err)
	assert.Len(t, list, 2, "Expected exactly 2 vacancies in the list")
}
