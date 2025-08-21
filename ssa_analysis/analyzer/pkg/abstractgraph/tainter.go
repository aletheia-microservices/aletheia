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

			if otherTaint.IsWrite() && currTaint.IsWrite() {
				if !currField.HasConstraintForeignKeyToField(otherField) && !otherField.HasConstraintForeignKeyToField(currField) {
					// 2nd condition is for sanity check
					// may happen when iterating queue.Push() --> queue.Pop()
					constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, currField, otherField)
					constraint.EnableMandatory()
					currField.AddConstraint(constraint)
					currDb.GetLastSchema().AddConstraint(constraint)
					fmt.Printf("\t\t[ITERATOR] [WRITE-WRITE] added new constraint: %s\n", constraint)
				}
			} else if otherTaint.IsRead() && currTaint.IsWrite() {
				if constraint := currField.GetConstraintForeignKeyToField(otherField); constraint != nil && constraint.IsMandatory() {
					constraint.DisableMandatory()
					fmt.Printf("\t\t[ITERATOR] [WRITE-READ] disabled mandatory: %s\n", constraint)
				} else if !currField.HasConstraintForeignKeyToField(otherField) && !otherField.HasConstraintForeignKeyToField(currField) {
					// 2nd condition is for sanity check
					// may happen when iterating queue.Push() --> queue.Pop()
					constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, currField, otherField)
					currField.AddConstraint(constraint)
					currDb.GetLastSchema().AddConstraint(constraint)
					fmt.Printf("\t\t[ITERATOR] [WRITE-READ] added new constraint: %s\n", constraint)
				}
			} else if otherTaint.IsWrite() && currTaint.IsRead() {
				// e.g.,
				// postnotification: TODO
				// 		=> FOREIGN_KEY notifications_queue.notification.PostID REFERENCES posts_db.post.PostID [MANDATORY]
				// 		=> the constraint already exists so condition ahead is skipped
				//
				// sockshop3: shippingservice.ship_db.write(shipping)* // shippingservice.ship_queue.push() --> queuemaster.ship_queue.pop()*
				// 		=> FOREIGN_KEY ship_db.shipments REFERENCES ship_queue.notification
				//
				// digota: orderservice.orders_db.write(order)* <-- orderservice.skuservice.get(ctx, item.parent) // skuservice.skus_db.read(parent)*
				// 		=> FOREIGN_KEY orders_db.orders.Items[*].Parent REFERENCES skus_db.skus.Id

				// NOTE: verify this
				// not sure if we shoud leave the following conditions ahead with "nothing to do"
				// to also capture foreign keys for other combinations of operatiions
				if !currField.HasConstraintForeignKeyToField(otherField) && !otherField.HasConstraintForeignKeyToField(currField) {
					// 2nd condition is for sanity check
					constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, otherField, currField)
					otherField.AddConstraint(constraint)
					otherDb.GetLastSchema().AddConstraint(constraint)
					fmt.Printf("\t\t[ITERATOR] [READ-WRITE] added new constraint: %s\n", constraint)
				}
			} else if otherTaint.IsRead() && currTaint.IsRead() {
				// nothing to do
			} else if otherTaint.IsDelete() && (currTaint.IsRead() || currTaint.IsWrite() || currTaint.IsDelete()) {
				// nothing to do
			} else if (otherTaint.IsRead() || otherTaint.IsWrite() || otherTaint.IsDelete()) && currTaint.IsDelete() {
				// nothing to do
			} else {
				log.Fatalf("\t\t[ABSTRACTGRAPH] unexpected taint mapping:\nOTHER TAINT: %s\nCURR TAINT:%s", otherTaint.LongString(), currTaint.LongString())
			}
		}
	}
}
