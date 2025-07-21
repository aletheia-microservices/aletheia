package tainter

import (
	"fmt"
	"strings"

	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/graph"
	"analyzer/pkg/utils"
)

type TaintInfo struct {
	database   string
	collection string
	path       string
	val        ssa.Value
}

func NewTaintInfo(database string, collection string) TaintInfo {
	return TaintInfo{
		database:   database,
		collection: collection,
		path:       "",
	}
}

func (t TaintInfo) updateValue(val ssa.Value) TaintInfo {
	t.val = val
	return t
}

func (t TaintInfo) updatePathPrefix(prefix string) TaintInfo {
	t.path = prefix + t.path
	return t
}

func doTaintNode(node *graph.Node, taintInfo TaintInfo) {
	// note that objfields/dbfields already have "." before them
	objField := "_obj" + taintInfo.path
	dbField := taintInfo.database + "." + taintInfo.collection
	node.AddTaintIfNotExists(objField, dbField)
}

func doTaintPointerToSets(g *graph.Graph, val ssa.Value, taintInfo TaintInfo, visited map[TaintInfo]bool) {
	node := g.GetNodeByName(val.Name())
	for _, edge := range g.GetEdgesFromNode(node) {
		if edge.GetType() == graph.EDGE_POINTS_TO {
			if edge.GetPath() != "" {
				// add before
				// note that both edge.path and objfields/dbfields already have "." before them
				taintInfo = taintInfo.updatePathPrefix(edge.GetPath())
			}
			doTaintNode(edge.GetToNode(), taintInfo)

			doTaintBackwards(g, edge.GetToNode().GetValue(), taintInfo, visited)
		}
	}
}

func doTaintBackwards(g *graph.Graph, val ssa.Value, taintInfo TaintInfo, visited map[TaintInfo]bool) {
	taintInfo = taintInfo.updateValue(val)

	fmt.Printf("visiting value %s: %s // %v // TAINT INFO = %v\n", val.Name(), val.String(), val, taintInfo)
	if visited[taintInfo] {
		fmt.Printf("\tskipping value %s: %s\n", val.Name(), val.String())
		return
	}
	visited[taintInfo] = true

	node := g.GetNodeByName(val.Name())
	doTaintNode(node, taintInfo)

	switch t := val.(type) {
	case *ssa.MakeInterface:
		doTaintBackwards(g, t.X, taintInfo, visited)
	case *ssa.UnOp:
		doTaintBackwards(g, t.X, taintInfo, visited)
	case *ssa.Phi:
		// includes values in t.Edges + other nodes pointing to
		for _, edge := range g.GetEdgesFromNode(g.GetNodeByName(t.Name())) {
			// in case it points to an instruction like store we need to fetch the value
			// (in this case, this corresponds to the variable where something is being stored, and NOT the value being stored)
			if edge.GetToNode().GetInstruction() != nil && edge.GetToNode().GetValue() == nil {
				doTaintNode(edge.GetToNode(), taintInfo)
				for _, edge2 := range g.GetEdgesToNode(edge.GetToNode()) {
					doTaintBackwards(g, edge2.GetFromNode().GetValue(), taintInfo, visited)
				}
			} /* else { // FIXME infinite loops
				doTaintBackwards(graph *Graph, edge.to.val, taintInfo, visited)
			} */
		}
	case *ssa.FieldAddr:
		fieldName := utils.FieldIndexToName(t)
		fmt.Printf("[INFO] field addr %s, tainting %s\n", fieldName, t.X.String())
		// add after
		taintInfoTmp := taintInfo
		taintInfoTmp = taintInfoTmp.updatePathPrefix("." + fieldName)
		doTaintBackwards(g, t.X, taintInfoTmp, visited)
	case *ssa.IndexAddr:
		// add after
		fmt.Printf("[INFO] index addr %s, tainting %s\n", t.Index.String(), t.X.String())
		taintInfoTmp := taintInfo
		taintInfoTmp = taintInfoTmp.updatePathPrefix("[*]")
		doTaintBackwards(g, t.X, taintInfoTmp, visited)
	case *ssa.Parameter, *ssa.Alloc:
		doTaintNode(node, taintInfo)
	default:
		fmt.Printf("[INFO] ignoring value: [%T] %v\n", val, val)
	}

	// if its fieldaddr then we use the objfield and dbfield
	// from the parameters and not the updated ones
	doTaintPointerToSets(g, val, taintInfo, visited)
}

func RunTaint(g *graph.Graph) {
	for _, node := range g.GetNodes() {
		call, _, ok := isMongoDBCall(node.GetInstruction())
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

			fmt.Printf("[INFO] tainting document %s.%s: %s\n", database, collection, valDocument.String())
			visited := make(map[TaintInfo]bool)
			taintInfo := NewTaintInfo(database, collection)
			doTaintBackwards(g, valDocument, taintInfo, visited)
		}
	}
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
