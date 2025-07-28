package tainter

import (
	"fmt"
	"go/types"
	"log"
	"slices"
	"strings"

	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/ssagraph"
	"analyzer/pkg/utils"
)

func doTaintNode(node *ssagraph.SSANode, taintInfo TaintInfo, taintMode TaintMode) {
	switch taintMode {
	case TAINT_BACKWARDS_MARK_AND_PROPAGATE:
		// note that objfields/dbfields already have "." before them
		fmt.Printf("[TAINT] [1] tainting node (%s) for objpath (%s) and dbfield (%s)\n", node.String(), taintInfo.getObjectFullPath(), taintInfo.getDatabaseField())
		ok := node.AddTaintIfNotExists(taintInfo.getObjectFullPath(), taintInfo.getDatabaseField(), taintInfo.getDbCall())
		if ok {
			fmt.Printf("\t[TAINT] OK!\n")
		}
	case TAINT_BACKWARDS_UPDATE_SUBPATHS_AND_FETCH:
		fmt.Printf("[TAINT] [2] tainting node (%s) for objpath (%s) and dbfield (%s)\n", node.String(), taintInfo.getObjectFullPath(), taintInfo.getDatabaseField()+taintInfo.getObjectPath())
		ok := node.AddTaintIfNotExists(taintInfo.getObjectFullPath(), taintInfo.getDatabaseField()+taintInfo.getObjectPath(), taintInfo.getDbCall())
		if ok {
			fmt.Printf("\t[TAINT] OK!\n")
		}
	}
}

func doTaintPointerToSets(graph *ssagraph.SSAGraph, val ssa.Value, taintInfo TaintInfo, visited map[TaintInfo]bool) {
	fmt.Printf("[TAINT|POINTERS] visiting %s: %s // TAINT INFO = (%s, %s)\n", val.Name(), val.String(), taintInfo.getPath(), taintInfo.getDatabaseField())
	node := graph.GetNodeByName(val.Name())
	for _, edge := range graph.GetEdgesFromNode(node) {
		if edge.GetType() == ssagraph.EDGE_POINTS_TO {
			if edge.GetPath() != "" {
				// add before
				// note that both edge.path and objfields/dbfields already have "." before them
				taintInfo = taintInfo.updatePathPrefix(edge.GetPath())
			}
			fmt.Printf("\t[TAINT|POINTERS] calling doTaintNode for pointed at: %s\n", edge.GetToNode().GetName())
			doTaintNode(edge.GetToNode(), taintInfo, TAINT_BACKWARDS_MARK_AND_PROPAGATE)

			backwardsAnalysis(graph, edge.GetToNode().GetValue(), taintInfo, visited, TAINT_BACKWARDS_MARK_AND_PROPAGATE, nil)
		}
	}
	fmt.Printf("\t[TAINT|POINTERS] exiting %s: %s\n", val.Name(), val.String())
}

func getObjectPathDiff(longPath1 string, shortPath2 string) string {
	longPath1 = strings.TrimPrefix(longPath1, "_obj")
	shortPath2 = strings.TrimPrefix(shortPath2, "_obj")
	// i.e., pathTop - pathBottomRel
	return strings.TrimPrefix(longPath1, shortPath2)
}

