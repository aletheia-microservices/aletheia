package unicity

import (
	"fmt"
	"slices"

	"github.com/golang-collections/collections/stack"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/datastores"
	"analyzer/pkg/detection/detection"
	"analyzer/pkg/frameworks/blueprint"
	"analyzer/pkg/logger"
	"analyzer/pkg/types/objects"
	"analyzer/pkg/utils"
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
	detection.Detector
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
	app.ResetAllDataflows()
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
	detector.checkInconsistencies()
	app.ResetAllDataflows()
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
	detector.onWriteOrUpdate(app, dbCall, child_idx)
}

func (detector *UnicityDetector) OnUpdate(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	detector.onWriteOrUpdate(app, dbCall, child_idx)
}

func (detector *UnicityDetector) onWriteOrUpdate(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, child_idx int) {
	schema := dbCall.DbInstance.GetDatastore().GetSchema()
	datastore := dbCall.DbInstance.GetDatastore()

	logger.Logger.Debugf("[UNICITY DETECTOR] onWriteOrUpdate: %s", dbCall.String())

	if !schema.HasConstraintsUnique() {
		return
	}

	reqInfo := detector.getCurrentRequestInfo()
	operationIdx := reqInfo.numOperations()
	writtenFieldNames := detection.GetWrittenFieldNamesForOperation(dbCall)
	writtenFields := dbCall.DbInstance.GetDatastore().GetSchema().GetFieldsByNames(writtenFieldNames)
	unicityConstraints := schema.GetConstraintsUniqueForFieldNames(writtenFieldNames)

	operation := NewOperation(operationIdx, dbCall, datastore, writtenFields, unicityConstraints)
	reqInfo.addOperation(operation)

	params := dbCall.GetParams()
	if blueprintBackendMethod := dbCall.GetParsedCall().GetMethod().(*blueprint.BackendMethod); blueprintBackendMethod != nil {
		switch datastore.Type {
		case datastores.Queue, datastores.NoSQL:
			obj := params[1]
			operation.addWrittenObjects(obj)

			abstractgraph.TaintDataflowWrite(app, obj, dbCall, datastore, "", true, child_idx)

		case datastores.Cache:
			key, value := params[1], params[2]
			operation.addWrittenObjects(key, value)

			abstractgraph.TaintDataflowWrite(app, key, dbCall, datastore, datastores.ROOT_FIELD_NAME_CACHE_KEY, false, child_idx)
			abstractgraph.TaintDataflowWrite(app, value, dbCall, datastore, datastores.ROOT_FIELD_NAME_CACHE_VALUE, false, child_idx)

		case datastores.RelationalDB:
			if blueprintBackendMethod.IsRelationalDBExecCall() {
				query, args := params[1], params[2:]
				operation.addWrittenObjects(args...)

				writtenFields, _ := abstractgraph.ParseSQLWrite(query, args)
				for _, field := range writtenFields {
					abstractgraph.TaintDataflowWrite(app, field.GetObject(), dbCall, datastore, field.GetName(), false, child_idx)
				}
			}

		default:
			logger.Logger.Fatalf("[SCHEMA] unknown type of datastore (%s) to parse call: %s", utils.GetType(datastore), dbCall.String())
		}
	}
}

func (detector *UnicityDetector) checkInconsistencies() {
	reqInfo := detector.getCurrentRequestInfo()
	for _, op := range detector.getCurrentRequestInfo().getOperations() {
		for _, obj := range op.getWrittenObjects() {
			if detector.checkConstrainedOperationsWithWriteOnField(op, obj) {
				reqInfo.flagInconsistency()
			}
		}
	}
}

// findConstrainedOperationsWithWriteOnField follows the same logic of abstractgraph.ReferenceTaintedDataflowForNestedField() for finding dataflows
// 1. search for other writes in the same request that used a given field (whose object is being written now)
// 2. for each found write-dataflow, find the corresponding operation saved in the requestInfo of the detector
// 3. if the operation was done against a unicity constraint (not necessarily on the current field), then it can affect our current operation, leading to inconsistencies
func (detector *UnicityDetector) checkConstrainedOperationsWithWriteOnField(currOp *Operation, writtenObj objects.Object) bool {
	var visited []*Operation
	objs := []objects.Object{writtenObj}

	currDb := currOp.getDatastore()
	if currDb.IsQueue() || currDb.IsNoSQLDatabase() {
		objs, _ = objects.GetReversedNestedFieldsAndNames(writtenObj, "", currDb.IsNoSQLDatabase(), currDb.IsQueue())
	}

	for _, obj := range objs {
		deps := obj.GetNestedDependencies(false)
		for _, dep := range deps {
			for _, df := range dep.GetVariableInfo().GetAllWriteDataflowsExceptDatastore(currDb.GetName()) {
				requestInfo := detector.getCurrentRequestInfo()

				for _, otherOp := range requestInfo.getOperations() {
					// ignore if the other operation
					// (1) is the current operation
					// (2) or was already visited
					if currOp == otherOp || slices.Contains(visited, otherOp) {
						continue
					}

					// other operation only affects current operation if the former
					// (1) is a write that uses the current field
					// (2) is constrained to UNIQUE or PK
					if otherOp.hasWrittenField(df.Field) && otherOp.isConstrained() {
						otherOp.addAffectedOpAndReferencedField(currOp, df.Field)
						visited = append(visited, otherOp)
					}
				}
			}

		}
	}
	return len(visited) > 0
}

func (detector *UnicityDetector) OnDelete(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	// no-op
}

func (detector *UnicityDetector) GetAnalysisTypeString() string {
	return "unicity"
}
