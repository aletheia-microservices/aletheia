package abstractcallgraph

import (
	"fmt"
	"log"

	"analyzer/pkg/app/backends"
	"analyzer/pkg/ssagraph"
	"analyzer/pkg/utils"
)

func (graph *AbstractCallGraph) ssaTaintToAbstractTaint(ssaTaintsMap map[string][]*ssagraph.SSATaint) map[string][]*AbstractTaint {
	abstractTaintsMap := make(map[string][]*AbstractTaint, len(ssaTaintsMap))
	for objPath, ssaTaints := range ssaTaintsMap {
		abstractTaints := make([]*AbstractTaint, len(ssaTaints))
		for i, ssaTaint := range ssaTaints {
			dbPath := ssaTaint.GetDbCall().GetDatabasePath()
			dbname := ssaTaint.GetDbCall().GetDatabaseName()
			dbNode := graph.GetNodeByNameIfExists(dbPath)
			if dbNode == nil {
				dbNode = NewAbstractNode(dbPath, NODE_DATABASE, "", "", dbname)
				graph.AddNode(dbPath, dbNode)

				if !graph.GetApp().HasDatabase(dbname) { // sanity check
					graph.GetApp().AddDatabase(backends.NewDatabase(dbname, backends.NewSchema()))
				}
			}

			abstractTaints[i] = NewAbstractTaint(ssaTaint.GetDbField(), ssaTaint.GetDbCall().GetID(), true)
		}
		abstractTaintsMap[objPath] = abstractTaints
	}
	return abstractTaintsMap
}

