package uniquenessconcurrency

import (
	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/detection"
)

type UniquenessConcurrencyDetector struct {
	detection.Detector
	requests            []*Request
	results             string
	vulnerableWriteSets map[*Request][]*VulnerableWriteSet
}

func NewDetector() *UniquenessConcurrencyDetector {
	return &UniquenessConcurrencyDetector{
		vulnerableWriteSets: make(map[*Request][]*VulnerableWriteSet),
	}
}

func (detector *UniquenessConcurrencyDetector) addVulnerableWriteSet(req *Request, writeSet *VulnerableWriteSet) {
	detector.vulnerableWriteSets[req] = append(detector.vulnerableWriteSets[req], writeSet)
}

func (detector *UniquenessConcurrencyDetector) findVulnerableWriteSetForOperation(req *Request, op *WriteOperation) *VulnerableWriteSet {
	for _, writeSet := range detector.vulnerableWriteSets[req] {
		if writeSet.constrainedOp == op {
			return writeSet
		}
	}
	return nil
}

func (detector *UniquenessConcurrencyDetector) getCurrentRequest() *Request {
	return detector.requests[len(detector.requests)-1]
}

func (detector *UniquenessConcurrencyDetector) GetTypeString() string {
	return "uniqueness-concurrency"
}

func (detector *UniquenessConcurrencyDetector) OnNewRun(app *app.App) {
	// nothing to do
}

func (detector *UniquenessConcurrencyDetector) OnEndRun(app *app.App) {
	// nothing to do
}

func (detector *UniquenessConcurrencyDetector) OnNewRequest(node *abstractgraph.AbstractNode, reqIdx int) {
	request := NewRequest(len(detector.requests), node)
	detector.requests = append(detector.requests, request)
}

func (detector *UniquenessConcurrencyDetector) OnEndRequest(app *app.App) {
	// nothing to do
}

func (detector *UniquenessConcurrencyDetector) OnNewNode(app *app.App, node *abstractgraph.AbstractNode) {
	// nothing to do
}

func (detector *UniquenessConcurrencyDetector) OnEndNode(app *app.App, node *abstractgraph.AbstractNode) {
	// nothing to do
}

func (detector *UniquenessConcurrencyDetector) OnRead(app *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	// nothing to do
}

func (detector *UniquenessConcurrencyDetector) OnWrite(app *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	op := NewWriteOperation(edge, edge.GetArguments())
	request := detector.getCurrentRequest()

	// must check inconsistency before adding read to request
	detector.checkInconsistency(app, request, op)
	request.AddOperation(op)
}

func (detector *UniquenessConcurrencyDetector) OnUpdate(app *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	// nothing to do
}

func (detector *UniquenessConcurrencyDetector) OnDelete(app *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	// nothing to do
}
