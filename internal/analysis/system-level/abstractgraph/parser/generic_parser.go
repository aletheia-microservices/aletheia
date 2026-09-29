package abstractgraphparser

import (
	"fmt"

	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph"
	abstractgraphinput "github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph/input"
)

// ParseFile loads an input model from a YAML file and parses it into graph (see ParseInput)
func ParseFile(graph *abstractgraph.AbstractCallGraph, path string) error {
	model, err := abstractgraphinput.LoadInputModel(path)
	if err != nil {
		return err
	}
	return ParseInput(graph, model)
}

// ParseInput is the language-independent equivalent of Parse: it builds graph from the functions
// of model, starting from each of its entrypoints. The graph is not mutated if model is invalid
func ParseInput(graph *abstractgraph.AbstractCallGraph, model *abstractgraphinput.InputModel) error {
	if graph == nil || graph.GetApp() == nil {
		return fmt.Errorf("abstract graph must belong to an app")
	}
	// validate the model and index its functions by func_short_path and calls by call_id
	index, err := model.Index()
	if err != nil {
		return err
	}
	// check external database references before mutating the graph
	for _, fn := range model.Functions {
		for _, call := range fn.Calls {
			if db := call.DatabaseCall; db != nil && !graph.GetApp().HasDatabase(db.Database) {
				return fmt.Errorf("function %q: call %q: database %q not found", fn.FuncShortPath, call.CallID, db.Database)
			}
		}
	}
	// build the graph starting from each entrypoint
	for _, entrypoint := range model.Entrypoints {
		parseInputFunction(graph, index, entrypoint, true)
	}
	return nil
}

func inputObject(graph *abstractgraph.AbstractCallGraph, node *abstractgraphinput.Node, calls map[string]*abstractgraphinput.Call) *abstractgraph.AbstractObject {
	taints := make(map[string][]*abstractgraph.AbstractTaint)
	traces := make(map[string][]*abstractgraph.AbstractTrace)
	for path, list := range node.Taints {
		for _, taint := range list {
			// timestamp of the call that originated the taint, prefixed with the caller's timestamp (if any)
			call := calls[taint.CallID]
			ts := call.CallTS
			if taint.CallerT != "" {
				ts = taint.CallerT + "." + ts
			}
			// database taints become abstract taints (and the database node is created if needed)
			if taint.TaintType == abstractgraphinput.TaintDatabase {
				getOrCreateDatabaseNode(graph, call.DatabaseCall.Database, call.DatabaseCall.Schema)
				op, _ := abstractgraphinput.OperationType(call.DatabaseCall.OperationType)
				var read abstractgraphinput.DatabaseTaint
				if taint.DatabaseTaint != nil {
					read = *taint.DatabaseTaint
				}
				taints[path] = append(taints[path], abstractgraph.NewAbstractTaint(ts, taint.Path, taint.CallID, op, true, false, read.ReadKey, read.ReadValue))
			} else {
				// service taints become abstract traces
				traces[path] = append(traces[path], abstractgraph.NewAbstractTrace(ts, taint.Path, taint.CallID))
			}
		}
	}
	return abstractgraph.NewAbstractObject(node.Name, taints, traces)
}

func inputObjects(graph *abstractgraph.AbstractCallGraph, nodes []*abstractgraphinput.Node, calls map[string]*abstractgraphinput.Call) []*abstractgraph.AbstractObject {
	var objs []*abstractgraph.AbstractObject
	for _, node := range nodes {
		objs = append(objs, inputObject(graph, node, calls))
	}
	return objs
}

func inputParams(graph *abstractgraph.AbstractCallGraph, index *abstractgraphinput.Index, funcShortPath string) func() []*abstractgraph.AbstractObject {
	return func() []*abstractgraph.AbstractObject {
		return inputObjects(graph, index.Functions[funcShortPath].Params, index.Calls[funcShortPath])
	}
}

func parseInputFunction(graph *abstractgraph.AbstractCallGraph, index *abstractgraphinput.Index, funcShortPath string, entrypoint bool) {
	fn := index.Functions[funcShortPath]
	calls := index.Calls[funcShortPath]
	// get or create the node of the function, and skip it if it was already parsed
	node, ok := visitFunction(graph, fn.Service, fn.Method, funcShortPath, entrypoint, inputParams(graph, index, funcShortPath))
	if !ok {
		return
	}

	// add the returns of the function (returns are optional in input models), merging the objects of all its return statements
	if len(fn.Returns) > 0 {
		var retsLst [][]*abstractgraph.AbstractObject
		for _, rets := range fn.Returns {
			retsLst = append(retsLst, inputObjects(graph, rets, calls))
		}
		addReturns(node, retsLst)
	}

	// add an edge for each call of the function
	for _, call := range fn.Calls {
		parseInputCall(graph, index, node, call, calls)
	}
	// then, recursively parse the functions called by RPC
	for _, call := range fn.Calls {
		if call.ServiceCall != nil {
			parseInputFunction(graph, index, call.ServiceCall.FuncShortPath, false)
		}
	}
}

func parseInputCall(graph *abstractgraph.AbstractCallGraph, index *abstractgraphinput.Index, node *abstractgraph.AbstractNode, call *abstractgraphinput.Call, calls map[string]*abstractgraphinput.Call) {
	// parse if service call
	if svc := call.ServiceCall; svc != nil {
		toNode := getOrCreateServiceNode(graph, svc.Service, svc.Method, inputParams(graph, index, svc.FuncShortPath))
		args := inputObjects(graph, call.Arguments, calls)
		rets := inputObjects(graph, svc.Returns, calls)
		addServiceCallEdge(graph, node, toNode, call.CallTS, call.CallID, svc.Method, args, rets)
		return
	}
	// or parse if database call
	db := call.DatabaseCall
	toNode := getOrCreateDatabaseNode(graph, db.Database, db.Schema)
	op, _ := abstractgraphinput.OperationType(db.OperationType)
	args := inputObjects(graph, call.Arguments, calls)
	addDatabaseCallEdge(graph, node, toNode, call.CallTS, call.CallID, db.Method, op, args)
}
