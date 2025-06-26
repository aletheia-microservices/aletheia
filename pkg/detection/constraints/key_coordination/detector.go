package key_coordination

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

type KeyCoordinationDetector struct {
	detection.Detector
	requestInfoStack *stack.Stack

	keyType string // 'primary_key' or 'foreign_key'

	// results
	results string
	summary string
	reads   []*ForeignKeyRead
}

func NewDetector(keyType string) *KeyCoordinationDetector {
	fmt.Println()
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println(" --------------------------------------- INITIALIZING KEY_COORD DETECTOR ---------------------------------------- ")
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println()
	return &KeyCoordinationDetector{
		keyType:          keyType,
		requestInfoStack: stack.New(),
	}
}

func (detector *KeyCoordinationDetector) keyTypeIsPrimaryKey() bool {
	return detector.keyType == "primary_key"
}

func (detector *KeyCoordinationDetector) GetSummary() string {
	return detector.summary
}

func (detector *KeyCoordinationDetector) SetSummary(summary string) {
	detector.summary = summary
}

func (detector *KeyCoordinationDetector) addForeignKeyRead(read *ForeignKeyRead) {
	detector.reads = append(detector.reads, read)
}

func (detector *KeyCoordinationDetector) getUsedForeignReferencesForFieldInDatastore(fieldName string, datastore *datastores.Datastore) []string {
	var foreignReference []string
	for _, read := range detector.reads {
		if read.refField.Datastore == datastore && read.refField.GetFullName() == fieldName {
			foreignReference = append(foreignReference, read.originField.GetFullName())
		}
	}
	return foreignReference
}

func (detector *KeyCoordinationDetector) OnNewRun(app *app.App) {
	app.ResetAllDataflows()
}

func (detector *KeyCoordinationDetector) OnEndRun(app *app.App) {
	//no-op
}

func (detector *KeyCoordinationDetector) OnNewRequest(entryNode *abstractgraph.AbstractServiceCall) {
	//no-op
}

func (detector *KeyCoordinationDetector) OnNewNode(app *app.App, node abstractgraph.AbstractNode) {
	//no-op
}

func (detector *KeyCoordinationDetector) OnEndNode(app *app.App, node abstractgraph.AbstractNode) {
	//no-op
}

func (detector *KeyCoordinationDetector) OnEndRequest(app *app.App) {
	app.ResetAllDataflows()
}

// FIXME:
// checkKeyRead gets all dependencies for the read object
// for each dependency, it iterates all previous read dataflows
// if the database field read on a previous (dependent) dataflow matches the current field
// then we detect a new foreignkey-based read
func (detector *KeyCoordinationDetector) checkKeyRead(app *app.App, currObj objects.Object, originFieldName string, datastore *datastores.Datastore, dbCall *abstractgraph.AbstractDatabaseCall) {
	currField := datastore.Schema.GetField(originFieldName)
	logger.Logger.Infof("[KEY_COORD 1] check KEY COORD read for origin field (%s) and object: %s", currField.String(), currObj.String())
	var visited []string
	deps := currObj.GetNestedDependencies(true)
	logger.Logger.Debugf("[KEY_COORD 2] \t deps: %v", deps)
	for _, dep := range deps {
		for _, df := range dep.GetVariableInfo().GetAllReadDataflowsExceptDatastore(dbCall.DbInstance.GetName()) { // except filter is just for sanity check
			logger.Logger.Debugf("[KEY_COORD 3.0] \t\t dep: %s", dep.String())
			logger.Logger.Debugf("[KEY_COORD 3.1] \t\t dataflow: %s", df.String())

			otherField := df.Field

			logger.Logger.Debugf("[KEY_COORD 3.2] \t\t other field: %s", otherField.GetFullName())

			// FIXME: we re-assign "otherField" because, for some reason, the original "otherField" is not
			// the one we are expecting with the reference, although in the schema there is a reference
			otherField = otherField.GetDatastore().GetSchema().GetFieldByFullName(otherField.GetFullName())

			logger.Logger.Debugf("[KEY_COORD 3.3] \t\t going to check inconsistencies...")

			var checkInconsistency = func(field1 *datastores.Field, field2 *datastores.Field) {
				// now that we know that the other field references the current one (e.g., Cart.Products references Product.ProductID)
				// we want to know if their association is mandatory (e.g., false for shoppingApp), meaning that they were written in the same request
				
				// this fine-grained approach should prevent a false positive flag in, for example,
				// the shopping_app where products can still be associated afterwards after the cart creation
				// but still allow a true positive flag in the post notification

				// check if there are any primary key constraints for field1
				// if not, then it's useless to search for inconsistencies involving field1
				if detector.keyTypeIsPrimaryKey() {
					if pkeyConstrains := field1.GetConstraints(datastores.ConstraintFilter{Primary: utils.BoolPtr(true)}); pkeyConstrains == nil {
						return
					}
				}
				
				filter := datastores.ConstraintFilter{
					Primary: utils.BoolPtr(false),
					Reference: utils.BoolPtr(true),
					// must be a mandatory reference, meaning that both fields (and corresponding objects) 
					// were presviously written in the same request 
					Mandatory: utils.BoolPtr(true),
				}

				for _, constraintField1 := range field1.GetConstraints(filter) {
					if constraintField1.FieldIsReferencingTo(field2) {

						// just to avoid duplicates
						if !slices.Contains(visited, field1.GetFullName()) {
							read := newForeignKeyRead(field1, field2, app.GetDataflowForObjectDataflow(df).GetDatabaseCall(), dbCall.ParsedCall)
							detector.addForeignKeyRead(read)

							logger.Logger.Debugf("[KEY_COORD] found new KEY COORD read:\n%s", read.String())
							visited = append(visited, field1.GetFullName())
						}
					}
				}
			}

			if otherField.HasReference(currField) {
				// example: postnotification w/ services NotifyService and StorageService
				// workflow: post_id = NotifyService.fetch_notification(notif) --> StorageService.read_post(post_id)
				checkInconsistency(otherField, currField)
			} else if currField.HasReference(otherField) {
				// example: postnotification w/ services StorageService and AnalyticsService
				// workflow: StorageService.read_post(post_id) --> AnalyticsService.read_analytics(post_id)
				// workflow: AnalyticsService.read_analytics(post_id) --> StorageService.read_post(post_id)
				checkInconsistency(currField, otherField)
			}
		}
	}
}

