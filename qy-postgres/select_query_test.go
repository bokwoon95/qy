package qy

import (
	"log"
	"os"
	"testing"

	"github.com/bokwoon95/qy/qx"
	"github.com/bokwoon95/qy/tables-postgres"
	_ "github.com/lib/pq"
	"github.com/matryer/is"
)

func TestSelectQuery_Get(t *testing.T) {
	type TT struct {
		DESCRIPTION string
		f           qx.Field
		wantQuery   string
		wantArgs    []interface{}
	}
	baseSelect := NewSelectQuery()
	baseSelect.Log = log.New(os.Stdout, "", 0)
	tests := []TT{
		{
			"SelectQuery Get (explicit alias)",
			baseSelect.From(tables.CUSTOMER()).As("selected_customers").Get("first_name"),
			"selected_customers.first_name",
			nil,
		},
		func() TT {
			DESCRIPTION := "SelectQuery Get (implicit alias)"
			q := baseSelect.From(tables.CUSTOMER())
			alias := q.GetAlias()
			if alias == "" {
				t.Fatalf("an alias should have automatically been assigned to SelectQuery, found empty alias")
			}
			return TT{DESCRIPTION, q.Get("first_name"), alias + ".first_name", nil}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.DESCRIPTION, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			gotQuery, gotArgs := tt.f.ToSQL(nil)
			is.Equal(tt.wantQuery, gotQuery)
			is.Equal(tt.wantArgs, gotArgs)
		})
	}
}

func TestSelectQuery_With(t *testing.T) {
	q := NewSelectQuery()
	wantQuery, wantArgs := "", []interface{}{}

	apac_customers := qx.NewCTE("apac_customers", func() qx.Query {
		cust, addr, city, coun := tables.CUSTOMER().As("cust"), tables.ADDRESS().As("addr"), tables.CITY(), tables.COUNTRY().As("coun")
		return Select(cust.FIRST_NAME, cust.LAST_NAME, addr.ADDRESS, city.CITY).
			From(cust).Join(addr, addr.ADDRESS_ID.Eq(cust.ADDRESS_ID)).
			LeftJoin(city, city.CITY_ID.Eq(addr.CITY_ID)).
			LeftJoin(coun, coun.COUNTRY_ID.Eq(city.COUNTRY_ID)).
			Where(Predicatef("? ILIKE ANY (ARRAY[?])", coun.COUNTRY, []string{"china", "japan", "australia", "korea"}))
	}())
	wantQuery += "WITH apac_customers AS" +
		" (SELECT cust.first_name, cust.last_name, addr.address, city.city" +
		" FROM customer AS cust JOIN address AS addr ON addr.address_id = cust.address_id" +
		" LEFT JOIN city ON city.city_id = addr.city_id" +
		" LEFT JOIN country AS coun ON coun.country_id = city.country_id" +
		" WHERE coun.country ILIKE ANY (ARRAY[$1, $2, $3, $4]))"
	wantArgs = append(wantArgs, "china", "japan", "australia", "korea")

	kids_films := qx.NewCTE("kids_films", func() qx.Query {
		film, fica, cate := tables.FILM(), tables.FILM_CATEGORY().As("fica"), tables.CATEGORY().As("cate")
		return Select(film.TITLE, film.RELEASE_YEAR, cate.NAME).
			From(film).Join(fica, fica.FILM_ID.Eq(film.FILM_ID)).
			LeftJoin(cate, cate.CATEGORY_ID.Eq(fica.CATEGORY_ID)).
			Where(Predicatef("? IN (?)", cate.NAME, []string{"Action", "Animation", "Children"}))
	}())
	wantQuery += ", kids_films AS" +
		" (SELECT film.title, film.release_year, cate.name" +
		" FROM film JOIN film_category AS fica ON fica.film_id = film.film_id" +
		" LEFT JOIN category AS cate ON cate.category_id = fica.category_id" +
		" WHERE cate.name IN ($5, $6, $7))"
	wantArgs = append(wantArgs, "Action", "Animation", "Children")

	films_stores := qx.NewCTE("films_stores", func() qx.Query {
		film, inve, stor := tables.FILM(), tables.INVENTORY().As("inve"), tables.STORE().As("stor")
		return Select(film.FILM_ID, film.TITLE, stor.ADDRESS_ID).
			From(film).Join(inve, inve.FILM_ID.Eq(film.FILM_ID)).
			Join(stor, stor.STORE_ID.Eq(inve.STORE_ID)).
			Where(stor.ADDRESS_ID.IsNotNull())
	}())
	wantQuery += ", films_stores AS" +
		" (SELECT film.film_id, film.title, stor.address_id" +
		" FROM film JOIN inventory AS inve ON inve.film_id = film.film_id" +
		" JOIN store AS stor ON stor.store_id = inve.store_id" +
		" WHERE stor.address_id IS NOT NULL)"
	wantArgs = append(wantArgs)

	fs1, fs2 := films_stores.As("fs1"), films_stores.As("fs2")
	q = q.With(qx.CTE{}, apac_customers, kids_films, films_stores).
		Select(
			apac_customers.Get("first_name"),
			kids_films.Get("release_year"),
			fs1.Get("title").As("fs1_title"),
			fs2.Get("title").As("fs2_title"),
		).
		From(apac_customers).
		CrossJoin(fs1).
		Join(fs2, fs2.Get("film_id").Eq(fs1.Get("film_id"))).
		CrossJoin(kids_films)
	wantQuery += " SELECT apac_customers.first_name, kids_films.release_year, fs1.title AS fs1_title, fs2.title AS fs2_title" +
		" FROM apac_customers" +
		" CROSS JOIN films_stores AS fs1" +
		" JOIN films_stores AS fs2 ON fs2.film_id = fs1.film_id" +
		" CROSS JOIN kids_films"
	wantArgs = append(wantArgs)

	is := is.New(t)
	gotQuery, gotArgs := q.ToSQL()
	is.Equal(wantQuery, gotQuery)
	is.Equal(wantArgs, gotArgs)
}

