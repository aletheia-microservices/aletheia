# Registering an Application in Blueprint

Follow these steps to analyze your own application. The [simpleshop tutorial](../2-getting-started/tutorial.md) goes through steps 4 to 7 on a small example.

**1. Port the application to Blueprint.** Skip this step if your application is already written with Blueprint. A Blueprint application has two parts: the _workflow_, where each service is a Go interface and its implementation, and the _wiring_, a spec that creates the services and databases and connects them. See Blueprint's guides on [workflow](https://github.com/Blueprint-uServices/blueprint/blob/main/docs/manual/workflow.md) and [wiring](https://github.com/Blueprint-uServices/blueprint/blob/main/docs/manual/wiring.md), and use the applications in `blueprint/examples/` as references (`simpleshop` is the smallest).

**2. Place the application in `blueprint/examples/{app}/`**, with the following structure:

```
blueprint/examples/{app}/
├── wiring/             # blueprint specification
├── workflow/           # blueprint workflow
│   └── {app}/          # microservices code
```

> [!NOTE]
> `blueprint/` is a git submodule, so your application's code belongs to that repository, not to Aletheia's.

**3. Follow Aletheia's naming conventions.** The current implementation makes a few assumptions about the application's code, described in [Conventions](../5-reference/conventions.md).

**4. Register the application** by adding an entry at the end of the `apps` list in `registry/apps.yaml`, replacing `{app}` with the name of your application's folder and `{spec}` with the name of its wiring spec:

```yaml title="registry/apps.yaml" hl_lines="4-12"
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

**7. Run the analysis and review the warnings.** Run `./bin/aletheia {app}`, then see [Reading Detection Warnings](../4-formats/reading-warnings.md) to understand the warnings and [Suppressing Detection Warnings](suppressing-warnings.md) to ignore false positives.
