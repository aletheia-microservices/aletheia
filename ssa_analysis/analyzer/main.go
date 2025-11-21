package main

import (
	"flag"
	"fmt"
	"go/token"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"gopkg.in/yaml.v2"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/detection"
	"analyzer/pkg/detection/constraints/foreignkeycascade"
	"analyzer/pkg/detection/constraints/foreignkeyconcurrency"
	"analyzer/pkg/detection/constraints/keycoordination"
	"analyzer/pkg/detection/constraints/unicityconcurrency"
	"analyzer/pkg/ssagraph"
	"analyzer/pkg/ssagraph/parser"
	"analyzer/pkg/ssagraph/registry"
	"analyzer/pkg/ssagraph/tainter"
	"analyzer/pkg/utils"
)

var EVAL = true

func main() {
	flag.BoolVar(&EVAL, "eval", false, "enable evaluation mode")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "usage: program [--eval] <appname>\n")
		fmt.Fprintln(os.Stderr, "available appnames:")
		fmt.Fprintln(os.Stderr, "- foobar")
		fmt.Fprintln(os.Stderr, "- postnotification_simple")
		fmt.Fprintln(os.Stderr, "- dsb_sn2")
		fmt.Fprintln(os.Stderr, "- digota")
		fmt.Fprintln(os.Stderr, "- sockshop3")
		fmt.Fprintln(os.Stderr, "- dsb_media_sql")
		fmt.Fprintln(os.Stderr, "- dsb_media_nosql")
		fmt.Fprintln(os.Stderr, "- dsb_hotel2")
		fmt.Fprintln(os.Stderr, "- train_ticket2")
		fmt.Fprintln(os.Stderr, "- large_scale_app")
		os.Exit(1)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.TimeOnly,
	})

	if EVAL {
		go func() {
			for {
				for _, r := range `-\|/` {
					fmt.Printf("\rRunning... %c", r)
					time.Sleep(1 * time.Second)
				}
			}
		}()
	}

	appname := flag.Arg(0)
	
	start := time.Now()
	
	logrus_ctx := logrus.WithField("app", appname)

	// ------------ PART 1
	logrus_ctx.Infof("[1/13] initializing program")

	apppath := utils.GetAppEntrypointPath(appname)
	app := app.NewApp(appname)
	app.Init()


	// ensure output sub directory exists
	err := os.MkdirAll(fmt.Sprintf("output/%s", appname), os.ModePerm)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	// ensure output sub directory for graphs exists
	err = os.MkdirAll(fmt.Sprintf("output/%s/ssagraphs/tainted", appname), os.ModePerm)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}
	err = os.MkdirAll(fmt.Sprintf("output/%s/ssagraphs/untainted", appname), os.ModePerm)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	// ensure output sub directory for graphs exists
	err = os.MkdirAll(fmt.Sprintf("output/%s/ssa", appname), os.ModePerm)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	// ------------ PART 2
	logrus_ctx.Infof("[2/13] building program")

	prog, pkgs, err := buildProgram(apppath)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	app.InitServiceFields(pkgs)
	app.ParseSQLSchemaFromUserFile()
	app.ParseNoSQLSchemaFromUserFile()

	// ------------ PART 3
	logrus_ctx.Infof("[3/13] initializing pointer analysis")
	result, err := parser.InitPointerAnalysis(prog, pkgs)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	funcGraphs := make(map[string]*ssagraph.SSAGraph)

	// ------------ PART 4
	logrus_ctx.Infof("[4/13] running SSA analysis")
	for _, pkg := range pkgs {
		parser.RunSSAAnalysis(app, prog, pkg, funcGraphs)
	}

	for _, ssagraph := range funcGraphs {
		ssagraph.Sort()
	}

	// ------------ PART 5
	logrus_ctx.Infof("[5/13] running pointer-to analysis")
	for _, pkg := range pkgs {
		parser.RunPointerToAnalysis(appname, prog, pkg, result, funcGraphs)
	}

	var graphsLst []*ssagraph.SSAGraph
	for _, graph := range funcGraphs {
		graphsLst = append(graphsLst, graph)
	}

	// ------------ PART 6
	logrus_ctx.Infof("[6/13] registering fields")
	registry.RegisterFields(app, graphsLst)

	if !EVAL {
		app.WriteAppToJSON()
		for fn, ssagraph := range funcGraphs {
			ssagraph.WriteToDOTFile(appname, fn, false)
		}
	}

	// ------------ PART 7
	logrus_ctx.Infof("[7/13] running SSA tainter")
	for _, ssagraph := range funcGraphs {
		tainter.RunTainter(ssagraph)
	}

	if !EVAL {
		for fn, ssagraph := range funcGraphs {
			ssagraph.WriteToDOTFile(appname, fn, true)
		}
	}

	// ------------ PART 8
	logrus_ctx.Infof("[8/13] creating new abstract call graph")
	absgraph := abstractgraph.NewAbstractCallGraph(app)
	for _, entrypoint := range app.GetEntrypointsShortPaths() {
		abstractgraph.Parse(absgraph, entrypoint, true, funcGraphs)
	}

	// ------------ PART 9
	logrus_ctx.Infof("[9/13] releasing memory associated with ssa graph")
	graphsLst = nil
	pkgs = nil
	for fn, ssagraph := range funcGraphs {
		if ssagraph != nil {
			ssagraph.Release()
		}
		delete(funcGraphs, fn)
	}
	funcGraphs = nil

	if !EVAL {
		absgraph.WriteVisited(appname)
	}

	elapsed_parsing := time.Since(start)

	detector1 := keycoordination.NewDetector(keycoordination.DETECTION_TYPE_PRIMARY_KEY)
	detector2 := keycoordination.NewDetector(keycoordination.DETECTION_TYPE_FOREIGN_KEY)
	detector3 := foreignkeycascade.NewDetector()
	detector4 := foreignkeyconcurrency.NewDetector()
	detector5 := unicityconcurrency.NewDetector()
	iterator := detection.NewIterator(app, absgraph, detector1, detector2, detector3, detector4, detector5)

	start_schema := time.Now()
	// ------------ PART 10
	logrus_ctx.Infof("[10/13] starting schema builder")
	iterator.Run(detection.PHASE_1_SCHEMA_BUILDER)
	// ------------ PART 11
	logrus_ctx.Infof("[11/13] starting schema builder (read only)")
	iterator.Run(detection.PHASE_1_SCHEMA_BUILDER_READ_ONLY)
	
	elapsed_schema := time.Since(start_schema)
	
	// ------------ PART 12
	logrus_ctx.Infof("[12/13] starting pattern detection")
	start_detection := time.Now()
	
	// phase 2: one pass for all detectors
	iterator.Run(detection.PHASE_2_PATTERN_DETECTOR)
	
	if !EVAL {
		// phase 0: dummy pass to generate dot files with taints for debugging
		iterator.Run(detection.PHASE_0_DEBUG)
		
		absgraph.WriteToDOTFile(appname, true)
		absgraph.WriteToDOTFile(appname, false)
		
		app.WriteAppToJSON()
		app.WriteSchemaToJSON()
	}
	
	elapsed_total := time.Since(start)
	elapsed_detection := time.Since(start_detection)
	
	// ------------ PART 13
	logrus_ctx.Infof("[13/13] saving results")
	results := detection.SaveResults(app, detector1, detector2, detector3, detector4, detector5)
	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Printf("Execution time (TOTAL):\t\t%.4f s\n", elapsed_total.Seconds())
	fmt.Printf("Execution time (PARSING):\t%.4f s\n", elapsed_parsing.Seconds())
	fmt.Printf("Execution time (SCHEMA):\t%.4f s\n", elapsed_schema.Seconds())
	fmt.Printf("Execution time (DETECTION):\t%.4f s\n", elapsed_detection.Seconds())

	if EVAL {
		times := AnalysisTimes{
			App:              app.GetName(),
			NumMicroservices: app.NumberOfMicroservices(),
			NumDatastores:    app.NumberOfDatastores(),
			Total:            elapsed_total.Seconds(),
			Parsing:          elapsed_parsing.Seconds(),
			Schema:           elapsed_schema.Seconds(),
			Detection:        elapsed_detection.Seconds(),
		}
		saveAnalysisTime(app, times)
	}
}

