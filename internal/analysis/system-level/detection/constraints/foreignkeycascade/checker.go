package foreignkeycascade

import (
	"slices"

	"github.com/sirupsen/logrus"

	"github.com/aletheia-microservices/aletheia/internal/app"
	"github.com/aletheia-microservices/aletheia/internal/app/backends"
)

type CascadeDelete struct {
	op            *DeleteOperation
	cascadingOps  []*DeleteOperation
	pendingFields []*backends.Field
}

// checks each delete in the request against already-pending cascade deletes, and registers
// any new pending cascade deletes of its own
func (detector *ForeignKeyCascadeDetector) checkInconsistenciesForRequest(app *app.App, request *Request) {
	for _, delete := range request.GetAllOperations() {
		// check if there was a write before involving the association
		// if thats the case, the bug is not flagged
		var do_not_flag bool
		database := app.GetDatabaseByName(delete.database)
		schema := database.GetSchemaByNameIfExists(delete.schema)
		for _, write := range request.GetAllWriteOperations() {
			for _, writtenField := range write.fields {
				for _, deletedField := range schema.GetAllFieldsLst() {
					if writtenField.HasConstraintForeignKeyNonMandatoryToField(deletedField) {
						do_not_flag = true
						logrus.Warnf("[FOREIGN KEY CASCADE | CHECKER] skipping cascade delete due to write in same request\n")
						break
					}
				}
			}
		}
		if !do_not_flag {
			// run for every delete, even leaf ones with no pending fields of their own:
			// a leaf delete can still be what satisfies an earlier delete's pending cascade
			detector.confirmPendingCascadeDelete(app, request, delete)
			
			cascadeDelete := detector.buildPendingCascadeDelete(app, delete)
			if cascadeDelete != nil {
				detector.addPendingCascadeDelete(request, cascadeDelete)
			}
		}
	}
}

func (detector *ForeignKeyCascadeDetector) buildPendingCascadeDelete(app *app.App, currOp *DeleteOperation) *CascadeDelete {
	var pendingFields []*backends.Field
	currDB := app.GetDatabaseByName(currOp.call.GetToNode().GetDatabaseName())

	for _, db := range app.GetAllDatabases() {
		// skip if it is current DB
		if db == currDB {
			continue
		}
		for _, schema := range db.GetSchemas() {
			for _, constraint := range schema.GetAllConstraints() {
				if constraint.IsForeignKey() {
					currField := constraint.GetFieldAt(1)
					if currField.GetDatabase() == currDB && currField.GetSchema().GetName() == currOp.schema {
						// found reference to current field
						otherField := constraint.GetFieldAt(0)

						// skip if other field is from a queue
						if otherField.GetDatabase().IsQueue() {
							continue
						}

						if !slices.Contains(pendingFields, otherField) {
							pendingFields = append(pendingFields, otherField)
						}
					}
				}
			}
		}
	}

	if pendingFields != nil {
		return &CascadeDelete{
			op:            currOp,
			pendingFields: pendingFields,
		}
	}
	return nil
}

func (detector *ForeignKeyCascadeDetector) confirmPendingCascadeDelete(app *app.App, request *Request, currOp *DeleteOperation) {
	currDB := app.GetDatabaseByName(currOp.call.GetToNode().GetDatabaseName())

	for _, prevCascadeDelete := range detector.getCascadeDeletesForRequest(request) {
		// skip if it is current operation
		if prevCascadeDelete.op == currOp {
			continue
		}

		dbsWithCascade := make(map[*backends.Database]bool)
		prevOp := prevCascadeDelete.op

		for _, somePendingField := range prevCascadeDelete.pendingFields {
			if currDB == somePendingField.GetDatabase() {
				// current operation is potential cascading delete of prevCascadeDelete
				// to make sure, we need to check if the current operation has a secondary taint resulting from the prev operation

				// same logic as in foreignkeycoordination but here we verify if secondaryTaint.IsDelete()
				for _, arg := range currOp.arguments {
					for _, secondaryTaint := range arg.GetSecondaryTaintsFlatList() {
						if secondaryTaint.GetDatabaseCallID() != currOp.GetCallID() && secondaryTaint.IsDelete() {
							otherOp := request.FindOperationByCallID(secondaryTaint.GetDatabaseCallID())
							if otherOp != nil && otherOp == prevOp {
								dbsWithCascade[somePendingField.GetDatabase()] = true
								prevCascadeDelete.cascadingOps = append(prevCascadeDelete.cascadingOps, currOp)
							}
						}
					}
				}
			}
		}

		// remove any pending fields whose cascading delete was found in their database
		var pendingFieldsToKeep []*backends.Field
		for _, field := range prevCascadeDelete.pendingFields {
			if _, exists := dbsWithCascade[field.GetDatabase()]; !exists {
				pendingFieldsToKeep = append(pendingFieldsToKeep, field)
			}
		}

		prevCascadeDelete.pendingFields = pendingFieldsToKeep
	}
}
