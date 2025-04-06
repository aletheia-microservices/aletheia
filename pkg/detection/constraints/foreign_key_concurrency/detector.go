package foreign_key_concurrency

import (
	"fmt"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/datastores"
	"analyzer/pkg/detection/detection"
	"analyzer/pkg/logger"
)

type ForeignKeyConcurrencyDetector struct {
	writesWithReference map[*datastores.Datastore][]*writeWithReference
	deletes             map[*datastores.Datastore][]*delete

	// results
	results                          string
	summary                          string
	numWritesWithAffectedConstraints int
	numConstraintsAffected           int
}

func (detector *ForeignKeyConcurrencyDetector) addWriteWithReference(ds *datastores.Datastore, w *writeWithReference) {
	detector.writesWithReference[ds] = append(detector.writesWithReference[ds], w)
}

func (detector *ForeignKeyConcurrencyDetector) getFirstDeleteToDatastoreIfExists(ds *datastores.Datastore) *delete {
	if dels, exists := detector.deletes[ds]; exists {
		return dels[0]
	}
	return nil
}

func (detector *ForeignKeyConcurrencyDetector) hasDeletesToDatastore(ds *datastores.Datastore) bool {
	_, exists := detector.deletes[ds]
	return exists
}

func (detector *ForeignKeyConcurrencyDetector) addDelete(ds *datastores.Datastore, d *delete) {
	detector.deletes[ds] = append(detector.deletes[ds], d)
}

func (detector *ForeignKeyConcurrencyDetector) GetSummary() string {
	return detector.summary
}

func (detector *ForeignKeyConcurrencyDetector) SetSummary(summary string) {
	detector.summary = summary
}

func NewDetector() *ForeignKeyConcurrencyDetector {
	fmt.Println()
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println(" ------------------------------------ INITIALIZING FOREIGN KEY CONCURRENCY -------------------------------------- ")
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println()
	return &ForeignKeyConcurrencyDetector{
		deletes: make(map[*datastores.Datastore][]*delete),
		writesWithReference: make(map[*datastores.Datastore][]*writeWithReference),
	}
}

func (detector *ForeignKeyConcurrencyDetector) OnNewRun(app *app.App) {
	app.ResetAllDataflows()
}

func (detector *ForeignKeyConcurrencyDetector) OnEndRun(app *app.App) {
	//no-op
}

func (detector *ForeignKeyConcurrencyDetector) OnNewRequest(entryNode *abstractgraph.AbstractServiceCall) {
	//no-op
}

func (detector *ForeignKeyConcurrencyDetector) OnEndRequest(app *app.App) {
	app.ResetAllDataflows()
}

func (detector *ForeignKeyConcurrencyDetector) OnNewNode(app *app.App, node abstractgraph.AbstractNode) {
	//no-op
}

func (detector *ForeignKeyConcurrencyDetector) OnEndNode(app *app.App, node abstractgraph.AbstractNode) {
	//no-op
}

func (detector *ForeignKeyConcurrencyDetector) OnRead(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	//no-op
}

func (detector *ForeignKeyConcurrencyDetector) OnWrite(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	detector.onWriteOrUpdate(dbCall)
}

func (detector *ForeignKeyConcurrencyDetector) OnUpdate(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	detector.onWriteOrUpdate(dbCall)
}

func (detector *ForeignKeyConcurrencyDetector) OnDelete(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	logger.Logger.Debugf("ON DELETE!")

	datastore := dbCall.DbInstance.GetDatastore()
	detector.addDelete(datastore, &delete{
		call:      dbCall,
		datastore: datastore,
	})

}

func (detector *ForeignKeyConcurrencyDetector) onWriteOrUpdate(dbCall *abstractgraph.AbstractDatabaseCall) {
	datastore := dbCall.DbInstance.GetDatastore()
	schema := datastore.GetSchema()

	if schema.HasConstraintsForeignKey() {
		writtenFieldNames := detection.GetWrittenFieldNamesForOperation(dbCall)
		foreignKeyConstraints := schema.GetConstraintsForeignKeyForFieldNames(writtenFieldNames)

		var fieldsWithReference []*fieldWithReference
		for field, constraints := range foreignKeyConstraints {
			fieldsWithReference = append(fieldsWithReference, &fieldWithReference{
				field:       field,
				constraints: constraints,
			})
		}
		write := &writeWithReference{
			call:      dbCall,
			datastore: datastore,
			fields:    fieldsWithReference,
		}

		detector.addWriteWithReference(datastore, write)
	}
}

func (detector *ForeignKeyConcurrencyDetector) checkInconsistencies() {
	for _, writes := range detector.writesWithReference {
		for _, w := range writes {
			var tmpResults string
			for _, f := range w.fields {
	
				datastoresBeingReferenced := make(map[*datastores.Datastore]bool)
				for _, r := range f.field.GetReferences() {
					datastoresBeingReferenced[r.GetDatastore()] = true
				}
	
				for ds := range datastoresBeingReferenced {
					if delete := detector.getFirstDeleteToDatastoreIfExists(ds); delete != nil {
						tmpResults += fmt.Sprintf("- %s\n", f.field.GetFullName())
						tmpResults += fmt.Sprintf("\t- foreign key:\t%s\n", f.field.GetReferenceForDatastore(ds.GetName()).GetFullName())
						tmpResults += fmt.Sprintf("\t- delete:\t%s\n", delete.call.ShortString())
						detector.numConstraintsAffected++
					}
				}
			}
			if tmpResults != "" {
				detector.numWritesWithAffectedConstraints++
				detector.results += fmt.Sprintf("[#%d] write with referential constraints affected by deletes\n", detector.numWritesWithAffectedConstraints)
				detector.results += tmpResults
			}
		}
	}
}


func (detector *ForeignKeyConcurrencyDetector) ComputeResults() {
	header := "------------------------------------------------------------\n"
	header += "------------ FOREIGN KEY CONCURRENCY ANALYSIS --------------\n"
	header += "------------------------------------------------------------\n"

	detector.checkInconsistencies()
	header += fmt.Sprintf(">> SUMMARY (# WRITES; # CONSTRAINTS AFFECTED BY DELETES):\n>> (%d;%d)\n", detector.numWritesWithAffectedConstraints, detector.numConstraintsAffected)
	detector.results = header + detector.results
}

func (detector *ForeignKeyConcurrencyDetector) GetAnalysisTypeString() string {
	return "foreign_key_concurrency"
}

func (detector *ForeignKeyConcurrencyDetector) GetResults() string {
	return detector.results
}
