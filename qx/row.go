package qx

import (
	"database/sql"
	"fmt"
	"runtime"
	"time"
)

type QxRow struct {
	Rows    *sql.Rows
	Index   int
	Fields  []Field
	Dest    []interface{}
	TmpDest []interface{}
}

/* custom */

func (r *QxRow) ScanInto(dest interface{}, field Field) {
	var nothing interface{}
	if r.Rows == nil {
		r.Fields = append(r.Fields, field)
		r.Dest = append(r.Dest, dest)
		return
	}
	if len(r.TmpDest) != len(r.Dest) {
		r.TmpDest = make([]interface{}, len(r.Dest))
		for i := range r.TmpDest {
			r.TmpDest[i] = &nothing
		}
	}
	r.TmpDest[r.Index] = dest
	err := r.Rows.Scan(r.TmpDest...)
	if err != nil {
		_, sourcefile, linenbr, _ := runtime.Caller(1)
		panic(fmt.Errorf("row.ScanInto failed on %s:%d: %w", sourcefile, linenbr, err))
	}
	r.TmpDest[r.Index] = &nothing
	r.Index++
}

/* bool */

func (r *QxRow) Bool(field BooleanField) bool {
	return r.NullBool_(field).Bool
}

func (r *QxRow) Bool_(field Field) bool {
	return r.NullBool_(field).Bool
}

func (r *QxRow) BoolValid(field BooleanField) bool {
	return r.NullBool_(field).Valid
}

func (r *QxRow) BoolValid_(field Field) bool {
	return r.NullBool_(field).Valid
}

func (r *QxRow) NullBool(field BooleanField) sql.NullBool {
	return r.NullBool_(field)
}

func (r *QxRow) NullBool_(field Field) sql.NullBool {
	if r.Rows == nil {
		r.Fields = append(r.Fields, field)
		r.Dest = append(r.Dest, &sql.NullBool{})
		return sql.NullBool{}
	}
	switch val := r.Dest[r.Index].(type) {
	case *sql.NullBool:
		r.Index++
		return *val
	default:
		panic("type mismatch")
	}
}

/* float64 */

func (r *QxRow) Float64(field NumberField) float64 {
	return r.NullFloat64_(field).Float64
}

func (r *QxRow) Float64_(field Field) float64 {
	return r.NullFloat64_(field).Float64
}

func (r *QxRow) Float64Valid(field NumberField) bool {
	return r.NullFloat64_(field).Valid
}

func (r *QxRow) Float64Valid_(field Field) bool {
	return r.NullFloat64_(field).Valid
}

func (r *QxRow) NullFloat64(field NumberField) sql.NullFloat64 {
	return r.NullFloat64_(field)
}

func (r *QxRow) NullFloat64_(field Field) sql.NullFloat64 {
	if r.Rows == nil {
		r.Fields = append(r.Fields, field)
		r.Dest = append(r.Dest, &sql.NullFloat64{})
		return sql.NullFloat64{}
	}
	switch val := r.Dest[r.Index].(type) {
	case *sql.NullFloat64:
		r.Index++
		return *val
	default:
		panic("type mismatch")
	}
}

/* int */

func (r *QxRow) Int(field NumberField) int {
	return int(r.NullInt64_(field).Int64)
}

func (r *QxRow) Int_(field Field) int {
	return int(r.NullInt64_(field).Int64)
}

func (r *QxRow) IntValid(field NumberField) bool {
	return r.NullInt64_(field).Valid
}

func (r *QxRow) IntValid_(field Field) bool {
	return r.NullInt64_(field).Valid
}

/* int64 */

func (r *QxRow) Int64(field NumberField) int64 {
	return r.NullInt64_(field).Int64
}

func (r *QxRow) Int64_(field Field) int64 {
	return r.NullInt64_(field).Int64
}

func (r *QxRow) Int64Valid(field NumberField) bool {
	return r.NullInt64_(field).Valid
}

func (r *QxRow) Int64Valid_(field Field) bool {
	return r.NullInt64_(field).Valid
}

func (r *QxRow) NullInt64(field NumberField) sql.NullInt64 {
	return r.NullInt64_(field)
}

func (r *QxRow) NullInt64_(field Field) sql.NullInt64 {
	if r.Rows == nil {
		r.Fields = append(r.Fields, field)
		r.Dest = append(r.Dest, &sql.NullInt64{})
		return sql.NullInt64{}
	}
	switch val := r.Dest[r.Index].(type) {
	case *sql.NullInt64:
		r.Index++
		return *val
	default:
		panic("type mismatch")
	}
}

/* string */

func (r *QxRow) String(field StringField) string {
	return r.NullString_(field).String
}

func (r *QxRow) String_(field Field) string {
	return r.NullString_(field).String
}

func (r *QxRow) StringValid(field StringField) bool {
	return r.NullString_(field).Valid
}

func (r *QxRow) StringValid_(field Field) bool {
	return r.NullString_(field).Valid
}

func (r *QxRow) NullString(field StringField) sql.NullString {
	return r.NullString_(field)
}

func (r *QxRow) NullString_(field Field) sql.NullString {
	if r.Rows == nil {
		r.Fields = append(r.Fields, field)
		r.Dest = append(r.Dest, &sql.NullString{})
		return sql.NullString{}
	}
	switch val := r.Dest[r.Index].(type) {
	case *sql.NullString:
		r.Index++
		return *val
	default:
		panic("type mismatch")
	}
}

/* time.Time */

func (r *QxRow) Time(field TimeField) time.Time {
	return r.NullTime_(field).Time
}

func (r *QxRow) Time_(field Field) time.Time {
	return r.NullTime_(field).Time
}

func (r *QxRow) TimeValid(field TimeField) bool {
	return r.NullTime_(field).Valid
}

func (r *QxRow) TimeValid_(field Field) bool {
	return r.NullTime_(field).Valid
}

func (r *QxRow) NullTime(field TimeField) sql.NullTime {
	return r.NullTime_(field)
}

func (r *QxRow) NullTime_(field Field) sql.NullTime {
	if r.Rows == nil {
		r.Fields = append(r.Fields, field)
		r.Dest = append(r.Dest, &sql.NullTime{})
		return sql.NullTime{}
	}
	switch val := r.Dest[r.Index].(type) {
	case *sql.NullTime:
		r.Index++
		return *val
	default:
		panic("type mismatch")
	}
}
