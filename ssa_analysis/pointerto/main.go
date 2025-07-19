package main

import (
	"fmt"
	"go/token"
	"go/types"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"golang.org/x/tools/go/types/typeutil"
)

// -------------------------------------
// ------------- CONSTANTS -------------
// -------------------------------------

var appname = "test2"

// -------------------------------------

type Graph struct {
	nodes []*Node
	edges []*Edge
	defs  map[string]*Node
}

func (graph *Graph) addEdge(edge *Edge) {
	graph.edges = append(graph.edges, edge)
}

func (graph *Graph) sortNodes() {
	sort.Slice(graph.nodes, func(i, j int) bool {
		/* ni, err1 := strconv.Atoi(strings.TrimPrefix(graph.nodes[i].name, "t"))
		nj, err2 := strconv.Atoi(strings.TrimPrefix(graph.nodes[j].name, "t"))
		if err1 != nil || err2 != nil {
			return graph.nodes[i].name < graph.nodes[j].name
		}
		return ni < nj */
		return graph.nodes[i].id < graph.nodes[j].id
	})
}

type Node struct {
	id    int
	name  string
	val   ssa.Value
	instr ssa.Instruction
	isdef bool

	// maps object to database field, e.g.:
	// key: Product    // value: prod_db.Product
	// key: Product.ID // value: prod_db.Product.ID
	// key: Product.ID // value: sku_db.Sku.ProductID
	taints map[string][]string
}

func (node *Node) isTainted() bool {
	return len(node.taints) > 0
}

func (node *Node) taintString() string {
	if len(node.taints) == 0 {
		return ""
	}
	var taintStr string
	for obj, dbfields := range node.taints {
		taintStr += fmt.Sprintf("%s: ", obj)
		for _, dbfield := range dbfields {
			taintStr += fmt.Sprintf("\t|--> %s\n", dbfield)
		}
	}
	return taintStr
}

func (node *Node) String() string {
	if node.val != nil {
		return node.name + ": " + node.val.String()
	}
	return node.instr.String()
}

func (node *Node) colorForSSA() string {
	switch node.instr.(type) {
	case *ssa.Store:
		return "blue"
	case *ssa.Alloc:
		return "orange"
	case *ssa.Return:
		return "yellow"
	case *ssa.Call:
		return "yellow"
	case *ssa.UnOp:
		return "red"
	case *ssa.FieldAddr, *ssa.IndexAddr:
		return "green"
	}
	return "black"
}

func RegisterNewNodeValue(graph *Graph, instr ssa.Instruction, val ssa.Value, id int) *Node {
	node := &Node{
		name:   val.Name(),
		val:    val,
		instr:  instr,
		isdef:  true,
		id:     id,
		taints: make(map[string][]string),
	}
	graph.nodes = append(graph.nodes, node)
	graph.defs[node.name] = node
	return node
}

func RegisterNewNode(graph *Graph, instr ssa.Instruction, id int) *Node {
	node := &Node{
		id:    id,
		instr: instr,
	}
	graph.nodes = append(graph.nodes, node)
	return node
}

type EdgeType int

const (
	EDGE_USAGE EdgeType = iota
	EDGE_STORE
	EDGE_LOAD
	EDGE_FIELD
	EDGE_INDEX
	EDGE_PARAMETER
	EDGE_POINTS_TO
	/* EDGE_VALUE
	EDGE_VERSION
	EDGE_INTERFACE
	EDGE_CONVERTED
	EDGE_PHI
	EDGE_RETURN
	EDGE_COPY */
)

type Edge struct {
	edgeType EdgeType
	from     *Node
	to       *Node

	index int
	param string

	path string //pointer only
}

func newEdge(from *Node, to *Node, edgeType EdgeType, index int, param string) *Edge {
	return &Edge{
		from:     from,
		to:       to,
		edgeType: edgeType,
		index:    index,
		param:    param,
	}
}

func (edge *Edge) typeString() string {
	switch edge.edgeType {
	case EDGE_USAGE:
		return "usage"
	case EDGE_STORE:
		return "store"
	case EDGE_LOAD:
		return "load"
	case EDGE_FIELD:
		return fmt.Sprintf("field(%s)", edge.param)
	case EDGE_INDEX:
		return fmt.Sprintf("index(%s)", edge.param)
	case EDGE_PARAMETER:
		return fmt.Sprintf("param(%d)", edge.index)
	}
	return ""
}

func (edge *Edge) HasFromNode(node *Node) bool {
	return edge.from == node
}

