package abstractgraph

import (
	"math"
	"sort"

	"github.com/sirupsen/logrus"

	"analyzer/pkg/app/backends"
	"analyzer/pkg/common"
	"analyzer/pkg/ssagraph"
	"analyzer/pkg/utils"
)

func ssaTaintDatabaseToAbstractTaint(graph *AbstractCallGraph, ssaTaintsMap map[string][]*ssagraph.SSATaint) map[string][]*AbstractTaint {
	abstractTaintsMap := make(map[string][]*AbstractTaint, 0)
	for objPath, ssaTaints := range ssaTaintsMap {
		var abstractTaints []*AbstractTaint
		for _, ssaTaint := range ssaTaints {
			if ssaTaint.IsDatabaseTaint() {
				dbPath := ssaTaint.GetDatabaseCall().GetDatabasePath()
				dbname := ssaTaint.GetDatabaseCall().GetDatabaseName()
				dbNode := graph.GetNodeByNameIfExists(dbPath)
				schemaName := ssaTaint.GetDatabaseCall().GetSchemaName()
				if dbNode == nil {
					dbNode = NewAbstractNode(dbPath, NODE_DATABASE, "", "", dbname, schemaName)
					graph.AddNode(dbPath, dbNode)

					if !graph.GetApp().HasDatabase(dbname) {
						logrus.Fatalf("database (%s) not found", dbname)
					}
					db := graph.GetApp().GetDatabaseByName(dbname)
					if !db.HasSchema(schemaName) {
						db.AddSchema(backends.NewSchema(schemaName, db))
					}
				}
				taint := NewAbstractTaint(
					ssaTaint.GetT(),
					ssaTaint.GetDatabasePath(),
					ssaTaint.GetDatabaseCall().GetID(),
					ssaTaint.GetDatabaseCall().GetOpType(),
					true, false, ssaTaint.IsReadKey(), ssaTaint.IsReadValue(),
				)

				abstractTaints = append(abstractTaints, taint)
				// EVAL: logrus.Tracef("[SSA TO ABSTRACT TAINT] new taint on object path (%s): %s\n", objPath, taint.LongString())
			}
		}
		if abstractTaints != nil {
			abstractTaintsMap[objPath] = abstractTaints
		}
	}
	return abstractTaintsMap
}

func ssaTaintServiceToAbstractTrace(graph *AbstractCallGraph, ssaTaintsMap map[string][]*ssagraph.SSATaint) map[string][]*AbstractTrace {
	abstractTaintsMap := make(map[string][]*AbstractTrace, 0)
	for objPath, ssaTaints := range ssaTaintsMap {
		var abstractTraces []*AbstractTrace
		for _, ssaTaint := range ssaTaints {
			if ssaTaint.IsServiceTaint() {
				trace := NewAbstractTrace(
					ssaTaint.GetT(),
					ssaTaint.GetServicePath(),
					ssaTaint.GetServiceCall().GetID(),
				)
				abstractTraces = append(abstractTraces, trace)
				// EVAL: logrus.Tracef("[SSA TO ABSTRACT TRACE] new trace on object path (%s): %s\n", objPath, trace.LongString())
			}
		}
		if abstractTraces != nil {
			abstractTaintsMap[objPath] = abstractTraces
		}
	}
	return abstractTaintsMap
}

func ComputeGraphStats(graph *AbstractCallGraph) {
	clientNode := graph.GetNodeByNameIfExists("client")
	if clientNode == nil {
		return
	}

	visited := make(map[*AbstractNode]int)
	maxDepth := 0

	for _, edge := range graph.GetEdgesFromNode(clientNode) {
		computeGraphStatsHelper(graph, edge.GetToNode(), 1, visited, &maxDepth)
	}

	logrus.Infof("Max call depth from client: %d", maxDepth)
}

func computeGraphStatsHelper(graph *AbstractCallGraph, node *AbstractNode, depth int, visited map[*AbstractNode]int, maxDepth *int) {
	if prevDepth, ok := visited[node]; ok && prevDepth >= depth {
		return
	}
	visited[node] = depth

	node.SetCallDepth(depth)
	if depth > *maxDepth {
		*maxDepth = depth
	}

	for _, edge := range graph.GetEdgesFromNode(node) {
		to := edge.GetToNode()

		if to.GetNodeType() == NODE_SERVICE {
			node.IncrStatelessFanout()
		}
		computeGraphStatsHelper(graph, to, depth+1, visited, maxDepth)

	}
}

func GatherGraphStats(graph *AbstractCallGraph) {
	gatherGraphStatsFanout(graph)
	gatherGraphStatsCounts(graph)
}

