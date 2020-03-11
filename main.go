package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/cohesion/query"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbDsn := os.Getenv("DB_DSN")
	if dbDsn == "" {
		fmt.Fprint(os.Stderr, "Please set env variable DB_DSN to a valid MySQL connection string")
		os.Exit(1)
	}
	temp := strings.SplitAfter(dbDsn, "/")
	dbName := temp[1]

	db, err := sql.Open("mysql", dbDsn)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	tables := getRows(db, query.ShowTableStatement)

	var tableName string
	for tables.Next() {
		err = tables.Scan(&tableName)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("## `", tableName, "`")
		fmt.Println("\n")

		getAutoIncrementQuery := fmt.Sprintf(query.GetAutoIncrementStatement, tableName, dbName)
		autoIncrement := getRows(db, getAutoIncrementQuery)

		var autoIncrementColumn string
		for autoIncrement.Next() {
			err = autoIncrement.Scan(&autoIncrementColumn)

			if err != nil {
				fmt.Println(err.Error())
			}

			getConstraintQuery := fmt.Sprintf(query.GetConstraintStatement, dbName, autoIncrementColumn, tableName)
			constraintResult := getRows(db, getConstraintQuery)

			var constraintTable, constraintKey []byte
			fmt.Println("| Table | Column |")
			fmt.Println("| ----- | ------ |")
			for constraintResult.Next() {
				err = constraintResult.Scan(&constraintTable, &constraintKey)

				if err != nil {
					fmt.Println(err.Error())
				}

				fmt.Println("| `", string(constraintTable), "` | `", string(constraintKey), "` |")
			}
		}
		fmt.Println("\n")
	}
}

func getRows(db *sql.DB, query string) *sql.Rows {
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
	}
	return rows
}
