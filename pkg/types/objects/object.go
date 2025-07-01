package objects

import (
	"fmt"
	"slices"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
	"analyzer/pkg/utils"
)

const VARIABLE_INLINE_ID int64 = 0
const VARIABLE_UNASSIGNED_ID int64 = -1

type Object interface {
	String() string
	LongString() string
	GetVariableInfo() *ObjectInfo
	GetDependencies() []Object
	AddReferenceWithID(reference Object, creator string)
	GetUnassaignedVariables() []Object
	Copy(force bool) Object
	DeepCopy() Object
	NewObject() Object
	NewVersion() Object
	GetId() int64
	GetType() gotypes.Type
	AssignVariable(variable Object)
	UpgradeFromPreviousInterface(interfaceVariable *InterfaceObject)
	GetNestedDependencies(includeRefBy bool) []Object
	GetElementAt(index int) Object
}

func GetAllParents(obj Object) []Object {
	var res []Object
	for _, parent := range obj.GetVariableInfo().GetParents() {
		// sanity check
		if obj == parent || obj.GetVariableInfo() == parent.GetVariableInfo() {
			logger.Logger.Fatalf("CYCLIC DEPENDENCY!! ABORTING!!\n\n PARENT TYPE = %s\n\n PARENT OBJ = %s\n\n CURRENT TYPE = %s\n\n CURRENT OBJ = %s", VariableTypeName(parent), parent.String(), VariableTypeName(obj), obj.String())
		}

		res = append(res, parent)
		res = append(res, GetAllParents(parent)...)
	}
	return res
}

func VariableTypeName(v Object) string {
	return utils.GetType(v) + " " + utils.GetType(v.GetType())
}

func GetDependenciesStringLst(deps ...Object) string {
	out := ""
	for _, d := range deps {
		out += fmt.Sprintf("\t\t\t\t\t - [%d] (%s) %s: %s\n", d.GetId(), VariableTypeName(d), d.GetVariableInfo().GetName(), d.GetVariableInfo().String())
	}
	return out
}

func UnwrapTupleIfSingleElement(variable Object) Object {
	if tupleVar, ok := variable.(*TupleObject); ok && len(tupleVar.Objects) == 1 {
		return tupleVar.Objects[0]
	}
	return variable
}

func UnwrapUntilBasicObjectIfExists(obj Object) (*BasicObject, bool) {
	if addressObj, ok := obj.(*AddressObject); ok {
		return UnwrapUntilBasicObjectIfExists(addressObj.AddressOf)
	} else if pointerObj, ok := obj.(*PointerObject); ok {
		return UnwrapUntilBasicObjectIfExists(pointerObj.PointerTo)
	} else if fieldObj, ok := obj.(*FieldObject); ok {
		return UnwrapUntilBasicObjectIfExists(fieldObj.GetWrappedVariable())
	} else if basicObj, ok := obj.(*BasicObject); ok {
		return basicObj, true
	}
	return nil, false
}

func WrapObjectToField(obj Object, structObj *StructObject, addTypeToStruct bool) {
	fieldObj, ok := obj.(*FieldObject)
	if !ok {
		fieldObj = NewFieldObject(NewObjectInfoInline(gotypes.NewFieldType(obj.GetType())), obj)
	}
	if addTypeToStruct {
		structObj.AddOrGetFieldVariableAndType(fieldObj)
	} else {
		structObj.AddOrGetFieldVariable(fieldObj)
	}
}

func UnwrapAddressVariable(variable Object) Object {
	if addressVar, ok := variable.(*AddressObject); ok {
		return UnwrapAddressVariable(addressVar.GetAddressOf())
	}
	return variable
}

func UnwrapPointerVariable(variable Object) Object {
	if pointerVar, ok := variable.(*PointerObject); ok {
		return UnwrapPointerVariable(pointerVar.GetPointerTo())
	}
	return variable
}

func UnwrapFieldVariable(variable Object) Object {
	if fieldVar, ok := variable.(*FieldObject); ok {
		return UnwrapPointerVariable(fieldVar.WrappedVariable)
	}
	return variable
}

func AddVariableInfoDependencies(variable Object, deps []Object) {
	variable.GetVariableInfo().Dependencies = append(variable.GetVariableInfo().Dependencies, deps...)
}

