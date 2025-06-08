package foreign_key_read

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

type ForeignKeyDetector struct {
	detection.Detector
	requestInfoStack *stack.Stack

	// results
	results string
	summary string
	reads   []*ForeignKeyRead
}

func NewDetector() *ForeignKeyDetector {
	fmt.Println()
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println(" --------------------------------------- INITIALIZING FOREIGN KEY DETECTOR ---------------------------------------- ")
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println()
	return &ForeignKeyDetector{
		requestInfoStack: stack.New(),
	}
}

func (detector *ForeignKeyDetector) GetSummary() string {
	return detector.summary
}

func (detector *ForeignKeyDetector) SetSummary(summary string) {
	detector.summary = summary
}

func (detector *ForeignKeyDetector) addForeignKeyRead(read *ForeignKeyRead) {
	detector.reads = append(detector.reads, read)
}

func (detector *ForeignKeyDetector) getUsedForeignReferencesForFieldInDatastore(fieldName string, datastore *datastores.Datastore) []string {
	var foreignReference []string
	for _, read := range detector.reads {
		if read.refField.Datastore == datastore && read.refField.GetFullName() == fieldName {
			foreignReference = append(foreignReference, read.originField.GetFullName())
		}
	}
	return foreignReference
}

func (detector *ForeignKeyDetector) OnNewRun(app *app.App) {
	app.ResetAllDataflows()
}

func (detector *ForeignKeyDetector) OnEndRun(app *app.App) {
	//no-op
}

func (detector *ForeignKeyDetector) OnNewRequest(entryNode *abstractgraph.AbstractServiceCall) {
	//no-op
}

func (detector *ForeignKeyDetector) OnNewNode(app *app.App, node abstractgraph.AbstractNode) {
	//no-op
}

func (detector *ForeignKeyDetector) OnEndNode(app *app.App, node abstractgraph.AbstractNode) {
	//no-op
}

func (detector *ForeignKeyDetector) OnEndRequest(app *app.App) {
	app.ResetAllDataflows()
}

// FIXME:
// checkForeignKeyRead gets all dependencies for the read object
// for each dependency, it iterates all previous read dataflows
// if the database field read on a previous (dependent) dataflow matches the current field
// then we detect a new foreignkey-based read
func (detector *ForeignKeyDetector) checkForeignKeyRead(app *app.App, currObj objects.Object, originFieldName string, datastore *datastores.Datastore, dbCall *abstractgraph.AbstractDatabaseCall) {
	currField := datastore.Schema.GetField(originFieldName)
	logger.Logger.Infof("[FK READ] check foreign key read for origin field (%s) and object: %s", currField.String(), currObj.String())
	var visited []string
	for _, dep := range currObj.GetNestedDependencies(true) {
		logger.Logger.Debugf("[FK READ] \t dep: %s", dep.String())
		for _, df := range dep.GetVariableInfo().GetAllReadDataflowsExceptDatastore(dbCall.DbInstance.GetName()) { // except filter is just for sanity check
			otherField := df.Field

			// FIXME: we re-assign "otherField" because, for some reason, the original "otherField" is not
			// the one we are expecting with the reference, although in the schema there is a reference
			otherField = otherField.GetDatastore().GetSchema().GetFieldByFullName(otherField.GetFullName())

			var checkInconsistency = func(field1 *datastores.Field, field2 *datastores.Field) {
				// now that we know that the other field references the current one (e.g., Cart.Products references Product.ProductID)
				// we want to know if their association is mandatory (false in this app), meaning that they were written in the same request
				// NOTE: this fine-grained approach should prevent a false positive flag in, for example,
				// the shopping_app where products can still be associated afterwards after the cart creation
				// but still allow a true positive flag in the post notification
				for _, constraint := range field1.GetConstraints(datastores.ConstraintFilter{Reference: utils.BoolPtr(true), Mandatory: utils.BoolPtr(true)}) {
					if constraint.FieldIsReferencingTo(field2) {
						if !slices.Contains(visited, field1.GetFullName()) {
							read := newForeignKeyRead(field1, field2, app.GetDataflowForObjectDataflow(df).GetDatabaseCall(), dbCall.ParsedCall)
							detector.addForeignKeyRead(read)

							logger.Logger.Warnf("[FK READ] found new foreign key read:\n%s", read.String())
							visited = append(visited, field1.GetFullName())
						}
					}
				}
			}

			if otherField.HasReference(currField) {
				// EXAMPLE: read(NOTIFICATION) (w/ other: fk on postid) ... read(POST) (w/ curr: postid)
				checkInconsistency(otherField, currField)
			} else if currField.HasReference(otherField) {
				// EXAMPLE: read(post) (w/ other: postid) ... read(analytics) (w/ curr: fk on postid)
				// THIS STILL NEEDS TO BE TESTED!!
				checkInconsistency(currField, otherField)
			}
		}
	}
}

