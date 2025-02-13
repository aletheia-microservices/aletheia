package abstractgraph

import (
	"fmt"
	"slices"
	"strings"

	"github.com/xwb1989/sqlparser"

	"analyzer/pkg/app"
	"analyzer/pkg/datastores"
	"analyzer/pkg/frameworks/blueprint"
	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
	"analyzer/pkg/types/objects"
	"analyzer/pkg/utils"
)

type DependencySet struct {
	Variable        objects.Object
	Dependencies    []objects.Object
	DependencyNames []string
}

func saveFieldToDatastore(variable objects.Object, entryName string, datastore *datastores.Datastore) {
	objType := variable.GetType()
	datastore.Schema.GetOrCreateField(entryName, objType.GetName(), variable.GetId(), datastore)
	logger.Logger.Infof("[SCHEMA] [%s] added entry (%s): %s", datastore.Name, entryName, objType.LongString())
}

func saveUnfoldedFieldsToDatastore(variable objects.Object, entryName string, datastore *datastores.Datastore) {
	objType := variable.GetType()
	datastore.Schema.GetOrCreateField(entryName, objType.LongString(), variable.GetId(), datastore)
	logger.Logger.Infof("[SCHEMA] [%s] added entry (%s): %s", datastore.Name, entryName, objType.LongString())

	datastore.Schema.GetOrCreateField(objType.GetName(), objType.LongString(), variable.GetId(), datastore)
	logger.Logger.Infof("[SCHEMA] [%s] added field (%s): %s", datastore.Name, objType.GetName(), objType.LongString())

	datastore.Schema.GetOrCreateUnfoldedField(objType.GetName(), objType.LongString(), variable.GetId(), datastore)
	logger.Logger.Infof("[SCHEMA] [%s] added unfolded (entry) field (%s): %s", datastore.Name, objType.GetName(), objType.LongString())

	// add nested unfolded types
	types, names := objType.GetNestedFieldTypes(objType.GetName(), datastore.IsNoSQLDatabase())
	for i, t := range types {
		name := names[i]
		datastore.Schema.GetOrCreateUnfoldedField(name, t.LongString(), 0, datastore)
		logger.Logger.Infof("[SCHEMA] [%s] added nested field (%s): %s", datastore.GetName(), name, t.LongString())
	}
}

func computeSchemaFieldName(object objects.Object, fieldName string) string {
	if fieldName == "" {
		return object.GetType().GetName()
	}
	return fieldName
}

func computeSchemaFieldNameRoot(datastore *datastores.Datastore, fieldName string) string {
	if fieldName == "" {
		return datastore.Schema.GetRootUnfoldedField().GetName()
	}
	return fieldName
}

func TaintDataflowWrite(app *app.App, variable objects.Object, call *AbstractDatabaseCall, datastore *datastores.Datastore, fieldName string, includeNestedFields bool, requestIdx int) {
	fmt.Printf("\n------------- TAINT WRITE DATAFLOW FOR CALL %s @ %s -------------\n\n", call.GetMethodStr(), datastore.Name)
	fmt.Println()

	// taint direct dataflow
	fieldName = computeSchemaFieldName(variable, fieldName)

	rootField := datastore.Schema.GetField(fieldName)
	logger.Logger.Infof("[TAINT WRITE] got root field for name (%s): %s", fieldName, rootField.String())
	df := variable.GetVariableInfo().SetDirectDataflow(datastore.Name, call.Service, variable, rootField, true, requestIdx)
	app.AddDataflow(df, call.ParsedCall)
	logger.Logger.Debugf("[TAINT WRITE DIRECT] %s ---> (%02d) %s [%s]", rootField.GetFullName(), variable.GetId(), variable.String(), utils.GetType(variable))
	if !slices.Contains(app.TaintedVariables[rootField.GetFullName()], variable) {
		app.TaintedVariables[rootField.GetFullName()] = append(app.TaintedVariables[rootField.GetFullName()], variable)
	}
	var taintedVariables []objects.Object

	// taint indirect dataflow
	fieldName = computeSchemaFieldNameRoot(datastore, fieldName)

	var vars []objects.Object
	var names []string
	if includeNestedFields {
		vars, names = objects.GetReversedNestedFieldsAndNames(variable, fieldName, datastore.IsNoSQLDatabase(), datastore.IsQueue())
	} else {
		vars = []objects.Object{variable}
		names = []string{fieldName}
	}

	for i, v := range vars {
		logger.Logger.Infof("[TENTATIVE TAINT WRITE VAR] [%s] %s", utils.GetType(v), v.LongString())
		dbField := datastore.Schema.GetField(names[i])
		deps := v.GetNestedDependencies(true)

		for _, dep := range deps {
			logger.Logger.Debugf("visiting dep: %s", dep.String())
			for _, ref := range dep.GetVariableInfo().GetReferences() {
				logger.Logger.Warnf("ref: %s", ref.String())
			}
			if !slices.Contains(taintedVariables, dep) {
				df := v.GetVariableInfo().SetIndirectDataflow(datastore.Name, call.Service, dep, variable, dbField, true, requestIdx)
				app.AddDataflow(df, call.ParsedCall)
				logger.Logger.Debugf("\t\t[TAINT WRITE INDIRECT] %s ---> (%02d) %s [%s]", dbField.GetFullName(), dep.GetId(), dep.String(), utils.GetType(dep))

				taintedVariables = append(taintedVariables, dep)
				app.AddTaintedVariableIfNotExists(dbField.GetFullName(), dep)
			}
		}
	}
	fmt.Println()
}

