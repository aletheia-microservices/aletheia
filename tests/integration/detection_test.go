package integration

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/aletheia-microservices/aletheia/tests/runner"
)

// regenerate the expected output with: go test ./tests/integration -run TestDetectionOutput -update
var update = flag.Bool("update", false, "update the expected output files in tests/expected")

// every realistic app with a blueprint spec registered in registry/apps.yaml
var analyzedApps = []string{
	"postnotification",
	"digota",
	"eshopmicroservices",
	"dsb_mediamicroservices",
	"dsb_socialnetwork",
	"sockshop",
	"trainticket",
}

// apps whose inferred constraints are known to vary between runs
var nondeterministicConstraints = map[string]string{
	"sockshop": "FOREIGN_KEY order_db.orders.Shipment.Name REFERENCES {ship_db.shipments.Name, ship_queue.notification.Name} " +
		"are only inferred on some runs (fields are attached to db.GetLastSchema(), whose order is not stable)",
}

var detectorTypes = []string{
	"foreign-key-cascade",
	"foreign-key-concurrency",
	"foreign-key-coordination",
	"primary-key-coordination",
	"uniqueness-concurrency",
}

// constraintsString returns all constraints inferred for the app, sorted and deduplicated
func constraintsString(a *runner.Analysis) string {
	var lines []string
	for _, db := range a.App.GetAllDatabases() {
		for _, schema := range db.GetSchemas() {
			for _, constraint := range schema.GetAllConstraints() {
				lines = append(lines, constraint.String())
			}
		}
	}
	sort.Strings(lines)
	lines = slices.Compact(lines)
	return strings.Join(lines, "\n") + "\n"
}

func checkExpectedOutput(t *testing.T, path string, got string) {
	t.Helper()
	if *update {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(got), 0644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading expected output (run with -update to create it): %v", err)
	}
	if got != string(want) {
		t.Errorf("output differs from %s (run with -update if the change is intended)\n--- got:\n%s\n--- want:\n%s", path, got, want)
	}
}

func runExpectedOutputTests(t *testing.T, appname string, configPath string, expectedDir string) {
	a := runner.GetWithConfig(t, appname, configPath)
	for _, detectorType := range detectorTypes {
		t.Run(detectorType, func(t *testing.T) {
			got, ok := a.Results[detectorType]
			if !ok {
				t.Fatalf("no results for detector %s", detectorType)
			}
			checkExpectedOutput(t, filepath.Join(expectedDir, detectorType+".txt"), got)
		})
	}
	t.Run("constraints", func(t *testing.T) {
		if reason, ok := nondeterministicConstraints[appname]; ok {
			t.Skip("known bug: " + reason)
		}
		checkExpectedOutput(t, filepath.Join(expectedDir, "constraints.txt"), constraintsString(a))
	})
}

// TestDetectionOutput compares the warnings of every detector (and the inferred constraints)
// against the expected output in tests/expected/{app}
func TestDetectionOutput(t *testing.T) {
	for _, appname := range analyzedApps {
		t.Run(appname, func(t *testing.T) {
			runExpectedOutputTests(t, appname, "", filepath.Join("tests", "expected", appname))
		})
	}
}

// TestDetectionOutputWithConfig is the same as TestDetectionOutput but suppresses
// warnings using the detection config files in config/{app}.yaml
func TestDetectionOutputWithConfig(t *testing.T) {
	for _, appname := range analyzedApps {
		configPath := filepath.Join("config", appname+".yaml")
		if _, err := os.Stat(configPath); err != nil {
			continue
		}
		t.Run(appname, func(t *testing.T) {
			runExpectedOutputTests(t, appname, configPath, filepath.Join("tests", "expected", appname+"_config"))
		})
	}
}

// ---------------------------------------------------------------------
// semantic checks: one known integrity violation per pattern
// ---------------------------------------------------------------------

func numWarnings(t *testing.T, result string) int {
	t.Helper()
	var n int
	for _, line := range strings.Split(result, "\n") {
		if _, err := fmt.Sscanf(line, "[NUM_WARNINGS = %d]", &n); err == nil {
			return n
		}
	}
	t.Fatalf("no warning count in result:\n%s", result)
	return -1
}

func assertWarnings(t *testing.T, a *runner.Analysis, detectorType string, want int, fragments ...string) {
	t.Helper()
	result := a.Results[detectorType]
	if got := numWarnings(t, result); got != want {
		t.Errorf("%s: %d warnings, want %d\n%s", detectorType, got, want, result)
	}
	for _, fragment := range fragments {
		if !strings.Contains(result, fragment) {
			t.Errorf("%s: missing %q in result:\n%s", detectorType, fragment, result)
		}
	}
}

func assertConstraint(t *testing.T, a *runner.Analysis, constraint string, want bool) {
	t.Helper()
	constraints := strings.Split(constraintsString(a), "\n")
	if got := slices.Contains(constraints, constraint); got != want {
		t.Errorf("constraint %q inferred = %v, want %v\nconstraints:\n%s", constraint, got, want, constraintsString(a))
	}
}

