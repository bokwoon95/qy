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

type InsertQuery struct {
	Nested bool
	Alias  string
	// WITH
	CTEs qx.CTEs
	// INSERT INTO
	IntoTable    qx.BaseTable
	InsertFields qx.Fields
	// VALUES
	ValuesList qx.ValuesList
	// SELECT
	SelectQuery *SelectQuery
	// ON CONFLICT
	HandleConflict       bool
	ConflictFields       qx.Fields
	ConflictPredicates   qx.VariadicPredicate
	ConflictConstraint   string
	Resolution           qx.FieldValueSets
	ResolutionPredicates qx.VariadicPredicate
	// RETURNING
	ReturningFields qx.Fields
	// DB
	DB          qx.DB
	Mapper      func(Row)
	Accumulator func()
	// Logging
	Log     qx.Logger
	LogFlag int
	LogSkip int
}

func (q InsertQuery) ToSQL() (string, []interface{}) {
	var buf = &strings.Builder{}
	var args []interface{}
	var excludeTableQualifiers []string
	// WITH
	q.CTEs.WriteSQL(buf, &args)
	{ // INSERT INTO
		intoQuery, intoArgs := "", []interface{}{}
		if q.IntoTable != nil {
			intoQuery, intoArgs = q.IntoTable.ToSQL()
			if q.IntoTable.GetAlias() != "" {
				excludeTableQualifiers = append(excludeTableQualifiers, q.IntoTable.GetAlias())
			} else if q.IntoTable.GetName() != "" {
				excludeTableQualifiers = append(excludeTableQualifiers, q.IntoTable.GetName())
			}
		}
		if intoQuery != "" {
			if buf.Len() > 0 {
				buf.WriteString(" ")
			}
			if q.IntoTable.GetAlias() != "" {
				buf.WriteString("INSERT INTO " + intoQuery + " AS " + q.IntoTable.GetAlias())
			} else {
				buf.WriteString("INSERT INTO " + intoQuery)
			}
			args = append(args, intoArgs...)
			q.InsertFields.WriteSQL(buf, &args, "(", ")", excludeTableQualifiers)
		}
	}
	// VALUES/SELECT
	switch {
	case len(q.ValuesList) > 0:
		q.ValuesList.WriteSQL(buf, &args, "VALUES ", "")
	case q.SelectQuery != nil:
		q.SelectQuery.Nested = true
		selectQuery, selectArgs := q.SelectQuery.ToSQL()
		if selectQuery != "" {
			if buf.Len() > 0 {
				buf.WriteString(" ")
			}
			buf.WriteString(selectQuery)
			args = append(args, selectArgs...)
		}
	}
	// ON CONFLICT
	var noConflict bool
	switch {
	case q.HandleConflict:
		switch {
		case q.ConflictConstraint != "":
			if buf.Len() > 0 {
				buf.WriteString(" ")
			}
			buf.WriteString("ON CONFLICT ON CONSTRAINT " + q.ConflictConstraint)
		case q.ConflictFields.WriteSQL(buf, &args, "ON CONFLICT (", ")", excludeTableQualifiers):
			q.ConflictPredicates.Toplevel = true
			q.ConflictPredicates.WriteSQL(buf, &args, "WHERE ", "", excludeTableQualifiers)
		default:
			if buf.Len() > 0 {
				buf.WriteString(" ")
			}
			buf.WriteString("ON CONFLICT")
		}
	default:
		noConflict = true
	}
	switch {
	case noConflict:
		// no-op
	case q.Resolution.WriteSQL(buf, &args, "DO UPDATE SET ", "", excludeTableQualifiers):
		q.ResolutionPredicates.Toplevel = true
		q.ResolutionPredicates.WriteSQL(buf, &args, "WHERE ", "", nil)
	default:
		if buf.Len() > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString("DO NOTHING")
	}
	// RETURNING
	q.ReturningFields.WriteSQLWithAlias(buf, &args, "RETURNING ", "", nil)
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
				logOutput = "Executing query: " + qx.PostgresInterpolateSQL(query, args...)
			default:
				logOutput = "Executing query: " + query + " " + fmt.Sprint(args)
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

func InsertInto(table qx.BaseTable) InsertQuery {
	return InsertQuery{
		IntoTable: table,
		Alias:     qx.RandomString(8),
	}
}

func (q InsertQuery) With(ctes ...qx.CTE) InsertQuery {
	q.CTEs = append(q.CTEs, ctes...)
	return q
}

func (q InsertQuery) InsertInto(table qx.BaseTable) InsertQuery {
	q.IntoTable = table
	return q
}

func (q InsertQuery) Columns(fields ...qx.Field) InsertQuery {
	q.InsertFields = append(q.InsertFields, fields...)
	return q
}

func (q InsertQuery) Values(values ...interface{}) InsertQuery {
	q.ValuesList = append(q.ValuesList, values)
	return q
}

func (q InsertQuery) InsertRow(sets ...qx.FieldValueSet) InsertQuery {
	fields, values := make([]qx.Field, len(sets)), make([]interface{}, len(sets))
	for i := range sets {
		fields[i] = sets[i].Field
		values[i] = sets[i].Value
	}
	if len(q.InsertFields) == 0 {
		q.InsertFields = fields
	}
	q.ValuesList = append(q.ValuesList, values)
	return q
}

func (q InsertQuery) Select(selectQuery SelectQuery) InsertQuery {
	q.SelectQuery = &selectQuery
	return q
}

func (q InsertQuery) OnConflict(fields ...qx.Field) insertConflict {
	q.HandleConflict = true
	q.ConflictFields = fields
	return insertConflict{insertQuery: &q}
}

func (q InsertQuery) OnConflictOnConstraint(name string) insertConflict {
	q.HandleConflict = true
	q.ConflictConstraint = name
	return insertConflict{insertQuery: &q}
}

type insertConflict struct{ insertQuery *InsertQuery }

func (c insertConflict) Where(predicates ...qx.Predicate) insertConflict {
	c.insertQuery.ConflictPredicates.Predicates = append(c.insertQuery.ConflictPredicates.Predicates, predicates...)
	return c
}

func (c insertConflict) DoNothing() InsertQuery {
	if c.insertQuery == nil {
		return InsertQuery{}
	}
	return *c.insertQuery
}

func (c insertConflict) DoUpdateSet(sets ...qx.FieldValueSet) InsertQuery {
	if c.insertQuery == nil {
		return InsertQuery{}
	}
	c.insertQuery.Resolution = append(c.insertQuery.Resolution, sets...)
	return *c.insertQuery
}

func Excluded(field qx.Field) qx.CustomField {
	return qx.CustomField{Format: "EXCLUDED." + field.GetName()}
}

func (q InsertQuery) Where(predicates ...qx.Predicate) InsertQuery {
	q.ResolutionPredicates.Predicates = append(q.ResolutionPredicates.Predicates, predicates...)
	return q
}

func (q InsertQuery) Returning(fields ...qx.Field) InsertQuery {
	q.ReturningFields = append(q.ReturningFields, fields...)
	return q
}

func (q InsertQuery) ReturningOne() InsertQuery {
	q.ReturningFields = qx.Fields{qx.FieldLiteral("1")}
	return q
}

func (q InsertQuery) Returningx(mapper func(Row), accumulator func()) InsertQuery {
	q.Mapper = mapper
	q.Accumulator = accumulator
	return q
}

func (q InsertQuery) ReturningRowx(mapper func(Row)) InsertQuery {
	q.Mapper = mapper
	return q
}

func (q InsertQuery) Fetch(db qx.DB) (err error) {
	q.LogSkip += 1
	return q.FetchContext(nil, db)
}

func (q InsertQuery) FetchContext(ctx context.Context, db qx.DB) (err error) {
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
		q.Mapper(r) // call the mapper once on the *Row to get all the selected that the user is interested in
	}
	q.ReturningFields = r.QxRow.Fields // then, transfer the selected collected by *Row to the InsertQuery
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
	return r.QxRow.Rows.Err()
}

func (q InsertQuery) Exec(db qx.DB) (sql.Result, error) {
	q.LogSkip += 1
	return q.ExecContext(nil, db)
}

func (q InsertQuery) ExecContext(ctx context.Context, db qx.DB) (sql.Result, error) {
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

func (q InsertQuery) As(alias string) InsertQuery {
	q.Alias = alias
	return q
}

func (q InsertQuery) GetAlias() string {
	return q.Alias
}

func (q InsertQuery) GetName() string {
	return ""
}

func (q InsertQuery) NestThis() qx.Query {
	q.Nested = true
	return q
}
