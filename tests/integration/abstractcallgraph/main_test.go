// Package abstractcallgraph checks the abstract call graph built for each app in blueprint/examples,
// with one file per app
package abstractcallgraph

import (
	"fmt"
	"os"
	"testing"

	"analyzer/tests/runner"
)

func TestMain(m *testing.M) {
	// the pipeline resolves apps and configs relative to the repository root
	if err := runner.ChdirToRepoRoot(); err != nil {
		fmt.Fprintf(os.Stderr, "could not find repository root: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}
