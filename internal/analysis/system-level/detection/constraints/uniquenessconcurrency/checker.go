package uniquenessconcurrency

import (
	"slices"

	"github.com/aletheia-microservices/aletheia/internal/app"
	"github.com/aletheia-microservices/aletheia/internal/app/backends"
	"github.com/aletheia-microservices/aletheia/internal/utils"
)

type VulnerableWriteSet struct {
	constrainedOp     *WriteOperation
	otherOps          []*WriteOperation
	constrainedFields []*backends.Field
}

func (writeSet *VulnerableWriteSet) addOtherOperation(op *WriteOperation) {
	writeSet.otherOps = append(writeSet.otherOps, op)
}

func (writeSet *VulnerableWriteSet) hasOtherOperation(op *WriteOperation) bool {
	return slices.Contains(writeSet.otherOps, op)
}

// checkInconsistenciesForRequest checks all the writes of a request once it ends, so that a related
// write is found whether it comes before or after the unique write
func (detector *UniquenessConcurrencyDetector) checkInconsistenciesForRequest(app *app.App, request *Request) {
	// 1. create a write set for each write of a unique field
	for _, op := range request.GetAllOperations() {
		if constrainedFields := computeConstrainedFields(app, op); constrainedFields != nil {
			writeSet := &VulnerableWriteSet{
				constrainedOp:     op,
				constrainedFields: constrainedFields,
			}
			detector.addVulnerableWriteSet(request, writeSet)
		}
	}

	// 2. add each write that carries a secondary taint from a write in step 1 to its write set
	// same logic as in foreignkeycoordination and foreignkeycascade
	// but here we verify if secondaryTaint.IsWrite()
	for _, currOp := range request.GetAllOperations() {
		for _, arg := range currOp.arguments {
			for _, secondaryTaint := range arg.GetSecondaryTaintsFlatList() {
				if secondaryTaint.GetDatabaseCallID() != currOp.GetCallID() && secondaryTaint.IsWrite() {
					otherOp := request.FindOperationByCallID(secondaryTaint.GetDatabaseCallID())
					if otherOp != nil {
						otherWriteSet := detector.findVulnerableWriteSetForOperation(request, otherOp)
						if otherWriteSet != nil && !otherWriteSet.hasOtherOperation(currOp) {
							otherWriteSet.addOtherOperation(currOp)
						}
					}
				}
			}
		}
	}
}

// computeConstrainedFields returns the unique fields written by op, or nil if it writes none
func computeConstrainedFields(app *app.App, op *WriteOperation) []*backends.Field {
	dbname := op.call.GetToNode().GetDatabaseName()
	db := app.GetDatabaseByName(dbname)

	var constrainedFields []*backends.Field
	for _, arg := range op.arguments {
		for _, taint := range arg.GetPrimaryTaintsFlatList() {
			fieldpath := taint.GetDatabasePath()

			// [TO BE IMPROVED]
			// there may be cases where primary taint is not related to this database
			// when services make more than one call to different databases
			//
			// in the future, we may just associate the taint with the call ID
			// and then just check if the IDs match
			if dbname == utils.ExtractDatabaseNameFromFieldPath(fieldpath) {
				field := app.ComputeDatabaseFieldFromPath(db, fieldpath)
				if field.HasContraintUniqueness() && !slices.Contains(constrainedFields, field) {
					constrainedFields = append(constrainedFields, field)
				}
			}
		}
	}
	return constrainedFields
}
