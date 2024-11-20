package handlers

import (
	"bytes"
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

// TestUpdateVacancyHandler_Success tests the successful update of a job vacancy.
func TestUpdateVacancyHandler_Success(t *testing.T) {
	testServer := SetupTestServer(t, func(router *httprouter.Router, container *tests.TestContainer) {
		router.HandlerFunc(http.MethodPost, "/v1/vacancies", container.CreateHandler.Get().Execute)
		router.HandlerFunc(http.MethodPatch, "/v1/vacancies/:id", container.UpdateHandler.Get().Execute)
	})
	defer testServer.Server.Close()

	// Create a vacancy to update
	payload := map[string]any{
		"title":       "Integration Test Engineer",
		"company":     "Tech Corp",
		"description": "Responsible for integration testing",
		"posted_at":   time.Now().Format(time.DateOnly),
		"location":    "Remote",
	}
	vacancyID := createVacancy(t, testServer, payload)

	// Prepare update payload
	updatePayload := map[string]any{
		"title":       "Senior Test Engineer",
		"description": "Responsible for advanced testing",
	}
	body, err := json.Marshal(updatePayload)
	require.NoError(t, err)

	// Make the HTTP PATCH request
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPatch, testServer.Server.URL+"/v1/vacancies/"+strconv.Itoa(vacancyID), bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}()

	// Assert successful update
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse and assert response body
	var response map[string]any
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, updatePayload["title"], response["title"])
	assert.Equal(t, updatePayload["description"], response["description"])
	assert.Equal(t, payload["company"], response["company"]) // Ensure unchanged fields are retained
}

// TestUpdateVacancyHandler_NotFound tests updating a non-existent vacancy.
func TestUpdateVacancyHandler_NotFound(t *testing.T) {
	testServer := SetupTestServer(t, func(router *httprouter.Router, container *tests.TestContainer) {
		router.HandlerFunc(http.MethodPatch, "/v1/vacancies/:id", container.UpdateHandler.Get().Execute)
	})
	defer testServer.Server.Close()

	// Prepare update payload
	updatePayload := map[string]any{
		"title":       "Senior Test Engineer",
		"description": "Responsible for advanced testing",
	}
	body, err := json.Marshal(updatePayload)
	require.NoError(t, err)

	// Make the HTTP PATCH request for a non-existent ID.
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPatch, testServer.Server.URL+"/v1/vacancies/9999", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}()

	// Assert not found response
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// TestUpdateVacancyHandler_ValidationFailure tests updating a vacancy with invalid data.
func TestUpdateVacancyHandler_ValidationFailure(t *testing.T) {
	testServer := SetupTestServer(t, func(router *httprouter.Router, container *tests.TestContainer) {
		router.HandlerFunc(http.MethodPost, "/v1/vacancies", container.CreateHandler.Get().Execute)
		router.HandlerFunc(http.MethodPatch, "/v1/vacancies/:id", container.UpdateHandler.Get().Execute)
	})
	defer testServer.Server.Close()

	// Create a vacancy to update
	payload := map[string]any{
		"title":       "Integration Test Engineer",
		"company":     "Tech Corp",
		"description": "Responsible for integration testing",
		"posted_at":   time.Now().Format(time.DateOnly),
		"location":    "Remote",
	}
	vacancyID := createVacancy(t, testServer, payload)

	// Prepare invalid update payload
	invalidPayload := map[string]any{
		"title": "", // Empty title is invalid
	}
	body, err := json.Marshal(invalidPayload)
	require.NoError(t, err)

	// Make the HTTP PATCH request
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPatch, testServer.Server.URL+"/v1/vacancies/"+strconv.Itoa(vacancyID), bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}()

	// Assert validation failure
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	// Parse and assert response body
	var response map[string]any
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Len(t, response["error"], 1)
	assert.Contains(t, response["error"], "title")
}

// TestUpdateVacancyHandler_PartialUpdate tests partial updates for a vacancy.
func TestUpdateVacancyHandler_PartialUpdate(t *testing.T) {
	testServer := SetupTestServer(t, func(router *httprouter.Router, container *tests.TestContainer) {
		router.HandlerFunc(http.MethodPost, "/v1/vacancies", container.CreateHandler.Get().Execute)
		router.HandlerFunc(http.MethodPatch, "/v1/vacancies/:id", container.UpdateHandler.Get().Execute)
	})
	defer testServer.Server.Close()

	// Create a vacancy to update
	payload := map[string]any{
		"title":       "Integration Test Engineer",
		"company":     "Tech Corp",
		"description": "Responsible for integration testing",
		"posted_at":   time.Now().Format(time.DateOnly),
		"location":    "Remote",
	}
	vacancyID := createVacancy(t, testServer, payload)

	// Prepare partial update payload
	partialPayload := map[string]any{
		"title": "Partial Update Engineer",
	}
	body, err := json.Marshal(partialPayload)
	require.NoError(t, err)

	// Make the HTTP PATCH request
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPatch, testServer.Server.URL+"/v1/vacancies/"+strconv.Itoa(vacancyID), bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}()

	// Assert successful update
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse and assert response body
	var response map[string]any
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, partialPayload["title"], response["title"])
	assert.Equal(t, payload["company"], response["company"])         // Ensure unchanged fields are retained
	assert.Equal(t, payload["description"], response["description"]) // Ensure unchanged fields are retained
}
