package tainter

import (
	"go/types"
	"slices"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/app"
	"analyzer/pkg/common"
	"analyzer/pkg/ssagraph"
	"analyzer/pkg/ssagraph/registry"
	"analyzer/pkg/utils"
)

const BLUEPRINT_BACKEND_PACKAGE = "github.com/blueprint-uservices/blueprint/runtime/core/backend"

var BLUEPRINT_BACKEND_CALLS_QUEUE = []string{"Push", "Pop"}
var BLUEPRINT_BACKEND_CALLS_NOSQLDATABASE = []string{"GetCollection"}

// TODO: UpdateMany (check foobar app)
var BLUEPRINT_BACKEND_CALLS_NOSQLCOLLECTION = []string{"InsertOne", "FindOne", "DeleteOne", "DeleteMany", "FindMany", "UpdateOne" /* , "UpdateMany" */, "Upsert", "ReplaceOne"}
var BLUEPRINT_BACKEND_CALLS_NOSQLCURSOR = []string{"One", "All"}
var BLUEPRINT_BACKEND_CALLS_RELATIONALDB = []string{"Exec", "Select", "Get"}
var BLUEPRINT_BACKEND_CALLS_CACHE = []string{"Get", "Put", "Mget"}

type ValFieldPath struct {
	val       ssa.Value
	fieldpath string

	// need to distinguish between filter keys and retrieved values
	// aplies to reads (any db) and updates (e.g., nosql)
	readKey   bool // aka filter keys
	readValue bool // aka returned values

	// cache only
	cacheMultiget bool

	// nosql only
	bsonFilterKey  string
	bsonCursorMany bool // destination objects from nosql cursor reads
	bsonFilterIn   bool // bson filter: $in
	bsonFilterEach bool // bson filter: $each
}

func isServiceCall(graph *ssagraph.SSAGraph, val ssa.Value) (string, string, string, []ssa.Value, *ssa.Call, bool) {
	if call, ok := val.(*ssa.Call); ok {
		logrus.Tracef("[CALLS BLUEPRINT] [SVC] checking for service call: %s\n", val.String())

		// Example:
		// t0 = &s.barService [#1]
		// t1 = *t0 								(*ssa.UnOp)
		// t2 = invoke t1.WriteBar(ctx, id, text) 	(*ssa.Call)
		if unOp, ok := call.Call.Value.(*ssa.UnOp); ok {
			if typesNamed, ok := unOp.Type().(*types.Named); ok {
				servicePath := typesNamed.String()
				if service := graph.GetApp().GetServiceWithPathIfExists(servicePath); service != nil {
					serviceName := service.GetName()
					method := call.Call.Method.Id()
					if service.HasMethod(method) {
						// NOTE: unOp.Type().String() does not contain "Impl" suffix here so GetShortFunctionPath will just ignore
						funcShortPath := utils.GetShortFunctionPath(unOp.Type().String() + "." + method)

						// return all args except context
						// NOTE: in this case (when call.Call.Value is UnOp) call.Call.Args does not contain the receiver
						return serviceName, method, funcShortPath, call.Call.Args[1:], call, true
					} else {
						logrus.Fatalf("[CALLS BLUEPRINT] [SVC] method (%s) not found for service (%s)", method, serviceName)
					}
				}
			}
		}
	}
	return "", "", "", nil, nil, false
}

func isMethodCall(instr ssa.Instruction, val ssa.Value) (string, string, []ssa.Value, []ssa.Value, *ssa.Call, *ssa.Function, bool) {
	if val != nil {
		if call, ok := val.(*ssa.Call); ok {
			if fn, ok := call.Call.Value.(*ssa.Function); ok {
				fnShortPath := utils.GetShortFunctionPath(fn.String())
				method := fn.Name()
				return method, fnShortPath, nil, call.Call.Args, call, fn, true
			}
			return "", "", nil, nil, nil, nil, false
		}
	} else {
		// go routine
		if goCall, ok := instr.(*ssa.Go); ok {
			if makeClosure, ok := goCall.Call.Value.(*ssa.MakeClosure); ok {
				if fn, ok := makeClosure.Fn.(*ssa.Function); ok {
					fnShortPath := utils.GetShortFunctionPath(fn.String())
					method := fn.Name()
					logrus.WithField("fn_short_path", fnShortPath).WithField("method", method).Warnf("[METHOD CALL] found call for go routine: %s", goCall.String())
					return method, fnShortPath, makeClosure.Bindings, goCall.Call.Args, nil, fn, true
				}
			}
		}
	}
	return "", "", nil, nil, nil, nil, false
}

func isDatabaseCall(graph *ssagraph.SSAGraph, val ssa.Value) (string, string, string, common.DatabaseOperationType, []ValFieldPath, bool) {
	if val == nil {
		return "", "", "", -1, nil, false
	}

	if call, ok := val.(*ssa.Call); ok {
		logrus.Tracef("[CALLS BLUEPRINT] [DB] checking for database call: %s\n", val.String())

		// --------------
		// blueprint apps
		// --------------
		if unOp, ok := call.Call.Value.(*ssa.UnOp); ok {
			if queue, topic, opType, valFieldPathLst, ok := isBlueprintQueueCall(graph, call, unOp); ok {
				// return all args except context
				// NOTE: in this case (when call.Call.Value is UnOp) call.Call.Args does not contain the receiver
				return queue, topic, call.Call.Method.Id(), opType, valFieldPathLst, true
			} else if cache, namespace, opType, valFieldPathLst, ok := isBlueprintCacheCall(graph, call, unOp); ok {
				return cache, namespace, call.Call.Method.Id(), opType, valFieldPathLst, true
			} else if database, collection, opType, valFieldPathLst, ok := isBlueprintRelationalDBCall(graph, call, unOp); ok {
				return database, collection, call.Call.Method.Id(), opType, valFieldPathLst, true
			} else if ok := isBlueprintNoSQLDatabaseCall(graph, call, unOp); ok {
				// call for NoSQLDatabase.GetCollection(...)
				// skip for now
				return "", "", "", -1, nil, false
			}
		}
		if extr, ok := call.Call.Value.(*ssa.Extract); ok {
			if database, collection, opType, valFieldPathLst, ok := isBlueprintNoSQLCollectionCall(graph, call, extr); ok {
				/* if opType == common.OP_READ {
					bsonFilter := call.Call.Args[]
				} */
				return database, collection, call.Call.Method.Id(), opType, valFieldPathLst, true
			}
		}
	}
	return "", "", "", -1, nil, false
}

