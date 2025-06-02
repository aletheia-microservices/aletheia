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

    --> w-tainted: write(cart_db.Cart.CartID) {1}       --> w-tainted: write(cart_db.Cart.CartID) {1} --> r-tainted: read(cart_db.Cart.CartID) {1}
[0] (BasicObject BasicType) cartID string
     --> w-tainted: write(cart_db.Cart.CartID) {1}         --> w-tainted: write(cart_db.Cart.CartID) {1} --> r-tainted: read(cart_db.Cart.CartID) {1}
[_1] (Reference BasicType) ref <cartID string> @ Frontend

    --> w-tainted: write(cart_db.Cart.LastProductID, cart_db.Cart.Products) {2}       --> w-tainted: write(cart_db.Cart.LastProductID, cart_db.Cart.Products) {2} --> r-tainted: read(product_db.Product.ProductID) {1}
[0] (BasicObject BasicType) productID string
     --> w-tainted: write(cart_db.Cart.LastProductID) {1}         --> w-tainted: write(cart_db.Cart.LastProductID) {1} --> r-tainted: read(product_db.Product.ProductID) {1}
[_1] (Reference BasicType) ref <productID string> @ Frontend

[0] (StructObject UserType) product shopping_simple.Product struct{ProductID string, Description string, PricePerUnit int, Category string}

[0] (StructObject UserType) product shopping_simple.Product struct{ProductID string, Description string, PricePerUnit int, Category string}
     --> r-tainted: read(product_db.Product) {1}
[_1] (Reference UserType) ref <product shopping_simple.Product struct{ProductID string, Description string, PricePerUnit int, Category string}> @ ProductService
[_1] (FieldObject FieldType) PricePerUnit int
[__2] (BasicObject BasicType) PricePerUnit int

[0] (InterfaceObject UserType) _ .error
[_1] (Reference UserType) ref <err .error> @ ProductService

    --> r-tainted: read(cart_db.Cart) {1}
[0] (StructObject UserType) cart shopping_simple.Cart struct{CartID string, LastProductID string, TotalQuantity int, Products []string}

[0] (StructObject UserType) cartProduct shopping_simple.CartProduct struct{CartID string, ProductID string, Quantity int, PricePerUnit int}

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = cart_database, collection = cart_collection}

[0] (InterfaceObject UserType) _ .error

[0] (SliceObject UserType) filter primitive.D
[_1] (StructObject StructType) struct{Key "cartid" string, Key "cartid" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "cartid" string
[___3] (BasicObject BasicType) "cartid" string
[__2] (FieldObject FieldType) Value string
       --> w-tainted: write(cart_db.Cart.CartID) {1}             --> w-tainted: write(cart_db.Cart.CartID) {1} --> r-tainted: read(cart_db.Cart.CartID) {1}
[___3] (BasicObject BasicType) cartID string
        --> w-tainted: write(cart_db.Cart.CartID) {1}               --> w-tainted: write(cart_db.Cart.CartID) {1} --> r-tainted: read(cart_db.Cart.CartID) {1}
[____4] (Reference BasicType) ref <cartID string> @ Frontend

    --> r-tainted: read(cart_db.Cart) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = cart_database, collection = cart_collection}
     --> r-tainted: read(cart_db.Cart) {1}
[_1] (StructObject UserType) cart shopping_simple.Cart struct{CartID string, LastProductID string, TotalQuantity int, Products []string}

[0] (InterfaceObject UserType) err .error

[0] (StructObject UserType) cartProduct shopping_simple.CartProduct struct{CartID string, ProductID string, Quantity 1 int, PricePerUnit int}
[_1] (FieldObject FieldType) CartID string
      --> w-tainted: write(cart_db.Cart.CartID) {1}           --> w-tainted: write(cart_db.Cart.CartID) {1} --> r-tainted: read(cart_db.Cart.CartID) {1}
[__2] (BasicObject BasicType) cartID string
       --> w-tainted: write(cart_db.Cart.CartID) {1}             --> w-tainted: write(cart_db.Cart.CartID) {1} --> r-tainted: read(cart_db.Cart.CartID) {1}
[___3] (Reference BasicType) ref <cartID string> @ Frontend
[_1] (FieldObject FieldType) PricePerUnit int
[__2] (BasicObject BasicType) PricePerUnit int
[_1] (FieldObject FieldType) ProductID string
      --> w-tainted: write(cart_db.Cart.LastProductID, cart_db.Cart.Products) {2}           --> w-tainted: write(cart_db.Cart.LastProductID, cart_db.Cart.Products) {2} --> r-tainted: read(product_db.Product.ProductID) {1}
[__2] (BasicObject BasicType) productID string
       --> w-tainted: write(cart_db.Cart.LastProductID) {1}             --> w-tainted: write(cart_db.Cart.LastProductID) {1} --> r-tainted: read(product_db.Product.ProductID) {1}
[___3] (Reference BasicType) ref <productID string> @ Frontend
[_1] (FieldObject FieldType) Quantity 1 int
[__2] (BasicObject BasicType) 1 int

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) _ .error

