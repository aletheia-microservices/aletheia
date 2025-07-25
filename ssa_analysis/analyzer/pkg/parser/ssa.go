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

	"analyzer/pkg/ssa_graph"
	"analyzer/pkg/utils"
)

func RunSSAAnalysis(appname string, prog *ssa.Program, pkg *ssa.Package, funcGraphs map[string]*ssa_graph.SSAGraph) {
	outfile1, err := os.Create(fmt.Sprintf("output/%s/%s.out", appname, pkg.Pkg.Name()))
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outfile1.Close()
	pkg.WriteTo(outfile1)

	outfile2, err := os.Create(fmt.Sprintf("output/%s/%s.ssa", appname, pkg.Pkg.Name()))
	if err != nil {
		log.Fatal(err)
	}
	defer outfile2.Close()

	for _, member := range pkg.Members {
		switch m := member.(type) {
		case *ssa.Function:
			iterateFunc(outfile2, m, nil, funcGraphs)

		case *ssa.Global:
			fmt.Fprintf(outfile2, "\tGlobal: %s, Type: %s\n", m.Name(), m.Type().String())

		case *ssa.Type:
			fmt.Fprintf(outfile2, "\tType: %s\n", m.Type())

			// this logic was copied from
			// package: golang.org/x/tools/go/ssa
			// file: print.go
			// function: func (p *Package) WriteTo(w io.Writer) (int64, error)
			for _, sel := range typeutil.IntuitiveMethodSet(m.Type(), &prog.MethodSets) {
				method := prog.MethodValue(sel)
				fmt.Fprintf(outfile2, "\tMethod: %v\n", sel.Obj().Type())
				if method != nil {
					iterateFunc(outfile2, method, m.Type(), funcGraphs)
				}
			}

			methods := prog.MethodSets.MethodSet(m.Type().Underlying())
			for i := 0; i < methods.Len(); i++ {
				sel := methods.At(i)
				fmt.Fprintf(outfile2, "\tMethod: %v\n", sel.Obj().Type())
				method := prog.MethodValue(sel)
				if method != nil {
					iterateFunc(outfile2, method, m.Type(), funcGraphs)
				}
			}

		default:
			fmt.Fprintf(outfile2, "\tUnknown member type: %T\n", m)
		}
	}
}