func TaintDataflowReadQueue(app *app.App, variable objects.Object, call *AbstractDatabaseCall, datastore *datastores.Datastore, fieldName string, requestIdx int) {
	fmt.Printf("\n------------- TAINT READ DATAFLOW FOR CALL %s @ %s -------------\n\n", call.GetName(), datastore.Name)
	fmt.Println()

	// taint direct dataflow
	rootField := datastore.Schema.GetField(fieldName)
	df := variable.GetVariableInfo().SetDirectDataflow(datastore.Name, call.Service, variable, rootField, false, requestIdx)
	app.AddDataflow(df, call.ParsedCall)
	logger.Logger.Debugf("[TAINT READ DIRECT] %s ---> (%02d) %s [%s]", rootField.GetFullName(), variable.GetId(), variable.String(), utils.GetType(variable))
	if !slices.Contains(app.TaintedVariables[rootField.GetFullName()], variable) {
		app.TaintedVariables[rootField.GetFullName()] = append(app.TaintedVariables[rootField.GetFullName()], variable)
	}
	var taintedVariables []objects.Object

	// taint indirect dataflow
	rootUnfoldedField := datastore.Schema.GetRootUnfoldedField()
	logger.Logger.Debugf("rootUnfoldedField = %v", rootUnfoldedField)
	vars, names := objects.GetReversedNestedFieldsAndNames(variable, rootUnfoldedField.GetName(), datastore.IsNoSQLDatabase(), datastore.IsQueue())

	for i, v := range vars {
		dbField := datastore.Schema.GetField(names[i])
		deps := v.GetNestedDependencies(false)
		logger.Logger.Infof("[TENTATIVE TAINT READ VAR] [%s] (%02d) (NUM DEPS = %d) %s", utils.GetType(v), v.GetId(), len(deps), v.LongString())
		for _, dep := range deps {
			if !slices.Contains(taintedVariables, dep) {
				df := v.GetVariableInfo().SetIndirectDataflow(datastore.Name, call.Service, dep, variable, dbField, false, requestIdx)
				app.AddDataflow(df, call.ParsedCall)
				logger.Logger.Debugf("\t\t[TAINT READ INDIRECT] %s ---> (%02d) %s [%s]", dbField.GetFullName(), dep.GetId(), dep.String(), utils.GetType(dep))

				taintedVariables = append(taintedVariables, dep)
				app.AddTaintedVariableIfNotExists(dbField.GetFullName(), dep)
			}
		}
	}
	fmt.Println()
}

