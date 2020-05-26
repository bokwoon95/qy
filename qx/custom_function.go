package qx

type CustomFunction struct {
	Alias  string
	Format string
	Values []interface{}

	// Each dialect-specific qy package (postgres, mysql, sqlite3) is expected
	// to provide their dialect-specific CustomSprintf function to CustomFunction.
	// If none is provided, it will fall back on using the the defaultSprintf
	// function in this package.
	CustomSprintf func(format string, values []interface{}, excludeTableQualifiers []string) (string, []interface{})
}

// ToSQL marshals a CustomFunction into an SQL query.
func (tbl CustomFunction) ToSQL() (string, []interface{}) {
	var query string
	var args []interface{}
	if tbl.CustomSprintf != nil {
		query, args = tbl.CustomSprintf(tbl.Format, tbl.Values, nil)
	} else {
		query, args = defaultSprintf(tbl.Format, tbl.Values, nil)
	}
	return query, args
}

// As returns a new CustomFunction with the new alias i.e. 'field AS Alias'.
func (tbl CustomFunction) As(alias string) CustomFunction {
	tbl.Alias = alias
	return tbl
}

// GetAlias implements the Table interface. It returns the alias of the
// CustomFunction.
func (tbl CustomFunction) GetAlias() string {
	return tbl.Alias
}

// GetName implements the Table interface. It returns the name of the
// CustomFunction.
func (tbl CustomFunction) GetName() string {
	name, _ := tbl.ToSQL()
	return name
}
