package abstractgraph

import (
	"container/list"
	"fmt"
	"sort"
	"strings"

	//"fmt"
	/* "strings"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/objects" */
	"analyzer/pkg/utils"
)

func (graph *AbstractGraph) Dump() {
	graph.dumpAbstractGraph()
	graph.dumpDiGraph()
}

func (graph *AbstractGraph) dumpAbstractGraph() {
	utils.DumpToJSONFile(graph.Nodes, graph.AppName, "abstractgraph")
}

func (graph *AbstractGraph) dumpDiGraph() {
	type node struct {
		Id   string `json:"id"`
		Type string `json:"type"`
	}
	type edge struct {
		Caller string `json:"caller"`
		Callee string `json:"callee"`
		Call   string `json:"call"`
		Depth  int    `json:"depth"`
	}
	type digraph struct {
		Nodes []node `json:"nodes"`
		Edges []edge `json:"edges"`
	}

	var nodeTypeOrder = map[string]int{
		"client":    1,
		"service":   2,
		"datastore": 3,
	}
	var sortNodes = func(nodes []node) {
		sort.Slice(nodes, func(i, j int) bool {
			if nodes[i].Type != nodes[j].Type {
				return nodeTypeOrder[nodes[i].Type] < nodeTypeOrder[nodes[j].Type]
			}
			return nodes[i].Id < nodes[j].Id
		})
	}
	var sortEdges = func(edges []edge) {
		sort.Slice(edges, func(i, j int) bool {
			if edges[i].Caller != edges[j].Caller {
				return edges[i].Caller < edges[j].Caller
			}
			if edges[i].Callee != edges[j].Callee {
				return edges[i].Callee < edges[j].Callee
			}
			if edges[i].Depth != edges[j].Depth {
				return edges[i].Depth < edges[j].Depth
			}
			return edges[i].Call < edges[j].Call
		})
	}

	global_nodes := []node{}
	global_edges := []edge{}

	for _, n := range graph.Nodes {
		var abstractNodesToVisit = list.New()
		var visitedAbstractNodes = make(map[AbstractNode]bool)

		abstractNodesToVisit.PushBack(n)
		nodes := []node{}
		edges := []edge{}
		appendedServiceNodes := make(map[string]bool)

		nodes = append(nodes, node{Id: "Client", Type: "client"})

		//fmt.Println()
		for abstractNodesToVisit.Len() > 0 {
			elem := abstractNodesToVisit.Front()
			n := elem.Value.(AbstractNode)
			abstractNodesToVisit.Remove(elem)

			if _, exists := visitedAbstractNodes[n]; !exists {
				caller := n.GetCallerStr()
				callee := n.GetCallee()
				call := n.GetName()

				if _, exists := appendedServiceNodes[callee]; !exists {
					nodes = append(nodes, node{Id: callee, Type: n.GetNodeType()})
					appendedServiceNodes[callee] = true
				}
				edges = append(edges, edge{Caller: caller, Callee: callee, Call: call, Depth: n.GetDepth()})
				visitedAbstractNodes[n] = true

				for _, c := range n.GetChildren() {
					if _, visited := visitedAbstractNodes[c]; !visited {
						abstractNodesToVisit.PushBack(c)
					}
				}
			}
		}
		global_nodes = append(global_nodes, nodes...)
		global_edges = append(global_edges, edges...)

		sortEdges(edges)
		sortNodes(nodes)
		outputGraph := digraph{Nodes: nodes, Edges: edges}
		utils.DumpToJSONFile(outputGraph, graph.AppName, fmt.Sprintf("graphs/output/trace_%s_%s", strings.ToLower(n.GetCallee()), strings.ToLower(n.GetName())))
	}

	sortEdges(global_edges)
	sortNodes(global_nodes)
	outputGraph := digraph{Nodes: global_nodes, Edges: global_edges}
	utils.DumpToJSONFile(outputGraph, graph.AppName, "graphs/output/call_graph")

}
