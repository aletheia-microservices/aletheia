package ssa_graph

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type ServiceCall struct {
	node    *SSANode
	args    []*SSANode
	service string
	method  string
}

func (call *ServiceCall) GetNode() *SSANode {
	return call.node
}

func (call *ServiceCall) GetArguments() []*SSANode {
	return call.args
}

type DatabaseCall struct {
	node *SSANode
	args []*SSANode
}

func (call *DatabaseCall) GetNode() *SSANode {
	return call.node
}

func (call *DatabaseCall) GetArguments() []*SSANode {
	return call.args
}

type SSAGraph struct {
	pkg      string
	fn       string
	nodes    []*SSANode
	edges    []*SSAEdge
	defs     map[string]*SSANode
	svcCalls []*ServiceCall
	dbCalls  []*DatabaseCall
}

func NewGraph(pkg string, fn string) *SSAGraph {
	return &SSAGraph{
		defs: make(map[string]*SSANode),
		pkg:  pkg,
		fn:   fn,
	}
}

func (graph *SSAGraph) GetPackageName() string {
	return graph.pkg
}

func (graph *SSAGraph) GetFunctionShortPath() string {
	return graph.fn
}

func (graph *SSAGraph) addEdge(edge *SSAEdge) {
	graph.edges = append(graph.edges, edge)
}

func (graph *SSAGraph) GetNodes() []*SSANode {
	return graph.nodes
}

func (graph *SSAGraph) GetEdges() []*SSAEdge {
	return graph.edges
}

func (graph *SSAGraph) AddServiceCall(node *SSANode, args []*SSANode, service string, method string) {
	graph.svcCalls = append(graph.svcCalls, &ServiceCall{
		node: node,
		args: args,
		service: service,
		method: method,
	})
}

func (graph *SSAGraph) HasServiceCalls() bool {
	return len(graph.svcCalls) > 0
}

func (graph *SSAGraph) GetServiceCalls() []*ServiceCall {
	return graph.svcCalls
}

func (graph *SSAGraph) AddDatabaseCall(node *SSANode, args []*SSANode) {
	graph.dbCalls = append(graph.dbCalls, &DatabaseCall{
		node: node,
		args: args,
	})
}

func (graph *SSAGraph) HasDatabaseCalls() bool {
	return len(graph.dbCalls) > 0
}

func (graph *SSAGraph) GetDatabaseCalls() []*DatabaseCall {
	return graph.dbCalls
}

func (graph *SSAGraph) GetEdgesFromNode(node *SSANode) []*SSAEdge {
	var edges []*SSAEdge
	for _, edge := range graph.edges {
		if edge.from == node {
			edges = append(edges, edge)
		}
	}
	return edges
}

func (graph *SSAGraph) GetEdgesToNode(node *SSANode) []*SSAEdge {
	var edges []*SSAEdge
	for _, edge := range graph.edges {
		if edge.to == node {
			edges = append(edges, edge)
		}
	}
	return edges
}

func (graph *SSAGraph) SortNodes() {
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

func (graph *SSAGraph) GetNodeByName(name string) *SSANode {
	if node, exists := graph.defs[name]; exists {
		return node
	}
	log.Fatalf("node with name (%s) not found in graph defs: %v\n", name, graph.defs)
	return nil
}

func (graph *SSAGraph) GetNodeByNameIfExists(name string) (*SSANode, bool) {
	node, exists := graph.defs[name]
	return node, exists
}

func (graph *SSAGraph) CreateAndAddNewEdge(from *SSANode, to *SSANode, edgeType EdgeType, index int, param string) (*SSAEdge, bool) {
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
	edge := &SSAEdge{
		from:     from,
		to:       to,
		edgeType: edgeType,
		index:    index,
		param:    param,
	}
	graph.addEdge(edge)
	return edge, true
}

func (graph *SSAGraph) WriteToDOTFile(appname string, fn string) error {
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
