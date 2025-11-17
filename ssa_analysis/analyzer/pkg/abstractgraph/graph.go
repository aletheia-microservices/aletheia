package abstractgraph

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"analyzer/pkg/app"
)

type AbstractCallGraph struct {
	app *app.App
	// can either be service (key is the service name) or database (key is the database path)
	nodes map[string]*AbstractNode
	// key is the id of the ssa instr name for the svc or db call on the callee side
	edges []*AbstractEdge
}

func NewAbstractCallGraph(app *app.App) *AbstractCallGraph {
	return &AbstractCallGraph{
		app:   app,
		nodes: make(map[string]*AbstractNode),
	}
}

func (graph *AbstractCallGraph) GetApp() *app.App {
	return graph.app
}

func (graph *AbstractCallGraph) AddNode(name string, node *AbstractNode) {
	if _, ok := graph.nodes[name]; ok {
		log.Fatalf("node with name (%s) already exists in graph: %v", name, graph)
	}
	graph.nodes[name] = node
}

func (graph *AbstractCallGraph) AddEdge(edge *AbstractEdge) {
	fmt.Printf("[ABSTRACTGRAPH] added new edge: %s\n", edge.String())
	graph.edges = append(graph.edges, edge)
}

func (graph *AbstractCallGraph) GetNodeByName(name string) *AbstractNode {
	if node, ok := graph.nodes[name]; ok {
		return node
	}
	log.Panicf("node with name (%s) not found in graph: %v", name, graph)
	return nil
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

func (graph *AbstractCallGraph) WriteToDOTFile(appname string, detailed bool) error {
	var filename string

	if detailed {
		filename = fmt.Sprintf("output/%s/abstractcallgraph_detailed.dot", appname)
	} else {
		filename = fmt.Sprintf("output/%s/abstractcallgraph.dot", appname)
	}
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

	var sortedNodes []*AbstractNode
	for _, node := range graph.GetNodes() {
		sortedNodes = append(sortedNodes, node)
	}

	sort.Slice(sortedNodes, func(i, j int) bool {
		return sortedNodes[i].name < sortedNodes[j].name
	})

	for _, node := range sortedNodes {
		nodeID := strings.ReplaceAll(node.GetName(), ".", "_")
		label := node.GetName()

		if detailed {
			for i, param := range node.GetParams() {
				if param.IsTainted() || param.IsTraced() {
					label += fmt.Sprintf("\n\n==== param %d (%s) tainted ====\n%s", i, param.name, param.Annotations())
				}

			}
			for i, ret := range node.GetReturns() {
				if ret.IsTainted() || ret.IsTraced() {
					label += fmt.Sprintf("\n\n==== ret %d (%s) tainted ====\n%s", i, ret.name, ret.Annotations())
				}

			}
			label = strings.ReplaceAll(label, `"`, `\"`)
		}

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

		label := edge.GetMethod() + "()"

		if detailed {
			for i, arg := range edge.GetArguments() {
				if arg.IsTainted() || arg.IsTraced() {
					label += fmt.Sprintf("\n\n==== arg %d (%s) tainted ====\n%s", i, arg.GetName(), arg.Annotations())
				}

			}
			for i, ret := range edge.GetReturns() {
				if ret.IsTainted() || ret.IsTraced() {
					label += fmt.Sprintf("\n\n==== ret %d (%s) tainted ====\n%s", i, ret.GetName(), ret.Annotations())
				}
			}
			label = strings.ReplaceAll(label, `"`, `\"`)
		}

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
		fmt.Fprintf(file, "\tN_%s -> N_%s [label=\"%s\", color=%s];\n", fromNodeID, toNodeID, label, color)
	}

	/* for _, node := range graph.GetNodes() {
		fmt.Printf("on node: %v\n", node)
		for _, edge := range graph.GetEdgesFromNode(node) {
			fmt.Printf("\t on edge: %v\n", edge)
		}
	} */

	fmt.Fprintln(file, "}")
	return nil
}