// example:
// t3 = &m.movieIdDB [#0]
// t4 = *t3
// t5 = new [2]any (varargs)
// t6 = &t5[0:int]
// t7 = make any <- string (movieID)
// *t6 = t7
// t8 = &t5[1:int]
// t9 = make any <- string (title)
// *t8 = t9
// t10 = slice t5[:]
// t11 = invoke t4.Exec(ctx, "INSERT INTO movie...":string, t10...)
// t12 = extract t11 #0
// t13 = extract t11 #1
func isBlueprintRelationalDBCall(graph *ssagraph.SSAGraph, call *ssa.Call, unOp *ssa.UnOp) (string, string, common.DatabaseOperationType, []ValFieldPath, bool) {
	var opType common.DatabaseOperationType
	if typeNamed, ok := unOp.Type().(*types.Named); ok {
		if typeNamed.String() == BLUEPRINT_BACKEND_PACKAGE+".RelationalDB" {
			if !slices.Contains(BLUEPRINT_BACKEND_CALLS_RELATIONALDB, call.Call.Method.Name()) {
				return "", "", -1, nil, false
			}
			logrus.Tracef("[CALLS BLUEPRINT] [RELDB] found RelationalDB call: %v\n", call)

			var dstVal, stmtVal, sliceArgsVal ssa.Value
			if call.Call.Method.Name() == "Select" || call.Call.Method.Name() == "Get" {
				// e.g., Select(ctx, &movieId, "SELECT * FROM movieid WHERE movieid = ?", movieID)
				// TODO: add support for more than 1 dst val
				dstVal = call.Call.Args[1]
				stmtVal = call.Call.Args[2]
				sliceArgsVal = call.Call.Args[3]
			} else if call.Call.Method.Name() == "Exec" {
				// e.g., Exec(ctx, "INSERT INTO movieid(movieid, title) VALUES (?, ?);", movieID, title)
				// DELETE FROM sock_tag WHERE sock_tag.sock_id=?;
				stmtVal = call.Call.Args[1]
				sliceArgsVal = call.Call.Args[2]
			}

			stmt, ok := utils.ExtractStringFromValue(stmtVal)
			if !ok {
				return "", "", -1, nil, false
			}

			switch call.Call.Method.Name() {
			case "Select", "Get":
				opType = common.OP_READ
			case "Exec": // can also be update or read
				if strings.HasPrefix(stmt, "INSERT") {
					opType = common.OP_WRITE
				} else if strings.HasPrefix(stmt, "DELETE") {
					opType = common.OP_DELETE
				} else {
					logrus.WithField("graph", graph.String()).Fatalf("[CALLS BLUEPRINT] [RELDB] unsupported SQL statement: %s\n", stmt)
				}
			default:
				logrus.WithField("graph", graph.String()).Fatalf("[CALLS BLUEPRINT] [RELDB] unknown method name for queue call: %s\n", call.String())
			}

			database, ok := extractDatabaseNameFromUnOp(graph, unOp)
			if !ok {
				return "", "", -1, nil, false
			}

			var argVals []ssa.Value
			if slice, ok := sliceArgsVal.(*ssa.Slice); ok {
				logrus.Tracef("[CALLS BLUEPRINT] [RELDB] on ssa slice: %v\n", slice)
				if alloc, ok := slice.X.(*ssa.Alloc); ok {
					// example:
					// t5 = new [2]any (varargs)
					// t6 = &t5[0:int]
					// t7 = make any <- string (movieID)
					// *t6 = t7
					// t10 = slice t5[:]
					//
					// TODO: also need to consider case when slice is declared earlier and
					// reused for another write in possbly another DB
					allocNode := graph.GetNodeByName(alloc.Name())
					logrus.Tracef("[CALLS BLUEPRINT] [RELDB] on alloc node: %v\n", allocNode)
					for _, edge := range graph.GetEdgesFromNode(allocNode) {
						if edge.GetType() == ssagraph.EDGE_INDEX {
							idxNode := edge.GetToNode()
							logrus.Tracef("[CALLS BLUEPRINT] [RELDB] on idx node: %v\n", idxNode)
							for _, edge := range graph.GetEdgesFromNode(idxNode) {
								if edge.GetType() == ssagraph.EDGE_STORE_ADDRESS {
									storeNode := edge.GetToNode()
									logrus.Tracef("[CALLS BLUEPRINT] [RELDB] on store node: %v\n", storeNode)
									storeInstr, _ := storeNode.GetInstruction().(*ssa.Store)
									argVals = append(argVals, storeInstr.Val)

									if storeInstr.Val == nil {
										logrus.Fatalf("unexpected nil val for storeinstr: %v\n", storeInstr)
									}
								}
							}
						}
					}
				}
			}

			var fields []string
			var tableName string
			var filterFields []string

			if opType == common.OP_READ {
				var tables []string
				logrus.Tracef("[CALLS BLUEPRINT] [SQL] [READ] parsing stmt: %s\n", stmt)
				filterFields, fields, tables, ok = app.ParseSQLRead(database, stmt)
				if !ok {
					return "", "", -1, nil, false
				}
				logrus.Tracef("[CALLS BLUEPRINT] [SQL] [READ] got filter fields: %v\n", filterFields)
				tableName = tables[0]

				// sanity check
				if len(argVals) != len(fields) {
					logrus.Fatalf("[CALLS BLUEPRINT] [RELDB] length of arg vals (%d) does not match length fields (%d)\n", len(argVals), len(fields))
				}
			} else if opType == common.OP_WRITE {
				logrus.Tracef("[CALLS BLUEPRINT] [SQL] [WRITE] parsing stmt: %s\n", stmt)
				var ok bool
				fields, _, tableName, ok = app.ParseSQLWrite(database, stmt)
				if !ok {
					return "", "", -1, nil, false
				}
				logrus.Tracef("[CALLS BLUEPRINT] [SQL] [WRITE] got written fields: %v\n", fields)
			} else if opType == common.OP_DELETE {
				logrus.Tracef("[CALLS BLUEPRINT] [SQL] [DELETE] parsing stmt: %s\n", stmt)
				var tables []string
				var ok bool
				filterFields, tables, ok = app.ParseSQLDelete(database, stmt)
				if !ok {
					return "", "", -1, nil, false
				}
				logrus.Tracef("[CALLS BLUEPRINT] [SQL] [DELETE] got filter fields: %v\n", filterFields)
				tableName = tables[0]

				// sanity check
				if len(argVals) != len(filterFields) {
					logrus.Fatalf("[CALLS BLUEPRINT] [RELDB] length of arg vals (%d) does not match length fields (%d)\n", len(argVals), len(fields))
				}
			}

			var valFieldPathLst []ValFieldPath
			for i, field := range fields {
				argVal := argVals[i]
				if argVals[i] == nil {
					logrus.Fatalf("field argvals[i] is nil")
				}
				valFieldPathLst = append(valFieldPathLst, ValFieldPath{
					val:       argVal,
					fieldpath: field,
					readKey:   true,
				})
			}

			if opType == common.OP_READ {
				// for SQL Selects on all fields (i.e., '*') the readFields length is 1
				// and the readField has format <database>.<table>
				if call.Call.Method.Name() == "Select" && len(filterFields) > 0 {
					// select method reads entire row
					filterField := filterFields[0]
					valFieldPathLst = append(valFieldPathLst, ValFieldPath{
						val:       dstVal,
						fieldpath: filterField,
						readValue: true,
					})

					if dstVal == nil {
						logrus.Fatalf("dstval is nil")
					}
				}
			} else if opType == common.OP_DELETE {
				// for SQL Selects on all fields (i.e., '*') the readFields length is 1
				// and the readField has format <database>.<table>
				for i, filterField := range filterFields {
					valFieldPathLst = append(valFieldPathLst, ValFieldPath{
						val:       argVals[i],
						fieldpath: filterField,
						readKey:   true,
					})

					if argVals[i] == nil {
						logrus.Fatalf("is nil")
					}
				}
			}

			return database, tableName, opType, valFieldPathLst, true
		}
	}
	return "", "", -1, nil, false
}

