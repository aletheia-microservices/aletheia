// Package integration runs the full Aletheia pipeline on the Blueprint apps in blueprint/examples
// and checks the SSA taints, the abstract call graph, the inferred schema and the detected
// integrity violations
package integration

import (
	"fmt"
	"os"
	"testing"

	"github.com/aletheia-microservices/aletheia/tests/runner"
)

func TestMain(m *testing.M) {
	// the pipeline resolves apps, configs and expected output relative to the repository root
	if err := runner.ChdirToRepoRoot(); err != nil {
		fmt.Fprintf(os.Stderr, "could not find repository root: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}
