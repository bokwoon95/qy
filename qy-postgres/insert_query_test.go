package qy

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/bokwoon95/qy/qx"
	"github.com/bokwoon95/qy/tables-postgres"
	"github.com/matryer/is"
)

type TestUser struct {
	Valid    bool
	Uid      int64
	Name     string
	Email    string
	Password sql.NullString
}

func TestInsertTemp1(t *testing.T) {
	wantQuery, wantArgs := "", []interface{}{}

	A := qx.NewCTE("A", func() qx.Query {
		coun := tables.COUNTRY().As("coun")
		return Select(coun.COUNTRY_ID).From(coun).Where(coun.COUNTRY.EqString("Canada"))
	}())
	wantQuery += "WITH A AS (SELECT coun.country_id FROM country AS coun WHERE coun.country = $1)"
	wantArgs = append(wantArgs, "Canada")

	B := qx.NewCTE("B", func() qx.Query {
		city := tables.CITY()
		return Select(city.CITY_ID).From(city).Where(Predicatef("? IN (SELECT country_id FROM ?)", city.COUNTRY_ID, A))
	}())
	wantQuery += ", B AS (SELECT city.city_id FROM city WHERE city.country_id IN (SELECT country_id FROM A))"

	C := qx.NewCTE("C", func() qx.Query {
		addr := tables.ADDRESS().As("addr")
		return Select(addr.ADDRESS_ID).From(addr).Where(Predicatef("? IN (SELECT city_id FROM ?)", addr.CITY_ID, B))
	}())
	wantQuery += ", C AS (SELECT addr.address_id FROM address AS addr WHERE addr.city_id IN (SELECT city_id FROM B))"

	D := qx.NewCTE("D", func() qx.Query {
		cust := tables.CUSTOMER().As("cust")
		return Select(
			cust.CUSTOMER_ID,
			cust.STORE_ID,
			cust.FIRST_NAME,
			cust.LAST_NAME,
			cust.ADDRESS_ID,
		).From(cust).Where(Predicatef("? IN (SELECT address_id FROM ?)", cust.ADDRESS_ID, C))
	}())
	wantQuery += ", D AS (SELECT cust.customer_id, cust.store_id, cust.first_name, cust.last_name, cust.address_id" +
		" FROM customer AS cust WHERE cust.address_id IN (SELECT address_id FROM C))"

	cust := tables.CUSTOMER().As("cust")
	q := NewInsertQuery().With(A, B, C, D)
	q.Log = log.New(os.Stdout, "", 0)

	// v1
	v1 := q.InsertInto(cust).
		Columns(cust.CUSTOMER_ID, cust.STORE_ID, cust.FIRST_NAME, cust.LAST_NAME, cust.ADDRESS_ID).
		Select(
			Select(
				D.Get("customer_id"),
				D.Get("store_id"),
				D.Get("first_name"),
				D.Get("last_name"),
				D.Get("address_id"),
			).From(D),
		).
		OnConflict(cust.CUSTOMER_ID).
		DoUpdateSet(
			cust.CUSTOMER_ID.Set(Excluded(cust.CUSTOMER_ID)),
			cust.STORE_ID.Set(Excluded(cust.STORE_ID)),
			cust.FIRST_NAME.Set(Excluded(cust.FIRST_NAME)),
			cust.LAST_NAME.Set(Excluded(cust.LAST_NAME)),
			cust.ADDRESS_ID.Set(Excluded(cust.ADDRESS_ID)),
		).
		Returning(Fieldf("*"))
	v1WantQuery := wantQuery + " INSERT INTO customer AS cust (customer_id, store_id, first_name, last_name, address_id)" +
		" SELECT D.customer_id, D.store_id, D.first_name, D.last_name, D.address_id FROM D" +
		" ON CONFLICT (customer_id) DO UPDATE SET" +
		" customer_id = EXCLUDED.customer_id" +
		", store_id = EXCLUDED.store_id" +
		", first_name = EXCLUDED.first_name" +
		", last_name = EXCLUDED.last_name" +
		", address_id = EXCLUDED.address_id" +
		" RETURNING *"
	v1WantArgs := wantArgs

	// v2
	v2 := q.InsertInto(cust).
		Columns(cust.CUSTOMER_ID, cust.STORE_ID, cust.FIRST_NAME, cust.LAST_NAME, cust.ADDRESS_ID).
		Values(1, 1, "bob", "the builder", 1).
		OnConflict(cust.CUSTOMER_ID).
		DoNothing().
		Returning(Fieldf("*"))
	v2WantQuery := wantQuery + " INSERT INTO customer AS cust (customer_id, store_id, first_name, last_name, address_id)" +
		" VALUES ($2, $3, $4, $5, $6) ON CONFLICT (customer_id) DO NOTHING RETURNING *"
	v2WantArgs := wantArgs
	v2WantArgs = append(v2WantArgs, 1, 1, "bob", "the builder", 1)

	is := is.New(t)
	v1GotQuery, v1GotArgs := v1.ToSQL()
	is.Equal(v1WantQuery, v1GotQuery)
	is.Equal(v1WantArgs, v1GotArgs)
	v2GotQuery, v2GotArgs := v2.ToSQL()
	is.Equal(v2WantQuery, v2GotQuery)
	is.Equal(v2WantArgs, v2GotArgs)
}

