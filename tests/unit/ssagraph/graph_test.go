package ssagraph_test

import (
	"testing"

	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/analysis/common"
	"analyzer/pkg/analysis/service-level/ssagraph"
)

func TestCreateAndAddNewEdgeDeduplicates(t *testing.T) {
	graph := newTestGraph()
	a := newValNode(graph, 1)
	b := newValNode(graph, 2)

	edge, created := graph.CreateAndAddNewEdge(a, b, ssagraph.EDGE_FIELD, 0, "PostID")
	if !created || edge == nil {
		t.Fatalf("first edge must be created")
	}
	again, created := graph.CreateAndAddNewEdge(a, b, ssagraph.EDGE_USAGE, 0, "")
	if created || again != edge {
		t.Errorf("edge between the same nodes must be reused")
	}
	if _, created := graph.CreateAndAddNewEdge(nil, b, ssagraph.EDGE_USAGE, 0, ""); created {
		t.Errorf("edge from nil node must not be created")
	}
	if _, created := graph.CreateAndAddNewEdge(b, a, ssagraph.EDGE_USAGE, 0, ""); !created {
		t.Errorf("edge in the opposite direction must be created")
	}

	if got := len(graph.GetEdges()); got != 2 {
		t.Fatalf("graph has %d edges, want 2", got)
	}
	if got := graph.GetEdgesTypedFrom(a, ssagraph.EDGE_FIELD); len(got) != 1 || got[0].GetParam() != "PostID" {
		t.Errorf("GetEdgesTypedFrom(a, FIELD) = %v", got)
	}
	if got := graph.GetEdgesTypedTo(b, ssagraph.EDGE_USAGE); len(got) != 0 {
		t.Errorf("GetEdgesTypedTo(b, USAGE) = %v, want none", got)
	}
	if got := len(graph.GetAllNodeEdges(a)); got != 2 {
		t.Errorf("node a has %d edges, want 2", got)
	}
}

func TestRegisterNewNodeInstrAddsNodeOnce(t *testing.T) {
	graph := newTestGraph()
	ssagraph.RegisterNewNodeInstr(graph, &ssa.Jump{}, "jump")
	if got := len(graph.GetNodes()); got != 1 {
		t.Errorf("graph has %d nodes, want 1", got)
	}
}

func TestGraphIdentity(t *testing.T) {
	graph := newTestGraph()
	if graph.GetServiceWithMethod() != "StorageService.StorePost" || graph.GetService() != "StorageService" || graph.GetMethodName() != "StorePost" {
		t.Errorf("service = %q, method = %q", graph.GetService(), graph.GetMethodName())
	}
	if graph.GetPackageName() != "postnotification" || graph.GetFunctionShortPath() != "postnotification.StorageService.StorePost" || graph.String() != graph.GetFunctionShortPath() {
		t.Errorf("package = %q, function path = %q", graph.GetPackageName(), graph.GetFunctionShortPath())
	}
	if graph.IsGoRoutine() {
		t.Errorf("new graph must not be a go routine")
	}
	graph.EnableGoRoutine()
	if !graph.IsGoRoutine() {
		t.Errorf("EnableGoRoutine must mark the graph as a go routine")
	}
}

func TestGraphNodeLookup(t *testing.T) {
	graph := newTestGraph()
	node := newValNode(graph, 1)

	if graph.GetNodeByName(node.GetName()) != node {
		t.Errorf("GetNodeByName must find registered nodes")
	}
	if _, ok := graph.GetNodeByNameIfExists("missing"); ok {
		t.Errorf("missing node must not be found")
	}
}

func TestGraphParamsReturnsAndFreeVars(t *testing.T) {
	graph := newTestGraph()
	s, ctx, reqID, text := newValNode(graph, 1), newValNode(graph, 2), newValNode(graph, 3), newValNode(graph, 4)
	for _, param := range []*ssagraph.SSANode{s, ctx, reqID, text} {
		graph.AddParameter(param)
	}

	if graph.GetParamAt(3) != text || graph.GetIndexOfParameter(reqID) != 2 {
		t.Errorf("params = %v", graph.GetParams())
	}
	// service methods take (receiver, context, ...)
	if got := graph.GetFuncParametersExceptMemberAndContext(); len(got) != 2 || got[0] != reqID || got[1] != text {
		t.Errorf("params except receiver and context = %v, want [reqID text]", got)
	}
	onlyContext := newTestGraph()
	onlyContext.AddParameter(newValNode(onlyContext, 1))
	onlyContext.AddParameter(newValNode(onlyContext, 2))
	if got := onlyContext.GetFuncParametersExceptMemberAndContext(); got != nil {
		t.Errorf("params except receiver and context = %v, want none", got)
	}

	graph.AddReturnsToLst([]*ssagraph.SSANode{reqID})
	graph.AddReturnsToLst([]*ssagraph.SSANode{text})
	if got := graph.GetReturnsLst(); len(got) != 2 || got[1][0] != text {
		t.Errorf("returns = %v", got)
	}
	graph.AddFreeVar(ctx)
	if got := graph.GetFreeVars(); len(got) != 1 || got[0] != ctx {
		t.Errorf("free vars = %v", got)
	}
}

