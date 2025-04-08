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
			bsonKey := updateBsonStruct.GetFieldByKey("Key")
			if basicWrappedObj, ok := bsonKey.GetWrappedVariable().(*objects.BasicObject); ok {
				op := basicWrappedObj.GetBasicType().GetBasicValue()

				if op == "$pull" {
					bsonValue := updateBsonStruct.GetFieldByKey("Value")
					
					if bsonValueSlice, ok := bsonValue.GetWrappedVariable().(*objects.SliceObject); ok {
						elemsToUpdate := bsonValueSlice.GetElements()
						
						for j, elemToUpdate := range elemsToUpdate {
							elemToUpdateStruct := elemToUpdate.(*objects.StructObject)
							elemKey := elemToUpdateStruct.GetFieldByKeyIfExists("Key").GetWrappedVariable().(*objects.BasicObject).GetBasicType().GetBasicValue()
							elemValue := elemToUpdateStruct.GetFieldByKeyIfExists("Value").GetWrappedVariable().(*objects.BasicObject).GetBasicType().GetBasicValue()
							
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
