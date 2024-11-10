package criteria

// SearchCriteria struct to hold multiple filters and the logical operator.
type SearchCriteria struct {
	Filters         []Filter // A list of filters to be applied.
	LogicalOperator string   // LogicalOperator specifies how filters are combined (e.g., "AND", "OR").
}