// example (t11 is UnOp and t10 is UnOp.X):
// t10 = &u.notificationsQueue [#1]
// t11 = *t10
// t14 = invoke t11.Push(ctx, t13)
func isBlueprintQueueCall(graph *ssagraph.SSAGraph, call *ssa.Call, unOp *ssa.UnOp) (string, string, common.DatabaseOperationType, []ValFieldPath, bool) {
	var opType common.DatabaseOperationType
	if typeNamed, ok := unOp.Type().(*types.Named); ok {
		if typeNamed.String() == BLUEPRINT_BACKEND_PACKAGE+".Queue" {
			if !slices.Contains(BLUEPRINT_BACKEND_CALLS_QUEUE, call.Call.Method.Name()) {
				return "", "", -1, nil, false
			}

			switch call.Call.Method.Name() {
			case "Pop":
				opType = common.OP_READ
			case "Push":
				opType = common.OP_WRITE
			default:
				logrus.Fatalf("[CALLS BLUEPRINT] [QUEUE] unknown method name for queue call: %s\n", call.String())
			}

			// e.g., t10 = &u.notificationsQueue [#1]
			if database, ok := extractDatabaseNameFromUnOp(graph, unOp); ok {
				topic := "notification"
				valFieldPathLst := make([]ValFieldPath, 1)
				docVal := call.Call.Args[1]
				valFieldPathLst[0] = ValFieldPath{
					val:       docVal,
					fieldpath: database + "." + topic,
					readValue: opType == common.OP_READ,
				}

				return database, topic, opType, valFieldPathLst, true
			}
		}

	}
	return "", "", -1, nil, false
}