// schema extraction: the queue message stores the post id and request id written to posts_db,
// so both fields must be inferred as (mandatory) foreign keys
func TestSchemaForeignKeysFromDataFlow(t *testing.T) {
	a := runner.Get(t, "postnotification")
	assertConstraint(t, a, "FOREIGN_KEY notifications_queue.notification.PostID REFERENCES posts_db.post.PostID [MANDATORY]", true)
	assertConstraint(t, a, "FOREIGN_KEY notifications_queue.notification.ReqID REFERENCES posts_db.post.ReqID [MANDATORY]", true)
	// references must follow the data flow direction
	assertConstraint(t, a, "FOREIGN_KEY posts_db.post.PostID REFERENCES notifications_queue.notification.PostID [MANDATORY]", false)

	a = runner.Get(t, "digota")
	assertConstraint(t, a, "FOREIGN_KEY skus_db.skus.Parent REFERENCES products_db.products.Id", true)
	assertConstraint(t, a, "FOREIGN_KEY orders_db.orders.Items[*].Parent REFERENCES skus_db.skus.Id", true)
}

func TestSchemaIgnoreForeignKeysConfig(t *testing.T) {
	a := runner.GetWithConfig(t, "postnotification", "config/postnotification.yaml")
	assertConstraint(t, a, "FOREIGN_KEY notifications_queue.notification.ReqID REFERENCES posts_db.post.ReqID [MANDATORY]", false)
	assertConstraint(t, a, "FOREIGN_KEY notifications_queue.notification.PostID REFERENCES posts_db.post.PostID [MANDATORY]", true)
}

// RI-3: NotifyService reads the post referenced by the queued message without coordination
// with UploadService, which writes the post before pushing the message
func TestDetectForeignKeyCoordination(t *testing.T) {
	a := runner.Get(t, "postnotification")
	assertWarnings(t, a, "foreign-key-coordination", 1,
		"entry request: UploadService.UploadPost()",
		"READ (FOREIGN KEY): NotifyService.Run() ... notifications_queue.notification.Pop()",
		"- constraint: FOREIGN_KEY notifications_queue.notification.PostID REFERENCES posts_db.post.PostID [MANDATORY]",
		"READ (ORIGIN): StorageService.ReadPost() ... posts_db.post.FindOne()",
	)
	for _, detectorType := range []string{"foreign-key-cascade", "foreign-key-concurrency", "primary-key-coordination", "uniqueness-concurrency"} {
		assertWarnings(t, a, detectorType, 0)
	}
}

// RI-1: deleting products/skus leaves dangling references in skus/orders
func TestDetectForeignKeyCascade(t *testing.T) {
	a := runner.Get(t, "digota")
	assertWarnings(t, a, "foreign-key-cascade", 2,
		"delete: ProductService.Delete() ... products_db.products.DeleteOne()\n\tmissing cascade #1: database={skus_db}, entity={skus}, pending_fields={Parent}",
		"delete: SkuService.Delete() ... skus_db.skus.DeleteOne()\n\tmissing cascade #2: database={orders_db}, entity={orders}, pending_fields={Items[*].Parent}",
	)
}

func TestDetectForeignKeyCascadeIgnoreConfig(t *testing.T) {
	a := runner.Get(t, "sockshop")
	assertWarnings(t, a, "foreign-key-cascade", 3, "database={order_db}, entity={orders}, pending_fields={Items}")

	// config/sockshop.yaml ignores missing cascades on order_db.orders triggered by cart_db.carts deletes
	a = runner.GetWithConfig(t, "sockshop", "config/sockshop.yaml")
	assertWarnings(t, a, "foreign-key-cascade", 0)
	// other detectors are not affected by the cascade config
	assertWarnings(t, a, "foreign-key-concurrency", 3)
}

// RI-2: a referenced object can be deleted concurrently with a write that references it
func TestDetectForeignKeyConcurrency(t *testing.T) {
	a := runner.Get(t, "digota")
	assertWarnings(t, a, "foreign-key-concurrency", 2,
		"delete: ProductService.Delete() ... products_db.products.DeleteOne()\n\twrite #1: SkuService.New() ... SkuService.New() ... skus_db.skus.InsertOne()",
		"- database={orders_db}, entity={orders}, written_fields={Items[*].Parent}",
	)
}

// EI-1: APIService.ReadPage reads replicated objects by primary key from two services without coordination
func TestDetectPrimaryKeyCoordination(t *testing.T) {
	a := runner.Get(t, "dsb_mediamicroservices")
	assertConstraint(t, a, "PRIMARY KEY (movie_id_db.movie._id)", true)
	assertWarnings(t, a, "primary-key-coordination", 1,
		"entry request: APIService.ReadPage()",
		"READ: MovieInfoService.ReadMovieInfo() ... movie_info_db.movie_info.FindOne()",
		"READ: MovieIdService.ReadMovieId() ... movie_id_db.movie.FindOne()",
	)
}

// Un-1: concurrent RegisterMovie requests may write the same (unique) movie title
func TestDetectUniquenessConcurrency(t *testing.T) {
	a := runner.Get(t, "dsb_mediamicroservices")
	assertConstraint(t, a, "UNIQUE (movie_id_db.movie.Title)", true)
	assertWarnings(t, a, "uniqueness-concurrency", 1,
		"entry request: APIService.RegisterMovie()",
		"- field (constrained): movie_id_db.movie.Title (UNIQUE)",
		"affected write #1: MovieInfoService.WriteMovieInfo() ... movie_info_db.movie_info.InsertOne()",
	)
}
