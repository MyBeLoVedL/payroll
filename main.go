package main

import (
	"context"
	"database/sql"
	"fmt"

	db "payroll/db"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ctx := context.Background()

	dbIns, err := sql.Open("mysql", "root:Lp262783@/payroll")
	if err != nil {
		panic(err)
	}

	queries := db.New(dbIns)

	// create an author
	err = queries.CreateEmployee(ctx, db.CreateEmployeeParams{
		Type:                  "hour",
		Mail:                  "17752254783@163.com",
		SocialSecurityNumber:  "111",
		StandardTaxDeductions: "0",
		OtherDuductions:       "0",
		PhoneNumber:           "17752254783",
		Rate:                  "0.15",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Word Done\n")

}
