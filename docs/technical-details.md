# Technical Details

## Detection of Code Patterns

Integrity violations are detected by searching for the following patterns, formalized in our paper and implemented in [`internal/analysis/system-level/detection/constraints/`](../internal/analysis/system-level/detection/constraints/):

| ID   | Constraint            | Problematic Pattern          | Implementation Package                                                       |
| ---- | --------------------- | ---------------------------- | ---------------------------------------------------------------------------- |
| RI-1 | Referential integrity | Absence of cascading effects | `internal/analysis/system-level/detection/constraints/foreignkeycascade`     |
| RI-2 | Referential integrity | Concurrent operations        | `internal/analysis/system-level/detection/constraints/foreignkeyconcurrency` |
| RI-3 | Referential integrity | Uncoordinated replication    | `internal/analysis/system-level/detection/constraints/keycoordination`       |
| EI-1 | Entity integrity      | Uncoordinated replication    | `internal/analysis/system-level/detection/constraints/keycoordination`       |
| Un-1 | Uniqueness            | Conflicting writes           | `internal/analysis/system-level/detection/constraints/uniquenessconcurrency` |

> [!NOTE]
> Refer to our paper for the formal definitions of these patterns.

## Cross-Microservice Foreign Key Inference

Data associations across microservices are inferred from taints propagated through related objects used in database operations in the _abstract call graph_. The inference is performed according to the rules implemented in [`internal/analysis/system-level/abstractgraph/tainter/tainter.go`](../internal/analysis/system-level/abstractgraph/tainter/tainter.go).

Rules are applied for each pair of operation (`op_1`, `op_2`) where `op_i` is either a `read` or a `write`. We use `field_1` and `field_2` to denote fields accessed (tainted) by the same object in `op_1` and `op_2`, respectively.

| Operation Pair   | Foreign Key Direction        |
| ---------------- | ---------------------------- |
| `(write, write)` | `field2` references `field1` |
| `(read, write)`  | `field2` references `field1` |
| `(write, read)`  | `field1` references `field2` |

In read operations, the `read_key` denotes that the propagated object is used as a filter in the read operation, while `read_val` denotes that the propagated object is returned from the read operation.

| Operation Pair         | Foreign Key Direction        |
| ---------------------- | ---------------------------- |
| `(read_key, read_key)` | `field2` references `field1` |
| `(read_val, read_key)` | `field1` references `field2` |

The `(read_val, read_key)` rule is disabled by default. To enable it, set `CreateReferencesFromReadReadPairAndValKey` in [`internal/config/config.go`](../internal/config/config.go).

> [!NOTE]
> Refer to our paper for a detailed explanation of these inference rules.

## Current Limitations

See [technical-assumptions.md](./technical-assumptions.md) for the current analysis assumptions and limitations.
