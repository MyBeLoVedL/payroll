up:
	migrate -path db/migration/ -database "mysql://root:Lp262783@tcp(127.0.0.1:3306)/payroll" up

down:
	migrate -path db/migration/ -database "mysql://root:Lp262783@tcp(127.0.0.1:3306)/payroll" down

create:
	migrate create -dir db/migration/ -ext sql -seq  init_schema