// aka TaintDataflowReadUnnamed
func TaintDataflowNoSQL(app *app.App, obj objects.Object, call *AbstractDatabaseCall, datastore *datastores.Datastore, fieldName string, queryField bool, write bool, requestIdx int) {
	fmt.Printf("\n------------- TAINT READ (DOC) DATAFLOW FOR CALL %s @ %s -------------\n\n", call.GetMethodStr(), datastore.Name)
	fmt.Println()

	var field *datastores.Field
	if queryField {
		field = datastore.Schema.GetOrCreateUnfoldedField(fieldName, obj.GetType().GetName(), obj.GetId(), datastore)
	} else { // cursor
		field = datastore.Schema.GetOrCreateRootField(datastores.ROOT_FIELD_NAME_NOSQL, datastores.UNKNOWN_FIELD_TYPE, -1, datastore)
	}

	// taint direct dataflow
	df := obj.GetVariableInfo().SetDirectDataflow(datastore.Name, call.Service, obj, field, false, requestIdx)
	app.AddDataflow(df, call.ParsedCall)
	logger.Logger.Debugf("[TAINT READ (DOC) DIRECT] %s ---> (%02d) %s [%s]", field.GetFullName(), obj.GetId(), obj.String(), utils.GetType(obj))
	if !slices.Contains(app.TaintedVariables[field.GetFullName()], obj) {
		app.TaintedVariables[field.GetFullName()] = append(app.TaintedVariables[field.GetFullName()], obj)
	}

	var taintedVariables []objects.Object
	logger.Logger.Infof("[TENTATIVE TAINT READ (DOC) VAR] [%s] (%02d) @ %s", utils.GetType(obj), obj.GetId(), obj.LongString())
	obj = getTargetVariableIfNoSQLCursorRead(datastore, obj)
	deps := obj.GetNestedDependencies(false)
	logger.Logger.Infof("[TENTATIVE TAINT READ (DOC) VAR] [%s] NUM DEPS = %d @ %s", utils.GetType(obj), len(deps), obj.LongString())
	var prefix string
	for _, dep := range deps {
		if !slices.Contains(taintedVariables, dep) {
			typeName := prefix + dep.GetType().GetName()

			taintDataflowNoSQLHelper(app, obj, dep, field, call, datastore, typeName, queryField, write, requestIdx)
			taintedVariables = append(taintedVariables, dep)

			if fieldObj, ok := dep.(*objects.FieldObject); ok {
				// if it is a nested field then it was captured in "deps" and will be visited eventually
				// otherwise, we want to capture any other types e.g. BasicObjects to ensure the typeName aka fieldName assigned is the same as the upper field
				if _, underlyingIsNestedField := fieldObj.WrappedVariable.(*objects.FieldObject); !underlyingIsNestedField {
					taintDataflowNoSQLHelper(app, obj, fieldObj.WrappedVariable, field, call, datastore, typeName, queryField, write, requestIdx)
					taintedVariables = append(taintedVariables, fieldObj.WrappedVariable)
				}
			}

			if structObj, ok := dep.(*objects.StructObject); ok {
				if userType, ok := structObj.GetType().(*gotypes.UserType); ok {
					prefix += userType.GetName() + "."
				}
			}
		}
	}

	if call.GetName() == "FindOne" && datastore.Name == "product_db" {
		logger.Logger.Debug("HERE!")
	}

	fmt.Println()
}

func taintDataflowNoSQLHelper(app *app.App, obj objects.Object, dep objects.Object, field *datastores.Field, call *AbstractDatabaseCall, datastore *datastores.Datastore, typeName string, queryField bool, write bool, requestIdx int) {
	if queryField { // query
		//logger.Logger.Debugf("query field? YES! for typename = %s", typeName)
		df := obj.GetVariableInfo().SetIndirectDataflow(datastore.Name, call.Service, dep, obj, field, write, requestIdx)
		app.AddDataflow(df, call.ParsedCall)
	} else { // cursor
		//logger.Logger.Debugf("query field? NO! for typename = %s", typeName)
		entry := datastores.NewField(typeName, typeName, 0, datastore)
		df := obj.GetVariableInfo().SetIndirectDataflow(datastore.Name, call.Service, dep, obj, entry, write, requestIdx)
		app.AddDataflow(df, call.ParsedCall)
	}
	app.AddTaintedVariableIfNotExists(dep.GetType().GetName(), dep)
}

// aka TaintDataflowReadUnnamed
func TaintDataflowReadCache(app *app.App, obj objects.Object, fieldName string, call *AbstractDatabaseCall, datastore *datastores.Datastore, requestIdx int) {
	fmt.Printf("\n------------- [CACHE %s] TAINT READ DATAFLOW FOR CALL %s -------------\n\n\n", datastore.GetName(), call.GetMethodStr())
	field := datastore.Schema.GetOrCreateField(fieldName, datastores.UNKNOWN_FIELD_TYPE, -1, datastore)
	taintDataflowReadCacheSQLHelper(app, obj, call, datastore, field, requestIdx)
}

func TaintDataflowReadSQL(app *app.App, obj objects.Object, fieldName string, call *AbstractDatabaseCall, datastore *datastores.Datastore, requestIdx int, isSelectAll bool) {
	fmt.Printf("\n------------- [SQL %s] TAINT READ FOR CALL %s -------------\n\n\n", datastore.GetName(), call.GetMethodStr())
	fieldType := datastores.UNKNOWN_FIELD_TYPE
	if isSelectAll {
		fieldType = "SQL Table"
	}
	field := datastore.Schema.GetOrCreateField(fieldName, fieldType, -1, datastore)
	taintDataflowReadCacheSQLHelper(app, obj, call, datastore, field, requestIdx)
}

