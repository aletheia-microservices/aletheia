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

		// filter out any constraint that is not foreign key
		var fkConstraints []*datastores.Constraint
		for _, constraint := range constraints {
			if constraint.IsForeignKey() {
				fkConstraints = append(fkConstraints, constraint)
			}
		}

		// iterate all previous delete operations with pending deletes to mark as found if current one is a cascading delete
		if fkConstraints != nil {
			detector.markAsCascading(datastore, constraints)
		}
	}
}

// FIXME: to be more precise, we can check which object was deleted and if that specific reference to that object was deleted aswell
func (detector *CascadeDetector) OnDelete(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	datastore := dbCall.DbInstance.GetDatastore()
	deleteOp := newDeleteOperation(dbCall, datastore)

	// find pending deletes for datastores that hold a reference to the current one
	detector.findPendingCascadingDeletes(app, deleteOp, datastore)

	// iterate all previous delete operations with pending deletes to mark as found if current one is a cascading delete
	fkConstraints := datastore.GetSchema().GetConstraintsForeignKey()
	if fkConstraints != nil {
		detector.markAsCascading(datastore, fkConstraints)
	}

	detector.getCurrentRequestInfo().addDeleteOperation(deleteOp)
}

// findPendingCascadingDeletes finds pending deletes for datastores that hold a reference to the current one
func (detector *CascadeDetector) findPendingCascadingDeletes(app *app.App, deleteOp *deleteOperation, datastore *datastores.Datastore) {
	// iterate all datastores (except queues) that have fields referencing the current one
	for _, depDatastore := range app.GetDatabasesReferencingCurrentExceptQueues(datastore) {
		// we can have more than one service if multiple services share the same datastore
		depServices := app.GetServicesUsingDatastore(depDatastore)
		// mark a "pending delete" on all pairs of services-datastore that have any object referencing the one to be deleted in the current datastore
		for _, depForeignKeyConstraint := range depDatastore.GetSchema().GetConstraintsForeignKeyToDatastore(datastore) {
			pendingDel := newPendingDelete(depDatastore, depForeignKeyConstraint, depServices)
			// conditional "if not exists" just to avoid duplicates
			deleteOp.addPendingDeleteIfNotExists(pendingDel)
		}
	}
}

// markAsCascading iterates all previous delete operations with pending deletes to mark as found if current one is a cascading delete
func (detector *CascadeDetector) markAsCascading(datastore *datastores.Datastore, constraints []*datastores.Constraint) {
	// iterates all previous delete operations with pending deletes
	for _, prevOp := range detector.getCurrentRequestInfo().getDeleteOperations() {
		for _, prevPendingDel := range prevOp.getPendingDeletes() {

			// check if current datastore was associated with a previous pending delete
			// this is just to filter out in case it's false
			if prevPendingDel.isOnDatastore(datastore) {

				for _, fkConstraint := range constraints {
					// check if the foreign key constraint of the schema of the datastore involved the current operation
					// was previously associated with a pending delete (i.e., awaiting for cascade delete)
					if prevPendingDel.isOnConstraint(fkConstraint) {
						prevPendingDel.setCascading(true)
					}
				}
			}
		}
	}
}

func (detector *CascadeDetector) GetAnalysisTypeString() string {
	return "fkey_cascade"
}
