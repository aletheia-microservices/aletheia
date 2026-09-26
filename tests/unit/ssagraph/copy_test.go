package ssagraph_test

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/internal/analysis/service-level/ssagraph"
)

func TestSimpleCopyClonesStructureWithoutTaints(t *testing.T) {
	graph := newTestGraph()
	param := newValNode(graph, 1)
	field := newValNode(graph, 2)
	ret := newValNode(graph, 3)
	graph.AddParameter(param)
	graph.AddReturnsToLst([]*ssagraph.SSANode{ret})
	graph.CreateAndAddNewEdge(param, field, ssagraph.EDGE_FIELD, 0, "PostID")
	graph.CreateAndAddNewEdge(field, ret, ssagraph.EDGE_USAGE, 0, "")
	graph.EnableGoRoutine()

	dbCall := newTestDatabaseCall(graph, newValNode(graph, 21), common.OP_WRITE)
	param.AddDatabaseTaintIfNotExists("_obj", "posts_db.post", dbCall, false, false, "")

	cp := graph.SimpleCopy()

	if len(cp.GetNodes()) != len(graph.GetNodes()) || len(cp.GetEdges()) != len(graph.GetEdges()) {
		t.Fatalf("copy has %d nodes / %d edges, want %d / %d", len(cp.GetNodes()), len(cp.GetEdges()), len(graph.GetNodes()), len(graph.GetEdges()))
	}
	for i, node := range cp.GetNodes() {
		if node == graph.GetNodes()[i] {
			t.Fatalf("copy must not share node pointers with original")
		}
		if node.IsTainted() {
			t.Errorf("copied node %s must not carry taints", node.GetName())
		}
	}
	if !param.IsTainted() {
		t.Errorf("copying must not clear taints from the original graph")
	}

	cpParam := cp.GetParamAt(0)
	if cpParam == param || cpParam.GetName() != param.GetName() {
		t.Errorf("copy parameter must be a new node with the same name")
	}
	if cp.GetIndexOfParameter(cpParam) != 0 {
		t.Errorf("copy parameter must be found by index")
	}
	edges := cp.GetEdgesFromNode(cpParam)
	if len(edges) != 1 || edges[0].GetToNode().GetName() != field.GetName() || edges[0].GetParam() != "PostID" {
		t.Errorf("copied edges must connect the copied nodes: %v", edges)
	}
	if cp.GetReturnsLst()[0][0].GetName() != ret.GetName() || cp.GetReturnsLst()[0][0] == ret {
		t.Errorf("copy returns must be remapped to copied nodes")
	}
	if !cp.IsGoRoutine() {
		t.Errorf("copy must preserve go routine flag")
	}
	if cp.GetMethodName() != graph.GetMethodName() || cp.GetFunctionShortPath() != graph.GetFunctionShortPath() {
		t.Errorf("copy must preserve method name and function path")
	}
}

func TestSimpleCopyPreservesServiceName(t *testing.T) {
	graph := newTestGraph()
	cp := graph.SimpleCopy()
	if cp.GetService() != graph.GetService() {
		t.Errorf("copy service = %q, want %q", cp.GetService(), graph.GetService())
	}
}

func TestSimpleCopyRemapsFreeVars(t *testing.T) {
	graph := newTestGraph()
	fv := newValNode(graph, 1)
	graph.AddFreeVar(fv)

	cp := graph.SimpleCopy()
	if got := cp.GetFreeVars(); len(got) != 1 || got[0] == fv || got[0].GetName() != fv.GetName() {
		t.Errorf("copy free vars must be remapped to the copied nodes: %v", got)
	}
}

func TestCombinedGraphBookkeeping(t *testing.T) {
	caller := newTestGraph()
	callee := newTestGraph()
	call := ssagraph.NewMethodCall("StorageService.StorePost.5:int", newValNode(caller, 5), nil, nil, "helper", "fn")
	other := ssagraph.NewMethodCall("StorageService.StorePost.6:int", newValNode(caller, 6), nil, nil, "helper", "fn")

	if caller.GetCombinedGraphForMethodCallIfExists(call) != nil {
		t.Errorf("graph without combined graphs must return nil")
	}
	caller.AddCombinedGraph(callee, call)

	if got := caller.GetAllCombinedGraphs(); len(got) != 1 || got[0] != callee {
		t.Errorf("combined graphs = %v", got)
	}
	if caller.GetMethodCallForCombinedGraph(callee) != call || caller.GetCombinedGraphForMethodCallIfExists(call) != callee {
		t.Errorf("combined graph and call must be mapped to each other")
	}
	if caller.GetCombinedGraphForMethodCallIfExists(other) != nil {
		t.Errorf("other calls must not be mapped")
	}
}