// func TestInsertTemp2(t *testing.T) {
// 	wantQuery, wantArgs := "", []interface{}{}
//
// 	q := NewInsertQuery()
// 	q.Log = log.New(os.Stdout, "", 0)
// 	u := tables.USERS().As("u")
// 	q = q.InsertInto(u).
// 		InsertRow(
// 			u.UID.SetInt(1),
// 			u.DISPLAYNAME.SetString("aaa"),
// 			u.EMAIL.SetString("aaa@email.com"),
// 		).
// 		InsertRow(
// 			u.UID.SetInt(2),
// 			u.DISPLAYNAME.SetString("bbb"),
// 			u.EMAIL.SetString("bbb@email.com"),
// 		).
// 		InsertRow(
// 			u.UID.SetInt(3),
// 			u.DISPLAYNAME.SetString("ccc"),
// 			u.EMAIL.SetString("ccc@email.com"),
// 		).
// 		OnConflictOnConstraint("users_email_key").
// 		Where(u.EMAIL.NeString("")).
// 		DoNothing().
// 		Where(
// 			u.UID.GtInt(0),
// 			u.DISPLAYNAME.NeString(""),
// 		).
// 		Returning(u.UID, u.DISPLAYNAME, u.EMAIL)
// 	wantQuery += "INSERT INTO public.users AS u (uid, displayname, email)" +
// 		" VALUES ($1, $2, $3), ($4, $5, $6), ($7, $8, $9)" +
// 		" ON CONFLICT ON CONSTRAINT users_email_key DO NOTHING" +
// 		" RETURNING u.uid, u.displayname, u.email"
// 	wantArgs = append(wantArgs, 1, "aaa", "aaa@email.com", 2, "bbb", "bbb@email.com",
// 		3, "ccc", "ccc@email.com")
//
// 	is := is.New(t)
// 	gotQuery, gotArgs := q.ToSQL()
// 	is.Equal(wantQuery, gotQuery)
// 	is.Equal(wantArgs, gotArgs)
// }