func taintDataflowReadCacheSQLHelper(app *app.App, obj objects.Object, call *AbstractDatabaseCall, datastore *datastores.Datastore, field *datastores.Field, requestIdx int) {
	// taint direct dataflow
	df := obj.GetVariableInfo().SetDirectDataflow(datastore.Name, call.Service, obj, field, false, requestIdx)
	app.AddDataflow(df, call.ParsedCall)
	logger.Logger.Debugf("[%s - TAINT READ DIRECT] %s ---> (%02d) %s [%s]", datastore.GetTypeString(), field.GetFullName(), obj.GetId(), obj.String(), utils.GetType(obj))
	if !slices.Contains(app.TaintedVariables[field.GetFullName()], obj) {
		app.TaintedVariables[field.GetFullName()] = append(app.TaintedVariables[field.GetFullName()], obj)
	}

	var taintedVariables []objects.Object
	logger.Logger.Infof("[%s - TENTATIVE TAINT READ VAR] [%s] (%02d) @ %s", datastore.GetTypeString(), utils.GetType(obj), obj.GetId(), obj.LongString())

	deps := obj.GetNestedDependencies(false)
	logger.Logger.Infof("[%s - TENTATIVE TAINT READ VAR] [%s] NUM DEPS = %d @ %s", datastore.GetTypeString(), utils.GetType(obj), len(deps), obj.LongString())
	for _, dep := range deps {
		if !slices.Contains(taintedVariables, dep) {
			typeName := dep.GetType().GetName()

			entry := datastores.NewField(typeName, typeName, 0, datastore)
			df := obj.GetVariableInfo().SetIndirectDataflow(datastore.Name, call.Service, dep, obj, entry, false, requestIdx)
			app.AddDataflow(df, call.ParsedCall)

			logger.Logger.Debugf("\t\t[%s - TAINT READ INDIRECT] <unnamed> ---> (%02d) %s [%s]", datastore.GetTypeString(), dep.GetId(), dep.String(), utils.GetType(dep))

			taintedVariables = append(taintedVariables, dep)
			app.AddTaintedVariableIfNotExists(dep.GetType().GetName(), dep)
		}
	}

	fmt.Println()
}

func getTargetVariableIfNoSQLCursorRead(datastore *datastores.Datastore, v objects.Object) objects.Object {
	if datastore.IsNoSQLDatabase() {
		if blueprintVariable, ok := v.(*blueprint.BlueprintBackendObject); ok && blueprintVariable.IsNoSQLCursor() {
			targetVariable := blueprintVariable.GetTargetObject()
			if ptrTargetVariable, ok := targetVariable.(*objects.PointerObject); ok {
				return ptrTargetVariable.PointerTo
			} else if ifaceTargetVariable, ok := targetVariable.(*objects.InterfaceObject); ok {
				logger.Logger.Fatalf("TODO!!!! %s", ifaceTargetVariable.String())
			} else { //FIXME: cursor.One(ctx, existing) --> this is happening in trainticket where "existing" object passed as a struct and not the ptr to it, but is this even possible
				logger.Logger.Warnf("????? [%s] %s", utils.GetType(targetVariable), targetVariable.String())
				return targetVariable
			}
		}
	}
	return v
}

type NoSQLQueryDocument struct {
	FieldName string
	Object    objects.Object
}

func (obj *NoSQLQueryDocument) String() string {
	return fmt.Sprintf("%s: %s", obj.FieldName, obj.Object.String())
}

func GetNoSQLQueryDocument_DEPRECATED(datastore *datastores.Datastore, variable objects.Object) []NoSQLQueryDocument {
	// should be a bson.D which is a slice with many (inline) structures
	// e.g. query := bson.D{{Key: "productid", Value: productID}}
	if sliceVariable, ok := variable.(*objects.SliceObject); ok {
		var queryObjects []NoSQLQueryDocument
		for i, elem := range sliceVariable.GetElements() {
			if structVariable, ok := elem.(*objects.StructObject); ok {
				logger.Logger.Warnf("MAP: %v", structVariable.GetFieldsMap())
				logger.Logger.Warnf("LIST: %v", structVariable.GetFieldsList())
				logger.Logger.Warnf("STRUCT: %v", structVariable.LongString())
				var key, val objects.Object
				if structVariable.NumFieldsList() != 0 {
					key = structVariable.GetFieldAt(0).GetWrappedVariable()
					val = structVariable.GetFieldAt(1).GetWrappedVariable()
				} else {
					key = structVariable.GetFieldByKey("Key").GetWrappedVariable()
					val = structVariable.GetFieldByKey("Value").GetWrappedVariable()
				}
				queryObj := NoSQLQueryDocument{
					FieldName: datastore.Schema.GetRootFieldName() + "." + key.GetType().GetBasicValue(),
					Object:    val,
				}
				queryObjects = append(queryObjects, queryObj)
				logger.Logger.Infof("[QUERY OBJ #%d] %s", i, queryObj.String())
				return queryObjects
			}
		}
		return nil
	}
	return nil
}

