package tainter

import (
	"fmt"
	"go/types"
	"log"
	"slices"

	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/app"
	"analyzer/pkg/common"
	"analyzer/pkg/ssagraph"
	"analyzer/pkg/utils"
)

const BLUEPRINT_BACKEND_PACKAGE = "github.com/blueprint-uservices/blueprint/runtime/core/backend"

var BLUEPRINT_BACKEND_QUEUE_CALLS = []string{"Push", "Pop"}
var BLUEPRINT_BACKEND_NOSQLDATABASE_CALLS = []string{"GetCollection"}
var BLUEPRINT_BACKEND_NOSQLCOLLECTION_CALLS = []string{"InsertOne", "FindOne", "DeleteOne"}
var BLUEPRINT_BACKEND_RELATIONALDB_CALLS = []string{"Exec", "Select"}

type ValFieldPath struct {
	val       ssa.Value
	fieldpath string
}

func isServiceCall(graph *ssagraph.SSAGraph, instr ssa.Instruction) (string, string, string, []ssa.Value, *ssa.Call, bool) {
	if call, ok := instr.(*ssa.Call); ok {
		fmt.Printf("[TAINT] checking for service call: %s\n", instr.String())

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
						log.Fatalf("method (%s) not found for service (%s)", method, serviceName)
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
		fmt.Printf("[TAINT] checking for database call: %s\n", val.String())

		// --------------
		// blueprint apps
		// --------------
		if unOp, ok := call.Call.Value.(*ssa.UnOp); ok {
			if queue, topic, opType, valFieldPathLst, ok := isBlueprintQueueCall(graph, call, unOp); ok {
				// return all args except context
				// NOTE: in this case (when call.Call.Value is UnOp) call.Call.Args does not contain the receiver
				return queue, topic, call.Call.Method.Id(), opType, valFieldPathLst, true
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

func isBlueprintNoSQLDatabaseCall(graph *ssagraph.SSAGraph, call *ssa.Call, unOp *ssa.UnOp) bool {
	if typeNamed, ok := unOp.Type().(*types.Named); ok {
		if typeNamed.String() == BLUEPRINT_BACKEND_PACKAGE+".NoSQLDatabase" {
			if !slices.Contains(BLUEPRINT_BACKEND_NOSQLDATABASE_CALLS, call.Call.Method.Name()) {
				return false
			}
			// call for NoSQLDatabase.GetCollection(...)
			// skip for now
			return true
		}
	}
	return false
}

// TODO: get database name (not the db name of mongodb!)
func isBlueprintNoSQLCollectionCall(graph *ssagraph.SSAGraph, call *ssa.Call, extr *ssa.Extract) (string, string, common.DatabaseOperationType, []ValFieldPath, bool) {
	var opType common.DatabaseOperationType
	if typeNamed, ok := extr.Type().(*types.Named); ok {
		if typeNamed.String() == BLUEPRINT_BACKEND_PACKAGE+".NoSQLCollection" {
			if !slices.Contains(BLUEPRINT_BACKEND_NOSQLCOLLECTION_CALLS, call.Call.Method.Name()) {
				return "", "", -1, nil, false
			}

			switch call.Call.Method.Name() {
			case "FindOne":
				opType = common.OP_READ
			case "InsertOne":
				opType = common.OP_WRITE
			case "DeleteOne":
				opType = common.OP_DELETE
			default:
				log.Fatalf("unknown method name for queue call: %s\n", call.String())
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
						log.Fatalf("database (%s) not found for app: %s", database, graph.GetApp().String())
					}

					valFieldPathLst := make([]ValFieldPath, 1)
					docVal := call.Call.Args[1]
					valFieldPathLst[0] = ValFieldPath{val: docVal, fieldpath: database + "." + collection}
					return database, collection, opType, valFieldPathLst, true
				}
			}
		}
	}
	return "", "", -1, nil, false
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
			if !slices.Contains(BLUEPRINT_BACKEND_RELATIONALDB_CALLS, call.Call.Method.Name()) {
				return "", "", -1, nil, false
			}
			fmt.Printf("[CALLS RELATIONALDB] found RelationalDB call: %v\n", call)

			switch call.Call.Method.Name() {
			case "Select":
				opType = common.OP_READ
			case "Exec": // can also be update or read
				opType = common.OP_WRITE
			default:
				log.Fatalf("unknown method name for queue call: %s\n", call.String())
			}

			
			var stmtVal, sliceArgsVal ssa.Value
			if opType == common.OP_READ {
				// e.g., Select(ctx, &movieId, "SELECT * FROM movieid WHERE movieid = ?", movieID)
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
				fmt.Printf("[CALLS RELATIONALDB] on ssa slice: %v\n", slice)
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
					fmt.Printf("[CALLS RELATIONALDB] on alloc node: %v\n", allocNode)
					for _, edge := range graph.GetEdgesFromNode(allocNode) {
						if edge.GetType() == ssagraph.EDGE_INDEX {
							idxNode := edge.GetToNode()
							fmt.Printf("[CALLS RELATIONALDB] on idx node: %v\n", idxNode)
							for _, edge := range graph.GetEdgesFromNode(idxNode) {
								if edge.GetType() == ssagraph.EDGE_STORE_ADDRESS {
									storeNode := edge.GetToNode()
									fmt.Printf("[CALLS RELATIONALDB] on store node: %v\n", storeNode)
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
				fmt.Printf("[CALLS RELATIONALDB] got filter fields: %v\n", fields)
				tableName = tables[0]
			} else if opType == common.OP_WRITE {
				fields, _, tableName = app.ParseSQLWrite(database, stmt)
				fmt.Printf("[CALLS RELATIONALDB] got written fields: %v\n", fields)
			}

			if len(argVals) != len(fields) {
				log.Fatalf("[CALLS RELATIONALDB] length of arg vals (%d) does not match length fields (%d)\n", len(argVals), len(fields))
			}

			valFieldPathLst := make([]ValFieldPath, len(argVals))
			for i, field := range fields {
				argVal := argVals[i]
				valFieldPath := ValFieldPath{val: argVal, fieldpath: field}
				valFieldPathLst[i] = valFieldPath
			}

			if opType == common.OP_READ {
				if call.Call.Method.Name() == "Select" {
					// select method reads entire row
					fetchToVal := call.Call.Args[2]
					readField := readFields[0]
					valFieldPathLst = append(valFieldPathLst, ValFieldPath{val: fetchToVal, fieldpath: readField})
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
			if !slices.Contains(BLUEPRINT_BACKEND_QUEUE_CALLS, call.Call.Method.Name()) {
				return "", "", -1, nil, false
			}

			switch call.Call.Method.Name() {
			case "Pop":
				opType = common.OP_READ
			case "Push":
				opType = common.OP_WRITE
			default:
				log.Fatalf("unknown method name for queue call: %s\n", call.String())
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

// extractDatabaseNameFromUnOp can be used for RelationalDB and Queue calls
// it cannot be used for NoSQLDatabase calls because the collection is extracted beforehand
func extractDatabaseNameFromUnOp(graph *ssagraph.SSAGraph, unOp *ssa.UnOp) (string, bool) {
	if ssaFieldAddr, ok := unOp.X.(*ssa.FieldAddr); ok {
		if ssaParam, ok := ssaFieldAddr.X.(*ssa.Parameter); ok {
			fmt.Printf("[TAINT - QUEUE] queue loaded from parameter (%d)\n", ssaFieldAddr.Field)
			if typesPointer, ok := ssaParam.Type().(*types.Pointer); ok {
				if typeNamed, ok := typesPointer.Elem().(*types.Named); ok {
					// e.g., github.com/blueprint-uservices/blueprint/examples/postnotification_simple/workflow/postnotification_simple.NotifyServiceImpl
					serviceImplPath := typeNamed.String()
					service := graph.GetApp().GetServiceWithImplPath(serviceImplPath)
					field := service.GetFieldAt(ssaFieldAddr.Field)

					database := field.GetWiringName()

					// sanity check
					// keep this while database logic is not complete
					if !graph.GetApp().HasDatabase(database) {
						log.Fatalf("database (%s) not found for app: %s", database, graph.GetApp().String())
					}

					return database, true
				}
			}
		}
	}
	return "", false
}

// TODO
func parseArgumentsForMongoDBFilter(graph *ssagraph.SSAGraph, bsonFilter ssa.Value) ([]ssa.Value, []string) {
	var args []ssa.Value
	var keys []string
	bsonFilterNode := graph.GetNodeByName(bsonFilter.Name())
	bsonFilterAllocNode := graph.GetEdgesToNodeExceptPointerTo(bsonFilterNode)[0].GetFromNode()
	elemNode := graph.GetEdgesFromNodeExceptPointerTo(bsonFilterAllocNode)[0].GetToNode()
	bsonFilterKeyNode := graph.GetEdgesFromNode(elemNode)[0].GetToNode()
	// only 1 expected
	edge := recurseEdgesForwardUntilStoreAddress(graph, bsonFilterKeyNode, nil, make(map[*ssagraph.SSANode]bool))[0]
	key := edge.GetToNode().GetInstruction().(*ssa.Store).Val.(*ssa.Const).Value.ExactString()
	keys = append(keys, "."+key)
	arg := graph.GetEdgesFromNode(elemNode)[1].GetToNode().GetValue()
	args = append(args, arg)
	return args, keys
}

func recurseEdgesForwardUntilStoreAddress(graph *ssagraph.SSAGraph, node *ssagraph.SSANode, storeEdges []*ssagraph.SSAEdge, visited map[*ssagraph.SSANode]bool) []*ssagraph.SSAEdge {
	if _, ok := visited[node]; ok {
		return storeEdges
	}
	visited[node] = true

	for _, edge := range graph.GetEdgesFromNode(node) {
		if edge.GetType() == ssagraph.EDGE_STORE_ADDRESS {
			storeEdges = append(storeEdges, edge)
		} else if edge.GetType() == ssagraph.EDGE_FIELD || edge.GetType() == ssagraph.EDGE_INDEX || edge.GetType() == ssagraph.EDGE_USAGE {
			storeEdges = append(storeEdges, recurseEdgesForwardUntilStoreAddress(graph, edge.GetToNode(), storeEdges, visited)...)
		}
	}
	return storeEdges
}
