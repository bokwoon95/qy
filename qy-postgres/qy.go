package qy

import (
	"database/sql"
	"time"

	"github.com/bokwoon95/qy/qx"
	"github.com/lib/pq"
)

type Row interface {
	ScanArray(array interface{}, field qx.Field)
	ScanInto(dest interface{}, field qx.Field)
	// bool
	Bool(qx.BooleanField) bool
	Bool_(qx.Field) bool
	BoolValid(qx.BooleanField) bool
	BoolValid_(qx.Field) bool
	NullBool(qx.BooleanField) sql.NullBool
	NullBool_(qx.Field) sql.NullBool
	// float64
	Float64(qx.NumberField) float64
	Float64_(qx.Field) float64
	Float64Valid(qx.NumberField) bool
	Float64Valid_(qx.Field) bool
	NullFloat64(qx.NumberField) sql.NullFloat64
	NullFloat64_(qx.Field) sql.NullFloat64
	// int
	Int(qx.NumberField) int
	Int_(qx.Field) int
	IntValid(qx.NumberField) bool
	IntValid_(qx.Field) bool
	// int64
	Int64(qx.NumberField) int64
	Int64_(qx.Field) int64
	Int64Valid(qx.NumberField) bool
	Int64Valid_(qx.Field) bool
	NullInt64(qx.NumberField) sql.NullInt64
	NullInt64_(qx.Field) sql.NullInt64
	// string
	String(qx.StringField) string
	String_(qx.Field) string
	StringValid(qx.StringField) bool
	StringValid_(qx.Field) bool
	NullString(qx.StringField) sql.NullString
	NullString_(qx.Field) sql.NullString
	// time.Time
	Time(qx.TimeField) time.Time
	Time_(qx.Field) time.Time
	TimeValid(qx.TimeField) bool
	TimeValid_(qx.Field) bool
	NullTime(qx.TimeField) sql.NullTime
	NullTime_(qx.Field) sql.NullTime
}

// QyRow is a wrapper around QxRow that additionally implements the scanning of
// postgres arrays into go slices. Only []bool, []float64, []int64 or
// []string slices are supported.
type QyRow struct {
	*qx.QxRow
}

// ScanArray implements Row interface. It received a pointer to a slice and
// scans a postgres array into that slice. Only []bool, []float64, []int64 or
// []string slices are supported.
func (r *QyRow) ScanArray(array interface{}, f qx.Field) {
	nothing := &sql.RawBytes{}
	if r.QxRow.Rows == nil {
		r.QxRow.Fields = append(r.QxRow.Fields, f)
		r.QxRow.Dest = append(r.QxRow.Dest, nothing)
		return
	}
	if len(r.QxRow.TmpDest) != len(r.QxRow.Dest) {
		r.QxRow.TmpDest = make([]interface{}, len(r.QxRow.Dest))
		for i := range r.QxRow.TmpDest {
			r.QxRow.TmpDest[i] = nothing
		}
	}
	r.TmpDest[r.Index] = pq.Array(array)
	r.Rows.Scan(r.TmpDest...)
	r.TmpDest[r.Index] = nothing
	r.Index++
}

func Fieldf(format string, values ...interface{}) qx.CustomField {
	return qx.CustomField{
		Format: format,
		Values: values,
	}
}

func Predicatef(format string, values ...interface{}) qx.CustomPredicate {
	return qx.CustomPredicate{
		Format: format,
		Values: values,
	}
}

func Tablef(format string, values ...interface{}) qx.CustomTable {
	return qx.CustomTable{
		Format: format,
		Values: values,
	}
}

func Queryf(format string, values ...interface{}) qx.CustomQuery {
	return qx.CustomQuery{
		Postgres: true,
		Format:   format,
		Values:   values,
	}
}

type BaseQuery struct {
	DB   qx.DB
	Log  qx.Logger
	CTEs qx.CTEs
}

func WithLog(logger qx.Logger) BaseQuery {
	return BaseQuery{
		Log: logger,
	}
}

func WithDB(db qx.DB) BaseQuery {
	return BaseQuery{
		DB: db,
	}
}

func With(CTEs ...qx.CTE) BaseQuery {
	return BaseQuery{
		CTEs: CTEs,
	}
}

