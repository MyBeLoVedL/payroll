package db

import (
	"context"
	"database/sql"
	"log"
	"payroll/db/util"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	r "github.com/stretchr/testify/require"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreateEmployee(t *testing.T) {
	arg := AddEmployeeParams{
		Type:                  EmployeesType(util.RandType()),
		Mail:                  util.RandDigits(10) + "@gmail.com",
		SocialSecurityNumber:  util.RandDigits(15),
		StandardTaxDeductions: "0.12",
		OtherDeductions:       "10000.0",
		PhoneNumber:           util.RandDigits(11),
		SalaryRate:            "0.25",
	}

	res, err := queries.AddEmployee(context.Background(), arg)
	check(err)
	id, err := res.LastInsertId()
	check(err)
	fetched, err := queries.SelectEmployeeById(context.Background(), id)
	check(err)
	r.Equal(t, arg.Type, fetched.Type)
	r.Equal(t, arg.Mail, fetched.Mail)
	r.Equal(t, arg.SocialSecurityNumber, fetched.SocialSecurityNumber)
	r.Equal(t, arg.PhoneNumber, fetched.PhoneNumber)
}

func TestUpdatePassword(t *testing.T) {
	id := 1
	pass := "new password"
	err := queries.UpdatePassword(context.Background(), UpdatePasswordParams{Password: sql.NullString{String: pass, Valid: true}, ID: int64(id)})
	r.NoError(t, err)

	f, err := queries.SelectEmployeeById(context.Background(), int64(id))
	r.NoError(t, err)

	r.Equal(t, pass, f.Password.String)
}

func TestDeleteEmployee(t *testing.T) {
	err := queries.DeleteEmployee(context.Background(), 1)
	r.NoError(t, err)
	row, err := queries.SelectEmployeeById(context.Background(), 1)
	r.NoError(t, err)
	r.Equal(t, int32(1), row.Deleted.Int32)
}

func TestUpdateEmployee(t *testing.T) {
	before, err := queries.SelectEmployeeById(context.Background(), 2)
	r.NoError(t, err)
	err = queries.UpdateEmployee(context.Background(), UpdateEmployeeParams{
		ID:                    before.ID,
		Type:                  before.Type,
		Mail:                  before.Mail,
		SocialSecurityNumber:  before.SocialSecurityNumber,
		StandardTaxDeductions: "0.23",
		OtherDeductions:       "20000.0",
		PhoneNumber:           before.PhoneNumber,
		SalaryRate:            before.SalaryRate,
		HourLimit:             before.HourLimit,
	})
	r.NoError(t, err)
	after, err := queries.SelectEmployeeById(context.Background(), 2)
	r.NoError(t, err)
	r.Equal(t, "0.23", after.StandardTaxDeductions)
	r.Equal(t, "20000.00", after.OtherDeductions)
}

func TestAddTimecards(t *testing.T) {
	_, err := queries.AddTimecard(context.Background(), 1)
	check(err)
}

func TestAddTimecardRecord(t *testing.T) {
	arg := AddTimecardRecordParams{
		ChargeNumber: 793,
		CardID:       1,
		Hours:        20,
		Date:         time.Now(),
	}
	err := queries.AddTimecardRecord(context.Background(), arg)
	r.NoError(t, err)
}
