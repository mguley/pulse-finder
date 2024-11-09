package handlers

import (
	"application/auth"
	"domain/auth/entity"
	"interfaces/api/auth/dto"
	"interfaces/api/utils"
	"net/http"
	"time"
)

// JwtTokenHandler handles HTTP requests for JWT token operations.
type JwtTokenHandler struct {
	*utils.Handler // HTTP handler utility
	*utils.Errors  // Error handling utility
	*auth.Service  // Jwt auth service
}

// NewJwtTokenHandler creates a new JwtTokenHandler instance.
func NewJwtTokenHandler(
	handler *utils.Handler,
	errors *utils.Errors,
	service *auth.Service,
) *JwtTokenHandler {
	return &JwtTokenHandler{
		Handler: handler,
		Errors:  errors,
		Service: service,
	}
}

// Execute processes a request to issue a JWT token.
func (h *JwtTokenHandler) Execute(w http.ResponseWriter, r *http.Request) {
	claims := entity.GetTokenClaims()
	defer claims.Release()

	// Set up token claims
	claims.SetIssuer("api.pulse-finder")
	claims.SetScope([]string{"read"})
	claims.SetExpiresAt(time.Now().Add(24 * time.Hour).Unix())

	token, err := h.Service.Generate(claims)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	response := dto.GetResponse()
	defer response.Release()
	response.FromToken(token)
	if err = h.WriteJson(w, http.StatusOK, response, nil); err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}
