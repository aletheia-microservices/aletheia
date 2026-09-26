package abstractcallgraph

import (
	"testing"

	"analyzer/pkg/analysis/common"
	"analyzer/tests/runner"
)

// digota: order, payment, product and sku services, each with its own database; orders reference
// skus, and skus reference products

func TestDigotaEntrypoints(t *testing.T) {
	g := runner.Get(t, "digota").AbsGraph
	assertList(t, "entrypoints", entrypoints(g), []string{
		"OrderService.Get",
		"OrderService.List",
		"OrderService.New",
		"OrderService.Pay",
		"OrderService.Return",
		"PaymentService.Get",
		"PaymentService.List",
		"PaymentService.NewCharge",
		"PaymentService.RefundCharge",
		"ProductService.Delete",
		"ProductService.Get",
		"ProductService.List",
		"ProductService.New",
		"ProductService.Update",
		"SkuService.Delete",
		"SkuService.Get",
		"SkuService.List",
		"SkuService.New",
		"SkuService.Update",
	})
}

func TestDigotaCalls(t *testing.T) {
	g := runner.Get(t, "digota").AbsGraph
	assertList(t, "rpc calls", rpcCalls(g), []string{
		"OrderService.New -> SkuService.Get",
		"OrderService.Pay -> PaymentService.NewCharge",
		"OrderService.Pay -> PaymentService.RefundCharge",
		"OrderService.Pay -> SkuService.Get",
		"OrderService.Return -> PaymentService.RefundCharge",
		"OrderService.Return -> SkuService.Get",
		"SkuService.New -> ProductService.Get",
		"SkuService.Update -> ProductService.Get",
	})
	assertList(t, "database calls", databaseCalls(g, nil), []string{
		"delete ProductService.Delete -> products_db.products.DeleteOne",
		"delete SkuService.Delete -> skus_db.skus.DeleteOne",
		"read OrderService.Get -> orders_db.orders.FindOne",
		"read OrderService.List -> orders_db.orders.FindMany",
		"read OrderService.Pay -> orders_db.orders.FindOne",
		"read OrderService.Return -> orders_db.orders.FindOne",
		"read PaymentService.Get -> payments_db.payments.FindOne",
		"read PaymentService.List -> payments_db.payments.FindMany",
		"read PaymentService.RefundCharge -> payments_db.payments.FindOne",
		"read ProductService.Get -> products_db.products.FindOne",
		"read ProductService.List -> products_db.products.FindMany",
		"read ProductService.Update -> products_db.products.FindOne",
		"read SkuService.Get -> skus_db.skus.FindOne",
		"read SkuService.List -> skus_db.skus.FindMany",
		"read SkuService.Update -> skus_db.skus.FindOne",
		"update OrderService.Pay -> orders_db.orders.ReplaceOne",
		"update OrderService.Return -> orders_db.orders.ReplaceOne",
		"update ProductService.Update -> products_db.products.ReplaceOne",
		"update SkuService.Update -> skus_db.skus.ReplaceOne",
		"write OrderService.New -> orders_db.orders.InsertOne",
		"write PaymentService.NewCharge -> payments_db.payments.InsertOne",
		"write ProductService.New -> products_db.products.InsertOne",
		"write SkuService.New -> skus_db.skus.InsertOne",
	})
}

func TestDigotaCounts(t *testing.T) {
	assertCounts(t, runner.Get(t, "digota").AbsGraph, map[string]int{
		"service nodes": 19, "database nodes": 4,
		"entry": 19, "rpc": 8, "read": 13, "write": 4, "update": 4, "delete": 2,
	})
}

// every service only accesses its own database
func TestDigotaDatabaseOwners(t *testing.T) {
	assertOwners(t, databaseOwners(runner.Get(t, "digota").AbsGraph, nil), map[string][]string{
		"orders_db.orders":     {"OrderService"},
		"payments_db.payments": {"PaymentService"},
		"products_db.products": {"ProductService"},
		"skus_db.skus":         {"SkuService"},
	})
}

// database calls inside internal helpers are attached to the service method that calls them
func TestDigotaHelperCallsAreInlined(t *testing.T) {
	g := runner.Get(t, "digota").AbsGraph

	// storageGetOne and storageUpdate are helpers of OrderService
	if find := getEdge(t, g, "OrderService.Pay", "orders_db.orders", "FindOne"); find.GetOpType() != common.OP_READ {
		t.Errorf("Pay FindOne op = %s, want read", common.OperationTypeToString(find.GetOpType()))
	}
	if replace := getEdge(t, g, "OrderService.Pay", "orders_db.orders", "ReplaceOne"); replace.GetOpType() != common.OP_UPDATE {
		t.Errorf("Pay ReplaceOne op = %s, want update", common.OperationTypeToString(replace.GetOpType()))
	}
	// rpc made from an internal helper (getLockedOrderItems)
	getEdge(t, g, "OrderService.Pay", "SkuService.Get", "Get")

	id := g.GetNodeByName("OrderService.Pay").GetParameterByNameIfExists("id")
	assertPrimaryTaint(t, id, "_obj", "orders_db.orders.Id", common.OP_READ)
	assertPrimaryTaint(t, id, "_obj", "orders_db.orders.Id", common.OP_UPDATE)
}

// SkuService.New checks the parent product with ProductService.Get and stores it in the sku,
// which is how skus_db.skus.Parent later references products_db.products.Id
func TestDigotaSkuParentFlow(t *testing.T) {
	g := runner.Get(t, "digota").AbsGraph

	getProduct := getEdge(t, g, "SkuService.New", "ProductService.Get", "Get")
	insert := getEdge(t, g, "SkuService.New", "skus_db.skus", "InsertOne")

	parent := getProduct.GetArgumentAt(0)
	assertPrimaryTaint(t, parent, "_obj", "skus_db.skus.Parent", common.OP_WRITE)
	assertTrace(t, insert.GetArgumentAt(0), "_obj.Parent", getProduct, "ProductService.Get.parent")

	// on the callee side, the product id is read from products_db
	assertPrimaryTaint(t, g.GetNodeByName("ProductService.Get").GetParamAt(0), "_obj", "products_db.products.Id", common.OP_READ)
}

// parameters that never reach a database are not tainted
func TestDigotaUnusedParameterIsNotTainted(t *testing.T) {
	g := runner.Get(t, "digota").AbsGraph

	newSku := g.GetNodeByName("SkuService.New")
	assertNotTainted(t, newSku.GetParameterByNameIfExists("image"))
	assertPrimaryTaint(t, newSku.GetParameterByNameIfExists("price"), "_obj", "skus_db.skus.Price", common.OP_WRITE)
}
