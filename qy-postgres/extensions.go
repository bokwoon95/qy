package qy

import (
	"errors"

	"github.com/bokwoon95/qy/qx"
)

func Exists(query qx.Query, db qx.DB) (exists bool, err error) {
	var dbV2 qx.DB
	var logger qx.Logger
	switch q := query.(type) {
	case SelectQuery:
		q.SelectFields = []qx.Field{Fieldf("1")}
		dbV2 = q.DB
		logger = q.Log
		query = q
	case InsertQuery:
		q.ReturningFields = []qx.Field{Fieldf("1")}
		dbV2 = q.DB
		logger = q.Log
		query = q
	case UpdateQuery:
		q.ReturningFields = []qx.Field{Fieldf("1")}
		dbV2 = q.DB
		logger = q.Log
		query = q
	case DeleteQuery:
		q.ReturningFields = []qx.Field{Fieldf("1")}
		dbV2 = q.DB
		logger = q.Log
		query = q
	default:
		return exists, errors.New("query is not a SelectQuery, InsertQuery, UpdateQuery or DeleteQuery")
	}
	if db == nil && dbV2 != nil {
		db = dbV2
	}
	if db == nil {
		return exists, errors.New("DB is not set")
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
