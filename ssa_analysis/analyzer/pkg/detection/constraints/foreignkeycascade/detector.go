package foreignkeycascade

import (
	"github.com/sirupsen/logrus"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/detection"
)

type ForeignKeyCascadeDetector struct {
	detection.Detector
	requests       []*Request
	results        string
	cascadeDeletes map[*Request][]*CascadeDelete
}

func NewDetector() *ForeignKeyCascadeDetector {
	// EVAL: fmt.Println()
	// EVAL: fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	// EVAL: fmt.Println(" ------------------------------------ INITIALIZING FOREIGN KEY CASCADE DETECTOR ----------------------------------- ")
	// EVAL: fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	// EVAL: fmt.Println()
	return &ForeignKeyCascadeDetector{
		cascadeDeletes: make(map[*Request][]*CascadeDelete),
	}
}

func (detector *ForeignKeyCascadeDetector) addCascadeDelete(req *Request, cascadeDelete *CascadeDelete) {
	detector.cascadeDeletes[req] = append(detector.cascadeDeletes[req], cascadeDelete)
}

func (detector *ForeignKeyCascadeDetector) getCascadeDeletesForRequest(req *Request) []*CascadeDelete {
	return detector.cascadeDeletes[req]
}

func (detector *ForeignKeyCascadeDetector) getCurrentRequest() *Request {
	return detector.requests[len(detector.requests)-1]
}

func (detector *ForeignKeyCascadeDetector) GetTypeString() string {
	return "foreign-key-cascade"
}

func (detector *ForeignKeyCascadeDetector) OnNewRun(app *app.App) {
	// nothing to do
}

func (detector *ForeignKeyCascadeDetector) OnEndRun(app *app.App) {
	detector.checkInconsistencies(app)
}

func (detector *ForeignKeyCascadeDetector) OnNewRequest(node *abstractgraph.AbstractNode, reqIdx int) {
	request := NewRequest(len(detector.requests), node)
	detector.requests = append(detector.requests, request)
	// EVAL: fmt.Printf("[DETECTOR - FOREIGN KEY CASCADE] on new request\n")
}

func (detector *ForeignKeyCascadeDetector) OnEndRequest(app *app.App) {
	// nothing to do
}

func (detector *ForeignKeyCascadeDetector) OnNewNode(app *app.App, node *abstractgraph.AbstractNode) {
	// nothing to do
}

func (detector *ForeignKeyCascadeDetector) OnEndNode(app *app.App, node *abstractgraph.AbstractNode) {
	// nothing to do
}

func (detector *ForeignKeyCascadeDetector) OnRead(app *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	// nothing to do
}

func (detector *ForeignKeyCascadeDetector) OnWrite(app *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	// nothing to do
}

func (detector *ForeignKeyCascadeDetector) OnUpdate(app *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	// nothing to do
}

func (detector *ForeignKeyCascadeDetector) OnDelete(app *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	op := NewDeleteOperation(edge, edge.GetArguments(), edge.GetToNode().GetDatabaseName(), edge.GetToNode().GetSchemaName())
	request := detector.getCurrentRequest()
	request.AddOperation(op)
	logrus.WithField("request", request.entry.String()).
		Debugf("[DETECTOR - FOREIGN KEY CASCADE] added new delete: %v\n", op.call.String())
}
