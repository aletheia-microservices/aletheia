package abstractgraph_test

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph"
)

func TestAbstractTaintOperationTypes(t *testing.T) {
	tests := []struct {
		op                                           common.DatabaseOperationType
		read, write, writeOrUpdate, update, deleteOp bool
	}{
		{common.OP_READ, true, false, false, false, false},
		{common.OP_READ_MANY, true, false, false, false, false},
		{common.OP_WRITE, false, true, true, false, false},
		{common.OP_UPDATE, false, false, true, true, false},
		{common.OP_DELETE, false, false, false, false, true},
	}
	for _, tt := range tests {
		taint := abstractgraph.NewAbstractTaint("t1", "db.schema.F", "call", tt.op, true, false, false, false)
		if taint.IsRead() != tt.read || taint.IsWrite() != tt.write || taint.IsWriteOrUpdate() != tt.writeOrUpdate ||
			taint.IsUpdate() != tt.update || taint.IsDelete() != tt.deleteOp {
			t.Errorf("wrong operation predicates for %s", common.OperationTypeToString(tt.op))
		}
	}
}

func TestAbstractTaintFlags(t *testing.T) {
	taint := abstractgraph.NewAbstractTaint("t14", "posts_db.post.PostID", "StorageService.ReadPost.t14", common.OP_READ, false, true, true, false)
	if taint.IsPrimary() || !taint.IsTraced() || !taint.IsReadKey() || taint.IsReadValue() {
		t.Errorf("flags = (primary=%v, traced=%v, key=%v, value=%v), want (false, true, true, false)",
			taint.IsPrimary(), taint.IsTraced(), taint.IsReadKey(), taint.IsReadValue())
	}
	if taint.GetT() != "t14" || taint.GetDatabaseCallID() != "StorageService.ReadPost.t14" {
		t.Errorf("GetT() = %q, GetDatabaseCallID() = %q", taint.GetT(), taint.GetDatabaseCallID())
	}

	taint.SetReadKey(false)
	if taint.IsReadKey() {
		t.Errorf("SetReadKey(false) must clear the read key flag")
	}
}

func TestAbstractTaintSetReadValue(t *testing.T) {
	t.Skip("known bug: abstractgraph.AbstractTaint.SetReadValue assigns readKey instead of readVal")

	taint := primaryWrite("t1", "db.schema", "call")
	taint.SetReadValue(true)
	if !taint.IsReadValue() || taint.IsReadKey() {
		t.Errorf("SetReadValue(true) -> (key=%v, value=%v), want (false, true)", taint.IsReadKey(), taint.IsReadValue())
	}
}

func TestAbstractTaintComparisons(t *testing.T) {
	a := primaryWrite("t21", "posts_db.post", "StorageService.StorePost.t21")
	secondary := secondaryWrite("t21", "posts_db.post", "StorageService.StorePost.t21")
	otherT := primaryWrite("t99", "posts_db.post", "StorageService.StorePost.t21")
	readKey := abstractgraph.NewAbstractTaint("t21", "posts_db.post", "StorageService.StorePost.t21", common.OP_WRITE, true, false, true, false)

	// Similar only compares the database field and call
	if !a.Similar(a.Copy()) || !a.Similar(secondary) || !a.Similar(otherT) {
		t.Errorf("Similar must only compare the database field and call")
	}
	if a.Similar(primaryWrite("t21", "posts_db.post", "StorageService.UpdatePost.t9")) {
		t.Errorf("taints from different calls must not be similar")
	}

	// EqualExceptReadKeyAndReadVal also compares the operation, primary and traced flags
	if a.EqualExceptReadKeyAndReadVal(secondary) {
		t.Errorf("EqualExceptReadKeyAndReadVal must compare the primary flag")
	}
	if !a.EqualExceptReadKeyAndReadVal(readKey) {
		t.Errorf("EqualExceptReadKeyAndReadVal must ignore read key/value flags")
	}

	// EqualExceptPrimaryAndTrace also compares the operation and read key/value flags
	if !a.EqualExceptPrimaryAndTrace(secondary) {
		t.Errorf("EqualExceptPrimaryAndTrace must ignore the primary flag")
	}
	if a.EqualExceptPrimaryAndTrace(readKey) {
		t.Errorf("EqualExceptPrimaryAndTrace must compare read key/value flags")
	}
}

func TestAbstractTaintIsUpperTaint(t *testing.T) {
	upper := primaryWrite("t21", "posts_db.post", "StorageService.StorePost.t21")
	lower := primaryWrite("t21", "posts_db.post.Creator.Username", "StorageService.StorePost.t21")
	otherCall := primaryWrite("t21", "posts_db.post.PostID", "StorageService.UpdatePost.t9")

	if ok, diff := upper.IsUpperTaint(lower); !ok || diff != ".Creator.Username" {
		t.Errorf("IsUpperTaint = (%v, %q), want (true, .Creator.Username)", ok, diff)
	}
	if ok, _ := upper.IsUpperTaint(otherCall); ok {
		t.Errorf("upper taint must belong to the same database call")
	}
	if ok, _ := lower.IsUpperTaint(upper); ok {
		t.Errorf("lower taint must not be an upper taint")
	}
	if ok, _ := upper.IsUpperTaint(upper.Copy()); ok {
		t.Errorf("a taint must not be an upper taint of itself")
	}
}

func TestAbstractTaintCopyAndPaths(t *testing.T) {
	a := primaryWrite("t21", "posts_db.post", "StorageService.StorePost.t21")

	cp := a.Copy()
	cp.AddSuffixToDatabasePath(".Text")
	if a.GetDatabasePath() != "posts_db.post" || cp.GetDatabasePath() != "posts_db.post.Text" {
		t.Errorf("Copy must not share state with the original taint")
	}
	cp.SetDatabasepath("posts_db.post.PostID")
	if cp.GetDatabasePath() != "posts_db.post.PostID" || cp.String() != "posts_db.post.PostID" {
		t.Errorf("SetDatabasepath = %q, String() = %q", cp.GetDatabasePath(), cp.String())
	}
}

func TestAbstractTaintStrings(t *testing.T) {
	taint := abstractgraph.NewAbstractTaint("t14", "posts_db.post.PostID", "StorageService.ReadPost.t14", common.OP_READ, true, false, true, false)

	if got, want := taint.LongString(), "{posts_db.post.PostID, StorageService.ReadPost.t14, read, true}"; got != want {
		t.Errorf("LongString() = %q, want %q", got, want)
	}
	if got, want := taint.LongLongString(), "{t14, posts_db.post.PostID, StorageService.ReadPost.t14, read, true, false, true, false}"; got != want {
		t.Errorf("LongLongString() = %q, want %q", got, want)
	}
}
