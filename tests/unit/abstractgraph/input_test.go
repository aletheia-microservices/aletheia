package abstractgraph_test

import (
	"fmt"
	"go/constant"
	"go/types"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"golang.org/x/tools/go/ssa"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/internal/analysis/service-level/ssagraph"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph"
	abstractgraphinput "github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph/input"
	abstractgraphparser "github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph/parser"
	"github.com/aletheia-microservices/aletheia/internal/app"
	"github.com/aletheia-microservices/aletheia/internal/app/backends"
)

// Shop.List calls Storage.Fetch and the store database, and Storage.Fetch calls the store database
// (see inputSSAGraphs for the equivalent SSA graphs)
const inputYAML = `entrypoints: [shop.Shop.List]
functions:
  - func_short_path: shop.Shop.List
    service: Shop
    method: List
    params:
      - name: req
        taints:
          _obj:
            - taint_type: TAINT_DATABASE
              call_id: db
              path: store.items.ID
              database_taint:
                read_key: true
                read_value: false
    returns:
      - - name: int
          taints:
            _obj:
              - taint_type: TAINT_SERVICE
                call_id: rpc
                path: Storage.Fetch.result
      - - name: int
          taints:
            _obj:
              - taint_type: TAINT_DATABASE
                call_id: db
                path: store.items
                database_taint:
                  read_key: false
                  read_value: true
    calls:
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
                    path: Storage.Fetch.result
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
  - func_short_path: shop.Storage.Fetch
    service: Storage
    method: Fetch
    params:
      - name: id
        taints:
          _obj:
            - taint_type: TAINT_DATABASE
              call_id: fetch
              path: store.items.ID
              database_taint:
                read_key: true
                read_value: false
    returns:
      - - name: int
    calls:
      - call_id: fetch
        call_ts: t3
        call_type: DB
        arguments:
          - name: id
            taints:
              _obj:
                - taint_type: TAINT_DATABASE
                  call_id: fetch
                  path: store.items.ID
                  database_taint:
                    read_key: true
                    read_value: false
        database_call:
          operation_type: read
          database: store
          schema: items
          method: FindOne
`

func inputGraph() *abstractgraph.AbstractCallGraph {
	a := app.NewApp("input")
	a.AddDatabase(backends.NewDatabase("store", "NoSQLDatabase"))
	return abstractgraph.NewAbstractCallGraph(a)
}

func writeInput(t *testing.T, data string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "input.yaml")
	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

// ssaValue is an SSA value with a given name, since constants are named after their value
type ssaValue struct {
	*ssa.Const
	name string
}

func (v ssaValue) Name() string {
	return v.name
}

func newSSANode(graph *ssagraph.SSAGraph, name string) *ssagraph.SSANode {
	return ssagraph.RegisterNewNodeVal(graph, nil, ssaValue{ssa.NewConst(constant.MakeInt64(0), types.Typ[types.Int]), name}, name)
}

