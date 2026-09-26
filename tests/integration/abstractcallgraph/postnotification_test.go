package abstractcallgraph

import (
	"slices"
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/tests/runner"
)

// postnotification: UploadService stores a post with StorageService and pushes a notification to a
// queue, which NotifyService pops to read the post back

func TestPostNotificationEntrypoints(t *testing.T) {
	g := runner.Get(t, "postnotification").AbsGraph
	assertList(t, "entrypoints", entrypoints(g), []string{
		"NotifyService.Run",
		"UploadService.DeletePost",
		"UploadService.UploadPost",
	})
	if got := g.ComputeAndGetNumCallGraphs(); got != 3 {
		t.Errorf("call graphs = %d, want 3", got)
	}
}

func TestPostNotificationCalls(t *testing.T) {
	g := runner.Get(t, "postnotification").AbsGraph
	assertList(t, "rpc calls", rpcCalls(g), []string{
		"NotifyService.Run -> StorageService.ReadPost",
		"UploadService.DeletePost -> StorageService.DeletePost",
		"UploadService.UploadPost -> StorageService.StorePost",
	})
	assertList(t, "database calls", databaseCalls(g, nil), []string{
		"delete StorageService.DeletePost -> posts_db.post.DeleteOne",
		"read NotifyService.Run -> notifications_queue.notification.Pop",
		"read StorageService.ReadPost -> posts_db.post.FindOne",
		"write StorageService.StorePost -> posts_db.post.InsertOne",
		"write UploadService.UploadPost -> notifications_queue.notification.Push",
	})
	// the rpc count includes the entrypoints
	if g.GetRPCCount() != 6 || g.GetDBAccessCount() != 5 {
		t.Errorf("rpc count = %d, db access count = %d, want 6 and 5", g.GetRPCCount(), g.GetDBAccessCount())
	}
}

func TestPostNotificationDatabaseOwners(t *testing.T) {
	g := runner.Get(t, "postnotification").AbsGraph
	assertOwners(t, databaseOwners(g, nil), map[string][]string{
		"posts_db.post": {"StorageService"},
		// the queue is shared by the producer and the consumer
		"notifications_queue.notification": {"NotifyService", "UploadService"},
	})
}

// the node parameters and edge arguments of the abstract call graph are built from the SSA taints
func TestPostNotificationTaintsFromSSA(t *testing.T) {
	g := runner.Get(t, "postnotification").AbsGraph

	storePost := g.GetNodeByName("StorageService.StorePost")
	if got := len(storePost.GetParams()); got != 2 {
		t.Fatalf("StorePost has %d params, want 2 (context and receiver are excluded)", got)
	}
	assertPrimaryTaint(t, storePost.GetParamAt(0), "_obj", "posts_db.post.ReqID", common.OP_WRITE)
	assertPrimaryTaint(t, storePost.GetParamAt(1), "_obj", "posts_db.post.Text", common.OP_WRITE)

	insert := getEdge(t, g, "StorageService.StorePost", "posts_db.post", "InsertOne")
	post := insert.GetArgumentAt(0)
	for _, field := range []string{"", ".PostID", ".Creator.Username", ".Mentions[*]"} {
		taint := assertPrimaryTaint(t, post, "_obj"+field, "posts_db.post"+field, common.OP_WRITE)
		if taint != nil && taint.GetDatabaseCallID() != insert.GetID() {
			t.Errorf("taint %s refers to call %s, want %s", taint.GetDatabasePath(), taint.GetDatabaseCallID(), insert.GetID())
		}
	}
	if got := post.GetAffectedDatabaseFieldsForCall(insert.GetID()); !slices.Contains(got, "posts_db.post.PostID") {
		t.Errorf("fields affected by InsertOne = %v, want posts_db.post.PostID", got)
	}

	readPost := g.GetNodeByName("StorageService.ReadPost")
	if taint := assertPrimaryTaint(t, readPost.GetParamAt(1), "_obj", "posts_db.post.PostID", common.OP_READ); taint != nil && !taint.IsReadKey() {
		t.Errorf("ReadPost(postID) must be a read key")
	}
	// username is never stored
	assertNotTainted(t, g.GetNodeByName("UploadService.UploadPost").GetParameterByNameIfExists("username"))
}

// the queued message links the post created by StorePost to the post read by NotifyService
func TestPostNotificationQueueMessageTracesPost(t *testing.T) {
	g := runner.Get(t, "postnotification").AbsGraph

	storePost := getEdge(t, g, "UploadService.UploadPost", "StorageService.StorePost", "StorePost")
	push := getEdge(t, g, "UploadService.UploadPost", "notifications_queue.notification", "Push")
	msg := push.GetArgumentAt(0)
	assertPrimaryTaint(t, msg, "_obj.PostID", "notifications_queue.notification.PostID", common.OP_WRITE)
	assertTrace(t, msg, "_obj.PostID", storePost, "")

	readPost := getEdge(t, g, "NotifyService.Run", "StorageService.ReadPost", "ReadPost")
	assertPrimaryTaint(t, readPost.GetArgumentAt(1), "_obj", "notifications_queue.notification.PostID", common.OP_READ)
}
