package tainter

import (
	"fmt"
	"go/types"
	"log"
	"slices"
	"strings"

	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/app"
	"analyzer/pkg/common"
	"analyzer/pkg/ssagraph"
	"analyzer/pkg/utils"
)

const BLUEPRINT_BACKEND_PACKAGE = "github.com/blueprint-uservices/blueprint/runtime/core/backend"

var BLUEPRINT_BACKEND_CALLS_QUEUE = []string{"Push", "Pop"}
var BLUEPRINT_BACKEND_CALLS_NOSQLDATABASE = []string{"GetCollection"}
var BLUEPRINT_BACKEND_CALLS_NOSQLCOLLECTION = []string{"InsertOne", "FindOne", "DeleteOne", "FindMany"}
var BLUEPRINT_BACKEND_CALLS_NOSQLCURSOR = []string{"One", "All"}
var BLUEPRINT_BACKEND_CALLS_RELATIONALDB = []string{"Exec", "Select"}
var BLUEPRINT_BACKEND_CALLS_CACHE = []string{"Get", "Put"}

type ValFieldPath struct {
	val            ssa.Value
	fieldpath      string
	bsonCursorMany bool // destination objects from nosql cursor reads
	bsonFilterIn   bool // bson filter: $in
}

func isServiceCall(graph *ssagraph.SSAGraph, val ssa.Value) (string, string, string, []ssa.Value, *ssa.Call, bool) {
	if call, ok := val.(*ssa.Call); ok {
		fmt.Printf("[CALLS BLUEPRINT] [SVC] checking for service call: %s\n", val.String())

		// --------------
		// blueprint apps
		// --------------
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
						log.Fatalf("[CALLS BLUEPRINT] [SVC] method (%s) not found for service (%s)", method, serviceName)
					}
				}
			}
		}
	}
	return "", "", "", nil, nil, false
}

