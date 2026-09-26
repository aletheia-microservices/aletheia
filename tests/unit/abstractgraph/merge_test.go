package abstractgraph_test

import (
	"slices"
	"testing"

	"analyzer/pkg/analysis/common"
	"analyzer/pkg/analysis/system-level/abstractgraph"
	abstractgraphtainter "analyzer/pkg/analysis/system-level/abstractgraph/tainter"
)

// MergeTaints is the core of inter-service taint propagation: merging a (secondary) taint into an
// object that already holds a primary taint on an upper path must map the primary taint (extended
// with the sub path) to the new taint, which is later used to infer foreign keys
func TestMergeTaintsMapsPrimaryUpperTaintToNewTaint(t *testing.T) {
	obj := newObject(map[string][]*abstractgraph.AbstractTaint{
		"_obj": {primaryWrite("t21", "posts_db.post", "StorageService.StorePost.t21")},
	})
	incoming := map[string][]*abstractgraph.AbstractTaint{
		"_obj.PostID": {primaryWrite("t13", "notifications_queue.notification.PostID", "UploadService.UploadPost.t13")},
	}

	tm := abstractgraphtainter.MergeTaints(obj, incoming, nil, abstractgraphtainter.MERGE_MODE_TAINT, "t3", false)

	added := obj.GetTaintsForObjectPath("_obj.PostID")
	if len(added) != 1 {
		t.Fatalf("_obj.PostID has %d taints, want 1", len(added))
	}
	if added[0].IsPrimary() || added[0].GetT() != "t3" || added[0].GetDatabasePath() != "notifications_queue.notification.PostID" {
		t.Errorf("merged taint = %s, want secondary taint at t3", added[0].LongLongString())
	}

	want := []string{"posts_db.post.PostID -> notifications_queue.notification.PostID"}
	if got := mappingPairs(tm); !slices.Equal(got, want) {
		t.Fatalf("mapping = %v, want %v", got, want)
	}
	if !tm.GetMappingKeys()[0].IsPrimary() {
		t.Errorf("mapping key must be the primary taint")
	}

	// merging the same taints again must be idempotent
	tm = abstractgraphtainter.MergeTaints(obj, incoming, nil, abstractgraphtainter.MERGE_MODE_TAINT, "t3", false)
	if len(obj.GetTaintsForObjectPath("_obj.PostID")) != 1 || len(tm.GetMappingKeys()) != 0 {
		t.Errorf("second merge must not add taints nor mappings")
	}
}

func TestMergeTaintsMapsPrimaryTaintOnSamePath(t *testing.T) {
	obj := newObject(map[string][]*abstractgraph.AbstractTaint{
		"_obj": {primaryWrite("t21", "posts_db.post.PostID", "StorageService.StorePost.t21")},
	})
	incoming := map[string][]*abstractgraph.AbstractTaint{
		"_obj": {primaryWrite("t13", "notifications_queue.notification.PostID", "UploadService.UploadPost.t13")},
	}

	tm := abstractgraphtainter.MergeTaints(obj, incoming, nil, abstractgraphtainter.MERGE_MODE_TAINT, "t3", false)

	want := []string{"posts_db.post.PostID -> notifications_queue.notification.PostID"}
	if got := mappingPairs(tm); !slices.Equal(got, want) {
		t.Errorf("mapping = %v, want %v", got, want)
	}
}

func TestMergeTaintsIgnoresSecondaryTaintsInTaintMode(t *testing.T) {
	obj := newObject(map[string][]*abstractgraph.AbstractTaint{
		"_obj": {secondaryWrite("t5", "posts_db.post", "StorageService.StorePost.t21")},
	})
	incoming := map[string][]*abstractgraph.AbstractTaint{
		"_obj.PostID": {primaryWrite("t13", "notifications_queue.notification.PostID", "UploadService.UploadPost.t13")},
	}

	tm := abstractgraphtainter.MergeTaints(obj, incoming, nil, abstractgraphtainter.MERGE_MODE_TAINT, "t3", false)

	if got := mappingPairs(tm); len(got) != 0 {
		t.Errorf("mapping = %v, want none (only primary taints are mapped)", got)
	}
	if len(obj.GetTaintsForObjectPath("_obj.PostID")) != 1 {
		t.Errorf("the new taint must still be added")
	}
}

