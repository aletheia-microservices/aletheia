package tainter

import (
	"log"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/ssagraph"
)

var seenTaint map[TaintInfo]bool

// REMINDER:
// sometimes there are can be taints such as: _obj[*][*], usertimeline_db.usertimeline.Posts[*].PostID
// this may happen when any[] type is using, for example, // EVAL: fmt.Print calls
// e.g., dsb_sn2 at UserTimelineService.ReadUserTimeline:
// > post_ids = append(new_post_ids, post_ids...)
// > // EVAL: fmt.Println(post_ids)
// ----------------------------
// t56: append(t55, t31...)
// t59: make any <- []int64 (t56)
// *t58 = t59
// t58: &t57[0:int]
// t57: new [1]any (vargs)
// t60: slice t57[:]
// t61: // EVAL: fmt.Println(t60...)
// ----------------------------
// t56 taint is:		 _obj[*], usertimeline_db.usertimeline.Posts[*].PostID
// t57 taint becomes: _obj[*][*], usertimeline_db.usertimeline.Posts[*].PostID
// ----------------------------
func doTaintNode(node *ssagraph.SSANode, taintInfo TaintInfo, taintMode TaintMode) {
	var ok bool
	logrus.WithField("taint_info", taintInfo.String()).
		WithField("node", node.String()).
		Debugf("taint node")

	switch taintMode {
	case TAINT_MODE_NEARBY:
		if taintInfo.isTypeDatabase() {
			logrus.Tracef("[TAINT NEARBY] [DATABASE] tainting node (%s) for objpath (%s) and dbfield (%s)\n", node.String(), taintInfo.getObjectFullPath(), taintInfo.getDatabasePath())
			ok = node.AddDatabaseTaintIfNotExists(taintInfo.getObjectFullPath(), taintInfo.getDatabasePath(), taintInfo.getDatabaseCall(), taintInfo.isReadKey(), taintInfo.isReadValue(), taintInfo.getCallerT())
		} else if taintInfo.isTypeService() {
			logrus.Tracef("[TAINT NEARBY] [SERVICE] tainting node (%s) for objpath (%s) and dbfield (%s)\n", node.String(), taintInfo.getObjectFullPath(), taintInfo.getDatabasePath())
			ok = node.AddServiceTaintIfNotExists(taintInfo.getObjectFullPath(), taintInfo.getServicePath(), taintInfo.getServiceCall(), taintInfo.getCallerT())
		} else {
			logrus.Fatalf("unexpected type of taint info: %s\n", taintInfo.String())
		}
	case TAINT_MODE_FETCH_UPWARDS:
		if taintInfo.isTypeDatabase() {
			logrus.Tracef("[TAINT FETCH] tainting node (%s) for objpath (%s) and dbfield (%s)\n", node.String(), taintInfo.getObjectFullPath(), taintInfo.getDatabasePath()+taintInfo.getObjectPath())
			ok = node.AddDatabaseTaintIfNotExists(taintInfo.getObjectFullPath(), taintInfo.getDatabasePath()+taintInfo.getObjectPath(), taintInfo.getDatabaseCall(), taintInfo.isReadKey(), taintInfo.isReadValue(), taintInfo.getCallerT())
		} else if taintInfo.isTypeService() {
			logrus.Tracef("[TAINT FETCH] [SERVICE] tainting node (%s) for objpath (%s) and dbfield (%s)\n", node.String(), taintInfo.getObjectFullPath(), taintInfo.getDatabasePath())
			ok = node.AddServiceTaintIfNotExists(taintInfo.getObjectFullPath(), taintInfo.getServicePath(), taintInfo.getServiceCall(), taintInfo.getCallerT())
		} else {
			logrus.Fatalf("unexpected type of taint info: %s\n", taintInfo.String())
		}
	default:
		logrus.Fatalf("unexpected taint mode: %v\n", taintMode)
	}
	if ok {
		logrus.Tracef("\t[TAINT] OK!\n")
	}
}

func getObjectPathDiff(longPath1 string, shortPath2 string) string {
	longPath1 = strings.TrimPrefix(longPath1, "_obj")
	shortPath2 = strings.TrimPrefix(shortPath2, "_obj")
	// i.e., pathTop - pathBottomRel
	return strings.TrimPrefix(longPath1, shortPath2)
}

