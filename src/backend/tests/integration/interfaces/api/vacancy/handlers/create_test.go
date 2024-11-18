package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"tests"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T) (container *tests.TestContainer, server *httptest.Server) {
	container = tests.NewTestContainer()

	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/v1/vacancies", container.CreateHandler.Get().Execute)
	server = httptest.NewServer(router)

	t.Cleanup(func() {
		server.Close()
		teardown(container.DB.Get())
	})
	return
}

func teardown(db *pgxpool.Pool) {
	ctx := context.Background()
	_, err := db.Exec(ctx, "TRUNCATE TABLE job_vacancies RESTART IDENTITY CASCADE;")
	if err != nil {
		log.Fatalf("failed to truncate job_vacancies: %v", err)
	}
}

func TestCreateVacancyHandler_Success(t *testing.T) {
	_, server := setupTestServer(t)

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
	resp, err := http.Post(server.URL+"/v1/vacancies", "application/json", bytes.NewBuffer(body))
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

func TestCreateVacancyHandler_ValidationFailure(t *testing.T) {
	_, server := setupTestServer(t)

	// Define an invalid request payload (missing required fields)
	payload := map[string]any{
		"company": "Tech Corp", // Missing title, description, etc.
	}
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	// Make the HTTP POST request
	resp, err := http.Post(server.URL+"/v1/vacancies", "application/json", bytes.NewBuffer(body))
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
