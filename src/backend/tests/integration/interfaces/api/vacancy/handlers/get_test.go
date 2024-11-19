package handlers

import (
	"context"
	"domain/vacancy/entity"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"testing"
	"tests"
	"time"

	"github.com/julienschmidt/httprouter"
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

// TestGetVacancyHandler_Success tests successfully fetching a job vacancy by ID.
func TestGetVacancyHandler_Success(t *testing.T) {
	testServer := SetupTestServer(t, func(router *httprouter.Router, container *tests.TestContainer) {
		router.HandlerFunc(http.MethodGet, "/v1/vacancies/:id", container.GetHandler.Get().Execute)
	})
	defer testServer.Server.Close()

	// Prepopulate the database with a test vacancy
	v := newVacancy("Integration Test Engineer", "Tech Corp", "Responsible for integration testing", "Remote")
	err := testServer.Container.VacancyRepository.Get().Save(context.Background(), v)
	require.NoError(t, err, "Failed to save vacancy to the database")

	// Make the HTTP GET request
	resp, err := http.Get(testServer.Server.URL + "/v1/vacancies/" + strconv.FormatInt(v.GetId(), 10))
	require.NoError(t, err)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}()

	// Assert response status and body
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response should be OK")

	// Parse and assert response body
	var response map[string]any
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, v.GetTitle(), response["title"], "Title should match")
	assert.Equal(t, v.GetCompany(), response["company"], "Company should match")
	assert.Equal(t, v.GetDescription(), response["description"], "Description should match")
	assert.Equal(t, v.GetLocation(), response["location"], "Location should match")
}

// TestGetVacancyHandler_NotFound tests fetching a vacancy that does not exist.
func TestGetVacancyHandler_NotFound(t *testing.T) {
	testServer := SetupTestServer(t, func(router *httprouter.Router, container *tests.TestContainer) {
		router.HandlerFunc(http.MethodGet, "/v1/vacancies/:id", container.GetHandler.Get().Execute)
	})
	defer testServer.Server.Close()

	// Make the HTTP GET request for a non-existent ID
	resp, err := http.Get(testServer.Server.URL + "/v1/vacancies/9999")
	require.NoError(t, err)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}()

	// Assert response status and error message
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Response should be Not Found")

	// Parse and assert response body
	var response map[string]any
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Contains(t, response["error"], "Not Found")
}
