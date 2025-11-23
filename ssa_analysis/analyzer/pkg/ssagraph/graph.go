package ssagraph

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"

	"analyzer/pkg/app"
)

type SSAGraph struct {
	app         *app.App
	pkgName     string
	fnShortPath string
	serviceName string
	methodName  string

	nodes []*SSANode
	edges []*SSAEdge
	defs  map[string]*SSANode

	params  []*SSANode
	returns [][]*SSANode // can have multiple return tuples depending on controlflow

	// managed by tainter.go
	svcCalls    []*ServiceCall
	methodCalls []*MethodCall
	dbCalls     []*DatabaseCall
	allCalls    []interface{}

	// managed by combiner.go
	combinedGraphs             []*SSAGraph
	callerT                    string
	combinedGraphsToMethodCall map[*SSAGraph]*MethodCall
	methodCallToCombinedGraphs map[*MethodCall]*SSAGraph
}

func NewGraph(app *app.App, pkg string, fn string, service string, method string) *SSAGraph {
	return &SSAGraph{
		app:         app,
		defs:        make(map[string]*SSANode),
		pkgName:     pkg,
		fnShortPath: fn,
		serviceName: service,
		methodName:  method,
	}
}

func (graph *SSAGraph) Release() {
	graph.nodes = nil
	graph.edges = nil
	graph.defs = nil
	graph.svcCalls = nil
	graph.methodCalls = nil
	graph.dbCalls = nil
	graph.allCalls = nil
	graph.params = nil
	graph.returns = nil
	graph.combinedGraphs = nil
}

func (graph *SSAGraph) SimpleCopy() *SSAGraph {
	var copyNodesMap = make(map[string]*SSANode, 0)
	var copyNodes []*SSANode
	for _, node := range graph.nodes {
		copyNode := node.SimpleCopy()
		copyNodes = append(copyNodes, copyNode)
		copyNodesMap[copyNode.String()] = copyNode
	}

	var newDefs = make(map[string]*SSANode, 0)
	for _, node := range copyNodes {
		newDefs[node.GetName()] = node
	}

	var newEdges []*SSAEdge
	for _, edge := range graph.edges {
		var newFromNode, newToNode *SSANode
		if edge.GetFromNode().inDefs {
			newFromNode = newDefs[edge.GetFromNode().GetName()]
		} else {
			newFromNode = copyNodesMap[edge.GetFromNode().String()]
		}
		if edge.GetToNode().inDefs {
			newToNode = newDefs[edge.GetToNode().GetName()]
		} else {
			newToNode = copyNodesMap[edge.GetToNode().String()]
		}
		newEdge := NewEdge(edge.edgeType, newFromNode, newToNode, edge.index, edge.param)
		newEdges = append(newEdges, newEdge)
	}

	var newParams []*SSANode
	for _, param := range graph.params {
		newNode := newDefs[param.GetName()]
		newParams = append(newParams, newNode)
	}

	var newReturns [][]*SSANode
	for _, retLst := range graph.returns {
		var newReturnsLst []*SSANode
		for _, ret := range retLst {
			newNode := newDefs[ret.GetName()]
			newReturnsLst = append(newReturnsLst, newNode)
		}
		newReturns = append(newReturns, newReturnsLst)
	}

	return &SSAGraph{
		app:         graph.app,
		pkgName:     graph.pkgName,
		fnShortPath: graph.fnShortPath,
		serviceName: graph.fnShortPath,
		methodName:  graph.methodName,
		nodes:       copyNodes,
		edges:       newEdges,
		defs:        newDefs,
		params:      newParams,
		returns:     newReturns,
	}
}

func (graph *SSAGraph) String() string {
	return graph.fnShortPath
}

func (graph *SSAGraph) SetCallerT(t string) {
	graph.callerT = t
}