// inputSSAGraphs builds the SSA graphs equivalent to inputYAML
func inputSSAGraphs(a *app.App) map[string]*ssagraph.SSAGraph {
	fetch := ssagraph.NewGraph(a, "shop", "shop.Storage.Fetch", "Storage", "Fetch")
	id := newSSANode(fetch, "id")
	fetch.AddParameter(newSSANode(fetch, "s"))
	fetch.AddParameter(newSSANode(fetch, "ctx"))
	fetch.AddParameter(id)
	fetchCall := ssagraph.NewDatabaseCall("fetch", newSSANode(fetch, "t3"), []*ssagraph.SSANode{id}, "store", "items", "FindOne", common.OP_READ)
	id.AddDatabaseTaintIfNotExists("_obj", "store.items.ID", fetchCall, true, false, "")
	fetch.AddCall(fetchCall)
	fetch.AddDatabaseCall(fetchCall)
	fetch.AddReturnsToLst([]*ssagraph.SSANode{newSSANode(fetch, "out")})

	list := ssagraph.NewGraph(a, "shop", "shop.Shop.List", "Shop", "List")
	req, arg, result, key, item := newSSANode(list, "req"), newSSANode(list, "id"), newSSANode(list, "result"), newSSANode(list, "key"), newSSANode(list, "item")
	list.AddParameter(newSSANode(list, "s"))
	list.AddParameter(newSSANode(list, "ctx"))
	list.AddParameter(req)
	rpc := ssagraph.NewServiceCall("rpc", newSSANode(list, "t1"), []*ssagraph.SSANode{arg}, []*ssagraph.SSANode{result}, "Storage", "Fetch", "shop.Storage.Fetch")
	db := ssagraph.NewDatabaseCall("db", newSSANode(list, "t2"), []*ssagraph.SSANode{key}, "store", "items", "Find", common.OP_READ_MANY)
	req.AddDatabaseTaintIfNotExists("_obj", "store.items.ID", db, true, false, "")
	arg.AddDatabaseTaintIfNotExists("_obj", "store.items.ID", db, true, false, "t0")
	result.AddServiceTaintIfNotExists("_obj", "Storage.Fetch.result", rpc, "")
	key.AddDatabaseTaintIfNotExists("_obj", "store.items.ID", db, true, false, "")
	item.AddDatabaseTaintIfNotExists("_obj", "store.items", db, false, true, "")
	list.AddCall(rpc)
	list.AddServiceCall(rpc)
	list.AddCall(db)
	list.AddDatabaseCall(db)
	list.AddReturnsToLst([]*ssagraph.SSANode{result})
	list.AddReturnsToLst([]*ssagraph.SSANode{item})

	return map[string]*ssagraph.SSAGraph{"shop.Shop.List": list, "shop.Storage.Fetch": fetch}
}

// describeGraph lists the nodes and edges of graph with their objects, to compare graphs built by different parsers
func describeGraph(graph *abstractgraph.AbstractCallGraph) []string {
	var lines []string
	for _, name := range slices.Sorted(maps.Keys(graph.GetNodes())) {
		node := graph.GetNodes()[name]
		lines = append(lines, fmt.Sprintf("node %s (parsed=%t)", name, node.IsParsed()))
		lines = append(lines, describeObjects("param", node.GetParams())...)
		lines = append(lines, describeObjects("ret", node.GetReturns())...)
	}
	for _, edge := range graph.GetEdges() {
		lines = append(lines, fmt.Sprintf("edge %d %s %s %s(): %s -> %s (%s)", edge.GetEdgeType(), edge.GetID(), edge.GetT(), edge.GetMethod(), edge.GetFromNode(), edge.GetToNode(), common.OperationTypeToString(edge.GetOpType())))
		lines = append(lines, describeObjects("arg", edge.GetArguments())...)
		lines = append(lines, describeObjects("ret", edge.GetReturns())...)
	}
	return append(lines, fmt.Sprintf("rpcs=%d db_accesses=%d", graph.GetRPCCount(), graph.GetDBAccessCount()))
}

func describeObjects(kind string, objs []*abstractgraph.AbstractObject) []string {
	var lines []string
	for _, obj := range objs {
		lines = append(lines, "\t"+kind+" "+obj.GetName())
		for _, path := range slices.Sorted(maps.Keys(obj.GetTaints())) {
			for _, taint := range obj.GetTaints()[path] {
				lines = append(lines, "\t\t"+path+" @ "+taint.LongLongString())
			}
		}
		for _, path := range slices.Sorted(maps.Keys(obj.GetTraces())) {
			for _, trace := range obj.GetTraces()[path] {
				lines = append(lines, "\t\t"+path+" @ "+trace.GetT()+" "+trace.LongString())
			}
		}
	}
	return lines
}

