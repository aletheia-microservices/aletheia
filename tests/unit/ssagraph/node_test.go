package ssagraph_test

import (
	"testing"

	"analyzer/pkg/analysis/common"
	"analyzer/pkg/analysis/service-level/ssagraph"
)

func TestRegisterNewNodeValAddsDefinition(t *testing.T) {
	graph := newTestGraph()
	node := newValNode(graph, 1)

	got, ok := graph.GetNodeByNameIfExists(node.GetName())
	if !ok || got != node {
		t.Fatalf("node %q not registered in graph defs", node.GetName())
	}
	if len(graph.GetNodes()) != 1 {
		t.Fatalf("graph has %d nodes, want 1", len(graph.GetNodes()))
	}
	if node.IsTainted() {
		t.Fatalf("new node must not be tainted")
	}
}

func TestAddDatabaseTaintDeduplicatesOnPathOpTypeAndReadFlags(t *testing.T) {
	graph := newTestGraph()
	node := newValNode(graph, 1)
	write := newTestDatabaseCall(graph, newValNode(graph, 2), common.OP_WRITE)
	read := newTestDatabaseCall(graph, newValNode(graph, 3), common.OP_READ)

	if !node.AddDatabaseTaintIfNotExists("_obj.PostID", "posts_db.post.PostID", read, true, false, "") {
		t.Fatalf("first taint must be added")
	}
	if node.AddDatabaseTaintIfNotExists("_obj.PostID", "posts_db.post.PostID", read, true, false, "") {
		t.Errorf("taint with same db path, op type and read flags must not be added twice")
	}
	// the same field can be reached as a filter key and as part of the read value
	if !node.AddDatabaseTaintIfNotExists("_obj.PostID", "posts_db.post.PostID", read, false, true, "") {
		t.Errorf("taint with same db path and op type but different read flags must be added")
	}
	if !node.AddDatabaseTaintIfNotExists("_obj.PostID", "posts_db.post.PostID", write, false, false, "") {
		t.Errorf("taint with same db path but different op type must be added")
	}
	if !node.AddDatabaseTaintIfNotExists("_obj", "posts_db.post", write, false, false, "") {
		t.Errorf("taint on a different object path must be added")
	}

	if got := len(node.GetTaintsForPath("_obj.PostID")); got != 3 {
		t.Errorf("_obj.PostID has %d taints, want 3", got)
	}
	if got := len(node.GetTaints()); got != 2 {
		t.Errorf("node has %d object paths, want 2", got)
	}
	if node.GetTaintsForPath("_obj.Missing") != nil {
		t.Errorf("missing path must return nil")
	}
	if _, ok := node.GetTaints()["_obj.Missing"]; ok {
		t.Errorf("GetTaintsForPath must not create new entries")
	}
}

func TestAddServiceTaintDeduplicatesOnPath(t *testing.T) {
	graph := newTestGraph()
	node := newValNode(graph, 1)
	callNode := newValNode(graph, 3)
	svcCall := ssagraph.NewServiceCall("UploadService.UploadPost.t3", callNode, nil, nil, "StorageService", "StorePost", "postnotification.StorageService.StorePost")

	if !node.AddServiceTaintIfNotExists("_obj", "StorageService.StorePost.text", svcCall, "") {
		t.Fatalf("first service taint must be added")
	}
	if node.AddServiceTaintIfNotExists("_obj", "StorageService.StorePost.text", svcCall, "") {
		t.Errorf("service taint with same path must not be added twice")
	}

	// a database taint with the same path does not collide with the service taint
	dbCall := newTestDatabaseCall(graph, newValNode(graph, 4), common.OP_WRITE)
	if !node.AddDatabaseTaintIfNotExists("_obj", "StorageService.StorePost.text", dbCall, false, false, "") {
		t.Errorf("database taint must not be deduplicated against a service taint")
	}
}

func TestSSATaintAccessorsDependOnTaintType(t *testing.T) {
	graph := newTestGraph()
	node := newValNode(graph, 1)
	dbCall := newTestDatabaseCall(graph, newValNode(graph, 2), common.OP_READ)
	svcCall := ssagraph.NewServiceCall("id", newValNode(graph, 3), nil, nil, "StorageService", "ReadPost", "fn")

	node.AddDatabaseTaintIfNotExists("_obj", "posts_db.post", dbCall, false, true, "")
	node.AddServiceTaintIfNotExists("_obj", "StorageService.ReadPost.t19", svcCall, "")

	taints := node.GetTaintsForPath("_obj")
	if len(taints) != 2 {
		t.Fatalf("got %d taints, want 2", len(taints))
	}
	db, svc := taints[0], taints[1]

	if !db.IsDatabaseTaint() || db.IsServiceTaint() {
		t.Errorf("first taint must be a database taint")
	}
	if db.GetDatabasePath() != "posts_db.post" || db.GetServicePath() != "" {
		t.Errorf("database taint paths = (%q, %q)", db.GetDatabasePath(), db.GetServicePath())
	}
	if db.GetDatabaseCall() != dbCall || db.GetServiceCall() != nil {
		t.Errorf("database taint must only expose its database call")
	}
	if !db.IsReadValue() || db.IsReadKey() {
		t.Errorf("database taint read flags = (key=%v, value=%v), want (false, true)", db.IsReadKey(), db.IsReadValue())
	}

	if !svc.IsServiceTaint() || svc.IsDatabaseTaint() {
		t.Errorf("second taint must be a service taint")
	}
	if svc.GetServicePath() != "StorageService.ReadPost.t19" || svc.GetDatabasePath() != "" {
		t.Errorf("service taint paths = (%q, %q)", svc.GetDatabasePath(), svc.GetServicePath())
	}
	if svc.GetServiceCall() != svcCall || svc.GetDatabaseCall() != nil {
		t.Errorf("service taint must only expose its service call")
	}
}

