package objects

import (
	"encoding/json"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
	"analyzer/pkg/utils"
)

type MapObject struct {
	Object
	ObjectInfo       *ObjectInfo
	KeyValues        map[Object]Object
	DynamicKeyValues map[Object]Object
}

func NewMapObject(info *ObjectInfo) *MapObject {
	return &MapObject{
		ObjectInfo: info,
		KeyValues: make(map[Object]Object),
		DynamicKeyValues: make(map[Object]Object),
	}
}

/* func (v *MapVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ObjectInfo *ObjectInfo         `json:"map"`
		KeyValues    map[Variable]Variable `json:"map_kvs,omitempty"`
	}{
		ObjectInfo: v.ObjectInfo,
		KeyValues:    v.KeyValues,
	})
} */

func (v *MapObject) MarshalJSON() ([]byte, error) {
	type kvstruct struct {
		Key   Object
		Value Object
	}

	var kvs []kvstruct
	for key, value := range v.KeyValues {
		kvs = append(kvs, kvstruct{
			Key:   key,
			Value: value,
		})
	}

	return json.MarshalIndent(&struct {
		ObjectInfo *ObjectInfo `json:"map"`
		KeyValues  []kvstruct  `json:"key_values"`
	}{
		ObjectInfo: v.ObjectInfo,
		KeyValues:  kvs,
	}, "", " ")
}

func (v *MapObject) String() string {
	return v.ObjectInfo.String()
}

func (v *MapObject) LongString() string {
	return v.ObjectInfo.LongString()
}

func (v *MapObject) GetKeyValueIfExists(targetKey Object) Object {
	if !v.GetMapType().GetKeyType().IsSameType(targetKey.GetType()) {
		logger.Logger.Fatalf("[VARS MAP] provided key (%s) with type (%s) does not match expected type (%s)", targetKey.String(), utils.GetType(targetKey.GetType()), utils.GetType(v.GetMapType().GetKeyType()))
	}
	switch t := targetKey.GetType().(type) {
	case *gotypes.BasicType:
		for key, value := range v.KeyValues {
			if key.GetType().(*gotypes.BasicType).Value == t.Value {
				return value
			}
		}
	default:
		logger.Logger.Fatalf("unknown gotype (%s) for target key (%s)", utils.GetType(targetKey.GetType()), targetKey.String())
	}
	logger.Logger.Warnf("??????? could not get or create key value pairs for target key (%s) in map (%v)", targetKey.String(), v.String())
	return v.KeyValues[targetKey]
}

func (v *MapObject) AddKeyValue(key Object, value Object) {
	v.KeyValues[key] = value
	value.GetVariableInfo().AddParent(value, v)
}

func (v *MapObject) GetId() int64 {
	return v.ObjectInfo.GetId()
}

// this only makes sense when key is an object with BasicType
func (v *MapObject) AddDynamicKeyValue(key Object, value Object) {
	if v.DynamicKeyValues == nil {
		v.DynamicKeyValues = make(map[Object]Object)
	}
	v.DynamicKeyValues[key] = value
}

func (v *MapObject) GetType() gotypes.Type {
	if v.ObjectInfo.GetType() == nil {
		logger.Logger.Fatalf("[VARS ADDRESS] unexpected nil type for map variable: %s", v.String())
	}
	return v.ObjectInfo.GetType()
}

func (v *MapObject) GetMapType() *gotypes.MapType {
	if userType, ok := v.ObjectInfo.GetType().(*gotypes.UserType); ok {
		return userType.UserType.(*gotypes.MapType)
	}
	return v.ObjectInfo.GetType().(*gotypes.MapType)
}

func (v *MapObject) GetVariableInfo() *ObjectInfo {
	return v.ObjectInfo
}

func (v *MapObject) AddKeyValuePair(key Object, value Object) {
	v.KeyValues[key] = value
	value.GetVariableInfo().AddParent(value, v)
}

func (v *MapObject) GetDependencies() []Object {
	var dependencies []Object
	for _, value := range v.KeyValues {
		dependencies = append(dependencies, value)
	}
	dependencies = append(dependencies, v.GetVariableInfo().GetDependencies()...)
	return dependencies
}

