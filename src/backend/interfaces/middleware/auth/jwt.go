package auth

import (
	"application/auth"
	"interfaces/api/utils"
	"net/http"
	"strings"
)

// JwtAuthMiddleware handles JWT-based authorization.
type JwtAuthMiddleware struct {
	*auth.Service // Jwt auth service
	*utils.Errors // Error handling utility
}

// NewJwtAuthMiddleware creates a new instance of JwtAuthMiddleware.
func NewJwtAuthMiddleware(
	service *auth.Service,
	errors *utils.Errors,
) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{
		Service: service,
		Errors:  errors,
	}
}

// Handle checks for a valid JWT token and verifies the issuer.
func (m *JwtAuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		h := r.Header.Get("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			m.Errors.Unauthorized(w, r)
			return
		}

		token := strings.TrimPrefix(h, "Bearer ")
		claims, err := m.Service.Verify(token)
		if err != nil {
			m.Errors.Unauthorized(w, r)
			return
		}

		expectedIssuer := "api.pulse-finder"
		if claims.GetIssuer() != expectedIssuer {
			m.Errors.Unauthorized(w, r)
			return
		}

		// Token is valid
		next.ServeHTTP(w, r)
	})
}
