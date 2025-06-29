package lookup

import (
	golangtypes "go/types"

	"analyzer/pkg/logger"
	"analyzer/pkg/types"
	"analyzer/pkg/types/gotypes"
	"analyzer/pkg/utils"
)

func isBlueprintPackagePath(name string) bool {
	return name == utils.BLUEPRINT_PATH_CORE_BACKEND
}

func SaveObjectToPackage(pkg *types.Package, obj golangtypes.Object) gotypes.Type {
	if constObj, ok := obj.(*golangtypes.Const); ok {
		t := LookupAndComputeTypesForGoTypes(pkg, obj.Type())
		t.(*gotypes.BasicType).Value = constObj.Val().String()
		v := CreateObjectFromType(obj.Id(), t)
		if pkg.HasPath(obj.Pkg().Path()) {
			pkg.AddConstant(v)
		} else {
			pkg.AddImportedConstant(v, obj.Pkg().Path())
		}
		return t
	} else if _, ok := obj.(*golangtypes.Var); ok {
		t := LookupAndComputeTypesForGoTypes(pkg, obj.Type())
		v := CreateObjectFromType(obj.Id(), t)
		if pkg.HasPath(obj.Pkg().Path()) {
			pkg.AddVariable(v)
		} else {
			pkg.AddImportedConstant(v, obj.Pkg().Path())
		}
		return t
	}
	return nil
}

func ResolveTypeAndAddToPackage(pkg *types.Package, object golangtypes.Object, goType golangtypes.Type, visitedNamedTypes map[*golangtypes.Named]bool, typeNameToFuncs map[string][]*golangtypes.Func, serviceTypes map[string]*gotypes.ServiceType) gotypes.Type {
	var typeName string
	var objectPackage *golangtypes.Package
	var objectPackagePath string
	var namedGoType *golangtypes.Named
	var ok bool

	if namedGoType, ok = goType.(*golangtypes.Named); ok {
		if visitedNamedTypes != nil {
			if visitedNamedTypes[namedGoType] {
				return nil
			}
			visitedNamedTypes[namedGoType] = true
		}

		typeName = namedGoType.Obj().Name()
		objectPackage = namedGoType.Obj().Pkg()
		objectPackagePath = ""

		logger.Logger.Tracef("[APP NAMED GOTYPE] [%s] %s", utils.GetType(namedGoType), namedGoType)
	} else if aliasGoType, ok := goType.(*golangtypes.Alias); ok {
		logger.Logger.Warnf("[APP GOTYPES] recursing alias for [%T] (%s); RHS = [%T] %s", aliasGoType, aliasGoType.String(), aliasGoType.Rhs(), aliasGoType.Rhs())
		return ResolveTypeAndAddToPackage(pkg, object, aliasGoType.Rhs(), visitedNamedTypes, typeNameToFuncs, serviceTypes)
	}

	if namedGoType == nil {
		logger.Logger.Fatalf("[APP GOTYPES] skipping %v with type %s", goType, utils.GetType(goType))
	}

	// built-in (e.g. error, make, println) golang types have no package
	if objectPackage != nil {
		objectPackagePath = objectPackage.Path()
	}

	if pkg.IsAppPackage() {
		// different package implies imported object
		if !pkg.HasPath(objectPackagePath) {
			if isBlueprintPackagePath(objectPackagePath) {
				bpPackage := pkg.GetImportedPackage(objectPackagePath)
				declaredType := bpPackage.GetDeclaredType(typeName).DeepCopy()
				pkg.AddImportedTypeIfNotExists(declaredType)
				pkg.AddDatastoreType(declaredType)

				logger.Logger.Infof("[APP GOTYPES] added imported blueprint type (%s)", declaredType.String())
				return declaredType
			} else if serviceType, ok := serviceTypes[typeName]; ok && serviceType.HasPackagePath(objectPackagePath) {
				pkg.AddImportedTypeIfNotExists(serviceType)
				pkg.AddServiceType(serviceType)

				logger.Logger.Infof("[APP GOTYPES] added imported service type (%s)", serviceType.String())
				return serviceType
			} else {
				userType := &gotypes.UserType{
					Name:        typeName,
					PackagePath: objectPackagePath, // this is the real package name
				}
				pkg.AddImportedTypeIfNotExists(userType)

				userType.UserType = ComputeTypesForGoTypes(pkg, namedGoType.Underlying(), true, visitedNamedTypes, typeNameToFuncs, serviceTypes)
				addMethodsIfStructOrInterface(userType, namedGoType)
				logger.Logger.Infof("[APP GOTYPES] added imported user type (%s)", userType.String())
				return userType
			}
		} else if serviceType, ok := serviceTypes[typeName]; ok && serviceType.HasPackagePath(objectPackagePath) {
			pkg.AddDeclaredType(serviceType)
			pkg.AddServiceType(serviceType)
			logger.Logger.Infof("[APP GOTYPES] added declared service type (%s)", serviceType.String())
			return serviceType
		} else {
			// new type defined in the current package
			userType := &gotypes.UserType{
				Name:        typeName,
				PackagePath: pkg.Name,
			}
			// since we can have nested structures with definitions out of order
			// there, here we only add an entry and latter we generate the subtypes
			pkg.AddDeclaredType(userType)

			userType.UserType = ComputeTypesForGoTypes(pkg, namedGoType.Underlying(), true, visitedNamedTypes, typeNameToFuncs, serviceTypes)
			addMethodsIfStructOrInterface(userType, namedGoType)
			logger.Logger.Infof("[APP GOTYPES] added declared named type (%s)", userType.String())
			return userType
		}
	} else if pkg.HasPath(objectPackagePath) {
		// new type defined in the current package
		userType := &gotypes.UserType{
			Name:        typeName,
			PackagePath: pkg.Name,
		}
		// since we can have nested structures with definitions out of order
		// there, here we only add an entry and latter we generate the subtypes
		pkg.AddDeclaredType(userType)

		userType.UserType = ComputeTypesForGoTypes(pkg, namedGoType.Underlying(), true, visitedNamedTypes, typeNameToFuncs, serviceTypes)
		addMethodsIfStructOrInterface(userType, namedGoType)
		logger.Logger.Infof("[APP GOTYPES] added declared named type %s", userType.String())
		return userType
	}
	return nil
}

