package fkey_cascade

import (
	"fmt"

	"github.com/golang-collections/collections/stack"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/datastores"
	"analyzer/pkg/detection/detection"
	"analyzer/pkg/frameworks/blueprint"
)

type CascadeDetector struct {
	detection.Detector
	requestInfoStack           *stack.Stack
	results                    string
	summary                    string
	numDeletes                 int
	numMissingCascadingDeletes int
}

func (detector *CascadeDetector) getCurrentRequestInfo() *RequestInfo {
	return detector.requestInfoStack.Peek().(*RequestInfo)
}

func (detector *CascadeDetector) GetSummary() string {
	return detector.summary
}

func (detector *CascadeDetector) SetSummary(summary string) {
	detector.summary = summary
}

func NewDetector() *CascadeDetector {
	fmt.Println()
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println(" ----------------------------------------- INITIALIZING CASCADE DETECTOR ------------------------------------------ ")
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println()
	return &CascadeDetector{
		requestInfoStack: stack.New(),
	}
}

func (detector *CascadeDetector) OnNewRun(app *app.App) {
	app.ResetAllDataflows()
}

func (detector *CascadeDetector) OnEndRun(app *app.App) {
	//no-op
}

func (detector *CascadeDetector) OnNewRequest(entryNode *abstractgraph.AbstractServiceCall) {
	detector.requestInfoStack.Push(&RequestInfo{
		index: detector.requestInfoStack.Len(),
		entry: entryNode,
	})
}

func (detector *CascadeDetector) OnEndRequest(app *app.App) {
	app.ResetAllDataflows()
}

func (detector *CascadeDetector) OnNewNode(app *app.App, node abstractgraph.AbstractNode) {
	//no-op
}

func (detector *CascadeDetector) OnEndNode(app *app.App, node abstractgraph.AbstractNode) {
	//no-op
}

func (detector *CascadeDetector) OnRead(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	//no-op
}

func (detector *CascadeDetector) OnWrite(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	//no-op
}

func (detector *CascadeDetector) OnUpdate(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	datastore := dbCall.DbInstance.GetDatastore()

	if dbCall.DbInstance.GetDatastore().IsNoSQLDatabase() {
		update := dbCall.GetParam(2)
		method := dbCall.ParsedCall.Method.(*blueprint.BackendMethod)
		constraints, _ := parseNoSQLUpdateOnRemovedFields(method, datastore.GetSchema(), update)
		if constraints != nil {
			detector.markAsCascading(datastore, constraints)
		}
	}
}

func (detector *CascadeDetector) OnDelete(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	datastore := dbCall.DbInstance.GetDatastore()
	foreignKeyConstraints := datastore.GetSchema().GetConstraintsForeignKey()

	// 1. check if other datastores hold a foreign key referencing the deleted object
	// 2. for each one of these services, check if they were notified about the deletion of the object
	// either before the deletion (which does not make sense, but is still possible) or after deleting
	// the "deletion notification" is expected to contain some information about the object that was just deleted and is done throught one of the following:
	// (i) message broker (queue)
	// (ii) RPC
	// NOTE: for now, we ony consider message brokers

	// 3. TODO LATER: to be more precise, we can check which object was deleted and if that specific reference to that object was deleted aswell

	deleteOp := newDeleteOperation(dbCall, datastore)

	detector.findPendingCascadingDeletes(app, deleteOp, datastore)

	detector.markAsCascading(datastore, foreignKeyConstraints)

	detector.getCurrentRequestInfo().addDeleteOperation(deleteOp)
}

func (detector *CascadeDetector) findPendingCascadingDeletes(app *app.App, deleteOp *deleteOperation, datastore *datastores.Datastore) {
	// iterate all datastores (except queues!!) that have fields referencing the current one
	for _, dependentDatastore := range app.GetDatabasesReferencingCurrent(datastore) {
		depServices := app.GetServicesUsingDatastore(dependentDatastore)
		for _, constraint := range dependentDatastore.GetSchema().GetConstraintsForeignKeyToDatastore(datastore) {
			pendingDel := newPendingDelete(dependentDatastore, constraint, depServices)
			deleteOp.addPendingDeleteIfNotExists(pendingDel)
		}
	}
}

func (detector *CascadeDetector) markAsCascading(datastore *datastores.Datastore, constraints []*datastores.Constraint) {
	for _, op := range detector.getCurrentRequestInfo().getDeleteOperations() {
		for _, pendingDel := range op.getPendingDeletes() {
			for _, constraint := range constraints {
				if pendingDel.isOnDatastore(datastore) && pendingDel.isOnConstraint(constraint) {
					pendingDel.setCascading(true)
				}
			}
		}
	}
}

func (detector *CascadeDetector) GetAnalysisTypeString() string {
	return "fkey_cascade"
}
