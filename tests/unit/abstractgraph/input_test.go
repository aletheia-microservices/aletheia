package abstractgraph_test

import (
	"analyzer/pkg/analysis/service-level/ssagraph"
	"go/constant"
	"go/types"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/analysis/common"
	"analyzer/pkg/analysis/system-level/abstractgraph"
	abstractgraphparser "analyzer/pkg/analysis/system-level/abstractgraph/parser"
	"analyzer/pkg/app"
	"analyzer/pkg/app/backends"
)

const inputYAML = `calls:
  - call_id: rpc
    call_ts: t1
    call_type: RPC
    arguments:
      - name: id
        taints:
          _obj:
            - taint_type: TAINT_DATABASE
              call_id: db
              caller_t: t0
              path: store.items.ID
              database_taint:
                read_key: true
                read_value: false
    service_call:
      service: Storage
      method: Fetch
      func_short_path: shop.Storage.Fetch
      returns:
        - name: result
          taints:
            _obj:
              - taint_type: TAINT_SERVICE
                call_id: rpc
                path: Storage.Fetch.t1
  - call_id: db
    call_ts: t2
    call_type: DB
    arguments:
      - name: key
        taints:
          _obj:
            - taint_type: TAINT_DATABASE
              call_id: db
              path: store.items.ID
              database_taint:
                read_key: true
                read_value: false
    database_call:
      operation_type: read_many
      database: store
      schema: items
      method: Find
`

func inputGraph() (*abstractgraph.AbstractCallGraph, *abstractgraph.AbstractNode) {
	a := app.NewApp("input")
	a.AddDatabase(backends.NewDatabase("store", "NoSQLDatabase"))
	graph := abstractgraph.NewAbstractCallGraph(a)
	source := abstractgraph.NewAbstractNode("Shop.List", abstractgraph.NODE_SERVICE, "Shop", "List", "", "")
	graph.AddNode(source.GetName(), source)
	return graph, source
}

func writeInput(t *testing.T, data string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "calls.yaml")
	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestParseInputFile(t *testing.T) {
	graph, source := inputGraph()
	if err := abstractgraphparser.ParseFile(graph, source, writeInput(t, inputYAML)); err != nil {
		t.Fatal(err)
	}
	edges := graph.GetEdges()
	if len(edges) != 2 || graph.GetRPCCount() != 1 || graph.GetDBAccessCount() != 1 {
		t.Fatalf("unexpected graph counts: %d edges", len(edges))
	}
	if edges[0].GetFromNode() != source || edges[0].GetToNode().GetName() != "Storage.Fetch" {
		t.Fatal("incorrect RPC endpoints")
	}
	taints := edges[0].GetArgumentAt(0).GetPrimaryTaints()["_obj"]
	if len(taints) != 1 || taints[0].GetT() != "t0.t2" || !taints[0].IsReadKey() || taints[0].IsReadValue() {
		t.Fatalf("incorrect database taints: %v", taints)
	}
	if !edges[0].GetReturns()[0].IsTraced() {
		t.Fatal("missing service return trace")
	}
	if edges[1].GetOpType() != common.OP_READ_MANY {
		t.Fatal("lost read_many operation")
	}
	schema := graph.GetApp().GetDatabaseByName("store").GetSchemaByNameIfExists("items")
	if schema == nil || !schema.HasField("store.items.ID") {
		t.Fatal("missing database schema or field")
	}
}

func TestInputRejectsInvalidYAML(t *testing.T) {
	for name, data := range map[string]string{
		"unknown field":     strings.Replace(inputYAML, "call_ts:", "timestamp:", 1),
		"unknown reference": strings.Replace(inputYAML, "call_id: db", "call_id: missing", 1),
		"duplicate ID":      strings.Replace(inputYAML, "  - call_id: db", "  - call_id: rpc", 1),
		"wrong payload":     strings.Replace(inputYAML, "call_type: RPC", "call_type: DB", 1),
		"unknown operation": strings.Replace(inputYAML, "read_many", "invalid", 1),
		"nil call":          "calls: [null]",
		"nil node":          "calls: [{call_id: rpc, call_type: RPC, service_call: {service: S, method: M}, arguments: [null]}]",
	} {
		t.Run(name, func(t *testing.T) {
			graph, source := inputGraph()
			if err := abstractgraphparser.ParseFile(graph, source, writeInput(t, data)); err == nil {
				t.Fatal("expected error")
			}
			if len(graph.GetNodes()) != 1 || len(graph.GetEdges()) != 0 {
				t.Fatal("invalid input mutated graph")
			}
		})
	}
}

func TestInputMissingDatabaseDoesNotMutateGraph(t *testing.T) {
	graph, source := inputGraph()
	data := strings.ReplaceAll(inputYAML, "store", "missing")
	if err := abstractgraphparser.ParseFile(graph, source, writeInput(t, data)); err == nil {
		t.Fatal("expected missing database error")
	}
	if len(graph.GetNodes()) != 1 || len(graph.GetEdges()) != 0 {
		t.Fatal("invalid input mutated graph")
	}
}

// Exercise the legacy SSA entry point without Blueprint's external IR setup.
func TestSSAParserPreservesDatabaseCallTaints(t *testing.T) {
	graph, _ := inputGraph()
	source := ssagraph.NewGraph(graph.GetApp(), "shop", "shop.Storage.Fetch", "Storage", "Fetch")
	value := ssa.NewConst(constant.MakeInt64(1), types.Typ[types.Int])
	arg := ssagraph.RegisterNewNodeVal(source, nil, value, "arg")
	call := ssagraph.NewDatabaseCall("lookup", arg, []*ssagraph.SSANode{arg}, "store", "items", "Find", common.OP_READ_MANY)
	arg.AddDatabaseTaintIfNotExists("_obj", "store.items.ID", call, true, false, "t3")
	source.AddCall(call)
	source.AddDatabaseCall(call)
	source.AddReturnsToLst(nil)
	abstractgraphparser.Parse(graph, "shop.Storage.Fetch", true, map[string]*ssagraph.SSAGraph{"shop.Storage.Fetch": source})
	edges := graph.GetEdges()
	if len(edges) != 2 || edges[0].GetEdgeType() != abstractgraph.EDGE_SERVICE_ENTRYPOINT {
		t.Fatal("missing entrypoint or database edge")
	}
	edge := edges[1]
	if edge.GetOpType() != common.OP_READ_MANY || edge.GetFromNode().GetName() != "Storage.Fetch" {
		t.Fatal("incorrect SSA edge")
	}
	taints := edge.GetArgumentAt(0).GetPrimaryTaints()["_obj"]
	if len(taints) != 1 || taints[0].GetT() != "t3."+call.GetT() || !taints[0].IsReadKey() {
		t.Fatalf("incorrect SSA taints: %v", taints)
	}
}
