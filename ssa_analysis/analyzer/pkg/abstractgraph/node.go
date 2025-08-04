package abstractgraph

import "log"

type NodeType int

const (
	NODE_SERVICE NodeType = iota
	NODE_DATABASE
	NODE_CLIENT
)

type AbstractNode struct {
	t      NodeType
	name   string
	parsed bool

	// for service nodes only
	service string
	method  string
	params  []*AbstractObject
	rets    []*AbstractObject

	// for database nodes only
	dbname string
}

func (node *AbstractNode) String() string {
	return node.name
}

func (node *AbstractNode) IsParsed() bool {
	return node.parsed
}

func (node *AbstractNode) SetParsed() {
	node.parsed = true
}

func (node *AbstractNode) GetName() string {
	return node.name
}

func (node *AbstractNode) GetMethod() string {
	return node.method
}

func (node *AbstractNode) GetDatabaseName() string {
	return node.dbname
}

func (node *AbstractNode) GetNodeType() NodeType {
	return node.t
}

func (node *AbstractNode) GetParams() []*AbstractObject {
	return node.params
}

func (node *AbstractNode) GetParamAt(i int) *AbstractObject {
	if i > len(node.params)-1 {
		log.Fatalf("index (%d) out of bounds for node params: %v\n", i, node.params)
	}
	return node.params[i]
}

func (node *AbstractNode) AddParam(param *AbstractObject) {
	node.params = append(node.params, param)
}

func (node *AbstractNode) AddReturn(obj *AbstractObject) {
	node.rets = append(node.rets, obj)
}

func (node *AbstractNode) GetReturns() []*AbstractObject {
	return node.rets
}

func (node *AbstractNode) GetReturnAt(i int) *AbstractObject {
	if i > len(node.rets)-1 {
		log.Fatalf("index (%d) out of bounds for node returns: %v\n", i, node.rets)
	}
	return node.rets[i]
}

func NewAbstractNode(name string, t NodeType, service string, method string, dbname string) *AbstractNode {
	return &AbstractNode{
		name:    name,
		t:       t,
		service: service,
		method:  method,
		dbname:  dbname,
	}
}
