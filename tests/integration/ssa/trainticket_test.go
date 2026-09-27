package ssa

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/tests/runner"
)

// AssuranceService.Create builds an assurance for an order and inserts it: the order id is tainted
// with the field it was assigned to
func TestTrainTicketWriteFieldFromParameter(t *testing.T) {
	a := runner.Get(t, "trainticket")
	graph := getSSAGraph(t, a, "trainticket.AssuranceService.Create")

	call := getDatabaseCall(t, graph, "assurance_db.assurance", "InsertOne")
	assertDatabaseTaint(t, getParam(t, graph, "orderid"), dbTaint{objpath: "_obj", dbpath: "assurance_db.assurance.OrderID", op: common.OP_WRITE, t: call.GetT()})

	assurance := call.GetArguments()[0]
	for _, field := range []string{"", ".ID", ".OrderID"} {
		assertDatabaseTaint(t, assurance, dbTaint{objpath: "_obj" + field, dbpath: "assurance_db.assurance" + field, op: common.OP_WRITE})
	}
}

// AssuranceService.FindAssuranceByOrderId reads the assurance with the filter {OrderID: order_id}
func TestTrainTicketReadByForeignKey(t *testing.T) {
	a := runner.Get(t, "trainticket")
	graph := getSSAGraph(t, a, "trainticket.AssuranceService.FindAssuranceByOrderId")

	call := getDatabaseCall(t, graph, "assurance_db.assurance", "FindOne")
	assertDatabaseTaint(t, getParam(t, graph, "order_id"), dbTaint{objpath: "_obj", dbpath: "assurance_db.assurance.OrderID", op: common.OP_READ, readKey: true, t: call.GetT()})
	assertDatabaseTaint(t, call.GetArguments()[2], dbTaint{objpath: "_obj", dbpath: "assurance_db.assurance", op: common.OP_READ, readVal: true})
}

// UserService.DeleteUser deletes the user with the filter {UserID: userID}: the user id is the key
// of the delete
func TestTrainTicketDeleteKey(t *testing.T) {
	a := runner.Get(t, "trainticket")
	graph := getSSAGraph(t, a, "trainticket.UserService.DeleteUser")

	call := getDatabaseCall(t, graph, "user_db.user", "DeleteOne")
	if call.GetOpType() != common.OP_DELETE {
		t.Fatalf("DeleteOne op type = %s, want delete", common.OperationTypeToString(call.GetOpType()))
	}
	assertDatabaseTaint(t, getParam(t, graph, "userID"), dbTaint{objpath: "_obj", dbpath: "user_db.user.UserID", op: common.OP_DELETE, readKey: true, t: call.GetT()})
}
