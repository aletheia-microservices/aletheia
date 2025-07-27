package abstractcallgraph

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type AbstractCallGraph struct {
	// can either be service (key is the service name) or database (key is the database path)
	nodes map[string]*AbstractNode
	// key is the id of the ssa instr name for the svc or db call on the callee side
	edges map[string]*AbstractEdge
}

func NewAbstractCallGraph() *AbstractCallGraph {
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
	fmt.Printf("[ABSTRACTGRAPH] added new edge: %s\n", edge.String())
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

func (graph *AbstractCallGraph) WriteToDOTFile(appname string) error {
	filename := fmt.Sprintf("output/%s/abstractcallgraph.dot", appname)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "digraph G {")
	fmt.Fprintln(file, "\trankdir=TD;")
	fmt.Fprintln(file, "\tranksep=1.5;")
	fmt.Fprintln(file, "\tnodesep=1;")

	services := []string{}
	databases := []string{}
	clients := []string{}
	others := []string{}

	for _, node := range graph.GetNodes() {
		nodeID := strings.ReplaceAll(node.GetName(), ".", "_")
		label := node.GetName()
		var color, shape string

		switch node.GetNodeType() {
		case NODE_SERVICE:
			color = "blue"
			shape = "box"
			services = append(services, fmt.Sprintf("\tN_%s [label=\"%s\", style=bold, shape=%s, color=\"%s\"];", nodeID, label, shape, color))
		case NODE_DATABASE:
			color = "green"
			shape = "cylinder"
			databases = append(databases, fmt.Sprintf("\tN_%s [label=\"%s\", style=bold, shape=%s, color=\"%s\"];", nodeID, label, shape, color))
		case NODE_CLIENT:
			color = "orange"
			shape = "invhouse"
			clients = append(clients, fmt.Sprintf("\tN_%s [label=\"%s\", style=bold, shape=%s, color=\"%s\"];", nodeID, label, shape, color))
		default:
			color = "black"
			shape = "ellipse"
			others = append(others, fmt.Sprintf("\tN_%s [label=\"%s\", style=bold, shape=%s, color=\"%s\"];", nodeID, label, shape, color))
		}
	}

	//fmt.Fprintln(file, "\tsubgraph cluster_clients {\n\t\tlabel = \"Clients\";")
	fmt.Fprintln(file, "\tsubgraph cluster_clients {\n\t\tstyle=invis;")
	for _, line := range clients {
		fmt.Fprintln(file, "\t"+line)
	}
	fmt.Fprintln(file, "\t}")

	//fmt.Fprintln(file, "\tsubgraph cluster_services {\n\t\tlabel = \"Services\";")
	fmt.Fprintln(file, "\tsubgraph cluster_services {\n\t\tstyle=invis;")
	for _, line := range services {
		fmt.Fprintln(file, "\t"+line)
	}
	fmt.Fprintln(file, "\t}")

	//fmt.Fprintln(file, "\tsubgraph cluster_databases {\n\t\tlabel = \"Databases\";")
	fmt.Fprintln(file, "\tsubgraph cluster_databases {\n\t\tstyle=invis;")
	for _, line := range databases {
		fmt.Fprintln(file, "\t"+line)
	}
	fmt.Fprintln(file, "\t}")

	for _, line := range others {
		fmt.Fprintln(file, "\t"+line)
	}

	for _, edge := range graph.GetEdges() {
		fromNodeID := strings.ReplaceAll(edge.GetFromNode().GetName(), ".", "_")
		toNodeID := strings.ReplaceAll(edge.GetToNode().GetName(), ".", "_")

		var color string
		switch edge.GetEdgeType() {
		case EDGE_SERVICE_RPC:
			color = "blue"
		case EDGE_DATABASE_CALL:
			color = "green"
		case EDGE_SERVICE_ENTRYPOINT:
			color = "orange"
		default:
			color = "black"
		}
		fmt.Fprintf(file, "\tN_%s -> N_%s [label=\"%s\", color=%s];\n", fromNodeID, toNodeID, edge.GetMethod()+"()", color)
	}

	fmt.Fprintln(file, "}")
	return nil
}
