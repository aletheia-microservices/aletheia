package numerical

import (
	"fmt"

	"github.com/golang-collections/collections/stack"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/datastores"
	"analyzer/pkg/detection/detector"
	"analyzer/pkg/frameworks/blueprint"
	"analyzer/pkg/logger"
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
	detector.Detector
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
	detector.onWriteOrUpdate(dbCall)
}

func (detector *NumericalDetector) OnUpdate(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	detector.onWriteOrUpdate(dbCall)
}

func (detector *NumericalDetector) onWriteOrUpdate(dbCall *abstractgraph.AbstractDatabaseCall) {
	schema := dbCall.DbInstance.GetDatastore().GetSchema()
	datastore := dbCall.DbInstance.GetDatastore()
	if schema.HasConstraintsNumerical() {
		if blueprintBackendMethod := dbCall.ParsedCall.Method.(*blueprint.BackendMethod); blueprintBackendMethod != nil {
			var writtenFieldNames []string
			params := dbCall.GetParams()

			switch datastore.Type {
			case datastores.NoSQL:
				obj := params[1]
				objType := obj.GetType()
				logger.Logger.Infof("[NUMERICAL DETECTOR] found WRITE/UPDATE on database (%s)", dbCall.DbInstance.GetName())
				_, writtenFieldNames = objType.GetNestedFieldTypes(objType.GetName(), datastore.IsNoSQLDatabase())
			case datastores.RelationalDB:
				if blueprintBackendMethod.IsRelationalDBExecCall() {
					query, args := params[1], params[2:]
					writtenFields, _ := abstractgraph.ParseSQLWrite(query, args)
					for _, field := range writtenFields {
						writtenFieldNames = append(writtenFieldNames, field.GetName())
					}
				} else {
					logger.Logger.Fatalf("[NUMERICAL DETECTOR] TODO on write/update for RelationalDB (%s) call: %s", datastore.GetName(), dbCall.LongString())
				}
			default:
				logger.Logger.Fatalf("[NUMERICAL DETECTOR] TODO on write/update for datastore (%s) call: %s", datastore.GetName(), dbCall.LongString())
			}

			var numericalConstraints []*datastores.Constraint
			for _, writtenFieldName := range writtenFieldNames {
				numericalConstraint := schema.GetConstraintsNumericalForFieldName(writtenFieldName)
				numericalConstraints = append(numericalConstraints, numericalConstraint...)
			}
			logger.Logger.Warnf("[NUMERICAL DETECTOR] WRITE/UPDATE in (%s) against numerical constraints:", dbCall.DbInstance.GetName())
			for _, uc := range numericalConstraints {
				logger.Logger.Warn("\t\t\t - " + uc.String())
			}

			requestInfo := detector.getCurrentRequestInfo()
			if len(numericalConstraints) > 0 { // operation that may be affecting the following ones
				operation := NewOperationOnNumericalConstraint(dbCall, datastore)
				requestInfo.addOperation(operation)
			} else if requestInfo.hasOperations() { // operation that is affected by previous operations
				operation := NewOperation(dbCall, datastore)
				requestInfo.addOperation(operation)
			}
		}
	}
}

func (detector *NumericalDetector) OnDelete(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	// no-op
}

func (detector *NumericalDetector) ComputeResults() {
	header := "------------------------------------------------------------\n"
	header += "--------------------- NUMERICAL ANALYSIS -------------------\n"
	header += "------------------------------------------------------------\n"

	var numRequests, numOps int

	for detector.requestInfoStack.Len() > 0 {
		requestInfo := detector.requestInfoStack.Pop().(*RequestInfo)
		if requestInfo.hasPotentialInconsistencies() {
			detector.results += fmt.Sprintf("\n[ENTRY] %s\n", requestInfo.entry.String())
			numRequests++
			for _, op := range requestInfo.getOperations() {
				if op.onNumericalConstraint {
					detector.results += "\t* "
				} else {
					detector.results += "\t- "
				}
				detector.results += fmt.Sprintf("(%s, %s) -> %s\n", op.call.Service, op.datastore.GetName(), op.call.String())
				numOps++
			}
		}
	}

	header += fmt.Sprintf(">> SUMMARY (# END-TO-END REQUESTS; # AFFECTED OPERATIONS):\n>> (%d;%d)\n", numRequests, numOps)
	detector.results = header + detector.results
}

func (detector *NumericalDetector) GetAnalysisTypeString() string {
	return "numerical"
}

func (detector *NumericalDetector) GetResults() string {
	return detector.results
}
