package keycoordination

import (
	"fmt"

	"analyzer/pkg/app"
	"analyzer/pkg/app/backends"
	"analyzer/pkg/utils"
)

type ForeignRead struct {
	op1    *ReadOperation // read1 has the secondary taint referencing read2
	op2    *ReadOperation
	field1 *backends.Field
	field2 *backends.Field // field used as foreign key in op2
}

func (detector *KeyCoordinationDetector) checkInconsistency(app *app.App, request *Request) {
	allOps := request.GetAllOperations()
	// check in reverse
	i := len(allOps) - 1
	for i >= 0 {
		detector.checkInconsistencyForOp(app, request, allOps[i])
		i--
	}

}

func (detector *KeyCoordinationDetector) checkInconsistencyForOp(app *app.App, request *Request, currOp *ReadOperation) {
	// same logic as in foreignkeycascade but here we verify if secondaryTaint.IsDelete()
	fmt.Printf("[%s| CHECKER] on op: %v\n", detector.GetTypeStringUpper(), currOp.call.String())
	for _, arg := range currOp.arguments {
		fmt.Printf("[%s | CHECKER] arg (%s)\n", detector.GetTypeStringUpper(), arg.String())
		fmt.Printf("\t[%s | CHECKER] arg (%s) w/ all taints: %v\n", arg.String(), detector.GetTypeStringUpper(), arg.GetAllTaints())
		for _, secondaryTaint := range arg.GetSecondaryTaintsFlatList() {
			fmt.Printf("\t[%s | CHECKER] arg (%s) w/ secondary taint: %v\n", detector.GetTypeStringUpper(), arg.String(), secondaryTaint.LongString())

			if secondaryTaint.GetDatabaseCallID() != currOp.GetCallID() && secondaryTaint.IsRead() {
				otherOp := request.FindOperationByCallID(secondaryTaint.GetDatabaseCallID())

				if currOp == otherOp {
					continue
				}

				if otherOp != nil && !detector.hasForeignRead(request, currOp, otherOp) {
					otherFieldpath := secondaryTaint.GetDatabasePath()
					otherDatabase := app.GetDatabaseByName(utils.ExtractDatabaseNameFromFieldPath(otherFieldpath))
					otherField := app.ComputeDatabaseFieldFromPath(otherDatabase, otherFieldpath)

					// [TO BE IMPROVED]
					// this is a bit hardcoded for now but for the future
					// we should associate the schema to the call before the analysis
					var field *backends.Field
					if constraints := otherField.GetConstraintForeignKey(); constraints != nil {
						field = constraints[0].GetFieldAt(1)
					} else if fields := app.GetAllFieldsReferencingCurrent(otherField); fields != nil {
						field = fields[0]
					}

					if (detector.isTypePrimaryKey() && field.IsPrimaryKey() && otherField.IsPrimaryKey()) ||
						(detector.isTypeForeignKey() && (!field.IsPrimaryKey() || !otherField.IsPrimaryKey())) {
						foreignRead := &ForeignRead{op1: currOp, op2: otherOp, field1: field, field2: otherField}
						detector.addForeignRead(request, foreignRead)
					}

				}
			}
		}
	}
}
