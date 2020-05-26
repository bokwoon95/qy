package qx

type FunctionInfo struct {
	Alias     string
	Query     string
	Arguments []interface{}

	// Each dialect-specific qy package (postgres, mysql, sqlite3) is expected
	// to provide their dialect-specific CustomSprintf function to FunctionInfo.
	// If none is provided, it will fall back on using the the defaultSprintf
	// function in this package.
	CustomSprintf func(format string, values []interface{}, excludeTableQualifiers []string) (string, []interface{})
}

// ToSQL marshals a FunctionInfo into an SQL query.
func (tbl FunctionInfo) ToSQL() (string, []interface{}) {
	return tbl.ToSQLExclude(nil)
}

// ToSQL marshals a FunctionInfo into an SQL query.
func (tbl FunctionInfo) ToSQLExclude(excludeTableQualifiers []string) (string, []interface{}) {
	var query string
	var args []interface{}
	if tbl.CustomSprintf != nil {
		query, args = tbl.CustomSprintf(tbl.Query, tbl.Arguments, excludeTableQualifiers)
	} else {
		query, args = defaultSprintf(tbl.Query, tbl.Arguments, excludeTableQualifiers)
	}
	return query, args
}

// As returns a new FunctionInfo with the new alias i.e. 'field AS Alias'.
func (tbl FunctionInfo) As(alias string) FunctionInfo {
	tbl.Alias = alias
	return tbl
}

// GetAlias implements the Table interface. It returns the alias of the
// FunctionInfo.
func (tbl FunctionInfo) GetAlias() string {
	return tbl.Alias
}

// GetName implements the Table interface. It returns the name of the
// FunctionInfo.
func (tbl FunctionInfo) GetName() string {
	name, _ := tbl.ToSQL()
	return name
}
