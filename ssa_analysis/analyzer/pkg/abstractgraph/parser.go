package abstractgraph

import (
	"fmt"
	"log"
	"sort"

	"analyzer/pkg/app/backends"
	"analyzer/pkg/ssagraph"
	"analyzer/pkg/utils"
)

func ssaTaintToAbstractTaint(graph *AbstractCallGraph, ssaTaintsMap map[string][]*ssagraph.SSATaint) map[string][]*AbstractTaint {
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

			abstractTaints[i] = NewAbstractTaint(ssaTaint.GetDbField(), ssaTaint.GetDbCall().GetID(), true, ssaTaint.GetDbCall().IsWrite())
		}
		abstractTaintsMap[objPath] = abstractTaints
	}
	return abstractTaintsMap
}

func Parse(graph *AbstractCallGraph, funcshortpath string, funcGraphs map[string]*ssagraph.SSAGraph) {
	// dummy node
	clientNode := graph.GetNodeByNameIfExists("client")
	if clientNode == nil {
		clientNode = NewAbstractNode("client", NODE_CLIENT, "", "", "")
		graph.AddNode("client", clientNode)
	}

	ssaGraph := funcGraphs[funcshortpath]

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
			obj := NewAbstractObject(funcParam.GetName(), ssaTaintToAbstractTaint(graph, (funcParam.GetTaints())))
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
		obj := NewAbstractObject(ret.GetValue().Type().String(), ssaTaintToAbstractTaint(graph, (ret.GetTaints())))
		node.AddReturn(obj)
		retsObjs = append(retsObjs, obj)
		fmt.Printf("\t[ABSTRACTGRAPH] [index=%d] added new return object (%s)\n", i, obj.String())
	}
	// then, merge taints with corresponding object in the remaining set of returns
	if len(retsLst) > 1 {
		for _, rets := range retsLst[1:] {
			for i, ret := range rets {
				obj := retsObjs[i]
				MergeTaints(obj, ssaTaintToAbstractTaint(graph, ret.GetTaints()), true)
				fmt.Printf("\t\t[ABSTRACTGRAPH] [index=%d] merged taints from (%s) to (%s)\n", i, ret.GetName(), obj.String())
			}
		}
	}
	// debug
	for i, obj := range node.GetReturns() {
		fmt.Printf("\t[ABSTRACTGRAPH] [index=%d] final taints for object (%s):\n%s\n", i, obj.String(), obj.TaintLongString())
	}

	// build dummy edges for entrypoints
	// FIXME: must not always happen!
	edge := NewAbstractEdge(funcshortpath, utils.ExtractMethodNameFromShortFunctionPath(funcshortpath), clientNode, node, false, EDGE_SERVICE_ENTRYPOINT)
	for _, funcParam := range ssaGraph.GetFuncParametersExceptMemberAndContext() {
		arg := NewAbstractObject(funcParam.GetName(), make(map[string][]*AbstractTaint))
		edge.AddArgument(arg)
	}
	graph.AddEdge(edge)

	var edges []*AbstractEdge

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
					param := NewAbstractObject(funcParam.GetName(), ssaTaintToAbstractTaint(graph, (funcParam.GetTaints())))
					toNode.AddParam(param)
					fmt.Printf("[debug] (2) added param (%s) to node (%s)\n", param.String(), toNode.String())
				}
			}

			edge := NewAbstractEdge(call.GetID(), call.GetMethod(), node, toNode, false, EDGE_SERVICE_RPC)

			// create call arguments
			for _, callArg := range call.GetArguments() {
				arg := NewAbstractObject(callArg.GetName(), ssaTaintToAbstractTaint(graph, (callArg.GetTaints())))
				edge.AddArgument(arg)
			}

			// create call returns
			for _, callRet := range call.GetReturns() {
				ret := NewAbstractObject(callRet.GetName(), ssaTaintToAbstractTaint(graph, (callRet.GetTaints())))
				fmt.Printf("[ABSTRACTGRAPH] [%s] added return object (%s) with taints: %v\n", node.String(), ret.String(), callRet.GetTaints())
				edge.AddReturn(ret)
			}

			edges = append(edges, edge)
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

			edge := NewAbstractEdge(call.GetID(), call.GetMethod(), node, toNode, call.IsWrite(), EDGE_DATABASE_CALL)

			for _, callArg := range call.GetArguments() {
				arg := NewAbstractObject(callArg.GetName(), ssaTaintToAbstractTaint(graph, (callArg.GetTaints())))
				edge.AddArgument(arg)
			}

			// create fields if they do not exist yet
			registerDatabaseFields(graph, edge.GetArguments())

			// propagate taints to databases (forward): args (from) >>> params (to)
			taintMapping := NewTaintMapping()
			for i, toParam := range toNode.GetParams() {
				fromArg := edge.GetArgumentAt(i)
				taintMappingTmp := MergeTaints(toParam, fromArg.GetPrimaryTaints(), true)
				taintMapping.Merge(taintMappingTmp)
			}
			propagateNewTaintsToDatabases(graph, taintMapping)

			edges = append(edges, edge)
		}
		fmt.Println()

		// at the end, we need to sort edges by ID (which also includes original ssa ID)
		// this is because the tainter first checks database calls and then service calls
		// so their order is not the real one after parsing them here
		sort.Slice(edges, func(i, j int) bool {
			return edges[i].GetIDNumber() < edges[j].GetIDNumber()
		})

		for _, edge := range edges {
			graph.AddEdge(edge)
		}

		for _, call := range ssaGraph.GetServiceCalls() {
			Parse(graph, call.GetFuncShortPath(), funcGraphs)
		}
	}
}