func AddUnderlyingDepsFromFuncCall(obj Object, deps []Object) bool {
	objType, _ := gotypes.UnwrapUserAndFieldTypes(obj.GetType())

	if _, ok := objType.(*gotypes.GenericType); ok {
		genericVariable, _ := obj.(*GenericObject)
		genericVariable.Params = append(genericVariable.Params, deps...)
		return true
	} else if _, ok := objType.(*gotypes.InterfaceType); ok {
		interfaceVariable, _ := obj.(*InterfaceObject)
		if interfaceVariable.UnderlyingVariable != nil || len(deps) > 1 {
			return false
		}
		interfaceVariable.UnderlyingVariable = deps[0]
		return true
		/* } else if sliceType, ok := objType.(*gotypes.SliceType); ok {
		sliceVariable, _ := obj.(*SliceObject)
		for _, dep := range deps {
			depType, _ := gotypes.UnwrapUserAndFieldTypes(dep.GetType())
			if !sliceType.UnderlyingType.IsSameType(depType) {
				logger.Logger.Fatalf("[CFG CALLS] (%s) --- (%s) %s", utils.GetType(sliceType.UnderlyingType), VariableTypeName(dep), dep.String())
				return false
			}
		}
		sliceVariable.Elements = append(sliceVariable.Elements, deps...)
		return true
		*/
		/* } else if _, ok := objType.(*gotypes.BasicType); ok {
		basicVariable, _ := obj.(*BasicObject)
		basicVariable.UnderlyingObjects = append(basicVariable.UnderlyingObjects, deps...)
		return true
		*/
	} else {
		obj.GetVariableInfo().AddDependencies(deps)
	}
	return false
}

func GetReversedNestedFieldsAndNames(variable Object, fieldName string, noSQL bool, queue bool) (nestedObjects []Object, nestedNames []string) {
	logger.Logger.Debugf("[VARS] get reversed nested fields and names for variable (%s): %s", utils.GetType(variable), variable.String())

	// e.g. in queue.pop we need to explicitly specify the fieldname otherwise it will be wrong because we are reading from an inteface
	// also applies to nosql
	logger.Logger.Debugf("TOP FIELD NAME BEFORE = %s", fieldName)
	if fieldName == "" {
		fieldName = variable.GetType().GetName()
	}

	logger.Logger.Debugf("TOP FIELD NAME AFTER = %s", fieldName)

	if structVariable, ok := variable.(*StructObject); ok {
		variables := []Object{structVariable}
		names := []string{fieldName}
		nestedObjects, nestedNames = structVariable.GetNestedFieldVariablesWithReferences(fieldName, noSQL)

		variables = append(variables, nestedObjects...)
		names = append(names, nestedNames...)
		slices.Reverse(variables)
		slices.Reverse(names)
		return variables, names
	} else if mapVariable, ok := variable.(*MapObject); ok {
		variables := []Object{mapVariable}
		names := []string{fieldName}
		nestedObjects, nestedNames = mapVariable.GetNestedFieldVariablesWithReferences(fieldName, noSQL)

		variables = append(variables, nestedObjects...)
		names = append(names, nestedNames...)
		slices.Reverse(variables)
		slices.Reverse(names)
		return variables, names
	} else if arrayObj, ok := variable.(*ArrayObject); ok { //FIXME: fix the names that are provided for the elements
		objs := []Object{arrayObj}
		names := []string{fieldName}

		for idx, elem := range arrayObj.Elements {
			logger.Logger.Debugf("[REVERSE NESTED FIELDS W/ NAMES] appending array elem at index %d: %s", idx, elem)
			objs = append(objs, elem)
			names = append(names,  fieldName + ".*")
		}
		for _, dynamicElem := range arrayObj.DynamicElements {
			logger.Logger.Debugf("[REVERSE NESTED FIELDS W/ NAMES] appending dynamic array elem at index: %s", dynamicElem)
			objs = append(objs, dynamicElem)
			names = append(names, fieldName + ".*")
		}
		slices.Reverse(objs)
		slices.Reverse(names)
		logger.Logger.Debugf("[REVERSE NESTED FIELDS W/ NAMES] returning objects and names:\n \t\t\t- %s\n \t\t\t- %s", objs, names)
		return objs, names
	} else if addressVariable, ok := variable.(*AddressObject); ok {
		return GetReversedNestedFieldsAndNames(addressVariable.AddressOf, fieldName, noSQL, queue)
	} else if pointerVariable, ok := variable.(*PointerObject); ok {
		return GetReversedNestedFieldsAndNames(pointerVariable.PointerTo, fieldName, noSQL, queue)
	} else {
		logger.Logger.Fatalf("unexpected type (%s) for variable (%s)", utils.GetType(variable), variable.String())
	}
	return nil, nil
}
