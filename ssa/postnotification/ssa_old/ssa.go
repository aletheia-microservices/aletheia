package main

import (
	"fmt"
	"go/token"
	"log"
	"os"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
)

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
	pkgs, err := packages.Load(cfg, "../workflow/postnotification/")
	if err != nil {
		log.Fatal(err)
	}

	fset := token.NewFileSet()
	prog := ssa.NewProgram(fset, ssa.PrintFunctions)

	ssaPkgs := make([]*ssa.Package, len(pkgs))
	for i, pkg := range pkgs {
		if _, ok := createdPkgs[pkg]; ok {
			continue
		}
		ssaPkgs[i] = prog.CreatePackage(pkg.Types, pkg.Syntax, pkg.TypesInfo, false)
		createdPkgs[pkg] = true
		for _, impt := range pkg.Imports {
			recurse(prog, impt)
		}
	}

	prog.Build()

	outFile, err := os.Create("ssa.out")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	for _, ssaPkg := range ssaPkgs {
		if ssaPkg == nil || ssaPkg.Pkg == nil {
			continue
		}
		if ssaPkg.Pkg.Name() != "postnotification" {
			continue
		}

		fmt.Fprintf(outFile, "Package: %s\n", ssaPkg.Pkg.Name())

		for _, member := range ssaPkg.Members {
			switch m := member.(type) {
			case *ssa.Function:
				fmt.Fprintf(outFile, "\tFunction: %s\n", m.Name())
				for i, block := range m.Blocks {
					fmt.Fprintf(outFile, "\t\tBlock #%d: %s\n", i, block.Comment)
					for j, instr := range block.Instrs {
						fmt.Fprintf(outFile, "\t\t\tInst #%d: %s\n", j, instr.String())
					}
				}

			case *ssa.Global:
				fmt.Fprintf(outFile, "\tGlobal: %s, Type: %s\n", m.Name(), m.Type().String())

			case *ssa.Type:
				fmt.Fprintf(outFile, "\tType: %s\n", m.Type().String())

				// print methods on the type
				methods := prog.MethodSets.MethodSet(m.Type())
				for i := 0; i < methods.Len(); i++ {
					sel := methods.At(i)
					fmt.Fprintf(outFile, "\t\tMethod: %t // %v\n", sel.Obj(), sel.Obj().Type())

					method := prog.MethodValue(sel)
					if method != nil {
						fmt.Fprintf(outFile, "\t\tMethod: %s\n", method.String())
						for i, block := range method.Blocks {
							fmt.Fprintf(outFile, "\t\t\tBlock #%d: %s\n", i, block.Comment)
							for j, instr := range block.Instrs {
								fmt.Fprintf(outFile, "\t\t\t\tInst #%d: %s\n", j, instr.String())
							}
						}
					}
				}

			default:
				fmt.Fprintf(outFile, "\tUnknown member type: %T\n", m)
			}
		}
	}
}

func printFunction(fn *ssa.Function, out *os.File) {
	fmt.Fprintf(out, "\tFunction: %s\n", fn.Name())
	for i, block := range fn.Blocks {
		fmt.Fprintf(out, "\t  Block #%d: %s\n", i, block.Comment)
		for j, instr := range block.Instrs {
			fmt.Fprintf(out, "\t    Inst #%d: %s\n", j, instr.String())
		}
	}
}
