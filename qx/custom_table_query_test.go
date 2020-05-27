package qx

import "testing"

func TestCustomTable_GameTheNumbers(t *testing.T) {
	CustomTable{}.ToSQL()
	tbl := CustomTable{}
	tbl.ToSQL()
	tbl.As("tbl")
	tbl.GetAlias()
	tbl.GetName()
}

func TestCustomQuery_GameTheNumbers(t *testing.T) {
	CustomQuery{}.ToSQL()
	q := CustomQuery{Postgres: true}
	q.ToSQL()
	q.As("q")
	q.GetAlias()
	q.GetName()
	q.NestThis()
}
