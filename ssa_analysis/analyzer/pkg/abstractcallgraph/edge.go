package abstractcallgraph

type EdgeType int

const (
	EDGE_SERVICE_RPC EdgeType = iota
	EDGE_DATABASE_CALL
)

type AbstractEdge struct {
	t    EdgeType
	str  string
	from *AbstractNode
	to   *AbstractNode
}

func (edge *AbstractEdge) GetEdgeType() EdgeType {
	return edge.t
}

func (edge *AbstractEdge) GetFromNode() *AbstractNode {
	return edge.from
}

func (edge *AbstractEdge) GetToNode() *AbstractNode {
	return edge.to
}

func NewAbstractEdge(str string, from *AbstractNode, to *AbstractNode, t EdgeType) *AbstractEdge {
	return &AbstractEdge{
		t:    t,
		str:  str,
		from: from,
		to:   to,
	}
}