func TestMergeTaintsMapsGrandparentPath(t *testing.T) {
	t.Skip("known bug: MergeTaints climbs upper paths with only the last path segment as sub path, so a taint on " +
		"_obj.Creator.Username merged into an object with a primary taint on _obj is mapped from posts_db.post.Creator " +
		"instead of posts_db.post.Creator.Username")

	obj := newObject(map[string][]*abstractgraph.AbstractTaint{
		"_obj":         {primaryWrite("t21", "posts_db.post", "c1")},
		"_obj.Creator": {primaryWrite("t21", "posts_db.post.Creator", "c1")},
	})
	incoming := map[string][]*abstractgraph.AbstractTaint{
		"_obj.Creator.Username": {primaryWrite("t13", "users_db.user.Username", "c2")},
	}

	tm := abstractgraphtainter.MergeTaints(obj, incoming, nil, abstractgraphtainter.MERGE_MODE_TAINT, "t3", false)

	want := []string{"posts_db.post.Creator.Username -> users_db.user.Username"}
	if got := mappingPairs(tm); !slices.Equal(got, want) {
		t.Errorf("mapping = %v, want %v", got, want)
	}
}

func TestMergeTaintsAddsTaintsToTracedLowerLocations(t *testing.T) {
	obj := newTracedObject("t0", nil, map[string][]*abstractgraph.AbstractTrace{
		"_obj.ID": {abstractgraph.NewAbstractTrace("t4", "StorageService.ReadPost.t0.ID", "rpc")},
	})
	incoming := map[string][]*abstractgraph.AbstractTaint{
		"_obj": {primaryWrite("t2", "posts_db.post", "c")},
	}

	abstractgraphtainter.MergeTaints(obj, incoming, nil, abstractgraphtainter.MERGE_MODE_TAINT, "t9", false)

	// the location only annotated by a trace also receives the taint, extended with its sub path
	got := obj.GetTaintsForObjectPath("_obj.ID")
	if len(got) != 1 || got[0].GetDatabasePath() != "posts_db.post.ID" || got[0].GetT() != "t9" {
		t.Errorf("_obj.ID taints = %v, want posts_db.post.ID at t9", got)
	}
	if len(obj.GetTaintsForObjectPath("_obj")) != 1 {
		t.Errorf("_obj must receive the merged taint")
	}
}

func TestMergeTaintsParseModeKeepsTaintsPrimary(t *testing.T) {
	obj := newObject(nil)
	incoming := map[string][]*abstractgraph.AbstractTaint{
		"_obj": {abstractgraph.NewAbstractTaint("t14", "posts_db.post.PostID", "c", common.OP_READ, false, false, true, false)},
	}

	tm := abstractgraphtainter.MergeTaints(obj, incoming, nil, abstractgraphtainter.MERGE_MODE_PARSE, "", false)

	got := obj.GetTaintsForObjectPath("_obj")
	if len(got) != 1 || !got[0].IsPrimary() || got[0].GetT() != "t14" || !got[0].IsReadKey() {
		t.Errorf("parsed taint = %v, want primary read key taint at t14", got)
	}
	if len(tm.GetMappingKeys()) != 0 {
		t.Errorf("parse mode must not create taint mappings")
	}
}

