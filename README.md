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

### Running Aletheia

Aletheia analyzes applications located in `blueprint/examples/`. Some examples include:

- digota
- sockshop
- eshopmicroservices
- postnotification
- dsb_socialnetwork
- dsb_mediamicroservices
- trainticket

To analyze an application, run Aletheia and specify the application name as the `app` parameter:

```zsh
go run main.go {app}
```

Example:

```zsh
go run main.go postnotification
```

The warnings related to integrity violations are saved in `output/postnotification/analysis/`. The information about the application dependencies (microservices and datastores used) and the schema are saved in `output/postnotification/app.json` and `output/postnotification/schema.json`, respectively. The application's SSA code are saved in `output/postnotification/ssa/`.

You can also specify the `--debug` flag to obtain tainted *ssa graphs* and *abstract call graph* in `.dot` format saved under `output/postnotification/abstractcallgraph` and `output/postnotification/ssagraphs`, which can then be visualized in, for example, [Graphivz](https://dreampuf.github.io/GraphvizOnline/).

```zsh
go run main.go --debug postnotification
```

You can also specify which warnings should be suppressed by passing the `--detection_config` flag followed by the file path:

```zsh
go run main.go --detection_config config/sockshop.yaml sockshop
```

### Registering new Applications

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

### Registering and Analyzing a Simple Application

We now demonstrate how to run Aletheia to analyze a simple application (`simpleshop`) provided in `blueprint/examples/simpleshop/`. The application is composed of two microservices, Product Service and Inventory Service and allows clients to register new products and their respective inventory, as well as delete products.

Add a new entry for the `simpleshop` application in the `aletheia/registry/apps.yaml`. This will tell Aletheia how to properly import and analyze the application:

```yaml
- name: simpleshop
  app_root: github.com/blueprint-uservices/blueprint/examples/simpleshop
  package_path: simpleshop/workflow/simpleshop
  spec_name: simpleshop_docker
  spec_path: github.com/blueprint-uservices/blueprint/examples/simpleshop/wiring/specs
```

Go to `aletheia` directory:

```zsh
cd aletheia
```

Generate the application registry:

```zsh
go run scripts/gen_app_registry/main.go
```

Now, you can run the analysis:

```zsh
go run main.go simpleshop
```

This command prints the analysis results and saves them in `aletheia/output/simpleshop/`. The warnings related to integrity violations are saved in `aletheia/output/simpleshop/analysis/`. Information about the application dependencies (microservices and datastores used) and the schema are saved in `aletheia/output/simpleshop/app.json` and `aletheia/output/simpleshop/schema.json`, respectively.

The output should contain a referential integrity warning indicating a missing cascading delete. In this case, when a product is deleted, the effect is not propagated to the Inventory Service, leaving a dangling inventory record.

```txt
[NUM_WARNINGS = 1]
delete: ProductService.DeleteProduct() ... product_db.product.DeleteOne()
	missing cascade #1: database={inventory_db}, entity={inventory}, pending_fields={ID}
```
