package detection_test

import (
	"testing"

	"analyzer/pkg/analysis/common"
	"analyzer/pkg/analysis/system-level/abstractgraph"
	"analyzer/pkg/analysis/system-level/detection"
	"analyzer/pkg/analysis/system-level/detection/constraints/foreignkeycascade"
	"analyzer/pkg/app"
)

// newSimpleshopApp builds the simpleshop schema, where inventory_db.inventory.ID references product_db.product.ID
func newSimpleshopApp() *app.App {
	a := newApp(map[string]string{"product_db": "NoSQLDatabase", "inventory_db": "NoSQLDatabase"})
	addForeignKey(field(a, "inventory_db.inventory.ID"), field(a, "product_db.product.ID"))
	return a
}

// deleteProduct returns ProductService.DeleteProduct and its delete on product_db.product
func deleteProduct() (*abstractgraph.AbstractNode, *abstractgraph.AbstractEdge) {
	entry := serviceNode("ProductService", "DeleteProduct")
	call := databaseCall("delete_product", "DeleteOne", entry, databaseNode("product_db", "product"), common.OP_DELETE,
		object(primary("product_db.product.ID", "delete_product", common.OP_DELETE)))
	return entry, call
}

// deleteInventory returns a delete on inventory_db.inventory whose argument comes from the product delete
func deleteInventory() *abstractgraph.AbstractEdge {
	return databaseCall("delete_inventory", "DeleteOne", serviceNode("InventoryService", "DeleteInventory"), databaseNode("inventory_db", "inventory"), common.OP_DELETE,
		object(
			primary("inventory_db.inventory.ID", "delete_inventory", common.OP_DELETE),
			secondary("product_db.product.ID", "delete_product", common.OP_DELETE),
		))
}

func TestCascadeReportsMissingDelete(t *testing.T) {
	a := newSimpleshopApp()
	d := foreignkeycascade.NewDetector()
	entry, del := deleteProduct()

	runRequest(d, a, 0, entry, del)

	got := results(d, a)
	wantWarnings(t, got, 1)
	wantLines(t, got,
		"delete: ProductService.DeleteProduct() ... product_db.product.DeleteOne()",
		"\tmissing cascade #1: database={inventory_db}, entity={inventory}, pending_fields={ID}",
	)
}

func TestCascadeDeleteInReferencingDatabaseClearsWarning(t *testing.T) {
	a := newSimpleshopApp()
	d := foreignkeycascade.NewDetector()
	entry, del := deleteProduct()

	runRequest(d, a, 0, entry, del, deleteInventory())

	wantWarnings(t, results(d, a), 0)
}

func TestCascadeDeleteInAnotherRequestDoesNotClearWarning(t *testing.T) {
	a := newSimpleshopApp()
	d := foreignkeycascade.NewDetector()
	entry, del := deleteProduct()

	runRequest(d, a, 0, entry, del)
	runRequest(d, a, 1, serviceNode("InventoryService", "DeleteInventory"), deleteInventory())

	wantWarnings(t, results(d, a), 1)
}

func TestCascadeSkippedWhenSameRequestWritesReference(t *testing.T) {
	a := newSimpleshopApp()
	d := foreignkeycascade.NewDetector()
	entry, del := deleteProduct()
	write := databaseCall("insert_inventory", "InsertOne", serviceNode("InventoryService", "AddInventory"), databaseNode("inventory_db", "inventory"), common.OP_WRITE,
		object(primary("inventory_db.inventory.ID", "insert_inventory", common.OP_WRITE)))

	runRequest(d, a, 0, entry, write, del)

	wantWarnings(t, results(d, a), 0)
}

func TestCascadeIgnoresReferencesFromQueues(t *testing.T) {
	a := newApp(map[string]string{"posts_db": "NoSQLDatabase", "notifications_queue": "Queue"})
	addForeignKey(field(a, "notifications_queue.notification.PostID"), field(a, "posts_db.post.PostID"))
	d := foreignkeycascade.NewDetector()
	entry := serviceNode("UploadService", "DeletePost")
	del := databaseCall("delete_post", "DeleteOne", entry, databaseNode("posts_db", "post"), common.OP_DELETE,
		object(primary("posts_db.post.PostID", "delete_post", common.OP_DELETE)))

	runRequest(d, a, 0, entry, del)

	wantWarnings(t, results(d, a), 0)
}

func TestCascadeIgnoredByDetectionConfig(t *testing.T) {
	previous := detection.Config
	t.Cleanup(func() { detection.Config = previous })

	tests := []struct {
		name  string
		entry detection.IgnoreCascadeEntry
		want  int
	}{
		{"entity", detection.IgnoreCascadeEntry{Database: "inventory_db", Entity: "inventory"}, 0},
		{"entity and trigger", detection.IgnoreCascadeEntry{Database: "inventory_db", Entity: "inventory", TriggerDatabase: "product_db", TriggerEntity: "product"}, 0},
		{"other trigger", detection.IgnoreCascadeEntry{Database: "inventory_db", Entity: "inventory", TriggerDatabase: "product_db", TriggerEntity: "category"}, 1},
		// a trigger takes effect only when both trigger fields are set
		{"trigger database only", detection.IgnoreCascadeEntry{Database: "inventory_db", Entity: "inventory", TriggerDatabase: "other_db"}, 0},
		{"other entity", detection.IgnoreCascadeEntry{Database: "inventory_db", Entity: "stock"}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detection.Config = detection.InputConfig{App: "test", IgnoreCascade: []detection.IgnoreCascadeEntry{tt.entry}}
			a := newSimpleshopApp()
			d := foreignkeycascade.NewDetector()
			entry, del := deleteProduct()

			runRequest(d, a, 0, entry, del)

			wantWarnings(t, results(d, a), tt.want)
		})
	}
}

func TestCascadeTypeString(t *testing.T) {
	if got := foreignkeycascade.NewDetector().GetTypeString(); got != "foreign-key-cascade" {
		t.Errorf("GetTypeString() = %q, want foreign-key-cascade", got)
	}
}

func TestCascadeSkipOnlyAppliesToReferencedDelete(t *testing.T) {
	a := newApp(map[string]string{"product_db": "NoSQLDatabase", "inventory_db": "NoSQLDatabase", "category_db": "NoSQLDatabase"})
	addForeignKey(field(a, "inventory_db.inventory.ID"), field(a, "product_db.product.ID"))
	addForeignKey(field(a, "product_db.product.CategoryID"), field(a, "category_db.category.ID"))
	d := foreignkeycascade.NewDetector()
	entry, delProduct := deleteProduct()
	write := databaseCall("insert_inventory", "InsertOne", serviceNode("InventoryService", "AddInventory"), databaseNode("inventory_db", "inventory"), common.OP_WRITE,
		object(primary("inventory_db.inventory.ID", "insert_inventory", common.OP_WRITE)))
	delCategory := databaseCall("delete_category", "DeleteOne", serviceNode("CategoryService", "DeleteCategory"), databaseNode("category_db", "category"), common.OP_DELETE,
		object(primary("category_db.category.ID", "delete_category", common.OP_DELETE)))

	// the write of inventory.ID skips the product delete, but not the category delete, which leaves
	// products pointing to a deleted category
	runRequest(d, a, 0, entry, write, delProduct, delCategory)

	got := results(d, a)
	wantWarnings(t, got, 1)
	wantLines(t, got, "\tmissing cascade #1: database={product_db}, entity={product}, pending_fields={CategoryID}")
}
