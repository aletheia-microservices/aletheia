package tainter

import (
	"fmt"
	"go/types"
	"log"
	"slices"
	"strings"

	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/ssagraph"
	"analyzer/pkg/utils"
)

const BLUEPRINT_BACKEND_PACKAGE = "github.com/blueprint-uservices/blueprint/runtime/core/backend"

var BLUEPRINT_BACKEND_QUEUE_CALLS = []string{"Push", "Pop"}
var BLUEPRINT_BACKEND_NOSQLDATABASE_CALLS = []string{"GetCollection"}
var BLUEPRINT_BACKEND_NOSQLCOLLECTION_CALLS = []string{"InsertOne", "FindOne"}

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

		// ------------
		// example apps
		// ------------
		if fn, ok := call.Call.Value.(*ssa.Function); ok && len(fn.Params) > 0 {
			maybeRcv := fn.Params[0]
			if maybeRcv.Type().String() == "*main.ShippingService" && fn.Name() == "NewShipment" {
				// return all args except receiver and context
				return "ShippingService", "NewShipment", "", call.Call.Args[2:], call, true
			}
			if maybeRcv.Type().String() == "*main.SkuService" && fn.Name() == "GetSku" {
				// return all args except receiver and context
				return "SkuService", "GetSku", "", call.Call.Args[2:], call, true
			}
			if maybeRcv.Type().String() == "*main.AnalyticsService" && fn.Name() == "UpdateAnalytics" {
				// return all args except receiver and context
				return "AnalyticsService", "UpdateAnalytics", "", call.Call.Args[2:], call, true
			}
			if slices.Contains([]string{
				"StorePost", "ReadPost", "DeletePost", // storage
				"ReadAnalytics",                                     // analytics
				"UploadPost", "DeletePost", "ReadPostWithAnalytics", // upload
			}, fn.Name()) {
				log.Fatal("EXIT!")
				// return all args except receiver and context
				return "", "", "", call.Call.Args[2:], call, true
			}
		}
	}
	return "", "", "", nil, nil, false
}