func backwardsAnalysis(graph *ssagraph.SSAGraph, val ssa.Value, taintInfo TaintInfo, visited map[TaintInfo]bool, taintMode TaintMode, checkTaintInfo *CheckTaintInfo) {
	taintInfo = taintInfo.updateValue(val)

	fmt.Printf("[TAINT|BACKWARD] visiting %s: %s // TAINT INFO (%s, %s)\n", val.Name(), val.String(), taintInfo.getPath(), taintInfo.getDatabaseField())
	if visited[taintInfo] {
		fmt.Printf("\t[TAINT|BACKWARD] skipping value %s: %s\n", val.Name(), val.String())
		return
	}
	visited[taintInfo] = true

	node := graph.GetNodeByName(val.Name())

	switch taintMode {
	case TAINT_BACKWARDS_MARK_AND_PROPAGATE:
		doTaintNode(node, taintInfo, taintMode)
	case TAINT_BACKWARDS_UPDATE_SUBPATHS_AND_FETCH:
		fmt.Printf("\t[TAINT|BACKWARD] checking upper taints: %v\n", node.GetTaints())
		// 1. taint "subpaths" for current variable and save to later taint the corresponding "subobjects" that requested the upper taint
		for objPath, taints := range node.GetTaints() {

			fmt.Printf("\t[TAINT|BACKWARD] comparing prefixes:\n\t - tainted obj path:\t %s\n\t - bottom to upper:\t %s\n", objPath, taintInfo.getObjectFullPath())

			if strings.HasPrefix(taintInfo.getObjectFullPath(), objPath) && taintInfo.getObjectFullPath() != objPath {
				// e.graph.,
				// existing path: 	_obj
				// current path: 	_obj.Shipping
				//
				// in this case, '_obj.Shipping' has prefix '_obj'
				// as result, we may get:
				//
				// existing taint: 	_obj			@ order_db.order
				// potential taint: _obj.Shipping 	@ order_db.order.Shipping
				for _, taint := range taints {
					// save the taint in the upper node
					taintInfoTmp := taintInfo
					taintInfoTmp.dbTaint.dbfield = taint.GetDbField()
					taintInfoTmp.dbTaint.dbcall = taint.GetDbCall()
					doTaintNode(node, taintInfoTmp, taintMode)

					// so that we can later taint the bottom node
					dbFieldIndirect := taintInfoTmp.getDatabaseField() + taintInfo.getObjectPath()
					if taintInfoTmp.getDbCall() == nil {
						log.Fatalf("[4] nil db call for taint info tmp: %v\n", taintInfoTmp)
					}
					checkTaintInfo.addToIndirectTaints(dbFieldIndirect, taintInfoTmp.getDbCall())
				}
				break
			} else if strings.HasPrefix(objPath, taintInfo.getObjectFullPath()) { // also true if strings are equal
				// e.graph.,
				// upper's taint: 		_obj.PostID @ posts_db.post.PostID
				// bottom's path: 		_obj.PostID
				// => bottom's taint: 	_obj		@ posts_db.post.PostID

				pathDiff := getObjectPathDiff(objPath, taintInfo.getObjectFullPath())
				for _, taint := range taints {
					checkTaintInfo.addToInheritedTaints(pathDiff, taint.GetDbField(), taint.GetDbCall())
				}
			}
		}
		// 2. taint forward propagation
		// propagate/update the taints for/of objects that use the current one
		// e.g., in postnotification.StorageService.ReadPost():
		// after "t8: new primitive.E (slicelit)" taint is updated with a subpath for the postid value of the bson filter
		// then we need to go forward again and propagate the taints for "t13: slice t8[...]"
		for _, edge := range graph.GetEdgesFromNode(node) {
			toVal := edge.GetToNode().GetValue()
			switch toVal.(type) {
				case *ssa.MakeInterface, *ssa.Slice:
					backwardsAnalysis(graph, toVal, taintInfo, visited, taintMode, checkTaintInfo)
				// TODO: maybe we also need to do this for:
				// (i) nodes whose pointerto set have the current node
				// (ii) nodes within the pointerto set of the current node
				// (iii) load and store instrs?
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
				if taintMode == TAINT_BACKWARDS_MARK_AND_PROPAGATE {
					doTaintNode(edge.GetToNode(), taintInfo, taintMode)
					for _, edge2 := range graph.GetEdgesToNode(edge.GetToNode()) {
						backwardsAnalysis(graph, edge2.GetFromNode().GetValue(), taintInfo, visited, taintMode, checkTaintInfo)
					}
				}
			}
		}
	case *ssa.FieldAddr:
		fieldName := utils.FieldIndexToName(t)
		fmt.Printf("\t[TAINT|BACKWARD] field addr %s, tainting %s\n", fieldName, t.X.String())
		// add after
		taintInfoTmp := taintInfo
		taintInfoTmp = taintInfoTmp.updatePathPrefix("." + fieldName)
		backwardsAnalysis(graph, t.X, taintInfoTmp, visited, taintMode, checkTaintInfo)
	case *ssa.IndexAddr:
		// add after
		fmt.Printf("\t[TAINT|BACKWARD] index addr %s, tainting %s\n", t.Index.String(), t.X.String())
		taintInfoTmp := taintInfo
		taintInfoTmp = taintInfoTmp.updatePathPrefix("[*]")
		backwardsAnalysis(graph, t.X, taintInfoTmp, visited, taintMode, checkTaintInfo)
	case *ssa.Slice:
		fmt.Printf("\t[TAINT|BACKWARD] slice of: %s\n", t.X.Name())
		// usually t.X is already contained in the set of pointers of the current one
		// note that objects in the pointer set are already tainted in the beginning of this function
		backwardsAnalysis(graph, t.X, taintInfo, visited, taintMode, checkTaintInfo)
	case *ssa.Alloc:
		/* fmt.Printf("\t[TAINT|BACKWARD] alloc used by: %s\n")
		// usually t.X is already contained in the set of pointers of the current one
		// note that objects in the pointer set are already tainted in the beginning of this function
		switch taintMode {
		case TAINT_BACKWARDS_UPDATE_SUBPATHS_AND_FETCH:
			backwardsAnalysis(graph, t.X, taintInfo, visited, taintMode, checkTaintInfo)
		} */
	default:
		fmt.Printf("\t[TAINT|BACKWARD] ignoring value: [%T] %v\n", val, val)
	}

	if taintMode == TAINT_BACKWARDS_MARK_AND_PROPAGATE {
		// if its fieldaddr then we use the objfield and dbfield
		// from the parameters and not the updated ones
		doTaintPointerToSets(graph, val, taintInfo, visited)
	}

	fmt.Printf("\t[TAINT|BACKWARD] exiting %s: %s\n", val.Name(), val.String())
}

