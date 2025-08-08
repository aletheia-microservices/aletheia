package unicityconcurrency

import (
	"fmt"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/detection"
)

type UnicityConcurrencyDetector struct {
	detection.Detector
	requests            []*Request
	summary             string
	vulnerableWriteSets map[*Request][]*VulnerableWriteSet
}

func NewDetector() *UnicityConcurrencyDetector {
	fmt.Println()
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println(" ------------------------------------ INITIALIZING UNICITY CONCURRENCY DETECTOR ----------------------------------- ")
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println()
	return &UnicityConcurrencyDetector{
		vulnerableWriteSets: make(map[*Request][]*VulnerableWriteSet),
	}
}

func (detector *UnicityConcurrencyDetector) addVulnerableWriteSet(req *Request, writeSet *VulnerableWriteSet) {
	detector.vulnerableWriteSets[req] = append(detector.vulnerableWriteSets[req], writeSet)
}

func (detector *UnicityConcurrencyDetector) findVulnerableWriteSetForOperation(req *Request, op *WriteOperation) *VulnerableWriteSet {
	for _, writeSet := range detector.vulnerableWriteSets[req] {
		if writeSet.constrainedOp == op {
			return writeSet
		}
	}
	return nil
}

func (detector *UnicityConcurrencyDetector) getCurrentRequest() *Request {
	return detector.requests[len(detector.requests)-1]
}

func (detector *UnicityConcurrencyDetector) GetTypeString() string {
	return "unicity-concurrency"
}

func (detector *UnicityConcurrencyDetector) OnNewRun(app *app.App) {
	// nothing to do
}

func (detector *UnicityConcurrencyDetector) OnEndRun(app *app.App) {
	// nothing to do
}

func (detector *UnicityConcurrencyDetector) OnNewRequest(node *abstractgraph.AbstractNode) {
	request := NewRequest(len(detector.requests), node)
	detector.requests = append(detector.requests, request)
	fmt.Printf("[DETECTOR - UNICITY CONCURRENCY] on new request\n")
}

func (detector *UnicityConcurrencyDetector) OnEndRequest(app *app.App) {
	// nothing to do
}

func (detector *UnicityConcurrencyDetector) OnNewNode(app *app.App, node *abstractgraph.AbstractNode) {
	// nothing to do
}

func (detector *UnicityConcurrencyDetector) OnEndNode(app *app.App, node *abstractgraph.AbstractNode) {
	// nothing to do
}

func (detector *UnicityConcurrencyDetector) OnRead(app *app.App, edge *abstractgraph.AbstractEdge) {
	// nothing to do
}

func (detector *UnicityConcurrencyDetector) OnWrite(app *app.App, edge *abstractgraph.AbstractEdge) {
	op := NewWriteOperation(edge, edge.GetArguments())
	request := detector.getCurrentRequest()

	// must check inconsistency before adding read to request
	detector.checkInconsistency(app, request, op)
	request.AddOperation(op)
	fmt.Printf("[DETECTOR - UNICITY CONCURRENCY] added new write: %v\n", op)
}

func (detector *UnicityConcurrencyDetector) OnUpdate(app *app.App, edge *abstractgraph.AbstractEdge) {
	// nothing to do
}

func (detector *UnicityConcurrencyDetector) OnDelete(app *app.App, edge *abstractgraph.AbstractEdge) {
	// nothing to do
}