func (graph *SSAGraph) AddCombinedGraph(toGraph *SSAGraph, methodCall *MethodCall) {
	if graph.combinedGraphsToMethodCall == nil {
		graph.combinedGraphsToMethodCall = make(map[*SSAGraph]*MethodCall)
		graph.methodCallToCombinedGraphs = make(map[*MethodCall]*SSAGraph)
	}
	toGraph.SetCallerT(methodCall.GetT())
	graph.combinedGraphs = append(graph.combinedGraphs, toGraph)
	graph.combinedGraphsToMethodCall[toGraph] = methodCall
	graph.methodCallToCombinedGraphs[methodCall] = toGraph
}

func (graph *SSAGraph) GetAllCombinedGraphs() []*SSAGraph {
	return graph.combinedGraphs
}

func (graph *SSAGraph) GetMethodCallForCombinedGraph(other *SSAGraph) *MethodCall {
	return graph.combinedGraphsToMethodCall[other]
}

func (graph *SSAGraph) GetCombinedGraphForMethodCallIfExists(call *MethodCall) *SSAGraph {
	if m, ok := graph.methodCallToCombinedGraphs[call]; ok {
		return m
	}
	logrus.Tracef("[SSA GRAPH] combined graph not found for call (%s)\n", call.String())
	return nil
}

func (graph *SSAGraph) GetIndexOfParameter(expParam *SSANode) int {
	for i, param := range graph.params {
		if param == expParam {
			return i
		}
	}
	logrus.Fatalf("[SSA GRAPH] could not find parameter (%s) in graph for method (%s)", expParam, graph.GetMethodName())
	return -1
}

func (graph *SSAGraph) GetApp() *app.App {
	return graph.app
}

func (graph *SSAGraph) AddNode(node *SSANode) {
	graph.nodes = append(graph.nodes, node)
}

func (graph *SSAGraph) AddNodeDef(node *SSANode) {
	graph.defs[node.GetName()] = node
}

func (graph *SSAGraph) GetServiceWithMethod() string {
	return graph.serviceName + "." + graph.methodName
}

func (graph *SSAGraph) GetService() string {
	return graph.serviceName
}

func (graph *SSAGraph) GetMethodName() string {
	return graph.methodName
}

func (graph *SSAGraph) GetParamAt(idx int) *SSANode {
	return graph.params[idx]
}

func (graph *SSAGraph) GetParams() []*SSANode {
	return graph.params
}

func (graph *SSAGraph) GetPackageName() string {
	return graph.pkgName
}

