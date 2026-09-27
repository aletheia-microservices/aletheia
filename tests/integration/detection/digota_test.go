package detection

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/tests/runner"
)

// SkuService.New stores the parent product checked with ProductService.Get, and orders store the
// skus of their items, so both are inferred as foreign keys
func TestDigotaForeignKeysFromDataFlow(t *testing.T) {
	a := runner.Get(t, "digota")
	assertConstraint(t, a, "FOREIGN_KEY skus_db.skus.Parent REFERENCES products_db.products.Id", true)
	assertConstraint(t, a, "FOREIGN_KEY orders_db.orders.Items[*].Parent REFERENCES skus_db.skus.Id", true)
}

// RI-1: deleting products/skus leaves dangling references in skus/orders
func TestDigotaForeignKeyCascade(t *testing.T) {
	a := runner.Get(t, "digota")
	assertWarnings(t, a, "foreign-key-cascade", 2,
		"delete: ProductService.Delete() ... products_db.products.DeleteOne()\n\tmissing cascade #1: database={skus_db}, entity={skus}, pending_fields={Parent}",
		"delete: SkuService.Delete() ... skus_db.skus.DeleteOne()\n\tmissing cascade #2: database={orders_db}, entity={orders}, pending_fields={Items[*].Parent}",
	)
}

// RI-2: a referenced object can be deleted concurrently with a write that references it
func TestDigotaForeignKeyConcurrency(t *testing.T) {
	a := runner.Get(t, "digota")
	assertWarnings(t, a, "foreign-key-concurrency", 2,
		"delete: ProductService.Delete() ... products_db.products.DeleteOne()\n\twrite #1: SkuService.New() ... SkuService.New() ... skus_db.skus.InsertOne()",
		"- database={orders_db}, entity={orders}, written_fields={Items[*].Parent}",
	)
}
