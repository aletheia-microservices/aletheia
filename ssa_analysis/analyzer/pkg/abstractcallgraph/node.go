package abstractcallgraph

type NodeType int

const (
	NODE_SERVICE NodeType = iota
	NODE_DATABASE
)

type AbstractNode struct {
	t    NodeType
	name string
}

func (node *AbstractNode) GetNodeType() NodeType {
	return node.t
}

func NewAbstractNode(name string) *AbstractNode {
	return &AbstractNode{
		name: name,
	}
}
