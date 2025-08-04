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
	"analyzer/pkg/detection/constraints/foreignkeycoordination"
	"analyzer/pkg/ssagraph"
	"analyzer/pkg/ssagraph/parser"
	"analyzer/pkg/ssagraph/registry"
	"analyzer/pkg/ssagraph/tainter"
)

const (
	APP_PATH_POSTNOTIFICATION = "github.com/blueprint-uservices/blueprint/examples/postnotification_simple/workflow/postnotification_simple"
	APP_PATH_DIGOTA           = "github.com/blueprint-uservices/blueprint/examples/digota/workflow/digota"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: program <appname>\n")
		fmt.Fprintln(os.Stderr, "available appnames:")
		fmt.Fprintln(os.Stderr, "- postnotification examples/postnotification")
		fmt.Fprintln(os.Stderr, "- shoppingcart examples/shoppingcart")
		fmt.Fprintln(os.Stderr, "- postnotification_simple blueprint/postnotification/notifyservice_run")
		fmt.Fprintln(os.Stderr, "- digota blueprint/digota/skuservice_get")
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

	app.InitFields(pkgs)

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
		abstractgraph.Parse(absgraph, entrypoint, funcGraphs)
	}

	detector := foreignkeycoordination.NewDetector("foreign-key")
	iterator := detection.NewIterator(app, absgraph, detector)
	iterator.Run()

	absgraph.WriteToDOTFile(appname, true)
	absgraph.WriteToDOTFile(appname, false)
	app.WriteAppToJSON()
	app.WriteSchemaToJSON()

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

	detector.ComputeResults()
	res := detection.SaveResults(app, detector)
	fmt.Println(res)
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
			if pkg.Pkg.Path() == APP_PATH_POSTNOTIFICATION {
				pkgs = append(pkgs, pkg)
			} else if pkg.Pkg.Path() == APP_PATH_DIGOTA {
				pkgs = append(pkgs, pkg)
			} else {
				fmt.Printf("skipping... %s\n", pkg.Pkg.Path())
			}
		}
	}
	return prog, pkgs, nil
}
