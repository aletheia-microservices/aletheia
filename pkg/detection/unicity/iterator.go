package unicity

import (
	"fmt"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/logger"
)

type Iterator struct {
	app      *app.App
	graph    *abstractgraph.AbstractGraph
	detector Detector
}

func (iterator *Iterator) getGraph() *abstractgraph.AbstractGraph {
	return iterator.graph
}

func NewIterator(app *app.App, graph *abstractgraph.AbstractGraph, detector Detector) *Iterator {
	return &Iterator{
		app:      app,
		graph:    graph,
		detector: detector,
	}
}

func (iterator *Iterator) Run() {
	for idx, entry := range iterator.getGraph().Nodes {
		entryServiceCall := entry.(*abstractgraph.AbstractServiceCall)
		iterator.detector.onNewRequest(entryServiceCall)
		iterator.transverseNode(idx, entryServiceCall, entry)
	}
}

type Detector interface {
	onNewRequest(entryNode *abstractgraph.AbstractServiceCall)
	onWrite(dbCall *abstractgraph.AbstractDatabaseCall)
	onUpdate(dbCall *abstractgraph.AbstractDatabaseCall)
	onDelete(dbCall *abstractgraph.AbstractDatabaseCall)
}

func (iterator *Iterator) transverseNode(child_idx int, lastServiceCallNode *abstractgraph.AbstractServiceCall, node abstractgraph.AbstractNode) {
	if svcCall, ok := node.(*abstractgraph.AbstractServiceCall); ok {
		lastServiceCallNode = svcCall
	}

	if dbCall, ok := node.(*abstractgraph.AbstractDatabaseCall); ok {
		fmt.Println()
		logger.Logger.Debugf("[ITERATOR #%d] %s", child_idx, dbCall.String())
		if dbCall.ParsedCall.Method.IsWrite() {
			iterator.detector.onWrite(dbCall)
		} else if dbCall.ParsedCall.Method.IsUpdate() {
			iterator.detector.onUpdate(dbCall)
		} else if dbCall.ParsedCall.Method.IsDelete() {
			iterator.detector.onDelete(dbCall)
		}
	}

	for idx, child := range node.GetChildren() {
		iterator.transverseNode(idx, lastServiceCallNode, child)
	}
}
