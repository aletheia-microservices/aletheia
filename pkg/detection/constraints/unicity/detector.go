package unicity

import (
	"fmt"

	"github.com/golang-collections/collections/stack"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/datastores"
	"analyzer/pkg/detection/detector"
	"analyzer/pkg/logger"
)

func NewDetector() *UnicityDetector {
	fmt.Println()
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println(" ----------------------------------------- INITIALIZING UNICITY DETECTOR ------------------------------------------ ")
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println()

	return &UnicityDetector{
		requestInfoStack: stack.New(),
	}
}

type UnicityDetector struct {
	detector.Detector
	results          string
	summary          string
	requestInfoStack *stack.Stack
}

func (detector *UnicityDetector) GetSummary() string {
	return detector.summary
}

func (detector *UnicityDetector) SetSummary(summary string) {
	detector.summary = summary
}

func (detector *UnicityDetector) getCurrentRequestInfo() *RequestInfo {
	return detector.requestInfoStack.Peek().(*RequestInfo)
}

func (detector *UnicityDetector) OnNewRun(app *app.App) {
	//no-op
}

func (detector *UnicityDetector) OnEndRun(app *app.App) {
	//no-op
}

func (detector *UnicityDetector) OnNewRequest(entryNode *abstractgraph.AbstractServiceCall) {
	detector.requestInfoStack.Push(&RequestInfo{
		entry: entryNode,
	})
}

func (detector *UnicityDetector) OnEndRequest(app *app.App) {
	//no-op
}

func (detector *UnicityDetector) OnNewNode(app *app.App, node abstractgraph.AbstractNode) {
	//no-op
}

func (detector *UnicityDetector) OnEndNode(app *app.App, node abstractgraph.AbstractNode) {
	//no-op
}

func (detector *UnicityDetector) OnRead(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	//no-op
}

func (detector *UnicityDetector) OnWrite(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	detector.onWriteOrUpdate(dbCall)
}

func (detector *UnicityDetector) OnUpdate(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	detector.onWriteOrUpdate(dbCall)
}

func (detector *UnicityDetector) onWriteOrUpdate(dbCall *abstractgraph.AbstractDatabaseCall) {
	schema := dbCall.DbInstance.GetDatastore().GetSchema()
	datastore := dbCall.DbInstance.GetDatastore()
	if schema.HasConstraintsUnique() {
		if datastore.IsNoSQLDatabase() {
			update := dbCall.GetParam(1)
			updateType := update.GetType()
			logger.Logger.Infof("[UNICITY DETECTOR] found WRITE/UPDATE on database (%s)", dbCall.DbInstance.GetName())
			_, fieldNames := updateType.GetNestedFieldTypes(updateType.GetName(), datastore.IsNoSQLDatabase())

			var unicityConstraints []*datastores.Constraint
			for _, fieldName := range fieldNames {
				unicityConstraint := schema.GetConstraintsUniqueForFieldName(fieldName)
				unicityConstraints = append(unicityConstraints, unicityConstraint...)
			}
			logger.Logger.Warnf("[UNICITY DETECTOR] WRITE/UPDATE in (%s) against unicity constraints:", dbCall.DbInstance.GetName())
			for _, uc := range unicityConstraints {
				logger.Logger.Warn("\t\t\t - " + uc.String())
			}

			requestInfo := detector.getCurrentRequestInfo()
			if len(unicityConstraints) > 0 {
				operation := NewOperationOnUnicityConstraint(dbCall, datastore)
				requestInfo.addOperation(operation)
			} else if requestInfo.hasOperations() {
				operation := NewOperation(dbCall, datastore)
				requestInfo.addOperation(operation)
			}
		}
	}
}

func (detector *UnicityDetector) OnDelete(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	// no-op
}

func (detector *UnicityDetector) ComputeResults() {
	header := "------------------------------------------------------------\n"
	header += "---------------------- UNICITY ANALYSIS --------------------\n"
	header += "------------------------------------------------------------\n"

	var numRequests, numOps int

	for detector.requestInfoStack.Len() > 0 {
		requestInfo := detector.requestInfoStack.Pop().(*RequestInfo)
		if requestInfo.hasPotentialInconsistencies() {
			detector.results += "\n[ENTRY] " + requestInfo.entry.String() + "\n"
			numRequests++
			for _, op := range requestInfo.getOperations() {
				detector.results += "\t- OPERATION @ " + op.call.Service + ": " + op.call.String() + "\n"
				numOps++
			}
		}
	}

	header += fmt.Sprintf(">> SUMMARY (# END-TO-END REQUESTS; # AFFECTED OPERATIONS):\n>> (%d;%d)\n", numRequests, numOps)
	detector.results = header + detector.results
}

func (detector *UnicityDetector) GetAnalysisTypeString() string {
	return "unicity"
}

func (detector *UnicityDetector) GetResults() string {
	return detector.results
}
