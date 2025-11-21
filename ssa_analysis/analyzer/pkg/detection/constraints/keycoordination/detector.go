package keycoordination

import (
	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/detection"
)

type DetectionType int

const (
	DETECTION_TYPE_PRIMARY_KEY DetectionType = iota
	DETECTION_TYPE_FOREIGN_KEY
)

type KeyCoordinationDetector struct {
	detection.Detector
	keyType      DetectionType
	requests     []*Request
	results      string
	foreignReads map[*Request][]*ForeignRead
}

func NewDetector(keyType DetectionType) *KeyCoordinationDetector {
	detector := &KeyCoordinationDetector{
		keyType:      keyType,
		foreignReads: make(map[*Request][]*ForeignRead),
	}
	// EVAL: fmt.Println()
	// EVAL: fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	// EVAL: fmt.Printf(" --------------------------------- INITIALIZING %s DETECTOR --------------------------------- \n", detector.GetTypeStringUpper())
	// EVAL: fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	// EVAL: fmt.Println()
	return detector
}

func (detector *KeyCoordinationDetector) isTypePrimaryKey() bool {
	return detector.keyType == DETECTION_TYPE_PRIMARY_KEY
}

func (detector *KeyCoordinationDetector) isTypeForeignKey() bool {
	return detector.keyType == DETECTION_TYPE_FOREIGN_KEY
}

func (detector *KeyCoordinationDetector) addForeignRead(req *Request, foreignread *ForeignRead) {
	detector.foreignReads[req] = append(detector.foreignReads[req], foreignread)
}

func (detector *KeyCoordinationDetector) hasForeignRead(req *Request, op1 *ReadOperation, op2 *ReadOperation) bool {
	for _, foreignRead := range detector.foreignReads[req] {
		if foreignRead.op1 == op1 && foreignRead.op2 == op2 || foreignRead.op1 == op2 && foreignRead.op2 == op1 {
			return true
		}

		if foreignRead.op1.call == op1.call && foreignRead.op2.call == op2.call || foreignRead.op1.call == op2.call && foreignRead.op2.call == op1.call {
			return true
		}
	}
	return false
}

func (detector *KeyCoordinationDetector) getCurrentRequest() *Request {
	return detector.requests[len(detector.requests)-1]
}

func (detector *KeyCoordinationDetector) GetTypeStringUpper() string {
	if detector.keyType == DETECTION_TYPE_PRIMARY_KEY {
		return "PRIMARY KEY COORDINATION"
	}
	return "FOREIGN KEY COORDINATION"
}

func (detector *KeyCoordinationDetector) GetTypeString() string {
	if detector.keyType == DETECTION_TYPE_PRIMARY_KEY {
		return "primary-key-coordination"
	}
	return "foreign-key-coordination"
}

func (detector *KeyCoordinationDetector) OnNewRun(app *app.App) {
	// nothing to do
}

func (detector *KeyCoordinationDetector) OnEndRun(app *app.App) {
	// nothing to do
}

func (detector *KeyCoordinationDetector) OnNewRequest(node *abstractgraph.AbstractNode, reqIdx int) {
	request := NewRequest(len(detector.requests), node)
	detector.requests = append(detector.requests, request)
	// EVAL: fmt.Printf("[%s | DETECTOR] on new request\n", detector.GetTypeStringUpper())
}

func (detector *KeyCoordinationDetector) OnEndRequest(app *app.App) {
	detector.checkInconsistency(app, detector.getCurrentRequest())
}

func (detector *KeyCoordinationDetector) OnNewNode(app *app.App, node *abstractgraph.AbstractNode) {
	// nothing to do
}

func (detector *KeyCoordinationDetector) OnEndNode(app *app.App, node *abstractgraph.AbstractNode) {
	// nothing to do
}

func (detector *KeyCoordinationDetector) OnRead(app *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	read := NewReadOperation(edge, edge.GetArguments(), reqIdx)
	request := detector.getCurrentRequest()
	request.AddOperation(read)
	// EVAL: fmt.Printf("[%s | DETECTOR] added new read: %v\n", detector.GetTypeStringUpper(), read)
}

func (detector *KeyCoordinationDetector) OnWrite(app *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	// nothing to do
}

func (detector *KeyCoordinationDetector) OnUpdate(app *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	// nothing to do
}

func (detector *KeyCoordinationDetector) OnDelete(app *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	// nothing to do
}
