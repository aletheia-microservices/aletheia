package main

import (
	"fmt"
	"go/token"
	"log"
	"os"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/detection"
	"analyzer/pkg/detection/constraints/foreignkeycascade"
	"analyzer/pkg/detection/constraints/foreignkeycoordination"
	"analyzer/pkg/detection/constraints/unicityconcurrency"
	"analyzer/pkg/ssagraph"
	"analyzer/pkg/ssagraph/parser"
	"analyzer/pkg/ssagraph/registry"
	"analyzer/pkg/ssagraph/tainter"
	"analyzer/pkg/utils"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: program <appname>\n")
		fmt.Fprintln(os.Stderr, "available appnames:")
		fmt.Fprintln(os.Stderr, "- postnotification_simple blueprint/postnotification/notifyservice_run")
		fmt.Fprintln(os.Stderr, "- dsb_media_sql blueprint/dsb_media_sql/api_readmovie")
		fmt.Fprintln(os.Stderr, "- digota blueprint/digota/skuservice_get")
		fmt.Fprintln(os.Stderr, "- sockshop3 blueprint/sockshop3/userservice_login")
		fmt.Fprintln(os.Stderr, "- dsb_sn blueprint/dsb_sn/poststorageservice_storepost")
		os.Exit(1)
	}

	appname := os.Args[1]
	apppath := os.Args[2]
	app := app.NewApp(appname)
	app.Init()

	// ensure output sub directory exists
	err := os.MkdirAll(fmt.Sprintf("output/%s", appname), os.ModePerm)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	// ensure output sub directory for graphs exists
	err = os.MkdirAll(fmt.Sprintf("output/%s/ssagraphs", appname), os.ModePerm)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	// ensure output sub directory for graphs exists
	err = os.MkdirAll(fmt.Sprintf("output/%s/ssa", appname), os.ModePerm)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	prog, pkgs, err := buildProgram(apppath)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	app.InitServiceFields(pkgs)
	app.ParseSchemaFromUserFile()

	fmt.Println("[INFO] running analysis for packages:")
	for _, pkg := range pkgs {
		fmt.Printf("\t- %s\n", pkg.String())
	}

	result, err := parser.InitPointerAnalysis(prog, pkgs)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	funcGraphs := make(map[string]*ssagraph.SSAGraph)

	for _, pkg := range pkgs {
		parser.RunSSAAnalysis(app, prog, pkg, funcGraphs)
	}

	for _, ssagraph := range funcGraphs {
		ssagraph.SortNodes()
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

	for _, ssagraph := range funcGraphs {
		tainter.RunTainter(ssagraph)
	}

	fmt.Print("\n\n ========== NODES ========== \n\n")
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
	}

	for fn, ssagraph := range funcGraphs {
		fmt.Printf("[MAIN] go ssa graph for (%s): %v\n", fn, ssagraph)
	}

	fmt.Print("\n\n ========== TAINTS ========== \n\n")
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
	}

	for fn, ssagraph := range funcGraphs {
		ssagraph.WriteToDOTFile(appname, fn)
	}

	fmt.Println("\n[INFO] successfully analyzed app (" + appname + ")\n")

	absgraph := abstractgraph.NewAbstractCallGraph(app)
	for _, entrypoint := range app.GetEntrypointsShortPaths() {
		abstractgraph.Parse(absgraph, entrypoint, true, funcGraphs)
	}

	detector1 := foreignkeycoordination.NewDetector("foreign-key")
	detector2 := foreignkeycascade.NewDetector()
	detector3 := unicityconcurrency.NewDetector()
	iterator := detection.NewIterator(app, absgraph, detector1, detector2, detector3)
	iterator.Run()

	absgraph.WriteToDOTFile(appname, true)
	absgraph.WriteToDOTFile(appname, false)
	app.WriteAppToJSON()
	app.WriteSchemaToJSON()

	fmt.Print("\n\n ========== SERVICE CALLS ========== \n\n")
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
	}

	fmt.Print("\n\n ========== DATABASE CALLS ========== \n\n")
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
	}

	fmt.Print("\n\n ========== APP ========== \n\n")
	fmt.Println(app.String())

	detector1.ComputeResults()
	res1 := detection.SaveResults(app, detector1)
	fmt.Println(res1)

	detector2.ComputeResults()
	res2 := detection.SaveResults(app, detector2)
	fmt.Println(res2)

	detector3.ComputeResults()
	res3 := detection.SaveResults(app, detector3)
	fmt.Println(res3)
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
