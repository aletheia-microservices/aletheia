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
				dbNode = NewAbstractNode(dbPath, NODE_DATABASE, "", "")
				graph.AddNode(dbPath, dbNode)
			}

			abstractTaints[i] = NewAbstractTaint(ssaTaint.GetDbField(), ssaTaint.GetDbCall().GetID(), true)
		}
		abstractTaintsMap[objPath] = abstractTaints
	}
	return abstractTaintsMap
}

func (graph *AbstractCallGraph) Parse(entryPoints []string, funcGraphs map[string]*ssagraph.SSAGraph) {
	// dummy node
	clientNode := NewAbstractNode("client", NODE_CLIENT, "", "")
	graph.AddNode("client", clientNode)

	for _, ssaGraph := range funcGraphs {
		funcShortPath := ssaGraph.GetFunctionShortPath()
		if !slices.Contains(entryPoints, funcShortPath) {
			continue
		}

		name := ssaGraph.GetServiceWithMethod()
		node := graph.GetNodeByNameIfExists(name)
		if node == nil {
			node = NewAbstractNode(name, NODE_SERVICE, ssaGraph.GetService(), ssaGraph.GetMethodName())
			graph.AddNode(name, node)

			fmt.Printf("[ABSTRACTGRAPH] creating node with (%d) params: %s\n", len(ssaGraph.GetFuncParametersExceptMemberAndContext()), node)
			for _, funcParam := range ssaGraph.GetFuncParametersExceptMemberAndContext() {
				param := NewAbstractObject(funcParam.GetName(), graph.ssaTaintToAbstractTaint(funcParam.GetTaints()))
				node.AddParam(param)
			}
		}

		// 1. build dummy edges for entrypoints
		edge := NewAbstractEdge(funcShortPath, utils.ExtractMethodNameFromShortFunctionPath(funcShortPath), clientNode, node, EDGE_SERVICE_ENTRYPOINT)
		for _, funcParam := range ssaGraph.GetFuncParametersExceptMemberAndContext() {
			arg := NewAbstractObject(funcParam.GetName(), make(map[string][]*AbstractTaint))
			edge.AddArgument(arg)
		}
		graph.AddEdge(edge.GetID(), edge)

		// 2. build edges for service/database RPCs/calls
		if ssaGraph.HasServiceCalls() {
			fmt.Printf("[ABSTRACTGRAPH] [%s] found function (%s) with service calls\n", ssaGraph.GetService(), funcShortPath)
			for _, call := range ssaGraph.GetServiceCalls() {
				toName := call.GetServiceWithMethod()
				toNode := graph.GetNodeByNameIfExists(toName)

				toSSAGraph := funcGraphs[call.GetFuncShortPath()]
				if toSSAGraph == nil {
					log.Fatalf("could not find ssa graph for short func path (%s)", call.GetFuncShortPath())
				}

				// create node for the first time
				if toNode == nil {
					toNode = NewAbstractNode(toName, NODE_SERVICE, call.GetService(), call.GetMethod())
					graph.AddNode(toName, toNode)

					fmt.Printf("[ABSTRACTGRAPH] creating toNode with (%d) params: %s\n", len(toSSAGraph.GetFuncParametersExceptMemberAndContext()), toNode)
					for _, funcParam := range toSSAGraph.GetFuncParametersExceptMemberAndContext() {
						param := NewAbstractObject(funcParam.GetName(), graph.ssaTaintToAbstractTaint(funcParam.GetTaints()))
						toNode.AddParam(param)
					}
				}

				edge := NewAbstractEdge(call.GetID(), call.GetMethod(), node, toNode, EDGE_SERVICE_RPC)

				// create call arguments
				for _, callArg := range call.GetArguments() {
					arg := NewAbstractObject(callArg.GetName(), graph.ssaTaintToAbstractTaint(callArg.GetTaints()))
					edge.AddArgument(arg)
				}

				// propagate taints (indirect): fromArgs >>> toParams
				for i, toParam := range toNode.GetParams() {
					fromArg := edge.GetArgumentAt(i)
					toParam.AddSecondaryTaints(fromArg.GetPrimaryTaints())
				}

				// propagate taints (indirect): fromArgs <<< toParams
				for i, fromArg := range edge.GetArguments() {
					toParam := toNode.GetParamAt(i)
					fromArg.AddSecondaryTaints(toParam.GetPrimaryTaints())
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
					toNode = NewAbstractNode(toDatabasePath, NODE_DATABASE, "", "")
					graph.AddNode(toDatabasePath, toNode)
				}

				edge := NewAbstractEdge(call.GetID(), call.GetMethod(), node, toNode, EDGE_DATABASE_CALL)

				for _, callArg := range call.GetArguments() {
					arg := NewAbstractObject(callArg.GetName(), graph.ssaTaintToAbstractTaint(callArg.GetTaints()))
					edge.AddArgument(arg)
				}

				graph.AddEdge(edge.GetID(), edge)
			}
			fmt.Println()
		}
	}
}
