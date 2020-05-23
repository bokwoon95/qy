package qx

import (
	"strings"
)

// JoinTable represents an SQL join.
type JoinTable struct {
	JoinType     JoinType
	Table        Table
	OnPredicates VariadicPredicate
}

// Join constructs a new JoinTable. Meant to be used if you want to do a custom
// join, like NATURAL JOIN.
func Join(joinType string, table Table, predicates ...Predicate) JoinTable {
	return JoinTable{
		JoinType: JoinType(joinType),
		Table:    table,
		OnPredicates: VariadicPredicate{
			Toplevel:   true,
			Predicates: predicates,
		},
	}
}

// JoinTables is a list of JoinTables.
type JoinTables []JoinTable

// WriteSQL will write the JOIN clause into the buffer and args. If there are
// no JoinTables it simply writes nothing into the buffer. It returns a flag
// indicating whether anything was written into the buffer.
func (joins JoinTables) WriteSQL(buf *strings.Builder, args *[]interface{}) (written bool) {
	for i := range joins {
		if joins[i].Table == nil {
			continue
		}
		tableQuery, tableArgs := joins[i].Table.ToSQL()
		if tableQuery == "" {
			continue
		}
		if _, ok := joins[i].Table.(Query); ok {
			tableQuery = "(" + tableQuery + ")"
		}
		if joins[i].JoinType == "" {
			joins[i].JoinType = JoinTypeDefault
		}
		if buf.Len() > 0 {
			buf.WriteString(" ")
		}
		if joins[i].Table.GetAlias() != "" {
			buf.WriteString(string(joins[i].JoinType) + " " + tableQuery + " AS " + joins[i].Table.GetAlias())
		} else {
			buf.WriteString(string(joins[i].JoinType) + " " + tableQuery)
		}
		*args = append(*args, tableArgs...)
		written = true
		joins[i].OnPredicates.Toplevel = true
		joins[i].OnPredicates.WriteSQL(buf, args, "ON ", "", nil)
	}
	return written
}
