package abstractgraph

import (
	"fmt"
	"log"
	"os"

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
				fmt.Printf("[SSA TO ABSTRACT TAINT] new taint on object path (%s): %s\n", objPath, taint.LongString())
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
				fmt.Printf("[SSA TO ABSTRACT TRACE] new trace on object path (%s): %s\n", objPath, trace.LongString())
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
	fmt.Printf("[ABSTRACTGRAPH] got ssa graph for (%s): %v\n", funcshortpath, ssaGraph)

	name := ssaGraph.GetServiceWithMethod()
	node := graph.GetNodeByNameIfExists(name)

	if node != nil && node.IsParsed() {
		fmt.Printf("[ABSTRACTGRAPH] ignoring node already visited: %s\n", node.String())
		return
	}

	if node == nil {
		node = NewAbstractNode(name, NODE_SERVICE, ssaGraph.GetService(), ssaGraph.GetMethodName(), "")
		graph.AddNode(name, node)

		fmt.Printf("[ABSTRACTGRAPH] creating node with (%d) params: %s\n", len(ssaGraph.GetFuncParametersExceptMemberAndContext()), node)
		for _, funcParam := range ssaGraph.GetFuncParametersExceptMemberAndContext() {
			obj := NewAbstractObject(funcParam.GetName(), ssaTaintDatabaseToAbstractTaint(graph, funcParam.GetTaints()), ssaTaintServiceToAbstractTrace(graph, funcParam.GetTaints()))
			fmt.Printf("[debug] (1) added param (%s) to node (%s)\n", obj.String(), node.String())
			node.AddParam(obj)
		}
	}

	node.SetParsed()

	// finalize parsing
	fmt.Printf("[ABSTRACTGRAPH] parsing returns for node: %s\n", node.String())
	retsLst := ssaGraph.GetReturnsLst()
	var retsObjs []*AbstractObject
	// first, just create new abstract objects using the first set of returns (could be any other)
	for i, ret := range retsLst[0] {
		obj := NewAbstractObject(ret.GetValue().Type().String(), ssaTaintDatabaseToAbstractTaint(graph, ret.GetTaints()), ssaTaintServiceToAbstractTrace(graph, ret.GetTaints()))
		node.AddReturn(obj)
		retsObjs = append(retsObjs, obj)
		fmt.Printf("\t[ABSTRACTGRAPH] [index=%d] added new return object (%s)\n", i, obj.String())
	}
	// then, merge taints with corresponding object in the remaining set of returns
	if len(retsLst) > 1 {
		for _, rets := range retsLst[1:] {
			for i, ret := range rets {
				obj := retsObjs[i]
				MergeTaints(obj, ssaTaintDatabaseToAbstractTaint(graph, ret.GetTaints()), nil, true, false)
				fmt.Printf("\t\t[ABSTRACTGRAPH] [index=%d] merged taints from (%s) to (%s)\n", i, ret.GetName(), obj.String())
			}
		}
	}
	// debug
	/* for i, obj := range node.GetReturns() {
		fmt.Printf("\t[ABSTRACTGRAPH] [index=%d] final taints for object (%s):\n%s\n", i, obj.String(), obj.TaintLongString())
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

	for _, call := range ssaGraph.GetAllCalls() {
		if serviceCall, ok := call.(*ssagraph.ServiceCall); ok {
			fmt.Printf("[ABSTRACTGRAPH] [%s] found function (%s) with service calls\n", ssaGraph.GetService(), funcshortpath)
			toName := serviceCall.GetServiceWithMethod()
			toNode := graph.GetNodeByNameIfExists(toName)

			toSSAGraph := funcGraphs[serviceCall.GetFuncShortPath()]
			if toSSAGraph == nil {
				log.Fatalf("could not find ssa graph for short func path (%s)", serviceCall.GetFuncShortPath())
			}

			// create node for the first time
			if toNode == nil {
				toNode = NewAbstractNode(toName, NODE_SERVICE, serviceCall.GetService(), serviceCall.GetMethod(), "")
				graph.AddNode(toName, toNode)

				fmt.Printf("[ABSTRACTGRAPH] creating toNode with (%d) params: %s\n", len(toSSAGraph.GetFuncParametersExceptMemberAndContext()), toNode)
				for _, funcParam := range toSSAGraph.GetFuncParametersExceptMemberAndContext() {
					param := NewAbstractObject(funcParam.GetName(), ssaTaintDatabaseToAbstractTaint(graph, funcParam.GetTaints()), ssaTaintServiceToAbstractTrace(graph, funcParam.GetTaints()))
					toNode.AddParam(param)
					fmt.Printf("[debug] (2) added param (%s) to node (%s)\n", param.String(), toNode.String())
				}
			}

			edge := NewAbstractEdge(serviceCall.GetID(), serviceCall.GetMethod(), node, toNode, common.OP_UNDEFINED, EDGE_SERVICE_RPC)

			// create call arguments
			for _, callArg := range serviceCall.GetArguments() {
				arg := NewAbstractObject(callArg.GetName(), ssaTaintDatabaseToAbstractTaint(graph, callArg.GetTaints()), ssaTaintServiceToAbstractTrace(graph, callArg.GetTaints()))
				edge.AddArgument(arg)
			}

			// create call returns
			for _, callRet := range serviceCall.GetReturns() {
				ret := NewAbstractObject(callRet.GetName(), ssaTaintDatabaseToAbstractTaint(graph, callRet.GetTaints()), ssaTaintServiceToAbstractTrace(graph, callRet.GetTaints()))
				fmt.Printf("[ABSTRACTGRAPH] [%s] added return object (%s) with taints: %v\n", node.String(), ret.String(), callRet.GetTaints())
				edge.AddReturn(ret)
			}

			edges = append(edges, edge)
			fmt.Printf("[ABSTRACT GRAPH] [SERVICE CALL] added edge: %v\n", edge)
		}

		if databaseCall, ok := call.(*ssagraph.DatabaseCall); ok {
			fmt.Printf("[ABSTRACTGRAPH] found [%s] function (%s) with database calls\n", ssaGraph.GetService(), funcshortpath)

			toDatabasePath := databaseCall.GetDatabasePath()
			toNode := graph.GetNodeByNameIfExists(toDatabasePath)
			dbname := databaseCall.GetDatabaseName()
			if toNode == nil {
				toNode = NewAbstractNode(toDatabasePath, NODE_DATABASE, "", "", dbname)
				graph.AddNode(toDatabasePath, toNode)
				schemaName := databaseCall.GetSchemaName()

				if !graph.GetApp().HasDatabase(dbname) {
					log.Fatalf("database (%s) not found", dbname)
				}

				db := graph.GetApp().GetDatabaseByName(dbname)
				if !db.HasSchema(schemaName) {
					db.AddSchema(backends.NewSchema(schemaName))
				}
			}

			edge := NewAbstractEdge(databaseCall.GetID(), databaseCall.GetMethod(), node, toNode, databaseCall.GetOpType(), EDGE_DATABASE_CALL)

			for _, callArg := range databaseCall.GetArguments() {
				arg := NewAbstractObject(callArg.GetName(), ssaTaintDatabaseToAbstractTaint(graph, callArg.GetTaints()), ssaTaintServiceToAbstractTrace(graph, callArg.GetTaints()))
				edge.AddArgument(arg)
			}

			// create fields if they do not exist yet
			registerDatabaseFields(graph, edge.GetArguments())

			// propagate taints to databases (forward): args (from) >>> params (to)
			taintMapping := NewTaintMapping()
			for i, toParam := range toNode.GetParams() {
				fromArg := edge.GetArgumentAt(i)
				taintMappingTmp, _ := MergeTaints(toParam, fromArg.GetPrimaryTaints(), nil, true, false)
				taintMapping.Merge(taintMappingTmp, true)
			}
			PropagateNewTaintsToDatabaseSchemas(graph, -1, taintMapping)

			edges = append(edges, edge)
			fmt.Printf("[ABSTRACT GRAPH] [DATABASE CALL] added edge: %v\n", edge)
		}
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

func (graph *AbstractCallGraph) WriteVisited(appname string) {
	filename := fmt.Sprintf("output/%s/abstractcallgraph.visited", appname)

	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	clientNode := graph.GetNodeByName("client")

	fmt.Printf("starting...")
	for i, edge := range graph.GetEdgesFromNode(clientNode) {
		fmt.Fprintf(file, "\n\n%d: %s", i, edge.to.String())
		for _, edge2 := range graph.GetEdgesFromNode(edge.to) {
			fmt.Fprintf(file, "\n\t-> "+edge2.to.String())
			for _, edge3 := range graph.GetEdgesFromNode(edge2.to) {
				fmt.Fprintf(file, "\n\t\t-> "+edge3.to.String())
				for _, edge4 := range graph.GetEdgesFromNode(edge3.to) {
					fmt.Fprintf(file, "\n\t\t\t-> "+edge4.to.String())
					for _, edge5 := range graph.GetEdgesFromNode(edge4.to) {
						fmt.Fprintf(file, "\n\t\t\t\t-> "+edge5.to.String())
					}
				}
			}
		}
	}
	fmt.Printf("end...")
}