func isMongoDBInsert(instr ssa.Instruction) (*ssa.Call, *ssa.Function, bool) {
	if call, ok := instr.(*ssa.Call); ok {
		if fn, ok := call.Call.Value.(*ssa.Function); ok && len(fn.Params) > 0 {
			maybeRcv := fn.Params[0]
			if maybeRcv.Type().String() == "*main.MongoDB" && fn.Name() == "Insert" {
				return call, fn, true
			}
		}
	}
	return nil, nil, false
}

func doTaintNode(graph *Graph, val ssa.Value, dbField string) {
	node := graph.defs[val.Name()]
	node.taints["."] = append(node.taints["."], dbField)
}

func fieldIndexToName(t *ssa.FieldAddr) string {
	return t.X.Type().Underlying().(*types.Pointer).Elem().(*types.Named).Underlying().(*types.Struct).Field(t.Field).Name()
}

func doTaintBackwards(graph *Graph, val ssa.Value, dbField string) {
	switch t := val.(type) {
	case *ssa.MakeInterface:
		doTaintBackwards(graph, t.X, dbField)
	case *ssa.UnOp:
		doTaintBackwards(graph, t.X, dbField)
	case *ssa.FieldAddr:
		fieldName := fieldIndexToName(t)
		dbField += "." + fieldName
		doTaintBackwards(graph, t.X, dbField)
	case *ssa.Parameter, *ssa.Alloc:
		doTaintNode(graph, val, dbField)
	default:
		fmt.Printf("[INFO] ignoring value: [%T] %v\n", val, val)
	}
}

func doTaint(graph *Graph) {
	for _, node := range graph.nodes {
		call, _, ok := isMongoDBInsert(node.instr)
		if ok {
			var database, collection string

			valDatabase := call.Call.Args[2]
			if c, ok := valDatabase.(*ssa.Const); ok {
				database = strings.Trim(c.Value.ExactString(), "\"")
			}

			valCollection := call.Call.Args[3]
			if c, ok := valCollection.(*ssa.Const); ok {
				collection = strings.Trim(c.Value.ExactString(), "\"")
			}

			valDocument := call.Call.Args[4]

			field := database + "." + collection
			doTaintBackwards(graph, valDocument, field)
		}
	}
}

func parseInstr(graph *Graph, instr ssa.Instruction, id int) {
	fmt.Printf("[A] %02d [%T] %v\n", id, instr, instr.String())

	if val, ok := instr.(ssa.Value); ok {
		parseValue(graph, instr, id, val)
		return
	}
	node := RegisterNewNode(graph, instr, id)

	switch t := instr.(type) {
	case *ssa.Store:
		// 04 [store] *t1 = currency
		addrNode := parseValue(graph, instr, id, t.Addr)
		valNode := parseValue(graph, instr, id, t.Val)

		edge1 := newEdge(addrNode, node, EDGE_STORE, 0, "")
		graph.addEdge(edge1)

		if valNode != nil { // nil whever an *ssa.Const is parsed
			edge2 := newEdge(valNode, node, EDGE_USAGE, 0, "")
			graph.addEdge(edge2)
		}
	case *ssa.Return:
		for _, res := range t.Results {
			resNode := parseValue(graph, instr, id, res)
			if resNode != nil { // nil whever an *ssa.Const is parsed
				edge := newEdge(resNode, node, EDGE_STORE, 0, "")
				graph.addEdge(edge)
			}
		}
	case *ssa.Jump:
		//fmt.Printf("[A] skipping... %02d [%T] %v\n", id, instr, instr.String())
	default:
		//fmt.Printf("[1] ignoring... %02d [%T] %v\n", id, instr, instr.String())
	}
}