func gatherGraphStatsCounts(graph *AbstractCallGraph) {
	var datastoreCount int
	var entrypointsCount int
	var services map[string]bool = make(map[string]bool)
	for _, node := range graph.GetNodes() {
		if node.GetName() != "client" {
			if node.GetNodeType() == NODE_SERVICE {
				if _, ok := services[node.GetServiceName()]; !ok {
					services[node.GetServiceName()] = true
				}
			} else {
				datastoreCount++
			}
		}
	}
	clientNode := graph.GetNodeByNameIfExists("client")
	entrypointsCount = len(graph.GetEdgesFromNode(clientNode))
	logrus.Infof("[INFO] [ABSTRACT GRAPH] #reqs. stateful (#rpcs)=%d, #reqs. stateless (#db accesses)=%d, #services=%d, #datastores=%d, #callgraphs=%d\n", graph.GetRPCCount(), graph.GetDBAccessCount(), len(services), datastoreCount, entrypointsCount)
}

func gatherGraphStatsFanout(graph *AbstractCallGraph) {
	var fanouts []int
	for _, node := range graph.GetNodes() {
		f := node.GetFanout()
		if f > 0 {
			fanouts = append(fanouts, f)
		}
	}
	if len(fanouts) == 0 {
		logrus.Warn("no nodes with fanout > 0 found")
		return
	}
	// sort ascending for median / percentiles
	sort.Ints(fanouts)

	// --- average ---
	var sum int
	for _, f := range fanouts {
		sum += f
	}
	avg := float64(sum) / float64(len(fanouts))
	// --- median ---
	n := len(fanouts)
	var median float64
	if n%2 == 1 {
		// odd
		median = float64(fanouts[n/2])
	} else {
		// even
		median = float64(fanouts[n/2-1]+fanouts[n/2]) / 2.0
	}
	// --- p90 (90th percentile, 1-based Ceil index) ---
	// idx = ceil(0.9 * n) - 1
	idx := int(math.Ceil(0.9*float64(n))) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= n {
		idx = n - 1
	}
	p90 := fanouts[idx]

	logrus.Infof(
		"[INFO] [ABSTRACT GRAPH] fanout stats (>0): count=%d, avg=%.2f, median=%.2f, p90=%d\n",
		n, avg, median, p90,
	)
}

