# Aletheia: Automated Detection of Data Integrity Violations in Microservices

Aletheia is a static analysis framework for detecting data integrity violations in microservice-based applications written in [Blueprint](https://github.com/Blueprint-uServices/blueprint).

Aletheia operates in four steps:

1. Intra-procedural analysis based on SSA graphs
2. Inter-microservice analysis through construction of an abstract call graph
3. Schema extraction for objects stored across databases
4. Detection of code sections that violate integrity constraints, including entity integrity, referential integrity, and uniqueness

Integrity violations are detected by searching for the following patterns formalized in our paper:

| ID | Constraint Type | Violation Pattern |
| --- | --- | --- |
| RI-1 | Referential integrity | Absence of cascading effects |
| RI-2 | Referential integrity | Concurrent operations |
| RI-3 | Referential integrity | Uncoordinated replication |
| EI-1 | Entity integrity | Uncoordinated replication |
| Un-2 | Uniqueness | Conflicting writes |

## Overview

The `aletheia/pkg` directory contains the core packages that implement the functionality:

```
aletheia/pkg/
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

After analyzing an application, the output will be stored in `aletheia/output/{app}` according to the following structure:

```
aletheia/output/{app}/
├── app.json                            # High-level dependencies with services metadata (packages, fields, methods, etc.) and databases
├── schema.json                         # Extracted data schema
└── analysis/                           # Results for each analyzed pattern
    ├── foreign-key-cascade.txt         # RI-1 pattern
    ├── foreign-key-concurrency.txt     # RI-2 pattern
    ├── foreign-key-coordination.txt    # RI-3 pattern
    ├── primary-key-coordination.txt    # EI-1 pattern
    └── unicity-concurrency.txt         # Un-2 pattern
```

## Getting Started

After cloning the repository, initialize the Blueprint submodule:

```zsh
git submodule update --init --recursive
```

## Requirements

- [Golang](https://go.dev/doc/install) >= 1.24.5

## Supported Applications

By default, Aletheia supports analysis for the following applications:

 - digota
 - sockshop
 - eshopmicroservices
 - postnotification
 - dsb_socialnetwork
 - dsb_mediamicroservices
 - trainticket

Additional applications can be analyzed by extending the framework-specific extraction logic and configuration.

## Running Aletheia

Run Aletheia to analyze the application specified by the `app` parameter:

```zsh
cd aletheia
go run main.go {app}
```

Example:

```zsh
go run main.go sockshop
```

This will generate the analysis results under `aletheia/output/sockshop`
