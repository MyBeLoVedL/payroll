package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	db "payroll/db"

	_ "github.com/go-sql-driver/mysql"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	ctx := context.Background()

	dbIns, err := sql.Open("mysql", "root:Lp262783@/payroll")
	if err != nil {
		panic(err)
	}

	queries := db.New(dbIns)

	// create an author
	queries.CreateEmployee(ctx, db.CreateEmployeeParams{
		Type:                  "salaried",
		Mail:                  "17752254783@gmail.com",
		SocialSecurityNumber:  "222",
		StandardTaxDeductions: "0.8",
		OtherDuductions:       "0",
		PhoneNumber:           "17752254783",
		Rate:                  "0.25",
	})
	rows, err := queries.ListEmployees(ctx)
	check(err)
	for _, row := range rows {
		fmt.Printf("%v\n", row)
	}

	fmt.Printf("Word Done\n")

}
