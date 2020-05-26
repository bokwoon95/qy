package qx

import "strings"

type ArrayField struct {
	// ArrayField will be one of the following:

	// 1) Literal array value (only []bool, []float64, []int64 or []string
	// slices are supported.) Nested slices are also not supported even though
	// both Go and Postgres support nested slices/arrays because I'm not even
	// sure if it's possible to convert between the two with lib/pq.
	// Addtionally []int is supported, but note that it only works when
	// converting from Go slices to Postgres arrays. When converting from
	// postgres arrays to Go slices, you have to use []int64 instead.
	// Examples of literal array values:
	// | query             | args                    |
	// |-------------------|-------------------------|
	// | ARRAY[?, ?, ?, ?] | 1, 2, 3, 4              |
	// | ARRAY[?, ?, ?]    | 22.7, 3.15, 4.0         |
	// | ARRAY[?, ?, ?]    | apple, banana, cucumber |
	value interface{}

	// 2) Array column
	// Examples of boolean columns:
	// | query                 | args |
	// |-----------------------|------|
	// | film.special_features |      |
	// | special_features      |      |
	alias      string
	table      Table
	name       string
	descending *bool
	nullsfirst *bool
}

func (f ArrayField) ToSQLExclude(excludeTableQualifiers []string) (string, []interface{}) {
	// 1) Literal array value
	if f.value != nil {
		var query string
		var args []interface{}
		switch array := f.value.(type) {
		case []bool:
			if len(array) == 0 {
				query, args = "ARRAY[]::BOOLEAN[]", nil
			} else {
				query = "ARRAY[?" + strings.Repeat(", ?", len(array)-1) + "]"
				args = make([]interface{}, len(array))
				for i := range array {
					args[i] = array[i]
				}
			}
		case []float64:
			if len(array) == 0 {
				query, args = "ARRAY[]::FLOAT[]", nil
			} else {
				query = "ARRAY[?" + strings.Repeat(", ?", len(array)-1) + "]"
				args = make([]interface{}, len(array))
				for i := range array {
					args[i] = array[i]
				}
			}
		case []int:
			if len(array) == 0 {
				query, args = "ARRAY[]::INT[]", nil
			} else {
				query = "ARRAY[?" + strings.Repeat(", ?", len(array)-1) + "]"
				args = make([]interface{}, len(array))
				for i := range array {
					args[i] = array[i]
				}
			}
		case []int64:
			if len(array) == 0 {
				query, args = "ARRAY[]::BIGINT[]", nil
			} else {
				query = "ARRAY[?" + strings.Repeat(", ?", len(array)-1) + "]"
				args = make([]interface{}, len(array))
				for i := range array {
					args[i] = array[i]
				}
			}
		case []string:
			if len(array) == 0 {
				query, args = "ARRAY[]::TEXT[]", nil
			} else {
				query = "ARRAY[?" + strings.Repeat(", ?", len(array)-1) + "]"
				args = make([]interface{}, len(array))
				for i := range array {
					args[i] = array[i]
				}
			}
		default:
			query = "(unknown array type: only []bool/[]float64/[]int64/[]string/[]int slices are supported.)"
		}
		return query, args
	}

	// 2) Array column
	var tableQualifier string
	if f.table != nil {
		if f.table.GetAlias() != "" {
			tableQualifier = f.table.GetAlias() + "."
		} else if f.table.GetName() != "" {
			tableQualifier = f.table.GetName() + "."
		}
	}
	for i := range excludeTableQualifiers {
		if tableQualifier == excludeTableQualifiers[i]+"." {
			tableQualifier = ""
			break
		}
	}
	columnName := tableQualifier + f.name
	if f.descending != nil {
		if *f.descending {
			columnName = columnName + " DESC"
		} else {
			columnName = columnName + " ASC"
		}
	}
	if f.nullsfirst != nil {
		if *f.nullsfirst {
			columnName = columnName + " NULLS FIRST"
		} else {
			columnName = columnName + " NULLS LAST"
		}
	}
	return columnName, nil
}

// NewArrayField returns a new ArrayField representing an array column.
func NewArrayField(name string, table Table) ArrayField {
	return ArrayField{
		name:  name,
		table: table,
	}
}

// Array returns a new ArrayField representing a literal string value.
func Array(slice interface{}) ArrayField {
	return ArrayField{
		value: slice,
	}
}

// SetArray returns a FieldValueSet associating the ArrayField to the value
// i.e. 'SET field = value'. It only accepts ArrayField.
func (f ArrayField) Set(value ArrayField) FieldValueSet {
	return FieldValueSet{
		Field: f,
		Value: value,
	}
}

