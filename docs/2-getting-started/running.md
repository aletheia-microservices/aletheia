# Running Aletheia

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
- `constraints.txt`: inferred constraints (foreign keys, primary keys and unique fields), grouped by database
- `analysis/`: warnings related to integrity violations

To analyze every application registered in `registry/apps.yaml`, one after the other, use the `--all` flag instead of an application name:

```zsh
./bin/aletheia --all
```

> [!TIP]
> You can also specify the `--debug` flag to obtain tainted _ssa graphs_ and _abstract call graph_ in `.dot` format saved under `output/postnotification/abstractcallgraph` and `output/postnotification/ssagraphs`, which can then be visualized in, for example, [Graphivz](https://dreampuf.github.io/GraphvizOnline/).
>
> ```zsh
> ./bin/aletheia --debug postnotification
> ```
