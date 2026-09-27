package ssa

import (
	"strings"
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/internal/analysis/service-level/ssagraph"
	"github.com/aletheia-microservices/aletheia/tests/runner"
)

// StorePost builds a Post struct and writes it with InsertOne: the object and all of its (nested)
// fields must be tainted with the corresponding database fields, and the taints must be propagated
// upwards to the function parameters that were assigned to those fields
func TestPostNotificationTaintWriteObjectFields(t *testing.T) {
	a := runner.Get(t, "postnotification")
	graph := getSSAGraph(t, a, "postnotification.StorageService.StorePost")

	if got := len(graph.GetDatabaseCalls()); got != 1 {
		t.Fatalf("StorePost has %d database calls, want 1", got)
	}
	call := getDatabaseCall(t, graph, "posts_db.post", "InsertOne")
	if call.GetOpType() != common.OP_WRITE {
		t.Fatalf("InsertOne op type = %s, want write", common.OperationTypeToString(call.GetOpType()))
	}
	if got := len(call.GetArguments()); got != 1 {
		t.Fatalf("InsertOne has %d tracked arguments, want 1 (the post)", got)
	}

	post := call.GetArguments()[0]
	for _, field := range []string{"", ".PostID", ".ReqID", ".Text", ".Timestamp", ".Mentions", ".Mentions[*]", ".Creator", ".Creator.Username"} {
		assertDatabaseTaint(t, post, dbTaint{objpath: "_obj" + field, dbpath: "posts_db.post" + field, op: common.OP_WRITE, t: call.GetT()})
	}

	// fields are propagated upwards to the values they were assigned from
	assertDatabaseTaint(t, getParam(t, graph, "text"), dbTaint{objpath: "_obj", dbpath: "posts_db.post.Text", op: common.OP_WRITE})
	assertDatabaseTaint(t, getParam(t, graph, "reqID"), dbTaint{objpath: "_obj", dbpath: "posts_db.post.ReqID", op: common.OP_WRITE})
	// ... but never to unrelated fields
	assertNotTaintedBy(t, getParam(t, graph, "text"), "posts_db.post.ReqID")
	assertNotTaintedBy(t, getParam(t, graph, "ctx"), "posts_db")
}

// ReadPost filters by PostID (read key) and decodes the result into a Post (read value)
func TestPostNotificationTaintReadKeyAndValue(t *testing.T) {
	a := runner.Get(t, "postnotification")
	graph := getSSAGraph(t, a, "postnotification.StorageService.ReadPost")

	call := getDatabaseCall(t, graph, "posts_db.post", "FindOne")
	if call.GetOpType() != common.OP_READ {
		t.Fatalf("FindOne op type = %s, want read", common.OperationTypeToString(call.GetOpType()))
	}

	postID := getParam(t, graph, "postID")
	assertDatabaseTaint(t, postID, dbTaint{objpath: "_obj", dbpath: "posts_db.post.PostID", op: common.OP_READ, readKey: true, t: call.GetT()})
	assertNotTaintedBy(t, getParam(t, graph, "reqID"), "posts_db")

	var foundValue bool
	for _, arg := range call.GetArguments() {
		if taint := findDatabaseTaint(arg, dbTaint{objpath: "_obj", dbpath: "posts_db.post", op: common.OP_READ}); taint != nil {
			foundValue = true
			if !taint.IsReadValue() || taint.IsReadKey() {
				t.Errorf("decoded post must be tainted as read value")
			}
		}
	}
	if !foundValue {
		t.Errorf("FindOne has no argument tainted as the decoded post")
	}
}

func TestPostNotificationTaintDeleteKey(t *testing.T) {
	a := runner.Get(t, "postnotification")
	graph := getSSAGraph(t, a, "postnotification.StorageService.DeletePost")

	call := getDatabaseCall(t, graph, "posts_db.post", "DeleteOne")
	if call.GetOpType() != common.OP_DELETE {
		t.Fatalf("DeleteOne op type = %s, want delete", common.OperationTypeToString(call.GetOpType()))
	}
	assertDatabaseTaint(t, getParam(t, graph, "postID"), dbTaint{objpath: "_obj", dbpath: "posts_db.post.PostID", op: common.OP_DELETE, readKey: true})
}

