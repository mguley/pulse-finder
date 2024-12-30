package handler

import (
	"application/auth"
	"context"
	"domain/auth/entity"
	authv1 "infrastructure/proto/auth/gen"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service implements the gRPC AuthServiceServer.
// It handles requests related to JWT token generation.
type Service struct {
	authv1.UnimplementedAuthServiceServer               // Ensures forward compatibility with the gRPC interface.
	authService                           *auth.Service // Dependency for handling JWT token operations.
}

// NewService creates a new instance of the gRPC Service handler.
func NewService(authService *auth.Service) *Service {
	return &Service{authService: authService}
}

// GenerateToken handles gRPC requests to generate a new JWT token.
func (s *Service) GenerateToken(
	ctx context.Context,
	req *authv1.GenerateTokenRequest,
) (*authv1.GenerateTokenResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "context canceled")
	default:
		return s.process(req)
	}
}

// process performs validation and token generation.
func (s *Service) process(req *authv1.GenerateTokenRequest) (*authv1.GenerateTokenResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	// Generate token claims.
	claims := entity.GetTokenClaims()
	defer claims.Release()

	claims.SetIssuer(req.GetIssuer())
	claims.SetScope(req.GetScopes())
	claims.SetExpiresAt(time.Now().Add(5 * time.Minute).Unix())

	// Generate the token.
	token, err := s.authService.Generate(claims)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}
	return &authv1.GenerateTokenResponse{Token: token}, nil
}

// validateRequest ensures the request has all necessary fields.
func (s *Service) validateRequest(req *authv1.GenerateTokenRequest) error {
	if req.GetIssuer() == "" {
		return status.Errorf(codes.InvalidArgument, "issuer (iss) must not be empty")
	}
	if len(req.GetScopes()) == 0 {
		return status.Errorf(codes.InvalidArgument, "scope must not be empty")
	}
	return nil
}
