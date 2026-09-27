package detection

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/tests/runner"
)

func TestSockshopForeignKeyCascadeIgnoreConfig(t *testing.T) {
	a := runner.Get(t, "sockshop")
	assertWarnings(t, a, "foreign-key-cascade", 3, "database={order_db}, entity={orders}, pending_fields={Items}")

	// config/sockshop.yaml ignores missing cascades on order_db.orders triggered by cart_db.carts deletes
	a = runner.GetWithConfig(t, "sockshop", "config/sockshop.yaml")
	assertWarnings(t, a, "foreign-key-cascade", 0)
	// other detectors are not affected by the cascade config
	assertWarnings(t, a, "foreign-key-concurrency", 3)
}

// orders store the user's address and card and the items of the cart they were created from
func TestSockshopOrdersReferenceUserAndCart(t *testing.T) {
	a := runner.Get(t, "sockshop")
	assertConstraint(t, a, "FOREIGN_KEY order_db.orders.Address REFERENCES user_db.address.Address", true)
	assertConstraint(t, a, "FOREIGN_KEY order_db.orders.Card REFERENCES user_db.card.Card", true)
	assertConstraint(t, a, "FOREIGN_KEY order_db.orders.Items REFERENCES cart_db.carts.Items", true)
	assertConstraint(t, a, "FOREIGN_KEY cart_db.carts.Items[*].ID REFERENCES catalogue_db.sock.SockID", true)
}

// RI-1: deleting a cart (directly or when merging carts at login and register) leaves the orders
// created from it with the items of the deleted cart
func TestSockshopDeleteCartCascade(t *testing.T) {
	a := runner.Get(t, "sockshop")
	assertWarnings(t, a, "foreign-key-cascade", 3,
		"delete: Frontend.DeleteCart() ... CartService.DeleteCart() ... cart_db.carts.DeleteMany()\n\tmissing cascade #1: database={order_db}, entity={orders}, pending_fields={Items}",
		"delete: Frontend.Login() ... CartService.MergeCarts() ... cart_db.carts.DeleteOne()",
		"delete: Frontend.Register() ... CartService.MergeCarts() ... cart_db.carts.DeleteOne()",
	)
}
