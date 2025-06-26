package fkey_concurrency

import (
	"fmt"

	"github.com/golang-collections/collections/stack"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/datastores"
	"analyzer/pkg/detection/detection"
	"analyzer/pkg/frameworks/blueprint"
	"analyzer/pkg/logger"
	"analyzer/pkg/utils"
)

type IterationPhase int

const (
	IterationPhaseOne_CheckDeletes IterationPhase = iota
	IterationPhaseTwo_CheckWritesAndUpdates
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
	if detector.iter == IterationPhaseOne_CheckDeletes {
		detector.iter = IterationPhaseTwo_CheckWritesAndUpdates
	}
}

func (detector *ForeignKeyConcurrencyDetector) getFirstDeleteToDatastoreIfExists(ds *datastores.Datastore) *delete {
	if dels, exists := detector.deletes[ds]; exists {
		return dels[0]
	}
	return nil
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

func (detector *ForeignKeyConcurrencyDetector) getCurrentRequest() *request {
	return detector.requestInfoStack.Peek().(*request)
}

func (detector *ForeignKeyConcurrencyDetector) OnNewRun(app *app.App) {
	app.ResetAllDataflows()
}

func (detector *ForeignKeyConcurrencyDetector) OnEndRun(app *app.App) {
	detector.iter++
}

func (detector *ForeignKeyConcurrencyDetector) OnNewRequest(entryNode *abstractgraph.AbstractServiceCall) {
	detector.requestInfoStack.Push(&request{
		index: detector.getNewRequestInfoIndex(),
		entry: entryNode,
	})
}

func (detector *ForeignKeyConcurrencyDetector) OnEndRequest(app *app.App) {
	detector.checkInconsistencies()
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
	if detector.iter != IterationPhaseTwo_CheckWritesAndUpdates {
		return
	}
	logger.Logger.Debugf("[FK CONCURRENCY DETECTOR] onWriteOrUpdate: %s", dbCall.String())

	datastore := dbCall.DbInstance.GetDatastore()
	//schema := datastore.GetSchema()

	/* if !schema.HasConstraintsForeignKey() {
		return
	} */

	req := detector.getCurrentRequest()
	operationIdx := req.numOperations()
	writtenFieldNames := detection.GetWrittenFieldNamesForOperation(dbCall)
	writtenFields := dbCall.DbInstance.GetDatastore().GetSchema().GetFieldsByNames(writtenFieldNames)

	operation := NewOperation(operationIdx, dbCall, datastore, writtenFields)
	req.addOperation(operation)

	params := dbCall.GetParams()
	if blueprintBackendMethod := dbCall.GetParsedCall().GetMethod().(*blueprint.BackendMethod); blueprintBackendMethod != nil {
		switch datastore.Type {
		case datastores.Queue, datastores.NoSQL:
			obj := dbCall.GetParam(blueprintBackendMethod.GetWrittenObjectIndex())
			operation.addWrittenObjects(obj)

			abstractgraph.TaintDataflowWrite(app, obj, dbCall, datastore, "", true, child_idx)

		case datastores.Cache:
			key := dbCall.GetParam(blueprintBackendMethod.GetWrittenKeyIndex())
			value := dbCall.GetParam(blueprintBackendMethod.GetWrittenObjectIndex())
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
			logger.Logger.Fatalf("[FK CONCURRENCY DETECTOR] unknown type of datastore (%s) to parse call: %s", utils.GetType(datastore), dbCall.String())
		}
	}
}

func (detector *ForeignKeyConcurrencyDetector) checkInconsistencies() {
	for _, op := range detector.getCurrentRequest().getOperations() {
		for _, currField := range op.getWrittenFields() {
			filter1 := datastores.ConstraintFilter{
				Reference: utils.BoolPtr(true),
				Mandatory: utils.BoolPtr(true),
			}
			
			// cannot be a mandatory reference, meaning that both fields (and corresponding objects)
			// were not previously written in **some** (not all) request
			mandatoryForeignKeyConstraints := currField.GetConstraints(filter1)
			if len(mandatoryForeignKeyConstraints) > 0 {
				continue
			}

			filter2 := datastores.ConstraintFilter{
				Reference: utils.BoolPtr(true),
				// now we get all non-mandatory references
				// so we can add a flag concerning that field
				Mandatory: utils.BoolPtr(false),
			}

			for _, constraint := range currField.GetConstraints(filter2) {
				refField := constraint.GetReferenceToField()
				refDatastore := refField.GetDatastore()
				if del := detector.getFirstDeleteToDatastoreIfExists(refDatastore); del != nil {
					/* if !detector.isMandatoryAssociation(op, del, refField, currFieldIdx) {
						del.flagAffectedWriteOnField(op.getDbCall(), currField, constraint)
					} */
					del.flagAffectedWriteOnField(op.getDbCall(), currField, constraint)
				}
			}
		}
	}
}

// DEPRECATED: canFlagInconsistency checks if foreign key association is not mandatory
func (detector *ForeignKeyConcurrencyDetector) isMandatoryAssociation(op *write, del *delete, refField *datastores.Field, idx int) bool {
	// condition for a more fine-grained detection:
	// in the current request, there can't be a write to the original record that is being referenced
	// ensuring that we only flag inconsistencies for cases of 1:N associations and not 1:1
	for _, otherOp := range detector.getCurrentRequest().getOperations() {
		if otherOp == op {
			continue
		}

		// high-level verification
		// checks based on database fields
		if otherOp.getDatastore() == del.getDatastore() {
			if ok, _ := otherOp.writesToField(refField); ok {

				// low-level verification
				// checks based on objects dataflow
				obj := op.getWrittenObjectAt(idx)
				for _, dep := range obj.GetNestedDependencies(false) {
					// find any dependencies of the current object that were also used in the other write
					for _, df := range dep.GetVariableInfo().GetAllWriteDataflowsForDatastore(otherOp.getDatastore().GetName()) {
						// found the field of the other operation that the current field is referencing to
						if df.Field == refField {
							logger.Logger.Debugf("[FK CONCURRENCY] cannot flag: %s vs. %s", df.Field.String(), refField.String())
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func (detector *ForeignKeyConcurrencyDetector) OnUpdate(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	//no-op
}

func (detector *ForeignKeyConcurrencyDetector) OnDelete(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	if detector.iter != IterationPhaseOne_CheckDeletes {
		return
	}
	datastore := dbCall.DbInstance.GetDatastore()
	detector.addDelete(datastore, &delete{
		call:      dbCall,
		datastore: datastore,
	})

}

func (detector *ForeignKeyConcurrencyDetector) GetAnalysisTypeString() string {
	return "fkey_concurrency"
}