func TestSelectQuery_Select(t *testing.T) {
	type TT struct {
		DESCRIPTION string
		q           SelectQuery
		wantQuery   string
		wantArgs    []interface{}
	}
	baseSelect := NewSelectQuery()
	tests := []TT{
		func() TT {
			DESCRIPTION := "basic select"
			cust := tables.CUSTOMER()
			q := baseSelect.Select(cust.FIRST_NAME, cust.LAST_NAME, cust.EMAIL)
			wantQuery := "SELECT customer.first_name, customer.last_name, customer.email"
			return TT{DESCRIPTION, q, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "multiple select calls also work"
			cust := tables.CUSTOMER()
			q := baseSelect.Select(cust.FIRST_NAME).Select(cust.LAST_NAME).Select(cust.EMAIL)
			wantQuery := "SELECT customer.first_name, customer.last_name, customer.email"
			return TT{DESCRIPTION, q, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "column aliases work"
			acto := tables.ACTOR()
			q := baseSelect.Select(acto.FIRST_NAME, acto.ACTOR_ID.As("id"), acto.LAST_UPDATE.As("update"))
			wantQuery := "SELECT actor.first_name, actor.actor_id AS id, actor.last_update AS update"
			return TT{DESCRIPTION, q, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "select distinct"
			cust := tables.CUSTOMER().As("cust")
			q := baseSelect.SelectDistinct(cust.FIRST_NAME, cust.EMAIL)
			wantQuery := "SELECT DISTINCT cust.first_name, cust.email"
			return TT{DESCRIPTION, q, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "select distinct on"
			cust := tables.CUSTOMER().As("cust")
			q := baseSelect.SelectDistinctOn(cust.EMAIL, cust.STORE_ID)(cust.FIRST_NAME, cust.LAST_NAME)
			wantQuery := "SELECT DISTINCT ON (cust.email, cust.store_id) cust.first_name, cust.last_name"
			return TT{DESCRIPTION, q, wantQuery, nil}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.DESCRIPTION, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			gotQuery, gotArgs := tt.q.ToSQL()
			is.Equal(tt.wantQuery, gotQuery)
			is.Equal(tt.wantArgs, gotArgs)
		})
	}
}

func TestSelectQuery_From(t *testing.T) {
	type TT struct {
		DESCRIPTION string
		q           SelectQuery
		wantQuery   string
		wantArgs    []interface{}
	}
	baseSelect := NewSelectQuery()
	tests := []TT{
		func() TT {
			DESCRIPTION := "no table alias"
			cust := tables.CUSTOMER()
			q := baseSelect.From(cust).Select(cust.CUSTOMER_ID)
			wantQuery := "SELECT customer.customer_id FROM customer"
			return TT{DESCRIPTION, q, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "with table alias"
			cust := tables.CUSTOMER().As("cust")
			q := baseSelect.From(cust).Select(cust.CUSTOMER_ID)
			wantQuery := "SELECT cust.customer_id FROM customer AS cust"
			return TT{DESCRIPTION, q, wantQuery, nil}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.DESCRIPTION, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			gotQuery, gotArgs := tt.q.ToSQL()
			is.Equal(tt.wantQuery, gotQuery)
			is.Equal(tt.wantArgs, gotArgs)
		})
	}
}

func TestSelectQuery_Joins(t *testing.T) {
	type TT struct {
		DESCRIPTION string
		q           SelectQuery
		wantQuery   string
		wantArgs    []interface{}
	}
	baseSelect := NewSelectQuery()
	baseSelect.Log = log.New(os.Stdout, "", 0)
	tests := []TT{
		func() TT {
			DESCRIPTION := "all joinGroups"
			cust, addr := tables.CUSTOMER().As("cust"), tables.ADDRESS().As("addr")
			q := baseSelect.From(cust).
				Join(addr, addr.ADDRESS_ID.Eq(cust.ADDRESS_ID)).
				LeftJoin(addr, addr.ADDRESS_ID.Eq(cust.ADDRESS_ID)).
				RightJoin(addr, addr.ADDRESS_ID.Eq(cust.ADDRESS_ID)).
				CrossJoin(addr)
			// You can't join tables with the same alias in SQL, this is just an example
			wantQuery := "FROM customer AS cust" +
				" JOIN address AS addr ON addr.address_id = cust.address_id" +
				" LEFT JOIN address AS addr ON addr.address_id = cust.address_id" +
				" RIGHT JOIN address AS addr ON addr.address_id = cust.address_id" +
				" CROSS JOIN address AS addr"
			return TT{DESCRIPTION, q, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "joining a subquery with explicit alias"
			cust, addr := tables.CUSTOMER().As("cust"), tables.ADDRESS().As("addr")
			subquery := Select(cust.ADDRESS_ID).From(cust).Where(cust.STORE_ID.GtInt(4)).As("subquery")
			q := baseSelect.From(addr).Join(subquery, subquery.Get("address_id").Eq(addr.ADDRESS_ID))
			wantQuery := "FROM address AS addr" +
				" JOIN (SELECT cust.address_id FROM customer AS cust WHERE cust.store_id > $1) AS subquery" +
				" ON subquery.address_id = addr.address_id"
			return TT{DESCRIPTION, q, wantQuery, []interface{}{4}}
		}(),
		func() TT {
			DESCRIPTION := "joining a subquery with implicit alias"
			cust, addr := tables.CUSTOMER().As("cust"), tables.ADDRESS().As("addr")
			subquery := Select(cust.ADDRESS_ID).From(cust).Where(cust.STORE_ID.GtInt(4))
			q := baseSelect.From(addr).Join(subquery, subquery.Get("address_id").Eq(addr.ADDRESS_ID))
			wantQuery := "FROM address AS addr" +
				" JOIN (SELECT cust.address_id FROM customer AS cust WHERE cust.store_id > $1) AS " + subquery.GetAlias() +
				" ON " + subquery.GetAlias() + ".address_id = addr.address_id"
			return TT{DESCRIPTION, q, wantQuery, []interface{}{4}}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.DESCRIPTION, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			gotQuery, gotArgs := tt.q.ToSQL()
			is.Equal(tt.wantQuery, gotQuery)
			is.Equal(tt.wantArgs, gotArgs)
		})
	}
}

func TestSelectQuery_Where(t *testing.T) {
	type TT struct {
		DESCRIPTION string
		q           SelectQuery
		wantQuery   string
		wantArgs    []interface{}
	}
	baseSelect := NewSelectQuery()
	tests := []TT{
		func() TT {
			DESCRIPTION := "basic where (implicit and)"
			cust := tables.CUSTOMER()
			q := baseSelect.Where(
				cust.CUSTOMER_ID.EqInt(22),
				cust.FIRST_NAME.ILikeString("%bob%"),
				cust.EMAIL.IsNotNull(),
			)
			wantQuery := "WHERE customer.customer_id = $1 AND customer.first_name ILIKE $2 AND customer.email IS NOT NULL"
			return TT{DESCRIPTION, q, wantQuery, []interface{}{22, "%bob%"}}
		}(),
		func() TT {
			DESCRIPTION := "basic where (explicit and)"
			cust := tables.CUSTOMER()
			q := baseSelect.Where(
				qx.And(
					cust.CUSTOMER_ID.EqInt(22),
					cust.FIRST_NAME.ILikeString("%bob%"),
					cust.EMAIL.IsNotNull(),
				),
			)
			wantQuery := "WHERE customer.customer_id = $1 AND customer.first_name ILIKE $2 AND customer.email IS NOT NULL"
			return TT{DESCRIPTION, q, wantQuery, []interface{}{22, "%bob%"}}
		}(),
		func() TT {
			DESCRIPTION := "basic where (explicit or)"
			cust := tables.CUSTOMER()
			q := baseSelect.Where(
				qx.Or(
					cust.CUSTOMER_ID.EqInt(22),
					cust.FIRST_NAME.ILikeString("%bob%"),
					cust.EMAIL.IsNotNull(),
				),
			)
			wantQuery := "WHERE customer.customer_id = $1 OR customer.first_name ILIKE $2 OR customer.email IS NOT NULL"
			return TT{DESCRIPTION, q, wantQuery, []interface{}{22, "%bob%"}}
		}(),
		func() TT {
			DESCRIPTION := "complex predicate"
			c1, c2 := tables.CUSTOMER().As("c1"), tables.CUSTOMER().As("c2")
			q := baseSelect.Where(
				c1.CUSTOMER_ID.EqInt(69),
				c1.FIRST_NAME.ILikeString("%bob%"),
				c1.EMAIL.IsNotNull(),
				qx.Or(
					c2.CUSTOMER_ID.EqInt(420),
					c2.FIRST_NAME.LikeString("%virgil%"),
					c2.EMAIL.IsNull(),
				),
			)
			wantQuery := "WHERE c1.customer_id = $1 AND c1.first_name ILIKE $2 AND c1.email IS NOT NULL" +
				" AND (c2.customer_id = $3 OR c2.first_name LIKE $4 OR c2.email IS NULL)"
			return TT{DESCRIPTION, q, wantQuery, []interface{}{69, "%bob%", 420, "%virgil%"}}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.DESCRIPTION, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			gotQuery, gotArgs := tt.q.ToSQL()
			is.Equal(tt.wantQuery, gotQuery)
			is.Equal(tt.wantArgs, gotArgs)
		})
	}
}

func TestSelectQuery_GroupBy(t *testing.T) {
	type TT struct {
		DESCRIPTION string
		q           SelectQuery
		wantQuery   string
		wantArgs    []interface{}
	}
	baseSelect := NewSelectQuery()
	tests := []TT{
		func() TT {
			DESCRIPTION := "basic group by"
			cust := tables.CUSTOMER()
			q := baseSelect.GroupBy(cust.CUSTOMER_ID, cust.STORE_ID, cust.ACTIVE)
			wantQuery := "GROUP BY customer.customer_id, customer.store_id, customer.active"
			return TT{DESCRIPTION, q, wantQuery, nil}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.DESCRIPTION, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			gotQuery, gotArgs := tt.q.ToSQL()
			is.Equal(tt.wantQuery, gotQuery)
			is.Equal(tt.wantArgs, gotArgs)
		})
	}
}

func TestSelectQuery_Having(t *testing.T) {
	type TT struct {
		DESCRIPTION string
		q           SelectQuery
		wantQuery   string
		wantArgs    []interface{}
	}
	baseSelect := NewSelectQuery()
	tests := []TT{
		func() TT {
			DESCRIPTION := "basic having (implicit and)"
			cust := tables.CUSTOMER()
			q := baseSelect.Having(
				cust.CUSTOMER_ID.EqInt(22),
				cust.FIRST_NAME.ILikeString("%bob%"),
				cust.EMAIL.IsNotNull(),
			)
			wantQuery := "HAVING customer.customer_id = $1 AND customer.first_name ILIKE $2 AND customer.email IS NOT NULL"
			return TT{DESCRIPTION, q, wantQuery, []interface{}{22, "%bob%"}}
		}(),
		func() TT {
			DESCRIPTION := "basic having (explicit and)"
			cust := tables.CUSTOMER()
			q := baseSelect.Having(
				qx.And(
					cust.CUSTOMER_ID.EqInt(22),
					cust.FIRST_NAME.ILikeString("%bob%"),
					cust.EMAIL.IsNotNull(),
				),
			)
			wantQuery := "HAVING customer.customer_id = $1 AND customer.first_name ILIKE $2 AND customer.email IS NOT NULL"
			return TT{DESCRIPTION, q, wantQuery, []interface{}{22, "%bob%"}}
		}(),
		func() TT {
			DESCRIPTION := "basic having (explicit or)"
			cust := tables.CUSTOMER()
			q := baseSelect.Having(
				qx.Or(
					cust.CUSTOMER_ID.EqInt(22),
					cust.FIRST_NAME.ILikeString("%bob%"),
					cust.EMAIL.IsNotNull(),
				),
			)
			wantQuery := "HAVING customer.customer_id = $1 OR customer.first_name ILIKE $2 OR customer.email IS NOT NULL"
			return TT{DESCRIPTION, q, wantQuery, []interface{}{22, "%bob%"}}
		}(),
		func() TT {
			DESCRIPTION := "complex predicate"
			c1, c2 := tables.CUSTOMER().As("c1"), tables.CUSTOMER().As("c2")
			q := baseSelect.Having(
				c1.CUSTOMER_ID.EqInt(69),
				c1.FIRST_NAME.ILikeString("%bob%"),
				c1.EMAIL.IsNotNull(),
				qx.Or(
					c2.CUSTOMER_ID.EqInt(420),
					c2.FIRST_NAME.LikeString("%virgil%"),
					c2.EMAIL.IsNull(),
				),
			)
			wantQuery := "HAVING c1.customer_id = $1 AND c1.first_name ILIKE $2 AND c1.email IS NOT NULL" +
				" AND (c2.customer_id = $3 OR c2.first_name LIKE $4 OR c2.email IS NULL)"
			return TT{DESCRIPTION, q, wantQuery, []interface{}{69, "%bob%", 420, "%virgil%"}}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.DESCRIPTION, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			gotQuery, gotArgs := tt.q.ToSQL()
			is.Equal(tt.wantQuery, gotQuery)
			is.Equal(tt.wantArgs, gotArgs)
		})
	}
}