func GetNoSQLQueryDocument(datastore *datastores.Datastore, variable objects.Object) []NoSQLQueryDocument {
	// should be a bson.D which is a slice with many (inline) structures
	// e.g. query := bson.D{{Key: "productid", Value: productID}}
	if sliceVariable, ok := variable.(*objects.SliceObject); ok {
		var queryObjects []NoSQLQueryDocument
		for _, elem := range sliceVariable.GetElements() {
			if structVariable, ok := elem.(*objects.StructObject); ok {
				logger.Logger.Infof("[DOC_ELEM] [%s] %s", utils.GetType(structVariable), structVariable.LongString())

				key := structVariable.GetFieldAt(0).GetWrappedVariable()

				if key.GetType().GetBasicValue() == "$and" {
					// detect and handle the "$and" clause as an array of conditions
					//
					// e.g.
					// query := bson.D{{"$and", bson.A{
					//		bson.D{{"startstation", start}},
					//		bson.D{{"endstation", end}},
					// }}}
					// handle $and as an array of conditions
					andArray := structVariable.GetFieldAt(1).GetWrappedVariable()
					if andSlice, ok := andArray.(*objects.SliceObject); ok {
						logger.Logger.Warnf("[$AND_SLICE] [%s] %s", utils.GetType(andSlice), andSlice.LongString())
						for i, andElem := range andSlice.GetElements() {
							logger.Logger.Warnf("[$AND_ELEM #%d] [%s] %s", i, utils.GetType(andElem), andElem.LongString())
							for j, andNestedElem := range GetNoSQLQueryDocument(datastore, andElem) {
								logger.Logger.Warnf("[$AND_NESTED_ELEM #%d.%d] [%s] %s", i, j, utils.GetType(andNestedElem.Object), andNestedElem.Object.LongString())
								queryObjects = append(queryObjects, andNestedElem)
							}
						}
					}
				} else {
					// normal key-value pairs
					val := structVariable.GetFieldAt(1).GetWrappedVariable()
					queryObj := NoSQLQueryDocument{
						FieldName: datastore.Schema.GetRootFieldName() + "." + key.GetType().GetBasicValue(),
						Object:    val,
					}
					queryObjects = append(queryObjects, queryObj)
					logger.Logger.Debugf("[DOC QUERY OBJ] %s", queryObj.String())
				}
			}
		}
		return queryObjects
	}
	return nil
}

func referenceTaintedDataflowForNestedField(writtenVariable objects.Object, datastore *datastores.Datastore, fieldName string, requestIdx int) {
	fmt.Printf("\n------------- REFERENCE TAINTED DATAFLOW FOR WRITTEN VARIABLE %s @ %s -------------\n\n", writtenVariable.GetType().GetName(), datastore.Name)
	fmt.Println()
	dbField := datastore.Schema.GetField(fieldName)
	deps := writtenVariable.GetNestedDependencies(false)
	logger.Logger.Infof("[TENTATIVE REF TAINTED VAR] [%s] (%02d) %s", utils.GetType(writtenVariable), writtenVariable.GetId(), writtenVariable.LongString())
	for _, dep := range deps {
		for _, df := range dep.GetVariableInfo().GetAllDataflows() {
			if df.Datastore != datastore.Name {
				var mandatoryPrefix string
				var isMandatory bool
				if df.InRequestIndex(requestIdx) {
					mandatoryPrefix = "[+ MANDATORY]"
					isMandatory = true
				}

				if !dbField.HasReference(df.Field) {
					constraint := dbField.CreateAndAddReference(df.Field, isMandatory)
					datastore.GetSchema().AddConstraint(constraint)

				}

				datastore.AddReferencingDatastoreIfNotExists(df.Field.GetDatastore())
				logger.Logger.Debugf("[SCHEMA REFERENCE] %s (%s -> %s) from dependency [%s]: %s", mandatoryPrefix, dbField.GetFullName(), df.Field.GetFullName(), utils.GetType(dep), dep.String())
			}
		}

	}
	fmt.Println()
}

func referenceTaintedDataflowForAllNestedFields(writtenVariable objects.Object, datastore *datastores.Datastore, requestIdx int) {
	fmt.Printf("\n------------- REFERENCE TAINTED DATAFLOW FOR WRITTEN VARIABLE %s @ %s -------------\n\n", writtenVariable.GetType().GetName(), datastore.Name)
	fmt.Println()
	vars, names := objects.GetReversedNestedFieldsAndNames(writtenVariable, "", datastore.IsNoSQLDatabase(), datastore.IsQueue())
	for i, variable := range vars {
		referenceTaintedDataflowForNestedField(variable, datastore, names[i], requestIdx)
	}
}

func BuildSchema(app *app.App, frontends []string, entryNodes []AbstractNode) {
	visited := make(map[AbstractNode]bool, 0)
	for _, frontend := range frontends {
		for requestIdx, exposedMethod := range app.Services[frontend].ExposedMethodsLst {
			logger.Logger.Infof("[BUILD SCHEMA] iterating exposed method %d: %s", requestIdx, exposedMethod.String())
			for _, entry := range entryNodes {
				if entry.GetName() == exposedMethod.GetName() {
					if _, isVisited := visited[entry]; !isVisited {
						//app.ResetAllDataflows() //FIXME!!!!!
						doBuildSchema(app, entry, requestIdx)
						visited[entry] = true
					}
				}
			}

		}
	}
}

