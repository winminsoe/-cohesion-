package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
)

func main() {
	dbDsn := os.Getenv("DB_DSN")
	if dbDsn == "" {
		fmt.Fprint(os.Stderr, "Please set env variable DB_DSN to a valid MySQL connection string")
		os.Exit(1)
	}
	temp := strings.SplitAfter(dbDsn,"/")
	dbName := temp[1]

	db, err := sql.Open("mysql", dbDsn)
	if err != nil {
    	fmt.Println(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	tables, err := db.Query("show tables")
	if err != nil {
		fmt.Println(err.Error())
	}

	var tableName string
	for tables.Next() {
		err = tables.Scan(&tableName)

		if err != nil {
			fmt.Println(err.Error()) 
		}

		fmt.Println("## `", tableName, "`");
		fmt.Println("\n");

		getAutoIncrementStatement := `SELECT COLUMN_NAME
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_NAME = '%s'
		AND TABLE_SCHEMA = '%s'
		AND EXTRA like 'auto_increment'
		`
		getAutoIncrementQuery := fmt.Sprintf(getAutoIncrementStatement, tableName, dbName)
    	autoIncrement, err := db.Query(getAutoIncrementQuery)
		if err != nil {
			fmt.Println(err.Error())
		}

		var autoIncrementColumn string
		for autoIncrement.Next() {
			err = autoIncrement.Scan(&autoIncrementColumn)

			if err != nil {
				fmt.Println(err.Error()) 
			}
			getConstriantStatement := `SELECT TABLE_NAME, COLUMN_NAME
			FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
			WHERE
			REFERENCED_TABLE_SCHEMA = '%s'
			AND REFERENCED_COLUMN_NAME = '%s'
			AND REFERENCED_TABLE_NAME = '%s'
			`
			getConstriantQuery := fmt.Sprintf(getConstriantStatement, dbName, autoIncrementColumn, tableName)
			constriantResult, err := db.Query(getConstriantQuery)

			if err != nil {
				fmt.Println(err.Error())
			}
			
			var constriantTable, constraintKey []byte
			fmt.Println("| Table | Column |")
			fmt.Println("| ----- | ------ |")
			for constriantResult.Next() {
				err = constriantResult.Scan(&constriantTable, &constraintKey)

				if err != nil {
					fmt.Println(err.Error()) 
				}

				fmt.Println("| `",string(constriantTable), "` | `", string(constraintKey), "` |")
			}
		}
		fmt.Println("\n");
	}
}