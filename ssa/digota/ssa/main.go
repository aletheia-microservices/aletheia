package main

import (
	"fmt"
	"go/token"
	"go/types"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/types/typeutil"
)

// -------------------------------------
// ------------- CONSTANTS -------------
// -------------------------------------
const (
	inputpackagepath = "../workflow/digota/"
	outfilename      = "./ssa-simple.out"
)

// -------------------------------------

// -------------------------------------
// --------------- GRAPH ---------------
// -------------------------------------

type Graph struct {
	nodes []*Node
	edges []*Edge
	pos   map[string]*Node
}

func NewAbstractGraph() *Graph {
	return &Graph{pos: make(map[string]*Node)}
}

func (graph *Graph) String() string {
	/* str := "\nGRAPH:\n\n"
	str += "\n\tNODES:\n" */

	keys := make([]string, 0, len(graph.pos))
	for k := range graph.pos {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		// remove the "t" prefix and convert to int
		ni, _ := strconv.Atoi(strings.TrimPrefix(keys[i], "t"))
		nj, _ := strconv.Atoi(strings.TrimPrefix(keys[j], "t"))
		return ni < nj
	})

	/* for _, pos := range keys {
		node := graph.pos[pos]
		str += fmt.Sprintf("%s: %s\n", pos, node.str)
	} */

	//str += "\n\tNODES W/ EDGES:\n"

	str := ""
	for _, pos := range keys {
		node := graph.pos[pos]
		if node.tainted {
			str += "*"
		}
		str += fmt.Sprintf("%s (%s): %s\n", node.name, strings.Join(node.pointedBy, ","), node.str)
		for _, edge := range graph.GetEdges() {
			if edge.HasFromNode(node) {
				str += fmt.Sprintf("\t|--%s--> %s (%s): %s\n", edge.Type(), edge.to.name, strings.Join(edge.to.pointedBy, ","), edge.to.str)
			}
		}
		str += "\n"
	}

	return str
}

func (graph *Graph) GetNodeByPosIfExists(pos string) *Node {
	return graph.pos[pos]
}

func (graph *Graph) GetNodeByPos(pos string) *Node {
	node, ok := graph.pos[pos]
	if !ok {
		log.Fatalf("could not find node for pos (%s)", pos)
	}
	return node
}

func (graph *Graph) AddExistingNodeToNewPos(node *Node, pos string) {
	graph.pos[pos] = node
}

func (graph *Graph) AddNode(node *Node, pos string) {
	graph.nodes = append(graph.nodes, node)
	graph.pos[pos] = node
}

func (graph *Graph) AddEdge(edge *Edge) {
	graph.edges = append(graph.edges, edge)
}

func (graph *Graph) EdgeExists(from *Node, to *Node, fieldOrIndex bool, fieldVal int, indexVal string) bool {
	for _, edge := range graph.edges {
		if edge.from.String() == from.String() && edge.to.String() == to.String() && (edge.isField == fieldOrIndex || edge.isIndex == fieldOrIndex) && (edge.field == fieldVal || edge.index == indexVal) {
			return true
		}
	}
	return false
}

func (graph *Graph) GetNodes() []*Node {
	return graph.nodes
}

func (graph *Graph) GetEdges() []*Edge {
	return graph.edges
}

func (graph *Graph) GetEdgesFromNode(node *Node) []*Edge {
	var edges []*Edge
	for _, edge := range graph.edges {
		if edge.from == node {
			edges = append(edges, edge)
		}
	}
	return edges
}

func (graph *Graph) GetValueEdgesFromNode(node *Node) []*Edge {
	var edges []*Edge
	for _, edge := range graph.edges {
		if edge.from == node && edge.isValue {
			edges = append(edges, edge)
		}
	}
	return edges
}

func (graph *Graph) GetFieldEdgesFromNode(node *Node) []*Edge {
	var edges []*Edge
	for _, edge := range graph.edges {
		if edge.from == node && edge.isField {
			edges = append(edges, edge)
		}
	}
	return edges
}

func (graph *Graph) NodeHasValueEdge(node *Node) bool {
	for _, edge := range graph.GetEdgesFromNode(node) {
		if edge.isValue {
			return true
		}
	}
	return false
}

