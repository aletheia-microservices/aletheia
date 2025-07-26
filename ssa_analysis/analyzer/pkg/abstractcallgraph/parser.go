package abstractcallgraph

import (
	"fmt"
	"log"
	"slices"

	"analyzer/pkg/ssagraph"
)

func (graph *AbstractCallGraph) ssaTaintToAbstractTaint(ssaTaintsMap map[string][]*ssagraph.SSATaint) map[string][]*AbstractTaint {
	abstractTaintsMap := make(map[string][]*AbstractTaint, len(ssaTaintsMap))
	for objPath, ssaTaints := range ssaTaintsMap {
		abstractTaints := make([]*AbstractTaint, len(ssaTaints))
		for i, ssaTaint := range ssaTaints {
			dbPath := ssaTaint.GetDbCall().GetDatabasePath()
			dbNode := graph.GetNodeByNameIfExists(dbPath)
			if dbNode == nil {
				dbNode = NewAbstractNode(dbPath, NODE_DATABASE)
				graph.AddNode(dbPath, dbNode)
			}
			
			abstractTaints[i] = NewAbstractTaint(ssaTaint.GetDbField(), ssaTaint.GetDbCall().GetID())
		}
		abstractTaintsMap[objPath] = abstractTaints
	}
	return abstractTaintsMap
}

func (graph *AbstractCallGraph) Parse(entryPoints []string, funcGraphs map[string]*ssagraph.SSAGraph) {
	for _, ssaGraph := range funcGraphs {
		if !slices.Contains(entryPoints, ssaGraph.GetFunctionShortPath()) {
			continue
		}

		fromService := ssaGraph.GetServiceName()
		fromNode := graph.GetNodeByNameIfExists(fromService)
		if fromNode == nil {
			fromNode = NewAbstractNode(fromService, NODE_SERVICE)
			graph.AddNode(fromService, fromNode)
		}

		if ssaGraph.HasServiceCalls() {
			fmt.Printf("[ABSTRACTGRAPH] [%s] found function (%s) with service calls\n", ssaGraph.GetServiceName(), ssaGraph.GetFunctionShortPath())
			for _, call := range ssaGraph.GetServiceCalls() {
				toService := call.GetService()
				toNode := graph.GetNodeByNameIfExists(toService)
				if toNode == nil {
					toNode = NewAbstractNode(toService, NODE_SERVICE)
					graph.AddNode(toService, toNode)
				}

				edge := NewAbstractEdge(call.GetID(), call.GetMethod(), fromNode, toNode, EDGE_SERVICE_RPC)

				toFuncGraph := funcGraphs[call.GetFuncShortPath()]
				fmt.Printf("[ABSTRACTGRAPH] ssaGraph = (%s) // toFuncGraph = (%s)\n", ssaGraph.GetFunctionShortPath(), toFuncGraph.GetFunctionShortPath())
				if toFuncGraph == nil {
					log.Fatalf("could not find ssa graph for short func path (%s)", call.GetFuncShortPath())
				}

				methodParams := toFuncGraph.GetParametersExceptMemberAndContext()
				for i, callArg := range call.GetArguments() {
					arg := NewAbstractArgument(graph.ssaTaintToAbstractTaint(callArg.GetTaints()), graph.ssaTaintToAbstractTaint(methodParams[i].GetTaints()), callArg.String())
					edge.AddCallArgument(arg)
				}

				callArgs := call.GetArguments()
				for i, methodParam := range toFuncGraph.GetParametersExceptMemberAndContext() {
					param := NewAbstractArgument(graph.ssaTaintToAbstractTaint(methodParam.GetTaints()), graph.ssaTaintToAbstractTaint(callArgs[i].GetTaints()), methodParam.String())
					edge.AddMethodParameter(param)
				}

				graph.AddEdge(call.GetID(), edge)
			}
			fmt.Println()
		}
		if ssaGraph.HasDatabaseCalls() {
			fmt.Printf("[ABSTRACTGRAPH] found [%s] function (%s) with database calls\n", ssaGraph.GetServiceName(), ssaGraph.GetFunctionShortPath())
			for _, call := range ssaGraph.GetDatabaseCalls() {
				toDatabasePath := call.GetDatabasePath()
				toNode := graph.GetNodeByNameIfExists(toDatabasePath)
				if toNode == nil {
					toNode = NewAbstractNode(toDatabasePath, NODE_DATABASE)
					graph.AddNode(toDatabasePath, toNode)
				}

				edge := NewAbstractEdge(call.GetID(), call.GetMethod(), fromNode, toNode, EDGE_SERVICE_RPC)

				for _, callArg := range call.GetArguments() {
					arg := NewAbstractArgument(graph.ssaTaintToAbstractTaint(callArg.GetTaints()), nil, callArg.String())
					edge.AddCallArgument(arg)
				}

				graph.AddEdge(call.GetID(), edge)
			}
			fmt.Println()
		}
	}
}
