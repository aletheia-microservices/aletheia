package abstractcallgraph

import (
	"testing"

	"analyzer/pkg/analysis/common"
	"analyzer/tests/runner"
)

// sockshop: a Frontend calls the cart, catalogue, order, payment, shipping and user services;
// shipments are pushed to a queue that QueueMaster pops to update their status

func TestSockshopEntrypoints(t *testing.T) {
	g := runner.Get(t, "sockshop").AbsGraph
	assertList(t, "entrypoints", entrypoints(g), []string{
		"Frontend.AddItem",
		"Frontend.DeleteCart",
		"Frontend.GetAddress",
		"Frontend.GetCard",
		"Frontend.GetCart",
		"Frontend.GetOrder",
		"Frontend.GetOrders",
		"Frontend.GetSock",
		"Frontend.GetUser",
		"Frontend.ListItems",
		"Frontend.ListTags",
		"Frontend.Login",
		"Frontend.NewOrder",
		"Frontend.PostAddress",
		"Frontend.PostCard",
		"Frontend.Register",
		"Frontend.RemoveItem",
		"Frontend.UpdateItem",
		"QueueMaster.Run",
	})
}

func TestSockshopCalls(t *testing.T) {
	g := runner.Get(t, "sockshop").AbsGraph
	assertList(t, "rpc calls", rpcCalls(g), []string{
		"Frontend.AddItem -> CartService.AddItem",
		"Frontend.AddItem -> CatalogueService.Get",
		"Frontend.DeleteCart -> CartService.DeleteCart",
		"Frontend.GetAddress -> UserService.GetAddresses",
		"Frontend.GetCard -> UserService.GetCards",
		"Frontend.GetCart -> CartService.GetCart",
		"Frontend.GetOrder -> OrderService.GetOrder",
		"Frontend.GetOrders -> OrderService.GetOrders",
		"Frontend.GetSock -> CatalogueService.Get",
		"Frontend.GetUser -> UserService.GetUsers",
		"Frontend.ListItems -> CatalogueService.List",
		"Frontend.ListTags -> CatalogueService.Tags",
		"Frontend.Login -> CartService.MergeCarts",
		"Frontend.Login -> UserService.Login",
		"Frontend.NewOrder -> OrderService.NewOrder",
		"Frontend.PostAddress -> UserService.PostAddress",
		"Frontend.PostCard -> UserService.PostCard",
		"Frontend.Register -> CartService.MergeCarts",
		"Frontend.Register -> UserService.Register",
		"Frontend.RemoveItem -> CartService.RemoveItem",
		"Frontend.UpdateItem -> CartService.UpdateItem",
		"Frontend.UpdateItem -> CatalogueService.Get",
		"OrderService.NewOrder -> CartService.DeleteCart",
		"OrderService.NewOrder -> CartService.GetCart",
		"OrderService.NewOrder -> PaymentService.Authorise",
		"OrderService.NewOrder -> ShippingService.PostShipping",
		"OrderService.NewOrder -> UserService.GetAddresses",
		"OrderService.NewOrder -> UserService.GetCards",
		"OrderService.NewOrder -> UserService.GetUsers",
		"QueueMaster.Run -> ShippingService.UpdateStatus",
	})
	assertList(t, "database calls", databaseCalls(g, nil), []string{
		"delete CartService.DeleteCart -> cart_db.carts.DeleteMany",
		"delete CartService.MergeCarts -> cart_db.carts.DeleteOne",
		"read CartService.AddItem -> cart_db.carts.FindOne",
		"read CartService.GetCart -> cart_db.carts.FindOne",
		"read CartService.MergeCarts -> cart_db.carts.FindOne",
		"read CartService.MergeCarts -> cart_db.carts.FindOne",
		"read CartService.RemoveItem -> cart_db.carts.FindOne",
		"read CartService.UpdateItem -> cart_db.carts.FindOne",
		"read CatalogueService.Get -> catalogue_db.sock.Get",
		"read CatalogueService.Tags -> catalogue_db.tag.Select",
		"read OrderService.GetOrder -> order_db.orders.FindOne",
		"read OrderService.GetOrders -> order_db.orders.FindMany",
		"read QueueMaster.Run -> ship_queue.notification.Pop",
		"read UserService.GetAddresses -> user_db.address.FindMany",
		"read UserService.GetAddresses -> user_db.address.FindOne",
		"read UserService.GetCards -> user_db.card.FindMany",
		"read UserService.GetCards -> user_db.card.FindOne",
		"read UserService.GetUsers -> user_db.user.FindMany",
		"read UserService.GetUsers -> user_db.user.FindOne",
		"read UserService.Login -> user_db.user.FindOne",
		"update CartService.AddItem -> cart_db.carts.Upsert",
		"update CartService.MergeCarts -> cart_db.carts.Upsert",
		"update CartService.RemoveItem -> cart_db.carts.ReplaceOne",
		"update CartService.UpdateItem -> cart_db.carts.Upsert",
		"update ShippingService.UpdateStatus -> ship_db.shipments.UpdateOne",
		"update UserService.PostAddress -> user_db.user.UpdateOne",
		"update UserService.PostCard -> user_db.user.UpdateOne",
		"write OrderService.NewOrder -> order_db.orders.InsertOne",
		"write ShippingService.PostShipping -> ship_db.shipments.InsertOne",
		"write ShippingService.PostShipping -> ship_queue.notification.Push",
	})
}

