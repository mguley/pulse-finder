package handler

import (
	"context"
	"domain/auth/entity"
	"fmt"
	vacancyv1 "infrastructure/proto/vacancy/gen"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// TestVacancyService_CreateVacancy tests the CreateVacancy method of the VacancyServiceServer.
//
// This test covers the following scenarios:
// 1. A valid request with a valid token should successfully create a vacancy.
// 2. A request without an Authorization header should return an Unauthenticated error.
// 3. A request with an invalid token format should return an Unauthenticated error.
// 4. A request with an expired token should return an Unauthenticated error.
// 5. A request missing the required "Title" field should return an InvalidArgument error.
// 6. A request with an invalid date format for the "PostedAt" field should return an InvalidArgument error.
//
// The test uses the SetupTestContainer function to initialize the test dependencies and ensures proper cleanup of resources.
func TestVacancyService_CreateVacancy(t *testing.T) {
	client, jwtService := SetupTestContainer(t)

	// Generate a valid token
	validClaims := entity.GetTokenClaims().
		SetIssuer("test-issuer").
		SetScope([]string{"test"}).
		SetExpiresAt(time.Now().Add(time.Hour).Unix())
	validToken, err := jwtService.Generate(validClaims)
	require.NoError(t, err, "could not generate token")

	// Generate an expired token
	expiredClaims := entity.GetTokenClaims().
		SetIssuer("test-issuer").
		SetScope([]string{"test"}).
		SetExpiresAt(time.Now().Add(-time.Hour).Unix())
	expiredToken, err := jwtService.Generate(expiredClaims)
	require.NoError(t, err, "could not generate expired token")

	// Define test cases
	tests := []struct {
		name        string
		token       string
		request     *vacancyv1.CreateVacancyRequest
		expectedErr bool
		errCode     codes.Code
	}{
		{
			name:  "Valid Request with Valid Token",
			token: validToken,
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
			name:        "Missing Authorization Header",
			token:       "",
			request:     &vacancyv1.CreateVacancyRequest{},
			expectedErr: true,
			errCode:     codes.Unauthenticated,
		},
		{
			name:        "Invalid Token Format",
			token:       "InvalidTokenFormat",
			request:     &vacancyv1.CreateVacancyRequest{},
			expectedErr: true,
			errCode:     codes.Unauthenticated,
		},
		{
			name:  "Expired Token",
			token: expiredToken,
			request: &vacancyv1.CreateVacancyRequest{
				Title:       "Expired Token Test",
				Company:     "Tech Co.",
				Description: "Exciting opportunity in tech.",
				PostedAt:    "2025-01-01",
				Location:    "New York",
			},
			expectedErr: true,
			errCode:     codes.Unauthenticated,
		},
		{
			name:  "Missing Title",
			token: validToken,
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
			name:  "Invalid Date Format",
			token: validToken,
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

			// Add the token to the metadata
			md := map[string]string{}
			if tc.token != "" {
				md["authorization"] = fmt.Sprintf("Bearer %s", tc.token)
			}
			ctx = metadata.NewOutgoingContext(ctx, metadata.New(md))

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

// TestVacancyService_DeleteVacancy tests the DeleteVacancy method of the VacancyServiceServer.
//
// This test covers the following scenarios:
// 1. A valid request with a correct ID and valid token should delete the vacancy successfully.
// 2. A request with an invalid ID (negative value) should return an InvalidArgument error.
// 3. A request without an Authorization header should return an Unauthenticated error.
// 4. A request with a non-existent ID should return an Internal error.
//
// The test uses the SetupTestContainer function to initialize the test dependencies and ensures proper cleanup of resources.
func TestVacancyService_DeleteVacancy(t *testing.T) {
	client, jwtService := SetupTestContainer(t)

	// Generate a valid token
	claims := entity.GetTokenClaims().
		SetIssuer("test-issuer").
		SetScope([]string{"test"}).
		SetExpiresAt(time.Now().Add(time.Hour).Unix())
	validToken, err := jwtService.Generate(claims)
	require.NoError(t, err, "could not generate token")

	// Add the token to the metadata
	md := map[string]string{}
	md["authorization"] = fmt.Sprintf("Bearer %s", validToken)
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(md))

	// Add a test vacancy to delete
	createResp, err := client.CreateVacancy(ctx, &vacancyv1.CreateVacancyRequest{
		Title:       "Software Engineer",
		Company:     "Tech Co.",
		Description: "Exciting opportunity in tech.",
		PostedAt:    "2025-01-01",
		Location:    "New York",
	})
	require.NoError(t, err, "could not create vacancy")

	// Define test cases
	tests := []struct {
		name        string
		token       string
		request     *vacancyv1.DeleteVacancyRequest
		expectedErr bool
		errCode     codes.Code
	}{
		{
			name:  "Valid Deletion",
			token: validToken,
			request: &vacancyv1.DeleteVacancyRequest{
				Id: createResp.Id,
			},
			expectedErr: false,
		},
		{
			name:  "Invalid ID",
			token: validToken,
			request: &vacancyv1.DeleteVacancyRequest{
				Id: -1,
			},
			expectedErr: true,
			errCode:     codes.InvalidArgument,
		},
		{
			name:  "Missing Authorization Header",
			token: "",
			request: &vacancyv1.DeleteVacancyRequest{
				Id: createResp.Id,
			},
			expectedErr: true,
			errCode:     codes.Unauthenticated,
		},
		{
			name:  "Vacancy not found",
			token: validToken,
			request: &vacancyv1.DeleteVacancyRequest{
				Id: createResp.Id + 1_000, // Non-existent ID
			},
			expectedErr: true,
			errCode:     codes.Internal,
		},
	}

	// Execute test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx = metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+tc.token))

			// Make the gRPC call
			resp, err := client.DeleteVacancy(ctx, tc.request)

			if tc.expectedErr {
				// Validate the expected error
				require.Error(t, err, "expected an error but got none")
				st, ok := status.FromError(err)
				require.True(t, ok, "error is not a gRPC status")
				assert.Equal(t, tc.errCode, st.Code(), "unexpected gRPC status code")
			} else {
				require.NoError(t, err, "unexpected error")
				assert.Equal(t, fmt.Sprintf("Vacancy with ID %d deleted", tc.request.Id), resp.Message)
			}
		})
	}
}
