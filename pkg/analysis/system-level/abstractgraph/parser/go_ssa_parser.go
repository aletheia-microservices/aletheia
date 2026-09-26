package abstractgraphparser

import (
	"github.com/sirupsen/logrus"

	"analyzer/pkg/analysis/common"
	"analyzer/pkg/analysis/service-level/ssagraph"
	"analyzer/pkg/analysis/system-level/abstractgraph"
	abstractgraphtainter "analyzer/pkg/analysis/system-level/abstractgraph/tainter"
	"analyzer/pkg/app/backends"
	"analyzer/pkg/utils"
)

func ssaTaintDatabaseToAbstractTaint(graph *abstractgraph.AbstractCallGraph, ssaTaintsMap map[string][]*ssagraph.SSATaint) map[string][]*abstractgraph.AbstractTaint {
	abstractTaintsMap := make(map[string][]*abstractgraph.AbstractTaint, 0)
	for objPath, ssaTaints := range ssaTaintsMap {
		var abstractTaints []*abstractgraph.AbstractTaint
		for _, ssaTaint := range ssaTaints {
			if ssaTaint.IsDatabaseTaint() {
				dbPath := ssaTaint.GetDatabaseCall().GetDatabasePath()
				dbname := ssaTaint.GetDatabaseCall().GetDatabaseName()
				dbNode := graph.GetNodeByNameIfExists(dbPath)
				schemaName := ssaTaint.GetDatabaseCall().GetSchemaName()
				if dbNode == nil {
					dbNode = abstractgraph.NewAbstractNode(dbPath, abstractgraph.NODE_DATABASE, "", "", dbname, schemaName)
					graph.AddNode(dbPath, dbNode)

					if !graph.GetApp().HasDatabase(dbname) {
						logrus.Fatalf("database (%s) not found", dbname)
					}
					db := graph.GetApp().GetDatabaseByName(dbname)
					if !db.HasSchema(schemaName) {
						db.AddSchema(backends.NewSchema(schemaName, db))
					}
				}
				taint := abstractgraph.NewAbstractTaint(
					ssaTaint.GetT(),
					ssaTaint.GetDatabasePath(),
					ssaTaint.GetDatabaseCall().GetID(),
					ssaTaint.GetDatabaseCall().GetOpType(),
					true, false, ssaTaint.IsReadKey(), ssaTaint.IsReadValue(),
				)

				abstractTaints = append(abstractTaints, taint)
			}
		}
		if abstractTaints != nil {
			abstractTaintsMap[objPath] = abstractTaints
		}
	}
	return abstractTaintsMap
}

func ssaTaintServiceToAbstractTrace(graph *abstractgraph.AbstractCallGraph, ssaTaintsMap map[string][]*ssagraph.SSATaint) map[string][]*abstractgraph.AbstractTrace {
	abstractTaintsMap := make(map[string][]*abstractgraph.AbstractTrace, 0)
	for objPath, ssaTaints := range ssaTaintsMap {
		var abstractTraces []*abstractgraph.AbstractTrace
		for _, ssaTaint := range ssaTaints {
			if ssaTaint.IsServiceTaint() {
				trace := abstractgraph.NewAbstractTrace(
					ssaTaint.GetT(),
					ssaTaint.GetServicePath(),
					ssaTaint.GetServiceCall().GetID(),
				)
				abstractTraces = append(abstractTraces, trace)
			}
		}
		if abstractTraces != nil {
			abstractTaintsMap[objPath] = abstractTraces
		}
	}
	return abstractTaintsMap
}

