# Tests

## Overview

Aletheia analyzes an app in four stages, in this order (see `internal/pipeline/pipeline.go`):

1. **SSA taint propagation:** marks values in the code with the database fields and RPCs they flow to.
2. **Abstract call graph:** links services, databases and the calls between them.
3. **Schema building:** infers constraints such as foreign keys.
4. **Detection:** searches for the five integrity-violation patterns.

The **integration tests** (`integration/`) run all four stages on the apps in `blueprint/examples/` and check the result of each stage. Check out the `blueprint` submodule before running them.

There is one folder per stage, each with one file per app (e.g., `ssa/digota_test.go`), a `helpers_test.go` with the shared helpers and a `main_test.go`:

- `ssa/` (stage 1): checks that values are marked with the right database fields and RPCs.
- `abstractcallgraph/` (stage 2): checks that the graph has the expected services, databases and calls, and that values flow between them as in the app's code.
- `detection/` (stages 3 and 4): checks that the expected constraints and violations are found. `detection/output_test.go` also checks that every app's results match `expected/`, and `detection/input_models_test.go` checks the results of the input models in `input-models/`, which must match the `expected/` output of the Blueprint app they describe, if any (it skips stage 1, so it also runs with `-short`).

The integration tests don't run the `aletheia` binary, because it only writes the final results to files. Instead, they call the **runner** (`runner/`), which runs the same pipeline as the binary (`pipeline.Run`) without writing to `output/` and keeps what each stage produced in memory (SSA graphs, abstract call graph, schema, detector results), so the tests can check them directly. The four stages run only once per app in each `go test` run, and all tests reuse that result.

The **unit tests** (`unit/`) check single functions on small hand-built inputs, without loading an app.

## Running

From the repository root:

```zsh
make test                # same as go test ./tests/...
make test-short          # same as go test -short ./tests/...
go test ./tests/...
go test -short ./tests/...
go test ./tests/integration/ssa          # a single stage
go test ./tests/... -run TestDigota      # a single app
```

Add `-count=1` after editing an app in `blueprint/examples/`, since Go's test cache doesn't notice those changes.

## Expected output

If a change is correct and meant to alter the expected results, regenerate the expected output and review the diff:

```zsh
go test ./tests/integration/detection -run TestDetectionOutput -update
git diff tests/expected
```
