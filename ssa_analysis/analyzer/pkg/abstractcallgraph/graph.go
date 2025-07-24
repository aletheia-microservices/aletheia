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
	graph.edges = append(graph.edges, edge)
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
		
		fnShortPath := ssaGraph.GetFunctionShortPath()
		node := graph.GetNodeByNameIfExists(fnShortPath)
		if node == nil {
			node = NewAbstractNode(fnShortPath)
			graph.AddNode(fnShortPath, node)
		}

		if ssaGraph.HasServiceCalls() {
			fmt.Printf("[INFO] function (%s) has service calls:\n", ssaGraph.GetFunctionShortPath())
			for _, call := range ssaGraph.GetServiceCalls() {
				fmt.Printf("\t - call %s\n", call.GetNode().String())

				edge := NewAbstractEdge(call.GetNode().String(), node, nil, EDGE_SERVICE_RPC)
				graph.AddEdge(edge)

				for _, arg := range call.GetArguments() {
					fmt.Printf("\t\t - arg %s\n", arg.String())
					for obj, dbfields := range arg.GetTaints() {
						fmt.Printf("\t\t\t - taint %s @ %s\n", obj, strings.Join(dbfields, ", "))
					}
				}
			}
			fmt.Println()
		}
		if ssaGraph.HasDatabaseCalls() {
			fmt.Printf("[INFO] function (%s) has database calls:\n", ssaGraph.GetFunctionShortPath())
			for _, call := range ssaGraph.GetDatabaseCalls() {
				fmt.Printf("\t - %s\n", call.GetNode().String())
				for _, arg := range call.GetArguments() {
					fmt.Printf("\t\t - arg %s\n", arg.String())
					for obj, dbfields := range arg.GetTaints() {
						fmt.Printf("\t\t\t - taint %s @ %s\n", obj, strings.Join(dbfields, ", "))
					}
				}
			}
			fmt.Println()
		}
	}
}