func Parse(graph *AbstractCallGraph, funcshortpath string, entrypoint bool, funcGraphs map[string]*ssagraph.SSAGraph) {
	// dummy node
	clientNode := graph.GetNodeByNameIfExists("client")
	if clientNode == nil {
		clientNode = NewAbstractNode("client", NODE_CLIENT, "", "", "", "")
		graph.AddNode("client", clientNode)
	}

	ssaGraph := funcGraphs[funcshortpath]
	// EVAL: logrus.Tracef("[ABSTRACTGRAPH] got ssa graph for (%s): %v\n", funcshortpath, ssaGraph)

	name := ssaGraph.GetServiceWithMethod()
	node := graph.GetNodeByNameIfExists(name)

	var created bool
	if node == nil {
		created = true
		node = NewAbstractNode(name, NODE_SERVICE, ssaGraph.GetService(), ssaGraph.GetMethodName(), "", "")
		graph.AddNode(name, node)

		// EVAL: logrus.Tracef("[ABSTRACTGRAPH] creating node with (%d) params: %s\n", len(ssaGraph.GetFuncParametersExceptMemberAndContext()), node)
		for _, funcParam := range ssaGraph.GetFuncParametersExceptMemberAndContext() {
			obj := NewAbstractObject(funcParam.GetName(), ssaTaintDatabaseToAbstractTaint(graph, funcParam.GetTaints()), ssaTaintServiceToAbstractTrace(graph, funcParam.GetTaints()))
			// EVAL: logrus.Tracef("[debug] (1) added param (%s) to node (%s)\n", obj.String(), node.String())
			node.AddParam(obj)
		}
	}

	// build dummy edges for entrypoints
	if entrypoint {
		edge := NewAbstractEdge("", funcshortpath, utils.ExtractMethodNameFromShortFunctionPath(funcshortpath), clientNode, node, common.OP_UNDEFINED, EDGE_SERVICE_ENTRYPOINT)
		for _, funcParam := range ssaGraph.GetFuncParametersExceptMemberAndContext() {
			arg := NewAbstractObject(funcParam.GetName(), make(map[string][]*AbstractTaint), make(map[string][]*AbstractTrace))
			edge.AddArgument(arg)
		}
		graph.AddEdge(edge)
		graph.rpcs++
	}

	if !created && node != nil && node.IsParsed() {
		// EVAL: logrus.Tracef("[ABSTRACTGRAPH] ignoring parsed node: %s\n", node.String())
		return
	}

	node.SetParsed()

	// finalize parsing
	// EVAL: logrus.Tracef("[ABSTRACTGRAPH] parsing returns for node: %s\n", node.String())
	retsLst := ssaGraph.GetReturnsLst()
	var retsObjs []*AbstractObject
	// first, just create new abstract objects using the first set of returns (could be any other)
	for _, ret := range retsLst[0] {
		obj := NewAbstractObject(ret.GetValue().Type().String(), ssaTaintDatabaseToAbstractTaint(graph, ret.GetTaints()), ssaTaintServiceToAbstractTrace(graph, ret.GetTaints()))
		obj.addToAllNames(ret.GetValue().Type().String())
		node.AddReturn(obj)
		retsObjs = append(retsObjs, obj)
		// EVAL: logrus.Tracef("\t[ABSTRACTGRAPH] [index=%d] added new return object (%s)\n", obj.String())
	}
	// then, merge taints with corresponding object in the remaining set of returns
	if len(retsLst) > 1 {
		for _, rets := range retsLst[1:] {
			for i, ret := range rets {
				obj := retsObjs[i]
				obj.addToAllNames(ret.GetValue().Type().String())

				// EVAL: logrus.Tracef("\t\t[ABSTRACTGRAPH] ret = %s\n", ret.String())
				MergeTaints(obj, ssaTaintDatabaseToAbstractTaint(graph, ret.GetTaints()), nil, MERGE_MODE_PARSE, "", false)
				MergeTraces(obj, ssaTaintServiceToAbstractTrace(graph, ret.GetTaints()))
				// EVAL: logrus.Tracef("\t\t[ABSTRACTGRAPH] [index=%d] merged taints from (%s) to (%s)\n", i, ret.GetName(), obj.String())
			}
		}
	}

	for _, call := range ssaGraph.GetAllCalls() {
		if serviceCall, ok := call.(*ssagraph.ServiceCall); ok {
			parseServiceCall(graph, node, serviceCall, funcGraphs)
		}

		if databaseCall, ok := call.(*ssagraph.DatabaseCall); ok {
			parseDatabaseCall(graph, node, databaseCall)
		}

		if methodCall, ok := call.(*ssagraph.MethodCall); ok {
			parseMethodCall(graph, node, ssaGraph, methodCall, funcGraphs)
		}
	}
	for _, call := range ssaGraph.GetServiceCalls() {
		Parse(graph, call.GetFuncShortPath(), false, funcGraphs)
	}
}

func parseServiceCall(graph *AbstractCallGraph, node *AbstractNode, serviceCall *ssagraph.ServiceCall, funcGraphs map[string]*ssagraph.SSAGraph) {
	logrus.WithField("node", node.String()).Tracef("[ABSTRACTGRAPH] found service call: %s\n", serviceCall.String())
	toName := serviceCall.GetServiceWithMethod()
	toNode := graph.GetNodeByNameIfExists(toName)

	toSSAGraph := funcGraphs[serviceCall.GetFuncShortPath()]
	if toSSAGraph == nil {
		logrus.Fatalf("could not find ssa graph for short func path (%s)", serviceCall.GetFuncShortPath())
	}

	// create node for the first time
	if toNode == nil {
		toNode = NewAbstractNode(toName, NODE_SERVICE, serviceCall.GetService(), serviceCall.GetMethod(), "", "")
		graph.AddNode(toName, toNode)

		// EVAL: logrus.Tracef("[ABSTRACTGRAPH] creating toNode with (%d) params: %s\n", len(toSSAGraph.GetFuncParametersExceptMemberAndContext()), toNode)
		for _, funcParam := range toSSAGraph.GetFuncParametersExceptMemberAndContext() {
			param := NewAbstractObject(funcParam.GetName(), ssaTaintDatabaseToAbstractTaint(graph, funcParam.GetTaints()), ssaTaintServiceToAbstractTrace(graph, funcParam.GetTaints()))
			toNode.AddParam(param)
			// EVAL: logrus.Tracef("[debug] (2) added param (%s) to node (%s)\n", param.String(), toNode.String())
		}
	}

	edge := NewAbstractEdge(serviceCall.GetT(), serviceCall.GetID(), serviceCall.GetMethod(), node, toNode, common.OP_UNDEFINED, EDGE_SERVICE_RPC)

	// create call arguments
	for _, callArg := range serviceCall.GetArguments() {
		arg := NewAbstractObject(callArg.GetName(), ssaTaintDatabaseToAbstractTaint(graph, callArg.GetTaints()), ssaTaintServiceToAbstractTrace(graph, callArg.GetTaints()))
		edge.AddArgument(arg)
	}

	// create call returns
	for _, callRet := range serviceCall.GetReturns() {
		ret := NewAbstractObject(callRet.GetName(), ssaTaintDatabaseToAbstractTaint(graph, callRet.GetTaints()), ssaTaintServiceToAbstractTrace(graph, callRet.GetTaints()))
		// EVAL: logrus.Tracef("[ABSTRACTGRAPH] [%s] added return object (%s) with taints: %v\n", node.String(), ret.String(), callRet.GetTaints())
		edge.AddReturn(ret)
	}

	// EVAL: logrus.Tracef("[ABSTRACT GRAPH] [SERVICE CALL] added edge: %v\n", edge)
	graph.AddEdge(edge)
	graph.rpcs++
}