// in trace mode, existing secondary taints are also mapped (e.g., gateway services that forward
// the same object to two services), except when they have the same t as the new taint
func TestMergeTaintsTraceMode(t *testing.T) {
	obj := newObject(map[string][]*abstractgraph.AbstractTaint{
		"_obj": {secondaryWrite("t34", "movie_info_db.movie_info.Casts[*].CastInfoID", "c1")},
	})
	incoming := map[string][]*abstractgraph.AbstractTaint{
		"_obj": {primaryWrite("t55", "cast_info_db.cast.CastInfoID", "c2")},
	}

	tm := abstractgraphtainter.MergeTaints(obj, incoming, nil, abstractgraphtainter.MERGE_MODE_TRACE, "t55", false)

	want := []string{"movie_info_db.movie_info.Casts[*].CastInfoID -> cast_info_db.cast.CastInfoID"}
	if got := mappingPairs(tm); !slices.Equal(got, want) {
		t.Errorf("mapping = %v, want %v", got, want)
	}
	var added *abstractgraph.AbstractTaint
	for _, taint := range obj.GetTaintsForObjectPath("_obj") {
		if taint.GetDatabasePath() == "cast_info_db.cast.CastInfoID" {
			added = taint
		}
	}
	if added == nil || !added.IsTraced() || added.IsPrimary() || added.GetT() != "t55" {
		t.Errorf("merged taint = %v, want traced (non primary) taint at t55", added)
	}

	// secondary taints with the same t come from the same source and are not mapped
	obj = newObject(map[string][]*abstractgraph.AbstractTaint{
		"_obj": {secondaryWrite("t55", "movie_info_db.movie_info.Casts[*].CastInfoID", "c1")},
	})
	tm = abstractgraphtainter.MergeTaints(obj, incoming, nil, abstractgraphtainter.MERGE_MODE_TRACE, "t55", false)
	if got := mappingPairs(tm); len(got) != 0 {
		t.Errorf("mapping = %v, want none for taints with the same t", got)
	}
}

func TestMergeTaintsDebugModeDoesNotModify(t *testing.T) {
	obj := newTracedObject("t0",
		map[string][]*abstractgraph.AbstractTaint{"_obj": {primaryWrite("t21", "posts_db.post", "c1")}},
		map[string][]*abstractgraph.AbstractTrace{"_obj.ID": {abstractgraph.NewAbstractTrace("t4", "S.M.t0.ID", "rpc")}},
	)
	incoming := map[string][]*abstractgraph.AbstractTaint{
		"_obj.PostID": {primaryWrite("t13", "notifications_queue.notification.PostID", "c2")},
	}

	tm := abstractgraphtainter.MergeTaints(obj, incoming, nil, abstractgraphtainter.MERGE_MODE_DEBUG, "t3", false)

	if len(obj.GetAllTaintsFlatList()) != 1 || len(tm.GetMappingKeys()) != 0 {
		t.Errorf("debug mode must not add taints nor mappings: %v", obj.GetAllTaints())
	}
}

func TestMergeTraces(t *testing.T) {
	obj := newTracedObject("t0", nil, map[string][]*abstractgraph.AbstractTrace{
		"_obj": {abstractgraph.NewAbstractTrace("t3", "StorageService.StorePost.t0", "UploadService.UploadPost.t3")},
	})

	abstractgraphtainter.MergeTraces(obj, map[string][]*abstractgraph.AbstractTrace{
		"_obj": {
			// same path and call: skipped even with a different t
			abstractgraph.NewAbstractTrace("t9", "StorageService.StorePost.t0", "UploadService.UploadPost.t3"),
			abstractgraph.NewAbstractTrace("t5", "StorageService.StorePost.t0", "UploadService.UploadPost.t5"),
		},
		"_obj.PostID": {abstractgraph.NewAbstractTrace("t3", "StorageService.StorePost.t0.PostID", "UploadService.UploadPost.t3")},
	})

	if got := len(obj.GetTracesForObjectPath("_obj")); got != 2 {
		t.Errorf("_obj has %d traces, want 2", got)
	}
	if got := len(obj.GetTracesForObjectPath("_obj.PostID")); got != 1 {
		t.Errorf("_obj.PostID has %d traces, want 1", got)
	}
}
