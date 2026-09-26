package ssagraph_test

import (
	"testing"

	"analyzer/pkg/analysis/common"
	"analyzer/pkg/analysis/service-level/ssagraph"
)

func TestCallIdentifiers(t *testing.T) {
	graph := newTestGraph()
	node := newValNode(graph, 21)

	dbCall := newTestDatabaseCall(graph, node, common.OP_WRITE)
	if got, want := dbCall.GetID(), "StorageService.StorePost.21:int"; got != want {
		t.Errorf("database call id = %q, want %q", got, want)
	}
	if got, want := dbCall.GetT(), node.GetName(); got != want {
		t.Errorf("database call t = %q, want %q", got, want)
	}
	if got, want := dbCall.GetDatabasePath(), "posts_db.post"; got != want {
		t.Errorf("database path = %q, want %q", got, want)
	}
	if got, want := dbCall.String(), "posts_db.post.InsertOne(...)"; got != want {
		t.Errorf("database call string = %q, want %q", got, want)
	}

	svcCall := ssagraph.NewServiceCall("UploadService.UploadPost.t3", node, nil, nil, "StorageService", "StorePost", "postnotification.StorageService.StorePost")
	if got, want := svcCall.GetServiceWithMethod(), "StorageService.StorePost"; got != want {
		t.Errorf("service call = %q, want %q", got, want)
	}

	// every call kind is accessible through the ssagraph.Call interface
	for _, call := range []ssagraph.Call{dbCall, svcCall, ssagraph.NewMethodCall("id", node, nil, nil, "m", "pkg.m")} {
		graph.AddCall(call)
	}
	if len(graph.GetAllCalls()) != 3 {
		t.Fatalf("graph has %d calls, want 3", len(graph.GetAllCalls()))
	}
}

func TestHasMethodCallMatchesByID(t *testing.T) {
	graph := newTestGraph()
	node := newValNode(graph, 1)
	graph.AddMethodCall(ssagraph.NewMethodCallGoRoutine("pkg.Service.Method$1", node.GetID(), node, nil, nil, nil, "Method$1", "pkg.Service.Method$1"))

	if !graph.HasMethodCall("pkg.Service.Method$1") {
		t.Errorf("HasMethodCall must find registered go routine by id")
	}
	if graph.HasMethodCall("pkg.Service.Other") {
		t.Errorf("HasMethodCall must not match other ids")
	}
}

func TestMethodCallAccessors(t *testing.T) {
	graph := newTestGraph()
	node := newValNode(graph, 5)
	recv, arg := newValNode(graph, 1), newValNode(graph, 2)
	ret := newValNode(graph, 3)
	call := ssagraph.NewMethodCall("StorageService.StorePost.5:int", node, []*ssagraph.SSANode{recv, arg}, []*ssagraph.SSANode{ret}, "helper", "postnotification.StorageService.helper")

	if call.GetArgumentAt(1) != arg || call.GetReturnAt(0) != ret || call.TryGetReturnAt(0) != ret {
		t.Errorf("arguments = %v, returns = %v", call.GetArguments(), call.GetReturns())
	}
	if call.TryGetReturnAt(1) != nil {
		t.Errorf("TryGetReturnAt out of range must return nil")
	}
	if call.GetT() != node.GetName() || call.GetMethod() != "helper" || call.GetNode() != node {
		t.Errorf("t = %q, method = %q", call.GetT(), call.GetMethod())
	}
	if call.String() != "postnotification.StorageService.helper" || call.GetFuncShortPath() != call.String() {
		t.Errorf("String() = %q, want the function path", call.String())
	}
}

func TestGoRoutineMethodCall(t *testing.T) {
	graph := newTestGraph()
	node := newValNode(graph, 5)
	bind := newValNode(graph, 1)
	call := ssagraph.NewMethodCallGoRoutine("postnotification.Spawn$1", "val_go_12", node, []*ssagraph.SSANode{bind}, nil, nil, "Spawn$1", "postnotification.Spawn$1")

	// go routines are ordered by the id of the go instruction
	if call.GetT() != "val_go_12" || call.GetID() != "postnotification.Spawn$1" {
		t.Errorf("t = %q, id = %q", call.GetT(), call.GetID())
	}
	if call.GetBindAt(0) != bind {
		t.Errorf("binding = %v, want the captured variable", call.GetBindAt(0))
	}
}

func TestServiceCallAccessors(t *testing.T) {
	graph := newTestGraph()
	node := newValNode(graph, 3)
	arg, ret := newValNode(graph, 1), newValNode(graph, 2)
	call := ssagraph.NewServiceCall("UploadService.UploadPost.3:int", node, []*ssagraph.SSANode{arg}, []*ssagraph.SSANode{ret}, "StorageService", "StorePost", "postnotification.StorageService.StorePost")

	if call.GetService() != "StorageService" || call.GetMethod() != "StorePost" || call.String() != "StorageService.StorePost" {
		t.Errorf("service = %q, method = %q, String() = %q", call.GetService(), call.GetMethod(), call.String())
	}
	if call.GetFuncShortPath() != "postnotification.StorageService.StorePost" {
		t.Errorf("function path = %q", call.GetFuncShortPath())
	}
	if len(call.GetArguments()) != 1 || call.GetArguments()[0] != arg || len(call.GetReturns()) != 1 || call.GetReturns()[0] != ret {
		t.Errorf("arguments = %v, returns = %v", call.GetArguments(), call.GetReturns())
	}
	if call.GetT() != node.GetName() {
		t.Errorf("t = %q, want the call value name", call.GetT())
	}
}

func TestDatabaseCallAccessors(t *testing.T) {
	graph := newTestGraph()
	arg := newValNode(graph, 1)
	call := ssagraph.NewDatabaseCall("StorageService.StorePost.21:int", newValNode(graph, 21), []*ssagraph.SSANode{arg}, "posts_db", "post", "InsertOne", common.OP_WRITE)

	if call.GetDatabaseName() != "posts_db" || call.GetSchemaName() != "post" || call.GetDatabasePath() != "posts_db.post" {
		t.Errorf("database = %q, schema = %q", call.GetDatabaseName(), call.GetSchemaName())
	}
	if call.GetOpType() != common.OP_WRITE || call.GetMethod() != "InsertOne" {
		t.Errorf("op = %s, method = %q", common.OperationTypeToString(call.GetOpType()), call.GetMethod())
	}
	if len(call.GetArguments()) != 1 || call.GetArguments()[0] != arg {
		t.Errorf("arguments = %v", call.GetArguments())
	}
}
