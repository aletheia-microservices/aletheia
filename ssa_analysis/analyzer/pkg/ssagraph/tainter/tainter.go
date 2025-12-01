package tainter

import (
	"log"
	"sort"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/ssagraph"
	"analyzer/pkg/utils"
)

var seenTaint map[TaintInfoData]bool

func RunTainter(graph *ssagraph.SSAGraph) {
	logrus.WithField("graph", graph.String()).Debugf("[SSA TAINTER] running SSA tainter...")
	databaseCallRegistry := make(map[*ssagraph.DatabaseCall][]ValFieldPath)
	registerCalls(graph, databaseCallRegistry)
	runTainterOnCalls(graph, databaseCallRegistry)
}

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

func generateRootTaintInfoFromTaint(node *ssagraph.SSANode, taint *ssagraph.SSATaint) TaintInfo {
	if taint.IsDatabaseTaint() {
		return NewTaintInfoDatabase(taint.GetDatabasePath(), "", node.GetValue(), taint.GetDatabaseCall(), taint.IsReadKey(), taint.IsReadValue())
	} else if taint.IsServiceTaint() {
		return NewTaintInfoService(taint.GetServicePath(), "", node.GetValue(), taint.GetServiceCall())
	}
	log.Panicf("[TAINT INFO FROM TAINT] unexpected type for taint: %s\n", taint.String())
	return TaintInfo{}
}

func propagateTaintNearby(graph *ssagraph.SSAGraph, recurse bool, val ssa.Value, taintInfo TaintInfo, visited map[ssa.Value]bool, upwards bool) {
	if val == nil {
		log.Panicf("[TAINT NEARBY] unexpected nil val // TAINT INFO (_obj%s, %s)\n", taintInfo.getObjectPath(), taintInfo.getDatabasePath())
	}

	taintInfo = taintInfo.updateValue(val)

	var prevValStr string
	if taintInfo.prevval != nil {
		prevValStr = taintInfo.prevval.Name() + ": " + taintInfo.prevval.String()
	}

	logrus.WithField("graph", graph.String()).
		WithField("prev", prevValStr).
		WithField("val", val.Name()).
		WithField("taint_info", taintInfo.String()).
		Debugf("[TAINT NEARBY] visiting")

	if seenTaint[taintInfo.TaintInfoData] {
		logrus.WithField("val", val.Name).WithField("taint_info", taintInfo.String()).
			Tracef("[TAINT NEARBY] skipping (seen taint)...")
		return
	}
	seenTaint[taintInfo.TaintInfoData] = true

	if visited[val] {
		logrus.WithField("val", val.Name).WithField("taint_info", taintInfo.String()).
			Tracef("[TAINT NEARBY] skipping (visited val)...")
		return
	}
	visited[val] = true

	node := graph.GetNodeByName(val.Name())
	doTaintNode(node, taintInfo, TAINT_MODE_NEARBY)

	logrus.Debugf("[TAINT NEARBY] [PART_1] [ROOT=%t] [RECURSE=%t] current node: %v\n", taintInfo.objroot, recurse, node)
	taintInfo.prevval = node.GetValue()
	propagateTaintNearbyFromNode(graph, node, recurse, taintInfo, visited, upwards)
	propagateTaintNearbyToNode(graph, node, recurse, taintInfo, visited, upwards)

	if node.IsUsedInBson() {
		return
	}
	if ok, _ := ssaValueIsUsedInMongoBsonFilter(graph, node.GetValue()); ok {
		node.EnableUsedInBson()
		return // skip
	}
}

