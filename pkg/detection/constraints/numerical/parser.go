package numerical

import (
	"fmt"

	"analyzer/pkg/datastores"
	"analyzer/pkg/frameworks/blueprint"
	"analyzer/pkg/types/objects"
)

func parseNoSQLUpdate(blueprintBackendMethod *blueprint.BackendMethod, schema *datastores.Schema, update objects.Object) ([]*datastores.Constraint, string) {
	var affectedConstraints []*datastores.Constraint
	operationRepr := "{"

	if updateBsonSlice, ok := update.(*objects.SliceObject); ok {
		updateBsonElements := updateBsonSlice.GetElements()
		for i, updateBson := range updateBsonElements {
			updateBsonStruct := updateBson.(*objects.StructObject)
			bsonKey := updateBsonStruct.GetFieldByKey("Key")
			if basicWrappedObj, ok := bsonKey.GetWrappedVariable().(*objects.BasicObject); ok {
				op := basicWrappedObj.GetBasicType().GetBasicValue()

				if op == "$inc" || op == "$dec" {
					bsonValue := updateBsonStruct.GetFieldByKey("Value")

					if bsonValueSlice, ok := bsonValue.GetWrappedVariable().(*objects.SliceObject); ok {
						elemsToUpdate := bsonValueSlice.GetElements()

						for j, elemToUpdate := range elemsToUpdate {
							elemToUpdateStruct := elemToUpdate.(*objects.StructObject)
							elemKey := elemToUpdateStruct.GetFieldByKeyIfExists("Key").GetWrappedVariable().(*objects.BasicObject).GetBasicType().GetBasicValue()
							elemValue := elemToUpdateStruct.GetFieldByKeyIfExists("Value").GetWrappedVariable().(*objects.BasicObject).GetBasicType().GetBasicValue()
							collection := blueprintBackendMethod.GetCalledBackendType().NoSQLComponent.Collection
							fieldName := collection + "." + elemKey

							for _, constraint := range schema.GetConstraintsNumericalForFieldName(fieldName) {
								if op == "$inc" && constraint.GetNumerical().ComparesIfLowerOrEqual() ||
									op == "$dec" && constraint.GetNumerical().ComparesIfGreaterOrEqual() {
									affectedConstraints = append(affectedConstraints, constraint)
								}
							}

							operationRepr += fmt.Sprintf("%s(%s, %s)", op, elemKey, elemValue)
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
	return affectedConstraints, operationRepr
}
