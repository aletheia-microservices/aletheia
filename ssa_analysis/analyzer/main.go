package main

import (
	"flag"
	"fmt"
	"go/token"
	"log"
	"os"
	"time"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"

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

	start := time.Now()
	
	fmt.Printf("[%s] [1/13] building program\n", time.Now().Format(time.TimeOnly))
	prog, pkgs, err := buildProgram(apppath)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}
	
	fmt.Printf("[%s] [2/13] initializing service fields\n", time.Now().Format(time.TimeOnly))
	app.InitServiceFields(pkgs)
	fmt.Printf("[%s] [3/13] parsing SQL schema from user file\n", time.Now().Format(time.TimeOnly))
	app.ParseSQLSchemaFromUserFile()
	fmt.Printf("[%s] [4/13] parsing NoSQL schema from user file\n", time.Now().Format(time.TimeOnly))
	app.ParseNoSQLSchemaFromUserFile()

	fmt.Printf("[%s] [5/13] initializing pointer analysis\n", time.Now().Format(time.TimeOnly))
	result, err := parser.InitPointerAnalysis(prog, pkgs)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	funcGraphs := make(map[string]*ssagraph.SSAGraph)

	fmt.Printf("[%s] [6/13] running SSA analysis\n", time.Now().Format(time.TimeOnly))
	for _, pkg := range pkgs {
		parser.RunSSAAnalysis(app, prog, pkg, funcGraphs)
	}

	for _, ssagraph := range funcGraphs {
		ssagraph.Sort()
	}

	fmt.Printf("[%s] [7/13] initializing pointer-to analysis\n", time.Now().Format(time.TimeOnly))
	for _, pkg := range pkgs {
		parser.RunPointerToAnalysis(appname, prog, pkg, result, funcGraphs)
	}

	var graphsLst []*ssagraph.SSAGraph
	for _, graph := range funcGraphs {
		graphsLst = append(graphsLst, graph)
	}

	fmt.Printf("[%s] [8/13] registering fields\n", time.Now().Format(time.TimeOnly))
	registry.RegisterFields(app, graphsLst)

	if !EVAL {
		app.WriteAppToJSON()
		for fn, ssagraph := range funcGraphs {
			ssagraph.WriteToDOTFile(appname, fn, false)
		}
	}

	fmt.Printf("[%s] [9/13] running SSA tainter\n", time.Now().Format(time.TimeOnly))
	for _, ssagraph := range funcGraphs {
		tainter.RunTainter(ssagraph)
	}

	if !EVAL {
		for fn, ssagraph := range funcGraphs {
			ssagraph.WriteToDOTFile(appname, fn, true)
		}
	}

	fmt.Printf("[%s] [10/13] creating new abstract call graph\n", time.Now().Format(time.TimeOnly))
	absgraph := abstractgraph.NewAbstractCallGraph(app)
	for _, entrypoint := range app.GetEntrypointsShortPaths() {
		abstractgraph.Parse(absgraph, entrypoint, true, funcGraphs)
	}

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
	// phase 1: two passes
	fmt.Printf("[%s] [11/13] starting schema builder\n", time.Now().Format(time.TimeOnly))
	iterator.Run(detection.PHASE_1_SCHEMA_BUILDER)
	fmt.Printf("[%s] [12/13] starting schema builder (read only)\n", time.Now().Format(time.TimeOnly))
	iterator.Run(detection.PHASE_1_SCHEMA_BUILDER_READ_ONLY)

	elapsed_schema := time.Since(start_schema)

	fmt.Printf("[%s] [13/13] starting pattern detection\n", time.Now().Format(time.TimeOnly))
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

	results := detection.SaveResults(app, detector1, detector2, detector3, detector4, detector5)
	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Printf("Execution time (TOTAL):\t%.4f s\n", elapsed_total.Seconds())
	fmt.Printf("Execution time (PARSING):\t%.4f s\n", elapsed_parsing.Seconds())
	fmt.Printf("Execution time (SCHEMA):\t%.4f s\n", elapsed_schema.Seconds())
	fmt.Printf("Execution time (DETECTION):\t%.4f s\n", elapsed_detection.Seconds())
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
