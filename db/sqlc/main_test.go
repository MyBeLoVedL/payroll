package db

import (
	"database/sql"
	"testing"
)

var queries *Queries

func TestMain(m *testing.M) {
	dbIns, err := sql.Open("mysql", "root:Lp262783@/payroll")
	if err != nil {
		panic(err)
	}
	queries = New(dbIns)
	m.Run()
}
