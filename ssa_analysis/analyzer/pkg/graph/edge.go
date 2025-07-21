package graph

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

type Edge struct {
	edgeType EdgeType
	from     *Node
	to       *Node

	index int
	param string

	path string //pointer only
}

func (edge *Edge) GetType() EdgeType {
	return edge.edgeType
}

func (edge *Edge) HasFromNode(node *Node) bool {
	return edge.from == node
}

func (edge *Edge) GetFromNode() *Node {
	return edge.from
}

func (edge *Edge) GetToNode() *Node {
	return edge.to
}

func (edge *Edge) GetPath() string {
	return edge.path
}

func (edge *Edge) SetPath(path string) {
	edge.path = path
}
