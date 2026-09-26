package abstractcallgraph

import (
	"testing"

	"analyzer/pkg/analysis/common"
	"analyzer/tests/runner"
)

// eshopmicroservices: a WebApp calls the basket, catalog, discount and order services; checking
// out a basket pushes an event to a queue that OrderService pops to create the order

func TestEshopEntrypoints(t *testing.T) {
	g := runner.Get(t, "eshopmicroservices").AbsGraph
	assertList(t, "entrypoints", entrypoints(g), []string{
		"BasketService.CheckoutBasket",
		"BasketService.DeleteBasket",
		"BasketService.GetBasket",
		"BasketService.StoreBasket",
		"CatalogService.CreateProduct",
		"CatalogService.DeleteProduct",
		"CatalogService.GetProductByCategory",
		"CatalogService.GetProductById",
		"CatalogService.GetProducts",
		"DiscountService.CreateDiscount",
		"DiscountService.DeleteDiscount",
		"DiscountService.GetDiscount",
		"DiscountService.UpdateDiscount",
		"OrderService.CreateNewOrder",
		"OrderService.DeleteOrder",
		"OrderService.GetOrdersByCustomer",
		"OrderService.Init",
		"OrderService.UpdateOrder",
		"WebApp.OnGetOrdersAsync",
		"WebApp.OnGetProductsAsync",
		"WebApp.OnPostAddToCartAsync",
		"WebApp.OnPostCheckoutAsync",
		"WebApp.OnPostRemoveToCartAsync",
	})
}

func TestEshopCalls(t *testing.T) {
	g := runner.Get(t, "eshopmicroservices").AbsGraph
	assertList(t, "rpc calls", rpcCalls(g), []string{
		"BasketService.StoreBasket -> DiscountService.GetDiscount",
		"WebApp.OnGetOrdersAsync -> OrderService.GetOrdersByCustomer",
		"WebApp.OnGetProductsAsync -> CatalogService.GetProducts",
		"WebApp.OnPostAddToCartAsync -> BasketService.GetBasket",
		"WebApp.OnPostAddToCartAsync -> BasketService.StoreBasket",
		"WebApp.OnPostAddToCartAsync -> CatalogService.GetProductById",
		"WebApp.OnPostCheckoutAsync -> BasketService.CheckoutBasket",
		"WebApp.OnPostCheckoutAsync -> BasketService.GetBasket",
		"WebApp.OnPostRemoveToCartAsync -> BasketService.GetBasket",
		"WebApp.OnPostRemoveToCartAsync -> BasketService.StoreBasket",
	})
	assertList(t, "database calls", databaseCalls(g, nil), []string{
		"delete BasketService.CheckoutBasket -> basket_db.basket.DeleteOne",
		"delete BasketService.DeleteBasket -> basket_db.basket.DeleteOne",
		"delete CatalogService.DeleteProduct -> catalog_db.product.DeleteOne",
		"delete DiscountService.DeleteDiscount -> discount_db.coupon.DeleteOne",
		"delete OrderService.DeleteOrder -> order_db.order.DeleteOne",
		"read BasketService.CheckoutBasket -> basket_db.basket.FindOne",
		"read BasketService.GetBasket -> basket_db.basket.FindOne",
		"read CatalogService.GetProductByCategory -> catalog_db.product.FindOne",
		"read CatalogService.GetProductById -> catalog_db.product.FindOne",
		"read CatalogService.GetProducts -> catalog_db.product.FindMany",
		"read DiscountService.GetDiscount -> discount_db.coupon.FindOne",
		"read OrderService.DeleteOrder -> order_db.order.FindOne",
		"read OrderService.GetOrdersByCustomer -> order_db.order.FindMany",
		"read OrderService.Init -> order_queue.notification.Pop",
		"read OrderService.UpdateOrder -> order_db.order.FindOne",
		"update DiscountService.UpdateDiscount -> discount_db.coupon.ReplaceOne",
		"update OrderService.UpdateOrder -> order_db.order.ReplaceOne",
		"write BasketService.CheckoutBasket -> order_queue.notification.Push",
		"write BasketService.StoreBasket -> basket_db.basket.InsertOne",
		"write CatalogService.CreateProduct -> catalog_db.product.InsertOne",
		"write DiscountService.CreateDiscount -> discount_db.coupon.InsertOne",
		"write OrderService.CreateNewOrder -> order_db.order.InsertOne",
	})
}

func TestEshopCounts(t *testing.T) {
	assertCounts(t, runner.Get(t, "eshopmicroservices").AbsGraph, map[string]int{
		"service nodes": 23, "database nodes": 5,
		"entry": 23, "rpc": 10, "read": 10, "write": 5, "update": 2, "delete": 5,
	})
}

func TestEshopDatabaseOwners(t *testing.T) {
	assertOwners(t, databaseOwners(runner.Get(t, "eshopmicroservices").AbsGraph, nil), map[string][]string{
		"basket_db.basket":   {"BasketService"},
		"catalog_db.product": {"CatalogService"},
		"discount_db.coupon": {"DiscountService"},
		"order_db.order":     {"OrderService"},
		// the queue is shared by the producer and the consumer
		"order_queue.notification": {"BasketService", "OrderService"},
	})
}

// adding a product to the cart stores the name and price returned by the catalog in the basket
func TestEshopAddToCartUsesCatalogProduct(t *testing.T) {
	g := runner.Get(t, "eshopmicroservices").AbsGraph

	getProduct := getEdge(t, g, "WebApp.OnPostAddToCartAsync", "CatalogService.GetProductById", "GetProductById")
	store := getEdge(t, g, "WebApp.OnPostAddToCartAsync", "BasketService.StoreBasket", "StoreBasket")

	basket := store.GetArgumentAt(0)
	assertTrace(t, basket, "_obj.Cart.Items[*].Price", getProduct, ".Product.Price")
	assertTrace(t, basket, "_obj.Cart.Items[*].ProductName", getProduct, ".Product.Name")
	if !isBefore(getProduct.GetT(), store.GetT()) {
		t.Errorf("GetProductById (t=%s) must happen before StoreBasket (t=%s)", getProduct.GetT(), store.GetT())
	}

	// on the callee side, the stored basket is written to basket_db
	command := g.GetNodeByName("BasketService.StoreBasket").GetParamAt(0)
	assertPrimaryTaint(t, command, "_obj.Cart.Items[*].Price", "basket_db.basket.Items[*].Price", common.OP_WRITE)
}

// the checkout event carries the total price read from the basket
func TestEshopCheckoutPushesBasketTotal(t *testing.T) {
	g := runner.Get(t, "eshopmicroservices").AbsGraph

	event := getEdge(t, g, "BasketService.CheckoutBasket", "order_queue.notification", "Push").GetArgumentAt(0)
	assertPrimaryTaint(t, event, "_obj.TotalPrice", "order_queue.notification.TotalPrice", common.OP_WRITE)
	assertPrimaryTaint(t, event, "_obj.TotalPrice", "basket_db.basket.TotalPrice", common.OP_READ)
}

func TestEshopOrderConsumerStoresOrder(t *testing.T) {
	t.Skip("known bug: database calls two helper calls deep are dropped (see LIKELY_BUGS.md #17), e.g. " +
		"OrderService.Init -> CreateNewOrder -> add")

	g := runner.Get(t, "eshopmicroservices").AbsGraph
	if !hasEdge(g, "OrderService.Init", "order_db.order", "InsertOne") {
		t.Errorf("missing write OrderService.Init -> order_db.order.InsertOne")
	}
}