func getPkgPathForGolangFuncType(funcType *golangtypes.Func) string {
	path := ""
	if funcType.Pkg() != nil { // e.g. "error" never has an object
		path = funcType.Pkg().Path()
	}
	return path
}

func addMethodsIfStructOrInterface(userType *gotypes.UserType, namedGoType *golangtypes.Named) {

	switch e := userType.UserType.(type) {
	case *gotypes.StructType:
		// in struct types, methods are directly defined for the named type
		for i := 0; i < namedGoType.NumMethods(); i++ {
			funcType := namedGoType.Method(i)
			funcTypePkgPath := getPkgPathForGolangFuncType(funcType)
			e.AddMethod(funcType.Name(), funcTypePkgPath)
			logger.Logger.Debugf("[COMPUTE GOTYPES STRUCT] saved struct function %s.%s(...) in package (%s)", userType.GetName(), funcType.Name(), funcTypePkgPath)
		}
	case *gotypes.InterfaceType:
		// in interface types, methods are defined for the underlying type of the named type
		for i := 0; i < namedGoType.Underlying().(*golangtypes.Interface).NumMethods(); i++ {
			funcType := namedGoType.Underlying().(*golangtypes.Interface).Method(i)
			funcTypePkgPath := getPkgPathForGolangFuncType(funcType)
			e.AddMethod(funcType.Name(), funcTypePkgPath)
			logger.Logger.Debugf("[COMPUTE GOTYPES STRUCT] saved interface function %s.%s(...) in package (%s)", userType.GetName(), funcType.Name(), funcTypePkgPath)
		}
	}
}

func LookupTypesForGoTypes(p *types.Package, goType golangtypes.Type) gotypes.Type {
	return ComputeTypesForGoTypes(p, goType, false, nil, nil, nil)
}

func LookupAndComputeTypesForGoTypes(p *types.Package, goType golangtypes.Type) gotypes.Type {
	// note that variable can be either inline and thus we need to compute if type not found
	// in case of declarations with assignments, we don't reach this function
	// e.g. inline "bson.D{{"id", payment.ID}}" in coll.FindOne(ctx, bson.D{{"id", payment.ID}})
	return ComputeTypesForGoTypes(p, goType, true, nil, nil, nil)
}

