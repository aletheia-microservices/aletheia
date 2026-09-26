package ssagraph_test

import (
	"go/ast"
	"go/constant"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"sort"
	"testing"

	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/internal/analysis/service-level/ssagraph"
	ssaparser "github.com/aletheia-microservices/aletheia/internal/analysis/service-level/ssagraph/parser"
	"github.com/aletheia-microservices/aletheia/internal/analysis/service-level/ssagraph/tainter"
	"github.com/aletheia-microservices/aletheia/internal/app"
)

// newTestGraph returns an empty graph for a service method (no app is needed for these tests)
func newTestGraph() *ssagraph.SSAGraph {
	return ssagraph.NewGraph(nil, "postnotification", "postnotification.StorageService.StorePost", "StorageService", "StorePost")
}

// newValNode registers a node backed by a constant whose ssa name is "<n>:int"
func newValNode(graph *ssagraph.SSAGraph, n int64) *ssagraph.SSANode {
	val := ssa.NewConst(constant.MakeInt64(n), types.Typ[types.Int])
	return ssagraph.RegisterNewNodeVal(graph, nil, val, "val_"+val.Name())
}

func newTestDatabaseCall(graph *ssagraph.SSAGraph, node *ssagraph.SSANode, opType common.DatabaseOperationType) *ssagraph.DatabaseCall {
	return ssagraph.NewDatabaseCall(ssagraph.ComputeCallID(graph, node), node, nil, "posts_db", "post", "InsertOne", opType)
}

// buildGraphs compiles src (a Go file of package shop, without imports) to SSA and parses it into SSA
// graphs, indexed by function path (e.g., "shop.Service.Upload" for (*ServiceImpl).Upload)
func buildGraphs(t *testing.T, src string) map[string]*ssagraph.SSAGraph {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "shop.go", src, 0)
	if err != nil {
		t.Fatalf("parsing source: %v", err)
	}
	// function paths are computed relative to "workflow/", as in Blueprint apps
	pkg := types.NewPackage("unittest/workflow/shop", "shop")
	ssaPkg, _, err := ssautil.BuildPackage(&types.Config{}, fset, pkg, []*ast.File{file}, ssa.SanityCheckFunctions)
	if err != nil {
		t.Fatalf("building ssa: %v", err)
	}

	// the parser dumps the ssa code to output/{app}/ssa, so run it in a temp dir
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)

	graphs := make(map[string]*ssagraph.SSAGraph)
	ssaparser.RunSSAAnalysis(app.NewApp("unittest"), ssaPkg.Prog, ssaPkg, graphs)
	for _, graph := range graphs {
		graph.Sort()
	}
	return graphs
}

// buildTaintedGraphs is the same as buildGraphs but also runs the SSA tainter on every graph,
// which registers the calls of each graph
func buildTaintedGraphs(t *testing.T, src string) map[string]*ssagraph.SSAGraph {
	t.Helper()
	graphs := buildGraphs(t, src)
	for _, graph := range graphs {
		tainter.RunTainter(graph)
	}
	return graphs
}

func getGraph(t *testing.T, graphs map[string]*ssagraph.SSAGraph, fnShortPath string) *ssagraph.SSAGraph {
	t.Helper()
	graph, ok := graphs[fnShortPath]
	if !ok {
		var names []string
		for name := range graphs {
			names = append(names, name)
		}
		sort.Strings(names)
		t.Fatalf("graph %s not found, got %v", fnShortPath, names)
	}
	return graph
}

// seedWriteTaint taints the node as if it was written to dbpath by a database call
func seedWriteTaint(node *ssagraph.SSANode, dbpath string) {
	scratch := ssagraph.NewGraph(nil, "", "", "", "")
	call := ssagraph.NewDatabaseCall("seed.99:int", newValNode(scratch, 99), nil, "posts_db", "post", "InsertOne", common.OP_WRITE)
	node.AddDatabaseTaintIfNotExists("_obj", dbpath, call, false, false, "")
}

// dbTaints returns the database taints of a node as "<object path> @ <db path>", sorted
func dbTaints(node *ssagraph.SSANode) []string {
	var taints []string
	for objpath, lst := range node.GetTaints() {
		for _, taint := range lst {
			if taint.IsDatabaseTaint() {
				taints = append(taints, objpath+" @ "+taint.GetDatabasePath())
			}
		}
	}
	sort.Strings(taints)
	return taints
}

// fieldNode returns the node for the field access "<from>.<field>"
func fieldNode(t *testing.T, graph *ssagraph.SSAGraph, from *ssagraph.SSANode, field string) *ssagraph.SSANode {
	t.Helper()
	for _, edge := range graph.GetEdgesTypedFrom(from, ssagraph.EDGE_FIELD) {
		if edge.GetParam() == field {
			return edge.GetToNode()
		}
	}
	t.Fatalf("field %s.%s not found in %s", from.GetName(), field, graph.String())
	return nil
}
