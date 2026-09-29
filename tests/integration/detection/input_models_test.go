package detection

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aletheia-microservices/aletheia/tests/runner"
)

// TestInputModelOutput checks that the input model of each app in input-models/ produces the same warnings and
// constraints as the Blueprint app it describes (without -update, since both share the expected output)
func TestInputModelOutput(t *testing.T) {
	for input, appname := range map[string]string{
		"input-models/postnotification": "postnotification",
	} {
		t.Run(appname, func(t *testing.T) {
			a := runner.GetFromInput(t, input)
			expectedDir := filepath.Join("tests", "expected", appname)
			got := map[string]string{"constraints": a.App.ConstraintsString()}
			for _, detectorType := range detectorTypes {
				got[detectorType] = a.Results[detectorType]
			}
			for name, result := range got {
				want, err := os.ReadFile(filepath.Join(expectedDir, name+".txt"))
				if err != nil {
					t.Fatal(err)
				}
				if result != string(want) {
					t.Errorf("%s differs from %s\n--- got:\n%s\n--- want:\n%s", name, expectedDir, result, want)
				}
			}
		})
	}
}

// the tinyshop frontend only calls services, but the product id it reads from ProductService and passes to
// CartService must be inferred as a foreign key, so that deleting a product warns about its cart items
func TestInputModelTinyShop(t *testing.T) {
	a := runner.GetFromInput(t, "input-models/tinyshop")
	assertConstraint(t, a, "FOREIGN_KEY cart_db.cart.ProductID REFERENCES product_db.product.ProductID", true)
	assertWarnings(t, a, "foreign-key-cascade", 1,
		"delete: Frontend.DeleteProduct() ... ProductService.DeleteProduct() ... product_db.product.DeleteOne()",
		"missing cascade #1: database={cart_db}, entity={cart}, pending_fields={ProductID}",
	)
	assertWarnings(t, a, "foreign-key-concurrency", 1,
		"write #1: Frontend.AddToCart() ... CartService.AddItem() ... cart_db.cart.InsertOne()",
	)
	for _, detectorType := range []string{"foreign-key-coordination", "primary-key-coordination", "uniqueness-concurrency"} {
		assertWarnings(t, a, detectorType, 0)
	}
}
