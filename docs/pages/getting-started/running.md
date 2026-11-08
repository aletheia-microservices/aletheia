# Running Aletheia

Aletheia analyzes the applications in `blueprint/examples/` (see [Getting Started → Example Applications](example-applications.md)). Pass the application name to Aletheia:

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
- `constraints.txt`: inferred constraints (foreign keys, primary keys and unique fields), grouped by database
- `analysis/`: warnings related to integrity violations

To analyze every application registered in `registry/apps.yaml`, one after the other, use the `--all` flag instead of an application name:

```zsh
./bin/aletheia --all
```

> [!TIP]
> Add `--debug` to save the tainted SSA graphs and abstract call graph as `.dot` files in `output/postnotification/ssagraphs/` and `output/postnotification/abstractcallgraph/`. You can view them with [Graphviz Online](https://dreampuf.github.io/GraphvizOnline/).
>
> ```zsh
> ./bin/aletheia --debug postnotification
> ```
