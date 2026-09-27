package ssa

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/tests/runner"
)

// OrderService.Pay calls the internal helpers storageGetOne and storageUpdate, which access the
// database: the helper graphs must be copied and inlined into the caller, with taints scoped by
// the caller timestamp (<caller t>.<callee t>) and propagated back to the caller objects
func TestDigotaInlinesHelperDatabaseCalls(t *testing.T) {
	a := runner.Get(t, "digota")
	caller := getSSAGraph(t, a, "digota.OrderService.Pay")

	getOne := getMethodCall(t, caller, "digota.OrderService.storageGetOne")
	callee := caller.GetInlinedGraphForMethodCallIfExists(getOne)
	if callee == nil {
		t.Fatalf("storageGetOne is not inlined into Pay")
	}
	original := getSSAGraph(t, a, "digota.OrderService.storageGetOne")
	if callee == original {
		t.Fatalf("inlined graph must be a copy of the original callee graph")
	}
	if caller.GetMethodCallForInlinedGraph(callee) != getOne {
		t.Errorf("inlined graph must be mapped back to its method call")
	}

	// callee database call inside the inlined copy is scoped by the caller timestamp,
	// e.g. FindOne at t14 in storageGetOne, called at t4 in Pay, is at t4.t14 in Pay
	find := getDatabaseCall(t, callee, "orders_db.orders", "FindOne")
	scopedT := getOne.GetT() + "." + find.GetT()
	if find.GetScopedT() != scopedT {
		t.Errorf("inlined call FindOne has scoped t=%s, want %s", find.GetScopedT(), scopedT)
	}
	var sawScoped bool
	for _, arg := range find.GetArguments() {
		for _, taints := range arg.GetTaints() {
			for _, taint := range taints {
				// the arguments also carry taints from other calls on the same order, which keep their own
				// scope, e.g. storageUpdate's ReplaceOne (t106.t15) or Pay's NewCharge (t43)
				if taint.IsDatabaseTaint() && taint.GetDatabaseCall() == find {
					sawScoped = true
					if taint.GetT() != scopedT {
						t.Errorf("inlined taint %s has t=%s, want %s", taint.GetDatabasePath(), taint.GetT(), scopedT)
					}
				}
			}
		}
	}
	if !sawScoped {
		t.Fatalf("no scoped taints from FindOne found in the inlined graph")
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
	// (see TestDigotaReadKeyFlagIsDeterministic)
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
	if caller.GetInlinedGraphForMethodCallIfExists(update) == nil {
		t.Fatalf("storageUpdate is not inlined into Pay")
	}
	replace := getDatabaseCall(t, caller.GetInlinedGraphForMethodCallIfExists(update), "orders_db.orders", "ReplaceOne")
	assertDatabaseTaint(t, order, dbTaint{objpath: "_obj", dbpath: "orders_db.orders", op: common.OP_UPDATE, t: update.GetT() + "." + replace.GetT()})
}

func TestDigotaInlineScopesAllCalleeTaints(t *testing.T) {
	a := runner.Get(t, "digota")
	caller := getSSAGraph(t, a, "digota.OrderService.Pay")
	getOne := getMethodCall(t, caller, "digota.OrderService.storageGetOne")
	find := getDatabaseCall(t, caller.GetInlinedGraphForMethodCallIfExists(getOne), "orders_db.orders", "FindOne")
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

func TestDigotaReadKeyFlagIsDeterministic(t *testing.T) {
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

// inlined graphs propagate service taints too: the charge amount passed to
// PaymentService.NewCharge comes from the order read by the helper
func TestDigotaInlinePropagatesServiceTaintsToHelperObjects(t *testing.T) {
	a := runner.Get(t, "digota")
	caller := getSSAGraph(t, a, "digota.OrderService.Pay")

	getOne := getMethodCall(t, caller, "digota.OrderService.storageGetOne")
	newCharge := getServiceCall(t, caller, "PaymentService.NewCharge")
	order := getOne.GetArgumentAt(2)

	assertServiceTaint(t, order, "_obj.Amount", newCharge)
	assertServiceTaint(t, order, "_obj.Currency", newCharge)
	assertDatabaseTaint(t, order, dbTaint{objpath: "_obj.Amount", dbpath: "orders_db.orders.Amount", op: common.OP_READ, readVal: true})
}
