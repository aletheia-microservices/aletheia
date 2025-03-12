package unicity

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

	if blueprintBackendMethod := dbCall.ParsedCall.Method.(*blueprint.BackendMethod); blueprintBackendMethod != nil {
		if schema.HasConstraintsUnique() {
			var writtenFieldNames []string
			params := dbCall.GetParams()
			
			switch datastore.Type {
			case datastores.NoSQL:
				obj := params[1]
				objType := obj.GetType()
				logger.Logger.Infof("[UNICITY DETECTOR] found WRITE/UPDATE on database (%s)", dbCall.DbInstance.GetName())

				// FIXME: this is getting all fields for the structure and not actually the ones that are written
				_, writtenFieldNames = objType.GetNestedFieldTypes(objType.GetName(), datastore.IsNoSQLDatabase())
				logger.Logger.Debugf("[UNICITY DETECTOR] got written fieldnames: %v", writtenFieldNames)

			case datastores.RelationalDB:
				if blueprintBackendMethod.IsRelationalDBExecCall() {
					query, args := params[1], params[2:]
					writtenFields, _ := abstractgraph.ParseSQLWrite(query, args)
					for _, field := range writtenFields {
						writtenFieldNames = append(writtenFieldNames, field.GetName())
					}
					logger.Logger.Debugf("[UNICITY DETECTOR] got written fields: %v", writtenFields)
				} else {
					logger.Logger.Fatalf("[UNICITY DETECTOR] TODO on write/update for RelationalDB (%s) call: %s", datastore.GetName(), dbCall.LongString())
				}
			default:
				logger.Logger.Fatalf("[UNICITY DETECTOR] TODO on write/update for datastore (%s) call: %s", datastore.GetName(), dbCall.LongString())
			}

			unicityConstraints := schema.GetConstraintsUniqueForFieldNames(writtenFieldNames)

			logger.Logger.Warnf("[UNICITY DETECTOR] WRITE/UPDATE in (%s) against unicity constraints:", dbCall.DbInstance.GetName())
			for _, uc := range unicityConstraints {
				logger.Logger.Warn("\t\t\t - " + uc.String())
			}

			requestInfo := detector.getCurrentRequestInfo()
			if len(unicityConstraints) > 0 {
				operation := NewOperationOnUnicityConstraint(dbCall, datastore, unicityConstraints)
				requestInfo.addOperation(operation)
				requestInfo.writeOnConstraint = true

			} else {
				operation := NewOperation(dbCall, datastore)
				requestInfo.addOperation(operation)
			}
		}
	}
}

func (detector *UnicityDetector) OnDelete(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	// no-op
}

func (detector *UnicityDetector) GetAnalysisTypeString() string {
	return "unicity"
}
