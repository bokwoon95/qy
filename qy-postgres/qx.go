package qy

import (
	"time"

	"github.com/bokwoon95/qy/qx"
)

func NewCTE(name string, query qx.Query) qx.CTE {
	return qx.CTE{
		Name:  name,
		Query: query,
	}
}

func And(predicates ...qx.Predicate) qx.Predicate { return qx.And(predicates...) }
func Or(predicates ...qx.Predicate) qx.Predicate  { return qx.Or(predicates...) }

func Array(slice interface{}) qx.ArrayField { return qx.Array(slice) }
func Bytes(b []byte) qx.BinaryField         { return qx.Bytes(b) }
func Bool(b bool) qx.BooleanField           { return qx.Bool(b) }
func Int(num int) qx.NumberField            { return qx.Int(num) }
func Int64(num int64) qx.NumberField        { return qx.Int64(num) }
func Float64(num float64) qx.NumberField    { return qx.Float64(num) }
func String(s string) qx.StringField        { return qx.String(s) }
func Time(t time.Time) qx.TimeField         { return qx.Time(t) }

type Table = qx.Table
type Query = qx.Query
type BaseTable = qx.BaseTable
type Predicate = qx.Predicate
type Field = qx.Field
type Fields = qx.Fields
type FieldLiteral = qx.FieldLiteral
type ValuesList = qx.ValuesList
type Queryer = qx.Queryer
type QueryerContext = qx.QueryerContext
type Logger = qx.Logger
