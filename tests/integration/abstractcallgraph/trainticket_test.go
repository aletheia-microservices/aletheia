package abstractcallgraph

import (
	"slices"
	"strings"
	"testing"

	"analyzer/pkg/analysis/common"
	"analyzer/tests/runner"
)

// trainticket: the largest app, with 40+ services; PreserveService books a ticket by calling most
// of them, and CancelService / RebookService change existing orders

func TestTrainTicketCounts(t *testing.T) {
	g := runner.Get(t, "trainticket").AbsGraph
	assertCounts(t, g, map[string]int{
		"service nodes": 126, "database nodes": 17,
		"entry": 58, "rpc": 112, "read": 36, "write": 14, "update": 15, "delete": 9,
	})
	if got := g.ComputeAndGetNumCallGraphs(); got != 58 {
		t.Errorf("call graphs = %d, want 58", got)
	}
}

// clients only call the Dashboard, the admin services and the delivery worker
func TestTrainTicketEntrypoints(t *testing.T) {
	g := runner.Get(t, "trainticket").AbsGraph
	frontends := []string{"Dashboard", "AdminBasicInfoService", "AdminOrderService", "AdminRouteService", "AdminTravelService", "AdminUserService", "DeliveryService"}
	for _, entry := range entrypoints(g) {
		if !slices.Contains(frontends, strings.Split(entry, ".")[0]) {
			t.Errorf("unexpected entrypoint %s", entry)
		}
	}
	assertContains(t, "entrypoints", entrypoints(g),
		"Dashboard.PreserveTicketConfirm",
		"Dashboard.CancelOrder",
		"Dashboard.Rebook",
		"DeliveryService.Run",
	)
	assertContains(t, "rpc calls", rpcCalls(g),
		"Dashboard.PreserveTicketConfirm -> PreserveService.Preserve",
		"Dashboard.CancelOrder -> CancelService.CancelOrder",
		"Dashboard.Rebook -> RebookService.Rebook",
	)
}

// booking a ticket calls every service involved in an order
func TestTrainTicketPreserveCalls(t *testing.T) {
	g := runner.Get(t, "trainticket").AbsGraph
	assertList(t, "rpc calls from Preserve", withPrefix(rpcCalls(g), "PreserveService.Preserve ->"), []string{
		"PreserveService.Preserve -> AssuranceService.Create",
		"PreserveService.Preserve -> BasicService.QueryForTravel",
		"PreserveService.Preserve -> ConsignService.InsertConsign",
		"PreserveService.Preserve -> ContactsService.FindContactsById",
		"PreserveService.Preserve -> FoodService.CreateFoodOrder",
		"PreserveService.Preserve -> OrderService.CreateNewOrder",
		"PreserveService.Preserve -> SeatService.DistributeSeat",
		"PreserveService.Preserve -> TravelService.GetTripAllDetailInfo",
		"PreserveService.Preserve -> UserService.FindByUserID",
	})
	assertList(t, "rpc calls from CancelOrder", withPrefix(rpcCalls(g), "CancelService.CancelOrder ->"), []string{
		"CancelService.CancelOrder -> InsidePaymentService.Drawback",
		"CancelService.CancelOrder -> OrderService.GetOrderById",
		"CancelService.CancelOrder -> OrderService.SaveOrderInfo",
		"CancelService.CancelOrder -> UserService.FindByUserID",
	})
}

