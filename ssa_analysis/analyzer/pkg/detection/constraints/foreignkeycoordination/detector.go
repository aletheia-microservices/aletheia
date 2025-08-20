package foreignkeycoordination

import (
	"fmt"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/detection"
)

type ForeignKeyCoordinationDetector struct {
	detection.Detector
	keyType      string // 'primary-key' or 'foreign-key'
	requests     []*Request
	results      string
	foreignReads map[*Request][]*ForeignRead
}

func NewDetector(keyType string) *ForeignKeyCoordinationDetector {
	fmt.Println()
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println(" --------------------------------- INITIALIZING FOREIGN KEY COORDINATION DETECTOR --------------------------------- ")
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println()
	return &ForeignKeyCoordinationDetector{
		keyType:      keyType,
		foreignReads: make(map[*Request][]*ForeignRead),
	}
}

func (detector *ForeignKeyCoordinationDetector) addForeignRead(req *Request, foreignread *ForeignRead) {
	detector.foreignReads[req] = append(detector.foreignReads[req], foreignread)
}

func (detector *ForeignKeyCoordinationDetector) hasForeignRead(req *Request, op1 *ReadOperation, op2 *ReadOperation) bool {
	for _, foreignRead := range detector.foreignReads[req] {
		if foreignRead.op1 == op1 && foreignRead.op2 == op2 {
			return true
		}
	}
	return false
}

func (detector *ForeignKeyCoordinationDetector) getCurrentRequest() *Request {
	return detector.requests[len(detector.requests)-1]
}

func (detector *ForeignKeyCoordinationDetector) GetTypeString() string {
	return detector.keyType + "-coordination"
}

func (detector *ForeignKeyCoordinationDetector) OnNewRun(app *app.App) {
	// nothing to do
}

func (detector *ForeignKeyCoordinationDetector) OnEndRun(app *app.App) {
	// nothing to do
}

func (detector *ForeignKeyCoordinationDetector) OnNewRequest(node *abstractgraph.AbstractNode) {
	request := NewRequest(len(detector.requests), node)
	detector.requests = append(detector.requests, request)
	fmt.Printf("[FOREIGN KEY COORDINATION | DETECTOR] on new request\n")
}

func (detector *ForeignKeyCoordinationDetector) OnEndRequest(app *app.App) {
	// nothing to do
}

func (detector *ForeignKeyCoordinationDetector) OnNewNode(app *app.App, node *abstractgraph.AbstractNode) {
	// nothing to do
}

func (detector *ForeignKeyCoordinationDetector) OnEndNode(app *app.App, node *abstractgraph.AbstractNode) {
	// nothing to do
}

func (detector *ForeignKeyCoordinationDetector) OnRead(app *app.App, edge *abstractgraph.AbstractEdge) {
	read := NewReadOperation(edge, edge.GetArguments())
	request := detector.getCurrentRequest()

	// must check inconsistency before adding read to request
	detector.checkInconsistency(request, read)
	request.AddOperation(read)
	fmt.Printf("[FOREIGN KEY COORDINATION | DETECTOR] added new read: %v\n", read)
}

func (detector *ForeignKeyCoordinationDetector) OnWrite(app *app.App, edge *abstractgraph.AbstractEdge) {
	// nothing to do
}

func (detector *ForeignKeyCoordinationDetector) OnUpdate(app *app.App, edge *abstractgraph.AbstractEdge) {
	// nothing to do
}

func (detector *ForeignKeyCoordinationDetector) OnDelete(app *app.App, edge *abstractgraph.AbstractEdge) {
	// nothing to do
}
