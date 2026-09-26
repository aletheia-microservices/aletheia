package integration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/internal/analysis/service-level/ssagraph"

	"github.com/aletheia-microservices/aletheia/tests/runner"
)

// ---------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------

func getSSAGraph(t *testing.T, a *runner.Analysis, fnShortPath string) *ssagraph.SSAGraph {
	t.Helper()
	graph, ok := a.FuncGraphs[fnShortPath]
	if !ok {
		t.Fatalf("ssa graph %s not found", fnShortPath)
	}
	return graph
}

func getDatabaseCall(t *testing.T, graph *ssagraph.SSAGraph, dbpath string, method string) *ssagraph.DatabaseCall {
	t.Helper()
	for _, call := range graph.GetDatabaseCalls() {
		if call.GetDatabasePath() == dbpath && call.GetMethod() == method {
			return call
		}
	}
	t.Fatalf("database call %s.%s not found in %s", dbpath, method, graph.String())
	return nil
}

func getServiceCall(t *testing.T, graph *ssagraph.SSAGraph, serviceWithMethod string) *ssagraph.ServiceCall {
	t.Helper()
	for _, call := range graph.GetServiceCalls() {
		if call.GetServiceWithMethod() == serviceWithMethod {
			return call
		}
	}
	t.Fatalf("service call %s not found in %s", serviceWithMethod, graph.String())
	return nil
}

func getMethodCall(t *testing.T, graph *ssagraph.SSAGraph, fnShortPath string) *ssagraph.MethodCall {
	t.Helper()
	for _, call := range graph.GetMethodCalls() {
		if call.GetFuncShortPath() == fnShortPath {
			return call
		}
	}
	t.Fatalf("method call %s not found in %s", fnShortPath, graph.String())
	return nil
}

func getParam(t *testing.T, graph *ssagraph.SSAGraph, name string) *ssagraph.SSANode {
	t.Helper()
	for _, param := range graph.GetParams() {
		if param.GetName() == name {
			return param
		}
	}
	t.Fatalf("parameter %s not found in %s", name, graph.String())
	return nil
}

// dbTaint describes an expected database taint at an object path of a node
type dbTaint struct {
	objpath string
	dbpath  string
	op      common.DatabaseOperationType
	readKey bool
	readVal bool
	t       string // optional: expected (caller scoped) timestamp
}

func (d dbTaint) String() string {
	return fmt.Sprintf("%s -> [%s] (key=%v, value=%v) @ %s", d.objpath, common.OperationTypeToString(d.op), d.readKey, d.readVal, d.dbpath)
}

func findDatabaseTaint(node *ssagraph.SSANode, want dbTaint) *ssagraph.SSATaint {
	for _, taint := range node.GetTaintsForPath(want.objpath) {
		if taint.IsDatabaseTaint() && taint.GetDatabasePath() == want.dbpath && taint.GetDatabaseCall().GetOpType() == want.op {
			return taint
		}
	}
	return nil
}

func assertDatabaseTaint(t *testing.T, node *ssagraph.SSANode, want dbTaint) {
	t.Helper()
	taint := findDatabaseTaint(node, want)
	if taint == nil {
		t.Errorf("node %s: missing taint %s\ngot:\n%s", node.GetName(), want, node.TaintAndTraceString())
		return
	}
	if taint.IsReadKey() != want.readKey || taint.IsReadValue() != want.readVal {
		t.Errorf("node %s: taint %s has (key=%v, value=%v)", node.GetName(), want, taint.IsReadKey(), taint.IsReadValue())
	}
	if want.t != "" && taint.GetT() != want.t {
		t.Errorf("node %s: taint %s has t=%s, want %s", node.GetName(), want, taint.GetT(), want.t)
	}
}

func assertNotTaintedBy(t *testing.T, node *ssagraph.SSANode, dbprefix string) {
	t.Helper()
	for objpath, taints := range node.GetTaints() {
		for _, taint := range taints {
			if taint.IsDatabaseTaint() && strings.HasPrefix(taint.GetDatabasePath(), dbprefix) {
				t.Errorf("node %s: unexpected taint %s @ %s", node.GetName(), objpath, taint.GetDatabasePath())
			}
		}
	}
}

