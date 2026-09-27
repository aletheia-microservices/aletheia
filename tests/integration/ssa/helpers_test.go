package ssa

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/internal/analysis/service-level/ssagraph"
	"github.com/aletheia-microservices/aletheia/tests/runner"
)

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
