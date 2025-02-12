package constraints

import (
	"fmt"
	"strings"

	"analyzer/pkg/app"
	"analyzer/pkg/datastores"
)

const TEXT_BOLD_LIGHT_RED = "\033[1;31m"
const TEXT_RESET_COLOR = "\033[0m"

func ParseUserUniqueConstraints(app *app.App, targetFieldsByDatastore map[string][]string) {
	fmt.Printf("\n\nSetting up unicity constraints for available schema:\n\n")
	for _, dbInstance := range app.Databases {
		for _, unfoldedField := range dbInstance.GetDatastore().Schema.UnfoldedFields {
			fmt.Println("- " + unfoldedField.GetFullName())
		}
		fmt.Println()

	}

	for db, targetFields := range targetFieldsByDatastore {
		dbInstance := app.GetDatastoreInstance(strings.ToLower(db))
		schema := dbInstance.GetDatastore().Schema
		for _, targetField := range targetFields {
			var fields []*datastores.Field

			for _, targetFieldSplit := range strings.Split(targetField, ",") {
				field := schema.GetFieldByFullName(targetFieldSplit)
				fields = append(fields, field)
			}

			constraint := datastores.NewConstraintUnique(fields...)
			schema.AddConstraint(constraint)
			for _, field := range fields {
				field.AddConstraint(constraint)
			}
		}
	}

	for _, db := range app.GetDbInstances() {
		schema := db.GetDatastore().Schema
		fmt.Printf("\n%s[WARNING] The following unicity constraints were added:\n", TEXT_BOLD_LIGHT_RED)

		for _, uc := range schema.GetConstraints() {
			fmt.Println("- " + uc.String())
		}
		fmt.Print(TEXT_RESET_COLOR)
	}
}
