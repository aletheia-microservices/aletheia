package tainter

import (
	"fmt"
	"go/token"
	"log"
	"strings"

	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/ssagraph"
	"analyzer/pkg/utils"
)

// REMINDER:
// sometimes there are can be taints such as: _obj[*][*], usertimeline_db.usertimeline.Posts[*].PostID
// this may happen when any[] type is using, for example, fmt.Print calls
// e.g., dsb_sn2 at UserTimelineService.ReadUserTimeline:
// > post_ids = append(new_post_ids, post_ids...)
// > fmt.Println(post_ids)
// ----------------------------
// t56: append(t55, t31...)
// t59: make any <- []int64 (t56)
// *t58 = t59
// t58: &t57[0:int]
// t57: new [1]any (vargs)
// t60: slice t57[:]
// t61: fmt.Println(t60...)
// ----------------------------
// t56 taint is:		 _obj[*], usertimeline_db.usertimeline.Posts[*].PostID
// t57 taint becomes: _obj[*][*], usertimeline_db.usertimeline.Posts[*].PostID
// ----------------------------
func doTaintNode(node *ssagraph.SSANode, taintInfo TaintInfo, taintMode TaintMode) {
	// sanity check for dsb_hotel2 app
	if (strings.Contains(taintInfo.path, ".HId[*]")) && taintInfo.dbTaint.call.GetDatabaseName() == "recommendation_db" {
		log.Fatalf("[TAINT] [DSB_HOTEL2] unexpected taint info path (%s): %v\n", taintInfo.path, taintInfo)
	}

	if taintInfo.isTypeDatabase() && taintInfo.getDatabaseCall() == nil {
		// FIXME: verify this
		fmt.Printf("[TAINT] [4] nil db call for taint info: %v\n", taintInfo)
		return
	}
	if taintInfo.isTypeService() && taintInfo.getServiceCall() == nil {
		// FIXME: verify this
		fmt.Printf("[TAINT] [4] nil sv call for taint info: %v\n", taintInfo)
		return
	}

	var ok bool
	switch taintMode {
	case TAINT_MODE_NEARBY:
		if taintInfo.isTypeDatabase() {
			fmt.Printf("[TAINT NEARBY] [DATABASE] tainting node (%s) for objpath (%s) and dbfield (%s)\n", node.String(), taintInfo.getObjectFullPath(), taintInfo.getDatabasePath())
			ok = node.AddDatabaseTaintIfNotExists(taintInfo.getObjectFullPath(), taintInfo.getDatabasePath(), taintInfo.getDatabaseCall())
		} else if taintInfo.isTypeService() {
			fmt.Printf("[TAINT NEARBY] [SERVICE] tainting node (%s) for objpath (%s) and dbfield (%s)\n", node.String(), taintInfo.getObjectFullPath(), taintInfo.getDatabasePath())
			ok = node.AddServiceTaintIfNotExists(taintInfo.getObjectFullPath(), taintInfo.getServicePath(), taintInfo.getServiceCall())
		}
	case TAINT_MODE_FETCH_UPWARDS:
		if taintInfo.isTypeDatabase() {
			fmt.Printf("[TAINT FETCH] tainting node (%s) for objpath (%s) and dbfield (%s)\n", node.String(), taintInfo.getObjectFullPath(), taintInfo.getDatabasePath()+taintInfo.getObjectPath())
			ok = node.AddDatabaseTaintIfNotExists(taintInfo.getObjectFullPath(), taintInfo.getDatabasePath()+taintInfo.getObjectPath(), taintInfo.getDatabaseCall())
		} else if taintInfo.isTypeService() {
			fmt.Printf("[TAINT FETCH] [SERVICE] tainting node (%s) for objpath (%s) and dbfield (%s)\n", node.String(), taintInfo.getObjectFullPath(), taintInfo.getDatabasePath())
			ok = node.AddServiceTaintIfNotExists(taintInfo.getObjectFullPath(), taintInfo.getServicePath(), taintInfo.getServiceCall())
		}
	}
	if ok {
		fmt.Printf("\t[TAINT] OK!\n")
	}
}

