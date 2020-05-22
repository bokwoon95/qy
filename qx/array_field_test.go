package qx

import (
	"fmt"
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
		func() TT {
			DESCRIPTION := "unspported type"
			field := Array("yeeehaw")
			wantQuery := "(unknown array type: only []bool/[]float64/[]int64/[]string/[]int slices are supported.)"
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
			DESCRIPTION := "basic ArrayField"
			field := FILM().SPECIAL_FEATURES
			wantQuery := "film.special_features"
			return TT{DESCRIPTION, field, nil, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "respect excludeTableQualifiers"
			film := FILM()
			field := film.SPECIAL_FEATURES
			wantQuery := "special_features"
			excludeTableQualifiers := []string{film.GetAlias(), film.GetName()}
			return TT{DESCRIPTION, field, excludeTableQualifiers, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "ArrayField with table alias (and column alias)"
			film := FILM().As("f")
			field := film.SPECIAL_FEATURES.As("speshul_fittures")
			is := is.New(t)
			is.Equal("speshul_fittures", field.GetAlias())
			wantQuery := "f.special_features"
			return TT{DESCRIPTION, field, nil, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "ArrayField DESC NULLS FIRST"
			field := FILM().SPECIAL_FEATURES.Desc().NullsFirst()
			wantQuery := "film.special_features DESC NULLS FIRST"
			return TT{DESCRIPTION, field, nil, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "ArrayField ASC NULLS LAST"
			field := FILM().SPECIAL_FEATURES.Asc().NullsLast()
			wantQuery := "film.special_features ASC NULLS LAST"
			return TT{DESCRIPTION, field, nil, wantQuery, nil}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.DESCRIPTION, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			gotQuery, gotArgs := tt.field.ToSQL(tt.excludeTableQualifiers)
			is.Equal(tt.wantQuery, gotQuery)
			is.Equal(tt.wantArgs, gotArgs)
		})
	}
}

func TestArrayField_GameTheNumbers(t *testing.T) {
	is := is.New(t)
	var p Predicate
	f := FILM().SPECIAL_FEATURES
	is.Equal("", f.GetAlias())
	is.Equal("special_features", f.GetName())
	var _ Field = f
	// Set
	f.Set(Array([]int{1, 2, 3}))
	// Predicates
	stringify := func(p Predicate) string {
		query, args := p.ToSQL(nil)
		return MySQLInterpolateSQL(query, args...)
	}
	f.IsNull()
	f.IsNotNull()
	f.Eq(f)
	f.Ne(f)
	f.Gt(f)
	f.Ge(f)
	f.Lt(f)
	f.Le(f)
	p = f.Contains(Array([]string{"Trailers", "Behind the Scenes"}))
	is.Equal("film.special_features @> ARRAY['Trailers', 'Behind the Scenes']", stringify(p))
	p = Array([]string{"Trailers", "Behind the Scenes"}).ContainedBy(f)
	is.Equal("ARRAY['Trailers', 'Behind the Scenes'] <@ film.special_features", stringify(p))
	p = f.Overlaps(Array([]string{"Trailers", "Behind the Scenes"}))
	is.Equal("film.special_features && ARRAY['Trailers', 'Behind the Scenes']", stringify(p))
	p = f.Concat(Array([]string{"Trailers", "Behind the Scenes"}))
	is.Equal("film.special_features || ARRAY['Trailers', 'Behind the Scenes']", stringify(p))
	// fmt.Stringer
	fmt.Println(f)
}
