package ssa_graph

type EdgeType int

const (
	EDGE_USAGE EdgeType = iota
	EDGE_STORE
	EDGE_LOAD
	EDGE_FIELD
	EDGE_INDEX
	EDGE_PARAMETER
	EDGE_POINTS_TO
)

type SSAEdge struct {
	edgeType EdgeType
	from     *SSANode
	to       *SSANode

	index int
	param string

	path string //pointer only
}

func (edge *SSAEdge) GetType() EdgeType {
	return edge.edgeType
}

func (edge *SSAEdge) HasFromNode(node *SSANode) bool {
	return edge.from == node
}

func (edge *SSAEdge) GetFromNode() *SSANode {
	return edge.from
}

func (edge *SSAEdge) GetToNode() *SSANode {
	return edge.to
}

func (edge *SSAEdge) GetPath() string {
	return edge.path
}

func (edge *SSAEdge) SetPath(path string) {
	edge.path = path
}
