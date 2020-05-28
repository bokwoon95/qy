package qy

import (
	"errors"

	"github.com/bokwoon95/qy/qx"
)

func Exists(query qx.Query, db qx.Queryer, logger qx.Logger) (exists bool, err error) {
	switch q := query.(type) {
	case SelectQuery:
		q.SelectFields = []qx.Field{Fieldf("1")}
		query = q
	case InsertQuery:
		q.ReturningFields = []qx.Field{Fieldf("1")}
		query = q
	case UpdateQuery:
		q.ReturningFields = []qx.Field{Fieldf("1")}
		query = q
	case DeleteQuery:
		q.ReturningFields = []qx.Field{Fieldf("1")}
		query = q
	default:
		return false, errors.New("query is not a SelectQuery, InsertQuery, UpdateQuery or DeleteQuery")
	}
	if db == nil {
		return exists, errors.New("db cannot be nil")
	}
	queryString, args := query.ToSQL()
	queryString = "SELECT EXISTS(" + queryString + ")"
	rows, err := db.Query(queryString, args...)
	if logger != nil {
		interpolatedQuery := qx.PostgresInterpolateSQL(queryString, args...)
		logger.Output(1, interpolatedQuery)
	}
	if err != nil {
		return exists, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&exists)
		if err != nil {
			return exists, err
		}
		break
	}
	if err = rows.Close(); err != nil {
		return exists, err
	}
	return exists, rows.Err()
}
