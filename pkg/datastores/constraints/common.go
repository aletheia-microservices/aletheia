package constraints

import (
	"fmt"

	"analyzer/pkg/app"
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

func ParseConstraints(app *app.App, autofill bool) {
	for _, dbInstance := range app.Databases {
		for _, unfoldedField := range dbInstance.GetDatastore().Schema.UnfoldedFields {
			fmt.Println("- " + unfoldedField.GetFullName())
		}
		fmt.Println()

	}

	if ok, input := utils.GetAppDatabaseDocPaths(app.Name, autofill); ok {
		jsonSchemas := parseAppDatabaseDocSchemas(input)
		for _, jsonSchema := range jsonSchemas {
			dbInstance := app.GetDatastoreInstance(jsonSchema.Schema.Database)
			parseDocJSON(dbInstance.GetDatastore(), jsonSchema)
		}
	}

	if ok, input := utils.GetAppDatabaseSQLPaths(app.Name, autofill); ok {
		dbStmts := parseAppDatabaseSQLStmts(input)
		for _, dbStmt := range dbStmts {
			dbInstance := app.GetDatastoreInstance(dbStmt.db)
			parseSQLStatement(dbInstance.GetDatastore(), dbStmt.stmt)
		}
	}

	if ok, input := utils.GetAppDatabaseUnicityConstraintFromUserInput(app.Name, autofill); ok {
		targetFieldsByDatastore := parseAppDatabaseSQLUserInput(input)
		parseUserUniqueConstraints(app, targetFieldsByDatastore)
	}
}
