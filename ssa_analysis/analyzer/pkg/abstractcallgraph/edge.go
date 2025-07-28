package abstractcallgraph

import (
	"fmt"
	"log"
)

type EdgeType int

const (
	EDGE_SERVICE_RPC EdgeType = iota
	EDGE_DATABASE_CALL
	EDGE_SERVICE_ENTRYPOINT
)

type AbstractEdge struct {
	t EdgeType

	// format: <func_short_path>_<ssa_instr_name>
	// except for entrypoint edges where format is just <func_short_path>
	id     string
	method string
	from   *AbstractNode
	to     *AbstractNode
	args   []*AbstractObject // caller side
}

func (edge *AbstractEdge) String() string {
	return fmt.Sprintf("(%s) --> (%s).%s(...)", edge.from.String(), edge.to.String(), edge.method)
}

func (edge *AbstractEdge) GetID() string {
	return edge.id
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

func (edge *AbstractEdge) GetMethod() string {
	return edge.method
}

func (edge *AbstractEdge) GetArguments() []*AbstractObject {
	return edge.args
}

func (edge *AbstractEdge) GetArgumentAt(i int) *AbstractObject {
	if i > len(edge.args)-1 {
		log.Fatalf("index (%d) out of bounds for edge arguments: %v\n", i, edge.args)
	}
	return edge.args[i]
}

func (edge *AbstractEdge) AddArgument(arg *AbstractObject) {
	edge.args = append(edge.args, arg)
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