func TestParseInputFile(t *testing.T) {
	graph := inputGraph()
	if err := abstractgraphparser.ParseFile(graph, writeInput(t, inputYAML)); err != nil {
		t.Fatal(err)
	}
	edges := graph.GetEdges()
	if len(edges) != 4 || graph.GetRPCCount() != 2 || graph.GetDBAccessCount() != 2 {
		t.Fatalf("unexpected graph counts: %d edges", len(edges))
	}
	list, fetch := graph.GetNodeByName("Shop.List"), graph.GetNodeByName("Storage.Fetch")

	entry := edges[0]
	if entry.GetEdgeType() != abstractgraph.EDGE_SERVICE_ENTRYPOINT || entry.GetFromNode().GetName() != "client" || entry.GetToNode() != list {
		t.Fatal("incorrect entrypoint edge")
	}
	if entry.GetID() != "shop.Shop.List" || len(entry.GetArguments()) != 1 || entry.GetArgumentAt(0).GetName() != "req" || entry.GetArgumentAt(0).IsTainted() {
		t.Fatal("incorrect entrypoint arguments")
	}
	if !list.IsParsed() || len(list.GetParams()) != 1 || !list.GetParamAt(0).IsTainted() {
		t.Fatal("incorrect entrypoint params")
	}
	// the returns of the second return statement are merged into the first
	if len(list.GetReturns()) != 1 || !list.GetReturnAt(0).IsTraced() {
		t.Fatal("missing entrypoint return trace")
	}
	if taints := list.GetReturnAt(0).GetTaints()["_obj"]; len(taints) != 1 || taints[0].GetDatabasePath() != "store.items" || !taints[0].IsReadValue() {
		t.Fatalf("returns not merged: %v", taints)
	}

	if edges[1].GetFromNode() != list || edges[1].GetToNode() != fetch {
		t.Fatal("incorrect RPC endpoints")
	}
	taints := edges[1].GetArgumentAt(0).GetPrimaryTaints()["_obj"]
	if len(taints) != 1 || taints[0].GetT() != "t0.t2" || !taints[0].IsReadKey() || taints[0].IsReadValue() {
		t.Fatalf("incorrect database taints: %v", taints)
	}
	if !edges[1].GetReturns()[0].IsTraced() {
		t.Fatal("missing service return trace")
	}
	if edges[2].GetOpType() != common.OP_READ_MANY {
		t.Fatal("lost read_many operation")
	}

	// the callee is parsed from its own function
	if !fetch.IsParsed() || len(fetch.GetParams()) != 1 || fetch.GetParamAt(0).GetName() != "id" || len(fetch.GetReturns()) != 1 {
		t.Fatal("incorrect callee node")
	}
	if edges[3].GetFromNode() != fetch || edges[3].GetID() != "fetch" || edges[3].GetOpType() != common.OP_READ {
		t.Fatal("incorrect callee database edge")
	}

	schema := graph.GetApp().GetDatabaseByName("store").GetSchemaByNameIfExists("items")
	if schema == nil || !schema.HasField("store.items.ID") {
		t.Fatal("missing database schema or field")
	}
}

// splitInput splits inputYAML into one model with Shop.List (which calls Storage.Fetch) and another with Storage.Fetch
func splitInput(t *testing.T) (string, string) {
	t.Helper()
	i := strings.Index(inputYAML, "  - func_short_path: shop.Storage.Fetch")
	if i < 0 {
		t.Fatal("Storage.Fetch not found in input")
	}
	return inputYAML[:i], "functions:\n" + inputYAML[i:]
}

