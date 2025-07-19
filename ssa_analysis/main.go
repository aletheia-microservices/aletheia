package main

import (
	"fmt"
	"go/token"
	"go/types"
	"log"
	"os"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/types/typeutil"
)

// -------------------------------------
// ------------- CONSTANTS -------------
// -------------------------------------

const packagePath = "./examples/test2"

//const packagePath = "../blueprint/examples/postnotification_simple/workflow/postnotification_simple/"

// -------------------------------------

var createdPkgs map[*packages.Package]bool

func recurse(prog *ssa.Program, pkg *packages.Package) {
	if _, ok := createdPkgs[pkg]; ok {
		return
	}
	prog.CreatePackage(pkg.Types, pkg.Syntax, pkg.TypesInfo, false)
	createdPkgs[pkg] = true
	for _, impt := range pkg.Imports {
		recurse(prog, impt)
	}
}

func main() {
	createdPkgs = make(map[*packages.Package]bool)
	cfg := &packages.Config{Mode: packages.LoadAllSyntax}
	pkgs, err := packages.Load(cfg, packagePath)
	if err != nil {
		log.Fatal(err)
	}

	fset := token.NewFileSet()
	prog := ssa.NewProgram(fset, 0)

	ssaPkgs := make([]*ssa.Package, len(pkgs))
	for i, pkg := range pkgs {
		if _, ok := createdPkgs[pkg]; !ok {
			prog.CreatePackage(pkg.Types, pkg.Syntax, pkg.TypesInfo, false)
			createdPkgs[pkg] = true
			for _, impt := range pkg.Imports {
				recurse(prog, impt)
			}
		}
		ssaPkgs[i] = prog.Package(pkg.Types)
	}

	prog.Build()

	var appPkgs []*ssa.Package
	for _, ssaPkg := range ssaPkgs {
		if ssaPkg == nil || ssaPkg.Pkg == nil {
			continue
		}
		/* if ssaPkg.Pkg.Name() != "postnotification_simple" {
			continue
		} */
		/* if ssaPkg.Func("main") == nil && ssaPkg.Func("init") == nil {
			continue
		} */
		appPkgs = append(appPkgs, ssaPkg)
	}

	ssaAnalysis(prog, appPkgs)
}

func iterateFunc(outFile *os.File, fn *ssa.Function, memberType types.Type) {
	namedMemberType, ok := memberType.(*types.Named)

	fmt.Printf("=============================\n")
	if ok {
		fmt.Printf("%s.%s()\n", namedMemberType.Obj().Name(), fn.Name())
	} else {
		fmt.Printf("%s()\n", fn.Name())
	}
	fmt.Printf("=============================\n")

	fmt.Fprintf(outFile, "Function: %s\n", fn.Name())
	fmt.Printf("\n--------------- Function: %s\n", fn.Name())
	for i, block := range fn.Blocks {
		fmt.Fprintf(outFile, "Block #%d: %s.%s\n", i, fn.Name(), block.Comment)
		fmt.Printf("----- Block #%d: %s.%s\n", i, fn.Name(), block.Comment)

		for j, instr := range block.Instrs {
			if val, ok := instr.(ssa.Value); ok {
				fmt.Fprintf(outFile, "\t\t\t%02d: %s = %s\n", j, val.Name(), instr.String())
			} else {
				fmt.Fprintf(outFile, "\t\t\t%02d: %s\n", j, instr.String())
			}
		}
	}

	fmt.Println()
	fmt.Println()
}

func ssaAnalysis(prog *ssa.Program, pkgs []*ssa.Package) {
	//outFile, err := os.Create(fmt.Sprintf("%s/app.ssa", packagePath))
	outFile, err := os.Create("output/app.ssa")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	for _, ssaPkg := range pkgs {
		outfile, err := os.Create(fmt.Sprintf("./output/%s.out", ssaPkg.Pkg.Name()))
		if err != nil {
			log.Fatalf("failed to create output file: %v", err)
		}
		defer outfile.Close()
		ssaPkg.WriteTo(outfile)

		for _, member := range ssaPkg.Members {
			switch m := member.(type) {
			case *ssa.Function:
				iterateFunc(outFile, m, nil)

			case *ssa.Global:
				fmt.Fprintf(outFile, "\tGlobal: %s, Type: %s\n", m.Name(), m.Type().String())

			case *ssa.Type:
				fmt.Fprintf(outFile, "\tType: %s\n", m.Type())

				// this logic was copied from
				// package: golang.org/x/tools/go/ssa
				// file: print.go
				// function: func (p *Package) WriteTo(w io.Writer) (int64, error)
				for _, sel := range typeutil.IntuitiveMethodSet(m.Type(), &prog.MethodSets) {
					method := prog.MethodValue(sel)
					fmt.Fprintf(outFile, "\tMethod: %v\n", sel.Obj().Type())
					if method != nil {
						iterateFunc(outFile, method, m.Type())
					}
				}

				methods := prog.MethodSets.MethodSet(m.Type().Underlying())
				for i := 0; i < methods.Len(); i++ {
					sel := methods.At(i)
					fmt.Fprintf(outFile, "\tMethod: %v\n", sel.Obj().Type())
					method := prog.MethodValue(sel)
					if method != nil {
						iterateFunc(outFile, method, m.Type())
					}
				}

			default:
				fmt.Fprintf(outFile, "\tUnknown member type: %T\n", m)
			}
		}
	}
}