func (v *MapObject) GetNestedDependencies(includeRefBy bool) []Object {
	var deps = []Object{v}
	if v.GetVariableInfo().HasReferences() {
		deps = append(deps, v.GetVariableInfo().GetReferencesNestedDependencies(includeRefBy, v)...)
	}
	if includeRefBy && v.GetVariableInfo().IsReferencedBy() {
		deps = append(deps, v.GetVariableInfo().GetNestedRefByDependencies(nil)...)
	}
	for _, elem := range v.GetDependencies() { // to include underlying dependencies from variable info
		logger.Logger.Debugf("[MAP OBJECT] GOT NESTED DEP (INC/ VINFO) FOR ELEM %s (%s)", elem.String(), VariableTypeName(elem))
		deps = append(deps, elem.GetNestedDependencies(includeRefBy)...)
	}
	return deps
}

func (v *MapObject) NewVersion() Object {
	copy := &MapObject{
		KeyValues:        make(map[Object]Object, 0),
		DynamicKeyValues: make(map[Object]Object),
		ObjectInfo:       v.ObjectInfo.Copy(true),
	}
	for k, v := range v.KeyValues { // FIXME: this is not 100% correct
		copy.KeyValues[k] = v
		copy.KeyValues[k].GetVariableInfo().AddParent(copy.KeyValues[k], copy)
	}
	return copy
}

func (v *MapObject) NewObject() Object {
	newObject := NewMapObject(NewObjectInfo(v.GetType()))
	// new map points to the same key-value pairs as the original one
	for k, v := range v.KeyValues {
		newObject.AddKeyValue(k, v)
		v.GetVariableInfo().AddParent(v, newObject)
	}
	for k, v := range v.DynamicKeyValues {
		newObject.AddKeyValue(k, v)
		v.GetVariableInfo().AddParent(v, newObject)
	}
	return newObject
}

func (v *MapObject) Copy(force bool) Object {
	copy := &MapObject{
		KeyValues:        make(map[Object]Object, 0),
		DynamicKeyValues: make(map[Object]Object),
		ObjectInfo:       v.ObjectInfo.Copy(force),
	}
	for k, v := range v.KeyValues {
		copy.KeyValues[k] = v.Copy(force)
		copy.KeyValues[k].GetVariableInfo().AddParent(copy.KeyValues[k], copy)
	}
	return copy
}

func (v *MapObject) DeepCopy() Object {
	logger.Logger.Debugf("[VARS MAP - DEEP COPY] (%s) %s", VariableTypeName(v), v.String())
	copy := &MapObject{
		KeyValues:        make(map[Object]Object, 0),
		DynamicKeyValues: make(map[Object]Object),
		ObjectInfo:       v.ObjectInfo.DeepCopy(),
	}
	for k, v := range v.KeyValues {
		copy.KeyValues[k] = v.DeepCopy()
		copy.KeyValues[k].GetVariableInfo().AddParent(copy.KeyValues[k], copy)
	}
	return copy
}