// assertServiceTaint checks that the node holds a service taint (trace) created by the service call
func assertServiceTaint(t *testing.T, node *ssagraph.SSANode, objpath string, call *ssagraph.ServiceCall) {
	t.Helper()
	for _, taint := range node.GetTaintsForPath(objpath) {
		if taint.IsServiceTaint() && taint.GetServiceCall() == call {
			if !strings.HasPrefix(taint.GetServicePath(), call.String()+".") {
				t.Errorf("node %s: service path %s must start with %s", node.GetName(), taint.GetServicePath(), call.String())
			}
			return
		}
	}
	t.Errorf("node %s: missing service taint %s -> %s\ngot:\n%s", node.GetName(), objpath, call.String(), node.TaintAndTraceString())
}

// ---------------------------------------------------------------------
// intra-procedural SSA taint propagation (postnotification)
// ---------------------------------------------------------------------

// StorePost builds a Post struct and writes it with InsertOne: the object and all of its (nested)
// fields must be tainted with the corresponding database fields, and the taints must be propagated
// upwards to the function parameters that were assigned to those fields
func TestSSATaintWriteObjectFields(t *testing.T) {
	a := runner.Get(t, "postnotification")
	graph := getSSAGraph(t, a, "postnotification.StorageService.StorePost")

	if got := len(graph.GetDatabaseCalls()); got != 1 {
		t.Fatalf("StorePost has %d database calls, want 1", got)
	}
	call := getDatabaseCall(t, graph, "posts_db.post", "InsertOne")
	if call.GetOpType() != common.OP_WRITE {
		t.Fatalf("InsertOne op type = %s, want write", common.OperationTypeToString(call.GetOpType()))
	}
	if got := len(call.GetArguments()); got != 1 {
		t.Fatalf("InsertOne has %d tracked arguments, want 1 (the post)", got)
	}

	post := call.GetArguments()[0]
	for _, field := range []string{"", ".PostID", ".ReqID", ".Text", ".Timestamp", ".Mentions", ".Mentions[*]", ".Creator", ".Creator.Username"} {
		assertDatabaseTaint(t, post, dbTaint{objpath: "_obj" + field, dbpath: "posts_db.post" + field, op: common.OP_WRITE, t: call.GetT()})
	}

	// fields are propagated upwards to the values they were assigned from
	assertDatabaseTaint(t, getParam(t, graph, "text"), dbTaint{objpath: "_obj", dbpath: "posts_db.post.Text", op: common.OP_WRITE})
	assertDatabaseTaint(t, getParam(t, graph, "reqID"), dbTaint{objpath: "_obj", dbpath: "posts_db.post.ReqID", op: common.OP_WRITE})
	// ... but never to unrelated fields
	assertNotTaintedBy(t, getParam(t, graph, "text"), "posts_db.post.ReqID")
	assertNotTaintedBy(t, getParam(t, graph, "ctx"), "posts_db")
}

// ReadPost filters by PostID (read key) and decodes the result into a Post (read value)
func TestSSATaintReadKeyAndValue(t *testing.T) {
	a := runner.Get(t, "postnotification")
	graph := getSSAGraph(t, a, "postnotification.StorageService.ReadPost")

	call := getDatabaseCall(t, graph, "posts_db.post", "FindOne")
	if call.GetOpType() != common.OP_READ {
		t.Fatalf("FindOne op type = %s, want read", common.OperationTypeToString(call.GetOpType()))
	}

	postID := getParam(t, graph, "postID")
	assertDatabaseTaint(t, postID, dbTaint{objpath: "_obj", dbpath: "posts_db.post.PostID", op: common.OP_READ, readKey: true, t: call.GetT()})
	assertNotTaintedBy(t, getParam(t, graph, "reqID"), "posts_db")

	var foundValue bool
	for _, arg := range call.GetArguments() {
		if taint := findDatabaseTaint(arg, dbTaint{objpath: "_obj", dbpath: "posts_db.post", op: common.OP_READ}); taint != nil {
			foundValue = true
			if !taint.IsReadValue() || taint.IsReadKey() {
				t.Errorf("decoded post must be tainted as read value")
			}
		}
	}
	if !foundValue {
		t.Errorf("FindOne has no argument tainted as the decoded post")
	}
}

