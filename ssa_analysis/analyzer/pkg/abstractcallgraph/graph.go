package abstractcallgraph

import (
	"fmt"
	"log"
)

type AbstractCallGraph struct {
	// can either be service (key is the service name) or database (key is the database path)
	nodes map[string]*AbstractNode
	// key is the id of the ssa instr name for the svc or db call on the callee side
	edges map[string]*AbstractEdge
}

func NewAbstractGraph() *AbstractCallGraph {
	return &AbstractCallGraph{
		nodes: make(map[string]*AbstractNode),
		edges: make(map[string]*AbstractEdge),
	}
}

func (graph *AbstractCallGraph) AddNode(name string, node *AbstractNode) {
	if _, ok := graph.nodes[name]; ok {
		log.Fatalf("node with name (%s) already exists in graph: %v", name, graph)
	}
	graph.nodes[name] = node
}

func taintsListToString(taints []*AbstractTaint) string {
	var str string
	for i, taint := range taints {
		str += taint.String()
		if i < len(taints)-1 {
			str += ", "
		}
	}
	return str
}

func (graph *AbstractCallGraph) AddEdge(id string, edge *AbstractEdge) {
	fmt.Printf("[ABSTRACT GRAPH] added new edge: %s\n", edge.String())
	graph.edges[id] = edge

	for i, arg := range edge.callArgs {
		fmt.Printf("\t\t - CALL ARG #%d: %s\n", i, arg.SSAString())
		for obj, directTaints := range arg.GetDirectTaints() {
			fmt.Printf("\t\t\t - TAINT: %s @ %s\n", obj, taintsListToString(directTaints))
		}
		for obj, indirectTaints := range arg.GetIndirectTaints() {
			fmt.Printf("\t\t\t - TAINT (INDIRECT): %s @ %s\n", obj, taintsListToString(indirectTaints))
		}
	}

	for i, param := range edge.methodParams {
		fmt.Printf("\t\t - METHOD PARAM #%d: %s\n", i, param.SSAString())
		for obj, directTaints := range param.GetDirectTaints() {
			fmt.Printf("\t\t\t - TAINT: %s @ %s\n", obj, taintsListToString(directTaints))
		}
		for obj, indirectTaints := range param.GetIndirectTaints() {
			fmt.Printf("\t\t\t - TAINT (INDIRECT): %s @ %s\n", obj, taintsListToString(indirectTaints))
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

func (graph *AbstractCallGraph) GetEdges() map[string]*AbstractEdge {
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
