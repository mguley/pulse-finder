package entity

import "sync"

// tokenClaimsInstance is the instance of getTokenClaimsPool function to access the pool.
var tokenClaimsInstance = getTokenClaimsPool()

// getTokenClaimsPool returns a singleton instance of sync.Pool used to manage TokenClaims entities.
// It ensures efficient memory use by reusing TokenClaims instances.
func getTokenClaimsPool() func() *sync.Pool {
	var once sync.Once
	var pool *sync.Pool

	return func() *sync.Pool {
		once.Do(func() {
			pool = &sync.Pool{
				New: func() interface{} {
					return &TokenClaims{}
				},
			}
		})
		return pool
	}
}

// TokenClaims represents the claims included in a JWT token.
type TokenClaims struct {
	issuer    string   // The issuer of the token, typically the API name or identifier.
	scope     []string // The permissions or scope associated with the token.
	expiresAt int64    // The expiration time of the token as a UNIX timestamp.
}

// Reset resets the fields of the TokenClaims to their zero values and returns the updated TokenClaims.
func (t *TokenClaims) Reset() *TokenClaims {
	t.issuer = ""
	t.scope = []string{}
	t.expiresAt = 0
	return t
}

// Release releases the TokenClaims instance back to the pool after resetting it.
func (t *TokenClaims) Release() {
	tokenClaimsInstance().Put(t.Reset())
}

// GetTokenClaims retrieves a new or recycled TokenClaims instance from the pool.
// It resets the fields to zero values before returning to ensure a clean instance.
func GetTokenClaims() *TokenClaims {
	return tokenClaimsInstance().Get().(*TokenClaims).Reset()
}

// GetIssuer returns the issuer of the token.
func (t *TokenClaims) GetIssuer() string {
	return t.issuer
}

// SetIssuer sets the issuer of the token.
func (t *TokenClaims) SetIssuer(issuer string) *TokenClaims {
	t.issuer = issuer
	return t
}

// GetScope returns the scope associated with the token.
func (t *TokenClaims) GetScope() []string {
	return t.scope
}

// SetScope sets the scope associated with the token.
func (t *TokenClaims) SetScope(scope []string) *TokenClaims {
	t.scope = scope
	return t
}

// GetExpiresAt returns the expiration time of the token.
func (t *TokenClaims) GetExpiresAt() int64 {
	return t.expiresAt
}

// SetExpiresAt sets the expiration time of the token.
func (t *TokenClaims) SetExpiresAt(expiresAt int64) *TokenClaims {
	t.expiresAt = expiresAt
	return t
}
