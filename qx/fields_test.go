package qx

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestFieldLiteral(t *testing.T) {
	is := is.New(t)
	var query string
	var args []interface{}
	const abc FieldLiteral = "abc"
	var _ Field = abc
	is.Equal("abc", abc.GetName())
	is.Equal("", abc.GetAlias())
	query, args = abc.ToSQLExclude(nil)
	is.Equal("abc", query)
	is.Equal(nil, args)
}

func TestFields_WriteSQL(t *testing.T) {
	type TT struct {
		DESCRIPTION            string
		fields                 Fields
		excludeTableQualifiers []string
		wantWritten            bool
		wantQuery              string
		wantArgs               []interface{}
	}
	tests := []TT{
		func() TT {
			DESCRIPTION := "0 fields"
			fields := Fields{}
			wantWritten := false
			wantQuery := ""
			return TT{DESCRIPTION, fields, nil, wantWritten, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "1 field"
			fields := Fields{
				CUSTOMER().CUSTOMER_ID,
			}
			wantWritten := true
			wantQuery := "customer.customer_id"
			return TT{DESCRIPTION, fields, nil, wantWritten, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "3 fields"
			fields := Fields{
				CUSTOMER().CUSTOMER_ID,
				CUSTOMER().FIRST_NAME,
				String("somebody_once_told_me"),
				nil,
				CustomField{},
			}
			wantWritten := true
			wantQuery := "customer.customer_id, customer.first_name, ?, NULL"
			wantArgs := []interface{}{"somebody_once_told_me"}
			return TT{DESCRIPTION, fields, nil, wantWritten, wantQuery, wantArgs}
		}(),
		func() TT {
			DESCRIPTION := "3 fields + excludeTableQualifiers"
			fields := Fields{
				CUSTOMER().CUSTOMER_ID,
				CUSTOMER().FIRST_NAME,
				String("somebody_once_told_me"),
				nil,
				CustomField{},
			}
			excludeTableQualifiers := []string{CUSTOMER().GetName()}
			wantWritten := true
			wantQuery := "customer_id, first_name, ?, NULL"
			wantArgs := []interface{}{"somebody_once_told_me"}
			return TT{DESCRIPTION, fields, excludeTableQualifiers, wantWritten, wantQuery, wantArgs}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.DESCRIPTION, func(t *testing.T) {
			t.Parallel()
			var is = is.New(t)
			var buf = &strings.Builder{}
			var args []interface{}
			buf.WriteString("xyz")
			var written = tt.fields.WriteSQL(buf, &args, "", "", tt.excludeTableQualifiers)
			is.Equal(tt.wantWritten, written)
			if written {
				is.Equal("xyz"+" "+tt.wantQuery, buf.String())
			} else {
				is.Equal("xyz"+tt.wantQuery, buf.String())
			}
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestFields_WriteSQLWithAlias(t *testing.T) {
	type TT struct {
		DESCRIPTION            string
		fields                 Fields
		excludeTableQualifiers []string
		wantWritten            bool
		wantQuery              string
		wantArgs               []interface{}
	}
	tests := []TT{
		func() TT {
			DESCRIPTION := "0 fields"
			fields := Fields{}
			wantWritten := false
			wantQuery := ""
			return TT{DESCRIPTION, fields, nil, wantWritten, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "1 field"
			fields := Fields{
				CUSTOMER().CUSTOMER_ID.As("c_id"),
			}
			wantWritten := true
			wantQuery := "customer.customer_id AS c_id"
			return TT{DESCRIPTION, fields, nil, wantWritten, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "multiple fields"
			fields := Fields{
				CUSTOMER().CUSTOMER_ID.As("c_id"),
				CUSTOMER().FIRST_NAME.As("c_first_name"),
				String("somebody_once_told_me").As("smsh_mth"),
				nil,
				CustomField{},
			}
			wantWritten := true
			wantQuery := "customer.customer_id AS c_id, customer.first_name AS c_first_name, ? AS smsh_mth, NULL"
			wantArgs := []interface{}{"somebody_once_told_me"}
			return TT{DESCRIPTION, fields, nil, wantWritten, wantQuery, wantArgs}
		}(),
		func() TT {
			DESCRIPTION := "multiple fields + excludeTableQualifiers"
			fields := Fields{
				CUSTOMER().CUSTOMER_ID.As("c_id"),
				CUSTOMER().FIRST_NAME.As("c_first_name"),
				String("somebody_once_told_me").As("smsh_mth"),
				nil,
				CustomField{},
			}
			excludeTableQualifiers := []string{CUSTOMER().GetName()}
			wantWritten := true
			wantQuery := "customer_id AS c_id, first_name AS c_first_name, ? AS smsh_mth, NULL"
			wantArgs := []interface{}{"somebody_once_told_me"}
			return TT{DESCRIPTION, fields, excludeTableQualifiers, wantWritten, wantQuery, wantArgs}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.DESCRIPTION, func(t *testing.T) {
			t.Parallel()
			var is = is.New(t)
			var buf = &strings.Builder{}
			var args []interface{}
			buf.WriteString("xyz")
			var written = tt.fields.WriteSQLWithAlias(buf, &args, "", "", tt.excludeTableQualifiers)
			is.Equal(tt.wantWritten, written)
			if written {
				is.Equal("xyz"+" "+tt.wantQuery, buf.String())
			} else {
				is.Equal("xyz"+tt.wantQuery, buf.String())
			}
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestFieldsValueSets_WriteSQL(t *testing.T) {
	type TT struct {
		DESCRIPTION            string
		sets                   FieldValueSets
		excludeTableQualifiers []string
		wantWritten            bool
		wantQuery              string
		wantArgs               []interface{}
	}
	tests := []TT{
		func() TT {
			DESCRIPTION := "0 sets"
			sets := FieldValueSets{}
			wantWritten := false
			wantQuery := ""
			return TT{DESCRIPTION, sets, nil, wantWritten, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "1 set"
			sets := FieldValueSets{
				CUSTOMER().CUSTOMER_ID.SetInt(5),
			}
			wantWritten := true
			wantQuery := "customer.customer_id = ?"
			wantArgs := []interface{}{5}
			return TT{DESCRIPTION, sets, nil, wantWritten, wantQuery, wantArgs}
		}(),
		func() TT {
			DESCRIPTION := "multiple sets"
			sets := FieldValueSets{
				CUSTOMER().CUSTOMER_ID.SetInt(5),
				CUSTOMER().FIRST_NAME.SetString("norman"),
				String("somebody_once_told_me").Set(String("smsh_mth")),
			}
			wantWritten := true
			wantQuery := "customer.customer_id = ?, customer.first_name = ?, ? = ?"
			wantArgs := []interface{}{5, "norman", "somebody_once_told_me", "smsh_mth"}
			return TT{DESCRIPTION, sets, nil, wantWritten, wantQuery, wantArgs}
		}(),
		func() TT {
			DESCRIPTION := "multiple sets + excludeTableQualifiers"
			sets := FieldValueSets{
				CUSTOMER().CUSTOMER_ID.SetInt(5),
				CUSTOMER().FIRST_NAME.SetString("norman"),
				String("somebody_once_told_me").Set(String("smsh_mth")),
				FieldValueSet{},
				FieldValueSet{Field: CustomField{}},
			}
			excludeTableQualifiers := []string{CUSTOMER().GetName()}
			wantWritten := true
			wantQuery := "customer_id = ?, first_name = ?, ? = ?"
			wantArgs := []interface{}{5, "norman", "somebody_once_told_me", "smsh_mth"}
			return TT{DESCRIPTION, sets, excludeTableQualifiers, wantWritten, wantQuery, wantArgs}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.DESCRIPTION, func(t *testing.T) {
			t.Parallel()
			var is = is.New(t)
			var buf = &strings.Builder{}
			var args []interface{}
			buf.WriteString("xyz")
			var written = tt.sets.WriteSQL(buf, &args, "", "", tt.excludeTableQualifiers)
			is.Equal(tt.wantWritten, written)
			if written {
				is.Equal("xyz"+" "+tt.wantQuery, buf.String())
			} else {
				is.Equal("xyz"+tt.wantQuery, buf.String())
			}
			is.Equal(tt.wantArgs, args)
		})
	}
}