func nodeHasTaintInfo(node *ssagraph.SSANode, objpath string, info TaintInfo) bool {
	for _, taint := range node.GetTaintsForPath(objpath) {
		if taint.IsDatabaseTaint() && info.isTypeDatabase() && taint.GetDatabaseCall() == info.getDatabaseCall() {
			return true
		} else if taint.IsServiceTaint() && info.isTypeService() && taint.GetServiceCall() == info.getServiceCall() {
			return true
		}
	}
	return false
}

func generateRootTaintInfoFromTaint(node *ssagraph.SSANode, taint *ssagraph.SSATaint) TaintInfo {
	if taint.IsDatabaseTaint() {
		return NewTaintInfoDatabase(taint.GetDatabasePath(), "", node.GetValue(), taint.GetDatabaseCall(), taint.IsReadKey(), taint.IsReadValue())
	} else if taint.IsServiceTaint() {
		return NewTaintInfoService(taint.GetServicePath(), "", node.GetValue(), taint.GetServiceCall())
	}
	log.Panicf("[TAINT INFO FROM TAINT] unexpected type for taint: %s\n", taint.String())
	return TaintInfo{}
}

func propagateTaintNearby(graph *ssagraph.SSAGraph, recurse bool, val ssa.Value, taintInfo TaintInfo, visited map[ssa.Value]bool, checkTaintInfo *CheckTaintInfo, upwards bool) {
	if val == nil {
		log.Panicf("[TAINT NEARBY] unexpected nil val // TAINT INFO (_obj%s, %s)\n", taintInfo.getObjectPath(), taintInfo.getDatabasePath())
	}

	taintInfo = taintInfo.updateValue(val)

	var prevVal *ssa.Value
	var prevValStr string
	if taintInfo.prevval != nil {
		prevValStr = taintInfo.prevval.Name() + ": " + taintInfo.prevval.String()
	}

	if strings.Contains(taintInfo.objpath, ".MapKey.Key.MapKey") {
		logrus.WithField("graph", graph.String()).
			WithField("prev", prevValStr).
			WithField("val", val.Name()).
			WithField("taint_info", taintInfo.String()).
			Fatalf("[TAINT NEARBY] stop!")
	}

	logrus.WithField("graph", graph.String()).
		WithField("prev", prevValStr).
		WithField("val", val.Name()).
		WithField("taint_info", taintInfo.String()).
		Debugf("[TAINT NEARBY] visiting")

	if seenTaint[taintInfo] {
		logrus.WithField("val", val.Name).WithField("taint_info", taintInfo.String()).
			Tracef("[TAINT NEARBY] skipping (seen taint)...")
		return
	}
	seenTaint[taintInfo] = true

	if visited[val] {
		logrus.WithField("val", val.Name).WithField("taint_info", taintInfo.String()).
			Tracef("[TAINT NEARBY] skipping (visited val)...")
		return
	}
	visited[val] = true

	node := graph.GetNodeByName(val.Name())

	// avoid infinite recursion
	if nodeHasTaintInfo(node, "_obj"+taintInfo.objpath, taintInfo) {
		logrus.Debugf("skipping taint (node=%s) (taint_info=%s)\n", node.String(), taintInfo.String())
		return
	}

	doTaintNode(node, taintInfo, TAINT_MODE_NEARBY)

	logrus.Debugf("[TAINT NEARBY] [PART_1] [ROOT=%t] [RECURSE=%t] current node: %v\n", taintInfo.objroot, recurse, node)
	taintInfo.prevval = node.GetValue()
	propagateTaintNearbyFromNode(graph, node, prevVal, recurse, taintInfo, visited, checkTaintInfo, upwards)
	propagateTaintNearbyToNode(graph, node, recurse, taintInfo, visited, checkTaintInfo, upwards)

	if node.IsUsedInBson() {
		return
	}
	if ok, _ := ssaValueIsUsedInMongoBsonFilter(graph, node.GetValue()); ok {
		node.EnableUsedInBson()
		return // skip
	}
}

