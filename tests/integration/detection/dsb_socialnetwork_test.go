package detection

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/tests/runner"
)

// timelines (in the databases and in the caches) store the ids of posts written by PostStorageService,
// and are always written after the post, so the references are mandatory
func TestSocialNetworkTimelinesReferencePosts(t *testing.T) {
	a := runner.Get(t, "dsb_socialnetwork")
	assertConstraint(t, a, "FOREIGN_KEY usertimeline_db.usertimeline.Posts[*].PostID REFERENCES post_db.post.PostID [MANDATORY]", true)
	assertConstraint(t, a, "FOREIGN_KEY usertimeline_cache.*.Value[*].PostID REFERENCES post_db.post.PostID [MANDATORY]", true)
	assertConstraint(t, a, "FOREIGN_KEY hometimeline_cache.*.Value[*].PostID REFERENCES post_db.post.PostID [MANDATORY]", true)
}

// the social graph stores the ids of the users who follow each other
func TestSocialNetworkSocialGraphReferencesUsers(t *testing.T) {
	a := runner.Get(t, "dsb_socialnetwork")
	assertConstraint(t, a, "FOREIGN_KEY socialgraph_db.socialgraph.UserID REFERENCES user_db.user.UserID", true)
	assertConstraint(t, a, "FOREIGN_KEY socialgraph_db.socialgraph.followers REFERENCES user_db.user.UserID", true)
	assertConstraint(t, a, "FOREIGN_KEY socialgraph_db.socialgraph.followees REFERENCES user_db.user.UserID", true)
}

// RI-3: HomeTimelineService.ReadHomeTimeline reads the post ids from its cache and then reads the posts
// with PostStorageService.ReadPosts, without coordination with the writes of the posts
// (UserTimelineService.ReadUserTimeline does the same with both its cache and its database)
func TestSocialNetworkForeignKeyCoordination(t *testing.T) {
	a := runner.Get(t, "dsb_socialnetwork")
	assertWarnings(t, a, "foreign-key-coordination", 3,
		"entry request: Wrk2APIService.ReadHomeTimeline()",
		"READ (FOREIGN KEY): HomeTimelineService.ReadHomeTimeline() ... hometimeline_cache.*.Get()",
		"entry request: Wrk2APIService.ReadUserTimeline()",
		"READ (FOREIGN KEY): UserTimelineService.ReadUserTimeline() ... usertimeline_db.usertimeline.FindOne()",
		"READ (ORIGIN): PostStorageService.ReadPosts() ... post_db.post.FindMany()",
	)
}

// dsb_socialnetwork never deletes a referenced object and has no unique fields
func TestSocialNetworkNoDeleteOrUniquenessWarnings(t *testing.T) {
	a := runner.Get(t, "dsb_socialnetwork")
	for _, detectorType := range []string{"foreign-key-cascade", "foreign-key-concurrency", "primary-key-coordination", "uniqueness-concurrency"} {
		assertWarnings(t, a, detectorType, 0)
	}
}
