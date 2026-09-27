package detection

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/aletheia-microservices/aletheia/tests/runner"
)

// regenerate the expected output with: go test ./tests/integration/detection -run TestDetectionOutput -update
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
var nondeterministicConstraints = map[string]string{}

var detectorTypes = []string{
	"foreign-key-cascade",
	"foreign-key-concurrency",
	"foreign-key-coordination",
	"primary-key-coordination",
	"uniqueness-concurrency",
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
		checkExpectedOutput(t, filepath.Join(expectedDir, "constraints.txt"), a.App.ConstraintsString())
	})
}

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
	constraints := strings.Split(a.App.ConstraintsString(), "\n")
	if got := slices.Contains(constraints, constraint); got != want {
		t.Errorf("constraint %q inferred = %v, want %v\nconstraints:\n%s", constraint, got, want, a.App.ConstraintsString())
	}
}