func (graph *Graph) GetNodeForSSAValue(v ssa.Value) *Node {
	for _, node := range graph.nodes {
		if (node.name == v.Name() || slices.Contains(node.pointedBy, v.Name())) && node.str == v.String() {
			return node
		}
	}
	return nil
}

func (graph *Graph) GetNodeForSSAValue2(v ssa.Value) *Node {
	for _, node := range graph.nodes {
		if node.str == v.String() {
			return node
		}
	}
	return nil
}

// -------------------------------------
// -------------- DATAFLOW -------------
// -------------------------------------

// -------------------------------------
// ---------------- NODE ---------------
// -------------------------------------

type Node struct { // objects
	name      string
	str       string
	fn        string
	pointedBy []string //variables
	edgesTo   []*Edge

	isFunction    bool
	isPlaceholder bool
	isParameter   bool
	isPhi         bool

	tainted bool
}

func (node *Node) String() string {
	return fmt.Sprintf("(%s) %s: %s", node.fn, node.name, node.str)
}

func (node *Node) NewNodeVersion(name string) *Node {
	return &Node{
		name:        name,
		pointedBy:   []string{name},
		str:         node.str,
		fn:          node.fn,
		isParameter: node.isParameter, // is this possible?
	}
}

func NewNode(val ssa.Value) *Node {
	return &Node{
		name:      val.Name(),
		str:       val.String(),
		fn:        val.Parent().Name(),
		pointedBy: []string{val.Name()},
	}
}

func NewNodePhi(val ssa.Value) *Node {
	return &Node{
		name:      val.Name(),
		str:       val.String(),
		fn:        val.Parent().Name(),
		pointedBy: []string{val.Name()},
		isPhi:     true,
	}
}

func NewNodeFunction(val ssa.Value) *Node {
	return &Node{
		name:       val.Name(),
		str:        val.String(),
		fn:         val.Parent().Name(),
		pointedBy:  []string{val.Name()},
		isFunction: true,
	}
}

func NewNodePlaceholder(val ssa.Value) *Node {
	return &Node{
		name:          val.Name(),
		str:           val.String(),
		fn:            val.Parent().Name(),
		pointedBy:     []string{val.Name()},
		isPlaceholder: true,
	}
}

func NewNodeParameter(val ssa.Value) *Node {
	var fn string
	if val.Parent() != nil {
		fn = val.Parent().Name()
	}
	return &Node{
		name:        val.Name(),
		str:         val.String(),
		fn:          fn,
		pointedBy:   []string{val.Name()},
		isParameter: true,
	}
}

func (node *Node) AddToPointedBy(pointedBy string) {
	node.pointedBy = append(node.pointedBy, pointedBy)
}

func (node *Node) AddToEdgesTo(edge *Edge) {
	node.edgesTo = append(node.edgesTo, edge)
}

func (node *Node) HasName(name string) bool {
	return node.name == name
}

func (node *Node) HasFn(fn string) bool {
	return node.fn == fn
}

// -------------------------------------
// ---------------- EDGE ---------------
// -------------------------------------

type Edge struct {
	from *Node
	to   *Node

	isField     bool
	isIndex     bool
	isParameter bool
	isValue     bool
	isVersion   bool
	isInterface bool
	isConverted bool
	isPhi       bool
	isReturn    bool

	field     int
	index     string
	parameter int
}

func (edge *Edge) Type() string {
	if edge.isField {
		return fmt.Sprintf("field(%d)", edge.field)
	} else if edge.isIndex {
		return fmt.Sprintf("index(%s)", edge.index)
	} else if edge.isParameter {
		return fmt.Sprintf("parameter(%d)", edge.parameter)
	} else if edge.isValue {
		return fmt.Sprint("value")
	} else if edge.isVersion {
		return fmt.Sprint("version")
	} else if edge.isInterface {
		return fmt.Sprint("interface of")
	} else if edge.isInterface {
		return fmt.Sprint("converted to")
	} else if edge.isInterface {
		return fmt.Sprint("phi")
	} else if edge.isReturn {
		return fmt.Sprintf("returns(%d)", edge.parameter)
	}
	return ""
}

