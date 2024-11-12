package dto

import (
	"domain/vacancy/entity"
	"sync"
	"time"
)

// responsePoolInstance is the instance of the getResponsePool function to access the pool.
var responsePoolInstance = getResponsePool()

// getResponsePool returns a singleton instance of sync.Pool used to manage Response objects.
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

// Response represents the data transfer object for a job vacancy response.
type Response struct {
	ID          *int64  `json:"id,omitempty"`          // ID is the unique identifier of the job vacancy.
	Title       *string `json:"title,omitempty"`       // Title of the job vacancy.
	Company     *string `json:"company,omitempty"`     // Company offering the job.
	Description *string `json:"description,omitempty"` // Description of the job vacancy.
	PostedAt    *string `json:"posted_at,omitempty"`   // PostedAt is the timestamp when the job was posted.
	Location    *string `json:"location,omitempty"`    // Location of the job.
}

// Reset resets the fields of the Response to their zero values and returns the updated Response.
func (r *Response) Reset() *Response {
	r.ID = nil
	r.Title = nil
	r.Company = nil
	r.Description = nil
	r.PostedAt = nil
	r.Location = nil
	return r
}

// Release releases the Response instance back to the pool after resetting it.
func (r *Response) Release() {
	responsePoolInstance().Put(r.Reset())
}

// GetResponse retrieves a Response object from the pool, resetting it before use.
// If no Response is available in the pool, a new one is created.
func GetResponse() *Response {
	return responsePoolInstance().Get().(*Response).Reset()
}

// FromEntity maps the Vacancy entity fields to the Response fields.
func (r *Response) FromEntity(e *entity.Vacancy) *Response {
	id, title, company, description, postedAt, location := e.GetId(), e.GetTitle(), e.GetCompany(), e.GetDescription(),
		e.GetPostedAt(), e.GetLocation()

	r.ID = &id
	r.Title = &title
	r.Company = &company
	r.Description = &description
	if !postedAt.IsZero() {
		v := postedAt.Format(time.DateOnly)
		r.PostedAt = &v
	}
	r.Location = &location
	return r
}

// ToList converts a slice of Vacancy entities to a slice of Response objects.
func (r *Response) ToList(list []*entity.Vacancy) *[]Response {
	items := make([]Response, len(list))
	for i, item := range list {
		items[i] = *r.Reset().FromEntity(item)
	}
	return &items
}
