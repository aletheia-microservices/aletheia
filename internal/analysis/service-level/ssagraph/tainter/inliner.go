package tainter

import (
	"slices"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ssa"

	"github.com/aletheia-microservices/aletheia/internal/analysis/service-level/ssagraph"
)

// InlineMethodGraphs inlines the graphs of the methods called by callerGraph. For each method call in callerGraph,
// it copies the graph of the called method (the callee) from graphsByFunc, runs the tainter on the copy, stores
// the copy as a inlined graph of callerGraph, and propagates taints between the call and the copy in both directions:
//   - caller <<< callee: taints are scoped by the call's t (t14 becomes t4.t14)
//   - caller >>> callee: taints keep their own scope
//
// e.g., Pay calls storageGetOne(ctx, order) at t4, and storageGetOne reads the order with FindOne at t14:
//   - caller <<< callee: the callee param's taint [t14] orders_db.orders is added to Pay's order as [t4.t14]
//   - caller >>> callee: Pay's order taints, e.g., [t4.t14] and storageUpdate's [t106.t15], are added to
//     the callee param of the inlined storageGetOne with those same scopes
//
// the copies then inline the graphs of their own helper methods, recursively
func InlineMethodGraphs(callerGraph *ssagraph.SSAGraph, graphsByFunc map[string]*ssagraph.SSAGraph) {
	for _, methodCall := range callerGraph.GetMethodCalls() {
		// 1. copy the callee graph and store it as a inlined graph of the caller
		// (calls to functions without a graph are skipped)
		calleeGraph := graphsByFunc[methodCall.GetFuncShortPath()]
		if calleeGraph == nil {
			continue
		}
		calleeGraph = calleeGraph.SimpleCopy()
		callerGraph.AddInlinedGraph(calleeGraph, methodCall)

		// 2. run the tainter on the copy, so it has the taints of its own calls
		var callerT string
		RunTainter(calleeGraph)
		callerT = methodCall.GetT()

		// 3. caller <<< callee: propagate the callee taints to the caller, scoped by callerT (e.g. t4)

		// propagation: caller binds <<< callee free vars (goroutines)
		if calleeGraph.IsGoRoutine() {
			for i, callee_freevar := range calleeGraph.GetFreeVars() {
				caller_var := methodCall.GetBindAt(i)
				callee_taints := callee_freevar.GetTaints()
				propagateTaints(callerGraph, caller_var, callee_taints, callerT)
			}
			// TODO rets
		}

		// propagation: caller args <<< callee params
		// TODO: upper/lower taints
		for i, callee_param := range calleeGraph.GetParams() {
			caller_arg := methodCall.GetArgumentAt(i)
			callee_taints := callee_param.GetTaints()
			propagateTaints(callerGraph, caller_arg, callee_taints, callerT)
		}

		// propagation: caller rets <<< callee rets
		// TODO: upper/lower taints
		for _, callee_rets := range calleeGraph.GetReturnsLst() {
			for i, callee_ret := range callee_rets {
				caller_ret := methodCall.TryGetReturnAt(i)
				if caller_ret != nil {
					callee_taints := callee_ret.GetTaints()
					propagateTaints(callerGraph, caller_ret, callee_taints, callerT)
				}
			}
		}

		// 4. scope the service and database calls inside the copy, and the taints on their arguments and
		// returns, by callerT (e.g. t14 becomes t4.t14), so they are ordered in the caller's timeline
		for _, call := range calleeGraph.GetServiceCalls() {
			call.SetCallerT(callerT)
		}
		for _, call := range calleeGraph.GetDatabaseCalls() {
			call.SetCallerT(callerT)
		}
		var callee_objs []*ssagraph.SSANode
		for _, call := range calleeGraph.GetServiceCalls() {
			for _, obj := range call.GetArguments() {
				if !slices.Contains(callee_objs, obj) {
					callee_objs = append(callee_objs, obj)
				}
			}
			for _, obj := range call.GetReturns() {
				if !slices.Contains(callee_objs, obj) {
					callee_objs = append(callee_objs, obj)
				}
			}
		}
		for _, call := range calleeGraph.GetDatabaseCalls() {
			for _, obj := range call.GetArguments() {
				if !slices.Contains(callee_objs, obj) {
					callee_objs = append(callee_objs, obj)
				}
			}
		}
		for _, obj := range callee_objs {
			for _, taintLst := range obj.GetTaints() {
				for _, taint := range taintLst {
					taint.SetCallerT(callerT)
				}
			}
		}
	}

	// 5. caller >>> callee: propagate the caller taints back to each copy, keeping their own scope
	// this runs after all method calls are inlined, so each copy also gets the taints that other callees
	// added to the caller's objects (e.g. storageUpdate's t106.t15 on the order passed to storageGetOne)

	// propagation: caller args >>> callee params
	// TODO: upper/lower taints
	for _, calleeGraph := range callerGraph.GetAllInlinedGraphs() {
		methodCall := callerGraph.GetMethodCallForInlinedGraph(calleeGraph)
		for i, arg := range methodCall.GetArguments() {
			callee_param := calleeGraph.GetParamAt(i)
			caller_taints := arg.GetTaints()
			propagateTaints(calleeGraph, callee_param, caller_taints, "")
		}

		// propagation: caller binds >>> callee free vars (goroutines)
		if calleeGraph.IsGoRoutine() {
			for i, callee_freevar := range calleeGraph.GetFreeVars() {
				caller_var := methodCall.GetBindAt(i)
				caller_taints := caller_var.GetTaints()
				propagateTaints(calleeGraph, callee_freevar, caller_taints, "")
			}
			// TODO rets
		}
	}

	// propagation: caller rets >>> callee rets
	// TODO: upper/lower taints
	for _, calleeGraph := range callerGraph.GetAllInlinedGraphs() {
		methodCall := callerGraph.GetMethodCallForInlinedGraph(calleeGraph)
		for i, ret := range methodCall.GetReturns() {
			for _, callee_rets := range calleeGraph.GetReturnsLst() {
				if i < len(callee_rets) { // sanity check
					callee_ret := callee_rets[i]
					caller_taints := ret.GetTaints()
					propagateTaints(calleeGraph, callee_ret, caller_taints, "")
				}
			}
		}
	}

	// 6. inline the method graphs of each copy, recursively
	for _, calleeGraph := range callerGraph.GetAllInlinedGraphs() {
		if callerGraph.GetFunctionShortPath() == calleeGraph.GetFunctionShortPath() {
			// skip to avoid recursion
			continue
		}
		InlineMethodGraphs(calleeGraph, graphsByFunc)
	}
}

