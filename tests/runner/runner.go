// Package runner runs Aletheia's stages (same as main.go) for tests, keeping the result of every
// stage (SSA graphs, abstract call graph, schema and detector results) in memory for inspection
package runner

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"

	"analyzer/pkg/analysis/service-level/ssagraph"
	"analyzer/pkg/analysis/service-level/ssagraph/parser"
	"analyzer/pkg/analysis/service-level/ssagraph/registry"
	"analyzer/pkg/analysis/service-level/ssagraph/tainter"
	"analyzer/pkg/analysis/system-level/abstractgraph"
	abstractgraphparser "analyzer/pkg/analysis/system-level/abstractgraph/parser"
	"analyzer/pkg/analysis/system-level/detection"
	"analyzer/pkg/analysis/system-level/detection/constraints/foreignkeycascade"
	"analyzer/pkg/analysis/system-level/detection/constraints/foreignkeyconcurrency"
	"analyzer/pkg/analysis/system-level/detection/constraints/keycoordination"
	"analyzer/pkg/analysis/system-level/detection/constraints/uniquenessconcurrency"
	"analyzer/pkg/app"
	appparser "analyzer/pkg/app/parser"
	"analyzer/pkg/config"
	"analyzer/pkg/utils"
)

// Analysis holds the state of every stage of the pipeline for one app so that
// tests can inspect intermediate results (unlike main, SSA graphs are not released)
type Analysis struct {
	App        *app.App
	FuncGraphs map[string]*ssagraph.SSAGraph
	AbsGraph   *abstractgraph.AbstractCallGraph
	Detectors  []detection.Detector
	// results indexed by detector type string (e.g., "foreign-key-cascade")
	Results map[string]string
}

var (
	analysesMu sync.Mutex
	analyses   = map[string]*Analysis{}
)

// Get runs the full pipeline for appname once per test binary and caches the result
func Get(t *testing.T, appname string) *Analysis {
	t.Helper()
	return GetWithConfig(t, appname, "")
}

// GetWithConfig is the same as Get but loads the detection config at
// configPath (same as the --detection_config flag) when it is not empty
func GetWithConfig(t *testing.T, appname string, configPath string) *Analysis {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping full pipeline in -short mode")
	}
	analysesMu.Lock()
	defer analysesMu.Unlock()
	key := appname + "|" + configPath
	if a, ok := analyses[key]; ok {
		return a
	}
	a := runAnalysis(t, appname, configPath)
	analyses[key] = a
	return a
}

// runAnalysis mirrors the stages of main.go without writing to output/
func runAnalysis(t *testing.T, appname string, configPath string) *Analysis {
	t.Helper()
	logrus.SetLevel(logrus.ErrorLevel)
	// the detection config is global state, so reset it for every run
	detection.Config = detection.InputConfig{}
	if configPath != "" {
		detection.LoadInputConfig(appname, configPath)
	}

	a := &Analysis{App: app.NewApp(appname), FuncGraphs: make(map[string]*ssagraph.SSAGraph)}
	appparser.Init(a.App, false)

	prog, pkgs, err := utils.BuildProgram(utils.GetAppRootPackagePath(appname))
	if err != nil {
		t.Fatalf("building program for %s: %v", appname, err)
	}
	appparser.InitServiceFields(a.App, pkgs)
	appparser.ParseSQLSchemaFromUserFile(a.App)
	appparser.ParseNoSQLSchemaFromUserFile(a.App)

	// the SSA parser always dumps the ssa code to output/{app}/ssa, so redirect it to a temp dir
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	for _, pkg := range pkgs {
		parser.RunSSAAnalysis(a.App, prog, pkg, a.FuncGraphs)
	}
	if err := os.Chdir(cwd); err != nil {
		t.Fatal(err)
	}

	var graphsLst []*ssagraph.SSAGraph
	for _, graph := range a.FuncGraphs {
		graph.Sort()
		graphsLst = append(graphsLst, graph)
	}
	registry.RegisterFields(a.App, graphsLst)

	for _, graph := range a.FuncGraphs {
		tainter.RunTainter(graph)
	}
	for _, graph := range a.FuncGraphs {
		tainter.Combine(graph, a.FuncGraphs)
	}

	a.AbsGraph = abstractgraph.NewAbstractCallGraph(a.App)
	for _, entrypoint := range a.App.GetEntrypointsShortPaths() {
		abstractgraphparser.Parse(a.AbsGraph, entrypoint, true, a.FuncGraphs)
	}

	a.Detectors = []detection.Detector{
		keycoordination.NewDetector(keycoordination.DETECTION_TYPE_PRIMARY_KEY),
		keycoordination.NewDetector(keycoordination.DETECTION_TYPE_FOREIGN_KEY),
		foreignkeycascade.NewDetector(),
		foreignkeyconcurrency.NewDetector(),
		uniquenessconcurrency.NewDetector(),
	}
	iterator := detection.NewIterator(a.App, a.AbsGraph, a.Detectors...)
	iterator.Run(detection.PHASE_1_SCHEMA_BUILDER)
	if config.Global.DualPassSchemaBuilding {
		iterator.Run(detection.PHASE_1_SCHEMA_BUILDER_READ_ONLY)
	}
	iterator.Run(detection.PHASE_2_PATTERN_DETECTOR)

	a.Results = make(map[string]string)
	for _, detector := range a.Detectors {
		detector.ComputeResults(a.App)
		a.Results[detector.GetTypeString()] = detector.GetResults()
	}
	return a
}

// ChdirToRepoRoot changes the working directory to the root of the repository (where go.mod is),
// since the pipeline resolves the app registry, configs and outputs relative to it
func ChdirToRepoRoot() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return os.Chdir(dir)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return os.ErrNotExist
		}
		dir = parent
	}
}
