package abstractcallgraph

import "log"

type NodeType int

const (
	NODE_SERVICE NodeType = iota
	NODE_DATABASE
	NODE_CLIENT
)

type AbstractNode struct {
	t       NodeType
	name    string
	service string
	method  string
	params  []*AbstractObject
}

func (node *AbstractNode) String() string {
	return node.name
}

func (node *AbstractNode) GetName() string {
	return node.name
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

func NewAbstractNode(name string, t NodeType, service string, method string) *AbstractNode {
	return &AbstractNode{
		name:    name,
		t:       t,
		service: service,
		method:  method,
	}
}
