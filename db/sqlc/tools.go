package db

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var q *Queries

func init() {
	dbIns, err := sql.Open("mysql", "root:Lp262783@/payroll")
	if err != nil {
		panic(err)
	}
	q = New(dbIns)
}

func ValidateUser(idOrMail string, passwd string) error {
	var res Employee
	var err error
	if strings.Contains(idOrMail, "@") {
		res, err = q.SelectEmployeeByMail(context.Background(), idOrMail)
		if err != nil {
			return errors.New("no such user")
		}
	} else {
		id, err := strconv.Atoi(idOrMail)
		if err != nil {
			return errors.New("invalid user id")
		}
		res, err = q.SelectEmployeeById(context.Background(), int64(id))
		if err != nil {
			return errors.New("no such user")
		}
	}

	if res.Password.String != passwd {
		return errors.New("invalid password")
	}
	return nil
}
