package validators

import "sync"

// validatorPoolInstance is the instance of the getValidatorPool function to access the pool.
var validatorPoolInstance = getValidatorPool()

// getValidatorPool returns a singleton instance of sync.Pool used to manage Validator objects.
func getValidatorPool() func() *sync.Pool {
	var once sync.Once
	var pool *sync.Pool

	return func() *sync.Pool {
		once.Do(func() {
			pool = &sync.Pool{
				New: func() interface{} {
					return &Validator{Errors: make(map[string]string)}
				},
			}
		})
		return pool
	}
}

// Validator accumulates validation errors in a map, associating each error with a field key.
type Validator struct {
	Errors map[string]string // Errors holds validation error messages associated with field keys.
}

// ClearErrors resets the Errors map, making it reusable without residual errors.
func (v *Validator) ClearErrors() *Validator {
	for key := range v.Errors {
		delete(v.Errors, key)
	}
	return v
}

// Valid checks if the Errors map contains no validation errors.
// Returns true if there are no errors; otherwise, false.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError stores an error message associated with a specific field key if it does not already exist.
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check performs a validation check and adds an error message to Errors if the check fails.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Release resets the Validator instance and releases it back to the pool for reuse.
func (v *Validator) Release() {
	validatorPoolInstance().Put(v.ClearErrors())
}

// GetValidator retrieves a Validator object from the pool, resetting it before use.
// If no Validator is available in the pool, a new one is created.
func GetValidator() *Validator {
	return validatorPoolInstance().Get().(*Validator).ClearErrors()
}
