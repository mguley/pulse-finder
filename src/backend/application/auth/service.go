package auth

import (
	"application/config"
	"domain/auth/entity"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Service provides application services for managing JWT tokens.
type Service struct {
	config    *config.Configuration
	secretKey []byte
}

// NewService initializes a new Service.
func NewService(c *config.Configuration) *Service {
	return &Service{config: c, secretKey: []byte(c.Jwt.Secret)}
}

// Generate creates a JWT token for the provided entity.TokenClaims.
func (s *Service) Generate(claims *entity.TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":   claims.GetIssuer(),
		"scope": claims.GetScope(),
		"exp":   claims.GetExpiresAt(),
	})

	// Sign the token using the secret key
	signedToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// Verify checks the provided token string, ensuring it is valid.
func (s *Service) Verify(t string) (*entity.TokenClaims, error) {
	token, err := jwt.Parse(t, s.keyFunc)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims structure")
	}

	return s.extractClaims(claims)
}

// keyFunc retrieves the signing key and validates the signing method.
func (s *Service) keyFunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return s.secretKey, nil
}

// extractClaims extracts and validates each claim.
func (s *Service) extractClaims(claims jwt.MapClaims) (*entity.TokenClaims, error) {
	issuer, err := s.getIssuer(claims)
	if err != nil {
		return nil, err
	}

	scope, err := s.getScope(claims)
	if err != nil {
		return nil, err
	}

	exp, err := s.getExpiration(claims)
	if err != nil {
		return nil, err
	}

	return entity.GetTokenClaims().SetIssuer(issuer).SetScope(scope).SetExpiresAt(exp), err
}

// getIssuer retrieves and validates the issuer claim.
func (s *Service) getIssuer(claims jwt.MapClaims) (string, error) {
	issuer, ok := claims["iss"].(string)
	if !ok || issuer == "" {
		return "", fmt.Errorf("missing or invalid issuer claim")
	}
	return issuer, nil
}

// getScope retrieves and validates the scope claim.
func (s *Service) getScope(claims jwt.MapClaims) ([]string, error) {
	rawScope, ok := claims["scope"].([]any)
	if !ok {
		return nil, fmt.Errorf("missing or invalid scope claim")
	}

	var scope []string
	for _, v := range rawScope {
		if s, ok := v.(string); ok {
			scope = append(scope, s)
		} else {
			return nil, fmt.Errorf("invalid scope value type")
		}
	}
	return scope, nil
}

// getExpiration retrieves and validates the expiration claim.
func (s *Service) getExpiration(claims jwt.MapClaims) (int64, error) {
	exp, ok := claims["exp"].(float64)
	if !ok {
		return 0, fmt.Errorf("missing or invalid expiration claim")
	}
	expiration := int64(exp)
	if expiration < time.Now().Unix() {
		return 0, fmt.Errorf("token expired")
	}
	return expiration, nil
}