func propagateTaintNearbyFromNode(graph *ssagraph.SSAGraph, node *ssagraph.SSANode, prevVal *ssa.Value, recurse bool, taintInfo TaintInfo, visited map[ssa.Value]bool, checkTaintInfo *CheckTaintInfo, upwards bool) {
	logrus.Debugf("[TAINT NEARBY] [PART_1] [ROOT=%t] [RECURSE=%t] current node: %v\n", taintInfo.objroot, recurse, node)
	for _, edge := range graph.GetEdgesFromNode(node) {
		toNode := edge.GetToNode()
		logrus.Debugf("\t[TAINT NEARBY] [PART_1] [r=%t] [FromNode %s] edge (%s) to node: %v\n", recurse, node.GetValue().Name(), edge.GetTypeString(), toNode)

		switch edge.GetType() {

		case ssagraph.EDGE_FIELD:
			propagateTaintNearbyFromNodeOnField(graph, edge, node, toNode, taintInfo, visited, checkTaintInfo, upwards)

		case ssagraph.EDGE_INDEX:
			propagateTaintNearbyFromNodeOnIndex(graph, edge, node, toNode, taintInfo, visited, checkTaintInfo, upwards)

		case ssagraph.EDGE_MAP_UPDATE, ssagraph.EDGE_MAP_KEY, ssagraph.EDGE_MAP_VALUE, ssagraph.EDGE_LOOKUP_MAP, ssagraph.EDGE_LOOKUP_MAP_INDEX:
			propagateTaintNearbyFromNodeOnMap(graph, edge, node, toNode, taintInfo, visited, checkTaintInfo, upwards)

		case ssagraph.EDGE_BINOP_X, ssagraph.EDGE_BINOP_Y:
			propagateTaintNearbyFromNodeOnBinOp(graph, edge, node, toNode, taintInfo, visited, checkTaintInfo, upwards)

		case ssagraph.EDGE_ARG_ON_CALL:
			if call, ok := toNode.GetValue().(*ssa.Call); ok {
				propagateTaintNearbyFromNodeOnCall(graph, node, toNode, taintInfo, visited, checkTaintInfo, upwards, call)
			}

		case ssagraph.EDGE_STORE_ADDRESS:
			val := toNode.GetInstruction().(*ssa.Store).Val
			valNode := graph.GetNodeByName(val.Name())
			propagateTaintNearby(graph, true, valNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)

		case ssagraph.EDGE_STORE_VALUE:
			addr := toNode.GetInstruction().(*ssa.Store).Addr
			addrNode := graph.GetNodeByName(addr.Name())
			propagateTaintNearby(graph, true, addrNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)

		case ssagraph.EDGE_RETURN_ON, ssagraph.EDGE_EXTRACT:
			// skip

		default:
			propagateTaintNearby(graph, true, toNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
		}
	}
}

func propagateTaintNearbyToNode(graph *ssagraph.SSAGraph, node *ssagraph.SSANode, recurse bool, taintInfo TaintInfo, visited map[ssa.Value]bool, checkTaintInfo *CheckTaintInfo, upwards bool) {
	logrus.Debugf("[TAINT NEARBY] [PART_2] [ROOT=%t] [RECURSE=%t] current node: %v\n", taintInfo.objroot, recurse, node)
	for _, edge := range graph.GetEdgesToNode(node) {
		fromNode := edge.GetFromNode()
		if fromNode.IsUsedInBson() {
			return
		}
		if ok, _ := ssaValueIsUsedInMongoBsonFilter(graph, fromNode.GetValue()); ok {
			fromNode.EnableUsedInBson()
			return // skip
		}

		switch edge.GetType() {

		case ssagraph.EDGE_FIELD:
			propagateTaintNearbyToNodeOnField(graph, edge, node, fromNode, taintInfo, visited, checkTaintInfo, upwards)

		case ssagraph.EDGE_INDEX:
			propagateTaintNearbyToNodeOnIndex(graph, edge, node, fromNode, taintInfo, visited, checkTaintInfo, upwards)

		case ssagraph.EDGE_MAP_UPDATE, ssagraph.EDGE_MAP_KEY, ssagraph.EDGE_MAP_VALUE, ssagraph.EDGE_LOOKUP_MAP, ssagraph.EDGE_LOOKUP_MAP_INDEX:
			propagateTaintNearbyToNodeOnMap(graph, edge, node, fromNode, taintInfo, visited, checkTaintInfo, upwards)

		case ssagraph.EDGE_BINOP_X, ssagraph.EDGE_BINOP_Y:
			propagateTaintNearbyToNodeOnBinOp(graph, edge, node, fromNode, taintInfo, visited, checkTaintInfo, upwards)

		case ssagraph.EDGE_ARG_ON_CALL:
			if call, ok := node.GetValue().(*ssa.Call); ok {
				propagateTaintNearbyToNodeOnCall(graph, taintInfo, visited, checkTaintInfo, upwards, call)
			}

		case ssagraph.EDGE_STORE_ADDRESS:
			val := fromNode.GetInstruction().(*ssa.Store).Val
			valNode := graph.GetNodeByName(val.Name())
			propagateTaintNearby(graph, true, valNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)

		case ssagraph.EDGE_STORE_VALUE:
			addr := fromNode.GetInstruction().(*ssa.Store).Addr
			addrNode := graph.GetNodeByName(addr.Name())
			propagateTaintNearby(graph, true, addrNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)

		case ssagraph.EDGE_USAGE, ssagraph.EDGE_PHI_ON:
			propagateTaintNearby(graph, true, fromNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)

		case ssagraph.EDGE_RETURN_ON, ssagraph.EDGE_EXTRACT:
			// skip

		default:
			propagateTaintNearby(graph, true, fromNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
		}
	}
}

func propagateTaintFetchUpwards(graph *ssagraph.SSAGraph, val ssa.Value, taintInfo TaintInfo, visited map[ssa.Value]bool, checkTaintInfo *CheckTaintInfo, upwards bool) {
	taintInfo = taintInfo.updateValue(val)

	logrus.Tracef("[TAINT FETCH] visiting %s: %s // TAINT INFO (%s, %s)\n", val.Name(), val.String(), taintInfo.getObjectPath(), taintInfo.getDatabasePath())
	if visited[val] {
		logrus.Tracef("\t[TAINT FETCH] skipping value %s: %s\n", val.Name(), val.String())
		return
	}
	visited[val] = true

	node := graph.GetNodeByName(val.Name())
	logrus.Tracef("\t[TAINT FETCH] checking upper taints: %v\n", node.GetTaints())

	if node.IsUsedInBson() {
		return
	}
	if ok, _ := ssaValueIsUsedInMongoBsonFilter(graph, node.GetValue()); ok {
		node.EnableUsedInBson()
		return // skip
	}

	// 1. taint "subpaths" for current variable and save to later taint the corresponding "subobjects" that requested the upper taint
	for objPath, taints := range node.GetTaints() {

		logrus.Tracef("\t[TAINT FETCH] comparing prefixes:\n\t - tainted obj path:\t %s\n\t - bottom to upper:\t %s\n", objPath, taintInfo.getObjectFullPath())

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
				taintInfoTmp.dbTaint.dbpath = taint.GetDatabasePath()
				taintInfoTmp.dbTaint.dbcall = taint.GetDatabaseCall()
				taintInfoTmp.setReadKey(taint.IsReadKey())
				taintInfoTmp.setReadValue(taint.IsReadValue())
				doTaintNode(node, taintInfoTmp, TAINT_MODE_FETCH_UPWARDS)

				// so that we can later taint the bottom node
				dbFieldIndirect := taintInfoTmp.getDatabasePath() + taintInfo.getObjectPath()
				if taintInfoTmp.getDatabaseCall() == nil {
					// FIXME: verify this
					logrus.Tracef("[TAINT FETCH] [4] nil db call for taint info tmp: %v\n", taintInfoTmp)
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
		}

	}
}

func RunTainter(graph *ssagraph.SSAGraph) {
	logrus.WithField("graph", graph.String()).Debugf("[SSA TAINTER] running SSA tainter...")
	databaseCallRegistry := make(map[*ssagraph.DatabaseCall][]ValFieldPath)
	registerCalls(graph, databaseCallRegistry)
	runTainterOnCalls(graph, databaseCallRegistry)
	runTainterOnParameters(graph)
	runTainterOnReturns(graph)
}

func registerCalls(graph *ssagraph.SSAGraph, databaseCallRegistry map[*ssagraph.DatabaseCall][]ValFieldPath) {
	for _, node := range graph.GetNodes() {
		ok := registerDatabaseCall(graph, node, databaseCallRegistry)
		if ok {
			logrus.WithField("graph", graph.String()).WithField("node", node.String()).Debugf("[SSA TAINTER] database call")
		} else {
			ok = registerServiceCall(graph, node)
			if ok {
				logrus.WithField("graph", graph.String()).WithField("node", node.String()).Debugf("[SSA TAINTER] service call")
			} else {
				registerMethodCall(graph, node)
			}
		}
	}
}

func registerDatabaseCall(graph *ssagraph.SSAGraph, node *ssagraph.SSANode, databaseCallRegistry map[*ssagraph.DatabaseCall][]ValFieldPath) bool {
	if database, collectionOrTopic, method, opType, valFieldPathLst, ok := isDatabaseCall(graph, node.GetValue()); ok {
		var argNodes []*ssagraph.SSANode
		for _, valFieldPath := range valFieldPathLst {
			argNodes = append(argNodes, graph.GetNodeByName(valFieldPath.val.Name()))
		}
		callId := ssagraph.ComputeCallID(graph, node)
		dbCall := ssagraph.NewDatabaseCall(callId, node, argNodes, database, collectionOrTopic, method, opType)
		graph.AddDatabaseCall(dbCall)
		graph.AddCall(dbCall)
		databaseCallRegistry[dbCall] = valFieldPathLst
		return true
	}
	return false
}

func registerServiceCall(graph *ssagraph.SSAGraph, node *ssagraph.SSANode) bool {
	// keep track of arguments passed in service RPCs
	// so that we can mark their indirect taints
	if service, method, funcShortPath, args, call, ok := isServiceCall(graph, node.GetValue()); ok {
		// keep track of objects passed as arguments
		var argNodes []*ssagraph.SSANode
		for _, arg := range args {
			argNodes = append(argNodes, graph.GetNodeByName(arg.Name()))
		}

		// keep track of objects extracted from returns
		var retNodes []*ssagraph.SSANode
		for _, val := range getReturnSSAValuesFromCall(graph, call) {
			retNodes = append(retNodes, graph.GetNodeByName(val.Name()))
		}

		callId := ssagraph.ComputeCallID(graph, node)
		svcCall := ssagraph.NewServiceCall(callId, node, argNodes, retNodes, service, method, funcShortPath)
		graph.AddServiceCall(svcCall)
		graph.AddCall(svcCall)
		return true
	}
	return false
}

func registerMethodCall(graph *ssagraph.SSAGraph, node *ssagraph.SSANode) bool {
	if method, fnShortPath, args, call, ok := isMethodCall(graph, node.GetValue()); ok {
		var argNodes []*ssagraph.SSANode
		for _, arg := range args {
			argNodes = append(argNodes, graph.GetNodeByName(arg.Name()))
		}
		var retNodes []*ssagraph.SSANode
		for _, val := range getReturnSSAValuesFromCall(graph, call) {
			retNode := graph.GetNodeByName(val.Name())
			retNodes = append(retNodes, retNode)
		}

		callId := ssagraph.ComputeCallID(graph, node)
		methodCall := ssagraph.NewMethodCall(callId, node, argNodes, retNodes, method, fnShortPath)
		graph.AddMethodCall(methodCall)
		graph.AddCall(methodCall)
		return true
	}
	return false
}

func getReturnSSAValuesFromCall(graph *ssagraph.SSAGraph, call *ssa.Call) []ssa.Value {
	var vals []ssa.Value
	callNode := graph.GetNodeByName(call.Name())
	if call.Call.Signature().Results().Len() > 1 {
		for _, edge := range graph.GetEdgesFromNode(callNode) {
			if edge.GetType() == ssagraph.EDGE_EXTRACT {
				vals = append(vals, edge.GetToNode().GetValue())
			}
		}
	} else {
		// when there is only one return value then there
		// are no extract instructions and the value is just
		// the one declared when invoking the function
		vals = append(vals, callNode.GetValue())
	}
	return vals
}

func runTainterOnCalls(graph *ssagraph.SSAGraph, databaseCallRegistry map[*ssagraph.DatabaseCall][]ValFieldPath) {
	for _, call := range graph.GetAllCalls() {
		if dbCall, ok := call.(*ssagraph.DatabaseCall); ok {
			valFieldPathLst := databaseCallRegistry[dbCall]
			taintOnDatabaseCall(graph, dbCall, valFieldPathLst)
		} else if svcCall, ok := call.(*ssagraph.ServiceCall); ok {
			taintOnServiceCall(graph, svcCall)
		}
	}
}

func taintOnDatabaseCall(graph *ssagraph.SSAGraph, dbCall *ssagraph.DatabaseCall, valFieldPathLst []ValFieldPath) {
	var nodesToVisit []*ssagraph.SSANode
	for _, valFieldPath := range valFieldPathLst {
		dbfield := valFieldPath.fieldpath
		obj := valFieldPath.val

		taintInfo := NewTaintInfoDatabase(dbfield, "", nil, dbCall, valFieldPath.readKey, valFieldPath.readValue)

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
		if valFieldPath.bsonCursorMany || valFieldPath.bsonFilterIn || valFieldPath.cacheMultiget {
			taintInfo.objpath += "[*]"
			taintInfo.objroot = false
		}

		visited := make(map[ssa.Value]bool)
		seenTaint = make(map[TaintInfo]bool)
		propagateTaintNearby(graph, false, obj, taintInfo, visited, nil, false)
		seenTaint = nil
		objNode := graph.GetNodeByName(obj.Name())
		nodesToVisit = append(nodesToVisit, objNode)
	}

	checkUpperTaintsForObjects(graph, nodesToVisit)
}

func taintOnServiceCall(graph *ssagraph.SSAGraph, svcCall *ssagraph.ServiceCall) {
	var nodesToVisit []*ssagraph.SSANode
	for _, argNode := range svcCall.GetArguments() {
		nodesToVisit = append(nodesToVisit, argNode)
		arg := argNode.GetValue()
		svpath := svcCall.String() + "." + arg.Name()
		taintInfo := NewTaintInfoService(svpath, "", nil, svcCall)
		visited := make(map[ssa.Value]bool)
		seenTaint = make(map[TaintInfo]bool)
		propagateTaintNearby(graph, false, arg, taintInfo, visited, nil, false)
		seenTaint = nil
	}

	for _, retNode := range svcCall.GetReturns() {
		nodesToVisit = append(nodesToVisit, retNode)
		ret := retNode.GetValue()
		svpath := svcCall.String() + "." + ret.Name()
		taintInfo := NewTaintInfoService(svpath, "", nil, svcCall)
		visited := make(map[ssa.Value]bool)
		seenTaint = make(map[TaintInfo]bool)
		propagateTaintNearby(graph, false, ret, taintInfo, visited, nil, false)
		seenTaint = nil
	}
	checkUpperTaintsForObjects(graph, nodesToVisit)
}

func runTainterOnParameters(graph *ssagraph.SSAGraph) {
	paramNodes := graph.GetFuncParametersExceptMemberAndContext()
	checkUpperTaintsForObjects(graph, paramNodes)
}

func runTainterOnReturns(graph *ssagraph.SSAGraph) {
	retNodesLst := graph.GetReturnsLst()
	for _, retNodes := range retNodesLst {
		checkUpperTaintsForObjects(graph, retNodes)
	}
}

func checkUpperTaintsForObjects(graph *ssagraph.SSAGraph, nodesToVisit []*ssagraph.SSANode) {
	return
	// check for upper taints affecting the current database/service calls objects
	for _, originNode := range nodesToVisit {
		// EVAL: fmt.Println()
		logrus.Tracef("[TAINT] check upper taints for node: %v\n", originNode.String())
		visited := make(map[ssa.Value]bool)
		taintInfo := NewTaintInfoDatabase("", "", nil, nil, false, false)
		checkTaintInfo := NewCheckTaintInfo()
		propagateTaintFetchUpwards(graph, originNode.GetValue(), taintInfo, visited, checkTaintInfo, false)
		node := graph.GetNodeByName(originNode.GetValue().Name())

		// indirect taints
		for _, taint := range checkTaintInfo.indirectTaints {
			if taint.dbcall == nil {
				logrus.Fatalf("[1] nil db call for taint: %v\n", taint)
			}
			taintInfo := NewTaintInfoDatabase(taint.dbpath, "", originNode.GetValue(), taint.dbcall, taint.readKey, taint.readVal)
			doTaintNode(node, taintInfo, TAINT_MODE_NEARBY)
		}

		// inherited taints
		for objpath, taints := range checkTaintInfo.inheritedTaints {
			logrus.Tracef("[TAINT] check inherited taints for objpath (%s): %v\n", objpath, taints)
			for _, taint := range taints {
				if taint.dbcall == nil {
					// FIXME: verify this
					logrus.Tracef("[2] nil db call for taint: %v\n", taint)
				} else {
					taintInfo := NewTaintInfoDatabase(taint.dbpath, objpath, originNode.GetValue(), taint.dbcall, taint.readKey, taint.readVal)
					doTaintNode(node, taintInfo, TAINT_MODE_NEARBY)
				}
			}
		}
	}
}
