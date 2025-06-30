package service

import (
	"go/ast"
	"slices"

	"analyzer/pkg/datastores"
	"analyzer/pkg/frameworks/blueprint"
	"analyzer/pkg/logger"
	"analyzer/pkg/types"
	"analyzer/pkg/types/gotypes"
	"analyzer/pkg/types/objects"
	"analyzer/pkg/utils"
)

func (service *Service) isMethodExposedByService(methodName string) bool {
	// check if the function declaration being implemented by the Service structure
	// is actually an exposed method to other services
	_, ok := service.ExposedMethods[methodName]
	return ok
}

func (service *Service) methodReceiverImplementsService(parsedMethod *types.ParsedMethod) bool {
	methodReceiver := parsedMethod.GetReceiverIfExists()
	if methodReceiver == nil {
		return false
	}
	return service.ImplType.IsCurrentServiceType(methodReceiver.GetType())
}

func (service *Service) methodIsServiceConstructor(parsedMethod *types.ParsedMethod) bool {
	return parsedMethod.Name == service.ConstructorName
}

func (service *Service) attachConstructor(parsedMethod *types.ParsedMethod) {
	service.Constructor = parsedMethod
}

func (service *Service) AttachParsedMethods() []*types.ParsedMethod {
	var serviceImplementedMethods []*types.ParsedMethod
	logger.Logger.Tracef("[SERVICE] inspecting exposed methods for service implementation %s\n", service.ImplName)

	for _, parsedMethod := range service.GetPackage().GetAllParsedMethods() {
		if service.methodReceiverImplementsService(parsedMethod) {
			if service.isMethodExposedByService(parsedMethod.Name) {
				service.attachExportedMethod(parsedMethod)
			} else {
				service.attachInternalMethod(parsedMethod)
			}
			if service.HasQueueField() {
				if instance := service.getQueueInstanceIfMethodHandlesQueue(parsedMethod); instance != nil {
					logger.Logger.Tracef("FOUND: METHOD (%s) HANDLES QUEUE? INSTANCE = %v", parsedMethod.Name, instance.GetName())
					service.attachQueueHandlerMethod(parsedMethod, instance)
				}
			}
			serviceImplementedMethods = append(serviceImplementedMethods, parsedMethod)
		} else if service.methodIsServiceConstructor(parsedMethod) {
			service.attachConstructor(parsedMethod)
		} else {
			logger.Logger.Warnf("[SERVICE] [%s] ignoring attachment parsed method (%s)", service.GetName(), parsedMethod.String())
		}
	}
	return serviceImplementedMethods
}

func (service *Service) AttachAllPackageMethods() {
	for _, parsedMethod := range service.GetPackage().GetAllParsedMethods() {
		if !parsedMethod.HasAttachedService() {
			service.PackageMethods[parsedMethod.Name] = parsedMethod
			logger.Logger.Warnf("[PARSER] [%s] attached package method: %s", service.Name, parsedMethod.LongString())
		}
	}
}

// RegisterFields parses the service's fields from the previously parsed impl user type declared in the package
func (service *Service) RegisterFields() {
	logger.Logger.Tracef("[PARSER] [%s] registering fields...\n", service.Name)

	if structVar, ok := service.GetImplVariable().(*objects.StructObject); ok {
		for idx, fieldVariable := range structVar.GetFieldVariablesLst() {
			fieldName := fieldVariable.GetFieldType().GetFieldName()
			field := service.computeServiceFieldFromVariable(fieldVariable, fieldName, idx)
			service.AddOrGetField(fieldName, field)
		}
	}

	var fieldsStr string
	i := 0
	for _, field := range service.Fields {
		fieldsStr += "\t\t\t\t\t - " + field.LongString()
		if i < len(service.Fields)-1 {
			fieldsStr += "\n"
		}
		i++
	}
	logger.Logger.Infof("[PARSER] [%s] registered %d fields:\n%s", service.GetName(), len(service.Fields), fieldsStr)
}

func (service *Service) computeServiceFieldFromVariable(fieldVariable *objects.FieldObject, paramName string, idx int) types.Field {
	wrappedType := fieldVariable.GetFieldType().GetWrappedType()
	wrappedObject := fieldVariable.GetWrappedVariable()
	fieldInfo := types.FieldInfo{
		Name:   paramName,
		Idx:    idx,
		Object: wrappedObject,
		Type:   wrappedType,
	}
	switch t := wrappedType.(type) {
	case *blueprint.BlueprintBackendType:
		dbField := &types.DatabaseField{
			FieldInfo:         fieldInfo,
			Queue:             t.IsQueue(),
			DatastoreInstance: t.DatastoreInstance,
		}
		if t.IsQueue() {
			service.EnableQueueField()
		}
		return dbField
	case *gotypes.ServiceType:
		return &types.ServiceField{
			FieldInfo: fieldInfo,
		}
	}
	return &types.GenericField{
		FieldInfo: fieldInfo,
	}
}

