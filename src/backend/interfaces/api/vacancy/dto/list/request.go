package list

import "sync"

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

// Request represents the data transfer object for listing job vacancies with filtering options.
type Request struct {
	Title   *string `json:"title,omitempty"`   // Title specifies a filter by job vacancy title.
	Company *string `json:"company,omitempty"` // Company specifies a filter by company offering the job.
	Filters Filters `json:"filters"`           // Filters define pagination and sorting options for listing.
}

// Filters defines pagination and sorting options.
type Filters struct {
	Page      *int    `json:"page,omitempty"`       // Page represents the current page number for pagination.
	PageSize  *int    `json:"page_size,omitempty"`  // PageSize represents the number of items per page.
	SortField *string `json:"sort_field,omitempty"` // SortField defines the field and direction for sorting results.
	SortOrder *string `json:"sort_order,omitempty"` // SortOrder defines the field and how to sort data (ASC, DESC).
}

// Reset resets the fields of the Request to their zero values and returns the updated Request.
func (r *Request) Reset() *Request {
	r.Title = nil
	r.Company = nil
	r.Filters.Reset()
	return r
}

// Reset resets the fields of the Filters to their zero values and returns the updated Filters.
func (f *Filters) Reset() *Filters {
	f.Page = nil
	f.PageSize = nil
	f.SortField = nil
	f.SortOrder = nil
	return f
}

// Release resets the Request struct and places it back into the sync.Pool for reuse.
func (r *Request) Release() {
	requestPoolInstance().Put(r.Reset())
}

// GetRequest retrieves a new or recycled Request struct from the sync.Pool and resets it for reuse.
func GetRequest() *Request {
	return requestPoolInstance().Get().(*Request).Reset()
}
