package qx

type BinaryField struct {
	// BinaryField will be one of the following:

	// 1) Literal []byte value
	value *[]byte

	// 2) BYTEA/BLOB column
	alias string
	table Table
	name  string
}

// ToSQL marshals a BinaryField into an SQL query and args (as described in the
// BinaryField internal struct comments). If the BinaryField's table name
// appears in the excludeTableQualifiers list, the output column name will not
// be table qualified.
func (f BinaryField) ToSQLExclude(excludeTableQualifiers []string) (string, []interface{}) {
	// 1) Literal []byte value
	if f.value != nil {
		return "?", []interface{}{f.value}
	}

	// 3) BYTEA/BLOB column
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
	return columnName, nil
}

// NewBinaryField returns a new BinaryField representing a BYTEA/BLOB column.
func NewBinaryField(name string, table Table) BinaryField {
	return BinaryField{
		name:  name,
		table: table,
	}
}

// Bytes returns a new BinaryField representing a literal []byte value.
func Bytes(b []byte) BinaryField {
	return BinaryField{
		value: &b,
	}
}

// Set returns a FieldValueSet associating the BinaryField to the value i.e.
// 'SET field = value'.
func (f BinaryField) Set(val interface{}) FieldValueSet {
	return FieldValueSet{
		Field: f,
		Value: val,
	}
}

// SetBytes returns a FieldValueSet associating the BinaryField to the int value
// i.e. 'SET field = value'.
func (f BinaryField) SetBytes(b []byte) FieldValueSet {
	return FieldValueSet{
		Field: f,
		Value: b,
	}
}

// GetAlias implements the Field interface. It returns the Alias of the
// BinaryField.
func (f BinaryField) GetAlias() string {
	return f.alias
}

// GetName implements the Field interface. It returns the Name of the
// BinaryField.
func (f BinaryField) GetName() string {
	return f.name
}
