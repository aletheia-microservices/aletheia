# Register Blueprint Application

Follow these steps to analyze your own application. The [simpleshop tutorial](../getting-started/tutorial.md) goes through steps 3 to 6 on a small example.

> [!TIP]
> See [Guidelines for Writing Blueprint Applications](#guidelines-for-writing-blueprint-applications) at the end of this page for the naming and code conventions Aletheia expects in Blueprint applications.

#### 1. Port Application to Blueprint

 Skip this step if your application already uses Blueprint. A Blueprint application has two components. The *workflow* defines each service as a Go interface and its implementation, and the *wiring* specifies the services and databases and how they are related.
 
 See the [Blueprint repository](https://github.com/Blueprint-uServices/blueprint) for its documentation, and use the applications in `blueprint/examples/` as references (`simpleshop` is the smallest).

#### 2. Place Application in Blueprint Folder

Place the application in `blueprint/examples/{app}/`, with the following structure:

<div class="file-tree" markdown>

- :material-folder-outline: `blueprint/examples/{app}/`
    - :material-folder-outline: `wiring/` Blueprint wiring spec
    - :material-folder-outline: `workflow/` Blueprint workflow
        - :material-folder-outline: `{app}/` Code of the services

</div>

> [!NOTE]
> `blueprint/` is a git submodule, so your application's code belongs to that repository, not to Aletheia's.

#### 3. Register Application in Aletheia

Register the application by adding an entry at the end of the `apps` list in `registry/apps.yaml`, replacing `{app}` with the name of your application's folder and `{spec}` with the name of its wiring spec:

```yaml title="registry/apps.yaml"
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

The last two fields are optional and declare constraints that Aletheia cannot automatically infer from the code:

- `sql_tables`: SQL files with `CREATE TABLE` statements, whose `PRIMARY KEY` and `UNIQUE` columns become constraints.
- `nosql_path`: folder that holds one JSON file per NoSQL collection, whose `uniqueItems` become uniqueness constraints (see `blueprint/examples/dsb_mediamicroservices/workflow/mediamicroservices/database/`). Blueprint's database interfaces do not include index creation, so these constraints must be declared explicitly.

#### 4. Generate Application Registry

Generate the application registry, which updates `go.mod` and `internal/frameworks/blueprint/apps/apps.go` and rebuilds `bin/aletheia`:

```zsh
make registry
```

#### 5. Check Application Wiring

The `--init` flag only loads the application's wiring and exits, which is a quick way to find registration errors:

```zsh
./bin/aletheia --init {app}
```

#### 6. Run Analysis

Run the analysis:

```zsh
./bin/aletheia {app}
```

## Guidelines for Writing Blueprint Applications

The current implementation of Aletheia assumes that Blueprint applications follow the naming and code conventions below. These are implementation limitations that future versions may remove.

### Database Naming

The name of the database used in the wiring specification (e.g., `mongodb.Container`) must exactly match the name passed in service operations (e.g., `GetCollection`).

```Go
// wiring
posts_db := mongodb.Container(spec, "posts_db")
// workflow
s.postsDb.GetCollection(ctx, "posts_db", ...)
```

### Service Structure

Name the service struct `<service-interface>Impl`.

```Go
type StorageService interface {...}
type StorageServiceImpl struct {...}
```

Service implementation functions must return the interface type (e.g., `StorageService`) rather than the concrete service struct (e.g., `StorageServiceImpl`).

```Go
func NewStorageServiceImpl(ctx context.Context, ...) (StorageService, error) {
  return &StorageServiceImpl{...}, nil
}
```

### Schema Fields

For NoSQL databases, BSON tags and filter keys must match the corresponding Go struct field names. The exception is `_id`, which Aletheia uses to infer primary keys.

```Go
type Movie struct {
	MovieID string `bson:"_id"`
	Title   string `bson:"Title"`
}
...
movie := Movie{...}
collection.InsertOne(ctx, movie)
...
query := bson.D{{Key: "Title", Value: title}}
collection.FindOne(ctx, query)
```
