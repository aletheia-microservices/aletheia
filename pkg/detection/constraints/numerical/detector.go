package numerical

import (
	"fmt"

	"github.com/golang-collections/collections/stack"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/datastores"
	"analyzer/pkg/detection/detection"
	"analyzer/pkg/frameworks/blueprint"
)

func NewDetector() *NumericalDetector {
	fmt.Println()
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println(" ---------------------------------------- INITIALIZING NUMERICAL DETECTOR ----------------------------------------- ")
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println()

	return &NumericalDetector{
		requestInfoStack: stack.New(),
	}
}

type NumericalDetector struct {
	detection.Detector
	results          string
	summary          string
	requestInfoStack *stack.Stack
}

func (detector *NumericalDetector) GetSummary() string {
	return detector.summary
}

func (detector *NumericalDetector) SetSummary(summary string) {
	detector.summary = summary
}

func (detector *NumericalDetector) getCurrentRequestInfo() *RequestInfo {
	return detector.requestInfoStack.Peek().(*RequestInfo)
}

func (detector *NumericalDetector) OnNewRun(app *app.App) {
	//no-op
}

func (detector *NumericalDetector) OnEndRun(app *app.App) {
	//no-op
}

func (detector *NumericalDetector) OnNewRequest(entryNode *abstractgraph.AbstractServiceCall) {
	detector.requestInfoStack.Push(&RequestInfo{
		entry: entryNode,
	})
}

func (detector *NumericalDetector) OnEndRequest(app *app.App) {
	//no-op
}

func (detector *NumericalDetector) OnNewNode(app *app.App, node abstractgraph.AbstractNode) {
	//no-op
}

func (detector *NumericalDetector) OnEndNode(app *app.App, node abstractgraph.AbstractNode) {
	//no-op
}

func (detector *NumericalDetector) OnRead(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	//no-op
}

func (detector *NumericalDetector) OnWrite(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	datastore := dbCall.DbInstance.GetDatastore()
	requestInfo := detector.getCurrentRequestInfo()

	operation := NewOperation(dbCall, datastore)
	requestInfo.addOperation(operation)
}

func (detector *NumericalDetector) OnUpdate(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	schema := dbCall.DbInstance.GetDatastore().GetSchema()
	datastore := dbCall.DbInstance.GetDatastore()
	if blueprintBackendMethod := dbCall.ParsedCall.Method.(*blueprint.BackendMethod); blueprintBackendMethod != nil &&
		schema.HasConstraintsNumerical() && datastore.Type == datastores.NoSQL {

		params := dbCall.GetParams()
		update := params[2]
		affectedConstraints, operationRepr := parseNoSQLUpdate(blueprintBackendMethod, schema, update)

		if len(affectedConstraints) > 0 {
			requestInfo := detector.getCurrentRequestInfo()
			operation := NewOperationOnNumericalConstraint(dbCall, datastore, affectedConstraints, operationRepr)
			requestInfo.addOperation(operation)
			requestInfo.writeOnConstraint = true
			return
		}
	}

	requestInfo := detector.getCurrentRequestInfo()
	operation := NewOperation(dbCall, datastore)
	requestInfo.addOperation(operation)
}

func (detector *NumericalDetector) OnDelete(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	datastore := dbCall.DbInstance.GetDatastore()
	requestInfo := detector.getCurrentRequestInfo()

	operation := NewOperation(dbCall, datastore)
	requestInfo.addOperation(operation)
}

func (detector *NumericalDetector) GetAnalysisTypeString() string {
	return "numerical"
}