// NotifyService pops a message from the queue and forwards its fields to StorageService.ReadPost:
// the RPC arguments must carry both the queue taints and the service taint of the RPC
func TestPostNotificationQueueReadFlowsIntoRPCArguments(t *testing.T) {
	a := runner.Get(t, "postnotification")
	graph := getSSAGraph(t, a, "postnotification.NotifyService.Run")

	pop := getDatabaseCall(t, graph, "notifications_queue.notification", "Pop")
	if pop.GetOpType() != common.OP_READ {
		t.Fatalf("Pop op type = %s, want read", common.OperationTypeToString(pop.GetOpType()))
	}
	msg := pop.GetArguments()[0]
	assertDatabaseTaint(t, msg, dbTaint{objpath: "_obj", dbpath: "notifications_queue.notification", op: common.OP_READ, readVal: true})
	assertDatabaseTaint(t, msg, dbTaint{objpath: "_obj.PostID", dbpath: "notifications_queue.notification.PostID", op: common.OP_READ, readVal: true})

	rpc := getServiceCall(t, graph, "StorageService.ReadPost")
	if got := len(rpc.GetArguments()); got != 2 {
		t.Fatalf("ReadPost rpc has %d arguments, want 2 (reqID, postID)", got)
	}
	reqID, postID := rpc.GetArguments()[0], rpc.GetArguments()[1]
	assertDatabaseTaint(t, reqID, dbTaint{objpath: "_obj", dbpath: "notifications_queue.notification.ReqID", op: common.OP_READ, readVal: true})
	assertDatabaseTaint(t, postID, dbTaint{objpath: "_obj", dbpath: "notifications_queue.notification.PostID", op: common.OP_READ, readVal: true})
	assertServiceTaint(t, postID, "_obj", rpc)
	// the rpc taint is also propagated back to the field of the queue message
	assertServiceTaint(t, msg, "_obj.PostID", rpc)
}

// UploadPost pushes a message built from the return of StorageService.StorePost:
// the message field must carry a service taint pointing to that RPC
func TestPostNotificationRPCReturnFlowsIntoQueueWrite(t *testing.T) {
	a := runner.Get(t, "postnotification")
	graph := getSSAGraph(t, a, "postnotification.UploadService.UploadPost")

	rpc := getServiceCall(t, graph, "StorageService.StorePost")
	if got := len(rpc.GetReturns()); got != 2 {
		t.Fatalf("StorePost rpc has %d returns, want 2 (postID, err)", got)
	}
	push := getDatabaseCall(t, graph, "notifications_queue.notification", "Push")
	msg := push.GetArguments()[0]

	assertDatabaseTaint(t, msg, dbTaint{objpath: "_obj.PostID", dbpath: "notifications_queue.notification.PostID", op: common.OP_WRITE})
	assertServiceTaint(t, msg, "_obj.PostID", rpc)
	assertServiceTaint(t, msg, "_obj.ReqID", rpc)

	// reqID is both passed to StorePost and written to the queue
	reqID := rpc.GetArguments()[0]
	assertDatabaseTaint(t, reqID, dbTaint{objpath: "_obj", dbpath: "notifications_queue.notification.ReqID", op: common.OP_WRITE})
	assertServiceTaint(t, reqID, "_obj", rpc)
}

func TestPostNotificationCallsAreRegisteredInProgramOrder(t *testing.T) {
	a := runner.Get(t, "postnotification")
	graph := getSSAGraph(t, a, "postnotification.UploadService.UploadPost")

	var kinds []string
	for _, call := range graph.GetAllCalls() {
		switch c := call.(type) {
		case *ssagraph.ServiceCall:
			kinds = append(kinds, "rpc:"+c.GetServiceWithMethod())
		case *ssagraph.DatabaseCall:
			kinds = append(kinds, "db:"+c.GetDatabasePath()+"."+c.GetMethod())
		}
	}
	want := []string{"rpc:StorageService.StorePost", "db:notifications_queue.notification.Push"}
	if strings.Join(kinds, ",") != strings.Join(want, ",") {
		t.Errorf("calls = %v, want %v", kinds, want)
	}
}