var visitedNodes []AbstractNode
var writtenDatastores = make(map[string]bool, 0)

func doBuildSchema(app *app.App, node AbstractNode, requestIdx int) bool {
	if dbCall, ok := node.(*AbstractDatabaseCall); ok && dbCall.ParsedCall.Method.IsRead() {
		datastore := dbCall.DbInstance.GetDatastore()
		params := dbCall.GetParams()
		returns := dbCall.GetReturns()

		if blueprintBackendMethod := dbCall.ParsedCall.Method.(*blueprint.BackendMethod); blueprintBackendMethod != nil {
			switch datastore.Type {
			case datastores.Queue:
				msg := params[1]
				TaintDataflowReadQueue(app, msg, dbCall, datastore, datastores.ROOT_FIELD_NAME_QUEUE, requestIdx)

			case datastores.NoSQL:
				cursor, query := returns[0], params[1]
				TaintDataflowNoSQL(app, cursor, dbCall, datastore, datastores.ROOT_FIELD_NAME_NOSQL, false, false, requestIdx)
				queryObjs := GetNoSQLQueryDocument(datastore, query)
				for _, v := range queryObjs {
					TaintDataflowNoSQL(app, v.Object, dbCall, datastore, v.FieldName, true, false, requestIdx)
				}

			case datastores.Cache:
				key, value := params[1], params[2]
				TaintDataflowReadCache(app, key, datastores.ROOT_FIELD_NAME_CACHE_KEY, dbCall, datastore, requestIdx)
				TaintDataflowReadCache(app, value, datastores.ROOT_FIELD_NAME_CACHE_VALUE, dbCall, datastore, requestIdx)

			case datastores.SQL:
				if blueprintBackendMethod.IsRelationalDBSelectCall() {
					target, query, args := params[1], params[2], params[3:]
					selectedFieldNames, filterFieldNames, filterFieldObjs := parseSQLReadSelect(query, args)
					for idx, fieldName := range filterFieldNames {
						fieldObj := filterFieldObjs[idx]
						TaintDataflowReadSQL(app, fieldObj, fieldName, dbCall, datastore, requestIdx, false)
					}
					TaintDataflowReadSQL(app, target, selectedFieldNames[0], dbCall, datastore, requestIdx, true)
				} else if blueprintBackendMethod.IsRelationalDBQueryCall() {
					logger.Logger.Fatalf("TODO!! implement cursor for sql similar to nosql mongodb")
				}

			default:
				logger.Logger.Fatalf("[SCHEMA] unknown type of datastore (%s) to parse call: %s", utils.GetType(datastore), dbCall.String())
			}
		}

	} else if dbCall, ok := node.(*AbstractDatabaseCall); ok && dbCall.ParsedCall.Method.IsWrite() {
		datastore := dbCall.DbInstance.GetDatastore()
		if found := writtenDatastores[datastore.Name]; !found {
			writtenDatastores[datastore.Name] = true
		}
		params := dbCall.GetParams()
		logger.Logger.Infof("[SCHEMA] [%s] building schema based on abstract node (%s)", datastore.Name, dbCall.GetName())

		if blueprintBackendMethod := dbCall.ParsedCall.Method.(*blueprint.BackendMethod); blueprintBackendMethod != nil {
			switch datastore.Type {
			case datastores.Queue:
				msg := params[1]
				saveUnfoldedFieldsToDatastore(msg, datastores.ROOT_FIELD_NAME_QUEUE, datastore)
				TaintDataflowWrite(app, msg, dbCall, datastore, "", true, requestIdx)
				referenceTaintedDataflowForAllNestedFields(msg, datastore, requestIdx)

			case datastores.NoSQL:
				doc := params[1]
				saveUnfoldedFieldsToDatastore(doc, datastores.ROOT_FIELD_NAME_NOSQL, datastore)
				TaintDataflowWrite(app, doc, dbCall, datastore, "", true, requestIdx)
				referenceTaintedDataflowForAllNestedFields(doc, datastore, requestIdx)

			case datastores.Cache:
				key, value := params[1], params[2]
				saveFieldToDatastore(key, datastores.ROOT_FIELD_NAME_CACHE_KEY, datastore)
				saveFieldToDatastore(value, datastores.ROOT_FIELD_NAME_CACHE_VALUE, datastore)
				TaintDataflowWrite(app, key, dbCall, datastore, datastores.ROOT_FIELD_NAME_CACHE_KEY, false, requestIdx)
				TaintDataflowWrite(app, value, dbCall, datastore, datastores.ROOT_FIELD_NAME_CACHE_VALUE, false, requestIdx)
				referenceTaintedDataflowForNestedField(key, datastore, datastores.ROOT_FIELD_NAME_CACHE_KEY, requestIdx)
				referenceTaintedDataflowForNestedField(value, datastore, datastores.ROOT_FIELD_NAME_CACHE_VALUE, requestIdx)

			case datastores.SQL:
				if blueprintBackendMethod.IsRelationalDBExecCall() {
					query, args := params[1], params[2:]
					newFieldNames, newFieldObjs, filterFieldNames, filterFieldObjs := parseSQLWrite(query, args)
					for idx, fieldName := range newFieldNames {
						fieldObj := newFieldObjs[idx]
						saveFieldToDatastore(fieldObj, fieldName, datastore)
						TaintDataflowWrite(app, fieldObj, dbCall, datastore, fieldName, false, requestIdx)
						referenceTaintedDataflowForNestedField(fieldObj, datastore, fieldName, requestIdx)
					}
					for idx, fieldName := range filterFieldNames {
						fieldObj := filterFieldObjs[idx]
						TaintDataflowReadSQL(app, fieldObj, fieldName, dbCall, datastore, requestIdx, false)
						referenceTaintedDataflowForNestedField(fieldObj, datastore, fieldName, requestIdx)
					}
				}

			default:
				logger.Logger.Fatalf("[SCHEMA] unknown type of datastore (%s) to parse call: %s", utils.GetType(datastore), dbCall.String())
			}
		}

	}

	visitedNodes = append(visitedNodes, node)

	for _, child := range node.GetChildren() {
		doBuildSchema(app, child, requestIdx)
	}
	return true
}

