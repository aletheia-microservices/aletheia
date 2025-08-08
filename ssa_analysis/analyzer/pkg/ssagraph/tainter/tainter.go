package tainter

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/ssagraph"
	"analyzer/pkg/utils"
)

func doTaintNode(node *ssagraph.SSANode, taintInfo TaintInfo, taintMode TaintMode) {
	if taintInfo.getDbCall() == nil {
		// FIXME: verify this
		fmt.Printf("[4] nil db call for taint info: %v\n", taintInfo)
		return
	}
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

	fmt.Printf("[TAINT - BACKWARD] visiting %s: %s // TAINT INFO (%s, %s)\n", val.Name(), val.String(), taintInfo.getPath(), taintInfo.getDatabaseField())
	if visited[taintInfo] {
		fmt.Printf("\t[TAINT - BACKWARD] skipping value %s: %s\n", val.Name(), val.String())
		return
	}
	visited[taintInfo] = true

	node := graph.GetNodeByName(val.Name())

	switch taintMode {
	case TAINT_BACKWARDS_MARK_AND_PROPAGATE:
		doTaintNode(node, taintInfo, taintMode)
	case TAINT_BACKWARDS_UPDATE_SUBPATHS_AND_FETCH:
		fmt.Printf("\t[TAINT - BACKWARD] checking upper taints: %v\n", node.GetTaints())
		// 1. taint "subpaths" for current variable and save to later taint the corresponding "subobjects" that requested the upper taint
		for objPath, taints := range node.GetTaints() {

			fmt.Printf("\t[TAINT - BACKWARD] comparing prefixes:\n\t - tainted obj path:\t %s\n\t - bottom to upper:\t %s\n", objPath, taintInfo.getObjectFullPath())

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
						// FIXME: verify this
						fmt.Printf("[4] nil db call for taint info tmp: %v\n", taintInfoTmp)
					} else {
						checkTaintInfo.addToIndirectTaints(dbFieldIndirect, taintInfoTmp.getDbCall())
					}
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
		fmt.Printf("\t[TAINT - BACKWARD] field addr %s, tainting %s\n", fieldName, t.X.String())
		// add after
		taintInfoTmp := taintInfo
		taintInfoTmp = taintInfoTmp.updatePathPrefix("." + fieldName)
		backwardsAnalysis(graph, t.X, taintInfoTmp, visited, taintMode, checkTaintInfo)
	case *ssa.IndexAddr:
		// add after
		fmt.Printf("\t[TAINT - BACKWARD] index addr %s, tainting %s\n", t.Index.String(), t.X.String())
		taintInfoTmp := taintInfo
		prefix := "[*]"
		if index, ok := utils.ExtractStringFromValue(t.Index); ok {
			prefix = "[" + index + "]"
		}
		taintInfoTmp = taintInfoTmp.updatePathPrefix(prefix)
		backwardsAnalysis(graph, t.X, taintInfoTmp, visited, taintMode, checkTaintInfo)
	case *ssa.Slice:
		fmt.Printf("\t[TAINT - BACKWARD] slice of: %s\n", t.X.Name())
		// usually t.X is already contained in the set of pointers of the current one
		// note that objects in the pointer set are already tainted in the beginning of this function
		backwardsAnalysis(graph, t.X, taintInfo, visited, taintMode, checkTaintInfo)
	case *ssa.Alloc:
		/* fmt.Printf("\t[TAINT - BACKWARD] alloc used by: %s\n")
		// usually t.X is already contained in the set of pointers of the current one
		// note that objects in the pointer set are already tainted in the beginning of this function
		switch taintMode {
		case TAINT_BACKWARDS_UPDATE_SUBPATHS_AND_FETCH:
			backwardsAnalysis(graph, t.X, taintInfo, visited, taintMode, checkTaintInfo)
		} */
	default:
		fmt.Printf("\t[TAINT - BACKWARD] ignoring value: [%T] %v\n", val, val)
	}

	if taintMode == TAINT_BACKWARDS_MARK_AND_PROPAGATE {
		// if its fieldaddr then we use the objfield and dbfield
		// from the parameters and not the updated ones
		doTaintPointerToSets(graph, val, taintInfo, visited)
	}

	if taintMode == TAINT_BACKWARDS_UPDATE_SUBPATHS_AND_FETCH {
		// TODO: verify this

		node := graph.GetNodeByName(val.Name())
		var allEdges []*ssagraph.SSAEdge
		var allNodes []*ssagraph.SSANode
		numEdgesPointsTo := 0
		for _, edge := range graph.GetEdgesFromNode(node) {
			if edge.GetType() == ssagraph.EDGE_POINTS_TO {
				allEdges = append(allEdges, edge)
				allNodes = append(allNodes, edge.GetToNode())
				numEdgesPointsTo++
			}
		}
		numEdgesPointedBy := 0
		for _, edge := range graph.GetEdgesToNode(node) {
			if edge.GetType() == ssagraph.EDGE_POINTS_TO {
				allEdges = append(allEdges, edge)
				allNodes = append(allNodes, edge.GetFromNode())
				numEdgesPointedBy++
			}
		}

		for i, edge := range allEdges {
			if edge.GetPath() != "" && i < numEdgesPointsTo {
				// add before
				// note that both edge.path and objfields/dbfields already have "." before them
				taintInfo = taintInfo.updatePathPrefix(edge.GetPath())
			} else if edge.GetPath() != "" && i >= numEdgesPointsTo {
				continue
			}
			node := allNodes[i]
			fmt.Printf("\t[TAINT - BACKWARD] calling doTaintNode for pointed to/by: %s\n", node.GetName())
			backwardsAnalysis(graph, node.GetValue(), taintInfo, visited, TAINT_BACKWARDS_UPDATE_SUBPATHS_AND_FETCH, checkTaintInfo)
		}

		// TODO: verify this
		// e.g. digota.SkuService.Delete
		// fetch delete taint (skus_db.skus[*].Value) of t6 (slice) to later propagate to t19 (queue message)
		//
		// problem:
		// it is adding a new incorrect constraint at postnotification
		// FOREIGN_KEY posts_db.post.ReqID REFERENCES notifications_queue.notification.ReqID
		for _, edge := range graph.GetEdgesFromNode(node) {
			if edge.GetType() != ssagraph.EDGE_POINTS_TO {
				if edge.GetPath() == "" {
					if edge.GetToNode().GetValue() != nil {
						if val, ok := edge.GetToNode().GetValue().(*ssa.MakeInterface); ok {
							backwardsAnalysis(graph, val, taintInfo, visited, TAINT_BACKWARDS_UPDATE_SUBPATHS_AND_FETCH, checkTaintInfo)
						}
					} else if edge.GetType() == ssagraph.EDGE_STORE_VALUE {
						if instr, ok := edge.GetToNode().GetInstruction().(*ssa.Store); ok {
							backwardsAnalysis(graph, instr.Addr, taintInfo, visited, TAINT_BACKWARDS_UPDATE_SUBPATHS_AND_FETCH, checkTaintInfo)
						}
					}
				}
			}
		}
	}

	fmt.Printf("\t[TAINT - BACKWARD] exiting %s: %s\n", val.Name(), val.String())
}

