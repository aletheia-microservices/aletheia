package constraints

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"analyzer/pkg/app"
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

func ParseConstraints(app *app.App) {
	fmt.Printf("\n\nSetting up unicity constraints for available schema:\n\n")
	for _, dbInstance := range app.Databases {
		for _, unfoldedField := range dbInstance.GetDatastore().Schema.UnfoldedFields {
			fmt.Println("- " + unfoldedField.GetFullName())
		}
		fmt.Println()

	}

	var input string
	var err error
	var targetFields []string
	var targetDbPaths []string

	input = ""
	fmt.Printf("\nPlease specify path(s) to mysql files (.sql) if existent (delimiter is ';', format is <db_name>:<path>):\n> ")

	if app.Name == "coupons_app_sql" {
		input = "coupons_db:blueprint/examples/coupons_app_sql/workflow/coupons_app_sql/database/coupons.sql;students_db:blueprint/examples/coupons_app_sql/workflow/coupons_app_sql/database/students.sql"
	} else if app.Name == "coupons_app" {
		//skip
	} else {
		reader := bufio.NewReader(os.Stdin)
		input, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
	}

	if app.Name != "coupons_app" {
		targetDbPaths = strings.Split(input, ";")
		for _, dbPath := range targetDbPaths {
			splits := strings.Split(dbPath, ":")
			db := splits[0]
			sqlStmt := splits[1]
			sqlBytes, err := os.ReadFile(sqlStmt)
			if err != nil {
				fmt.Println("Error reading database sql files:", err)
				return
			}
			sqlStmts := strings.Split(string(sqlBytes), ";")
			dbInstance := app.GetDatastoreInstance(db)
			for _, stmt := range sqlStmts {
				if stmt == "\n" {
					continue
				}
				parseSQLStatement(app, dbInstance.GetDatastore(), stmt)
			}
		}
	}

	input = ""
	fmt.Printf("\nPlease specify fields to enforce unicity constraint (delimiter is ';', composed uniqueness is within '(...)'):\n> ")

	if app.Name == "coupons_app" {
		input = "(STUDENTS_DB.Student.StudentID);(COUPONS_DB.Coupon.CouponID);(COUPONS_DB.ClaimedCoupon.CouponID,COUPONS_DB.ClaimedCoupon.UserID)"
	} else if app.Name == "coupons_app_sql" {
		//skip
	} else {
		reader := bufio.NewReader(os.Stdin)
		input, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
	}

	if app.Name != "coupons_app_sql" {
		input = strings.TrimSpace(input)
		targetFields = strings.Split(input, ";")
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

		parseUserUniqueConstraints(app, targetFieldsByDatastore)
	}
}
