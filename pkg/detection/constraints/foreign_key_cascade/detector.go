package foreign_key_cascade

import (
	"fmt"

	"github.com/golang-collections/collections/stack"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/datastores"
	"analyzer/pkg/detection/detection"
	"analyzer/pkg/frameworks/blueprint"
	"analyzer/pkg/logger"
	"analyzer/pkg/utils"
)

type CascadeDetector struct {
	detection.Detector
	requestInfoStack           *stack.Stack
	results                    string
	summary                    string
	deleteOperations           []*deleteOperation
	numDeletes                 int
	numMissingCascadingDeletes int
}

func (detector *CascadeDetector) addDeleteOperation(op *deleteOperation) {
	detector.deleteOperations = append(detector.deleteOperations, op)
}

func (detector *CascadeDetector) getCurrentRequestInfo() *RequestInfo {
	return detector.requestInfoStack.Peek().(*RequestInfo)
}

func (detector *CascadeDetector) getDeleteOperations() []*deleteOperation {
	return detector.deleteOperations
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
	params := dbCall.GetParams()
	update := params[2]
	method := dbCall.ParsedCall.Method.(*blueprint.BackendMethod)

	if dbCall.DbInstance.GetDatastore().IsNoSQLDatabase() {
		if constraints, _ := parseNoSQLUpdateOnRemovedFields(method, datastore.GetSchema(), update); constraints != nil {
			detector.markAsCascading(dbCall.DbInstance.GetDatastore(), constraints)
		}
	}
}

func (detector *CascadeDetector) OnDelete(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	datastore := dbCall.DbInstance.GetDatastore()

	// 1. check if other datastores hold a foreign key referencing the deleted object
	// 2. for each one of these services, check if they were notified about the deletion of the object
	// either before the deletion (which does not make sense, but is still possible) or after deleting
	// the "deletion notification" is expected to contain some information about the object that was just deleted and is done throught one of the following:
	// (i) message broker (queue)
	// (ii) RPC
	// NOTE: for now, we ony consider message brokers

	// 3. TODO LATER: to be more precise, we can check which object was deleted and if that specific reference to that object was deleted aswell

	deleteOp := newDeleteOperation(dbCall, datastore)

	logger.Logger.Infof("[CASCADE DETECTOR] searching dependencies for datastore (%s)", dbCall.DbInstance.GetDatastore().GetName())
	// iterate all datastores (except caches and queues!!) that have fields referencing the current one
	for _, dependentDatastore := range app.GetDatabasesReferencingCurrent(datastore) {
		depServices := app.GetServicesUsingDatastore(dependentDatastore)
		for _, constraint := range dependentDatastore.GetSchema().GetConstraintsForeignKeyToDatastore(datastore) {
			pendingDel := newPendingDelete(dependentDatastore, constraint, depServices)
			deleteOp.addPendingDeleteIfNotExists(pendingDel)
		}
	}

	detector.getCurrentRequestInfo().addDeleteOperation(deleteOp)
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

// DEPRECATED
func (detector *CascadeDetector) searchCascadingDeletes(deleteOp *deleteOperation, lastServiceCallNode *abstractgraph.AbstractServiceCall, deleteCall *abstractgraph.AbstractDatabaseCall, child_idx int) {
	logger.Logger.Infof("[CASCADE DETECTOR] searching for cascading deletes originating @ (%s, %s): %s", deleteCall.GetCallerStr(), deleteCall.DbInstance.GetDatastore().GetName(), deleteCall.LongString())
	numServiceCalls := len(lastServiceCallNode.GetChildren())
	for idx := child_idx + 1; idx < numServiceCalls; idx++ {
		node := lastServiceCallNode.GetChildAt(idx)
		// found call that pushes notifications
		if call, ok := node.(*abstractgraph.AbstractDatabaseCall); ok && call.IsPushToQueue() {
			logger.Logger.Debugf("[CASCADE DETECTOR] found push call @ (%s, %s): %s", call.GetCallerStr(), call.DbInstance.GetDatastore().GetName(), call.String())
			// check all calls that follow the read of the queue
			for _, queueReadNode := range call.GetChildren() {
				// check if after reading the queue, there is a delete operation to the same original database that is being referenced
				for _, childDbCall := range queueReadNode.GetDatabaseCalls() {
					if childDbCall.IsUpdateOrDelete() {
						logger.Logger.Debugf("[CASCADE DETECTOR] found child update/delete call @ (%s, %s): %s", childDbCall.Service, childDbCall.DbInstance, childDbCall.LongString())
						if deleteDep := deleteOp.getDependency(childDbCall.DbInstance.GetDatastore()); deleteDep != nil {
							logger.Logger.Debugf("[CASCADE DETECTOR] found cascading action!")
							deleteDep.setCascading(true)
						}
					}
				}
			}
		}
	}
}

func (detector *CascadeDetector) checkInconsistencies() {
	detector.numDeletes = len(detector.getDeleteOperations())

	/* for i, op := range detector.getDeleteOperations() { */
	/* } */

	for detector.requestInfoStack.Len() > 0 {
		requestInfo := detector.requestInfoStack.Pop().(*RequestInfo)
		for i, op := range requestInfo.getDeleteOperations() {
			depsWithMissingCascading := op.getDependenciesWithMissingCascade()
			detector.results += fmt.Sprintf("[%d] delete with %d missing cascades:\n", i+1, len(depsWithMissingCascading))
			detector.results += fmt.Sprintf("%s: %s\n", op.getCall().GetCallerStr(), op.call.ShortString())
			for _, dep := range depsWithMissingCascading {
				if !dep.cascading {
					detector.results += fmt.Sprintf("\t- %s\n", dep.LongString())
					detector.numMissingCascadingDeletes++
				}
			}
		}
	}
}

func (detector *CascadeDetector) ComputeResults() {
	header := "---------------------------------------------------------------------\n"
	header += "-------------------- FOREIGN KEY CASCADE ANALYSIS -------------------\n"
	header += "---------------------------------------------------------------------\n"

	detector.checkInconsistencies()

	header += fmt.Sprintf(">> (# DELETES ON REFERENCED OBJECT; # ABSENCE OF CASCADING DELETES):\n>> (%d;%d)\n", detector.numDeletes, detector.numMissingCascadingDeletes)
	detector.results = header + "---------------------------------------------------------------------\n" + utils.TEXT_RESET_COLOR + detector.results
}

func (detector *CascadeDetector) GetAnalysisTypeString() string {
	return "foreign_key_cascade"
}

func (detector *CascadeDetector) GetResults() string {
	return detector.results
}
