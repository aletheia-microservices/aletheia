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
	"analyzer/pkg/detection/cascade"
	"analyzer/pkg/detection/foreign_key"
	"analyzer/pkg/detection/specialization"
	"analyzer/pkg/detection/unicity"
	"analyzer/pkg/detection/xcy"
	"analyzer/pkg/frameworks/blueprint"
	"analyzer/pkg/logger"
	"analyzer/pkg/utils"
)

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
		fmt.Println()
		fmt.Println(" -------------------------------------------------------------------------------------------------------------- ")
		fmt.Println(" --------------------------------------- CHECK XCY - TAINTED APPROACH  ---------------------------------------- ")
		fmt.Println(" -------------------------------------------------------------------------------------------------------------- ")
		fmt.Println()

		detectorSet := xcy.NewDetectorSet(app, abstractGraph.Nodes)
		var cumulativeDatastoreOps map[*datastores.Datastore][]*xcy.Operation
		for _, detector := range detectorSet.GetAllDetectors() {
			request := detector.InitRequest(cumulativeDatastoreOps)
			detector.InitXCYRequestTransversal(request)
			cumulativeDatastoreOps = detector.GetDatastoreOps()
		}

		fmt.Println()
		results := detectorSet.Results()
		summary += results + "\n\n"
	}

	app.ResetAllDataflows()

	if analysis.cascadeDetection {
		fmt.Println()
		fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
		fmt.Println(" --------------------------------------- CHECK ABSENCE OF CASCADING DELETE ---------------------------------------- ")
		fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
		fmt.Println()
		detector := cascade.NewDetector(app, abstractGraph)
		detector.Run()
		results := detector.Results()
		summary += results + "\n\n"
	}

	if analysis.foreignKeyDetection {
		fmt.Println()
		fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
		fmt.Println(" ----------------------------------- CHECK INTEGRITY ANOMALIES FOR FOREIGN KEYS ----------------------------------- ")
		fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
		fmt.Println()
		detector := foreign_key.NewDetector(app, abstractGraph)
		detector.Run()
		results := detector.Results()
		summary += results + "\n\n"

		detector.CompactSchema()
		app.DumpYamlSchema(true)
	}

	if analysis.specializationDetection {
		fmt.Println()
		fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
		fmt.Println(" ----------------------------------- CHECK REMOVALS IN MANDATORY SPECIALIZATIONS ---------------------------------- ")
		fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
		fmt.Println()
		detector := specialization.NewDetector(app, abstractGraph)
		detector.Run()
		results := detector.Results()
		summary += results + "\n\n"

		app.DumpYamlSchema(true)
	}

	bold_light_red := "\033[1;31m"
	reset_color := "\033[0m"

	if analysis.unicityIndividualDetection {
		parseUniqueConstaintsFromUserInput(app)

		fmt.Println()
		fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
		fmt.Println(" ----------------------------------- INCONSISTENCY DETECTOR - UNICITY CONSTRAINTS --------------------------------- ")
		fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
		fmt.Println()

		detector := unicity.NewDetector()
		iterator := unicity.NewIterator(app, abstractGraph, detector)
		iterator.Run()

		fmt.Println(bold_light_red)
		res := detector.Results()
		fmt.Println(res)
		fmt.Println(reset_color)
	}

	fmt.Println(bold_light_red + summary)
}

func parseUniqueConstaintsFromUserInput(app *app.App) {
	bold_light_red := "\033[1;31m"
	reset_color := "\033[0m"

	fmt.Printf("\n\nSetting up unicity constraints for available schema:\n\n")
	for _, dbInstance := range app.Databases {
		for _, unfoldedField := range dbInstance.GetDatastore().Schema.UnfoldedFields {
			fmt.Println("- " + unfoldedField.GetFullName())
		}
		fmt.Println()

	}

	fmt.Printf("\nPlease specify fields to enforce unicity constraint (delimiter is ';', composed uniqueness is within '(...)'):\n> ")

	var input string
	var err error
	if app.Name == "coupons_app" {
		input = "(STUDENTS_DB.Student.StudentID);(COUPONS_DB.Coupon.CouponID);(COUPONS_DB.ClaimedCoupon.CouponID,COUPONS_DB.ClaimedCoupon.UserID)"
	} else {
		reader := bufio.NewReader(os.Stdin)
		input, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
	}
	
	input = strings.TrimSpace(input)
	targetFields := strings.Split(input, ";")
	targetFieldsByDatastore := make(map[string][]string)
	fmt.Printf("\n%s[WARNING] Unicity constraint will be added to each of the following fields:\n", bold_light_red)

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
		fmt.Println(reset_color)
	}


	for db, targetFields := range targetFieldsByDatastore {
		dbInstance := app.GetDatastoreInstance(strings.ToLower(db))
		schema := dbInstance.GetDatastore().Schema
		for _, targetField := range targetFields {
			var fields []datastores.Field

			for _, targetFieldSplit := range strings.Split(targetField, ",") {
				field := schema.GetFieldByFullName(targetFieldSplit)
				fields = append(fields, field)
			}

			schema.CreateUniqueConstraint(fields...)
		}
	}

	for _, db := range app.GetDbInstances() {
		schema := db.GetDatastore().Schema
		fmt.Printf("\n%s[WARNING] The following unicity constraints were added:\n", bold_light_red)

		for _, uc := range schema.UniqueConstraints {
			fmt.Println("- " + uc.String())
		}
		fmt.Print(reset_color)
	}
}
