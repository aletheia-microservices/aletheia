package main

import (
	"fmt"
	"go/token"
	"go/types"
	"log"
	"math/rand"
	"os"
	"slices"
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

var appname = "new3"

// -------------------------------------

type Graph struct {
	nodes []*Node
	edges []*Edge
	defs  map[string]*Node
}

func (graph *Graph) addEdge(edge *Edge) {
	graph.edges = append(graph.edges, edge)
}

func (graph *Graph) getEdgesFromNode(node *Node) []*Edge {
	var edges []*Edge
	for _, edge := range graph.edges {
		if edge.from == node {
			edges = append(edges, edge)
		}
	}
	return edges
}

func (graph *Graph) getEdgesToNode(node *Node) []*Edge {
	var edges []*Edge
	for _, edge := range graph.edges {
		if edge.to == node {
			edges = append(edges, edge)
		}
	}
	return edges
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

func (graph *Graph) getNodeByName(name string) *Node {
	if node, exists := graph.defs[name]; exists {
		return node
	}
	log.Fatalf("node with name (%s) not found in graph defs: %v\n", name, graph.defs)
	return nil
}

type Node struct {
	id    string
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

func (node *Node) addTaintIfNotExists(taintInfo TaintInfo) {
	// note that objfields/dbfields already have "." before them
	objField := taintInfo.objPrefix + taintInfo.objField
	dbField := taintInfo.dbfieldPrefix + taintInfo.dbField
	if !slices.Contains(node.taints[objField], dbField) {
		node.taints[objField] = append(node.taints[objField], dbField)
	}
}

func (node *Node) taintString() string {
	if len(node.taints) == 0 {
		return ""
	}
	var taintStr string
	for obj, dbfields := range node.taints {
		taintStr += fmt.Sprintf("\n%s\n", obj)
		for _, dbfield := range dbfields {
			taintStr += fmt.Sprintf("@ %s\n", dbfield)
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

func RegisterNewNodeValue(graph *Graph, instr ssa.Instruction, val ssa.Value, id string) *Node {
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

func RegisterNewNode(graph *Graph, instr ssa.Instruction, id string) *Node {
	node := &Node{
		id:     id,
		instr:  instr,
		taints: make(map[string][]string),
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
)

type Edge struct {
	edgeType EdgeType
	from     *Node
	to       *Node

	index int
	param string

	path string //pointer only
}

func (graph *Graph) createAndAddNewEdge(from *Node, to *Node, edgeType EdgeType, index int, param string) (*Edge, bool) {
	// 1st is for sanity check; 2nd is for nodes obtained from *ssa.Const
	if from == nil || to == nil {
		return nil, false
	}
	for _, edge := range graph.getEdgesFromNode(from) {
		if edge.to == to {
			return edge, false
		}
	}
	for _, edge := range graph.getEdgesToNode(to) {
		if edge.from == from {
			return edge, false
		}
	}
	edge := &Edge{
		from:     from,
		to:       to,
		edgeType: edgeType,
		index:    index,
		param:    param,
	}
	graph.addEdge(edge)
	return edge, true
}

/* func (edge *Edge) typeString() string {
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
} */

func (edge *Edge) HasFromNode(node *Node) bool {
	return edge.from == node
}

func isMongoDBCall(instr ssa.Instruction) (*ssa.Call, *ssa.Function, bool) {
	if call, ok := instr.(*ssa.Call); ok {
		if fn, ok := call.Call.Value.(*ssa.Function); ok && len(fn.Params) > 0 {
			maybeRcv := fn.Params[0]
			if maybeRcv.Type().String() == "*main.MongoDB" && fn.Name() == "Insert" || fn.Name() == "Find" {
				return call, fn, true
			}
		}
	}
	return nil, nil, false
}

func doTaintNode(graph *Graph, val ssa.Value, taintInfo TaintInfo) {
	node := graph.getNodeByName(val.Name())
	node.addTaintIfNotExists(taintInfo)
}

func doTaintPointerToSets(graph *Graph, val ssa.Value, taintInfo TaintInfo, visited map[TaintInfo]bool) {
	node := graph.getNodeByName(val.Name())
	for _, edge := range graph.getEdgesFromNode(node) {
		if edge.edgeType == EDGE_POINTS_TO {
			if edge.path != "" {
				// add before
				// note that both edge.path and objfields/dbfields already have "." before them
				taintInfo.objField = edge.path + taintInfo.objField
			}
			edge.to.addTaintIfNotExists(taintInfo)

			doTaintBackwards(graph, edge.to.val, taintInfo, visited)
		}
	}
}

func fieldIndexToName(t *ssa.FieldAddr) string {
	return t.X.Type().Underlying().(*types.Pointer).Elem().(*types.Named).Underlying().(*types.Struct).Field(t.Field).Name()
}

type TaintInfo struct {
	val           ssa.Value
	objField      string
	dbField       string
	objPrefix     string
	dbfieldPrefix string
}

func doTaintBackwards(graph *Graph, val ssa.Value, taintInfo TaintInfo, visited map[TaintInfo]bool) {
	fmt.Printf("visiting value %s: %s // TAINT INFO = %v\n", val.Name(), val.String(), taintInfo)
	taintInfo.val = val
	if visited[taintInfo] {
		fmt.Printf("\tskipping value %s: %s\n", val.Name(), val.String())
		return
	}
	visited[taintInfo] = true

	node := graph.getNodeByName(val.Name())
	node.addTaintIfNotExists(taintInfo)

	switch t := val.(type) {
	case *ssa.MakeInterface:
		doTaintBackwards(graph, t.X, taintInfo, visited)
	case *ssa.UnOp:
		doTaintBackwards(graph, t.X, taintInfo, visited)
	case *ssa.Phi:
		// includes values in t.Edges + other nodes pointing to
		for _, edge := range graph.getEdgesFromNode(graph.getNodeByName(t.Name())) {
			// in case it points to an instruction like store we need to fetch the value
			// (in this case, this corresponds to the variable where something is being stored, and NOT the value being stored)
			if edge.to.instr != nil && edge.to.val == nil {
				edge.to.addTaintIfNotExists(taintInfo)
				for _, edge2 := range graph.getEdgesToNode(edge.to) {
					doTaintBackwards(graph, edge2.from.val, taintInfo, visited)
				}
			} /* else { // FIXME infinite loops
				doTaintBackwards(graph, edge.to.val, taintInfo, visited)
			} */
		}
	case *ssa.FieldAddr:
		fieldName := fieldIndexToName(t)
		// add after
		taintInfoTmp := taintInfo
		taintInfoTmp.objField = "." + fieldName + taintInfoTmp.objField
		doTaintBackwards(graph, t.X, taintInfoTmp, visited)
	case *ssa.IndexAddr:
		// add after
		taintInfoTmp := taintInfo
		taintInfoTmp.objField = "[*]" + taintInfoTmp.objField
		doTaintBackwards(graph, t.X, taintInfoTmp, visited)
	case *ssa.Parameter, *ssa.Alloc:
		doTaintNode(graph, val, taintInfo)
	default:
		fmt.Printf("[INFO] ignoring value: [%T] %v\n", val, val)
	}

	// if its fieldaddr then we use the objfield and dbfield
	// from the parameters and not the updated ones
	doTaintPointerToSets(graph, val, taintInfo, visited)
}

func doTaint(graph *Graph) {
	for _, node := range graph.nodes {
		call, _, ok := isMongoDBCall(node.instr)
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

			visited := make(map[TaintInfo]bool)
			taintInfo := TaintInfo{
				objField:      "",
				dbField:       "",
				objPrefix:     "_obj",
				dbfieldPrefix: database + "." + collection,
			}
			doTaintBackwards(graph, valDocument, taintInfo, visited)
		}
	}
}

func getInstructionID(instr ssa.Instruction) string {
	if !instr.Pos().IsValid() { // meaning there is no position
		return ""
	}
	return "instr_" + instrString(instr) + "_" + fmt.Sprintf("%d", instr.Pos())
}

func instrString(instr ssa.Instruction) string {
	switch instr.(type) {
	case *ssa.Store:
		return "store"
	case *ssa.Return:
		return "ret"
	}
	return ""
}

func getValueID(val ssa.Value) string {
	if !val.Pos().IsValid() { // meaning there is no position
		if c, ok := val.(*ssa.Const); ok {
			if c.IsNil() {
				return fmt.Sprintf("nil_%d", rand.Int())
			}
		}
		return "val_" + valString(val) + val.Name()
	}
	return "val_" + valString(val) + "_" + fmt.Sprintf("%d", val.Pos())
}

func valString(val ssa.Value) string {
	switch val.(type) {
	case *ssa.Call:
		return "call"
	case *ssa.Alloc:
		return "alloc"
	case *ssa.Slice:
		return "slice"
	case *ssa.FieldAddr:
		return "field"
	case *ssa.IndexAddr:
		return "index"
	case *ssa.UnOp:
		return "unary"
	case *ssa.MakeInterface:
		return "iface"
	case *ssa.Convert:
		return "conv"
	case *ssa.Parameter:
		return "param"
	case *ssa.Global:
		return "glob"
	case *ssa.Phi:
		return "phi"
	case *ssa.Extract:
		return "extr"
	case *ssa.BinOp:
		return "binary"
	}
	return ""
}

func parseInstr(graph *Graph, instr ssa.Instruction, instrIdx int, visited map[ssa.Value]bool) *Node {
	fmt.Printf("[A] %02d [%T] %v\n", instrIdx, instr, instr.String())

	id := getInstructionID(instr)
	if id == "" { // e.g., conditions or jumps (instructions and not values)
		log.Printf("skipping instruction with invalid id: %v\n", instr)
		return nil
	}

	if val, ok := instr.(ssa.Value); ok {
		return parseValue(graph, instr, instrIdx, val, visited)
	}
	node := RegisterNewNode(graph, instr, id)

	switch t := instr.(type) {
	case *ssa.Store:
		// 04 [store] *t1 = currency
		addrNode := parseValue(graph, instr, instrIdx, t.Addr, visited)
		valNode := parseValue(graph, instr, instrIdx, t.Val, visited)

		graph.createAndAddNewEdge(addrNode, node, EDGE_STORE, 0, "")
		graph.createAndAddNewEdge(valNode, node, EDGE_USAGE, 0, "")

		fmt.Printf("ADDING EDGE FOR ADDR NDOE AND VAL NODE: %v // %v \n", t.Addr, t.Val)
	case *ssa.Return:
		for _, res := range t.Results {
			resNode := parseValue(graph, instr, instrIdx, res, visited)
			graph.createAndAddNewEdge(resNode, node, EDGE_STORE, 0, "")
		}
	default:
		fmt.Printf("[1] ignoring... %02d [%T] %v\n", instrIdx, instr, instr.String())
	}

	return node
}

func parseValue(graph *Graph, instr ssa.Instruction, instrIdx int, val ssa.Value, visited map[ssa.Value]bool) *Node {
	fmt.Printf("[B] %02d [%T] %v\n", instrIdx, val, val.String())

	if visited[val] {
		return graph.defs[val.Name()]
	}
	visited[val] = true

	id := getValueID(val)
	if id == "" { // sanity check
		log.Fatalf("unexpected invalid id for value: %v\n", val)
		return nil
	}

	node, exists := graph.defs[val.Name()]
	if !exists {
		node = RegisterNewNodeValue(graph, instr, val, id)
	}

	switch t := val.(type) {
	case *ssa.Call:
		for _, arg := range t.Call.Args {
			for _, edges := range graph.getEdgesFromNode(node) {
				if edges.to.name == arg.Name() {
					fmt.Printf("[INFO] skipping arg edge for %s\n", t.Name())
					continue
				}
			}
			for _, edges := range graph.getEdgesToNode(node) {
				if edges.from.name == arg.Name() {
					fmt.Printf("[INFO] skipping arg edge for %s\n", t.Name())
					continue
				}
			}
			argNode := parseValue(graph, instr, instrIdx, arg, visited)
			graph.createAndAddNewEdge(argNode, node, EDGE_STORE, 0, "")
		}
	case *ssa.Alloc:
		// nothing to do
	case *ssa.Slice:
		// nothing to do
		targetNode := parseValue(graph, instr, instrIdx, t.X, visited)
		graph.createAndAddNewEdge(targetNode, node, EDGE_USAGE, 0, "")
	case *ssa.FieldAddr:
		// 00 [field] t27 = &t0.Items [#3]
		targetNode := parseValue(graph, instr, instrIdx, t.X, visited)
		graph.createAndAddNewEdge(targetNode, node, EDGE_FIELD, 0, fieldIndexToName(t))
	case *ssa.IndexAddr:
		targetNode := parseValue(graph, instr, instrIdx, t.X, visited)
		//FIXME: should parse value for t.Index
		graph.createAndAddNewEdge(targetNode, node, EDGE_FIELD, 0, t.Index.String())
	case *ssa.UnOp:
		// 01 [unary] t14 = *t13
		// 05 [unary] t31 = *t30
		targetNode := parseValue(graph, instr, instrIdx, t.X, visited)
		graph.createAndAddNewEdge(targetNode, node, EDGE_LOAD, 0, "")

	case *ssa.MakeInterface: // same as *ssa.UnOp
		targetNode := parseValue(graph, instr, instrIdx, t.X, visited)
		graph.createAndAddNewEdge(targetNode, node, EDGE_USAGE, 0, "")
	case *ssa.Convert:
		targetNode := parseValue(graph, instr, instrIdx, t.X, visited)
		graph.createAndAddNewEdge(targetNode, node, EDGE_USAGE, 0, "")

	case *ssa.Parameter:
		// nothing to do

	case *ssa.Global:
		// nothing to do

	case *ssa.Phi:
		for _, phiEdge := range t.Edges {
			for _, edges := range graph.getEdgesFromNode(node) {
				if edges.to.name == phiEdge.Name() {
					fmt.Printf("[INFO] skipping phi edge for %s\n", t.Name())
					continue
				}
			}
			for _, edges := range graph.getEdgesToNode(node) {
				if edges.from.name == phiEdge.Name() {
					fmt.Printf("[INFO] skipping phi edge for %s\n", t.Name())
					continue
				}
			}
			edgeNode := parseValue(graph, instr, instrIdx, phiEdge, visited)
			graph.createAndAddNewEdge(edgeNode, node, EDGE_STORE, 0, "")
		}

	case *ssa.Extract:
		extractFromNode := parseValue(graph, instr, instrIdx, t.Tuple, visited)
		graph.createAndAddNewEdge(extractFromNode, node, EDGE_USAGE, t.Index, "")

	case *ssa.BinOp:
		xNode := parseValue(graph, instr, instrIdx, t.X, visited)
		yNode := parseValue(graph, instr, instrIdx, t.Y, visited)
		graph.createAndAddNewEdge(xNode, node, EDGE_STORE, 0, "")
		graph.createAndAddNewEdge(yNode, node, EDGE_USAGE, 0, "")

	default:
		fmt.Printf("[2] ignoring... %02d [%T] %s = %v\n", id, val, val.Name(), val.String())
	}
	return node
}

func (g *Graph) writeToDOTFile(fn string) error {
	filename := fmt.Sprintf("output/%s/graphs/%s.dot", appname, fn)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "digraph G {")
	fmt.Fprintln(file, "\trankdir=TD;")

	for _, node := range g.nodes {
		str := node.String()
		if node.isTainted() {
			str += "\n\n==== tainted ====\n" + node.taintString()
		}
		label := strings.ReplaceAll(str, `"`, `\"`)
		nodecolor := node.colorForSSA()

		shape := "ellipse"
		if node.isTainted() {
			shape = "box"
		}

		color := "black"
		if nodecolor != "" {
			color = nodecolor
		}

		fmt.Fprintf(file, "\tN_%s [label=\"%s\", style=bold, shape=%s, color=\"%s\"];\n", node.id, label, shape, color)
	}

	for _, edge := range g.edges {
		if edge.edgeType == EDGE_POINTS_TO {
			path := strings.ReplaceAll(edge.path, `"`, `\"`)
			fmt.Fprintf(file, "\tN_%s -> N_%s [label=\"%s\", style=dashed, color=blue];\n", edge.from.id, edge.to.id, path)
		} else if edge.from != nil && edge.to != nil {
			fmt.Fprintf(file, "\tN_%s -> N_%s;\n", edge.from.id, edge.to.id)
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
				fmt.Printf("[%s] [%s] [%T] \t %s %v\n", fn, node.id, node.instr, prefix, node.instr.String())
			} else {
				fmt.Printf("[%s] [%s] [%T] \t %s %v\n", fn, node.id, node.val, prefix, node.val.String())
			}
		}
	}

	fmt.Print("\n\n ========== TAINTS ========== \n\n")
	for fn, graph := range funcGraphs {
		for _, node := range graph.nodes {
			if node.isTainted() {
				for obj, dbfields := range node.taints {
					fmt.Printf("[%s] %s [%s]: %s\n", fn, node.String(), node.name, obj)
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
	outfile1, err := os.Create(fmt.Sprintf("output/%s/app.ssa", appname))
	if err != nil {
		log.Fatal(err)
	}
	defer outfile1.Close()

	outfile2, err := os.Create(fmt.Sprintf("output/%s/%s.out", appname, pkg.Pkg.Name()))
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outfile2.Close()
	pkg.WriteTo(outfile2)

	for _, member := range pkg.Members {
		switch m := member.(type) {
		case *ssa.Function:
			iterateFunc(outfile1, m, nil, funcGraphs)

		case *ssa.Global:
			fmt.Fprintf(outfile1, "\tGlobal: %s, Type: %s\n", m.Name(), m.Type().String())

		case *ssa.Type:
			fmt.Fprintf(outfile1, "\tType: %s\n", m.Type())

			// this logic was copied from
			// package: golang.org/x/tools/go/ssa
			// file: print.go
			// function: func (p *Package) WriteTo(w io.Writer) (int64, error)
			for _, sel := range typeutil.IntuitiveMethodSet(m.Type(), &prog.MethodSets) {
				method := prog.MethodValue(sel)
				fmt.Fprintf(outfile1, "\tMethod: %v\n", sel.Obj().Type())
				if method != nil {
					iterateFunc(outfile1, method, m.Type(), funcGraphs)
				}
			}

			methods := prog.MethodSets.MethodSet(m.Type().Underlying())
			for i := 0; i < methods.Len(); i++ {
				sel := methods.At(i)
				fmt.Fprintf(outfile1, "\tMethod: %v\n", sel.Obj().Type())
				method := prog.MethodValue(sel)
				if method != nil {
					iterateFunc(outfile1, method, m.Type(), funcGraphs)
				}
			}

		default:
			fmt.Fprintf(outfile1, "\tUnknown member type: %T\n", m)
		}
	}
}

func iterateFunc(outFile *os.File, fn *ssa.Function, memberType types.Type, funcGraphs map[string]*Graph) {
	fullfuncname := fn.Package().Pkg.Name() + "." + fn.Name()
	graph := &Graph{
		defs: make(map[string]*Node),
	}
	if _, exists := funcGraphs[fn.Name()]; exists {
		log.Fatalf("graph for function (%s) already exists\n", fullfuncname)
	}
	funcGraphs[fn.Name()] = graph
	fmt.Printf("added new graph for function (%s)\n", fullfuncname)

	var visited = make(map[ssa.Value]bool)

	fmt.Fprintf(outFile, "\t\tParameters:\n")
	for i, param := range fn.Params {
		fmt.Fprintf(outFile, "\t\t\t%s = %s\n", param.Name(), param.String())
		parseValue(graph, nil, -i-1, param, visited)
	}

	fmt.Fprintf(outFile, "Function: %s\n", fullfuncname)
	for i, block := range fn.Blocks {
		fmt.Fprintf(outFile, "Block #%d: %s.%s\n", i, fullfuncname, block.Comment)
		for j, instr := range block.Instrs {
			parseInstr(graph, instr, j, visited)

			if val, ok := instr.(ssa.Value); ok {
				fmt.Fprintf(outFile, "\t\t\t%02d: %s = %s\n", j, val.Name(), instr.String())
			} else {
				fmt.Fprintf(outFile, "\t\t\t%02d: %s\n", j, instr.String())
			}
		}
	}
}

func initPackages() (*ssa.Program, *ssa.Package, *pointer.Result, error) {
	// e.g. "../apps/test2/main.go"
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
	outFile, err := os.Create(fmt.Sprintf("output/%s/app.ptrs", appname))
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
						edge, _ := graph.createAndAddNewEdge(node, pointsToNode, EDGE_POINTS_TO, 0, "")
						if edge != nil {
							edge.path = lbl.Path()
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
