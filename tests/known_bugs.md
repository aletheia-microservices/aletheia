# Known Bugs

Found while writing the test suite in `tests/`. Each item has a test that is skipped with a `known bug:` message. To check a fix, remove the `t.Skip(...)` line and run the test. All of these tests fail today when unskipped. Fixed items are removed from this list.

| # | Area | Summary | Impact | Test |
|---|------|---------|--------|------|
| 1 | Abstract graph | `GetPrimaryTaints` / `GetSecondaryTaints` / `GetWriteTaints` do not filter | High | `TestAbstractObjectPrimaryAndSecondaryFilters` |
| 2 | Schema | `GetLastSchema()` returns `schemas[0]`, whose order is random, so fields and constraints land in the wrong collection | High | `TestDetectionOutput/sockshop/constraints` |
| 3 | Detection | `TaintMapping.Clear()` does nothing, so the backward phase of an RPC reuses the forward mapping | Medium | `TestTaintMappingClear` |
| 4 | SSA taints | SQL calls whose statement is built at runtime are silently ignored | Medium | `TestSockshopCatalogueListReadsSocks` |
| 5 | Abstract graph | Database and RPC calls two or more helper calls deep are dropped | High | `TestSockshopRemoveItemDeletesEmptyCart`, `TestEshopOrderConsumerStoresOrder`, `TestSocialNetworkUnfollowWithUsernameUpdatesGraph` |

Items 1 and 2 are ordered by estimated impact. Item 3 was found later, while adding unit tests for `pkg/abstractgraph`, and items 4 and 5 while adding the abstract call graph tests in `tests/integration/abstractcallgraph`.

---

## 1. `GetPrimaryTaints`, `GetSecondaryTaints` and `GetWriteTaints` do not filter

**Where:** `pkg/abstractgraph/object.go:238`, `:264`, `:276`

Each function builds a filtered map and then returns `obj.taints`:

```go
func (obj *AbstractObject) GetPrimaryTaints() map[string][]*AbstractTaint {
	primaryTaints := make(map[string][]*AbstractTaint, 0)
	for objpath, taints := range obj.taints {
		for _, taint := range taints {
			if taint.IsPrimary() {
				primaryTaints[objpath] = append(primaryTaints[objpath], taint)
			}
		}
	}
	return obj.taints // should be primaryTaints
}
```

**Impact:** every caller receives secondary taints as well:
- `abstractgraph.parseDatabaseCall`: args → database params.
- `abstractgraph.registerDatabaseFields`: creates schema fields.
- `detection.Iterator.transverseQueue`: queue push → pop propagation.
- `AbstractObject.GetAffectedDatabaseFieldsForCall`: used by the cascade and concurrency detectors to decide which fields a write affects.

During the detection phase, objects hold secondary taints propagated from other requests. Writes can then look as if they touch fields they don't. This could produce false positives, and it can also create extra schema fields.

**Fix:** return the filtered map. Some detection results will change. Review the diff (`go test ./tests/integration -run TestDetectionOutput`) before running `-update`.

## 2. `GetLastSchema()` returns the first schema, in random order

**Where:** `pkg/app/backends/database.go:78`

```go
func (database *Database) GetLastSchema() *Schema {
	return database.schemas[0]
}
```

**Callers:** `pkg/abstractgraph/parser.go:266` and `pkg/abstractgraph/tainter.go:301, 311, 414, 454, 510, 544, 563`.

- **Wrong collection:** whatever the name suggests, fields and constraints always go to the first schema, never the one matching the taint's collection. On databases with several collections they land in the wrong one.
- **Random order:** schemas are created in `registry.RegisterFields` (`pkg/ssagraph/registry/nosql_primary_keys.go:39`), which receives `graphsLst`. `main.go` builds that list by iterating the `funcGraphs` map, so which schema ends up at `schemas[0]` changes between runs.

**Evidence (sockshop):**
- `user_db` has the `user`, `address` and `card` collections. On some runs the `address` schema in `schema.json` contains `user_db.card.*` and `user_db.user.*` fields.
- `schema.json` differs from `output-expected/` on HEAD, and between repeated runs of the same binary.
- The inferred constraints vary too. `FOREIGN_KEY order_db.orders.Shipment.Name REFERENCES ship_db.shipments.Name` (and the same key into `ship_queue.notification.Name`) appeared in 7 of 12 runs.
- The detection output was identical in all 12 runs, but that looks like luck rather than something guaranteed.
- trainticket's `schema.json` also changes order between runs.

**Fix:** look up the schema from the field path, e.g. `db.GetSchemaByNameIfExists(utils.ExtractSchemaNameFromFieldPath(path))`. Separately, sort `graphsLst` (or the `funcGraphs` keys) in `main.go` so that the remaining order-dependent steps are reproducible. Once constraints are stable, remove `sockshop` from `nondeterministicConstraints` in `tests/integration/detection_test.go` and run `-update`.

