# Aletheia

Aletheia is a static analysis framework for detecting data integrity violations in microservice-based applications.

In microservice architectures, data is stored across heterogeneous systems, with data schemas partitioned and managed by separate services. Due to the complexity of microservices, it can be almost impossible for developers to have a comprehensive understanding of the entire system, making it challenging to reason about and maintain data integrity at the application level. 

Aletheia solves this problem through static analysis, identifying semantic violations in microservice ecosystems (i.e., service interactions and operations that break data integrity) for various types of integrity constraints, including entity integrity, referential integrity, and uniqueness.

Currently, Aletheia analyzes applications developed with [Blueprint](https://github.com/Blueprint-uServices/blueprint) and located in `blueprint/examples/`.

The framework operates in four steps:

1. **Intra-procedural analysis** based on Static Single-Assignment (SSA) graphs extracted from Go code to infer how values flow throughout execution by propagating taint information
2. **Inter-microservice analysis** based on a new _abstract call graph_ that represents possible call graphs containing microservice invocations and database operations, along with filtered taint information from the SSA analysis
3. **Schema extraction** for objects stored across databases
4. **Detection of code sections** that violate integrity constraints, including entity integrity, referential integrity, and uniqueness

Integrity violations are detected by searching for the following patterns formalized in our paper and implemented by Aletheia in `pkg/detection/constraints/`.

| ID   | Constraint Type       | Violation Pattern            | Implementation Package                            |
| ---- | --------------------- | ---------------------------- | ------------------------------------------------- |
| RI-1 | Referential integrity | Absence of cascading effects | `pkg/detection/constraints/foreignkeycascade`     |
| RI-2 | Referential integrity | Concurrent operations        | `pkg/detection/constraints/foreignkeyconcurrency` |
| RI-3 | Referential integrity | Uncoordinated replication    | `pkg/detection/constraints/keycoordination`       |
| EI-1 | Entity integrity      | Uncoordinated replication    | `pkg/detection/constraints/keycoordination`       |
| Un-1 | Uniqueness            | Conflicting writes           | `pkg/detection/constraints/uniquenessconcurrency` |

## Project Structure

The `pkg/` directory contains the packages that implement Aletheia and is organized as follows:

```
pkg/
├── abstractgraph/                  # Abstract call graph construction and analysis
├── app/                            # Application metadata with services, databases, schemas, and constraints
├── common/
├── config/
├── detection/                      # Detection for each pattern (RI-1, RI-2, RI-3, EI-1, Un-1)
├── frameworks/                     # Framework-specific parsing code (e.g., wiring specs for Blueprint apps)
│   ├── blueprint/
│   └── components/
├── ssagraph/                       # SSA graph construction and analysis
└── utils/
```

The `config/` folder contains YAML files that specify warnings to be ignored by Aletheia.

The `registry/` folder contains YAML files needed by Aletheia to properly import and analyze applications.

The `scripts/gen_app_registry/` folder contains a script that uses the YAML files in `registry/` to: (i) generate a Go file under `pkg/frameworks/blueprint/` defining how Aletheia locates applications and imports their corresponding Blueprint specs, and (ii) update `go.mod` with new entries so that Go can locate applications relative to Aletheia's path.

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
│   └── uniqueness-concurrency.txt      # Un-1 pattern
```

## Requirements

- [Golang](https://go.dev/doc/install) >= 1.24.5

## Getting Started

Clone the repository:

```zsh
git clone --recurse-submodules https://github.com/aletheia-microservices/aletheia.git
```

If you already cloned the repository without `--recurse-submodules`, make sure to initialize the submodules:

```zsh
git submodule update --init --recursive
```

### Running Aletheia

Aletheia analyzes applications located in `blueprint/examples/`. Some examples include:

- digota
- sockshop
- dsb_mediamicroservices
- dsb_socialnetwork
- eshopmicroservices
- postnotification
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

You can also specify the `--debug` flag to obtain tainted _ssa graphs_ and _abstract call graph_ in `.dot` format saved under `output/postnotification/abstractcallgraph` and `output/postnotification/ssagraphs`, which can then be visualized in, for example, [Graphivz](https://dreampuf.github.io/GraphvizOnline/).

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

See [assumptions.md](./assumptions.md) for the current analysis assumptions and limitations.

First, you will need to add a new application entry to `registry/apps.yaml`. You can use the existing entries as examples. The new entry should contain the following values:

- `name`: application name
- `package_path`: package path for application code
- `spec_name`: blueprint spec name with format `{app_name}_{spec_name}` (e.g., `foobar_docker` => spec `Docker` for application `foobar`)
- `spec_path`: package path for spec
- `sql_tables` (optional): primary keys and uniqueness constraints for sql databases
- `nosql_path` (optional): indexes constraints for nosql databases, which are then treated as uniqueness constraints

Then, generate the application registry:

```zsh
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

Now, you will need to generate the application registry according to the new entry added to `aletheia/registry/apps.yaml`. The following script will (i) generate a Go file under `pkg/frameworks/blueprint/` defining how Aletheia locates applications and imports their corresponding Blueprint specs, and (ii) update `go.mod` with new entries so that Go can locate applications relative to Aletheia's path.

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

If you want to suppress all warnings related to missing cascade deletes on the inventory, create a new YAML file at `aletheia/config/simpleshop.yaml` with the following content:

```yaml
app: simpleshop
ignore_cascade:
  - database: inventory_db
    entity: inventory
    # optional fields for more fine-grained control
    trigger_database: product_db
    trigger_entity: product

```

Then, you can run again the analysis and pass the `--detection_config` flag followed by the file path:

```zsh
go run main.go --detection_config config/simpleshop.yaml simpleshop
```
