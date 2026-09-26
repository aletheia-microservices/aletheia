package utils_test

import (
	"testing"

	"analyzer/pkg/utils"
)

func TestLessT(t *testing.T) {
	tests := []struct {
		t1, t2 string
		want   bool
	}{
		{"t3", "t4", true},
		{"t4", "t3", false},
		{"t4", "t4", false},
		// numeric, not lexicographic
		{"t9", "t10", true},
		{"t10", "t9", false},
		// nested timestamps from combined graphs: <caller t>.<callee t>
		{"t4.t7", "t4.t8", true},
		{"t4.t8", "t4.t8.t9", true},
		{"t4.t8.t9", "t4.t8", false},
		{"t4.t14", "t5", true},
		{"t5", "t4.t14", false},
		// the callee timestamp of a combined graph does not matter when callers differ
		{"t4.t99", "t106.t15", true},
		// values without 't' prefix (e.g., empty) are always considered smaller
		{"", "t0", true},
		{"t0", "", false},
	}
	for _, tt := range tests {
		if got := utils.LessT(tt.t1, tt.t2); got != tt.want {
			t.Errorf("utils.LessT(%q, %q) = %v, want %v", tt.t1, tt.t2, got, tt.want)
		}
	}
}

func TestGreaterAndEqualT(t *testing.T) {
	if !utils.GreaterT("t10", "t9") {
		t.Errorf("utils.GreaterT(t10, t9) = false, want true")
	}
	if utils.GreaterT("t4.t8", "t4.t8.t9") {
		t.Errorf("utils.GreaterT(t4.t8, t4.t8.t9) = true, want false")
	}
	if !utils.EqualT("t4.t8", "t4.t8") || utils.EqualT("t4", "t4.t8") {
		t.Errorf("utils.EqualT must only hold for identical timestamps")
	}
}

func TestIsUpperPath(t *testing.T) {
	tests := []struct {
		upper, lower string
		want         bool
		diff         string
	}{
		{"notification", "notification.PostID", true, ".PostID"},
		{"_obj", "_obj.Creator.Username", true, ".Creator.Username"},
		{"_obj", "_obj[*]", true, "[*]"},
		{"notification", "notification", false, ""},
		{"notification.PostID", "notification", false, ""},
		{"posts", "notification.PostID", false, ""},
	}
	for _, tt := range tests {
		ok, diff := utils.IsUpperPath(tt.upper, tt.lower)
		if ok != tt.want || diff != tt.diff {
			t.Errorf("utils.IsUpperPath(%q, %q) = (%v, %q), want (%v, %q)", tt.upper, tt.lower, ok, diff, tt.want, tt.diff)
		}
	}
}

func TestIsUpperOrEqualPath(t *testing.T) {
	if ok, diff := utils.IsUpperOrEqualPath("_obj", "_obj"); !ok || diff != "" {
		t.Errorf("utils.IsUpperOrEqualPath(_obj, _obj) = (%v, %q), want (true, \"\")", ok, diff)
	}
	if ok, diff := utils.IsUpperOrEqualPath("_obj", "_obj.ID"); !ok || diff != ".ID" {
		t.Errorf("utils.IsUpperOrEqualPath(_obj, _obj.ID) = (%v, %q), want (true, \".ID\")", ok, diff)
	}
	if ok, _ := utils.IsUpperOrEqualPath("_obj.ID", "_obj"); ok {
		t.Errorf("utils.IsUpperOrEqualPath(_obj.ID, _obj) = true, want false")
	}
}

func TestExtractUpperPath(t *testing.T) {
	upper, sub, ok := utils.ExtractUpperPath("_obj.Creator.Username")
	if !ok || upper != "_obj.Creator" || sub != ".Username" {
		t.Errorf("utils.ExtractUpperPath = (%q, %q, %v), want (_obj.Creator, .Username, true)", upper, sub, ok)
	}
	if _, _, ok := utils.ExtractUpperPath("_obj"); ok {
		t.Errorf("utils.ExtractUpperPath(_obj) must fail for root path")
	}
}

func TestGetShortFunctionPath(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{
			"(*github.com/blueprint-uservices/blueprint/examples/postnotification/workflow/postnotification.StorageServiceImpl).StorePost",
			"postnotification.StorageService.StorePost",
		},
		{
			"github.com/blueprint-uservices/blueprint/examples/postnotification/workflow/postnotification/common.Int64ToString",
			"postnotification/common.Int64ToString",
		},
	}
	for _, tt := range tests {
		if got := utils.GetShortFunctionPath(tt.in); got != tt.want {
			t.Errorf("utils.GetShortFunctionPath(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestExtractFromShortFunctionPath(t *testing.T) {
	path := "postnotification.StorageService.StorePost"
	if got := utils.ExtractServiceNameFromShortFunctionPath(path); got != "StorageService" {
		t.Errorf("service = %q, want StorageService", got)
	}
	if got := utils.ExtractMethodNameFromShortFunctionPath(path); got != "StorePost" {
		t.Errorf("method = %q, want StorePost", got)
	}
	// plain functions do not belong to any service
	if got := utils.ExtractServiceNameFromShortFunctionPath("postnotification/common.Int64ToString"); got != "" {
		t.Errorf("service for plain function = %q, want empty", got)
	}
}

func TestExtractFromFieldPath(t *testing.T) {
	if got := utils.ExtractDatabaseNameFromFieldPath("posts_db.post.PostID"); got != "posts_db" {
		t.Errorf("database = %q, want posts_db", got)
	}
	if got := utils.ExtractSchemaNameFromFieldPath("posts_db.post.PostID"); got != "post" {
		t.Errorf("schema = %q, want post", got)
	}
}

func TestInsertAfterFieldName(t *testing.T) {
	tests := []struct {
		path, prefix, want string
	}{
		// full path: <database>.<table>.<fieldname>[.<any sub path> or [<any sub path>]
		{"posts_db.post.Creator", ".Username", "posts_db.post.Creator.Username"},
		{"posts_db.post.Creator.Username", "[*]", "posts_db.post.Creator[*].Username"},
		{"movie_db.movie.Casts[*].CastInfoID", ".Info", "movie_db.movie.Casts.Info[*].CastInfoID"},
	}
	for _, tt := range tests {
		if got := utils.InsertAfterFieldName(tt.path, tt.prefix); got != tt.want {
			t.Errorf("utils.InsertAfterFieldName(%q, %q) = %q, want %q", tt.path, tt.prefix, got, tt.want)
		}
	}
}