type tableNameAlias struct {
	alias string
	name  string
}

func parseColumnName(compliantName string) (string, string) {
	splits := strings.SplitAfterN(compliantName, ".", 1)
	if len(splits) == 2 {
		return splits[0], splits[1]
	}
	return "", compliantName
}

func parseTableName(prefixTableName string, tableNameAliasLst []tableNameAlias) string {
	if prefixTableName == "" {
		return tableNameAliasLst[0].name
	}
	for _, t := range tableNameAliasLst {
		if t.alias != "" && prefixTableName == t.alias {
			return t.name
		}
		if prefixTableName == t.name {
			return t.name
		}
	}
	logger.Logger.Fatal("unexpected")
	return tableNameAliasLst[0].name
}

func parseSQLTableExprs(tableExprs sqlparser.TableExprs) []tableNameAlias {
	var tableNameAliasLst []tableNameAlias
	for _, table := range tableExprs {
		if aliasedTableExpr, ok := table.(*sqlparser.AliasedTableExpr); ok {
			if tableName, ok := aliasedTableExpr.Expr.(sqlparser.TableName); ok {
				tableNameAliasLst = append(tableNameAliasLst, tableNameAlias{alias: aliasedTableExpr.As.CompliantName(), name: tableName.Name.CompliantName()})
			}
		}
	}
	return tableNameAliasLst
}

func parseSQLWhere(args []objects.Object, stmtWhere *sqlparser.Where, tableNameAliasLst []tableNameAlias, argIdx int, filterFieldNames []string, filterFieldObjs []objects.Object) (int, []string, []objects.Object){
	if comparisonExpr, ok := stmtWhere.Expr.(*sqlparser.ComparisonExpr); ok {
		var leftFieldName string
		var rightFieldObj objects.Object
		if col, ok := comparisonExpr.Left.(*sqlparser.ColName); ok {
			prefixTableName, columnName := parseColumnName(string(col.Name.CompliantName()))
			tableName := parseTableName(prefixTableName, tableNameAliasLst)
			leftFieldName = tableName + "." + columnName
		}
		if sqlVal, ok := comparisonExpr.Right.(*sqlparser.SQLVal); ok {
			if sqlVal.Type == sqlparser.ValArg { //placeholder (e.g., '?' that is then parsed into ':v1', ':v2', etc.)
				rightFieldObj = args[argIdx]
				argIdx++
			}
		}
		logger.Logger.Infof("[SCHEMA] [WHERE]: %s -> %s", leftFieldName, rightFieldObj)
		filterFieldNames = append(filterFieldNames, leftFieldName)
		filterFieldObjs = append(filterFieldObjs, rightFieldObj)
	}
	return argIdx, filterFieldNames, filterFieldObjs
}

