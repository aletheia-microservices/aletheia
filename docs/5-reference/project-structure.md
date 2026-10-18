# Project Structure

```
aletheia/
├── cmd/aletheia/          # CLI entry point (flags, printing results, saving eval metrics)
├── internal/              # Packages that implement Aletheia (see below)
├── scripts/
│   ├── app-diagram/       # Draws the service dependency graph of an app (see scripts/app-diagram/README.md)
│   ├── gen-registry/      # Generates the app registry from registry/*.yaml (see scripts/gen-registry/README.md)
│   ├── schema-diagram/    # Draws the tables, fields and foreign keys of an app (see scripts/schema-diagram/README.md)
│   └── verify/            # Compares warning counts of every app against a baseline (see scripts/verify/README.md)
|
├── registry/              # Registered applications (apps.yaml)
├── config/                # Per-application detection configs that suppress warnings
├── input-models/          # Input models that describe applications without Blueprint, one folder per app (see input-models/README.md)
├── tests/                 # Unit and integration tests (see tests/README.md)
├── docs/                  # Documentation pages, built into a site with MkDocs (make docs-serve)
│   ├── 1-pipeline/        # Intra-procedural analysis, inter-microservice analysis, schema extraction, and detection
│   ├── 2-getting-started/ # Installation, running Aletheia, and the simpleshop tutorial
│   ├── 3-usage/           # Suppressing warnings, registering applications, and generating diagrams
│   ├── 4-formats/         # Writing input models and reading detection warnings
│   └── 5-reference/       # Project structure and conventions
├── overrides/             # MkDocs theme customizations (color schemes, text size buttons)
├── blueprint/             # Blueprint framework and example applications (git submodule)
├── mkdocs.yml             # MkDocs configuration (navigation, theme, plugins)
├── pyproject.toml         # Python dependencies for building the docs, installed with uv
└── Makefile
```

The `internal/` directory contains **the main implementation of Aletheia** and is organized as follows:

```
internal/
├── pipeline/                       # Runs every stage of the analysis for one app (used by the CLI and the tests)
├── analysis/
│   ├── common/                     # Database operation types shared across packages
│   ├── service-level/              # Implementation for intra-service analysis
│   │   └── ssagraph/               # SSA graph construction and taint propagation within each service
|   |
│   └── system-level/               # Implementation for inter-service analysis
│       ├── abstractgraph/          # Abstract call graph construction and analysis across services
│       └── detection/              # Detection for each pattern (RI-1, RI-2, RI-3, EI-1, Un-1)
|
├── app/                            # Application metadata with services, databases, schemas, and constraints
├── config/                         # Global analysis settings (config.Global)
├── frameworks/                     # Framework-specific parsing code (e.g., wiring specs for Blueprint apps)
│   ├── blueprint/
│   └── components/
└── utils/                          # Helpers for loading programs, parsing function and field paths, and comparing timestamps
```

The `config/` folder contains per-application detection config files that suppress warnings (see [Suppressing Detection Warnings](../3-usage/suppressing-warnings.md)). Not to be confused with `internal/config/`, which holds global analysis settings.

The `registry/` folder contains YAML files needed by Aletheia to properly import and analyze applications.

The `scripts/gen-registry/` generator (run with `make registry`) reads `registry/*.yaml` and: (i) writes `internal/frameworks/blueprint/apps/apps.go`, which imports each application's Blueprint wiring spec, and (ii) adds `replace` entries to `go.mod` that point each application to its folder in `blueprint/examples/`.

After analyzing an application, the output will be stored in `output/{app}` according to the following structure:

```
output/{app}/
├── app.json                            # High-level dependencies with services metadata (packages, fields, methods, etc.) and databases
├── schema.json                         # Extracted data schema
├── analysis/                           # Results for each analyzed pattern
│   ├── foreign-key-cascade.txt         # RI-1 pattern
│   ├── foreign-key-concurrency.txt     # RI-2 pattern
│   ├── foreign-key-coordination.txt    # RI-3 pattern
│   ├── primary-key-coordination.txt    # EI-1 pattern
│   └── uniqueness-concurrency.txt      # Un-1 pattern
└── diagrams/                           # App and schema diagrams, only if generated (see Generating Diagrams)
```
