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
	edgeType EdgeType

	opType common.DatabaseOperationType // for database calls
	t      string                       // format: <ssa_variable_name> (only when primary!!)

	// format: <func_short_path>_<ssa_instr_name>
	// except for entrypoint edges where format is just <func_short_path>
	id     string
	method string
	from   *AbstractNode
	to     *AbstractNode
	args   []*AbstractObject // caller side
	rets   []*AbstractObject // caller side
}

func NewAbstractEdge(t string, id string, method string, from *AbstractNode, to *AbstractNode, opType common.DatabaseOperationType, edgeType EdgeType) *AbstractEdge {
	return &AbstractEdge{
		t:        t,
		id:       id,
		edgeType: edgeType,
		method:   method,
		from:     from,
		to:       to,
		opType:   opType,
	}
}

func (edge *AbstractEdge) String() string {
	return fmt.Sprintf("%s() --> %s.%s()", edge.from.String(), edge.to.String(), edge.method)
}

func (edge *AbstractEdge) GetOpType() common.DatabaseOperationType {
	return edge.opType
}

func (edge *AbstractEdge) GetID() string {
	return edge.id
}

func (edge *AbstractEdge) GetT() string {
	return edge.t
}

func (edge *AbstractEdge) GetTNumber() int {
	n, _ := strconv.Atoi(edge.t[1:]) // assuming "t3", "t13", etc.
	return n
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
	return edge.edgeType
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

func (edge *AbstractEdge) GetArgumentByNameIfExists(name string) *AbstractObject {
	for _, arg := range edge.args {
		if arg.name == name {
			return arg
		}
	}
	return nil
}

func (edge *AbstractEdge) GetArgumentAt(i int) *AbstractObject {
	if i > len(edge.args)-1 {
		log.Fatalf("index (%d) out of bounds for edge arguments: %v\n", i, edge.args)
	}
	return edge.args[i]
}

func (edge *AbstractEdge) SetArguments(args []*AbstractObject) {
	edge.args = args
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
