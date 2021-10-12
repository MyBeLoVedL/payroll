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
var dbIns *sql.DB

func init() {
	var err error
	dbIns, err = sql.Open("mysql", "lee:@@Lp262783@/payroll?parseTime=true")
	if err != nil {
		panic(err)
	}
	q = New(dbIns)
}

func ValidateUser(idOrMail string, passwd string) (int64, error) {
	var res Employee
	var err error
	var id int
	if strings.Contains(idOrMail, "@") {
		res, err = q.SelectEmployeeByMail(context.Background(), idOrMail)
		if err != nil {
			return 0, errors.New("no such user")
		}
	} else {
		id, err = strconv.Atoi(idOrMail)
		if err != nil {
			return 0, errors.New("invalid user id")
		}
		res, err = q.SelectEmployeeById(context.Background(), int64(id))
		if err != nil {
			return 0, errors.New("no such user")
		}
	}

	if res.Password.String != passwd {
		return 0, errors.New("invalid password")
	}
	return int64(id), nil
}

func UpdatePayment(method string, id int64) error {
	return q.UpdatePaymentMethod(context.Background(), UpdatePaymentMethodParams{EmployeesPaymentMethod(method), id})
}

func UpdatePaymentWIthMail(id int64, mail string) error {
	return q.UpdatePaymentMethodWithMail(context.Background(), UpdatePaymentMethodWithMailParams{EmployeesPaymentMethod("mail"), mail, id})
}

func UpdatePaymentWithBank(id int64, bank, account string) error {
	store := NewStore(dbIns)
	err := store.UpdatePaymentMethodWithBank(context.Background(), UpdateBankParam{id, bank, account})
	return err
}

func GetUser(id int64) (Employee, error) {
	return q.GetUser(context.Background(), id)
}

func AddOrder(ctx context.Context, arg AddOrderParams) (int64, error) {
	store := NewStore(dbIns)
	return store.AddOrder(context.Background(), arg)
}

func UpdateOrder(ctx context.Context, arg UpdateOrderParams) error {
	store := NewStore(dbIns)
	err := store.UpdateOrder(context.Background(), arg)
	return err
}
