package db

import (
	"context"
	"log"
	"testing"

	"github.com/MyBeLoVedL/payroll/util"

	_ "github.com/go-sql-driver/mysql"

	r "github.com/stretchr/testify/require"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreateEmployee(t *testing.T) {
	arg := CreateEmployeeParams{
		Type:                  EmployeesType(util.RandType()),
		Mail:                  util.RandStr(21),
		SocialSecurityNumber:  util.RandStr(15),
		StandardTaxDeductions: "0.12",
		OtherDuductions:       "10000",
		PhoneNumber:           util.RandDigits(11),
		Rate:                  "0.25",
	}

	res, err := queries.CreateEmployee(context.Background(), arg)
	check(err)
	id, err := res.LastInsertId()
	r.Equal(t, 1, id)
}
