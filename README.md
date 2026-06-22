# Aletheia

**Aletheia: Automated Detection of Data Integrity Violations in Microservices**  
_Mafalda Sofia Ferreira, João Ferreira Loff, João Garcia, and Rodrigo Rodrigues_  
_INESC-ID, Instituto Superior Técnico, Universidade de Lisboa_  
_In Proceedings of the 20th USENIX Symposium on Operating Systems Design and Implementation (**OSDI ’26**)_

For instructions on how to reproduce the experiments from the paper, see the [Aletheia Artifact OSDI'26](https://github.com/aletheia-microservices/aletheia-artifact-osdi26) repository.

---

## Overview

Aletheia is a static analysis framework for detecting data integrity violations in microservice-based applications.

In microservice architectures, data is stored across heterogeneous systems, with data schemas partitioned and managed by separate services. Due to the complexity of microservices, it can be almost impossible for developers to have a comprehensive understanding of the entire system, making it challenging to reason about and maintain data integrity at the application level.

Aletheia solves this problem through static analysis by identifying semantic violations in microservice ecosystems (i.e., service interactions and operations that break data integrity) for various types of integrity constraints:

- **Entity integrity constraints** (defined through primary keys)
- **Referential integrity constraints** (defined through foreign keys)
- **Uniqueness constraints**

Aletheia analyzes applications targeting the [Blueprint](https://github.com/Blueprint-uServices/blueprint) compiler.

The code of all applications is located in `blueprint/examples/`, which includes a [README](https://github.com/aletheia-microservices/blueprint/blob/osdi26/examples/README.md) summarizing the source repositories and versions used for each application.

The framework operates in four steps:

1. **Intra-procedural analysis** based on Static Single-Assignment (SSA) graphs extracted from Go code to infer how values flow throughout execution by propagating taint information
2. **Inter-microservice analysis** based on a new _abstract call graph_ that represents possible call graphs containing microservice invocations and database operations, along with filtered taint information from the SSA analysis
3. **Schema extraction** for objects stored across databases
4. **Detection of problematic code sections** that violate integrity constraints, including entity integrity (primary keys), referential integrity (foreign keys), and uniqueness

### Detection of Code Patterns

Integrity violations are detected by searching for the following patterns, formalized in our paper and implemented in [`pkg/detection/constraints/`](./pkg/detection/constraints/):

| ID   | Constraint Type       | Violation Pattern            | Implementation Package                            |
| ---- | --------------------- | ---------------------------- | ------------------------------------------------- |
| RI-1 | Referential integrity | Absence of cascading effects | `pkg/detection/constraints/foreignkeycascade`     |
| RI-2 | Referential integrity | Concurrent operations        | `pkg/detection/constraints/foreignkeyconcurrency` |
| RI-3 | Referential integrity | Uncoordinated replication    | `pkg/detection/constraints/keycoordination`       |
| EI-1 | Entity integrity      | Uncoordinated replication    | `pkg/detection/constraints/keycoordination`       |
| Un-1 | Uniqueness            | Conflicting writes           | `pkg/detection/constraints/uniquenessconcurrency` |

> [!NOTE]
> Refer to our paper for the formal definitions of these patterns.

### Cross-Microservice Foreign Key Inference

Data associations across microservices are inferred from taints propagated through related objects used in database operations in the _abstract call graph_. The inference is performed according to the rules implemented in [`pkg/abstractgraph/tainter.go`](./pkg/abstractgraph/tainter.go).

Rules are applied for each pair of operation (`op_1`, `op_2`) where `op_i` is either a `read` or a `write`. We use `field_1` and `field_2` to denote fields accessed (tainted) by the same object in `op_1` and `op_2`, respectively.

| Operation Pair   | Foreign Key Direction        |
| ---------------- | ---------------------------- |
| `(write, write)` | `field2` references `field1` |
| `(read, write)`  | `field2` references `field1` |
| `(write, read)`  | `field1` references `field2` |

In read operations, the `read_key` denotes that the propagated object is used as a filter in the read operation, while `read_val` denotes that the propagated object is returned from the read operation.

| Operation Pair         | Foreign Key Direction        |
| ---------------------- | ---------------------------- |
| `(read_key, read_key)` | `field2` references `field1` |
| `(read_val, read_key)` | `field1` references `field2` |

> [!NOTE]
> Refer to our paper for a detailed explanation of these inference rules.

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
- `sql_tables` (optional): primary keys and uniqueness constraints for SQL databases
- `nosql_path` (optional): indexes constraints for NoSQL databases, which are then treated as uniqueness constraints

> [!NOTE]
> Go struct fields annotated with the `_id` BSON tag for NoSQL databases such as MongoDB are automatically treated as primary keys in the global application schema, since these fields are typically indexed by default:
>
> ```go
> type User struct {
>     ID string `bson:"_id"`
> }
> ```

Then, generate the application registry:

```zsh
go run scripts/gen_app_registry/main.go
```

### Registering and Analyzing a Simple Application

We now demonstrate how to run Aletheia to analyze a simple application (`simpleshop`) provided in `blueprint/examples/simpleshop/`. The application is composed of two microservices, Product Service and Inventory Service, and allows clients to register new products and their respective inventory, as well as delete products.

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

Then, you can run the analysis again and pass the `--detection_config` flag followed by the file path:

```zsh
go run main.go --detection_config config/simpleshop.yaml simpleshop
```