func parseArgumentsForMongoDBFilter(graph *ssagraph.SSAGraph, bsonFilter ssa.Value) ([]ssa.Value, []string) {
	var args []ssa.Value
	var keys []string
	bsonFilterNode := graph.GetNodeByName(bsonFilter.Name())
	bsonFilterAllocNode := graph.GetEdgesToNodeExceptPointerTo(bsonFilterNode)[0].GetFromNode()
	elemNode := graph.GetEdgesFromNodeExceptPointerTo(bsonFilterAllocNode)[0].GetToNode()
	bsonFilterKeyNode := graph.GetEdgesFromNode(elemNode)[0].GetToNode()
	// only 1 expected
	edge := recurseEdgesForwardUntilStoreAddress(graph, bsonFilterKeyNode, nil, make(map[*ssagraph.SSANode]bool))[0]
	key := edge.GetToNode().GetInstruction().(*ssa.Store).Val.(*ssa.Const).Value.ExactString()
	keys = append(keys, "."+key)
	arg := graph.GetEdgesFromNode(elemNode)[1].GetToNode().GetValue()
	args = append(args, arg)
	return args, keys
}

func RunTainter(graph *ssagraph.SSAGraph) {
	var nodes []*ssagraph.SSANode
	for _, node := range graph.GetNodes() {
		var foundDatabaseCall bool
		if database, collectionOrTopic, method, args, ok := isDatabaseCall(graph, node.GetInstruction()); ok {
			foundDatabaseCall = true
			/* if node.String() == "t14: invoke t4.FindOne(ctx, t13, nil:[]go.mongodb.org/mongo-driver/bson/primitive.D...)" {
				log.Fatal("EXIT!")
			} */
			if node.String() == "nil:[]go.mongodb.org/mongo-driver/bson/primitive.D: nil:[]go.mongodb.org/mongo-driver/bson/primitive.D" {
				//FIXME (this variable is nil because it is not passed in the call and is optional but for some reason it's assuming it is a db call)
				continue
			}

			var argNodes []*ssagraph.SSANode
			for _, arg := range args {
				argNodes = append(argNodes, graph.GetNodeByName(arg.Name()))
			}

			callId := ssagraph.ComputeCallID(graph, node)
			dbCall := ssagraph.NewDatabaseCall(callId, node, argNodes, database, collectionOrTopic, method)
			graph.AddDatabaseCall(dbCall)

			valDocumentOrMessage := args[0]

			visited := make(map[TaintInfo]bool)
			dbfield := database + "." + collectionOrTopic
			taintInfo := NewTaintInfo(dbfield, "", nil, dbCall)

			backwardsAnalysis(graph, valDocumentOrMessage, taintInfo, visited, TAINT_BACKWARDS_MARK_AND_PROPAGATE, nil)

			node := graph.GetNodeByName(valDocumentOrMessage.Name())
			nodes = append(nodes, node)
		}

		// check for common taints
		for _, originNode := range nodes {
			fmt.Printf("[TAINT] visiting node (origin): %v\n", originNode.String())
			for _, edge := range recurseEdgesBackwardsUntilLoadFrom(graph, originNode, nil, make(map[*ssagraph.SSANode]bool)) {
				// expecting only one node
				node := edge.GetFromNode()
				fmt.Printf("\t[TAINT] visiting node (load): %v\n", node.String())
				spreadTaintsInStorePoint(graph, node, true)
			}
		}

		// keep track of arguments passed in service RPCs so that we can get their indirect taints
		if service, method, funcShortPath, args, ok := isServiceCall(graph, node.GetInstruction()); ok {
			var argNodes []*ssagraph.SSANode
			for _, arg := range args {
				argNodes = append(argNodes, graph.GetNodeByName(arg.Name()))
			}

			callId := ssagraph.ComputeCallID(graph, node)
			svcCall := ssagraph.NewServiceCall(callId, node, argNodes, service, method, funcShortPath)
			graph.AddServiceCall(svcCall)

			fmt.Printf("[TAINT] added service call (%s) --> (%s)\n", graph.GetFunctionShortPath(), funcShortPath)
			for _, arg := range args {
				fmt.Printf("[TAINT] checking taint for service call with arg: %s\n", arg.String())
				node := graph.GetNodeByName(arg.Name())
				nodes = append(nodes, node)
			}
		}

		// mark the parameters of the current function so that we can get their indirect taints
		// NOTE: currently not adding to nodes array
		if foundDatabaseCall {
			params := graph.GetFuncParametersExceptMemberAndContext()
			for _, param := range params {
				spreadTaintsInStorePoint(graph, param, false)
			}
		}

		// check for upper taints affecting the current database/service calls objects
		for _, originNode := range nodes {
			fmt.Println()
			fmt.Printf("[TAINT] check upper taints for node: %v\n", originNode.String())
			visited := make(map[TaintInfo]bool)
			taintInfo := NewTaintInfo("", "", nil, nil)
			checkTaintInfo := NewCheckTaintInfo()
			backwardsAnalysis(graph, originNode.GetValue(), taintInfo, visited, TAINT_BACKWARDS_UPDATE_SUBPATHS_AND_FETCH, checkTaintInfo)
			node = graph.GetNodeByName(originNode.GetValue().Name())

			// indirect taints
			for _, taint := range checkTaintInfo.indirectTaints {
				if taint.dbcall == nil {
					log.Fatalf("[1] nil db call for taint: %v\n", taint)
				}
				taintInfo := NewTaintInfo(taint.dbfield, "", originNode.GetValue(), taint.dbcall)
				doTaintNode(node, taintInfo, TAINT_BACKWARDS_MARK_AND_PROPAGATE)
			}

			// inherited taints
			for objpath, taints := range checkTaintInfo.inheritedTaints {
				fmt.Printf("[TAINT] check inherited taints for objpath (%s): %v\n", objpath, taints)
				for _, taint := range taints {
					if taint.dbcall == nil {
						log.Fatalf("[2] nil db call for taint: %v\n", taint)
					}
					taintInfo := NewTaintInfo(taint.dbfield, objpath, originNode.GetValue(), taint.dbcall)
					doTaintNode(node, taintInfo, TAINT_BACKWARDS_MARK_AND_PROPAGATE)
				}
			}
		}
	}
}

