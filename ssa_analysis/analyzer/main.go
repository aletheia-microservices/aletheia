package main

import (
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

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: program <appname>\n")
		fmt.Fprintln(os.Stderr, "available appnames:")
		fmt.Fprintln(os.Stderr, "- foobar")
		fmt.Fprintln(os.Stderr, "- postnotification_simple")
		fmt.Fprintln(os.Stderr, "- dsb_sn2")
		fmt.Fprintln(os.Stderr, "- digota")
		fmt.Fprintln(os.Stderr, "- sockshop3")
		fmt.Fprintln(os.Stderr, "- dsb_media_sql")
		fmt.Fprintln(os.Stderr, "- dsb_hotel2")
		fmt.Fprintln(os.Stderr, "- train_ticket2")
		os.Exit(1)
	}

	appname := os.Args[1]

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

	prog, pkgs, err := buildProgram(apppath)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	app.InitServiceFields(pkgs)
	app.ParseSchemaFromUserFile()

	fmt.Println("[INFO] running analysis for packages:")
	/* for _, pkg := range pkgs {
		fmt.Printf("\t- %s\n", pkg.String())
	} */

	result, err := parser.InitPointerAnalysis(prog, pkgs)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	funcGraphs := make(map[string]*ssagraph.SSAGraph)

	for _, pkg := range pkgs {
		parser.RunSSAAnalysis(app, prog, pkg, funcGraphs)
	}

	for _, ssagraph := range funcGraphs {
		ssagraph.Sort()
	}

	for _, pkg := range pkgs {
		parser.RunPointerToAnalysis(appname, prog, pkg, result, funcGraphs)
	}

	var graphsLst []*ssagraph.SSAGraph
	for _, graph := range funcGraphs {
		graphsLst = append(graphsLst, graph)
	}

	registry.RegisterFields(app, graphsLst)
	app.WriteAppToJSON()

	for fn, ssagraph := range funcGraphs {
		ssagraph.WriteToDOTFile(appname, fn, false)
	}

	for _, ssagraph := range funcGraphs {
		tainter.RunTainter(ssagraph)
	}

	/* fmt.Print("\n\n ========== NODES ========== \n\n")
	for fn, ssagraph := range funcGraphs {
		for _, node := range ssagraph.GetNodes() {
			var prefix string
			if node.GetName() != "" {
				prefix = node.GetName() + ":"
			} else {
				prefix = "\t"
			}
			if node.GetInstruction() != nil {
				fmt.Printf("[%s] [%s] [%T] \t %s %v\n", fn, node.GetID(), node.GetInstruction(), prefix, node.GetInstruction().String())
			} else {
				fmt.Printf("[%s] [%s] [%T] \t %s %v\n", fn, node.GetID(), node.GetValue(), prefix, node.GetValue().String())
			}
		}
	} */

	/* for fn, ssagraph := range funcGraphs {
		fmt.Printf("[MAIN] go ssa graph for (%s): %v\n", fn, ssagraph)
	} */

	/* fmt.Print("\n\n ========== TAINTS ========== \n\n")
	for fn, ssagraph := range funcGraphs {
		for _, node := range ssagraph.GetNodes() {
			if node.IsTainted() {
				for obj, taints := range node.GetTaints() {
					fmt.Printf("[%s] %s [%s]: %s\n", fn, node.String(), node.GetName(), obj)
					for _, taint := range taints {
						fmt.Printf("\t\t |--> %s\n", taint.String())
					}
				}
			}
		}
	} */

	for fn, ssagraph := range funcGraphs {
		ssagraph.WriteToDOTFile(appname, fn, true)
	}

	fmt.Println("\n[INFO] successfully analyzed app (" + appname + ")\n")

	absgraph := abstractgraph.NewAbstractCallGraph(app)
	for _, entrypoint := range app.GetEntrypointsShortPaths() {
		abstractgraph.Parse(absgraph, entrypoint, true, funcGraphs)
	}

	absgraph.WriteVisited(appname)

	elapsed_parsing := time.Since(start)

	detector1 := keycoordination.NewDetector(keycoordination.DETECTION_TYPE_PRIMARY_KEY)
	detector2 := keycoordination.NewDetector(keycoordination.DETECTION_TYPE_FOREIGN_KEY)
	detector3 := foreignkeycascade.NewDetector()
	detector4 := foreignkeyconcurrency.NewDetector()
	detector5 := unicityconcurrency.NewDetector()
	iterator := detection.NewIterator(app, absgraph, detector1, detector2, detector3, detector4, detector5)

	start_schema := time.Now()
	// phase 1: two passes
	iterator.Run(detection.PHASE_1_SCHEMA_BUILDER)
	iterator.Run(detection.PHASE_1_SCHEMA_BUILDER_READ_ONLY)

	elapsed_schema := time.Since(start_schema)
	start_detection := time.Now()

	// phase 2: one pass for all detectors
	iterator.Run(detection.PHASE_2_PATTERN_DETECTOR)

	// phase 0: dummy pass to generate dot files with taints for debugging
	iterator.Run(detection.PHASE_0_DEBUG)

	absgraph.WriteToDOTFile(appname, true)
	absgraph.WriteToDOTFile(appname, false)

	app.WriteAppToJSON()
	app.WriteSchemaToJSON()

	/* fmt.Print("\n\n ========== SERVICE CALLS ========== \n\n")
	for _, node := range absgraph.GetNodes() {
		if node.GetNodeType() == abstractgraph.NODE_SERVICE {
			for _, edge := range absgraph.GetEdgesToNode(node) {
				fmt.Printf("SERVICE CALL: %s\n", edge.String())
				for i, arg := range edge.GetArguments() {
					fmt.Printf("ARG %d (%s) w/ TAINTS:\n%s", i, arg.String(), arg.TaintLongString())
				}
				fmt.Println("--")
				for i, param := range edge.GetToNode().GetParams() {
					fmt.Printf("PARAM %d (%s) w/ TAINTS:\n%s", i, param.String(), param.TaintLongString())
				}
				fmt.Println("--")
				for i, ret := range edge.GetReturns() {
					fmt.Printf("RET (EDGE) %d (%s) w/ TAINTS:\n%s", i, ret.String(), ret.TaintLongString())
				}
				fmt.Println("--")
				for i, ret := range edge.GetToNode().GetReturns() {
					fmt.Printf("RET (NODE) %d (%s) w/ TAINTS:\n%s", i, ret.String(), ret.TaintLongString())
				}
				fmt.Println()
			}
		}
	} */

	/* fmt.Print("\n\n ========== DATABASE CALLS ========== \n\n")
	for _, node := range absgraph.GetNodes() {
		if node.GetNodeType() == abstractgraph.NODE_DATABASE {
			for _, edge := range absgraph.GetEdgesToNode(node) {
				fmt.Printf("DATABASE CALL: %s\n", edge.String())
				for i, arg := range edge.GetArguments() {
					fmt.Printf("ARG %d (%s) w/ TAINTS:\n%s", i, arg.String(), arg.TaintLongString())
				}
				fmt.Println()
			}
		}
	} */

	elapsed_total := time.Since(start)
	elapsed_detection := time.Since(start_detection)

	fmt.Print("\n\n ========== APP ========== \n\n")
	fmt.Println(app.String())

	results := detection.SaveResults(app, detector1, detector2, detector3, detector4, detector5)
	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Printf("Execution time (TOTAL): %.4f s\n", elapsed_total.Seconds())
	fmt.Printf("Execution time (PARSING): %.4f s\n", elapsed_parsing.Seconds())
	fmt.Printf("Execution time (SCHEMA): %.4f s\n", elapsed_schema.Seconds())
	fmt.Printf("Execution time (DETECTION): %.4f s\n", elapsed_detection.Seconds())
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
		fmt.Println("parse error:", err)
		return nil, nil, err
	}
	conf.CreateFromFiles("main", file)

	iprog, err := conf.Load()
	if err != nil {
		fmt.Println("type error:", err)
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
			} else {
				fmt.Printf("skipping... %s\n", pkg.Pkg.Path())
			}
		}
	}
	return prog, pkgs, nil
}
