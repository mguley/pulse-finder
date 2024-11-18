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

// createVacancy creates a vacancy.
func createVacancy(t *testing.T, testServer *TestServer, payload map[string]any) int {
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	resp, err := http.Post(testServer.Server.URL+"/v1/vacancies", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response map[string]any
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	id, ok := response["id"].(float64)
	require.True(t, ok, "expected valid ID in response")
	return int(id)
}

// TestDeleteVacancyHandler_Success tests the successful deletion of a job vacancy.
func TestDeleteVacancyHandler_Success(t *testing.T) {
	testServer := SetupTestServer(t, func(router *httprouter.Router, container *tests.TestContainer) {
		router.HandlerFunc(http.MethodPost, "/v1/vacancies", container.CreateHandler.Get().Execute)
		router.HandlerFunc(http.MethodDelete, "/v1/vacancies/:id", container.DeleteHandler.Get().Execute)
	})
	defer testServer.Server.Close()

	// Create a vacancy to delete
	payload := map[string]any{
		"title":       "Integration Test Engineer",
		"company":     "Tech Corp",
		"description": "Responsible for integration testing",
		"posted_at":   time.Now().Format(time.DateOnly),
		"location":    "Remote",
	}
	vacancyID := createVacancy(t, testServer, payload)

	// Make the HTTP DELETE request
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, testServer.Server.URL+"/v1/vacancies/"+strconv.Itoa(vacancyID), nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}()

	// Assert successful deletion
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Verify the vacancy no longer exists
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}()

	// Assert not found response
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// TestDeleteVacancyHandler_NotFound tests attempting to delete a non-existent vacancy.
func TestDeleteVacancyHandler_NotFound(t *testing.T) {
	testServer := SetupTestServer(t, func(router *httprouter.Router, container *tests.TestContainer) {
		router.HandlerFunc(http.MethodDelete, "/v1/vacancies/:id", container.DeleteHandler.Get().Execute)
	})
	defer testServer.Server.Close()

	// Attempt to delete a non-existent vacancy
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, testServer.Server.URL+"/v1/vacancies/9999", nil)
	require.NoError(t, err)

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

// TestDeleteVacancyHandler_InvalidID tests attempting to delete a vacancy with an invalid ID.
func TestDeleteVacancyHandler_InvalidID(t *testing.T) {
	testServer := SetupTestServer(t, func(router *httprouter.Router, container *tests.TestContainer) {
		router.HandlerFunc(http.MethodDelete, "/v1/vacancies/:id", container.DeleteHandler.Get().Execute)
	})
	defer testServer.Server.Close()

	// Attempt to delete with an invalid ID
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, testServer.Server.URL+"/v1/vacancies/invalid-id", nil)
	require.NoError(t, err)

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
