package interceptors

import (
	"application/auth"
	"context"
	"domain/auth/entity"
	"errors"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// contextKey is a custom type to avoid collisions in context keys.
type contextKey string

const (
	userClaimsKey       contextKey = "userClaims"
	authorizationHeader string     = "authorization"
	bearerPrefix        string     = "Bearer "
)

// JwtVacancyInterceptor validates JWT tokens and adds claims to the context.
func JwtVacancyInterceptor(jwtService *auth.Service) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Extract and validate the token.
		claims, err := extractAndValidateToken(ctx, jwtService)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: %v", err)
		}

		// Add the claims to the context for downstream handlers.
		ctxWithClaims := context.WithValue(ctx, userClaimsKey, claims)
		return handler(ctxWithClaims, req)
	}
}

// extractAndValidateToken extracts the JWT token from incoming metadata and validates it.
func extractAndValidateToken(ctx context.Context, jwtService *auth.Service) (*entity.TokenClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata in request")
	}

	token, err := extractBearerToken(md)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token extraction: %v", err)
	}

	claims, err := jwtService.Verify(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}
	return claims, nil
}

// extractBearerToken retrieves and validates the Bearer token from the Authorization header.
func extractBearerToken(md metadata.MD) (string, error) {
	authHeaderValues := md.Get(authorizationHeader)
	if len(authHeaderValues) == 0 {
		return "", errors.New("authorization header is not provided")
	}

	// We only check the first value (common case).
	tokenString := authHeaderValues[0]
	if !strings.HasPrefix(tokenString, bearerPrefix) {
		return "", errors.New("invalid token format, expected 'Bearer <token>'")
	}

	// Trim out the Bearer part
	return strings.TrimPrefix(tokenString, bearerPrefix), nil
}

// ClaimsFromContext retrieves the claims from the request context.
func ClaimsFromContext(ctx context.Context) (*entity.TokenClaims, bool) {
	claims, ok := ctx.Value(userClaimsKey).(*entity.TokenClaims)
	return claims, ok
}
