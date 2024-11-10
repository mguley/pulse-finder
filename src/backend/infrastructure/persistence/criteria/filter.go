package criteria

// Filter struct to represent a single filtering condition.
type Filter struct {
	Field    string // Field represents the column name to filter by.
	Operator string // Operator is the comparison operator (e.g., "=", "LIKE", ">", "<").
	Value    any    // Value is the value to compare against.
}
