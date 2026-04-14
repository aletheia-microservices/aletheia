package foreignkeycascade

import (
	"slices"

	"github.com/sirupsen/logrus"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/app/backends"
	"analyzer/pkg/detection"
)

type ForeignKeyCascadeDetector struct {
	detection.Detector
	requests       []*Request
	results        string
	cascadeDeletes map[*Request][]*CascadeDelete
}

func NewDetector() *ForeignKeyCascadeDetector {
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

func getFields(app *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) []*backends.Field {
	database := app.GetDatabaseByName(edge.GetToNode().GetDatabaseName())
	var fields []*backends.Field
	for _, arg := range edge.GetArguments() {
		for _, fieldpath := range arg.GetAffectedDatabaseFieldsForCall(edge.GetID()) {
			writtenField := app.ComputeDatabaseFieldFromPath(database, fieldpath)
			for _, field := range app.GetAllDatabaseFieldsWithPrefixPath(writtenField, true) {
				if field.HasConstraintForeignKeyNonMandatory() && !slices.Contains(fields, field) {
					fields = append(fields, field)
				}
			}
		}
	}
	return fields
}

func (detector *ForeignKeyCascadeDetector) OnWrite(app *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	op := NewWriteOperation(edge, edge.GetArguments(), edge.GetToNode().GetDatabaseName(), edge.GetToNode().GetSchemaName())
	request := detector.getCurrentRequest()
	request.AddWriteOperation(op)
	logrus.WithField("request", request.entry.String()).
		Debugf("[DETECTOR - FOREIGN KEY CASCADE] added write: %v\n", op.call.String())
	op.fields = getFields(app, reqIdx, edge)
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