// func TestInsertReal(t *testing.T) {
// 	is := is.New(t)
// 	db, err := sql.Open("txdb", qx.RandomString(8))
// 	is.NoErr(err)
// 	defer db.Close()
// 	q := NewInsertQuery()
// 	q.Log = log.New(os.Stdout, "[BIG BREIN] ", 0)
//
// 	u := tables.USERS()
// 	names := []string{"aaa", "bbb", "ccc"}
// 	emails := []string{"aaa@email.com", "bbb@email.com", "ccc@email.com"}
// 	insert := q.InsertInto(u).
// 		OnConflict(u.EMAIL).
// 		Where(Predicatef("(? + ? - ?) <> ?", 5, 99999, u.UID, 0)).
// 		DoUpdateSet(u.DISPLAYNAME.Set(Excluded(u.DISPLAYNAME))).
// 		Where(u.EMAIL.NeString(""))
// 	for i := range names {
// 		if i == 0 {
// 			insert = insert.Values(names[i], emails[i])
// 			continue
// 		}
// 		insert = insert.InsertRow(
// 			u.DISPLAYNAME.SetString(names[i]),
// 			u.EMAIL.SetString(emails[i]),
// 		)
// 	}
// 	var user TestUser
// 	var users []TestUser
// 	err = insert.Returningx(func(row Row) {
// 		user.Valid = row.IntValid(u.UID)
// 		user.Uid = row.Int64(u.UID)
// 		user.Name = row.String(u.DISPLAYNAME)
// 		user.Email = row.String(u.EMAIL)
// 	}, func() {
// 		users = append(users, user)
// 	}).Exec(db)
// 	is.NoErr(err)
// 	fmt.Println(users)
// 	is.Equal(len(names), len(users))
// 	for i := range users {
// 		is.True(users[i].Valid)
// 		is.True(users[i].Uid != 0)
// 		is.Equal(names[i], users[i].Name)
// 		is.Equal(emails[i], users[i].Email)
// 	}
//
// 	selectCount := func(db qx.Queryer, q InsertQuery) (count int, err error) {
// 		q.ReturningFields = []qx.Field{Fieldf("1")}
// 		stmt := qx.NewCTE("stmt", q)
// 		q2 := NewSelectQuery()
// 		q2.Log = q.Log
// 		err = q2.With(stmt).From(stmt).SelectRowx(func(row Row) {
// 			count = row.Int_(Fieldf("COUNT(*)"))
// 		}).Exec(db)
// 		return count, err
// 	}
// 	count, err := selectCount(db,
// 		q.InsertInto(u).
// 			Columns(u.UID, u.DISPLAYNAME, u.EMAIL).
// 			Values(users[0].Uid, users[0].Name, users[0].Email).
// 			Values(users[1].Uid, users[1].Name, users[1].Email).
// 			Values(users[2].Uid, users[2].Name, users[2].Email).
// 			OnConflict().DoNothing(),
// 	)
// 	is.Equal(0, count)
// 	is.NoErr(err)
// 	var randStrs []string
// 	for i := 0; i < 10; i++ {
// 		randStrs = append(randStrs, qx.RandomString(8))
// 	}
// 	q = q.InsertInto(u).Columns(u.UID, u.DISPLAYNAME, u.EMAIL).OnConflict().DoNothing()
// 	for i := range randStrs {
// 		q = q.Values(math.MaxInt32-i, randStrs[i], randStrs[i])
// 	}
// 	count, err = selectCount(db, q)
// 	is.NoErr(err)
// 	is.Equal(len(randStrs), count)
// }

// func TestInsertNull(t *testing.T) {
// 	is := is.New(t)
// 	db, err := sql.Open("txdb", qx.RandomString(8))
// 	is.NoErr(err)
// 	defer db.Close()
// 	const DEFAULT qx.FieldLiteral = "DEFAULT"
// 	q := NewInsertQuery()
// 	q.Log = log.New(os.Stdout, "", 0)
//
// 	u := tables.USERS().As("u")
// 	names := []string{"aaa", "bbb", "ccc"}
// 	emails := []string{"aaa@email.com", "bbb@email.com", "ccc@email.com"}
// 	insert := q.InsertInto(u).
// 		OnConflict(u.EMAIL).
// 		Where(Predicatef("(? + ? - ?) <> ?", 5, 99999, u.UID, 0)).
// 		Where(Predicatef("? = ?", u.DISPLAYNAME, nil)).
// 		DoUpdateSet(u.DISPLAYNAME.Set(Excluded(u.DISPLAYNAME))).
// 		Where(u.EMAIL.NeString(""))
// 	for i := range names {
// 		if i == 0 {
// 			insert = insert.Values(names[i], emails[i], "yohoho")
// 			continue
// 		}
// 		if i == 1 {
// 			insert = insert.Values(names[i], emails[i], DEFAULT)
// 			continue
// 		}
// 		insert = insert.InsertRow(
// 			u.DISPLAYNAME.SetString(names[i]),
// 			u.EMAIL.SetString(emails[i]),
// 			u.PASSWORD.Set(nil),
// 		)
// 	}
// 	var user TestUser
// 	var users []TestUser
// 	err = insert.Returningx(func(row Row) {
// 		uid := row.NullInt64(u.UID)
// 		user.Valid = uid.Valid
// 		user.Uid = uid.Int64
// 		user.Name = row.String(u.DISPLAYNAME)
// 		user.Email = row.String(u.EMAIL)
// 		user.Password = row.NullString(u.PASSWORD)
// 	}, func() {
// 		users = append(users, user)
// 	}).Exec(db)
// 	fmt.Println(users)
// }
