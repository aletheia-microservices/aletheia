[0] (PointerObject PointerType) s (*shopping_app.ProductServiceImpl struct{product_db NoSQLDatabase})
[_1] (StructObject UserType) shopping_app.ProductServiceImpl struct{product_db NoSQLDatabase}
[__2] (FieldObject FieldType) product_db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) product_db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) productID string

[0] (BasicObject BasicType) description string

[0] (BasicObject BasicType) pricePerUnit int

[0] (BasicObject BasicType) category string

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = product_database, collection = product_collection}

[0] (InterfaceObject UserType) _ .error

[0] (StructObject UserType) product shopping_app.Product struct{ProductID string, Description string, PricePerUnit int, Category string}
[_1] (FieldObject FieldType) Category string
[__2] (BasicObject BasicType) category string
[_1] (FieldObject FieldType) Description string
[__2] (BasicObject BasicType) description string
[_1] (FieldObject FieldType) PricePerUnit int
[__2] (BasicObject BasicType) pricePerUnit int
[_1] (FieldObject FieldType) ProductID string
[__2] (BasicObject BasicType) productID string

[0] (InterfaceObject UserType) err .error

