package qx

import (
	"testing"

	"github.com/matryer/is"
)

func TestArrayField_Literal_ToSQL(t *testing.T) {
	type TT struct {
		DESCRIPTION            string
		field                  ArrayField
		excludeTableQualifiers []string
		wantQuery              string
		wantArgs               []interface{}
	}
	tests := []TT{
		func() TT {
			DESCRIPTION := "empty []bool literal"
			field := Array([]bool{})
			wantQuery := "ARRAY[]::BOOLEAN[]"
			return TT{DESCRIPTION, field, nil, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "[]bool literal"
			field := Array([]bool{true, true, false, true})
			wantQuery := "ARRAY[?, ?, ?, ?]"
			wantArgs := []interface{}{true, true, false, true}
			return TT{DESCRIPTION, field, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			DESCRIPTION := "empty []float64 literal"
			field := Array([]float64{})
			wantQuery := "ARRAY[]::FLOAT[]"
			return TT{DESCRIPTION, field, nil, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "[]float64 literal"
			field := Array([]float64{22.7, 3.15, 4.0})
			wantQuery := "ARRAY[?, ?, ?]"
			wantArgs := []interface{}{22.7, 3.15, 4.0}
			return TT{DESCRIPTION, field, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			DESCRIPTION := "empty []int literal"
			field := Array([]int{})
			wantQuery := "ARRAY[]::INT[]"
			return TT{DESCRIPTION, field, nil, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "[]int literal"
			field := Array([]int{1, 2, 3, 4})
			wantQuery := "ARRAY[?, ?, ?, ?]"
			wantArgs := []interface{}{1, 2, 3, 4}
			return TT{DESCRIPTION, field, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			DESCRIPTION := "empty []int64 literal"
			field := Array([]int64{})
			wantQuery := "ARRAY[]::BIGINT[]"
			return TT{DESCRIPTION, field, nil, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "[]int64 literal"
			field := Array([]int64{1, 2, 3, 4})
			wantQuery := "ARRAY[?, ?, ?, ?]"
			wantArgs := []interface{}{int64(1), int64(2), int64(3), int64(4)}
			return TT{DESCRIPTION, field, nil, wantQuery, wantArgs}
		}(),
		func() TT {
			DESCRIPTION := "empty []string literal"
			field := Array([]string{})
			wantQuery := "ARRAY[]::TEXT[]"
			return TT{DESCRIPTION, field, nil, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "[]string literal"
			field := Array([]string{"apple", "banana", "cucumber"})
			wantQuery := "ARRAY[?, ?, ?]"
			wantArgs := []interface{}{"apple", "banana", "cucumber"}
			return TT{DESCRIPTION, field, nil, wantQuery, wantArgs}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.DESCRIPTION, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			gotQuery, gotArgs := tt.field.ToSQL(nil)
			is.Equal(tt.wantQuery, gotQuery)
			is.Equal(tt.wantArgs, gotArgs)
		})
	}
}

func TestArrayField_Column_ToSQL(t *testing.T) {
	type TT struct {
		DESCRIPTION            string
		field                  ArrayField
		excludeTableQualifiers []string
		wantQuery              string
		wantArgs               []interface{}
	}
	tests := []TT{
		func() TT {
			DESCRIPTION := "empty []bool literal"
			field := Array([]bool{})
			wantQuery := "ARRAY[]::BOOLEAN[]"
			return TT{DESCRIPTION, field, nil, wantQuery, nil}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.DESCRIPTION, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			gotQuery, gotArgs := tt.field.ToSQL(nil)
			is.Equal(tt.wantQuery, gotQuery)
			is.Equal(tt.wantArgs, gotArgs)
		})
	}
}