func (graph *SSAGraph) GetFunctionShortPath() string {
	return graph.fnShortPath
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

func (graph *SSAGraph) AddCall(call interface{}) {
	graph.allCalls = append(graph.allCalls, call)
}

func (graph *SSAGraph) AddServiceCall(call *ServiceCall) {
	graph.svcCalls = append(graph.svcCalls, call)
}

func (graph *SSAGraph) AddMethodCall(call *MethodCall) {
	graph.methodCalls = append(graph.methodCalls, call)
}

func (graph *SSAGraph) GetAllCalls() []interface{} {
	return graph.allCalls
}

func (graph *SSAGraph) GetServiceCalls() []*ServiceCall {
	return graph.svcCalls
}

func (graph *SSAGraph) AddParameter(param *SSANode) {
	graph.params = append(graph.params, param)
}

func (graph *SSAGraph) AddReturnsToLst(rets []*SSANode) {
	graph.returns = append(graph.returns, rets)
}

func (graph *SSAGraph) GetReturnsLst() [][]*SSANode {
	return graph.returns
}

func (graph *SSAGraph) GetFuncParametersExceptMemberAndContext() []*SSANode {
	// EVAL: fmt.Printf("[SSAGRAPH] filtered func parameters: %v\n", graph.params)
	if len(graph.params) <= 2 {
		return nil
	}
	return graph.params[2:]
}

func (graph *SSAGraph) AddDatabaseCall(call *DatabaseCall) {
	graph.dbCalls = append(graph.dbCalls, call)
}

func (graph *SSAGraph) HasDatabaseCalls() bool {
	return len(graph.dbCalls) > 0
}

func (graph *SSAGraph) GetDatabaseCalls() []*DatabaseCall {
	return graph.dbCalls
}

func (graph *SSAGraph) GetMethodCalls() []*MethodCall {
	return graph.methodCalls
}

func (graph *SSAGraph) GetEdgesTypedTo(node *SSANode, t EdgeType) []*SSAEdge {
	var edges []*SSAEdge
	for _, edge := range graph.edges {
		if edge.GetType() == t && edge.to == node {
			edges = append(edges, edge)
		}
	}
	return edges
}

func (graph *SSAGraph) GetEdgesTypedFrom(node *SSANode, t EdgeType) []*SSAEdge {
	var edges []*SSAEdge
	for _, edge := range graph.edges {
		if edge.GetType() == t && edge.from == node {
			edges = append(edges, edge)
		}
	}
	return edges
}

func (graph *SSAGraph) GetFirstEdgeTypedFrom(node *SSANode, t EdgeType) *SSAEdge {
	for _, edge := range graph.edges {
		if edge.GetType() == t && edge.from == node {
			return edge
		}
	}
	return nil
}

func (graph *SSAGraph) GetFirstEdgeToNode(node *SSANode) *SSAEdge {
	for _, edge := range graph.edges {
		if edge.to == node {
			return edge
		}
	}
	return nil
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

func (graph *SSAGraph) GetAllNodeEdges(node *SSANode) []*SSAEdge {
	var edges []*SSAEdge
	for _, edge := range graph.edges {
		if edge.from == node || edge.to == node {
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

func (graph *SSAGraph) Sort() {
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
	log.Panicf("node with name (%s) not found in graph defs: %v\n", name, graph.defs)
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
			return edge, false
		}
	}
	for _, edge := range graph.GetEdgesToNode(to) {
		if edge.from == from /* && edge.edgeType == edgeType */ {
			return edge, false
		}
	}
	edge := NewEdge(edgeType, from, to, index, param)
	graph.addEdge(edge)
	return edge, true
}

func safeLabel(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	return s
}

func safeID(id string) string {
	// replace anything that's not a letter, number, or underscore with underscore
	re := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	return re.ReplaceAllString(id, "_")
}

func (graph *SSAGraph) WriteToDOTFile(appname string, fn string, tainted bool) {
	stage := "untainted"
	if tainted {
		stage = "tainted"
	}
	dirname := fmt.Sprintf("output/%s/ssagraphs/%s", appname, stage)
	filename := fmt.Sprintf("%s/%s.dot", dirname, fn)

	err := os.MkdirAll(dirname, os.ModePerm)
	if err != nil {
		logrus.Fatalf("error: %s", err.Error())
	}

	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		logrus.Fatalf("error: %s", err.Error())
	}

	fmt.Fprintln(file, "digraph G {")
	fmt.Fprintln(file, "\trankdir=TD;")

	for _, node := range graph.nodes {
		label := safeLabel(node.String())
		if nodeLabels := node.LabelsString(); nodeLabels != "" {
			label += "\n" + nodeLabels
		}
		if node.IsTainted() {
			label += "\n\n==== tainted ====\n" + node.TaintAndTraceString()
		}
		nodecolor := node.colorForSSA()

		shape := "ellipse"
		if node.IsTainted() {
			shape = "box"
		}

		color := "black"
		if nodecolor != "" {
			color = nodecolor
		}

		fmt.Fprintf(file, "\tN_%s [label=\"%s\", style=bold, shape=%s, color=\"%s\"];\n", safeID(node.id), label, shape, color)
	}

	for _, edge := range graph.edges {
		if edge.edgeType == EDGE_POINTS_TO {
			path := strings.ReplaceAll(edge.path, `"`, `\"`)
			fmt.Fprintf(file, "\tN_%s -> N_%s [label=\"%s\", style=dashed, color=blue];\n", safeID(edge.from.id), safeID(edge.to.id), path)
		} else if edge.from != nil && edge.to != nil {
			fmt.Fprintf(file, "\tN_%s -> N_%s;\n", safeID(edge.from.id), safeID(edge.to.id))
		}
	}

	fmt.Fprintln(file, "}")
}
