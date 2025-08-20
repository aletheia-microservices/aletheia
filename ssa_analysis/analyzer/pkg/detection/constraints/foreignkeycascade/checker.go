package foreignkeycascade

import (
	"slices"

	"analyzer/pkg/app"
	"analyzer/pkg/app/backends"
)

type CascadeDelete struct {
	op            *DeleteOperation
	cascadingOps  []*DeleteOperation
	pendingFields []*backends.Field
}

func (detector *ForeignKeyCascadeDetector) registerFutureCascadeDelete(app *app.App, currOp *DeleteOperation) *CascadeDelete {
	cascadeDelete := &CascadeDelete{op: currOp}
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
					if currField.GetDatabase() == currDB {
						// found reference to current field
						otherField := constraint.GetFieldAt(0)

						// skip if other field is from a queue
						if otherField.GetDatabase().IsQueue() {
							continue
						}

						if !slices.Contains(cascadeDelete.pendingFields, otherField) {
							cascadeDelete.pendingFields = append(cascadeDelete.pendingFields, otherField)
						}
					}
				}
			}
		}
	}

	return cascadeDelete
}

func (detector *ForeignKeyCascadeDetector) markCascadingDelete(app *app.App, request *Request, currOp *DeleteOperation) {
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
