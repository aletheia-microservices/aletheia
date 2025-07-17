package main

import (
	"fmt"
	"go/token"
	"io/ioutil"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

type Node struct { // objects
	name      string
	str       string
	fn        string
	pointedBy []string //variables
	edgesTo   []*Edge
}

func (node *Node) String() string {
	return fmt.Sprintf("(%s) %s: %s", node.fn, node.name, node.str)
}

func (node *Node) LongString() string {
	return fmt.Sprintf("(%s) %v: %s", node.fn, node.pointedBy, node.str)
}

type Edge struct {
	from    *Node
	to      *Node
	isField bool
	isIndex bool
	field   int
	index   string
}

func (edge *Edge) String() string {
	if edge.isField {
		return fmt.Sprintf("%s --(field = %d)--> %s", edge.from.LongString(), edge.field, edge.to.LongString())
	}
	if edge.isIndex {
		return fmt.Sprintf("%s --(index = %s)--> %s", edge.from.LongString(), edge.index, edge.to.LongString())
	}
	return fmt.Sprintf("%s --> %s", edge.from.LongString(), edge.to.LongString())
}

type Graph struct {
	nodes []*Node
	edges []*Edge
}

var graph = Graph{}

func getNode(v ssa.Value) *Node {
	for _, node := range graph.nodes {
		if (node.name == v.Name() || slices.Contains(node.pointedBy, v.Name())) && node.str == v.String() {
			return node
		}
	}
	return nil
}

func edgeExists(from *Node, to *Node, fieldOrIndex bool, fieldVal int, indexVal string) bool {
	for _, edge := range graph.edges {
		if edge.from.String() == from.String() && edge.to.String() == to.String() && (edge.isField == fieldOrIndex || edge.isIndex == fieldOrIndex) && (edge.field == fieldVal || edge.index == indexVal) {
			return true
		}
	}
	return false
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
	return fmt.Sprintf("%s (@%s): (%T)", v.Name(), fn.Name(), v)
}

func save(v ssa.Value) {
	switch vt := v.(type) {
	case *ssa.Call:
		for _, arg := range vt.Call.Args {
			fmt.Printf("ON CALL: %s\n", vt.Call.String())
			for _, node := range graph.nodes {
				if slices.Contains(node.pointedBy, arg.Name()) {
					fmt.Printf("\t NODE: %s", node.LongString())
					for _, edge := range node.edgesTo {
						fmt.Printf("\t\t EDGE: %s", edge.String())
					}
				}
			}
		}
	case *ssa.Alloc:
		node := &Node{
			name:      vt.Name(),
			str:       vt.String(),
			fn:        v.Parent().Name(),
			pointedBy: []string{vt.Name()},
		}
		graph.nodes = append(graph.nodes, node)
	case *ssa.Slice:
		node := getNode(vt.X)
		if node == nil {
			save(vt.X)
		}
		node = getNode(vt.X)
		node.pointedBy = append(node.pointedBy, v.String())
	case *ssa.FieldAddr:
		node := getNode(v)
		if node == nil {
			node = &Node{
				name:      vt.Name(),
				str:       vt.String(),
				pointedBy: []string{vt.Name()},
			}
		}
		for _, gnode := range graph.nodes {
			if gnode.name == vt.X.Name() && gnode.fn == vt.X.Parent().Name() {
				if !edgeExists(gnode, node, true, vt.Field, "") {
					edge := &Edge{
						from:    gnode,
						to:      node,
						isField: true,
						field:   vt.Field,
					}
					gnode.edgesTo = append(gnode.edgesTo, edge)
					graph.edges = append(graph.edges, edge)
				}
			}
		}
	case *ssa.IndexAddr:
		node := getNode(v)
		if node == nil {
			node = &Node{
				name:      vt.Name(),
				str:       vt.String(),
				pointedBy: []string{vt.Name()},
			}
		}
		for _, gnode := range graph.nodes {
			if gnode.name == vt.X.Name() && gnode.fn == vt.X.Parent().Name() {
				if !edgeExists(gnode, node, true, 0, vt.Index.String()) {
					edge := &Edge{
						from:    gnode,
						to:      node,
						isIndex: true,
						index:   vt.Index.String(),
					}
					gnode.edgesTo = append(gnode.edgesTo, edge)
					graph.edges = append(graph.edges, edge)
				}
			}
		}

	}
}

func main() {
	filename := "../examples-go-simple/main.go"
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read file %s: %v\n", filename, err)
		os.Exit(1)
	}

	// setup loader
	var conf loader.Config
	fset := token.NewFileSet()
	conf.Fset = fset
	file, err := conf.ParseFile(filename, string(source))
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}
	conf.CreateFromFiles("main", file)

	iprog, err := conf.Load()
	if err != nil {
		fmt.Println("type error:", err)
		return
	}

	prog := ssautil.CreateProgram(iprog, 0)
	mainPkg := prog.Package(iprog.Created[0].Pkg)
	prog.Build()

	config := &pointer.Config{
		Mains:          []*ssa.Package{mainPkg},
		BuildCallGraph: true,
	}

	for fn := range ssautil.AllFunctions(prog) {
		if fn == nil || fn.Pkg == nil || fn.Pkg.Pkg.Path() != "main" {
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
			for i, instr := range b.Instrs {
				fmt.Printf("%d: analyzing inst [%T]: %v\n", i, instr, instr)
				switch v := instr.(type) {
				case ssa.Value:
					switch vv := v.(type) {
					case *ssa.MakeInterface:
						fmt.Printf("MakeInterface → %v\n", vv)
						fmt.Printf("  X (inner value) → [%T] %v\n", vv.X, vv.X)
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

	pkgPath := mainPkg.Pkg.Path()

	callGraph(pkgPath, result)
	pointsTo(prog, pkgPath, result)

	/* for fn := range ssautil.AllFunctions(prog) {
		if fn == nil || fn.Pkg == nil || fn.Pkg.Pkg.Path() != "main" {
			continue
		}
		fmt.Printf("### Function: %s ###\n", fn.String())
		for _, b := range fn.Blocks {
			for i, instr := range b.Instrs {
				switch v := instr.(type) {
				case *ssa.Store:
					fmt.Printf("%d: analyzing inst [%T]: %v\n", i, instr, instr)
					fmt.Printf("\tREFERRERS = %v\n\n", v.Referrers())
				case *ssa.UnOp:
					fmt.Printf("%d: analyzing inst [%T]: %v\n", i, instr, instr)
					fmt.Printf("\tREFERRERS = %v\n\n", v.Referrers())

				}
			}
		}
	} */

	fmt.Println("\nGRAPH:\n")
	for _, node := range graph.nodes {
		for _, edge := range graph.edges {
			if edge.from == node {
				fmt.Println(edge.String())
			}
		}
	}

	//pointsToReverse(prog, pkgPath, result)
}

func callGraph(pkgPath string, result *pointer.Result) {
	fmt.Println("Call graph edges:")
	var edges []string
	callgraph.GraphVisitEdges(result.CallGraph, func(edge *callgraph.Edge) error {
		caller := edge.Caller.Func
		if caller.Pkg == nil || caller.Pkg.Pkg.Path() != pkgPath {
			return nil
		}
		edges = append(edges, fmt.Sprintf("%s --> %s", caller, edge.Callee.Func))
		return nil
	})
	sort.Strings(edges)
	for _, edge := range edges {
		fmt.Println("  ", edge)
	}
}

func pointsTo(prog *ssa.Program, pkgPath string, result *pointer.Result) {
	fmt.Println("\nPoints-to sets by function:\n")

	type ptsEntry struct {
		pos  token.Position
		desc string
		name string
	}

	entriesByFunc := make(map[*ssa.Function][]ptsEntry)

	for value, pts := range result.Queries {
		fn := valueParent(value)
		if fn == nil || fn.Pkg == nil || fn.Pkg.Pkg.Path() != pkgPath {
			continue
		}

		pos := prog.Fset.Position(value.Pos())
		desc := valueDesc(fn, value) + "\n"
		save(value)
		name := value.Name()
		switch t := value.(type) {
			case *ssa.Call:
				for i, arg := range t.Call.Args {
					desc += fmt.Sprintf("\t\targ #%d → %s\n", i, valueDescShort(arg.Parent(), arg))
				}
			case *ssa.Parameter:
				name = "_" + value.Name()
			case *ssa.Phi:
				for i, edge := range t.Edges {
					desc += fmt.Sprintf("\t\tphi(#%d) ← %s\n", i, valueDescShort(edge.Parent(), edge))
				}
		}
		for _, lbl := range pts.PointsTo().Labels() {
			//poslbl := prog.Fset.Position(lbl.Pos())
			desc += fmt.Sprintf("\t → %s\n", valueDescShort(lbl.Value().Parent(), lbl.Value()))
			switch t := lbl.Value().(type) {
			case *ssa.MakeInterface:
				desc += fmt.Sprintf("\t\tX (interface of) → %s\n", valueDescShort(fn, t.X))
			/* case *ssa.UnOp:
				desc += fmt.Sprintf("\t\tX (loaded from) → %s [%T] %v\n", t.X.Name(), t.X, t.X) */
			}
			save(lbl.Value())
		}
		entriesByFunc[fn] = append(entriesByFunc[fn], ptsEntry{pos, desc, name})
	}

	var fns []*ssa.Function
	for fn, _ := range entriesByFunc {
		fns = append(fns, fn)
	}

	for _, fn := range fns {
		for _, block := range fn.Blocks {
			for _, instr := range block.Instrs {
				var name string
				var desc string
				var pos token.Position
				if val, ok := instr.(*ssa.UnOp); ok {
					name = val.Name()
					desc = valueDesc(val.Parent(), val) + "\n"
					pos = prog.Fset.Position(val.Pos())
				}
				if store, ok := instr.(*ssa.Store); ok {
					name = store.Addr.Name()
					desc = valueDesc(store.Addr.Parent(), store.Addr) + "\n"
					pos = prog.Fset.Position(store.Addr.Pos())
				}
				if store, ok := instr.(ssa.Value); ok {
					name = store.Name()
					desc = valueDesc(store.Parent(), store) + "\n"
					pos = prog.Fset.Position(store.Pos())
				}
				if name == "" {
					continue
				}
				var found bool
				for _, entry := range entriesByFunc[fn] {
					if entry.name == name {
						found = true
						break
					}
				}
				if !found {
					entriesByFunc[fn] = append(entriesByFunc[fn], ptsEntry{pos, desc, name})
				}
			}
		}
	}

	var funcs []*ssa.Function
	for fn := range entriesByFunc {
		funcs = append(funcs, fn)
	}
	sort.Slice(funcs, func(i, j int) bool {
		return funcs[i].String() < funcs[j].String()
	})

	for _, fn := range funcs {
		fmt.Printf("### Function: %s ###\n", fn.String())
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
			fmt.Println(e.desc)
		}
		fmt.Println()
	}
}

func pointsToReverse(prog *ssa.Program, pkgPath string, result *pointer.Result) {
	fmt.Println("\nReverse Points-to sets:")

	type labelEntry struct {
		pos   token.Position
		label *pointer.Label
		from  []ssa.Value
	}

	labelToPointers := make(map[*pointer.Label][]ssa.Value)

	for value, pts := range result.Queries {
		fn := valueParent(value)
		if fn == nil || fn.Pkg == nil || fn.Pkg.Pkg.Path() != pkgPath {
			continue
		}
		for _, label := range pts.PointsTo().Labels() {
			labelToPointers[label] = append(labelToPointers[label], value)
		}
	}

	// collect and sort by source position
	var entries2 []labelEntry
	for label, values := range labelToPointers {
		pos := prog.Fset.Position(label.Pos())
		entries2 = append(entries2, labelEntry{pos, label, values})
	}

	sort.Slice(entries2, func(i, j int) bool {
		a, b := entries2[i].pos, entries2[j].pos
		if a.Filename != b.Filename {
			return a.Filename < b.Filename
		}
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		return a.Column < b.Column
	})

	for _, e := range entries2 {
		fmt.Printf("%s: %s\n", e.pos, e.label)
		for _, v := range e.from {
			fmt.Printf("    ← %v\n", v)
		}
	}
}

// valueParent safely gets the *ssa.Function that owns a value
func valueParent(v ssa.Value) *ssa.Function {
	switch val := v.(type) {
	case *ssa.Function:
		return val
	default:
		return v.Parent()
	}
}
