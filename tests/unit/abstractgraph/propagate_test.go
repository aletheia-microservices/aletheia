package abstractgraph_test

import (
	"slices"
	"testing"

	"analyzer/pkg/analysis/system-level/abstractgraph"
	abstractgraphtainter "analyzer/pkg/analysis/system-level/abstractgraph/tainter"
)

// new taints mapped to a primary taint must reach the database call arguments that hold that
// primary taint (or an upper one), at the matching object path
func TestPropagateNewTaintsToDatabaseCallObjects(t *testing.T) {
	graph, edges := newPostNotificationGraph()
	upload := graph.GetNodeByName("UploadService.UploadPost")

	// the queued message holds the primary taint of the push
	msg := newTracedObject("t12", map[string][]*abstractgraph.AbstractTaint{
		"_obj": {primaryWrite("t13", "notifications_queue.notification", edges["push"].GetID())},
	}, nil)
	edges["push"].AddArgument(msg)
	// the rpc argument holds the same taint, but only database calls are updated
	rpcArg := newTracedObject("t0", map[string][]*abstractgraph.AbstractTaint{
		"_obj": {primaryWrite("t13", "notifications_queue.notification", edges["push"].GetID())},
	}, nil)
	edges["rpc"].AddArgument(rpcArg)

	tm := abstractgraph.NewTaintMapping()
	tm.AddIfNotExists(
		*primaryWrite("t13", "notifications_queue.notification.PostID", edges["push"].GetID()),
		*secondaryWrite("t3", "posts_db.post.PostID", edges["insert"].GetID()),
		true, false)

	abstractgraphtainter.PropagateNewTaintsToDatabaseCallObjects(graph, upload, tm, false)

	got := msg.GetTaintsForObjectPath("_obj.PostID")
	if len(got) != 1 || got[0].GetDatabasePath() != "posts_db.post.PostID" || got[0].GetDatabaseCallID() != edges["insert"].GetID() {
		t.Errorf("message _obj.PostID taints = %v, want posts_db.post.PostID from the insert", got)
	}
	if len(rpcArg.GetTaintsForObjectPath("_obj.PostID")) != 0 {
		t.Errorf("rpc arguments must not be updated")
	}

	// propagating again does not duplicate taints
	abstractgraphtainter.PropagateNewTaintsToDatabaseCallObjects(graph, upload, tm, false)
	if got := len(msg.GetTaintsForObjectPath("_obj.PostID")); got != 1 {
		t.Errorf("message _obj.PostID has %d taints after propagating twice, want 1", got)
	}
}

// forward: the taints of a node parameter reach the argument of a call that traces it
func TestPropagateTaintsFromNodeParamsToTracedCallArguments(t *testing.T) {
	graph, edges := newPostNotificationGraph()
	upload := graph.GetNodeByName("UploadService.UploadPost")
	rpc := edges["rpc"]

	// UploadPost(text) passes text to StorePost(text)
	text := newTracedObject("text",
		map[string][]*abstractgraph.AbstractTaint{"_obj": {primaryWrite("t21", "posts_db.post.Text", edges["insert"].GetID())}},
		map[string][]*abstractgraph.AbstractTrace{"_obj": {abstractgraph.NewAbstractTrace("t3", "StorageService.StorePost.text", rpc.GetID())}},
	)
	upload.AddParam(text)
	// the argument already holds a primary taint, so a mapping between both fields is created
	arg := newTracedObject("text", map[string][]*abstractgraph.AbstractTaint{
		"_obj": {primaryWrite("t30", "search_db.index.Text", "SearchService.Index.t30")},
	}, nil)
	rpc.AddArgument(arg)
	// an argument with the same name in another call is not traced by the parameter
	otherArg := newTracedObject("text", nil, nil)
	edges["push"].AddArgument(otherArg)

	tm := abstractgraph.NewTaintMapping()
	abstractgraphtainter.PropagateTaintsToServiceCallObjects(graph, upload, tm, nil, true, false)

	var propagated *abstractgraph.AbstractTaint
	for _, taint := range arg.GetTaintsForObjectPath("_obj") {
		if taint.GetDatabasePath() == "posts_db.post.Text" {
			propagated = taint
		}
	}
	if propagated == nil || !propagated.IsTraced() || propagated.IsPrimary() || propagated.GetT() != "t3" {
		t.Fatalf("rpc argument taints = %v, want traced posts_db.post.Text at t3", arg.GetTaintsForObjectPath("_obj"))
	}
	if otherArg.IsTainted() {
		t.Errorf("arguments of calls that are not traced must not be tainted")
	}
	want := []string{"search_db.index.Text -> posts_db.post.Text"}
	if got := mappingPairs(tm); !slices.Equal(got, want) {
		t.Errorf("mapping = %v, want %v", got, want)
	}
}

// backward: the taints of a call argument (received from the callee) reach the parameter of
// the caller that has the same trace
func TestPropagateTaintsFromCallArgumentsToTracedNodeParams(t *testing.T) {
	graph, edges := newPostNotificationGraph()
	upload := graph.GetNodeByName("UploadService.UploadPost")
	rpc := edges["rpc"]
	trace := func() *abstractgraph.AbstractTrace {
		return abstractgraph.NewAbstractTrace("t3", "StorageService.StorePost.t0", rpc.GetID())
	}

	param := newTracedObject("reqID", nil, map[string][]*abstractgraph.AbstractTrace{"_obj": {trace()}})
	upload.AddParam(param)
	untraced := newTracedObject("username", nil, nil)
	upload.AddParam(untraced)

	// after StorePost is visited, the argument holds the taint of the callee parameter
	arg := newTracedObject("t0",
		map[string][]*abstractgraph.AbstractTaint{"_obj": {primaryWrite("t21", "posts_db.post.ReqID", edges["insert"].GetID())}},
		map[string][]*abstractgraph.AbstractTrace{"_obj": {trace()}},
	)
	rpc.AddArgument(arg)

	abstractgraphtainter.PropagateTaintsToServiceCallObjects(graph, upload, abstractgraph.NewTaintMapping(), rpc, false, false)

	got := param.GetTaintsForObjectPath("_obj")
	if len(got) != 1 || got[0].GetDatabasePath() != "posts_db.post.ReqID" || !got[0].IsTraced() {
		t.Errorf("reqID taints = %v, want traced posts_db.post.ReqID", got)
	}
	if untraced.IsTainted() {
		t.Errorf("parameters without the trace must not be tainted")
	}
}
