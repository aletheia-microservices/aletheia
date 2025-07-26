package parser

import (
	"fmt"
	"go/token"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"

	"analyzer/pkg/ssagraph"
	"analyzer/pkg/utils"
)

func InitPointerAnalysis(prog *ssa.Program, pkgs []*ssa.Package) (*pointer.Result, error) {

	config := &pointer.Config{
		Mains:          pkgs,
		BuildCallGraph: true,
	}

	for fn := range ssautil.AllFunctions(prog) {
		if fn == nil || fn.Pkg == nil || !slices.Contains(pkgs, fn.Pkg) {
			continue
		}
		for _, param := range fn.Params {
			if pointer.CanPoint(param.Type()) {
				config.AddQuery(param)
			}
		}
		for _, fv := range fn.FreeVars {
			if pointer.CanPoint(fv.Type()) {
				config.AddQuery(fv)
			}
		}
		for _, lcl := range fn.Locals {
			if pointer.CanPoint(lcl.Type()) {
				config.AddQuery(lcl)
			}
		}
		for _, b := range fn.Blocks {
			for _, instr := range b.Instrs {
				switch v := instr.(type) {
				case ssa.Value:
					switch vv := v.(type) {
					case *ssa.MakeInterface:
						if pointer.CanPoint(vv.X.Type()) {
							config.AddQuery(vv.X)
						}
						if pointer.CanPoint(vv.Type()) {
							config.AddQuery(vv)
						}
					default:
						if pointer.CanPoint(v.Type()) {
							config.AddQuery(v)
						}
					}
				case *ssa.Store:
					if pointer.CanPoint(v.Addr.Type()) {
						config.AddQuery(v.Addr)
					}

				case *ssa.Return:
					for _, r := range v.Results {
						if pointer.CanPoint(r.Type()) {
							config.AddQuery(r)
						}
					}
				}
			}
		}
	}

	result, err := pointer.Analyze(config)
	if err != nil {
		panic(err)
	}
	return result, nil
}

func RunPointerToAnalysis(appname string, prog *ssa.Program, pkg *ssa.Package, result *pointer.Result, funcGraphs map[string]*ssagraph.SSAGraph) {
	fmt.Printf("\n[PTA] running pointer analysis for package: %s\n", pkg.String())
	outFile, err := os.Create(fmt.Sprintf("output/%s/%s.ptrs", appname, pkg.Pkg.Name()))
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	type ptsEntry struct {
		pos  token.Position
		desc string
		name string
	}

	entriesByFunc := make(map[*ssa.Function][]ptsEntry)

	for value, pts := range result.Queries {
		fn := valueParent(value)
		if fn == nil {
			continue
		}

		shortFuncPath := utils.GetShortFunctionPath(fn.String())
		fmt.Println()
		fmt.Printf("\t[PTA] [%s] analyzing value: %v // pointers = %v\n", shortFuncPath, value, pts)
		if fn.Pkg == nil {
			continue
		}

		graph := funcGraphs[shortFuncPath]
		if graph == nil {
			fmt.Printf("skipping graph not found for name (%s)\n", shortFuncPath)
			continue
		}

		pos := prog.Fset.Position(value.Pos())
		desc := valueDesc(fn, value) + "\n"
		name := value.Name()
		node, ok := graph.GetNodeByNameIfExists(name)
		if !ok {
			fmt.Printf("skipping node not found for name (%s)\n", name)
			continue
		}
		//fmt.Printf("points to set of [%T] %v @ %v:\n", value, value, value.Parent())
		for _, lbl := range pts.PointsTo().Labels() {
			lblFn := lbl.Value().Parent()
			lblFnLongName := utils.GetShortFunctionPath(lblFn.String())
			fmt.Printf("\t\t- [%s] got label: [%T] %s\n", lblFnLongName, lbl.Value(), lbl.Value())
			desc += fmt.Sprintf("\t → %s\n", valueDescShort(lbl.Value().Parent(), lbl.Value()))

			if lbl.Value().Parent() == fn {
				pointsToNode, _ := graph.GetNodeByNameIfExists(lbl.Value().Name())

				if node != nil && pointsToNode != nil && node != pointsToNode {
					var exists bool
					/* for _, edge := range graph.GetEdges() {
						// this is reverse on purpose for field and index addresses
						//if edge.from == pointsToNode && edge.to == node {
						if edge.GetFromNode() == node && edge.GetToNode() == pointsToNode {
							exists = true
						}
					} */
					if !exists {
						fmt.Printf("creating edge\n")

						edge, _ := graph.CreateAndAddNewEdge(node, pointsToNode, ssagraph.EDGE_POINTS_TO, 0, "")
						if edge != nil {
							edge.SetPath(lbl.Path())

							fmt.Printf("created edge from: %v\n", edge.GetFromNode())
							for _, edge := range graph.GetEdgesFromNode(edge.GetFromNode()) {
								fmt.Printf("- edge to (%v): %v\n", edge.GetType(), edge.GetToNode().String())
							}
						}
					}
				}
			}
		}
		entriesByFunc[fn] = append(entriesByFunc[fn], ptsEntry{pos, desc, name})
	}

	var funcs []*ssa.Function
	for fn := range entriesByFunc {
		funcs = append(funcs, fn)
	}
	sort.Slice(funcs, func(i, j int) bool {
		return funcs[i].String() < funcs[j].String()
	})

	for _, fn := range funcs {
		fmt.Fprintf(outFile, "### Function: %s ###\n", fn.String())
		entries := entriesByFunc[fn]
		sort.Slice(entries, func(i, j int) bool {
			// sort by name, then position as a fallback
			ni, _ := strconv.Atoi(strings.TrimPrefix(entries[i].name, "t"))
			nj, _ := strconv.Atoi(strings.TrimPrefix(entries[j].name, "t"))
			if ni != nj {
				return ni < nj
			}

			a, b := entries[i].pos, entries[j].pos
			if a.Filename != b.Filename {
				return a.Filename < b.Filename
			}
			if a.Line != b.Line {
				return a.Line < b.Line
			}
			return a.Column < b.Column
		})
		for _, e := range entries {
			fmt.Fprintln(outFile, e.desc)
		}
		fmt.Fprintln(outFile)
	}
}

func valueDesc(fn *ssa.Function, v ssa.Value) string {
	if v == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s (@%s): (%T) (%s)", v.Name(), fn.Name(), v, v)
}

func valueDescShort(fn *ssa.Function, v ssa.Value) string {
	if v == nil {
		return "<nil>"
	}
	fnname := "?"
	if fn != nil {
		fnname = fn.Name()
	}
	return fmt.Sprintf("%s (@%s): (%T)", v.Name(), fnname, v)
}

func valueParent(v ssa.Value) *ssa.Function {
	switch val := v.(type) {
	case *ssa.Function:
		return val
	default:
		return v.Parent()
	}
}
