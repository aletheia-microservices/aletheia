package abstractgraph

import (
	"fmt"
	"log"

	"analyzer/pkg/app/backends"
	"analyzer/pkg/utils"
)

func PropagateNewTaintsToDatabases(graph *AbstractCallGraph, taintMapping *TaintMapping) {
	for currTaint, otherTaintsLst := range taintMapping.mapping {
		currDb := graph.GetApp().GetDatabaseByName(utils.ExtractDatabaseNameFromFieldPath(currTaint.GetField()))
		currField := currDb.GetSchema().GetOrCreateField(currDb, currTaint.GetField())

		for _, otherTaint := range otherTaintsLst {
			otherDb := graph.GetApp().GetDatabaseByName(utils.ExtractDatabaseNameFromFieldPath(otherTaint.GetField()))
			otherField := otherDb.GetSchema().GetOrCreateField(otherDb, otherTaint.GetField())

			if currTaint.IsWrite() && otherTaint.IsWrite() || currTaint.IsWrite() && otherTaint.IsRead() {
				if !currField.HasConstraintForeignKeyToField(otherField) {
					constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, currField, otherField)
					currField.AddConstraint(constraint)
					currDb.GetSchema().AddConstraint(constraint)
					fmt.Printf("\t\t[ITERATOR] [WRITE] added new constraint: %s\n", constraint)
				}
			} else if currTaint.IsRead() && otherTaint.IsRead() || currTaint.IsRead() && otherTaint.IsWrite() || currTaint.IsWrite() && otherTaint.IsDelete() {
				if !otherField.HasConstraintForeignKeyToField(currField) {
					constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, otherField, currField)
					otherField.AddConstraint(constraint)
					otherDb.GetSchema().AddConstraint(constraint)
					fmt.Printf("\t\t[ITERATOR] [READ] added new constraint: %s\n", constraint)
				}
			} else if currTaint.IsDelete() && otherTaint.IsDelete() {
				// nothing to do
			} else {
				log.Fatalf("\t\t[ABSTRACTGRAPH] unexpected taint mapping:\nCURR TAINT: %s\nOTHER TAINT:%s", currTaint.LongString(), otherTaint.LongString())
			}
		}
	}
}
