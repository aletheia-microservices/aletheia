package abstractcallgraph

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/tests/runner"
)

// dsb_socialnetwork (DeathStarBench): a Wrk2APIService composes posts, follows users and reads
// timelines; posts are fanned out to the user and home timelines

func TestSocialNetworkEntrypoints(t *testing.T) {
	g := runner.Get(t, "dsb_socialnetwork").AbsGraph
	assertList(t, "entrypoints", entrypoints(g), []string{
		"Wrk2APIService.ComposePost",
		"Wrk2APIService.Follow",
		"Wrk2APIService.ReadHomeTimeline",
		"Wrk2APIService.ReadUserTimeline",
		"Wrk2APIService.Register",
		"Wrk2APIService.Unfollow",
	})
}

func TestSocialNetworkCalls(t *testing.T) {
	g := runner.Get(t, "dsb_socialnetwork").AbsGraph
	assertList(t, "rpc calls", rpcCalls(g), []string{
		"ComposePostService.ComposePost -> HomeTimelineService.WriteHomeTimeline",
		"ComposePostService.ComposePost -> MediaService.ComposeMedia",
		"ComposePostService.ComposePost -> PostStorageService.StorePost",
		"ComposePostService.ComposePost -> TextService.ComposeText",
		"ComposePostService.ComposePost -> UniqueIdService.ComposeUniqueId",
		"ComposePostService.ComposePost -> UserService.ComposeCreatorWithUserId",
		"ComposePostService.ComposePost -> UserTimelineService.WriteUserTimeline",
		"HomeTimelineService.ReadHomeTimeline -> PostStorageService.ReadPosts",
		"HomeTimelineService.WriteHomeTimeline -> SocialGraphService.GetFollowers",
		"SocialGraphService.FollowWithUsername -> UserIDService.GetUserId",
		"SocialGraphService.FollowWithUsername -> UserIDService.GetUserId",
		"SocialGraphService.UnfollowWithUsername -> UserIDService.GetUserId",
		"SocialGraphService.UnfollowWithUsername -> UserIDService.GetUserId",
		"TextService.ComposeText -> UrlShortenService.ComposeUrls",
		"TextService.ComposeText -> UserMentionService.ComposeUserMentions",
		"UserService.RegisterUserWithId -> SocialGraphService.InsertUser",
		"UserTimelineService.ReadUserTimeline -> PostStorageService.ReadPosts",
		"Wrk2APIService.ComposePost -> ComposePostService.ComposePost",
		"Wrk2APIService.Follow -> SocialGraphService.Follow",
		"Wrk2APIService.Follow -> SocialGraphService.FollowWithUsername",
		"Wrk2APIService.ReadHomeTimeline -> HomeTimelineService.ReadHomeTimeline",
		"Wrk2APIService.ReadUserTimeline -> UserTimelineService.ReadUserTimeline",
		"Wrk2APIService.Register -> UserService.RegisterUserWithId",
		"Wrk2APIService.Unfollow -> SocialGraphService.Unfollow",
		"Wrk2APIService.Unfollow -> SocialGraphService.UnfollowWithUsername",
	})
	assertList(t, "database calls", databaseCalls(g, nil), []string{
		"read HomeTimelineService.ReadHomeTimeline -> hometimeline_cache.*.Get",
		"read HomeTimelineService.WriteHomeTimeline -> hometimeline_cache.*.Get",
		"read PostStorageService.ReadPosts -> post_cache.*.Mget",
		"read PostStorageService.ReadPosts -> post_db.post.FindMany",
		"read SocialGraphService.Follow -> socialgraph_cache.followees.Get",
		"read SocialGraphService.Follow -> socialgraph_cache.followers.Get",
		"read SocialGraphService.FollowWithUsername -> socialgraph_cache.followees.Get",
		"read SocialGraphService.FollowWithUsername -> socialgraph_cache.followers.Get",
		"read SocialGraphService.GetFollowers -> socialgraph_cache.followers.Get",
		"read SocialGraphService.GetFollowers -> socialgraph_db.socialgraph.FindOne",
		"read UserIDService.GetUserId -> user_cache.UserID.Get",
		"read UserIDService.GetUserId -> user_db.user.FindOne",
		"read UserService.RegisterUserWithId -> user_db.user.FindOne",
		"read UserTimelineService.ReadUserTimeline -> usertimeline_cache.*.Get",
		"read UserTimelineService.ReadUserTimeline -> usertimeline_db.usertimeline.FindOne",
		"read UserTimelineService.WriteUserTimeline -> usertimeline_cache.*.Get",
		"read UserTimelineService.WriteUserTimeline -> usertimeline_db.usertimeline.FindMany",
		"update SocialGraphService.Follow -> socialgraph_db.socialgraph.UpdateOne",
		"update SocialGraphService.Follow -> socialgraph_db.socialgraph.UpdateOne",
		"update SocialGraphService.FollowWithUsername -> socialgraph_db.socialgraph.UpdateOne",
		"update SocialGraphService.FollowWithUsername -> socialgraph_db.socialgraph.UpdateOne",
		"update SocialGraphService.Unfollow -> socialgraph_db.socialgraph.UpdateOne",
		"update SocialGraphService.Unfollow -> socialgraph_db.socialgraph.UpdateOne",
		"update UserTimelineService.WriteUserTimeline -> usertimeline_db.usertimeline.UpdateMany",
		"write HomeTimelineService.WriteHomeTimeline -> hometimeline_cache.*.Put",
		"write PostStorageService.ReadPosts -> post_cache.*.Put",
		"write PostStorageService.StorePost -> post_db.post.InsertOne",
		"write SocialGraphService.Follow -> socialgraph_cache.followees.Put",
		"write SocialGraphService.Follow -> socialgraph_cache.followers.Put",
		"write SocialGraphService.FollowWithUsername -> socialgraph_cache.followees.Put",
		"write SocialGraphService.FollowWithUsername -> socialgraph_cache.followers.Put",
		"write SocialGraphService.GetFollowers -> socialgraph_cache.followers.Put",
		"write SocialGraphService.InsertUser -> socialgraph_db.socialgraph.InsertOne",
		"write UserIDService.GetUserId -> user_cache.UserID.Put",
		"write UserService.RegisterUserWithId -> user_db.user.InsertOne",
		"write UserTimelineService.ReadUserTimeline -> usertimeline_cache.*.Put",
		"write UserTimelineService.WriteUserTimeline -> usertimeline_cache.*.Put",
		"write UserTimelineService.WriteUserTimeline -> usertimeline_db.usertimeline.InsertOne",
	})
}