func (graph *AbstractCallGraph) Parse(funcshortpath string, funcGraphs map[string]*ssagraph.SSAGraph) {
	// dummy node
	clientNode := graph.GetNodeByNameIfExists("client")
	if clientNode == nil {
		clientNode = NewAbstractNode("client", NODE_CLIENT, "", "", "")
		graph.AddNode("client", clientNode)
	}

	ssaGraph := funcGraphs[funcshortpath]

	name := ssaGraph.GetServiceWithMethod()
	node := graph.GetNodeByNameIfExists(name)
	if node == nil {
		node = NewAbstractNode(name, NODE_SERVICE, ssaGraph.GetService(), ssaGraph.GetMethodName(), "")
		graph.AddNode(name, node)

		fmt.Printf("[ABSTRACTGRAPH] creating node with (%d) params: %s\n", len(ssaGraph.GetFuncParametersExceptMemberAndContext()), node)
		for _, funcParam := range ssaGraph.GetFuncParametersExceptMemberAndContext() {
			obj := NewAbstractObject(funcParam.GetName(), graph.ssaTaintToAbstractTaint(funcParam.GetTaints()))
			node.AddParam(obj)
		}
	}

	fmt.Printf("[ABSTRACTGRAPH] parsing returns for node: %s\n", node.String())
	retsLst := ssaGraph.GetReturnsLst()
	var retsObjs []*AbstractObject
	// first, just create new abstract objects using the first set of returns (could be any other)
	for i, ret := range retsLst[0] {
		obj := NewAbstractObject(ret.GetValue().Type().String(), graph.ssaTaintToAbstractTaint(ret.GetTaints()))
		node.AddReturn(obj)
		retsObjs = append(retsObjs, obj)
		fmt.Printf("\t[ABSTRACTGRAPH] [index=%d] added new return object (%s)\n", i, obj.String())
	}
	// then, merge taints with corresponding object in the remaining set of returns
	if len(retsLst) > 1 {
		for _, rets := range retsLst[1:] {
			for i, ret := range rets {
				obj := retsObjs[i]
				obj.MergeTaints(graph.ssaTaintToAbstractTaint(ret.GetTaints()), true)
				fmt.Printf("\t\t[ABSTRACTGRAPH] [index=%d] merged taints from (%s) to (%s)\n", i, ret.GetName(), obj.String())
			}
		}
	}
	// debug
	for i, obj := range node.GetReturns() {
		fmt.Printf("\t[ABSTRACTGRAPH] [index=%d] final taints for object (%s):\n%s\n", i, obj.String(), obj.taintString())
	}

	// build dummy edges for entrypoints
	// FIXME: must not always happen!
	edge := NewAbstractEdge(funcshortpath, utils.ExtractMethodNameFromShortFunctionPath(funcshortpath), clientNode, node, EDGE_SERVICE_ENTRYPOINT)
	for _, funcParam := range ssaGraph.GetFuncParametersExceptMemberAndContext() {
		arg := NewAbstractObject(funcParam.GetName(), make(map[string][]*AbstractTaint))
		edge.AddArgument(arg)
	}
	graph.AddEdge(edge.GetID(), edge)

	// build edges for service/database RPCs/calls
	if ssaGraph.HasServiceCalls() {
		fmt.Printf("[ABSTRACTGRAPH] [%s] found function (%s) with service calls\n", ssaGraph.GetService(), funcshortpath)
		for _, call := range ssaGraph.GetServiceCalls() {
			toName := call.GetServiceWithMethod()
			toNode := graph.GetNodeByNameIfExists(toName)

			toSSAGraph := funcGraphs[call.GetFuncShortPath()]
			if toSSAGraph == nil {
				log.Fatalf("could not find ssa graph for short func path (%s)", call.GetFuncShortPath())
			}

			// create node for the first time
			if toNode == nil {
				toNode = NewAbstractNode(toName, NODE_SERVICE, call.GetService(), call.GetMethod(), "")
				graph.AddNode(toName, toNode)

				fmt.Printf("[ABSTRACTGRAPH] creating toNode with (%d) params: %s\n", len(toSSAGraph.GetFuncParametersExceptMemberAndContext()), toNode)
				for _, funcParam := range toSSAGraph.GetFuncParametersExceptMemberAndContext() {
					param := NewAbstractObject(funcParam.GetName(), graph.ssaTaintToAbstractTaint(funcParam.GetTaints()))
					toNode.AddParam(param)
				}
			}

			edge := NewAbstractEdge(call.GetID(), call.GetMethod(), node, toNode, EDGE_SERVICE_RPC)

			// create call arguments
			for _, callArg := range call.GetArguments() {
				arg := NewAbstractObject(callArg.GetName(), graph.ssaTaintToAbstractTaint(callArg.GetTaints()))
				edge.AddArgument(arg)
			}

			// create call returns
			for _, callRet := range call.GetReturns() {
				ret := NewAbstractObject(callRet.GetName(), graph.ssaTaintToAbstractTaint(callRet.GetTaints()))
				fmt.Printf("[ABSTRACTGRAPH] [%s] added return object (%s) with taints: %v\n", node.String(), ret.String(), callRet.GetTaints())
				edge.AddReturn(ret)
			}

			// propagate taints across services (forwards): fromArgs >>> toParams
			for i, toParam := range toNode.GetParams() {
				fromArg := edge.GetArgumentAt(i)
				toParam.MergeTaints(fromArg.GetPrimaryTaints(), false)
			}

			// propagate taints across services (backwards): fromArgs <<< toParams
			for i, fromArg := range edge.GetArguments() {
				toParam := toNode.GetParamAt(i)
				fromArg.MergeTaints(toParam.GetPrimaryTaints(), false)
			}

			graph.AddEdge(edge.GetID(), edge)
		}
		fmt.Println()
	}

	if ssaGraph.HasDatabaseCalls() {
		fmt.Printf("[ABSTRACTGRAPH] found [%s] function (%s) with database calls\n", ssaGraph.GetService(), funcshortpath)

		for _, call := range ssaGraph.GetDatabaseCalls() {
			toDatabasePath := call.GetDatabasePath()
			toNode := graph.GetNodeByNameIfExists(toDatabasePath)
			dbname := call.GetDatabaseName()
			if toNode == nil {
				toNode = NewAbstractNode(toDatabasePath, NODE_DATABASE, "", "", dbname)
				graph.AddNode(toDatabasePath, toNode)

				if !graph.GetApp().HasDatabase(dbname) { // sanity check
					graph.GetApp().AddDatabase(backends.NewDatabase(dbname, backends.NewSchema()))
				}
			}

			edge := NewAbstractEdge(call.GetID(), call.GetMethod(), node, toNode, EDGE_DATABASE_CALL)

			for _, callArg := range call.GetArguments() {
				arg := NewAbstractObject(callArg.GetName(), graph.ssaTaintToAbstractTaint(callArg.GetTaints()))
				edge.AddArgument(arg)
			}

			// propagate taints (indirect): fromParams >>> toArgs (NOT THE SAME INDEX!!)
			// this propagates any new secondary taints that were newly added when calling this service
			//
			// keys are the fieldpath of the current database
			// values are a list of corresponding fieldpaths in other databases (future candidates to foreign key)
			taintpairsmap := make(map[string][]string)
			currdbname := toNode.GetDatabaseName()
			for _, fromParam := range node.GetParams() {
				for _, taintLst := range fromParam.GetAllTaints() {
					var otherFieldPaths []string
					var currFieldPath string
					for _, taint := range taintLst {
						fieldpath := taint.GetDbField()
						otherdbname := utils.ExtractDatabaseNameFromFieldPath(fieldpath)
						if currdbname != otherdbname {
							otherFieldPaths = append(otherFieldPaths, fieldpath)
						} else {
							currFieldPath = fieldpath
						}
					}
					if currFieldPath != "" {
						taintpairsmap[currFieldPath] = otherFieldPaths
					}
				}
			}
			for _, toArg := range edge.GetArguments() {
				// same logic as standard taints
				// key is the object path
				// value is a list of fieldpaths
				newtaints := make(map[string][]*AbstractTaint)
				for objpath, argTaintLst := range toArg.GetAllTaints() {
					for _, argTaint := range argTaintLst {
						// get list of corresponding fieldpaths of the other databases for field in current taint
						taintpairslst, ok := taintpairsmap[argTaint.dbfield]
						if !ok {
							continue
						}

						fmt.Printf("GOT TAINT PAIR (%v) FOR DB FIELD (%s)\n", taintpairslst, argTaint.dbfield)
						for _, dbfield := range taintpairslst {
							for _, argTaint2 := range argTaintLst {
								// skip if there is already a taint with the fieldpath of the other database
								if argTaint2.GetDbField() == dbfield {
									continue
								}
							}

							// associate current object path to fieldpath in other database
							newtaints[objpath] = append(newtaints[objpath], &AbstractTaint{
								dbfield: dbfield,
								primary: false,
							})
						}
					}
				}

				fmt.Printf("GOT NEW TAINTS: %v\n", newtaints)
				toArg.MergeTaints(newtaints, false)
			}

			fmt.Printf("[ABSTRACTGRAPH] getting current database for name (%s)\n", toNode.GetDatabaseName())
			currDb := graph.GetApp().GetDatabaseByName(toNode.GetDatabaseName())
			//TODO: deal with paths that are prefixes and dbpath is not the same
			for _, arg := range edge.GetArguments() {
				if arg.IsTainted() {
					for _, taints := range arg.GetAllTaints() {
						var otherFields []*backends.Field
						var currField *backends.Field
						for _, taint := range taints {
							dbfieldpath := taint.GetDbField()
							dbname := utils.ExtractDatabaseNameFromFieldPath(dbfieldpath)
							fmt.Printf("[ABSTRACTGRAPH] getting other database for name (%s)\n", dbname)
							otherDb := graph.GetApp().GetDatabaseByName(dbname)

							// create new field does not exist yet
							var field *backends.Field
							if !otherDb.GetSchema().HasField(dbfieldpath) {
								field = backends.NewField(dbfieldpath, otherDb)
								otherDb.GetSchema().AddField(field)
							} else {
								field = otherDb.GetSchema().GetFieldByPath(dbfieldpath)
							}
							if currDb == otherDb {
								currField = field
							} else {
								otherFields = append(otherFields, field)
							}
						}

						fmt.Printf("[ABSTRACTGRAPH] curr field = (%v) // other fields = (%v)\n", currField, otherFields)

						// create new foreign key constraint
						if currField != nil {
							for _, otherField := range otherFields {
								if !currField.HasConstraintForeignKeyToField(otherField) {
									constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, currField, otherField)
									if constraint.GetFields()[1] == nil {
										log.Fatal("unexpected!")
									}
									currField.AddConstraint(constraint)
									currField.GetDatabase().GetSchema().AddConstraint(constraint)
								}
							}
						}
					}
				}
			}

			graph.AddEdge(edge.GetID(), edge)
		}
		fmt.Println()

		for _, call := range ssaGraph.GetServiceCalls() {
			graph.Parse(call.GetFuncShortPath(), funcGraphs)
		}
	}
}