func isDatabaseCall(graph *ssagraph.SSAGraph, instr ssa.Instruction) (string, string, string, []ssa.Value, bool, bool) {
	if instr == nil {
		return "", "", "", nil, false, false
	}

	if call, ok := instr.(*ssa.Call); ok {
		fmt.Printf("[TAINT] checking for database call: %s\n", instr.String())

		// --------------
		// blueprint apps
		// --------------
		if unOp, ok := call.Call.Value.(*ssa.UnOp); ok {
			if ok, queue, topic, isWrite := isBlueprintQueueCall(graph, call, unOp); ok {
				// return all args except context
				// NOTE: in this case (when call.Call.Value is UnOp) call.Call.Args does not contain the receiver
				return queue, topic, call.Call.Method.Id(), call.Call.Args[1:], isWrite, true
			} else if ok := isBlueprintNoSQLDatabaseCall(graph, call, unOp); ok {
				// call for NoSQLDatabase.GetCollection(...)
				// skip for now
				return "", "", "", nil, false, false
			}
		}
		if extr, ok := call.Call.Value.(*ssa.Extract); ok {
			if ok, database, collection, isWrite := isBlueprintNoSQLCollectionCall(graph, call, extr); ok {
				return database, collection, call.Call.Method.Id(), call.Call.Args[1:], isWrite, true
			}
		}

		// ------------
		// example apps
		// ------------
		if fn, ok := call.Call.Value.(*ssa.Function); ok && len(fn.Params) > 0 {
			fmt.Printf("[TAINT] [1] found call: %v\n", call)
			maybeRcv := fn.Params[0]
			var isWrite bool
			if maybeRcv.Type().String() == "*main.MongoDB" && fn.Name() == "Insert" || fn.Name() == "Find" {
				if fn.Name() == "Insert" {
					isWrite = true
				}
				// return arg without receiver and context
				return "mydb", "mycollection", call.Call.Method.Id(), call.Call.Args[2:], isWrite, true
			}
			if maybeRcv.Type().String() == "*main.RabbitMQ" && fn.Name() == "Push" {
				isWrite = true
				// return arg without receiver and context
				return "mydb", "mycollection", call.Call.Method.Id(), call.Call.Args[2:], isWrite, true
			}
		}
	}
	return "", "", "", nil, false, false
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

func isBlueprintNoSQLCollectionCall(graph *ssagraph.SSAGraph, call *ssa.Call, extr *ssa.Extract) (bool, string, string, bool) {
	if typeNamed, ok := extr.Type().(*types.Named); ok {
		if typeNamed.String() == BLUEPRINT_BACKEND_PACKAGE+".NoSQLCollection" {
			if !slices.Contains(BLUEPRINT_BACKEND_NOSQLCOLLECTION_CALLS, call.Call.Method.Name()) {
				return false, "", "", false
			}
			var isWrite bool
			if call.Call.Method.Name() == "InsertOne" {
				isWrite = true
			}
			// e.g.
			// t29 = invoke t28.GetCollection(ctx, "posts_db":string, "post":string)
			// t30 = extract t29 #0
			// t34 = invoke t30.InsertOne(ctx, t33)
			if ssaExtract, ok := call.Call.Value.(*ssa.Extract); ok {
				if ssaCall, ok := ssaExtract.Tuple.(*ssa.Call); ok {
					databaseVal := ssaCall.Call.Args[1]
					collectionVal := ssaCall.Call.Args[2]
					var database, collection string
					if c, ok := databaseVal.(*ssa.Const); ok {
						database = strings.Trim(c.Value.ExactString(), "\"")
					}
					if c, ok := collectionVal.(*ssa.Const); ok {
						collection = strings.Trim(c.Value.ExactString(), "\"")
					}

					// sanity check
					if !graph.GetApp().HasDatabase(database) {
						log.Fatalf("database (%s) not found for app: %s", database, graph.GetApp().String())
					}

					return true, database, collection, isWrite
				}
			}
		}
	}
	return false, "", "", false
}

// e.g. where t11 is unaryOp and t10 is its unaryOp.X
// t10 = &u.notificationsQueue [#1]
// t11 = *t10
// t14 = invoke t11.Push(ctx, t13)
func isBlueprintQueueCall(graph *ssagraph.SSAGraph, call *ssa.Call, unOp *ssa.UnOp) (bool, string, string, bool) {
	if typeNamed, ok := unOp.Type().(*types.Named); ok {
		if typeNamed.String() == BLUEPRINT_BACKEND_PACKAGE+".Queue" {
			if !slices.Contains(BLUEPRINT_BACKEND_QUEUE_CALLS, call.Call.Method.Name()) {
				return false, "", "", false
			}
			var isWrite bool
			if call.Call.Method.Name() == "Push" {
				isWrite = true
			}

			/* fmt.Printf("\tcall (sign):\t [%T] [%T] %s\n", call.Call.Method.Type(), call.Call.Method, call.Call.Method)
			fmt.Printf("\tcall (rets):\t [%T] [%T] %s\n", call.Type(), call, call)
			fmt.Printf("\tunary op:\t [%T] [%T] %s\n", unOp.Type(), unOp, unOp)
			fmt.Printf("\tunary op X:\t [%T] [%T] %s\n", unOp.X.Type(), unOp.X, unOp.X)
			fmt.Printf("\tunary op NAME:\t %s\n", unOp.Name()) */

			// e.g., t10 = &u.notificationsQueue [#1]
			if ssaFieldAddr, ok := unOp.X.(*ssa.FieldAddr); ok {
				if ssaParam, ok := ssaFieldAddr.X.(*ssa.Parameter); ok {
					fmt.Printf("[TAINT - QUEUE] queue loaded from parameter (%d)\n", ssaFieldAddr.Field)
					if typesPointer, ok := ssaParam.Type().(*types.Pointer); ok {
						if typeNamed, ok := typesPointer.Elem().(*types.Named); ok {
							// e.g., github.com/blueprint-uservices/blueprint/examples/postnotification_simple/workflow/postnotification_simple.NotifyServiceImpl
							serviceImplPath := typeNamed.String()
							serviceImplName := typeNamed.Obj().Id()
							fmt.Printf("[TAINT - QUEUE] got type name (%s) @ (%s)\n", serviceImplName, serviceImplPath)

							// dummy logic
							database := "notifications_queue"
							schema := "notification"
							// sanity check
							if !graph.GetApp().HasDatabase(database) {
								log.Fatalf("database (%s) not found for app: %s", database, graph.GetApp().String())
							}

							return true, database, schema, isWrite
						}
					}
				}
			}
		}

	}
	return false, "", "", false
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