func NewEdgeParameter(from *Node, to *Node, parameter int) *Edge {
	return &Edge{
		from:        from,
		to:          to,
		isParameter: true,
		parameter:   parameter,
	}
}

func NewEdgeField(from *Node, to *Node, field int) *Edge {
	return &Edge{
		from:    from,
		to:      to,
		isField: true,
		field:   field,
	}
}

func NewEdgeIndex(from *Node, to *Node, index string) *Edge {
	return &Edge{
		from:    from,
		to:      to,
		isIndex: true,
		index:   index,
	}
}

func NewEdgeValue(from *Node, to *Node) *Edge {
	return &Edge{
		from:    from,
		to:      to,
		isValue: true,
	}
}

func NewEdgeVersion(from *Node, to *Node) *Edge {
	return &Edge{
		from:      from,
		to:        to,
		isVersion: true,
	}
}

func NewEdgePhi(from *Node, to *Node) *Edge {
	return &Edge{
		from:  from,
		to:    to,
		isPhi: true,
	}
}

func NewEdgeReturn(from *Node, to *Node, idx int) *Edge {
	return &Edge{
		from:      from,
		to:        to,
		parameter: idx,
		isReturn:  true,
	}
}

func NewEdgeInterface(from *Node, to *Node) *Edge {
	return &Edge{
		from:        from,
		to:          to,
		isInterface: true,
	}
}

func NewEdgeConvert(from *Node, to *Node) *Edge {
	return &Edge{
		from:        from,
		to:          to,
		isConverted: true,
	}
}

func (edge *Edge) HasFromNode(node *Node) bool {
	return edge.from == node
}

// -------------------------------------
// -------------- PARSER ---------------
// -------------------------------------

func ParseInstr(graph *Graph, instr ssa.Instruction, idx int) {
	if val, ok := instr.(ssa.Value); ok {
		parseValue(graph, instr, idx, val)
		return
	}

	switch t := instr.(type) {
	case *ssa.Store:
		fmt.Printf("%02d [store] %v\n", idx, instr.String())
		addrNode := parseValue(graph, instr, idx, t.Addr)
		valNode := parseValue(graph, instr, idx, t.Val)

		/* fmt.Printf("HERE!!! %v \n", graph.String())	 */
		if graph.NodeHasValueEdge(addrNode) {
			newAddrNode := addrNode.NewNodeVersion(t.Addr.Name())
			graph.AddNode(newAddrNode, t.Addr.Name())
			edge := NewEdgeVersion(addrNode, newAddrNode)
			graph.AddEdge(edge)
			addrNode = newAddrNode
		}

		edge := NewEdgeValue(addrNode, valNode)
		graph.AddEdge(edge)
	case *ssa.Return:
		fmt.Printf("[A] skipping... %02d [%T] %v\n", idx, instr, instr.String())
	case *ssa.Jump:
		fmt.Printf("[A] skipping... %02d [%T] %v\n", idx, instr, instr.String())
	default:
		fmt.Printf("[1] ignoring... %02d [%T] %v\n", idx, instr, instr.String())
	}
}

func recurseTaint(graph *Graph, node *Node) {
	fmt.Printf("visiting node: %v\n", node.String())
	if node.tainted == false {
		node.tainted = true
		fmt.Printf("\ttainting node: %v\n", node.String())
		for _, edge := range graph.GetEdgesFromNode(node) {
			recurseTaint(graph, edge.to)
		}
	}
}