func (service *Service) RegisterConstructor() {
	method := service.GetPackage().GetParsedMethod(service.ConstructorName, "")
	method.EnableConstructorFlag()
	service.SetConstructor(method)
	logger.Logger.Infof("[PARSER] [%s] registered constructor (%s)", service.GetName(), service.ConstructorName)
}

func (service *Service) RegisterImplName() {
	var exportedMethods []string
	for name := range service.ExposedMethods {
		exportedMethods = append(exportedMethods, name)
	}
	// get impl name
	ast.Inspect(service.GetFile().GetAst(), func(n ast.Node) bool {
		if service.HasImplName() {
			return false
		}
		if funcDecl, ok := n.(*ast.FuncDecl); ok && slices.Contains(exportedMethods, funcDecl.Name.Name) {
			if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
				if recv, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
					if ident, ok := recv.X.(*ast.Ident); ok {
						service.SetImplName(ident.Name)
						service.SetImplTypeName(ident.Name)
						return false
					}
				}
			}
		}
		return true
	})
}

func (service *Service) getQueueInstanceIfMethodHandlesQueue(parsedMethod *types.ParsedMethod) datastores.DatabaseInstance {
	var instance datastores.DatabaseInstance

	// inspect methods
	ast.Inspect(parsedMethod.Ast, func(n ast.Node) bool {
		ok, methodIdent, fieldIdent := getFuncCallIfSelectedServiceField(n, parsedMethod.Ast.Recv.List[0].Names[0])
		if !ok || fieldIdent == nil {
			return true
		}
		if datastoreField := service.GetDatastoreFieldIfExists(fieldIdent.Name); datastoreField != nil {
			if datastoreField.IsQueue() && methodIdent.Name == "Pop" {
				instance = datastoreField.GetDatastoreInstance()
				logger.Logger.Tracef("[%s.%s] GOT FIELD (%s): %s // INSTANCE = %v", service.GetName(), parsedMethod.GetName(), utils.GetType(datastoreField), datastoreField.LongString(), instance)
				return false
			}
		}
		return true
	})
	return instance
}

func getFuncCallIfSelectedServiceField(service ast.Node, expectedRecvIdent *ast.Ident) (ok bool, methodIdent *ast.Ident, fieldIdent *ast.Ident) {
	hasCurrentServiceRecvIdent := func(selectorExpr *ast.SelectorExpr, expectedRecvIdent *ast.Ident) bool {
		if serviceRecvIdent, ok := selectorExpr.X.(*ast.Ident); ok {
			if serviceRecvIdent.Name == expectedRecvIdent.Name {
				return true
			}
		}
		return false
	}

	// e.g. f.queue.Push
	//    ^ident2 ^ident ^method
	// e.g. f.storageService.StorePost
	//      ^ident2  ^ident    ^method
	if funcCall, ok := service.(*ast.CallExpr); ok {
		if methodSel, ok := funcCall.Fun.(*ast.SelectorExpr); ok {
			// get first selector (e.g. storageService)
			if fieldSel, ok := methodSel.X.(*ast.SelectorExpr); ok && hasCurrentServiceRecvIdent(fieldSel, expectedRecvIdent) {
				return true, methodSel.Sel, fieldSel.Sel
			} else if hasCurrentServiceRecvIdent(methodSel, expectedRecvIdent) {
				return true, methodSel.Sel, nil
			}
		}
	}
	return false, nil, nil
}

func (service *Service) attachQueueHandlerMethod(parsedMethod *types.ParsedMethod, dbInstance datastores.DatabaseInstance) {
	parsedMethod.DbInstances = append(parsedMethod.DbInstances, dbInstance)
	service.QueueHandlerMethods[parsedMethod.Name] = parsedMethod
	logger.Logger.Infof("[PARSER] [%s] attached queue handler method for instance (%s): %s", service.Name, dbInstance.String(), parsedMethod.String())
	//logger.Logger.Warnf("[PARSER] [%s] queue handler methods list: %v", service.Name, service.QueueHandlerMethods)
}

func (service *Service) attachInternalMethod(parsedMethod *types.ParsedMethod) {
	service.InternalMethods[parsedMethod.Name] = parsedMethod
	parsedMethod.AttachService(service.GetName())
	logger.Logger.Infof("[PARSER] [%s] attached internal method: %s", service.Name, parsedMethod.String())
	//logger.Logger.Warnf("[PARSER] [%s] internal methods list: %v", service.Name, service.InternalMethods)
}

func (service *Service) attachExportedMethod(parsedMethod *types.ParsedMethod) {
	service.ExposedMethods[parsedMethod.Name] = parsedMethod
	parsedMethod.AttachService(service.GetName())
	service.ExposedMethodsLst = append(service.ExposedMethodsLst, parsedMethod)
	logger.Logger.Infof("[PARSER] [%s] attached exposed method: %s", service.Name, parsedMethod.String())
	//logger.Logger.Warnf("[PARSER] [%s] exposed methods list: %v", service.Name, service.ExposedMethods)
}
