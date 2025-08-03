package parser

import (
	"crypto/rand"
	"fmt"
	"go/types"
	"log"
	"math/big"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/types/typeutil"

	"analyzer/pkg/app"
	"analyzer/pkg/ssagraph"
	"analyzer/pkg/utils"
)

func RunSSAAnalysis(app *app.App, prog *ssa.Program, pkg *ssa.Package, funcGraphs map[string]*ssagraph.SSAGraph) {
	path1 := fmt.Sprintf("output/%s/ssa/%s.out", app.GetName(), pkg.Pkg.Name())
	if err := os.MkdirAll(filepath.Dir(path1), 0755); err != nil {
		log.Fatal(err)
	}
	outfile1, err := os.Create(path1)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile1.Close()

	path2 := fmt.Sprintf("output/%s/ssa/%s.ssa", app.GetName(), pkg.Pkg.Name())
	if err := os.MkdirAll(filepath.Dir(path2), 0755); err != nil {
		log.Fatal(err)
	}
	outfile2, err := os.Create(path2)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile2.Close()

	for _, member := range pkg.Members {
		switch m := member.(type) {
		case *ssa.Function:
			iterateFunc(app, outfile2, m, nil, funcGraphs)

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
					iterateFunc(app, outfile2, method, m.Type(), funcGraphs)
				}
			}

			methods := prog.MethodSets.MethodSet(m.Type().Underlying())
			for i := 0; i < methods.Len(); i++ {
				sel := methods.At(i)
				fmt.Fprintf(outfile2, "\tMethod: %v\n", sel.Obj().Type())
				method := prog.MethodValue(sel)
				if method != nil {
					iterateFunc(app, outfile2, method, m.Type(), funcGraphs)
				}
			}

		default:
			fmt.Fprintf(outfile2, "\tUnknown member type: %T\n", m)
		}
	}
}

