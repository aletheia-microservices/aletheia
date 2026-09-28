# Registry Generator

Generates the code that lets Aletheia import the Blueprint wiring spec of each app registered in `registry/*.yaml`.

For each `registry/{name}.yaml`, the generator:

1. Writes `internal/frameworks/blueprint/apps/{name}.go` (e.g., `apps.yaml` → `apps.go`), which imports each app's wiring spec and maps its `spec_name` to that spec. Don't edit this file by hand.
2. Adds `require` and `replace` entries to `go.mod`, so that Go finds each app's modules in `blueprint/examples/{app}/`.

## When to Run

Run it after you add or remove an app in `registry/apps.yaml`, or change an app's `app_root`, `spec_name` or `spec_path`.

You don't need to run it after changing the other fields (`package_path`, `sql_tables`, `nosql_path`), because Aletheia reads those from the YAML file each time it runs.

## Running

From the repository root:

```zsh
make registry                                      # every registry/*.yaml, then rebuilds bin/aletheia
go run ./scripts/gen-registry                      # every registry/*.yaml, without rebuilding
go run ./scripts/gen-registry -config apps.yaml    # a single file, relative to registry/
```

If you use `go run`, rebuild the binary afterwards with `make build`.

The generator exits with an error unless it runs from a directory named `aletheia`.

## Registry Format

See step 4 of [Analyzing Your Own Application](../../README.md#analyzing-your-own-application) for how to add an app. The generator uses these fields:

- `build_tag`: the Go build constraint written at the top of the generated file (default `!eval`).
- `app_root`: the module path of the app. The `go.mod` entries point it, and its `/workflow` and `/wiring` modules, to `./blueprint/examples/{last segment of app_root}`.
- `spec_path`: the import path of the package that defines the wiring spec.
- `spec_name`: must be `{name}_{suffix}`. The generated code uses the spec named after the capitalized suffix (e.g., `digota_docker` → `specs.Docker`, `dsb_hotel_original` → `specs.Original`).

Apps without a `spec_name` or `spec_path` are left out of the generated file. The generator fails if no app in the file has both.

## Notes

- The generator only adds the `go.mod` entries that are missing, and never removes any. When you remove an app, delete its `require` and `replace` entries from `go.mod` by hand.
- To undo the changes, restore the tracked files and rebuild:

  ```zsh
  git restore registry/apps.yaml go.mod internal/frameworks/blueprint/apps/apps.go
  make build
  ```
