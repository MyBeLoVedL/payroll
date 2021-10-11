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
UPDATE employees SET type = ?,mail = ?,social_security_number=?,
	standard_tax_deductions=?,other_deductions=?,phone_number = ?,
	salary_rate=?,hour_limit=? where id = ?;


-- name: SelectActiveTimecard :one
SELECT * FROM timecard WHERE emp_id = ?;

-- name: AddTimecard :execresult
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



-- name: GetUser :one
SELECT * FROM employees WHERE id = ?;


-- name: AddOrderInfo :exec
INSERT INTO order_info(order_id,product_id,amount) 
  VALUES(?,?,?); 


-- name: AddPurchaseOrder :execresult
INSERT INTO purchase_order(emp_id,customer_contact,customer_address,date)
VALUES(?,?,?,?);


-- name: UpdateOrderInfo :exec
UPDATE order_info SET product_id = ?,amount = ?  where order_id = ?;


-- name: UpdatePurchaseOrder :exec
UPDATE purchase_order SET customer_contact = ?,customer_address =? , date = ? WHERE id = ?;


-- name: SelectOrderById :one
select * from purchase_order where id = ?;

-- name: SelectOrderInfoById :one
select * from order_info where order_id = ?;


-- name: DeletePurchaseOrderById :exec
DELETE FROM purchase_order WHERE id = ?;

-- name: DeleteOrderInfoById :exec
DELETE FROM order_info WHERE order_id = ?;