func parseDatabaseCall(graph *AbstractCallGraph, node *AbstractNode, databaseCall *ssagraph.DatabaseCall) {
	// EVAL: logrus.Tracef("[ABSTRACTGRAPH] found function with database call: %s\n", databaseCall.String())
	toDatabasePath := databaseCall.GetDatabasePath()
	toNode := graph.GetNodeByNameIfExists(toDatabasePath)
	dbname := databaseCall.GetDatabaseName()
	schema := databaseCall.GetSchemaName()

	if toNode == nil {
		toNode = NewAbstractNode(toDatabasePath, NODE_DATABASE, "", "", dbname, schema)
		graph.AddNode(toDatabasePath, toNode)

		schemaName := databaseCall.GetSchemaName()

		if !graph.GetApp().HasDatabase(dbname) {
			logrus.Fatalf("database (%s) not found", dbname)
		}

		db := graph.GetApp().GetDatabaseByName(dbname)
		if !db.HasSchema(schemaName) {
			db.AddSchema(backends.NewSchema(schemaName, db))
		}
	}

	edge := NewAbstractEdge(databaseCall.GetT(), databaseCall.GetID(), databaseCall.GetMethod(), node, toNode, databaseCall.GetOpType(), EDGE_DATABASE_CALL)

	for _, callArg := range databaseCall.GetArguments() {
		arg := NewAbstractObject(callArg.GetName(), ssaTaintDatabaseToAbstractTaint(graph, callArg.GetTaints()), ssaTaintServiceToAbstractTrace(graph, callArg.GetTaints()))
		edge.AddArgument(arg)
	}

	// create fields if they do not exist yet
	registerDatabaseFields(graph, edge.GetArguments())

	// propagate taints to databases (forward): args (from) >>> params (to)
	for i, toParam := range toNode.GetParams() {
		fromArg := edge.GetArgumentAt(i)
		MergeTaints(toParam, fromArg.GetPrimaryTaints(), nil, MERGE_MODE_PARSE, "", false)
	}

	// EVAL: logrus.Tracef("[ABSTRACT GRAPH] [DATABASE CALL] added edge: %v\n", edge)
	graph.AddEdge(edge)
	graph.dbaccesses++
}

func parseMethodCall(graph *AbstractCallGraph, node *AbstractNode, fromSSAGraph *ssagraph.SSAGraph, methodCall *ssagraph.MethodCall, funcGraphs map[string]*ssagraph.SSAGraph) {
	toSSAGraph := fromSSAGraph.GetCombinedGraphForMethodCallIfExists(methodCall)
	if toSSAGraph == nil {
		// should never happen
		return
	}

	for _, call := range toSSAGraph.GetAllCalls() {
		if serviceCall, ok := call.(*ssagraph.ServiceCall); ok {
			parseServiceCall(graph, node, serviceCall, funcGraphs)
		}

		if databaseCall, ok := call.(*ssagraph.DatabaseCall); ok {
			parseDatabaseCall(graph, node, databaseCall)
		}

		if methodCall, ok := call.(*ssagraph.MethodCall); ok {
			parseMethodCall(graph, node, fromSSAGraph, methodCall, funcGraphs)
		}
	}
	for _, call := range toSSAGraph.GetServiceCalls() {
		Parse(graph, call.GetFuncShortPath(), false, funcGraphs)
	}
}

func registerDatabaseFields(graph *AbstractCallGraph, args []*AbstractObject) {
	for _, arg := range args {
		for _, taintLst := range arg.GetPrimaryTaints() {
			for _, taint := range taintLst {
				db := graph.GetApp().GetDatabaseByName(utils.ExtractDatabaseNameFromFieldPath(taint.GetDatabasePath()))
				latestSchema := db.GetLastSchema()
				if !latestSchema.HasField(taint.GetDatabasePath()) {
					field := backends.NewField(taint.GetDatabasePath(), db, latestSchema)
					latestSchema.AddField(field)
				}
			}
		}
	}
}
