package events

// VacancyDeletedEvent represents the event that occurs when an existing vacancy is deleted from the system.
type VacancyDeletedEvent struct {
	VacancyId int64 // Unique identifier of the deleted vacancy.
	// todo add payload
}

// NewVacancyDeletedEvent initializes a new VacancyDeletedEvent.
func NewVacancyDeletedEvent(id int64) *VacancyDeletedEvent {
	return &VacancyDeletedEvent{VacancyId: id}
}

// EventType returns the type of the event.
func (e *VacancyDeletedEvent) EventType() string {
	return string(VacancyDeleted)
}

// AggregateId returns the unique identifier of the vacancy associated with this event.
func (e *VacancyDeletedEvent) AggregateId() int64 {
	return e.VacancyId
}