func isBlueprintCacheCall(graph *ssagraph.SSAGraph, call *ssa.Call, unOp *ssa.UnOp) (string, string, common.DatabaseOperationType, []ValFieldPath, bool) {
	var opType common.DatabaseOperationType
	if typeNamed, ok := unOp.Type().(*types.Named); ok {
		if typeNamed.String() == BLUEPRINT_BACKEND_PACKAGE+".Cache" {
			if !slices.Contains(BLUEPRINT_BACKEND_CALLS_CACHE, call.Call.Method.Name()) {
				return "", "", -1, nil, false
			}

			switch call.Call.Method.Name() {
			case "Get":
				opType = common.OP_READ
			case "Mget":
				opType = common.OP_READ_MANY
			case "Put":
				opType = common.OP_WRITE
			default:
				logrus.Fatalf("[CALLS BLUEPRINT] [CACHE] unknown method name for queue call: %s\n", call.String())
			}

			// e.g., t10 = &u.notificationsQueue [#1]
			if database, ok := extractDatabaseNameFromUnOp(graph, unOp); ok {
				namespace := "*"
				var valFieldPathLst []ValFieldPath
				cacheKeyVal := call.Call.Args[1]
				cacheValueVal := call.Call.Args[2]
				logrus.Tracef("[CALLS BLUEPRINT] [CACHE] cache key [%T]: %v\n", cacheKeyVal, cacheKeyVal)
				logrus.Tracef("[CALLS BLUEPRINT] [CACHE] cache value [%T]: %v\n", cacheValueVal, cacheValueVal)

				var keyField = "Key"
				var valField = "Value"

				// track cache key
				if _, ok := utils.ExtractStringFromValue(cacheKeyVal); ok {
					valFieldPathLst = append(valFieldPathLst, ValFieldPath{
						val:           cacheKeyVal,
						fieldpath:     database + "." + namespace + "." + keyField,
						cacheMultiget: opType == common.OP_READ_MANY,
						readKey:       true,
					})
				} else if binOp, ok := cacheKeyVal.(*ssa.BinOp); ok {
					if suffix, ok := utils.ExtractStringFromValue(binOp.Y); ok {
						namespace, _ = strings.CutPrefix(suffix, ":")
						// real cache key
						valFieldPathLst = append(valFieldPathLst, ValFieldPath{
							val:           binOp.X,
							fieldpath:     database + "." + namespace + "." + keyField,
							cacheMultiget: opType == common.OP_READ_MANY,
							readKey:       true,
						})
					}
				} else if call, ok := cacheKeyVal.(*ssa.Call); ok {
					for _, arg := range call.Call.Args {
						valFieldPathLst = append(valFieldPathLst, ValFieldPath{
							val:           arg,
							fieldpath:     database + "." + namespace + "." + keyField,
							cacheMultiget: opType == common.OP_READ_MANY,
							readKey:       true,
						})
					}
				} else {
					logrus.WithField("graph", graph.String()).
						Warnf("[CALLS BLUEPRINT] [CACHE] unknown cache key (%s)", cacheKeyVal.String())
					valFieldPathLst = append(valFieldPathLst, ValFieldPath{
						val:           cacheKeyVal,
						fieldpath:     database + "." + namespace + "." + keyField,
						cacheMultiget: opType == common.OP_READ_MANY,
						readKey:       true,
					})
				}

				if valFieldPathLst == nil {
					// [TO BE IMPROVED]
					valFieldPathLst = append(valFieldPathLst, ValFieldPath{
						val:           cacheKeyVal,
						fieldpath:     database + "." + namespace + "." + keyField,
						cacheMultiget: opType == common.OP_READ_MANY,
						readKey:       true,
					})
					logrus.Tracef("[CALLS CACHE] [%s] could not save any cache key for call: %v\n", graph.String(), call)
				}

				// track cache value
				valFieldPathLst = append(valFieldPathLst, ValFieldPath{
					val:           cacheValueVal,
					fieldpath:     database + "." + namespace + "." + valField,
					cacheMultiget: opType == common.OP_READ_MANY,
					readValue:     true,
				})

				return database, namespace, opType, valFieldPathLst, true
			}
		}

	}
	return "", "", -1, nil, false
}

func isBlueprintNoSQLDatabaseCall(graph *ssagraph.SSAGraph, call *ssa.Call, unOp *ssa.UnOp) bool {
	if typeNamed, ok := unOp.Type().(*types.Named); ok {
		if typeNamed.String() == BLUEPRINT_BACKEND_PACKAGE+".NoSQLDatabase" {
			if !slices.Contains(BLUEPRINT_BACKEND_CALLS_NOSQLDATABASE, call.Call.Method.Name()) {
				return false
			}
			// call for NoSQLDatabase.GetCollection(...)
			// skip for now
			return true
		}
	}
	return false
}

// e.g.
// cursor, err := collection.FindOne(ctx, query)
func isBlueprintNoSQLCursorCall(graph *ssagraph.SSAGraph, call *ssa.Call, extr *ssa.Extract) (ssa.Value, bool) {
	if typeNamed, ok := extr.Type().(*types.Named); ok {
		if typeNamed.String() == BLUEPRINT_BACKEND_PACKAGE+".NoSQLCursor" {
			if !slices.Contains(BLUEPRINT_BACKEND_CALLS_NOSQLCURSOR, call.Call.Method.Name()) {
				return nil, false
			}
			dst := call.Call.Args[1]

			// just to be more precise
			// e.g., FooBar: FooService.ReadFoo()
			// t0: new Foo
			// t17: make interface{} <- *Foo (t0)
			// t18: invoke t14.One(ctx, t17)
			if iface, ok := dst.(*ssa.MakeInterface); ok {
				return iface.X, true
			}

			// return original dst value
			return dst, true
		}
	}
	return nil, false
}

