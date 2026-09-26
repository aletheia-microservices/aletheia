package abstractgraph_test

import (
	"testing"

	"analyzer/pkg/analysis/common"
	"analyzer/pkg/analysis/system-level/abstractgraph"
)

// newPostNotificationGraph builds a small part of the postnotification abstract call graph:
// client -> UploadService.UploadPost -> StorageService.StorePost -> posts_db.post
// and UploadService.UploadPost -> notifications_queue.notification
func newPostNotificationGraph() (*abstractgraph.AbstractCallGraph, map[string]*abstractgraph.AbstractEdge) {
	graph := abstractgraph.NewAbstractCallGraph(nil)
	client := abstractgraph.NewAbstractNode("client", abstractgraph.NODE_CLIENT, "", "", "", "")
	upload := newServiceNode("UploadService", "UploadPost")
	store := newServiceNode("StorageService", "StorePost")
	posts := newDatabaseNode("posts_db", "post")
	queue := newDatabaseNode("notifications_queue", "notification")
	for _, node := range []*abstractgraph.AbstractNode{client, upload, store, posts, queue} {
		graph.AddNode(node.GetName(), node)
	}

	edges := map[string]*abstractgraph.AbstractEdge{
		"entry":  abstractgraph.NewAbstractEdge("", "postnotification.UploadService.UploadPost", "UploadPost", client, upload, common.OP_UNDEFINED, abstractgraph.EDGE_SERVICE_ENTRYPOINT),
		"rpc":    abstractgraph.NewAbstractEdge("t3", "UploadService.UploadPost.t3", "StorePost", upload, store, common.OP_UNDEFINED, abstractgraph.EDGE_SERVICE_RPC),
		"insert": abstractgraph.NewAbstractEdge("t21", "StorageService.StorePost.t21", "InsertOne", store, posts, common.OP_WRITE, abstractgraph.EDGE_DATABASE_CALL),
		"push":   abstractgraph.NewAbstractEdge("t13", "UploadService.UploadPost.t13", "Push", upload, queue, common.OP_WRITE, abstractgraph.EDGE_DATABASE_CALL),
	}
	for _, name := range []string{"entry", "rpc", "insert", "push"} {
		graph.AddEdge(edges[name])
	}
	return graph, edges
}

func TestAbstractCallGraphNodes(t *testing.T) {
	graph, _ := newPostNotificationGraph()

	if got := len(graph.GetNodes()); got != 5 {
		t.Errorf("graph has %d nodes, want 5", got)
	}
	if node := graph.GetNodeByNameIfExists("StorageService.StorePost"); node == nil || node.GetMethod() != "StorePost" {
		t.Errorf("GetNodeByNameIfExists(StorageService.StorePost) = %v", node)
	}
	if graph.GetNodeByName("posts_db.post").GetNodeType() != abstractgraph.NODE_DATABASE {
		t.Errorf("posts_db.post must be a database node")
	}
	if graph.GetNodeByNameIfExists("missing") != nil {
		t.Errorf("missing node must return nil")
	}
}

func TestAbstractCallGraphEdges(t *testing.T) {
	graph, edges := newPostNotificationGraph()
	upload := graph.GetNodeByName("UploadService.UploadPost")

	// edges are returned in insertion order
	from := graph.GetEdgesFromNode(upload)
	if len(from) != 2 || from[0] != edges["rpc"] || from[1] != edges["push"] {
		t.Errorf("edges from UploadPost = %v", from)
	}
	to := graph.GetEdgesToNode(upload)
	if len(to) != 1 || to[0] != edges["entry"] {
		t.Errorf("edges to UploadPost = %v", to)
	}
	if got := graph.GetEdgesFromNode(graph.GetNodeByName("posts_db.post")); len(got) != 0 {
		t.Errorf("database nodes must not have outgoing edges: %v", got)
	}
	if got := len(graph.GetEdges()); got != 4 {
		t.Errorf("graph has %d edges, want 4", got)
	}
}

func TestAbstractCallGraphNumCallGraphs(t *testing.T) {
	graph, _ := newPostNotificationGraph()
	if got := graph.ComputeAndGetNumCallGraphs(); got != 1 {
		t.Errorf("call graphs = %d, want 1 (one entrypoint)", got)
	}

	// one call graph per entrypoint edge from the client
	client := graph.GetNodeByName("client")
	notify := newServiceNode("NotifyService", "Run")
	graph.AddNode(notify.GetName(), notify)
	graph.AddEdge(abstractgraph.NewAbstractEdge("", "postnotification.NotifyService.Run", "Run", client, notify, common.OP_UNDEFINED, abstractgraph.EDGE_SERVICE_ENTRYPOINT))
	if got := graph.ComputeAndGetNumCallGraphs(); got != 2 {
		t.Errorf("call graphs = %d, want 2", got)
	}

	if got := abstractgraph.NewAbstractCallGraph(nil).ComputeAndGetNumCallGraphs(); got != 0 {
		t.Errorf("graph without client has %d call graphs, want 0", got)
	}
}
