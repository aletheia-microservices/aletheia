package controlflow

import (
	"go/ast"
	"go/token"
	golangtypes "go/types"
	"strconv"
	"strings"

	"analyzer/pkg/logger"
	"analyzer/pkg/lookup"
	"analyzer/pkg/service"
	"analyzer/pkg/types"
	"analyzer/pkg/types/gotypes"
	"analyzer/pkg/types/objects"
	"analyzer/pkg/utils"
)

func computeArrayIndexFromExpr(expr ast.Expr, block *types.Block) (int, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		index, err := strconv.Atoi(e.Value)
		if err != nil {
			logger.Logger.Warnf("error converting index (%s): %s", e.Value, err.Error())
		}
		return index, true
	case *ast.Ident:
		obj := block.GetObject(e.Name) //FIXME: should actually look for all possible objects in package to include global ones
		if basicObj, ok := getUnderlyingBasicObjectIfExists(obj); ok {
			valStr := basicObj.GetBasicType().Value
			index, err := strconv.Atoi(valStr)
			if err != nil {
				logger.Logger.Warnf("error converting index (%s) in expr (%v): %s", valStr, expr, err.Error())
			}
			return index, true
		}
	}
	logger.Logger.Warnf("could not compute index for expr (%s) with type (%s)", expr, utils.GetType(expr))
	return -1, false
}

func computeArrayIndexFromObject(obj objects.Object) (int, bool) {
	valStr := obj.GetType().GetBasicValue()
	index, err := strconv.Atoi(valStr)
	if err != nil {
		return index, true
	}
	logger.Logger.Warnf("error converting index (%s) from object (%s): %s", valStr, obj.String(), err.Error())
	return -1, false
}

func getUnderlyingBasicObjectIfExists(obj objects.Object) (*objects.BasicObject, bool) {
	if addressObj, ok := obj.(*objects.AddressObject); ok {
		return getUnderlyingBasicObjectIfExists(addressObj.AddressOf)
	} else if pointerObj, ok := obj.(*objects.PointerObject); ok {
		return getUnderlyingBasicObjectIfExists(pointerObj.PointerTo)
	} else if basicObj, ok := obj.(*objects.BasicObject); ok {
		return basicObj, true
	}
	return nil, false
}

func lookupVariableFromIdentIfExists(ctx *ControlflowContext, block *types.Block, ident *ast.Ident) objects.Object {
	logger.Logger.Debugf("[CFG - LOOKUP IDENT] (%s) looking up variable for ident (%s)", ctx.String(), ident.Name)

	if utils.IsBuiltInGoFunc(ident.Name) {
		logger.Logger.Warnf("[CFG - LOOKUP IDENT] FIXME: IDENTIFIED BUILT IN FUNC FOR (%s)", ident.Name)
		return nil
	}

	if utils.IsBuiltInGoType(ident.Name) {
		basicType := gotypes.NewBasicType(ident.Name, "")
		variable := &objects.BasicObject{
			ObjectInfo: &objects.ObjectInfo{
				Type: basicType,
				Id:   objects.VARIABLE_INLINE_ID,
			},
		}
		logger.Logger.Debugf("variable is builtin go value type: %s", variable.String())
		return variable
	}

	if utils.IsBuiltInConstValue(ident.Name) {
		typeName := utils.GetBuiltInConstTypeName(ident.Name)
		basicType := gotypes.NewBasicType(typeName, ident.Name)
		variable := &objects.BasicObject{
			ObjectInfo: &objects.ObjectInfo{
				Type: basicType,
				Id:   objects.VARIABLE_INLINE_ID,
			},
		}
		logger.Logger.Debugf("variable is builtin go const type: %s", variable.String())
		return variable
	}

	variable := block.GetLastestVariableIfExists(ident.Name)
	if variable != nil {
		logger.Logger.Debugf("variable found in block declarations: %s", variable.String())
		return variable
	}
	// if its not a variable in the block then it can be either
	// 1. a declared variable in the package
	// 2. a ident from a import (which is dealt with in the switch case for the selectorExpr)
	variable = ctx.GetPackage().GetDeclaredVariableOrConstIfExists(ident.Name)
	if variable != nil {
		logger.Logger.Debugf("variable found in package declarations: %s", variable.String())
		return variable
	}
	return nil
}

