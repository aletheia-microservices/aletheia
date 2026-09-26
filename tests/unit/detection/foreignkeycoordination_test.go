package detection_test

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/detection/constraints/keycoordination"
	"github.com/aletheia-microservices/aletheia/internal/app"
)

// newPostnotificationApp builds the postnotification schema, where
// notifications_queue.notification.PostID references posts_db.post.PostID
func newPostnotificationApp(mandatoryReqs ...int) *app.App {
	a := newApp(map[string]string{"posts_db": "NoSQLDatabase", "notifications_queue": "Queue"})
	addForeignKey(field(a, "notifications_queue.notification.PostID"), field(a, "posts_db.post.PostID"), mandatoryReqs...)
	return a
}

// readNotificationAndPost returns the notifier's read of a notification from the queue, followed
// by the read of the post it points to
func readNotificationAndPost() (*abstractgraph.AbstractNode, []*abstractgraph.AbstractEdge) {
	pop := databaseCall("pop_notification", "Pop", serviceNode("NotifyService", "Run"), databaseNode("notifications_queue", "notification"), common.OP_READ,
		object(primary("notifications_queue.notification.PostID", "pop_notification", common.OP_READ)))
	find := databaseCall("find_post", "FindOne", serviceNode("StorageService", "ReadPost"), databaseNode("posts_db", "post"), common.OP_READ,
		object(
			primary("posts_db.post.PostID", "find_post", common.OP_READ),
			secondary("notifications_queue.notification.PostID", "pop_notification", common.OP_READ),
		))
	return serviceNode("UploadService", "UploadPost"), []*abstractgraph.AbstractEdge{pop, find}
}

func TestForeignKeyCoordinationReportsReadsOfMandatoryReference(t *testing.T) {
	// the foreign key is created in request 0 and read in request 1
	a := newPostnotificationApp(0)
	d := keycoordination.NewDetector(keycoordination.DETECTION_TYPE_FOREIGN_KEY)
	entry, reads := readNotificationAndPost()

	runRequest(d, a, 1, entry, reads...)

	got := results(d, a)
	wantWarnings(t, got, 1)
	wantLines(t, got,
		"entry request: UploadService.UploadPost()",
		"\tFOREIGN KEY READS #1:",
		"\t\tREAD (FOREIGN KEY): NotifyService.Run() ... notifications_queue.notification.Pop()",
		"\t\t\t- field: notifications_queue.notification.PostID",
		"\t\t\t- constraint: FOREIGN_KEY notifications_queue.notification.PostID REFERENCES posts_db.post.PostID [MANDATORY]",
		"\t\tREAD (ORIGIN): StorageService.ReadPost() ... posts_db.post.FindOne()",
		"\t\t\t- field: posts_db.post.PostID",
	)
}

func TestForeignKeyCoordinationIgnoresReferenceCreatedInSameRequest(t *testing.T) {
	a := newPostnotificationApp(0)
	d := keycoordination.NewDetector(keycoordination.DETECTION_TYPE_FOREIGN_KEY)
	entry, reads := readNotificationAndPost()

	runRequest(d, a, 0, entry, reads...)

	wantWarnings(t, results(d, a), 0)
}

func TestForeignKeyCoordinationIgnoresNonMandatoryReference(t *testing.T) {
	a := newPostnotificationApp()
	d := keycoordination.NewDetector(keycoordination.DETECTION_TYPE_FOREIGN_KEY)
	entry, reads := readNotificationAndPost()

	runRequest(d, a, 1, entry, reads...)

	wantWarnings(t, results(d, a), 0)
}

func TestForeignKeyCoordinationTypeString(t *testing.T) {
	if got := keycoordination.NewDetector(keycoordination.DETECTION_TYPE_FOREIGN_KEY).GetTypeString(); got != "foreign-key-coordination" {
		t.Errorf("GetTypeString() = %q, want foreign-key-coordination", got)
	}
}
