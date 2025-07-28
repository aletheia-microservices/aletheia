package abstractcallgraph

import (
	"fmt"
	"log"
	"slices"

	"analyzer/pkg/ssagraph"
	"analyzer/pkg/utils"
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
	// dummy node
	clientNode := NewAbstractNode("client", NODE_CLIENT)
	graph.AddNode("client", clientNode)

	for _, ssaGraph := range funcGraphs {
		funcShortPath := ssaGraph.GetFunctionShortPath()
		if !slices.Contains(entryPoints, funcShortPath) {
			continue
		}

		service := ssaGraph.GetService()
		node := graph.GetNodeByNameIfExists(service)
		if node == nil {
			node = NewAbstractNode(service, NODE_SERVICE)
			graph.AddNode(service, node)
		}

		// 1. build edges for entrypoints
		edge := NewAbstractEdge(funcShortPath, utils.ExtractMethodNameFromShortFunctionPath(funcShortPath), clientNode, node, EDGE_SERVICE_ENTRYPOINT)
		for _, methodParam := range ssaGraph.GetFuncParametersExceptMemberAndContext() {
			param := NewAbstractArgument(graph.ssaTaintToAbstractTaint(methodParam.GetTaints()), make(map[string][]*AbstractTaint), methodParam.String())
			edge.AddMethodParameter(param)
		}
		graph.AddEdge(edge.GetID(), edge)

		// 2. build edges for service/database RPCs/calls
		if ssaGraph.HasServiceCalls() {
			fmt.Printf("[ABSTRACTGRAPH] [%s] found function (%s) with service calls\n", ssaGraph.GetService(), funcShortPath)
			for _, call := range ssaGraph.GetServiceCalls() {
				toService := call.GetService()
				toNode := graph.GetNodeByNameIfExists(toService)
				if toNode == nil {
					toNode = NewAbstractNode(toService, NODE_SERVICE)
					graph.AddNode(toService, toNode)
				}

				edge := NewAbstractEdge(call.GetID(), call.GetMethod(), node, toNode, EDGE_SERVICE_RPC)

				toSSAGraph := funcGraphs[call.GetFuncShortPath()]
				fmt.Printf("[ABSTRACTGRAPH] ssaGraph = (%s) // toFuncGraph = (%s)\n", funcShortPath, toSSAGraph.GetFunctionShortPath())
				if toSSAGraph == nil {
					log.Fatalf("could not find ssa graph for short func path (%s)", call.GetFuncShortPath())
				}

				methodParams := toSSAGraph.GetFuncParametersExceptMemberAndContext()
				for i, callArg := range call.GetArguments() {
					arg := NewAbstractArgument(graph.ssaTaintToAbstractTaint(callArg.GetTaints()), graph.ssaTaintToAbstractTaint(methodParams[i].GetTaints()), callArg.String())
					edge.AddCallArgument(arg)
				}

				callArgs := call.GetArguments()
				for i, methodParam := range toSSAGraph.GetFuncParametersExceptMemberAndContext() {
					param := NewAbstractArgument(graph.ssaTaintToAbstractTaint(methodParam.GetTaints()), graph.ssaTaintToAbstractTaint(callArgs[i].GetTaints()), methodParam.String())
					edge.AddMethodParameter(param)
				}

				graph.AddEdge(edge.GetID(), edge)
			}
			fmt.Println()
		}
		if ssaGraph.HasDatabaseCalls() {
			fmt.Printf("[ABSTRACTGRAPH] found [%s] function (%s) with database calls\n", ssaGraph.GetService(), funcShortPath)
			for _, call := range ssaGraph.GetDatabaseCalls() {
				toDatabasePath := call.GetDatabasePath()
				toNode := graph.GetNodeByNameIfExists(toDatabasePath)
				if toNode == nil {
					toNode = NewAbstractNode(toDatabasePath, NODE_DATABASE)
					graph.AddNode(toDatabasePath, toNode)
				}

				edge := NewAbstractEdge(call.GetID(), call.GetMethod(), node, toNode, EDGE_DATABASE_CALL)

				for _, callArg := range call.GetArguments() {
					arg := NewAbstractArgument(graph.ssaTaintToAbstractTaint(callArg.GetTaints()), nil, callArg.String())
					edge.AddCallArgument(arg)
				}

				graph.AddEdge(edge.GetID(), edge)
			}
			fmt.Println()
		}
	}
}
