package ssagraph_test

import (
	"slices"
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/service-level/ssagraph"
)

const parserSrc = `package shop

type Creator struct { Username string }
type Post struct { ID string; Text string; Creator Creator }

type ServiceImpl struct{}

func (s *ServiceImpl) Upload(id string, text string) (string, error) {
	p := &Post{ID: id, Text: text}
	return s.helper(p), nil
}

func (s *ServiceImpl) helper(p *Post) string { return p.Creator.Username }

func Pick(c bool, a string, b string) string {
	x := a
	if c {
		x = b
	}
	return x
}

func Either(c bool, a string, b string) string {
	if c {
		return a
	}
	return b
}

func Pair(a string) (string, int) { return a, 1 }
func UsePair(a string) string { s, _ := Pair(a); return s }

func Put(m map[string]string, v string) { m["key"] = v }

func Spawn(id string) { go func() { _ = id }() }
`

func nodeNames(nodes []*ssagraph.SSANode) []string {
	var names []string
	for _, node := range nodes {
		names = append(names, node.GetName())
	}
	return names
}

func TestParserCreatesOneGraphPerFunction(t *testing.T) {
	graphs := buildGraphs(t, parserSrc)

	upload := getGraph(t, graphs, "shop.Service.Upload")
	if upload.GetService() != "Service" || upload.GetMethodName() != "Upload" || upload.GetServiceWithMethod() != "Service.Upload" {
		t.Errorf("(*ServiceImpl).Upload: service = %q, method = %q", upload.GetService(), upload.GetMethodName())
	}
	if upload.GetPackageName() != "shop" {
		t.Errorf("package = %q, want shop", upload.GetPackageName())
	}
	// unexported methods and plain functions also get a graph
	getGraph(t, graphs, "shop.Service.helper")
	if pick := getGraph(t, graphs, "shop.Pick"); pick.GetService() != "" || pick.GetMethodName() != "Pick" {
		t.Errorf("plain function: service = %q, method = %q", pick.GetService(), pick.GetMethodName())
	}
}

func TestParserRegistersParametersInOrder(t *testing.T) {
	upload := getGraph(t, buildGraphs(t, parserSrc), "shop.Service.Upload")

	// the receiver is the first parameter
	if got, want := nodeNames(upload.GetParams()), []string{"s", "id", "text"}; !slices.Equal(got, want) {
		t.Errorf("params = %v, want %v", got, want)
	}
	if got := upload.GetIndexOfParameter(upload.GetParamAt(2)); got != 2 {
		t.Errorf("GetIndexOfParameter(text) = %d, want 2", got)
	}
}

func TestParserFieldAccessEdges(t *testing.T) {
	helper := getGraph(t, buildGraphs(t, parserSrc), "shop.Service.helper")

	// p.Creator.Username is two field accesses followed by a load
	p := helper.GetParamAt(1)
	creator := fieldNode(t, helper, p, "Creator")
	username := fieldNode(t, helper, creator, "Username")
	loads := helper.GetEdgesTypedFrom(username, ssagraph.EDGE_LOAD)
	if len(loads) != 1 {
		t.Fatalf("p.Creator.Username has %d loads, want 1", len(loads))
	}
	if ret := helper.GetReturnsLst()[0][0]; ret != loads[0].GetToNode() {
		t.Errorf("returned value = %s, want the loaded username", ret.String())
	}
}

func TestParserStoreEdges(t *testing.T) {
	upload := getGraph(t, buildGraphs(t, parserSrc), "shop.Service.Upload")

	// &Post{ID: id} stores the parameter id at the address of field ID
	id := upload.GetParamAt(1)
	values := upload.GetEdgesTypedFrom(id, ssagraph.EDGE_STORE_VALUE)
	if len(values) != 1 {
		t.Fatalf("id is stored %d times, want 1", len(values))
	}
	store := values[0].GetToNode()
	addresses := upload.GetEdgesTypedTo(store, ssagraph.EDGE_STORE_ADDRESS)
	if len(addresses) != 1 {
		t.Fatalf("store has %d addresses, want 1", len(addresses))
	}
	// the address is the field ID of the new post
	fields := upload.GetEdgesTypedTo(addresses[0].GetFromNode(), ssagraph.EDGE_FIELD)
	if len(fields) != 1 || fields[0].GetParam() != "ID" {
		t.Errorf("stored address = %s, want field ID", addresses[0].GetFromNode().String())
	}
}

func TestParserPhiEdges(t *testing.T) {
	pick := getGraph(t, buildGraphs(t, parserSrc), "shop.Pick")

	phi := pick.GetReturnsLst()[0][0]
	var from []string
	for _, edge := range pick.GetEdgesTypedTo(phi, ssagraph.EDGE_PHI_ON) {
		from = append(from, edge.GetFromNode().GetName())
	}
	slices.Sort(from)
	if want := []string{"a", "b"}; !slices.Equal(from, want) {
		t.Errorf("phi values = %v, want %v", from, want)
	}
}

