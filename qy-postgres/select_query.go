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
	"time"

	"github.com/bokwoon95/qy/qx"
)

type SelectQuery struct {
	Nested bool
	Alias  string
	// WITH
	CTEs qx.CTEs
	// SELECT
	SelectType   qx.SelectType
	DistinctOn   qx.Fields
	SelectFields qx.Fields
	// FROM
	FromTable  qx.Table
	JoinGroups qx.JoinTables
	// WHERE
	WherePredicates qx.VariadicPredicate
	// GROUP BY
	GroupByFields qx.Fields
	// HAVING
	HavingPredicates qx.VariadicPredicate
	// ORDER BY
	OrderByFields qx.Fields
	// LIMIT
	LimitValue *uint64
	// OFFSET
	OffsetValue *uint64
	// DB
	DB          qx.DB
	Mapper      func(Row)
	Accumulator func()
	// Logging
	Log     qx.Logger
	LogFlag int
	LogSkip int
}

func (q SelectQuery) ToSQL() (string, []interface{}) {
	var buf = &strings.Builder{}
	var args []interface{}
	// WITH
	q.CTEs.WriteSQL(buf, &args)
	{ // SELECT
		tempBuf, tempArgs := &strings.Builder{}, []interface{}{}
		if q.SelectFields.WriteSQLWithAlias(tempBuf, &tempArgs, "", "", nil) {
			if q.SelectType == "" {
				q.SelectType = qx.SelectTypeDefault
			}
			if buf.Len() > 0 {
				buf.WriteString(" ")
			}
			switch q.SelectType {
			case qx.SelectTypeDistinctOn:
				q.DistinctOn.WriteSQL(buf, &args, string(q.SelectType)+" (", ") "+tempBuf.String(), nil)
			default:
				buf.WriteString(string(q.SelectType) + " " + tempBuf.String())
			}
			args = append(args, tempArgs...)
		}
	}
	{ // FROM
		fromQuery, fromArgs := "", []interface{}{}
		if q.FromTable != nil {
			fromQuery, fromArgs = q.FromTable.ToSQL()
		}
		if fromQuery != "" {
			if buf.Len() > 0 {
				buf.WriteString(" ")
			}
			if _, ok := q.FromTable.(qx.Query); ok {
				fromQuery = "(" + fromQuery + ")"
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
	// GROUP BY
	q.GroupByFields.WriteSQL(buf, &args, "GROUP BY ", "", nil)
	// HAVING
	q.HavingPredicates.Toplevel = true
	q.HavingPredicates.WriteSQL(buf, &args, "HAVING ", "", nil)
	// ORDER BY
	q.OrderByFields.WriteSQL(buf, &args, "ORDER BY ", "", nil)
	// LIMIT
	if q.LimitValue != nil {
		if buf.Len() > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString("LIMIT ?")
		args = append(args, *q.LimitValue)
	}
	// OFFSET
	if q.OffsetValue != nil {
		if buf.Len() > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString("OFFSET ?")
		args = append(args, *q.OffsetValue)
	}
	query := buf.String()
	if !q.Nested {
		query = qx.MySQLToPostgresPlaceholders(query)
		if q.Log != nil {
			var logOutput string
			switch {
			case LStats&q.LogFlag != 0:
				logOutput = "\n----[ Executing query ]----\n" + query + " " + fmt.Sprint(args) +
					"\n----[ with bind values ]----\n" + qx.PostgresInterpolateSQL(query, args...)
			case LInterpolate&q.LogFlag != 0:
				logOutput = qx.PostgresInterpolateSQL(query, args...)
			default:
				logOutput = query + " " + fmt.Sprint(args)
			}
			switch q.Log.(type) {
			case *log.Logger:
				q.Log.Output(q.LogSkip+2, logOutput)
			default:
				q.Log.Output(q.LogSkip+1, logOutput)
			}
		}
	}
	return query, args
}

func From(table qx.Table) SelectQuery {
	return SelectQuery{
		FromTable: table,
		Alias:     qx.RandomString(8),
	}
}

func Select(fields ...qx.Field) SelectQuery {
	return SelectQuery{
		SelectFields: fields,
		Alias:        qx.RandomString(8),
	}
}

func SelectOne(fields ...qx.Field) SelectQuery {
	return SelectQuery{
		SelectFields: qx.Fields{qx.FieldLiteral("1")},
		Alias:        qx.RandomString(8),
	}
}

func SelectDistinct(fields ...qx.Field) SelectQuery {
	return SelectQuery{
		SelectType:   qx.SelectTypeDistinct,
		SelectFields: fields,
		Alias:        qx.RandomString(8),
	}
}

func SelectDistinctOn(distinctFields ...qx.Field) func(...qx.Field) SelectQuery {
	return func(fields ...qx.Field) SelectQuery {
		return SelectQuery{
			SelectType:   qx.SelectTypeDistinctOn,
			DistinctOn:   distinctFields,
			SelectFields: fields,
			Alias:        qx.RandomString(8),
		}
	}
}

func Selectx(mapper func(Row), accumulator func()) SelectQuery {
	return SelectQuery{
		Mapper:      mapper,
		Accumulator: accumulator,
		Alias:       qx.RandomString(8),
	}
}

func SelectRowx(mapper func(Row)) SelectQuery {
	return SelectQuery{
		Mapper: mapper,
		Alias:  qx.RandomString(8),
	}
}

func (q SelectQuery) With(ctes ...qx.CTE) SelectQuery {
	q.CTEs = append(q.CTEs, ctes...)
	return q
}

func (q SelectQuery) Select(fields ...qx.Field) SelectQuery {
	q.SelectFields = append(q.SelectFields, fields...)
	return q
}

func (q SelectQuery) SelectOne() SelectQuery {
	q.SelectFields = qx.Fields{qx.FieldLiteral("1")}
	return q
}

func (q SelectQuery) SelectAll() SelectQuery {
	q.SelectFields = qx.Fields{qx.FieldLiteral("*")}
	return q
}

func (q SelectQuery) SelectCount() SelectQuery {
	q.SelectFields = qx.Fields{qx.FieldLiteral("COUNT(*)")}
	return q
}

func (q SelectQuery) SelectDistinct(fields ...qx.Field) SelectQuery {
	q.SelectType = qx.SelectTypeDistinct
	return q.Select(fields...)
}

func (q SelectQuery) SelectDistinctOn(distinctFields ...qx.Field) func(...qx.Field) SelectQuery {
	return func(fields ...qx.Field) SelectQuery {
		q.SelectType = qx.SelectTypeDistinctOn
		q.DistinctOn = append(q.DistinctOn, distinctFields...)
		return q.Select(fields...)
	}
}

func (q SelectQuery) From(table qx.Table) SelectQuery {
	q.FromTable = table
	return q
}

func (q SelectQuery) Join(table qx.Table, predicate qx.Predicate, predicates ...qx.Predicate) SelectQuery {
	predicates = append([]qx.Predicate{predicate}, predicates...)
	q.JoinGroups = append(q.JoinGroups, qx.JoinTable{
		JoinType:     qx.JoinTypeDefault,
		Table:        table,
		OnPredicates: qx.VariadicPredicate{Predicates: predicates},
	})
	return q
}

func (q SelectQuery) LeftJoin(table qx.Table, predicate qx.Predicate, predicates ...qx.Predicate) SelectQuery {
	predicates = append([]qx.Predicate{predicate}, predicates...)
	q.JoinGroups = append(q.JoinGroups, qx.JoinTable{
		JoinType:     qx.JoinTypeLeft,
		Table:        table,
		OnPredicates: qx.VariadicPredicate{Predicates: predicates},
	})
	return q
}

func (q SelectQuery) RightJoin(table qx.Table, predicate qx.Predicate, predicates ...qx.Predicate) SelectQuery {
	predicates = append([]qx.Predicate{predicate}, predicates...)
	q.JoinGroups = append(q.JoinGroups, qx.JoinTable{
		JoinType:     qx.JoinTypeRight,
		Table:        table,
		OnPredicates: qx.VariadicPredicate{Predicates: predicates},
	})
	return q
}

func (q SelectQuery) FullJoin(table qx.Table, predicate qx.Predicate, predicates ...qx.Predicate) SelectQuery {
	predicates = append([]qx.Predicate{predicate}, predicates...)
	q.JoinGroups = append(q.JoinGroups, qx.JoinTable{
		JoinType:     qx.JoinTypeFull,
		Table:        table,
		OnPredicates: qx.VariadicPredicate{Predicates: predicates},
	})
	return q
}

func (q SelectQuery) CrossJoin(table qx.Table) SelectQuery {
	q.JoinGroups = append(q.JoinGroups, qx.JoinTable{
		JoinType: qx.JoinTypeCross,
		Table:    table,
	})
	return q
}

func (q SelectQuery) Where(predicates ...qx.Predicate) SelectQuery {
	q.WherePredicates.Predicates = append(q.WherePredicates.Predicates, predicates...)
	return q
}

func (q SelectQuery) GroupBy(fields ...qx.Field) SelectQuery {
	q.GroupByFields = append(q.GroupByFields, fields...)
	return q
}

func (q SelectQuery) Having(predicates ...qx.Predicate) SelectQuery {
	q.HavingPredicates.Predicates = append(q.HavingPredicates.Predicates, predicates...)
	return q
}

func (q SelectQuery) OrderBy(fields ...qx.Field) SelectQuery {
	q.OrderByFields = append(q.OrderByFields, fields...)
	return q
}

func (q SelectQuery) Limit(limit int) SelectQuery {
	if limit < 0 {
		limit = -limit
	}
	num := uint64(limit)
	q.LimitValue = &num
	return q
}

func (q SelectQuery) Offset(offset int) SelectQuery {
	if offset < 0 {
		offset = -offset
	}
	num := uint64(offset)
	q.OffsetValue = &num
	return q
}

func (q SelectQuery) Selectx(mapper func(Row), accumulator func()) SelectQuery {
	q.Mapper = mapper
	q.Accumulator = accumulator
	return q
}

func (q SelectQuery) SelectRowx(mapper func(Row)) SelectQuery {
	q.Mapper = mapper
	return q
}

func (q SelectQuery) Fetch(db qx.DB) error {
	q.LogSkip += 1
	return q.FetchContext(nil, db)
}

func (q SelectQuery) FetchContext(ctx context.Context, db qx.DB) (err error) {
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
	logBuf := &strings.Builder{}
	var rowcount int
	defer func() func() {
		var logskip int
		switch q.Log.(type) {
		case *log.Logger:
			logskip = q.LogSkip + 2
		default:
			logskip = q.LogSkip + 1
		}
		start := time.Now()
		return func() {
			elapsed := time.Since(start)
			if LResults&q.LogFlag != 0 && q.Log != nil && rowcount > 5 {
				logBuf.WriteString("\n...")
			}
			if LStats&q.LogFlag != 0 && q.Log != nil {
				logBuf.WriteString("\n(Fetched " + strconv.Itoa(rowcount) + " rows in " + elapsed.String() + ")")
			}
			if logBuf.Len() > 0 && q.Log != nil {
				q.Log.Output(logskip, logBuf.String())
			}
		}
	}()()
	if db == nil {
		if q.DB == nil {
			return errors.New("DB cannot be nil")
		}
		db = q.DB
	}
	r := &QyRow{QxRow: &qx.QxRow{}}
	if q.Mapper != nil {
		q.Mapper(r)                     // call the mapper once on the *Row to get all the selected that the user is interested in
		q.SelectFields = r.QxRow.Fields // then, transfer the selected collected by *Row to the SelectQuery
		if len(q.SelectFields) == 0 {
			q.SelectFields = append(q.SelectFields, Fieldf("1"))
		}
	}
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
		if LResults&q.LogFlag != 0 && q.Log != nil && rowcount <= 5 {
			logBuf.WriteString("\n----[ Row " + strconv.Itoa(rowcount) + " ]----")
			for i := range r.QxRow.Dest {
				q, a := r.QxRow.Fields[i].ToSQLExclude(nil)
				logBuf.WriteString("\n" + qx.MySQLInterpolateSQL(q, a...) + ": " + qx.ArgToStringV2(r.QxRow.Dest[i]))
			}
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
	if e := r.QxRow.Rows.Close(); e != nil {
		return e
	}
	return r.QxRow.Rows.Err()
}

func (q SelectQuery) Exec(db qx.DB) (sql.Result, error) {
	q.LogSkip += 1
	return q.ExecContext(nil, db)
}

func (q SelectQuery) ExecContext(ctx context.Context, db qx.DB) (sql.Result, error) {
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

func (q SelectQuery) As(alias string) SelectQuery {
	q.Alias = alias
	return q
}

func (q SelectQuery) Get(fieldName string) qx.CustomField {
	return Fieldf(q.Alias + "." + fieldName)
}

func (q SelectQuery) GetAlias() string {
	return q.Alias
}

func (q SelectQuery) GetName() string {
	return ""
}

func (q SelectQuery) NestThis() qx.Query {
	q.Nested = true
	return q
}
