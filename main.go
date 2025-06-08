package main

import (
	"flag"
	"fmt"
	"slices"
	"time"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/datastores"
	"analyzer/pkg/datastores/constraints"
	"analyzer/pkg/detection/constraints/foreign_key_cascade"
	"analyzer/pkg/detection/constraints/foreign_key_concurrency"
	"analyzer/pkg/detection/constraints/key_coordination"
	"analyzer/pkg/detection/constraints/numerical"
	"analyzer/pkg/detection/constraints/specialization"
	"analyzer/pkg/detection/constraints/unicity"
	"analyzer/pkg/detection/constraints/xcy"
	"analyzer/pkg/detection/detection"
	"analyzer/pkg/detection/iterator"
	"analyzer/pkg/frameworks/blueprint"
	"analyzer/pkg/logger"
	"analyzer/pkg/utils"
)

const TEXT_BOLD_LIGHT_YELLOW = "\033[1;38;5;179m"
const TEXT_BOLD_LIGHT_RED = "\033[1;31m"
const TEXT_RESET_COLOR = "\033[0m"
const TEXT_BOLD_LIGHT_BLUE = "\033[1;38;5;75m"
const TEXT_BOLD_LIGHT_GREEN = "\033[1;32m"

type analysisConfig struct {
	allFlag                         string
	appName                         string
	autofill                        bool
	detectOnly                      bool
	xcyDetection                    bool
	primaryKeyCoordinationDetection  bool
	foreignKeyCoordinationDetection bool
	foreignKeyConcurrencyDetection  bool
	foreignKeyCascadeDetection      bool
	compactSchema                   bool
	unicityDetection                bool
	numericalDetection              bool
	specializationDetection         bool
}

func main() {
	allFlag := flag.String("all", "", fmt.Sprintf("Run analyzer for all applications: %v", utils.Apps))
	appName := flag.String("app", "", fmt.Sprintf("The name of the application to be analyzed: %v", utils.Apps))
	autofill := flag.Bool("auto", false, "Autofills additional user input information")
	detectOnly := flag.Bool("detect_only", false, "Only perform detection (assume parsing is already done)")
	xcyDetection := flag.Bool("xcy", false, "Enable detection of xcy dependencies and inconsistencies")
	primaryKeyCoordinationDetection := flag.Bool("pk_coordination", false, "Enable detection of anomalies for uncoordinated reads in primary key constraints")
	foreignKeyCoordinationDetection := flag.Bool("fk_coordination", false, "Enable detection of anomalies for uncoordinated reads in foreign key constraints")
	foreignKeyConcurrencyDetection := flag.Bool("fk_concurrency", false, "Enable detection of concurrency anomalies in foreign key constraints")
	foreignKeyCascadeDetection := flag.Bool("fk_cascade", false, "Enable detection of the absence of cascading delete logic")
	compactSchema := flag.Bool("compact_schema", false, "Enable schema compaction (only available when `fk_coordination` is also enabled)")
	unicityDetection := flag.Bool("unicity", false, "Enable detection of inconsistencies for unicity constraints")
	numericalDetection := flag.Bool("numerical", false, "Enable detection of inconsistencies for numerical constraints")
	specializationDetection := flag.Bool("specialization", false, "Enable detection of removals in mandatory specializations")

	flag.Parse()
	analysis := analysisConfig{
		allFlag:                         *allFlag,
		appName:                         *appName,
		autofill:                        *autofill,
		detectOnly:                      *detectOnly,
		xcyDetection:                    *xcyDetection,
		primaryKeyCoordinationDetection: *primaryKeyCoordinationDetection,
		foreignKeyCoordinationDetection: *foreignKeyCoordinationDetection,
		foreignKeyConcurrencyDetection:  *foreignKeyConcurrencyDetection,
		foreignKeyCascadeDetection:      *foreignKeyCascadeDetection,
		compactSchema:                   *compactSchema,
		unicityDetection:                *unicityDetection,
		numericalDetection:              *numericalDetection,
		specializationDetection:         *specializationDetection,
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
		logger.Logger.Fatalf("error initializing app: %s", err.Error())
	}
	app.RegisterPackages()
	app.RegisterDatastores(databaseInstances)
	app.RegisterServices(servicesInfo)
	app.BuildServices()
	app.ParseMethods()
	app.Dump()

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

	prepSchema(analysis, app, abstractGraph)

	abstractgraph.BuildSchema(app, frontends, abstractGraph.Nodes)
	app.Dump()
	abstractGraph.Dump()
	fmt.Println()

	prepAnalysis(analysis, app)
	runAnalysis(analysis, app, abstractGraph)
	endAnalysis(analysis, app)
}