func TestParserOneReturnListPerReturnStatement(t *testing.T) {
	either := getGraph(t, buildGraphs(t, parserSrc), "shop.Either")

	var returned []string
	for _, rets := range either.GetReturnsLst() {
		returned = append(returned, nodeNames(rets)...)
	}
	slices.Sort(returned)
	if want := []string{"a", "b"}; len(either.GetReturnsLst()) != 2 || !slices.Equal(returned, want) {
		t.Errorf("returns = %v, want one list per return statement with %v", either.GetReturnsLst(), want)
	}
}

func TestParserExtractEdges(t *testing.T) {
	usePair := getGraph(t, buildGraphs(t, parserSrc), "shop.UsePair")

	// s, _ := Pair(a) extracts both results of the call
	call := usePair.GetEdgesTypedFrom(usePair.GetParamAt(0), ssagraph.EDGE_ARG_ON_CALL)[0].GetToNode()
	extracts := usePair.GetEdgesTypedFrom(call, ssagraph.EDGE_EXTRACT)
	if len(extracts) != 2 || extracts[0].GetIndex() != 0 || extracts[1].GetIndex() != 1 {
		t.Fatalf("extract edges = %v, want indexes 0 and 1", extracts)
	}
	if ret := usePair.GetReturnsLst()[0][0]; ret != extracts[0].GetToNode() {
		t.Errorf("returned value = %s, want the first result", ret.String())
	}
}

func TestParserMapUpdateUsesConstantKey(t *testing.T) {
	put := getGraph(t, buildGraphs(t, parserSrc), "shop.Put")

	m, v := put.GetParamAt(0), put.GetParamAt(1)
	update := put.GetFirstEdgeTypedFrom(m, ssagraph.EDGE_MAP_UPDATE)
	value := put.GetFirstEdgeTypedFrom(v, ssagraph.EDGE_MAP_VALUE)
	if update == nil || value == nil || update.GetToNode() != value.GetToNode() {
		t.Fatalf("m and v must point to the same map update")
	}
	if update.GetParam() != "key" || value.GetParam() != "key" {
		t.Errorf("map update index = (%q, %q), want the constant key", update.GetParam(), value.GetParam())
	}
}

func TestParserGoRoutines(t *testing.T) {
	graphs := buildGraphs(t, parserSrc)

	closure := getGraph(t, graphs, "shop.Spawn$1")
	if !closure.IsGoRoutine() {
		t.Errorf("closure started with go must be marked as a go routine")
	}
	// captured variables are free variables, not parameters
	if len(closure.GetParams()) != 0 || len(closure.GetFreeVars()) != 1 || closure.GetFreeVars()[0].GetName() != "id" {
		t.Errorf("closure params = %v, free vars = %v", nodeNames(closure.GetParams()), nodeNames(closure.GetFreeVars()))
	}
	if getGraph(t, graphs, "shop.Spawn").IsGoRoutine() {
		t.Errorf("the function starting the go routine must not be marked as a go routine")
	}
}

// the tainter registers static calls to other functions as method calls
func TestTainterRegistersMethodCalls(t *testing.T) {
	graphs := buildTaintedGraphs(t, parserSrc)

	upload := getGraph(t, graphs, "shop.Service.Upload")
	if len(upload.GetMethodCalls()) != 1 || len(upload.GetAllCalls()) != 1 {
		t.Fatalf("Upload has %d method calls, want 1", len(upload.GetMethodCalls()))
	}
	call := upload.GetMethodCalls()[0]
	if call.GetFuncShortPath() != "shop.Service.helper" || call.GetMethod() != "helper" {
		t.Errorf("call = %s (%s)", call.GetFuncShortPath(), call.GetMethod())
	}
	if call.GetID() != "Service.Upload."+call.GetT() {
		t.Errorf("call id = %q, want Service.Upload.%s", call.GetID(), call.GetT())
	}
	// the receiver is the first argument
	if got := nodeNames(call.GetArguments()); len(got) != 2 || got[0] != "s" {
		t.Errorf("arguments = %v, want receiver first", got)
	}
	if len(call.GetReturns()) != 1 {
		t.Errorf("returns = %v, want the call value", call.GetReturns())
	}

	// calls with several results return the extracted values
	pair := getGraph(t, graphs, "shop.UsePair").GetMethodCalls()[0]
	if len(pair.GetReturns()) != 2 || pair.GetReturnAt(0) != getGraph(t, graphs, "shop.UsePair").GetReturnsLst()[0][0] {
		t.Errorf("Pair returns = %v, want both extracted values", nodeNames(pair.GetReturns()))
	}

	// go routines are registered once, with the captured variables as bindings
	spawn := getGraph(t, graphs, "shop.Spawn")
	if len(spawn.GetMethodCalls()) != 1 {
		t.Fatalf("Spawn has %d method calls, want 1", len(spawn.GetMethodCalls()))
	}
	goCall := spawn.GetMethodCalls()[0]
	if goCall.GetID() != "shop.Spawn$1" || goCall.GetBindAt(0) == nil {
		t.Errorf("go routine call id = %q", goCall.GetID())
	}
}
