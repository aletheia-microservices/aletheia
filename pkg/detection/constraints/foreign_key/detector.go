package foreign_key

import (
	"fmt"
	"slices"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/datastores"
	"analyzer/pkg/detection/detector"
	"analyzer/pkg/logger"
	"analyzer/pkg/types/objects"
	"analyzer/pkg/utils"
)

type ForeignKeyDetector struct {
	detector.Detector
	results string
	reads   []*ForeignKeyRead
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

func NewDetector() *ForeignKeyDetector {
	fmt.Println()
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println(" --------------------------------------- INITIALIZING FOREIGN KEY DETECTOR ---------------------------------------- ")
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println()
	return &ForeignKeyDetector{}
}

func (detector *ForeignKeyDetector) OnRun(app *app.App) {
	app.ResetAllDataflows()
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
func (detector *ForeignKeyDetector) checkForeignKeyRead(app *app.App, obj objects.Object, originField *datastores.Entry, dbCall *abstractgraph.AbstractDatabaseCall) {
	logger.Logger.Infof("[FOREIGN KEY] check foreign key read for origin field (%s) and object: %s", originField.String(), obj.String())
	var savedOriginFieldName []string
	//datastore := dbCall.DbInstance.GetDatastore()
	for _, dep := range obj.GetNestedDependencies(true) {
		logger.Logger.Debugf("[FOREIGN KEY] \t dep: %s", dep.String())
		for _, df := range dep.GetVariableInfo().GetAllReadDataflowsExceptDatastore(dbCall.DbInstance.GetName()) { // except filter is just for sanity check
			//FIXME: for some reason there are some "loose" fields that
			// are not associated anymore with the (un)folded fields of the datastore schema
			// and which also unexpectedly do not have any References
			// so right now we are just getting the full name of the field that we want for the current dataflow
			// and then we get the field that is actually attached to the schema to get the correct References
			refField := df.Field.(*datastores.Entry)
			for _, field := range df.Field.GetDatastore().Schema.UnfoldedFields {
				if field.GetFullName() == refField.GetFullName() {
					attachedRefField := field.(*datastores.Entry)
					for _, refTarget := range attachedRefField.References {
						if !slices.Contains(savedOriginFieldName, originField.GetFullName()) && refTarget == originField {
							read := newForeignKeyRead(attachedRefField, originField, app.GetDataflowForObjectDataflow(df).GetDatabaseCall(), dbCall.ParsedCall)
							detector.addForeignKeyRead(read)
							savedOriginFieldName = append(savedOriginFieldName, originField.GetFullName())
						}
					}
				}
			}
		}
	}
}

func (detector *ForeignKeyDetector) OnRead(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	logger.Logger.Infof("[FOREIGN KEY] analyzing %s @ %s: %s", utils.GetType(dbCall), dbCall.DbInstance.GetName(), dbCall.String())
	datastore := dbCall.DbInstance.GetDatastore()
	params := dbCall.GetParams()
	returns := dbCall.GetReturns()

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
			field := datastore.Schema.GetField(qObj.FieldName).(*datastores.Entry)
			detector.checkForeignKeyRead(app, qObj.Object, field, dbCall)
		}

		abstractgraph.TaintDataflowNoSQL(app, cursor, dbCall, datastore, datastores.ROOT_FIELD_NAME_NOSQL, false, false, child_idx)
		for _, obj := range queryObjs {
			logger.Logger.Infof("[FOREIGN KEY - QUERY OBJ] %s", obj.String())
			abstractgraph.TaintDataflowNoSQL(app, obj.Object, dbCall, datastore, obj.FieldName, true, false, child_idx)
		}
	case datastores.Cache:
		key, value := params[1], params[2]

		field := datastore.Schema.GetField(datastores.ROOT_FIELD_NAME_CACHE_KEY).(*datastores.Entry)
		detector.checkForeignKeyRead(app, key, field, dbCall)

		abstractgraph.TaintDataflowReadCache(app, key, datastores.ROOT_FIELD_NAME_CACHE_KEY, dbCall, datastore, child_idx)
		abstractgraph.TaintDataflowReadCache(app, value, datastores.ROOT_FIELD_NAME_CACHE_VALUE, dbCall, datastore, child_idx)

	default:
		logger.Logger.Fatalf("[FOREIGN KEY] TODO: %s", dbCall.String())
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

func (detector *ForeignKeyDetector) ComputeResults() {
	detector.results = "------------------------------------------------------------\n"
	detector.results += "------------------- FOREIGN KEY ANALYSIS -------------------\n"
	detector.results += "------------------------------------------------------------\n"
	for i, read := range detector.reads {
		detector.results += fmt.Sprintf("foreign key read #%d:\n%s\n", i, read.String())
		if i < len(detector.reads)-1 {
			detector.results += "\n" // enforce empty line between each foreign key read result
		}
	}
}

func (detector *ForeignKeyDetector) GetAnalysisTypeString() string {
	return "foreign_key"
}

func (detector *ForeignKeyDetector) GetResults() string {
	return detector.results
}

func (detector *ForeignKeyDetector) CompactSchema(app *app.App) {
	for _, ds := range app.Databases {
		for _, unfoldedField := range ds.GetDatastore().Schema.GetAllFields() {
			var refsToKeep []datastores.Field
			foreignReferences := detector.getUsedForeignReferencesForFieldInDatastore(unfoldedField.GetFullName(), ds.GetDatastore())
			for _, ref := range unfoldedField.(*datastores.Entry).References {
				if slices.Contains(foreignReferences, ref.GetFullName()) {
					refsToKeep = append(refsToKeep, ref)
				}
			}
			unfoldedField.(*datastores.Entry).References = refsToKeep
		}
	}
}
