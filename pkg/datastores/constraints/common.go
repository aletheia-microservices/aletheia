package constraints

import (
	"fmt"
	"os"
	"strings"

	"analyzer/pkg/app"
	"analyzer/pkg/logger"
	"analyzer/pkg/utils"
)

const TEXT_BOLD_LIGHT_RED = "\033[1;31m"
const TEXT_RESET_COLOR = "\033[0m"

func summarize(app *app.App, prefix string) {
	for _, db := range app.GetDbInstances() {
		schema := db.GetDatastore().Schema
		fmt.Printf("\n%s[%s] The following unicity constraints were added:\n", TEXT_BOLD_LIGHT_RED, prefix)

		for _, uc := range schema.GetAllConstraints() {
			fmt.Println("- " + uc.String())
		}
		fmt.Print(TEXT_RESET_COLOR)
	}
}

type DbStmts struct {
	db   string
	stmt string
}

func ParseConstraints(app *app.App) {
	for _, dbInstance := range app.Databases {
		for _, unfoldedField := range dbInstance.GetDatastore().Schema.UnfoldedFields {
			fmt.Println("- " + unfoldedField.GetFullName())
		}
		fmt.Println()

	}

	if ok, input := utils.GetAppDatabaseSQLPaths(app.Name); ok {
		dbStmts := parseAppDatabaseSQLPaths(input)
		for _, dbStmt := range dbStmts {
			dbInstance := app.GetDatastoreInstance(dbStmt.db)
			parseSQLStatement(dbInstance.GetDatastore(), dbStmt.stmt)
		}
	}

	if ok, input := utils.GetAppDatabaseSQLUserInput(app.Name); ok {
		targetFieldsByDatastore := parseAppDatabaseSQLUserInput(input)
		parseUserUniqueConstraints(app, targetFieldsByDatastore)
	}
}

func parseAppDatabaseSQLPaths(input string) []DbStmts {
	var dbStmts []DbStmts
	targetDbPaths := strings.Split(input, ";")
	for _, dbPath := range targetDbPaths {
		splits := strings.Split(dbPath, ":")
		db := splits[0]
		sqlStmt := splits[1]
		sqlBytes, err := os.ReadFile(sqlStmt)
		if err != nil {
			logger.Logger.Fatalf("error reading sql files: %s", err.Error())
			return nil
		}
		sqlStmts := strings.Split(string(sqlBytes), ";")
		for _, stmt := range sqlStmts {
			if stmt == "\n" {
				continue
			}
			dbStmts = append(dbStmts, DbStmts{db, stmt})
		}
	}
	return dbStmts
}

func parseAppDatabaseSQLUserInput(input string) map[string][]string {
	input = strings.TrimSpace(input)
	targetFields := strings.Split(input, ";")
	targetFieldsByDatastore := make(map[string][]string)
	fmt.Printf("\n%s[WARNING] Unicity constraint will be added to each of the following fields:\n", TEXT_BOLD_LIGHT_RED)

	if len(targetFields) > 0 && strings.TrimSpace(targetFields[0]) != "" {
		for _, targetField := range targetFields {
			// remove parentheses
			targetField = targetField[1 : len(targetField)-1]

			splits := strings.SplitN(targetField, ".", 2)
			dbName := splits[0]
			fieldName := targetField
			targetFieldsByDatastore[dbName] = append(targetFieldsByDatastore[dbName], fieldName)

			fmt.Println("- " + targetField)
		}
		fmt.Println(TEXT_RESET_COLOR)
	}
	return targetFieldsByDatastore
}
