# Writing Input Models

An input model is a YAML file that describes an application for Aletheia to analyze with `--input` (see [Registering an Application without Blueprint](../3-usage/without-blueprint.md), which also has a complete example in [tinyshop](../3-usage/without-blueprint.md#example-tinyshop)). It has four top-level fields:

```yaml
app: myapp                   # name of the app (the results are saved in output/myapp/)
databases: [...]             # databases of the app
entrypoints: [...]           # func_short_path of each function called by clients
functions: [...]             # every function reachable from the entrypoints
```

## Databases

```yaml
databases:
  - name: users_db
    type: RelationalDB       # NoSQLDatabase, RelationalDB, Cache or Queue
    schemas:                 # optional: keys that cannot be inferred from the calls
      - name: user           # collection, table or topic
        primary_key: [UserID]
        unique: [Email]      # each field is unique by itself
```

These are the equivalent of the backends in the Blueprint wiring, and of the `bson:"_id"` tags and SQL files that declare keys. Schemas without keys do not need to be listed: they are created from the calls.

## Functions

A function is a service method that is either an entrypoint or called by another function.

```yaml
functions:
  - func_short_path: shop.Frontend.Checkout   # unique ID of the function
    service: Frontend
    method: Checkout
    params: [...]            # objects, excluding the receiver and the context
    returns: [...]           # one list of objects for each return statement
    calls: [...]             # calls to other services and databases, in the order they are made
```

The objects in the same position of each return statement are merged into one return of the function. Calls to internal (non-service) methods are not listed as such: their calls are listed as calls of the function.

## Calls

```yaml
calls:
  - call_id: Frontend.Checkout.GetUser   # unique across the whole app
    call_ts: t4                          # position of the call in the function (t<number>, e.g., t4 or t4.t7)
    call_type: RPC                       # RPC or DB
    arguments: [...]                     # objects
    service_call:                        # only for RPC
      service: UserService
      method: GetUser
      func_short_path: shop.UserService.GetUser
      returns: [...]                     # objects
  - call_id: Frontend.Checkout.InsertOne
    call_ts: t9
    call_type: DB
    arguments: [...]
    database_call:                       # only for DB
      operation_type: write              # write, update, read, read_many or delete
      database: orders_db
      schema: order
      method: InsertOne
```

## Objects and Taints

An object has a name and the taints of each of its paths, where `_obj` is the object itself and, for example, `_obj.UserID` or `_obj.Items[*].ItemID` are its fields.

```yaml
- name: order
  taints:
    _obj.UserID:
      - {taint_type: TAINT_DATABASE, call_id: Frontend.Checkout.InsertOne, path: orders_db.order.UserID}
      - {taint_type: TAINT_SERVICE, call_id: Frontend.Checkout.GetUser, path: UserService.GetUser.user.ID}
```

- A **database taint** (`TAINT_DATABASE`) marks a value that is written to, or read from, the database field `path` (`<database>.<schema>[.<field>]`) by the DB call `call_id`. For reads, `database_taint: {read_key: true}` marks a value used to filter the read, and `database_taint: {read_value: true}` marks a value that holds what was read.
- A **service taint** (`TAINT_SERVICE`) marks a value that is passed to, or returned by, the RPC `call_id`, where `path` is `<service>.<method>.<name>[.<field>]` and `<name>` is the name of one of the arguments or returns of that call. This is how Aletheia follows data across services, so object names must match.

Taints can only reference calls of the same function. The order of each taint is the `call_ts` of its call, prefixed by `caller_t` when the call was made inside an internal method called at `caller_t` (e.g., `caller_t: t4` for a call at `t7` gives `t4.t7`).

> [!NOTE]
> The input model is validated before the analysis starts, so any error (e.g., an unknown `call_id`, `func_short_path` or database) is reported without analyzing the app.
