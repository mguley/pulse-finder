package event

import "domain"

// Dispatcher defines an interface for dispatching events that occur within the domain.
type Dispatcher interface {
	// Dispatch sends a domain event to the appropriate subscribers.
	// It returns an error if the dispatching fails.
	Dispatch(event domain.Event) error
}
