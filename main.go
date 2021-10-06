package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	db "payroll/db/sqlc"

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
	rows, err := queries.ListEmployees(ctx)
	check(err)
	for _, row := range rows {
		fmt.Printf("%v\n", row)
	}

	fmt.Printf("Word Done\n")

}
