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

type UpdateQuery struct {
	Nested bool
	Alias  string
	// WITH
	CTEs qx.CTEs
	// UPDATE
	UpdateTable qx.BaseTable
	// SET
	SetFields qx.FieldValueSets
	// FROM
	FromTable  qx.Table
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

func (q UpdateQuery) ToSQL() (string, []interface{}) {
	var buf = &strings.Builder{}
	var args []interface{}
	var excludeTableQualifiers []string
	// WITH
	q.CTEs.WriteSQL(buf, &args)
	{ // UPDATE
		updateQuery, updateArgs := "", []interface{}{}
		if q.UpdateTable != nil {
			updateQuery, updateArgs = q.UpdateTable.ToSQL()
			if q.UpdateTable.GetAlias() != "" {
				excludeTableQualifiers = append(excludeTableQualifiers, q.UpdateTable.GetAlias())
			} else if q.UpdateTable.GetName() != "" {
				excludeTableQualifiers = append(excludeTableQualifiers, q.UpdateTable.GetName())
			}
		}
		if updateQuery != "" {
			if buf.Len() > 0 {
				buf.WriteString(" ")
			}
			if q.UpdateTable.GetAlias() != "" {
				buf.WriteString("UPDATE " + updateQuery + " AS " + q.UpdateTable.GetAlias())
			} else {
				buf.WriteString("UPDATE " + updateQuery)
			}
			args = append(args, updateArgs...)
		}
	}
	// SET
	q.SetFields.WriteSQL(buf, &args, "SET ", "", excludeTableQualifiers)
	{ // FROM
		fromQuery, fromArgs := "", []interface{}{}
		if q.FromTable != nil {
			fromQuery, fromArgs = q.FromTable.ToSQL()
		}
		if fromQuery != "" {
			if buf.Len() > 0 {
				buf.WriteString(" ")
			}
			if q.FromTable.GetAlias() != "" {
				buf.WriteString("FROM " + fromQuery + " AS " + q.FromTable.GetAlias())
			} else {
				buf.WriteString("FROM " + fromQuery)
			}
			args = append(args, fromArgs...)
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

func (q UpdateQuery) GetAlias() string {
	return q.Alias
}

func (q UpdateQuery) GetName() string {
	return ""
}

func (q UpdateQuery) NestThis() qx.Query {
	q.Nested = true
	return q
}

func (q UpdateQuery) As(alias string) UpdateQuery {
	q.Alias = alias
	return q
}

func Update(table qx.BaseTable) UpdateQuery {
	return UpdateQuery{
		UpdateTable: table,
		Alias:       qx.RandomString(8),
	}
}

func (q UpdateQuery) With(cteList ...qx.CTE) UpdateQuery {
	q.CTEs = append(q.CTEs, cteList...)
	return q
}

func (q UpdateQuery) Update(tbl qx.BaseTable) UpdateQuery {
	q.UpdateTable = tbl
	return q
}

func (q UpdateQuery) Set(sets ...qx.FieldValueSet) UpdateQuery {
	q.SetFields = append(q.SetFields, sets...)
	return q
}

func (q UpdateQuery) From(tbl qx.Table) UpdateQuery {
	q.FromTable = tbl
	return q
}

func (q UpdateQuery) Join(tbl qx.Table, pred qx.Predicate, preds ...qx.Predicate) UpdateQuery {
	preds = append([]qx.Predicate{pred}, preds...)
	q.JoinGroups = append(q.JoinGroups, qx.JoinTable{
		JoinType:     qx.JoinTypeDefault,
		Table:        tbl,
		OnPredicates: qx.VariadicPredicate{Predicates: preds},
	})
	return q
}

func (q UpdateQuery) LeftJoin(tbl qx.Table, pred qx.Predicate, preds ...qx.Predicate) UpdateQuery {
	preds = append([]qx.Predicate{pred}, preds...)
	q.JoinGroups = append(q.JoinGroups, qx.JoinTable{
		JoinType:     qx.JoinTypeLeft,
		Table:        tbl,
		OnPredicates: qx.VariadicPredicate{Predicates: preds},
	})
	return q
}

func (q UpdateQuery) RightJoin(tbl qx.Table, pred qx.Predicate, preds ...qx.Predicate) UpdateQuery {
	preds = append([]qx.Predicate{pred}, preds...)
	q.JoinGroups = append(q.JoinGroups, qx.JoinTable{
		JoinType:     qx.JoinTypeRight,
		Table:        tbl,
		OnPredicates: qx.VariadicPredicate{Predicates: preds},
	})
	return q
}

func (q UpdateQuery) FullJoin(tbl qx.Table, pred qx.Predicate, preds ...qx.Predicate) UpdateQuery {
	preds = append([]qx.Predicate{pred}, preds...)
	q.JoinGroups = append(q.JoinGroups, qx.JoinTable{
		JoinType:     qx.JoinTypeFull,
		Table:        tbl,
		OnPredicates: qx.VariadicPredicate{Predicates: preds},
	})
	return q
}

func (q UpdateQuery) CrossJoin(tbl qx.Table) UpdateQuery {
	q.JoinGroups = append(q.JoinGroups, qx.JoinTable{
		JoinType: qx.JoinTypeCross,
		Table:    tbl,
	})
	return q
}

func (q UpdateQuery) Where(preds ...qx.Predicate) UpdateQuery {
	q.WherePredicates.Predicates = append(q.WherePredicates.Predicates, preds...)
	return q
}

func (q UpdateQuery) Returning(fields ...qx.Field) UpdateQuery {
	q.ReturningFields = append(q.ReturningFields, fields...)
	return q
}

func (q UpdateQuery) ReturningOne() UpdateQuery {
	q.ReturningFields = qx.Fields{qx.FieldLiteral("1")}
	return q
}

func (q UpdateQuery) Returningx(mapper func(Row), accumulator func()) UpdateQuery {
	q.Mapper = mapper
	q.Accumulator = accumulator
	return q
}

func (q UpdateQuery) ReturningRowx(mapper func(Row)) UpdateQuery {
	q.Mapper = mapper
	return q
}

func (q UpdateQuery) Fetch(db qx.DB) (err error) {
	q.LogSkip += 1
	return q.FetchContext(nil, db)
}

func (q UpdateQuery) FetchContext(ctx context.Context, db qx.DB) (err error) {
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
		if q.DB == nil {
			return errors.New("DB cannot be nil")
		}
		db = q.DB
	}
	r := &QyRow{QxRow: &qx.QxRow{}}
	if q.Mapper != nil {
		q.Mapper(r) // call the mapper once on the *Row to get all the selected that the user is interested in
	}
	q.ReturningFields = r.QxRow.Fields // then, transfer the selected collected by *Row to the UpdateQuery
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

func (q UpdateQuery) Exec(db qx.DB) (sql.Result, error) {
	q.LogSkip += 1
	return q.ExecContext(nil, db)
}

func (q UpdateQuery) ExecContext(ctx context.Context, db qx.DB) (sql.Result, error) {
	var res sql.Result
	var err error
	if db == nil {
		if q.DB == nil {
			return res, errors.New("DB cannot be nil")
		}
		db = q.DB
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
