package detection

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/tests/runner"
)

// orders store the trip they were made for (its id as train number and its start time as travel time)
func TestTrainTicketOrdersReferenceTrips(t *testing.T) {
	a := runner.Get(t, "trainticket")
	assertConstraint(t, a, "FOREIGN_KEY order_db.order.TrainNumber REFERENCES travel_db.trip.TripID", true)
	assertConstraint(t, a, "FOREIGN_KEY order_db.order.TravelTime REFERENCES travel_db.trip.StartTime", true)
}

// RI-1: deleting an order leaves the assurance, consign, delivery and food records made for it
func TestTrainTicketDeleteOrderCascade(t *testing.T) {
	a := runner.Get(t, "trainticket")
	assertWarnings(t, a, "foreign-key-cascade", 14,
		"delete: AdminOrderService.DeleteOrder() ... OrderService.DeleteOrder() ... order_db.order.DeleteOne()",
		"database={assurance_db}, entity={assurance}, pending_fields={OrderID}",
		"database={consign_db}, entity={consign_record}, pending_fields={OrderID, TargetDate}",
		"database={delivery_db}, entity={delivery}, pending_fields={OrderID}",
		"database={food_db}, entity={food_order}, pending_fields={OrderID}",
	)
}

// RI-1: UserService.DeleteUser only deletes the user, so its orders and consign records keep its account id
func TestTrainTicketDeleteUserCascade(t *testing.T) {
	a := runner.Get(t, "trainticket")
	assertWarnings(t, a, "foreign-key-cascade", 14,
		"delete: AdminUserService.DeleteUser() ... UserService.DeleteUser() ... user_db.user.DeleteOne()",
		"database={consign_db}, entity={consign_record}, pending_fields={AccountID}",
		"database={order_db}, entity={order}, pending_fields={AccountID}",
	)
}

// RI-3: Dashboard.QueryOrderWithAllInfo reads the order and then the records that reference it,
// without coordination with the writes of those records
func TestTrainTicketForeignKeyCoordination(t *testing.T) {
	a := runner.Get(t, "trainticket")
	assertWarnings(t, a, "foreign-key-coordination", 4,
		"entry request: Dashboard.QueryOrderWithAllInfo()",
		"READ (FOREIGN KEY): AssuranceService.FindAssuranceByOrderId() ... assurance_db.assurance.FindOne()",
		"- constraint: FOREIGN_KEY assurance_db.assurance.OrderID REFERENCES order_db.order.ID [MANDATORY]",
		"READ (ORIGIN): OrderService.GetOrderById() ... order_db.order.FindOne()",
	)
}

func TestTrainTicketNoPrimaryKeyOrUniquenessWarnings(t *testing.T) {
	a := runner.Get(t, "trainticket")
	assertWarnings(t, a, "primary-key-coordination", 0)
	assertWarnings(t, a, "uniqueness-concurrency", 0)
}
