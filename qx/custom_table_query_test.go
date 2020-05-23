package qx

import "testing"

func TestCustomTable_GameTheNumbers(t *testing.T) {
	CustomTable{
		CustomSprintf: func(string, []interface{}, []string) (string, []interface{}) {
			return "", nil
		},
	}.ToSQL()
	tbl := CustomTable{}
	tbl.ToSQL()
	tbl.As("tbl")
	tbl.GetAlias()
	tbl.GetName()
}

func TestCustomQuery_GameTheNumbers(t *testing.T) {
	CustomQuery{
		CustomSprintf: func(string, []interface{}, []string) (string, []interface{}) {
			return "", nil
		},
	}.ToSQL()
	q := CustomQuery{}
	q.ToSQL()
	q.As("q")
	q.GetAlias()
	q.GetName()
	q.NestThis()
}
