package events

// VacancyCreatedEvent represents the event that occurs when a new vacancy is created in the system.
type VacancyCreatedEvent struct {
	VacancyId int64 // Unique identifier of the created vacancy.
	// todo add payload
}

// NewVacancyCreatedEvent initializes a new VacancyCreatedEvent.
func NewVacancyCreatedEvent(id int64) *VacancyCreatedEvent {
	return &VacancyCreatedEvent{VacancyId: id}
}

// EventType returns the type of the event.
func (e *VacancyCreatedEvent) EventType() string {
	return string(VacancyCreated)
}

// AggregateId returns the unique identifier of the vacancy associated with this event.
func (e *VacancyCreatedEvent) AggregateId() int64 {
	return e.VacancyId
}
