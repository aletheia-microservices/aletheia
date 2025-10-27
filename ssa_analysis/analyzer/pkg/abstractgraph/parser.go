package abstractgraph

import (
	"log"
	"sort"

	"analyzer/pkg/app/backends"
	"analyzer/pkg/common"
	"analyzer/pkg/ssagraph"
	"analyzer/pkg/utils"
)

func ssaTaintDatabaseToAbstractTaint(graph *AbstractCallGraph, ssaTaintsMap map[string][]*ssagraph.SSATaint) map[string][]*AbstractTaint {
	abstractTaintsMap := make(map[string][]*AbstractTaint, 0)
	for objPath, ssaTaints := range ssaTaintsMap {
		var abstractTaints []*AbstractTaint
		for _, ssaTaint := range ssaTaints {
			if ssaTaint.IsDatabaseTaint() {
				dbPath := ssaTaint.GetDatabaseCall().GetDatabasePath()
				dbname := ssaTaint.GetDatabaseCall().GetDatabaseName()
				dbNode := graph.GetNodeByNameIfExists(dbPath)
				schemaName := ssaTaint.GetDatabaseCall().GetSchemaName()
				if dbNode == nil {
					dbNode = NewAbstractNode(dbPath, NODE_DATABASE, "", "", dbname)
					graph.AddNode(dbPath, dbNode)

					if !graph.GetApp().HasDatabase(dbname) {
						log.Fatalf("database (%s) not found", dbname)
					}
					db := graph.GetApp().GetDatabaseByName(dbname)
					if !db.HasSchema(schemaName) {
						db.AddSchema(backends.NewSchema(schemaName))
					}
				}
				taint := NewAbstractTaint(
					ssaTaint.GetDatabasePath(),
					ssaTaint.GetDatabaseCall().GetID(),
					ssaTaint.GetDatabaseCall().GetOpType(),
					true, false,
				)
				abstractTaints = append(abstractTaints, taint)
				//EVAL - fmt.Printf("[SSA TO ABSTRACT TAINT] new taint on object path (%s): %s\n", objPath, taint.LongString())
			}
		}
		if abstractTaints != nil {
			abstractTaintsMap[objPath] = abstractTaints
		}
	}
	return abstractTaintsMap
}

func ssaTaintServiceToAbstractTrace(graph *AbstractCallGraph, ssaTaintsMap map[string][]*ssagraph.SSATaint) map[string][]*AbstractTrace {
	abstractTaintsMap := make(map[string][]*AbstractTrace, 0)
	for objPath, ssaTaints := range ssaTaintsMap {
		var abstractTraces []*AbstractTrace
		for _, ssaTaint := range ssaTaints {
			if ssaTaint.IsServiceTaint() {
				trace := NewAbstractTrace(
					ssaTaint.GetServicePath(),
					ssaTaint.GetServiceCall().GetID(),
				)
				abstractTraces = append(abstractTraces, trace)
				//EVAL - fmt.Printf("[SSA TO ABSTRACT TRACE] new trace on object path (%s): %s\n", objPath, trace.LongString())
			}
		}
		if abstractTraces != nil {
			abstractTaintsMap[objPath] = abstractTraces
		}
	}
	return abstractTaintsMap
}