func TestSSATaintDeleteKey(t *testing.T) {
	a := runner.Get(t, "postnotification")
	graph := getSSAGraph(t, a, "postnotification.StorageService.DeletePost")

	call := getDatabaseCall(t, graph, "posts_db.post", "DeleteOne")
	if call.GetOpType() != common.OP_DELETE {
		t.Fatalf("DeleteOne op type = %s, want delete", common.OperationTypeToString(call.GetOpType()))
	}
	assertDatabaseTaint(t, getParam(t, graph, "postID"), dbTaint{objpath: "_obj", dbpath: "posts_db.post.PostID", op: common.OP_DELETE, readKey: true})
}

// NotifyService pops a message from the queue and forwards its fields to StorageService.ReadPost:
// the RPC arguments must carry both the queue taints and the service taint of the RPC
func TestSSATaintQueueReadFlowsIntoRPCArguments(t *testing.T) {
	a := runner.Get(t, "postnotification")
	graph := getSSAGraph(t, a, "postnotification.NotifyService.Run")

	pop := getDatabaseCall(t, graph, "notifications_queue.notification", "Pop")
	if pop.GetOpType() != common.OP_READ {
		t.Fatalf("Pop op type = %s, want read", common.OperationTypeToString(pop.GetOpType()))
	}
	msg := pop.GetArguments()[0]
	assertDatabaseTaint(t, msg, dbTaint{objpath: "_obj", dbpath: "notifications_queue.notification", op: common.OP_READ, readVal: true})
	assertDatabaseTaint(t, msg, dbTaint{objpath: "_obj.PostID", dbpath: "notifications_queue.notification.PostID", op: common.OP_READ, readVal: true})

	rpc := getServiceCall(t, graph, "StorageService.ReadPost")
	if got := len(rpc.GetArguments()); got != 2 {
		t.Fatalf("ReadPost rpc has %d arguments, want 2 (reqID, postID)", got)
	}
	reqID, postID := rpc.GetArguments()[0], rpc.GetArguments()[1]
	assertDatabaseTaint(t, reqID, dbTaint{objpath: "_obj", dbpath: "notifications_queue.notification.ReqID", op: common.OP_READ, readVal: true})
	assertDatabaseTaint(t, postID, dbTaint{objpath: "_obj", dbpath: "notifications_queue.notification.PostID", op: common.OP_READ, readVal: true})
	assertServiceTaint(t, postID, "_obj", rpc)
	// the rpc taint is also propagated back to the field of the queue message
	assertServiceTaint(t, msg, "_obj.PostID", rpc)
}

// UploadPost pushes a message built from the return of StorageService.StorePost:
// the message field must carry a service taint pointing to that RPC
func TestSSATaintRPCReturnFlowsIntoQueueWrite(t *testing.T) {
	a := runner.Get(t, "postnotification")
	graph := getSSAGraph(t, a, "postnotification.UploadService.UploadPost")

	rpc := getServiceCall(t, graph, "StorageService.StorePost")
	if got := len(rpc.GetReturns()); got != 2 {
		t.Fatalf("StorePost rpc has %d returns, want 2 (postID, err)", got)
	}
	push := getDatabaseCall(t, graph, "notifications_queue.notification", "Push")
	msg := push.GetArguments()[0]

	assertDatabaseTaint(t, msg, dbTaint{objpath: "_obj.PostID", dbpath: "notifications_queue.notification.PostID", op: common.OP_WRITE})
	assertServiceTaint(t, msg, "_obj.PostID", rpc)
	assertServiceTaint(t, msg, "_obj.ReqID", rpc)

	// reqID is both passed to StorePost and written to the queue
	reqID := rpc.GetArguments()[0]
	assertDatabaseTaint(t, reqID, dbTaint{objpath: "_obj", dbpath: "notifications_queue.notification.ReqID", op: common.OP_WRITE})
	assertServiceTaint(t, reqID, "_obj", rpc)
}

