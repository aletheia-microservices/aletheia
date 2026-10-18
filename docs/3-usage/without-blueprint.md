# Registering an Application without Blueprint

Aletheia can analyze an application described by a YAML input model instead of a Blueprint application registered in `registry/apps.yaml`. The input model replaces what Aletheia otherwise extracts from Blueprint and from the SSA analysis of the Go code: the databases of the application, and the calls made by each service method with the taints of their arguments and returns. From there, the analysis is the same: Aletheia builds the abstract call graph, infers the constraints and detects the integrity violations.

This means the application does not need to be written in Go or Blueprint, as long as its services and their calls can be described in this format (e.g., by hand, or with a tool that analyzes another language). **See [Writing Input Models](../4-formats/writing-input-models.md) for the format.**

Each app has its own folder, which usually holds a single input model (optionally, it can hold several input models that are combined, see [Combining Input Models](#combining-input-models)):

```
input-models/
└── {app}/
    └── model.yaml
```

The folders in `input-models/` are examples of input models:

| Example | Description |
|---|---|
| [`postnotification`](https://github.com/aletheia-microservices/aletheia/blob/main/input-models/postnotification/blueprint.yaml) | Equivalent of `blueprint/examples/postnotification`, with the same results as its Blueprint analysis (see `tests/expected/postnotification`) |
| [`tinyshop`](https://github.com/aletheia-microservices/aletheia/tree/main/input-models/tinyshop) | Small shop split into a [frontend](https://github.com/aletheia-microservices/aletheia/blob/main/input-models/tinyshop/frontend.yaml) that only calls services and a [backend](https://github.com/aletheia-microservices/aletheia/blob/main/input-models/tinyshop/backend.yaml) with the services that call databases (see [Combining Input Models](#combining-input-models)) |

## Example: postnotification

From the repository root, pass the folder of the app (or the path of a single input model) to `--input`:

```zsh
./bin/aletheia --input input-models/postnotification
# or, without building
go run ./cmd/aletheia --input input-models/postnotification
```

The results are saved in `output/{app}/`, where `{app}` is the `app` of the input model (`postnotification_example`), with the same structure as for Blueprint applications (see [Reading Detection Warnings](../4-formats/reading-warnings.md)). The `--debug`, `--eval`, `--refs` and `--detection_config` flags also work with `--input`, as long as the app in the detection config matches the app of the input model.

## Combining Input Models

The folder of an app can hold more than one input model (`*.yaml` or `*.yml`), for example when different parts of the application are described separately:

```
input-models/
└── {app}/
    ├── frontend.yaml        # services called by clients, which only call other services
    └── backend.yaml         # services that call other services and databases
```

All input models in the folder are combined into one before the analysis, so an input model can call functions of another one (e.g., in [`tinyshop`](https://github.com/aletheia-microservices/aletheia/tree/main/input-models/tinyshop), `frontend.yaml` calls functions described in `backend.yaml`). Only the combined model has to be valid, where:

- the `databases`, `entrypoints` and `functions` of all input models are combined
- the `app` must be the same in all input models that name it (at least one must)
- a database can be declared in more than one input model with the same type, and its schemas are combined
- each function, entrypoint and `call_id` must still be unique across all input models

### Example: tinyshop

[`tinyshop`](https://github.com/aletheia-microservices/aletheia/tree/main/input-models/tinyshop) is split into a frontend, which only calls other services, and a backend, with the services that call the databases:

```
input-models/tinyshop/
├── frontend.yaml            # Frontend.AddToCart and Frontend.DeleteProduct (the entrypoints)
└── backend.yaml             # ProductService (product_db) and CartService (cart_db)
```

<details markdown>
<summary>Input models: <code>frontend.yaml</code> and <code>backend.yaml</code></summary>

```yaml title="input-models/tinyshop/frontend.yaml"
--8<-- "input-models/tinyshop/frontend.yaml"
```

```yaml title="input-models/tinyshop/backend.yaml"
--8<-- "input-models/tinyshop/backend.yaml"
```

</details>

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

> [!WARNING]
> Running only one of the input models fails, since neither is valid by itself: `frontend.yaml` calls functions that are only described in `backend.yaml`, and `backend.yaml` has no entrypoints.