func Parse(graph *AbstractCallGraph, funcshortpath string, entrypoint bool, funcGraphs map[string]*ssagraph.SSAGraph) {
	// dummy node
	clientNode := graph.GetNodeByNameIfExists("client")
	if clientNode == nil {
		clientNode = NewAbstractNode("client", NODE_CLIENT, "", "", "")
		graph.AddNode("client", clientNode)
	}

	ssaGraph := funcGraphs[funcshortpath]
	//EVAL - fmt.Printf("[ABSTRACTGRAPH] got ssa graph for (%s): %v\n", funcshortpath, ssaGraph)

	name := ssaGraph.GetServiceWithMethod()
	node := graph.GetNodeByNameIfExists(name)

	if node != nil && node.IsParsed() {
		//EVAL - fmt.Printf("[ABSTRACTGRAPH] ignoring node already visited: %s\n", node.String())
		return
	}

	if node == nil {
		node = NewAbstractNode(name, NODE_SERVICE, ssaGraph.GetService(), ssaGraph.GetMethodName(), "")
		graph.AddNode(name, node)

		//EVAL - fmt.Printf("[ABSTRACTGRAPH] creating node with (%d) params: %s\n", len(ssaGraph.GetFuncParametersExceptMemberAndContext()), node)
		for _, funcParam := range ssaGraph.GetFuncParametersExceptMemberAndContext() {
			obj := NewAbstractObject(funcParam.GetName(), ssaTaintDatabaseToAbstractTaint(graph, funcParam.GetTaints()), ssaTaintServiceToAbstractTrace(graph, funcParam.GetTaints()))
			//EVAL - fmt.Printf("[debug] (1) added param (%s) to node (%s)\n", obj.String(), node.String())
			node.AddParam(obj)
		}
	}

	node.SetParsed()

	// finalize parsing
	//EVAL - fmt.Printf("[ABSTRACTGRAPH] parsing returns for node: %s\n", node.String())
	retsLst := ssaGraph.GetReturnsLst()
	var retsObjs []*AbstractObject
	// first, just create new abstract objects using the first set of returns (could be any other)
	for _, ret := range retsLst[0] {
		obj := NewAbstractObject(ret.GetValue().Type().String(), ssaTaintDatabaseToAbstractTaint(graph, ret.GetTaints()), ssaTaintServiceToAbstractTrace(graph, ret.GetTaints()))
		node.AddReturn(obj)
		retsObjs = append(retsObjs, obj)
		//EVAL - fmt.Printf("\t[ABSTRACTGRAPH] [index=%d] added new return object (%s)\n", i, obj.String())
	}
	// then, merge taints with corresponding object in the remaining set of returns
	if len(retsLst) > 1 {
		for _, rets := range retsLst[1:] {
			for i, ret := range rets {
				obj := retsObjs[i]
				MergeTaints(obj, ssaTaintDatabaseToAbstractTaint(graph, ret.GetTaints()), true, false)
				//EVAL - fmt.Printf("\t\t[ABSTRACTGRAPH] [index=%d] merged taints from (%s) to (%s)\n", i, ret.GetName(), obj.String())
			}
		}
	}
	// debug
	/* for i, obj := range node.GetReturns() {
		//EVAL - fmt.Printf("\t[ABSTRACTGRAPH] [index=%d] final taints for object (%s):\n%s\n", i, obj.String(), obj.TaintLongString())
	} */

	// build dummy edges for entrypoints
	if entrypoint {
		edge := NewAbstractEdge(funcshortpath, utils.ExtractMethodNameFromShortFunctionPath(funcshortpath), clientNode, node, common.OP_UNDEFINED, EDGE_SERVICE_ENTRYPOINT)
		for _, funcParam := range ssaGraph.GetFuncParametersExceptMemberAndContext() {
			arg := NewAbstractObject(funcParam.GetName(), make(map[string][]*AbstractTaint), make(map[string][]*AbstractTrace))
			edge.AddArgument(arg)
		}
		graph.AddEdge(edge)
	}

	var edges []*AbstractEdge

	// build edges for service/database RPCs/calls
	if ssaGraph.HasServiceCalls() {
		//EVAL - fmt.Printf("[ABSTRACTGRAPH] [%s] found function (%s) with service calls\n", ssaGraph.GetService(), funcshortpath)
		for _, call := range ssaGraph.GetServiceCalls() {
			toName := call.GetServiceWithMethod()
			toNode := graph.GetNodeByNameIfExists(toName)

			/* if len(ssaGraph.GetServiceCalls()) > 2 {
				for _, call := range ssaGraph.GetServiceCalls() {
					//EVAL - fmt.Printf("call: %v\n", call.String())
				}
				log.Fatalf("EXIT!")
			} */

			toSSAGraph := funcGraphs[call.GetFuncShortPath()]
			if toSSAGraph == nil {
				log.Fatalf("could not find ssa graph for short func path (%s)", call.GetFuncShortPath())
			}

			// create node for the first time
			if toNode == nil {
				toNode = NewAbstractNode(toName, NODE_SERVICE, call.GetService(), call.GetMethod(), "")
				graph.AddNode(toName, toNode)

				//EVAL - fmt.Printf("[ABSTRACTGRAPH] creating toNode with (%d) params: %s\n", len(toSSAGraph.GetFuncParametersExceptMemberAndContext()), toNode)
				for _, funcParam := range toSSAGraph.GetFuncParametersExceptMemberAndContext() {
					param := NewAbstractObject(funcParam.GetName(), ssaTaintDatabaseToAbstractTaint(graph, funcParam.GetTaints()), ssaTaintServiceToAbstractTrace(graph, funcParam.GetTaints()))
					toNode.AddParam(param)
					//EVAL - fmt.Printf("[debug] (2) added param (%s) to node (%s)\n", param.String(), toNode.String())
				}
			}

			edge := NewAbstractEdge(call.GetID(), call.GetMethod(), node, toNode, common.OP_UNDEFINED, EDGE_SERVICE_RPC)

			// create call arguments
			for _, callArg := range call.GetArguments() {
				arg := NewAbstractObject(callArg.GetName(), ssaTaintDatabaseToAbstractTaint(graph, callArg.GetTaints()), ssaTaintServiceToAbstractTrace(graph, callArg.GetTaints()))
				edge.AddArgument(arg)
			}

			// create call returns
			for _, callRet := range call.GetReturns() {
				ret := NewAbstractObject(callRet.GetName(), ssaTaintDatabaseToAbstractTaint(graph, callRet.GetTaints()), ssaTaintServiceToAbstractTrace(graph, callRet.GetTaints()))
				//EVAL - fmt.Printf("[ABSTRACTGRAPH] [%s] added return object (%s) with taints: %v\n", node.String(), ret.String(), callRet.GetTaints())
				edge.AddReturn(ret)
			}

			edges = append(edges, edge)
			//EVAL - fmt.Printf("[ABSTRACT GRAPH] added edge: %v\n", edge)
		}
		//EVAL - fmt.Println()
	}

	if ssaGraph.HasDatabaseCalls() {
		//EVAL - fmt.Printf("[ABSTRACTGRAPH] found [%s] function (%s) with database calls\n", ssaGraph.GetService(), funcshortpath)

		for _, call := range ssaGraph.GetDatabaseCalls() {
			toDatabasePath := call.GetDatabasePath()
			toNode := graph.GetNodeByNameIfExists(toDatabasePath)
			dbname := call.GetDatabaseName()
			if toNode == nil {
				toNode = NewAbstractNode(toDatabasePath, NODE_DATABASE, "", "", dbname)
				graph.AddNode(toDatabasePath, toNode)
				schemaName := call.GetSchemaName()

				if !graph.GetApp().HasDatabase(dbname) {
					log.Fatalf("database (%s) not found", dbname)
				}

				db := graph.GetApp().GetDatabaseByName(dbname)
				if !db.HasSchema(schemaName) {
					db.AddSchema(backends.NewSchema(schemaName))
				}
			}

			edge := NewAbstractEdge(call.GetID(), call.GetMethod(), node, toNode, call.GetOpType(), EDGE_DATABASE_CALL)

			for _, callArg := range call.GetArguments() {
				arg := NewAbstractObject(callArg.GetName(), ssaTaintDatabaseToAbstractTaint(graph, callArg.GetTaints()), ssaTaintServiceToAbstractTrace(graph, callArg.GetTaints()))
				edge.AddArgument(arg)
			}

			// create fields if they do not exist yet
			registerDatabaseFields(graph, edge.GetArguments())

			// propagate taints to databases (forward): args (from) >>> params (to)
			taintMapping := NewTaintMapping()
			for i, toParam := range toNode.GetParams() {
				fromArg := edge.GetArgumentAt(i)
				taintMappingTmp := MergeTaints(toParam, fromArg.GetPrimaryTaints(), true, false)
				taintMapping.Merge(taintMappingTmp)
			}
			PropagateNewTaintsToDatabaseSchemas(graph, -1, taintMapping)

			edges = append(edges, edge)
		}
		//EVAL - fmt.Println()

		// at the end, we need to sort edges by ID (which also includes original ssa ID)
		// this is because the tainter first checks database calls and then service calls
		// so their order is not the real one after parsing them here
		sort.Slice(edges, func(i, j int) bool {
			return edges[i].GetIDNumber() < edges[j].GetIDNumber()
		})
	}

	for _, edge := range edges {
		graph.AddEdge(edge)
	}

	for _, call := range ssaGraph.GetServiceCalls() {
		Parse(graph, call.GetFuncShortPath(), false, funcGraphs)
	}
}

func registerDatabaseFields(graph *AbstractCallGraph, args []*AbstractObject) {
	for _, arg := range args {
		for _, taintLst := range arg.GetPrimaryTaints() {
			for _, taint := range taintLst {
				db := graph.GetApp().GetDatabaseByName(utils.ExtractDatabaseNameFromFieldPath(taint.GetDatabasePath()))
				latestSchema := db.GetLastSchema()
				if !latestSchema.HasField(taint.GetDatabasePath()) {
					field := backends.NewField(taint.GetDatabasePath(), db, latestSchema)
					latestSchema.AddField(field)
				}
			}
		}
	}
}
