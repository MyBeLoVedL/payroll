-- name: AddEmployee :execresult
INSERT INTO
  employees (
    type,
    mail,
    social_security_number,
    standard_tax_deductions,
    other_deductions,
    phone_number,
    salary_rate
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?);
-- name: ListEmployees :many
SELECT
  *
FROM
  employees
ORDER BY
  id;


-- name: SelectEmployeeById :one
SELECT
  *
from
  employees
where
  id = ?
LIMIT
  1;


-- name: SelectEmployeeByMail :one
SELECT * FROM employees WHERE mail = ?;



-- name: UpdatePassword :exec
UPDATE employees SET password = ?
WHERE id = ?;



-- name: DeleteEmployee :exec
UPDATE employees SET deleted = 1 where id = ?;


-- name: UpdateEmployee :exec
UPDATE employees SET type = ?,mail = ?,social_security_number=?,standard_tax_deductions=?,other_deductions=?,phone_number = ?,salary_rate=?,hour_limit=? where id = ?;


-- name: AddTimecard :exec
INSERT INTO timecard(emp_id) VALUES (?);


-- name: AddTimecardRecord :exec
INSERT INTO timecard_record(charge_number,card_id,hours,date) VALUES (?,?,?,?);


-- name: UpdatePaymentMethod :exec
UPDATE employees SET payment_method = ? where id = ?;


-- name: UpdatePaymentMethodWithMail :exec
UPDATE employees SET payment_method = ?,mail = ?  where id = ?;


-- name: InsertBank :exec
INSERT INTO employee_account(id,bank_name,account_number)
	VALUES (?,?,?);

