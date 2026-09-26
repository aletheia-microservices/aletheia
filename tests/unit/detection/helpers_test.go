package detection_test

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/detection"
	"github.com/aletheia-microservices/aletheia/internal/app"
	"github.com/aletheia-microservices/aletheia/internal/app/backends"
	"github.com/aletheia-microservices/aletheia/internal/utils"
)

// newApp returns an app with the given databases, mapped from name to type (e.g., "NoSQLDatabase", "Queue")
func newApp(databases map[string]string) *app.App {
	a := app.NewApp("test")
	for name, typeString := range databases {
		a.AddDatabase(backends.NewDatabase(name, typeString))
	}
	return a
}

// field returns the field at path (database.entity.field), creating its entity and field if needed
func field(a *app.App, path string) *backends.Field {
	db := a.GetDatabaseByName(utils.ExtractDatabaseNameFromFieldPath(path))
	schema := db.GetOrCreateSchema(utils.ExtractSchemaNameFromFieldPath(path))
	return schema.GetOrCreateField(db, path)
}

// addForeignKey declares that from references to; the foreign key is mandatory in the given requests
func addForeignKey(from *backends.Field, to *backends.Field, mandatoryReqs ...int) {
	constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, from, to)
	// mandatory flags must be set before adding the constraint to the schema
	for _, reqIdx := range mandatoryReqs {
		constraint.EnableMandatory(reqIdx)
	}
	from.AddConstraint(constraint)
	from.GetSchema().AddConstraint(constraint)
}

func addConstraint(t backends.ConstraintType, f *backends.Field) {
	constraint := backends.NewConstraint(t, f)
	f.AddConstraint(constraint)
	f.GetSchema().AddConstraint(constraint)
}

func serviceNode(service string, method string) *abstractgraph.AbstractNode {
	return abstractgraph.NewAbstractNode(service+"."+method, abstractgraph.NODE_SERVICE, service, method, "", "")
}

func databaseNode(database string, schema string) *abstractgraph.AbstractNode {
	return abstractgraph.NewAbstractNode(database+"."+schema, abstractgraph.NODE_DATABASE, "", "", database, schema)
}

func databaseCall(id string, method string, from *abstractgraph.AbstractNode, to *abstractgraph.AbstractNode, op common.DatabaseOperationType, args ...*abstractgraph.AbstractObject) *abstractgraph.AbstractEdge {
	edge := abstractgraph.NewAbstractEdge("t1", id, method, from, to, op, abstractgraph.EDGE_DATABASE_CALL)
	for _, arg := range args {
		edge.AddArgument(arg)
	}
	return edge
}

func primary(dbpath string, callID string, op common.DatabaseOperationType) *abstractgraph.AbstractTaint {
	return abstractgraph.NewAbstractTaint("t0", dbpath, callID, op, true, false, false, false)
}

func secondary(dbpath string, callID string, op common.DatabaseOperationType) *abstractgraph.AbstractTaint {
	return abstractgraph.NewAbstractTaint("t0", dbpath, callID, op, false, false, false, false)
}

// object returns an argument holding all taints on the same object path
func object(taints ...*abstractgraph.AbstractTaint) *abstractgraph.AbstractObject {
	return abstractgraph.NewAbstractObject("t0", map[string][]*abstractgraph.AbstractTaint{"t0": taints}, map[string][]*abstractgraph.AbstractTrace{})
}

// runRequest passes a request to the detector the same way the iterator does: the entry node,
// then each database call in order, then the end of the request
func runRequest(d detection.Detector, a *app.App, reqIdx int, entry *abstractgraph.AbstractNode, calls ...*abstractgraph.AbstractEdge) {
	d.OnNewRequest(entry, reqIdx)
	for _, call := range calls {
		switch call.GetOpType() {
		case common.OP_READ, common.OP_READ_MANY:
			d.OnRead(a, reqIdx, call)
		case common.OP_WRITE:
			d.OnWrite(a, reqIdx, call)
		case common.OP_UPDATE:
			d.OnUpdate(a, reqIdx, call)
		case common.OP_DELETE:
			d.OnDelete(a, reqIdx, call)
		}
	}
	d.OnEndRequest(a)
}

// results ends the run and returns the detector's results
func results(d detection.Detector, a *app.App) string {
	d.OnEndRun(a)
	d.ComputeResults(a)
	return d.GetResults()
}

var numWarningsRegex = regexp.MustCompile(`\[NUM_WARNINGS = (\d+)\]`)

func wantWarnings(t *testing.T, results string, want int) {
	t.Helper()
	match := numWarningsRegex.FindStringSubmatch(results)
	if match == nil {
		t.Fatalf("results have no [NUM_WARNINGS = N] line:\n%s", results)
	}
	if got, _ := strconv.Atoi(match[1]); got != want {
		t.Errorf("got %d warnings, want %d:\n%s", got, want, results)
	}
}

// wantLines checks that each line appears in the results
func wantLines(t *testing.T, results string, lines ...string) {
	t.Helper()
	got := strings.Split(results, "\n")
	for _, line := range lines {
		found := false
		for _, gotLine := range got {
			if gotLine == line {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("results are missing line %q:\n%s", line, results)
		}
	}
}
