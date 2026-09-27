package detection

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/tests/runner"
)

// the queue message stores the post id and request id written to posts_db,
// so both fields must be inferred as (mandatory) foreign keys
func TestPostNotificationForeignKeysFromDataFlow(t *testing.T) {
	a := runner.Get(t, "postnotification")
	assertConstraint(t, a, "FOREIGN_KEY notifications_queue.notification.PostID REFERENCES posts_db.post.PostID [MANDATORY]", true)
	assertConstraint(t, a, "FOREIGN_KEY notifications_queue.notification.ReqID REFERENCES posts_db.post.ReqID [MANDATORY]", true)
	// references must follow the data flow direction
	assertConstraint(t, a, "FOREIGN_KEY posts_db.post.PostID REFERENCES notifications_queue.notification.PostID [MANDATORY]", false)
}

func TestPostNotificationIgnoreForeignKeysConfig(t *testing.T) {
	a := runner.GetWithConfig(t, "postnotification", "config/postnotification.yaml")
	assertConstraint(t, a, "FOREIGN_KEY notifications_queue.notification.ReqID REFERENCES posts_db.post.ReqID [MANDATORY]", false)
	assertConstraint(t, a, "FOREIGN_KEY notifications_queue.notification.PostID REFERENCES posts_db.post.PostID [MANDATORY]", true)
}

// RI-3: NotifyService reads the post referenced by the queued message without coordination
// with UploadService, which writes the post before pushing the message
func TestPostNotificationForeignKeyCoordination(t *testing.T) {
	a := runner.Get(t, "postnotification")
	assertWarnings(t, a, "foreign-key-coordination", 1,
		"entry request: UploadService.UploadPost()",
		"READ (FOREIGN KEY): NotifyService.Run() ... notifications_queue.notification.Pop()",
		"- constraint: FOREIGN_KEY notifications_queue.notification.PostID REFERENCES posts_db.post.PostID [MANDATORY]",
		"READ (ORIGIN): StorageService.ReadPost() ... posts_db.post.FindOne()",
	)
	for _, detectorType := range []string{"foreign-key-cascade", "foreign-key-concurrency", "primary-key-coordination", "uniqueness-concurrency"} {
		assertWarnings(t, a, detectorType, 0)
	}
}
