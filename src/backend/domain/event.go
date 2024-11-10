package domain

// Event represents a domain event.
// Each event has a unique type and is associated with a specific aggregate (entity) ID.
type Event interface {
	// EventType returns a string representing the type of the event.
	EventType() string

	// AggregateId returns the ID of the entity (aggregate) that the event is related to.
	AggregateId() int64
}
