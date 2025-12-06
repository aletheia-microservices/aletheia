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
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> w-tainted: write(cart_db.Cart.CartID) {1}
[0] (BasicObject BasicType) cartID string
     --> w-tainted: write(cart_db.Cart.CartID) {1}
[_1] (Reference BasicType) ref <cartID string> @ Frontend

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = cart_database, collection = cart_database}

[0] (InterfaceObject UserType) _ .error

    --> w-tainted: write(cart_db.Cart) {1}
[0] (StructObject UserType) cart shopping_simple.Cart struct{CartID string, LastProductID string, TotalQuantity int, Products []string}
     --> w-tainted: write(cart_db.Cart.CartID) {1}
[_1] (FieldObject FieldType) CartID string
      --> w-tainted: write(cart_db.Cart.CartID) {1}
[__2] (BasicObject BasicType) cartID string
       --> w-tainted: write(cart_db.Cart.CartID) {1}
[___3] (Reference BasicType) ref <cartID string> @ Frontend

[0] (InterfaceObject UserType) err .error

