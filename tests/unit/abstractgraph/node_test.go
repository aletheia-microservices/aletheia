package abstractgraph_test

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph"
)

func TestAbstractNodeService(t *testing.T) {
	node := abstractgraph.NewAbstractNode("StorageService.StorePost", abstractgraph.NODE_SERVICE, "StorageService", "StorePost", "", "")

	if node.GetNodeType() != abstractgraph.NODE_SERVICE || node.GetName() != "StorageService.StorePost" || node.String() != node.GetName() {
		t.Errorf("node = (%d, %q)", node.GetNodeType(), node.GetName())
	}
	if node.GetServiceName() != "StorageService" || node.GetMethod() != "StorePost" || node.GetServiceWithMethod() != "StorageService.StorePost" {
		t.Errorf("service = %q, method = %q", node.GetServiceName(), node.GetMethod())
	}

	if node.IsParsed() {
		t.Errorf("new node must not be parsed")
	}
	node.SetParsed()
	if !node.IsParsed() {
		t.Errorf("SetParsed must mark the node as parsed")
	}
}

func TestAbstractNodeParamsAndReturns(t *testing.T) {
	node := abstractgraph.NewAbstractNode("StorageService.StorePost", abstractgraph.NODE_SERVICE, "StorageService", "StorePost", "", "")
	reqID := newTracedObject("reqID", nil, nil)
	text := newTracedObject("text", nil, nil)
	ret := newTracedObject("int64", nil, nil)
	node.AddParam(reqID)
	node.AddParam(text)
	node.AddReturn(ret)

	if len(node.GetParams()) != 2 || node.GetParamAt(1) != text {
		t.Errorf("params = %v", node.GetParams())
	}
	if node.GetParameterByNameIfExists("text") != text || node.GetParameterByNameIfExists("missing") != nil {
		t.Errorf("GetParameterByNameIfExists must find parameters by name")
	}
	if len(node.GetReturns()) != 1 || node.GetReturnAt(0) != ret {
		t.Errorf("returns = %v", node.GetReturns())
	}
}

func TestAbstractNodeDatabase(t *testing.T) {
	node := abstractgraph.NewAbstractNode("posts_db.post", abstractgraph.NODE_DATABASE, "", "", "posts_db", "post")

	if node.GetNodeType() != abstractgraph.NODE_DATABASE || node.GetDatabaseName() != "posts_db" || node.GetSchemaName() != "post" {
		t.Errorf("database node = (%d, %q, %q)", node.GetNodeType(), node.GetDatabaseName(), node.GetSchemaName())
	}
}
