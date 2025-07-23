package graph

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type Graph struct {
	nodes []*Node
	edges []*Edge
	defs  map[string]*Node
}

func NewGraph() *Graph {
	return &Graph{
		defs: make(map[string]*Node),
	}
}

func (graph *Graph) addEdge(edge *Edge) {
	graph.edges = append(graph.edges, edge)
}

func (graph *Graph) GetNodes() []*Node {
	return graph.nodes
}

func (graph *Graph) GetEdges() []*Edge {
	return graph.edges
}

func (graph *Graph) GetEdgesFromNode(node *Node) []*Edge {
	var edges []*Edge
	for _, edge := range graph.edges {
		if edge.from == node {
			edges = append(edges, edge)
		}
	}
	return edges
}

func (graph *Graph) GetEdgesToNode(node *Node) []*Edge {
	var edges []*Edge
	for _, edge := range graph.edges {
		if edge.to == node {
			edges = append(edges, edge)
		}
	}
	return edges
}

func (graph *Graph) SortNodes() {
	sort.Slice(graph.nodes, func(i, j int) bool {
		/* ni, err1 := strconv.Atoi(strings.TrimPrefix(graph.nodes[i].name, "t"))
		nj, err2 := strconv.Atoi(strings.TrimPrefix(graph.nodes[j].name, "t"))
		if err1 != nil || err2 != nil {
			return graph.nodes[i].name < graph.nodes[j].name
		}
		return ni < nj */
		return graph.nodes[i].id < graph.nodes[j].id
	})
}

func (graph *Graph) GetNodeByName(name string) *Node {
	if node, exists := graph.defs[name]; exists {
		return node
	}
	log.Fatalf("node with name (%s) not found in graph defs: %v\n", name, graph.defs)
	return nil
}

func (graph *Graph) GetNodeByNameIfExists(name string) (*Node, bool) {
	node, exists := graph.defs[name]
	return node, exists
}

func (graph *Graph) CreateAndAddNewEdge(from *Node, to *Node, edgeType EdgeType, index int, param string) (*Edge, bool) {
	// 1st is for sanity check; 2nd is for nodes obtained from *ssa.Const
	if from == nil || to == nil {
		return nil, false
	}
	for _, edge := range graph.GetEdgesFromNode(from) {
		if edge.to == to /* && edge.edgeType == edgeType */ {
			fmt.Printf("[GRAPH] [1] found existing edge with type: %v\n", edge.edgeType)
			return edge, false
		}
	}
	for _, edge := range graph.GetEdgesToNode(to) {
		if edge.from == from /* && edge.edgeType == edgeType */ {
			fmt.Printf("[GRAPH] [2] found existing edge with type: %v\n", edge.edgeType)
			return edge, false
		}
	}
	edge := &Edge{
		from:     from,
		to:       to,
		edgeType: edgeType,
		index:    index,
		param:    param,
	}
	graph.addEdge(edge)
	return edge, true
}

func (graph *Graph) WriteToDOTFile(appname string, fn string) error {
	filename := fmt.Sprintf("output/%s/graphs/%s.dot", appname, fn)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "digraph G {")
	fmt.Fprintln(file, "\trankdir=TD;")

	for _, node := range graph.nodes {
		str := node.String()
		if node.IsTainted() {
			str += "\n\n==== tainted ====\n" + node.taintString()
		}
		label := strings.ReplaceAll(str, `"`, `\"`)
		nodecolor := node.colorForSSA()

		shape := "ellipse"
		if node.IsTainted() {
			shape = "box"
		}

		color := "black"
		if nodecolor != "" {
			color = nodecolor
		}

		fmt.Fprintf(file, "\tN_%s [label=\"%s\", style=bold, shape=%s, color=\"%s\"];\n", node.id, label, shape, color)
	}

	for _, edge := range graph.edges {
		if edge.edgeType == EDGE_POINTS_TO {
			path := strings.ReplaceAll(edge.path, `"`, `\"`)
			fmt.Fprintf(file, "\tN_%s -> N_%s [label=\"%s\", style=dashed, color=blue];\n", edge.from.id, edge.to.id, path)
		} else if edge.from != nil && edge.to != nil {
			fmt.Fprintf(file, "\tN_%s -> N_%s;\n", edge.from.id, edge.to.id)
		}
	}

	fmt.Fprintln(file, "}")
	return nil
}
