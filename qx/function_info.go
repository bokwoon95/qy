package qx

import "strings"

type FunctionInfo struct {
	Schema    string
	Name      string
	Alias     string
	Arguments []interface{}
}

// ToSQL marshals a FunctionInfo into an SQL query.
func (f *FunctionInfo) ToSQL() (string, []interface{}) {
	return f.ToSQLExclude(nil)
}

// ToSQL marshals a FunctionInfo into an SQL query.
func (f *FunctionInfo) ToSQLExclude(excludeTableQualifiers []string) (string, []interface{}) {
	var query string
	var args []interface{}
	schema := f.Schema + "."
	if f.Schema == "public." {
		schema = ""
	}
	switch len(f.Arguments) {
	case 0:
		query = schema + f.Name + "()"
	default:
		query = schema + f.Name + "(?" + strings.Repeat(", ?", len(f.Arguments)-1) + ")"
	}
	query, args = defaultSprintf(query, f.Arguments, excludeTableQualifiers)
	return query, args
}

// GetAlias implements the Table interface. It returns the alias of the
// FunctionInfo.
func (f *FunctionInfo) GetAlias() string {
	return f.Alias
}

// GetName implements the Table interface. It returns the name of the
// FunctionInfo.
func (f *FunctionInfo) GetName() string {
	return f.Name
}
