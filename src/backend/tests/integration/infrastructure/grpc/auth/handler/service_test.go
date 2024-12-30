package handler

import (
	"context"
	authv1 "infrastructure/proto/auth/gen"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestAuthServiceServer_GenerateToken tests the GenerateToken method of the AuthServiceServer.
//
// This test covers the following scenarios:
// 1. A valid request with proper inputs returns a valid token.
// 2. A request with a missing issuer (Issuer) returns an InvalidArgument error.
// 3. A request with an empty scope returns an InvalidArgument error.
//
// It uses the SetupTestContainer to initialize dependencies and ensures proper cleanup of resources after the test.
func TestAuthServiceServer_GenerateToken(t *testing.T) {
	client := SetupTestContainer(t)

	// Define test cases
	tests := []struct {
		name        string
		request     *authv1.GenerateTokenRequest
		expectedErr bool
		errCode     codes.Code
	}{
		{
			name: "Valid Request",
			request: &authv1.GenerateTokenRequest{
				Issuer: "test-issuer",
				Scopes: []string{"read", "write"},
			},
			expectedErr: false,
		},
		{
			name: "Missing Issuer",
			request: &authv1.GenerateTokenRequest{
				Scopes: []string{"read", "write"},
			},
			expectedErr: true,
			errCode:     codes.InvalidArgument,
		},
		{
			name: "Empty Scope",
			request: &authv1.GenerateTokenRequest{
				Issuer: "",
			},
			expectedErr: true,
			errCode:     codes.InvalidArgument,
		},
	}

	// Execute the test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			// Make the gRPC call
			resp, err := client.GenerateToken(ctx, tc.request)

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
				assert.NotEmpty(t, resp.Token, "token should not be empty")
			}
		})
	}
}