func propagateTaintNearbyFromNode(graph *ssagraph.SSAGraph, node *ssagraph.SSANode, recurse bool, taintInfo TaintInfo, visited map[ssa.Value]bool, upwards bool) {
	logrus.Debugf("[TAINT NEARBY] [PART_1] [ROOT=%t] [RECURSE=%t] current node: %v\n", taintInfo.objroot, recurse, node)
	for _, edge := range graph.GetEdgesFromNode(node) {
		toNode := edge.GetToNode()
		logrus.Debugf("\t[TAINT NEARBY] [PART_1] [r=%t] [FromNode %s] edge (%s) to node: %v\n", recurse, node.GetValue().Name(), edge.GetTypeString(), toNode)

		switch edge.GetType() {

		case ssagraph.EDGE_FIELD:
			propagateTaintNearbyFromNodeOnField(graph, edge, node, toNode, taintInfo, visited, upwards)

		case ssagraph.EDGE_INDEX:
			propagateTaintNearbyFromNodeOnIndex(graph, edge, node, toNode, taintInfo, visited, upwards)

		case ssagraph.EDGE_MAP_UPDATE, ssagraph.EDGE_MAP_KEY, ssagraph.EDGE_MAP_VALUE, ssagraph.EDGE_LOOKUP_MAP, ssagraph.EDGE_LOOKUP_MAP_INDEX:
			propagateTaintNearbyFromNodeOnMap(graph, edge, node, toNode, taintInfo, visited, upwards)

		case ssagraph.EDGE_BINOP_X, ssagraph.EDGE_BINOP_Y:
			propagateTaintNearbyFromNodeOnBinOp(graph, edge, node, toNode, taintInfo, visited, upwards)

		case ssagraph.EDGE_ARG_ON_CALL:
			if call, ok := toNode.GetValue().(*ssa.Call); ok {
				propagateTaintNearbyFromNodeOnCall(graph, node, toNode, taintInfo, visited, upwards, call)
			}

		case ssagraph.EDGE_STORE_ADDRESS:
			val := toNode.GetInstruction().(*ssa.Store).Val
			valNode := graph.GetNodeByName(val.Name())
			propagateTaintNearby(graph, true, valNode.GetValue(), taintInfo, visited, upwards)

		case ssagraph.EDGE_STORE_VALUE:
			addr := toNode.GetInstruction().(*ssa.Store).Addr
			addrNode := graph.GetNodeByName(addr.Name())
			propagateTaintNearby(graph, true, addrNode.GetValue(), taintInfo, visited, upwards)

		case ssagraph.EDGE_RETURN_ON, ssagraph.EDGE_EXTRACT:
			// skip

		default:
			propagateTaintNearby(graph, true, toNode.GetValue(), taintInfo, visited, upwards)
		}
	}
}