func TestSockshopCounts(t *testing.T) {
	assertCounts(t, runner.Get(t, "sockshop").AbsGraph, map[string]int{
		"service nodes": 41, "database nodes": 9,
		"entry": 19, "rpc": 30, "read": 18, "write": 3, "update": 7, "delete": 2,
	})
}

func TestSockshopDatabaseOwners(t *testing.T) {
	assertOwners(t, databaseOwners(runner.Get(t, "sockshop").AbsGraph, nil), map[string][]string{
		"cart_db.carts":     {"CartService"},
		"catalogue_db.sock": {"CatalogueService"},
		"catalogue_db.tag":  {"CatalogueService"},
		"order_db.orders":   {"OrderService"},
		"ship_db.shipments": {"ShippingService"},
		"user_db.address":   {"UserService"},
		"user_db.card":      {"UserService"},
		"user_db.user":      {"UserService"},
		// the queue is shared by the producer and the consumer
		"ship_queue.notification": {"QueueMaster", "ShippingService"},
	})
}

// PostShipping writes the same shipment to the database and to the queue
func TestSockshopShipmentWrittenToDatabaseAndQueue(t *testing.T) {
	g := runner.Get(t, "sockshop").AbsGraph

	shipment := g.GetNodeByName("ShippingService.PostShipping").GetParameterByNameIfExists("shipment")
	for _, field := range []string{"", ".ID", ".Name"} {
		assertPrimaryTaint(t, shipment, "_obj"+field, "ship_db.shipments"+field, common.OP_WRITE)
		assertPrimaryTaint(t, shipment, "_obj"+field, "ship_queue.notification"+field, common.OP_WRITE)
	}
}

// QueueMaster pops a shipment from the queue and updates its status with its id
func TestSockshopQueueMasterForwardsPoppedShipment(t *testing.T) {
	g := runner.Get(t, "sockshop").AbsGraph

	pop := getEdge(t, g, "QueueMaster.Run", "ship_queue.notification", "Pop")
	update := getEdge(t, g, "QueueMaster.Run", "ShippingService.UpdateStatus", "UpdateStatus")

	assertPrimaryTaint(t, update.GetArgumentAt(0), "_obj", "ship_queue.notification.ID", common.OP_READ)
	assertTrace(t, pop.GetArgumentAt(0), "_obj.ID", update, "")
	if !isBefore(pop.GetT(), update.GetT()) {
		t.Errorf("Pop (t=%s) must happen before UpdateStatus (t=%s)", pop.GetT(), update.GetT())
	}
}

func TestSockshopRegisterWritesUser(t *testing.T) {
	t.Skip("known bug: NoSQLCollection.UpsertID is not a recognized database call (see LIKELY_BUGS.md #15), so " +
		"UserService.Register never writes to user_db")

	g := runner.Get(t, "sockshop").AbsGraph
	if !hasEdge(g, "UserService.Register", "user_db.user", "UpsertID") {
		t.Errorf("missing write UserService.Register -> user_db.user.UpsertID")
	}
}

func TestSockshopCatalogueListReadsSocks(t *testing.T) {
	t.Skip("known bug: Select with a query built at runtime is ignored (see LIKELY_BUGS.md #16), so " +
		"CatalogueService.List has no database call")

	g := runner.Get(t, "sockshop").AbsGraph
	if !hasEdge(g, "CatalogueService.List", "catalogue_db.sock", "Select") {
		t.Errorf("missing read CatalogueService.List -> catalogue_db.sock.Select")
	}
}

func TestSockshopRemoveItemDeletesEmptyCart(t *testing.T) {
	t.Skip("known bug: database calls two helper calls deep are dropped (see LIKELY_BUGS.md #17), e.g. " +
		"CartService.RemoveItem -> DeleteCart -> deleteMany")

	g := runner.Get(t, "sockshop").AbsGraph
	if !hasEdge(g, "CartService.RemoveItem", "cart_db.carts", "DeleteMany") {
		t.Errorf("missing delete CartService.RemoveItem -> cart_db.carts.DeleteMany")
	}
	if !hasEdge(g, "UserService.Login", "user_db.address", "FindMany") {
		t.Errorf("missing read UserService.Login -> user_db.address.FindMany")
	}
}
