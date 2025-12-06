[0] (PointerObject PointerType) c (*shopping_app.CartServiceImpl struct{order_service shopping_app.OrderService, product_service shopping_app.ProductService, cart_db NoSQLDatabase})
[_1] (StructObject UserType) shopping_app.CartServiceImpl struct{order_service shopping_app.OrderService, product_service shopping_app.ProductService, cart_db NoSQLDatabase}
[__2] (FieldObject FieldType) cart_db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) cart_db NoSQLDatabase
[__2] (FieldObject FieldType) order_service shopping_app.OrderService
[___3] (ServiceObject ServiceType) order_service shopping_app.OrderService
[__2] (FieldObject FieldType) product_service shopping_app.ProductService
[___3] (ServiceObject ServiceType) product_service shopping_app.ProductService

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) cartID string

[0] (StructObject UserType) cart shopping_app.Cart struct{CartID string, UserID string, ProductID string, Quantity int, PricePerUnit int, Status string}
[_1] (FieldObject FieldType) PricePerUnit int
[__2] (BasicObject BasicType) PricePerUnit int
[_1] (FieldObject FieldType) ProductID string
[__2] (BasicObject BasicType) ProductID string
[_1] (FieldObject FieldType) Quantity int
[__2] (BasicObject BasicType) Quantity int
[_1] (FieldObject FieldType) UserID string
[__2] (BasicObject BasicType) UserID string

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = cart_database, collection = cart_collection}

[0] (InterfaceObject UserType) _ .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "cartID" string, Key "cartID" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "cartID" string
[___3] (BasicObject BasicType) "cartID" string
[__2] (FieldObject FieldType) Value string
[___3] (BasicObject BasicType) cartID string

[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = cart_database, collection = cart_collection}
[_1] (StructObject UserType) cart shopping_app.Cart struct{CartID string, UserID string, ProductID string, Quantity int, PricePerUnit int, Status string}
[__2] (FieldObject FieldType) PricePerUnit int
[___3] (BasicObject BasicType) PricePerUnit int
[__2] (FieldObject FieldType) ProductID string
[___3] (BasicObject BasicType) ProductID string
[__2] (FieldObject FieldType) Quantity int
[___3] (BasicObject BasicType) Quantity int
[__2] (FieldObject FieldType) UserID string
[___3] (BasicObject BasicType) UserID string

[0] (InterfaceObject UserType) _ .error

[0] (SliceObject UserType) filter primitive.D
[_1] (StructObject StructType) struct{Key "cartID" string, Key "cartID" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "cartID" string
[___3] (BasicObject BasicType) "cartID" string
[__2] (FieldObject FieldType) Value string
[___3] (BasicObject BasicType) cartID string

[0] (SliceObject UserType) update primitive.D
[_1] (StructObject StructType) struct{Key "status" string, Key "status" string, Value "completed" string, Value "completed" string}
[__2] (FieldObject FieldType) Key "status" string
[___3] (BasicObject BasicType) "status" string
[__2] (FieldObject FieldType) Value "completed" string
[___3] (BasicObject BasicType) "completed" string

