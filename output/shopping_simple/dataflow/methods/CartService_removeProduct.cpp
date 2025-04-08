[0] (PointerObject PointerType) s (*shopping_simple.CartServiceImpl struct{product_service shopping_simple.ProductService, cart_db NoSQLDatabase, product_queue Queue, num_workers int})
[_1] (StructObject UserType) shopping_simple.CartServiceImpl struct{product_service shopping_simple.ProductService, cart_db NoSQLDatabase, product_queue Queue, num_workers 4 int}
[__2] (FieldObject FieldType) cart_db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) cart_db NoSQLDatabase
[__2] (FieldObject FieldType) num_workers 4 int
[___3] (BasicObject BasicType) 4 int
[__2] (FieldObject FieldType) product_queue Queue
[___3] (BlueprintBackendObject BlueprintBackendType) product_queue Queue
[__2] (FieldObject FieldType) product_service shopping_simple.ProductService
[___3] (ServiceObject ServiceType) product_service shopping_simple.ProductService

[0] (InterfaceObject UserType) ctx context.Context

[0] (StructObject UserType) message shopping_simple.ProductQueueMessage struct{ProductID string, Remove bool}
[_1] (FieldObject FieldType) ProductID string
[__2] (BasicObject BasicType) ProductID string

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = cart_database, collection = cart_database}

[0] (InterfaceObject UserType) _ .error

[0] (SliceObject UserType) filter primitive.D

[0] (SliceObject UserType) update primitive.D
[_1] (StructObject StructType) struct{Key "$dec" string, Key "$dec" string, Value primitive.D, Value primitive.D}
[__2] (FieldObject FieldType) Key "$dec" string
[___3] (BasicObject BasicType) "$dec" string
[__2] (FieldObject FieldType) Value primitive.D
[___3] (SliceObject UserType) primitive.D
[____4] (StructObject StructType) struct{Key "totalquantity" string, Key "totalquantity" string, Value 1 int, Value 1 int}
[_____5] (FieldObject FieldType) Key "totalquantity" string
[______6] (BasicObject BasicType) "totalquantity" string
[_____5] (FieldObject FieldType) Value 1 int
[______6] (BasicObject BasicType) 1 int
[_1] (StructObject StructType) struct{Key "$pull" string, Key "$pull" string, Value primitive.D, Value primitive.D}
[__2] (FieldObject FieldType) Key "$pull" string
[___3] (BasicObject BasicType) "$pull" string
[__2] (FieldObject FieldType) Value primitive.D
[___3] (SliceObject UserType) primitive.D
[____4] (StructObject StructType) struct{Key "Products" string, Key "Products" string, Value string, Value string}
[_____5] (FieldObject FieldType) Key "Products" string
[______6] (BasicObject BasicType) "Products" string
[_____5] (FieldObject FieldType) Value string
[______6] (BasicObject BasicType) ProductID string
[_1] (StructObject StructType) struct{Key "$pull" string, Key "$pull" string, Value primitive.D, Value primitive.D}
[__2] (FieldObject FieldType) Key "$pull" string
[___3] (BasicObject BasicType) "$pull" string
[__2] (FieldObject FieldType) Value primitive.D
[___3] (SliceObject UserType) primitive.D
[____4] (StructObject StructType) struct{Key "LastProductID" string, Key "LastProductID" string, Value string, Value string}
[_____5] (FieldObject FieldType) Key "LastProductID" string
[______6] (BasicObject BasicType) "LastProductID" string
[_____5] (FieldObject FieldType) Value string
[______6] (BasicObject BasicType) ProductID string

[0] (BasicObject BasicType) updated int

[0] (InterfaceObject UserType) err .error

