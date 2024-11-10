package events

// VacancyUpdatedEvent represents the event that occurs when an existing vacancy is updated in the system.
type VacancyUpdatedEvent struct {
	VacancyId int64 // Unique identifier of the updated vacancy.
	// todo add payload
}

// NewVacancyUpdatedEvent initializes a new VacancyUpdatedEvent.
func NewVacancyUpdatedEvent(id int64) *VacancyUpdatedEvent {
	return &VacancyUpdatedEvent{VacancyId: id}
}

// EventType returns the type of the event.
func (e *VacancyUpdatedEvent) EventType() string {
	return string(VacancyUpdated)
}

// AggregateId returns the unique identifier of the vacancy associated with this event.
func (e *VacancyUpdatedEvent) AggregateId() int64 {
	return e.VacancyId
}