func lookupFieldForVariable(variable objects.Object, fieldName string, inAssignment bool) objects.Object {
	switch v := variable.(type) {
	case *objects.PointerObject:
		return lookupFieldForVariable(v.PointerTo, fieldName, inAssignment)
	case *objects.StructObject:
		var field *objects.FieldObject
		field = v.GetFieldByKeyIfExists(fieldName)
		if field == nil {
			logger.Logger.Infof("[LOOKUP FIELD] struct object: %s", variable.String())
			fieldType := v.GetStructType().GetFieldTypeByName(fieldName)
			field = lookup.CreateObjectFromType(fieldName, fieldType).(*objects.FieldObject)
			v.SetFieldByKey(fieldName, field)
			logger.Logger.Warnf("[LOOKUP FIELD - REVIEW!!] added new variable (%s) for field (%s) in structure variable (%s) -- fields: %v", field.String(), fieldName, variable.String(), v.GetFieldsMap())
		}
		if !inAssignment {
			return objects.UnwrapFieldVariable(field)
		}
		return field
	case *objects.GenericObject:
		logger.Logger.Fatalf("[LOOKUP FIELD] ignoring generic variable (%s)", v.String())
	default:
		logger.Logger.Fatalf("[LOOKUP FIELD] unknown type (%s) for variable: %v", objects.VariableTypeName(variable), variable.String())
	}
	return nil
}

func lookupImportedPackageFromIdent(service *service.Service, ident *ast.Ident) *gotypes.PackageType {
	if impt := service.GetFile().GetImportIfExists(ident.Name); impt != nil {
		t := &gotypes.PackageType{
			Alias: ident.Name,
			Name:  impt.PackageName,
			Path:  impt.PackagePath,
		}
		return t
	}
	return nil
}

func objOnExpr(expr ast.Expr, block *types.Block) bool {
	switch e := expr.(type) {
	case *ast.UnaryExpr:
		if ident, ok := e.X.(*ast.Ident); ok && block.HasVariable(ident.Name) {
			return true
		}
	case *ast.BinaryExpr:
		if ident, ok := e.X.(*ast.Ident); ok && block.HasVariable(ident.Name) {
			return true
		}
	case *ast.StarExpr:
		if ident, ok := e.X.(*ast.Ident); ok && block.HasVariable(ident.Name) {
			return true
		}
	}
	return false
}

// FIXME: THIS IS DUPLICATED
func findTypeFromSelectedImportedPackage(pkg *types.Package, typeExpr ast.Expr) gotypes.Type {
	goType := pkg.GetTypeInfo(typeExpr)
	logger.Logger.Infof("GO GOTYPE [%s]: %s", utils.GetType(goType), goType.String())

	if goNamedType, ok := goType.(*golangtypes.Named); ok {
		t := pkg.GetImportedTypeIfExists(goNamedType.String())
		if t != nil {
			logger.Logger.Debugf("GOT TYPE [%s]: %s", utils.GetType(t), t.String())
			return t
		}
	}
	return nil
}

