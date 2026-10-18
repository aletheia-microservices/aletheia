# Installation

## Requirements

- [Golang](https://go.dev/doc/install) >= 1.26

## Building Aletheia

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

> [!WARNING]
> Always run Aletheia from the repository root, since it loads the applications through the `replace` directives in `go.mod` and reads `registry/` and `config/` relative to the working directory.

These docs use `./bin/aletheia`, but either of these alternatives also works:

- `go run ./cmd/aletheia {app}`: skips the build step and always runs your latest changes.
- `make install`, then `aletheia {app}`: installs into `$GOBIN` (or `$(go env GOPATH)/bin`), which must be on your `PATH`. Rerun `make install` after changing the code or running `make registry`.

To get started, [run one of the included apps](running.md), then follow the [simpleshop tutorial](tutorial.md), which walks through registering an app, running the analysis, reading the warning, and suppressing it.
