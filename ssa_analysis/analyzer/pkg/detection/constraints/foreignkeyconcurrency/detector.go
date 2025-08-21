package foreignkeyconcurrency

import (
	"fmt"
	"slices"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/app/backends"
	"analyzer/pkg/detection"
)

type ForeignKeyConcurrencyDetector struct {
	detection.Detector
	requests         []*Request
	dangerousDeletes map[*Request][]*DangerousDelete
	results          string
}

func NewDetector() *ForeignKeyConcurrencyDetector {
	fmt.Println()
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println(" ---------------------------------- INITIALIZING FOREIGN KEY CONCURRENCY DETECTOR --------------------------------- ")
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println()
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

func (detector *ForeignKeyConcurrencyDetector) OnNewRequest(node *abstractgraph.AbstractNode) {
	request := NewRequest(len(detector.requests), node)
	detector.requests = append(detector.requests, request)
	fmt.Printf("[FOREIGN KEY CONCURRENCY | DETECTOR] on new request\n")
}

func (detector *ForeignKeyConcurrencyDetector) OnEndRequest(app *app.App) {
	detector.FinalizeOperationsFields(app)
}

func (detector *ForeignKeyConcurrencyDetector) OnNewNode(app *app.App, node *abstractgraph.AbstractNode) {
	// nothing to do
}

func (detector *ForeignKeyConcurrencyDetector) OnEndNode(app *app.App, node *abstractgraph.AbstractNode) {
	// nothing to do
}

func (detector *ForeignKeyConcurrencyDetector) OnRead(app *app.App, edge *abstractgraph.AbstractEdge) {
	// nothing to do
}

func (detector *ForeignKeyConcurrencyDetector) OnWrite(app *app.App, edge *abstractgraph.AbstractEdge) {
	database := app.GetDatabaseByName(edge.GetToNode().GetDatabaseName())
	write := NewWriteOperation(edge, database)
	request := detector.getCurrentRequest()
	request.addWriteOperation(write)
	fmt.Printf("[FOREIGN KEY CONCURRENCY | DETECTOR] added new write: %v\n", write)
}

func (detector *ForeignKeyConcurrencyDetector) OnUpdate(app *app.App, edge *abstractgraph.AbstractEdge) {
	// nothing to do
}

func (detector *ForeignKeyConcurrencyDetector) OnDelete(app *app.App, edge *abstractgraph.AbstractEdge) {
	database := app.GetDatabaseByName(edge.GetToNode().GetDatabaseName())
	delete := NewDeleteOperation(edge, database)
	request := detector.getCurrentRequest()
	request.addDeleteOperation(delete)
	fmt.Printf("[FOREIGN KEY CONCURRENCY | DETECTOR] added new delete: %v\n", delete)
}

// logic very similar for write operations and delete operations
func (detector *ForeignKeyConcurrencyDetector) FinalizeOperationsFields(app *app.App) {
	request := detector.getCurrentRequest()

	// search for fields:
	// fields in current database with constraint foreign key + mandatory
	for _, write := range request.getAllWriteOperations() {
		fmt.Printf("[FOREIGN KEY CONCURRENCY | DETECTOR] write = %s\n", write.call.String())
		var fields []*backends.Field
		for _, arg := range write.call.GetArguments() {
			for _, fieldpath := range arg.GetAffectedDatabaseFieldsForCall(write.call.GetID()) {
				writtenField := app.ComputeDatabaseFieldFromPath(write.database, fieldpath)
				for _, field := range app.GetAllDatabaseFieldsWithPrefixPath(writtenField, true) {
					fmt.Printf("\t[FOREIGN KEY CONCURRENCY | DETECTOR] field = %s\n", field.String())
					if field.HasConstraintForeignKeyNonMandatory() && !slices.Contains(fields, field) {
						fmt.Printf("\t\t[FOREIGN KEY CONCURRENCY | DETECTOR] OK!\n")
						fields = append(fields, field)
					}
				}
			}
		}
		write.SetFields(fields)
	}

	// search for pending fields:
	// fields in other databases with constraint foreign key + mandatory
	// that reference some field in the current database
	for _, delete := range request.getAllDeleteOperations() {
		fmt.Printf("[FOREIGN KEY CONCURRENCY | DETECTOR] delete = %s\n", delete.call.String())
		for _, arg := range delete.call.GetArguments() {
			for _, fieldpath := range arg.GetAffectedDatabaseFieldsForCall(delete.call.GetID()) {
				// [TO BE IMPROVED]
				// just associated the schema to the call when parsing it...
				field := app.ComputeDatabaseFieldFromPath(delete.database, fieldpath)
				schema := field.GetSchema()
				delete.setSchema(schema)
				break
			}
		}
	}
}
