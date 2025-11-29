package tainter

import (
	"slices"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/ssagraph"
)

func Combine(graph *ssagraph.SSAGraph, graphs map[string]*ssagraph.SSAGraph) {
	for _, methodCall := range graph.GetMethodCalls() {
		toGraph := graphs[methodCall.GetFuncShortPath()]
		if toGraph == nil {
			// it's from an external package
			logrus.Tracef("[SSA GRAPH] toGraph not found for methodCall (%s)\n", methodCall.GetFuncShortPath())
			continue
		}
		toGraph = toGraph.SimpleCopy()
		graph.AddCombinedGraph(toGraph, methodCall)
		RunTainter(toGraph)
		callerT := methodCall.GetT()

		logrus.Debugf("combining SSA graphs (caller=%s) (at=%s) (callee=%s)\n", graph.String(), methodCall.GetID(), toGraph.String())

		// propagation: caller args <<< callee params
		// TODO: upper/lower taints
		for i, callee_params := range toGraph.GetParams() {
			caller_arg := methodCall.GetArgumentAt(i)
			callee_taints := callee_params.GetTaints()
			propagateTaints(graph, caller_arg, callee_taints, callerT)
		}

		// propagation: caller rets <<< callee rets
		// TODO: upper/lower taints
		for _, callee_rets := range toGraph.GetReturnsLst() {
			for i, callee_ret := range callee_rets {
				caller_ret := methodCall.GetReturnAt(i)
				callee_taints := callee_ret.GetTaints()
				propagateTaints(graph, caller_ret, callee_taints, callerT)
			}
		}

		var callee_objs []*ssagraph.SSANode
		for _, call := range toGraph.GetServiceCalls() {
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
		for _, call := range toGraph.GetDatabaseCalls() {
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

	// propagation: caller args >>> callee params
	// TODO: upper/lower taints
	for _, toGraph := range graph.GetAllCombinedGraphs() {
		methodCall := graph.GetMethodCallForCombinedGraph(toGraph)
		for i, arg := range methodCall.GetArguments() {
			callee_param := toGraph.GetParamAt(i)
			caller_taints := arg.GetTaints()
			propagateTaints(toGraph, callee_param, caller_taints, "")
		}
	}

	// propagation: caller rets >>> callee rets
	// TODO: upper/lower taints
	for _, toGraph := range graph.GetAllCombinedGraphs() {
		methodCall := graph.GetMethodCallForCombinedGraph(toGraph)
		for i, ret := range methodCall.GetReturns() {
			for _, callee_rets := range toGraph.GetReturnsLst() {
				if i < len(callee_rets) { // sanity check
					callee_ret := callee_rets[i]
					caller_taints := ret.GetTaints()
					propagateTaints(toGraph, callee_ret, caller_taints, "")
				}
			}
		}
	}

	for _, toGraph := range graph.GetAllCombinedGraphs() {
		if graph.GetFunctionShortPath() == toGraph.GetFunctionShortPath() {
			// skip to avoid recursion
			continue
		}
		Combine(toGraph, graphs)
	}
}

func propagateTaints(graph *ssagraph.SSAGraph, to_obj *ssagraph.SSANode, from_taints map[string][]*ssagraph.SSATaint, callerT string) {
	for objpath, taintsLst := range from_taints {
		for _, taint := range taintsLst {
			visited := make(map[ssa.Value]bool)
			var taintInfo TaintInfo
			path, ok := strings.CutPrefix(objpath, "_obj")
			if !ok {
				logrus.Fatalf("objpath (%s) does not have '_obj' prefix", objpath)
			}
			if taint.IsDatabaseTaint() {
				taintInfo = NewTaintInfoDatabase(taint.GetDatabasePath(), path, nil, taint.GetDatabaseCall(), taint.IsReadKey(), taint.IsReadValue())
			} else if taint.IsServiceTaint() {
				taintInfo = NewTaintInfoService(taint.GetServicePath(), path, nil, taint.GetServiceCall())
			} else {
				logrus.Fatalf("unexpected type of taint: %s\n", taint.String())
			}
			if callerT != "" {
				taintInfo.callerT = callerT
			}
			seenTaint = make(map[TaintInfo]bool)
			propagateTaintNearby(graph, false, to_obj.GetValue(), taintInfo, visited, false)
		}
	}
}
