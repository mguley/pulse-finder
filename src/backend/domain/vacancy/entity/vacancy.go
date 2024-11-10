package entity

import (
	"sync"
	"time"
)

// vacancyInstance is the instance of getVacancyPool function to access the pool.
var vacancyInstance = getVacancyPool()

// getVacancyPool returns a singleton instance of sync.Pool used to manage Vacancy entities.
// It ensures efficient memory use by reusing Vacancy instances.
func getVacancyPool() func() *sync.Pool {
	var once sync.Once
	var pool *sync.Pool

	return func() *sync.Pool {
		once.Do(func() {
			pool = &sync.Pool{
				New: func() interface{} {
					return &Vacancy{}
				},
			}
		})
		return pool
	}
}

// Vacancy represents a job vacancy with relevant details.
type Vacancy struct {
	id          int64     // Unique identifier for the vacancy.
	title       string    // The job title.
	company     string    // The company offering the job.
	description string    // A description of the job.
	postedAt    time.Time // The timestamp when the job was posted.
	location    string    // The location of the job.
	version     int32     // The version number of the job entry, useful for optimistic concurrency control.
}

// Reset resets the fields of the Vacancy to their zero values and returns the updated Vacancy.
func (v *Vacancy) Reset() *Vacancy {
	v.id = 0
	v.title = ""
	v.company = ""
	v.description = ""
	v.postedAt = time.Time{}
	v.location = ""
	v.version = 0
	return v
}

// Release releases the Vacancy instance back to the pool after resetting it.
func (v *Vacancy) Release() {
	vacancyInstance().Put(v.Reset())
}

// GetVacancy retrieves a new or recycled Vacancy instance from the pool.
// It resets the fields to zero values before returning to ensure a clean instance.
func GetVacancy() *Vacancy {
	return vacancyInstance().Get().(*Vacancy).Reset()
}

// GetId returns the unique identifier for the vacancy.
func (v *Vacancy) GetId() int64 {
	return v.id
}

// SetId sets the unique identifier for the vacancy.
func (v *Vacancy) SetId(id int64) *Vacancy {
	v.id = id
	return v
}

// GetTitle returns the job title.
func (v *Vacancy) GetTitle() string {
	return v.title
}

// SetTitle sets the job title.
func (v *Vacancy) SetTitle(title string) *Vacancy {
	v.title = title
	return v
}

// GetCompany returns the company offering the job.
func (v *Vacancy) GetCompany() string {
	return v.company
}

// SetCompany sets the company offering the job.
func (v *Vacancy) SetCompany(company string) *Vacancy {
	v.company = company
	return v
}

// GetDescription returns the description of the job.
func (v *Vacancy) GetDescription() string {
	return v.description
}

// SetDescription sets the description of the job.
func (v *Vacancy) SetDescription(description string) *Vacancy {
	v.description = description
	return v
}

// GetPostedAt returns the timestamp when the job was posted.
func (v *Vacancy) GetPostedAt() time.Time {
	return v.postedAt
}

// SetPostedAt sets the timestamp when the job was posted.
func (v *Vacancy) SetPostedAt(postedAt time.Time) *Vacancy {
	v.postedAt = postedAt
	return v
}

// GetLocation returns the location of the job.
func (v *Vacancy) GetLocation() string {
	return v.location
}

// SetLocation sets the location of the job.
func (v *Vacancy) SetLocation(location string) *Vacancy {
	v.location = location
	return v
}

// GetVersion returns the version number of the job entry.
func (v *Vacancy) GetVersion() int32 {
	return v.version
}

// SetVersion sets the version number of the job entry.
func (v *Vacancy) SetVersion(version int32) *Vacancy {
	v.version = version
	return v
}
