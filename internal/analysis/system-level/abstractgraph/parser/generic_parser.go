package abstractgraphparser

import (
	"fmt"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph"
	abstractgraphinput "github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph/input"
	abstractgraphtainter "github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph/tainter"
	"github.com/aletheia-microservices/aletheia/internal/app/backends"
)

// ParseFile loads calls from YAML and adds their edges from source to graph.
// Service definitions and entrypoint edges can be supplied by the caller independently.
func ParseFile(graph *abstractgraph.AbstractCallGraph, source *abstractgraph.AbstractNode, path string) error {
	model, err := abstractgraphinput.LoadInputModel(path)
	if err != nil {
		return err
	}
	return ParseInput(graph, source, model)
}

// ParseInput builds call edges without depending on SSA or a source language.
// source must already belong to graph. All taint call references must occur in model.
func ParseInput(graph *abstractgraph.AbstractCallGraph, source *abstractgraph.AbstractNode, model *abstractgraphinput.InputModel) error {
	if graph == nil || source == nil || graph.GetNodeByNameIfExists(source.GetName()) != source {
		return fmt.Errorf("source must belong to the abstract graph")
	}
	calls, err := model.Index()
	if err != nil {
		return err
	}
	// Check external database references before mutating the graph.
	for _, call := range model.Calls {
		if db := call.DatabaseCall; db != nil {
			if graph.GetApp() == nil || !graph.GetApp().HasDatabase(db.Database) {
				return fmt.Errorf("call %q: database %q not found", call.CallID, db.Database)
			}
		}
	}
	for _, call := range model.Calls {
		parseInputCall(graph, source, call, calls)
	}
	return nil
}

func inputDatabaseNode(graph *abstractgraph.AbstractCallGraph, db *abstractgraphinput.DatabaseCall) *abstractgraph.AbstractNode {
	path := db.Database + "." + db.Schema
	node := graph.GetNodeByNameIfExists(path)
	if node == nil {
		node = abstractgraph.NewAbstractNode(path, abstractgraph.NODE_DATABASE, "", "", db.Database, db.Schema)
		graph.AddNode(path, node)
	}
	database := graph.GetApp().GetDatabaseByName(db.Database)
	if !database.HasSchema(db.Schema) {
		database.AddSchema(backends.NewSchema(db.Schema, database))
	}
	return node
}

func inputObject(graph *abstractgraph.AbstractCallGraph, node *abstractgraphinput.Node, calls map[string]*abstractgraphinput.Call) *abstractgraph.AbstractObject {
	taints := make(map[string][]*abstractgraph.AbstractTaint)
	traces := make(map[string][]*abstractgraph.AbstractTrace)
	for path, list := range node.Taints {
		for _, taint := range list {
			call := calls[taint.CallID]
			ts := call.CallTS
			if taint.CallerT != "" {
				ts = taint.CallerT + "." + ts
			}
			if taint.TaintType == abstractgraphinput.TaintDatabase {
				inputDatabaseNode(graph, call.DatabaseCall)
				op, _ := abstractgraphinput.OperationType(call.DatabaseCall.OperationType)
				taints[path] = append(taints[path], abstractgraph.NewAbstractTaint(ts, taint.Path, taint.CallID, op, true, false, taint.DatabaseTaint.ReadKey, taint.DatabaseTaint.ReadValue))
			} else {
				traces[path] = append(traces[path], abstractgraph.NewAbstractTrace(ts, taint.Path, taint.CallID))
			}
		}
	}
	return abstractgraph.NewAbstractObject(node.Name, taints, traces)
}

func parseInputCall(graph *abstractgraph.AbstractCallGraph, source *abstractgraph.AbstractNode, call *abstractgraphinput.Call, calls map[string]*abstractgraphinput.Call) {
	var edge *abstractgraph.AbstractEdge
	if svc := call.ServiceCall; svc != nil {
		name := svc.Service + "." + svc.Method
		target := graph.GetNodeByNameIfExists(name)
		if target == nil {
			target = abstractgraph.NewAbstractNode(name, abstractgraph.NODE_SERVICE, svc.Service, svc.Method, "", "")
			graph.AddNode(name, target)
		}
		edge = abstractgraph.NewAbstractEdge(call.CallTS, call.CallID, svc.Method, source, target, common.OP_UNDEFINED, abstractgraph.EDGE_SERVICE_RPC)
		for _, ret := range svc.Returns {
			edge.AddReturn(inputObject(graph, ret, calls))
		}
		graph.IncrRPCs()
	} else {
		db := call.DatabaseCall
		target := inputDatabaseNode(graph, db)
		op, _ := abstractgraphinput.OperationType(db.OperationType)
		edge = abstractgraph.NewAbstractEdge(call.CallTS, call.CallID, db.Method, source, target, op, abstractgraph.EDGE_DATABASE_CALL)
		graph.IncrDBAccesses()
	}
	for _, arg := range call.Arguments {
		edge.AddArgument(inputObject(graph, arg, calls))
	}
	if call.CallType == abstractgraphinput.CallTypeDB {
		registerDatabaseFields(graph, edge.GetArguments())
		for i, param := range edge.GetToNode().GetParams() {
			if i < len(edge.GetArguments()) {
				abstractgraphtainter.MergeTaints(param, edge.GetArgumentAt(i).GetPrimaryTaints(), nil, abstractgraphtainter.MERGE_MODE_PARSE, "", false)
			}
		}
	}
	graph.AddEdge(edge)
}
