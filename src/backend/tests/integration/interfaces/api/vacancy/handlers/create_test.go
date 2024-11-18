package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"tests"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateVacancyHandler_Success tests the successful creation of a job vacancy.
func TestCreateVacancyHandler_Success(t *testing.T) {
	testServer := SetupTestServer(t, func(router *httprouter.Router, container *tests.TestContainer) {
		router.HandlerFunc(http.MethodPost, "/v1/vacancies", container.CreateHandler.Get().Execute)
	})
	defer testServer.Server.Close()

	// Define a valid request payload
	payload := map[string]any{
		"title":       "Integration Test Engineer",
		"company":     "Tech Corp",
		"description": "Responsible for integration testing",
		"posted_at":   time.Now().Format(time.DateOnly),
		"location":    "Remote",
	}
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	// Make the HTTP POST request
	resp, err := http.Post(testServer.Server.URL+"/v1/vacancies", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}()

	// Assert response status and body
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Parse and assert response body
	var response map[string]any
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, payload["title"], response["title"])
	assert.Equal(t, payload["company"], response["company"])
	assert.Equal(t, payload["description"], response["description"])
	assert.Equal(t, payload["location"], response["location"])
}

// TestCreateVacancyHandler_ValidationFailure tests failure when required fields are missing.
func TestCreateVacancyHandler_ValidationFailure(t *testing.T) {
	testServer := SetupTestServer(t, func(router *httprouter.Router, container *tests.TestContainer) {
		router.HandlerFunc(http.MethodPost, "/v1/vacancies", container.CreateHandler.Get().Execute)
	})
	defer testServer.Server.Close()

	// Define an invalid request payload (missing required fields)
	payload := map[string]any{
		"company": "Tech Corp", // Missing title, description, etc.
	}
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	// Make the HTTP POST request
	resp, err := http.Post(testServer.Server.URL+"/v1/vacancies", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}()

	// Assert response status and error message
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	// Parse and assert response body
	var response map[string]any
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Len(t, response["error"], 4)
	assert.Contains(t, response["error"], "description")
	assert.Contains(t, response["error"], "location")
	assert.Contains(t, response["error"], "posted_at")
	assert.Contains(t, response["error"], "title")
}
