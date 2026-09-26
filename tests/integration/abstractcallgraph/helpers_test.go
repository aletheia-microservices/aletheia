package abstractcallgraph

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph"
	"github.com/aletheia-microservices/aletheia/internal/utils"
)

func edgeKind(edge *abstractgraph.AbstractEdge) string {
	switch edge.GetEdgeType() {
	case abstractgraph.EDGE_SERVICE_ENTRYPOINT:
		return "entry"
	case abstractgraph.EDGE_SERVICE_RPC:
		return "rpc"
	case abstractgraph.EDGE_DATABASE_CALL:
		return common.OperationTypeToString(edge.GetOpType())
	}
	return "unknown"
}

// entrypoints returns the service methods called by the client, sorted
func entrypoints(graph *abstractgraph.AbstractCallGraph) []string {
	var lst []string
	for _, edge := range graph.GetEdges() {
		if edge.GetEdgeType() == abstractgraph.EDGE_SERVICE_ENTRYPOINT {
			lst = append(lst, edge.GetToNode().GetName())
		}
	}
	sort.Strings(lst)
	return lst
}

// rpcCalls returns every RPC as "<caller> -> <callee>", sorted
func rpcCalls(graph *abstractgraph.AbstractCallGraph) []string {
	var lst []string
	for _, edge := range graph.GetEdges() {
		if edge.GetEdgeType() == abstractgraph.EDGE_SERVICE_RPC {
			lst = append(lst, edge.GetFromNode().GetName()+" -> "+edge.GetToNode().GetName())
		}
	}
	sort.Strings(lst)
	return lst
}

// databaseCalls returns every database call as "<op> <caller> -> <database>.<schema>.<method>", sorted,
// optionally keeping only the databases for which keep returns true
func databaseCalls(graph *abstractgraph.AbstractCallGraph, keep func(database string) bool) []string {
	var lst []string
	for _, edge := range graph.GetEdges() {
		if edge.GetEdgeType() != abstractgraph.EDGE_DATABASE_CALL {
			continue
		}
		if keep != nil && !keep(edge.GetToNode().GetDatabaseName()) {
			continue
		}
		lst = append(lst, fmt.Sprintf("%s %s -> %s.%s", edgeKind(edge), edge.GetFromNode().GetName(), edge.GetToNode().GetName(), edge.GetMethod()))
	}
	sort.Strings(lst)
	return lst
}

// databaseOwners returns, for each database node, the sorted services that call it
func databaseOwners(graph *abstractgraph.AbstractCallGraph, keep func(database string) bool) map[string][]string {
	owners := make(map[string][]string)
	for _, edge := range graph.GetEdges() {
		if edge.GetEdgeType() != abstractgraph.EDGE_DATABASE_CALL {
			continue
		}
		if keep != nil && !keep(edge.GetToNode().GetDatabaseName()) {
			continue
		}
		db, service := edge.GetToNode().GetName(), edge.GetFromNode().GetServiceName()
		if !slices.Contains(owners[db], service) {
			owners[db] = append(owners[db], service)
		}
	}
	for db := range owners {
		sort.Strings(owners[db])
	}
	return owners
}

// counts returns the number of nodes of each type and of edges of each kind
func counts(graph *abstractgraph.AbstractCallGraph) map[string]int {
	c := make(map[string]int)
	for _, node := range graph.GetNodes() {
		switch node.GetNodeType() {
		case abstractgraph.NODE_SERVICE:
			c["service nodes"]++
		case abstractgraph.NODE_DATABASE:
			c["database nodes"]++
		}
	}
	for _, edge := range graph.GetEdges() {
		c[edgeKind(edge)]++
	}
	return c
}

// withPrefix keeps the elements of lst that start with one of the prefixes
func withPrefix(lst []string, prefixes ...string) []string {
	var filtered []string
	for _, elem := range lst {
		for _, prefix := range prefixes {
			if strings.HasPrefix(elem, prefix) {
				filtered = append(filtered, elem)
				break
			}
		}
	}
	return filtered
}

// isBefore reports whether timestamp t1 comes before t2 in the same request
func isBefore(t1 string, t2 string) bool {
	return utils.LessT(t1, t2)
}

func isNotCache(database string) bool {
	return !strings.HasSuffix(database, "_cache")
}

