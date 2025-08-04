package abstractgraph

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"analyzer/pkg/common"
)

type EdgeType int

const (
	EDGE_SERVICE_RPC EdgeType = iota
	EDGE_DATABASE_CALL
	EDGE_SERVICE_ENTRYPOINT
)

type AbstractEdge struct {
	t EdgeType

	opType common.DatabaseOperationType // for database calls

	// format: <func_short_path>_<ssa_instr_name>
	// except for entrypoint edges where format is just <func_short_path>
	id     string
	method string
	from   *AbstractNode
	to     *AbstractNode
	args   []*AbstractObject // caller side
	rets   []*AbstractObject // caller side
}

func (edge *AbstractEdge) String() string {
	return fmt.Sprintf("(%s) --> (%s).%s(...)", edge.from.String(), edge.to.String(), edge.method)
}

func (edge *AbstractEdge) IsRead() bool {
	return edge.opType == common.OP_READ
}

func (edge *AbstractEdge) IsWrite() bool {
	return edge.opType == common.OP_WRITE
}

func (edge *AbstractEdge) IsUpdate() bool {
	return edge.opType == common.OP_UPDATE
}

func (edge *AbstractEdge) IsDelete() bool {
	return edge.opType == common.OP_DELETE
}

func (edge *AbstractEdge) GetID() string {
	return edge.id
}

// some exceptions can be:
// ProductService.New.nil:*github.com/blueprint-uservices/blueprint/examples/digota/workflow/digota.PackageDimensions
// ProductService.New."image":string
func (edge *AbstractEdge) GetIDNumber() int {
	parts := strings.Split(edge.id, ".t")
	if len(parts) != 2 {
		return -1
	}

	n, err := strconv.Atoi(parts[1])
	if err != nil {
		return -1
	}
	return n
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

func (edge *AbstractEdge) GetReturns() []*AbstractObject {
	return edge.rets
}

func (edge *AbstractEdge) GetReturnAt(i int) *AbstractObject {
	if i > len(edge.rets)-1 {
		log.Fatalf("index (%d) out of bounds for edge returns: %v\n", i, edge.rets)
	}
	return edge.rets[i]
}

func (edge *AbstractEdge) AddReturn(ret *AbstractObject) {
	edge.rets = append(edge.rets, ret)
}

func NewAbstractEdge(id string, method string, from *AbstractNode, to *AbstractNode, opType common.DatabaseOperationType, t EdgeType) *AbstractEdge {
	return &AbstractEdge{
		id:     id,
		t:      t,
		method: method,
		from:   from,
		to:     to,
		opType: opType,
	}
}