func TestSSACallsAreRegisteredInProgramOrder(t *testing.T) {
	a := runner.Get(t, "postnotification")
	graph := getSSAGraph(t, a, "postnotification.UploadService.UploadPost")

	var kinds []string
	for _, call := range graph.GetAllCalls() {
		switch c := call.(type) {
		case *ssagraph.ServiceCall:
			kinds = append(kinds, "rpc:"+c.GetServiceWithMethod())
		case *ssagraph.DatabaseCall:
			kinds = append(kinds, "db:"+c.GetDatabasePath()+"."+c.GetMethod())
		}
	}
	want := []string{"rpc:StorageService.StorePost", "db:notifications_queue.notification.Push"}
	if strings.Join(kinds, ",") != strings.Join(want, ",") {
		t.Errorf("calls = %v, want %v", kinds, want)
	}
}

// ---------------------------------------------------------------------
// inter-procedural SSA taint propagation through combined graphs (digota)
// ---------------------------------------------------------------------

// OrderService.Pay calls the internal helpers storageGetOne and storageUpdate, which access the
// database: the helper graphs must be copied and inlined into the caller, with taints scoped by
// the caller timestamp (<caller t>.<callee t>) and propagated back to the caller objects
func TestSSACombineInlinesHelperDatabaseCalls(t *testing.T) {
	a := runner.Get(t, "digota")
	caller := getSSAGraph(t, a, "digota.OrderService.Pay")

	getOne := getMethodCall(t, caller, "digota.OrderService.storageGetOne")
	callee := caller.GetCombinedGraphForMethodCallIfExists(getOne)
	if callee == nil {
		t.Fatalf("storageGetOne is not combined into Pay")
	}
	original := getSSAGraph(t, a, "digota.OrderService.storageGetOne")
	if callee == original {
		t.Fatalf("combined graph must be a copy of the original callee graph")
	}
	if caller.GetMethodCallForCombinedGraph(callee) != getOne {
		t.Errorf("combined graph must be mapped back to its method call")
	}

	// callee database call inside the combined copy is scoped by the caller timestamp
	find := getDatabaseCall(t, callee, "orders_db.orders", "FindOne")
	scopedT := getOne.GetT() + "." + find.GetT()
	var sawScoped bool
	for _, arg := range find.GetArguments() {
		for _, taints := range arg.GetTaints() {
			for _, taint := range taints {
				if taint.GetCallerT() == "" {
					continue // see TestSSACombineScopesAllCalleeTaints
				}
				if taint.GetCallerT() != getOne.GetT() {
					t.Errorf("combined taint %s has caller t=%s, want %s", taint.GetPath(), taint.GetCallerT(), getOne.GetT())
				}
				if taint.GetDatabaseCall() == find {
					sawScoped = true
					if taint.GetT() != scopedT {
						t.Errorf("combined taint %s has t=%s, want %s", taint.GetDatabasePath(), taint.GetT(), scopedT)
					}
				}
			}
		}
	}
	if !sawScoped {
		t.Fatalf("no scoped taints from FindOne found in the combined graph")
	}

	// the original graph is left untouched (no caller timestamp)
	for _, arg := range getDatabaseCall(t, original, "orders_db.orders", "FindOne").GetArguments() {
		for _, taints := range arg.GetTaints() {
			for _, taint := range taints {
				if taint.GetCallerT() != "" {
					t.Errorf("original graph taint %s has caller t %s", taint.GetDatabasePath(), taint.GetCallerT())
				}
			}
		}
	}

	// propagation: caller args <<< callee params
	order := getOne.GetArgumentAt(2)
	assertDatabaseTaint(t, order, dbTaint{objpath: "_obj", dbpath: "orders_db.orders", op: common.OP_READ, readVal: true, t: scopedT})
	// the read key/value flags of order.Id are not checked because they depend on the propagation order
	// (see TestSSAReadKeyFlagIsDeterministic)
	if findDatabaseTaint(order, dbTaint{objpath: "_obj.Id", dbpath: "orders_db.orders.Id", op: common.OP_READ}) == nil {
		t.Errorf("order must be tainted with orders_db.orders.Id read\ngot:\n%s", order.TaintAndTraceString())
	}

	// the caller parameter assigned to order.Id is tainted by the helper read
	id := getParam(t, caller, "id")
	if findDatabaseTaint(id, dbTaint{objpath: "_obj", dbpath: "orders_db.orders.Id", op: common.OP_READ}) == nil {
		t.Errorf("Pay(id) must be tainted by orders_db.orders.Id read in storageGetOne\ngot:\n%s", id.TaintAndTraceString())
	}

	// the second helper (storageUpdate) writes the same object later in the caller
	update := getMethodCall(t, caller, "digota.OrderService.storageUpdate")
	if caller.GetCombinedGraphForMethodCallIfExists(update) == nil {
		t.Fatalf("storageUpdate is not combined into Pay")
	}
	replace := getDatabaseCall(t, caller.GetCombinedGraphForMethodCallIfExists(update), "orders_db.orders", "ReplaceOne")
	assertDatabaseTaint(t, order, dbTaint{objpath: "_obj", dbpath: "orders_db.orders", op: common.OP_UPDATE, t: update.GetT() + "." + replace.GetT()})
}

