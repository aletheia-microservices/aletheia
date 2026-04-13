# Aletheia: Automated Detection of Data Integrity Violations in Microservices

Aletheia is a static analysis framework for detecting data integrity violations in microservice-based applications written in [Blueprint](https://github.com/Blueprint-uServices/blueprint).

Aletheia operates in four steps:

1. Intra-procedural analysis based on SSA graphs
2. Inter-microservice analysis through construction of an abstract call graph
3. Schema extraction for objects stored across databases
4. Detection of code sections that violate integrity constraints, including entity integrity, referential integrity, and uniqueness

Integrity violations are detected by searching for the following patterns formalized in our paper:

| ID   | Constraint Type       | Violation Pattern            |
| ---- | --------------------- | ---------------------------- |
| RI-1 | Referential integrity | Absence of cascading effects |
| RI-2 | Referential integrity | Concurrent operations        |
| RI-3 | Referential integrity | Uncoordinated replication    |
| EI-1 | Entity integrity      | Uncoordinated replication    |
| Un-2 | Uniqueness            | Conflicting writes           |

## Overview

The `pkg` directory contains the core packages that implement the functionality:

```
pkg/
├── abstractgraph/                  # Abstract call graph construction and analysis
├── app/                            # Application model with services and databases
├── common/
├── config/
├── detection/                      # Data integrity violation detection for each pattern
│   └── constraints
│       ├── foreignkeycascade/      # RI-1 pattern
│       ├── foreignkeyconcurrency/  # RI-2 pattern
│       ├── keycoordination/        # RI-3 and EI-1 patterns
│       └── unicityconcurrency/     # Un-2 pattern
├── frameworks/                     # Framework-specific analysis code
│   ├── blueprint/
│   └── components/
├── ssagraph/                       # SSA graph construction and analysis
└── utils/
```

After analyzing an application, the output will be stored in `output/{app}` according to the following structure:

```
output/{app}/
├── app.json                            # High-level dependencies with services metadata (packages, fields, methods, etc.) and databases
├── schema.json                         # Extracted data schema
└── analysis/                           # Results for each analyzed pattern
│   ├── foreign-key-cascade.txt         # RI-1 pattern
│   ├── foreign-key-concurrency.txt     # RI-2 pattern
│   ├── foreign-key-coordination.txt    # RI-3 pattern
│   ├── primary-key-coordination.txt    # EI-1 pattern
│   └── unicity-concurrency.txt         # Un-2 pattern
```

The `config/` folder contains YAML files that specify warnings to be ignored by Aletheia.

The `registry/` folder contains YAML files needed by Aletheia to properly import and analyze applications.

The `scripts/gen_app_registry/` folder contains a script that uses the YAML files in `registry/` to: (i) generate a Go file under `pkg/frameworks/blueprint/` defining how Aletheia locates applications and imports their corresponding Blueprint specs, and (ii) update `go.mod` with new entries so that Go can locate applications relative to Aletheia's path.

## Requirements

- [Golang](https://go.dev/doc/install) >= 1.24.5

## Getting Started

After cloning the repository, initialize the Blueprint submodule:

```zsh
git submodule update --init --recursive
```

### Registering new Applications

By default, Aletheia supports analysis for the following applications inside `blueprint/examples/`:

- digota
- sockshop
- eshopmicroservices
- postnotification
- dsb_socialnetwork
- dsb_mediamicroservices
- trainticket

If you want to analyze your own application written in Blueprint, make sure it is placed in `blueprint/examples/`. The expected structure is:

```
blueprint/examples/{app}/
├── wiring/             # blueprint specification
├── workflow/           # blueprint workflow
│   └── {app}/          # microservices code
```

Then, add a new application entry to `registry/apps.yaml`. You can use the existing entries as examples. The new entry should contain the following values:
- `name`: application name
- `package_path`: package path for application code
- `spec_name`: blueprint spec name with format `{app_name}_{spec_name}` (e.g., `foobar_docker` => spec `Docker` for application `foobar`)
- `spec_path`: package path for spec
- `sql_tables` (optional): primary keys and uniqueness constraints for sql databases
- `nosql_path` (optional): indexes constraints for nosql databases

Generate the application registry:

```zsh
# make sure you run from this directory so paths for the new Go file and the Go mod file are resolved correctly
go run scripts/gen_app_registry/main.go
```

### Running Aletheia

Run Aletheia to analyze the application specified by the `app` parameter:

```zsh
go run main.go [--detection_config <filepath.yaml>] {app}
```

Examples:

```zsh
go run main.go postnotification
go run main.go --detection_config config/sockshop.yaml sockshop
```

The warnings related to integrity violations will be saved in `aletheia/output/{app}/analysis/`.

The information about the application dependencies (microservices and datastores used) and the schema are saved in `aletheia/output/{app}/app.json` and `aletheia/output/{app}/schema.json`, respectively.
