package service

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/golang-collections/collections/stack"

	"analyzer/pkg/datastores"
	"analyzer/pkg/logger"
	"analyzer/pkg/types"
	"analyzer/pkg/types/gotypes"
	"analyzer/pkg/types/objects"
	"analyzer/pkg/utils"
)

type Context struct {
	Block   *types.Block
	Package *types.Package
}

func (context *Context) GetBlock() *types.Block {
	return context.Block
}
func (context *Context) GetPackage() *types.Package {
	return context.Package
}

type Service struct {
	Name            string
	ConstructorName string
	File            *types.File
	ContextStack    *stack.Stack

	//TODO maybe use variable instead of ImplType + Fields
	//Impl   objects.Variable
	ImplName     string
	ImplVariable objects.Object       // variable is either StructVariable or InterfaceVariable, and its type is UserType whose underlying type is the following ImplType:
	ImplType     *gotypes.ServiceType // type is either StructType or InterfaceType
	Fields       map[string]types.Field

	// the map key is the service type (e.g. StorageService in 'storageService StorageService')
	Services   map[string]*Service
	Datastores map[string]datastores.DatabaseInstance
	// safe because methods are unique since Golang does not allow overloading
	// also this captures all exposed methods because they must be defined within the service struct file
	ExposedMethods      map[string]*types.ParsedMethod
	QueueHandlerMethods map[string]*types.ParsedMethod
	InternalMethods     map[string]*types.ParsedMethod
	PackageMethods      map[string]*types.ParsedMethod
	Constructor         *types.ParsedMethod

	QueueField bool

	// helps tracking which ones appear first so we can then build schema
	// -> THIS COULD BE IMPROVED!! BECAUSE IT IS TRACKING THE ORDER OF IMPLEMENTATION AND NOT THE ORDER OF DEFINITION
	ExposedMethodsLst []*types.ParsedMethod
}

func (node *Service) GetFile() *types.File {
	return node.File
}

func (node *Service) HasImplName() bool {
	return node.ImplName != ""
}

func (node *Service) SetImplName(name string) {
	node.ImplName = name
}

func (node *Service) SetImplTypeName(name string) {
	node.ImplType.ImplName = name
}

func (node *Service) GetMethods() map[string]*types.ParsedMethod {
	methods := make(map[string]*types.ParsedMethod, 0)
	maps.Copy(methods, node.ExposedMethods)
	maps.Copy(methods, node.InternalMethods) // internal methods already contain queue handlers methods
	maps.Copy(methods, node.PackageMethods)
	return methods
}

func (node *Service) SetConstructor(constructor *types.ParsedMethod) {
	node.Constructor = constructor
}

func (node *Service) HasDatastore(datastore *datastores.Datastore) bool {
	for _, datastoreInstance := range node.Datastores {
		if datastoreInstance.GetDatastore() == datastore {
			return true
		}
	}
	return false
}

func (node *Service) GetDatastoreFieldIfExists(name string) *types.DatabaseField {
	if field, ok := node.Fields[name]; ok {
		if datastoreField, isDbField := field.(*types.DatabaseField); isDbField {
			return datastoreField
		}
	}
	logger.Logger.Warnf("[SERVICE] [%s] datastore field not found for name (%s) in service fields map: %v", node.GetName(), name, node.Fields)
	return nil
}

func (node *Service) GetFieldIfExists(name string) types.Field {
	if field, ok := node.Fields[name]; ok {
		return field
	}
	logger.Logger.Warnf("[SERVICE] [%s] field not found for name (%s) in service fields map: %v", node.GetName(), name, node.Fields)
	return nil
}

func (node *Service) AddOrGetField(name string, field types.Field) {
	node.Fields[name] = field
}

func (node *Service) AddDatastore(datastore datastores.DatabaseInstance) {
	node.Datastores[datastore.GetName()] = datastore
}

func (node *Service) EnableQueueField() {
	node.QueueField = true
}

