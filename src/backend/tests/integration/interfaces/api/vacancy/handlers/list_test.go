package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"testing"
	"tests"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestListVacancyHandler_Success tests successfully listing all job vacancies.
func TestListVacancyHandler_Success(t *testing.T) {
	testServer := SetupTestServer(t, func(router *httprouter.Router, container *tests.TestContainer) {
		router.HandlerFunc(http.MethodGet, "/v1/vacancies", container.ListHandler.Get().Execute)
	})
	defer testServer.Server.Close()

	ctx := context.Background()

	// Prepopulate the database with test vacancies
	v1 := newVacancy("Integration Test Engineer", "Tech Corp", "Responsible for integration testing", "Remote")
	v2 := newVacancy("Software Developer", "Innovatech", "Develop cutting-edge software", "New York")
	err := testServer.Container.VacancyRepository.Get().Save(ctx, v1)
	require.NoError(t, err)
	err = testServer.Container.VacancyRepository.Get().Save(ctx, v2)
	require.NoError(t, err)

	// Make the HTTP GET request to list vacancies
	resp, err := http.Get(testServer.Server.URL + "/v1/vacancies")
	require.NoError(t, err)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}()

	// Assert response status and body
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse and assert response body
	var response []map[string]any
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	// Assert the returned vacancies
	assert.Len(t, response, 2)
	assert.Equal(t, v2.GetTitle(), response[0]["title"])
	assert.Equal(t, v1.GetTitle(), response[1]["title"])
}

// TestListVacancyHandler_ValidationFailure tests validation errors.
func TestListVacancyHandler_ValidationFailure(t *testing.T) {
	testServer := SetupTestServer(t, func(router *httprouter.Router, container *tests.TestContainer) {
		router.HandlerFunc(http.MethodGet, "/v1/vacancies", container.ListHandler.Get().Execute)
	})
	defer testServer.Server.Close()

	// Make a request with an invalid page size
	resp, err := http.Get(testServer.Server.URL + "/v1/vacancies?page_size=-1")
	require.NoError(t, err)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}()

	// Assert response status and error message
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	var response map[string]any
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "size")
}

// TestListVacancyHandler_FilterByTitle tests filtering vacancies by title.
func TestListVacancyHandler_FilterByTitle(t *testing.T) {
	testServer := SetupTestServer(t, func(router *httprouter.Router, container *tests.TestContainer) {
		router.HandlerFunc(http.MethodGet, "/v1/vacancies", container.ListHandler.Get().Execute)
	})
	defer testServer.Server.Close()

	ctx := context.Background()

	// Prepopulate the database with test vacancies
	v1 := newVacancy("Integration Test Engineer", "Tech Corp", "Responsible for integration testing", "Remote")
	v2 := newVacancy("Software Developer", "Innovatech", "Develop cutting-edge software", "New York")
	err := testServer.Container.VacancyRepository.Get().Save(ctx, v1)
	require.NoError(t, err)
	err = testServer.Container.VacancyRepository.Get().Save(ctx, v2)
	require.NoError(t, err)

	// Make the HTTP GET request to filter by title (partial match)
	resp, err := http.Get(testServer.Server.URL + "/v1/vacancies?title=gin")
	require.NoError(t, err)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}()

	// Assert response status and body
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response []map[string]any
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	// Assert the returned vacancies
	assert.Len(t, response, 1)
	assert.Equal(t, v1.GetTitle(), response[0]["title"])
}

// TestListVacancyHandler_SortByTitle tests sorting vacancies by title.
func TestListVacancyHandler_SortByTitle(t *testing.T) {
	testServer := SetupTestServer(t, func(router *httprouter.Router, container *tests.TestContainer) {
		router.HandlerFunc(http.MethodGet, "/v1/vacancies", container.ListHandler.Get().Execute)
	})
	defer testServer.Server.Close()

	ctx := context.Background()

	// Prepopulate the database with test vacancies
	v1 := newVacancy("Software Developer", "Innovatech", "Develop cutting-edge software", "New York")
	v2 := newVacancy("Integration Test Engineer", "Tech Corp", "Responsible for integration testing", "Remote")
	err := testServer.Container.VacancyRepository.Get().Save(ctx, v1)
	require.NoError(t, err)
	err = testServer.Container.VacancyRepository.Get().Save(ctx, v2)
	require.NoError(t, err)

	// Make the HTTP GET request to sort by title
	resp, err := http.Get(testServer.Server.URL + "/v1/vacancies?sort_field=title&sort_order=asc")
	require.NoError(t, err)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}()

	// Assert response status and body
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response []map[string]any
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	// Assert the order of returned vacancies
	assert.Len(t, response, 2)
	assert.Equal(t, v2.GetTitle(), response[0]["title"]) // Alphabetically first
	assert.Equal(t, v1.GetTitle(), response[1]["title"])
}

// TestListVacancyHandler_Pagination tests pagination of results.
func TestListVacancyHandler_Pagination(t *testing.T) {
	testServer := SetupTestServer(t, func(router *httprouter.Router, container *tests.TestContainer) {
		router.HandlerFunc(http.MethodGet, "/v1/vacancies", container.ListHandler.Get().Execute)
	})
	defer testServer.Server.Close()

	// Prepopulate the database with multiple vacancies
	for i := 1; i <= 20; i++ {
		v := newVacancy(
			"Job "+strconv.Itoa(i),
			"Company "+strconv.Itoa(i),
			"Description for job "+strconv.Itoa(i),
			"Location "+strconv.Itoa(i),
		)
		err := testServer.Container.VacancyRepository.Get().Save(context.Background(), v)
		require.NoError(t, err)
	}

	// Make the HTTP GET request with pagination
	resp, err := http.Get(testServer.Server.URL + "/v1/vacancies?page=2&page_size=5&sort_order=asc")
	require.NoError(t, err)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}()

	// Assert response status and body
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response []map[string]any
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	// Assert the size and content of the returned page
	assert.Len(t, response, 5)
	assert.Equal(t, "Job 6", response[0]["title"])
	assert.Equal(t, "Job 10", response[4]["title"])
}

// TestListVacancyHandler_MaxPageSize tests that the maximum allowed page_size is accepted.
func TestListVacancyHandler_MaxPageSize(t *testing.T) {
	testServer := SetupTestServer(t, func(router *httprouter.Router, container *tests.TestContainer) {
		router.HandlerFunc(http.MethodGet, "/v1/vacancies", container.ListHandler.Get().Execute)
	})
	defer testServer.Server.Close()

	// Make the HTTP GET request with the maximum allowed page_size
	resp, err := http.Get(testServer.Server.URL + "/v1/vacancies?page=1&page_size=750")
	require.NoError(t, err)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}()

	// Assert response status and body
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response []map[string]any
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	// Assert the size of the returned page matches the request
	assert.LessOrEqual(t, len(response), 750)
}