func parseValue(graph *Graph, instr ssa.Instruction, idx int, val ssa.Value) *Node {
	if node := graph.GetNodeByPosIfExists(val.Name()); node != nil {
		return node
	}

	switch t := val.(type) {
	case *ssa.Call:
		fmt.Printf("%02d [call] %s = %v\n", idx, val.Name(), val.String())
		node := NewNodeFunction(val)

		graph.AddNode(node, val.Name())
		var argNodes []*Node
		for i, arg := range t.Call.Args {
			argNode := parseValue(graph, instr, idx, arg)
			argNodes = append(argNodes, argNode)
			edge := NewEdgeParameter(node, argNode, i)
			graph.AddEdge(edge)
		}

		if fn, ok := t.Call.Value.(*ssa.Builtin); ok {
			fmt.Printf("1. CALLING: %v\n", fn.Name())
			for _, argNode := range argNodes {
				fmt.Printf("\t ARGNODE = %v\n", argNode)
			}
		} else if t.Call.Method != nil &&
			t.Call.Method.Signature().Recv().Type().String() == "github.com/blueprint-uservices/blueprint/runtime/core/backend.NoSQLCollection" &&
			t.Call.Method.Name() == "InsertOne" {
			fmt.Printf("2. CALLING: [%T] %v //[%T] %v\n", t.Call.Method, t.Call.Method, t.Call.Value, t.Call.Value)
			for _, argNode := range argNodes {
				recurseTaint(graph, argNode)
			}
		}

		return node
	case *ssa.Alloc:
		fmt.Printf("%02d [alloc] %s = %v\n", idx, val.Name(), val.String())
		node := NewNode(t)
		graph.AddNode(node, t.Name())
		return node
	case *ssa.Slice:
		fmt.Printf("%02d [slice] %s = %v\n", idx, val.Name(), val.String())
		node := parseValue(graph, instr, idx, t.X)
		node.AddToPointedBy(val.String())
		return node
	case *ssa.FieldAddr:
		// 00 [field] t27 = &t0.Items [#3]
		fmt.Printf("%02d [field] %s = %v\n", idx, val.Name(), val.String())
		node := graph.GetNodeForSSAValue2(val)
		if node == nil {
			node = NewNode(t)
			graph.AddNode(node, val.Name())
		} else {
			node.AddToPointedBy(val.Name())
			graph.AddExistingNodeToNewPos(node, val.Name())
		}
		topNode := graph.GetNodeByPosIfExists(t.X.Name())
		if topNode == nil {
			// e.g. 00 [field] t36 = &s.skuService [#0]
			// ignore for now
			return nil
		}
		if !graph.EdgeExists(topNode, node, true, t.Field, "") {
			edge := NewEdgeField(topNode, node, t.Field)
			graph.AddEdge(edge)
			topNode.AddToEdgesTo(edge)
		}
		return node
	case *ssa.IndexAddr:
		fmt.Printf("%02d [index] %s = %v\n", idx, val.Name(), val.String())
		node := graph.GetNodeForSSAValue2(val)
		if node == nil {
			node = NewNode(t)
			graph.AddNode(node, t.Name())
		} else {
			node.AddToPointedBy(t.Name())
			graph.AddExistingNodeToNewPos(node, t.Name())
		}
		topNode := graph.GetNodeByPosIfExists(t.X.Name())
		if topNode == nil {
			// ignore for now
			return nil
		}
		if !graph.EdgeExists(topNode, node, true, 0, t.Index.String()) {
			edge := NewEdgeIndex(topNode, node, t.Index.String())
			graph.AddEdge(edge)
			topNode.AddToEdgesTo(edge)
		}
		return node
	case *ssa.UnOp:
		// e.g.,
		// 01 [unary] t14 = *t13
		// 05 [unary] t31 = *t30
		fmt.Printf("%02d [unary] %s = %v\n", idx, val.Name(), val.String())
		xNode := parseValue(graph, instr, idx, t.X)

		node := NewNode(t)
		graph.AddNode(node, t.Name())

		valueEdges := graph.GetValueEdgesFromNode(xNode)
		fieldEdges := graph.GetFieldEdgesFromNode(xNode)
		if valueEdges != nil {
			// 1. get the current value of the current address
			// 2. create new address and assign that value
			xNodeValue := valueEdges[0].to
			edge := NewEdgeValue(node, xNodeValue)
			graph.AddEdge(edge)
		} else if fieldEdges != nil { //FIXME
			edge := NewEdgeValue(node, xNode)
			graph.AddEdge(edge)
		} else {
			// 1. assign value for the first time if it does not exist yet
			edge := NewEdgeValue(xNode, node)
			graph.AddEdge(edge)
		}
		return node
	case *ssa.MakeInterface: // same as *ssa.UnOp
		fmt.Printf("%02d [interface] %s = %v\n", idx, val.Name(), val.String())
		xNode := parseValue(graph, instr, idx, t.X)

		node := NewNode(t)
		graph.AddNode(node, t.Name())

		edge := NewEdgeInterface(node, xNode)
		graph.AddEdge(edge)
		return node
	case *ssa.Convert: // same as *ssa.UnOp and *ssa.MakeInterface
		fmt.Printf("%02d [convert] %s = %v\n", idx, val.Name(), val.String())
		xNode := parseValue(graph, instr, idx, t.X)

		node := NewNode(t)
		graph.AddNode(node, t.Name())

		edge := NewEdgeConvert(node, xNode)
		graph.AddEdge(edge)
		return node
	case *ssa.Parameter: // dynamic
		fmt.Printf("%02d [parameter] %s = %v\n", idx, val.Name(), val.String())
		node := NewNodePlaceholder(val)
		graph.AddNode(node, val.Name())
		return node
	case *ssa.Const:
		fmt.Printf("%02d [const] %s = %v\n", idx, val.Name(), val.String())
		return NewNodeParameter(val)
	case *ssa.Phi:
		fmt.Printf("%02d [phi] %s = %v\n", idx, val.Name(), val.String())
		node := NewNodePhi(val)
		graph.AddNode(node, val.Name())
		for _, phiEdge := range t.Edges {
			fmt.Printf("HERE FOR PHI EDGE: %v\n", phiEdge)
			otherNode := parseValue(graph, instr, idx, phiEdge)
			edge := NewEdgePhi(node, otherNode)
			graph.AddEdge(edge)
		}
		return node
	case *ssa.Extract:
		fmt.Printf("%02d [extract] %s = %v\n", idx, val.Name(), val.String())
		extractFromNode := parseValue(graph, instr, idx, t.Tuple)
		node := NewNode(t)
		graph.AddNode(node, val.Name())
		edge := NewEdgeReturn(extractFromNode, node, t.Index)
		graph.AddEdge(edge)
		return node

	case *ssa.BinOp, *ssa.Global: //FIXME
		fmt.Printf("[B] skipping... %02d [%T] %s = %v\n", idx, val, val.Name(), val.String())
		node := NewNode(val)
		graph.AddNode(node, val.Name())
		return node
	default:
		fmt.Printf("[2] ignoring... %02d [%T] %s = %v\n", idx, val, val.Name(), val.String())
	}
	log.Fatal("returning nil node")
	return nil
}

