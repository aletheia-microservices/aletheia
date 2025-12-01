package parser

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/types/typeutil"

	"analyzer/pkg/app"
	"analyzer/pkg/ssagraph"
	"analyzer/pkg/utils"
)

func RunSSAAnalysis(app *app.App, prog *ssa.Program, pkg *ssa.Package, funcGraphs map[string]*ssagraph.SSAGraph) {
	path1 := fmt.Sprintf("output/%s/ssa/%s.out", app.GetName(), pkg.Pkg.Name())
	if err := os.MkdirAll(filepath.Dir(path1), 0755); err != nil {
		panic(err)
	}
	outfile1, err := os.Create(path1)
	if err != nil {
		panic(err)
	}
	defer outfile1.Close()

	path2 := fmt.Sprintf("output/%s/ssa/%s.ssa", app.GetName(), pkg.Pkg.Name())
	if err := os.MkdirAll(filepath.Dir(path2), 0755); err != nil {
		panic(err)
	}
	outfile2, err := os.Create(path2)
	if err != nil {
		panic(err)
	}
	defer outfile2.Close()

	for _, member := range pkg.Members {
		logrus.Tracef("[SSA] [%T] member: %v\n", member, member)
		switch m := member.(type) {
		case *ssa.Function:
			iterateFunc(app, outfile2, m, funcGraphs, false)

		case *ssa.Global:
			fmt.Fprintf(outfile2, "\tGlobal: %s, Type: %s\n", m.Name(), m.Type().String())

		case *ssa.Type:
			fmt.Fprintf(outfile2, "\tType: %s\n", m.Type())

			// this logic was copied from
			// package: golang.org/x/tools/go/ssa
			// file: print.go
			// function: func (p *Package) WriteTo(w io.Writer) (int64, error)
			for _, sel := range typeutil.IntuitiveMethodSet(m.Type(), &prog.MethodSets) {
				logrus.Tracef("\t[SSA] [INTUITIVE METHOD SET] [%T] (index=%v, indirect=%t): %v\n", sel, sel.Index(), sel.Indirect(), sel)
				method := prog.MethodValue(sel)
				if method != nil {
					fmt.Fprintf(outfile2, "\tMethod: %v\n", sel.Obj().Type())
					logrus.Tracef("\t[SSA] [INTUITIVE METHOD SET] [%T]: %v\n", method, method)
					if len(sel.Index()) != 1 {
						// when a structure has an embedded field its methods are promoted and
						// will appear for the current structure
						//
						// e.g. in dsb socialnetwork:
						// type claimsT struct {
						//		Username  string
						//		UserID    string
						//		Timestamp int64
						//		jwt.StandardClaims
						// }
						// where jwt.StandardClaims has methods Valid(), VerifyAudience(), etc.
						//
						// WORKAROUND: just ignore them for now
						logrus.Tracef("\t[SSA] [INTUITIVE METHOD SET] [%T]: skipping...\n", method)
						continue
					}
					iterateFunc(app, outfile2, method, funcGraphs, false)
				}
			}

			methods := prog.MethodSets.MethodSet(m.Type().Underlying())
			for i := 0; i < methods.Len(); i++ {
				sel := methods.At(i)
				logrus.Tracef("\t[SSA] [METHOD SET] [%T] (index=%v, indirect=%t): %v\n", sel, sel.Index(), sel.Indirect(), sel)
				fmt.Fprintf(outfile2, "\tMethod: %v\n", sel.Obj().Type())
				method := prog.MethodValue(sel)
				if method != nil {
					logrus.Tracef("\t[SSA] [METHOD SET] [%T]: %v\n", method, method)
					if len(sel.Index()) != 1 {
						// same reason as above when iterating IntuitiveMethodSet
						logrus.Tracef("\t[SSA] [METHOD SET] [%T]: skipping...\n", method)
						continue
					}
					iterateFunc(app, outfile2, method, funcGraphs, false)
				}
			}

		default:
			fmt.Fprintf(outfile2, "\tUnknown member type: %T\n", m)
		}
	}
}

