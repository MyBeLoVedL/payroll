// Code generated by sqlc. DO NOT EDIT.
// source: employee.sql

package db

import (
	"context"
)

const createEmployee = `-- name: CreateEmployee :exec
INSERT INTO employees (
    type,
    mail,
    social_security_number,
    standard_tax_deductions,
    other_duductions,
    phone_number,
    rate
  )
VALUES (?,?,?,?,?,?,?)
`

type CreateEmployeeParams struct {
	Type                  EmployeesType `json:"type"`
	Mail                  string        `json:"mail"`
	SocialSecurityNumber  string        `json:"social_security_number"`
	StandardTaxDeductions string        `json:"standard_tax_deductions"`
	OtherDuductions       string        `json:"other_duductions"`
	PhoneNumber           string        `json:"phone_number"`
	Rate                  string        `json:"rate"`
}

func (q *Queries) CreateEmployee(ctx context.Context, arg CreateEmployeeParams) error {
	_, err := q.db.ExecContext(ctx, createEmployee,
		arg.Type,
		arg.Mail,
		arg.SocialSecurityNumber,
		arg.StandardTaxDeductions,
		arg.OtherDuductions,
		arg.PhoneNumber,
		arg.Rate,
	)
	return err
}
