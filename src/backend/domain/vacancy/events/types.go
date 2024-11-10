package events

// EventType represents the type of the event occurring in the system.
type EventType string

const (
	// VacancyCreated indicates that a new vacancy has been created.
	VacancyCreated EventType = "VacancyCreated"
	// VacancyUpdated indicates that an existing vacancy has been updated.
	VacancyUpdated EventType = "VacancyUpdated"
	// VacancyDeleted indicates that an existing vacancy has been deleted.
	VacancyDeleted EventType = "VacancyDeleted"
)