// TODO: get database name (not the db name of mongodb!)
func isBlueprintNoSQLCollectionCall(graph *ssagraph.SSAGraph, call *ssa.Call, extr *ssa.Extract) (string, string, common.DatabaseOperationType, []ValFieldPath, bool) {
	var opType common.DatabaseOperationType
	if typeNamed, ok := extr.Type().(*types.Named); ok {
		if typeNamed.String() == BLUEPRINT_BACKEND_PACKAGE+".NoSQLCollection" {
			if !slices.Contains(BLUEPRINT_BACKEND_CALLS_NOSQLCOLLECTION, call.Call.Method.Name()) {
				return "", "", -1, nil, false
			}

			switch call.Call.Method.Name() {
			case "FindOne":
				opType = common.OP_READ
			case "FindMany":
				opType = common.OP_READ_MANY
			case "InsertOne":
				opType = common.OP_WRITE
			case "UpdateOne", "UpdateMany", "ReplaceOne", "Upsert":
				opType = common.OP_UPDATE
			case "DeleteOne", "DeleteMany":
				opType = common.OP_DELETE
			default:
				logrus.Fatalf("[CALLS BLUEPRINT] [NOSQL] unknown method name for queue call: %s\n", call.String())
			}

			// e.g.
			// t29 = invoke t28.GetCollection(ctx, "posts_db":string, "post":string)
			// t30 = extract t29 #0
			// t34 = invoke t30.InsertOne(ctx, t33)
			if ssaExtract, ok := call.Call.Value.(*ssa.Extract); ok {
				if ssaCall, ok := ssaExtract.Tuple.(*ssa.Call); ok {
					databaseVal := ssaCall.Call.Args[1]
					collectionVal := ssaCall.Call.Args[2]
					database, _ := utils.ExtractStringFromValue(databaseVal)
					collection, _ := utils.ExtractStringFromValue(collectionVal)

					// sanity check
					// keep this while database logic is not complete
					if !graph.GetApp().HasDatabase(database) {
						logrus.Fatalf("[CALLS BLUEPRINT] [NOSQL] database (%s) extracted from value (%s) not found for app with databases: %v", database, databaseVal.String(), graph.GetApp().GetAllDatabases())
					}

					var valFieldPathLst []ValFieldPath
					if opType == common.OP_WRITE {
						docVal := call.Call.Args[1]

						registry.RegisterNoSQLPrimaryKey(graph.GetApp(), database, collection, docVal)

						valFieldPathLst = append(valFieldPathLst, ValFieldPath{
							val:       docVal,
							fieldpath: database + "." + collection,
						})
					} else { // reads, updates, or deletes
						filterVal := call.Call.Args[1]
						bsonNode := findBsonNode(graph, filterVal)
						if bsonNode != nil {
							filterKeyToValues := computeNoSQLFilterKeyToValues(graph, bsonNode, nil, false)
							for filter, vals := range filterKeyToValues {
								for _, val := range vals {
									// sanity check
									if filter != "" {
										val.fieldpath = database + "." + collection + "." + filter
									} else {
										val.fieldpath = database + "." + collection + ".*"
									}
									val.readKey = true
									valFieldPathLst = append(valFieldPathLst, val)
								}
							}
						}
					}

					var projections []string
					// NoSQL operation uses projection
					if opType == common.OP_READ || opType == common.OP_READ_MANY {
						if len(call.Call.Args) > 2 {
							projection := call.Call.Args[2]
							bsonNode := findBsonNode(graph, projection)
							if bsonNode != nil {
								filterKeyToValues := computeNoSQLFilterKeyToValues(graph, bsonNode, nil, true)
								for projectionValue := range filterKeyToValues {
									// sanity check
									if projectionValue != "" {
										projections = append(projections, projectionValue)
									}
								}
							}
						}
					}

					if opType == common.OP_UPDATE {
						if call.Call.Method.Name() == "UpdateOne" {
							//FIXME: change to call.Call.Args[2] (test with SockShop app)
							updateVal := call.Call.Args[1]
							bsonNode := findBsonNode(graph, updateVal)
							if bsonNode != nil {
								filterKeyToValues := computeNoSQLFilterKeyToValues(graph, bsonNode, nil, false)
								for filter, vals := range filterKeyToValues {
									for _, val := range vals {
										// sanity check
										if filter != "" {
											val.fieldpath = database + "." + collection + "." + filter
										} else {
											val.fieldpath = database + "." + collection + ".*"
										}
										valFieldPathLst = append(valFieldPathLst, val)
									}
								}
							}
						} else if call.Call.Method.Name() == "UpdateMany" {
							updateVal := call.Call.Args[2]
							bsonNode := findBsonNode(graph, updateVal)
							if bsonNode != nil {
								filterKeyToValues := computeNoSQLFilterKeyToValues(graph, bsonNode, nil, false)
								for filter, vals := range filterKeyToValues {
									for _, val := range vals {
										// sanity check
										if filter != "" {
											val.fieldpath = database + "." + collection + "." + filter
										} else {
											val.fieldpath = database + "." + collection + ".*"
										}
										valFieldPathLst = append(valFieldPathLst, val)
									}
								}
							}
						} else if call.Call.Method.Name() == "ReplaceOne" || call.Call.Method.Name() == "Upsert" {
							docVal := call.Call.Args[2]
							valFieldPathLst = append(valFieldPathLst, ValFieldPath{
								val:       docVal,
								fieldpath: database + "." + collection,
							})
						}
					}

					// propagata taint to future reads on cursor
					if opType == common.OP_READ || opType == common.OP_READ_MANY {
						for cursorCall, extr := range getNoSQLCursorCallsFromCollectionCall(graph, call) {
							if dstVal, ok := isBlueprintNoSQLCursorCall(graph, cursorCall, extr); ok {
								fieldpath := database + "." + collection
								if opType == common.OP_READ_MANY && len(projections) > 0 {
									// NoSQL FindMany returns many documents for the projection
									// so we need to add the projectionValue directly to the path
									if len(projections) == 1 {
										fieldpath += "." + projections[0]
									} else {
										// TODO: add support for multiple projection values
										logrus.Fatalf("TODO! projections = %v\n", projections)
									}
								} else {
									// skip
									// NosQL FindOne returns the document that includes the projection
									// so the fieldpath does not need to contain the projectionValue now
								}
								// distinguish objects used as arguments for any operation
								// from objects used as destination for reads (isDst=true)
								valFieldPathLst = append(valFieldPathLst, ValFieldPath{
									fieldpath:      fieldpath,
									val:            dstVal,
									bsonCursorMany: opType == common.OP_READ_MANY,
									readValue:      true,
								})
							}
						}
					}

					return database, collection, opType, valFieldPathLst, true
				}
			}
		}
	}
	return "", "", -1, nil, false
}