type AnalysisTimes struct {
	App              string  `yaml:"app"`
	NumMicroservices int     `yaml:"ms_count"`
	NumDatastores    int     `yaml:"ds_count"`
	Total            float64 `yaml:"total_s"`
	Parsing          float64 `yaml:"parsing_s"`
	Schema           float64 `yaml:"schema_s"`
	Detection        float64 `yaml:"detection_s"`
}

func saveAnalysisTime(app *app.App, times AnalysisTimes) {
	ts := time.Now().Format("2006-01-02_15-04-05")
	dir := "analysis_times"
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}

	filepath := fmt.Sprintf("%s/%s_%s.yaml", dir, app.GetName(), ts)

	out, err := yaml.Marshal(times)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(filepath, out, 0644)
	if err != nil {
		panic(err)
	}
}

func buildProgram(apppath string) (*ssa.Program, []*ssa.Package, error) {
	// e.graph. "../apps/test2/main.go"
	filepath := "apps/" + apppath + "/main.go"
	source, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read file %s: %v\n", filepath, err)
		os.Exit(1)
	}

	// setup loader
	var conf loader.Config
	fset := token.NewFileSet()
	conf.Fset = fset
	file, err := conf.ParseFile(filepath, string(source))
	if err != nil {
		return nil, nil, err
	}
	conf.CreateFromFiles("main", file)

	iprog, err := conf.Load()
	if err != nil {
		return nil, nil, err
	}

	prog := ssautil.CreateProgram(iprog, 0)
	mainPkg := prog.Package(iprog.Created[0].Pkg)

	prog.Build()

	var pkgs = []*ssa.Package{mainPkg}

	for _, pkg := range prog.AllPackages() {
		if pkg.Pkg.Path() != "main" { // skip the synthetic main if needed
			if utils.IsAppPackagePath(pkg.Pkg.Path()) {
				pkgs = append(pkgs, pkg)
			}
		}
	}
	return prog, pkgs, nil
}