func prepSchema(analysis analysisConfig, app *app.App, abstractGraph *abstractgraph.AbstractGraph) {
	abstractGraph.AttachDatabaseFieldsToEntryArgs(app, analysis.autofill)
}

func prepAnalysis(analysis analysisConfig, app *app.App) {
	constraints.ParseConstraints(app, analysis.autofill)
	app.DumpYamlSchema(false)
}

func runAnalysis(analysis analysisConfig, app *app.App, abstractGraph *abstractgraph.AbstractGraph) {
	var results, summary string

	if analysis.xcyDetection {
		xcyDetectorGroup := xcy.NewDetectorGroup(abstractGraph.Nodes)
		var cumulativeDatastoreOps map[*datastores.Datastore][]*xcy.Operation
		for _, xcyDetector := range xcyDetectorGroup.GetAllDetectors() {
			xcyDetector.InitRequest(cumulativeDatastoreOps)

			iterator := iterator.NewIterator(app, abstractGraph, xcyDetector)
			iterator.Run()
			cumulativeDatastoreOps = xcyDetector.GetDatastoreOps()
		}

		results += detection.SaveResults(app, xcyDetectorGroup)
		summary += xcyDetectorGroup.GetSummary()
	}

	app.ResetAllDataflows()

	if analysis.primaryKeyCoordinationDetection {
		keyCoordinationDetector := key_coordination.NewDetector("primary_key")
		iterator := iterator.NewIterator(app, abstractGraph, keyCoordinationDetector)
		iterator.Run()

		results += detection.SaveResults(app, keyCoordinationDetector)
		summary += keyCoordinationDetector.GetSummary()
	}

	if analysis.foreignKeyCoordinationDetection {
		keyCoordinationDetector := key_coordination.NewDetector("foreign_key")
		iterator := iterator.NewIterator(app, abstractGraph, keyCoordinationDetector)
		iterator.Run()

		results += detection.SaveResults(app, keyCoordinationDetector)
		summary += keyCoordinationDetector.GetSummary()

		if analysis.compactSchema {
			keyCoordinationDetector.CompactSchema(app)
		}
	}

	if analysis.foreignKeyConcurrencyDetection {
		foreignKeyConcurrencyDetector := foreign_key_concurrency.NewDetector()
		iterator := iterator.NewIterator(app, abstractGraph, foreignKeyConcurrencyDetector)
		iterator.Run()
		foreignKeyConcurrencyDetector.NextIterationPhase()
		iterator.Run()
		results += detection.SaveResults(app, foreignKeyConcurrencyDetector)
		summary += foreignKeyConcurrencyDetector.GetSummary()
	}

	if analysis.foreignKeyCascadeDetection {
		foreignKeyCascadeDetector := foreign_key_cascade.NewDetector()
		iterator := iterator.NewIterator(app, abstractGraph, foreignKeyCascadeDetector)
		iterator.Run()
		results += detection.SaveResults(app, foreignKeyCascadeDetector)
		summary += foreignKeyCascadeDetector.GetSummary()
	}

	if analysis.specializationDetection {
		specializationDetector := specialization.NewDetector()
		iterator := iterator.NewIterator(app, abstractGraph, specializationDetector)
		iterator.Run()
		results += detection.SaveResults(app, specializationDetector)
		summary += specializationDetector.GetSummary()
	}

	if analysis.unicityDetection {
		unicityDetector := unicity.NewDetector()
		iterator := iterator.NewIterator(app, abstractGraph, unicityDetector)
		iterator.Run()
		results += detection.SaveResults(app, unicityDetector)
		summary += unicityDetector.GetSummary()
	}

	if analysis.numericalDetection {
		numericalDetector := numerical.NewDetector()
		iterator := iterator.NewIterator(app, abstractGraph, numericalDetector)
		iterator.Run()
		results += detection.SaveResults(app, numericalDetector)
		summary += numericalDetector.GetSummary()
	}

	fmt.Println("\n--------- RESULTS ---------\n" + results)
	fmt.Println("\n--------- SUMMARY ---------\n" + summary)
}

func endAnalysis(analysis analysisConfig, app *app.App) {
	if analysis.foreignKeyCoordinationDetection || analysis.specializationDetection {
		app.DumpYamlSchema(true)
	}
}
