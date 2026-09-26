package detection_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/detection"
	"github.com/aletheia-microservices/aletheia/internal/app"
)

// resetDetectionConfig clears the global detection config and restores it when the test ends
func resetDetectionConfig(t *testing.T) {
	previous := detection.Config
	detection.Config = detection.InputConfig{}
	t.Cleanup(func() { detection.Config = previous })
}

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadInputConfig(t *testing.T) {
	resetDetectionConfig(t)
	path := writeConfigFile(t, `app: sockshop
ignore_foreignkeys:
  - order_db.orders.CustomerID
ignore_cascade:
  - database: order_db
    entity: orders
    trigger_database: cart_db
    trigger_entity: carts
`)

	detection.LoadInputConfig("sockshop", path)

	want := detection.InputConfig{
		App:               "sockshop",
		IgnoreForeignKeys: []string{"order_db.orders.CustomerID"},
		IgnoreCascade: []detection.IgnoreCascadeEntry{
			{Database: "order_db", Entity: "orders", TriggerDatabase: "cart_db", TriggerEntity: "carts"},
		},
	}
	if !reflect.DeepEqual(detection.Config, want) {
		t.Errorf("Config = %+v, want %+v", detection.Config, want)
	}
}

type fatalExit struct{ code int }

func TestLoadInputConfigRejectsOtherApp(t *testing.T) {
	resetDetectionConfig(t)
	path := writeConfigFile(t, "app: sockshop\n")

	// logrus.Fatalf exits the process, so make it panic instead
	logger := logrus.StandardLogger()
	previousExit := logger.ExitFunc
	logger.ExitFunc = func(code int) { panic(fatalExit{code}) }
	t.Cleanup(func() { logger.ExitFunc = previousExit })

	defer func() {
		if _, ok := recover().(fatalExit); !ok {
			t.Errorf("LoadInputConfig must stop when the config is for another app")
		}
	}()
	detection.LoadInputConfig("postnotification", path)
}

// fixedDetector returns fixed results
type fixedDetector struct {
	detection.Detector
	results string
}

func (d *fixedDetector) ComputeResults(*app.App) {}
func (d *fixedDetector) GetResults() string      { return d.results }
func (d *fixedDetector) GetTypeString() string   { return "fixed" }

func TestSaveResults(t *testing.T) {
	t.Chdir(t.TempDir())
	a := app.NewApp("postnotification")
	d := &fixedDetector{results: "[NUM_WARNINGS = 0]\n"}
	path := filepath.Join("output", "postnotification", "analysis", "fixed.txt")

	for i, step := range []struct {
		results string
		want    string
	}{
		{"[NUM_WARNINGS = 0]\n", "(modified)"},
		{"[NUM_WARNINGS = 0]\n", "(unmodified)"},
		{"[NUM_WARNINGS = 1]\n", "(modified)"},
	} {
		d.results = step.results
		printed := detection.SaveResults(a, d)

		if len(printed) != 1 || !strings.Contains(printed[0], step.want) {
			t.Errorf("run %d: SaveResults printed %q, want it to contain %q", i, printed, step.want)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("run %d: %v", i, err)
		}
		if string(content) != step.results {
			t.Errorf("run %d: %s contains %q, want %q", i, path, content, step.results)
		}
	}
}

// recordingDetector records the calls it receives from the iterator
type recordingDetector struct {
	detection.Detector
	calls []string
}

func (d *recordingDetector) record(format string, args ...any) {
	d.calls = append(d.calls, fmt.Sprintf(format, args...))
}

func (d *recordingDetector) OnNewRun(*app.App) { d.record("new run") }
func (d *recordingDetector) OnEndRun(*app.App) { d.record("end run") }
func (d *recordingDetector) OnNewRequest(node *abstractgraph.AbstractNode, reqIdx int) {
	d.record("new request %d: %s", reqIdx, node)
}
func (d *recordingDetector) OnEndRequest(*app.App) { d.record("end request") }
func (d *recordingDetector) OnNewNode(_ *app.App, node *abstractgraph.AbstractNode) {
	d.record("new node: %s", node)
}
func (d *recordingDetector) OnEndNode(_ *app.App, node *abstractgraph.AbstractNode) {
	d.record("end node: %s", node)
}
func (d *recordingDetector) OnRead(_ *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	d.record("read %d: %s", reqIdx, edge)
}
func (d *recordingDetector) OnWrite(_ *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	d.record("write %d: %s", reqIdx, edge)
}
func (d *recordingDetector) OnUpdate(_ *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	d.record("update %d: %s", reqIdx, edge)
}
func (d *recordingDetector) OnDelete(_ *app.App, reqIdx int, edge *abstractgraph.AbstractEdge) {
	d.record("delete %d: %s", reqIdx, edge)
}