func (v *MapObject) AddReferenceWithID(target Object, creator string) {
	v.ObjectInfo.AddReferenceWithID(v, target, creator)
	if targetMap, ok := target.(*MapObject); ok {
		for key, value := range v.KeyValues {
			targetValue, ok := targetMap.KeyValues[key]
			if !ok {
				logger.Logger.Warnf("[VARS MAP - REF] invalid target key %s in map: %v", key, targetMap)
			} else {
				targetValue.AddReferenceWithID(value, creator)
			}
		}
		logger.Logger.Debugf("[VARS MAP - REF] added reference (%s) -> (%s) with id = %d (creator: %s)", v.ObjectInfo.Name, target.GetVariableInfo().GetName(), v.ObjectInfo.Id, creator)
		return
	}
	// exception: map[string]interface{} --> struct
	if _, keyIsMapType := v.GetMapType().KeyType.(*gotypes.BasicType); keyIsMapType {
		if _, valueIsInterfaceType := v.GetMapType().ValueType.(*gotypes.InterfaceType); valueIsInterfaceType {
			if targetStructVariable, ok := target.(*StructObject); ok {
				for key, value := range v.KeyValues {
					targetValue, ok := targetStructVariable.fields[key.GetType().GetBasicValue()]
					if !ok {
						logger.Logger.Warnf("[VARS MAP - REF] skipping target key (%s) (%s) in struct (%s) with fields:", utils.GetType(key), key.GetType().GetBasicValue(), targetStructVariable)
						for k, f := range targetStructVariable.fields {
							logger.Logger.Warnf("\t\t [VARS MAP - REF] - %s: (%s) %s", k, VariableTypeName(f), f.String())
						}
					} else {
						logger.Logger.Debugf("[VARS MAP - REF] target (%s) got referenced to value (%s)", targetValue.String(), value.String())
						targetValue.AddReferenceWithID(value, creator)
					}
				}
				logger.Logger.Debugf("[VARS MAP - REF] added reference (%s) (%s) -> (%s) (%s) with id = %d (creator: %s)", v.ObjectInfo.Name, utils.GetType(v), target.GetVariableInfo().GetName(), utils.GetType(target), v.ObjectInfo.Id, creator)
				return
			}
		}
	}
	logger.Logger.Fatalf("[VARS MAP - REF] attempted to reference variables with different types (%s vs %s) (%s vs %s)", v.String(), target.String(), utils.GetType(v), utils.GetType(target))
}

func (v *MapObject) GetNestedFieldVariables(prefix string, noSQL bool) ([]Object, []string) {
	var nestedVariables []Object
	var nestedIDs []string

	logger.Logger.Debugf("[VARS MAP] found (%d) field VARIABLES for (%s): %v", len(v.KeyValues), v.String(), v.KeyValues)
	for _, value := range v.KeyValues {
		if fieldVariable, ok := value.(*FieldObject); ok {
			nestedFieldObjects, nestedFieldIDs := fieldVariable.GetNestedFieldVariables(prefix, noSQL)
			nestedVariables = append(nestedVariables, nestedFieldObjects...)
			nestedIDs = append(nestedIDs, nestedFieldIDs...)
		} else {
			logger.Logger.Warnf("[VARS MAP] ignoring field key typed (%s) for value variable (%s) in map variable (%s): (%s)", utils.GetType(value), value.String(), VariableTypeName(v), v.String())
		}
	}
	return nestedVariables, nestedIDs
}

func (v *MapObject) GetNestedFieldVariablesWithReferences(prefix string, noSQL bool) ([]Object, []string) {
	logger.Logger.Debugf("[VARS MAP] HAS REFERENCE????? %v", v.GetVariableInfo().GetReferences())
	nestedVariables, nestedIDs := v.GetNestedFieldVariables(prefix, noSQL)
	for _, reference := range v.GetVariableInfo().GetReferences() {
		logger.Logger.Debugf("[VARS MAP] HEREEEEE FOR REFERENCE %s", reference.String())
		if referenceMapVar, ok := reference.Object.(*MapObject); ok {
			nestedVariablesRef, nestedIDsRef := referenceMapVar.GetNestedFieldVariablesWithReferences(prefix, noSQL)
			nestedVariables = append(nestedVariables, nestedVariablesRef...)
			nestedIDs = append(nestedIDs, nestedIDsRef...)
		} else if referenceStructVar, ok := reference.Object.(*StructObject); ok {
			nestedVariablesRef, nestedIDsRef := referenceStructVar.GetNestedFieldVariablesWithReferences(prefix, noSQL)
			nestedVariables = append(nestedVariables, nestedVariablesRef...)
			nestedIDs = append(nestedIDs, nestedIDsRef...)
		} else {
			logger.Logger.Warnf("[VARS MAP] ignoring reference typed (%s) for reference variable (%s) in map variable (%s): (%s)", utils.GetType(reference.Object), reference.Object.String(), VariableTypeName(v), v.String())
		}
	}
	return nestedVariables, nestedIDs
}

func (v *MapObject) GetUnassaignedVariables() []Object {
	var variables []Object
	if v.GetVariableInfo().IsUnassigned() {
		variables = append(variables, v)
		for k, v := range v.KeyValues {
			variables = append(variables, k.GetUnassaignedVariables()...)
			variables = append(variables, v.GetUnassaignedVariables()...)
		}
	}
	return variables
}
