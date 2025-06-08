[0] (PointerObject PointerType) s (*shopping_simple.ProductServiceImpl struct{product_db NoSQLDatabase, product_queue Queue, num_workers int})
[_1] (StructObject UserType) shopping_simple.ProductServiceImpl struct{product_db NoSQLDatabase, product_queue Queue, num_workers int}
[__2] (FieldObject FieldType) product_db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) product_db NoSQLDatabase
[__2] (FieldObject FieldType) product_queue Queue
[___3] (BlueprintBackendObject BlueprintBackendType) product_queue Queue

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> w-tainted: write(product_queue.ProductQueueMessage.ProductID) {1}       --> w-tainted: write(product_queue.ProductQueueMessage.ProductID) {1} --> r-tainted: read(product_queue.ProductQueueMessage.ProductID) {1}
[0] (BasicObject BasicType) productID string
     --> w-tainted: write(product_queue.ProductQueueMessage.ProductID) {1}         --> w-tainted: write(product_queue.ProductQueueMessage.ProductID) {1} --> r-tainted: read(product_queue.ProductQueueMessage.ProductID) {1}
[_1] (Reference BasicType) ref <productID string> @ Frontend

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = product_database, collection = product_collection}

[0] (InterfaceObject UserType) _ .error

[0] (SliceObject UserType) filter primitive.D
[_1] (StructObject StructType) struct{Key "productid" string, Key "productid" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "productid" string
[___3] (BasicObject BasicType) "productid" string
[__2] (FieldObject FieldType) Value string
       --> w-tainted: write(product_queue.ProductQueueMessage.ProductID) {1}             --> w-tainted: write(product_queue.ProductQueueMessage.ProductID) {1} --> r-tainted: read(product_queue.ProductQueueMessage.ProductID) {1}
[___3] (BasicObject BasicType) productID string
        --> w-tainted: write(product_queue.ProductQueueMessage.ProductID) {1}               --> w-tainted: write(product_queue.ProductQueueMessage.ProductID) {1} --> r-tainted: read(product_queue.ProductQueueMessage.ProductID) {1}
[____4] (Reference BasicType) ref <productID string> @ Frontend

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) val bool

    --> w-tainted: write(product_queue.ProductQueueMessage) {1}       --> w-tainted: write(product_queue.ProductQueueMessage) {1} --> r-tainted: read(product_queue.ProductQueueMessage) {1}
[0] (StructObject UserType) message shopping_simple.ProductQueueMessage struct{ProductID string, Remove true bool}
     --> w-tainted: write(product_queue.ProductQueueMessage.ProductID) {1}         --> w-tainted: write(product_queue.ProductQueueMessage.ProductID) {1} --> r-tainted: read(product_queue.ProductQueueMessage.ProductID) {1}
[_1] (FieldObject FieldType) ProductID string
      --> w-tainted: write(product_queue.ProductQueueMessage.ProductID) {1}           --> w-tainted: write(product_queue.ProductQueueMessage.ProductID) {1} --> r-tainted: read(product_queue.ProductQueueMessage.ProductID) {1}
[__2] (BasicObject BasicType) productID string
       --> w-tainted: write(product_queue.ProductQueueMessage.ProductID) {1}             --> w-tainted: write(product_queue.ProductQueueMessage.ProductID) {1} --> r-tainted: read(product_queue.ProductQueueMessage.ProductID) {1}
[___3] (Reference BasicType) ref <productID string> @ Frontend
     --> w-tainted: write(product_queue.ProductQueueMessage.Remove) {1}         --> w-tainted: write(product_queue.ProductQueueMessage.Remove) {1} --> r-tainted: read(product_queue.ProductQueueMessage.Remove) {1}
[_1] (FieldObject FieldType) Remove true bool
      --> w-tainted: write(product_queue.ProductQueueMessage.Remove) {1}           --> w-tainted: write(product_queue.ProductQueueMessage.Remove) {1} --> r-tainted: read(product_queue.ProductQueueMessage.Remove) {1}
[__2] (BasicObject BasicType) true bool

[0] (BasicObject BasicType) val bool

[0] (InterfaceObject UserType) err .error