// assertList compares sorted lists and reports missing and unexpected elements
func assertList(t *testing.T, what string, got []string, want []string) {
	t.Helper()
	want = slices.Clone(want)
	sort.Strings(want)
	if slices.Equal(got, want) {
		return
	}
	var missing, unexpected []string
	remaining := slices.Clone(got)
	for _, w := range want {
		if i := slices.Index(remaining, w); i >= 0 {
			remaining = slices.Delete(remaining, i, i+1)
		} else {
			missing = append(missing, w)
		}
	}
	unexpected = remaining
	t.Errorf("%s differ\nmissing:\n\t%s\nunexpected:\n\t%s", what, strings.Join(missing, "\n\t"), strings.Join(unexpected, "\n\t"))
}

// assertContains checks that every element of want is in got
func assertContains(t *testing.T, what string, got []string, want ...string) {
	t.Helper()
	for _, w := range want {
		if !slices.Contains(got, w) {
			t.Errorf("%s: missing %q", what, w)
		}
	}
}

func assertCounts(t *testing.T, graph *abstractgraph.AbstractCallGraph, want map[string]int) {
	t.Helper()
	got := counts(graph)
	if !maps.Equal(got, want) {
		t.Errorf("counts = %v, want %v", got, want)
	}
}

func assertOwners(t *testing.T, got map[string][]string, want map[string][]string) {
	t.Helper()
	for db, services := range want {
		if !slices.Equal(got[db], services) {
			t.Errorf("%s is called by %v, want %v", db, got[db], services)
		}
	}
	for db, services := range got {
		if _, ok := want[db]; !ok {
			t.Errorf("unexpected database %s called by %v", db, services)
		}
	}
}

func getEdge(t *testing.T, graph *abstractgraph.AbstractCallGraph, from string, to string, method string) *abstractgraph.AbstractEdge {
	t.Helper()
	for _, edge := range graph.GetEdges() {
		if edge.GetFromNode().GetName() == from && edge.GetToNode().GetName() == to && edge.GetMethod() == method {
			return edge
		}
	}
	t.Fatalf("edge %s -> %s.%s not found", from, to, method)
	return nil
}

func hasEdge(graph *abstractgraph.AbstractCallGraph, from string, to string, method string) bool {
	for _, edge := range graph.GetEdges() {
		if edge.GetFromNode().GetName() == from && edge.GetToNode().GetName() == to && edge.GetMethod() == method {
			return true
		}
	}
	return false
}

func taintOpType(taint *abstractgraph.AbstractTaint) common.DatabaseOperationType {
	switch {
	case taint.IsRead():
		return common.OP_READ
	case taint.IsWrite():
		return common.OP_WRITE
	case taint.IsUpdate():
		return common.OP_UPDATE
	case taint.IsDelete():
		return common.OP_DELETE
	}
	return common.OP_UNDEFINED
}

// assertPrimaryTaint checks that obj has a primary taint for dbpath with the given operation at objpath
func assertPrimaryTaint(t *testing.T, obj *abstractgraph.AbstractObject, objpath string, dbpath string, op common.DatabaseOperationType) *abstractgraph.AbstractTaint {
	t.Helper()
	for _, taint := range obj.GetTaintsForObjectPath(objpath) {
		if taint.IsPrimary() && taint.GetDatabasePath() == dbpath && taintOpType(taint) == op {
			return taint
		}
	}
	t.Errorf("object %s: missing primary %s taint %s @ %s\ngot:\n%s", obj.GetName(), common.OperationTypeToString(op), objpath, dbpath, obj.Annotations())
	return nil
}

func assertNotTainted(t *testing.T, obj *abstractgraph.AbstractObject) {
	t.Helper()
	if obj.IsTainted() {
		t.Errorf("object %s must not be tainted, got:\n%s", obj.GetName(), obj.Annotations())
	}
}

// assertTrace checks that obj has, at objpath, a trace of the RPC call whose service path ends with suffix
// (e.g., the field of a value returned by that call)
func assertTrace(t *testing.T, obj *abstractgraph.AbstractObject, objpath string, call *abstractgraph.AbstractEdge, suffix string) {
	t.Helper()
	for _, trace := range obj.GetTracesForObjectPath(objpath) {
		if trace.GetServiceCallID() == call.GetID() && strings.HasSuffix(trace.GetServicePath(), suffix) {
			return
		}
	}
	t.Errorf("object %s: missing trace %s from %s (...%s)\ngot:\n%s", obj.GetName(), objpath, call.String(), suffix, obj.Annotations())
}
