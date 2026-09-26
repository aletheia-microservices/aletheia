// Package runner runs Aletheia's pipeline for tests, keeping the result of every stage (SSA
// graphs, abstract call graph, schema and detector results) in memory for inspection
package runner

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/aletheia-microservices/aletheia/internal/pipeline"
)

// Analysis holds the state of every stage of the pipeline for one app so that
// tests can inspect intermediate results (unlike the CLI, SSA graphs are not released)
type Analysis = pipeline.Result

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

// runAnalysis runs the same pipeline as the CLI without writing to output/
func runAnalysis(t *testing.T, appname string, configPath string) *Analysis {
	t.Helper()
	logrus.SetLevel(logrus.ErrorLevel)
	a, err := pipeline.Run(pipeline.Options{
		App:             appname,
		DetectionConfig: configPath,
		KeepGraphs:      true,
	})
	if err != nil {
		t.Fatalf("running pipeline for %s: %v", appname, err)
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