func parseSQLReadSelect(query objects.Object, args []objects.Object) ([]string, []string, []objects.Object) {
	sql := query.GetType().GetBasicValue()
	logger.Logger.Infof("[SCHEMA] parsing sql stmt: %s", sql)

	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		logger.Logger.Fatalf("[SCHEMA] unable to parse sql query (%s): %s", sql, err.Error())
	}

	argIdx := 0
	var selectedFieldNames []string
	var filterFieldNames []string
	var filterFieldObjs []objects.Object
	var tableNameAliasLst []tableNameAlias

	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		tableNameAliasLst = parseSQLTableExprs(stmt.From)

		var readAllFields bool
		for _, expr := range stmt.SelectExprs {
			if _, ok := expr.(*sqlparser.StarExpr); ok {
				readAllFields = true
				selectedFieldNames = append(selectedFieldNames, tableNameAliasLst[0].name+".*")
				logger.Logger.Tracef("[SCHEMA] found sqlparser.StarExpr (%t)", readAllFields)
			} else if aliasedExpr, ok := expr.(*sqlparser.AliasedExpr); ok {
				if valTuple, ok := aliasedExpr.Expr.(sqlparser.ValTuple); ok {
					for rowIdx, expr := range valTuple {
						if col, ok := expr.(*sqlparser.ColName); ok {
							prefixTableName, columnName := parseColumnName(string(col.Name.CompliantName()))
							tableName := parseTableName(prefixTableName, tableNameAliasLst)
							fieldName := tableName + "." + columnName

							logger.Logger.Infof("[SCHEMA] [SELECT record %d/%d]: %s", rowIdx+1, len(valTuple), fieldName)
							selectedFieldNames = append(selectedFieldNames, fieldName)
						}
					}
				}
			}
		}

		_, filterFieldNames, filterFieldObjs = parseSQLWhere(args, stmt.Where, tableNameAliasLst, argIdx, filterFieldNames, filterFieldObjs)

	default:
		logger.Logger.Fatalf("[SCHEMA] Unsupported SQL statement: %s", sql)
	}

	return selectedFieldNames, filterFieldNames, filterFieldObjs
}

func parseSQLWrite(query objects.Object, args []objects.Object) ([]string, []objects.Object, []string, []objects.Object) {
	sql := query.GetType().GetBasicValue()
	logger.Logger.Infof("[SCHEMA] parsing sql stmt: %s", sql)

	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		logger.Logger.Fatalf("[SCHEMA] unable to parse sql query (%s): %s", sql, err.Error())
	}

	argIdx := 0
	var writtenFieldNames []string
	var writtenFieldObjs []objects.Object
	var filterFieldNames []string
	var filterFieldObjs []objects.Object

	switch stmt := stmt.(type) {
	case *sqlparser.Insert:
		if values, ok := stmt.Rows.(sqlparser.Values); ok {
			for rowIdx, tuple := range values {
				for colIdx, expr := range tuple {
					if sqlVal, ok := expr.(*sqlparser.SQLVal); ok {
						stmt.Table.Name.CompliantName()
						if sqlVal.Type == sqlparser.ValArg { // placeholder (e.g., '?' that is then parsed into ':v1', ':v2', etc.)
							fieldName := stmt.Table.Name.CompliantName() + "." + stmt.Columns[colIdx].CompliantName()
							fieldObj := args[argIdx]
							placeholderVal := string(sqlVal.Val)
							logger.Logger.Infof("[SCHEMA] (record %d/%d) INSERT %s = (%s) -> %s", rowIdx+1, len(values), fieldName, placeholderVal, fieldObj)
							writtenFieldNames = append(writtenFieldNames, fieldName)
							writtenFieldObjs = append(writtenFieldObjs, fieldObj)
							argIdx++
						}
					}
				}
			}
		} else {
			logger.Logger.Fatalf("[SCHEMA] unexpected type %T for rows in sql insert: %s", stmt.Rows, sql)
		}
	case *sqlparser.Update:
		tableNameAliasLst := parseSQLTableExprs(stmt.TableExprs)

		for _, expr := range stmt.Exprs {

			prefixTableName, columnName := parseColumnName(string(expr.Name.Name.CompliantName()))
			tableName := parseTableName(prefixTableName, tableNameAliasLst)
			fieldName := tableName + "." + columnName

			if sqlVal, ok := expr.Expr.(*sqlparser.SQLVal); ok {
				if sqlVal.Type == sqlparser.ValArg { // placeholder (e.g., '?', parsed as ':v1', ':v2')
					fieldObj := args[argIdx]
					placeholderVal := string(sqlVal.Val)
					logger.Logger.Infof("[SCHEMA] SET %s = (%s) -> %s", fieldName, placeholderVal, fieldObj)
					writtenFieldNames = append(writtenFieldNames, fieldName)
					writtenFieldObjs = append(writtenFieldObjs, fieldObj)
					argIdx++
				}
			}
		}

		_, filterFieldNames, filterFieldObjs = parseSQLWhere(args, stmt.Where, tableNameAliasLst, argIdx, filterFieldNames, filterFieldObjs)
	}
	return writtenFieldNames, writtenFieldObjs, filterFieldNames, filterFieldObjs
}