func spreadTaintsInStorePoint(graph *ssagraph.SSAGraph, node *ssagraph.SSANode, valToAddr bool) {
	var edges []*ssagraph.SSAEdge

	if valToAddr { // addr <<< val
		edges = recurseEdgesForwardUntilStoreAddress(graph, node, nil, make(map[*ssagraph.SSANode]bool))
	} else { // addr >>> val
		edges = recurseEdgesForwardUntilStoreValue(graph, node, nil, make(map[*ssagraph.SSANode]bool))
	}
	for _, edge := range edges {
		// if valToAddr is true, then srcNode is the Value and dstNode is the Address
		// if valToAddr is false, then srcNode is the Address and dstNode is the Value
		var dstNode, storeNode, srcNode *ssagraph.SSANode

		dstNode = edge.GetFromNode()
		storeNode = edge.GetToNode()

		var srcNodes []*ssagraph.SSANode // THIS IS NOT NECESSARY??
		if sr, ok := storeNode.GetInstruction().(*ssa.Store); ok {
			if valToAddr {
				srcNode = graph.GetNodeByName(sr.Val.Name())
			} else {
				srcNode = graph.GetNodeByName(sr.Addr.Name())
			}
			// sanity check
			if !slices.Contains(srcNodes, srcNode) { // THIS IS NOT NECESSARY!
				srcNodes = append(srcNodes, srcNode)
			}
		}

		for _, srcNode := range srcNodes {
			visited := make(map[TaintInfo]bool)
			taintInfo := NewTaintInfo("", "", nil, nil)
			checkTaintInfo := NewCheckTaintInfo()

			// go up to fetch all possible indirect taints for the current node
			backwardsAnalysis(graph, srcNode.GetValue(), taintInfo, visited, TAINT_BACKWARDS_UPDATE_SUBPATHS_AND_FETCH, checkTaintInfo)

			// indirect taints
			for _, taint := range checkTaintInfo.indirectTaints {
				if taint.dbcall == nil {
					log.Fatalf("[3] nil db call for taint: %v\n", taint)
				}
				taintInfo := NewTaintInfo(taint.dbfield, "", srcNode.GetValue(), taint.dbcall)

				// taint current node with all possible indirect taints
				doTaintNode(srcNode, taintInfo, TAINT_BACKWARDS_MARK_AND_PROPAGATE)

				// not needed but helps in visualization ssagraph
				doTaintNode(storeNode, taintInfo, TAINT_BACKWARDS_MARK_AND_PROPAGATE)

				// now "spread" the previous obtained taints to the addrNode
				visited2 := make(map[TaintInfo]bool)
				taintInfo2 := NewTaintInfo(taint.dbfield, "", nil, taint.dbcall)
				backwardsAnalysis(graph, dstNode.GetValue(), taintInfo2, visited2, TAINT_BACKWARDS_MARK_AND_PROPAGATE, nil)
			}
		}
	}
}