func TestTrainTicketWritesAndDeletes(t *testing.T) {
	g := runner.Get(t, "trainticket").AbsGraph
	assertList(t, "writes and deletes", withPrefix(databaseCalls(g, nil), "write ", "delete "), []string{
		"delete ConfigService.Delete -> config_db.config.DeleteOne",
		"delete ContactsService.Delete -> contacts_db.contacts.DeleteOne",
		"delete OrderService.DeleteOrder -> order_db.order.DeleteOne",
		"delete PriceService.DeletePriceConfig -> price_db.price_config.DeleteOne",
		"delete RouteService.DeleteRoute -> route_db.route.DeleteOne",
		"delete StationService.DeleteStation -> station_db.station.DeleteOne",
		"delete TrainService.Delete -> train_db.train.DeleteOne",
		"delete TravelService.DeleteTrip -> travel_db.trip.DeleteOne",
		"delete UserService.DeleteUser -> user_db.user.DeleteOne",
		"write AssuranceService.Create -> assurance_db.assurance.InsertOne",
		"write ConfigService.Create -> config_db.config.InsertOne",
		"write ConsignService.InsertConsign -> consign_db.consign_record.InsertOne",
		"write ContactsService.CreateContacts -> contacts_db.contacts.InsertOne",
		"write DeliveryService.Run -> delivery_db.delivery.InsertOne",
		"write FoodService.CreateFoodOrder -> delivery_queue.notification.Push",
		"write FoodService.CreateFoodOrder -> food_db.food_order.InsertOne",
		"write OrderService.CreateNewOrder -> order_db.order.InsertOne",
		"write PriceService.CreateNewPriceConfig -> price_db.price_config.InsertOne",
		"write RouteService.CreateAndModify -> route_db.route.InsertOne",
		"write StationService.CreateStation -> station_db.station.InsertOne",
		"write TrainService.Create -> train_db.train.InsertOne",
		"write TravelService.CreateTrip -> travel_db.trip.InsertOne",
		"write UserService.SaveUser -> user_db.user.InsertOne",
	})
}

// every service only accesses its own database
func TestTrainTicketDatabaseOwners(t *testing.T) {
	assertOwners(t, databaseOwners(runner.Get(t, "trainticket").AbsGraph, nil), map[string][]string{
		"assurance_db.assurance":    {"AssuranceService"},
		"config_db.config":          {"ConfigService"},
		"consign_db.consign_record": {"ConsignService"},
		"consignprice_db.consign":   {"ConsignPriceService"},
		"contacts_db.contacts":      {"ContactsService"},
		"delivery_db.delivery":      {"DeliveryService"},
		"food_db.food_order":        {"FoodService"},
		"inside_payment_db.money":   {"InsidePaymentService"},
		"order_db.order":            {"OrderService"},
		"payment_db.payment":        {"PaymentService"},
		"price_db.price_config":     {"PriceService"},
		"route_db.route":            {"RouteService"},
		"station_db.station":        {"StationService"},
		"train_db.train":            {"TrainService"},
		"travel_db.trip":            {"TravelService"},
		"user_db.user":              {"UserService"},
		// the queue is shared by the producer and the consumer
		"delivery_queue.notification": {"DeliveryService", "FoodService"},
	})
}

// the order created by Preserve is built from the values returned by other services
func TestTrainTicketPreserveBuildsOrderFromOtherServices(t *testing.T) {
	g := runner.Get(t, "trainticket").AbsGraph

	contacts := getEdge(t, g, "PreserveService.Preserve", "ContactsService.FindContactsById", "FindContactsById")
	seat := getEdge(t, g, "PreserveService.Preserve", "SeatService.DistributeSeat", "DistributeSeat")
	trip := getEdge(t, g, "PreserveService.Preserve", "TravelService.GetTripAllDetailInfo", "GetTripAllDetailInfo")
	create := getEdge(t, g, "PreserveService.Preserve", "OrderService.CreateNewOrder", "CreateNewOrder")

	order := create.GetArgumentAt(0)
	assertTrace(t, order, "_obj.ContactsName", contacts, ".Name")
	assertTrace(t, order, "_obj.SeatNumber", seat, ".SeatNo")
	assertTrace(t, order, "_obj.FromStation", trip, ".From")
	for _, call := range []string{contacts.GetT(), seat.GetT(), trip.GetT()} {
		if !isBefore(call, create.GetT()) {
			t.Errorf("call at t=%s must happen before CreateNewOrder (t=%s)", call, create.GetT())
		}
	}

	// on the callee side, the order is written to order_db as a whole
	assertPrimaryTaint(t, g.GetNodeByName("OrderService.CreateNewOrder").GetParamAt(0), "_obj", "order_db.order", common.OP_WRITE)
}