func (graph *AbstractCallGraph) Visit(node *AbstractNode) {
	fmt.Printf("[ABSTRACTCALLGRAPH] visiting service node: %s\n", node.String())
	for _, edge := range graph.GetEdgesFromNode(node) {
		if edge.GetEdgeType() == EDGE_SERVICE_ENTRYPOINT {
			graph.Visit(edge.GetToNode())
		}
		
		if edge.GetEdgeType() == EDGE_SERVICE_RPC {
			fmt.Printf("\t[ABSTRACTGRAPH] visiting service call: %s\n", edge.String())
			toNode := edge.GetToNode()
			taintMapping := &TaintMapping{mapping: make(map[string][]string)}
	
			// propagate taints across services (backwards): args (from) <<< params (to)
			for i, fromArg := range edge.GetArguments() {
				toParam := toNode.GetParamAt(i)
				taintMappingTmp := fromArg.MergeTaints(toParam.GetPrimaryTaints(), false)
				taintMapping.Merge(taintMappingTmp)
				fmt.Printf("\t\t[ABSTRACTGRAPH] [ARGS] [index=%d] taint mapping for arg (%s): %s\n", i, fromArg.String(), taintMappingTmp.String())
			}
			
			// propagate taints across services (backwards): rets (from) <<< rets (to)
			for i, fromRet := range edge.GetReturns() {
				toRet := toNode.GetReturnAt(i)
				fmt.Printf("to ret: (%s): %v\n", toRet.String(), toRet.GetPrimaryTaints())
				taintMappingTmp := fromRet.MergeTaints(toRet.GetPrimaryTaints(), false)
				taintMapping.Merge(taintMappingTmp)
				fmt.Printf("\t\t[ABSTRACTGRAPH] [RETS] [index=%d] taint mapping for ret (%s): %s\n", i, fromRet.String(), taintMappingTmp.String())
			}

			fmt.Printf("\t\t[ABSTRACTGRAPH] final taint mapping: %s\n", taintMapping.String())

			// propagate taints into databases
			for currFieldPath, otherFieldsPathLst := range taintMapping.mapping {
				currDb := graph.GetApp().GetDatabaseByName(utils.ExtractDatabaseNameFromFieldPath(currFieldPath))

				var currField *backends.Field
				if !currDb.GetSchema().HasField(currFieldPath) {
					currField = backends.NewField(currFieldPath, currDb)
					currDb.GetSchema().AddField(currField)
				} else {
					currField = currDb.GetSchema().GetFieldByPath(currFieldPath)
				}

				for _, otherFieldPath := range otherFieldsPathLst {
					otherField := graph.GetApp().GetDatabaseByName(utils.ExtractDatabaseNameFromFieldPath(otherFieldPath)).GetSchema().GetFieldByPath(otherFieldPath)
					constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, currField, otherField)
					currField.AddConstraint(constraint)
					currDb.GetSchema().AddConstraint(constraint)
					fmt.Printf("\t\t[ABSTRACTGRAPH] added new constraint: %s\n", constraint)
				}
			}

			fmt.Println()
			graph.Visit(edge.GetToNode())
		}
	}
}
