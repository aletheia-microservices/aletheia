# Register Abstract Model

An application can be analyzed without porting it to Blueprint by describing it in a YAML abstract model. The model specifies the application's databases and service methods, including their calls and the taints of their arguments and return values.

It is called *abstract* because it already has what Aletheia would otherwise extract from Blueprint and the Go code, so Aletheia can build the abstract call graph from it directly.

The application can be in any language, and the model can be written by hand or generated with a tool for that language.

Unlike Blueprint applications, abstract models are not registered in `registry/apps.yaml`. Aletheia reads them directly from the path passed to `--input`.

<div class="file-tree" markdown>

- :material-folder-outline: `input-models/`
    - :material-folder-outline: `{app}/`
        - :material-file-outline: `model.yaml`

</div>

> [!TIP]
> See [Guidelines for Writing Abstract Models](#guidelines-for-writing-abstract-models) for a detailed description of each field.

The repository includes two examples in `input-models/`:

| Example | Description |
|---|---|
| `postnotification` | Equivalent to `blueprint/examples/postnotification` (see [Example: Post Notification](#example-post-notification)) |
| `tinyshop` | Small shop split into frontend and backend abstract models (see [Example: Tiny Shop](#example-tiny-shop)) |

## Example: Post Notification

The Post Notification application is described by a single abstract model, equivalent to `blueprint/examples/postnotification`:

<div class="file-tree" markdown>

- :material-folder-outline: `input-models/postnotification/`
    - :material-file-outline: `blueprint.yaml` `UploadService`, `NotifyService` and `StorageService` (`posts_db` and `notifications_queue`)

</div>

<details markdown>
<summary>Abstract model: <code>blueprint.yaml</code> (excerpt)</summary>

```yaml title="input-models/postnotification/blueprint.yaml"
--8<-- "input-models/postnotification/blueprint.yaml:16:16"

--8<-- "input-models/postnotification/blueprint.yaml:19:23"

--8<-- "input-models/postnotification/blueprint.yaml:27:30"

--8<-- "input-models/postnotification/blueprint.yaml:32:32"
  # ... UploadService, NotifyService and the other methods of StorageService
--8<-- "input-models/postnotification/blueprint.yaml:304:331"
```

</details>

From the repository root, pass the application's folder or the path to a single abstract model using `--input`:

```zsh
./bin/aletheia --input input-models/postnotification
```

Similar to analyzing Blueprint applications, the results are saved in `output/{app}/`, where `{app}` corresponds to the `app` field of the abstract model (`postnotification_example`). 

The `--debug`, `--eval`, `--refs`, and `--detection_config` flags can also be used with `--input`, provided that the application specified in the detection configuration matches the abstract model.

## Example: Tiny Shop

An application's folder can contain multiple abstract models (`*.yaml` or `*.yml`), allowing different parts of the application to be described separately. 


The Tiny Shop application is split into a frontend that calls other services and a backend with services that access databases:

<div class="file-tree" markdown>

- :material-folder-outline: `input-models/tinyshop/`
    - :material-file-outline: `frontend.yaml` `Frontend.AddToCart` and `Frontend.DeleteProduct` (entrypoints)
    - :material-file-outline: `backend.yaml` `ProductService` (`product_db`) and `CartService` (`cart_db`)

</div>

<details markdown>
<summary>Abstract model: <code>frontend.yaml</code> (excerpt)</summary>

```yaml title="input-models/tinyshop/frontend.yaml"
--8<-- "input-models/tinyshop/frontend.yaml:16:16"

--8<-- "input-models/tinyshop/frontend.yaml:18:20"

--8<-- "input-models/tinyshop/frontend.yaml:22:22"
  # ... Frontend.AddToCart
--8<-- "input-models/tinyshop/frontend.yaml:96:126"
```

</details>

<details markdown>
<summary>Abstract model: <code>backend.yaml</code> (excerpt)</summary>

```yaml title="input-models/tinyshop/backend.yaml"
--8<-- "input-models/tinyshop/backend.yaml:4:4"

--8<-- "input-models/tinyshop/backend.yaml:6:10"

--8<-- "input-models/tinyshop/backend.yaml:18:18"
  # ... ProductService.GetProduct and CartService.AddItem
--8<-- "input-models/tinyshop/backend.yaml:70:97"
```

</details>

All abstract models in the folder are combined before the analysis, allowing functions defined in one abstract model to call functions defined in another, as `frontend.yaml` does with `backend.yaml`. Only the combined model needs to be valid, subject to the following requirements:

- The `databases`, `entrypoints`, and `functions` from all abstract models are combined.
- The `app` field must have the same value in all abstract models that define it, and at least one must define it.
- A database can be declared in multiple abstract models with the same type, and its schemas are combined.
- Each function, entrypoint, and `call_id` must be unique across all abstract models.

To analyze the application, pass its folder to `--input` so that both abstract models are combined:

```zsh
./bin/aletheia --input input-models/tinyshop
```

> [!WARNING]
> Neither file can be analyzed alone, as `frontend.yaml` calls functions defined in `backend.yaml`, and `backend.yaml` has no entrypoints.

### Reading the Output

The frontend establishes the relationship between the two backend databases, even though it does not access either directly. `Frontend.AddToCart` reads a `ProductID` through `ProductService.GetProduct` and passes it to `CartService.AddItem`. Aletheia tracks this value and infers the following foreign key in `output/tinyshop/constraints.txt`:

``` { .txt .warnings title="output/tinyshop/constraints.txt" }
FOREIGN_KEY cart_db.cart.ProductID REFERENCES product_db.product.ProductID
```

The backend alone does not provide enough information to infer this relationship. `CartService.AddItem` writes the `ProductID` received as an argument to `cart_db`, but `backend.yaml` does not indicate that this value originates from `product_db`. Only the frontend establishes this connection.

Based on the inferred foreign key, Aletheia detects two potential integrity violations. Deleting a product leaves behind the cart items that reference it (RI-1), and a product can be deleted while it is being added to a cart (RI-2):

``` { .txt .warnings title="output/tinyshop/analysis/foreign-key-cascade.txt" }
delete: Frontend.DeleteProduct() ... ProductService.DeleteProduct() ... product_db.product.DeleteOne()
	missing cascade #1: database={cart_db}, entity={cart}, pending_fields={ProductID}
```

``` { .txt .warnings title="output/tinyshop/analysis/foreign-key-concurrency.txt" }
delete: Frontend.DeleteProduct() ... ProductService.DeleteProduct() ... product_db.product.DeleteOne()
   write #1: Frontend.AddToCart() ... CartService.AddItem() ... cart_db.cart.InsertOne()
       - database={cart_db}, entity={cart}, written_fields={ProductID}
```

## Guidelines for Writing Abstract Models

An abstract model has four top-level fields:

```yaml title="input-models/myapp/model.yaml"
app: myapp                   # name of the app (the results are saved in output/myapp/)
databases: [...]             # databases of the app
entrypoints: [...]           # func_short_path of each function called by clients
functions: [...]             # every function reachable from the entrypoints
```

### Databases

```yaml title="input-models/myapp/model.yaml"
databases:
  - name: users_db
    type: RelationalDB       # NoSQLDatabase, RelationalDB, Cache or Queue
    schemas:                 # optional: keys that cannot be inferred from the calls
      - name: user           # collection, table or topic
        primary_key: [UserID]
        unique: [Email]      # each field is unique by itself
```

These are the equivalent of the backends in the Blueprint wiring, and of the `bson:"_id"` tags and SQL files that declare keys. Schemas without keys do not need to be listed, since they are created from the calls.

### Functions

A function is a service method that is either an entrypoint or called by another function.

```yaml title="input-models/myapp/model.yaml"
functions:
  - func_short_path: shop.Frontend.Checkout   # unique ID of the function
    service: Frontend
    method: Checkout
    params: [...]            # objects, excluding the receiver and the context
    returns: [...]           # one list of objects for each return statement
    calls: [...]             # calls to other services and databases, in the order they are made
```

The objects in the same position of each return statement are merged into one return of the function. Calls to internal (non-service) methods are not listed as such. Instead, their calls are listed as calls of the function.

### Calls

```yaml title="input-models/myapp/model.yaml"
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

```yaml title="input-models/myapp/model.yaml"
- name: order
  taints:
    _obj.UserID:
      - {taint_type: TAINT_DATABASE, call_id: Frontend.Checkout.InsertOne, path: orders_db.order.UserID}
      - {taint_type: TAINT_SERVICE, call_id: Frontend.Checkout.GetUser, path: UserService.GetUser.user.ID}
```

> [!TIP]
> Each taint is either a database taint (`TAINT_DATABASE`) or a service taint (`TAINT_SERVICE`).
>
> See [Aletheia Internals → Taint Representation](../pipeline/taint-representation.md) for how database taints and service taints are represented.

For database taints, `database_taint: {read_key: true}` marks a value used to filter the read, and `database_taint: {read_value: true}` marks a value that holds what was read.

Taints can only reference calls of the same function. The order of each taint is the `call_ts` of its call, prefixed by `caller_t` when the call was made inside an internal method called at `caller_t` (e.g., `caller_t: t4` for a call at `t7` gives `t4.t7`).

> [!NOTE]
> The abstract model is validated before the analysis starts, so any error (e.g., an unknown `call_id`, `func_short_path` or database) is reported without analyzing the app.
