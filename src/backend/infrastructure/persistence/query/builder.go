package query

import (
	"fmt"
	"infrastructure/persistence/criteria"
	"strings"
	"sync"
)

// queryPoolInstance is the instance of the getQueryPool function to access the pool.
var queryPoolInstance = getQueryPool()

// getQueryPool returns a singleton instance of sync.Pool used to manage Builder objects.
func getQueryPool() func() *sync.Pool {
	var once sync.Once
	var pool *sync.Pool

	return func() *sync.Pool {
		once.Do(func() {
			pool = &sync.Pool{
				New: func() interface{} {
					return &Builder{}
				},
			}
		})
		return pool
	}
}

// Builder assists in building dynamic SQL queries based on criteria.SearchCriteria.
type Builder struct {
	baseQuery  string   // The base SQL query (e.g., "SELECT * FROM table_name")
	conditions []string // SQL conditions to apply
	args       []any    // Arguments for SQL placeholders
	orderBy    string   // Sorting clause
	limit      int      // Pagination limit
	offset     int      // Pagination offset
}

// GetBuilder retrieves a Builder object from the pool, resetting it before use.
// If no Builder is available in the pool, a new one is created.
func GetBuilder(baseQuery string) *Builder {
	b := queryPoolInstance().Get().(*Builder).Reset()
	b.baseQuery = baseQuery
	return b
}

// ApplySearchCriteria adds filters from criteria.SearchCriteria to Builder.
func (b *Builder) ApplySearchCriteria(c criteria.SearchCriteria) {
	for _, filter := range c.Filters {
		placeholder := fmt.Sprintf("$%d", len(b.args)+1)
		condition := fmt.Sprintf("%s %s %s", filter.Field, filter.Operator, placeholder)
		b.conditions = append(b.conditions, condition)
		b.args = append(b.args, filter.Value)
	}
}

// SetOrderBy specifies the ORDER BY clause.
func (b *Builder) SetOrderBy(field, order string) {
	b.orderBy = fmt.Sprintf("ORDER BY %s %s", field, order)
}

// SetPagination configures LIMIT and OFFSET for pagination.
func (b *Builder) SetPagination(page, pageSize int) {
	b.limit = pageSize
	b.offset = (page - 1) * pageSize
}

// Build constructs the final query with applied conditions, sorting and pagination.
func (b *Builder) Build(c criteria.SearchCriteria) (query string, args []any) {
	query = b.baseQuery

	// Combine conditions with the specified logical operator (e.g., "AND", "OR")
	if len(b.conditions) > 0 {
		conditions := strings.Join(b.conditions, " "+c.LogicalOperator+" ")
		query += " WHERE " + conditions
	}

	// Add sorting and pagination
	if b.orderBy != "" {
		query += " " + b.orderBy
	}
	if b.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", b.limit)
	}
	if b.offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", b.offset)
	}
	return query, b.args
}

// Reset resets the fields of the Builder.
func (b *Builder) Reset() *Builder {
	b.conditions = b.conditions[:0]
	b.args = b.args[:0]
	b.orderBy = ""
	b.limit = 0
	b.offset = 0
	return b
}

// Release releases Builder instance back to the pool after resetting it.
func (b *Builder) Release() {
	b.baseQuery = ""
	queryPoolInstance().Put(b.Reset())
}
