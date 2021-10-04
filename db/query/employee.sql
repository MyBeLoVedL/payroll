-- name: CreateEmployee :execresult
INSERT INTO
  employees (
    type,
    mail,
    social_security_number,
    standard_tax_deductions,
    other_duductions,
    phone_number,
    rate
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