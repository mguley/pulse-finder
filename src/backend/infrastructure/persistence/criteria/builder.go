package criteria

import "sync"

// criteriaPoolInstance is the instance of the getCriteriaPool function to access the pool.
var criteriaPoolInstance = getCriteriaPool()

// getCriteriaPool returns a singleton instance of sync.Pool used to manage SearchCriteriaBuilder objects.
func getCriteriaPool() func() *sync.Pool {
	var once sync.Once
	var pool *sync.Pool

	return func() *sync.Pool {
		once.Do(func() {
			pool = &sync.Pool{
				New: func() interface{} {
					return &SearchCriteriaBuilder{}
				},
			}
		})
		return pool
	}
}

// SearchCriteriaBuilder provides methods to construct SearchCriteria.
type SearchCriteriaBuilder struct {
	criteria SearchCriteria
}

// GetSearchCriteriaBuilder retrieves a SearchCriteriaBuilder object from the pool, resetting it before use.
// If no SearchCriteriaBuilder is available in the pool, a new one is created.
func GetSearchCriteriaBuilder() *SearchCriteriaBuilder {
	return criteriaPoolInstance().Get().(*SearchCriteriaBuilder).Reset()
}

// AddFilter adds a new filter to the SearchCriteria.
func (b *SearchCriteriaBuilder) AddFilter(field, operator string, value any) *SearchCriteriaBuilder {
	b.criteria.Filters = append(b.criteria.Filters, Filter{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
	return b
}

// SetLogicalOperator sets the logical operator (AND/OR) for the criteria.
func (b *SearchCriteriaBuilder) SetLogicalOperator(operator string) *SearchCriteriaBuilder {
	b.criteria.LogicalOperator = operator
	return b
}

// Build finalizes and returns the constructed SearchCriteria.
func (b *SearchCriteriaBuilder) Build() SearchCriteria {
	return b.criteria
}

// Reset resets the fields of the SearchCriteriaBuilder.
func (b *SearchCriteriaBuilder) Reset() *SearchCriteriaBuilder {
	b.criteria.Filters = b.criteria.Filters[:0]
	b.criteria.LogicalOperator = ""
	return b
}

// Release releases the SearchCriteriaBuilder instance back to the pool after resetting it.
func (b *SearchCriteriaBuilder) Release() {
	criteriaPoolInstance().Put(b.Reset())
}
