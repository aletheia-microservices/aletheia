# Input Models

Aletheia can analyze an application described by a YAML input model instead of a Blueprint application registered in `registry/apps.yaml`. The input model replaces what Aletheia otherwise extracts from Blueprint and from the SSA analysis of the Go code: the databases of the application, and the calls made by each service method with the taints of their arguments and returns. From there, the analysis is the same: Aletheia builds the abstract call graph, infers the constraints and detects the integrity violations.

This means the application does not need to be written in Go or Blueprint, as long as its services and their calls can be described in this format (e.g., by hand, or with a tool that analyzes another language). **See [Format](#format) for how to write an input model.**

Each app has its own folder, which usually holds a single input model (optionally, it can hold several input models that are combined, see [Combining Input Models](#combining-input-models)):

```
input-models/
└── {app}/
    └── model.yaml
```

The folders in `input-models/` are examples of input models:

| Example | Description |
|---|---|
| [`postnotification`](postnotification/blueprint.yaml) | Equivalent of `blueprint/examples/postnotification`, with the same results as its Blueprint analysis (see `tests/expected/postnotification`) |
| [`tinyshop`](tinyshop/) | Small shop split into a [frontend](tinyshop/frontend.yaml) that only calls services and a [backend](tinyshop/backend.yaml) with the services that call databases (see [Combining Input Models](#combining-input-models)) |

## Running an Example: PostNotification

From the repository root, pass the folder of the app (or the path of a single input model) to `--input`:

```zsh
./bin/aletheia --input input-models/postnotification
# or, without building
go run ./cmd/aletheia --input input-models/postnotification
```

The results are saved in `output/{app}/`, where `{app}` is the `app` of the input model (`postnotification_example`), with the same structure as for Blueprint applications (see [Reading the Output](../README.md#reading-the-output)). The `--debug`, `--eval`, `--refs` and `--detection_config` flags also work with `--input`, as long as the app in the detection config matches the app of the input model.

## Combining Input Models

The folder of an app can hold more than one input model (`*.yaml` or `*.yml`), for example when different parts of the application are described separately:

```
input-models/
└── {app}/
    ├── frontend.yaml        # services called by clients, which only call other services
    └── backend.yaml         # services that call other services and databases
```

All input models in the folder are combined into one before the analysis, so an input model can call functions of another one (e.g., in [`tinyshop`](tinyshop/), `frontend.yaml` calls functions described in `backend.yaml`). Only the combined model has to be valid, where:

- the `databases`, `entrypoints` and `functions` of all input models are combined
- the `app` must be the same in all input models that name it (at least one must)
- a database can be declared in more than one input model with the same type, and its schemas are combined
- each function, entrypoint and `call_id` must still be unique across all input models

### Example: tinyshop

[`tinyshop`](tinyshop/) is split into a frontend, which only calls other services, and a backend, with the services that call the databases:

```
input-models/tinyshop/
├── frontend.yaml            # Frontend.AddToCart and Frontend.DeleteProduct (the entrypoints)
└── backend.yaml             # ProductService (product_db) and CartService (cart_db)
```

To analyze it, pass the folder of the app so that both input models are combined:

```zsh
./bin/aletheia --input input-models/tinyshop
```

**The foreign key between the backend databases is only inferred because the frontend is included.** Even though the frontend does not access any database, Aletheia follows the `ProductID` that `Frontend.AddToCart` reads from `ProductService.GetProduct` and passes to `CartService.AddItem`, and infers that the items of each cart reference products (in `output/tinyshop/constraints.txt`):

```
FOREIGN_KEY cart_db.cart.ProductID REFERENCES product_db.product.ProductID
```

The backend alone cannot show this: `CartService.AddItem` writes the `ProductID` it receives as a parameter to `cart_db`, and nothing in `backend.yaml` says that this value was read from `product_db`. Only the frontend connects the two services.

From this foreign key, Aletheia warns that deleting a product leaves the cart items that reference it (RI-1, in `output/tinyshop/analysis/foreign-key-cascade.txt`), and that a product can be deleted while it is added to a cart (RI-2, in `output/tinyshop/analysis/foreign-key-concurrency.txt`):

```
delete: Frontend.DeleteProduct() ... ProductService.DeleteProduct() ... product_db.product.DeleteOne()
	missing cascade #1: database={cart_db}, entity={cart}, pending_fields={ProductID}

delete: Frontend.DeleteProduct() ... ProductService.DeleteProduct() ... product_db.product.DeleteOne()
	write #1: Frontend.AddToCart() ... CartService.AddItem() ... cart_db.cart.InsertOne()
		- database={cart_db}, entity={cart}, written_fields={ProductID}
```

Running only one of the input models fails, since neither is valid by itself: `frontend.yaml` calls functions that are only described in `backend.yaml`, and `backend.yaml` has no entrypoints.

## Format

```yaml
app: myapp                   # name of the app (the results are saved in output/myapp/)
databases: [...]             # databases of the app
entrypoints: [...]           # func_short_path of each function called by clients
functions: [...]             # every function reachable from the entrypoints
```

### Databases

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

### Functions

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

### Calls

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

### Objects and Taints

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

The input model is validated before the analysis starts, so any error (e.g., an unknown `call_id`, `func_short_path` or database) is reported without analyzing the app.
