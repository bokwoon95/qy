package qx

// CustomTable is a Table that can render itself in an arbitrary way as defined
// by its Format string. Values are interpolated into the Format string as
// described in the (CustomTable).CustomSprintf function.
type CustomTable struct {
	Alias  string
	Format string
	Values []interface{}
}

// ToSQL marshals a CustomTable into an SQL query.
func (tbl CustomTable) ToSQL() (string, []interface{}) {
	query, args := FormatPreprocessor(tbl.Format, tbl.Values, nil)
	return query, args
}

// As returns a new CustomTable with the new alias i.e. 'field AS Alias'.
func (tbl CustomTable) As(alias string) CustomTable {
	tbl.Alias = alias
	return tbl
}

// GetAlias implements the Table interface. It returns the alias of the
// CustomTable.
func (tbl CustomTable) GetAlias() string {
	return tbl.Alias
}

// GetName implements the Table interface. It returns the name of the
// CustomTable.
func (tbl CustomTable) GetName() string {
	name, _ := tbl.ToSQL()
	return name
}

// CustomQuery is a Query that can render itself in an arbitrary way as defined
// by its Format string. Values are interpolated into the Format string as
// described in the (CustomQuery).CustomSprintf function.
//
// The difference between CustomTable and CustomQuery is that CustomTable is
// not meant for writing full queries, because it does not do any form of
// placeholder ?, ?, ? -> $1, $2, $3 etc rebinding.
type CustomQuery struct {
	// Postgres flag determines whether we need to rebind ?, ?, ? to $1, $2,
	// $3.
	Postgres bool

	Nested bool
	Alias  string
	Format string
	Values []interface{}
}

// ToSQL marshals a CustomQuery into an SQL query.
func (q CustomQuery) ToSQL() (string, []interface{}) {
	query, args := FormatPreprocessor(q.Format, q.Values, nil)
	if !q.Nested && q.Postgres {
		query = MySQLToPostgresPlaceholders(query)
	}
	return query, args
}

// As returns a new CustomQuery with the new alias i.e. 'field AS Alias'.
func (q CustomQuery) As(alias string) CustomQuery {
	q.Alias = alias
	return q
}

// GetAlias implements the Table interface. It returns the alias of the
// CustomQuery.
func (q CustomQuery) GetAlias() string {
	return q.Alias
}

// GetName implements the Table interface. It returns the name of the
// CustomQuery.
func (q CustomQuery) GetName() string {
	name, _ := q.ToSQL()
	return name
}

// NestThis implements the Query interfaces.
func (q CustomQuery) NestThis() Query {
	q.Nested = true
	return q
}