func isDatabaseCall(graph *ssagraph.SSAGraph, val ssa.Value) (string, string, string, common.DatabaseOperationType, []ValFieldPath, bool) {
	if val == nil {
		return "", "", "", -1, nil, false
	}

	if call, ok := val.(*ssa.Call); ok {
		fmt.Printf("[CALLS BLUEPRINT] [DB] checking for database call: %s\n", val.String())

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
			fmt.Printf("[CALLS BLUEPRINT] [RELDB] found RelationalDB call: %v\n", call)

			switch call.Call.Method.Name() {
			case "Select":
				opType = common.OP_READ
			case "Exec": // can also be update or read
				opType = common.OP_WRITE
			default:
				log.Fatalf("[CALLS BLUEPRINT] [RELDB] unknown method name for queue call: %s\n", call.String())
			}

			var dstVal, stmtVal, sliceArgsVal ssa.Value
			if opType == common.OP_READ {
				// e.g., Select(ctx, &movieId, "SELECT * FROM movieid WHERE movieid = ?", movieID)
				dstVal = call.Call.Args[1]
				stmtVal = call.Call.Args[2]
				sliceArgsVal = call.Call.Args[3]
			} else if opType == common.OP_WRITE {
				// e.g., Exec(ctx, "INSERT INTO movieid(movieid, title) VALUES (?, ?);", movieID, title)
				stmtVal = call.Call.Args[1]
				sliceArgsVal = call.Call.Args[2]
			}

			stmt, _ := utils.ExtractStringFromValue(stmtVal)

			database, ok := extractDatabaseNameFromUnOp(graph, unOp)
			if !ok {
				return "", "", -1, nil, false
			}

			var argVals []ssa.Value
			if slice, ok := sliceArgsVal.(*ssa.Slice); ok {
				fmt.Printf("[CALLS BLUEPRINT] [RELDB] on ssa slice: %v\n", slice)
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
					fmt.Printf("[CALLS BLUEPRINT] [RELDB] on alloc node: %v\n", allocNode)
					for _, edge := range graph.GetEdgesFromNode(allocNode) {
						if edge.GetType() == ssagraph.EDGE_INDEX {
							idxNode := edge.GetToNode()
							fmt.Printf("[CALLS BLUEPRINT] [RELDB] on idx node: %v\n", idxNode)
							for _, edge := range graph.GetEdgesFromNode(idxNode) {
								if edge.GetType() == ssagraph.EDGE_STORE_ADDRESS {
									storeNode := edge.GetToNode()
									fmt.Printf("[CALLS BLUEPRINT] [RELDB] on store node: %v\n", storeNode)
									storeInstr, _ := storeNode.GetInstruction().(*ssa.Store)
									argVals = append(argVals, storeInstr.Val)
								}
							}
						}
					}
				}
			}

			var fields []string
			var tableName string
			var readFields []string

			if opType == common.OP_READ {
				var tables []string
				readFields, fields, tables = app.ParseSQLRead(database, stmt)
				fmt.Printf("[CALLS BLUEPRINT] [RELDB] got filter fields: %v\n", fields)
				tableName = tables[0]
			} else if opType == common.OP_WRITE {
				fields, _, tableName = app.ParseSQLWrite(database, stmt)
				fmt.Printf("[CALLS BLUEPRINT] [RELDB] got written fields: %v\n", fields)
			}

			if len(argVals) != len(fields) {
				log.Fatalf("[CALLS BLUEPRINT] [RELDB] length of arg vals (%d) does not match length fields (%d)\n", len(argVals), len(fields))
			}

			valFieldPathLst := make([]ValFieldPath, len(argVals))
			for i, field := range fields {
				argVal := argVals[i]
				valFieldPath := ValFieldPath{val: argVal, fieldpath: field}
				valFieldPathLst[i] = valFieldPath
			}

			if opType == common.OP_READ {
				// for SQL Selects on all fields (i.e., '*') the readFields length is 1
				// and the readField has format <database>.<table>
				if call.Call.Method.Name() == "Select" {
					// select method reads entire row
					readField := readFields[0]
					valFieldPathLst = append(valFieldPathLst, ValFieldPath{val: dstVal, fieldpath: readField})
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
				log.Fatalf("[CALLS BLUEPRINT] [QUEUE] unknown method name for queue call: %s\n", call.String())
			}

			// e.g., t10 = &u.notificationsQueue [#1]
			if database, ok := extractDatabaseNameFromUnOp(graph, unOp); ok {
				topic := "notification"
				valFieldPathLst := make([]ValFieldPath, 1)
				docVal := call.Call.Args[1]
				valFieldPathLst[0] = ValFieldPath{val: docVal, fieldpath: database + "." + topic}

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
			case "Put":
				opType = common.OP_WRITE
			default:
				log.Fatalf("[CALLS BLUEPRINT] [CACHE] unknown method name for queue call: %s\n", call.String())
			}

			// e.g., t10 = &u.notificationsQueue [#1]
			if database, ok := extractDatabaseNameFromUnOp(graph, unOp); ok {
				namespace := "*"
				var valFieldPathLst []ValFieldPath
				cacheKeyVal := call.Call.Args[1]
				cacheValueVal := call.Call.Args[2]
				fmt.Printf("[CALLS BLUEPRINT] [CACHE] cache key [%T]: %v\n", cacheKeyVal, cacheKeyVal)
				fmt.Printf("[CALLS BLUEPRINT] [CACHE] cache value [%T]: %v\n", cacheValueVal, cacheValueVal)

				// track cache key
				if _, ok := utils.ExtractStringFromValue(cacheKeyVal); ok {
					valFieldPathLst = append(valFieldPathLst, ValFieldPath{val: cacheKeyVal, fieldpath: database + "." + namespace + ".Key"})
				} else if binOp, ok := cacheKeyVal.(*ssa.BinOp); ok {
					if suffix, ok := utils.ExtractStringFromValue(binOp.Y); ok {
						namespace, _ = strings.CutPrefix(suffix, ":")
						// real cache key
						valFieldPathLst = append(valFieldPathLst, ValFieldPath{val: binOp.X, fieldpath: database + "." + namespace + ".Key"})
					}
				} else if call, ok := cacheKeyVal.(*ssa.Call); ok {
					for _, arg := range call.Call.Args {
						valFieldPathLst = append(valFieldPathLst, ValFieldPath{val: arg, fieldpath: database + "." + namespace + ".Key"})
					}
				}

				if valFieldPathLst == nil {
					// [TO BE IMPROVED]
					valFieldPathLst = append(valFieldPathLst, ValFieldPath{val: cacheKeyVal, fieldpath: database + "." + namespace + ".Key"})
					fmt.Printf("[CALLS CACHE] [%s] could not save any cache key for call: %v\n", graph.String(), call)
				}

				// track cache value
				valFieldPathLst = append(valFieldPathLst, ValFieldPath{val: cacheValueVal, fieldpath: database + "." + namespace + ".Value"})

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
			case "DeleteOne":
				opType = common.OP_DELETE
			default:
				log.Fatalf("[CALLS BLUEPRINT] [NOSQL] unknown method name for queue call: %s\n", call.String())
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
						log.Fatalf("[CALLS BLUEPRINT] [NOSQL] database (%s) extracted from value (%s) not found for app with databases: %v", database, databaseVal.String(), graph.GetApp().GetAllDatabases())
					}

					var valFieldPathLst []ValFieldPath

					if opType == common.OP_WRITE {
						docVal := call.Call.Args[1]
						valFieldPathLst = append(valFieldPathLst, ValFieldPath{
							val:       docVal,
							fieldpath: database + "." + collection,
						})
					} else { // reads, updates, or deletes
						filterVal := call.Call.Args[1]
						for filter, vals := range computeNoSQLFilterKeyToValues(graph, filterVal) {
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

					/* if database == "reservation_db" && call.Call.Method.Name() == "FindMany" {
						fmt.Printf("CALL: %s\n", call.Name())
						for _, val := range valFieldPathLst {
							fmt.Printf("(%s) --> (%s)\n", val.fieldpath, val.val.String())
						}
						log.Fatalf("HERE!")
					} */

					if opType == common.OP_READ || opType == common.OP_READ_MANY {
						for cursorCall, extr := range getNoSQLCursorCallsFromCollectionCall(graph, call) {
							if dstVal, ok := isBlueprintNoSQLCursorCall(graph, cursorCall, extr); ok {
								fieldpath := database + "." + collection
								// distinguish objects used as arguments for any operation
								// from objects used as destination for reads (isDst=true)
								valFieldPathLst = append(valFieldPathLst, ValFieldPath{
									fieldpath:      fieldpath,
									val:            dstVal,
									bsonCursorMany: opType == common.OP_READ_MANY,
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

func computeNoSQLFilterKeyToValues(graph *ssagraph.SSAGraph, bsonVal ssa.Value) map[string][]ValFieldPath {
	filterKeyToValues := make(map[string][]ValFieldPath)
	if slice, ok := bsonVal.(*ssa.Slice); ok && ssaValueIsUsedInMongoBsonFilter(graph, slice) {
		bsonSliceNode := graph.GetNodeByName(slice.X.Name())

		for _, edge := range graph.GetEdgesTypedTo(bsonSliceNode, ssagraph.EDGE_POINTS_TO) {
			if edge.HasPath("[*]") {
				var filterField string
				bsonElemNode := edge.GetFromNode()

				// objects that are excluded from taint:
				// - bson slice
				// - bson slice elems
				// - bson slice elem key
				// - bson slice elem value
				var filterObjs = []ValFieldPath{}

				for _, edge := range graph.GetEdgesTypedFrom(bsonElemNode, ssagraph.EDGE_FIELD) {
					if edge.GetParam() == "Key" {
						// track objects used as value in store instructions for current bson key
						bsonKeyNode := edge.GetToNode()
						for _, edge := range graph.GetEdgesTypedFrom(bsonKeyNode, ssagraph.EDGE_STORE_ADDRESS) {
							storeInstr := edge.GetToNode().GetInstruction().(*ssa.Store)
							filterObjs = append(filterObjs, ValFieldPath{
								val: storeInstr.Val,
							})
							// should only occur once but we keep iterating just for sanity check
							if filterFieldTmp, ok := utils.ExtractStringFromValue(storeInstr.Val); ok {
								filterField = filterFieldTmp
							}
						}
					} else if edge.GetParam() == "Value" {
						// track objects used as value in store instructions for current bson value
						bsonValueNode := edge.GetToNode()
						for _, edge := range graph.GetEdgesTypedFrom(bsonValueNode, ssagraph.EDGE_STORE_ADDRESS) {
							storeInstr := edge.GetToNode().GetInstruction().(*ssa.Store)
							filterObjs = append(filterObjs, ValFieldPath{
								val: storeInstr.Val,
							})
							if iface, ok := storeInstr.Val.(*ssa.MakeInterface); ok {
								for filterFieldTmp, filterObjsTmp := range computeNoSQLFilterKeyToValues(graph, iface.X) {
									switch filterFieldTmp {
									case "$in":
										for _, filterObjTmp := range filterObjsTmp {
											filterObjTmp.bsonFilterIn = true
											filterObjs = append(filterObjs, filterObjTmp)
										}
									default:
										log.Fatalf("[CALLS BLUEPRINT] [BSON] unexpected filter key (%s) for objects: %v", filterFieldTmp, filterObjsTmp)
									}
								}
							}
						}
					}
				}
				if filterField == "" {
					log.Fatalf("[CALLS BLUEPRINT] [BSON] empty filter field for bsonVal (%s) and bsonElem (%s)\n", bsonVal.Name(), bsonElemNode.GetValue().Name())
				}
				filterKeyToValues[filterField] = filterObjs
			}
		}
	}
	fmt.Printf("[CALLS BLUEPRINT] [NOSQL FILTER] returning lst: %v\n", filterKeyToValues)
	return filterKeyToValues
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
					if named.Obj().Pkg().Path() == "go.mongodb.org/mongo-driver/bson/primitive" && named.Obj().Id() == "E" {
						return true
					}
				}
			}
		}
	}

	return false
}

func ssaValueIsUsedInMongoBsonFilter(graph *ssagraph.SSAGraph, val ssa.Value) bool {
	if val == nil {
		return false
	}
	if ok := ssaValueIsMongoBsonFilter(val); ok {
		return true
	}
	if slice, ok := val.(*ssa.Slice); ok {
		// e.g.,
		// [ssa.Slice] t73: slice t60[:]
		if named, ok := slice.Type().(*types.Named); ok {
			// this is the alias type: go.mongodb.org/mongo-driver/bson/primitive.D
			if named.Obj().Pkg().Path() == "go.mongodb.org/mongo-driver/bson/primitive" && named.Obj().Id() == "D" {
				return true
			}
		}
	}
	for _, edge := range graph.GetEdgesTypedFrom(graph.GetNodeByName(val.Name()), ssagraph.EDGE_POINTS_TO) {
		if ssaValueIsUsedInMongoBsonFilter(graph, edge.GetToNode().GetValue()) {
			return true
		}
	}
	return false
}

// extractDatabaseNameFromUnOp can be used for RelationalDB and Queue calls
// it cannot be used for NoSQLDatabase calls because the collection is extracted beforehand
func extractDatabaseNameFromUnOp(graph *ssagraph.SSAGraph, unOp *ssa.UnOp) (string, bool) {
	if ssaFieldAddr, ok := unOp.X.(*ssa.FieldAddr); ok {
		fmt.Printf("[CALLS BLUEPRINT] [QUEUE] ssa field addr (field=%d): %s\n", ssaFieldAddr.Field, ssaFieldAddr.String())
		if ssaParam, ok := ssaFieldAddr.X.(*ssa.Parameter); ok {
			fmt.Printf("[CALLS BLUEPRINT] [QUEUE] queue loaded from parameter (%d)\n", ssaFieldAddr.Field)
			if typesPointer, ok := ssaParam.Type().(*types.Pointer); ok {
				if typeNamed, ok := typesPointer.Elem().(*types.Named); ok {
					// e.g., github.com/blueprint-uservices/blueprint/examples/postnotification_simple/workflow/postnotification_simple.NotifyServiceImpl
					serviceImplPath := typeNamed.String()
					service := graph.GetApp().GetServiceWithImplPath(serviceImplPath)
					fmt.Printf("[CALLS BLUEPRINT] [QUEUE] service fields: %v\n", service.GetAllFields())
					field := service.GetFieldAt(ssaFieldAddr.Field)
					fmt.Printf("[CALLS BLUEPRINT] [QUEUE] field: %s\n", field.String())

					database := field.GetWiringName()

					if database == "" {
						log.Fatalf("[CALLS BLUEPRINT] [QUEUE] empty database name!\n")
					}

					// sanity check
					// keep this while database logic is not complete
					if !graph.GetApp().HasDatabase(database) {
						log.Fatalf("[CALLS BLUEPRINT] [QUEUE] database (%s) not found for app with databases: %v", database, graph.GetApp().GetAllDatabases())
					}

					return database, true
				}
			}
		}
	}
	return "", false
}
