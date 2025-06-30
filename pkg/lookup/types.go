package lookup

import (
	"analyzer/pkg/frameworks/blueprint"
	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
	"analyzer/pkg/types/objects"
	"analyzer/pkg/utils"
)

func CreateDynamicObjectFromType(t gotypes.Type) objects.Object {
	obj := CreateObjectFromType("*", t)
	obj.GetVariableInfo().SetDynamic()
	return obj
}

func CreateObjectFromType(name string, t gotypes.Type) objects.Object {
	logger.Logger.Infof("[LOOKUP] creating variable (%s) for type [%s]: %v", name, utils.GetType(t), t)
	info := &objects.ObjectInfo{
		Name: name,
		Type: t,
		Id:   objects.VARIABLE_UNASSIGNED_ID,
	}
	if name == "" {
		info.Id = objects.VARIABLE_INLINE_ID
	}

	switch e := t.(type) {
	case *gotypes.UserType:
		if e.UserType != nil {
			underlyingVariable := CreateObjectFromType(name, e.UserType)
			if underlyingVariable == nil {
				logger.Logger.Fatalf("[LOOKUP] unexpected nil underlying variable named (%s) for user type (%s)", name, e.Name)
			}
			e.UserType = underlyingVariable.GetType()
			underlyingVariable.GetVariableInfo().Type = e
			logger.Logger.Debugf("[LOOKUP] got underlying variable (%s) for user type named (%s): %s", underlyingVariable.String(), name, e.String())
			//logger.Logger.Warnf("[LOOKUP] returning user type variable (%s) with underlying type (%s)", underlyingVariable.String(), utils.GetType(underlyingVariable.GetType()))
			return underlyingVariable
		}
		// e.g. context.Context is nil
		logger.Logger.Warnf("[LOOKUP] user type (%s) with nil underlying type", e.String())
		return &objects.GenericObject{ObjectInfo: info}
	case *gotypes.ChanType:
		return &objects.ChanObject{ObjectInfo: info}
	case *gotypes.BasicType:
		return &objects.BasicObject{ObjectInfo: info}
	case *gotypes.InterfaceType:
		return &objects.InterfaceObject{ObjectInfo: info}
	case *gotypes.ArrayType:
		return &objects.ArrayObject{ObjectInfo: info}
	case *gotypes.AddressType:
		addressOfVariable := CreateObjectFromType("", e.AddressOf)
		addressVariable := &objects.AddressObject{
			ObjectInfo: info, AddressOf: CreateObjectFromType("", e.AddressOf),
		}
		addressOfVariable.GetVariableInfo().AddParent(addressOfVariable, addressVariable)
		return addressVariable
	case *gotypes.PointerType:
		v := &objects.PointerObject{
			ObjectInfo: info, PointerTo: CreateObjectFromType("", e.PointerTo),
		}
		return v
	case *gotypes.StructType:
		obj := objects.NewStructObject(info)
		logger.Logger.Debugf("[LOOKUP - %s] created struct object: %s", utils.GetType(obj), obj.String())
		return obj
	case *gotypes.SliceType:
		return &objects.SliceObject{ObjectInfo: info}
	case *gotypes.MapType:
		return &objects.MapObject{ObjectInfo: info, KeyValues: make(map[objects.Object]objects.Object, 0)}
	case *gotypes.FieldType:
		info.Name = e.FieldName
		wrappedFieldVariable := CreateObjectFromType(name, e.WrappedType)
		fieldVariable := &objects.FieldObject{
			ObjectInfo:      info,
			WrappedVariable: wrappedFieldVariable,
		}
		wrappedFieldVariable.GetVariableInfo().AddParent(wrappedFieldVariable, fieldVariable)
		return fieldVariable
	case *blueprint.BlueprintBackendType:
		info.Type = e.Copy(true)
		if e.IsNoSQLCollection() {
			// blueprint NoSQL type (blueprint.NoSQLCollectionType) is later assigned
			return &blueprint.BlueprintBackendObject{ObjectInfo: info}
		}
		if e.IsNoSQLCursor() {
			// blueprint NoSQL type (blueprint.NoSQLCursorType) is later assigned
			return &blueprint.BlueprintBackendObject{ObjectInfo: info}
		}
		return &blueprint.BlueprintBackendObject{ObjectInfo: info}
	case *gotypes.ServiceType:
		return &objects.ServiceObject{ObjectInfo: info}
	}
	logger.Logger.Fatalf("[LOOKUP] could not create variable with name (%s) for type (%v) (type = %s)", name, t, utils.GetType(t))
	return nil
}
