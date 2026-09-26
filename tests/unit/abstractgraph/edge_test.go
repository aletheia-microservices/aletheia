package abstractgraph_test

import (
	"testing"

	"analyzer/pkg/analysis/common"
	"analyzer/pkg/analysis/system-level/abstractgraph"
)

func newServiceNode(service string, method string) *abstractgraph.AbstractNode {
	return abstractgraph.NewAbstractNode(service+"."+method, abstractgraph.NODE_SERVICE, service, method, "", "")
}

func newDatabaseNode(database string, schema string) *abstractgraph.AbstractNode {
	return abstractgraph.NewAbstractNode(database+"."+schema, abstractgraph.NODE_DATABASE, "", "", database, schema)
}

func TestAbstractEdgeDatabaseCall(t *testing.T) {
	upload := newServiceNode("UploadService", "UploadPost")
	queue := newDatabaseNode("notifications_queue", "notification")
	edge := abstractgraph.NewAbstractEdge("t13", "UploadService.UploadPost.t13", "Push", upload, queue, common.OP_WRITE, abstractgraph.EDGE_DATABASE_CALL)

	if edge.GetEdgeType() != abstractgraph.EDGE_DATABASE_CALL || edge.GetOpType() != common.OP_WRITE {
		t.Errorf("edge type = %d, op = %s", edge.GetEdgeType(), common.OperationTypeToString(edge.GetOpType()))
	}
	if edge.GetFromNode() != upload || edge.GetToNode() != queue || edge.GetMethod() != "Push" || edge.GetT() != "t13" {
		t.Errorf("edge = %s (t=%s)", edge.String(), edge.GetT())
	}
	if got, want := edge.String(), "UploadService.UploadPost() ... notifications_queue.notification.Push()"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestAbstractEdgeArgumentsAndReturns(t *testing.T) {
	edge := abstractgraph.NewAbstractEdge("t3", "UploadService.UploadPost.t3", "StorePost",
		newServiceNode("UploadService", "UploadPost"), newServiceNode("StorageService", "StorePost"), common.OP_UNDEFINED, abstractgraph.EDGE_SERVICE_RPC)
	reqID := newTracedObject("t0", nil, nil)
	text := newTracedObject("text", nil, nil)
	ret := newTracedObject("t4", nil, nil)
	edge.AddArgument(reqID)
	edge.AddArgument(text)
	edge.AddReturn(ret)

	if len(edge.GetArguments()) != 2 || edge.GetArgumentAt(1) != text {
		t.Errorf("arguments = %v", edge.GetArguments())
	}
	if edge.GetArgumentByNameIfExists("text") != text || edge.GetArgumentByNameIfExists("missing") != nil {
		t.Errorf("GetArgumentByNameIfExists must find arguments by name")
	}
	if len(edge.GetReturns()) != 1 || edge.GetReturnAt(0) != ret {
		t.Errorf("returns = %v", edge.GetReturns())
	}

	edge.SetArguments([]*abstractgraph.AbstractObject{text})
	if len(edge.GetArguments()) != 1 || edge.GetArgumentAt(0) != text {
		t.Errorf("SetArguments must replace the arguments")
	}
}

func TestAbstractEdgeIDNumber(t *testing.T) {
	tests := []struct {
		id   string
		want int
	}{
		{"UploadService.UploadPost.t13", 13},
		{"digota.OrderService.storageGetOne.storageGetOne.t14", 14},
		// entrypoint edges only have the function path
		{"postnotification.UploadService.UploadPost", -1},
		// ids of constant values
		{"ProductService.New.nil:*github.com/blueprint-uservices/blueprint/examples/digota/workflow/digota.PackageDimensions", -1},
		{`ProductService.New."image":string`, -1},
	}
	for _, tt := range tests {
		edge := abstractgraph.NewAbstractEdge("", tt.id, "m", nil, nil, common.OP_UNDEFINED, abstractgraph.EDGE_SERVICE_RPC)
		if got := edge.GetIDNumber(); got != tt.want {
			t.Errorf("GetIDNumber(%q) = %d, want %d", tt.id, got, tt.want)
		}
	}
}
