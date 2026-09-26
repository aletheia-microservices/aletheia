package ssagraph_test

import (
	"strings"
	"testing"

	"analyzer/pkg/analysis/service-level/ssagraph"
)

func TestEdgeTypeString(t *testing.T) {
	for edgeType := ssagraph.EDGE_USAGE; edgeType <= ssagraph.EDGE_ITERATOR_OF; edgeType++ {
		edge := ssagraph.NewEdge(edgeType, nil, nil, 0, "")
		if s := edge.GetTypeString(); s == "" || strings.Contains(s, "UNKNOWN") {
			t.Errorf("edge type %d has no string representation", edgeType)
		}
	}
}

func TestEdgeAccessors(t *testing.T) {
	graph := newTestGraph()
	from, to := newValNode(graph, 1), newValNode(graph, 2)
	edge := ssagraph.NewEdge(ssagraph.EDGE_EXTRACT, from, to, 1, "")

	if edge.GetType() != ssagraph.EDGE_EXTRACT || !edge.IsType(ssagraph.EDGE_EXTRACT) || edge.IsType(ssagraph.EDGE_FIELD) {
		t.Errorf("type = %s", edge.GetTypeString())
	}
	if edge.GetFromNode() != from || edge.GetToNode() != to || !edge.HasFromNode(from) || edge.HasFromNode(to) {
		t.Errorf("edge must go from %s to %s", from.GetName(), to.GetName())
	}
	if edge.GetIndex() != 1 || edge.GetParam() != "" {
		t.Errorf("index = %d, param = %q", edge.GetIndex(), edge.GetParam())
	}

	edge.SetPath("_obj.ID")
	if edge.GetPath() != "_obj.ID" || !edge.HasPath("_obj.ID") || edge.HasPath("_obj") {
		t.Errorf("path = %q", edge.GetPath())
	}
}

func TestEdgeTypeStringValues(t *testing.T) {
	tests := []struct {
		edgeType ssagraph.EdgeType
		want     string
	}{
		{ssagraph.EDGE_FIELD, "FIELD"},
		{ssagraph.EDGE_STORE_ADDRESS, "STORE_ADDRESS"},
		{ssagraph.EDGE_PHI_ON, "PHI_ON"},
		// arguments and receivers are both edges to the call
		{ssagraph.EDGE_ARG_ON_CALL, "CALL_ON"},
		{ssagraph.EDGE_RECEIVER_ON_CALL, "CALL_ON"},
		{ssagraph.EdgeType(-1), "UNKNOWN"},
	}
	for _, tt := range tests {
		if got := ssagraph.NewEdge(tt.edgeType, nil, nil, 0, "").GetTypeString(); got != tt.want {
			t.Errorf("edge type %d = %q, want %q", tt.edgeType, got, tt.want)
		}
	}
}
