package abstractcallgraph

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"analyzer/pkg/ssa_graph"
)

type AbstractCallGraph struct {
	nodes map[string]*AbstractNode
	edges []*AbstractEdge
}

func NewAbstractGraph() *AbstractCallGraph {
	return &AbstractCallGraph{
		nodes: make(map[string]*AbstractNode),
	}
}

func (graph *AbstractCallGraph) AddNode(name string, node *AbstractNode) {
	if _, ok := graph.nodes[name]; ok {
		log.Fatalf("node with name (%s) already exists in graph: %v", name, graph)
	}
	graph.nodes[name] = node
}

func (graph *AbstractCallGraph) AddEdge(edge *AbstractEdge) {
	fmt.Printf("[ABSTRACT GRAPH] added new edge: %s\n", edge.String())
	graph.edges = append(graph.edges, edge)

	for i, arg := range edge.args {
		fmt.Printf("\t\t - ARG #%d: %s\n", i, arg.SSAString())
		for obj, dbfields := range arg.GetTaints() {
			fmt.Printf("\t\t\t - TAINT: %s @ %s\n", obj, strings.Join(dbfields, ", "))
		}
	}
}

func (graph *AbstractCallGraph) GetNodeByNameIfExists(name string) *AbstractNode {
	if node, ok := graph.nodes[name]; ok {
		return node
	}
	return nil
}

func (graph *AbstractCallGraph) GetNodes() map[string]*AbstractNode {
	return graph.nodes
}

func (graph *AbstractCallGraph) GetEdges() []*AbstractEdge {
	return graph.edges
}

func (graph *AbstractCallGraph) GetEdgesFromNode(node *AbstractNode) []*AbstractEdge {
	var edges []*AbstractEdge
	for _, edge := range graph.edges {
		if edge.from == node {
			edges = append(edges, edge)
		}
	}
	return edges
}

func (graph *AbstractCallGraph) GetEdgesToNode(node *AbstractNode) []*AbstractEdge {
	var edges []*AbstractEdge
	for _, edge := range graph.edges {
		if edge.to == node {
			edges = append(edges, edge)
		}
	}
	return edges
}

func (graph *AbstractCallGraph) Init(entryPoints []string, funcGraphs map[string]*ssa_graph.SSAGraph) {
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
			fmt.Printf("[INFO] [%s] found function (%s) with service calls\n", ssaGraph.GetServiceName(), ssaGraph.GetFunctionShortPath())
			for _, call := range ssaGraph.GetServiceCalls() {
				toService := call.GetService()
				toNode := graph.GetNodeByNameIfExists(toService)
				if toNode == nil {
					toNode = NewAbstractNode(toService, NODE_SERVICE)
					graph.AddNode(toService, toNode)
				}

				edge := NewAbstractEdge(call.GetMethod(), fromNode, toNode, EDGE_SERVICE_RPC)
				
				for _, callArg := range call.GetArguments() {
					arg := NewAbstractArgument(callArg.GetTaints(), callArg.String())
					edge.AddArgument(arg)
				}
				
				graph.AddEdge(edge)
			}
			fmt.Println()
		}
		if ssaGraph.HasDatabaseCalls() {
			fmt.Printf("[INFO] found [%s] function (%s) with database calls\n", ssaGraph.GetServiceName(), ssaGraph.GetFunctionShortPath())
			for _, call := range ssaGraph.GetDatabaseCalls() {
				toDatabasePath := call.GetDatabase() + "." + call.GetCollectionOrTopic()
				toNode := graph.GetNodeByNameIfExists(toDatabasePath)
				if toNode == nil {
					toNode = NewAbstractNode(toDatabasePath, NODE_DATABASE)
					graph.AddNode(toDatabasePath, toNode)
				}

				edge := NewAbstractEdge(call.GetMethod(), fromNode, toNode, EDGE_SERVICE_RPC)
				
				for _, callArg := range call.GetArguments() {
					arg := NewAbstractArgument(callArg.GetTaints(), callArg.String())
					edge.AddArgument(arg)
				}
				
				graph.AddEdge(edge)
			}
			fmt.Println()
		}
	}
}