// propagateTaints copies the taints of an object in one graph (fromTaints) to the matching object of
// another graph (toObj), and then spreads them from toObj to the rest of that graph, as the tainter
// does for taints created from a call
//
// the object path of each taint (e.g., _obj.Status) is kept, so the taint lands on the same field of toObj
//
// scopeT scopes the copied taints (t14 becomes t4.t14) when propagating from a callee to its caller;
// when empty (from a caller to a callee), each taint keeps its own caller scope (see InlineMethodGraphs)
func propagateTaints(graph *ssagraph.SSAGraph, toObj *ssagraph.SSANode, fromTaints map[string][]*ssagraph.SSATaint, scopeT string) {
	for objpath, taintsLst := range fromTaints {
		for _, taint := range taintsLst {
			visited := make(map[ssa.Value]bool)
			var taintInfo TaintInfo
			fieldpath, ok := strings.CutPrefix(objpath, "_obj")
			if !ok {
				logrus.Fatalf("objpath (%s) does not have '_obj' prefix", objpath)
			}
			if taint.IsDatabaseTaint() {
				taintInfo = NewTaintInfoDatabase(taint.GetDatabasePath(), fieldpath, nil, taint.GetDatabaseCall(), taint.IsReadKey(), taint.IsReadValue())
			} else if taint.IsServiceTaint() {
				taintInfo = NewTaintInfoService(taint.GetServicePath(), fieldpath, nil, taint.GetServiceCall())
			} else {
				logrus.Fatalf("unexpected type of taint: %s\n", taint.String())
			}
			// scopeT is the scope to apply, set by InlineMethodGraphs depending on the direction:
			//
			// - caller <<< callee (if): InlineMethodGraphs passes the method call's t (e.g. t4), so the taint is scoped by it,
			//   e.g., the callee's taint [t14] is added to the caller as [t4.t14]
			//
			// - caller >>> callee (else): InlineMethodGraphs passes "", so the taint keeps its own scope,
			//   e.g., the caller's taint [t4.t14] (added earlier through the if) is added back to the callee as [t4.t14],
			//   not [t14], and a taint from the caller's own call [t43] stays [t43]
			if scopeT != "" {
				taintInfo.callerT = scopeT
			} else {
				taintInfo.callerT = taint.GetCallerT()
			}
			seenTaint = make(map[TaintInfoData]bool)
			if fieldpath != "" {
				taintInfo = taintInfo.disableObjectRoot()
			}
			propagateTaintNearby(graph, false, toObj.GetValue(), taintInfo, visited, false)
			seenTaint = nil
		}
	}
}