func (detector *KeyCoordinationDetector) OnRead(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	logger.Logger.Infof("[KEY COORD] analyzing %s @ %s: %s", utils.GetType(dbCall), dbCall.DbInstance.GetName(), dbCall.String())
	datastore := dbCall.DbInstance.GetDatastore()
	params := dbCall.GetParams()
	returns := dbCall.GetReturns()

	if blueprintBackendMethod := dbCall.GetParsedCall().GetMethod().(*blueprint.BackendMethod); blueprintBackendMethod != nil {
		switch datastore.Type {
		case datastores.Queue:
			msg := params[1]
			logger.Logger.Infof("[KEY COORD - QUEUE MESSAGE] %s", msg.String())
			for _, df := range msg.GetVariableInfo().GetAllDataflows() {
				logger.Logger.Infof("[df] %s", df.String())
			}
			abstractgraph.TaintDataflowReadQueue(app, msg, dbCall, datastore, datastores.ROOT_FIELD_NAME_QUEUE, child_idx)
		case datastores.NoSQL:
			cursor, query := returns[0], params[1]

			queryObjs := abstractgraph.GetNoSQLQueryDocument(datastore, query)
			for _, qObj := range queryObjs {
				logger.Logger.Infof("[KEY COORD - QUERY OBJ] %s", qObj.String())
				detector.checkKeyRead(app, qObj.Object, qObj.FieldName, datastore, dbCall)
			}

			abstractgraph.TaintDataflowReadNoSQL(app, cursor, dbCall, datastore, datastores.ROOT_FIELD_NAME_NOSQL, false, child_idx)
			for _, obj := range queryObjs {
				logger.Logger.Infof("[KEY COORD - QUERY OBJ] %s", obj.String())
				abstractgraph.TaintDataflowReadNoSQL(app, obj.Object, dbCall, datastore, obj.FieldName, true, child_idx)
			}
		case datastores.Cache:
			key, value := params[1], params[2]

			detector.checkKeyRead(app, key, datastores.ROOT_FIELD_NAME_CACHE_KEY, datastore, dbCall)
			abstractgraph.TaintDataflowReadCache(app, key, datastores.ROOT_FIELD_NAME_CACHE_KEY, dbCall, datastore, child_idx)
			abstractgraph.TaintDataflowReadCache(app, value, datastores.ROOT_FIELD_NAME_CACHE_VALUE, dbCall, datastore, child_idx)

		case datastores.RelationalDB:
			if blueprintBackendMethod.IsRelationalDBSelectCall() {
				targetObj, queryObj, argsObjs := params[1], params[2], params[3:]
				selectedFields, filterFields := abstractgraph.ParseSQLReadSelect(targetObj, queryObj, argsObjs)
				for _, field := range filterFields {
					// fieldName already contains the ROOT FIELD '*' in SQL if needed
					detector.checkKeyRead(app, field.GetObject(), field.GetName(), datastore, dbCall)
					abstractgraph.TaintDataflowReadSQL(app, field.GetObject(), field.GetName(), dbCall, datastore, child_idx, false)
				}

				// fieldName already contains the ROOT FIELD '*' in SQL if needed
				detector.checkKeyRead(app, selectedFields[0].GetObject(), selectedFields[0].GetName(), datastore, dbCall)
				abstractgraph.TaintDataflowReadSQL(app, selectedFields[0].GetObject(), selectedFields[0].GetName(), dbCall, datastore, child_idx, true)
			} else if blueprintBackendMethod.IsRelationalDBQueryCall() {
				logger.Logger.Fatalf("TODO!! implement cursor for sql similar to nosql mongodb")
			}

		default:
			logger.Logger.Fatalf("[KEY COORD] TODO: %s", dbCall.String())
		}
	}
}

func (detector *KeyCoordinationDetector) OnWrite(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	//no-op
}

func (detector *KeyCoordinationDetector) OnUpdate(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	//no-op
}

func (detector *KeyCoordinationDetector) OnDelete(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	//no-op
}

func (detector *KeyCoordinationDetector) GetAnalysisTypeString() string {
	return detector.keyType + "_coordination"
}

func (detector *KeyCoordinationDetector) CompactSchema(app *app.App) {
	for _, ds := range app.Databases {
		for _, unfoldedField := range ds.GetDatastore().Schema.GetAllFields() {
			var refsToKeep []*datastores.Field
			foreignReferences := detector.getUsedForeignReferencesForFieldInDatastore(unfoldedField.GetFullName(), ds.GetDatastore())
			for _, ref := range unfoldedField.GetReferences() {
				if slices.Contains(foreignReferences, ref.GetFullName()) {
					refsToKeep = append(refsToKeep, ref)
				}
			}
			unfoldedField.CompactReferences(refsToKeep)
		}
	}
}