func Parse(graph *abstractgraph.AbstractCallGraph, funcshortpath string, entrypoint bool, funcGraphs map[string]*ssagraph.SSAGraph) {
	// dummy node
	clientNode := graph.GetNodeByNameIfExists("client")
	if clientNode == nil {
		clientNode = abstractgraph.NewAbstractNode("client", abstractgraph.NODE_CLIENT, "", "", "", "")
		graph.AddNode("client", clientNode)
	}

	ssaGraph := funcGraphs[funcshortpath]
	name := ssaGraph.GetServiceWithMethod()
	node := graph.GetNodeByNameIfExists(name)

	var created bool
	if node == nil {
		created = true
		node = abstractgraph.NewAbstractNode(name, abstractgraph.NODE_SERVICE, ssaGraph.GetService(), ssaGraph.GetMethodName(), "", "")
		graph.AddNode(name, node)

		for _, funcParam := range ssaGraph.GetFuncParametersExceptMemberAndContext() {
			obj := abstractgraph.NewAbstractObject(funcParam.GetName(), ssaTaintDatabaseToAbstractTaint(graph, funcParam.GetTaints()), ssaTaintServiceToAbstractTrace(graph, funcParam.GetTaints()))
			node.AddParam(obj)
		}
	}

	// build dummy edges for entrypoints
	if entrypoint {
		edge := abstractgraph.NewAbstractEdge("", funcshortpath, utils.ExtractMethodNameFromShortFunctionPath(funcshortpath), clientNode, node, common.OP_UNDEFINED, abstractgraph.EDGE_SERVICE_ENTRYPOINT)
		for _, funcParam := range ssaGraph.GetFuncParametersExceptMemberAndContext() {
			arg := abstractgraph.NewAbstractObject(funcParam.GetName(), make(map[string][]*abstractgraph.AbstractTaint), make(map[string][]*abstractgraph.AbstractTrace))
			edge.AddArgument(arg)
		}
		graph.AddEdge(edge)
		graph.IncrRPCs()
	}

	if !created && node != nil && node.IsParsed() {
		return
	}

	node.SetParsed()

	// finalize parsing
	retsLst := ssaGraph.GetReturnsLst()
	var retsObjs []*abstractgraph.AbstractObject
	// first, just create new abstract objects using the first set of returns (could be any other)
	for _, ret := range retsLst[0] {
		obj := abstractgraph.NewAbstractObject(ret.GetValue().Type().String(), ssaTaintDatabaseToAbstractTaint(graph, ret.GetTaints()), ssaTaintServiceToAbstractTrace(graph, ret.GetTaints()))
		obj.AddToAllNames(ret.GetValue().Type().String())
		node.AddReturn(obj)
		retsObjs = append(retsObjs, obj)
	}
	// then, merge taints with corresponding object in the remaining set of returns
	if len(retsLst) > 1 {
		for _, rets := range retsLst[1:] {
			for i, ret := range rets {
				obj := retsObjs[i]
				obj.AddToAllNames(ret.GetValue().Type().String())

				abstractgraphtainter.MergeTaints(obj, ssaTaintDatabaseToAbstractTaint(graph, ret.GetTaints()), nil, abstractgraphtainter.MERGE_MODE_PARSE, "", false)
				abstractgraphtainter.MergeTraces(obj, ssaTaintServiceToAbstractTrace(graph, ret.GetTaints()))
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

func parseServiceCall(graph *abstractgraph.AbstractCallGraph, node *abstractgraph.AbstractNode, serviceCall *ssagraph.ServiceCall, funcGraphs map[string]*ssagraph.SSAGraph) {
	logrus.WithField("node", node.String()).Tracef("[ABSTRACTGRAPH] found service call: %s\n", serviceCall.String())
	toName := serviceCall.GetServiceWithMethod()
	toNode := graph.GetNodeByNameIfExists(toName)

	toSSAGraph := funcGraphs[serviceCall.GetFuncShortPath()]
	if toSSAGraph == nil {
		logrus.Fatalf("could not find ssa graph for short func path (%s)", serviceCall.GetFuncShortPath())
	}

	// create node for the first time
	if toNode == nil {
		toNode = abstractgraph.NewAbstractNode(toName, abstractgraph.NODE_SERVICE, serviceCall.GetService(), serviceCall.GetMethod(), "", "")
		graph.AddNode(toName, toNode)

		for _, funcParam := range toSSAGraph.GetFuncParametersExceptMemberAndContext() {
			param := abstractgraph.NewAbstractObject(funcParam.GetName(), ssaTaintDatabaseToAbstractTaint(graph, funcParam.GetTaints()), ssaTaintServiceToAbstractTrace(graph, funcParam.GetTaints()))
			toNode.AddParam(param)
		}
	}

	edge := abstractgraph.NewAbstractEdge(serviceCall.GetT(), serviceCall.GetID(), serviceCall.GetMethod(), node, toNode, common.OP_UNDEFINED, abstractgraph.EDGE_SERVICE_RPC)

	// create call arguments
	for _, callArg := range serviceCall.GetArguments() {
		arg := abstractgraph.NewAbstractObject(callArg.GetName(), ssaTaintDatabaseToAbstractTaint(graph, callArg.GetTaints()), ssaTaintServiceToAbstractTrace(graph, callArg.GetTaints()))
		edge.AddArgument(arg)
	}

	// create call returns
	for _, callRet := range serviceCall.GetReturns() {
		ret := abstractgraph.NewAbstractObject(callRet.GetName(), ssaTaintDatabaseToAbstractTaint(graph, callRet.GetTaints()), ssaTaintServiceToAbstractTrace(graph, callRet.GetTaints()))
		edge.AddReturn(ret)
	}

	graph.AddEdge(edge)
	graph.IncrRPCs()
}

func parseDatabaseCall(graph *abstractgraph.AbstractCallGraph, node *abstractgraph.AbstractNode, databaseCall *ssagraph.DatabaseCall) {
	toDatabasePath := databaseCall.GetDatabasePath()
	toNode := graph.GetNodeByNameIfExists(toDatabasePath)
	dbname := databaseCall.GetDatabaseName()
	schema := databaseCall.GetSchemaName()

	if toNode == nil {
		toNode = abstractgraph.NewAbstractNode(toDatabasePath, abstractgraph.NODE_DATABASE, "", "", dbname, schema)
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

	edge := abstractgraph.NewAbstractEdge(databaseCall.GetT(), databaseCall.GetID(), databaseCall.GetMethod(), node, toNode, databaseCall.GetOpType(), abstractgraph.EDGE_DATABASE_CALL)

	for _, callArg := range databaseCall.GetArguments() {
		arg := abstractgraph.NewAbstractObject(callArg.GetName(), ssaTaintDatabaseToAbstractTaint(graph, callArg.GetTaints()), ssaTaintServiceToAbstractTrace(graph, callArg.GetTaints()))
		edge.AddArgument(arg)
	}

	// create fields if they do not exist yet
	registerDatabaseFields(graph, edge.GetArguments())

	// propagate taints to databases (forward): args (from) >>> params (to)
	for i, toParam := range toNode.GetParams() {
		fromArg := edge.GetArgumentAt(i)
		abstractgraphtainter.MergeTaints(toParam, fromArg.GetPrimaryTaints(), nil, abstractgraphtainter.MERGE_MODE_PARSE, "", false)
	}

	graph.AddEdge(edge)
	graph.IncrDBAccesses()
}

func parseMethodCall(graph *abstractgraph.AbstractCallGraph, node *abstractgraph.AbstractNode, fromSSAGraph *ssagraph.SSAGraph, methodCall *ssagraph.MethodCall, funcGraphs map[string]*ssagraph.SSAGraph) {
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

func registerDatabaseFields(graph *abstractgraph.AbstractCallGraph, args []*abstractgraph.AbstractObject) {
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
