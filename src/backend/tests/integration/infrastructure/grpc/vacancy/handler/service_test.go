package handler

import (
	"context"
	vacancyv1 "infrastructure/proto/vacancy/gen"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestVacancyService_CreateVacancy tests the CreateVacancy method of the VacancyServiceServer.
//
// This test covers the following scenarios:
// 1. A valid request with all required fields should successfully create a vacancy.
// 2. A request missing the required "Title" field should return an InvalidArgument error.
// 3. A request with an invalid date format for the "PostedAt" field should return an InvalidArgument error.
//
// The test uses the SetupTestContainer function to initialize the test dependencies and ensures proper cleanup of resources.
func TestVacancyService_CreateVacancy(t *testing.T) {
	client := SetupTestContainer(t)

	// Define test cases
	tests := []struct {
		name        string
		request     *vacancyv1.CreateVacancyRequest
		expectedErr bool
		errCode     codes.Code
	}{
		{
			name: "Valid Request",
			request: &vacancyv1.CreateVacancyRequest{
				Title:       "Software Engineer",
				Company:     "Tech Co.",
				Description: "Exciting opportunity in tech.",
				PostedAt:    "2025-01-01",
				Location:    "New York",
			},
			expectedErr: false,
		},
		{
			name: "Missing Title",
			request: &vacancyv1.CreateVacancyRequest{
				Company:     "Tech Co.",
				Description: "Exciting opportunity in tech.",
				PostedAt:    "2025-01-01",
				Location:    "Remote",
			},
			expectedErr: true,
			errCode:     codes.InvalidArgument,
		},
		{
			name: "Invalid Date Format",
			request: &vacancyv1.CreateVacancyRequest{
				Title:       "Software Engineer",
				Company:     "Tech Co.",
				Description: "Exciting opportunity in tech.",
				PostedAt:    "01-01-2025",
				Location:    "New York",
			},
			expectedErr: true,
			errCode:     codes.InvalidArgument,
		},
	}

	// Execute test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			// Make the gRPC call
			resp, err := client.CreateVacancy(ctx, tc.request)

			if tc.expectedErr {
				// Validate the expected error
				require.Error(t, err, "expected an error but got none")
				st, ok := status.FromError(err)
				require.True(t, ok, "error is not a gRPC status")
				assert.Equal(t, tc.errCode, st.Code(), "unexpected gRPC status code")
			} else {
				// Validate the response
				require.NoError(t, err, "unexpected error")
				assert.NotNil(t, resp, "response should not be nil")
				assert.NotEmpty(t, resp.Id, "response ID should not be empty")
			}
		})
	}
}
