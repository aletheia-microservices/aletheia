package abstractcallgraph

import "fmt"

type EdgeType int

const (
	EDGE_SERVICE_RPC EdgeType = iota
	EDGE_DATABASE_CALL
)

type AbstractEdge struct {
	t            EdgeType
	id           string // the ssa instr name for the db call on the callee side
	method       string
	from         *AbstractNode
	to           *AbstractNode
	callArgs     []*AbstractArgument // caller side
	methodParams []*AbstractArgument // callee side
}

func (edge *AbstractEdge) String() string {
	return fmt.Sprintf("(%s) --> (%s).%s(...)", edge.from.String(), edge.to.String(), edge.method)
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

func (edge *AbstractEdge) AddCallArgument(arg *AbstractArgument) {
	edge.callArgs = append(edge.callArgs, arg)
}

func (edge *AbstractEdge) AddMethodParameter(arg *AbstractArgument) {
	edge.methodParams = append(edge.methodParams, arg)
}

func NewAbstractEdge(id string, method string, from *AbstractNode, to *AbstractNode, t EdgeType) *AbstractEdge {
	return &AbstractEdge{
		id:     id,
		t:      t,
		method: method,
		from:   from,
		to:     to,
	}
}
