package detection_test

import (
	"testing"

	"analyzer/pkg/analysis/common"
	"analyzer/pkg/analysis/system-level/abstractgraph"
	"analyzer/pkg/analysis/system-level/detection/constraints/foreignkeyconcurrency"
	"analyzer/pkg/app"
	"analyzer/pkg/config"
)

// newDigotaApp builds the digota schema, where skus_db.skus.Parent references products_db.products.Id
func newDigotaApp(mandatoryReqs ...int) *app.App {
	a := newApp(map[string]string{"products_db": "NoSQLDatabase", "skus_db": "NoSQLDatabase"})
	addForeignKey(field(a, "skus_db.skus.Parent"), field(a, "products_db.products.Id"), mandatoryReqs...)
	return a
}

// deleteDigotaProduct returns ProductService.Delete and its delete on products_db.products
func deleteDigotaProduct() (*abstractgraph.AbstractNode, *abstractgraph.AbstractEdge) {
	entry := serviceNode("ProductService", "Delete")
	call := databaseCall("delete_product", "DeleteOne", entry, databaseNode("products_db", "products"), common.OP_DELETE,
		object(primary("products_db.products.Id", "delete_product", common.OP_DELETE)))
	return entry, call
}

// newSku returns SkuService.New and its write of skus_db.skus.Parent, as the given operation
func newSku(method string, op common.DatabaseOperationType) (*abstractgraph.AbstractNode, *abstractgraph.AbstractEdge) {
	entry := serviceNode("SkuService", "New")
	call := databaseCall("insert_sku", method, entry, databaseNode("skus_db", "skus"), op,
		object(primary("skus_db.skus.Parent", "insert_sku", op)))
	return entry, call
}

func TestConcurrencyReportsWriteInAnotherRequest(t *testing.T) {
	a := newDigotaApp()
	d := foreignkeyconcurrency.NewDetector()
	delEntry, del := deleteDigotaProduct()
	writeEntry, write := newSku("InsertOne", common.OP_WRITE)

	runRequest(d, a, 0, delEntry, del)
	runRequest(d, a, 1, writeEntry, write)

	got := results(d, a)
	wantWarnings(t, got, 1)
	wantLines(t, got,
		"delete: ProductService.Delete() ... products_db.products.DeleteOne()",
		"\twrite #1: SkuService.New() ... SkuService.New() ... skus_db.skus.InsertOne()",
		"\t\t- database={skus_db}, entity={skus}, written_fields={Parent}",
	)
}

func TestConcurrencyIgnoresWriteInSameRequest(t *testing.T) {
	a := newDigotaApp()
	d := foreignkeyconcurrency.NewDetector()
	entry, del := deleteDigotaProduct()
	_, write := newSku("InsertOne", common.OP_WRITE)

	runRequest(d, a, 0, entry, del, write)

	wantWarnings(t, results(d, a), 0)
}

func TestConcurrencyIgnoresMandatoryForeignKeys(t *testing.T) {
	a := newDigotaApp(1)
	d := foreignkeyconcurrency.NewDetector()
	delEntry, del := deleteDigotaProduct()
	writeEntry, write := newSku("InsertOne", common.OP_WRITE)

	runRequest(d, a, 0, delEntry, del)
	runRequest(d, a, 1, writeEntry, write)

	wantWarnings(t, results(d, a), 0)
}

func TestConcurrencyUpsertCountsAsWriteOnlyWhenEnabled(t *testing.T) {
	previous := config.Global.ForeignKeyConcurrencyDetectorIncludeOnUpdates
	t.Cleanup(func() { config.Global.ForeignKeyConcurrencyDetectorIncludeOnUpdates = previous })

	for _, tt := range []struct {
		enabled bool
		want    int
	}{{false, 0}, {true, 1}} {
		config.Global.ForeignKeyConcurrencyDetectorIncludeOnUpdates = tt.enabled
		a := newDigotaApp()
		d := foreignkeyconcurrency.NewDetector()
		delEntry, del := deleteDigotaProduct()
		writeEntry, upsert := newSku("Upsert", common.OP_UPDATE)

		runRequest(d, a, 0, delEntry, del)
		runRequest(d, a, 1, writeEntry, upsert)

		wantWarnings(t, results(d, a), tt.want)
	}
}

func TestConcurrencyTypeString(t *testing.T) {
	if got := foreignkeyconcurrency.NewDetector().GetTypeString(); got != "foreign-key-concurrency" {
		t.Errorf("GetTypeString() = %q, want foreign-key-concurrency", got)
	}
}
