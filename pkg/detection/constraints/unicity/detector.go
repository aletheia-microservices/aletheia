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

	writtenFieldNames := detection.GetWrittenFieldNamesForOperation(dbCall)
	unicityConstraints := schema.GetConstraintsUniqueForFieldNames(writtenFieldNames)
	writtenFields := dbCall.DbInstance.GetDatastore().GetSchema().GetFieldsByNames(writtenFieldNames)

	/* logger.Logger.Debugf("[UNICITY DETECTOR] written field names:")
	for _, f := range writtenFields {
		logger.Logger.Debugf("\t\t\t - %s", f.Name)
	}

	logger.Logger.Warnf("[UNICITY DETECTOR] WRITE/UPDATE in (%s) against unicity constraints:", dbCall.DbInstance.GetName())
	for _, uc := range unicityConstraints {
		logger.Logger.Warnf("\t\t\t - %s", uc.String())
	} */

	var found int
	reqInfo := detector.getCurrentRequestInfo()
	operation := NewOperation(reqInfo.numOperations(), dbCall, datastore, writtenFields, unicityConstraints)

	params := dbCall.GetParams()
	if blueprintBackendMethod := dbCall.GetParsedCall().GetMethod().(*blueprint.BackendMethod); blueprintBackendMethod != nil {
		switch datastore.Type {
		case datastores.Queue, datastores.NoSQL:
			obj := params[1]
			abstractgraph.TaintDataflowWrite(app, obj, dbCall, datastore, "", true, child_idx)
			found += detector.findConstrainedOperationsWithWriteOnField(operation, datastore, obj)

		case datastores.Cache:
			key, value := params[1], params[2]
			abstractgraph.TaintDataflowWrite(app, key, dbCall, datastore, datastores.ROOT_FIELD_NAME_CACHE_KEY, false, child_idx)
			abstractgraph.TaintDataflowWrite(app, value, dbCall, datastore, datastores.ROOT_FIELD_NAME_CACHE_VALUE, false, child_idx)
			found += detector.findConstrainedOperationsWithWriteOnField(operation, datastore, key)
			found += detector.findConstrainedOperationsWithWriteOnField(operation, datastore, value)

		case datastores.RelationalDB:
			if blueprintBackendMethod.IsRelationalDBExecCall() {
				query, args := params[1], params[2:]
				writtenFields, _ := abstractgraph.ParseSQLWrite(query, args)
				for _, field := range writtenFields {
					abstractgraph.TaintDataflowWrite(app, field.GetObject(), dbCall, datastore, field.GetName(), false, child_idx)
					found += detector.findConstrainedOperationsWithWriteOnField(operation, datastore, field.GetObject())
				}
			}

		default:
			logger.Logger.Fatalf("[SCHEMA] unknown type of datastore (%s) to parse call: %s", utils.GetType(datastore), dbCall.String())
		}
	}

	reqInfo.addOperation(operation)
	if found > 0 {
		reqInfo.affectedOps = true
	}
}

// findConstrainedOperationsWithWriteOnField follows the same logic of abstractgraph.ReferenceTaintedDataflowForNestedField() for finding dataflows
// 1. search for previous writes in the same request that used a given field (whose object is being written now)
// 2. for each found write-dataflow, find the corresponding operation saved in the requestInfo of the detector
// 3. if the operation was done against a unicity constraint (not necessarily on the current field), then it can affect our current operation, leading to inconsistencies
func (detector *UnicityDetector) findConstrainedOperationsWithWriteOnField(currOp *Operation, datastore *datastores.Datastore, writtenObj objects.Object) int {
	var foundOps []*Operation
	objs := []objects.Object{writtenObj}
	if datastore.IsQueue() || datastore.IsNoSQLDatabase() {
		objs, _ = objects.GetReversedNestedFieldsAndNames(writtenObj, "", datastore.IsNoSQLDatabase(), datastore.IsQueue())
	}
	
	for _, obj := range objs {
		deps := obj.GetNestedDependencies(false)
		for _, dep := range deps {
			for _, df := range dep.GetVariableInfo().GetAllWriteDataflowsExceptDatastore(datastore.GetName()) {
				requestInfo := detector.getCurrentRequestInfo()
	
				for _, prevOp := range requestInfo.getOperations() {
					if slices.Contains(foundOps, prevOp) {
						continue
					}

					if prevOp.HasWrittenField(df.Field) && prevOp.onUnicityConstraint {
						prevOp.AddAffectedOpAndReferencedField(currOp, df.Field)
						foundOps = append(foundOps, prevOp)
					}
				}
			}
	
		}
	}
	logger.Logger.Debugf("LEN = %d", len(foundOps))
	return len(foundOps)
}

func (detector *UnicityDetector) OnDelete(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	// no-op
}

func (detector *UnicityDetector) GetAnalysisTypeString() string {
	return "unicity"
}