func iterateFunc(app *app.App, outFile *os.File, fn *ssa.Function, memberType types.Type, funcGraphs map[string]*ssagraph.SSAGraph) {
	shortFuncPath := utils.GetShortFunctionPath(fn.String())
	serviceName := utils.ExtractServiceNameFromShortFunctionPath(shortFuncPath)
	methodName := utils.ExtractMethodNameFromShortFunctionPath(shortFuncPath)

	fmt.Printf("[SSA] iterating function %s\n", shortFuncPath)

	graph := ssagraph.NewGraph(app, fn.Pkg.Pkg.Name(), shortFuncPath, serviceName, methodName)
	if _, exists := funcGraphs[shortFuncPath]; exists {
		fmt.Printf("[SSA] ssagraph for function (%s) already exists\n", shortFuncPath)
		fmt.Println("[SSA] skipping...")
		return
	}
	funcGraphs[shortFuncPath] = graph
	fmt.Printf("[SSA] added new ssagraph for function (%s)\n", shortFuncPath)

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

func parseInstr(graph *ssagraph.SSAGraph, instr ssa.Instruction, instrIdx int, visited map[ssa.Value]bool) *ssagraph.SSANode {
	fmt.Printf("[SSA] [A] %02d [%T] %v\n", instrIdx, instr, instr.String())

	id := computeInstructionID(instr)
	if id == "" { // e.graph., conditions or jumps (instructions and not values)
		fmt.Printf("[SSA] skipping instruction with invalid id: %v\n", instr)
		return nil
	}

	if val, ok := instr.(ssa.Value); ok {
		return parseValue(graph, instr, instrIdx, val, visited)
	}
	node := ssagraph.RegisterNewNode(graph, instr, id)

	switch t := instr.(type) {
	case *ssa.Store:
		// 04 [store] *t1 = currency
		addrNode := parseValue(graph, instr, instrIdx, t.Addr, visited)
		valNode := parseValue(graph, instr, instrIdx, t.Val, visited)

		graph.CreateAndAddNewEdge(addrNode, node, ssagraph.EDGE_STORE_ADDRESS, 0, "")
		graph.CreateAndAddNewEdge(valNode, node, ssagraph.EDGE_STORE_VALUE, 0, "")
	case *ssa.Return:
		for _, res := range t.Results {
			resNode := parseValue(graph, instr, instrIdx, res, visited)
			graph.CreateAndAddNewEdge(resNode, node, ssagraph.EDGE_STORE_ADDRESS, 0, "")
		}
	default:
		fmt.Printf("[SSA] [1] ignoring... %02d [%T] %v\n", instrIdx, instr, instr.String())
	}

	return node
}

func parseValue(graph *ssagraph.SSAGraph, instr ssa.Instruction, instrIdx int, val ssa.Value, visited map[ssa.Value]bool) *ssagraph.SSANode {
	fmt.Printf("[SSA] [B] %02d [%T] %v\n", instrIdx, val, val.String())

	if visited[val] {
		return graph.GetNodeByName(val.Name())
	}
	visited[val] = true

	id := computeValueID(val)
	if id == "" { // sanity check
		log.Fatalf("unexpected invalid id for value: %v\n", val)
		return nil
	}

	node, exists := graph.GetNodeByNameIfExists(val.Name())
	if !exists {
		node = ssagraph.RegisterNewNodeValue(graph, instr, val, id)
	}

	switch t := val.(type) {
	case *ssa.Call:
		for _, arg := range t.Call.Args {
			for _, edges := range graph.GetEdgesFromNode(node) {
				if edges.GetToNode().GetName() == arg.Name() {
					fmt.Printf("[SSA] skipping arg edge for %s\n", t.Name())
					continue
				}
			}
			for _, edges := range graph.GetEdgesToNode(node) {
				if edges.GetFromNode().GetName() == arg.Name() {
					fmt.Printf("[SSA] skipping arg edge for %s\n", t.Name())
					continue
				}
			}
			argNode := parseValue(graph, instr, instrIdx, arg, visited)
			graph.CreateAndAddNewEdge(argNode, node, ssagraph.EDGE_STORE_ADDRESS, 0, "")
		}
	case *ssa.Alloc:
		// nothing to do
	case *ssa.Slice:
		// nothing to do
		targetNode := parseValue(graph, instr, instrIdx, t.X, visited)
		graph.CreateAndAddNewEdge(targetNode, node, ssagraph.EDGE_USAGE, 0, "")
	case *ssa.FieldAddr:
		// 00 [field] t27 = &t0.Items [#3]
		targetNode := parseValue(graph, instr, instrIdx, t.X, visited)
		graph.CreateAndAddNewEdge(targetNode, node, ssagraph.EDGE_FIELD, 0, utils.FieldIndexToName(t))
	case *ssa.IndexAddr:
		targetNode := parseValue(graph, instr, instrIdx, t.X, visited)
		//FIXME: should parse value for t.Index
		graph.CreateAndAddNewEdge(targetNode, node, ssagraph.EDGE_FIELD, 0, t.Index.String())
	case *ssa.UnOp:
		// 01 [unary] t14 = *t13
		// 05 [unary] t31 = *t30
		targetNode := parseValue(graph, instr, instrIdx, t.X, visited)
		graph.CreateAndAddNewEdge(targetNode, node, ssagraph.EDGE_LOAD, 0, "")

	case *ssa.MakeInterface: // same as *ssa.UnOp
		targetNode := parseValue(graph, instr, instrIdx, t.X, visited)
		graph.CreateAndAddNewEdge(targetNode, node, ssagraph.EDGE_USAGE, 0, "")
	case *ssa.Convert:
		targetNode := parseValue(graph, instr, instrIdx, t.X, visited)
		graph.CreateAndAddNewEdge(targetNode, node, ssagraph.EDGE_USAGE, 0, "")

	case *ssa.Parameter:
		graph.AddParameter(node)
		// nothing to do

	case *ssa.Global:
		// nothing to do

	case *ssa.Phi:
		for _, phiEdge := range t.Edges {
			for _, edges := range graph.GetEdgesFromNode(node) {
				if edges.GetToNode().GetName() == phiEdge.Name() {
					fmt.Printf("[SSA] skipping phi edge for %s\n", t.Name())
					continue
				}
			}
			for _, edges := range graph.GetEdgesToNode(node) {
				if edges.GetFromNode().GetName() == phiEdge.Name() {
					fmt.Printf("[SSA] skipping phi edge for %s\n", t.Name())
					continue
				}
			}
			edgeNode := parseValue(graph, instr, instrIdx, phiEdge, visited)
			graph.CreateAndAddNewEdge(edgeNode, node, ssagraph.EDGE_STORE_ADDRESS, 0, "")
		}

	case *ssa.Extract:
		extractFromNode := parseValue(graph, instr, instrIdx, t.Tuple, visited)
		graph.CreateAndAddNewEdge(extractFromNode, node, ssagraph.EDGE_EXTRACT, t.Index, "")

	case *ssa.BinOp:
		xNode := parseValue(graph, instr, instrIdx, t.X, visited)
		yNode := parseValue(graph, instr, instrIdx, t.Y, visited)
		graph.CreateAndAddNewEdge(xNode, node, ssagraph.EDGE_STORE_ADDRESS, 0, "")
		graph.CreateAndAddNewEdge(yNode, node, ssagraph.EDGE_USAGE, 0, "")

	default:
		fmt.Printf("[SSA] [2] ignoring... %s [%T] %s = %v\n", id, val, val.Name(), val.String())
	}
	return node
}

func computeInstructionID(instr ssa.Instruction) string {
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

func computeValueID(val ssa.Value) string {
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
