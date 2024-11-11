package dto

import (
	"domain/vacancy/entity"
	"sync"
	"time"
)

// requestPoolInstance is the instance of the getRequestPool function to access the pool.
var requestPoolInstance = getRequestPool()

// getRequestPool returns a singleton instance of sync.Pool used to manage Request objects.
func getRequestPool() func() *sync.Pool {
	var once sync.Once
	var pool *sync.Pool

	return func() *sync.Pool {
		once.Do(func() {
			pool = &sync.Pool{
				New: func() interface{} {
					return &Request{}
				},
			}
		})
		return pool
	}
}

// Request represents the data transfer object for a job vacancy request.
type Request struct {
	ID          *int64  `json:"id,omitempty"`          // ID is the unique identifier of the job vacancy.
	Title       *string `json:"title,omitempty"`       // Title of the job vacancy.
	Company     *string `json:"company,omitempty"`     // Company offering the job.
	Description *string `json:"description,omitempty"` // Description of the job vacancy.
	PostedAt    *string `json:"posted_at,omitempty"`   // PostedAt is the timestamp when the job was posted.
	Location    *string `json:"location,omitempty"`    // Location of the job.
}

// Reset resets the fields of the Request to their zero values and returns the updated Request.
func (r *Request) Reset() *Request {
	r.ID = nil
	r.Title = nil
	r.Company = nil
	r.Description = nil
	r.PostedAt = nil
	r.Location = nil
	return r
}

// Release releases the Request instance back to the pool after resetting it.
func (r *Request) Release() {
	requestPoolInstance().Put(r.Reset())
}

// GetRequest retrieves a Request object from the pool, resetting it before use.
// If no Request is available in the pool, a new one is created.
func GetRequest() *Request {
	return requestPoolInstance().Get().(*Request).Reset()
}

// ToEntity maps the Request fields to the provided Vacancy entity.
func (r *Request) ToEntity(e *entity.Vacancy) {
	if r.Title != nil {
		e.SetTitle(*r.Title)
	}
	if r.Company != nil {
		e.SetCompany(*r.Company)
	}
	if r.Description != nil {
		e.SetDescription(*r.Description)
	}
	if r.PostedAt != nil {
		parsedTime, err := time.Parse(time.RFC3339, *r.PostedAt)
		if err == nil {
			e.SetPostedAt(parsedTime)
		}
	}
	if r.Location != nil {
		e.SetLocation(*r.Location)
	}
}
