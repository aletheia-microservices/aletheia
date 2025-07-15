package main

import (
	"fmt"
	"go/token"
	"log"
	"os"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/types/typeutil"
)

// -------------------------------------
// ------------- CONSTANTS -------------
// -------------------------------------
const (
	inputpackagepath = "../workflow/postnotification/"
	outfilename      = "./ssa-simple.out"
)

// -------------------------------------

var ssaPkgs map[*packages.Package]bool

func recurse(prog *ssa.Program, pkg *packages.Package) {
	if _, ok := ssaPkgs[pkg]; ok {
		return
	}
	prog.CreatePackage(pkg.Types, pkg.Syntax, pkg.TypesInfo, false)
	ssaPkgs[pkg] = true
	for _, impt := range pkg.Imports {
		recurse(prog, impt)
	}
}

func main() {
	cfg := &packages.Config{Mode: packages.LoadAllSyntax}
	pkgs, err := packages.Load(cfg, inputpackagepath)
	if err != nil {
		log.Fatal(err)
	}

	fset := token.NewFileSet()
	//prog := ssa.NewProgram(fset, ssa.PrintFunctions)
	prog := ssa.NewProgram(fset, 0)

	ssaPkgs = make(map[*packages.Package]bool)
	ssaPkgsFiltered := make([]*ssa.Package, len(pkgs))
	for i, pkg := range pkgs {
		if _, ok := ssaPkgs[pkg]; ok {
			continue
		}
		ssaPkgsFiltered[i] = prog.CreatePackage(pkg.Types, pkg.Syntax, pkg.TypesInfo, false)
		ssaPkgs[pkg] = true
		for _, impt := range pkg.Imports {
			recurse(prog, impt)
		}
	}

	prog.Build()

	var appPkgs []*ssa.Package
	for _, ssaPkg := range ssaPkgsFiltered {
		if ssaPkg == nil || ssaPkg.Pkg == nil {
			continue
		}
		if ssaPkg.Pkg.Name() != "postnotification" {
			continue
		}
		/* if ssaPkg.Func("main") == nil && ssaPkg.Func("init") == nil {
			continue
		} */
		appPkgs = append(appPkgs, ssaPkg)
	}
	ssaAnalysis(prog, appPkgs)
}

func iterateFunc(outFile *os.File, fn *ssa.Function) {
	fmt.Fprintf(outFile, "\tFunction: %s\n", fn.Name())
	for i, block := range fn.Blocks {
		fmt.Fprintf(outFile, "\t\tBlock #%d: %s.%s\n", i, fn.Name(), block.Comment)
		for j, instr := range block.Instrs {
			// check if the instruction is also a Value (i.e. has a result)
			if val, ok := instr.(ssa.Value); ok {
				fmt.Fprintf(outFile, "\t\t\t%02d: %s = %s\n", j, val.Name(), instr.String())
			} else {
				fmt.Fprintf(outFile, "\t\t\t%02d: %s\n", j, instr.String())
			}
		}
	}
}

func ssaAnalysis(prog *ssa.Program, pkgs []*ssa.Package) {
	outFile, err := os.Create(outfilename)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	for _, ssaPkg := range pkgs {
		outfile, err := os.Create(fmt.Sprintf("%s.ssa", ssaPkg.Pkg.Name()))
		if err != nil {
			log.Fatalf("failed to create output file: %v", err)
		}
		defer outfile.Close()
		ssaPkg.WriteTo(outfile)

		for _, member := range ssaPkg.Members {
			switch m := member.(type) {
			case *ssa.Function:
				iterateFunc(outFile, m)

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
						iterateFunc(outFile, method)
					}
				}

				methods := prog.MethodSets.MethodSet(m.Type().Underlying())
				for i := 0; i < methods.Len(); i++ {
					sel := methods.At(i)
					fmt.Fprintf(outFile, "\tMethod: %v\n", sel.Obj().Type())
					method := prog.MethodValue(sel)
					if method != nil {
						iterateFunc(outFile, method)
					}
				}

			default:
				fmt.Fprintf(outFile, "\tUnknown member type: %T\n", m)
			}
		}
	}
}
