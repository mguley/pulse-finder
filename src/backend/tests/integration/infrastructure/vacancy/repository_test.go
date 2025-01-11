package vacancy

import (
	"context"
	"domain/vacancy/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newVacancy creates a new vacancy entity for testing purposes.
func newVacancy(title, company, description, location string) *entity.Vacancy {
	return (&entity.Vacancy{}).
		SetTitle(title).
		SetCompany(company).
		SetDescription(description).
		SetPostedAt(time.Now()).
		SetLocation(location)
}

// TestPgxVacancyRepository_Save tests the repository's ability to save a vacancy to the database.
func TestPgxVacancyRepository_Save(t *testing.T) {
	c := SetupTestDatabase(t)
	r := c.Container.VacancyRepository.Get()

	v := newVacancy("Integration Test Engineer", "Tech Corp", "Responsible for integration testing", "Remote")
	err := r.Save(context.Background(), v)
	require.NoError(t, err, "Failed to save vacancy to the database")

	assert.NotZero(t, v.GetId(), "Expected vacancy ID to be non-zero after save")
	assert.Equal(t, int32(1), v.GetVersion(), "Expected version to be 1 for newly created vacancy")
}

// TestPgxVacancyRepository_Get tests fetching a saved vacancy by ID.
func TestPgxVacancyRepository_Get(t *testing.T) {
	c := SetupTestDatabase(t)
	r := c.Container.VacancyRepository.Get()
	ctx := context.Background()

	v := newVacancy("Software Engineer", "Tech Innovations", "Develop cutting-edge software solutions", "New York")
	err := r.Save(ctx, v)
	require.NoError(t, err, "Failed to save vacancy to the database")

	// Fetch the saved item by ID
	item, err := r.Get(ctx, v.GetId())
	require.NoError(t, err, "Failed to get vacancy from the database")

	assert.Equal(t, v.GetId(), item.GetId(), "Expected fetched ID to match saved ID")
	assert.Equal(t, v.GetTitle(), item.GetTitle(), "Expected fetched title to match saved title")
	assert.Equal(t, v.GetCompany(), item.GetCompany(), "Expected fetched company to match saved company")
	assert.Equal(t, v.GetDescription(), item.GetDescription(), "Expected fetched description to match saved description")
	assert.Equal(t, v.GetLocation(), item.GetLocation(), "Expected fetched location to match saved location")
}

// TestPgxVacancyRepository_Update tests updating an existing vacancy.
func TestPgxVacancyRepository_Update(t *testing.T) {
	c := SetupTestDatabase(t)
	r := c.Container.VacancyRepository.Get()
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

// TestPgxVacancyRepository_Update_VersionMismatch tests behavior when updating with a version mismatch.
func TestPgxVacancyRepository_Update_VersionMismatch(t *testing.T) {
	c := SetupTestDatabase(t)
	r := c.Container.VacancyRepository.Get()
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

// TestPgxVacancyRepository_Delete tests deleting a vacancy.
func TestPgxVacancyRepository_Delete(t *testing.T) {
	c := SetupTestDatabase(t)
	r := c.Container.VacancyRepository.Get()
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

// TestPgxVacancyRepository_Delete_NonExistentID tests deleting a non-existent vacancy.
func TestPgxVacancyRepository_Delete_NonExistentID(t *testing.T) {
	c := SetupTestDatabase(t)
	r := c.Container.VacancyRepository.Get()

	// Attempt to delete a non-existent vacancy
	err := r.Delete(context.Background(), 9999) // Assuming this ID does not exist
	assert.Error(t, err, "Expected an error for non-existent ID")
	assert.Contains(t, err.Error(), "does not exist", "Expected the error message to indicate non-existent ID")
}

// TestPgxVacancyRepository_GetList tests retrieving a list of vacancies.
func TestPgxVacancyRepository_GetList(t *testing.T) {
	c := SetupTestDatabase(t)
	r := c.Container.VacancyRepository.Get()
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

// TestPgxVacancyRepository_GetList_EmptyDatabase tests fetching from an empty database.
func TestPgxVacancyRepository_GetList_EmptyDatabase(t *testing.T) {
	c := SetupTestDatabase(t)
	r := c.Container.VacancyRepository.Get()

	// Fetch the list of vacancies from an empty table
	list, err := r.GetList(context.Background())
	require.NoError(t, err)
	assert.Empty(t, list, "Expected the vacancy list to be empty")
}

// TestPgxVacancyRepository_GetFilteredList tests fetching vacancies with filters.
func TestPgxVacancyRepository_GetFilteredList(t *testing.T) {
	c := SetupTestDatabase(t)
	r := c.Container.VacancyRepository.Get()
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

// TestPgxVacancyRepository_GetFilteredList_NoFilters tests fetching vacancies with no filters applied.
func TestPgxVacancyRepository_GetFilteredList_NoFilters(t *testing.T) {
	c := SetupTestDatabase(t)
	r := c.Container.VacancyRepository.Get()
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

// TestPgxVacancyRepository_Purge tests the repository's ability to purge all vacancies from the database.
func TestPgxVacancyRepository_Purge(t *testing.T) {
	c := SetupTestDatabase(t)
	r := c.Container.VacancyRepository.Get()
	ctx := context.Background()

	// Insert multiple test vacancies
	v1 := newVacancy("Software Engineer", "Tech Innovations", "Develop cutting-edge software solutions", "New York")
	v2 := newVacancy("Data Scientist", "AI Corp", "Analyze data and build AI models", "San Francisco")
	v3 := newVacancy("Product Manager", "Innovative Tech", "Oversee product development", "Remote")

	err := r.Save(ctx, v1)
	require.NoError(t, err, "Failed to save first vacancy")
	err = r.Save(ctx, v2)
	require.NoError(t, err, "Failed to save second vacancy")
	err = r.Save(ctx, v3)
	require.NoError(t, err, "Failed to save third vacancy")

	// Verify that the vacancies exist
	list, err := r.GetList(ctx)
	require.NoError(t, err, "Failed to fetch vacancy list before purge")
	assert.Len(t, list, 3, "Expected exactly 3 vacancies in the list before purge")

	// Perform the purge operation
	err = r.Purge(ctx)
	require.NoError(t, err, "Failed to purge vacancies")

	// Verify that the vacancies have been purged
	list, err = r.GetList(ctx)
	require.NoError(t, err, "Failed to fetch vacancy list after purge")
	assert.Empty(t, list, "Expected no vacancies in the list after purge")

	// Verify that primary key sequence is reset by inserting a new vacancy
	v4 := newVacancy("New Vacancy", "Fresh Start", "A new opportunity", "Anywhere")
	err = r.Save(ctx, v4)
	require.NoError(t, err, "Failed to save new vacancy after purge")
	assert.Equal(t, int64(1), v4.GetId(), "Expected new vacancy ID to start from 1 after purge")
}