func TestSocialNetworkCounts(t *testing.T) {
	assertCounts(t, runner.Get(t, "dsb_socialnetwork").AbsGraph, map[string]int{
		"service nodes": 27, "database nodes": 10,
		"entry": 6, "rpc": 25, "read": 17, "write": 14, "update": 7,
	})
}

func TestSocialNetworkDatabaseOwners(t *testing.T) {
	assertOwners(t, databaseOwners(runner.Get(t, "dsb_socialnetwork").AbsGraph, nil), map[string][]string{
		"hometimeline_cache.*":         {"HomeTimelineService"},
		"post_cache.*":                 {"PostStorageService"},
		"post_db.post":                 {"PostStorageService"},
		"socialgraph_cache.followees":  {"SocialGraphService"},
		"socialgraph_cache.followers":  {"SocialGraphService"},
		"socialgraph_db.socialgraph":   {"SocialGraphService"},
		"user_cache.UserID":            {"UserIDService"},
		"usertimeline_cache.*":         {"UserTimelineService"},
		"usertimeline_db.usertimeline": {"UserTimelineService"},
		// user_db is wired to both services
		"user_db.user": {"UserIDService", "UserService"},
	})
}

// ComposePost stores the post with the id created by UniqueIdService
func TestSocialNetworkComposePostUsesUniqueId(t *testing.T) {
	g := runner.Get(t, "dsb_socialnetwork").AbsGraph

	uniqueID := getEdge(t, g, "ComposePostService.ComposePost", "UniqueIdService.ComposeUniqueId", "ComposeUniqueId")
	store := getEdge(t, g, "ComposePostService.ComposePost", "PostStorageService.StorePost", "StorePost")

	assertTrace(t, store.GetArgumentAt(1), "_obj.PostID", uniqueID, "")
	if !isBefore(uniqueID.GetT(), store.GetT()) {
		t.Errorf("ComposeUniqueId (t=%s) must happen before StorePost (t=%s)", uniqueID.GetT(), store.GetT())
	}
	// on the callee side, the post is written to post_db as a whole
	assertPrimaryTaint(t, g.GetNodeByName("PostStorageService.StorePost").GetParameterByNameIfExists("post"), "_obj", "post_db.post", common.OP_WRITE)
}

// ComposeUserMentions returns before any of its calls, so it has none
func TestSocialNetworkUnreachableCodeHasNoCalls(t *testing.T) {
	g := runner.Get(t, "dsb_socialnetwork").AbsGraph
	if edges := g.GetEdgesFromNode(g.GetNodeByName("UserMentionService.ComposeUserMentions")); len(edges) != 0 {
		t.Errorf("ComposeUserMentions has %d calls, want none", len(edges))
	}
}

func TestSocialNetworkUnfollowWithUsernameUpdatesGraph(t *testing.T) {
	t.Skip("known bug: database calls two helper calls deep are dropped (see LIKELY_BUGS.md #17), e.g. " +
		"SocialGraphService.UnfollowWithUsername -> Unfollow -> go routines")

	g := runner.Get(t, "dsb_socialnetwork").AbsGraph
	if !hasEdge(g, "SocialGraphService.UnfollowWithUsername", "socialgraph_db.socialgraph", "UpdateOne") {
		t.Errorf("missing update SocialGraphService.UnfollowWithUsername -> socialgraph_db.socialgraph.UpdateOne")
	}
}

func TestSocialNetworkComposeUrlsWritesUrls(t *testing.T) {
	t.Skip("known bug: NoSQLCollection.InsertMany is not a recognized database call (see LIKELY_BUGS.md #15), so " +
		"UrlShortenService.ComposeUrls never writes to urlshorten_db")

	g := runner.Get(t, "dsb_socialnetwork").AbsGraph
	if !hasEdge(g, "UrlShortenService.ComposeUrls", "urlshorten_db.urlshorten", "InsertMany") {
		t.Errorf("missing write UrlShortenService.ComposeUrls -> urlshorten_db.urlshorten.InsertMany")
	}
}
