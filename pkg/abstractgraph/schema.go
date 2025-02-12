package abstractgraph

import (
	"fmt"
	"slices"

	/* "github.com/golang-collections/collections/stack" */

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
		entry := datastores.NewEntry(typeName, typeName, 0, datastore)
		df := obj.GetVariableInfo().SetIndirectDataflow(datastore.Name, call.Service, dep, obj, entry, write, requestIdx)
		app.AddDataflow(df, call.ParsedCall)
	}
	app.AddTaintedVariableIfNotExists(dep.GetType().GetName(), dep)
}

// aka TaintDataflowReadUnnamed
func TaintDataflowReadCache(app *app.App, obj objects.Object, fieldName string, call *AbstractDatabaseCall, datastore *datastores.Datastore, requestIdx int) {
	fmt.Printf("\n------------- TAINT READ (KV) DATAFLOW FOR CALL %s @ %s -------------\n\n", call.GetMethodStr(), datastore.Name)
	fmt.Println()

	// taint direct dataflow
	rootField := datastore.Schema.GetOrCreateField(fieldName, datastores.UNKNOWN_FIELD_TYPE, -1, datastore)
	df := obj.GetVariableInfo().SetDirectDataflow(datastore.Name, call.Service, obj, rootField, false, requestIdx)
	app.AddDataflow(df, call.ParsedCall)
	logger.Logger.Debugf("[TAINT READ (KV) DIRECT] %s ---> (%02d) %s [%s]", rootField.GetFullName(), obj.GetId(), obj.String(), utils.GetType(obj))
	if !slices.Contains(app.TaintedVariables[rootField.GetFullName()], obj) {
		app.TaintedVariables[rootField.GetFullName()] = append(app.TaintedVariables[rootField.GetFullName()], obj)
	}

	var taintedVariables []objects.Object
	logger.Logger.Infof("[TENTATIVE TAINT READ (KV) VAR] [%s] (%02d) @ %s", utils.GetType(obj), obj.GetId(), obj.LongString())
	obj = getTargetVariableIfNoSQLCursorRead(datastore, obj)
	deps := obj.GetNestedDependencies(false)
	logger.Logger.Infof("[TENTATIVE TAINT READ (KV) VAR] [%s] NUM DEPS = %d @ %s", utils.GetType(obj), len(deps), obj.LongString())
	for _, dep := range deps {
		if !slices.Contains(taintedVariables, dep) {
			typeName := dep.GetType().GetName()

			entry := datastores.NewEntry(typeName, typeName, 0, datastore)
			df := obj.GetVariableInfo().SetIndirectDataflow(datastore.Name, call.Service, dep, obj, entry, false, requestIdx)
			app.AddDataflow(df, call.ParsedCall)

			logger.Logger.Debugf("\t\t[TAINT READ (KV) INDIRECT] <unnamed> ---> (%02d) %s [%s]", dep.GetId(), dep.String(), utils.GetType(dep))

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

							/* if andStruct, ok := andElem.(*objects.StructObject); ok {
								key := andStruct.GetFieldAt(0).GetWrappedVariable()
								val := andStruct.GetFieldAt(1).GetWrappedVariable()
								queryObj := NoSQLQueryDocument{
									FieldName: datastore.Schema.GetRootFieldName() + "." + key.GetType().GetBasicValue(),
									Object:    val,
								}
								queryObjects = append(queryObjects, queryObj)
								logger.Logger.Infof("[QUERY OBJ $AND] %s", queryObj.String())
							} */
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
				if df.InRequestIndex(requestIdx) {
					mandatoryPrefix = "* [MANDATORY]"
					dbField.AddMandatoryReference(df.Field)
				}
				datastore.Schema.AddForeignReferenceToField(dbField, df.Field)
				datastore.AddReferencingDatastoreIfNotExists(df.Field.GetDatastore())
				//logger.Logger.Fatalf("%s: %s", datastore.GetName(), datastore.ReferencedByDatastores[0].GetName())
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

/* var pendingNodes *stack.Stack

func BuildSchema2(app *app.App, entryNodes []AbstractNode) {
	pendingNodes = stack.New()
	for _, entry := range entryNodes {
		pendingNodes.Push(entry)
	}
	for {
		node := pendingNodes.Pop().(*AbstractNode)
		if pendingNodes.Len() == 0 {
			break
		}
	}
} */

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
	//app.ResetAllDataflows()

	/* logger.Logger.Debug()
	for i, pendingNode := range pendingNodes {
		logger.Logger.Debugf("[%d] pending node %s ", i, pendingNode.dbCall.String())
	}
	logger.Logger.Info()
	for i, node := range visitedNodes {
		logger.Logger.Infof("[%d] visited node %s ", i, node.String())
	} */
}

/* type pendingNode struct {
	dbCall *AbstractDatabaseCall
	done   bool
}

var pendingNodes []*pendingNode
var pendingDBs map[*datastores.Datastore]bool = make(map[*datastores.Datastore]bool) */

var visitedNodes []AbstractNode
var writtenDatastores = make(map[string]bool, 0)

func doBuildSchema(app *app.App, node AbstractNode, requestIdx int) bool {
	if dbCall, ok := node.(*AbstractDatabaseCall); ok && dbCall.ParsedCall.Method.IsRead() {
		datastore := dbCall.DbInstance.GetDatastore()
		params := dbCall.GetParams()
		returns := dbCall.GetReturns()
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
		}

	} else if dbCall, ok := node.(*AbstractDatabaseCall); ok && dbCall.ParsedCall.Method.IsWrite() {
		datastore := dbCall.DbInstance.GetDatastore()
		if found := writtenDatastores[datastore.Name]; !found {
			writtenDatastores[datastore.Name] = true
		}
		params := dbCall.GetParams()
		logger.Logger.Infof("[SCHEMA] [%s] building schema based on abstract node (%s)", datastore.Name, dbCall.GetName())

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
		}
	}

	visitedNodes = append(visitedNodes, node)

	/* if dbCall, ok := node.(*AbstractDatabaseCall); ok && dbCall.ParsedCall.Method.IsUpdate() {
		logger.Logger.Fatalf("TODO FOR DB CALL: %s", dbCall.String())
	} */

	for _, child := range node.GetChildren() {
		doBuildSchema(app, child, requestIdx)
	}
	return true
}