func propagateTaintNearbyToNode(graph *ssagraph.SSAGraph, node *ssagraph.SSANode, recurse bool, taintInfo TaintInfo, visited map[ssa.Value]bool, upwards bool) {
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
			propagateTaintNearbyToNodeOnField(graph, edge, node, fromNode, taintInfo, visited, upwards)

		case ssagraph.EDGE_INDEX:
			propagateTaintNearbyToNodeOnIndex(graph, edge, node, fromNode, taintInfo, visited, upwards)

		case ssagraph.EDGE_MAP_UPDATE, ssagraph.EDGE_MAP_KEY, ssagraph.EDGE_MAP_VALUE, ssagraph.EDGE_LOOKUP_MAP, ssagraph.EDGE_LOOKUP_MAP_INDEX:
			propagateTaintNearbyToNodeOnMap(graph, edge, node, fromNode, taintInfo, visited, upwards)

		case ssagraph.EDGE_BINOP_X, ssagraph.EDGE_BINOP_Y:
			propagateTaintNearbyToNodeOnBinOp(graph, edge, node, fromNode, taintInfo, visited, upwards)

		case ssagraph.EDGE_ARG_ON_CALL:
			if call, ok := node.GetValue().(*ssa.Call); ok {
				propagateTaintNearbyToNodeOnCall(graph, taintInfo, visited, upwards, call)
			}

		case ssagraph.EDGE_STORE_ADDRESS:
			val := fromNode.GetInstruction().(*ssa.Store).Val
			valNode := graph.GetNodeByName(val.Name())
			propagateTaintNearby(graph, true, valNode.GetValue(), taintInfo, visited, upwards)

		case ssagraph.EDGE_STORE_VALUE:
			addr := fromNode.GetInstruction().(*ssa.Store).Addr
			addrNode := graph.GetNodeByName(addr.Name())
			propagateTaintNearby(graph, true, addrNode.GetValue(), taintInfo, visited, upwards)

		case ssagraph.EDGE_USAGE, ssagraph.EDGE_PHI_ON:
			propagateTaintNearby(graph, true, fromNode.GetValue(), taintInfo, visited, upwards)

		case ssagraph.EDGE_RETURN_ON, ssagraph.EDGE_EXTRACT:
			// skip

		default:
			propagateTaintNearby(graph, true, fromNode.GetValue(), taintInfo, visited, upwards)
		}
	}
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
	if method, fnShortPath, bindings, args, call, fn, ok := isMethodCall(node.GetInstruction(), node.GetValue()); ok {
		var argNodes []*ssagraph.SSANode
		for _, arg := range args {
			argNodes = append(argNodes, graph.GetNodeByName(arg.Name()))
		}
		var retNodes []*ssagraph.SSANode
		if call != nil {
			for _, val := range getReturnSSAValuesFromCall(graph, call) {
				retNode := graph.GetNodeByName(val.Name())
				retNodes = append(retNodes, retNode)
			}
		}

		// go routine
		var bindNodes []*ssagraph.SSANode
		for _, b := range bindings {
			bindNodes = append(bindNodes, graph.GetNodeByName(b.Name()))
		}
		sort.Slice(bindNodes, func(i, j int) bool {
			return utils.LessT(bindNodes[i].GetName(), bindNodes[j].GetName())
		})

		var methodCall *ssagraph.MethodCall
		if call != nil {
			callId := ssagraph.ComputeCallID(graph, node)
			methodCall = ssagraph.NewMethodCall(callId, node, argNodes, retNodes, method, fnShortPath)
		} else {
			// go routine

			callId := utils.GetShortFunctionPath(fn.String())
			// because instructions may appear multiple times
			if graph.HasMethodCall(&ssagraph.MethodCall{ID: callId}) {
				return false
			}
			methodCall = ssagraph.NewMethodCallGoRoutine(callId, node.GetID(), node, bindNodes, argNodes, retNodes, method, fnShortPath)
		}
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
	seenTaint = make(map[TaintInfoData]bool)
	for _, call := range graph.GetAllCalls() {
		if dbCall, ok := call.(*ssagraph.DatabaseCall); ok {
			valFieldPathLst := databaseCallRegistry[dbCall]
			taintOnDatabaseCall(graph, dbCall, valFieldPathLst)
		} else if svcCall, ok := call.(*ssagraph.ServiceCall); ok {
			taintOnServiceCall(graph, svcCall)
		}
	}
	seenTaint = nil
}

func taintOnDatabaseCall(graph *ssagraph.SSAGraph, dbCall *ssagraph.DatabaseCall, valFieldPathLst []ValFieldPath) {
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
		propagateTaintNearby(graph, false, obj, taintInfo, visited, false)
	}
}

func taintOnServiceCall(graph *ssagraph.SSAGraph, svcCall *ssagraph.ServiceCall) {
	for _, argNode := range svcCall.GetArguments() {
		arg := argNode.GetValue()
		svpath := svcCall.String() + "." + arg.Name()
		taintInfo := NewTaintInfoService(svpath, "", nil, svcCall)
		visited := make(map[ssa.Value]bool)
		propagateTaintNearby(graph, false, arg, taintInfo, visited, false)
	}

	for _, retNode := range svcCall.GetReturns() {
		ret := retNode.GetValue()
		svpath := svcCall.String() + "." + ret.Name()
		taintInfo := NewTaintInfoService(svpath, "", nil, svcCall)
		visited := make(map[ssa.Value]bool)
		propagateTaintNearby(graph, false, ret, taintInfo, visited, false)
	}
}
