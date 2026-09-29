package abstractgraphparser

import (
	"github.com/sirupsen/logrus"

	"github.com/aletheia-microservices/aletheia/internal/analysis/service-level/ssagraph"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph"
)

func ssaTaintDatabaseToAbstractTaint(graph *abstractgraph.AbstractCallGraph, ssaTaintsMap map[string][]*ssagraph.SSATaint) map[string][]*abstractgraph.AbstractTaint {
	abstractTaintsMap := make(map[string][]*abstractgraph.AbstractTaint, 0)
	for objPath, ssaTaints := range ssaTaintsMap {
		var abstractTaints []*abstractgraph.AbstractTaint
		for _, ssaTaint := range ssaTaints {
			// only database taints become abstract taints (service taints become traces)
			if ssaTaint.IsDatabaseTaint() {
				getOrCreateDatabaseNode(graph, ssaTaint.GetDatabaseCall().GetDatabaseName(), ssaTaint.GetDatabaseCall().GetSchemaName())
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
		// skip paths without database taints
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
			// only service taints become abstract traces (database taints become taints)
			if ssaTaint.IsServiceTaint() {
				trace := abstractgraph.NewAbstractTrace(
					ssaTaint.GetT(),
					ssaTaint.GetServicePath(),
					ssaTaint.GetServiceCall().GetID(),
				)
				abstractTraces = append(abstractTraces, trace)
			}
		}
		// skip paths without service taints
		if abstractTraces != nil {
			abstractTaintsMap[objPath] = abstractTraces
		}
	}
	return abstractTaintsMap
}

func ssaObject(graph *abstractgraph.AbstractCallGraph, name string, ssaTaintsMap map[string][]*ssagraph.SSATaint) *abstractgraph.AbstractObject {
	return abstractgraph.NewAbstractObject(name, ssaTaintDatabaseToAbstractTaint(graph, ssaTaintsMap), ssaTaintServiceToAbstractTrace(graph, ssaTaintsMap))
}

func ssaObjects(graph *abstractgraph.AbstractCallGraph, ssaNodes []*ssagraph.SSANode) []*abstractgraph.AbstractObject {
	var objs []*abstractgraph.AbstractObject
	for _, ssaNode := range ssaNodes {
		objs = append(objs, ssaObject(graph, ssaNode.GetName(), ssaNode.GetTaints()))
	}
	return objs
}

func ssaParams(graph *abstractgraph.AbstractCallGraph, ssaGraph *ssagraph.SSAGraph) func() []*abstractgraph.AbstractObject {
	return func() []*abstractgraph.AbstractObject {
		return ssaObjects(graph, ssaGraph.GetFuncParametersExceptMemberAndContext())
	}
}

func Parse(graph *abstractgraph.AbstractCallGraph, funcshortpath string, entrypoint bool, funcGraphs map[string]*ssagraph.SSAGraph) {
	ssaGraph := funcGraphs[funcshortpath]
	// get or create the node of the function, and skip it if it was already parsed
	node, ok := visitFunction(graph, ssaGraph.GetService(), ssaGraph.GetMethodName(), funcshortpath, entrypoint, ssaParams(graph, ssaGraph))
	if !ok {
		return
	}

	// finalize parsing: add the returns of the function, merging the objects of all its return statements
	// (return values have no name, so each object is named after its type)
	var retsLst [][]*abstractgraph.AbstractObject
	for _, rets := range ssaGraph.GetReturnsLst() {
		var retsObjs []*abstractgraph.AbstractObject
		for _, ret := range rets {
			retsObjs = append(retsObjs, ssaObject(graph, ret.GetValue().Type().String(), ret.GetTaints()))
		}
		retsLst = append(retsLst, retsObjs)
	}
	addReturns(node, retsLst)

	// add an edge for each call of the function and recursively parse the functions called by RPC
	parseCalls(graph, node, ssaGraph, ssaGraph, funcGraphs)
}

// parseCalls parses the calls of ssaGraph, which is either fromSSAGraph or one of its inlined graphs
func parseCalls(graph *abstractgraph.AbstractCallGraph, node *abstractgraph.AbstractNode, fromSSAGraph *ssagraph.SSAGraph, ssaGraph *ssagraph.SSAGraph, funcGraphs map[string]*ssagraph.SSAGraph) {
	for _, call := range ssaGraph.GetAllCalls() {
		// RPC to another service: add a service edge
		if serviceCall, ok := call.(*ssagraph.ServiceCall); ok {
			parseServiceCall(graph, node, serviceCall, funcGraphs)
		}

		// database call: add a database edge
		if databaseCall, ok := call.(*ssagraph.DatabaseCall); ok {
			parseDatabaseCall(graph, node, databaseCall)
		}

		// internal method call: parse the calls of its inlined graph as calls of node
		if methodCall, ok := call.(*ssagraph.MethodCall); ok {
			parseMethodCall(graph, node, fromSSAGraph, methodCall, funcGraphs)
		}
	}
	// then, recursively parse the functions called by RPC
	for _, call := range ssaGraph.GetServiceCalls() {
		Parse(graph, call.GetFuncShortPath(), false, funcGraphs)
	}
}

func parseServiceCall(graph *abstractgraph.AbstractCallGraph, node *abstractgraph.AbstractNode, serviceCall *ssagraph.ServiceCall, funcGraphs map[string]*ssagraph.SSAGraph) {
	logrus.WithField("node", node.String()).Tracef("[ABSTRACTGRAPH] found service call: %s\n", serviceCall.String())
	// the graph of the callee is needed to build the params of its node
	toSSAGraph := funcGraphs[serviceCall.GetFuncShortPath()]
	if toSSAGraph == nil {
		logrus.Fatalf("could not find ssa graph for short func path (%s)", serviceCall.GetFuncShortPath())
	}
	toNode := getOrCreateServiceNode(graph, serviceCall.GetService(), serviceCall.GetMethod(), ssaParams(graph, toSSAGraph))
	// add the edge with the arguments and returns of the call
	args := ssaObjects(graph, serviceCall.GetArguments())
	rets := ssaObjects(graph, serviceCall.GetReturns())
	addServiceCallEdge(graph, node, toNode, serviceCall.GetT(), serviceCall.GetID(), serviceCall.GetMethod(), args, rets)
}

func parseDatabaseCall(graph *abstractgraph.AbstractCallGraph, node *abstractgraph.AbstractNode, databaseCall *ssagraph.DatabaseCall) {
	toNode := getOrCreateDatabaseNode(graph, databaseCall.GetDatabaseName(), databaseCall.GetSchemaName())
	// add the edge with the arguments of the call (which also propagates their taints to the database)
	args := ssaObjects(graph, databaseCall.GetArguments())
	addDatabaseCallEdge(graph, node, toNode, databaseCall.GetScopedT(), databaseCall.GetID(), databaseCall.GetMethod(), databaseCall.GetOpType(), args)
}

func parseMethodCall(graph *abstractgraph.AbstractCallGraph, node *abstractgraph.AbstractNode, fromSSAGraph *ssagraph.SSAGraph, methodCall *ssagraph.MethodCall, funcGraphs map[string]*ssagraph.SSAGraph) {
	toSSAGraph := fromSSAGraph.GetInlinedGraphForMethodCallIfExists(methodCall)
	if toSSAGraph == nil {
		// should never happen
		return
	}
	// the method belongs to the same service, so its calls are added as calls of node
	parseCalls(graph, node, fromSSAGraph, toSSAGraph, funcGraphs)
}