func findSliceForBsonAorD(graph *ssagraph.SSAGraph, bsonAVal ssa.Value) ssa.Value {
	// BSON A example (t28 -> t10):
	// t28: make interface{} <- go.mongodb.org/mongo-driver/bson/primitive.A (t27)
	// t27: slice t10[:]
	// t10: new [2]interface{} (slicelit)

	// BSON D example (t26 -> t20):
	// t26: make interface{} <- go.mongodb.org/mongo-driver/bson/primitive.D (t25)
	// t25: slice t20[:]
	// t20: new [1]go.mongodb.org/mongo-driver/bson/primitive.E{} (slicelit)

	// sanity check
	if _, ok := bsonAVal.(*ssa.MakeInterface); !ok {
		return nil
	}

	toNode := graph.GetNodeByName(bsonAVal.Name())
	fromNode := graph.GetFirstEdgeToNode(toNode).GetFromNode()
	if _, ok := fromNode.GetValue().(*ssa.Slice); !ok {
		logrus.Fatalf("unexpected type [%T]: %s\n", fromNode.GetValue(), fromNode.GetValue().String())
	}
	fromNode2 := graph.GetFirstEdgeToNode(fromNode).GetFromNode()
	if _, ok := fromNode2.GetValue().(*ssa.Alloc); !ok {
		logrus.Fatalf("unexpected type [%T]: %s\n", fromNode2.GetValue(), fromNode2.GetValue().String())
	}

	return fromNode2.GetValue()
}

func getValInStoreInstrForCurrentAddr(graph *ssagraph.SSAGraph, bsonValueNode *ssagraph.SSANode) ssa.Value {
	edge := graph.GetFirstEdgeTypedFrom(bsonValueNode, ssagraph.EDGE_STORE_ADDRESS)
	if edge == nil {
		logrus.Fatalf("unexpected nil edge")
	}

	storeInstr := edge.GetToNode().GetInstruction().(*ssa.Store)
	return storeInstr.Val
}

// bson filter: $and
//
// Example in Go:
//
//	query := bson.D{{Key: "$and", Value: bson.A{
//			bson.D{{Key: "RouteID", Value: routeID}},
//			bson.D{{Key: "TrainType", Value: trainType}},
//		}}}
//
// Example in SSA:
// t6: new [1]go.mongodb.org/mongo-driver/bson/primitive.E (slicelit)
// t7: &t5[0:int]
// t8: &t7.Key [#0]
// *t8 = "$and" :string
// t9: &t7.Value [#1]
// *t9 = t28
// t28: make interface{} <- go.mongodb.org/mongo-driver/bson/primitive.A (t27)
// t27: slice t10[:]
// t10: new [2]interface{} (slicelit)
//
// OBSERVATIONS:
// t9 is our bsonValueNode
// t28 is the bson.A in Go code
// t10 is the real slice
// we need to go from t9 to t10
func computeNoSQLFilterKeyToValues_AND(graph *ssagraph.SSAGraph, bsonVal ssa.Value, bsonElemNode *ssagraph.SSANode) []*ssagraph.SSANode {
	var elemSliceValsNodes []*ssagraph.SSANode

	for _, edge := range graph.GetEdgesTypedFrom(bsonElemNode, ssagraph.EDGE_FIELD) {
		if edge.GetParam() == "Value" {
			// track objects used as value in store instructions for current bson value
			bsonValueNode := edge.GetToNode()

			storeInstrVal := getValInStoreInstrForCurrentAddr(graph, bsonValueNode)
			sliceVal := findSliceForBsonAorD(graph, storeInstrVal)
			if sliceVal == nil {
				logrus.Fatalf("unexpected nil sliceVal")
			}

			// restrict to ssagraph.EDGE_FIELD just for sanity check
			// we are not expecting any other type of edges
			for _, edge := range graph.GetEdgesTypedFrom(graph.GetNodeByName(sliceVal.Name()), ssagraph.EDGE_INDEX) {
				elemNode := edge.GetToNode()
				elemStoreInstrVal := getValInStoreInstrForCurrentAddr(graph, elemNode)
				elemSliceVal := findSliceForBsonAorD(graph, elemStoreInstrVal)
				if elemSliceVal == nil {
					logrus.Fatalf("unexpected nil sliceVal")
				}
				elemSliceValsNodes = append(elemSliceValsNodes, graph.GetNodeByName(elemSliceVal.Name()))
			}
		}
	}

	return elemSliceValsNodes
}

