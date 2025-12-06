[0] (PointerObject PointerType) s (*shopping_simple.ProductServiceImpl struct{product_db NoSQLDatabase, product_queue Queue, num_workers int})
[_1] (StructObject UserType) shopping_simple.ProductServiceImpl struct{product_db NoSQLDatabase, product_queue Queue, num_workers int}
[__2] (FieldObject FieldType) product_db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) product_db NoSQLDatabase
[__2] (FieldObject FieldType) product_queue Queue
[___3] (BlueprintBackendObject BlueprintBackendType) product_queue Queue

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> r-tainted: read(product_db.[]shopping_simple.Product struct{ProductID string, Description string, PricePerUnit int, Category string}) {1}
[0] (ArrayObject ArrayType) products []shopping_simple.Product struct{ProductID string, Description string, PricePerUnit int, Category string}

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = product_database, collection = product_collection}

[0] (InterfaceObject UserType) _ .error

[0] (SliceObject UserType) filter primitive.D

    --> r-tainted: read(product_db.Product) {1}
[0] (BlueprintBackendObject BlueprintBackendType) cursor NoSQLCursor {database = product_database, collection = product_collection}
     --> r-tainted: read(product_db.[]shopping_simple.Product struct{ProductID string, Description string, PricePerUnit int, Category string}) {1}
[_1] (ArrayObject ArrayType) products []shopping_simple.Product struct{ProductID string, Description string, PricePerUnit int, Category string}

[0] (InterfaceObject UserType) _ .error

[0] (InterfaceObject UserType) err .error