func TestSelectQuery_OrderBy(t *testing.T) {
	type TT struct {
		DESCRIPTION string
		q           SelectQuery
		wantQuery   string
		wantArgs    []interface{}
	}
	baseSelect := NewSelectQuery()
	tests := []TT{
		func() TT {
			DESCRIPTION := "basic order by"
			cust := tables.CUSTOMER().As("cust")
			q := baseSelect.Select().OrderBy(cust.CUSTOMER_ID, cust.STORE_ID, cust.ACTIVE)
			wantQuery := "ORDER BY cust.customer_id, cust.store_id, cust.active"
			return TT{DESCRIPTION, q, wantQuery, nil}
		}(),
		func() TT {
			DESCRIPTION := "Asc, Desc, NullsFirst and NullsLast"
			cust := tables.CUSTOMER().As("cust")
			q := baseSelect.OrderBy(
				cust.CUSTOMER_ID,
				cust.CUSTOMER_ID.Asc(),
				cust.CUSTOMER_ID.Desc(),
				cust.CUSTOMER_ID.NullsLast(),
				cust.CUSTOMER_ID.Asc().NullsFirst(),
				cust.CUSTOMER_ID.Desc().NullsLast(),
			)
			wantQuery := "ORDER BY cust.customer_id" +
				", cust.customer_id ASC" +
				", cust.customer_id DESC" +
				", cust.customer_id NULLS LAST" +
				", cust.customer_id ASC NULLS FIRST" +
				", cust.customer_id DESC NULLS LAST"
			return TT{DESCRIPTION, q, wantQuery, nil}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.DESCRIPTION, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			gotQuery, gotArgs := tt.q.ToSQL()
			is.Equal(tt.wantQuery, gotQuery)
			is.Equal(tt.wantArgs, gotArgs)
		})
	}
}

func TestSelectQuery_Limit_Offset(t *testing.T) {
	type TT struct {
		DESCRIPTION string
		q           SelectQuery
		wantQuery   string
		wantArgs    []interface{}
	}
	baseSelect := NewSelectQuery()
	tests := []TT{
		func() TT {
			DESCRIPTION := "Limit and offset"
			q := baseSelect.Limit(10).Offset(20)
			wantQuery := "LIMIT $1 OFFSET $2"
			return TT{DESCRIPTION, q, wantQuery, []interface{}{uint64(10), uint64(20)}}
		}(),
		func() TT {
			DESCRIPTION := "negative numbers get abs-ed"
			q := baseSelect.Limit(-22).Offset(-34)
			wantQuery := "LIMIT $1 OFFSET $2"
			return TT{DESCRIPTION, q, wantQuery, []interface{}{uint64(22), uint64(34)}}
		}(),
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.DESCRIPTION, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			gotQuery, gotArgs := tt.q.ToSQL()
			is.Equal(tt.wantQuery, gotQuery)
			is.Equal(tt.wantArgs, gotArgs)
		})
	}
}