// newPostnotificationStorageGraph builds two postnotification requests:
// client -> UploadService.UploadPost -> StorageService.StorePost -> posts_db.post (write)
// client -> StorageService.ReadPost -> posts_db.post (read)
func newPostnotificationStorageGraph(a *app.App) *abstractgraph.AbstractCallGraph {
	graph := abstractgraph.NewAbstractCallGraph(a)
	client := abstractgraph.NewAbstractNode("client", abstractgraph.NODE_CLIENT, "", "", "", "")
	upload := serviceNode("UploadService", "UploadPost")
	store := serviceNode("StorageService", "StorePost")
	read := serviceNode("StorageService", "ReadPost")
	posts := databaseNode("posts_db", "post")
	for _, node := range []*abstractgraph.AbstractNode{client, upload, store, read, posts} {
		graph.AddNode(node.GetName(), node)
	}
	graph.AddEdge(abstractgraph.NewAbstractEdge("", "postnotification.UploadService.UploadPost", "UploadPost", client, upload, common.OP_UNDEFINED, abstractgraph.EDGE_SERVICE_ENTRYPOINT))
	graph.AddEdge(abstractgraph.NewAbstractEdge("", "postnotification.StorageService.ReadPost", "ReadPost", client, read, common.OP_UNDEFINED, abstractgraph.EDGE_SERVICE_ENTRYPOINT))
	graph.AddEdge(abstractgraph.NewAbstractEdge("t3", "UploadService.UploadPost.t3", "StorePost", upload, store, common.OP_UNDEFINED, abstractgraph.EDGE_SERVICE_RPC))
	graph.AddEdge(abstractgraph.NewAbstractEdge("t21", "StorageService.StorePost.t21", "InsertOne", store, posts, common.OP_WRITE, abstractgraph.EDGE_DATABASE_CALL))
	graph.AddEdge(abstractgraph.NewAbstractEdge("t7", "StorageService.ReadPost.t7", "FindOne", read, posts, common.OP_READ, abstractgraph.EDGE_DATABASE_CALL))
	return graph
}

func TestIteratorCallsDetectorsInOrder(t *testing.T) {
	a := newApp(map[string]string{"posts_db": "NoSQLDatabase"})
	d := &recordingDetector{}
	iterator := detection.NewIterator(a, newPostnotificationStorageGraph(a), d)

	iterator.Run(detection.PHASE_2_PATTERN_DETECTOR)

	want := []string{
		"new run",
		"new request 0: UploadService.UploadPost",
		"new node: UploadService.UploadPost",
		"new node: StorageService.StorePost",
		"write 0: StorageService.StorePost() ... posts_db.post.InsertOne()",
		"end node: StorageService.StorePost",
		"end node: UploadService.UploadPost",
		"end request",
		"new request 1: StorageService.ReadPost",
		"new node: StorageService.ReadPost",
		"read 1: StorageService.ReadPost() ... posts_db.post.FindOne()",
		"end node: StorageService.ReadPost",
		"end request",
		"end run",
	}
	if !reflect.DeepEqual(d.calls, want) {
		t.Errorf("detector calls:\n%s\nwant:\n%s", strings.Join(d.calls, "\n"), strings.Join(want, "\n"))
	}
}

func TestIteratorCallsDetectorsOnlyInPatternDetectionPhase(t *testing.T) {
	a := newApp(map[string]string{"posts_db": "NoSQLDatabase"})
	d := &recordingDetector{}
	iterator := detection.NewIterator(a, newPostnotificationStorageGraph(a), d)

	iterator.Run(detection.PHASE_1_SCHEMA_BUILDER)

	if len(d.calls) != 0 {
		t.Errorf("schema building must not call detectors, got:\n%s", strings.Join(d.calls, "\n"))
	}
}

func TestIteratorEndsEveryRequestItStarts(t *testing.T) {
	a := newApp(map[string]string{"posts_db": "NoSQLDatabase"})
	graph := abstractgraph.NewAbstractCallGraph(a)
	client := abstractgraph.NewAbstractNode("client", abstractgraph.NODE_CLIENT, "", "", "", "")
	run := serviceNode("NotifyService", "Run")
	read := serviceNode("StorageService", "ReadPost")
	posts := databaseNode("posts_db", "post")
	for _, node := range []*abstractgraph.AbstractNode{client, run, read, posts} {
		graph.AddNode(node.GetName(), node)
	}
	graph.AddEdge(abstractgraph.NewAbstractEdge("", "postnotification.NotifyService.Run", "Run", client, run, common.OP_UNDEFINED, abstractgraph.EDGE_SERVICE_ENTRYPOINT))
	graph.AddEdge(abstractgraph.NewAbstractEdge("", "postnotification.StorageService.ReadPost", "ReadPost", client, read, common.OP_UNDEFINED, abstractgraph.EDGE_SERVICE_ENTRYPOINT))
	graph.AddEdge(abstractgraph.NewAbstractEdge("t7", "StorageService.ReadPost.t7", "FindOne", read, posts, common.OP_READ, abstractgraph.EDGE_DATABASE_CALL))
	d := &recordingDetector{}

	detection.NewIterator(a, graph, d).Run(detection.PHASE_2_PATTERN_DETECTOR)

	var started, ended int
	for _, call := range d.calls {
		if strings.HasPrefix(call, "new request") {
			started++
		} else if call == "end request" {
			ended++
		}
	}
	if started != ended {
		t.Errorf("iterator started %d requests but ended %d:\n%s", started, ended, strings.Join(d.calls, "\n"))
	}
}
