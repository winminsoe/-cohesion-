package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/cohesion/query"
	_ "github.com/go-sql-driver/mysql"
)

const invalidParamMsg = "Please set env variable DB_DSN to a valid MySQL connection string"

type dependentTable struct {
	TableName  string `json:"table"`
	ColumnName string `json:"column"`
}

type mainTable struct {
	TableName       string           `json:"table"`
	DependentTables []dependentTable `json:"dependent_tables"`
}

func main() {
	dbDsn := os.Getenv("DB_DSN")
	if dbDsn == "" {
		fmt.Fprint(os.Stderr, invalidParamMsg)
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
		currentTable := mainTable{TableName: tableName}
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

			var constraintTable, constraintKey string
			for constraintResult.Next() {
				err = constraintResult.Scan(&constraintTable, &constraintKey)

				if err != nil {
					fmt.Println(err.Error())
				}

				currentDependentTable := dependentTable{TableName: constraintTable, ColumnName: constraintKey}
				currentTable.addDependentTable(currentDependentTable)
			}
		}
		result, _ := json.Marshal(currentTable)
		fmt.Println(string(result))
	}
}

func getRows(db *sql.DB, query string) *sql.Rows {
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
	}
	return rows
}

func (mainTable *mainTable) addDependentTable(dependentTable dependentTable) []dependentTable {
	mainTable.DependentTables = append(mainTable.DependentTables, dependentTable)
	return mainTable.DependentTables
}
