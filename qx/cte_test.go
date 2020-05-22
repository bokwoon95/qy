package qx

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestCTE_WriteSQL(t *testing.T) {
	type TT struct {
		DESCRIPTION string
		ctes        CTEs
		wantWritten bool
		wantQuery   string
		wantArgs    []interface{}
	}
	tests := []TT{}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.DESCRIPTION, func(t *testing.T) {
			t.Parallel()
			var is = is.New(t)
			var buf = &strings.Builder{}
			var args []interface{}
			var written = tt.ctes.WriteSQL(buf, &args)
			is.Equal(tt.wantWritten, written)
			is.Equal(tt.wantQuery, buf.String())
			is.Equal(tt.wantArgs, args)
		})
	}
}

func TestCTE_GameTheNumbers(t *testing.T) {
	is := is.New(t)
	var query string
	cte := NewCTE("my_cte", nil)
	cte.ToSQL()
	cte.GetAlias()
	cte.GetName()
	query, _ = cte.Get("xkcd").ToSQL(nil)
	is.Equal("my_cte.xkcd", query)
	query, _ = cte.Get("xkcd").ToSQL([]string{"my_cte"})
	is.Equal("my_cte.xkcd", query) // CTEs are not expected to respect excludeTableQualifiers, only Predicates and Fields are
	aliasedcte := cte.As("other_cte")
	is.Equal("other_cte", aliasedcte.GetAlias())
	is.Equal("my_cte", aliasedcte.GetName())
	query, _ = aliasedcte.Get("xkcd").ToSQL(nil)
	is.Equal("other_cte.xkcd", query)
}
