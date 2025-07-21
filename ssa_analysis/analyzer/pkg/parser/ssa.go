package parser

import (
	"crypto/rand"
	"fmt"
	"go/types"
	"log"
	"math/big"
	"os"

	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/types/typeutil"

	"analyzer/pkg/graph"
	"analyzer/pkg/utils"
)

func RunSSAAnalysis(appname string, prog *ssa.Program, pkg *ssa.Package, funcGraphs map[string]*graph.Graph) {
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

func iterateFunc(outFile *os.File, fn *ssa.Function, memberType types.Type, funcGraphs map[string]*graph.Graph) {
	fullfuncname := fn.Package().Pkg.Name() + "." + fn.Name()
	g := graph.NewGraph()
	if _, exists := funcGraphs[fn.Name()]; exists {
		log.Fatalf("graph for function (%s) already exists\n", fullfuncname)
	}
	funcGraphs[fn.Name()] = g
	fmt.Printf("added new graph for function (%s)\n", fullfuncname)

	var visited = make(map[ssa.Value]bool)

	fmt.Fprintf(outFile, "\t\tParameters:\n")
	for i, param := range fn.Params {
		fmt.Fprintf(outFile, "\t\t\t%s = %s\n", param.Name(), param.String())
		parseValue(g, nil, -i-1, param, visited)
	}

	fmt.Fprintf(outFile, "Function: %s\n", fullfuncname)
	for i, block := range fn.Blocks {
		fmt.Fprintf(outFile, "Block #%d: %s.%s\n", i, fullfuncname, block.Comment)
		for j, instr := range block.Instrs {
			parseInstr(g, instr, j, visited)

			if val, ok := instr.(ssa.Value); ok {
				fmt.Fprintf(outFile, "\t\t\t%02d: %s = %s\n", j, val.Name(), instr.String())
			} else {
				fmt.Fprintf(outFile, "\t\t\t%02d: %s\n", j, instr.String())
			}
		}
	}
}

func parseInstr(g *graph.Graph, instr ssa.Instruction, instrIdx int, visited map[ssa.Value]bool) *graph.Node {
	fmt.Printf("[A] %02d [%T] %v\n", instrIdx, instr, instr.String())

	id := getInstructionID(instr)
	if id == "" { // e.g., conditions or jumps (instructions and not values)
		log.Printf("skipping instruction with invalid id: %v\n", instr)
		return nil
	}

	if val, ok := instr.(ssa.Value); ok {
		return parseValue(g, instr, instrIdx, val, visited)
	}
	node := graph.RegisterNewNode(g, instr, id)

	switch t := instr.(type) {
	case *ssa.Store:
		// 04 [store] *t1 = currency
		addrNode := parseValue(g, instr, instrIdx, t.Addr, visited)
		valNode := parseValue(g, instr, instrIdx, t.Val, visited)

		g.CreateAndAddNewEdge(addrNode, node, graph.EDGE_STORE, 0, "")
		g.CreateAndAddNewEdge(valNode, node, graph.EDGE_USAGE, 0, "")

		fmt.Printf("ADDING EDGE FOR ADDR NDOE AND VAL NODE: %v // %v \n", t.Addr, t.Val)
	case *ssa.Return:
		for _, res := range t.Results {
			resNode := parseValue(g, instr, instrIdx, res, visited)
			g.CreateAndAddNewEdge(resNode, node, graph.EDGE_STORE, 0, "")
		}
	default:
		fmt.Printf("[1] ignoring... %02d [%T] %v\n", instrIdx, instr, instr.String())
	}

	return node
}

func parseValue(g *graph.Graph, instr ssa.Instruction, instrIdx int, val ssa.Value, visited map[ssa.Value]bool) *graph.Node {
	fmt.Printf("[B] %02d [%T] %v\n", instrIdx, val, val.String())

	if visited[val] {
		return g.GetNodeByName(val.Name())
	}
	visited[val] = true

	id := getValueID(val)
	if id == "" { // sanity check
		log.Fatalf("unexpected invalid id for value: %v\n", val)
		return nil
	}

	node, exists := g.GetNodeByIfExists(val.Name())
	if !exists {
		node = graph.RegisterNewNodeValue(g, instr, val, id)
	}

	switch t := val.(type) {
	case *ssa.Call:
		for _, arg := range t.Call.Args {
			for _, edges := range g.GetEdgesFromNode(node) {
				if edges.GetToNode().GetName() == arg.Name() {
					fmt.Printf("[INFO] skipping arg edge for %s\n", t.Name())
					continue
				}
			}
			for _, edges := range g.GetEdgesToNode(node) {
				if edges.GetFromNode().GetName() == arg.Name() {
					fmt.Printf("[INFO] skipping arg edge for %s\n", t.Name())
					continue
				}
			}
			argNode := parseValue(g, instr, instrIdx, arg, visited)
			g.CreateAndAddNewEdge(argNode, node, graph.EDGE_STORE, 0, "")
		}
	case *ssa.Alloc:
		// nothing to do
	case *ssa.Slice:
		// nothing to do
		targetNode := parseValue(g, instr, instrIdx, t.X, visited)
		g.CreateAndAddNewEdge(targetNode, node, graph.EDGE_USAGE, 0, "")
	case *ssa.FieldAddr:
		// 00 [field] t27 = &t0.Items [#3]
		targetNode := parseValue(g, instr, instrIdx, t.X, visited)
		g.CreateAndAddNewEdge(targetNode, node, graph.EDGE_FIELD, 0, utils.FieldIndexToName(t))
	case *ssa.IndexAddr:
		targetNode := parseValue(g, instr, instrIdx, t.X, visited)
		//FIXME: should parse value for t.Index
		g.CreateAndAddNewEdge(targetNode, node, graph.EDGE_FIELD, 0, t.Index.String())
	case *ssa.UnOp:
		// 01 [unary] t14 = *t13
		// 05 [unary] t31 = *t30
		targetNode := parseValue(g, instr, instrIdx, t.X, visited)
		g.CreateAndAddNewEdge(targetNode, node, graph.EDGE_LOAD, 0, "")

	case *ssa.MakeInterface: // same as *ssa.UnOp
		targetNode := parseValue(g, instr, instrIdx, t.X, visited)
		g.CreateAndAddNewEdge(targetNode, node, graph.EDGE_USAGE, 0, "")
	case *ssa.Convert:
		targetNode := parseValue(g, instr, instrIdx, t.X, visited)
		g.CreateAndAddNewEdge(targetNode, node, graph.EDGE_USAGE, 0, "")

	case *ssa.Parameter:
		// nothing to do

	case *ssa.Global:
		// nothing to do

	case *ssa.Phi:
		for _, phiEdge := range t.Edges {
			for _, edges := range g.GetEdgesFromNode(node) {
				if edges.GetToNode().GetName() == phiEdge.Name() {
					fmt.Printf("[INFO] skipping phi edge for %s\n", t.Name())
					continue
				}
			}
			for _, edges := range g.GetEdgesToNode(node) {
				if edges.GetFromNode().GetName() == phiEdge.Name() {
					fmt.Printf("[INFO] skipping phi edge for %s\n", t.Name())
					continue
				}
			}
			edgeNode := parseValue(g, instr, instrIdx, phiEdge, visited)
			g.CreateAndAddNewEdge(edgeNode, node, graph.EDGE_STORE, 0, "")
		}

	case *ssa.Extract:
		extractFromNode := parseValue(g, instr, instrIdx, t.Tuple, visited)
		g.CreateAndAddNewEdge(extractFromNode, node, graph.EDGE_USAGE, t.Index, "")

	case *ssa.BinOp:
		xNode := parseValue(g, instr, instrIdx, t.X, visited)
		yNode := parseValue(g, instr, instrIdx, t.Y, visited)
		g.CreateAndAddNewEdge(xNode, node, graph.EDGE_STORE, 0, "")
		g.CreateAndAddNewEdge(yNode, node, graph.EDGE_USAGE, 0, "")

	default:
		fmt.Printf("[2] ignoring... %s [%T] %s = %v\n", id, val, val.Name(), val.String())
	}
	return node
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
				n, err := rand.Int(rand.Reader, big.NewInt(1<<31))
				if err != nil {
					return "nil_rand_error"
				}
				return fmt.Sprintf("nil_%d", n)
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