func iterateFunc(app *app.App, outFile *os.File, fn *ssa.Function, funcGraphs map[string]*ssagraph.SSAGraph, goroutine bool) {
	shortFuncPath := utils.GetShortFunctionPath(fn.String())
	serviceName := utils.ExtractServiceNameFromShortFunctionPath(shortFuncPath)
	methodName := utils.ExtractMethodNameFromShortFunctionPath(shortFuncPath)

	logrus.Tracef("[SSA] iterating function (%s)\n", shortFuncPath)

	graph := ssagraph.NewGraph(app, fn.Pkg.Pkg.Name(), shortFuncPath, serviceName, methodName)
	if _, exists := funcGraphs[shortFuncPath]; exists {
		logrus.Tracef("[SSA] ssagraph for function (%s) already exists\n", shortFuncPath)
		logrus.Traceln("[SSA] skipping...")
		return
	}
	if goroutine {
		graph.EnableGoRoutine()
	}
	funcGraphs[shortFuncPath] = graph
	logrus.Tracef("[SSA] added new ssagraph for function (%s)\n", shortFuncPath)

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
			parseInstr(app, graph, instr, j, visited, outFile, funcGraphs)

			if val, ok := instr.(ssa.Value); ok {
				fmt.Fprintf(outFile, "\t\t\t%02d: %s = %s\n", j, val.Name(), instr.String())
			} else {
				fmt.Fprintf(outFile, "\t\t\t%02d: %s\n", j, instr.String())
			}
		}
	}
}

func parseInstr(app *app.App, graph *ssagraph.SSAGraph, instr ssa.Instruction, instrIdx int, visited map[ssa.Value]bool, outFile *os.File, funcGraphs map[string]*ssagraph.SSAGraph) *ssagraph.SSANode {
	logrus.Tracef("[SSA PARSE INSTR] %02d [%T] %v\n", instrIdx, instr, instr.String())

	id := utils.ComputeInstructionID(instr)
	if id == "" { // e.graph., conditions or jumps (instructions and not values)
		logrus.Tracef("[SSA PARSE INSTR] skipping instruction with invalid id: %v\n", instr)
		return nil
	}

	if val, ok := instr.(ssa.Value); ok {
		return parseValue(graph, instr, instrIdx, val, visited)
	}
	node := ssagraph.RegisterNewNodeInstr(graph, instr, id)

	switch t := instr.(type) {
	case *ssa.Store:
		// 04 [store] *t1 = currency
		addrNode := parseValue(graph, instr, instrIdx, t.Addr, visited)
		valNode := parseValue(graph, instr, instrIdx, t.Val, visited)

		graph.CreateAndAddNewEdge(addrNode, node, ssagraph.EDGE_STORE_ADDRESS, 0, "")
		graph.CreateAndAddNewEdge(valNode, node, ssagraph.EDGE_STORE_VALUE, 0, "")
	case *ssa.Return:
		var rets []*ssagraph.SSANode
		for _, ret := range t.Results {
			retNode := parseValue(graph, instr, instrIdx, ret, visited)
			rets = append(rets, retNode)
			graph.CreateAndAddNewEdge(retNode, node, ssagraph.EDGE_RETURN_ON, 0, "")
		}
		graph.AddReturnsToLst(rets)
	case *ssa.MapUpdate:
		mapNode := parseValue(graph, instr, instrIdx, t.Map, visited)
		keyNode := parseValue(graph, instr, instrIdx, t.Key, visited)
		valueNode := parseValue(graph, instr, instrIdx, t.Value, visited)

		index := "[*]"
		if val, ok := utils.ExtractStringFromValue(keyNode.GetValue()); ok {
			index = val
		}

		graph.CreateAndAddNewEdge(mapNode, node, ssagraph.EDGE_MAP_UPDATE, 0, index)
		graph.CreateAndAddNewEdge(keyNode, node, ssagraph.EDGE_MAP_KEY, 0, index)
		graph.CreateAndAddNewEdge(valueNode, node, ssagraph.EDGE_MAP_VALUE, 0, index)

	case *ssa.If, *ssa.Jump:
		// nothing to do

	case *ssa.Panic:
		// ignore

	case *ssa.Go:
		if makeClosure, ok := t.Call.Value.(*ssa.MakeClosure); ok {
			if fn, ok := makeClosure.Fn.(*ssa.Function); ok {
				fmt.Printf("make_closure_fn: %s\n", makeClosure.Fn)
				fmt.Printf("short_func_path: %s\n", utils.GetShortFunctionPath(fn.String()))
				logrus.WithField("instr", instr.String()).Warnf("[SSA PARSE INSTR] found *ssa.Go")
				iterateFunc(app, outFile, fn, funcGraphs, true)
			}
		}

	case *ssa.RunDefers, *ssa.Defer:
		// TODO
		logrus.Tracef("[SSA PARSE INSTR] ignoring... %02d [%T] %v\n", instrIdx, instr, instr.String())

	default:
		logrus.Fatalf("[SSA PARSE INSTR] ignoring... %02d [%T] %v\n", instrIdx, instr, instr.String())
	}

	return node
}

