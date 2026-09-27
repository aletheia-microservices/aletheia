package ssa

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/tests/runner"
)

// CatalogService.store inserts the product it receives as is: the parameter is tainted as the
// written collection
func TestEshopWriteParameterObject(t *testing.T) {
	a := runner.Get(t, "eshopmicroservices")
	graph := getSSAGraph(t, a, "catalog.CatalogService.store")

	call := getDatabaseCall(t, graph, "catalog_db.product", "InsertOne")
	if call.GetOpType() != common.OP_WRITE {
		t.Fatalf("InsertOne op type = %s, want write", common.OperationTypeToString(call.GetOpType()))
	}
	assertDatabaseTaint(t, getParam(t, graph, "product"), dbTaint{objpath: "_obj", dbpath: "catalog_db.product", op: common.OP_WRITE, t: call.GetT()})
}

// BasketService.getBasket reads the basket of a user with the filter {UserName: username}
func TestEshopReadBasketByUserName(t *testing.T) {
	a := runner.Get(t, "eshopmicroservices")
	graph := getSSAGraph(t, a, "basket.BasketService.getBasket")

	call := getDatabaseCall(t, graph, "basket_db.basket", "FindOne")
	assertDatabaseTaint(t, getParam(t, graph, "username"), dbTaint{objpath: "_obj", dbpath: "basket_db.basket.UserName", op: common.OP_READ, readKey: true, t: call.GetT()})
	assertDatabaseTaint(t, call.GetArguments()[2], dbTaint{objpath: "_obj", dbpath: "basket_db.basket", op: common.OP_READ, readVal: true})
}

// OrderService.Init pops the checkout messages pushed by BasketService.CheckoutBasket: every field
// of the popped message is a read value of the queue
func TestEshopQueuePopTaintsMessageFields(t *testing.T) {
	a := runner.Get(t, "eshopmicroservices")
	graph := getSSAGraph(t, a, "order.OrderService.Init")

	pop := getDatabaseCall(t, graph, "order_queue.notification", "Pop")
	if pop.GetOpType() != common.OP_READ {
		t.Fatalf("Pop op type = %s, want read", common.OperationTypeToString(pop.GetOpType()))
	}
	msg := pop.GetArguments()[0]
	for _, field := range []string{"", ".CustomerId", ".UserName", ".CardNumber", ".AddressLine"} {
		assertDatabaseTaint(t, msg, dbTaint{objpath: "_obj" + field, dbpath: "order_queue.notification" + field, op: common.OP_READ, readVal: true, t: pop.GetT()})
	}
}
