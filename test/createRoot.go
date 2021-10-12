package main

import (
	"context"
	"database/sql"
	"payroll/db/util"
	"testing"
)

type EmployeesType string

type AddEmployeeParams struct {
	Type                  EmployeesType `json:"type"`
	Mail                  string        `json:"mail"`
	SocialSecurityNumber  string        `json:"social_security_number"`
	StandardTaxDeductions string        `json:"standard_tax_deductions"`
	OtherDeductions       string        `json:"other_deductions"`
	PhoneNumber           string        `json:"phone_number"`
	SalaryRate            string        `json:"salary_rate"`
}

var queries *Queries

func TestMain(m *testing.M) {
}

func main() {
	dbIns, err := sql.Open("mysql", "lee:Lp262783@/payroll")
	if err != nil {
		panic(err)
	}
	arg := AddEmployeeParams{
		Type:                  EmployeesType(util.RandType()),
		Mail:                  util.RandDigits(10) + "@gmail.com",
		SocialSecurityNumber:  util.RandDigits(15),
		StandardTaxDeductions: "0.12",
		OtherDeductions:       "10000.0",
		PhoneNumber:           util.RandDigits(11),
		SalaryRate:            "0.25",
	}

	res, err := dbIns.QueryRow()
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