func (detector *ForeignKeyDetector) OnRead(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	logger.Logger.Infof("[FK READ] analyzing %s @ %s: %s", utils.GetType(dbCall), dbCall.DbInstance.GetName(), dbCall.String())
	datastore := dbCall.DbInstance.GetDatastore()
	params := dbCall.GetParams()
	returns := dbCall.GetReturns()

	if blueprintBackendMethod := dbCall.GetParsedCall().GetMethod().(*blueprint.BackendMethod); blueprintBackendMethod != nil {
		switch datastore.Type {
		case datastores.Queue:
			msg := params[1]
			logger.Logger.Infof("[FOREIGN KEY - QUEUE MESSAGE] %s", msg.String())
			for _, df := range msg.GetVariableInfo().GetAllDataflows() {
				logger.Logger.Infof("[df] %s", df.String())
			}
			abstractgraph.TaintDataflowReadQueue(app, msg, dbCall, datastore, datastores.ROOT_FIELD_NAME_QUEUE, child_idx)
		case datastores.NoSQL:
			cursor, query := returns[0], params[1]

			queryObjs := abstractgraph.GetNoSQLQueryDocument(datastore, query)
			for _, qObj := range queryObjs {
				logger.Logger.Infof("[FOREIGN KEY - QUERY OBJ] %s", qObj.String())
				detector.checkForeignKeyRead(app, qObj.Object, qObj.FieldName, datastore, dbCall)
			}

			abstractgraph.TaintDataflowReadNoSQL(app, cursor, dbCall, datastore, datastores.ROOT_FIELD_NAME_NOSQL, false, child_idx)
			for _, obj := range queryObjs {
				logger.Logger.Infof("[FOREIGN KEY - QUERY OBJ] %s", obj.String())
				abstractgraph.TaintDataflowReadNoSQL(app, obj.Object, dbCall, datastore, obj.FieldName, true, child_idx)
			}
		case datastores.Cache:
			key, value := params[1], params[2]

			detector.checkForeignKeyRead(app, key, datastores.ROOT_FIELD_NAME_CACHE_KEY, datastore, dbCall)
			abstractgraph.TaintDataflowReadCache(app, key, datastores.ROOT_FIELD_NAME_CACHE_KEY, dbCall, datastore, child_idx)
			abstractgraph.TaintDataflowReadCache(app, value, datastores.ROOT_FIELD_NAME_CACHE_VALUE, dbCall, datastore, child_idx)

		case datastores.RelationalDB:
			if blueprintBackendMethod.IsRelationalDBSelectCall() {
				targetObj, queryObj, argsObjs := params[1], params[2], params[3:]
				selectedFields, filterFields := abstractgraph.ParseSQLReadSelect(targetObj, queryObj, argsObjs)
				for _, field := range filterFields {
					// fieldName already contains the ROOT FIELD '*' in SQL if needed
					detector.checkForeignKeyRead(app, field.GetObject(), field.GetName(), datastore, dbCall)
					abstractgraph.TaintDataflowReadSQL(app, field.GetObject(), field.GetName(), dbCall, datastore, child_idx, false)
				}

				// fieldName already contains the ROOT FIELD '*' in SQL if needed
				detector.checkForeignKeyRead(app, selectedFields[0].GetObject(), selectedFields[0].GetName(), datastore, dbCall)
				abstractgraph.TaintDataflowReadSQL(app, selectedFields[0].GetObject(), selectedFields[0].GetName(), dbCall, datastore, child_idx, true)
			} else if blueprintBackendMethod.IsRelationalDBQueryCall() {
				logger.Logger.Fatalf("TODO!! implement cursor for sql similar to nosql mongodb")
			}

		default:
			logger.Logger.Fatalf("[FK READ] TODO: %s", dbCall.String())
		}
	}
}

func (detector *ForeignKeyDetector) OnWrite(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	//no-op
}

func (detector *ForeignKeyDetector) OnUpdate(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	//no-op
}

func (detector *ForeignKeyDetector) OnDelete(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	//no-op
}

func (detector *ForeignKeyDetector) GetAnalysisTypeString() string {
	return "foreign_key_read"
}

func (detector *ForeignKeyDetector) CompactSchema(app *app.App) {
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
