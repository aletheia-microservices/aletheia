package keycoordination

import (
	"slices"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/app/backends"
	"analyzer/pkg/config"
	"analyzer/pkg/utils"
)

type ForeignRead struct {
	op1         *ReadOperation // read1 has the secondary taint referencing read2
	op2         *ReadOperation
	field1      *backends.Field
	field2      *backends.Field // field used as foreign key in op2
	constraint1 *backends.Constraint
	constraint2 *backends.Constraint
}

func (fr *ForeignRead) FirstCallString() string {
	return fr.op1.call.String()
}

func (fr *ForeignRead) SecondCallString() string {
	return fr.op2.call.String()
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
	// EVAL: logrus.Tracef("[%s| CHECKER] on op: %s\n", detector.GetTypeStringUpper(), currOp.call.String())
	for _, arg := range currOp.arguments {
		// EVAL: logrus.Tracef("[%s | CHECKER] arg (%s)\n", detector.GetTypeStringUpper(), arg.String())

		for _, taintLst := range arg.GetAllTaints() {
			var currTaint *abstractgraph.AbstractTaint
			var otherTaints []*abstractgraph.AbstractTaint

			for _, taint := range taintLst {
				if taint.IsPrimary() && taint.GetDatabaseCallID() == currOp.GetCallID() {
					currTaint = taint
				} else if !taint.IsPrimary() && taint.GetDatabaseCallID() != currOp.GetCallID() && !slices.Contains(otherTaints, taint) {
					otherTaints = append(otherTaints, taint)
				}
			}

			if currTaint == nil {
				// EVAL: logrus.Tracef("\t[%s | CHECKER] skipping currTaint with otherTaints\n", detector.GetTypeStringUpper())
				continue
			}

			currFieldPath := currTaint.GetDatabasePath()
			currDatabase := app.GetDatabaseByName(utils.ExtractDatabaseNameFromFieldPath(currFieldPath))
			currField := app.ComputeDatabaseFieldFromPath(currDatabase, currFieldPath)

			for _, otherTaint := range otherTaints {
				// EVAL: logrus.Tracef("\t[%s | CHECKER] arg (%s) w/ secondary taint: %s\n", detector.GetTypeStringUpper(), arg.String(), otherTaint.LongString())
				otherOp := request.FindOperationByCallID(otherTaint.GetDatabaseCallID())

				// sanity checks
				if currOp == otherOp || otherOp == nil || currOp.call.GetToNode().GetDatabaseName() == otherOp.call.GetToNode().GetDatabaseName() {
					// EVAL: logrus.Tracef("\t[%s | CHECKER] skipping nil op for otherTaint (arg=%s): %s\n", detector.GetTypeStringUpper(), arg.String(), otherTaint.LongString())
					continue
				}

				if !detector.hasForeignRead(request, currOp, otherOp) {
					otherFieldpath := otherTaint.GetDatabasePath()
					otherDatabase := app.GetDatabaseByName(utils.ExtractDatabaseNameFromFieldPath(otherFieldpath))
					otherField := app.ComputeDatabaseFieldFromPath(otherDatabase, otherFieldpath)

					// EVAL: logrus.Tracef("\t\t[%s | CHECKER] currField={%s} // otherField={%s}\n", detector.GetTypeStringUpper(), currField.String(), otherField.String())

					foreignRead := &ForeignRead{
						op1:    currOp,
						op2:    otherOp,
						field1: currField,
						field2: otherField,
					}
					detector.addForeignRead(request, foreignRead)
				}
			}
		}
	}
}

func (detector *KeyCoordinationDetector) updateForeignReadConstraints(fread *ForeignRead) {
	if detector.isTypeForeignKey() {
		fread.constraint1 = fread.field1.GetConstraintForeignKeyToField(fread.field2)
		fread.constraint2 = fread.field2.GetConstraintForeignKeyToField(fread.field1)
	} else if detector.isTypePrimaryKey() {
		fread.constraint1 = fread.field1.GetConstraintPrimaryKey()
		fread.constraint2 = fread.field2.GetConstraintPrimaryKey()
	}
}

// isValidForeignRead can only be called after the entire schema is built at the end of the iteration
func (detector *KeyCoordinationDetector) isValidForeignRead(fread *ForeignRead) bool {
	if detector.isTypePrimaryKey() {
		if !fread.field1.IsPrimaryKey() || !fread.field2.IsPrimaryKey() {
			return false
		}
	}
	if detector.isTypeForeignKey() {
		if fread.field1.IsPrimaryKey() && fread.field2.IsPrimaryKey() {
			return false
		}
	}

	if detector.isTypeForeignKey() && config.Global.RestrictiveForeignKeyCoordinationAnalysis {
		// 1. restrict detection warnings to mandatory constraints
		// 2. filter out constraints that were created in the current request
		if fread.constraint1 != nil && fread.constraint1.IsMandatory() &&
			!fread.constraint1.HasRequestIndexOnMandatory(fread.op1.reqIdx) {
			return true
		} else if fread.constraint2 != nil && fread.constraint2.IsMandatory() &&
			!fread.constraint2.HasRequestIndexOnMandatory(fread.op2.reqIdx) {
			return true
		}
	} else if detector.isTypePrimaryKey() && config.Global.RestrictivePrimaryKeyCoordinationAnalysis {
		if constraint := fread.field1.GetConstraintForeignKeyToField(fread.field2); constraint != nil {
			if constraint.IsMandatory() && !constraint.HasRequestIndexOnMandatory(fread.op1.reqIdx) {
				return true
			}
		} else if constraint := fread.field2.GetConstraintForeignKeyToField(fread.field1); constraint != nil {
			if constraint.IsMandatory() && !constraint.HasRequestIndexOnMandatory(fread.op2.reqIdx) {
				return true
			}
		} else {
			return false
		}
	} else {
		return true
	}
	return false
}