func recurseEdgesBackwardsUntilLoadFrom(graph *ssagraph.SSAGraph, node *ssagraph.SSANode, storeEdges []*ssagraph.SSAEdge, visited map[*ssagraph.SSANode]bool) []*ssagraph.SSAEdge {
	if _, ok := visited[node]; ok {
		return storeEdges
	}
	visited[node] = true

	for _, edge := range graph.GetEdgesToNode(node) {
		if edge.GetType() == ssagraph.EDGE_LOAD {
			storeEdges = append(storeEdges, edge)
		} else {
			storeEdges = append(storeEdges, recurseEdgesBackwardsUntilLoadFrom(graph, edge.GetFromNode(), storeEdges, visited)...)
		}
	}
	return storeEdges
}

func recurseEdgesForwardUntilStoreAddress(graph *ssagraph.SSAGraph, node *ssagraph.SSANode, storeEdges []*ssagraph.SSAEdge, visited map[*ssagraph.SSANode]bool) []*ssagraph.SSAEdge {
	if _, ok := visited[node]; ok {
		return storeEdges
	}
	visited[node] = true

	for _, edge := range graph.GetEdgesFromNode(node) {
		if edge.GetType() == ssagraph.EDGE_STORE_ADDRESS {
			storeEdges = append(storeEdges, edge)
		} else if edge.GetType() == ssagraph.EDGE_FIELD || edge.GetType() == ssagraph.EDGE_INDEX || edge.GetType() == ssagraph.EDGE_USAGE {
			storeEdges = append(storeEdges, recurseEdgesForwardUntilStoreAddress(graph, edge.GetToNode(), storeEdges, visited)...)
		}
	}
	return storeEdges
}

