package tainter

import (
	"fmt"
	"slices"
	"strings"

	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/graph"
	"analyzer/pkg/utils"
)

type TaintMode int

const (
	TAINT_MARK_UPPER TaintMode = iota
	TAINT_CHECK_UPPER
)

type TaintInfo struct {
	dbfield string
	path    string
	val     ssa.Value
}

func NewTaintInfo(database string, collection string) TaintInfo {
	return TaintInfo{
		dbfield: database + "." + collection,
		path:    "",
	}
}

func NewTaintInfoWithDbField(dbfield string) TaintInfo {
	return TaintInfo{
		dbfield: dbfield,
		path:    "",
	}
}

func (t TaintInfo) objectFullPath() string {
	return "_obj" + t.path
}

func (t TaintInfo) objectPath() string {
	return t.path
}

func (t TaintInfo) databaseField() string {
	return t.dbfield
}

func (t TaintInfo) updateValue(val ssa.Value) TaintInfo {
	t.val = val
	return t
}

func (t TaintInfo) updatePathPrefix(prefix string) TaintInfo {
	t.path = prefix + t.path
	return t
}

type CheckTaintInfo struct {
	indirectTaints  []string
	inheritedTaints map[string][]string
}

func (t *CheckTaintInfo) addToInheritedTaints(objPath string, dbField string) {
	if !slices.Contains(t.inheritedTaints[objPath], dbField) {
		t.inheritedTaints[objPath] = append(t.inheritedTaints[objPath], dbField)
	}
}

func (t *CheckTaintInfo) addToIndirectTaints(field string) {
	if !slices.Contains(t.indirectTaints, field) {
		t.indirectTaints = append(t.indirectTaints, field)
	}
}

func doTaintNode(node *graph.Node, taintInfo TaintInfo, taintMode TaintMode) {
	switch taintMode {
	case TAINT_MARK_UPPER:
		// note that objfields/dbfields already have "." before them
		fmt.Printf("[TAINT] [1] tainting node (%s) for objpath (%s) and dbfield (%s)\n", node.String(), taintInfo.objectFullPath(), taintInfo.databaseField())
		ok := node.AddTaintIfNotExists(taintInfo.objectFullPath(), taintInfo.databaseField())
		if ok {
			fmt.Printf("\t[TAINT] OK!\n")
		}
	case TAINT_CHECK_UPPER:
		fmt.Printf("[TAINT] [2] tainting node (%s) for objpath (%s) and dbfield (%s)\n", node.String(), taintInfo.objectFullPath(), taintInfo.databaseField()+taintInfo.objectPath())
		ok := node.AddTaintIfNotExists(taintInfo.objectFullPath(), taintInfo.databaseField()+taintInfo.objectPath())
		if ok {
			fmt.Printf("\t[TAINT] OK!\n")
		}
	}
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
			doTaintNode(edge.GetToNode(), taintInfo, TAINT_MARK_UPPER)

			backwardsAnalysis(g, edge.GetToNode().GetValue(), taintInfo, visited, TAINT_MARK_UPPER, nil)
		}
	}
}

func getObjectPathDiff(longPath1 string, shortPath2 string) string {
	longPath1 = strings.TrimPrefix(longPath1, "_obj")
	shortPath2 = strings.TrimPrefix(shortPath2, "_obj")
	// i.e., pathTop - pathBottomRel
	return strings.TrimPrefix(longPath1, shortPath2)
}

