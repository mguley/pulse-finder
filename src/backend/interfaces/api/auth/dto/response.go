package dto

import "sync"

// responsePoolInstance is the instance of the getResponsePool function to access the pool.
var responsePoolInstance = getResponsePool()

// getResponsePool returns a singleton instance of sync.Pool used to manage Response objects.
// It uses a closure to keep the pool variable hidden from the package scope.
func getResponsePool() func() *sync.Pool {
	var once sync.Once
	var pool *sync.Pool

	return func() *sync.Pool {
		once.Do(func() {
			pool = &sync.Pool{
				New: func() interface{} {
					return &Response{}
				},
			}
		})
		return pool
	}
}

// Response defines the structure of the JSON response for JWT token data.
type Response struct {
	Token string `json:"token,omitempty"` // Token is the completed and signed JWT.
}

// Reset resets the fields of the Response to their zero values and returns the updated Response.
func (r *Response) Reset() *Response {
	r.Token = ""
	return r
}

// Release releases the Response instance back to the pool after resetting it.
func (r *Response) Release() {
	responsePoolInstance().Put(r.Reset())
}

// GetResponse retrieves a new or cycled Response instance from the pool.
// It resets the fields to zero values before returning to ensure a clean instance.
func GetResponse() *Response {
	return responsePoolInstance().Get().(*Response).Reset()
}

// FromToken populates the Response struct with data from the provided token.
func (r *Response) FromToken(t string) *Response {
	r.Token = t
	return r
}