// helper to process bson slice node
// if isProjection=true then we only want the key and not the value (which should be set to "true")
// REMINDER: we assume that all keys present have value set to true, otherwise we need
// more code logic to know which key matches which value
func computeNoSQLFilterKeyToValues(graph *ssagraph.SSAGraph, bsonArrayNode *ssagraph.SSANode, filterKeyToValues map[string][]ValFieldPath, isProjection bool) map[string][]ValFieldPath {
	if filterKeyToValues == nil {
		filterKeyToValues = make(map[string][]ValFieldPath)
	}
	var edges []*ssagraph.SSAEdge
	for _, edge := range graph.GetEdgesTypedFrom(bsonArrayNode, ssagraph.EDGE_INDEX) {
		bsonArrayElemNode := edge.GetToNode()
		edges = append(edges, graph.GetEdgesTypedFrom(bsonArrayElemNode, ssagraph.EDGE_FIELD)...)
		if edge.GetIndex() > 0 {
			logrus.WithField("index", edge.GetIndex()).Warnf("[CALLS BLUEPRINT] check bson edge\n")
		}
	}
	for _, edge := range edges {
		if edge.GetParam() != "Key" {
			continue
		}

		var filterField string
		bsonArrayElemNode := edge.GetFromNode()

		bsonArrayElemNode.EnableUsedInBson()

		// objects that are excluded from taint:
		// - bson slice
		// - bson slice elems
		// - bson slice elem key
		// - bson slice elem value

		var filterObjs = []ValFieldPath{}

		// track objects used as value in store instructions for current bson key
		bsonArrayElemKeyNode := edge.GetToNode()
		bsonArrayElemKeyNode.EnableUsedInBson()

		for _, edge := range graph.GetEdgesTypedFrom(bsonArrayElemKeyNode, ssagraph.EDGE_STORE_ADDRESS) {
			storeInstr := edge.GetToNode().GetInstruction().(*ssa.Store)
			filterObjs = append(filterObjs, ValFieldPath{
				val:           storeInstr.Val,
				bsonFilterKey: filterField,
			})
			// should only occur once
			if filterFieldTmp, ok := utils.ExtractStringFromValue(storeInstr.Val); ok {
				filterField = filterFieldTmp
				break
			}
		}
		if filterField == "$and" {
			bsonArrayElemSliceNodes := computeNoSQLFilterKeyToValues_AND(graph, bsonArrayNode.GetValue(), bsonArrayElemNode)
			// NOTE: the final appended filterKeyToValues will not contain the
			// filterObj for the "Key", which is good because we don't want taints with "$and"
			for _, node := range bsonArrayElemSliceNodes {
				filterKeyToValuesTmp := computeNoSQLFilterKeyToValues(graph, node, filterKeyToValues, isProjection)
				for k, lst := range filterKeyToValuesTmp {
					for _, v := range lst {
						if !slices.Contains(filterKeyToValues[k], v) {
							filterKeyToValues[k] = append(filterKeyToValues[k], v)
						}
					}
				}
			}
			return filterKeyToValues
		}

		if filterField == "" {
			filterField = "*"
			logrus.WithField("graph", graph.String).WithField("bson_array_node", bsonArrayNode.String()).Fatalf("empty filter field")
		}

		if !isProjection {
			for _, edge := range graph.GetEdgesTypedFrom(bsonArrayElemNode, ssagraph.EDGE_FIELD) {
				if edge.GetParam() == "Value" {
					// track objects used as value in store instructions for current bson value
					bsonArrayElemValNode := edge.GetToNode()
					bsonArrayElemValNode.EnableUsedInBson()

					for _, edge := range graph.GetEdgesTypedFrom(bsonArrayElemValNode, ssagraph.EDGE_STORE_ADDRESS) {
						storeInstr := edge.GetToNode().GetInstruction().(*ssa.Store)
						filterObjs = append(filterObjs, ValFieldPath{
							val:           storeInstr.Val,
							bsonFilterKey: filterField,
						})
						if iface, ok := storeInstr.Val.(*ssa.MakeInterface); ok {
							bsonNode := findBsonNode(graph, iface.X)
							if bsonNode != nil {
								filterKeyToValues := computeNoSQLFilterKeyToValues(graph, bsonNode, nil, false)
								for filterFieldTmp, filterObjsTmp := range filterKeyToValues {
									switch filterFieldTmp {
									case "$in":
										for _, filterObjTmp := range filterObjsTmp {
											filterObjTmp.bsonFilterIn = true
											filterObjTmp.bsonFilterKey = filterField
											filterObjs = append(filterObjs, filterObjTmp)
										}
										//logrus.Fatalf("[DEBUG] [BSON FILTER $in] FILTER OBJ TEMPS = %v\n", filterObjsTmp)
									case "$gt":
										for _, filterObjTmp := range filterObjsTmp {
											filterObjTmp.bsonFilterIn = true
											filterObjTmp.bsonFilterKey = filterField
											filterObjs = append(filterObjs, filterObjTmp)
										}
										//logrus.Fatalf("[DEBUG] [BSON FILTER $in] FILTER OBJ TEMPS = %v\n", filterObjsTmp)
									case "$each":
										// TODO
										logrus.Warnf("[CALLS BLUEPRINT] [BSON] unexpected filter key (%s) for objects: %v", filterFieldTmp, filterObjsTmp)
									case "$position":
										// TODO
										logrus.Warnf("[CALLS BLUEPRINT] [BSON] unexpected filter key (%s) for objects: %v", filterFieldTmp, filterObjsTmp)
									default:
										logrus.Fatalf("[CALLS BLUEPRINT] [BSON] unexpected filter key (%s) for objects: %v", filterFieldTmp, filterObjsTmp)
									}
								}
							}
						}
					}
				}
			}
		}
		if filterField == "" {
			logrus.Fatalf("[CALLS BLUEPRINT] [BSON] empty filter field for bsonVal (%s) and bsonElem (%s)\n", bsonArrayNode.GetValue().Name(), bsonArrayElemNode.GetValue().Name())
		}
		filterKeyToValues[filterField] = filterObjs
	}
	return filterKeyToValues
}

// example of bson filters:
//
//	(1) query_d := bson.D{{Key: "PostID", Value: bson.D{
//			{Key: "$in", Value: unique_pids},
//			}}}
//
//	(2) query := bson.D{{Key: "$and", Value: bson.A{
//			bson.D{{Key: "RouteID", Value: routeID}},
//			bson.D{{Key: "TrainType", Value: trainType}},
//		}}}
func findBsonNode(graph *ssagraph.SSAGraph, bsonSliceVal ssa.Value) *ssagraph.SSANode {
	if slice, ok := bsonSliceVal.(*ssa.Slice); ok {
		ok, isSliceOfSlice := ssaValueIsUsedInMongoBsonFilter(graph, slice)
		if ok {
			bsonSliceNode := graph.GetNodeByName(slice.X.Name())
			if isSliceOfSlice {
				// can be a slice of bson.D
				// e.g., projections where parameter is "projection ...bson.D"
				for _, edge := range graph.GetEdgesFromNode(bsonSliceNode) {
					if _, ok := edge.GetToNode().GetValue().(*ssa.IndexAddr); ok {
						for _, edge := range graph.GetEdgesFromNode(edge.GetToNode()) {
							if store, ok := edge.GetToNode().GetInstruction().(*ssa.Store); ok {
								if slice, ok := store.Val.(*ssa.Slice); ok {
									bsonSliceNode.EnableUsedInBson()
									bsonNode := graph.GetNodeByName(slice.X.Name())
									bsonNode.EnableUsedInBson()
									return bsonNode
								}
							}
						}
					}
				}
			} else {
				bsonSliceNode.EnableUsedInBson()
				bsonNode := graph.GetNodeByName(slice.X.Name())
				bsonNode.EnableUsedInBson()
				return bsonNode
			}
		} else {
			logrus.WithField("graph", graph.String()).WithField("slice", slice.Name()).Warnf("slice is not used in mongo bson filter")
		}
	}
	//logrus.WithField("graph", graph.String()).Warnf("nil bson node for val: %v\n", bsonSliceVal)
	return nil
}