func backwardsAnalysis(g *graph.Graph, val ssa.Value, taintInfo TaintInfo, visited map[TaintInfo]bool, taintMode TaintMode, checkTaintInfo *CheckTaintInfo) {
	taintInfo = taintInfo.updateValue(val)

	fmt.Printf("[BACKWARD] visiting %s: %s // %v // TAINT INFO (%s, %s)\n", val.Name(), val.String(), val, taintInfo.path, taintInfo.dbfield)
	if visited[taintInfo] {
		fmt.Printf("\t[BACKWARD] skipping value %s: %s\n", val.Name(), val.String())
		return
	}
	visited[taintInfo] = true

	node := g.GetNodeByName(val.Name())

	switch taintMode {
	case TAINT_MARK_UPPER:
		doTaintNode(node, taintInfo, taintMode)
	case TAINT_CHECK_UPPER:
		fmt.Printf("[BACKWARD] checking upper taints: %v\n", node.GetTaints())
		for objPath, dbFields := range node.GetTaints() {

			fmt.Printf("[BACKWARD] comparing prefixes:\n\t - tainted obj path:\t %s\n\t - bottom to upper:\t %s\n", objPath, taintInfo.objectFullPath())

			if strings.HasPrefix(taintInfo.objectFullPath(), objPath) && taintInfo.objectFullPath() != objPath {
				// e.g.,
				// existing path: 	_obj
				// current path: 	_obj.Shipping
				//
				// in this case, '_obj.Shipping' has prefix '_obj'
				// as result, we may get:
				//
				// existing taint: 	_obj			@ order_db.order
				// potential taint: _obj.Shipping 	@ order_db.order.Shipping
				for _, dbField := range dbFields {
					// save the taint in the upper node
					taintInfoTmp := taintInfo
					taintInfoTmp.dbfield = dbField
					doTaintNode(node, taintInfoTmp, taintMode)

					// so that we can later taint the bottom node
					dbFieldIndirect := taintInfoTmp.databaseField() + taintInfo.objectPath()
					checkTaintInfo.addToIndirectTaints(dbFieldIndirect)
				}
				break
			} else if strings.HasPrefix(objPath, taintInfo.objectFullPath()) { // also true if strings are equal
				// e.g.,
				// upper's taint: 		_obj.PostID @ posts_db.post.PostID
				// bottom's path: 		_obj.PostID
				// => bottom's taint: 	_obj		@ posts_db.post.PostID

				pathDiff := getObjectPathDiff(objPath, taintInfo.objectFullPath())
				for _, dbField := range dbFields {
					checkTaintInfo.addToInheritedTaints(pathDiff, dbField)
				}
			}
		}
	}

	switch t := val.(type) {
	case *ssa.MakeInterface:
		backwardsAnalysis(g, t.X, taintInfo, visited, taintMode, checkTaintInfo)
	case *ssa.UnOp:
		backwardsAnalysis(g, t.X, taintInfo, visited, taintMode, checkTaintInfo)
	case *ssa.Phi:
		// includes values in t.Edges + other nodes pointing to
		for _, edge := range g.GetEdgesFromNode(g.GetNodeByName(t.Name())) {
			// in case it points to an instruction like store we need to fetch the value
			// (in this case, this corresponds to the variable where something is being stored, and NOT the value being stored)
			if edge.GetToNode().GetInstruction() != nil && edge.GetToNode().GetValue() == nil {
				if taintMode == TAINT_MARK_UPPER {
					doTaintNode(edge.GetToNode(), taintInfo, taintMode)
					for _, edge2 := range g.GetEdgesToNode(edge.GetToNode()) {
						backwardsAnalysis(g, edge2.GetFromNode().GetValue(), taintInfo, visited, taintMode, checkTaintInfo)
					}
				}
			}
		}
	case *ssa.FieldAddr:
		fieldName := utils.FieldIndexToName(t)
		fmt.Printf("[BACKWARD] field addr %s, tainting %s\n", fieldName, t.X.String())
		// add after
		taintInfoTmp := taintInfo
		taintInfoTmp = taintInfoTmp.updatePathPrefix("." + fieldName)
		backwardsAnalysis(g, t.X, taintInfoTmp, visited, taintMode, checkTaintInfo)
	case *ssa.IndexAddr:
		// add after
		fmt.Printf("[BACKWARD] index addr %s, tainting %s\n", t.Index.String(), t.X.String())
		taintInfoTmp := taintInfo
		taintInfoTmp = taintInfoTmp.updatePathPrefix("[*]")
		backwardsAnalysis(g, t.X, taintInfoTmp, visited, taintMode, checkTaintInfo)
	case *ssa.Parameter, *ssa.Alloc:
		if taintMode == TAINT_MARK_UPPER {
			doTaintNode(node, taintInfo, taintMode)
		}
	default:
		fmt.Printf("[BACKWARD] ignoring value: [%T] %v\n", val, val)
	}

	if taintMode == TAINT_MARK_UPPER {
		// if its fieldaddr then we use the objfield and dbfield
		// from the parameters and not the updated ones
		doTaintPointerToSets(g, val, taintInfo, visited)
	}

	fmt.Printf("[BACKWARD] exit %s: %s // %v // TAINT INFO (%s, %s)\n", val.Name(), val.String(), val, taintInfo.path, taintInfo.dbfield)

}

