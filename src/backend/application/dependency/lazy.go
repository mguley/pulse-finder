package dependency

import "sync"

// LazyDependency encapsulates lazy initialization logic for any type T.
// It initializes dependencies only upon first access.
type LazyDependency[T any] struct {
	once     sync.Once // Ensures initialization only happens once
	value    T         // Holds the lazily initialized value
	InitFunc func() T  // Initialization function for the dependency
}

// Get initializes the dependency on the first call and returns it thereafter.
func (l *LazyDependency[T]) Get() T {
	l.once.Do(func() {
		l.value = l.InitFunc()
	})
	return l.value
}
