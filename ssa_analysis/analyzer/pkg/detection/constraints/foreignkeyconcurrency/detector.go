package foreignkeyconcurrency

import (
	"slices"

	"github.com/sirupsen/logrus"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/app/backends"
	"analyzer/pkg/config"
	"analyzer/pkg/detection"
)

type ForeignKeyConcurrencyDetector struct {
	detection.Detector
	requests         []*Request
	dangerousDeletes map[*Request][]*DangerousDelete
	results          string
}

func NewDetector() *ForeignKeyConcurrencyDetector {
	// EVAL: logrus.Traceln()
	// EVAL: logrus.Traceln(" ------------------------------------------------------------------------------------------------------------------ ")
	// EVAL: logrus.Traceln(" ---------------------------------- INITIALIZING FOREIGN KEY CONCURRENCY DETECTOR --------------------------------- ")
	// EVAL: logrus.Traceln(" ------------------------------------------------------------------------------------------------------------------ ")
	// EVAL: logrus.Traceln()
	return &ForeignKeyConcurrencyDetector{
		dangerousDeletes: make(map[*Request][]*DangerousDelete),
	}
}

func (detector *ForeignKeyConcurrencyDetector) addDangerousDelete(request *Request, dangerousDelete *DangerousDelete) {
	detector.dangerousDeletes[request] = append(detector.dangerousDeletes[request], dangerousDelete)
}

func (detector *ForeignKeyConcurrencyDetector) getCurrentRequest() *Request {
	return detector.requests[len(detector.requests)-1]
}

func (detector *ForeignKeyConcurrencyDetector) GetTypeString() string {
	return "foreign-key-concurrency"
}

func (detector *ForeignKeyConcurrencyDetector) OnNewRun(app *app.App) {
	// nothing to do
}

func (detector *ForeignKeyConcurrencyDetector) OnEndRun(app *app.App) {
	detector.checkInconsistencies()
}

func (detector *ForeignKeyConcurrencyDetector) OnNewRequest(node *abstractgraph.AbstractNode, reqIdx int) {
	request := NewRequest(len(detector.requests), node)
	detector.requests = append(detector.requests, request)
	// EVAL: logrus.Tracef("[FOREIGN KEY CONCURRENCY | DETECTOR] on new request\n")
}

func (detector *ForeignKeyConcurrencyDetector) OnEndRequest(app *app.App) {
	// nothing to do
}

func (detector *ForeignKeyConcurrencyDetector) OnNewNode(app *app.App, node *abstractgraph.AbstractNode) {
	// nothing to do
}

func (detector *ForeignKeyConcurrencyDetector) OnEndNode(app *app.App, node *abstractgraph.AbstractNode) {
	// nothing to do
}

func (detector *ForeignKeyConcurrencyDetector) OnRead(app *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	// nothing to do
}

func (detector *ForeignKeyConcurrencyDetector) OnWrite(app *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	database := app.GetDatabaseByName(edge.GetToNode().GetDatabaseName())
	request := detector.getCurrentRequest()
	write := NewWriteOperation(edge, database, request)

	// search for fields:
	// fields in current database with constraint foreign key + mandatory
	// EVAL: logrus.Tracef("[FOREIGN KEY CONCURRENCY | DETECTOR] write={%s}, entry={%s}\n", write.call.String(), write.request.entry.String())
	var fields []*backends.Field
	for _, arg := range write.call.GetArguments() {
		for _, fieldpath := range arg.GetAffectedDatabaseFieldsForCall(write.call.GetID()) {
			writtenField := app.ComputeDatabaseFieldFromPath(write.database, fieldpath)
			for _, field := range app.GetAllDatabaseFieldsWithPrefixPath(writtenField, true) {
				// EVAL: logrus.Tracef("\t[FOREIGN KEY CONCURRENCY | DETECTOR] field = %s\n", field.String())
				if field.HasConstraintForeignKeyNonMandatory() && !slices.Contains(fields, field) {
					// EVAL: logrus.Tracef("\t\t[FOREIGN KEY CONCURRENCY | DETECTOR] OK!\n")
					fields = append(fields, field)
				}
			}
		}
	}
	write.SetFields(fields)

	request.addWriteOperation(write)
	// EVAL: logrus.Tracef("[FOREIGN KEY CONCURRENCY | DETECTOR] added new write: %v\n", write)
}

func (detector *ForeignKeyConcurrencyDetector) OnUpdate(app *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	if config.Global.ForeignKeyConcurrencyDetectorIncludeOnUpdates {
		database := app.GetDatabaseByName(edge.GetToNode().GetDatabaseName())
		if database.IsNoSQL() && edge.GetMethod() == "Upsert" {
			detector.OnWrite(app, reqIdx, edge)
			return
		}
	}
	// nothing to do
}

func (detector *ForeignKeyConcurrencyDetector) OnDelete(app *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	database := app.GetDatabaseByName(edge.GetToNode().GetDatabaseName())
	delete := NewDeleteOperation(edge, database)
	request := detector.getCurrentRequest()
	var schema *backends.Schema

	// search for pending fields:
	// fields in other databases with constraint foreign key + mandatory
	// that reference some field in the current database
	// EVAL: logrus.Tracef("[FOREIGN KEY CONCURRENCY | DETECTOR] delete = %s\n", delete.call.String())
	for _, arg := range delete.call.GetArguments() {
		for _, fieldpath := range arg.GetAffectedDatabaseFieldsForCall(delete.call.GetID()) {
			// [TO BE IMPROVED]
			// just associated the schema to the call when parsing it...
			field := app.ComputeDatabaseFieldFromPath(delete.database, fieldpath)
			schema = field.GetSchema()
			break
		}
	}
	if schema == nil {
		logrus.WithField("delete_op", delete.call.String()).Fatalf("[FOREIGN KEY CONCURRENCY] [DETECTOR] unexpected nil schema")
	}
	delete.setSchema(schema)

	request.addDeleteOperation(delete)
	// EVAL: logrus.Tracef("[FOREIGN KEY CONCURRENCY | DETECTOR] added new delete: %v\n", delete)
}