func TestSSACombineScopesAllCalleeTaints(t *testing.T) {
	t.Skip("known bug: taints propagated back into a combined graph (caller args >>> callee params, with callerT=\"\") " +
		"lose the caller scope, so e.g. Pay/storageGetOne FindOne has both [t4.t14] and [t14] taints for the same call")

	a := runner.Get(t, "digota")
	caller := getSSAGraph(t, a, "digota.OrderService.Pay")
	getOne := getMethodCall(t, caller, "digota.OrderService.storageGetOne")
	find := getDatabaseCall(t, caller.GetCombinedGraphForMethodCallIfExists(getOne), "orders_db.orders", "FindOne")
	for _, arg := range find.GetArguments() {
		for objpath, taints := range arg.GetTaints() {
			for _, taint := range taints {
				if taint.GetDatabaseCall() == find && taint.GetT() != getOne.GetT()+"."+find.GetT() {
					t.Errorf("%s: taint %s has unscoped t=%s", objpath, taint.GetDatabasePath(), taint.GetT())
				}
			}
		}
	}
}

func TestSSAReadKeyFlagIsDeterministic(t *testing.T) {
	a := runner.Get(t, "digota")
	caller := getSSAGraph(t, a, "digota.OrderService.Pay")
	order := getMethodCall(t, caller, "digota.OrderService.storageGetOne").GetArgumentAt(2)
	// order.Id is used as the filter key of FindOne and is also part of the decoded order, so it
	// has both a read key and a read value taint, whatever the propagation order
	var key, val bool
	for _, taint := range order.GetTaintsForPath("_obj.Id") {
		if taint.IsDatabaseTaint() && taint.GetDatabasePath() == "orders_db.orders.Id" && taint.GetDatabaseCall().GetOpType() == common.OP_READ {
			key = key || taint.IsReadKey()
			val = val || taint.IsReadValue()
		}
	}
	if !key || !val {
		t.Errorf("order.Id must have both read key and read value taints (key=%v, value=%v)\ngot:\n%s", key, val, order.TaintAndTraceString())
	}
}

// combined graphs propagate service taints too: the charge amount passed to
// PaymentService.NewCharge comes from the order read by the helper
func TestSSACombinePropagatesServiceTaintsToHelperObjects(t *testing.T) {
	a := runner.Get(t, "digota")
	caller := getSSAGraph(t, a, "digota.OrderService.Pay")

	getOne := getMethodCall(t, caller, "digota.OrderService.storageGetOne")
	newCharge := getServiceCall(t, caller, "PaymentService.NewCharge")
	order := getOne.GetArgumentAt(2)

	assertServiceTaint(t, order, "_obj.Amount", newCharge)
	assertServiceTaint(t, order, "_obj.Currency", newCharge)
	assertDatabaseTaint(t, order, dbTaint{objpath: "_obj.Amount", dbpath: "orders_db.orders.Amount", op: common.OP_READ, readVal: true})
}