// As returns a new ArrayField with the new field Alias i.e. 'field AS Alias'.
func (f ArrayField) As(alias string) ArrayField {
	f.alias = alias
	return f
}

// Asc returns a new ArrayField indicating that it should be ordered in
// ascending order i.e. 'ORDER BY field ASC'.
func (f ArrayField) Asc() ArrayField {
	desc := false
	f.descending = &desc
	return f
}

// Desc returns a new ArrayField indicating that it should be ordered in
// descending order i.e. 'ORDER BY field DESC'.
func (f ArrayField) Desc() ArrayField {
	desc := true
	f.descending = &desc
	return f
}

// NullsFirst returns a new ArrayField indicating that it should be ordered
// with nulls first i.e. 'ORDER BY field NULLS FIRST'.
func (f ArrayField) NullsFirst() ArrayField {
	nullsfirst := true
	f.nullsfirst = &nullsfirst
	return f
}

// NullsLast returns a new ArrayField indicating that it should be ordered
// with nulls last i.e. 'ORDER BY field NULLS LAST'.
func (f ArrayField) NullsLast() ArrayField {
	nullsfirst := false
	f.nullsfirst = &nullsfirst
	return f
}

// IsNull returns an 'A IS NULL' Predicate.
func (f ArrayField) IsNull() Predicate {
	return UnaryPredicate{
		Operator: PredicateIsNull,
		Field:    f,
	}
}

// IsNotNull returns an 'A IS NOT NULL' Predicate.
func (f ArrayField) IsNotNull() Predicate {
	return UnaryPredicate{
		Operator: PredicateIsNotNull,
		Field:    f,
	}
}

// Eq returns an 'A = B' Predicate. It only accepts ArrayField.
func (f ArrayField) Eq(field ArrayField) Predicate {
	return BinaryPredicate{
		Operator:   PredicateEq,
		LeftField:  f,
		RightField: field,
	}
}

// Ne returns an 'A <> B' Predicate. It only accepts ArrayField.
func (f ArrayField) Ne(field ArrayField) Predicate {
	return BinaryPredicate{
		Operator:   PredicateNe,
		LeftField:  f,
		RightField: field,
	}
}

// Gt returns an 'A > B' Predicate. It only accepts ArrayField.
func (f ArrayField) Gt(field ArrayField) Predicate {
	return BinaryPredicate{
		Operator:   PredicateGt,
		LeftField:  f,
		RightField: field,
	}
}

// Ge returns an 'A >= B' Predicate. It only accepts ArrayField.
func (f ArrayField) Ge(field ArrayField) Predicate {
	return BinaryPredicate{
		Operator:   PredicateGe,
		LeftField:  f,
		RightField: field,
	}
}

// Lt returns an 'A < B' Predicate. It only accepts ArrayField.
func (f ArrayField) Lt(field ArrayField) Predicate {
	return BinaryPredicate{
		Operator:   PredicateLt,
		LeftField:  f,
		RightField: field,
	}
}

// Le returns an 'A <= B' Predicate. It only accepts ArrayField.
func (f ArrayField) Le(field ArrayField) Predicate {
	return BinaryPredicate{
		Operator:   PredicateLe,
		LeftField:  f,
		RightField: field,
	}
}

// Contains checks whether the subject ArrayField contains the object
// ArrayField.
func (f ArrayField) Contains(field ArrayField) Predicate {
	return CustomPredicate{
		Format: "? @> ?",
		Values: []interface{}{f, field},
	}
}

// Contains checks whether the subject ArrayField is contained by the object
// ArrayField.
func (f ArrayField) ContainedBy(field ArrayField) Predicate {
	return CustomPredicate{
		Format: "? <@ ?",
		Values: []interface{}{f, field},
	}
}

// Overlaps checks whether the subject ArrayField and the object ArrayField
// have any values in common.
func (f ArrayField) Overlaps(field ArrayField) Predicate {
	return CustomPredicate{
		Format: "? && ?",
		Values: []interface{}{f, field},
	}
}

// Concat concatenates the object ArrayField to the subject ArrayField.
func (f ArrayField) Concat(field ArrayField) Predicate {
	return CustomPredicate{
		Format: "? || ?",
		Values: []interface{}{f, field},
	}
}

// String implements the fmt.Stringer interface. It returns the string
// representation of an ArrayField.
func (f ArrayField) String() string {
	query, args := f.ToSQLExclude(nil)
	return MySQLInterpolateSQL(query, args...)
}

// GetAlias implements the Field interface. It returns the Alias of the
// ArrayField.
func (f ArrayField) GetAlias() string {
	return f.alias
}

// GetName implements the Field interface. It returns the Name of the
// ArrayField.
func (f ArrayField) GetName() string {
	return f.name
}
