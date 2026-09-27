package ssagraph_test

import (
	"slices"
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/service-level/ssagraph"
	"github.com/aletheia-microservices/aletheia/internal/analysis/service-level/ssagraph/tainter"
)

// each Call* function passes its arguments to a helper with a specific code shape
const inlineSrc = `package shop

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

type inlined struct {
	graphs map[string]*ssagraph.SSAGraph
	caller *ssagraph.SSAGraph
	call   *ssagraph.MethodCall
	callee *ssagraph.SSAGraph // copy of the callee graph inlined for the call
}

// inlineWithSeed taints the argument argIdx that caller passes to callee as if it was written to
// dbpath, inlines the method graphs of the caller, and returns the copy of the callee graph made for the call
func inlineWithSeed(t *testing.T, caller string, callee string, argIdx int, dbpath string) inlined {
	t.Helper()
	graphs := buildTaintedGraphs(t, inlineSrc)
	c := inlined{graphs: graphs, caller: getGraph(t, graphs, caller)}
	for _, call := range c.caller.GetMethodCalls() {
		if call.GetFuncShortPath() == callee {
			c.call = call
		}
	}
	if c.call == nil {
		t.Fatalf("%s does not call %s", caller, callee)
	}
	seedWriteTaint(c.call.GetArgumentAt(argIdx), dbpath)

	tainter.InlineMethodGraphs(c.caller, graphs)

	c.callee = c.caller.GetInlinedGraphForMethodCallIfExists(c.call)
	if c.callee == nil {
		t.Fatalf("%s is not inlined into %s", callee, caller)
	}
	return c
}

func assertDBTaints(t *testing.T, what string, node *ssagraph.SSANode, want ...string) {
	t.Helper()
	if got := dbTaints(node); !slices.Equal(got, want) {
		t.Errorf("%s taints = %v, want %v", what, got, want)
	}
}

func TestInlineCopiesCallee(t *testing.T) {
	c := inlineWithSeed(t, "shop.CallRead", "shop.readText", 0, "posts_db.post")

	original := getGraph(t, c.graphs, "shop.readText")
	if c.callee == original || c.callee.GetFunctionShortPath() != original.GetFunctionShortPath() {
		t.Errorf("inlined graph must be a copy of %s", original.String())
	}
	if c.caller.GetMethodCallForInlinedGraph(c.callee) != c.call {
		t.Errorf("inlined graph must be mapped back to its call")
	}
	// each call gets its own copy
	if got := len(c.caller.GetAllInlinedGraphs()); got != 2 {
		t.Errorf("CallRead has %d inlined graphs, want 2 (readText and readUsername)", got)
	}
	for _, node := range original.GetNodes() {
		if node.IsTainted() {
			t.Errorf("original graph must not be tainted: %s", node.String())
		}
	}
}

func TestInlinePropagatesThroughFields(t *testing.T) {
	c := inlineWithSeed(t, "shop.CallRead", "shop.readText", 0, "posts_db.post")

	p := c.callee.GetParamAt(0)
	assertDBTaints(t, "param p", p, "_obj @ posts_db.post", "_obj.Text @ posts_db.post.Text")
	assertDBTaints(t, "p.Text", fieldNode(t, c.callee, p, "Text"), "_obj @ posts_db.post.Text")
	assertDBTaints(t, "returned value", c.callee.GetReturnsLst()[0][0], "_obj @ posts_db.post.Text")

	// nested fields extend the database path at every level
	username := c.caller.GetInlinedGraphForMethodCallIfExists(c.caller.GetMethodCalls()[1])
	if username.GetFunctionShortPath() != "shop.readUsername" {
		t.Fatalf("second inlined graph = %s, want shop.readUsername", username.String())
	}
	assertDBTaints(t, "returned username", username.GetReturnsLst()[0][0], "_obj @ posts_db.post.Creator.Username")
}

func TestInlinePropagatesIntoBuiltStruct(t *testing.T) {
	c := inlineWithSeed(t, "shop.CallWrap", "shop.wrap", 0, "posts_db.post.ID")

	// the id stored in Post.ID marks the field of the new post
	assertDBTaints(t, "returned post", c.callee.GetReturnsLst()[0][0], "_obj.ID @ posts_db.post.ID")
}

func TestInlinePropagatesThroughSlices(t *testing.T) {
	c := inlineWithSeed(t, "shop.CallFirst", "shop.first", 0, "posts_db.post.Tags")

	assertDBTaints(t, "param tags", c.callee.GetParamAt(0), "_obj @ posts_db.post.Tags", "_obj[*] @ posts_db.post.Tags[*]")
	assertDBTaints(t, "tags[0]", c.callee.GetReturnsLst()[0][0], "_obj @ posts_db.post.Tags[*]")
}

func TestInlinePropagatesThroughMaps(t *testing.T) {
	c := inlineWithSeed(t, "shop.CallPut", "shop.put", 1, "posts_db.post.Text")

	// m["key"] = v marks the value stored at that key
	assertDBTaints(t, "param m", c.callee.GetParamAt(0), "_obj.key.Val @ posts_db.post.Text")
}

func TestInlinePropagatesThroughPhiAndBinOp(t *testing.T) {
	c := inlineWithSeed(t, "shop.CallPick", "shop.pick", 1, "posts_db.post.Text")
	assertDBTaints(t, "phi", c.callee.GetReturnsLst()[0][0], "_obj @ posts_db.post.Text")

	c = inlineWithSeed(t, "shop.CallGreet", "shop.greet", 0, "posts_db.post.Text")
	assertDBTaints(t, `"hi " + a`, c.callee.GetReturnsLst()[0][0], "_obj @ posts_db.post.Text")
}

func TestInlineNestedHelpers(t *testing.T) {
	c := inlineWithSeed(t, "shop.CallOuter", "shop.outer", 0, "posts_db.post")

	// outer is inlined into CallOuter, and inner into the copy of outer
	nested := c.callee.GetAllInlinedGraphs()
	if len(nested) != 1 || nested[0].GetFunctionShortPath() != "shop.inner" {
		t.Fatalf("inlined graphs of outer = %v, want [shop.inner]", nested)
	}
	assertDBTaints(t, "inner returned value", nested[0].GetReturnsLst()[0][0], "_obj @ posts_db.post.ID")
}
