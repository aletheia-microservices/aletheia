package iterator

import (
	"fmt"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/detection/detector"
	"analyzer/pkg/logger"
)

type Iterator struct {
	app      *app.App
	graph    *abstractgraph.AbstractGraph
	detector detector.Detector
}

func (iterator *Iterator) getGraph() *abstractgraph.AbstractGraph {
	return iterator.graph
}

func NewIterator(app *app.App, graph *abstractgraph.AbstractGraph, detector detector.Detector) *Iterator {
	return &Iterator{
		app:      app,
		graph:    graph,
		detector: detector,
	}
}

func (iterator *Iterator) Run() {
	iterator.detector.OnNewRun(iterator.app)
	for idx, entry := range iterator.getGraph().Nodes {
		entryServiceCall := entry.(*abstractgraph.AbstractServiceCall)
		iterator.detector.OnNewRequest(entryServiceCall)
		iterator.transverseNode(idx, entryServiceCall, entry)
		iterator.detector.OnEndRequest(iterator.app)
	}
	iterator.detector.OnEndRun(iterator.app)
}

func (iterator *Iterator) transverseNode(child_idx int, lastServiceCallNode *abstractgraph.AbstractServiceCall, node abstractgraph.AbstractNode) {
	if svcCall, ok := node.(*abstractgraph.AbstractServiceCall); ok {
		lastServiceCallNode = svcCall
	}

	iterator.detector.OnNewNode(iterator.app, node)

	if dbCall, ok := node.(*abstractgraph.AbstractDatabaseCall); ok {
		fmt.Println()
		logger.Logger.Debugf("[ITERATOR #%d] %s", child_idx, dbCall.String())
		if dbCall.ParsedCall.Method.IsRead() {
			iterator.detector.OnRead(iterator.app, dbCall, lastServiceCallNode, child_idx)
		} else if dbCall.ParsedCall.Method.IsWrite() {
			iterator.detector.OnWrite(iterator.app, dbCall, lastServiceCallNode, child_idx)
		} else if dbCall.ParsedCall.Method.IsUpdate() {
			iterator.detector.OnUpdate(iterator.app, dbCall, lastServiceCallNode, child_idx)
		} else if dbCall.ParsedCall.Method.IsDelete() {
			iterator.detector.OnDelete(iterator.app, dbCall, lastServiceCallNode, child_idx)
		}
	}

	for idx, child := range node.GetChildren() {
		iterator.transverseNode(idx, lastServiceCallNode, child)
	}

	iterator.detector.OnEndNode(iterator.app, node)
}
