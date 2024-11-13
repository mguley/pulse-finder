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

// Response defines the structure of the JSON response for health check operation.
// It contains fields for status, system information and a timestamp.
type Response struct {
	Status     *string    `json:"status,omitempty"`      // Status represents the health status of the system.
	SystemInfo SystemInfo `json:"system_info,omitempty"` // SystemInfo contains environment and version of the system.
	Timestamp  *string    `json:"timestamp,omitempty"`   // Timestamp is the time of the health check.
}

// SystemInfo holds information about the environment and version of the system.
type SystemInfo struct {
	Environment *string `json:"environment,omitempty"` // Environment specifies the environment.
	Version     *string `json:"version,omitempty"`     // Version indicates the version of the application.
}

// Reset resets the fields of the Response to their zero values and returns the updated Response.
// This prepares the Response for reuse from the pool by clearing all field values.
func (r *Response) Reset() *Response {
	r.Status = nil
	r.SystemInfo.Environment = nil
	r.SystemInfo.Version = nil
	r.Timestamp = nil
	return r
}

// Release resets the Response struct and places it back into the sync.Pool for reuse.
// This helps manage memory usage by reusing Response objects instead of creating new ones.
func (r *Response) Release() {
	responsePoolInstance().Put(r.Reset())
}

// GetResponse fetches a new or recycled Response struct from the sync.Pool and resets it for reuse.
func GetResponse() *Response {
	return responsePoolInstance().Get().(*Response).Reset()
}