func getNoSQLCursorCallsFromCollectionCall(graph *ssagraph.SSAGraph, collectionCall *ssa.Call) map[*ssa.Call]*ssa.Extract {
	var cursorCalls = make(map[*ssa.Call]*ssa.Extract)
	callNode := graph.GetNodeByName(collectionCall.Name())
	for _, edge := range graph.GetEdgesFromNode(callNode) {
		if edge.IsType(ssagraph.EDGE_EXTRACT) && edge.GetIndex() == 0 {
			// e.g., FooBar: FooService.ReadFoo()
			// cursor, err := collection.FindOne(ctx, query)
			cursorNode := edge.GetToNode()
			for _, edge := range graph.GetEdgesFromNode(cursorNode) {
				if edge.IsType(ssagraph.EDGE_RECEIVER_ON_CALL) {
					cursorCall := edge.GetToNode().GetValue().(*ssa.Call)

					// checking if Call.Value is Extract is the standard for every NoSQL call
					if extr, ok := cursorCall.Call.Value.(*ssa.Extract); ok {
						cursorCalls[cursorCall] = extr
					}
				}
			}
			break
		}
	}
	return cursorCalls
}

func ssaValueIsMongoBsonFilter(val ssa.Value) bool {
	if val == nil {
		return false
	}
	if alloc, ok := val.(*ssa.Alloc); ok {
		// e.g.,
		// [ssa.Alloc] t60: new [3]go.mongodb.org/mongo-driver/bson/primitive.E (slicelit)
		if ptr, ok := alloc.Type().(*types.Pointer); ok {
			if array, ok := ptr.Elem().Underlying().(*types.Array); ok {
				if named, ok := array.Elem().(*types.Named); ok {
					// e.g.,
					// this is the real type: go.mongodb.org/mongo-driver/bson/primitive.E
					if named.Obj().Pkg() == nil {
						return false
					}
					if named.Obj().Pkg().Path() == "go.mongodb.org/mongo-driver/bson/primitive" && named.Obj().Id() == "E" {
						return true
					}
				}
			}
		}
	}

	return false
}

func ssaValueIsUsedInMongoBsonFilter(graph *ssagraph.SSAGraph, val ssa.Value) (bool, bool) {
	if val == nil {
		return false, false
	}
	if ok := ssaValueIsMongoBsonFilter(val); ok {
		return true, false
	}
	if slice, ok := val.(*ssa.Slice); ok {
		// e.g.,
		// [ssa.Slice] t73: slice t60[:]
		if alias, ok := slice.Type().(*types.Alias); ok {
			if aliasSlice, ok := alias.Underlying().(*types.Slice); ok {
				if named, ok := aliasSlice.Elem().(*types.Named); ok {
					// this is the alias type: go.mongodb.org/mongo-driver/bson/primitive.E
					if named.Obj().Pkg().Path() == "go.mongodb.org/mongo-driver/bson/primitive" && named.Obj().Id() == "E" {
						return true, false
					}
				}
			}
		} else if slice2, ok := slice.Type().(*types.Slice); ok {
			if alias, ok := slice2.Elem().(*types.Alias); ok {
				if slice3, ok := alias.Underlying().(*types.Slice); ok {
					if named, ok := slice3.Elem().(*types.Named); ok {
						// can be a slice of bson.D
						// e.g., projections where parameter is "projection ...bson.D"
						// this is the alias type: go.mongodb.org/mongo-driver/bson/primitive.E
						if named.Obj().Pkg().Path() == "go.mongodb.org/mongo-driver/bson/primitive" && named.Obj().Id() == "E" {
							return true, true
						}
					}
				}
			}
		}
	}
	return false, false
}

// extractDatabaseNameFromUnOp can be used for RelationalDB and Queue calls
// it cannot be used for NoSQLDatabase calls because the collection is extracted beforehand
func extractDatabaseNameFromUnOp(graph *ssagraph.SSAGraph, unOp *ssa.UnOp) (string, bool) {
	if ssaFieldAddr, ok := unOp.X.(*ssa.FieldAddr); ok {
		logrus.Tracef("[CALLS BLUEPRINT] [QUEUE] ssa field addr (field=%d): %s\n", ssaFieldAddr.Field, ssaFieldAddr.String())
		if ssaParam, ok := ssaFieldAddr.X.(*ssa.Parameter); ok {
			logrus.Tracef("[CALLS BLUEPRINT] [QUEUE] queue loaded from parameter (%d)\n", ssaFieldAddr.Field)
			if typesPointer, ok := ssaParam.Type().(*types.Pointer); ok {
				if typeNamed, ok := typesPointer.Elem().(*types.Named); ok {
					// e.g., github.com/blueprint-uservices/blueprint/examples/postnotification_simple/workflow/postnotification_simple.NotifyServiceImpl
					serviceImplPath := typeNamed.String()
					service := graph.GetApp().GetServiceWithImplPath(serviceImplPath)
					logrus.Tracef("[CALLS BLUEPRINT] [QUEUE] service fields: %v\n", service.GetAllFields())
					field := service.GetFieldAt(ssaFieldAddr.Field)
					logrus.Tracef("[CALLS BLUEPRINT] [QUEUE] field: %s\n", field.String())

					database := field.GetWiringName()

					if database == "" {
						logrus.WithField("graph", graph.String()).WithField("field", field.String()).Fatalf("[CALLS BLUEPRINT] [QUEUE] empty database name!\n")
					}

					// sanity check
					// keep this while database logic is not complete
					if !graph.GetApp().HasDatabase(database) {
						logrus.Fatalf("[CALLS BLUEPRINT] [QUEUE] database (%s) not found for app with databases: %v", database, graph.GetApp().GetAllDatabases())
					}

					return database, true
				}
			}
		}
	}
	return "", false
}