// -------------------------------------
// ---------------- MAIN ---------------
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
		if ssaPkg.Pkg.Name() != "digota" {
			continue
		}
		/* if ssaPkg.Func("main") == nil && ssaPkg.Func("init") == nil {
			continue
		} */
		appPkgs = append(appPkgs, ssaPkg)
	}

	ssaAnalysis(prog, appPkgs)
}

func iterateFunc(outFile *os.File, fn *ssa.Function, memberType types.Type) {
	var graph = NewAbstractGraph()

	var filename string
	namedMemberType, ok := memberType.(*types.Named)

	if ok && namedMemberType.Obj().Name() != "SkuServiceImpl" && namedMemberType.Obj().Name() != "OrderServiceImpl" {
		return
	}

	fmt.Printf("=============================\n")
	if ok {
		filename = fmt.Sprintf("%s/%s.graph", namedMemberType.Obj().Name(), fn.Name())
		fmt.Printf("%s.%s()\n", namedMemberType.Obj().Name(), fn.Name())
	} else {
		filename = fmt.Sprintf("%s.graph", fn.Name())
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
			if filename == "OrderServiceImpl/New.graph" {
				ParseInstr(graph, instr, j)
			}
		}
	}

	if ok {
		outfile, err := os.Create(fmt.Sprintf("out/%s", filename))
		if err != nil {
			log.Fatalf("failed to create output file: %v", err)
		}
		outfile.WriteString(graph.String())
		defer outfile.Close()

	}
	fmt.Println()
	fmt.Println()
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