func RunTainter(graph *ssagraph.SSAGraph) {
	// for simplicity, always run taint on database calls before service calls due to cases where
	// 1st there is a service call that returns some value
	// 2nd that value is then used in a database call
	//
	// the taint is detected when analyzing the service call
	// only if the database call was analyzed before
	// because of the way we are checking upper taints with the "spreading of taints in store points"
	// at the service call
	//
	// this logic is close to the tainting process of returns and parameters
	// where we only taint them after all database calls have been processed
	runTainterOnDatabaseCalls(graph)
	runTainterOnServiceCalls(graph)
	runTainterOnParameters(graph)
	runTainterOnReturns(graph)
}

func runTainterOnParameters(graph *ssagraph.SSAGraph) {
	// mark the parameters of the current function so that we can get their indirect taints
	// NOTE: currently not adding to nodes array
	params := graph.GetFuncParametersExceptMemberAndContext()
	for _, param := range params {
		spreadTaintsInStorePoint(graph, param, false)
	}
}

func runTainterOnReturns(graph *ssagraph.SSAGraph) {
	var nodesToVisit []*ssagraph.SSANode
	for _, node := range graph.GetNodes() {
		// keep track of arguments returned in the current function possibly invoked by other services
		// so that we can mark their indirect taints
		if ret, ok := node.GetInstruction().(*ssa.Return); ok {
			var rets []*ssagraph.SSANode
			for _, res := range ret.Results {
				resNode := graph.GetNodeByName(res.Name())
				nodesToVisit = append(nodesToVisit, resNode)
				rets = append(rets, resNode)
			}
			graph.AddReturnsToLst(rets)
		}

		checkUpperTaintsForObjects(graph, nodesToVisit)
	}
}