func parseValue(graph *Graph, instr ssa.Instruction, id int, val ssa.Value) *Node {
	fmt.Printf("[B] %02d [%T] %v\n", id, val, val.String())

	if _, ok := val.(*ssa.Const); ok {
		return nil
	}

	if node, exists := graph.defs[val.Name()]; exists {
		return node
	}

	node := RegisterNewNodeValue(graph, instr, val, id)

	switch t := val.(type) {
	case *ssa.Call:
		for _, arg := range t.Call.Args {
			argNode := parseValue(graph, instr, id, arg)
			if argNode != nil { // nil whever an *ssa.Const is parsed
				edge := newEdge(argNode, node, EDGE_STORE, 0, "")
				graph.addEdge(edge)
			}
		}
	case *ssa.Alloc:
		// nothing to do
	case *ssa.Slice:
		// nothing to do
	case *ssa.FieldAddr:
		// 00 [field] t27 = &t0.Items [#3]
		targetNode := parseValue(graph, instr, id, t.X)
		edge := newEdge(targetNode, node, EDGE_FIELD, 0, fieldIndexToName(t))
		graph.addEdge(edge)
	case *ssa.IndexAddr:
		targetNode := parseValue(graph, instr, id, t.X)
		//FIXME: should parse value for t.Index
		edge := newEdge(targetNode, node, EDGE_FIELD, 0, t.Index.String())
		graph.addEdge(edge)
	case *ssa.UnOp:
		// 01 [unary] t14 = *t13
		// 05 [unary] t31 = *t30
		targetNode := parseValue(graph, instr, id, t.X)
		edge := newEdge(targetNode, node, EDGE_LOAD, 0, "")
		graph.addEdge(edge)

	case *ssa.MakeInterface: // same as *ssa.UnOp
		targetNode := parseValue(graph, instr, id, t.X)
		edge := newEdge(targetNode, node, EDGE_USAGE, 0, "")
		graph.addEdge(edge)
	case *ssa.Convert:
		targetNode := parseValue(graph, instr, id, t.X)
		edge := newEdge(targetNode, node, EDGE_USAGE, 0, "")
		graph.addEdge(edge)

	case *ssa.Parameter:
		// nothing to do

	case *ssa.Global:
		// nothing to do

	case *ssa.Phi:
		for _, phiEdge := range t.Edges {
			edgeNode := parseValue(graph, instr, id, phiEdge)
			edge := newEdge(edgeNode, node, EDGE_STORE, 0, "")
			graph.addEdge(edge)
		}

	case *ssa.Extract:
		extractFromNode := parseValue(graph, instr, id, t.Tuple)
		edge := newEdge(extractFromNode, node, EDGE_USAGE, t.Index, "")
		graph.addEdge(edge)

	case *ssa.BinOp:
		xNode := parseValue(graph, instr, id, t.X)
		yNode := parseValue(graph, instr, id, t.Y)

		edge1 := newEdge(xNode, node, EDGE_STORE, 0, "")
		graph.addEdge(edge1)

		edge2 := newEdge(yNode, node, EDGE_USAGE, 0, "")
		graph.addEdge(edge2)

	default:
		//fmt.Printf("[2] ignoring... %02d [%T] %s = %v\n", id, val, val.Name(), val.String())
	}
	return node
}

func (g *Graph) writeToDOTFile(fn string) error {
	filename := fmt.Sprintf("%s.dot", fn)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "digraph G {")
	fmt.Fprintln(file, "\trankdir=TD;")

	for _, node := range g.nodes {
		str := node.String() + "\n" + node.taintString()
		label := strings.ReplaceAll(str, `"`, `\"`)
		nodecolor := node.colorForSSA()

		shape := "box"
		if node.isTainted() {
			shape = "ellipse"
		}

		color := "black"
		if nodecolor != "" {
			color = nodecolor
		}

		fmt.Fprintf(file, "\tN%d [label=\"%s\", style=bold, shape=%s, color=\"%s\"];\n", node.id, label, shape, color)
	}

	for _, edge := range g.edges {
		if edge.edgeType == EDGE_POINTS_TO {
			path := strings.ReplaceAll(edge.path, `"`, `\"`)
			fmt.Fprintf(file, "\tN%d -> N%d [label=\"%s\", style=dashed, color=blue];\n", edge.from.id, edge.to.id, path)
		} else {
			fmt.Fprintf(file, "\tN%d -> N%d;\n", edge.from.id, edge.to.id)
		}
	}

	fmt.Fprintln(file, "}")
	return nil
}

func main() {
	if len(os.Args) >= 2 {
		appname = os.Args[1]
	}

	prog, pkg, result, err := initPackages()
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	funcGraphs := make(map[string]*Graph)

	runSSAAnalysis(prog, pkg, funcGraphs)

	for _, graph := range funcGraphs {
		graph.sortNodes()
	}

	runPointerToAnalysis(prog, pkg, result, funcGraphs)

	for _, graph := range funcGraphs {
		doTaint(graph)
	}

	fmt.Print("\n\n ========== NODES ========== \n\n")
	for fn, graph := range funcGraphs {
		for _, node := range graph.nodes {
			var prefix string
			if node.name != "" {
				prefix = node.name + ":"
			} else {
				prefix = "\t"
			}
			if node.instr != nil {
				fmt.Printf("[%s] [%02d] [%T] \t %s %v\n", fn, node.id, node.instr, prefix, node.instr.String())
			} else {
				fmt.Printf("[%s] [%02d] [%T] \t %s %v\n", fn, node.id, node.val, prefix, node.val.String())
			}
		}
	}

	fmt.Print("\n\n ========== TAINTS ========== \n\n")
	for fn, graph := range funcGraphs {
		for _, node := range graph.nodes {
			if node.isTainted() {
				for obj, dbfields := range node.taints {
					fmt.Printf("[%s] %s [%s]: %s\n", fn, node.val.String(), node.name, obj)
					for _, dbfield := range dbfields {
						fmt.Printf("\t\t |--> %s\n", dbfield)
					}
				}
			}
		}
	}

	for fn, graph := range funcGraphs {
		graph.writeToDOTFile(fn)
	}

	fmt.Println("\n[INFO] successfully analyzed app (" + appname + ")\n")
}

