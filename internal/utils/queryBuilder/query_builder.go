package queryBuilder

import (
	"fmt"
	"strings"

	"github.com/Alfian57/belajar-golang/internal/dto"
	"github.com/Alfian57/belajar-golang/internal/logger"
)

type QueryBuilder struct {
	baseQuery     string
	whereClause   []string
	args          []any
	orderBy       string
	orderType     string
	hasPagination bool
}

func NewQueryBuilder(baseQuery string) *QueryBuilder {
	return &QueryBuilder{
		baseQuery: baseQuery,
		args:      make([]any, 0),
	}
}

func (qb *QueryBuilder) Where(condition string, args ...any) *QueryBuilder {
	qb.whereClause = append(qb.whereClause, condition)
	qb.args = append(qb.args, args...)

	logger.Log.Infow("QueryBuilder Where condition", "args", args)
	return qb
}

func (qb *QueryBuilder) Search(column, searchTerm string) *QueryBuilder {
	if searchTerm != "" {
		paramNumber := len(qb.args) + 1
		qb.whereClause = append(qb.whereClause, fmt.Sprintf("%s LIKE $%d", column, paramNumber))
		qb.args = append(qb.args, "%"+searchTerm+"%")
	}
	return qb
}

func (qb *QueryBuilder) OrderBy(column, orderType string) *QueryBuilder {
	if column == "" {
		column = "created_at" // default order by
	}
	qb.orderBy = column
	qb.orderType = strings.ToUpper(orderType)
	if qb.orderType != "ASC" && qb.orderType != "DESC" {
		qb.orderType = "ASC"
	}
	return qb
}

func (qb *QueryBuilder) Paginate(pagination dto.PaginationRequest) *QueryBuilder {
	pagination.SetDefaults()

	qb.args = append(qb.args, pagination.Limit, pagination.GetOffset())
	qb.hasPagination = true
	return qb
}

// Build constructs the final SQL query string based on the base query, where clauses,
// order by clause, and pagination settings. It returns the complete query string and
// a slice of arguments to be used with parameterized queries. The method appends
// WHERE conditions if present, applies ORDER BY if specified, and adds LIMIT/OFFSET
// for pagination when enabled.
func (qb *QueryBuilder) Build() (string, []any) {
	query := qb.baseQuery

	if len(qb.whereClause) > 0 {
		query += " WHERE " + strings.Join(qb.whereClause, " AND ")
	}

	if qb.orderBy != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", qb.orderBy, qb.orderType)
	}

	if qb.hasPagination {
		limitParam := len(qb.args) - 1 // second to last arg is limit
		offsetParam := len(qb.args)    // last arg is offset
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", limitParam, offsetParam)
	}

	return query, qb.args
}

// BuildCount constructs a SQL count query by appending any WHERE clauses stored in the QueryBuilder.
// It returns the final query string and a slice of arguments to be used with the query.
func (qb *QueryBuilder) BuildCount(countQuery string) (string, []any) {
	query := countQuery

	if len(qb.whereClause) > 0 {
		query += " WHERE " + strings.Join(qb.whereClause, " AND ")
	}

	return query, qb.args
}
