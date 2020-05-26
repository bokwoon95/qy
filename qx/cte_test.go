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
	tests := []TT{
		func() TT {
			DESCRIPTION := "basic"
			ctes := CTEs{
				CTE{Name: "cte0"},
				CTE{Name: "cte1", Query: CustomQuery{
					Format: "Everyone was curious about the large white blimp that appeared overnight.",
				}},
				CTE{Name: "cte2", Query: CustomQuery{
					Format: "They did nothing as the raccoon attacked the lady’s bag of food.",
				}},
				CTE{Name: "cte3", Query: CustomQuery{
					Format: "There can never be too many cherries on an ice cream sundae.",
				}},
				CTE{Name: "cte4", Query: CustomQuery{}},
			}
			wantQuery := "WITH cte1 AS (Everyone was curious about the large white blimp that appeared overnight.)" +
				", cte2 AS (They did nothing as the raccoon attacked the lady’s bag of food.)" +
				", cte3 AS (There can never be too many cherries on an ice cream sundae.)"
			return TT{DESCRIPTION, ctes, true, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "nothing"
			ctes := CTEs{}
			wantQuery := ""
			return TT{DESCRIPTION, ctes, false, wantQuery, nil}
		}(),
	}
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
	query, _ = cte.Get("xkcd").ToSQLExclude(nil)
	is.Equal("my_cte.xkcd", query)
	query, _ = cte.Get("xkcd").ToSQLExclude([]string{"my_cte"})
	is.Equal("my_cte.xkcd", query) // CTEs are not expected to respect excludeTableQualifiers, only Predicates and Fields are
	aliasedcte := cte.As("other_cte")
	is.Equal("other_cte", aliasedcte.GetAlias())
	is.Equal("my_cte", aliasedcte.GetName())
	query, _ = aliasedcte.ToSQL()
	is.Equal("my_cte", query)
	query, _ = aliasedcte.Get("xkcd").ToSQLExclude(nil)
	is.Equal("other_cte.xkcd", query)
}
