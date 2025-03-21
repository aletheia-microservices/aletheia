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

    --> r-tainted: read(cart_db.Cart.CartID) {1}
[0] (BasicObject BasicType) cartID string
     --> r-tainted: read(cart_db.Cart.CartID) {1}
[_1] (Reference BasicType) ref <cartID string> @ Frontend

    --> r-tainted: read(cart_db.Cart, cart_db.Cart.Products) {2}
[0] (StructObject UserType) cart shopping_simple.Cart struct{CartID string, LastProductID string, TotalQuantity int, Products []string}
     --> r-tainted: read(cart_db.Cart.Products) {1}
[_1] (FieldObject FieldType) Products []string
      --> r-tainted: read(product_db.Product.ProductID, cart_db.Cart.Products) {2}
[__2] (SliceObject SliceType) Products []string

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = cart_database, collection = cart_database}

[0] (InterfaceObject UserType) _ .error

[0] (SliceObject UserType) filter primitive.D
[_1] (StructObject StructType) struct{Key "cartid" string, Key "cartid" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "cartid" string
[___3] (BasicObject BasicType) "cartid" string
[__2] (FieldObject FieldType) Value string
       --> r-tainted: read(cart_db.Cart.CartID) {1}
[___3] (BasicObject BasicType) cartID string
        --> r-tainted: read(cart_db.Cart.CartID) {1}
[____4] (Reference BasicType) ref <cartID string> @ Frontend

    --> r-tainted: read(cart_db.Cart) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = cart_database, collection = cart_database}
     --> r-tainted: read(cart_db.Cart, cart_db.Cart.Products) {2}
[_1] (StructObject UserType) cart shopping_simple.Cart struct{CartID string, LastProductID string, TotalQuantity int, Products []string}
      --> r-tainted: read(cart_db.Cart.Products) {1}
[__2] (FieldObject FieldType) Products []string
       --> r-tainted: read(product_db.Product.ProductID, cart_db.Cart.Products) {2}
[___3] (SliceObject SliceType) Products []string

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(product_db.Product.ProductID) {1}
[0] (BasicObject BasicType) p string
     --> r-tainted: read(product_db.Product.ProductID, cart_db.Cart.Products) {2}
[_1] (SliceObject SliceType) Products []string

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

    --> r-tainted: read(cart_db.Cart.CartID) {1}
[0] (BasicObject BasicType) cartID string
     --> r-tainted: read(cart_db.Cart.CartID) {1}
[_1] (Reference BasicType) ref <cartID string> @ Frontend

    --> r-tainted: read(cart_db.Cart, cart_db.Cart.Products) {2}
[0] (StructObject UserType) cart shopping_simple.Cart struct{CartID string, LastProductID string, TotalQuantity int, Products []string}
     --> r-tainted: read(cart_db.Cart.Products) {1}
[_1] (FieldObject FieldType) Products []string
      --> r-tainted: read(product_db.Product.ProductID, cart_db.Cart.Products) {2}
[__2] (SliceObject SliceType) Products []string

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = cart_database, collection = cart_database}

[0] (InterfaceObject UserType) _ .error

[0] (SliceObject UserType) filter primitive.D
[_1] (StructObject StructType) struct{Key "cartid" string, Key "cartid" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "cartid" string
[___3] (BasicObject BasicType) "cartid" string
[__2] (FieldObject FieldType) Value string
       --> r-tainted: read(cart_db.Cart.CartID) {1}
[___3] (BasicObject BasicType) cartID string
        --> r-tainted: read(cart_db.Cart.CartID) {1}
[____4] (Reference BasicType) ref <cartID string> @ Frontend

    --> r-tainted: read(cart_db.Cart) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = cart_database, collection = cart_database}
     --> r-tainted: read(cart_db.Cart, cart_db.Cart.Products) {2}
[_1] (StructObject UserType) cart shopping_simple.Cart struct{CartID string, LastProductID string, TotalQuantity int, Products []string}
      --> r-tainted: read(cart_db.Cart.Products) {1}
[__2] (FieldObject FieldType) Products []string
       --> r-tainted: read(product_db.Product.ProductID, cart_db.Cart.Products) {2}
[___3] (SliceObject SliceType) Products []string

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(product_db.Product.ProductID) {1}
[0] (BasicObject BasicType) p string
     --> r-tainted: read(product_db.Product.ProductID, cart_db.Cart.Products) {2}
[_1] (SliceObject SliceType) Products []string