func registerDatabaseFields(graph *AbstractCallGraph, args []*AbstractObject) {
	for _, arg := range args {
		for _, taintLst := range arg.GetPrimaryTaints() {
			for _, taint := range taintLst {
				db := graph.GetApp().GetDatabaseByName(utils.ExtractDatabaseNameFromFieldPath(taint.GetField()))
				if !db.GetSchema().HasField(taint.GetField()) {
					field := backends.NewField(taint.GetField(), db)
					db.GetSchema().AddField(field)
				}
			}
		}
	}
}

func propagateNewTaintsToDatabases(graph *AbstractCallGraph, taintMapping *TaintMapping) {
	for currTaint, otherTaintsLst := range taintMapping.mapping {
		currDb := graph.GetApp().GetDatabaseByName(utils.ExtractDatabaseNameFromFieldPath(currTaint.GetField()))
		currField := currDb.GetSchema().GetOrCreateField(currDb, currTaint.GetField())

		for _, otherTaint := range otherTaintsLst {
			otherDb := graph.GetApp().GetDatabaseByName(utils.ExtractDatabaseNameFromFieldPath(otherTaint.GetField()))
			otherField := otherDb.GetSchema().GetOrCreateField(otherDb, otherTaint.GetField())

			if currTaint.IsWrite() && otherTaint.IsWrite() {
				if currField.HasConstraintForeignKeyToField(otherField) {
					continue
				}
				constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, currField, otherField)
				currField.AddConstraint(constraint)
				currDb.GetSchema().AddConstraint(constraint)
				fmt.Printf("\t\t[ABSTRACTGRAPH] [WRITE] added new constraint: %s\n", constraint)
			} else if !currTaint.IsWrite() && !otherTaint.IsWrite() {
				if otherField.HasConstraintForeignKeyToField(currField) {
					continue
				}
				constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, otherField, currField)
				otherField.AddConstraint(constraint)
				otherDb.GetSchema().AddConstraint(constraint)
				fmt.Printf("\t\t[ABSTRACTGRAPH] [READ] added new constraint: %s\n", constraint)
			} else {
				// TODO
				log.Fatalf("\t\t[ABSTRACTGRAPH] unexpected taint mapping with write and read taints:\nCURR TAINT: %s\nOTHER TAINT:%s", currTaint.String(), otherTaint.String())
			}
		}
	}
}
