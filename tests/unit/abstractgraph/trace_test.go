package abstractgraph_test

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph"
)

func TestAbstractTraceArgument(t *testing.T) {
	tests := []struct {
		svpath, name, path string
	}{
		{"MovieIdService.RegisterMovieId.t4", "t4", "_obj"},
		{"MovieIdService.RegisterMovieId.t4.MovieId", "t4", "_obj.MovieId"},
		{"MovieInfoService.ReadMovieInfo.t4.Casts[*].CastInfoID", "t4", "_obj.Casts[*].CastInfoID"},
		{"CastInfoService.ReadCastInfos.t17[*]", "t17", "_obj[*]"},
	}
	for _, tt := range tests {
		trace := abstractgraph.NewAbstractTrace("t1", tt.svpath, "call")
		if got := trace.GetArgumentName(); got != tt.name {
			t.Errorf("GetArgumentName(%q) = %q, want %q", tt.svpath, got, tt.name)
		}
		if got := trace.GetArgumentPath(); got != tt.path {
			t.Errorf("GetArgumentPath(%q) = %q, want %q", tt.svpath, got, tt.path)
		}
	}
}

func TestAbstractTraceArgumentPathOnArrayWithSubPath(t *testing.T) {
	t.Skip("known bug: abstractgraph.AbstractTrace.GetArgumentPath drops '[*]' when the array is followed by a sub path " +
		"(e.g., 'CastInfoService.ReadCastInfos.t17[*].CastInfoID' yields '_obj.CastInfoID')")

	trace := abstractgraph.NewAbstractTrace("t1", "CastInfoService.ReadCastInfos.t17[*].CastInfoID", "call")
	if got := trace.GetArgumentPath(); got != "_obj[*].CastInfoID" {
		t.Errorf("GetArgumentPath = %q, want _obj[*].CastInfoID", got)
	}
}

func TestAbstractTraceComparisons(t *testing.T) {
	upper := abstractgraph.NewAbstractTrace("t3", "StorageService.StorePost.t0", "UploadService.UploadPost.t3")
	lower := abstractgraph.NewAbstractTrace("t3", "StorageService.StorePost.t0.PostID", "UploadService.UploadPost.t3")
	otherCall := abstractgraph.NewAbstractTrace("t3", "StorageService.StorePost.t0.PostID", "UploadService.UploadPost.t9")

	if !upper.Equals(abstractgraph.NewAbstractTrace("t9", "StorageService.StorePost.t0", "UploadService.UploadPost.t3")) {
		t.Errorf("traces with same path and call must be equal regardless of t")
	}
	if upper.Equals(abstractgraph.NewAbstractTrace("t3", "StorageService.StorePost.t0", "UploadService.UploadPost.t9")) {
		t.Errorf("traces from different calls must not be equal")
	}

	if ok, sub := upper.IsUpperPath(lower); !ok || sub != ".PostID" {
		t.Errorf("IsUpperPath = (%v, %q), want (true, .PostID)", ok, sub)
	}
	if ok, _ := upper.IsUpperPath(otherCall); ok {
		t.Errorf("upper trace must belong to the same call")
	}
	if ok, _ := lower.IsUpperPath(upper); ok {
		t.Errorf("lower trace must not be an upper trace")
	}
	if ok, _ := upper.IsUpperPath(upper); ok {
		t.Errorf("a trace must not be an upper trace of itself")
	}
}

func TestAbstractTraceStrings(t *testing.T) {
	trace := abstractgraph.NewAbstractTrace("t3", "StorageService.StorePost.text", "UploadService.UploadPost.t3")
	if trace.String() != "StorageService.StorePost.text" || trace.GetServicePath() != trace.String() {
		t.Errorf("String() = %q", trace.String())
	}
	if got, want := trace.LongString(), "{StorageService.StorePost.text, UploadService.UploadPost.t3, rpc}"; got != want {
		t.Errorf("LongString() = %q, want %q", got, want)
	}
	if trace.GetT() != "t3" || trace.GetServiceCallID() != "UploadService.UploadPost.t3" {
		t.Errorf("GetT() = %q, GetServiceCallID() = %q", trace.GetT(), trace.GetServiceCallID())
	}
}
