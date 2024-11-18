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

func newVacancy(title, company, description, location string) *entity.Vacancy {
	return (&entity.Vacancy{}).
		SetTitle(title).
		SetCompany(company).
		SetDescription(description).
		SetPostedAt(time.Now()).
		SetLocation(location)
}

func TestPgxVacancyRepository_Save(t *testing.T) {
	c := tests.NewTestContainer()
	db := c.DB.Get()
	defer teardown(db)

	r := c.VacancyRepository.Get()
	ctx := context.Background()
	v := newVacancy("Integration Test Engineer", "Tech Corp", "Responsible for integration testing", "Remote")

	err := r.Save(ctx, v)
	require.NoError(t, err, "Failed to save vacancy to the database")
	assert.NotZero(t, v.GetId(), "Expected vacancy ID to be non-zero after save")
	assert.Equal(t, int32(1), v.GetVersion(), "Expected version to be 1 for newly created vacancy")
}

func TestPgxVacancyRepository_Get(t *testing.T) {
	c := tests.NewTestContainer()
	db := c.DB.Get()
	defer teardown(db)

	r := c.VacancyRepository.Get()
	ctx := context.Background()
	v := newVacancy("Software Engineer", "Tech Innovations", "Develop cutting-edge software solutions", "New York")

	err := r.Save(ctx, v)
	require.NoError(t, err, "Failed to save test vacancy")

	// Fetch the saved item by ID
	result, err := r.Get(ctx, v.GetId())
	require.NoError(t, err, "Failed to fetch vacancy by ID")

	assert.Equal(t, v.GetId(), result.GetId(), "Expected fetched ID to match saved ID")
	assert.Equal(t, v.GetTitle(), result.GetTitle(), "Expected fetched title to match saved title")
	assert.Equal(t, v.GetCompany(), result.GetCompany(), "Expected fetched company to match saved company")
	assert.Equal(t, v.GetDescription(), result.GetDescription(), "Expected fetched description to match saved description")
	assert.Equal(t, v.GetLocation(), result.GetLocation(), "Expected fetched location to match saved location")
}

func TestPgxVacancyRepository_Update(t *testing.T) {
	c := tests.NewTestContainer()
	db := c.DB.Get()
	defer teardown(db)

	r := c.VacancyRepository.Get()
	ctx := context.Background()

	// Insert a test item into the database
	v := newVacancy("Software Engineer", "Tech Innovations", "Develop cutting-edge software solutions", "New York")
	err := r.Save(ctx, v)
	require.NoError(t, err, "Failed to save initial vacancy")

	// Update the vacancy's details
	v.SetTitle("Senior Software Engineer").
		SetCompany("Innovative Tech").
		SetDescription("Lead the development of innovative solutions").
		SetLocation("San Francisco")

	err = r.Update(ctx, v)
	require.NoError(t, err, "Failed to update the vacancy")

	// Fetch the updated vacancy
	updated, err := r.Get(ctx, v.GetId())
	require.NoError(t, err, "Failed to fetch updated vacancy")

	assert.Equal(t, "Senior Software Engineer", updated.GetTitle(), "Expected updated title to match")
	assert.Equal(t, "Innovative Tech", updated.GetCompany(), "Expected updated company to match")
	assert.Equal(t, "Lead the development of innovative solutions", updated.GetDescription(), "Expected updated description to match")
	assert.Equal(t, "San Francisco", updated.GetLocation(), "Expected updated location to match")
	assert.Equal(t, v.GetVersion(), updated.GetVersion(), "Expected version to increment after update")
}

func TestPgxVacancyRepository_Update_VersionMismatch(t *testing.T) {
	c := tests.NewTestContainer()
	db := c.DB.Get()
	defer teardown(db)

	r := c.VacancyRepository.Get()
	ctx := context.Background()

	// Insert a test vacancy into the database
	v := newVacancy("Software Engineer", "Tech Innovations", "Develop cutting-edge software solutions", "New York")
	err := r.Save(ctx, v)
	require.NoError(t, err, "Failed to save initial vacancy")

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
	require.NoError(t, err, "Failed to fetch original vacancy after failed update")
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
	v := newVacancy("Software Engineer", "Tech Innovations", "Develop cutting-edge software solutions", "New York")
	err := r.Save(ctx, v)
	require.NoError(t, err, "Failed to save vacancy for deletion test")

	// Delete the vacancy by ID
	err = r.Delete(ctx, v.GetId())
	require.NoError(t, err, "Failed to delete vacancy")

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
	v1 := newVacancy("Software Engineer", "Tech Innovations", "Develop cutting-edge software solutions", "New York")
	v2 := newVacancy("Data Scientist", "AI Corp", "Analyze data and build AI models", "San Francisco")

	err := r.Save(ctx, v1)
	require.NoError(t, err, "Failed to save first vacancy")
	err = r.Save(ctx, v2)
	require.NoError(t, err, "Failed to save second vacancy")

	// Fetch the list of vacancies
	list, err := r.GetList(ctx)
	require.NoError(t, err, "Failed to fetch vacancy list")
	assert.Len(t, list, 2, "Expected exactly 2 vacancies in the list")

	// Verify the contents of the list
	assert.Equal(t, v1.GetTitle(), list[0].GetTitle(), "Expected first vacancy title to match")
	assert.Equal(t, v2.GetTitle(), list[1].GetTitle(), "Expected second vacancy title to match")
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
	v1 := newVacancy("Software Engineer", "Tech Innovations", "Develop cutting-edge software solutions", "New York")
	v2 := newVacancy("Data Scientist", "AI Corp", "Analyze data and build AI models", "San Francisco")
	err := r.Save(ctx, v1)
	require.NoError(t, err, "Failed to save first vacancy")
	err = r.Save(ctx, v2)
	require.NoError(t, err, "Failed to save second vacancy")

	// Fetch filtered vacancies
	list, err := r.GetFilteredList(ctx, "Software", "Tech", 1, 10, "title", "ASC")
	require.NoError(t, err, "Failed to fetch filtered list")
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
	v1 := newVacancy("Software Engineer", "Tech Innovations", "Develop cutting-edge software solutions", "New York")
	v2 := newVacancy("Data Scientist", "AI Corp", "Analyze data and build AI models", "San Francisco")
	err := r.Save(ctx, v1)
	require.NoError(t, err, "Failed to save first vacancy")
	err = r.Save(ctx, v2)
	require.NoError(t, err, "Failed to save second vacancy")

	// Fetch all vacancies without filters
	list, err := r.GetFilteredList(ctx, "", "", 1, 10, "", "")
	require.NoError(t, err)
	assert.Len(t, list, 2, "Expected exactly 2 vacancies in the list")
}
