package foreign_key_concurrency

import (
	"fmt"

	"github.com/golang-collections/collections/stack"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/datastores"
	"analyzer/pkg/detection/detection"
	"analyzer/pkg/utils"
)

type IterationPhase int

const (
	IterationPhaseCheckDeletes IterationPhase = iota
	IterationPhaseCheckWritesAndUpdates
)

type ForeignKeyConcurrencyDetector struct {
	iter    IterationPhase
	deletes map[*datastores.Datastore][]*delete

	// to be used later
	requestInfoStack *stack.Stack

	// results
	results                  string
	summary                  string
	numDeletes               int
	numAffectedWrittenFields int
}

func (detector *ForeignKeyConcurrencyDetector) NextIterationPhase() {
	if detector.iter == IterationPhaseCheckDeletes {
		detector.iter = IterationPhaseCheckWritesAndUpdates
	}
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

func (detector *ForeignKeyConcurrencyDetector) getNewRequestInfoIndex() int {
	return detector.requestInfoStack.Len()
}

func NewDetector() *ForeignKeyConcurrencyDetector {
	fmt.Println()
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println(" ------------------------------------ INITIALIZING FOREIGN KEY CONCURRENCY -------------------------------------- ")
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println()
	return &ForeignKeyConcurrencyDetector{
		requestInfoStack: stack.New(),
		deletes:          make(map[*datastores.Datastore][]*delete),
	}
}

func (detector *ForeignKeyConcurrencyDetector) OnNewRun(app *app.App) {
	app.ResetAllDataflows()
}

func (detector *ForeignKeyConcurrencyDetector) OnEndRun(app *app.App) {
	detector.iter++
}

func (detector *ForeignKeyConcurrencyDetector) OnNewRequest(entryNode *abstractgraph.AbstractServiceCall) {
	detector.requestInfoStack.Push(&RequestInfo{
		index: detector.getNewRequestInfoIndex(),
		entry: entryNode,
	})
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
	if detector.iter != IterationPhaseCheckWritesAndUpdates {
		return
	}

	datastore := dbCall.DbInstance.GetDatastore()
	schema := datastore.GetSchema()

	if schema.HasConstraintsForeignKey() {
		writtenFieldNames := detection.GetWrittenFieldNamesForOperation(dbCall)
		fields := schema.GetFieldsByNames(writtenFieldNames)
		for _, field := range fields {
			for _, constraint := range field.GetConstraints(datastores.ConstraintFilter{Reference: utils.BoolPtr(true)}) {
				refField := constraint.GetReferencedByField()
				refDatastore := refField.GetDatastore()
				if del := detector.getFirstDeleteToDatastoreIfExists(refDatastore); del != nil {
					del.addAffectedWrittenField(dbCall, field, constraint)
				}
			}
		}
	}
}

func (detector *ForeignKeyConcurrencyDetector) OnUpdate(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	//no-op
}

func (detector *ForeignKeyConcurrencyDetector) OnDelete(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	if detector.iter != IterationPhaseCheckDeletes {
		return
	}

	datastore := dbCall.DbInstance.GetDatastore()
	detector.addDelete(datastore, &delete{
		call:      dbCall,
		datastore: datastore,
	})

}

func (detector *ForeignKeyConcurrencyDetector) GetAnalysisTypeString() string {
	return "foreign_key_concurrency"
}