func (node *Service) HasQueueField() bool {
	return node.QueueField
}

func (node *Service) NewContext() *Context {
	var newContext *Context
	node.ContextStack.Push(newContext)
	return newContext
}

func (node *Service) GetContext() *Context {
	return node.ContextStack.Peek().(*Context)
}

func (node *Service) GetConstructor() *types.ParsedMethod {
	return node.Constructor
}

func (node *Service) GetImplType() *gotypes.ServiceType {
	if node.ImplType == nil {
		logger.Logger.Fatalf("[SERVICE] [%s] unexpected nil impl type", node.GetName())
	}
	return node.ImplType
}

func (node *Service) GetImplVariable() objects.Object {
	if node.ImplVariable == nil {
		logger.Logger.Fatalf("[SERVICE] [%s] unexpected nil impl variable", node.GetName())
	}
	return node.ImplVariable
}

func (node *Service) SetImplVariableWithType(v objects.Object) {
	logger.Logger.Debugf("[SERVICE] [%s] set impl variable (%s): %s", node.GetName(), utils.GetType(v), v.LongString())
	v = objects.UnwrapAddressVariable(v)
	t := v.GetType()

	if ptrVariable, ok := v.(*objects.PointerObject); ok {
		node.ImplVariable = ptrVariable.GetPointerTo()
		t = ptrVariable.GetPointerType().GetPointerTo()
	} else {
		logger.Logger.Fatalf("[SERVICE] [%s] unexpected type for impl (%s) for variable (%s): %s", node.GetName(), utils.GetType(t), utils.GetType(v), t.String())
	}

	if userType, ok := t.(*gotypes.UserType); ok {
		t = userType.UserType
	} else {
		logger.Logger.Fatalf("[SERVICE] [%s] unexpected type for pointerTo (%s): %s", node.GetName(), utils.GetType(t), t.String())
	}

	if structType, ok := t.(*gotypes.StructType); ok {
		node.ImplType.StructType = structType
	} else if interfaceType, ok := t.(*gotypes.InterfaceType); ok {
		node.ImplType.InterfaceType = interfaceType
	} else {
		logger.Logger.Fatalf("[SERVICE] [%s] unexpected type for underlying userType (%s): %s", node.GetName(), utils.GetType(t), t.String())
	}
	logger.Logger.Debugf("[SERVICE] [%s] set up impl variable: %s", node.GetName(), v.LongString())
}

func (node *Service) PopContext() *Context {
	return node.ContextStack.Pop().(*Context)
}

func (node *Service) MarshalJSON() ([]byte, error) {
	var serviceKeys []string
	for k := range node.Services {
		serviceKeys = append(serviceKeys, k)
	}
	var databaseKeys []string
	for k := range node.Datastores {
		databaseKeys = append(databaseKeys, k)
	}
	fieldTypes := make(map[string]string)
	for name, field := range node.Fields {
		fieldTypes[name] = field.GetTypeName()
	}
	return json.MarshalIndent(&struct {
		Fields    map[string]string `json:"fields,omitempty"`
		Services  []string          `json:"services,omitempty"`
		Databases []string          `json:"databases,omitempty"`
	}{
		Fields:    fieldTypes,
		Services:  serviceKeys,
		Databases: databaseKeys,
	}, "", " ")
}

func (node *Service) GetName() string {
	return node.Name
}

// Remove "Service" at the end and start with lower case
// e.g. StorageService → storage
func (node *Service) GetShortName() string {
	if strings.HasSuffix(node.Name, "Service") {
		shortName := strings.TrimSuffix(node.Name, "Service")
		return strings.ToLower(shortName[:1]) + shortName[1:]
	}
	return strings.ToLower(node.Name)
}

// Add space before "Service"
// e.g. StorageService → Storage Service
func (node *Service) GetSpacedName() string {
	return utils.GetSpacedName(node.Name)
}

func (node *Service) GetPackage() *types.Package {
	return node.GetFile().GetPackage()
}