func getObjectPathDiff(longPath1 string, shortPath2 string) string {
	longPath1 = strings.TrimPrefix(longPath1, "_obj")
	shortPath2 = strings.TrimPrefix(shortPath2, "_obj")
	// i.e., pathTop - pathBottomRel
	return strings.TrimPrefix(longPath1, shortPath2)
}

func propagateTaintNearby(graph *ssagraph.SSAGraph, val ssa.Value, taintInfo TaintInfo, visited map[ssa.Value]bool, checkTaintInfo *CheckTaintInfo, upwards bool) {
	taintInfo = taintInfo.updateValue(val)

	fmt.Printf("[TAINT NEARBY] visiting %s: %s // TAINT INFO (_obj%s, %s)\n", val.Name(), val.String(), taintInfo.getObjectPath(), taintInfo.getDatabasePath())
	if visited[val] {
		fmt.Printf("\t[TAINT NEARBY] skipping value %s: %s\n", val.Name(), val.String())
		return
	}
	visited[val] = true

	if ssaValueIsUsedInMongoBsonFilter(graph, val) {
		return // skip
	}

	node := graph.GetNodeByName(val.Name())
	doTaintNode(node, taintInfo, TAINT_MODE_NEARBY)

	fmt.Printf("[TAINT NEARBY] [ROOT=%t] current node: %v\n", taintInfo.objroot, node)
	for _, edge := range graph.GetEdgesFromNode(node) {
		toNode := edge.GetToNode()
		fmt.Printf("\t[TAINT NEARBY] edge (%s) to node: %v\n", edge.GetTypeString(), toNode)

		switch edge.GetType() {
		case ssagraph.EDGE_FIELD:
			if upwards {
				break
			}
			taintInfoTmp := taintInfo.updateCallPathSuffix("." + edge.GetParam())
			propagateTaintNearby(graph, toNode.GetValue(), taintInfoTmp, visited, checkTaintInfo, upwards)
		case ssagraph.EDGE_INDEX:
			if upwards {
				break
			}
			var taintInfoTmp TaintInfo
			if taintInfo.isObjectRoot() {
				taintInfoTmp = taintInfo.updateCallPathSuffix("[" + edge.GetParam() + "]")
			} else {
				taintInfoTmp = taintInfo.enableObjectRoot()
				taintInfoTmp.path, _ = strings.CutSuffix(taintInfoTmp.path, "[*]")
			}

			propagateTaintNearby(graph, toNode.GetValue(), taintInfoTmp, visited, checkTaintInfo, upwards)
		case ssagraph.EDGE_STORE_ADDRESS:
			val := toNode.GetInstruction().(*ssa.Store).Val
			valNode := graph.GetNodeByName(val.Name())
			propagateTaintNearby(graph, valNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
		case ssagraph.EDGE_STORE_VALUE:
			addr := toNode.GetInstruction().(*ssa.Store).Addr
			addrNode := graph.GetNodeByName(addr.Name())
			propagateTaintNearby(graph, addrNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
		case ssagraph.EDGE_POINTS_TO:
			// ignore for now
		case ssagraph.EDGE_ARG_ON_CALL:
			if ok, builtin := utils.SSAValueIsBuiltinFuncCall(toNode.GetValue()); ok {
				if ok, taintFunc, _ := utils.SSABuiltinFuncIsDirect(builtin); ok {
					/* if funcName == "append" && toNode.GetValue().Name() == "t80" {
						fmt.Printf("NODE: %s\n", node.String())
						fmt.Printf("CALL: %s\n", toNode.GetValue().String())
						for _, edge := range graph.GetEdgesToNode(toNode) {
							funcArg := edge.GetFromNode()
							fmt.Printf("ARG: %v\n", funcArg.String())
							if funcArg == node {
								fmt.Printf("skipping...\n")
							}
						}
						log.Fatalf("HERE!")
					} */

					if taintFunc {
						// REMINDER:
						// builtin append() can safely taint its arguments because
						// the Go SSA transforms the second argument into a slice
						// e.g. in dsb_sn2 at UserTimelineService.ReadUserTimeline:
						// > new_post_ids = append(new_post_ids, post.PostID)
						// ----------------------------
						// t107: local PostInfo (post)
						// t122: &t107.PostID [#0]
						// t123: *t122
						// t124: new[1]int64 (varargs)
						// t125: &t124[0:int]
						// *t125 = t123
						// t126: slice t124[:]
						// t127: append(t111, t126...)
						// ----------------------------
						propagateTaintNearby(graph, toNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
					} else {
						for _, edge := range graph.GetEdgesToNode(toNode) {
							funcArg := edge.GetFromNode()
							if funcArg != node {
								propagateTaintNearby(graph, funcArg.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
							}
						}
					}

				}
			}
			// skip
		case ssagraph.EDGE_RETURN_ON, ssagraph.EDGE_EXTRACT:
			// skip
		case ssagraph.EDGE_BINOP_X:
			binOp := edge.GetToNode().GetValue().(*ssa.BinOp)
			if binOp.Op >= token.EQL && binOp.Op <= token.GTR || binOp.Op >= token.NEQ && binOp.Op <= token.GEQ {
				// skip (if conditions)
			} else {
				propagateTaintNearby(graph, toNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
			}
		case ssagraph.EDGE_BINOP_Y:
			binOp := edge.GetToNode().GetValue().(*ssa.BinOp)
			if binOp.Op >= token.EQL && binOp.Op <= token.GTR || binOp.Op >= token.NEQ && binOp.Op <= token.GEQ {
				// skip (if conditions)
			} else {
				propagateTaintNearby(graph, toNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
			}
		case ssagraph.EDGE_MAP_TARGET, ssagraph.EDGE_MAP_KEY, ssagraph.EDGE_MAP_VALUE:
			// [TO BE IMPROVED]
			// skip because toNode is instr and not value
		default:
			propagateTaintNearby(graph, toNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
		}
	}

	fmt.Printf("[TAINT NEARBY] [ROOT=%t] current node: %v\n", taintInfo.objroot, node)
	for _, edge := range graph.GetEdgesToNode(node) {
		fromNode := edge.GetFromNode()
		fmt.Printf("\t[TAINT NEARBY] edge (%s) from node: %v\n", edge.GetTypeString(), fromNode)

		if ssaValueIsUsedInMongoBsonFilter(graph, fromNode.GetValue()) {
			continue // skip
		}

		switch edge.GetType() {
		case ssagraph.EDGE_FIELD:
			visitedTmp := make(map[ssa.Value]bool)
			taintInfoTmp := taintInfo.updateObjectPathPrefix("." + edge.GetParam())
			propagateTaintNearby(graph, fromNode.GetValue(), taintInfoTmp, visitedTmp, checkTaintInfo, true)
		case ssagraph.EDGE_INDEX:
			visitedTmp := make(map[ssa.Value]bool)

			var taintInfoTmp TaintInfo
			if taintInfo.isObjectRoot() {
				taintInfoTmp = taintInfo.updateObjectPathPrefix("[" + edge.GetParam() + "]")
			} else {
				taintInfoTmp = taintInfo.enableObjectRoot()
			}

			propagateTaintNearby(graph, fromNode.GetValue(), taintInfoTmp, visitedTmp, checkTaintInfo, true)
		case ssagraph.EDGE_STORE_ADDRESS:
			val := fromNode.GetInstruction().(*ssa.Store).Val
			valNode := graph.GetNodeByName(val.Name())
			propagateTaintNearby(graph, valNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
		case ssagraph.EDGE_STORE_VALUE:
			addr := fromNode.GetInstruction().(*ssa.Store).Addr
			addrNode := graph.GetNodeByName(addr.Name())
			propagateTaintNearby(graph, addrNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
		case ssagraph.EDGE_POINTS_TO:
			// ignore for now
		case ssagraph.EDGE_RETURN_ON, ssagraph.EDGE_ARG_ON_CALL, ssagraph.EDGE_EXTRACT:
			// skip
		case ssagraph.EDGE_BINOP_X:
			binOp := edge.GetToNode().GetValue().(*ssa.BinOp)
			if binOp.Op >= token.EQL && binOp.Op <= token.GTR || binOp.Op >= token.NEQ && binOp.Op <= token.GEQ {
				// skip (if conditions)
			} else {
				propagateTaintNearby(graph, fromNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
			}
		case ssagraph.EDGE_BINOP_Y:
			binOp := edge.GetToNode().GetValue().(*ssa.BinOp)
			if binOp.Op >= token.EQL && binOp.Op <= token.GTR || binOp.Op >= token.NEQ && binOp.Op <= token.GEQ {
				// skip (if conditions)
			} else {
				propagateTaintNearby(graph, fromNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
			}
		default:
			propagateTaintNearby(graph, fromNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
		}
	}
	fmt.Printf("\t[TAINT NEARBY] exiting %s: %s\n", val.Name(), val.String())
}

func propagateTaintFetchUpwards(graph *ssagraph.SSAGraph, val ssa.Value, taintInfo TaintInfo, visited map[ssa.Value]bool, checkTaintInfo *CheckTaintInfo, upwards bool) {
	taintInfo = taintInfo.updateValue(val)

	fmt.Printf("[TAINT FETCH] visiting %s: %s // TAINT INFO (%s, %s)\n", val.Name(), val.String(), taintInfo.getObjectPath(), taintInfo.getDatabasePath())
	if visited[val] {
		fmt.Printf("\t[TAINT FETCH] skipping value %s: %s\n", val.Name(), val.String())
		return
	}
	visited[val] = true

	node := graph.GetNodeByName(val.Name())
	fmt.Printf("\t[TAINT FETCH] checking upper taints: %v\n", node.GetTaints())

	if ssaValueIsUsedInMongoBsonFilter(graph, val) {
		return // skip
	}

	// 1. taint "subpaths" for current variable and save to later taint the corresponding "subobjects" that requested the upper taint
	for objPath, taints := range node.GetTaints() {

		fmt.Printf("\t[TAINT FETCH] comparing prefixes:\n\t - tainted obj path:\t %s\n\t - bottom to upper:\t %s\n", objPath, taintInfo.getObjectFullPath())

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
				taintInfoTmp.dbTaint.path = taint.GetDatabasePath()
				taintInfoTmp.dbTaint.call = taint.GetDatabaseCall()
				doTaintNode(node, taintInfoTmp, TAINT_MODE_FETCH_UPWARDS)

				// so that we can later taint the bottom node
				dbFieldIndirect := taintInfoTmp.getDatabasePath() + taintInfo.getObjectPath()
				if taintInfoTmp.getDatabaseCall() == nil {
					// FIXME: verify this
					fmt.Printf("[TAINT FETCH] [4] nil db call for taint info tmp: %v\n", taintInfoTmp)
				} else {
					checkTaintInfo.addToIndirectTaints(dbFieldIndirect, taintInfoTmp.getDatabaseCall())
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
				checkTaintInfo.addToInheritedTaints(pathDiff, taint.GetDatabasePath(), taint.GetDatabaseCall())
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
			propagateTaintFetchUpwards(graph, toVal, taintInfo, visited, checkTaintInfo, upwards)
			// TODO: maybe we also need to do this for:
			// (i) nodes whose pointerto set have the current node
			// (ii) nodes within the pointerto set of the current node
			// (iii) load and store instrs?
		}

	}

	fmt.Printf("\t[TAINT FETCH] exiting %s: %s\n", val.Name(), val.String())
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
	params := graph.GetFuncParametersExceptMemberAndContext()
	checkUpperTaintsForObjects(graph, params)
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
				obj := valFieldPath.val

				visited := make(map[ssa.Value]bool)

				taintInfo := NewTaintInfoDatabase(dbfield, "", nil, dbCall)

				// e.g. dsb_hotel2: RecommendationService.LoadRecommendations()
				// where &hotels is the destination
				// -----------------------------------------------
				// var hotels []Hotel
				// cursor, err := collection.FindMany(ctx, filter)
				// cursor.All(ctx, &hotels)
				// -----------------------------------------------
				// REMINDER:
				// cannot do default propagateTaintNearby (root=true) for read operations fetched
				// into arrays/slices (e.g., FindMany) because the default taintinfo is the root path
				// when propagating to other objects
				if valFieldPath.bsonCursorMany || valFieldPath.bsonFilterIn {
					taintInfo.path += "[*]"
					taintInfo.objroot = false
				}

				propagateTaintNearby(graph, obj, taintInfo, visited, nil, false)
				objNode := graph.GetNodeByName(obj.Name())
				nodesToVisit = append(nodesToVisit, objNode)
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
		if service, method, funcShortPath, args, call, ok := isServiceCall(graph, node.GetValue()); ok {
			// keep track of objects passed as arguments
			var argNodes []*ssagraph.SSANode
			fmt.Printf("[TAINT] added service call (%s) --> (%s)\n", graph.GetFunctionShortPath(), funcShortPath)
			for _, arg := range args {
				argNodes = append(argNodes, graph.GetNodeByName(arg.Name()))
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

			for _, argNode := range argNodes {
				arg := argNode.GetValue()
				svpath := svcCall.String2() + "." + arg.Name()
				/* if svpath == "MovieInfoService.WriteMovieInfo.t4" {
					log.Fatalf("HERE")
				} */
				visited := make(map[ssa.Value]bool)
				taintInfo := NewTaintInfoService(svpath, "", nil, svcCall)
				propagateTaintNearby(graph, arg, taintInfo, visited, nil, false)
			}

			for _, retNode := range retNodes {
				ret := retNode.GetValue()
				svpath := svcCall.String2() + "." + ret.Name()
				/* if svpath == "MovieInfoService.WriteMovieInfo.t4" {
					log.Fatalf("HERE")
				} */
				visited := make(map[ssa.Value]bool)
				taintInfo := NewTaintInfoService(svpath, "", nil, svcCall)
				propagateTaintNearby(graph, ret, taintInfo, visited, nil, false)
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
		visited := make(map[ssa.Value]bool)
		taintInfo := NewTaintInfoDatabase("", "", nil, nil)
		checkTaintInfo := NewCheckTaintInfo()
		propagateTaintFetchUpwards(graph, originNode.GetValue(), taintInfo, visited, checkTaintInfo, false)
		node := graph.GetNodeByName(originNode.GetValue().Name())

		// indirect taints
		for _, taint := range checkTaintInfo.indirectTaints {
			if taint.call == nil {
				log.Fatalf("[1] nil db call for taint: %v\n", taint)
			}
			taintInfo := NewTaintInfoDatabase(taint.path, "", originNode.GetValue(), taint.call)
			doTaintNode(node, taintInfo, TAINT_MODE_NEARBY)
		}

		// inherited taints
		for objpath, taints := range checkTaintInfo.inheritedTaints {
			fmt.Printf("[TAINT] check inherited taints for objpath (%s): %v\n", objpath, taints)
			for _, taint := range taints {
				if taint.call == nil {
					// FIXME: verify this
					fmt.Printf("[2] nil db call for taint: %v\n", taint)
				} else {
					taintInfo := NewTaintInfoDatabase(taint.path, objpath, originNode.GetValue(), taint.call)
					doTaintNode(node, taintInfo, TAINT_MODE_NEARBY)
				}
			}
		}
	}
}
