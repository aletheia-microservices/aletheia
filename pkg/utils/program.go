package utils

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
)

// BuildProgram loads the packages of the app at apppath and builds their SSA program, returning
// the program and the SSA packages of the app
func BuildProgram(apppath string) (*ssa.Program, []*ssa.Package, error) {
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
		if IsAppPackagePath(progPkg.Pkg.Path()) && progPkg != nil && !pkgsSeen[progPkg] {
			pkgsSeen[progPkg] = true
			pkgs = append(pkgs, progPkg)
		}
	}
	return prog, pkgs, nil
}
