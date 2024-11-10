package criteria

// SearchCriteriaBuilder provides methods to construct SearchCriteria.
type SearchCriteriaBuilder struct {
	criteria SearchCriteria
}

// NewSearchCriteriaBuilder initializes a new SearchCriteriaBuilder.
func NewSearchCriteriaBuilder() *SearchCriteriaBuilder {
	return &SearchCriteriaBuilder{
		criteria: SearchCriteria{
			Filters:         []Filter{},
			LogicalOperator: "AND", // Default to "AND" if not specified.
		},
	}
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
