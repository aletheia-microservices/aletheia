package ssa

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/tests/runner"
)

// SocialGraphService.InsertUser inserts an empty social graph entry for the user: the user id is
// tainted with the field it was assigned to, and the entry with all of its fields
func TestSocialNetworkWriteNewEntry(t *testing.T) {
	a := runner.Get(t, "dsb_socialnetwork")
	graph := getSSAGraph(t, a, "socialnetwork.SocialGraphService.InsertUser")

	call := getDatabaseCall(t, graph, "socialgraph_db.socialgraph", "InsertOne")
	assertDatabaseTaint(t, getParam(t, graph, "userID"), dbTaint{objpath: "_obj", dbpath: "socialgraph_db.socialgraph.UserID", op: common.OP_WRITE, t: call.GetT()})
	assertNotTaintedBy(t, getParam(t, graph, "reqID"), "socialgraph_db")

	entry := call.GetArguments()[0]
	for _, field := range []string{"", ".UserID", ".Followers", ".Followees"} {
		assertDatabaseTaint(t, entry, dbTaint{objpath: "_obj" + field, dbpath: "socialgraph_db.socialgraph" + field, op: common.OP_WRITE})
	}
}

// UserIDService.GetUserId looks up the user id in the cache, reads it from the database on a miss
// and then stores it in the cache: the username is the key of the three calls
func TestSocialNetworkCacheAsideRead(t *testing.T) {
	a := runner.Get(t, "dsb_socialnetwork")
	graph := getSSAGraph(t, a, "socialnetwork.UserIDService.GetUserId")

	get := getDatabaseCall(t, graph, "user_cache.UserID", "Get")
	find := getDatabaseCall(t, graph, "user_db.user", "FindOne")
	put := getDatabaseCall(t, graph, "user_cache.UserID", "Put")
	username := getParam(t, graph, "username")
	assertDatabaseTaint(t, username, dbTaint{objpath: "_obj", dbpath: "user_cache.UserID.Key", op: common.OP_READ, readKey: true, t: get.GetT()})
	assertDatabaseTaint(t, username, dbTaint{objpath: "_obj", dbpath: "user_db.user.Username", op: common.OP_READ, readKey: true, t: find.GetT()})
	assertDatabaseTaint(t, username, dbTaint{objpath: "_obj", dbpath: "user_cache.UserID.Key", op: common.OP_WRITE, readKey: true, t: put.GetT()})

	// the user id read from the database is the value stored in the cache
	user := find.GetArguments()[2]
	assertDatabaseTaint(t, user, dbTaint{objpath: "_obj.UserID", dbpath: "user_db.user.UserID", op: common.OP_READ, readVal: true, t: find.GetT()})
	assertDatabaseTaint(t, user, dbTaint{objpath: "_obj.UserID", dbpath: "user_cache.UserID.Value", op: common.OP_WRITE, readVal: true, t: put.GetT()})
}

// PostStorageService.StorePost inserts the post it receives as is
func TestSocialNetworkWriteParameterObject(t *testing.T) {
	a := runner.Get(t, "dsb_socialnetwork")
	graph := getSSAGraph(t, a, "socialnetwork.PostStorageService.StorePost")

	call := getDatabaseCall(t, graph, "post_db.post", "InsertOne")
	assertDatabaseTaint(t, getParam(t, graph, "post"), dbTaint{objpath: "_obj", dbpath: "post_db.post", op: common.OP_WRITE, t: call.GetT()})
}
