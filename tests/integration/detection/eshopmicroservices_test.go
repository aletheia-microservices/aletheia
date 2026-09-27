package detection

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/tests/runner"
)

// WebApp.OnPostAddToCartAsync reads the product with CatalogService.GetProductById and copies its id,
// name and price into the basket item stored by BasketService.StoreBasket
func TestEshopBasketItemsReferenceCatalog(t *testing.T) {
	a := runner.Get(t, "eshopmicroservices")
	assertConstraint(t, a, "FOREIGN_KEY basket_db.basket.Items[*].ProductId REFERENCES catalog_db.product.Id", true)
	assertConstraint(t, a, "FOREIGN_KEY basket_db.basket.Items[*].ProductName REFERENCES catalog_db.product.Name", true)
	assertConstraint(t, a, "FOREIGN_KEY basket_db.basket.Items[*].Price REFERENCES catalog_db.product.Price", true)
	// discounts are looked up by product name
	assertConstraint(t, a, "FOREIGN_KEY basket_db.basket.Items[*].ProductName REFERENCES discount_db.coupon.ProductName", true)
}

// RI-1: deleting a product or a discount leaves baskets with items that still reference it
func TestEshopForeignKeyCascade(t *testing.T) {
	a := runner.Get(t, "eshopmicroservices")
	assertWarnings(t, a, "foreign-key-cascade", 2,
		"delete: CatalogService.DeleteProduct() ... catalog_db.product.DeleteOne()\n\tmissing cascade #1: database={basket_db}, entity={basket}, pending_fields={Items[*].Price, Items[*].ProductId, Items[*].ProductName}",
		"delete: DiscountService.DeleteDiscount() ... discount_db.coupon.DeleteOne()\n\tmissing cascade #2: database={basket_db}, entity={basket}, pending_fields={Items[*].ProductName}",
	)
}

// RI-2: a product can be deleted while a basket that references it is being stored, both directly
// and from the web app requests that add or remove items
func TestEshopForeignKeyConcurrency(t *testing.T) {
	a := runner.Get(t, "eshopmicroservices")
	assertWarnings(t, a, "foreign-key-concurrency", 6,
		"delete: CatalogService.DeleteProduct() ... catalog_db.product.DeleteOne()",
		"write #2: WebApp.OnPostAddToCartAsync() ... BasketService.StoreBasket() ... basket_db.basket.InsertOne()",
		"write #3: WebApp.OnPostRemoveToCartAsync() ... BasketService.StoreBasket() ... basket_db.basket.InsertOne()",
	)
}

// eshopmicroservices has no reads that follow a foreign key and no unique fields
func TestEshopNoCoordinationOrUniquenessWarnings(t *testing.T) {
	a := runner.Get(t, "eshopmicroservices")
	for _, detectorType := range []string{"foreign-key-coordination", "primary-key-coordination", "uniqueness-concurrency"} {
		assertWarnings(t, a, detectorType, 0)
	}
}