func TestSSATaintTIsScopedByCaller(t *testing.T) {
	graph := newTestGraph()
	node := newValNode(graph, 1)
	callNode := newValNode(graph, 14)
	dbCall := newTestDatabaseCall(graph, callNode, common.OP_READ)

	node.AddDatabaseTaintIfNotExists("_obj", "orders_db.orders", dbCall, false, true, "")
	taint := node.GetTaintsForPath("_obj")[0]
	if got := taint.GetT(); got != callNode.GetName() {
		t.Errorf("GetT() without caller = %q, want %q", got, callNode.GetName())
	}

	// combiner.go sets the caller timestamp once the callee graph is inlined
	taint.SetCallerT("t4")
	if got, want := taint.GetT(), "t4."+callNode.GetName(); got != want {
		t.Errorf("GetT() with caller = %q, want %q", got, want)
	}
	if got := taint.GetCallerT(); got != "t4" {
		t.Errorf("GetCallerT() = %q, want t4", got)
	}
}

func TestCombineTaintsMergesPaths(t *testing.T) {
	graph := newTestGraph()
	a := newValNode(graph, 1)
	b := newValNode(graph, 2)
	write := newTestDatabaseCall(graph, newValNode(graph, 3), common.OP_WRITE)
	read := newTestDatabaseCall(graph, newValNode(graph, 4), common.OP_READ)

	a.AddDatabaseTaintIfNotExists("_obj", "posts_db.post", write, false, false, "")
	b.AddDatabaseTaintIfNotExists("_obj", "posts_db.post", read, false, true, "")
	b.AddDatabaseTaintIfNotExists("_obj.PostID", "posts_db.post.PostID", read, false, true, "")

	a.CombineTaints(b.GetTaints())
	if got := len(a.GetTaintsForPath("_obj")); got != 2 {
		t.Errorf("_obj has %d taints after combine, want 2", got)
	}
	if got := len(a.GetTaintsForPath("_obj.PostID")); got != 1 {
		t.Errorf("_obj.PostID has %d taints after combine, want 1", got)
	}
}

func TestTaintAndTraceString(t *testing.T) {
	graph := newTestGraph()
	node := newValNode(graph, 1)
	dbCall := newTestDatabaseCall(graph, newValNode(graph, 21), common.OP_WRITE)
	readCall := newTestDatabaseCall(graph, newValNode(graph, 14), common.OP_READ)

	node.AddDatabaseTaintIfNotExists("_obj.PostID", "posts_db.post.PostID", dbCall, false, false, "")
	node.AddDatabaseTaintIfNotExists("_obj", "posts_db.post", readCall, true, false, "")

	want := "_obj\n[read] [K] [14:int] @ posts_db.post\n" +
		"_obj.PostID\n[write] [21:int] @ posts_db.post.PostID\n"
	if got := node.TaintAndTraceString(); got != want {
		t.Errorf("TaintAndTraceString() =\n%s\nwant\n%s", got, want)
	}
}

func TestNodeUsedInBson(t *testing.T) {
	node := newValNode(newTestGraph(), 1)
	if node.IsUsedInBson() || node.LabelsString() != "[ssa: const]" {
		t.Errorf("new node: bson = %v, labels = %q", node.IsUsedInBson(), node.LabelsString())
	}
	node.EnableUsedInBson()
	if !node.IsUsedInBson() || node.LabelsString() != "[ssa: const] [bson]" {
		t.Errorf("bson node: bson = %v, labels = %q", node.IsUsedInBson(), node.LabelsString())
	}
}

func TestNodeSimpleCopyDropsTaints(t *testing.T) {
	graph := newTestGraph()
	node := newValNode(graph, 1)
	node.EnableUsedInBson()
	node.AddDatabaseTaintIfNotExists("_obj", "posts_db.post", newTestDatabaseCall(graph, newValNode(graph, 2), common.OP_WRITE), false, false, "")

	cp := node.SimpleCopy()
	if cp == node || cp.GetName() != node.GetName() || cp.GetID() != node.GetID() || cp.GetValue() != node.GetValue() {
		t.Errorf("copy must be a new node with the same name, id and value")
	}
	if cp.IsTainted() || !cp.IsUsedInBson() {
		t.Errorf("copy: tainted = %v, bson = %v, want (false, true)", cp.IsTainted(), cp.IsUsedInBson())
	}
	if !node.IsTainted() {
		t.Errorf("copying must not clear the taints of the original node")
	}
}

func TestNodeString(t *testing.T) {
	node := newValNode(newTestGraph(), 7)
	if got, want := node.String(), "7:int: 7:int"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
	if node.TaintAndTraceString() != "" {
		t.Errorf("untainted node must have no taint string")
	}
}