func TestInputCombinesModelsInFolder(t *testing.T) {
	frontend, backend := splitInput(t)
	dir := t.TempDir()
	for name, data := range map[string]string{"frontend.yaml": frontend, "backend.yml": backend, "README.md": "not a model"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
	}
	combinedGraph := inputGraph()
	if err := abstractgraphparser.ParseFile(combinedGraph, dir); err != nil {
		t.Fatal(err)
	}
	singleGraph := inputGraph()
	if err := abstractgraphparser.ParseFile(singleGraph, writeInput(t, inputYAML)); err != nil {
		t.Fatal(err)
	}
	if got, want := describeGraph(combinedGraph), describeGraph(singleGraph); !slices.Equal(got, want) {
		t.Fatalf("combined graph:\n%s\n\nsingle graph:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}

	// each model is only valid when combined with the other one
	for name, data := range map[string]string{"frontend": frontend, "backend": backend} {
		if _, err := abstractgraphinput.LoadInputModel(writeInput(t, data)); err == nil {
			t.Fatalf("expected error for %s alone", name)
		}
	}
	if _, err := abstractgraphinput.LoadInputModel(t.TempDir()); err == nil {
		t.Fatal("expected error for folder without input models")
	}
}

func TestMergeInputModels(t *testing.T) {
	frontend := &abstractgraphinput.InputModel{
		App:         "shop",
		Databases:   []*abstractgraphinput.Database{{Name: "events", Type: "Queue", Schemas: []*abstractgraphinput.Schema{{Name: "order"}}}},
		Entrypoints: []string{"f"},
	}
	backend := &abstractgraphinput.InputModel{
		Databases: []*abstractgraphinput.Database{
			{Name: "events", Type: "Queue", Schemas: []*abstractgraphinput.Schema{{Name: "payment"}}},
			{Name: "orders", Type: "NoSQLDatabase"},
		},
		Functions: []*abstractgraphinput.Function{{FuncShortPath: "f", Service: "S", Method: "M"}},
	}
	merged, err := abstractgraphinput.MergeInputModels(frontend, backend)
	if err != nil {
		t.Fatal(err)
	}
	if merged.App != "shop" || len(merged.Entrypoints) != 1 || len(merged.Functions) != 1 || len(merged.Databases) != 2 {
		t.Fatalf("unexpected merged model: %+v", merged)
	}
	if events := merged.Databases[0]; events.Name != "events" || len(events.Schemas) != 2 || len(frontend.Databases[0].Schemas) != 1 {
		t.Fatal("database schemas not combined (or models were modified)")
	}
	if _, err := merged.Index(); err != nil {
		t.Fatal(err)
	}

	for name, other := range map[string]*abstractgraphinput.InputModel{
		"different apps":           {App: "other"},
		"different database types": {Databases: []*abstractgraphinput.Database{{Name: "events", Type: "Cache"}}},
	} {
		if _, err := abstractgraphinput.MergeInputModels(frontend, other); err == nil {
			t.Fatalf("%s: expected error", name)
		}
	}
}

func TestInputParsesFunctionsOnce(t *testing.T) {
	graph := inputGraph()
	data := strings.Replace(inputYAML, "entrypoints: [shop.Shop.List]", "entrypoints: [shop.Shop.List, shop.Storage.Fetch]", 1)
	if err := abstractgraphparser.ParseFile(graph, writeInput(t, data)); err != nil {
		t.Fatal(err)
	}
	edges := graph.GetEdges()
	if len(edges) != 5 || graph.GetRPCCount() != 3 || graph.GetDBAccessCount() != 2 {
		t.Fatalf("unexpected graph counts: %d edges", len(edges))
	}
	if edges[4].GetEdgeType() != abstractgraph.EDGE_SERVICE_ENTRYPOINT || edges[4].GetToNode().GetName() != "Storage.Fetch" {
		t.Fatal("missing entrypoint edge for already parsed function")
	}
}

func TestInputParserMatchesSSAParser(t *testing.T) {
	inputAbsGraph := inputGraph()
	if err := abstractgraphparser.ParseFile(inputAbsGraph, writeInput(t, inputYAML)); err != nil {
		t.Fatal(err)
	}
	ssaAbsGraph := inputGraph()
	abstractgraphparser.Parse(ssaAbsGraph, "shop.Shop.List", true, inputSSAGraphs(ssaAbsGraph.GetApp()))

	got, want := describeGraph(inputAbsGraph), describeGraph(ssaAbsGraph)
	if !slices.Equal(got, want) {
		t.Fatalf("input graph:\n%s\n\nSSA graph:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestInputRejectsInvalidYAML(t *testing.T) {
	for name, data := range map[string]string{
		"unknown field":        strings.Replace(inputYAML, "call_ts:", "timestamp:", 1),
		"unknown reference":    strings.Replace(inputYAML, "call_id: db", "call_id: missing", 1),
		"foreign reference":    strings.Replace(inputYAML, "call_id: fetch", "call_id: db", 1),
		"duplicate ID":         strings.Replace(inputYAML, "  - call_id: db", "  - call_id: rpc", 1),
		"wrong payload":        strings.Replace(inputYAML, "call_type: RPC", "call_type: DB", 1),
		"unknown operation":    strings.Replace(inputYAML, "read_many", "invalid", 1),
		"no entrypoints":       strings.Replace(inputYAML, "entrypoints: [shop.Shop.List]\n", "", 1),
		"unknown entrypoint":   strings.Replace(inputYAML, "shop.Shop.List]", "shop.Shop.Missing]", 1),
		"duplicate function":   strings.Replace(inputYAML, "- func_short_path: shop.Storage.Fetch", "- func_short_path: shop.Shop.List", 1),
		"unknown callee":       strings.Replace(inputYAML, "func_short_path: shop.Storage.Fetch", "func_short_path: shop.Storage.Missing", 1),
		"mismatched callee":    strings.Replace(inputYAML, "    service: Storage\n    method: Fetch\n", "    service: Storage\n    method: Get\n", 1),
		"uneven returns":       strings.Replace(inputYAML, "      - - name: int\n    calls:", "      - - name: int\n      - - name: int\n        - name: int\n    calls:", 1),
		"duplicate ID across":  strings.ReplaceAll(inputYAML, "call_id: fetch", "call_id: db"),
		"invalid call_ts":      strings.Replace(inputYAML, "call_ts: t1", "call_ts: 1", 1),
		"invalid caller_t":     strings.Replace(inputYAML, "caller_t: t0", "caller_t: x0", 1),
		"foreign trace":        strings.Replace(inputYAML, "path: Storage.Fetch.result", "path: Other.Fetch.result", 1),
		"unknown trace name":   strings.Replace(inputYAML, "path: Storage.Fetch.result", "path: Storage.Fetch.missing", 1),
		"unknown db type":      "databases: [{name: store, type: SQL}]\n" + inputYAML,
		"conflicting database": "databases: [{name: store, type: Queue}, {name: store, type: Cache}]\n" + inputYAML,
		"duplicate entrypoint": strings.Replace(inputYAML, "entrypoints: [shop.Shop.List]", "entrypoints: [shop.Shop.List, shop.Shop.List]", 1),
		"duplicate schema":     "databases: [{name: store, type: Queue, schemas: [{name: items}, {name: items}]}]\n" + inputYAML,
		"nil function":         "entrypoints: [f]\nfunctions: [null]",
		"nil call":             "entrypoints: [f]\nfunctions: [{func_short_path: f, service: S, method: M, calls: [null]}]",
		"nil node":             "entrypoints: [f]\nfunctions: [{func_short_path: f, service: S, method: M, params: [null]}]",
	} {
		t.Run(name, func(t *testing.T) {
			if data == inputYAML {
				t.Fatal("test input was not modified")
			}
			graph := inputGraph()
			err := abstractgraphparser.ParseFile(graph, writeInput(t, data))
			if err == nil {
				t.Fatal("expected error")
			}
			t.Log(err)
			if len(graph.GetNodes()) != 0 || len(graph.GetEdges()) != 0 {
				t.Fatal("invalid input mutated graph")
			}
		})
	}
}

func TestInputMissingDatabaseDoesNotMutateGraph(t *testing.T) {
	graph := inputGraph()
	data := strings.ReplaceAll(inputYAML, "store", "missing")
	if err := abstractgraphparser.ParseFile(graph, writeInput(t, data)); err == nil {
		t.Fatal("expected missing database error")
	}
	if len(graph.GetNodes()) != 0 || len(graph.GetEdges()) != 0 {
		t.Fatal("invalid input mutated graph")
	}
}

// Exercise the legacy SSA entry point without Blueprint's external IR setup
func TestSSAParserPreservesDatabaseCallTaints(t *testing.T) {
	graph := inputGraph()
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