func RunTaint(g *graph.Graph) {
	var nodes []*graph.Node
	for _, node := range g.GetNodes() {
		if call, _, ok := isDatabaseCall(node.GetInstruction()); ok {
			var database, collection string

			valDatabase := call.Call.Args[2]
			if c, ok := valDatabase.(*ssa.Const); ok {
				database = strings.Trim(c.Value.ExactString(), "\"")
			}

			valCollectionOrTopic := call.Call.Args[3]
			if c, ok := valCollectionOrTopic.(*ssa.Const); ok {
				collection = strings.Trim(c.Value.ExactString(), "\"")
			}

			valDocumentOrMessage := call.Call.Args[4]

			fmt.Printf("[RUN] tainting object %s.%s: %s\n", database, collection, valDocumentOrMessage.String())
			visited := make(map[TaintInfo]bool)
			taintInfo := NewTaintInfo(database, collection)
			backwardsAnalysis(g, valDocumentOrMessage, taintInfo, visited, TAINT_MARK_UPPER, nil)

			node := g.GetNodeByName(valDocumentOrMessage.Name())
			nodes = append(nodes, node)
		}

		// check for common taints
		for _, originNode := range nodes {
			fmt.Printf("[RUN] visiting node (origin): %v\n", originNode.String())
			for _, edge := range recurseEdgesBackwardsUntilLoadFrom(g, originNode, nil, make(map[*graph.Node]bool)) {
				// expecting only one node
				node := edge.GetFromNode()
				fmt.Printf("\t[RUN] visiting node (load): %v\n", node.String())

				for _, edge := range recurseEdgesForwardUntilStoreTo(g, node, nil, make(map[*graph.Node]bool)) {
					fromStoreNode := edge.GetFromNode()
					toStoreNode := edge.GetToNode()

					var valNodesOnStore []*graph.Node
					if sr, ok := toStoreNode.GetInstruction().(*ssa.Store); ok {
						valNode := g.GetNodeByName(sr.Val.Name())

						// avoid duplicates
						if !slices.Contains(valNodesOnStore, valNode) {
							valNodesOnStore = append(valNodesOnStore, valNode)
						}
					}

					for _, fromValNode := range valNodesOnStore {
						visited := make(map[TaintInfo]bool)
						taintInfo := TaintInfo{}
						checkTaintInfo := &CheckTaintInfo{inheritedTaints: make(map[string][]string)}
						backwardsAnalysis(g, fromValNode.GetValue(), taintInfo, visited, TAINT_CHECK_UPPER, checkTaintInfo)

						// indirect taints
						for _, dbfield := range checkTaintInfo.indirectTaints {
							taintInfo := TaintInfo{
								dbfield: dbfield,
								val:     fromValNode.GetValue(),
							}

							// not needed but helps in visualization graph
							doTaintNode(toStoreNode, taintInfo, TAINT_MARK_UPPER)
							doTaintNode(fromValNode, taintInfo, TAINT_MARK_UPPER)

							visited2 := make(map[TaintInfo]bool)
							taintInfo2 := NewTaintInfoWithDbField(dbfield)
							backwardsAnalysis(g, fromStoreNode.GetValue(), taintInfo2, visited2, TAINT_MARK_UPPER, nil)
						}
					}
				}
			}
		}

		if call, _, ok := isServiceCall(node.GetInstruction()); ok {
			for _, arg := range call.Call.Args[2:] {
				fmt.Printf("[RUN] checking taint for service call with arg: %s\n", arg.String())
				node := g.GetNodeByName(arg.Name())
				nodes = append(nodes, node)
			}
		}

		// check for upper taints affecting the current database/service calls objects
		for _, originNode := range nodes {
			fmt.Println()
			fmt.Printf("[RUN] check upper taints for node: %v\n", originNode.String())
			visited := make(map[TaintInfo]bool)
			taintInfo := TaintInfo{}
			checkTaintInfo := &CheckTaintInfo{inheritedTaints: make(map[string][]string)}
			backwardsAnalysis(g, originNode.GetValue(), taintInfo, visited, TAINT_CHECK_UPPER, checkTaintInfo)
			node = g.GetNodeByName(originNode.GetValue().Name())

			// indirect taints
			for _, dbfield := range checkTaintInfo.indirectTaints {
				taintInfo := TaintInfo{
					dbfield: dbfield,
					val:     originNode.GetValue(),
				}
				doTaintNode(node, taintInfo, TAINT_MARK_UPPER)
			}

			// inherited taints
			for objpath, dbfields := range checkTaintInfo.inheritedTaints {
				fmt.Printf("[RUN] check inherited taints for objpath (%s): %v\n", objpath, dbfields)
				for _, dbfield := range dbfields {
					taintInfo := TaintInfo{
						dbfield: dbfield,
						path:    objpath,
						val:     originNode.GetValue(),
					}
					doTaintNode(node, taintInfo, TAINT_MARK_UPPER)
				}
			}
		}
	}
}