func (node *Service) GetPackageName() string {
	return node.GetFile().GetPackage().GetName()
}

func (node *Service) GetQueueHandlersForDatabaseInstance(database datastores.DatabaseInstance) []*types.ParsedMethod {
	logger.Logger.Warnf("[SERVICE] [%s] getting queue handlers for datastore (%s): %v", node.GetName(), database.GetName(), node.QueueHandlerMethods)
	if _, ok := node.Datastores[database.GetName()]; !ok {
		return nil
	}

	var handlers []*types.ParsedMethod
	for _, handler := range node.QueueHandlerMethods {
		if slices.Contains(handler.DbInstances, database) {
			handlers = append(handlers, handler)
		}
	}
	return handlers
}

func (node *Service) GetQueueHandlersForDatastore(datastore *datastores.Datastore) []*types.ParsedMethod {
	var handlers []*types.ParsedMethod
	for _, handler := range node.QueueHandlerMethods {
		for _, dbInstances := range handler.DbInstances {
			if dbInstances.GetDatastore() == datastore {
				handlers = append(handlers, handler)
			}
		}
	}
	return handlers
}

func (node *Service) GetExportedMethod(name string) *types.ParsedMethod {
	if m, ok := node.ExposedMethods[name]; ok {
		return m
	}
	logger.Logger.Fatalf("[SERVICE] unknown exposed method (%s) for service (%s)", name, node.GetName())
	return nil
}

func (node *Service) GetInternalMethod(name string) *types.ParsedMethod {
	if m, ok := node.InternalMethods[name]; ok {
		return m
	}
	logger.Logger.Fatalf("[SERVICE] unknown internal method (%s) for service (%s)", name, node.GetName())
	return nil
}

func (node *Service) GetPackageMethod(name string) *types.ParsedMethod {
	if m, ok := node.PackageMethods[name]; ok {
		return m
	}
	logger.Logger.Fatalf("[SERVICE] unknown package method (%s) for service (%s)", name, node.GetName())
	return nil
}

func (node *Service) GetInternalOrPackageMethod(name string) *types.ParsedMethod {
	if internalMethod, ok := node.InternalMethods[name]; ok {
		return internalMethod
	}
	if packageMethod, ok := node.PackageMethods[name]; ok {
		return packageMethod
	}
	// exposed methods can also call themselves
	if exposedMethod, ok := node.ExposedMethods[name]; ok {
		return exposedMethod
	}
	internalMethodsStr := "\t\t\t\t\t\t\t Internal Methods:\n"
	for _, m := range node.InternalMethods {
		internalMethodsStr += fmt.Sprintf("\t\t\t\t\t\t\t - %s\n", m.String())
	}
	packageMethodStr := "\t\t\t\t\t\t\t Package Methods:\n"
	for _, m := range node.PackageMethods {
		packageMethodStr += fmt.Sprintf("\t\t\t\t\t\t\t - %s\n", m.String())
	}
	exposedMethodsStr := "\t\t\t\t\t\t\t Exposed Methods:\n"
	for _, m := range node.ExposedMethods {
		exposedMethodsStr += fmt.Sprintf("\t\t\t\t\t\t\t - %s\n", m.String())
	}
	logger.Logger.Fatalf("[SERVICE] unknown internal/package/exposed method (%s) for service (%s):\n%s%s%s", name, node.GetName(), internalMethodsStr, packageMethodStr, exposedMethodsStr)
	return nil
}

func (node *Service) GetQueueHandlerMethod(name string) *types.ParsedMethod {
	return node.QueueHandlerMethods[name]
}

func (node *Service) String() string {
	str, _ := node.MarshalJSON()
	return string(str)
}

func (node *Service) ServicePath() string {
	return node.GetFile().GetPackage().GetPackagePath() + "." + node.GetName()
}

func (node *Service) ServiceShortPath() string {
	return node.GetFile().GetPackage().GetName() + "." + node.GetName()
}
