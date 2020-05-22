package qy

// import (
// 	"log"
// 	"os"
// 	"testing"
//
// 	"github.com/bokwoon95/qy/qx"
// 	"github.com/bokwoon95/qy/tables-postgres"
// 	_ "github.com/lib/pq"
// 	"github.com/matryer/is"
// )
//
// func TestSelectQuery_Get(t *testing.T) {
// 	type TT struct {
// 		DESCRIPTION string
// 		f           qx.Field
// 		wantQuery   string
// 		wantArgs    []interface{}
// 	}
// 	baseSelect := NewSelectQuery()
// 	baseSelect.Log = log.New(os.Stdout, "", 0)
// 	tests := []TT{
// 		{
// 			"SelectQuery Get (explicit alias)",
// 			baseSelect.From(tables.EMPLOYEES()).As("123").Get("abcde"),
// 			"123.abcde",
// 			nil,
// 		},
// 		func() TT {
// 			DESCRIPTION := "SelectQuery Get (implicit alias)"
// 			q := baseSelect.From(tables.EMPLOYEES())
// 			alias := q.GetAlias()
// 			if alias == "" {
// 				t.Fatalf("an alias should have automatically been assigned to SelectQuery, found empty alias")
// 			}
// 			return TT{DESCRIPTION, q.Get("abcde"), alias + ".abcde", nil}
// 		}(),
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.DESCRIPTION, func(t *testing.T) {
// 			t.Parallel()
// 			is := is.New(t)
// 			gotQuery, gotArgs := tt.f.ToSQL(nil)
// 			is.Equal(tt.wantQuery, gotQuery)
// 			is.Equal(tt.wantArgs, gotArgs)
// 		})
// 	}
// }
//
// func TestSelectQuery_With(t *testing.T) {
// 	q := NewSelectQuery()
// 	wantQuery, wantArgs := "", []interface{}{}
//
// 	cte1 := qx.NewCTE("cte1", func() qx.Query {
// 		empl, orde, ship := tables.EMPLOYEES().As("empl"), tables.ORDERS().As("orde"), tables.SHIPPERS().As("ship")
// 		return Select(empl.FIRST_NAME, empl.LAST_NAME, orde.SHIP_NAME, orde.SHIP_ADDRESS, ship.COMPANY_NAME).
// 			From(empl).Join(orde, orde.EMPLOYEE_ID.Eq(empl.EMPLOYEE_ID)).
// 			LeftJoin(ship, ship.SHIPPER_ID.Eq(orde.SHIP_VIA)).
// 			Where(empl.TITLE_OF_COURTESY.ILikeString("%mr%"))
// 	}())
// 	wantQuery += "WITH cte1 AS" +
// 		" (SELECT empl.first_name, empl.last_name, orde.ship_name, orde.ship_address, ship.company_name" +
// 		" FROM employees AS empl JOIN orders AS orde ON orde.employee_id = empl.employee_id" +
// 		" LEFT JOIN shippers AS ship ON ship.shipper_id = orde.ship_via" +
// 		" WHERE empl.title_of_courtesy ILIKE $1)"
// 	wantArgs = append(wantArgs, "%mr%")
//
// 	cte2 := qx.NewCTE("cte2", func() qx.Query {
// 		prod, cate, supp := tables.PRODUCTS().As("prod"), tables.CATEGORIES().As("cate"), tables.SUPPLIERS().As("supp")
// 		return Select(prod.PRODUCT_NAME, prod.SUPPLIER_ID, supp.COMPANY_NAME, cate.CATEGORY_NAME).
// 			From(prod).Join(supp, supp.SUPPLIER_ID.Eq(prod.SUPPLIER_ID)).
// 			LeftJoin(cate, cate.CATEGORY_ID.Eq(prod.CATEGORY_ID)).
// 			Where(Predicatef("? ILIKE ANY (ARRAY[?])", supp.CITY, []string{"london", "singapore", "tokyo", "zurich"}))
// 	}())
// 	wantQuery += ", cte2 AS" +
// 		" (SELECT prod.product_name, prod.supplier_id, supplier.company_name, cate.category_name" +
// 		" FROM products AS prod JOIN suppliers AS supp ON supp.supplier_id = prod.supplier_id" +
// 		" LEFT JOIN categories AS cate ON cate.category_id = prod.category_id" +
// 		" WHERE supp.city ILIKE ANY (ARRAY[?, ?, ?, ?]))"
// 	wantArgs = append(wantArgs, "london", "singapore", "tokyo", "zurich")
//
// 	cte3 := qx.NewCTE("cte3", func() qx.Query {
// 		return Select(u.UID, ur.URID, u.DISPLAYNAME, u.EMAIL).
// 			From(u).Join(ur, ur.UID.Eq(u.UID)).
// 			Where(ur.ROLE.EqString("adviser"))
// 	}())
// 	wantQuery += ", advisers AS" +
// 		" (SELECT u.uid, ur.urid, u.displayname, u.email" +
// 		" FROM public.users AS u JOIN public.user_roles AS ur ON ur.uid = u.uid" +
// 		" WHERE ur.role = $3)"
// 	wantArgs = append(wantArgs, "adviser")
//
// 	stu1, stu2 := students.As("stu1"), students.As("stu2")
// 	q = q.With(qx.CTE{}, cte1, students, advisers).
// 		Select(cte1.Get("displayname"), stu1.Get("uid"), stu2.Get("uid"), advisers.Get("email")).
// 		From(cte1).
// 		CrossJoin(stu1).
// 		Join(stu2, stu2.Get("uid").Eq(stu1.Get("uid"))).
// 		CrossJoin(advisers)
// 	wantQuery += " SELECT applicants.displayname, stu1.uid, stu2.uid, advisers.email" +
// 		" FROM applicants" +
// 		" CROSS JOIN students AS stu1" +
// 		" JOIN students AS stu2 ON stu2.uid = stu1.uid" +
// 		" CROSS JOIN advisers"
// 	wantArgs = append(wantArgs)
//
// 	is := is.New(t)
// 	gotQuery, gotArgs := q.ToSQL()
// 	is.Equal(wantQuery, gotQuery)
// 	is.Equal(wantArgs, gotArgs)
// }
//
// func TestSelectQuery_Select(t *testing.T) {
// 	type TT struct {
// 		DESCRIPTION string
// 		q           SelectQuery
// 		wantQuery   string
// 		wantArgs    []interface{}
// 	}
// 	baseSelect := NewSelectQuery()
// 	tests := []TT{
// 		func() TT {
// 			DESCRIPTION := "basic select"
// 			u := tables.USERS().As("u")
// 			q := baseSelect.Select(u.UID, u.DISPLAYNAME, u.EMAIL)
// 			wantQuery := "SELECT u.uid, u.displayname, u.email"
// 			return TT{DESCRIPTION, q, wantQuery, nil}
// 		}(),
// 		func() TT {
// 			DESCRIPTION := "multiple select calls also work"
// 			u := tables.USERS().As("u")
// 			q := baseSelect.Select(u.UID).Select(u.DISPLAYNAME).Select(u.EMAIL)
// 			wantQuery := "SELECT u.uid, u.displayname, u.email"
// 			return TT{DESCRIPTION, q, wantQuery, nil}
// 		}(),
// 		func() TT {
// 			DESCRIPTION := "column aliases work"
// 			ur := tables.USER_ROLES().As("ur")
// 			q := baseSelect.Select(ur.URID, ur.COHORT, ur.ROLE.As("user_role"), ur.CREATED_AT.As("date_created"))
// 			wantQuery := "SELECT ur.urid, ur.cohort, ur.role AS user_role, ur.created_at AS date_created"
// 			return TT{DESCRIPTION, q, wantQuery, nil}
// 		}(),
// 		func() TT {
// 			DESCRIPTION := "select distinct"
// 			ur := tables.USER_ROLES().As("ur")
// 			q := baseSelect.SelectDistinct(ur.URID, ur.COHORT)
// 			wantQuery := "SELECT DISTINCT ur.urid, ur.cohort"
// 			return TT{DESCRIPTION, q, wantQuery, nil}
// 		}(),
// 		func() TT {
// 			DESCRIPTION := "select distinct on"
// 			ur := tables.USER_ROLES().As("ur")
// 			q := baseSelect.SelectDistinctOn(ur.URID, ur.COHORT)(ur.URID, ur.COHORT, ur.UPDATED_AT)
// 			wantQuery := "SELECT DISTINCT ON (ur.urid, ur.cohort) ur.urid, ur.cohort, ur.updated_at"
// 			return TT{DESCRIPTION, q, wantQuery, nil}
// 		}(),
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.DESCRIPTION, func(t *testing.T) {
// 			t.Parallel()
// 			is := is.New(t)
// 			gotQuery, gotArgs := tt.q.ToSQL()
// 			is.Equal(tt.wantQuery, gotQuery)
// 			is.Equal(tt.wantArgs, gotArgs)
// 		})
// 	}
// }
//
// func TestSelectQuery_From(t *testing.T) {
// 	type TT struct {
// 		DESCRIPTION string
// 		q           SelectQuery
// 		wantQuery   string
// 		wantArgs    []interface{}
// 	}
// 	baseSelect := NewSelectQuery()
// 	tests := []TT{
// 		func() TT {
// 			DESCRIPTION := "fields from unaliased tables are qualified by table name not table alias"
// 			u := tables.USERS()
// 			q := baseSelect.From(u).Select(u.UID)
// 			wantQuery := "SELECT users.uid FROM public.users"
// 			return TT{DESCRIPTION, q, wantQuery, nil}
// 		}(),
// 		func() TT {
// 			DESCRIPTION := "but setting an alias for the table manually also works"
// 			u := tables.USERS().As("u")
// 			q := baseSelect.From(u)
// 			wantQuery := "FROM public.users AS u"
// 			return TT{DESCRIPTION, q, wantQuery, nil}
// 		}(),
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.DESCRIPTION, func(t *testing.T) {
// 			t.Parallel()
// 			is := is.New(t)
// 			gotQuery, gotArgs := tt.q.ToSQL()
// 			is.Equal(tt.wantQuery, gotQuery)
// 			is.Equal(tt.wantArgs, gotArgs)
// 		})
// 	}
// }
//
// func TestSelectQuery_Joins(t *testing.T) {
// 	type TT struct {
// 		DESCRIPTION string
// 		q           SelectQuery
// 		wantQuery   string
// 		wantArgs    []interface{}
// 	}
// 	baseSelect := NewSelectQuery()
// 	baseSelect.Log = log.New(os.Stdout, "", 0)
// 	tests := []TT{
// 		func() TT {
// 			DESCRIPTION := "all joinGroups"
// 			u, ur := tables.USERS().As("u"), tables.USER_ROLES().As("ur")
// 			q := baseSelect.From(u).
// 				Join(ur, ur.UID.Eq(u.UID)).
// 				LeftJoin(ur, ur.UID.Eq(u.UID)).
// 				RightJoin(ur, ur.UID.Eq(u.UID)).
// 				CrossJoin(ur)
// 			// You can't join tables with the same alias in SQL, this is just an example
// 			wantQuery := "FROM public.users AS u" +
// 				" JOIN public.user_roles AS ur ON ur.uid = u.uid" +
// 				" LEFT JOIN public.user_roles AS ur ON ur.uid = u.uid" +
// 				" RIGHT JOIN public.user_roles AS ur ON ur.uid = u.uid" +
// 				" CROSS JOIN public.user_roles AS ur"
// 			return TT{DESCRIPTION, q, wantQuery, nil}
// 		}(),
// 		func() TT {
// 			DESCRIPTION := "joining a subquery with explicit alias"
// 			ur, c := tables.USER_ROLES().As("ur"), tables.COHORT_ENUM().As("c")
// 			latest_cohorts := Select(c.COHORT).From(c).Where(Predicatef("TRIM(?)::INT > ?", c.COHORT, qx.Int(2020))).As("latest_cohorts")
// 			q := baseSelect.From(ur).Join(latest_cohorts, latest_cohorts.Get("cohort").Eq(ur.COHORT))
// 			wantQuery := "FROM public.user_roles AS ur" +
// 				" JOIN (SELECT c.cohort FROM public.cohort_enum AS c WHERE TRIM(c.cohort)::INT > $1) AS latest_cohorts" +
// 				" ON latest_cohorts.cohort = ur.cohort"
// 			return TT{DESCRIPTION, q, wantQuery, []interface{}{2020}}
// 		}(),
// 		func() TT {
// 			DESCRIPTION := "joining a subquery with implicit alias"
// 			ur, c := tables.USER_ROLES().As("ur"), tables.COHORT_ENUM().As("c")
// 			latest_cohorts := Select(c.COHORT).From(c).Where(Predicatef("TRIM(?)::INT > ?", c.COHORT, qx.Int(2020)))
// 			alias := latest_cohorts.GetAlias()
// 			is := is.New(t)
// 			is.True(alias != "")
// 			q := baseSelect.From(ur).Join(latest_cohorts, latest_cohorts.Get("cohort").Eq(ur.COHORT))
// 			wantQuery := "FROM public.user_roles AS ur" +
// 				" JOIN (SELECT c.cohort FROM public.cohort_enum AS c WHERE TRIM(c.cohort)::INT > $1) AS " + alias +
// 				" ON " + alias + ".cohort = ur.cohort"
// 			return TT{DESCRIPTION, q, wantQuery, []interface{}{2020}}
// 		}(),
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.DESCRIPTION, func(t *testing.T) {
// 			t.Parallel()
// 			is := is.New(t)
// 			gotQuery, gotArgs := tt.q.ToSQL()
// 			is.Equal(tt.wantQuery, gotQuery)
// 			is.Equal(tt.wantArgs, gotArgs)
// 		})
// 	}
// }
//
// func TestSelectQuery_Where(t *testing.T) {
// 	type TT struct {
// 		DESCRIPTION string
// 		q           SelectQuery
// 		wantQuery   string
// 		wantArgs    []interface{}
// 	}
// 	baseSelect := NewSelectQuery()
// 	tests := []TT{
// 		func() TT {
// 			DESCRIPTION := "basic where (implicit and)"
// 			u := tables.USERS().As("u")
// 			q := baseSelect.Where(
// 				u.UID.EqInt(22),
// 				u.DISPLAYNAME.ILikeString("%bob%"),
// 				u.EMAIL.IsNotNull(),
// 			)
// 			wantQuery := "WHERE u.uid = $1 AND u.displayname ILIKE $2 AND u.email IS NOT NULL"
// 			return TT{DESCRIPTION, q, wantQuery, []interface{}{22, "%bob%"}}
// 		}(),
// 		func() TT {
// 			DESCRIPTION := "basic where (explicit and)"
// 			u := tables.USERS().As("u")
// 			q := baseSelect.Where(
// 				qx.And(
// 					u.UID.EqInt(22),
// 					u.DISPLAYNAME.ILikeString("%bob%"),
// 					u.EMAIL.IsNotNull(),
// 				),
// 			)
// 			wantQuery := "WHERE u.uid = $1 AND u.displayname ILIKE $2 AND u.email IS NOT NULL"
// 			return TT{DESCRIPTION, q, wantQuery, []interface{}{22, "%bob%"}}
// 		}(),
// 		func() TT {
// 			DESCRIPTION := "basic where (explicit or)"
// 			u := tables.USERS().As("u")
// 			q := baseSelect.Where(
// 				qx.Or(
// 					u.UID.EqInt(22),
// 					u.DISPLAYNAME.ILikeString("%bob%"),
// 					u.EMAIL.IsNotNull(),
// 				),
// 			)
// 			wantQuery := "WHERE u.uid = $1 OR u.displayname ILIKE $2 OR u.email IS NOT NULL"
// 			return TT{DESCRIPTION, q, wantQuery, []interface{}{22, "%bob%"}}
// 		}(),
// 		func() TT {
// 			DESCRIPTION := "complex predicate"
// 			u1, u2 := tables.USERS().As("u1"), tables.USERS().As("u2")
// 			q := baseSelect.Where(
// 				u1.UID.EqInt(69),
// 				u1.DISPLAYNAME.LikeString("%bob%"),
// 				u1.EMAIL.IsNull(),
// 				qx.Or(
// 					u2.UID.EqInt(420),
// 					u2.DISPLAYNAME.ILikeString("%virgil%"),
// 					u2.EMAIL.IsNotNull(),
// 				),
// 			)
// 			wantQuery := "WHERE u1.uid = $1 AND u1.displayname LIKE $2 AND u1.email IS NULL" +
// 				" AND (u2.uid = $3 OR u2.displayname ILIKE $4 OR u2.email IS NOT NULL)"
// 			return TT{DESCRIPTION, q, wantQuery, []interface{}{69, "%bob%", 420, "%virgil%"}}
// 		}(),
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.DESCRIPTION, func(t *testing.T) {
// 			t.Parallel()
// 			is := is.New(t)
// 			gotQuery, gotArgs := tt.q.ToSQL()
// 			is.Equal(tt.wantQuery, gotQuery)
// 			is.Equal(tt.wantArgs, gotArgs)
// 		})
// 	}
// }
//
// func TestSelectQuery_GroupBy(t *testing.T) {
// 	type TT struct {
// 		DESCRIPTION string
// 		q           SelectQuery
// 		wantQuery   string
// 		wantArgs    []interface{}
// 	}
// 	baseSelect := NewSelectQuery()
// 	tests := []TT{
// 		func() TT {
// 			DESCRIPTION := "basic group by"
// 			u := tables.USERS().As("u")
// 			q := baseSelect.GroupBy(u.UID, u.DISPLAYNAME, u.EMAIL)
// 			wantQuery := "GROUP BY u.uid, u.displayname, u.email"
// 			return TT{DESCRIPTION, q, wantQuery, nil}
// 		}(),
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.DESCRIPTION, func(t *testing.T) {
// 			t.Parallel()
// 			is := is.New(t)
// 			gotQuery, gotArgs := tt.q.ToSQL()
// 			is.Equal(tt.wantQuery, gotQuery)
// 			is.Equal(tt.wantArgs, gotArgs)
// 		})
// 	}
// }
//
// func TestSelectQuery_Having(t *testing.T) {
// 	type TT struct {
// 		DESCRIPTION string
// 		q           SelectQuery
// 		wantQuery   string
// 		wantArgs    []interface{}
// 	}
// 	baseSelect := NewSelectQuery()
// 	tests := []TT{
// 		func() TT {
// 			DESCRIPTION := "basic having (implicit and)"
// 			u := tables.USERS().As("u")
// 			q := baseSelect.Having(
// 				u.UID.EqInt(22),
// 				u.DISPLAYNAME.ILikeString("%bob%"),
// 				u.EMAIL.IsNotNull(),
// 			)
// 			wantQuery := "HAVING u.uid = $1 AND u.displayname ILIKE $2 AND u.email IS NOT NULL"
// 			return TT{DESCRIPTION, q, wantQuery, []interface{}{22, "%bob%"}}
// 		}(),
// 		func() TT {
// 			DESCRIPTION := "basic having (explicit and)"
// 			u := tables.USERS().As("u")
// 			q := baseSelect.Having(
// 				qx.And(
// 					u.UID.EqInt(22),
// 					u.DISPLAYNAME.ILikeString("%bob%"),
// 					u.EMAIL.IsNotNull(),
// 				),
// 			)
// 			wantQuery := "HAVING u.uid = $1 AND u.displayname ILIKE $2 AND u.email IS NOT NULL"
// 			return TT{DESCRIPTION, q, wantQuery, []interface{}{22, "%bob%"}}
// 		}(),
// 		func() TT {
// 			DESCRIPTION := "basic having (explicit or)"
// 			u := tables.USERS().As("u")
// 			q := baseSelect.Having(
// 				qx.Or(
// 					u.UID.EqInt(22),
// 					u.DISPLAYNAME.ILikeString("%bob%"),
// 					u.EMAIL.IsNotNull(),
// 				),
// 			)
// 			wantQuery := "HAVING u.uid = $1 OR u.displayname ILIKE $2 OR u.email IS NOT NULL"
// 			return TT{DESCRIPTION, q, wantQuery, []interface{}{22, "%bob%"}}
// 		}(),
// 		func() TT {
// 			DESCRIPTION := "complex predicate"
// 			u1, u2 := tables.USERS().As("u1"), tables.USERS().As("u2")
// 			q := baseSelect.Having(
// 				u1.UID.EqInt(69),
// 				u1.DISPLAYNAME.LikeString("%bob%"),
// 				u1.EMAIL.IsNull(),
// 				qx.Or(
// 					u2.UID.EqInt(420),
// 					u2.DISPLAYNAME.ILikeString("%virgil%"),
// 					u2.EMAIL.IsNotNull(),
// 				),
// 			)
// 			wantQuery := "HAVING u1.uid = $1 AND u1.displayname LIKE $2 AND u1.email IS NULL" +
// 				" AND (u2.uid = $3 OR u2.displayname ILIKE $4 OR u2.email IS NOT NULL)"
// 			return TT{DESCRIPTION, q, wantQuery, []interface{}{69, "%bob%", 420, "%virgil%"}}
// 		}(),
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.DESCRIPTION, func(t *testing.T) {
// 			t.Parallel()
// 			is := is.New(t)
// 			gotQuery, gotArgs := tt.q.ToSQL()
// 			is.Equal(tt.wantQuery, gotQuery)
// 			is.Equal(tt.wantArgs, gotArgs)
// 		})
// 	}
// }
//
// func TestSelectQuery_OrderBy(t *testing.T) {
// 	type TT struct {
// 		DESCRIPTION string
// 		q           SelectQuery
// 		wantQuery   string
// 		wantArgs    []interface{}
// 	}
// 	baseSelect := NewSelectQuery()
// 	tests := []TT{
// 		func() TT {
// 			DESCRIPTION := "basic order by"
// 			u := tables.USERS().As("u")
// 			q := baseSelect.Select().OrderBy(u.UID, u.DISPLAYNAME, u.EMAIL)
// 			wantQuery := "ORDER BY u.uid, u.displayname, u.email"
// 			return TT{DESCRIPTION, q, wantQuery, nil}
// 		}(),
// 		func() TT {
// 			DESCRIPTION := "Asc, DESCRIPTION, NullsFirst and NullsLast"
// 			u := tables.USERS().As("u")
// 			q := baseSelect.OrderBy(
// 				u.UID,
// 				u.UID.Asc(),
// 				u.UID.Desc(),
// 				u.UID.NullsLast(),
// 				u.UID.Asc().NullsFirst(),
// 				u.UID.Desc().NullsLast(),
// 			)
// 			wantQuery := "ORDER BY u.uid, u.uid ASC, u.uid DESC, u.uid NULLS LAST, u.uid ASC NULLS FIRST, u.uid DESC NULLS LAST"
// 			return TT{DESCRIPTION, q, wantQuery, nil}
// 		}(),
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.DESCRIPTION, func(t *testing.T) {
// 			t.Parallel()
// 			is := is.New(t)
// 			gotQuery, gotArgs := tt.q.ToSQL()
// 			is.Equal(tt.wantQuery, gotQuery)
// 			is.Equal(tt.wantArgs, gotArgs)
// 		})
// 	}
// }
//
// func TestSelectQuery_Limit_Offset(t *testing.T) {
// 	type TT struct {
// 		DESCRIPTION string
// 		q           SelectQuery
// 		wantQuery   string
// 		wantArgs    []interface{}
// 	}
// 	baseSelect := NewSelectQuery()
// 	tests := []TT{
// 		func() TT {
// 			DESCRIPTION := "Limit and offset"
// 			q := baseSelect.Limit(10).Offset(20)
// 			wantQuery := "LIMIT $1 OFFSET $2"
// 			return TT{DESCRIPTION, q, wantQuery, []interface{}{uint64(10), uint64(20)}}
// 		}(),
// 		func() TT {
// 			DESCRIPTION := "negative numbers get abs-ed"
// 			q := baseSelect.Limit(-22).Offset(-34)
// 			wantQuery := "LIMIT $1 OFFSET $2"
// 			return TT{DESCRIPTION, q, wantQuery, []interface{}{uint64(22), uint64(34)}}
// 		}(),
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.DESCRIPTION, func(t *testing.T) {
// 			t.Parallel()
// 			is := is.New(t)
// 			gotQuery, gotArgs := tt.q.ToSQL()
// 			is.Equal(tt.wantQuery, gotQuery)
// 			is.Equal(tt.wantArgs, gotArgs)
// 		})
// 	}
// }
