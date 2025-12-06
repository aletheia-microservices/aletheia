package abstracttypes

import (
	"encoding/json"
)

type AbstractNode interface {
	GetName() string
	GetEdges() []AbstractEdge
	IsVisited() bool
	GetDepth() int
	AddEdge(edge AbstractEdge)
}

type AbstractServiceNode struct {
	AbstractNode `json:"-"`
	Name         string
	Edges        []AbstractEdge
	Visited      bool
	Depth        int
}

func (node *AbstractServiceNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name  string         `json:"name"`
		Edges []AbstractEdge `json:"edges"`
	}{
		Name:  node.Name,
		Edges: node.Edges,
	})
}

func (node *AbstractServiceNode) GetName() string {
	return node.Name
}

func (node *AbstractServiceNode) GetEdges() []AbstractEdge {
	return node.Edges
}

func (node *AbstractServiceNode) IsVisited() bool {
	return node.Visited
}

func (node *AbstractServiceNode) GetDepth() int {
	return node.Depth
}

func (node *AbstractServiceNode) AddEdge(edge AbstractEdge) {
	node.Edges = append(node.Edges, edge)
}

type AbstractDatastoreNode struct {
	AbstractNode `json:"-"`
	Name         string
	Edges        []AbstractEdge
	Visited      bool
	Depth        int
}

func (node *AbstractDatastoreNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name  string         `json:"name"`
		Edges []AbstractEdge `json:"edges"`
	}{
		Name:  node.Name,
		Edges: node.Edges,
	})
}

func (node *AbstractDatastoreNode) GetName() string {
	return node.Name
}

func (node *AbstractDatastoreNode) GetEdges() []AbstractEdge {
	return node.Edges
}

func (node *AbstractDatastoreNode) IsVisited() bool {
	return node.Visited
}

func (node *AbstractDatastoreNode) GetDepth() int {
	return node.Depth
}

func (node *AbstractDatastoreNode) AddEdge(edge AbstractEdge) {
	node.Edges = append(node.Edges, edge)
}
