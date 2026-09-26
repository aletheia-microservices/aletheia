# Tests

## Overview

Aletheia analyzes an app in four stages, in this order (see `internal/pipeline/pipeline.go`):

1. **SSA taint propagation:** marks values in the code with the database fields and RPCs they flow to.
2. **Abstract call graph:** links services, databases and the calls between them.
3. **Schema building:** infers constraints such as foreign keys.
4. **Detection:** searches for the five integrity-violation patterns.

The **integration tests** (`integration/`) run all four stages on the apps in `blueprint/examples/` and check the result of each stage. Check out the `blueprint` submodule before running them.

- `ssa_taint_test.go` (stage 1): checks that values are marked with the right database fields and RPCs.
- `abstractcallgraph/` (stage 2): one file per app. Checks that the graph has the expected services, databases and calls, and that values flow between them as in the app's code.
- `detection_test.go` (stages 3 and 4): checks that the expected constraints and violations are found, and that every app's results match `expected/`.

The integration tests don't run the `aletheia` binary, because it only writes the final results to files. Instead, they call the **runner** (`runner/`), which runs the same pipeline as the binary (`pipeline.Run`) without writing to `output/` and keeps what each stage produced in memory (SSA graphs, abstract call graph, schema, detector results), so the tests can check them directly. The four stages run only once per app in each `go test` run, and all tests reuse that result.

The **unit tests** (`unit/`) check single functions on small hand-built inputs, without loading an app. There is one folder per package: `unit/abstractgraph`, `unit/detection`, `unit/ssagraph` and `unit/utils` test `internal/analysis/system-level/abstractgraph`, `internal/analysis/system-level/detection`, `internal/analysis/service-level/ssagraph` and `internal/utils`. For example, they check that the timestamp `t4.t14` sorts before `t5`, and that the same taint is not added to a value twice.

## Running

From the repository root:

```zsh
make test                # same as go test ./tests/...
make test-short          # same as go test -short ./tests/...
go test ./tests/...
go test -short ./tests/...
go test ./tests/... -run TestSSA
```

Add `-count=1` after editing an app in `blueprint/examples/`, since Go's test cache doesn't notice those changes.

## Expected output

If a change is correct and meant to alter the expected results, regenerate the expected output and review the diff:

```zsh
go test ./tests/integration -run TestDetectionOutput -update
git diff tests/expected
```