func iterateFunc(outFile *os.File, fn *ssa.Function, memberType types.Type, funcGraphs map[string]*ssa_graph.SSAGraph) {
	shortFuncPath := getShortFunctionPath(fn.String())
	serviceName := extractServiceNameFromShortFunctionPath(shortFuncPath)

	fmt.Printf("[SSA] iterating function %s\n", shortFuncPath)

	graph := ssa_graph.NewGraph(fn.Pkg.Pkg.Name(), shortFuncPath, serviceName)
	if _, exists := funcGraphs[shortFuncPath]; exists {
		log.Printf("ssa_graph for function (%s) already exists\n", shortFuncPath)
		log.Println("skipping...")
		return
	}
	funcGraphs[shortFuncPath] = graph
	fmt.Printf("added new ssa_graph for function (%s)\n", shortFuncPath)

	var visited = make(map[ssa.Value]bool)

	fmt.Fprintf(outFile, "\t\tParameters:\n")
	for i, param := range fn.Params {
		fmt.Fprintf(outFile, "\t\t\t%s = %s\n", param.Name(), param.String())
		parseValue(graph, nil, -i-1, param, visited)
	}

	fmt.Fprintf(outFile, "Function: %s\n", shortFuncPath)
	for i, block := range fn.Blocks {
		fmt.Fprintf(outFile, "Block #%d: %s.%s\n", i, shortFuncPath, block.Comment)
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

func parseInstr(graph *ssa_graph.SSAGraph, instr ssa.Instruction, instrIdx int, visited map[ssa.Value]bool) *ssa_graph.SSANode {
	fmt.Printf("[A] %02d [%T] %v\n", instrIdx, instr, instr.String())

	id := getInstructionID(instr)
	if id == "" { // e.graph., conditions or jumps (instructions and not values)
		log.Printf("skipping instruction with invalid id: %v\n", instr)
		return nil
	}

	if val, ok := instr.(ssa.Value); ok {
		return parseValue(graph, instr, instrIdx, val, visited)
	}
	node := ssa_graph.RegisterNewNode(graph, instr, id)

	switch t := instr.(type) {
	case *ssa.Store:
		// 04 [store] *t1 = currency
		addrNode := parseValue(graph, instr, instrIdx, t.Addr, visited)
		valNode := parseValue(graph, instr, instrIdx, t.Val, visited)

		graph.CreateAndAddNewEdge(addrNode, node, ssa_graph.EDGE_STORE, 0, "")
		graph.CreateAndAddNewEdge(valNode, node, ssa_graph.EDGE_USAGE, 0, "")

		fmt.Printf("ADDING EDGE FOR ADDR NDOE AND VAL NODE: %v // %v \n", t.Addr, t.Val)
	case *ssa.Return:
		for _, res := range t.Results {
			resNode := parseValue(graph, instr, instrIdx, res, visited)
			graph.CreateAndAddNewEdge(resNode, node, ssa_graph.EDGE_STORE, 0, "")
		}
	default:
		fmt.Printf("[1] ignoring... %02d [%T] %v\n", instrIdx, instr, instr.String())
	}

	return node
}

func parseValue(graph *ssa_graph.SSAGraph, instr ssa.Instruction, instrIdx int, val ssa.Value, visited map[ssa.Value]bool) *ssa_graph.SSANode {
	fmt.Printf("[B] %02d [%T] %v\n", instrIdx, val, val.String())

	if visited[val] {
		return graph.GetNodeByName(val.Name())
	}
	visited[val] = true

	id := getValueID(val)
	if id == "" { // sanity check
		log.Fatalf("unexpected invalid id for value: %v\n", val)
		return nil
	}

	node, exists := graph.GetNodeByNameIfExists(val.Name())
	if !exists {
		node = ssa_graph.RegisterNewNodeValue(graph, instr, val, id)
	}

	switch t := val.(type) {
	case *ssa.Call:
		for _, arg := range t.Call.Args {
			for _, edges := range graph.GetEdgesFromNode(node) {
				if edges.GetToNode().GetName() == arg.Name() {
					fmt.Printf("[INFO] skipping arg edge for %s\n", t.Name())
					continue
				}
			}
			for _, edges := range graph.GetEdgesToNode(node) {
				if edges.GetFromNode().GetName() == arg.Name() {
					fmt.Printf("[INFO] skipping arg edge for %s\n", t.Name())
					continue
				}
			}
			argNode := parseValue(graph, instr, instrIdx, arg, visited)
			graph.CreateAndAddNewEdge(argNode, node, ssa_graph.EDGE_STORE, 0, "")
		}
	case *ssa.Alloc:
		// nothing to do
	case *ssa.Slice:
		// nothing to do
		targetNode := parseValue(graph, instr, instrIdx, t.X, visited)
		graph.CreateAndAddNewEdge(targetNode, node, ssa_graph.EDGE_USAGE, 0, "")
	case *ssa.FieldAddr:
		// 00 [field] t27 = &t0.Items [#3]
		targetNode := parseValue(graph, instr, instrIdx, t.X, visited)
		graph.CreateAndAddNewEdge(targetNode, node, ssa_graph.EDGE_FIELD, 0, utils.FieldIndexToName(t))
	case *ssa.IndexAddr:
		targetNode := parseValue(graph, instr, instrIdx, t.X, visited)
		//FIXME: should parse value for t.Index
		graph.CreateAndAddNewEdge(targetNode, node, ssa_graph.EDGE_FIELD, 0, t.Index.String())
	case *ssa.UnOp:
		// 01 [unary] t14 = *t13
		// 05 [unary] t31 = *t30
		targetNode := parseValue(graph, instr, instrIdx, t.X, visited)
		graph.CreateAndAddNewEdge(targetNode, node, ssa_graph.EDGE_LOAD, 0, "")

	case *ssa.MakeInterface: // same as *ssa.UnOp
		targetNode := parseValue(graph, instr, instrIdx, t.X, visited)
		graph.CreateAndAddNewEdge(targetNode, node, ssa_graph.EDGE_USAGE, 0, "")
	case *ssa.Convert:
		targetNode := parseValue(graph, instr, instrIdx, t.X, visited)
		graph.CreateAndAddNewEdge(targetNode, node, ssa_graph.EDGE_USAGE, 0, "")

	case *ssa.Parameter:
		// nothing to do

	case *ssa.Global:
		// nothing to do

	case *ssa.Phi:
		for _, phiEdge := range t.Edges {
			for _, edges := range graph.GetEdgesFromNode(node) {
				if edges.GetToNode().GetName() == phiEdge.Name() {
					fmt.Printf("[INFO] skipping phi edge for %s\n", t.Name())
					continue
				}
			}
			for _, edges := range graph.GetEdgesToNode(node) {
				if edges.GetFromNode().GetName() == phiEdge.Name() {
					fmt.Printf("[INFO] skipping phi edge for %s\n", t.Name())
					continue
				}
			}
			edgeNode := parseValue(graph, instr, instrIdx, phiEdge, visited)
			graph.CreateAndAddNewEdge(edgeNode, node, ssa_graph.EDGE_STORE, 0, "")
		}

	case *ssa.Extract:
		extractFromNode := parseValue(graph, instr, instrIdx, t.Tuple, visited)
		graph.CreateAndAddNewEdge(extractFromNode, node, ssa_graph.EDGE_USAGE, t.Index, "")

	case *ssa.BinOp:
		xNode := parseValue(graph, instr, instrIdx, t.X, visited)
		yNode := parseValue(graph, instr, instrIdx, t.Y, visited)
		graph.CreateAndAddNewEdge(xNode, node, ssa_graph.EDGE_STORE, 0, "")
		graph.CreateAndAddNewEdge(yNode, node, ssa_graph.EDGE_USAGE, 0, "")

	default:
		fmt.Printf("[2] ignoring... %s [%T] %s = %v\n", id, val, val.Name(), val.String())
	}
	return node
}

func getInstructionID(instr ssa.Instruction) string {
	if !instr.Pos().IsValid() { // meaning there is no position
		n, err := rand.Int(rand.Reader, big.NewInt(1<<31))
		if err != nil {
			return ""
		}
		return "instr_" + instrString(instr) + "_" + fmt.Sprintf("%d", n)
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
					return ""
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
