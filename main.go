package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/datastores"
	"analyzer/pkg/detection/constraints/cascade"
	"analyzer/pkg/detection/constraints/foreign_key"
	"analyzer/pkg/detection/constraints/specialization"
	"analyzer/pkg/detection/constraints/unicity"
	"analyzer/pkg/detection/constraints/xcy"
	"analyzer/pkg/detection/detector"
	"analyzer/pkg/detection/iterator"
	"analyzer/pkg/frameworks/blueprint"
	"analyzer/pkg/logger"
	"analyzer/pkg/utils"
)

const TEXT_BOLD_LIGHT_RED = "\033[1;31m"
const TEXT_RESET_COLOR = "\033[0m"

type analysisConfig struct {
	allFlag                    string
	appName                    string
	detectOnly                 bool
	xcyDetection               bool
	foreignKeyDetection        bool
	cascadeDetection           bool
	unicityIndividualDetection bool
	unicityAggregateDetection  bool
	specializationDetection    bool
}

func main() {
	allFlag := flag.String("all", "", fmt.Sprintf("Run analyzer for all applications: %v", utils.Apps))
	appName := flag.String("app", "", fmt.Sprintf("The name of the application to be analyzed: %v", utils.Apps))
	detectOnly := flag.Bool("detect_only", false, "Only perform detection (assume parsing is already done)")
	xcyDetection := flag.Bool("xcy", false, "Enable detection of xcy dependencies and inconsistencies")
	foreignKeyDetection := flag.Bool("fk", false, "Enable detection of anomalies in foreign key constraints")
	cascadeDetection := flag.Bool("cascade", false, "Enable detection of the absence of cascading delete logic")
	unicityIndividualDetection := flag.Bool("unicity_individual", false, "Enable detection of inconsistencies for unicity constraints (individual)")
	unicityAggregateDetection := flag.Bool("unicity_aggregate", false, "Enable detection of inconsistencies for unicity constraints (aggregate)")
	specializationDetection := flag.Bool("specialization", false, "Enable detection of removals in mandatory specializations")

	flag.Parse()
	analysis := analysisConfig{
		allFlag:                    *allFlag,
		appName:                    *appName,
		detectOnly:                 *detectOnly,
		xcyDetection:               *xcyDetection,
		foreignKeyDetection:        *foreignKeyDetection,
		cascadeDetection:           *cascadeDetection,
		unicityIndividualDetection: *unicityIndividualDetection,
		unicityAggregateDetection:  *unicityAggregateDetection,
		specializationDetection:    *specializationDetection,
	}

	if analysis.allFlag == "true" || analysis.allFlag == "True" || analysis.allFlag == "1" {
		for _, app := range utils.Apps {
			logger.Logger.Infof(fmt.Sprintf("running analyzer for '%s'...", app))
			time.Sleep(1500 * time.Millisecond)
			initAnalyzer(analysis)
			fmt.Println()
			fmt.Println()
		}
		return
	}
	if !slices.Contains(utils.Apps, analysis.appName) {
		logger.Logger.Fatal(fmt.Sprintf("invalid app name (%s) must provide an application name using the -app flag for one of the available applications: %v", analysis.appName, utils.Apps))

	}
	initAnalyzer(analysis)
}

func initAnalyzer(analysis analysisConfig) {
	servicesInfo, databaseInstances, frontends := blueprint.LoadWiring(analysis.appName)

	app, err := app.InitApp(analysis.appName, servicesInfo)
	if err != nil {
		return
	}
	app.RegisterPackages()
	app.RegisterDatastoreInstances(databaseInstances)
	app.RegisterServiceNodes(servicesInfo)
	app.BuildServiceNodes()
	/* app.PreDump()
	logger.Logger.Fatalf("EXIT!") */

	fmt.Println()
	fmt.Println(" -------------------------------------------------------------------------------------------------------------- ")
	fmt.Println(" -------------------------------------------- BUILD ABSTRACT GRAPH -------------------------------------------- ")
	fmt.Println(" -------------------------------------------------------------------------------------------------------------- ")
	fmt.Println()

	abstractGraph := abstractgraph.Build(app, frontends)

	fmt.Println()
	fmt.Println(" ----------------------------------------------------------------------------------------------------------------- ")
	fmt.Println(" -------------------------------------------- BUILD DATASTORES SCHEMA -------------------------------------------- ")
	fmt.Println(" ----------------------------------------------------------------------------------------------------------------- ")
	fmt.Println()

	abstractgraph.BuildSchema(app, frontends, abstractGraph.Nodes)
	app.Dump()
	abstractGraph.Dump()
	fmt.Println()

	summary := "\n\n"

	if analysis.xcyDetection {
		xcyDetectorGroup := xcy.NewDetectorGroup(abstractGraph.Nodes)
		var cumulativeDatastoreOps map[*datastores.Datastore][]*xcy.Operation
		for _, xcyDetector := range xcyDetectorGroup.GetAllDetectors() {
			xcyDetector.InitRequest(cumulativeDatastoreOps)

			iterator := iterator.NewIterator(app, abstractGraph, xcyDetector)
			iterator.Run()
			cumulativeDatastoreOps = xcyDetector.GetDatastoreOps()
		}

		summary += detector.SaveResults(app, xcyDetectorGroup)
	}

	app.ResetAllDataflows()

	if analysis.cascadeDetection {
		cascadeDetector := cascade.NewDetector()
		iterator := iterator.NewIterator(app, abstractGraph, cascadeDetector)
		iterator.Run()
		summary += detector.SaveResults(app, cascadeDetector)
	}

	if analysis.foreignKeyDetection {
		foreignKeyDetector := foreign_key.NewDetector()
		iterator := iterator.NewIterator(app, abstractGraph, foreignKeyDetector)
		iterator.Run()
		summary += detector.SaveResults(app, foreignKeyDetector)

		foreignKeyDetector.CompactSchema(app)
		app.DumpYamlSchema(true)
	}

	if analysis.specializationDetection {
		specializationDetector := specialization.NewDetector()
		iterator := iterator.NewIterator(app, abstractGraph, specializationDetector)
		iterator.Run()
		summary += detector.SaveResults(app, specializationDetector)

		app.DumpYamlSchema(true)
	}

	if analysis.unicityIndividualDetection {
		parseUniqueConstaintsFromUserInput(app)

		unicityDetector := unicity.NewDetector()
		iterator := iterator.NewIterator(app, abstractGraph, unicityDetector)
		iterator.Run()
		summary += detector.SaveResults(app, unicityDetector)
	}

	fmt.Println(summary)
}

func parseUniqueConstaintsFromUserInput(app *app.App) {
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

	if input != "" {
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
				datastores.ParseSQLStatement(dbInstance.GetDatastore(), stmt)
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
