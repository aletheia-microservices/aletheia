package ssa

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/tests/runner"
)

// CatalogueService.Get reads a sock with a SQL statement (SELECT ... WHERE sock.SockID = ?):
// the id is the read key and the destination sock (and its fields) the read value
func TestSockshopSQLReadKeyAndValue(t *testing.T) {
	a := runner.Get(t, "sockshop")
	graph := getSSAGraph(t, a, "catalogue.CatalogueService.Get")

	call := getDatabaseCall(t, graph, "catalogue_db.sock", "Get")
	if call.GetOpType() != common.OP_READ {
		t.Fatalf("Get op type = %s, want read", common.OperationTypeToString(call.GetOpType()))
	}
	assertDatabaseTaint(t, getParam(t, graph, "id"), dbTaint{objpath: "_obj", dbpath: "catalogue_db.sock.SockID", op: common.OP_READ, readKey: true, t: call.GetT()})

	sock := call.GetArguments()[1]
	assertDatabaseTaint(t, sock, dbTaint{objpath: "_obj", dbpath: "catalogue_db.sock", op: common.OP_READ, readVal: true})
	assertDatabaseTaint(t, sock, dbTaint{objpath: "_obj.ImageURL1", dbpath: "catalogue_db.sock.ImageURL1", op: common.OP_READ, readVal: true})
}

// CartService.getCart reads the cart with the filter {ID: id}: the id is the read key
// and the decoded cart the read value
func TestSockshopNoSQLReadKeyAndValue(t *testing.T) {
	a := runner.Get(t, "sockshop")
	graph := getSSAGraph(t, a, "carts.CartService.getCart")

	call := getDatabaseCall(t, graph, "cart_db.carts", "FindOne")
	assertDatabaseTaint(t, getParam(t, graph, "id"), dbTaint{objpath: "_obj", dbpath: "cart_db.carts.ID", op: common.OP_READ, readKey: true, t: call.GetT()})
	assertDatabaseTaint(t, call.GetArguments()[2], dbTaint{objpath: "_obj", dbpath: "cart_db.carts", op: common.OP_READ, readVal: true})
}

// ShippingService.PostShipping pushes the shipment to the queue and then inserts it in the database:
// the same shipment (and its fields) carries one write taint per call, each with the t of its call
func TestSockshopSameObjectWrittenToQueueAndDatabase(t *testing.T) {
	a := runner.Get(t, "sockshop")
	graph := getSSAGraph(t, a, "shipping.ShippingService.PostShipping")

	push := getDatabaseCall(t, graph, "ship_queue.notification", "Push")
	insert := getDatabaseCall(t, graph, "ship_db.shipments", "InsertOne")
	shipment := getParam(t, graph, "shipment")
	for _, field := range []string{"", ".ID", ".Name"} {
		assertDatabaseTaint(t, shipment, dbTaint{objpath: "_obj" + field, dbpath: "ship_queue.notification" + field, op: common.OP_WRITE, t: push.GetT()})
		assertDatabaseTaint(t, shipment, dbTaint{objpath: "_obj" + field, dbpath: "ship_db.shipments" + field, op: common.OP_WRITE, t: insert.GetT()})
	}
}
