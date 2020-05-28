package qy

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/bokwoon95/qy/qx"
)

type DeleteQuery struct {
	Nested bool
	Alias  string
	// WITH
	CTEs qx.CTEs
	// DELETE FROM
	FromTable qx.BaseTable
	// USING
	UsingTable qx.Table
	JoinGroups qx.JoinTables
	// WHERE
	WherePredicates qx.VariadicPredicate
	// RETURNING
	ReturningFields qx.Fields
	// DB
	DB          qx.DB
	Mapper      func(Row)
	Accumulator func()
	// Logging
	Log     qx.Logger
	LogSkip int
}

func (q DeleteQuery) ToSQL() (string, []interface{}) {
	var buf = &strings.Builder{}
	var args []interface{}
	var excludeTableQualifiers []string
	// WITH
	q.CTEs.WriteSQL(buf, &args)
	{ // DELETE FROM
		deleteQuery, deleteArgs := "", []interface{}{}
		if q.FromTable != nil {
			deleteQuery, deleteArgs = q.FromTable.ToSQL()
			if q.FromTable.GetAlias() != "" {
				excludeTableQualifiers = append(excludeTableQualifiers, q.FromTable.GetAlias())
			} else if q.FromTable.GetName() != "" {
				excludeTableQualifiers = append(excludeTableQualifiers, q.FromTable.GetName())
			}
		}
		if deleteQuery != "" {
			if buf.Len() > 0 {
				buf.WriteString(" ")
			}
			if q.FromTable.GetAlias() != "" {
				buf.WriteString("DELETE FROM " + deleteQuery + " AS " + q.FromTable.GetAlias())
			} else {
				buf.WriteString("DELETE FROM " + deleteQuery)
			}
			args = append(args, deleteArgs...)
		}
	}
	{ // USING
		usingQuery, usingArgs := "", []interface{}{}
		if q.UsingTable != nil {
			usingQuery, usingArgs = q.UsingTable.ToSQL()
		}
		if usingQuery != "" {
			if buf.Len() > 0 {
				buf.WriteString(" ")
			}
			if q.UsingTable.GetAlias() != "" {
				buf.WriteString("USING " + usingQuery + " AS " + q.UsingTable.GetAlias())
			} else {
				buf.WriteString("USING " + usingQuery)
			}
			args = append(args, usingArgs...)
		}
	}
	// JOIN
	q.JoinGroups.WriteSQL(buf, &args)
	// WHERE
	q.WherePredicates.Toplevel = true
	q.WherePredicates.WriteSQL(buf, &args, "WHERE ", "", nil)
	// RETURNING
	q.ReturningFields.WriteSQLWithAlias(buf, &args, "RETURNING ", "", nil)
	query := buf.String()
	if !q.Nested {
		query = qx.MySQLToPostgresPlaceholders(query)
		switch q.Log.(type) {
		case nil:
			// do nothing
		case *log.Logger:
			q.Log.Output(q.LogSkip+2, qx.PostgresInterpolateSQL(query, args...))
		default:
			q.Log.Output(q.LogSkip+1, qx.PostgresInterpolateSQL(query, args...))
		}
	}
	return query, args
}

func (q DeleteQuery) GetAlias() string {
	return q.Alias
}

func (q DeleteQuery) GetName() string {
	return ""
}

func (q DeleteQuery) NestThis() qx.Query {
	q.Nested = true
	return q
}

func (q DeleteQuery) As(alias string) DeleteQuery {
	q.Alias = alias
	return q
}

func DeleteFrom(table qx.BaseTable) DeleteQuery {
	return DeleteQuery{
		FromTable: table,
		Alias:     qx.RandomString(8),
	}
}

func (q DeleteQuery) With(cteList ...qx.CTE) DeleteQuery {
	q.CTEs = append(q.CTEs, cteList...)
	return q
}

func (q DeleteQuery) DeleteFrom(tbl qx.BaseTable) DeleteQuery {
	q.FromTable = tbl
	return q
}

func (q DeleteQuery) Using(tbl qx.Table) DeleteQuery {
	q.UsingTable = tbl
	return q
}

func (q DeleteQuery) Join(tbl qx.Table, predicate qx.Predicate, predicates ...qx.Predicate) DeleteQuery {
	predicates = append([]qx.Predicate{predicate}, predicates...)
	q.JoinGroups = append(q.JoinGroups, qx.JoinTable{
		JoinType:     qx.JoinTypeDefault,
		Table:        tbl,
		OnPredicates: qx.VariadicPredicate{Predicates: predicates},
	})
	return q
}

