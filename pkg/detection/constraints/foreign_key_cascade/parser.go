package foreign_key_cascade

import (
	"fmt"

	"analyzer/pkg/datastores"
	"analyzer/pkg/frameworks/blueprint"
	"analyzer/pkg/types/objects"
)

func parseNoSQLUpdateOnRemovedFields(blueprintBackendMethod *blueprint.BackendMethod, schema *datastores.Schema, update objects.Object) ([]*datastores.Constraint, string) {
	var removedForeignKeysWithConstraints []*datastores.Constraint
	operationRepr := "{"

	if updateBsonSlice, ok := update.(*objects.SliceObject); ok {
		updateBsonElements := updateBsonSlice.GetElements()
		for i, updateBson := range updateBsonElements {
			updateBsonStruct := updateBson.(*objects.StructObject)
			bsonKeyField := updateBsonStruct.GetFieldByKeyIfExists("Key")
			if bsonKeyField == nil {
				bsonKeyField = updateBsonStruct.GetFieldAt(0)
			}
			bsonKey := bsonKeyField.GetWrappedVariable()
			if basicWrappedObj, ok := bsonKey.(*objects.BasicObject); ok {
				op := basicWrappedObj.GetBasicType().GetBasicValue()

				if op == "$pull" {
					bsonValue := updateBsonStruct.GetFieldByKey("Value")
					if bsonValue == nil {
						bsonValue = updateBsonStruct.GetFieldAt(1)
					}
					
					if bsonValueSlice, ok := bsonValue.GetWrappedVariable().(*objects.SliceObject); ok {
						elemsToUpdate := bsonValueSlice.GetElements()
						
						for j, elemToUpdate := range elemsToUpdate {
							elemToUpdateStruct := elemToUpdate.(*objects.StructObject)

							elemKeyField := elemToUpdateStruct.GetFieldByKeyIfExists("Key")
							if elemKeyField == nil {
								elemKeyField = elemToUpdateStruct.GetFieldAt(0)
							}

							elemValueField := elemToUpdateStruct.GetFieldByKeyIfExists("Value")
							if elemValueField == nil {
								elemValueField = elemToUpdateStruct.GetFieldAt(1)
							}


							elemKey := elemKeyField.GetWrappedVariable().(*objects.BasicObject).GetBasicType().GetBasicValue()
							elemValue := elemValueField.GetWrappedVariable().(*objects.BasicObject).GetBasicType().GetBasicValue()
							
							//collection := blueprintBackendMethod.GetCalledBackendType().NoSQLComponent.Collection
							//fieldName := collection + "." + elemKey

							fieldName := schema.GetRootFieldName() + "." + elemKey
							constraints := schema.GetConstraintsForeignKeyForFieldName(fieldName)
							removedForeignKeysWithConstraints = append(removedForeignKeysWithConstraints, constraints...)

							operationRepr += fmt.Sprintf("%s(%s, %s)", op, fieldName, elemValue)
							if i+j < len(updateBsonElements)-1+len(elemsToUpdate)-1 {
								operationRepr += ", "
							}
						}
					}
				}

			}
		}
	}
	operationRepr += "}"
	return removedForeignKeysWithConstraints, operationRepr
}
