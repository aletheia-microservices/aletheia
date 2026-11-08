# Codebase Overview

This page describes how the repository and the implementation of Aletheia are organized, and where each step of the analysis is implemented.

<div class="repo-group-nav">
  <span class="repo-group users">Using Aletheia</span>
  <span class="repo-group developers">Developing Aletheia</span>
</div>

## Repository Folders

Folders users work with when analyzing their own applications.

--8<-- "docs/snippets/user-folders.md"

Folders developers need to understand or change how Aletheia works.

<div class="grid cards compact developers" markdown>

-   :material-cogs: __`internal/`__

    Analysis packages, described [below](#internal).

-   :material-console: __`cmd/aletheia/`__

    CLI entry point, which handles flags, printing of results, and evaluation metrics.

-   :material-test-tube: __`tests/`__

    Unit and integration tests.

-   :material-book-open-variant: __`docs/`__

    This site, built with MkDocs (`make docs-serve`).

</div>

## Application Registry { #registry }

<div class="grid cards compact users" markdown>

-   :material-clipboard-list-outline: __`registry/*.yaml`__

    Tell Aletheia where to find each application and its wiring specs.

-   :material-cog-sync-outline: __`scripts/gen-registry/`__

    Reads `registry/*.yaml`, writes the Blueprint wiring spec imports to `internal/frameworks/blueprint/apps/apps.go`, and adds `replace` entries to `go.mod` for the local application folders in `blueprint/examples/`. Run it with `make registry`.

</div>

## Analysis Packages { #internal }

The analysis is implemented in `internal/`. The packages that implement each step of the pipeline are marked with its number (e.g., <span class="step-badge" title="Step 1 of the pipeline">1</span>), and the markers next to some packages open more details.

```mermaid
--8<-- "docs/snippets/pipeline.mmd"
    class ssa,acg,schema,detect step
```

<div class="file-tree annotate developers" markdown>

- :material-folder-outline: `internal/`
    - :material-folder-outline: `pipeline/` Runs all steps for one app (1)
    - :material-folder-outline: `analysis/`
        - :material-folder-outline: `common/` Shared database operation types
        - :material-folder-outline: `service-level/` Intra-service analysis
            - :material-folder-outline: `ssagraph/` <span class="step-badge" title="Step 1 of the pipeline">1</span> SSA graphs and taints (2)
        - :material-folder-outline: `system-level/` Inter-service analysis
            - :material-folder-outline: `abstractgraph/` <span class="step-badge" title="Step 2 of the pipeline">2</span> <span class="step-badge" title="Step 3 of the pipeline">3</span> Abstract call graph and foreign keys (3)
            - :material-folder-outline: `detection/` <span class="step-badge" title="Step 4 of the pipeline">4</span> One detector per pattern (4)
    - :material-folder-outline: `app/` Services, databases, and schemas (5)
    - :material-folder-outline: `config/` Global analysis settings (`config.Global`)
    - :material-folder-outline: `frameworks/` Framework parsing, e.g., Blueprint
        - :material-folder-outline: `blueprint/`
        - :material-folder-outline: `components/`
    - :material-folder-outline: `utils/` Program loading, field paths, timestamps

</div>

1.  Used by both the CLI (`cmd/aletheia/`) and the tests.
2.  Implements the [intra-procedural analysis](intra-procedural-analysis.md), which builds the SSA graph of each function and propagates database and service taints within each service.
3.  Implements the [inter-microservice analysis](inter-microservice-analysis.md), which builds and taints the abstract call graph. Its `tainter/` also infers the foreign keys of [schema extraction](schema-extraction.md), and `input/` parses [abstract models](../app-registry/abstract-model.md#guidelines-for-writing-abstract-models).
4.  Implements [pattern detection](detection-of-problematic-patterns.md), with one package per pattern in `constraints/` (`foreignkeycascade` for RI-1, `foreignkeyconcurrency` for RI-2, `keycoordination` for RI-3 and EI-1, and `uniquenessconcurrency` for Un-1).
5.  Represents the analyzed application, with its services, databases, schemas, and constraints, and saves them to `app.json`, `schema.json`, and `constraints.txt` in [`output/{app}/`](#output).

## Analysis Output { #output }

For each application, Aletheia saves its results in `output/{app}/`:

<div class="grid cards compact wide users" markdown>

-   :material-database-outline: __Application and schema__

    - `app.json`: services (packages, fields, methods) and databases
    - `schema.json`: extracted data schema
    - `constraints.txt`: primary keys, unique fields, and foreign keys, one per line

    See [Aletheia Internals → Schema Extraction](schema-extraction.md).

-   :material-alert-outline: __Detection warnings__ in `analysis/`

    - `foreign-key-cascade.txt`: RI-1
    - `foreign-key-concurrency.txt`: RI-2
    - `foreign-key-coordination.txt`: RI-3
    - `primary-key-coordination.txt`: EI-1
    - `uniqueness-concurrency.txt`: Un-1

    See [Usage → Read Detection Warnings](../usage/reading-warnings.md) for how to read them.

-   :material-sitemap-outline: __Diagrams__ in `diagrams/`

    App and schema diagrams, only if generated.

    See [Usage → Generate Diagrams](../usage/diagrams.md).

-   :material-bug-outline: __Intermediate results__

    - `ssa/`: SSA code of each package
    - `ssagraphs/`: tainted and untainted SSA graphs of each function, only with `--debug`
    - `abstractcallgraph/`: abstract call graph, only with `--debug`

    The graphs are saved as `.dot` files.

    See [Getting Started → Running Aletheia](../getting-started/running.md).

</div>
