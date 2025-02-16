[0] (PointerObject PointerType) s (*shopping_simple.ProductServiceImpl struct{product_db NoSQLDatabase, product_queue Queue, num_workers int})
[_1] (StructObject UserType) shopping_simple.ProductServiceImpl struct{product_db NoSQLDatabase, product_queue Queue, num_workers int}
[__2] (FieldObject FieldType) product_db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) product_db NoSQLDatabase
[__2] (FieldObject FieldType) product_queue Queue
[___3] (BlueprintBackendObject BlueprintBackendType) product_queue Queue

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend
[_1] (Reference UserType) ref <ctx context.Context> @ CartService
[__2] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> w-tainted: write(cart_db.Cart.LastProductID) {1}       --> w-tainted: write(cart_db.Cart.LastProductID) {1} --> r-tainted: read(product_db.Product.ProductID) {1}
[0] (BasicObject BasicType) productID string
     --> r-tainted: read(product_db.Product.ProductID) {1}
[_1] (Reference BasicType) ref <productID string> @ Frontend
     --> w-tainted: write(cart_db.Cart.LastProductID, cart_db.Cart.Products) {2}         --> w-tainted: write(cart_db.Cart.LastProductID, cart_db.Cart.Products) {2} --> r-tainted: read(product_db.Product.ProductID) {1}
[_1] (Reference BasicType) ref <productID string> @ CartService
      --> w-tainted: write(cart_db.Cart.LastProductID) {1}           --> w-tainted: write(cart_db.Cart.LastProductID) {1} --> r-tainted: read(product_db.Product.ProductID) {1}
[__2] (Reference BasicType) ref <productID string> @ Frontend

    --> r-tainted: read(product_db.Product) {1}
[0] (StructObject UserType) product shopping_simple.Product struct{ProductID string, Description string, PricePerUnit int, Category string}

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = product_database, collection = product_collection}

[0] (InterfaceObject UserType) _ .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "productid" string, Key "productid" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "productid" string
[___3] (BasicObject BasicType) "productid" string
[__2] (FieldObject FieldType) Value string
       --> w-tainted: write(cart_db.Cart.LastProductID) {1}             --> w-tainted: write(cart_db.Cart.LastProductID) {1} --> r-tainted: read(product_db.Product.ProductID) {1}
[___3] (BasicObject BasicType) productID string
        --> r-tainted: read(product_db.Product.ProductID) {1}
[____4] (Reference BasicType) ref <productID string> @ Frontend
        --> w-tainted: write(cart_db.Cart.LastProductID, cart_db.Cart.Products) {2}               --> w-tainted: write(cart_db.Cart.LastProductID, cart_db.Cart.Products) {2} --> r-tainted: read(product_db.Product.ProductID) {1}
[____4] (Reference BasicType) ref <productID string> @ CartService
         --> w-tainted: write(cart_db.Cart.LastProductID) {1}                 --> w-tainted: write(cart_db.Cart.LastProductID) {1} --> r-tainted: read(product_db.Product.ProductID) {1}
[_____5] (Reference BasicType) ref <productID string> @ Frontend

    --> r-tainted: read(product_db.Product) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = product_database, collection = product_collection}
     --> r-tainted: read(product_db.Product) {1}
[_1] (StructObject UserType) product shopping_simple.Product struct{ProductID string, Description string, PricePerUnit int, Category string}

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