func parseValue(graph *ssagraph.SSAGraph, instr ssa.Instruction, instrIdx int, val ssa.Value, visited map[ssa.Value]bool) *ssagraph.SSANode {
	logrus.Tracef("[SSA PARSE VALUE] %02d [%T] %v\n", instrIdx, val, val.String())

	if visited[val] {
		return graph.GetNodeByName(val.Name())
	}
	visited[val] = true

	id := computeValueID(val)
	if id == "" { // sanity check
		logrus.Fatalf("[SSA PARSE VALUE] unexpected invalid id for value: %v\n", val)
		return nil
	}

	node, exists := graph.GetNodeByNameIfExists(val.Name())
	if !exists {
		node = ssagraph.RegisterNewNodeVal(graph, instr, val, id)
	}

	switch t := val.(type) {
	case *ssa.Call:
		for _, arg := range t.Call.Args {
			logrus.Tracef("[SSA PARSE VALUE] [CALL: %s] ARG: %s\n", t.Name(), arg.Name())
			for _, edges := range graph.GetEdgesFromNode(node) {
				if edges.GetToNode().GetName() == arg.Name() {
					logrus.Tracef("[SSA PARSE VALUE] skipping arg edge for %s\n", t.Name())
					continue
				}
			}
			for _, edges := range graph.GetEdgesToNode(node) {
				if edges.GetFromNode().GetName() == arg.Name() {
					logrus.Tracef("[SSA PARSE VALUE] skipping arg edge for %s\n", t.Name())
					continue
				}
			}
			argNode := parseValue(graph, instr, instrIdx, arg, visited)
			graph.CreateAndAddNewEdge(argNode, node, ssagraph.EDGE_ARG_ON_CALL, 0, "")
		}
		if t.Call.IsInvoke() {
			rcv := t.Call.Value
			logrus.Tracef("[SSA PARSE VALUE] [CALL: %s] RCV: %s\n", t.Name(), rcv.Name())
			rcvNode := parseValue(graph, instr, instrIdx, rcv, visited)
			graph.CreateAndAddNewEdge(rcvNode, node, ssagraph.EDGE_RECEIVER_ON_CALL, 0, "")
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
		param := utils.FieldIndexToName(t)

		graph.CreateAndAddNewEdge(targetNode, node, ssagraph.EDGE_FIELD, 0, param)
	case *ssa.IndexAddr:
		targetNode := parseValue(graph, instr, instrIdx, t.X, visited)
		param := "*"
		/* if index, ok := utils.ExtractStringFromValue(t.Index); ok {
			param = index
		} else {
			param = "*"
		} */
		graph.CreateAndAddNewEdge(targetNode, node, ssagraph.EDGE_INDEX, 0, param)
	case *ssa.Field:
		// e.g., [*ssa.Field] t151 = t150.StartPlace [#1]
		// where t150 is a map value
		targetNode := parseValue(graph, instr, instrIdx, t.X, visited)
		param := "*"
		graph.CreateAndAddNewEdge(targetNode, node, ssagraph.EDGE_FIELD, 0, param)
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

	case *ssa.FreeVar:
		// in case of go routines (variables that are not passed as parameters but SSA assumes this)
		graph.AddFreeVar(node)

	case *ssa.Const:
		// nothing to do

	case *ssa.MakeMap:
		logrus.Tracef("[SSA PARSE VALUE] MAKE MAP! %v\n", t)
		// nothing to do

	case *ssa.Global:
		// nothing to do

	case *ssa.Phi:
		for _, phiEdge := range t.Edges {
			for _, edges := range graph.GetEdgesFromNode(node) {
				if edges.GetToNode().GetName() == phiEdge.Name() {
					logrus.Tracef("[SSA PARSE VALUE] skipping phi edge for %s\n", t.Name())
					continue
				}
			}
			for _, edges := range graph.GetEdgesToNode(node) {
				if edges.GetFromNode().GetName() == phiEdge.Name() {
					logrus.Tracef("[SSA PARSE VALUE] skipping phi edge for %s\n", t.Name())
					continue
				}
			}
			edgeNode := parseValue(graph, instr, instrIdx, phiEdge, visited)
			graph.CreateAndAddNewEdge(edgeNode, node, ssagraph.EDGE_PHI_ON, 0, "")
		}

	case *ssa.Extract:
		extractFromNode := parseValue(graph, instr, instrIdx, t.Tuple, visited)
		graph.CreateAndAddNewEdge(extractFromNode, node, ssagraph.EDGE_EXTRACT, t.Index, "")

	case *ssa.BinOp:
		xNode := parseValue(graph, instr, instrIdx, t.X, visited)
		yNode := parseValue(graph, instr, instrIdx, t.Y, visited)
		graph.CreateAndAddNewEdge(xNode, node, ssagraph.EDGE_BINOP_X, 0, "")
		graph.CreateAndAddNewEdge(yNode, node, ssagraph.EDGE_BINOP_Y, 0, "")

	case *ssa.Lookup:
		xNode := parseValue(graph, instr, instrIdx, t.X, visited)
		idxNode := parseValue(graph, instr, instrIdx, t.Index, visited)

		index := "[*]"
		if val, ok := utils.ExtractStringFromValue(idxNode.GetValue()); ok {
			index = val
		}

		graph.CreateAndAddNewEdge(xNode, node, ssagraph.EDGE_LOOKUP_MAP, 0, index)
		graph.CreateAndAddNewEdge(idxNode, node, ssagraph.EDGE_LOOKUP_MAP_INDEX, 0, index)

	case *ssa.Range:
		// e.g., dsb_sn2 at PostStorageService.ReadPosts:
		// ----------------------------------------
		// t0 = make map[int64]bool
		// t71 = range t0
		// ----------------------------------------
		// for k := range unique_post_ids {
		// 	  unique_pids = append(unique_pids, k)
		// }
		// ----------------------------------------
		xNode := parseValue(graph, instr, instrIdx, t.X, visited)
		graph.CreateAndAddNewEdge(xNode, node, ssagraph.EDGE_RANGE_OF, 0, "")
	case *ssa.Next:
		// e.g., dsb_sn2 at PostStorageService.ReadPosts:
		// ----------------------------------------
		// t0 = make map[int64]bool
		// t71 = range t0
		// t74 = next t71
		// ----------------------------------------
		// for k := range unique_post_ids {
		// 	  unique_pids = append(unique_pids, k)
		// }
		// ----------------------------------------
		iterNode := parseValue(graph, instr, instrIdx, t.Iter, visited)
		graph.CreateAndAddNewEdge(iterNode, node, ssagraph.EDGE_ITERATOR_OF, 0, "")

	case *ssa.MakeClosure, *ssa.Select, *ssa.MakeSlice, *ssa.ChangeInterface, *ssa.Index,
		*ssa.TypeAssert, *ssa.ChangeType: // dsb_sn2
		// TODO
		logrus.Tracef("[SSA PARSE VALUE] ignoring ssa.Value... %s [%T] %s = %v\n", id, val, val.Name(), val.String())

	default:
		logrus.Fatalf("[SSA PARSE VALUE] unknown ssa.Value... %s [%T] %s = %v\n", id, val, val.Name(), val.String())
	}
	return node
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