func runTainterOnDatabaseCalls(graph *ssagraph.SSAGraph) {
	for _, node := range graph.GetNodes() {
		var nodesToVisit []*ssagraph.SSANode
		if database, collectionOrTopic, method, opType, valFieldPathLst, ok := isDatabaseCall(graph, node.GetValue()); ok {
			var argNodes []*ssagraph.SSANode

			for _, valFieldPath := range valFieldPathLst {
				argNodes = append(argNodes, graph.GetNodeByName(valFieldPath.val.Name()))
			}

			callId := ssagraph.ComputeCallID(graph, node)
			dbCall := ssagraph.NewDatabaseCall(callId, node, argNodes, database, collectionOrTopic, method, opType)
			graph.AddDatabaseCall(dbCall)

			for _, valFieldPath := range valFieldPathLst {
				dbfield := valFieldPath.fieldpath
				arg := valFieldPath.val

				visited := make(map[TaintInfo]bool)
				taintInfo := NewTaintInfo(dbfield, "", nil, dbCall)

				backwardsAnalysis(graph, arg, taintInfo, visited, TAINT_BACKWARDS_MARK_AND_PROPAGATE, nil)

				node := graph.GetNodeByName(arg.Name())
				nodesToVisit = append(nodesToVisit, node)
			}

			// check for common taints
			for _, originNode := range nodesToVisit {
				fmt.Printf("[TAINT] visiting node (origin): %v\n", originNode.String())
				for _, edge := range recurseEdgesBackwardsUntilLoadFrom(graph, originNode, nil, make(map[*ssagraph.SSANode]bool)) {
					node := edge.GetFromNode()
					fmt.Printf("\t[TAINT] visiting node (load): %v\n", node.String())
					spreadTaintsInStorePoint(graph, node, true)
				}
			}

			checkUpperTaintsForObjects(graph, nodesToVisit)
		}
	}
}

func runTainterOnServiceCalls(graph *ssagraph.SSAGraph) {
	for _, node := range graph.GetNodes() {
		var nodesToVisit []*ssagraph.SSANode
		// keep track of arguments passed in service RPCs
		// so that we can mark their indirect taints
		if service, method, funcShortPath, args, call, ok := isServiceCall(graph, node.GetInstruction()); ok {
			// keep track of objects passed as arguments
			var argNodes []*ssagraph.SSANode
			for _, arg := range args {
				argNodes = append(argNodes, graph.GetNodeByName(arg.Name()))
			}

			fmt.Printf("[TAINT] added service call (%s) --> (%s)\n", graph.GetFunctionShortPath(), funcShortPath)
			for _, arg := range args {
				fmt.Printf("[TAINT] checking taint for service call with arg: %s\n", arg.String())
				node := graph.GetNodeByName(arg.Name())
				nodesToVisit = append(nodesToVisit, node)
			}

			// keep track of objects extracted from returns
			var retNodes []*ssagraph.SSANode
			callNode := graph.GetNodeByName(call.Name())
			if call.Call.Signature().Results().Len() > 1 {
				for _, edge := range graph.GetEdgesFromNode(callNode) {
					if edge.GetType() == ssagraph.EDGE_EXTRACT {
						nodesToVisit = append(nodesToVisit, edge.GetToNode())
						retNodes = append(retNodes, edge.GetToNode())
					}
				}
			} else {
				// when there is only one return value then there
				// are no extract instructions and the value is just
				// the one declared when invoking the function
				nodesToVisit = append(nodesToVisit, callNode)
				retNodes = append(retNodes, callNode)
			}
			fmt.Printf("[TAINT] [%s] got ret nodes: %v\n", node.String(), retNodes)

			callId := ssagraph.ComputeCallID(graph, node)
			svcCall := ssagraph.NewServiceCall(callId, node, argNodes, retNodes, service, method, funcShortPath)
			graph.AddServiceCall(svcCall)

			// check for common taints
			for _, retNode := range retNodes {
				fmt.Printf("[TAINT] spreading taints for ret node: %s\n", retNode)
				spreadTaintsInStorePoint(graph, retNode, false)
			}
			fmt.Printf("[TAINT] visiting nodes for call (%s) --> (%s)\n", graph.GetFunctionShortPath(), funcShortPath)
			checkUpperTaintsForObjects(graph, nodesToVisit)
		}

	}
}