func runSSAAnalysis(prog *ssa.Program, pkg *ssa.Package, funcGraphs map[string]*Graph) {
	outFile, err := os.Create("app.ssa")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	outfile, err := os.Create(fmt.Sprintf("%s.out", pkg.Pkg.Name()))
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outfile.Close()
	pkg.WriteTo(outfile)

	for _, member := range pkg.Members {
		switch m := member.(type) {
		case *ssa.Function:
			iterateFunc(outFile, m, nil, funcGraphs)

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
					iterateFunc(outFile, method, m.Type(), funcGraphs)
				}
			}

			methods := prog.MethodSets.MethodSet(m.Type().Underlying())
			for i := 0; i < methods.Len(); i++ {
				sel := methods.At(i)
				fmt.Fprintf(outFile, "\tMethod: %v\n", sel.Obj().Type())
				method := prog.MethodValue(sel)
				if method != nil {
					iterateFunc(outFile, method, m.Type(), funcGraphs)
				}
			}

		default:
			fmt.Fprintf(outFile, "\tUnknown member type: %T\n", m)
		}
	}
}

func iterateFunc(outFile *os.File, fn *ssa.Function, memberType types.Type, funcGraphs map[string]*Graph) {
	graph := &Graph{
		defs: make(map[string]*Node),
	}
	if _, exists := funcGraphs[fn.Name()]; exists {
		log.Fatalf("graph for function (%s) already exists\n", fn.Name())
	}
	funcGraphs[fn.Name()] = graph
	fmt.Printf("added new graph for function (%s)\n", fn.Name())

	for i, param := range fn.Params {
		parseValue(graph, nil, -i-1, param)
	}

	fmt.Fprintf(outFile, "Function: %s\n", fn.Name())
	for i, block := range fn.Blocks {
		fmt.Fprintf(outFile, "Block #%d: %s.%s\n", i, fn.Name(), block.Comment)
		for j, instr := range block.Instrs {
			parseInstr(graph, instr, j)
		}
	}
}

func initPackages() (*ssa.Program, *ssa.Package, *pointer.Result, error) {
	// e.g. "../examples/test2/main.go"
	filepath := "../examples/" + appname + "/main.go"
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
		return nil, nil, nil, err
	}
	conf.CreateFromFiles("main", file)

	iprog, err := conf.Load()
	if err != nil {
		fmt.Println("type error:", err)
		return nil, nil, nil, err
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
	return prog, mainPkg, result, nil
}

func runPointerToAnalysis(prog *ssa.Program, pkg *ssa.Package, result *pointer.Result, funcGraphs map[string]*Graph) {
	outFile, err := os.Create("app.ptrs")
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
		if fn == nil || fn.Pkg == nil || fn.Pkg.Pkg.Path() != pkg.Pkg.Path() {
			continue
		}

		graph := funcGraphs[fn.Name()]

		pos := prog.Fset.Position(value.Pos())
		desc := valueDesc(fn, value) + "\n"
		name := value.Name()
		node := graph.defs[name]
		for _, lbl := range pts.PointsTo().Labels() {
			desc += fmt.Sprintf("\t → %s\n", valueDescShort(lbl.Value().Parent(), lbl.Value()))

			if lbl.Value().Parent() == fn {
				pointsToNode := graph.defs[lbl.Value().Name()]

				if node != nil && pointsToNode != nil && node != pointsToNode {
					var exists bool
					for _, edge := range graph.edges {
						// this is reverse on purpose for field and index addresses
						//if edge.from == pointsToNode && edge.to == node {
						if edge.from == node && edge.to == pointsToNode {
							exists = true
						}
					}
					if !exists {
						edge := newEdge(node, pointsToNode, EDGE_POINTS_TO, 0, "")
						edge.path = lbl.Path()
						graph.addEdge(edge)
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
