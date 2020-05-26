package qx

import (
	"testing"

	"github.com/matryer/is"
)

func TestCustomField_ToSQL(t *testing.T) {
	type TT struct {
		DESCRIPTION            string
		field                  CustomField
		excludeTableQualifiers []string
		wantQuery              string
		wantArgs               []interface{}
	}
	tests := []TT{
		func() TT {
			DESCRIPTION := "nothing"
			field := CustomField{}
			wantQuery := ""
			return TT{DESCRIPTION, field, nil, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "if field empty, ASC/DESC/NULLS LAST/NULLS FIRST will not be added"
			field := CustomField{
				Format: "",
			}.Desc().NullsFirst()
			wantQuery := ""
			return TT{DESCRIPTION, field, nil, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "DESC NULLS FIRST"
			field := CustomField{
				Format: "lorem ipsum dolor sit amet",
			}.Desc().NullsFirst()
			wantQuery := "lorem ipsum dolor sit amet DESC NULLS FIRST"
			return TT{DESCRIPTION, field, nil, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "ASC NULLS FIRST"
			field := CustomField{
				Format: "lorem ipsum dolor sit amet",
			}.Asc().NullsLast()
			wantQuery := "lorem ipsum dolor sit amet ASC NULLS LAST"
			return TT{DESCRIPTION, field, nil, wantQuery, nil}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.DESCRIPTION, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			gotQuery, gotArgs := tt.field.ToSQLExclude(tt.excludeTableQualifiers)
			is.Equal(tt.wantQuery, gotQuery)
			is.Equal(tt.wantArgs, gotArgs)
		})
	}
}

func TestCustomField_GameTheNumbers(t *testing.T) {
	is := is.New(t)
	lipsum := "lorem ipsum dolor sit amet"
	field := CustomField{
		Format: lipsum,
		CustomSprintf: func(string, []interface{}, []string) (string, []interface{}) {
			return lipsum, nil
		},
	}.As("lipsum")
	query, _ := field.ToSQLExclude(nil)
	is.Equal(lipsum, query)
	field.IsNull()
	field.IsNotNull()
	field.Eq(field)
	field.Ne(field)
	field.Gt(field)
	field.Ge(field)
	field.Lt(field)
	field.Le(field)
	field.In(field)
	_ = field.String()
	is.Equal("lipsum", field.GetAlias())
	is.Equal(lipsum, field.GetName())
}
