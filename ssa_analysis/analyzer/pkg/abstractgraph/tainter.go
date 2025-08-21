package abstractgraph

import (
	"fmt"
	"log"

	"analyzer/pkg/app/backends"
	"analyzer/pkg/utils"
)

func PropagateNewTaintsToDatabases(graph *AbstractCallGraph, taintMapping *TaintMapping) {
	for currTaint, otherTaintsLst := range taintMapping.mapping {
		currDb := graph.GetApp().GetDatabaseByName(utils.ExtractDatabaseNameFromFieldPath(currTaint.GetDatabasePath()))
		currField := currDb.GetLastSchema().GetOrCreateField(currDb, currTaint.GetDatabasePath())

		for _, otherTaint := range otherTaintsLst {
			otherDb := graph.GetApp().GetDatabaseByName(utils.ExtractDatabaseNameFromFieldPath(otherTaint.GetDatabasePath()))

			if currDb == otherDb {
				// skip if its the same
				// may happen when iterating queue.Push() --> queue.Pop()
				continue
			}
			otherField := otherDb.GetLastSchema().GetOrCreateField(otherDb, otherTaint.GetDatabasePath())

			if currTaint.IsWrite() && otherTaint.IsWrite() {
				if !currField.HasConstraintForeignKeyToField(otherField) && !otherField.HasConstraintForeignKeyToField(currField) {
					// 2nd condition is for sanity check
					// may happen when iterating queue.Push() --> queue.Pop()
					constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, currField, otherField)
					constraint.EnableMandatory()
					currField.AddConstraint(constraint)
					currDb.GetLastSchema().AddConstraint(constraint)
					fmt.Printf("\t\t[ITERATOR] [WRITE-WRITE] added new constraint: %s\n", constraint)
				}
			} else if currTaint.IsWrite() && otherTaint.IsRead() {
				if constraint := currField.GetConstraintForeignKeyToField(otherField); constraint != nil && constraint.IsMandatory() {
					constraint.DisableMandatory()
				} else if !currField.HasConstraintForeignKeyToField(otherField) && !otherField.HasConstraintForeignKeyToField(currField) {
					// 2nd condition is for sanity check
					// may happen when iterating queue.Push() --> queue.Pop()
					constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, currField, otherField)
					currField.AddConstraint(constraint)
					currDb.GetLastSchema().AddConstraint(constraint)
					fmt.Printf("\t\t[ITERATOR] [WRITE-READ] added new constraint: %s\n", constraint)
				}
			} else if currTaint.IsRead() && otherTaint.IsWrite() {
				// NOTE: verify this
				// not sure if we shoud leave the following conditions ahead with "nothing to do"
				// to also capture foreign keys for other combinations of operatiions
				if constraint := otherField.GetConstraintForeignKeyToField(currField); constraint != nil && constraint.IsMandatory() {
					constraint.DisableMandatory()
				} else if !otherField.HasConstraintForeignKeyToField(currField) && !currField.HasConstraintForeignKeyToField(otherField) {
					// 2nd condition is for sanity check
					// may happen when iterating queue.Push() --> queue.Pop()
					constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, otherField, currField)
					otherField.AddConstraint(constraint)
					otherDb.GetLastSchema().AddConstraint(constraint)
					fmt.Printf("\t\t[ITERATOR] [READ-WRITE] added new constraint: %s\n", constraint)
				}
			} else if currTaint.IsRead() && otherTaint.IsRead() {
				// nothing to do
			} else if (currTaint.IsRead() || currTaint.IsWrite() || currTaint.IsDelete()) && otherTaint.IsDelete() {
				// nothing to do
			} else if currTaint.IsDelete() && (currTaint.IsRead() || currTaint.IsWrite() || currTaint.IsDelete()) {
				// nothing to do
			} else {
				log.Fatalf("\t\t[ABSTRACTGRAPH] unexpected taint mapping:\nCURR TAINT: %s\nOTHER TAINT:%s", currTaint.LongString(), otherTaint.LongString())
			}
		}
	}
}
