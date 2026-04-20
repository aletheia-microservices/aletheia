package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"gopkg.in/yaml.v2"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/config"
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

var (
	INIT             bool
	EVAL             bool
	SYNTHETIC        bool
	INPUT_REFS       bool
	DEBUG            bool
	DETECTION_CONFIG string
)

const EVAL_METRICS_BASE = "eval/metrics"

func main() {
	flag.BoolVar(&INIT, "init", false, "enable init mode")
	flag.BoolVar(&EVAL, "eval", false, "enable evaluation mode")
	flag.BoolVar(&SYNTHETIC, "synthetic", false, "enable synthetic app")
	flag.BoolVar(&INPUT_REFS, "refs", false, "enable input of references")
	flag.BoolVar(&DEBUG, "debug", false, "enable debug output")
	flag.StringVar(&DETECTION_CONFIG, "detection_config", "", "path to detection config yaml")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "usage: program [--eval] [--debug] [--detection_config filepath.yaml] <appname>\n")
		fmt.Fprintln(os.Stderr, "available appnames:")
		fmt.Fprintln(os.Stderr, "- postnotification")
		fmt.Fprintln(os.Stderr, "- dsb_socialnetwork")
		fmt.Fprintln(os.Stderr, "- digota")
		fmt.Fprintln(os.Stderr, "- eshopmicroservices")
		fmt.Fprintln(os.Stderr, "- sockshop")
		fmt.Fprintln(os.Stderr, "- dsb_mediamicroservices")
		fmt.Fprintln(os.Stderr, "- dsb_hotel2")
		fmt.Fprintln(os.Stderr, "- trainticket")
		fmt.Fprintln(os.Stderr, "- synthetic")
		os.Exit(1)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.TimeOnly,
	})
	logrus.SetLevel(logrus.InfoLevel)

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

	if DETECTION_CONFIG != "" {
		detection.LoadInputConfig(appname, DETECTION_CONFIG)
	}

	start := time.Now()

	logrus_ctx := logrus.WithField("app", appname)

	// ------------ PART 1
	logrus_ctx.WithField("synthetic", SYNTHETIC).Infof("[1/12] initializing program")

	apppath := utils.GetAppRootPackagePath(appname)
	app := app.NewApp(appname)
	app.Init(SYNTHETIC)

	if INIT {
		// load app from init and skip analysis
		return
	}

	elapsed_blueprint_compiler := time.Since(start)

	// ensure output sub directory exists
	err := os.MkdirAll(fmt.Sprintf("output/%s", appname), os.ModePerm)
	if err != nil {
		logrus.Fatalf("error: %s", err.Error())
	}

	// ensure output sub directory for graphs exists
	if DEBUG {
		err = os.MkdirAll(fmt.Sprintf("output/%s/ssagraphs/tainted", appname), os.ModePerm)
		if err != nil {
			logrus.Fatalf("error: %s", err.Error())
		}
		err = os.MkdirAll(fmt.Sprintf("output/%s/ssagraphs/untainted", appname), os.ModePerm)
		if err != nil {
			logrus.Fatalf("error: %s", err.Error())
		}
		err = os.MkdirAll(fmt.Sprintf("output/%s/abstractcallgraph", appname), os.ModePerm)
		if err != nil {
			logrus.Fatalf("error: %s", err.Error())
		}
		// ensure output sub directory for graphs exists
		err = os.MkdirAll(fmt.Sprintf("output/%s/ssa", appname), os.ModePerm)
		if err != nil {
			logrus.Fatalf("error: %s", err.Error())
		}
	}

	// ------------ PART 2
	logrus_ctx.Infof("[2/12] building program")

	prog, pkgs, err := buildProgram(apppath)
	if err != nil {
		logrus.Fatalf("error: %s", err.Error())
	}

	app.InitServiceFields(pkgs)
	app.ParseSQLSchemaFromUserFile()
	app.ParseNoSQLSchemaFromUserFile()
	if INPUT_REFS {
		logrus_ctx.Infof("reading input refs...")
		app.ParseUserInputReferences()
	}

	funcGraphs := make(map[string]*ssagraph.SSAGraph)

	// ------------ PART 3
	logrus_ctx.Infof("[3/12] running SSA analysis")
	for _, pkg := range pkgs {
		parser.RunSSAAnalysis(app, prog, pkg, funcGraphs)
	}

	for _, ssagraph := range funcGraphs {
		ssagraph.Sort()
	}

	var graphsLst []*ssagraph.SSAGraph
	for _, graph := range funcGraphs {
		graphsLst = append(graphsLst, graph)
	}

	// ------------ PART 4
	logrus_ctx.Infof("[4/12] registering fields")
	registry.RegisterFields(app, graphsLst)

	elapsed_ssa_parsing := time.Since(start)
	start_ssa_tainting := time.Now()

	if !EVAL {
		app.WriteAppToJSON()
		if DEBUG {
			for fn, ssagraph := range funcGraphs {
				ssagraph.WriteToDOTFile(appname, fn, false)
			}
		}
	}

	// ------------ PART 5
	logrus_ctx.Infof("[5/12] running SSA tainter for single graphs")
	for _, ssagraph := range funcGraphs {
		tainter.RunTainter(ssagraph)
	}

	elapsed_ssa_tainting := time.Since(start_ssa_tainting)

	// ------------ PART 6
	logrus_ctx.Infof("[6/12] combining SSA graphs")
	for _, ssagraph := range funcGraphs {
		tainter.Combine(ssagraph, funcGraphs)
	}

	if !EVAL && DEBUG {
		var written = make(map[string]bool)
		for fn, ssagraph := range funcGraphs {
			ssagraph.WriteToDOTFile(appname, fn, true)
			written[fn] = true
		}
		for fn, ssagraph := range funcGraphs {
			for _, toGraph := range ssagraph.GetAllCombinedGraphs() {
				newFn := fn + "." + toGraph.GetMethodName()
				if exists, _ := written[newFn]; !exists {
					toGraph.WriteToDOTFile(appname, newFn, true)
				}
				written[newFn] = true
			}
		}
	}

	// ------------ PART 7
	logrus_ctx.Infof("[7/12] creating new abstract call graph")
	absgraph := abstractgraph.NewAbstractCallGraph(app)
	for _, entrypoint := range app.GetEntrypointsShortPaths() {
		abstractgraph.Parse(absgraph, entrypoint, true, funcGraphs)
	}

	// ------------ PART 8
	logrus_ctx.Infof("[8/12] releasing memory associated with ssa graph")
	graphsLst = nil
	pkgs = nil
	for fn, ssagraph := range funcGraphs {
		if ssagraph != nil {
			ssagraph.Release()
		}
		delete(funcGraphs, fn)
	}
	funcGraphs = nil

	elapsed_parsing := time.Since(start)

	detector1 := keycoordination.NewDetector(keycoordination.DETECTION_TYPE_PRIMARY_KEY)
	detector2 := keycoordination.NewDetector(keycoordination.DETECTION_TYPE_FOREIGN_KEY)
	detector3 := foreignkeycascade.NewDetector()
	detector4 := foreignkeyconcurrency.NewDetector()
	detector5 := unicityconcurrency.NewDetector()
	iterator := detection.NewIterator(app, absgraph, detector1, detector2, detector3, detector4, detector5)

	start_schema := time.Now()
	// ------------ PART 9
	logrus_ctx.Infof("[9/12] starting schema builder")
	iterator.Run(detection.PHASE_1_SCHEMA_BUILDER)
	// ------------ PART 10
	if config.Global.DualPassSchemaBuilding {
		logrus_ctx.Infof("[10/12] starting schema builder (read only)")
		iterator.Run(detection.PHASE_1_SCHEMA_BUILDER_READ_ONLY)
	} else {
		logrus_ctx.Infof("[10/12] skipping schema builder (read only)...")
	}

	elapsed_schema := time.Since(start_schema)

	// ------------ PART 11
	logrus_ctx.Infof("[11/12] starting pattern detection")
	start_detection := time.Now()

	// phase 2: one pass for all detectors
	iterator.Run(detection.PHASE_2_PATTERN_DETECTOR)

	elapsed_total := time.Since(start)
	elapsed_detection := time.Since(start_detection)

	if !EVAL {
		if DEBUG {
			// phase 0: dummy pass to generate dot files with taints for debugging
			iterator.Run(detection.PHASE_0_DEBUG)
			absgraph.WriteToDOTFile(appname, true)
			absgraph.WriteToDOTFile(appname, false)
		}

		app.WriteAppToJSON()
		app.WriteSchemaToJSON()
	}

	// ------------ PART 12
	logrus_ctx.Infof("[12/12] saving results")
	results := detection.SaveResults(app, detector1, detector2, detector3, detector4, detector5)
	for _, result := range results {
		fmt.Println(result)
	}

	if EVAL || DEBUG {
		fmt.Printf("Execution time (TOTAL):\t\t%.4f s\n", elapsed_total.Seconds())
		fmt.Printf("Execution time (BLUEPRINT):\t%.4f s\n", elapsed_blueprint_compiler.Seconds())
		fmt.Printf("Execution time (PARSING):\t%.4f s\n", elapsed_parsing.Seconds())
		fmt.Printf("Execution time (SSA PARS):\t%.4f s\n", elapsed_ssa_parsing.Seconds())
		fmt.Printf("Execution time (SSA TAIN):\t%.4f s\n", elapsed_ssa_tainting.Seconds())
		fmt.Printf("Execution time (SCHEMA):\t%.4f s\n", elapsed_schema.Seconds())
		fmt.Printf("Execution time (DETECTION):\t%.4f s\n", elapsed_detection.Seconds())
	}

	app.WriteAppToJSON()
	app.WriteSchemaToJSON()

	if EVAL {
		times := AnalysisTimes{
			App:              app.GetName(),
			NumMicroservices: app.NumberOfMicroservices(),
			NumDatastores:    app.NumberOfDatastores(),
			NumCallGraphs:    absgraph.ComputeAndGetNumCallGraphs(),
			Blueprint:        elapsed_blueprint_compiler.Seconds(),
			Rpcs:             absgraph.GetRPCCount(),
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
	NumCallGraphs    int     `yaml:"callgraphs"`
	Rpcs             int     `yaml:"rpcs"`
	Blueprint        float64 `yaml:"blueprint"`
	Total            float64 `yaml:"total_s"`
	Parsing          float64 `yaml:"parsing_s"`
	Schema           float64 `yaml:"schema_s"`
	Detection        float64 `yaml:"detection_s"`
}

func saveAnalysisTime(app *app.App, times AnalysisTimes) {
	ts := time.Now().Unix()
	dir := path.Join(EVAL_METRICS_BASE, time.Now().Format(time.DateOnly))
	if SYNTHETIC {
		dir += "/synthetic"
	} else {
		dir += "/realistic"
	}
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}

	filepath := fmt.Sprintf("%s/%s_%d.yaml", dir, app.GetName(), ts)

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
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedDeps,
	}

	initialPkgs, err := packages.Load(cfg, apppath)
	if err != nil {
		return nil, nil, err
	}
	if len(initialPkgs) == 0 {
		return nil, nil, fmt.Errorf("no packages found for apppath (%s)", apppath)
	}

	typeToPkg := make(map[*types.Package]*packages.Package)
	var visitPkgs func(p *packages.Package)
	visitPkgs = func(p *packages.Package) {
		if p == nil || p.Types == nil || typeToPkg[p.Types] != nil {
			return
		}
		typeToPkg[p.Types] = p
		for _, imp := range p.Imports {
			visitPkgs(imp)
		}
	}
	for _, p := range initialPkgs {
		visitPkgs(p)
	}

	prog := ssa.NewProgram(initialPkgs[0].Fset, 0)

	// recursively create SSA packages following go/types imports
	seenTypes := make(map[*types.Package]bool)
	var createSSA func(tp *types.Package)
	createSSA = func(tp *types.Package) {
		if tp == nil || seenTypes[tp] {
			return
		}
		seenTypes[tp] = true

		pp := typeToPkg[tp]
		var files []*ast.File
		var info *types.Info
		if pp != nil && len(pp.Syntax) > 0 {
			files = pp.Syntax
			info = pp.TypesInfo
		}
		if prog.Package(tp) == nil {
			_ = prog.CreatePackage(tp, files, info, true)
		}

		for _, imp := range tp.Imports() {
			createSSA(imp)
		}
	}

	for _, p := range initialPkgs {
		if p.Types != nil {
			createSSA(p.Types)
		}
	}

	prog.Build()

	var pkgsSeen = make(map[*ssa.Package]bool)
	var pkgs []*ssa.Package
	for _, p := range initialPkgs {
		if p.Types != nil {
			if progPkg := prog.Package(p.Types); progPkg != nil && !pkgsSeen[progPkg] {
				pkgsSeen[progPkg] = true
				pkgs = append(pkgs, progPkg)
			}
		}
	}
	for _, progPkg := range prog.AllPackages() {
		if utils.IsAppPackagePath(progPkg.Pkg.Path()) && progPkg != nil && !pkgsSeen[progPkg] {
			pkgsSeen[progPkg] = true
			pkgs = append(pkgs, progPkg)
		}
	}
	return prog, pkgs, nil
}
