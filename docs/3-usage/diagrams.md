# Generating Diagrams

Two scripts draw diagrams of an analyzed app from the files in `output/{app}/` and save them in `output/{app}/diagrams/`:

- **[App diagram](#app-diagram)**: the services and databases of the app and the dependencies between them, drawn from `app.json`.
- **[Schema diagram](#schema-diagram)**: the tables of each database with their fields and keys, and the inferred foreign keys between them, drawn from `schema.json`.

## Requirements

Each app needs an `output/{app}/app.json` (or `schema.json` for `schema-diagram`). If it's missing, analyze the app first (e.g., `./bin/aletheia trainticket`).

Both scripts need [Graphviz](https://graphviz.org/download/) to draw the image (`brew install graphviz` on macOS, `apt install graphviz` on Debian/Ubuntu). Without it, they print a warning and only write the `.dot` file.

## Running

Analyze the app first, then run the scripts from the repository root:

```zsh
./bin/aletheia postnotification
go run ./scripts/app-diagram postnotification      # writes output/postnotification/diagrams/app.png
go run ./scripts/schema-diagram postnotification   # writes output/postnotification/diagrams/schema.png
```

Both scripts take the same arguments, shown here for `app-diagram`:

```zsh
go run ./scripts/app-diagram trainticket               # a single app
go run ./scripts/app-diagram trainticket sockshop      # several apps
go run ./scripts/app-diagram -format svg trainticket   # svg or pdf instead of png
go run ./scripts/app-diagram -dpi 600 trainticket      # sharper png (default: 300 dpi)
go run ./scripts/app-diagram                           # every app with an output/{app}/app.json
```

## Output

For each app, the scripts write to `output/{app}/diagrams/`:

- `app.dot` and `schema.dot`: the graphs in Graphviz format.
- `app.{format}` and `schema.{format}`: the rendered images (`app.png` and `schema.png` by default). PNGs and other bitmap formats are drawn at 300 dpi, or the value of `-dpi`. SVG and PDF are vector formats, so they ignore `-dpi` and stay sharp at any zoom.

To redraw an image after editing its `.dot` file by hand:

```zsh
dot -Tpng -Gdpi=300 output/trainticket/diagrams/app.dot -o output/trainticket/diagrams/app.png
```

> [!WARNING]
> `make verify` replaces `output/{app}/` for every app it checks, which deletes `diagrams/`, so run the scripts again afterwards.

## App Diagram

Draws the service dependency graph of an app from the `output/{app}/app.json` that Aletheia writes when it analyzes the app:

- **Services** are blue boxes
- **Databases** are green cylinders
- **Edges** are service → service dependencies (grey) and service → database dependencies (green)

The service dependency graph of `trainticket`:

![Service dependency graph of trainticket](https://raw.githubusercontent.com/aletheia-microservices/aletheia/main/scripts/app-diagram/examples/trainticket.svg)

The [examples](https://github.com/aletheia-microservices/aletheia/tree/main/scripts/app-diagram/examples) folder has the diagram of every app in `output/`.

> [!NOTE]
> - Queues (e.g., `email_queue`) are listed as databases in `app.json`, so they are drawn as databases.
> - The per-service `databases` list in `app.json` is currently always empty, so the script finds service → database dependencies from each service's `fields` instead. For example, a field `db #0 (assurance_db)` links the service to `assurance_db`, because `assurance_db` is in the app's top-level `databases` list.

## Schema Diagram

Draws the data schema of an app from the `output/{app}/schema.json` that Aletheia writes when it analyzes the app:

- **Databases** are green boxes that group their tables
- **Tables** list their fields, marked `PK` (primary key), `UQ` (unique) or `FK` (foreign key)
- **Edges** are foreign keys, from the referencing field to the referenced field: red if mandatory, grey if not, and blue if transitive. See [Reading Detection Warnings](../4-formats/reading-warnings.md#foreign-keys) for what these mean

The data schema of `trainticket`:

![Data schema of trainticket](https://raw.githubusercontent.com/aletheia-microservices/aletheia/main/scripts/schema-diagram/examples/trainticket.svg)

The [examples](https://github.com/aletheia-microservices/aletheia/tree/main/scripts/schema-diagram/examples) folder has the diagram of every app in `output/`.

> [!NOTE]
> - Databases without tables in `schema.json` are not drawn. The script lists them when it runs.
> - A constraint in `schema.json` can name a field that is missing from its table's `fields` list. The script still adds that field to the table so the constraint can be drawn. Constraints on a table that is not in `schema.json` are skipped with a warning.

## Updating the Examples

To update the example diagrams after changing a script, run it on every app and copy each diagram into its `examples/` folder:

```zsh
go run ./scripts/app-diagram -format svg
for d in output/*/diagrams; do cp $d/app.svg scripts/app-diagram/examples/$(basename $(dirname $d)).svg; done

go run ./scripts/schema-diagram -format svg
for d in output/*/diagrams; do cp $d/schema.svg scripts/schema-diagram/examples/$(basename $(dirname $d)).svg; done
```

The diagrams on this page are loaded from the `main` branch on GitHub, so they show the updated examples once they are pushed.