func checkUpperTaintsForObjects(graph *ssagraph.SSAGraph, nodesToVisit []*ssagraph.SSANode) {
	// check for upper taints affecting the current database/service calls objects
	for _, originNode := range nodesToVisit {
		fmt.Println()
		fmt.Printf("[TAINT] check upper taints for node: %v\n", originNode.String())
		visited := make(map[TaintInfo]bool)
		taintInfo := NewTaintInfo("", "", nil, nil)
		checkTaintInfo := NewCheckTaintInfo()
		backwardsAnalysis(graph, originNode.GetValue(), taintInfo, visited, TAINT_BACKWARDS_UPDATE_SUBPATHS_AND_FETCH, checkTaintInfo)
		node := graph.GetNodeByName(originNode.GetValue().Name())

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
					// FIXME: verify this
					fmt.Printf("[2] nil db call for taint: %v\n", taint)
				} else {
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

	// [TO BE IMPROVED] 
	// must declared visited here on top otherwise we have infinite recursion 
	// when recurseEdgesBackwardsUntilLoadFrom() includes edges with type EDGE_FIELD and EDGE_INDEX
	visited := make(map[TaintInfo]bool)
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
				taintInfo2 := NewTaintInfo(taint.dbfield, "", nil, taint.dbcall)
				backwardsAnalysis(graph, dstNode.GetValue(), taintInfo2, visited, TAINT_BACKWARDS_MARK_AND_PROPAGATE, nil)
			}
		}
	}
}

func recurseEdgesBackwardsUntilLoadFrom(graph *ssagraph.SSAGraph, node *ssagraph.SSANode, loadEdges []*ssagraph.SSAEdge, visited map[*ssagraph.SSANode]bool) []*ssagraph.SSAEdge {
	fmt.Printf("[TAINT] recurse edges backward until load from for current node (%s) with load edges list: %v\n", node.String(), loadEdges)
	if visited[node] {
		return loadEdges
	}
	visited[node] = true

	for _, edge := range graph.GetEdgesToNode(node) {
		// [TO BE IMPROVED]
		// include edges with type EDGE_FIELD and EDGE_INDEX
		// to make sure we fetch taints from objects children
		// e.g., digota.OrderService.Run() (for more info check .dot file):
		// where t11 (bson slice) must fetch the read taint from t3 (queue message)
		// by including edges like t14 (&t12.Value where t12: &t11[0])
		// we will later spreadTaintsInStorePoint (*t14 = t17) - in code this means adding the id to the bson filter
		if edge.GetType() == ssagraph.EDGE_LOAD || edge.GetType() == ssagraph.EDGE_FIELD || edge.GetType() == ssagraph.EDGE_INDEX {
			loadEdges = append(loadEdges, edge)
		} else {
			loadEdges = append(loadEdges, recurseEdgesBackwardsUntilLoadFrom(graph, edge.GetFromNode(), loadEdges, visited)...)
		}
	}
	return loadEdges
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