// + Compute
func ComputeTypesForGoTypes(p *types.Package, goType golangtypes.Type, computeIfNotFound bool, visitedNamedTypes map[*golangtypes.Named]bool, typeNameToFuncs map[string][]*golangtypes.Func, serviceTypes map[string]*gotypes.ServiceType) gotypes.Type {
	logger.Logger.Debugf("[LOOKUP GOTYPES] [%s] %v", utils.GetType(goType), goType)
	switch e := goType.(type) {
	case *golangtypes.Named:
		name := e.Obj().Name()
		path := ""
		if e.Obj().Pkg() != nil { // error never has an object
			path = e.Obj().Pkg().Path()
		}
		if namedType, ok := p.DeclaredTypes[name]; ok {
			return namedType.DeepCopy()
		}
		if serviceType, ok := p.ServiceTypes[name]; ok {
			return serviceType
		}
		if datastoreType, ok := p.DatastoreTypes[name]; ok {
			return datastoreType
		}
		if importedType, ok := p.ImportedTypes[types.ImportedTypeKey(path, name)]; ok {
			return importedType.DeepCopy()
		}

		if true {
			return ResolveTypeAndAddToPackage(p, e.Obj(), e, visitedNamedTypes, typeNameToFuncs, serviceTypes)
		}
	case *golangtypes.Struct:
		structType := &gotypes.StructType{Methods: make(map[string]string)}
		for i := 0; i < e.NumFields(); i++ {
			f := e.Field(i)
			name := f.Name()
			fieldType := &gotypes.FieldType{
				Origin:      structType,
				WrappedType: ComputeTypesForGoTypes(p, f.Type(), computeIfNotFound, visitedNamedTypes, typeNameToFuncs, serviceTypes),
				StructField: true,
				Embedded:    f.Embedded(),
				FieldName:   name,
				FieldTag:    e.Tag(i),
			}
			structType.AddOrGetFieldType(fieldType)
		}
		return structType
	case *golangtypes.Basic:
		return &gotypes.BasicType{
			Name: e.Name(),
		}
	case *golangtypes.Interface:
		return &gotypes.InterfaceType{
			Content: goType.String(),
			Methods: make(map[string]string),
		}
	case *golangtypes.Signature:
		signatureType := &gotypes.SignatureType{}
		for i := 0; i < e.Results().Len(); i++ {
			var v *golangtypes.Var = e.Results().At(i)
			signatureType.ReturnTypes = append(signatureType.ReturnTypes, ComputeTypesForGoTypes(p, v.Type(), computeIfNotFound, visitedNamedTypes, typeNameToFuncs, serviceTypes))
		}
		return signatureType
	case *golangtypes.Tuple:
		tupleType := &gotypes.TupleType{}
		for i := 0; i < e.Len(); i++ {
			var v *golangtypes.Var = e.At(i)
			tupleType.Types = append(tupleType.Types, ComputeTypesForGoTypes(p, v.Type(), computeIfNotFound, visitedNamedTypes, typeNameToFuncs, serviceTypes))
		}
		return tupleType
	case *golangtypes.Slice:
		sliceType := &gotypes.SliceType{}
		sliceType.UnderlyingType = ComputeTypesForGoTypes(p, e.Elem(), computeIfNotFound, visitedNamedTypes, typeNameToFuncs, serviceTypes)
		return sliceType
	case *golangtypes.Array:
		arrayType := &gotypes.ArrayType{}
		arrayType.ElementsType = ComputeTypesForGoTypes(p, e.Elem(), computeIfNotFound, visitedNamedTypes, typeNameToFuncs, serviceTypes)
		return arrayType
	case *golangtypes.Pointer:
		pointerType := &gotypes.PointerType{}
		pointerType.PointerTo = ComputeTypesForGoTypes(p, e.Elem(), computeIfNotFound, visitedNamedTypes, typeNameToFuncs, serviceTypes)
		return pointerType
	case *golangtypes.Chan:
		chanType := &gotypes.ChanType{}
		chanType.ChanType = ComputeTypesForGoTypes(p, e.Elem(), computeIfNotFound, visitedNamedTypes, typeNameToFuncs, serviceTypes)
		return chanType
	case *golangtypes.Map:
		chanType := &gotypes.MapType{}
		chanType.KeyType = ComputeTypesForGoTypes(p, e.Key(), computeIfNotFound, visitedNamedTypes, typeNameToFuncs, serviceTypes)
		chanType.ValueType = ComputeTypesForGoTypes(p, e.Elem(), computeIfNotFound, visitedNamedTypes, typeNameToFuncs, serviceTypes)
		return chanType
	case *golangtypes.Alias:
		logger.Logger.Warnf("[APP GOTYPES] recursing alias for [%T] (%s); RHS = [%T] %s", e, e.String(), e.Rhs(), e.Rhs())
		return ComputeTypesForGoTypes(p, e.Rhs(), computeIfNotFound, visitedNamedTypes, typeNameToFuncs, serviceTypes)
	default:
		if goType != nil {
			logger.Logger.Fatalf("unknown gotype %s for %v", utils.GetType(goType), goType)
		}
	}
	logger.Logger.Fatalf("unknown gotype (%s) for %v", utils.GetType(goType), goType)
	return nil
}
