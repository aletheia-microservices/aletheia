# Aletheia

**Aletheia: Automated Detection of Data Integrity Violations in Microservices** [[USENIX Page]](https://www.usenix.org/conference/osdi26/presentation/ferreira) [[Paper]](https://www.usenix.org/system/files/osdi26-ferreira.pdf)<br>
Mafalda Sofia Ferreira, João Ferreira Loff, João Garcia, and Rodrigo Rodrigues  
INESC-ID, Instituto Superior Técnico, Universidade de Lisboa  
*In Proceedings of the 20th USENIX Symposium on Operating Systems Design and Implementation (**OSDI ’26**)*

For instructions on how to reproduce the experiments from the paper, see the [Aletheia Artifact OSDI'26](https://github.com/aletheia-microservices/aletheia-artifact-osdi26) repository.

---

## Overview

In microservice architectures, data is stored across heterogeneous systems, with data schemas partitioned and managed by separate services. Due to the complexity of microservices, it can be almost impossible for developers to have a comprehensive understanding of the entire system, making it challenging to reason about and maintain data integrity at the application level.

Aletheia solves this problem through static analysis by identifying semantic violations in microservice ecosystems (i.e., service interactions and operations that break data integrity) for various types of integrity constraints:

- **Entity integrity constraints** (defined through **primary keys**)
- **Referential integrity constraints** (defined through **foreign keys**)
- **Uniqueness constraints**

Aletheia analyzes applications targeting the [Blueprint](https://github.com/Blueprint-uServices/blueprint) compiler.

The code of all applications is located in `blueprint/examples/`, which includes a [README](https://github.com/aletheia-microservices/blueprint/blob/osdi26/examples/README.md) summarizing the source repositories and versions used for each application.

The framework operates in four steps:

1. **Intra-procedural analysis** based on Static Single-Assignment (SSA) graphs extracted from Go code to infer how values flow throughout execution by propagating taint information
2. **Inter-microservice analysis** based on a new _abstract call graph_ that represents possible call graphs containing microservice invocations and database operations, along with filtered taint information from the SSA analysis
3. **Schema extraction** for objects stored across databases
4. **Detection of problematic code sections** that violate integrity constraints, including entity integrity, referential integrity, and uniqueness

### Integrity Violation Patterns

Aletheia searches the code for five patterns of operations that can break these constraints. Each pattern is formalized in Section 3 of our paper and illustrated in its Figure 2.

- **RI-1: Referential integrity, absence of cascading deletes**

  A request deletes a record in one database but not the records in other databases that reference it, so they point to a record that no longer exists.

  _Example (`simpleshop`):_ `ProductService.DeleteProduct()` deletes a product from `product_db` but leaves its record in `inventory_db`.

- **RI-2: Referential integrity, concurrent operations**

  One request reads a record and stores a reference to it in another database, while a concurrent request deletes that record.

  _Example (`eshopmicroservices`):_ `BasketService.StoreBasket()` saves product IDs in `basket_db` while `CatalogService.DeleteProduct()` can delete those products from `catalog_db` at the same time.

- **RI-3: Referential integrity, uncoordinated replication**

  A request writes a record to one database and a reference to it in another. Because the databases replicate independently, another request can see the reference before the record it points to.

  _Example (`postnotification`):_ `UploadService.UploadPost()` saves a post in `posts_db` and pushes a notification with its `PostID` to `notifications_queue`. `NotifyService` can read the notification before the post is visible.

- **EI-1: Entity integrity, uncoordinated replication**

  A request splits one entity across two databases under the same primary key. Another request that reads both by that key can find the entity in one database but not yet in the other.

  _Example (`dsb_mediamicroservices`):_ `APIService.RegisterMovie()` saves each movie in `movie_id_db` and `movie_info_db` under the same `MovieID`, and `APIService.ReadPage()` can find it in only one of them.

- **Un-1: Uniqueness, conflicting writes**

  A request writes a unique value to one database and related data to another. If the databases are replicated and two requests write the same unique value concurrently, the first database keeps only one write when replicas merge, but the second keeps both.

  _Example (`dsb_mediamicroservices`):_ two concurrent `APIService.RegisterMovie()` calls with the same `Title` can leave one movie in `movie_id_db`, where `Title` is unique, but two in `movie_info_db`.

## Project Structure

```
aletheia/
├── cmd/aletheia/          # CLI entry point (flags, printing results, saving eval metrics)
├── internal/              # Packages that implement Aletheia (see below)
├── scripts/
│   ├── gen-registry/      # Generates the app registry from registry/*.yaml
│   └── verify/            # Compares warning counts of every app against a baseline
├── registry/              # Registered applications (apps.yaml)
├── config/                # Per-application detection configs that suppress warnings
├── tests/                 # Unit and integration tests (see tests/README.md)
├── docs/                  # Technical details and assumptions
├── blueprint/             # Blueprint framework and example applications (git submodule)
└── Makefile
```

The `internal/` directory is organized as follows:

```
internal/
├── pipeline/                       # Runs every stage of the analysis for one app (used by the CLI and the tests)
├── analysis/
│   ├── common/                     # Database operation types shared across packages
│   ├── service-level/              # Implementation for intra-service analysis
│   │   └── ssagraph/               # SSA graph construction and taint propagation within each service
|   |
│   └── system-level/               # Implementation for intra-service analysis
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

The `config/` folder contains per-application detection config files that suppress warnings (see [Suppressing Detection Warnings](#suppressing-detection-warnings)). Not to be confused with `internal/config/`, which holds global analysis settings.

The `registry/` folder contains YAML files needed by Aletheia to properly import and analyze applications.

The `scripts/gen-registry/` generator (run with `make registry`) reads `registry/*.yaml` and: (i) writes `internal/frameworks/blueprint/apps/apps.go`, which imports each application's Blueprint wiring spec, and (ii) adds `replace` entries to `go.mod` that point each application to its folder in `blueprint/examples/`.

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

- [Golang](https://go.dev/doc/install) >= 1.26

## Getting Started

Clone the repository:

```zsh
git clone --recurse-submodules https://github.com/aletheia-microservices/aletheia.git
```

If you already cloned the repository without `--recurse-submodules`, make sure to initialize the submodules:

```zsh
git submodule update --init --recursive
```

Build Aletheia binary from the repository root:

```zsh
make build
```

Run `make help` to list the other targets (`test`, `verify`, `registry`, ...). 

Always run Aletheia from the repository root, since it loads the applications through the `replace` directives in `go.mod` and reads `registry/` and `config/` relative to the working directory.

The rest of this README uses `./bin/aletheia`, but either of these alternatives also works:

- `go run ./cmd/aletheia {app}`: skips the build step and always runs your latest changes.
- `make install`, then `aletheia {app}`: installs into `$GOBIN` (or `$(go env GOPATH)/bin`), which must be on your `PATH`. Rerun `make install` after changing the code or running `make registry`.

To get started, run one of the included apps as shown below, then follow the [simpleshop tutorial](#tutorial-analyzing-your-first-application-simpleshop), which walks through registering an app, running the analysis, reading the warning, and suppressing it.

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
./bin/aletheia {app}
```

Example:

```zsh
./bin/aletheia postnotification
```

The results are saved in `output/postnotification/`:

- `ssa/`: the application's SSA code
- `app.json`: application dependencies (microservices and datastores used)
- `schema.json`: inferred data schema
- `analysis/`: warnings related to integrity violations

You can also specify the `--debug` flag to obtain tainted _ssa graphs_ and _abstract call graph_ in `.dot` format saved under `output/postnotification/abstractcallgraph` and `output/postnotification/ssagraphs`, which can then be visualized in, for example, [Graphivz](https://dreampuf.github.io/GraphvizOnline/).

```zsh
./bin/aletheia --debug postnotification
```

### Reading the Output

Each file in `output/{app}/analysis/` starts with `[NUM_WARNINGS = N]`, followed by one block per warning. All files use the same notation:

- **Call paths** show which service method runs a database operation. For example, `StorageService.ReadPost() ... posts_db.post.FindOne()` means that `StorageService.ReadPost()` runs `FindOne()` on entity `post` in database `posts_db`. Calls in between are left out. Some paths start with a third element: the entry point of the request (e.g., `UploadService.UploadPost() ... StorageService.StorePost() ... posts_db.post.InsertOne()`).
- **Fields** are written as `database.entity.field`. For example, `posts_db.post.Creator.Username` is the `Username` of the post's `Creator`, and `posts_db.post.Mentions[*]` refers to every element of the list `Mentions`.
- **Foreign keys** are inferred by Aletheia. For example, `FOREIGN_KEY notifications_queue.notification.PostID REFERENCES posts_db.post.PostID [MANDATORY]` means that each notification's `PostID` points to a post in `posts_db`. Foreign keys can have two tags:
  - `[MANDATORY]`: every request that creates the reference also writes the record it points to. Here, `UploadService.UploadPost()` writes both the post and its notification.
  - `[T]`: transitive, derived from two other foreign keys (if X references Y and Y references Z, then X references Z).

Click a pattern below to see an example warning explained.

<details>
<summary><b>RI-1</b> <code>foreign-key-cascade.txt</code>: missing cascading deletes</summary>

```txt
delete: ProductService.Delete() ... products_db.products.DeleteOne()
    missing cascade #1: database={skus_db}, entity={skus}, pending_fields={Parent}
```

- `delete:` is the operation that deletes a record.
- `missing cascade #N` names a database and entity whose records can still reference the deleted record, through the fields in `pending_fields`.

</details>

<details>
<summary><b>RI-2</b> <code>foreign-key-concurrency.txt</code>: concurrent delete and write</summary>

```txt
delete: ProductService.Delete() ... products_db.products.DeleteOne()
    write #1: SkuService.New() ... SkuService.New() ... skus_db.skus.InsertOne()
        - database={skus_db}, entity={skus}, written_fields={Parent}
```

- `delete:` is the operation that deletes a record.
- `write #N` is an operation in another request that can store a reference to that record at the same time, through the fields in `written_fields`.

</details>

<details>
<summary><b>RI-3</b> <code>foreign-key-coordination.txt</code>: reference visible before the record</summary>

```txt
entry request: UploadService.UploadPost()
    FOREIGN KEY READS #1:
        READ (FOREIGN KEY): NotifyService.Run() ... notifications_queue.notification.Pop()
            - field: notifications_queue.notification.PostID
            - constraint: FOREIGN_KEY notifications_queue.notification.PostID REFERENCES posts_db.post.PostID [MANDATORY]
        READ (ORIGIN): StorageService.ReadPost() ... posts_db.post.FindOne()
            - field: posts_db.post.PostID
```

- `entry request` is the client request that leads to the two reads.
- `READ (FOREIGN KEY)` reads a record that holds a reference (`field`, following `constraint`).
- `READ (ORIGIN)` reads the record that the reference points to.

</details>

<details>
<summary><b>EI-1</b> <code>primary-key-coordination.txt</code>: entity split across databases</summary>

```txt
entry request: APIService.ReadPage()
    PRIMARY KEY READS #1:
        READ: MovieInfoService.ReadMovieInfo() ... movie_info_db.movie_info.FindOne()
            - field: movie_info_db.movie_info._id
            - constraint: PRIMARY KEY (movie_info_db.movie_info._id)
        READ: MovieIdService.ReadMovieId() ... movie_id_db.movie.FindOne()
            - field: movie_id_db.movie._id
            - constraint: PRIMARY KEY (movie_id_db.movie._id)
```

- `entry request` is the client request that leads to the two reads.
- Each `READ` reads one part of the entity from a different database, using the same primary key.

</details>

<details>
<summary><b>Un-1</b> <code>uniqueness-concurrency.txt</code>: conflicting writes of a unique value</summary>

```txt
entry request: APIService.RegisterMovie()
write (origin): APIService.RegisterMovie() ... MovieIdService.RegisterMovieId() ... movie_id_db.movie.InsertOne()
        - field (constrained): movie_id_db.movie.Title (UNIQUE)
    - affected write #1: MovieInfoService.WriteMovieInfo() ... movie_info_db.movie_info.InsertOne()
```

- `write (origin)` writes the unique field shown in `field (constrained)`.
- `affected write #N` is a related write to another database in the same request. Its effect remains even if the unique write is discarded.

To review a warning, follow its call path in the code. Fix real problems, and suppress false positives as described below.

</details>

#### Suppressing Detection Warnings

To suppress warnings, write a detection config file and pass its path with the `--detection_config` flag. A detection config file has the following format:

```yaml
app: <app>                        # must match the app name passed to Aletheia
ignore_foreignkeys:               # inferred foreign keys to ignore (all patterns)
  - <database>.<entity>.<field>   # the referencing field, as shown in schema.json
ignore_cascade:                   # missing cascading deletes to ignore (RI-1 only)
  - database: <database>          # database of the records left behind
    entity: <entity>
    trigger_database: <database>  # optional, together with trigger_entity: only when
    trigger_entity: <entity>      # the deleted record is in this database and entity
```

- `ignore_foreignkeys` stops Aletheia from inferring foreign keys from the listed fields, so those foreign keys do not appear in `schema.json` and do not lead to warnings.
- `ignore_cascade` hides RI-1 warnings about records left behind in the given `database` and `entity`. To hide them only when the delete happens in a specific database and entity, set both `trigger_database` and `trigger_entity`. If only one of them is set, it has no effect.

You can ignore specific inferred foreign keys. For example, `config/postnotification.yaml` ignores the foreign key on `notifications_queue.notification.ReqID`:

```zsh
./bin/aletheia --detection_config config/postnotification.yaml postnotification
```

Or you can suppress missing cascading delete warnings. For example, `config/sockshop.yaml` ignores missing cascades from `cart_db.carts` to `order_db.orders`:

```zsh
./bin/aletheia --detection_config config/sockshop.yaml sockshop
```

### Tutorial: Analyzing Your First Application (simpleshop)

We now demonstrate how to run Aletheia to analyze a simple application (`simpleshop`) provided in `blueprint/examples/simpleshop/`. The application is composed of two microservices, Product Service and Inventory Service, and allows clients to register new products and their respective inventory, as well as delete products.

Add a new entry for the `simpleshop` application at the end of the `apps` list in `registry/apps.yaml`. This will tell Aletheia how to properly import and analyze the application:

```yaml
apps:
  # ... existing entries ...

  - name: simpleshop
    app_root: github.com/blueprint-uservices/blueprint/examples/simpleshop
    package_path: simpleshop/workflow/simpleshop
    spec_name: simpleshop_docker
    spec_path: github.com/blueprint-uservices/blueprint/examples/simpleshop/wiring/specs
```

Now, you will need to generate the application registry according to the new entry added to `registry/apps.yaml`. The following command will (i) generate a Go file under `internal/frameworks/blueprint/apps/` defining how Aletheia locates applications and imports their corresponding Blueprint specs, (ii) update `go.mod` with new entries so that Go can locate applications relative to Aletheia's path, and (iii) rebuild `bin/aletheia` so that it includes the new application.

```zsh
make registry
```

Now, you can run the analysis:

```zsh
./bin/aletheia simpleshop
```

This command prints the analysis results and saves them in `output/simpleshop/`.

The output should contain a referential integrity warning indicating a missing cascading delete (pattern RI-1, see [Reading the Output](#reading-the-output)). In this case, when a product is deleted, the effect is not propagated to the Inventory Service, leaving a dangling inventory record.

```txt
[NUM_WARNINGS = 1]
delete: ProductService.DeleteProduct() ... product_db.product.DeleteOne()
	missing cascade #1: database={inventory_db}, entity={inventory}, pending_fields={ID}
```

If you want to suppress all warnings related to missing cascade deletes on the inventory, create a new YAML file at `config/simpleshop.yaml` with the following content:

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
./bin/aletheia --detection_config config/simpleshop.yaml simpleshop
```

The tutorial modifies files tracked by git. When you are done, you can undo these changes with:

```zsh
git restore registry/apps.yaml go.mod internal/frameworks/blueprint/apps/apps.go
make build
rm config/simpleshop.yaml
```

### Analyzing Your Own Application

Follow these steps to analyze your own application. The [simpleshop tutorial](#tutorial-analyzing-your-first-application-simpleshop) goes through steps 4 to 7 on a small example.

**1. Port the application to Blueprint.** Skip this step if your application is already written with Blueprint. A Blueprint application has two parts: the _workflow_, where each service is a Go interface and its implementation, and the _wiring_, a spec that creates the services and databases and connects them. See Blueprint's guides on [workflow](https://github.com/Blueprint-uServices/blueprint/blob/main/docs/manual/workflow.md) and [wiring](https://github.com/Blueprint-uServices/blueprint/blob/main/docs/manual/wiring.md), and use the applications in `blueprint/examples/` as references (`simpleshop` is the smallest).

**2. Place the application in `blueprint/examples/{app}/`**, with the following structure:

```
blueprint/examples/{app}/
├── wiring/             # blueprint specification
├── workflow/           # blueprint workflow
│   └── {app}/          # microservices code
```

Note that `blueprint/` is a git submodule, so your application's code belongs to that repository, not to Aletheia's.

**3. Follow Aletheia's naming conventions.** The current implementation makes a few assumptions about the application's code, described in [technical-assumptions.md](./docs/technical-assumptions.md).

**4. Register the application** by adding an entry at the end of the `apps` list in `registry/apps.yaml`, replacing `{app}` with the name of your application's folder and `{spec}` with the name of its wiring spec:

```yaml
apps:
  # ... existing entries ...

  - name: {app}
    app_root: github.com/blueprint-uservices/blueprint/examples/{app}
    package_path: {app}/workflow/{app}
    spec_name: {app}_{spec}
    spec_path: github.com/blueprint-uservices/blueprint/examples/{app}/wiring/specs
    # optional
    sql_tables:
      - "{database_name}:blueprint/examples/{app}/workflow/{app}/database/{database_filename}.sql"
    nosql_path: blueprint/examples/{app}/workflow/{app}/{path}
```

The last two fields are optional and declare constraints that Aletheia cannot infer from the code:

- `sql_tables`: SQL files with `CREATE TABLE` statements (PostgreSQL syntax), whose `PRIMARY KEY` and `UNIQUE` columns become constraints.
- `nosql_path`: folder that holds one JSON file per NoSQL collection, whose `uniqueItems` become uniqueness constraints (see `blueprint/examples/dsb_mediamicroservices/workflow/mediamicroservices/database/`).

**5. Generate the application registry**, which updates `go.mod` and `internal/frameworks/blueprint/apps/apps.go` and rebuilds `bin/aletheia`:

```zsh
make registry
```

**6. Check that the wiring loads.** The `--init` flag only loads the application's wiring and exits, which is a quick way to find registration errors:

```zsh
./bin/aletheia --init {app}
```

**7. Run the analysis and review the warnings.** Run `./bin/aletheia {app}`, then see [Reading the Output](#reading-the-output) to interpret the warnings and [Suppressing Detection Warnings](#suppressing-detection-warnings) to ignore false positives.

## Technical Details

See [technical-details.md](./docs/technical-details.md) for how each pattern maps to the code, how cross-microservice foreign keys are inferred, and the current limitations.

## Citation

If you use Aletheia in your work, please cite our paper:

```bibtex
@inproceedings{ferreira2026aletheia,
  author = {Mafalda Sofia Ferreira and Jo{\~a}o Ferreira Loff and Jo{\~a}o Garcia and Rodrigo Rodrigues},
  title = {Aletheia: Automated Detection of Data Integrity Violations in Microservices},
  booktitle = {20th USENIX Symposium on Operating Systems Design and Implementation (OSDI 26)},
  year = {2026},
  isbn = {978-1-939133-55-7},
  address = {Seattle, WA},
  pages = {721--737},
  url = {https://www.usenix.org/conference/osdi26/presentation/ferreira},
  publisher = {USENIX Association},
  month = jul
}
```