func TestGraphCallLists(t *testing.T) {
	graph := newTestGraph()
	dbCall := newTestDatabaseCall(graph, newValNode(graph, 21), common.OP_WRITE)
	svcCall := ssagraph.NewServiceCall("id", newValNode(graph, 3), nil, nil, "StorageService", "ReadPost", "fn")
	methodCall := ssagraph.NewMethodCall("StorageService.StorePost.5:int", newValNode(graph, 5), nil, nil, "helper", "fn")

	if graph.HasDatabaseCalls() {
		t.Errorf("new graph must not have database calls")
	}
	graph.AddDatabaseCall(dbCall)
	graph.AddServiceCall(svcCall)
	graph.AddMethodCall(methodCall)
	for _, call := range []ssagraph.Call{svcCall, dbCall, methodCall} {
		graph.AddCall(call)
	}

	if !graph.HasDatabaseCalls() || len(graph.GetDatabaseCalls()) != 1 || len(graph.GetServiceCalls()) != 1 || len(graph.GetMethodCalls()) != 1 {
		t.Errorf("calls per kind are wrong")
	}
	// all calls keep the order in which they were added (program order)
	if all := graph.GetAllCalls(); len(all) != 3 || all[0] != svcCall || all[1] != dbCall || all[2] != methodCall {
		t.Errorf("all calls = %v", all)
	}
}

func TestGraphSortOrdersNodesByID(t *testing.T) {
	graph := newTestGraph()
	b := ssagraph.RegisterNewNodeVal(graph, nil, newValNode(newTestGraph(), 1).GetValue(), "val_b")
	a := ssagraph.RegisterNewNodeVal(graph, nil, newValNode(newTestGraph(), 2).GetValue(), "val_a")

	graph.Sort()
	if nodes := graph.GetNodes(); nodes[0] != a || nodes[1] != b {
		t.Errorf("nodes = [%s %s], want sorted by id", nodes[0].GetID(), nodes[1].GetID())
	}
}

func TestGraphEdgeQueries(t *testing.T) {
	graph := newTestGraph()
	a, b, c := newValNode(graph, 1), newValNode(graph, 2), newValNode(graph, 3)
	ab, _ := graph.CreateAndAddNewEdge(a, b, ssagraph.EDGE_FIELD, 0, "ID")
	ac, _ := graph.CreateAndAddNewEdge(a, c, ssagraph.EDGE_FIELD, 0, "Text")
	bc, _ := graph.CreateAndAddNewEdge(b, c, ssagraph.EDGE_USAGE, 0, "")

	if got := graph.GetFirstEdgeTypedFrom(a, ssagraph.EDGE_FIELD); got != ab {
		t.Errorf("first field edge from a = %v, want a -> b", got)
	}
	if graph.GetFirstEdgeTypedFrom(a, ssagraph.EDGE_LOAD) != nil {
		t.Errorf("missing edge type must return nil")
	}
	if got := graph.GetFirstEdgeToNode(c); got != ac {
		t.Errorf("first edge to c = %v, want a -> c", got)
	}
	if got := graph.GetEdgesToNode(c); len(got) != 2 || got[1] != bc {
		t.Errorf("edges to c = %v", got)
	}
}

func TestGraphRelease(t *testing.T) {
	graph := newTestGraph()
	a, b := newValNode(graph, 1), newValNode(graph, 2)
	graph.CreateAndAddNewEdge(a, b, ssagraph.EDGE_USAGE, 0, "")
	graph.AddParameter(a)
	graph.AddDatabaseCall(newTestDatabaseCall(graph, b, common.OP_WRITE))

	graph.Release()
	if len(graph.GetNodes()) != 0 || len(graph.GetEdges()) != 0 || len(graph.GetParams()) != 0 || graph.HasDatabaseCalls() {
		t.Errorf("Release must drop nodes, edges, parameters and calls")
	}
	if graph.GetFunctionShortPath() == "" {
		t.Errorf("Release must keep the graph identity")
	}
}
