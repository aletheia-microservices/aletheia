package abstractgraphparser

import (
	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph"
	abstractgraphtainter "github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph/tainter"
	"github.com/aletheia-microservices/aletheia/internal/app/backends"
	"github.com/aletheia-microservices/aletheia/internal/utils"
)

// Graph building shared by all parsers (SSA and generic input): each parser converts its own
// objects and taints into abstract objects, and uses the functions below to add them to the graph

// visitFunction gets or creates the node of service.method (with the params built by newParams) and,
// for entrypoints, adds an edge from the client. It returns false if the node was already parsed,
// otherwise marks it as parsed so that the caller adds its returns and calls only once
func visitFunction(graph *abstractgraph.AbstractCallGraph, service string, method string, funcShortPath string, entrypoint bool, newParams func() []*abstractgraph.AbstractObject) (*abstractgraph.AbstractNode, bool) {
	// dummy node
	clientNode := graph.GetNodeByNameIfExists("client")
	if clientNode == nil {
		clientNode = abstractgraph.NewAbstractNode("client", abstractgraph.NODE_CLIENT, "", "", "", "")
		graph.AddNode("client", clientNode)
	}

	// get or create the node of the function
	node := getOrCreateServiceNode(graph, service, method, newParams)

	// build dummy edges for entrypoints
	if entrypoint {
		edge := abstractgraph.NewAbstractEdge("", funcShortPath, utils.ExtractMethodNameFromShortFunctionPath(funcShortPath), clientNode, node, common.OP_UNDEFINED, abstractgraph.EDGE_SERVICE_ENTRYPOINT)
		// the client passes one untainted argument for each param
		for _, param := range node.GetParams() {
			arg := abstractgraph.NewAbstractObject(param.GetName(), make(map[string][]*abstractgraph.AbstractTaint), make(map[string][]*abstractgraph.AbstractTrace))
			edge.AddArgument(arg)
		}
		graph.AddEdge(edge)
		graph.IncrRPCs()
	}

	// a node may be reached by multiple calls, but it is only parsed once
	if node.IsParsed() {
		return node, false
	}
	node.SetParsed()
	return node, true
}

// getOrCreateServiceNode returns the node of service.method, creating it with the params built by newParams if it does not exist yet
func getOrCreateServiceNode(graph *abstractgraph.AbstractCallGraph, service string, method string, newParams func() []*abstractgraph.AbstractObject) *abstractgraph.AbstractNode {
	name := service + "." + method
	node := graph.GetNodeByNameIfExists(name)
	if node == nil {
		node = abstractgraph.NewAbstractNode(name, abstractgraph.NODE_SERVICE, service, method, "", "")
		graph.AddNode(name, node)
		// params are only built when the node is created
		for _, param := range newParams() {
			node.AddParam(param)
		}
	}
	return node
}

// getOrCreateDatabaseNode returns the node of the schema in dbname, creating it (and registering the schema in the app) if it does not exist yet
func getOrCreateDatabaseNode(graph *abstractgraph.AbstractCallGraph, dbname string, schema string) *abstractgraph.AbstractNode {
	path := dbname + "." + schema
	node := graph.GetNodeByNameIfExists(path)
	if node == nil {
		node = abstractgraph.NewAbstractNode(path, abstractgraph.NODE_DATABASE, "", "", dbname, schema)
		graph.AddNode(path, node)

		// the database must exist in the app, but the schema is registered in it if it is new
		db := graph.GetApp().GetDatabaseByName(dbname)
		if !db.HasSchema(schema) {
			db.AddSchema(backends.NewSchema(schema, db))
		}
	}
	return node
}

// addReturns adds the returns of node, given the objects of each of its return statements
func addReturns(node *abstractgraph.AbstractNode, retsLst [][]*abstractgraph.AbstractObject) {
	// first, just use the abstract objects of the first set of returns (could be any other)
	for _, ret := range retsLst[0] {
		ret.AddToAllNames(ret.GetName())
		node.AddReturn(ret)
	}
	// then, merge taints with corresponding object in the remaining set of returns
	for _, rets := range retsLst[1:] {
		for i, ret := range rets {
			// objects are matched by their position in the return statement
			obj := node.GetReturnAt(i)
			obj.AddToAllNames(ret.GetName())

			abstractgraphtainter.MergeTaints(obj, ret.GetTaints(), nil, abstractgraphtainter.MERGE_MODE_PARSE, "", false)
			abstractgraphtainter.MergeTraces(obj, ret.GetTraces())
		}
	}
}

func addServiceCallEdge(graph *abstractgraph.AbstractCallGraph, node *abstractgraph.AbstractNode, toNode *abstractgraph.AbstractNode, t string, id string, method string, args []*abstractgraph.AbstractObject, rets []*abstractgraph.AbstractObject) {
	// build the edge from node to toNode with the arguments and returns of the call
	edge := abstractgraph.NewAbstractEdge(t, id, method, node, toNode, common.OP_UNDEFINED, abstractgraph.EDGE_SERVICE_RPC)
	for _, arg := range args {
		edge.AddArgument(arg)
	}
	for _, ret := range rets {
		edge.AddReturn(ret)
	}

	// add the edge and count the RPC
	graph.AddEdge(edge)
	graph.IncrRPCs()
}

func addDatabaseCallEdge(graph *abstractgraph.AbstractCallGraph, node *abstractgraph.AbstractNode, toNode *abstractgraph.AbstractNode, t string, id string, method string, opType common.DatabaseOperationType, args []*abstractgraph.AbstractObject) {
	// build the edge from node to the database node with the arguments of the call
	edge := abstractgraph.NewAbstractEdge(t, id, method, node, toNode, opType, abstractgraph.EDGE_DATABASE_CALL)
	for _, arg := range args {
		edge.AddArgument(arg)
	}

	// create fields if they do not exist yet
	registerDatabaseFields(graph, edge.GetArguments())

	// propagate taints to databases (forward): args (from) >>> params (to)
	for i, toParam := range toNode.GetParams() {
		fromArg := edge.GetArgumentAt(i)
		abstractgraphtainter.MergeTaints(toParam, fromArg.GetPrimaryTaints(), nil, abstractgraphtainter.MERGE_MODE_PARSE, "", false)
	}

	// add the edge and count the database access
	graph.AddEdge(edge)
	graph.IncrDBAccesses()
}

func registerDatabaseFields(graph *abstractgraph.AbstractCallGraph, args []*abstractgraph.AbstractObject) {
	for _, arg := range args {
		for _, taintLst := range arg.GetPrimaryTaints() {
			for _, taint := range taintLst {
				// each taint references a field by its path (<database>.<schema>.<field>)
				db := graph.GetApp().GetDatabaseByName(utils.ExtractDatabaseNameFromFieldPath(taint.GetDatabasePath()))
				schema := db.GetSchemaForFieldPath(taint.GetDatabasePath())
				// register the field in its schema if it is new
				if !schema.HasField(taint.GetDatabasePath()) {
					field := backends.NewField(taint.GetDatabasePath(), db, schema)
					schema.AddField(field)
				}
			}
		}
	}
}
