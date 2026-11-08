# Taint Representation

A taint records where a value is written to or read from. The [intra-procedural analysis](intra-procedural-analysis.md#tainting-ssa-graphs) annotates the arguments and return values of database operations and service invocations with taints and propagates them within each function, while [abstract models](../app-registry/abstract-model.md#objects-and-taints) declare them directly. Each taint is either a database taint or a service taint:

| Taint | Identifies | Example |
| ----- | ---------- | ------- |
| **Database taint** | The database operation and the schema field to or from which the value is written or read, as `<database>.<schema>.<field>` | `[write] [t21] @ posts_db.post.PostID` |
| **Service taint** | The service invocation and the interface parameter or return value to which the value corresponds, as `<service>.<method>.<name>`, where `<name>` is the name of one of the arguments or returns of that call | `[rpc] [t3] @ StorageService.StorePost.t4` |

In both cases, the taint starts with the type of the operation (e.g., `[write]` or `[rpc]`) and the SSA name of the call (e.g., `[t21]`), which also serves as its timestamp. Service taints mark where a value leaves or enters a service, which the [inter-microservice analysis](inter-microservice-analysis.md#tainting-the-abstract-call-graph) relies on to propagate taints across services.
