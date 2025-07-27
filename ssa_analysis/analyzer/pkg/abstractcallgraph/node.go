package abstractcallgraph

type NodeType int

const (
	NODE_SERVICE NodeType = iota
	NODE_DATABASE
	NODE_CLIENT
)

type AbstractNode struct {
	t    NodeType
	name string
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

func NewAbstractNode(name string, t NodeType) *AbstractNode {
	return &AbstractNode{
		name: name,
		t:    t,
	}
}