func (qy BaseQuery) WithLog(logger qx.Logger) BaseQuery {
	qy.Log = logger
	return qy
}

func (qy BaseQuery) WithDB(db qx.DB) BaseQuery {
	qy.DB = db
	return qy
}

func (qy BaseQuery) With(CTEs ...qx.CTE) BaseQuery {
	qy.CTEs = CTEs
	return qy
}

func (qy BaseQuery) From(table qx.Table) SelectQuery {
	return SelectQuery{
		FromTable: table,
		Alias:     qx.RandomString(8),
		CTEs:      qy.CTEs,
		DB:        qy.DB,
		Log:       qy.Log,
	}
}

func (qy BaseQuery) Select(fields ...qx.Field) SelectQuery {
	return SelectQuery{
		SelectFields: fields,
		Alias:        qx.RandomString(8),
		CTEs:         qy.CTEs,
		Log:          qy.Log,
		DB:           qy.DB,
	}
}

func (qy BaseQuery) SelectOne() SelectQuery {
	return SelectQuery{
		SelectFields: qx.Fields{qx.FieldLiteral("1")},
		Alias:        qx.RandomString(8),
		CTEs:         qy.CTEs,
		Log:          qy.Log,
		DB:           qy.DB,
	}
}

func (qy BaseQuery) SelectAll() SelectQuery {
	return SelectQuery{
		SelectFields: qx.Fields{qx.FieldLiteral("*")},
		Alias:        qx.RandomString(8),
		CTEs:         qy.CTEs,
		Log:          qy.Log,
		DB:           qy.DB,
	}
}

func (qy BaseQuery) SelectCount() SelectQuery {
	return SelectQuery{
		SelectFields: qx.Fields{qx.FieldLiteral("COUNT(*)")},
		Alias:        qx.RandomString(8),
		CTEs:         qy.CTEs,
		Log:          qy.Log,
		DB:           qy.DB,
	}
}

func (qy BaseQuery) SelectDistinct(fields ...qx.Field) SelectQuery {
	return SelectQuery{
		SelectType:   qx.SelectTypeDistinct,
		SelectFields: fields,
		Alias:        qx.RandomString(8),
		CTEs:         qy.CTEs,
		DB:           qy.DB,
		Log:          qy.Log,
	}
}

func (qy BaseQuery) SelectDistinctOn(distinctFields ...qx.Field) func(...qx.Field) SelectQuery {
	return func(fields ...qx.Field) SelectQuery {
		return SelectQuery{
			SelectType:   qx.SelectTypeDistinctOn,
			DistinctOn:   distinctFields,
			SelectFields: fields,
			Alias:        qx.RandomString(8),
			CTEs:         qy.CTEs,
			DB:           qy.DB,
			Log:          qy.Log,
		}
	}
}

func (qy BaseQuery) Selectx(mapper func(Row), accumulator func()) SelectQuery {
	return SelectQuery{
		Mapper:      mapper,
		Accumulator: accumulator,
		Alias:       qx.RandomString(8),
		CTEs:        qy.CTEs,
		DB:          qy.DB,
		Log:         qy.Log,
	}
}

func (qy BaseQuery) SelectRowx(mapper func(Row)) SelectQuery {
	return SelectQuery{
		Mapper: mapper,
		Alias:  qx.RandomString(8),
		CTEs:   qy.CTEs,
		DB:     qy.DB,
		Log:    qy.Log,
	}
}

func (qy BaseQuery) InsertInto(table qx.BaseTable) InsertQuery {
	return InsertQuery{
		IntoTable: table,
		Alias:     qx.RandomString(8),
		CTEs:      qy.CTEs,
		DB:        qy.DB,
		Log:       qy.Log,
	}
}

func (qy BaseQuery) Update(table qx.BaseTable) UpdateQuery {
	return UpdateQuery{
		UpdateTable: table,
		Alias:       qx.RandomString(8),
		CTEs:        qy.CTEs,
		DB:          qy.DB,
		Log:         qy.Log,
	}
}

func (qy BaseQuery) DeleteFrom(table qx.BaseTable) DeleteQuery {
	return DeleteQuery{
		FromTable: table,
		Alias:     qx.RandomString(8),
		CTEs:      qy.CTEs,
		DB:        qy.DB,
		Log:       qy.Log,
	}
}
