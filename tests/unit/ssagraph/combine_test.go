package ssagraph_test

import (
	"slices"
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/service-level/ssagraph"
	"github.com/aletheia-microservices/aletheia/internal/analysis/service-level/ssagraph/tainter"
)

// each Call* function passes its arguments to a helper with a specific code shape
const combineSrc = `package shop

type Creator struct { Username string }
type Post struct { ID string; Text string; Tags []string; Creator Creator }

func readText(p *Post) string { return p.Text }
func readUsername(p *Post) string { return p.Creator.Username }
func CallRead(p *Post) (string, string) { return readText(p), readUsername(p) }

func wrap(id string) *Post { return &Post{ID: id} }
func CallWrap(id string) *Post { return wrap(id) }

func first(tags []string) string { return tags[0] }
func CallFirst(tags []string) string { return first(tags) }

func put(m map[string]string, v string) { m["key"] = v }
func CallPut(m map[string]string, v string) { put(m, v) }

func pick(c bool, a string, b string) string {
	x := a
	if c {
		x = b
	}
	return x
}
func CallPick(c bool, a string, b string) string { return pick(c, a, b) }

func greet(a string) string { return "hi " + a }
func CallGreet(a string) string { return greet(a) }

func outer(p *Post) string { return inner(p) }
func inner(p *Post) string { return p.ID }
func CallOuter(p *Post) string { return outer(p) }
`

type combined struct {
	graphs map[string]*ssagraph.SSAGraph
	caller *ssagraph.SSAGraph
	call   *ssagraph.MethodCall
	callee *ssagraph.SSAGraph // copy of the callee graph inlined for the call
}

// combineWithSeed taints the argument argIdx that caller passes to callee as if it was written to
// dbpath, combines the caller graph, and returns the copy of the callee graph made for the call
func combineWithSeed(t *testing.T, caller string, callee string, argIdx int, dbpath string) combined {
	t.Helper()
	graphs := buildTaintedGraphs(t, combineSrc)
	c := combined{graphs: graphs, caller: getGraph(t, graphs, caller)}
	for _, call := range c.caller.GetMethodCalls() {
		if call.GetFuncShortPath() == callee {
			c.call = call
		}
	}
	if c.call == nil {
		t.Fatalf("%s does not call %s", caller, callee)
	}
	seedWriteTaint(c.call.GetArgumentAt(argIdx), dbpath)

	tainter.Combine(c.caller, graphs)

	c.callee = c.caller.GetCombinedGraphForMethodCallIfExists(c.call)
	if c.callee == nil {
		t.Fatalf("%s is not combined into %s", callee, caller)
	}
	return c
}

func assertDBTaints(t *testing.T, what string, node *ssagraph.SSANode, want ...string) {
	t.Helper()
	if got := dbTaints(node); !slices.Equal(got, want) {
		t.Errorf("%s taints = %v, want %v", what, got, want)
	}
}

func TestCombineInlinesCopyOfCallee(t *testing.T) {
	c := combineWithSeed(t, "shop.CallRead", "shop.readText", 0, "posts_db.post")

	original := getGraph(t, c.graphs, "shop.readText")
	if c.callee == original || c.callee.GetFunctionShortPath() != original.GetFunctionShortPath() {
		t.Errorf("combined graph must be a copy of %s", original.String())
	}
	if c.caller.GetMethodCallForCombinedGraph(c.callee) != c.call {
		t.Errorf("combined graph must be mapped back to its call")
	}
	// each call gets its own copy
	if got := len(c.caller.GetAllCombinedGraphs()); got != 2 {
		t.Errorf("CallRead has %d combined graphs, want 2 (readText and readUsername)", got)
	}
	for _, node := range original.GetNodes() {
		if node.IsTainted() {
			t.Errorf("original graph must not be tainted: %s", node.String())
		}
	}
}

func TestCombinePropagatesThroughFields(t *testing.T) {
	c := combineWithSeed(t, "shop.CallRead", "shop.readText", 0, "posts_db.post")

	p := c.callee.GetParamAt(0)
	assertDBTaints(t, "param p", p, "_obj @ posts_db.post", "_obj.Text @ posts_db.post.Text")
	assertDBTaints(t, "p.Text", fieldNode(t, c.callee, p, "Text"), "_obj @ posts_db.post.Text")
	assertDBTaints(t, "returned value", c.callee.GetReturnsLst()[0][0], "_obj @ posts_db.post.Text")

	// nested fields extend the database path at every level
	username := c.caller.GetCombinedGraphForMethodCallIfExists(c.caller.GetMethodCalls()[1])
	if username.GetFunctionShortPath() != "shop.readUsername" {
		t.Fatalf("second combined graph = %s, want shop.readUsername", username.String())
	}
	assertDBTaints(t, "returned username", username.GetReturnsLst()[0][0], "_obj @ posts_db.post.Creator.Username")
}

func TestCombinePropagatesIntoBuiltStruct(t *testing.T) {
	c := combineWithSeed(t, "shop.CallWrap", "shop.wrap", 0, "posts_db.post.ID")

	// the id stored in Post.ID marks the field of the new post
	assertDBTaints(t, "returned post", c.callee.GetReturnsLst()[0][0], "_obj.ID @ posts_db.post.ID")
}

func TestCombinePropagatesThroughSlices(t *testing.T) {
	c := combineWithSeed(t, "shop.CallFirst", "shop.first", 0, "posts_db.post.Tags")

	assertDBTaints(t, "param tags", c.callee.GetParamAt(0), "_obj @ posts_db.post.Tags", "_obj[*] @ posts_db.post.Tags[*]")
	assertDBTaints(t, "tags[0]", c.callee.GetReturnsLst()[0][0], "_obj @ posts_db.post.Tags[*]")
}

func TestCombinePropagatesThroughMaps(t *testing.T) {
	c := combineWithSeed(t, "shop.CallPut", "shop.put", 1, "posts_db.post.Text")

	// m["key"] = v marks the value stored at that key
	assertDBTaints(t, "param m", c.callee.GetParamAt(0), "_obj.key.Val @ posts_db.post.Text")
}

func TestCombinePropagatesThroughPhiAndBinOp(t *testing.T) {
	c := combineWithSeed(t, "shop.CallPick", "shop.pick", 1, "posts_db.post.Text")
	assertDBTaints(t, "phi", c.callee.GetReturnsLst()[0][0], "_obj @ posts_db.post.Text")

	c = combineWithSeed(t, "shop.CallGreet", "shop.greet", 0, "posts_db.post.Text")
	assertDBTaints(t, `"hi " + a`, c.callee.GetReturnsLst()[0][0], "_obj @ posts_db.post.Text")
}

func TestCombineNestedHelpers(t *testing.T) {
	c := combineWithSeed(t, "shop.CallOuter", "shop.outer", 0, "posts_db.post")

	// outer is combined into CallOuter, and inner into the copy of outer
	nested := c.callee.GetAllCombinedGraphs()
	if len(nested) != 1 || nested[0].GetFunctionShortPath() != "shop.inner" {
		t.Fatalf("combined graphs of outer = %v, want [shop.inner]", nested)
	}
	assertDBTaints(t, "inner returned value", nested[0].GetReturnsLst()[0][0], "_obj @ posts_db.post.ID")
}