func recurseEdgesBackwardsUntilLoadFrom(g *graph.Graph, node *graph.Node, storeEdges []*graph.Edge, visited map[*graph.Node]bool) []*graph.Edge {
	if _, ok := visited[node]; ok {
		return storeEdges
	}
	visited[node] = true

	for _, edge := range g.GetEdgesToNode(node) {
		if edge.GetType() == graph.EDGE_LOAD {
			storeEdges = append(storeEdges, edge)
		} else {
			storeEdges = append(storeEdges, recurseEdgesBackwardsUntilLoadFrom(g, edge.GetFromNode(), storeEdges, visited)...)
		}
	}
	return storeEdges
}

func recurseEdgesForwardUntilStoreTo(g *graph.Graph, node *graph.Node, storeEdges []*graph.Edge, visited map[*graph.Node]bool) []*graph.Edge {
	if _, ok := visited[node]; ok {
		return storeEdges
	}
	visited[node] = true

	for _, edge := range g.GetEdgesFromNode(node) {
		if edge.GetType() == graph.EDGE_STORE {
			storeEdges = append(storeEdges, edge)
		} else if edge.GetType() == graph.EDGE_FIELD || edge.GetType() == graph.EDGE_INDEX {
			storeEdges = append(storeEdges, recurseEdgesForwardUntilStoreTo(g, edge.GetToNode(), storeEdges, visited)...)
		}
	}
	return storeEdges
}

func isDatabaseCall(instr ssa.Instruction) (*ssa.Call, *ssa.Function, bool) {
	if call, ok := instr.(*ssa.Call); ok {
		if fn, ok := call.Call.Value.(*ssa.Function); ok && len(fn.Params) > 0 {
			maybeRcv := fn.Params[0]
			if maybeRcv.Type().String() == "*main.MongoDB" && fn.Name() == "Insert" || fn.Name() == "Find" {
				return call, fn, true
			}
			if maybeRcv.Type().String() == "*main.RabbitMQ" && fn.Name() == "Push" {
				return call, fn, true
			}
		}
	}
	return nil, nil, false
}

func isServiceCall(instr ssa.Instruction) (*ssa.Call, *ssa.Function, bool) {
	if call, ok := instr.(*ssa.Call); ok {
		if fn, ok := call.Call.Value.(*ssa.Function); ok && len(fn.Params) > 0 {
			maybeRcv := fn.Params[0]
			if maybeRcv.Type().String() == "*main.ShippingService" && fn.Name() == "NewShipment" {
				return call, fn, true
			}
			if maybeRcv.Type().String() == "*main.SkuService" && fn.Name() == "GetSku" {
				return call, fn, true
			}
			if maybeRcv.Type().String() == "*main.AnalyticsService" && fn.Name() == "UpdateAnalytics" {
				return call, fn, true
			}
		}
	}
	return nil, nil, false
}
