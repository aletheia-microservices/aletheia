package unicity

import (
	"github.com/golang-collections/collections/stack"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/datastores"
	"analyzer/pkg/logger"
)

func NewDetector() *UnicityDetector {
	return &UnicityDetector{
		requestInfoStack: stack.New(),
	}
}

type UnicityDetector struct {
	Detector
	requestInfoStack *stack.Stack
}

type RequestInfo struct {
	entry      *abstractgraph.AbstractServiceCall
	operations []*Operation
}

func (info *RequestInfo) addOperation(operation *Operation) {
	info.operations = append(info.operations, operation)
}

func (info *RequestInfo) hasOperations() bool {
	return len(info.operations) > 0
}

func (info *RequestInfo) hasPotentialInconsistencies() bool {
	return len(info.operations) > 1 // only if we have more than 2 ops
}

func (info *RequestInfo) getOperations() []*Operation {
	return info.operations
}

type Operation struct {
	call                *abstractgraph.AbstractDatabaseCall
	datastore           *datastores.Datastore
	onUnicityConstraint bool
}

func NewOperation(call *abstractgraph.AbstractDatabaseCall, datastore *datastores.Datastore) *Operation {
	return &Operation{
		call:      call,
		datastore: datastore,
	}
}

func NewOperationOnUnicityConstraint(call *abstractgraph.AbstractDatabaseCall, datastore *datastores.Datastore) *Operation {
	return &Operation{
		call:                call,
		datastore:           datastore,
		onUnicityConstraint: true,
	}
}

func (detector *UnicityDetector) getCurrentRequestInfo() *RequestInfo {
	return detector.requestInfoStack.Peek().(*RequestInfo)
}

func (detector *UnicityDetector) onNewRequest(entryNode *abstractgraph.AbstractServiceCall) {
	detector.requestInfoStack.Push(&RequestInfo{
		entry: entryNode,
	})
}

func (detector *UnicityDetector) onWrite(dbCall *abstractgraph.AbstractDatabaseCall) {
	schema := dbCall.DbInstance.GetDatastore().GetSchema()
	datastore := dbCall.DbInstance.GetDatastore()
	if schema.HasUnicityConstraints() {
		if datastore.IsNoSQLDatabase() {
			doc := dbCall.GetParam(1)
			docType := doc.GetType()
			logger.Logger.Infof("[UNICITY DETECTOR] found WRITE on database (%s)", dbCall.DbInstance.GetName())
			_, fieldNames := docType.GetNestedFieldTypes(docType.GetName(), datastore.IsNoSQLDatabase())

			var unicityConstraints []*datastores.UniqueConstraint
			for _, fieldName := range fieldNames {
				unicityConstraint := schema.GetUnicityConstraintsForFieldName(fieldName)
				unicityConstraints = append(unicityConstraints, unicityConstraint...)
			}
			logger.Logger.Warnf("[UNICITY DETECTOR] WRITE in (%s) against unicity constraints:", dbCall.DbInstance.GetName())
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

func (detector *UnicityDetector) onUpdate(dbCall *abstractgraph.AbstractDatabaseCall) {
	schema := dbCall.DbInstance.GetDatastore().GetSchema()
	datastore := dbCall.DbInstance.GetDatastore()
	if schema.HasUnicityConstraints() {
		if datastore.IsNoSQLDatabase() {
			update := dbCall.GetParam(1)
			updateType := update.GetType()
			logger.Logger.Infof("[UNICITY DETECTOR] found UPDATE on database (%s)", dbCall.DbInstance.GetName())
			_, fieldNames := updateType.GetNestedFieldTypes(updateType.GetName(), datastore.IsNoSQLDatabase())

			var unicityConstraints []*datastores.UniqueConstraint
			for _, fieldName := range fieldNames {
				unicityConstraint := schema.GetUnicityConstraintsForFieldName(fieldName)
				unicityConstraints = append(unicityConstraints, unicityConstraint...)
			}
			logger.Logger.Warnf("[UNICITY DETECTOR] UPDATE in (%s) against unicity constraints:", dbCall.DbInstance.GetName())
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

func (detector *UnicityDetector) onDelete(dbCall *abstractgraph.AbstractDatabaseCall) {
	// no-op
}

func (detector *UnicityDetector) Results() string {
	var res string
	res += "[UNICITY DETECTOR] ========= RESULTS =========\n"
	for detector.requestInfoStack.Len() > 0 {
		requestInfo := detector.requestInfoStack.Pop().(*RequestInfo)
		if requestInfo.hasPotentialInconsistencies() {
			res += "\n[ENTRY] " + requestInfo.entry.String() + "\n"
			for _, op := range requestInfo.getOperations() {
				res += "\t- OPERATION @ " + op.call.Service + ": " + op.call.String() + "\n"
			}
		}
	}
	return res
}
