package foreign_key_cascade

import (
	"fmt"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/detection/detection"
	"analyzer/pkg/logger"
)

type CascadeDetector struct {
	detection.Detector
	results                    string
	summary                    string
	deleteOperations           []*deleteOperation
	numDeletes                 int
	numMissingCascadingDeletes int
}

func (detector *CascadeDetector) addDeleteOperation(op *deleteOperation) {
	detector.deleteOperations = append(detector.deleteOperations, op)
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
	return &CascadeDetector{}
}

func (detector *CascadeDetector) OnNewRun(app *app.App) {
	app.ResetAllDataflows()
}

func (detector *CascadeDetector) OnEndRun(app *app.App) {
	//no-op
}

func (detector *CascadeDetector) OnNewRequest(entryNode *abstractgraph.AbstractServiceCall) {
	//no-op
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
	//no-op
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

	deleteOp := &deleteOperation{
		call:      dbCall,
		datastore: datastore,
	}

	logger.Logger.Infof("[CASCADE DETECTOR] searching dependencies for datastore (%s)", dbCall.DbInstance.GetDatastore().GetName())
	for _, datastoreInstance := range app.GetDbInstances() {
		dependentDatastore := datastoreInstance.GetDatastore()
		// discard current datastore
		if dependentDatastore != deleteOp.datastore {
			// found a datastore that has fields referencing the current one
			if dependentDatastore.IsReferencingDatastore(deleteOp.datastore) {
				logger.Logger.Infof("[CASCADE DETECTOR] found dependency for datastore (%s): %s", deleteOp.datastore.GetName(), dependentDatastore.GetName())
				for _, service := range app.GetServices() {
					if service.HasDatastore(dependentDatastore) {
						deleteDep := &deleteDependency{
							service:   service,
							datastore: dependentDatastore,
						}
						if !deleteOp.hasDependency(deleteDep) {
							deleteOp.addDependency(deleteDep)
							logger.Logger.Debugf("[CASCADE DETECTOR] added dependency %s to %s", deleteDep.String(), deleteOp.String())
						}
					}
				}
			}
		}
	}
	detector.addDeleteOperation(deleteOp)
	detector.searchCascadingDeletes(deleteOp, lastServiceCallNode, dbCall, child_idx)
}

func (detector *CascadeDetector) searchCascadingDeletes(deleteOp *deleteOperation, lastServiceCallNode *abstractgraph.AbstractServiceCall, deleteCall *abstractgraph.AbstractDatabaseCall, child_idx int) {
	logger.Logger.Infof("[CASCADE DETECTOR] searching for cascading deletes originating @ (%s, %s): %s", deleteCall.GetCallerStr(), deleteCall.DbInstance.GetDatastore().GetName(), deleteCall.LongString())
	numServiceCalls := len(lastServiceCallNode.GetChildren())
	for idx := child_idx + 1; idx < numServiceCalls; idx++ {
		node := lastServiceCallNode.GetChildren()[idx]
		// found call that pushes notifications
		if pushCall, ok := node.(*abstractgraph.AbstractDatabaseCall); ok && pushCall.DbInstance.GetDatastore().IsQueue() && pushCall.ParsedCall.Method.IsWrite() {
			logger.Logger.Debugf("[CASCADE DETECTOR] found push call @ (%s, %s): %s", pushCall.GetCallerStr(), pushCall.DbInstance.GetDatastore().GetName(), pushCall.String())
			queueReadNodes := pushCall.GetChildren()
			// check all calls that follow the read of the queue
			for _, queueReadNode := range queueReadNodes {
				// check if after reading the queue, there is a delete operation to the same original database that is being referenced
				for _, child := range queueReadNode.GetChildren() {
					if childDeleteCall, ok := child.(*abstractgraph.AbstractDatabaseCall); ok && childDeleteCall.ParsedCall.Method.IsUpdate() || childDeleteCall.ParsedCall.Method.IsDelete() {
						logger.Logger.Debugf("[CASCADE DETECTOR] found child update/delete call @ (%s, %s): %s", childDeleteCall.Service, childDeleteCall.DbInstance, childDeleteCall.LongString())
						if deleteDep := deleteOp.getDependency(childDeleteCall.Service, childDeleteCall.DbInstance.GetDatastore()); deleteDep != nil {
							logger.Logger.Debugf("[CASCADE DETECTOR] found cascading action!")
							deleteDep.cascading = true
						}
					}
				}
			}
		}
	}
}

func (detector *CascadeDetector) checkInconsistencies() {
	detector.numDeletes = len(detector.getDeleteOperations())

	for i, op := range detector.getDeleteOperations() {
		detector.results += fmt.Sprintf("[%d] %s: %s\n", i+1, op.getCall().GetCallerStr(), op.call.ShortString())
		detector.results += fmt.Sprintf("\tmissing %d cascading deletes\n", len(op.getDependencies()))
		for _, dep := range op.getDependencies() {
			if !dep.cascading {
				detector.results += fmt.Sprintf("\t- %s\n", dep.LongString())
				detector.numMissingCascadingDeletes++
			}
		}
	}
}

func (detector *CascadeDetector) ComputeResults() {
	header := "------------------------------------------------------------\n"
	header += "--------------- FOREIGN KEY CASCADE ANALYSIS ---------------\n"
	header += "------------------------------------------------------------\n"

	detector.checkInconsistencies()

	header += fmt.Sprintf(">> SUMMARY (# DELETES ON REFERENCED OBJECT; # ABSENCE OF CASCADING DELETES):\n>> (%d;%d)\n", detector.numDeletes, detector.numMissingCascadingDeletes)
	detector.results = header + detector.results
}

func (detector *CascadeDetector) GetAnalysisTypeString() string {
	return "foreign_key_cascade"
}

func (detector *CascadeDetector) GetResults() string {
	return detector.results
}
