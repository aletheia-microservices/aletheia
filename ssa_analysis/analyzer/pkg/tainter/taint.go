package tainter

import (
	"fmt"
	"go/types"
	"log"
	"slices"
	"strings"

	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/ssa_graph"
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

func doTaintNode(node *ssa_graph.SSANode, taintInfo TaintInfo, taintMode TaintMode) {
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

func doTaintPointerToSets(graph *ssa_graph.SSAGraph, val ssa.Value, taintInfo TaintInfo, visited map[TaintInfo]bool) {
	node := graph.GetNodeByName(val.Name())
	for _, edge := range graph.GetEdgesFromNode(node) {
		if edge.GetType() == ssa_graph.EDGE_POINTS_TO {
			if edge.GetPath() != "" {
				// add before
				// note that both edge.path and objfields/dbfields already have "." before them
				taintInfo = taintInfo.updatePathPrefix(edge.GetPath())
			}
			doTaintNode(edge.GetToNode(), taintInfo, TAINT_MARK_UPPER)

			backwardsAnalysis(graph, edge.GetToNode().GetValue(), taintInfo, visited, TAINT_MARK_UPPER, nil)
		}
	}
}

func getObjectPathDiff(longPath1 string, shortPath2 string) string {
	longPath1 = strings.TrimPrefix(longPath1, "_obj")
	shortPath2 = strings.TrimPrefix(shortPath2, "_obj")
	// i.e., pathTop - pathBottomRel
	return strings.TrimPrefix(longPath1, shortPath2)
}

func backwardsAnalysis(graph *ssa_graph.SSAGraph, val ssa.Value, taintInfo TaintInfo, visited map[TaintInfo]bool, taintMode TaintMode, checkTaintInfo *CheckTaintInfo) {
	taintInfo = taintInfo.updateValue(val)

	fmt.Printf("[BACKWARD] visiting %s: %s // %v // TAINT INFO (%s, %s)\n", val.Name(), val.String(), val, taintInfo.path, taintInfo.dbfield)
	if visited[taintInfo] {
		fmt.Printf("\t[BACKWARD] skipping value %s: %s\n", val.Name(), val.String())
		return
	}
	visited[taintInfo] = true

	node := graph.GetNodeByName(val.Name())

	switch taintMode {
	case TAINT_MARK_UPPER:
		doTaintNode(node, taintInfo, taintMode)
	case TAINT_CHECK_UPPER:
		fmt.Printf("[BACKWARD] checking upper taints: %v\n", node.GetTaints())
		for objPath, dbFields := range node.GetTaints() {

			fmt.Printf("[BACKWARD] comparing prefixes:\n\t - tainted obj path:\t %s\n\t - bottom to upper:\t %s\n", objPath, taintInfo.objectFullPath())

			if strings.HasPrefix(taintInfo.objectFullPath(), objPath) && taintInfo.objectFullPath() != objPath {
				// e.graph.,
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
				// e.graph.,
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
		backwardsAnalysis(graph, t.X, taintInfo, visited, taintMode, checkTaintInfo)
	case *ssa.UnOp:
		backwardsAnalysis(graph, t.X, taintInfo, visited, taintMode, checkTaintInfo)
	case *ssa.Phi:
		// includes values in t.Edges + other nodes pointing to
		for _, edge := range graph.GetEdgesFromNode(graph.GetNodeByName(t.Name())) {
			// in case it points to an instruction like store we need to fetch the value
			// (in this case, this corresponds to the variable where something is being stored, and NOT the value being stored)
			if edge.GetToNode().GetInstruction() != nil && edge.GetToNode().GetValue() == nil {
				if taintMode == TAINT_MARK_UPPER {
					doTaintNode(edge.GetToNode(), taintInfo, taintMode)
					for _, edge2 := range graph.GetEdgesToNode(edge.GetToNode()) {
						backwardsAnalysis(graph, edge2.GetFromNode().GetValue(), taintInfo, visited, taintMode, checkTaintInfo)
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
		backwardsAnalysis(graph, t.X, taintInfoTmp, visited, taintMode, checkTaintInfo)
	case *ssa.IndexAddr:
		// add after
		fmt.Printf("[BACKWARD] index addr %s, tainting %s\n", t.Index.String(), t.X.String())
		taintInfoTmp := taintInfo
		taintInfoTmp = taintInfoTmp.updatePathPrefix("[*]")
		backwardsAnalysis(graph, t.X, taintInfoTmp, visited, taintMode, checkTaintInfo)
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
		doTaintPointerToSets(graph, val, taintInfo, visited)
	}

	fmt.Printf("[BACKWARD] exit %s: %s // %v // TAINT INFO (%s, %s)\n", val.Name(), val.String(), val, taintInfo.path, taintInfo.dbfield)

}

func RunTaint(graph *ssa_graph.SSAGraph) {
	var nodes []*ssa_graph.SSANode
	for _, node := range graph.GetNodes() {
		if database, collectionOrTopic, method, args, ok := isDatabaseCall(graph, node.GetInstruction()); ok {
			/* if node.String() == "t14: invoke t4.FindOne(ctx, t13, nil:[]go.mongodb.org/mongo-driver/bson/primitive.D...)" {
				log.Fatal("EXIT!")
			} */
			if node.String() == "nil:[]go.mongodb.org/mongo-driver/bson/primitive.D: nil:[]go.mongodb.org/mongo-driver/bson/primitive.D" {
				//FIXME (this variable is nil because it is not passed in the call and is optional but for some reason it's assuming it is a db call)
				continue
			}

			var argNodes []*ssa_graph.SSANode
			for _, arg := range args {
				argNodes = append(argNodes, graph.GetNodeByName(arg.Name()))
			}
			graph.AddDatabaseCall(node, argNodes, database, collectionOrTopic, method)


			fmt.Printf("[TAINT] got func args: %v\n", args)
			valDocumentOrMessage := args[0]

			fmt.Printf("[RUN] tainting object %s.%s: %s\n", database, collectionOrTopic, valDocumentOrMessage.String())
			visited := make(map[TaintInfo]bool)
			taintInfo := NewTaintInfo(database, collectionOrTopic)
			backwardsAnalysis(graph, valDocumentOrMessage, taintInfo, visited, TAINT_MARK_UPPER, nil)

			node := graph.GetNodeByName(valDocumentOrMessage.Name())
			nodes = append(nodes, node)
		}

		// check for common taints
		for _, originNode := range nodes {
			fmt.Printf("[RUN] visiting node (origin): %v\n", originNode.String())
			for _, edge := range recurseEdgesBackwardsUntilLoadFrom(graph, originNode, nil, make(map[*ssa_graph.SSANode]bool)) {
				// expecting only one node
				node := edge.GetFromNode()
				fmt.Printf("\t[RUN] visiting node (load): %v\n", node.String())

				for _, edge := range recurseEdgesForwardUntilStoreTo(graph, node, nil, make(map[*ssa_graph.SSANode]bool)) {
					fromStoreNode := edge.GetFromNode()
					toStoreNode := edge.GetToNode()

					var valNodesOnStore []*ssa_graph.SSANode
					if sr, ok := toStoreNode.GetInstruction().(*ssa.Store); ok {
						valNode := graph.GetNodeByName(sr.Val.Name())

						// avoid duplicates
						if !slices.Contains(valNodesOnStore, valNode) {
							valNodesOnStore = append(valNodesOnStore, valNode)
						}
					}

					for _, fromValNode := range valNodesOnStore {
						visited := make(map[TaintInfo]bool)
						taintInfo := TaintInfo{}
						checkTaintInfo := &CheckTaintInfo{inheritedTaints: make(map[string][]string)}
						backwardsAnalysis(graph, fromValNode.GetValue(), taintInfo, visited, TAINT_CHECK_UPPER, checkTaintInfo)

						// indirect taints
						for _, dbfield := range checkTaintInfo.indirectTaints {
							taintInfo := TaintInfo{
								dbfield: dbfield,
								val:     fromValNode.GetValue(),
							}

							// not needed but helps in visualization ssa_graph
							doTaintNode(toStoreNode, taintInfo, TAINT_MARK_UPPER)
							doTaintNode(fromValNode, taintInfo, TAINT_MARK_UPPER)

							visited2 := make(map[TaintInfo]bool)
							taintInfo2 := NewTaintInfoWithDbField(dbfield)
							backwardsAnalysis(graph, fromStoreNode.GetValue(), taintInfo2, visited2, TAINT_MARK_UPPER, nil)
						}
					}
				}
			}
		}

		if call, service, method, ok := isServiceCall(graph, node.GetInstruction()); ok {
			args := call.Call.Args[2:]
			var argNodes []*ssa_graph.SSANode
			for _, arg := range args {
				argNodes = append(argNodes, graph.GetNodeByName(arg.Name()))
			}
			graph.AddServiceCall(node, argNodes, service, method)
			for _, arg := range args {
				fmt.Printf("[RUN] checking taint for service call with arg: %s\n", arg.String())
				node := graph.GetNodeByName(arg.Name())
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
			backwardsAnalysis(graph, originNode.GetValue(), taintInfo, visited, TAINT_CHECK_UPPER, checkTaintInfo)
			node = graph.GetNodeByName(originNode.GetValue().Name())

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

func recurseEdgesBackwardsUntilLoadFrom(graph *ssa_graph.SSAGraph, node *ssa_graph.SSANode, storeEdges []*ssa_graph.SSAEdge, visited map[*ssa_graph.SSANode]bool) []*ssa_graph.SSAEdge {
	if _, ok := visited[node]; ok {
		return storeEdges
	}
	visited[node] = true

	for _, edge := range graph.GetEdgesToNode(node) {
		if edge.GetType() == ssa_graph.EDGE_LOAD {
			storeEdges = append(storeEdges, edge)
		} else {
			storeEdges = append(storeEdges, recurseEdgesBackwardsUntilLoadFrom(graph, edge.GetFromNode(), storeEdges, visited)...)
		}
	}
	return storeEdges
}

func recurseEdgesForwardUntilStoreTo(graph *ssa_graph.SSAGraph, node *ssa_graph.SSANode, storeEdges []*ssa_graph.SSAEdge, visited map[*ssa_graph.SSANode]bool) []*ssa_graph.SSAEdge {
	if _, ok := visited[node]; ok {
		return storeEdges
	}
	visited[node] = true

	for _, edge := range graph.GetEdgesFromNode(node) {
		if edge.GetType() == ssa_graph.EDGE_STORE {
			storeEdges = append(storeEdges, edge)
		} else if edge.GetType() == ssa_graph.EDGE_FIELD || edge.GetType() == ssa_graph.EDGE_INDEX {
			storeEdges = append(storeEdges, recurseEdgesForwardUntilStoreTo(graph, edge.GetToNode(), storeEdges, visited)...)
		}
	}
	return storeEdges
}

func isDatabaseCall(graph *ssa_graph.SSAGraph, instr ssa.Instruction) (string, string, string, []ssa.Value, bool) {
	if call, ok := instr.(*ssa.Call); ok {
		
		// ------------
		// example apps
		// ------------
		if fn, ok := call.Call.Value.(*ssa.Function); ok && len(fn.Params) > 0 {
			fmt.Printf("[TAINT] [1] found call: %v\n", call)
			maybeRcv := fn.Params[0]
			if maybeRcv.Type().String() == "*main.MongoDB" && fn.Name() == "Insert" || fn.Name() == "Find" {
				// return arg without receiver and context
				return "mydb", "mycollection", call.Call.Method.Id(), call.Call.Args[2:], true
			}
			if maybeRcv.Type().String() == "*main.RabbitMQ" && fn.Name() == "Push" {
				// return arg without receiver and context
				return "mydb", "mycollection", call.Call.Method.Id(), call.Call.Args[2:], true
			}
		}

		// --------------
		// blueprint apps
		// --------------
		if unOp, ok := call.Call.Value.(*ssa.UnOp); ok {
			if unOp.Type().String() == "github.com/blueprint-uservices/blueprint/runtime/core/backend.Queue" {
				if slices.Contains([]string{"Push", "Pop"}, call.Call.Method.Name()) {
					fmt.Printf("[TAINT] [2] found %s() call: %v\n", call.Call.Method.Name(), call.Call.Method)
					if fieldAddr, ok := unOp.X.(*ssa.FieldAddr); ok {
						if ptr, ok := fieldAddr.X.Type().(*types.Pointer); ok {
							if _, ok := ptr.Elem().(*types.Named); ok {
								//service, _ := strings.CutSuffix(strings.ToLower(named.Obj().Id()), "serviceimpl")
								//queue := service + "_queue"
								queue := "queue"
								//topic := service + "_message"
								topic := "notification"
								// return arg without context
								return queue, topic, call.Call.Method.Id(), call.Call.Args[1:], true
							}
						}
					}
				}
			}
			if unOp.Type().String() == "github.com/blueprint-uservices/blueprint/runtime/core/backend.NoSQLDatabase" {
				// call for nosqldatabase.GetCollection(...)
				// skip for now
				return "", "", "", nil, false
			}
		}
		if extr, ok := call.Call.Value.(*ssa.Extract); ok {
			if slices.Contains([]string{"InsertOne", "FindOne"}, call.Call.Method.Name()) {
				fmt.Printf("[TAINT] [3] found %s() call: %v\n", call.Call.Method.Name(), call.Call.Method)
				getCollectionNodeCall := graph.GetNodeByName(extr.Tuple.Name())
				if colCal, ok := getCollectionNodeCall.GetInstruction().(*ssa.Call); ok {
					if _, ok := colCal.Call.Value.(*ssa.UnOp); ok {
						dbVal := colCal.Call.Args[1]
						colVal := colCal.Call.Args[2]
						var database, collection string
						if c, ok := dbVal.(*ssa.Const); ok {
							database = strings.Trim(c.Value.ExactString(), "\"")
						}
						if c, ok := colVal.(*ssa.Const); ok {
							collection = strings.Trim(c.Value.ExactString(), "\"")
						}
						// return arg without context
						return database, collection, call.Call.Method.Id(), call.Call.Args[1:], true
					}
				}
			}
		}
	}
	return "", "", "", nil, false
}

func isServiceCall(graph *ssa_graph.SSAGraph, instr ssa.Instruction) (*ssa.Call, string, string, bool) {
	if call, ok := instr.(*ssa.Call); ok {
		// ------------
		// example apps
		// ------------
		if fn, ok := call.Call.Value.(*ssa.Function); ok && len(fn.Params) > 0 {
			maybeRcv := fn.Params[0]
			if maybeRcv.Type().String() == "*main.ShippingService" && fn.Name() == "NewShipment" {
				return call, "ShippingService", "NewShipment", true
			}
			if maybeRcv.Type().String() == "*main.SkuService" && fn.Name() == "GetSku" {
				return call, "SkuService", "GetSku", true
			}
			if maybeRcv.Type().String() == "*main.AnalyticsService" && fn.Name() == "UpdateAnalytics" {
				return call, "AnalyticsService", "UpdateAnalytics", true
			}
			if slices.Contains([]string{
				"StorePost", "ReadPost", "DeletePost", // storage
				"ReadAnalytics",                                     // analytics
				"UploadPost", "DeletePost", "ReadPostWithAnalytics", // upload
			}, fn.Name()) {
				log.Fatal("EXIT!")
				return call, "", "", true
			}
		}

		// --------------
		// blueprint apps
		// --------------
		if unOp, ok := call.Call.Value.(*ssa.UnOp); ok {
			if unOp.Type().String() == "github.com/blueprint-uservices/blueprint/examples/postnotification_simple/workflow/postnotification_simple.UploadService" ||
				unOp.Type().String() == "github.com/blueprint-uservices/blueprint/examples/postnotification_simple/workflow/postnotification_simple.StorageService" ||
				unOp.Type().String() == "github.com/blueprint-uservices/blueprint/examples/postnotification_simple/workflow/postnotification_simple.NotifyService" {

					service := unOp.Type().String()
					var found bool
					service, found = strings.CutPrefix(service, "github.com/blueprint-uservices/blueprint/examples/postnotification_simple/workflow/postnotification_simple.")
					
					if !found {
						log.Fatalf("could not find prefix for service (%s)", service)
					}

					method := call.Call.Method.Id()

					return call, service, method, true
			}
		}
	}
	return nil, "", "", false
}
