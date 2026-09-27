// Package detection checks the constraints inferred (stage 3) and the integrity violations detected (stage 4)
// for each app in blueprint/examples, with one file per app
package detection

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