func (q DeleteQuery) LeftJoin(tbl qx.Table, predicate qx.Predicate, predicates ...qx.Predicate) DeleteQuery {
	predicates = append([]qx.Predicate{predicate}, predicates...)
	q.JoinGroups = append(q.JoinGroups, qx.JoinTable{
		JoinType:     qx.JoinTypeLeft,
		Table:        tbl,
		OnPredicates: qx.VariadicPredicate{Predicates: predicates},
	})
	return q
}

func (q DeleteQuery) RightJoin(tbl qx.Table, predicate qx.Predicate, predicates ...qx.Predicate) DeleteQuery {
	predicates = append([]qx.Predicate{predicate}, predicates...)
	q.JoinGroups = append(q.JoinGroups, qx.JoinTable{
		JoinType:     qx.JoinTypeRight,
		Table:        tbl,
		OnPredicates: qx.VariadicPredicate{Predicates: predicates},
	})
	return q
}

func (q DeleteQuery) FullJoin(tbl qx.Table, predicate qx.Predicate, predicates ...qx.Predicate) DeleteQuery {
	predicates = append([]qx.Predicate{predicate}, predicates...)
	q.JoinGroups = append(q.JoinGroups, qx.JoinTable{
		JoinType:     qx.JoinTypeFull,
		Table:        tbl,
		OnPredicates: qx.VariadicPredicate{Predicates: predicates},
	})
	return q
}

func (q DeleteQuery) CrossJoin(tbl qx.Table) DeleteQuery {
	q.JoinGroups = append(q.JoinGroups, qx.JoinTable{
		JoinType: qx.JoinTypeCross,
		Table:    tbl,
	})
	return q
}

func (q DeleteQuery) Returning(fields ...qx.Field) DeleteQuery {
	q.ReturningFields = append(q.ReturningFields, fields...)
	return q
}

func (q DeleteQuery) Returningx(mapper func(Row), accumulator func()) DeleteQuery {
	q.Mapper = mapper
	q.Accumulator = accumulator
	return q
}

func (q DeleteQuery) ReturningRowx(mapper func(Row)) DeleteQuery {
	q.Mapper = mapper
	return q
}

func (q DeleteQuery) Fetch(db qx.DB) (err error) {
	q.LogSkip += 1
	return q.FetchContext(nil, db)
}

func (q DeleteQuery) FetchContext(ctx context.Context, db qx.DB) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = v
			case string:
				err = errors.New(v)
			}
		}
	}()
	if db == nil {
		return errors.New("DB cannot be nil")
	}
	r := &QyRow{QxRow: &qx.QxRow{}}
	if q.Mapper != nil {
		q.Mapper(r) // call the mapper once on the *Row to get all the selected that the user is interested in
	}
	q.ReturningFields = r.QxRow.Fields // then, transfer the selected collected by *Row to the DeleteQuery
	q.LogSkip += 1
	query, args := q.ToSQL()
	if ctx == nil {
		r.QxRow.Rows, err = db.Query(query, args...)
	} else {
		r.QxRow.Rows, err = db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return err
	}
	defer r.QxRow.Rows.Close()
	var rowcount int
	if len(r.QxRow.Dest) == 0 {
		// If there's nothing to scan into, return early
		return nil
	}
	for r.QxRow.Rows.Next() {
		rowcount++
		err = r.QxRow.Rows.Scan(r.QxRow.Dest...)
		if err != nil {
			buf := &strings.Builder{}
			for i := range r.QxRow.Dest {
				query, args := r.QxRow.Fields[i].ToSQLExclude(nil)
				buf.WriteString("\n" +
					strconv.Itoa(i) + ") " +
					qx.MySQLInterpolateSQL(query, args...) + " => " +
					reflect.TypeOf(r.QxRow.Dest[i]).String())
			}
			return fmt.Errorf("Please check if your mapper function is correct:%s\n%w", buf.String(), err)
		}
		r.QxRow.Index = 0 // index must always be reset back to 0 before mapper is called
		q.Mapper(r)
		if q.Accumulator == nil {
			break
		}
		q.Accumulator()
	}
	if rowcount == 0 && q.Accumulator == nil {
		return sql.ErrNoRows
	}
	return r.QxRow.Rows.Err()
}

func (q DeleteQuery) Exec(db qx.DB) (sql.Result, error) {
	q.LogSkip += 1
	return q.ExecContext(nil, db)
}

func (q DeleteQuery) ExecContext(ctx context.Context, db qx.DB) (sql.Result, error) {
	var res sql.Result
	var err error
	if db == nil {
		return res, errors.New("DB cannot be nil")
	}
	q.LogSkip += 1
	query, args := q.ToSQL()
	if ctx == nil {
		res, err = db.Exec(query, args...)
	} else {
		res, err = db.ExecContext(ctx, query, args...)
	}
	return res, err
}
