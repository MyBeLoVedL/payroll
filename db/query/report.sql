


-- name: GetHoursByEmpID :one
select sum(hours)  from timecard_record 
	where card_id = 
		( select id from timecard where emp_id = 1);



-- name: GetIDByName :many
select id from employees where name = ?; 