func lookupObjectFromAstExpression(ctx *ControlflowContext, method *types.ParsedMethod, block *types.Block, expr ast.Expr, expectedType gotypes.Type, inAssignment bool) (variable objects.Object, packageType *gotypes.PackageType) {
	logger.Logger.Infof("[CFG - LOOKUP OBJ] (%T) visiting expression (%v)", expr, expr)
	if ctx.GetPackage() != nil && ctx.GetPackage().GetTypeInfo(expr) != nil || ctx.GetService() != nil && ctx.GetPackage().GetTypeInfo(expr) != nil {
		logger.Logger.Debugf("[CFG - LOOKUP OBJ GOTYPE] (%T) visiting go expression (%v)", ctx.GetPackage().GetTypeInfo(expr), ctx.GetPackage().GetTypeInfo(expr))
	}
	switch e := expr.(type) {
	case *ast.CallExpr:
		variable = parseAndSaveCall(ctx, method, block, e)
		variable = objects.UnwrapTupleIfSingleElement(variable)
		/* if _, ok := variable.(*objects.BasicVariable); ok {
			logger.Logger.Fatalf("GOT CALL EXPR IN METHOD (%s): %v", method.String(), e.Fun)
		} */
	case *ast.BasicLit:
		basicType := &gotypes.BasicType{
			Name:  strings.ToLower(e.Kind.String()),
			Value: e.Value,
		}
		variable = &objects.BasicObject{
			ObjectInfo: &objects.ObjectInfo{
				Type: basicType,
				Id:   objects.VARIABLE_INLINE_ID,
			},
		}
	case *ast.Ident:
		variable = lookupVariableFromIdentIfExists(ctx, block, e)
		if variable != nil {
			return variable, nil
		}

		packageType = lookupImportedPackageFromIdent(ctx.GetService(), e)
		if packageType != nil {
			return nil, packageType
		}
		logger.Logger.Fatalf("[CFG - LOOKUP VAR] unexpected nil variable (or package) with type (%s) for expr (%v)", utils.GetType(expr), expr)
		return nil, nil
	case *ast.SelectorExpr:
		variable, packageType = lookupObjectFromAstExpression(ctx, method, block, e.X, nil, inAssignment)
		if packageType != nil {
			t := findTypeFromSelectedImportedPackage(ctx.GetPackage(), expr)
			if t != nil {
				variable = lookup.CreateObjectFromType("", t)
				return variable, nil
			}

			importedPkg := ctx.GetPackage().GetImportedPackage(packageType.Path)
			if importedPkg.IsExternalPackage() {
				// note that variable can be either inline and thus we need to compute if type not found
				// in case of declarations with assignments, we don't reach this function
				// e.g. inline "bson.D{{"id", payment.ID}}" in coll.FindOne(ctx, bson.D{{"id", payment.ID}})
				t := lookup.LookupAndComputeTypesForGoTypes(ctx.GetPackage(), ctx.GetPackage().GetTypeInfo(e.Sel))
				variable = lookup.CreateObjectFromType("", t)
				return variable, nil
			}
			declaredType := importedPkg.GetDeclaredTypeIfExists(e.Sel)
			if declaredType != nil {
				newDeclaredType := declaredType.DeepCopy()
				variable = lookup.CreateObjectFromType("", newDeclaredType)
				return variable, nil
			}

			variable = importedPkg.GetDeclaredVariableOrConst(e.Sel)
			return variable, nil
		}
		if tupleVariable, ok := variable.(*objects.TupleObject); ok {
			if tupleVariable.NumVariables() == 1 {
				variable = tupleVariable.GetVariableAt(0)
			} else {
				logger.Logger.Fatalf("[CFG - LOOKUP VAR] unexpected number (%d) of variables in tuple: %s", tupleVariable.NumVariables(), tupleVariable.String())
			}
		}

		field := lookupFieldForVariable(variable, e.Sel.Name, inAssignment)
		return field, nil

	case *ast.KeyValueExpr:
		if key, ok := e.Key.(*ast.Ident); ok {
			variable, _ = lookupObjectFromAstExpression(ctx, method, block, e.Value, nil, false)
			fieldVariable := &objects.FieldObject{
				ObjectInfo: &objects.ObjectInfo{
					Name: key.Name,
					Type: &gotypes.FieldType{
						WrappedType: variable.GetType(),
						FieldName:   key.Name,
					},
					Id: objects.VARIABLE_INLINE_ID,
				},
				WrappedVariable: variable,
			}
			variable.GetVariableInfo().AddParent(variable, fieldVariable)
			logger.Logger.Debugf("KEY VALUE EXPR - RETURNING VARIABLE WITH TYPE (%s): %s", utils.GetType(variable), variable.String())
			return fieldVariable, nil
		}
		logger.Logger.Fatalf("[CFG - LOOKUP VAR] unsupported key type (%s) for expr: %v", utils.GetType(e.Key), e.Key)

	// examples:
	// (1) 	- expression:
	// 			query := bson.D{{Key: "postid", Value: postID}}
	// 		- composite lit 1 (e.Type is well defined):
	// 			 bson.D{{Key: "postid", Value: postID}}
	// 		- composite lit 2 (e.Type is nil):
	// 			{Key: "postid", Value: postID}
	// (2) 	- expression:
	// 			query := bson.D{{"index", 0}}
	// 		- composite lit 1 (e.Type is well defined):
	// 			 bson.D{{"index", 0}}
	// 		- composite lit 2 (e.Type is nil and e.Elts are unkeyed):
	// 			{"index", 0}
	// (3)  - expression:
	// 			mentions := []string{"alice", "bob"}
	// 		- composite lit:
	// 			[]string{"alice", "bob"}
	case *ast.CompositeLit:
		if e.Type == nil {
			logger.Logger.Debugf("[COMPOSITE LIT - TYPE NIL] e.Type (%v), and e.Elts (%v)", e.Type, e.Elts)
			structVariable := objects.NewStructObject(objects.NewObjectInfoInline(gotypes.NewStructType()))
			for _, elt := range e.Elts {
				eltVar, _ := lookupObjectFromAstExpression(ctx, method, block, elt, nil, false)
				logger.Logger.Debugf("[%s.%s] FOUND ELT VAR (%s) FOR COMPOSITE LIT (%v)", ctx.GetService().GetName(), method.GetName(), eltVar.String(), e)
				objects.WrapObjectToField(eltVar, structVariable, true)
			}
			logger.Logger.Infof("%s", structVariable.LongString())
			return structVariable, nil
		}

		logger.Logger.Debugf("[COMPOSITE LIT] e.Type (%v), and e.Elts (%v)", e.Type, e.Elts)
		eType := lookup.ComputeTypeForAstExpr(ctx.GetService().GetFile(), e.Type)
		logger.Logger.Debugf("BEFORE eType = %s", eType.String())
		eType, eTypeOrUserType := gotypes.UnwrapUserType(eType)
		logger.Logger.Debugf("AFTER eType = %s", eType.String())

		switch eType.(type) {
		case *gotypes.StructType:
			structVariable := objects.NewStructObject(objects.NewObjectInfoInline(eTypeOrUserType))
			logger.Logger.Debugf("????? struct Type = %s", structVariable.GetType().String())
			for _, elt := range e.Elts {
				eltVar, _ := lookupObjectFromAstExpression(ctx, method, block, elt, nil, false)
				objects.WrapObjectToField(eltVar, structVariable, false)
			}
			logger.Logger.Infof("GOT STRUCT VARIABLE: %s", structVariable.String())
			if userType, ok := eTypeOrUserType.(*gotypes.UserType); ok && userType.Name == "Post" && ctx.GetService().GetName() == "UploadService" && method.Name == "UploadPost" {
				logger.Logger.Warnf("!!!!! struct Type = %s", structVariable.GetType().String())
				logger.Logger.Warnf("[%s.%s()] eType [%s]: %s", ctx.GetService().GetName(), method.GetName(), utils.GetType(eType), eType.String())
			}
			return structVariable, nil
		case *gotypes.ArrayType:
			arrayVariable := &objects.ArrayObject{
				ObjectInfo: &objects.ObjectInfo{
					Type: eTypeOrUserType,
					Id:   objects.VARIABLE_INLINE_ID,
				},
			}
			for _, elt := range e.Elts {
				eltVar, _ := lookupObjectFromAstExpression(ctx, method, block, elt, nil, false)
				arrayVariable.AddElement(eltVar)
				eltVar.GetVariableInfo().AddParent(eltVar, arrayVariable)
			}
			return arrayVariable, nil
		case *gotypes.SliceType:
			sliceVariable := &objects.SliceObject{
				ObjectInfo: &objects.ObjectInfo{
					Type: eTypeOrUserType,
					Id:   objects.VARIABLE_INLINE_ID,
				},
			}
			for _, elt := range e.Elts {
				eltVar, _ := lookupObjectFromAstExpression(ctx, method, block, elt, nil, false)
				sliceVariable.AddElement(eltVar)
				eltVar.GetVariableInfo().AddParent(eltVar, sliceVariable)
			}
			return sliceVariable, nil
		default:
			logger.Logger.Warnf("[CFG - LOOKUP VAR] [%s.%s] unsupported type for eType (%s): %s", ctx.GetService().GetName(), method.GetName(), utils.GetType(eType), eType.String())
		}
		if selectorExpr, ok := e.Type.(*ast.SelectorExpr); ok {
			logger.Logger.Debugf("[CFG - LOOKUP VAR] [%s.%s] lookup up selector (%v.%v)", ctx.GetService().GetName(), method.GetName(), selectorExpr.X, selectorExpr.Sel)
			logger.Logger.Fatalf("[CFG - LOOKUP VAR] [%s.%s] lookup up selector (%v.%v)", ctx.GetService().GetName(), method.GetName(), selectorExpr.X, selectorExpr.Sel)
			variable, packageType = lookupObjectFromAstExpression(ctx, method, block, selectorExpr, nil, inAssignment)
			if variable == nil {
				logger.Logger.Fatalf("[CFG - LOOKUP VAR] unexpected nil variable for expr (%s): %v", utils.GetType(e.Type), e.Type)
			}
			for _, elt := range e.Elts {
				eltVar, _ := lookupObjectFromAstExpression(ctx, method, block, elt, nil, false)
				if sliceVariable, ok := variable.(*objects.SliceObject); ok {
					sliceVariable.AddElement(eltVar)
				} else if arrayVariable, ok := variable.(*objects.ArrayObject); ok {
					arrayVariable.AddElement(eltVar)
				} else {
					logger.Logger.Fatalf("[CFG - LOOKUP VAR] unknown type for composite lit (%s)", utils.GetType(variable))
				}
			}
			logger.Logger.Debugf("[CFG - LOOKUP VAR] LOOKUP HERE!!! got selector (%s): %s", utils.GetType(variable), variable.String())
		} else {
			logger.Logger.Fatalf("[CFG - LOOKUP VAR] nil variable for composite lit (e.Type = %s): %v", utils.GetType(e.Type), e)
		}
	case *ast.TypeAssertExpr:
		variable, packageType = lookupObjectFromAstExpression(ctx, method, block, e.X, nil, inAssignment)

		// in a switch case we have nil types for e.g. f.(type)
		var assertedType gotypes.Type
		if e.Type == nil {
			logger.Logger.Debugf("[CFG - LOOKUP VAR] found nil type in TypeAssertExpr: %v", e)
		} else {
			assertedType = lookup.ComputeTypeForAstExpr(ctx.GetService().GetFile(), e.Type)
		}
		if interfaceVariable, ok := variable.(*objects.InterfaceObject); ok {
			// FIXME: it is creating two duplicates:
			// (i) InterfaceVariable during lookup with e.g. BasicType after UpgradeToAssertedType
			// (ii) The e.g. BasicVariable after UpgradeFromPreviousInterface with e.g. the BasicType
			//interfaceVariable.UpgradeToAssertedType(assertedType)

			var newVariable objects.Object
			if assertedType != nil {
				newVariable = lookup.CreateObjectFromType(variable.GetVariableInfo().GetName(), assertedType)
				newVariable.UpgradeFromPreviousInterface(interfaceVariable) //idk why this is here
			} else {
				newType := interfaceVariable.GetInterfaceType() //FIXME: should we actually get the user type?
				newVariable = lookup.CreateObjectFromType(variable.GetVariableInfo().GetName(), newType)
			}
			block.AddObject(newVariable)
			logger.Logger.Warnf("[FIXME - ASSERT!!!!] CREATED NEW VARIABLE (%s): %s", objects.VariableTypeName(newVariable), newVariable.String())
			return newVariable, nil
		} else {
			logger.Logger.Fatalf("[CFG - LOOKUP VAR] unexpected type (%s) for variable (%s)", utils.GetType(variable), variable.String())
		}
	case *ast.IndexExpr:
		variable, packageType = lookupObjectFromAstExpression(ctx, method, block, e.X, nil, inAssignment)
		if variable == nil {
			logger.Logger.Fatalf("[CFG - LOOKUP VAR] nil variable for index expr: %v (e.X type = %s)", e, utils.GetType(e.X))
		}
		if arrayVar, ok := variable.(*objects.ArrayObject); ok {
			arrayIndex, ok := computeArrayIndexFromExpr(e.Index, block)
			if ok {
				variable = arrayVar.GetElementAtIfExists(arrayIndex)
				if variable == nil {
					ok = false
				}
			}

			if !ok {
				elemsType := arrayVar.GetElementsType()
				variable = lookup.CreateDynamicObjectFromType(elemsType)
				arrayVar.AddDynamicElement(variable)
				logger.Logger.Warnf("CREATED VARIABLE FROM NIL ARRAY ELEMENT %s | %s", variable.String(), arrayVar.String())
			}
		}
		if mapVar, ok := variable.(*objects.MapObject); ok {
			key, _ := lookupObjectFromAstExpression(ctx, method, block, e.Index, nil, false)
			variable, ok = mapVar.KeyValues[key]

			if !ok {
				valueType := mapVar.GetMapType().GetValueType()
				variable = lookup.CreateDynamicObjectFromType(valueType)
				mapVar.AddDynamicKeyValue(key, variable)
				logger.Logger.Warnf("[CFG - LOOKUP VAR] got map %v with unassigned value for key %s \n\t\t\t\t\t\t\t => created new variable for key named (%s): %s", mapVar.String(), key.String(), key.GetType().GetBasicValue(), variable.String())
			}
		}

	case *ast.StarExpr: // e.g. *post
		if !objOnExpr(e, block) {
			t := lookup.LookupTypesForGoTypes(ctx.GetPackage(), ctx.GetPackage().TypesInfo.Types[e].Type)
			return lookup.CreateObjectFromType("", t), nil
		}

		variable, packageType = lookupObjectFromAstExpression(ctx, method, block, e.X, nil, inAssignment)
		if ptrVariable, ok := variable.(*objects.PointerObject); ok {
			return ptrVariable.PointerTo, packageType
		}

		logger.Logger.Fatalf("unknown token (%v) for star expr %v: [%s] %v", e.Star, e, utils.GetType(ctx.GetPackage().TypesInfo.Types[e].Type), ctx.GetPackage().TypesInfo.Types[e])

	case *ast.UnaryExpr:
		logger.Logger.Debugf("[CFG - LOOKUP VAR] UNARY EXPR for [%s]: %v", utils.GetType(ctx.GetPackage().TypesInfo.Types[e].Type), ctx.GetPackage().TypesInfo.Types[e])

		if e.Op == token.AND { // e.g. &post
			// creates a pointer to the variable
			variable, packageType = lookupObjectFromAstExpression(ctx, method, block, e.X, nil, inAssignment)
			ptrType := &gotypes.PointerType{
				PointerTo: variable.GetType(),
			}
			ptrVariable := &objects.PointerObject{
				PointerTo: variable,
				ObjectInfo: &objects.ObjectInfo{
					Type: ptrType,
					Id:   objects.VARIABLE_INLINE_ID,
					Name: variable.GetVariableInfo().Name,
				},
			}
			variable.GetVariableInfo().AddParent(variable, ptrVariable)
			return ptrVariable, packageType
		} else if e.Op == token.MUL { // e.g. *post
			// dereferences the pointer and gets the value of the variable
			variable, packageType = lookupObjectFromAstExpression(ctx, method, block, e.X, nil, inAssignment)
			if ptrVariable, ok := variable.(*objects.PointerObject); ok {
				return ptrVariable.PointerTo, packageType
			}
			logger.Logger.Fatalf("[LOOKUP] [%s] unexpected variable (%s) for unary token '*': %s", ctx.GetService().GetName(), utils.GetType(variable), variable.LongString())
			/* addrType := &gotypes.AddressType{
				AddressOf: variable.GetType(),
			}
			pointerVariable := &objects.PointerVariable{
				PointerTo: variable,
				ObjectInfo: &objects.ObjectInfo{
					Type: addrType,
					Id:   objects.VARIABLE_INLINE_ID,
					Name: variable.GetVariableInfo().Name,
				},
			}
			variable.GetVariableInfo().AddParent(variable, pointerVariable)
			return pointerVariable, packageType */
		} else if e.Op == token.NOT { // e.g. !exists
			variable, packageType = lookupObjectFromAstExpression(ctx, method, block, e.X, nil, inAssignment)
			if basicVariable, ok := variable.(*objects.BasicObject); ok {
				basicVariable.GetBasicType().InvertBool()
			} else {
				logger.Logger.Fatalf("[LOOKUP] unexpected type (%s) for object: %v", utils.GetType(variable), variable.String())
			}
			return variable, packageType
		} else if e.Op == token.ARROW { // e.g. err_post_channel <- err
			variable, packageType = lookupObjectFromAstExpression(ctx, method, block, e.X, nil, inAssignment)
			logger.Logger.Warnf("[LOOKUP] ignoring unary expr with token (%v): %v", e.Op, e)
			return variable, packageType
		}

		if !objOnExpr(e, block) {
			t := lookup.LookupTypesForGoTypes(ctx.GetPackage(), ctx.GetPackage().TypesInfo.Types[e].Type)
			return lookup.CreateObjectFromType("", t), nil
		}

		logger.Logger.Fatalf("unknown token (%v) for unary expr %v: [%s] %v", e.Op, e, utils.GetType(ctx.GetPackage().TypesInfo.Types[e].Type), ctx.GetPackage().TypesInfo.Types[e])

	case *ast.BinaryExpr:
		if !objOnExpr(e, block) {
			t := lookup.LookupTypesForGoTypes(ctx.GetPackage(), ctx.GetPackage().TypesInfo.Types[e].Type)
			return lookup.CreateObjectFromType("", t), nil
		}

		if e.Op == token.ADD || e.Op == token.SUB || e.Op == token.QUO || e.Op == token.MUL || e.Op == token.AND { // "+", "-", "/", "*", "&"
			objX, _ := lookupObjectFromAstExpression(ctx, method, block, e.X, nil, inAssignment)
			objY, _ := lookupObjectFromAstExpression(ctx, method, block, e.Y, nil, inAssignment)
			typeX := objX.GetType()
			typeY := objX.GetType()
			if typeX.IsSameType(typeY) {
				basicTypeX, ok1 := objX.GetType().(*gotypes.BasicType)
				basicTypeY, ok2 := objY.GetType().(*gotypes.BasicType)
				if ok1 && ok2 {
					variable = &objects.BasicObject{
						//FIXME: this operations should differ for the types of x and y (e.g. integers, strings, etc)
						ObjectInfo:        objects.NewObjectInfoInline(gotypes.NewBasicType("", basicTypeX.GetBasicValue()+basicTypeY.GetBasicValue())),
						UnderlyingObjects: []objects.Object{objX, objY},
					}
				} else {
					logger.Logger.Fatalf("[LOOKUP] one of the types of objects X (%s) or Y (%s) in BinaryExpr are not BasicType: %v", utils.GetType(typeX), utils.GetType(typeY), e)
				}
			} else {
				logger.Logger.Fatalf("x expr (%v) and y expr (%v) differ types (%v vs %v)", e.X, e.Y, typeX, typeY)
			}
		} else if e.Op == token.LEQ || e.Op == token.GEQ || e.Op == token.NEQ || e.Op == token.EQL { // "<=", ">=", "!=", "=="

			// FIXME!!!
			variable_x, t_x := lookupObjectFromAstExpression(ctx, method, block, e.X, nil, inAssignment)
			address_variable_x := &objects.AddressObject{
				AddressOf: variable_x,
				ObjectInfo: &objects.ObjectInfo{
					Type: &gotypes.AddressType{
						AddressOf: variable_x.GetType(),
					},
					Id:   objects.VARIABLE_INLINE_ID,
					Name: variable_x.GetVariableInfo().Name,
				},
			}
			variable_x.GetVariableInfo().AddParent(variable_x, address_variable_x)

			variable_y, t_y := lookupObjectFromAstExpression(ctx, method, block, e.Y, nil, inAssignment)
			address_variable_y := &objects.AddressObject{
				AddressOf: variable_y,
				ObjectInfo: &objects.ObjectInfo{
					Type: &gotypes.AddressType{
						AddressOf: variable_y.GetType(),
					},
					Id:   objects.VARIABLE_INLINE_ID,
					Name: variable_y.GetVariableInfo().Name,
				},
			}
			variable_y.GetVariableInfo().AddParent(variable_y, address_variable_y)

			if !t_x.IsSameType(t_y) {
				logger.Logger.Fatalf("x expr (%v) and y expr (%v) differ types (%v vs %v)", e.X, e.Y, t_x, t_y)
			}

			packageType = t_x
			variable = address_variable_x

			logger.Logger.Warnf("FIXMEEEEEEE!!!!!!!!")
		} else {
			logger.Logger.Fatalf("unknown token %v for binary expr %v", e.Op, e)
			variable, packageType = lookupObjectFromAstExpression(ctx, method, block, e.X, nil, inAssignment)
		}
	case *ast.SliceExpr: //e.g. timestamp_hex[:10]
		//TODO
		variable, _ := lookupObjectFromAstExpression(ctx, method, block, e.X, nil, false)
		variable.GetVariableInfo().Id = objects.VARIABLE_INLINE_ID
		return variable, nil

	case *ast.ChanType:
		variable, _ := lookupObjectFromAstExpression(ctx, method, block, e.Value, nil, false)
		variable.GetVariableInfo().Id = objects.VARIABLE_INLINE_ID
		return variable, nil

	case *ast.ArrayType: //e.g. []rune
		elemtsType := lookup.ComputeTypeForAstExpr(ctx.GetService().GetFile(), e.Elt)

		// when length is not specified then e.len is nil
		var numElemsInt int
		if e.Len == nil {
			logger.Logger.Warnf("e.LEN IS NIL!")
			numElemsInt = -1
		} else {
			numElemsType := lookup.ComputeTypeForAstExpr(ctx.GetService().GetFile(), e.Len)
			numElemsInt, _ = strconv.Atoi(numElemsType.GetBasicValue())
		}

		variable = &objects.ArrayObject{
			ObjectInfo: &objects.ObjectInfo{
				Type: &gotypes.ArrayType{
					ElementsType: elemtsType,
					NumElems:     numElemsInt,
				},
				Id: objects.VARIABLE_UNASSIGNED_ID,
			},
		}
	case *ast.MapType: //e.g. make(map[string]PriceConfig)
		keyType := lookup.ComputeTypeForAstExpr(ctx.GetService().GetFile(), e.Key)
		valueType := lookup.ComputeTypeForAstExpr(ctx.GetService().GetFile(), e.Value)
		variable = &objects.MapObject{
			KeyValues: make(map[objects.Object]objects.Object, 0),
			ObjectInfo: &objects.ObjectInfo{
				Type: &gotypes.MapType{
					KeyType:   keyType,
					ValueType: valueType,
				},
				Id: objects.VARIABLE_UNASSIGNED_ID,
			},
		}
	case *ast.ParenExpr: // e.g. "(weight - cp.InitialWeight)" in: price += (weight - cp.InitialWeight) * cp.WithinPrice
		variable, _ := lookupObjectFromAstExpression(ctx, method, block, e.X, nil, false)
		variable.GetVariableInfo().Id = objects.VARIABLE_INLINE_ID
		return variable, nil
	case *ast.FuncLit:
		return &objects.FuncObject{
			ObjectInfo: &objects.ObjectInfo{
				Id: objects.VARIABLE_INLINE_ID,
				Type: &gotypes.FuncTypeType{
					Body:   e.Body,
					Params: e.Type.Params,
				},
			},
		}, nil

	default:
		logger.Logger.Fatalf("[CFG - LOOKUP VAR] cannot lookup unknown type (%s) for variable (%v)", utils.GetType(e), e)
	}
	if variable == nil && packageType == nil {
		logger.Logger.Fatalf("[CFG - LOOKUP VAR] unexpected nil variable (or package) with type (%s) for expr (%v)", utils.GetType(expr), expr)
	}
	return variable, nil
}