## 3. `TaintMapping.Clear()` does nothing

**Where:** `pkg/abstractgraph/taintmapping.go:33`, called from `pkg/detection/iterator.go:192`

```go
func (tm *TaintMapping) Clear() {
	tm = &TaintMapping{mapping: make(map[AbstractTaint][]AbstractTaint)} // only reassigns the local receiver
}
```

**Impact:** for each RPC, `Iterator.transverse` builds a mapping while propagating forward (caller args → callee params), visits the callee, and then calls `Clear()` before propagating backward (callee → caller). Because nothing is cleared, the backward phase still holds the forward pairs, and:
- `PropagateNewTaintsToDatabaseCallObjects` applies them to the database call arguments of the caller, although they were only meant for the callee's.
- `PropagateNewTaintsToDatabaseSchemas` applies them to the schemas a second time, which is redundant.

The first point can add taints on the caller side, and later constraints, that the backward phase alone would not create.

**Fix:** clear the fields in place (`tm.mapping = make(...)`, `tm.mappingKeys = nil`), or create a new mapping at `iterator.go:192` with `taintMapping = abstractgraph.NewTaintMapping()`. Note that `NewTaintMapping()` also resets the global `seenTaintPairsOnJoin`, so the two fixes differ slightly. Review the diff with `TestDetectionOutput`.

**Tried:** clearing the fields in place changes no warnings in the included apps, but removes one inferred foreign key in `dsb_mediamicroservices`: `movie_info_cache.*.Key REFERENCES movie_info_db.movie_info._id [T]`. That foreign key looks correct, since `MovieInfoService.ReadMovieInfo` uses the same `movieID` as the cache key and as the `_id` filter of `movie_info_db`. It is currently found only as a side effect of this bug, so find out why it is not inferred directly before fixing.

## 4. SQL calls whose statement is built at runtime are silently ignored

**Where:** `pkg/ssagraph/tainter/blueprint_calls.go:172`

`isBlueprintRelationalDBCall` reads the SQL statement with `utils.ExtractStringFromValue`, which only works for constants. When the statement is built at runtime, the call is treated as not being a database call, with no warning.

**Evidence:** sockshop `CatalogueService.List` builds its query from a base query and the requested tags (`catalogue/catalogueservice.go:47-72`) before calling `Select`, and has no database call in the abstract call graph.

**Impact:** medium. In the included apps it only hides one read, but a dynamic `INSERT` or `DELETE` would be hidden too.

**Fix:** at least log a warning. To support the call, fall back to the table of the base query (e.g., by following the `BinOp` that builds the string), or to the tables of the database schema.

## 5. Database and RPC calls two or more helper calls deep are dropped

**Where:** `pkg/abstractgraph/parser.go:253`

`parseMethodCall` adds the calls of a helper method to the service node. When the helper calls another helper, it recurses with the outer graph instead of the helper's own graph:

```go
toSSAGraph := fromSSAGraph.GetInlinedGraphForMethodCallIfExists(methodCall) // line 237
...
parseMethodCall(graph, node, fromSSAGraph, methodCall, funcGraphs) // should be toSSAGraph
```

The nested helper is registered as a inlined graph of `toSSAGraph`, not of `fromSSAGraph`, so the lookup at line 237 returns nil and the `// should never happen` branch returns without adding anything.

**Evidence:** 9 database calls two helpers deep are missing from the abstract call graph:
- sockshop (6): `CartService.RemoveItem` and `CartService.UpdateItem` → `DeleteCart` → `deleteMany` (`DeleteMany` on `cart_db.carts`); `UserService.Login` → `userdb_GetUserAttributes` → address and card reads; `UserService.Register` → `userdb_CreateUser` → address and card deletes.
- eshopmicroservices (1): `OrderService.Init` → `CreateNewOrder` → `add`, which inserts the order into `order_db`.
- dsb_socialnetwork (2): `SocialGraphService.UnfollowWithUsername` → `Unfollow` → the two go routines that update `socialgraph_db`.

RPCs made at that depth are dropped the same way.

**Impact:** high. Deletes and writes are missing from the requests that make them, so the cascade and concurrency detectors never see them (e.g., the cart deleted by `RemoveItem` in sockshop).

**Fix:** pass `toSSAGraph` in the recursive call. The counts and edge lists in `tests/integration/abstractcallgraph` and some detection results will change; review both diffs.

---

## Other observations

- A regression from the ongoing `Call` interface refactor: `tainter.go` called `ssagraph.NewMethodCall(callId, nil, ...)`, which dereferences the nil node. This crashed `dsb_socialnetwork` in `scripts/regression/regression.sh`. It is already fixed by `HasMethodCall(callID string)`.
- `pkg/utils/comparator.go:27`: `go vet` reports unreachable code, a `logrus.Fatalf` after `return nil, false` in `parseT`. As a result, timestamps without a `t` prefix are silently treated as smaller than any other.
