package main

import (
	"fmt"
	"go/token"
	"log"
	"os"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"

	"analyzer/pkg/parser"
	"analyzer/pkg/ssa_graph"
	"analyzer/pkg/abstractcallgraph"
	"analyzer/pkg/tainter"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: program <appname>\n")
		fmt.Fprintln(os.Stderr, "available appnames:")
		fmt.Fprintln(os.Stderr, "- examples/postnotification")
		fmt.Fprintln(os.Stderr, "- examples/shoppingcart")
		fmt.Fprintln(os.Stderr, "- blueprint/postnotification/storageservice_storepost")
		os.Exit(1)
	}

	appname := os.Args[1]

	// ensure output sub directory exists
	err := os.MkdirAll(fmt.Sprintf("output/%s", appname), os.ModePerm)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	// ensure output sub directory for graphs exists
	err = os.MkdirAll(fmt.Sprintf("output/%s/graphs", appname), os.ModePerm)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	prog, pkgs, err := buildProgram(appname)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	fmt.Println("[INFO] running analysis for packages:")
	for _, pkg := range pkgs {
		fmt.Printf("\t- %s\n", pkg.String())
	}

	result, err := parser.InitPointerAnalysis(prog, pkgs)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	funcGraphs := make(map[string]*ssa_graph.SSAGraph)

	for _, pkg := range pkgs {
		parser.RunSSAAnalysis(appname, prog, pkg, funcGraphs)
	}

	for _, ssa_graph := range funcGraphs {
		ssa_graph.SortNodes()
	}

	for _, pkg := range pkgs {
		parser.RunPointerToAnalysis(appname, prog, pkg, result, funcGraphs)
	}

	for _, ssa_graph := range funcGraphs {
		tainter.RunTaint(ssa_graph)
	}

	fmt.Print("\n\n ========== NODES ========== \n\n")
	for fn, ssa_graph := range funcGraphs {
		for _, node := range ssa_graph.GetNodes() {
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
	for fn, ssa_graph := range funcGraphs {
		for _, node := range ssa_graph.GetNodes() {
			if node.IsTainted() {
				for obj, dbfields := range node.GetTaints() {
					fmt.Printf("[%s] %s [%s]: %s\n", fn, node.String(), node.GetName(), obj)
					for _, dbfield := range dbfields {
						fmt.Printf("\t\t |--> %s\n", dbfield)
					}
				}
			}
		}
	}

	for fn, ssaGraph := range funcGraphs {
		ssaGraph.WriteToDOTFile(appname, fn)
	}

	fmt.Println("\n[INFO] successfully analyzed app (" + appname + ")\n")
	
	var entryPoints = []string{
		"postnotification_simple.UploadService.UploadPost",
		"postnotification_simple.StorageService.StorePost",
		"postnotification_simple.StorageService.ReadPost",
		"postnotification_simple.NotifyService.workerThread",
	}

	abstractGraph := abstractcallgraph.NewAbstractGraph()
	abstractGraph.Init(entryPoints, funcGraphs)
}

func buildProgram(appname string) (*ssa.Program, []*ssa.Package, error) {
	// e.graph. "../apps/test2/main.go"
	filepath := "apps/" + appname + "/main.go"
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
			if pkg.Pkg.Path() == "github.com/blueprint-uservices/blueprint/examples/postnotification_simple/workflow/postnotification_simple" {
				pkgs = append(pkgs, pkg)
			}
		}
	}
	return prog, pkgs, nil
}
