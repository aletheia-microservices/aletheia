# Suppressing Detection Warnings

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
linearizable_databases:           # databases that are not eventually consistent (Un-1 only)
  - <database>                    # database name, as shown in schema.json
```

- `ignore_foreignkeys` stops Aletheia from inferring foreign keys from the listed fields, so those foreign keys do not appear in `schema.json` and do not lead to warnings.
- `ignore_cascade` hides RI-1 warnings about records left behind in the given `database` and `entity`. To hide them only when the delete happens in a specific database and entity, set both `trigger_database` and `trigger_entity`.
- `linearizable_databases` marks the listed databases as linearizable: each write appears to take effect at a single point in time, and all clients observe writes in that order. By default, Aletheia assumes every database is eventually consistent, so replicas accept writes independently and merge them later. A linearizable database rejects the second of two concurrent writes of the same unique value, so Un-1 does not report unique writes to it.

> [!WARNING]
> Set `trigger_database` and `trigger_entity` together. If only one of them is set, it has no effect.

You can ignore specific inferred foreign keys. For example, `config/postnotification.yaml` ignores the foreign key on `notif_queue.notification.ReqID`:

```zsh
./bin/aletheia --detection_config config/postnotification.yaml postnotification
```

Or you can suppress missing cascading delete warnings. For example, `config/sockshop.yaml` ignores missing cascades from `cart_db.carts` to `order_db.orders`:

```zsh
./bin/aletheia --detection_config config/sockshop.yaml sockshop
```