func recurseEdgesForwardUntilStoreValue(graph *ssagraph.SSAGraph, node *ssagraph.SSANode, storeEdges []*ssagraph.SSAEdge, visited map[*ssagraph.SSANode]bool) []*ssagraph.SSAEdge {
	if _, ok := visited[node]; ok {
		return storeEdges
	}
	visited[node] = true
	for _, edge := range graph.GetEdgesFromNode(node) {
		if edge.GetType() == ssagraph.EDGE_STORE_VALUE {
			storeEdges = append(storeEdges, edge)
		} else if edge.GetType() == ssagraph.EDGE_FIELD || edge.GetType() == ssagraph.EDGE_INDEX || edge.GetType() == ssagraph.EDGE_USAGE {
			storeEdges = append(storeEdges, recurseEdgesForwardUntilStoreValue(graph, edge.GetToNode(), storeEdges, visited)...)
		}
	}
	return storeEdges
}

func isDatabaseCall(graph *ssagraph.SSAGraph, instr ssa.Instruction) (string, string, string, []ssa.Value, bool) {
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
								// return all args except context
								// NOTE: in this case (when call.Call.Value is UnOp) call.Call.Args does not contain the receiver
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
						// return all args except context
						// NOTE: in this case (when call.Call.Value is UnOp) call.Call.Args does not contain the receiver
						return database, collection, call.Call.Method.Id(), call.Call.Args[1:], true
					}
				}
			}
		}
	}
	return "", "", "", nil, false
}

func isServiceCall(graph *ssagraph.SSAGraph, instr ssa.Instruction) (string, string, string, []ssa.Value, bool) {
	if call, ok := instr.(*ssa.Call); ok {
		// ------------
		// example apps
		// ------------
		if fn, ok := call.Call.Value.(*ssa.Function); ok && len(fn.Params) > 0 {
			maybeRcv := fn.Params[0]
			if maybeRcv.Type().String() == "*main.ShippingService" && fn.Name() == "NewShipment" {
				// return all args except receiver and context
				return "ShippingService", "NewShipment", "", call.Call.Args[2:], true
			}
			if maybeRcv.Type().String() == "*main.SkuService" && fn.Name() == "GetSku" {
				// return all args except receiver and context
				return "SkuService", "GetSku", "", call.Call.Args[2:], true
			}
			if maybeRcv.Type().String() == "*main.AnalyticsService" && fn.Name() == "UpdateAnalytics" {
				// return all args except receiver and context
				return "AnalyticsService", "UpdateAnalytics", "", call.Call.Args[2:], true
			}
			if slices.Contains([]string{
				"StorePost", "ReadPost", "DeletePost", // storage
				"ReadAnalytics",                                     // analytics
				"UploadPost", "DeletePost", "ReadPostWithAnalytics", // upload
			}, fn.Name()) {
				log.Fatal("EXIT!")
				// return all args except receiver and context
				return "", "", "", call.Call.Args[2:], true
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

				// NOTE: unOp.Type().String() does not contain "Impl" suffix here so GetShortFunctionPath will just ignore
				funcShortPath := utils.GetShortFunctionPath(unOp.Type().String() + "." + method)

				// return all args except context
				// NOTE: in this case (when call.Call.Value is UnOp) call.Call.Args does not contain the receiver
				return service, method, funcShortPath, call.Call.Args[1:], true
			}
		}
	}
	return "", "", "", nil, false
}
