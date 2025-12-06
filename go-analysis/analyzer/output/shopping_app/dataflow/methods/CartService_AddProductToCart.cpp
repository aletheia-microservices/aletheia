[0] (PointerObject PointerType) c (*shopping_app.CartServiceImpl struct{order_service shopping_app.OrderService, product_service shopping_app.ProductService, cart_db NoSQLDatabase})
[_1] (StructObject UserType) shopping_app.CartServiceImpl struct{order_service shopping_app.OrderService, product_service shopping_app.ProductService, cart_db NoSQLDatabase}
[__2] (FieldObject FieldType) cart_db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) cart_db NoSQLDatabase
[__2] (FieldObject FieldType) order_service shopping_app.OrderService
[___3] (ServiceObject ServiceType) order_service shopping_app.OrderService
[__2] (FieldObject FieldType) product_service shopping_app.ProductService
[___3] (ServiceObject ServiceType) product_service shopping_app.ProductService

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> r-tainted: read(cart_db._.cartID) {1}
[0] (BasicObject BasicType) cartID string
     --> r-tainted: read(cart_db._.cartID) {1}
[_1] (Reference BasicType) ref <cartID string> @ Frontend

    --> r-tainted: read(product_db._.productID) {1}
[0] (BasicObject BasicType) productID string
     --> r-tainted: read(product_db._.productID) {1}
[_1] (Reference BasicType) ref <productID string> @ Frontend

[0] (StructObject UserType) product shopping_app.Product struct{ProductID string, Description string, PricePerUnit int, Category string}

[0] (StructObject UserType) product shopping_app.Product struct{ProductID string, Description string, PricePerUnit int, Category string}
     --> r-tainted: read(product_db.Product) {1}
[_1] (Reference UserType) ref <product shopping_app.Product struct{ProductID string, Description string, PricePerUnit int, Category string}> @ ProductService
[_1] (FieldObject FieldType) PricePerUnit int
      --> w-tainted: write(cart_db.Cart.PricePerUnit) {1}           --> w-tainted: write(cart_db.Cart.PricePerUnit) {1} --> r-tainted: read(cart_db.Cart.PricePerUnit) {1}
[__2] (BasicObject BasicType) PricePerUnit int
[_1] (FieldObject FieldType) ProductID string
      --> w-tainted: write(cart_db.Cart.ProductID) {1}           --> w-tainted: write(cart_db.Cart.ProductID) {1} --> r-tainted: read(cart_db.Cart.ProductID) {1}
[__2] (BasicObject BasicType) ProductID string

[0] (InterfaceObject UserType) _ .error
[_1] (Reference UserType) ref <err .error> @ ProductService

    --> w-tainted: write(cart_db.Cart) {1}       --> w-tainted: write(cart_db.Cart) {1} --> r-tainted: read(cart_db.Cart, cart_db.Cart.PricePerUnit, cart_db.Cart.ProductID, cart_db.Cart.Quantity) {4}
[0] (StructObject UserType) cart shopping_app.Cart struct{CartID string, UserID string, ProductID string, Quantity int, PricePerUnit int, Status string}
     --> w-tainted: write(cart_db.Cart.PricePerUnit) {1}         --> w-tainted: write(cart_db.Cart.PricePerUnit) {1} --> r-tainted: read(cart_db.Cart.PricePerUnit) {1}
[_1] (FieldObject FieldType) PricePerUnit int
      --> w-tainted: write(cart_db.Cart.PricePerUnit) {1}           --> w-tainted: write(cart_db.Cart.PricePerUnit) {1} --> r-tainted: read(cart_db.Cart.PricePerUnit) {1}
[__2] (BasicObject BasicType) PricePerUnit int
     --> w-tainted: write(cart_db.Cart.ProductID) {1}         --> w-tainted: write(cart_db.Cart.ProductID) {1} --> r-tainted: read(cart_db.Cart.ProductID) {1}
[_1] (FieldObject FieldType) ProductID string
      --> w-tainted: write(cart_db.Cart.ProductID) {1}           --> w-tainted: write(cart_db.Cart.ProductID) {1} --> r-tainted: read(cart_db.Cart.ProductID) {1}
[__2] (BasicObject BasicType) ProductID string
     --> w-tainted: write(cart_db.Cart.Quantity) {1}         --> w-tainted: write(cart_db.Cart.Quantity) {1} --> r-tainted: read(cart_db.Cart.Quantity) {1}
[_1] (FieldObject FieldType) Quantity int
      --> w-tainted: write(cart_db.Cart.Quantity) {1}           --> w-tainted: write(cart_db.Cart.Quantity) {1} --> r-tainted: read(cart_db.Cart.Quantity) {1}
[__2] (BasicObject BasicType) 1 int

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = cart_database, collection = cart_collection}

[0] (InterfaceObject UserType) _ .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "cartID" string, Key "cartID" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "cartID" string
[___3] (BasicObject BasicType) "cartID" string
[__2] (FieldObject FieldType) Value string
       --> r-tainted: read(cart_db._.cartID) {1}
[___3] (BasicObject BasicType) cartID string
        --> r-tainted: read(cart_db._.cartID) {1}
[____4] (Reference BasicType) ref <cartID string> @ Frontend

    --> r-tainted: read(cart_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = cart_database, collection = cart_collection}
     --> w-tainted: write(cart_db.Cart) {1}         --> w-tainted: write(cart_db.Cart) {1} --> r-tainted: read(cart_db.Cart, cart_db.Cart.PricePerUnit, cart_db.Cart.ProductID, cart_db.Cart.Quantity) {4}
[_1] (StructObject UserType) cart shopping_app.Cart struct{CartID string, UserID string, ProductID string, Quantity int, PricePerUnit int, Status string}
      --> w-tainted: write(cart_db.Cart.PricePerUnit) {1}           --> w-tainted: write(cart_db.Cart.PricePerUnit) {1} --> r-tainted: read(cart_db.Cart.PricePerUnit) {1}
[__2] (FieldObject FieldType) PricePerUnit int
       --> w-tainted: write(cart_db.Cart.PricePerUnit) {1}             --> w-tainted: write(cart_db.Cart.PricePerUnit) {1} --> r-tainted: read(cart_db.Cart.PricePerUnit) {1}
[___3] (BasicObject BasicType) PricePerUnit int
      --> w-tainted: write(cart_db.Cart.ProductID) {1}           --> w-tainted: write(cart_db.Cart.ProductID) {1} --> r-tainted: read(cart_db.Cart.ProductID) {1}
[__2] (FieldObject FieldType) ProductID string
       --> w-tainted: write(cart_db.Cart.ProductID) {1}             --> w-tainted: write(cart_db.Cart.ProductID) {1} --> r-tainted: read(cart_db.Cart.ProductID) {1}
[___3] (BasicObject BasicType) ProductID string
      --> w-tainted: write(cart_db.Cart.Quantity) {1}           --> w-tainted: write(cart_db.Cart.Quantity) {1} --> r-tainted: read(cart_db.Cart.Quantity) {1}
[__2] (FieldObject FieldType) Quantity int
       --> w-tainted: write(cart_db.Cart.Quantity) {1}             --> w-tainted: write(cart_db.Cart.Quantity) {1} --> r-tainted: read(cart_db.Cart.Quantity) {1}
[___3] (BasicObject BasicType) 1 int

[0] (InterfaceObject UserType) _ .error

[0] (PointerObject PointerType) c (*shopping_app.CartServiceImpl struct{order_service shopping_app.OrderService, product_service shopping_app.ProductService, cart_db NoSQLDatabase})
[_1] (StructObject UserType) shopping_app.CartServiceImpl struct{order_service shopping_app.OrderService, product_service shopping_app.ProductService, cart_db NoSQLDatabase}
[__2] (FieldObject FieldType) cart_db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) cart_db NoSQLDatabase
[__2] (FieldObject FieldType) order_service shopping_app.OrderService
[___3] (ServiceObject ServiceType) order_service shopping_app.OrderService
[__2] (FieldObject FieldType) product_service shopping_app.ProductService
[___3] (ServiceObject ServiceType) product_service shopping_app.ProductService

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> r-tainted: read(cart_db._.cartID) {1}
[0] (BasicObject BasicType) cartID string
     --> r-tainted: read(cart_db._.cartID) {1}
[_1] (Reference BasicType) ref <cartID string> @ Frontend

    --> r-tainted: read(product_db._.productID) {1}
[0] (BasicObject BasicType) productID string
     --> r-tainted: read(product_db._.productID) {1}
[_1] (Reference BasicType) ref <productID string> @ Frontend

[0] (StructObject UserType) product shopping_app.Product struct{ProductID string, Description string, PricePerUnit int, Category string}

[0] (StructObject UserType) product shopping_app.Product struct{ProductID string, Description string, PricePerUnit int, Category string}
     --> r-tainted: read(product_db.Product) {1}
[_1] (Reference UserType) ref <product shopping_app.Product struct{ProductID string, Description string, PricePerUnit int, Category string}> @ ProductService
[_1] (FieldObject FieldType) PricePerUnit int
      --> w-tainted: write(cart_db.Cart.PricePerUnit) {1}           --> w-tainted: write(cart_db.Cart.PricePerUnit) {1} --> r-tainted: read(cart_db.Cart.PricePerUnit) {1}
[__2] (BasicObject BasicType) PricePerUnit int
[_1] (FieldObject FieldType) ProductID string
      --> w-tainted: write(cart_db.Cart.ProductID) {1}           --> w-tainted: write(cart_db.Cart.ProductID) {1} --> r-tainted: read(cart_db.Cart.ProductID) {1}
[__2] (BasicObject BasicType) ProductID string

[0] (InterfaceObject UserType) _ .error
[_1] (Reference UserType) ref <err .error> @ ProductService

    --> w-tainted: write(cart_db.Cart) {1}       --> w-tainted: write(cart_db.Cart) {1} --> r-tainted: read(cart_db.Cart, cart_db.Cart.PricePerUnit, cart_db.Cart.ProductID, cart_db.Cart.Quantity) {4}
[0] (StructObject UserType) cart shopping_app.Cart struct{CartID string, UserID string, ProductID string, Quantity int, PricePerUnit int, Status string}
     --> w-tainted: write(cart_db.Cart.PricePerUnit) {1}         --> w-tainted: write(cart_db.Cart.PricePerUnit) {1} --> r-tainted: read(cart_db.Cart.PricePerUnit) {1}
[_1] (FieldObject FieldType) PricePerUnit int
      --> w-tainted: write(cart_db.Cart.PricePerUnit) {1}           --> w-tainted: write(cart_db.Cart.PricePerUnit) {1} --> r-tainted: read(cart_db.Cart.PricePerUnit) {1}
[__2] (BasicObject BasicType) PricePerUnit int
     --> w-tainted: write(cart_db.Cart.ProductID) {1}         --> w-tainted: write(cart_db.Cart.ProductID) {1} --> r-tainted: read(cart_db.Cart.ProductID) {1}
[_1] (FieldObject FieldType) ProductID string
      --> w-tainted: write(cart_db.Cart.ProductID) {1}           --> w-tainted: write(cart_db.Cart.ProductID) {1} --> r-tainted: read(cart_db.Cart.ProductID) {1}
[__2] (BasicObject BasicType) ProductID string
     --> w-tainted: write(cart_db.Cart.Quantity) {1}         --> w-tainted: write(cart_db.Cart.Quantity) {1} --> r-tainted: read(cart_db.Cart.Quantity) {1}
[_1] (FieldObject FieldType) Quantity int
      --> w-tainted: write(cart_db.Cart.Quantity) {1}           --> w-tainted: write(cart_db.Cart.Quantity) {1} --> r-tainted: read(cart_db.Cart.Quantity) {1}
[__2] (BasicObject BasicType) 1 int

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = cart_database, collection = cart_collection}

[0] (InterfaceObject UserType) _ .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "cartID" string, Key "cartID" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "cartID" string
[___3] (BasicObject BasicType) "cartID" string
[__2] (FieldObject FieldType) Value string
       --> r-tainted: read(cart_db._.cartID) {1}
[___3] (BasicObject BasicType) cartID string
        --> r-tainted: read(cart_db._.cartID) {1}
[____4] (Reference BasicType) ref <cartID string> @ Frontend

    --> r-tainted: read(cart_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = cart_database, collection = cart_collection}
     --> w-tainted: write(cart_db.Cart) {1}         --> w-tainted: write(cart_db.Cart) {1} --> r-tainted: read(cart_db.Cart, cart_db.Cart.PricePerUnit, cart_db.Cart.ProductID, cart_db.Cart.Quantity) {4}
[_1] (StructObject UserType) cart shopping_app.Cart struct{CartID string, UserID string, ProductID string, Quantity int, PricePerUnit int, Status string}
      --> w-tainted: write(cart_db.Cart.PricePerUnit) {1}           --> w-tainted: write(cart_db.Cart.PricePerUnit) {1} --> r-tainted: read(cart_db.Cart.PricePerUnit) {1}
[__2] (FieldObject FieldType) PricePerUnit int
       --> w-tainted: write(cart_db.Cart.PricePerUnit) {1}             --> w-tainted: write(cart_db.Cart.PricePerUnit) {1} --> r-tainted: read(cart_db.Cart.PricePerUnit) {1}
[___3] (BasicObject BasicType) PricePerUnit int
      --> w-tainted: write(cart_db.Cart.ProductID) {1}           --> w-tainted: write(cart_db.Cart.ProductID) {1} --> r-tainted: read(cart_db.Cart.ProductID) {1}
[__2] (FieldObject FieldType) ProductID string
       --> w-tainted: write(cart_db.Cart.ProductID) {1}             --> w-tainted: write(cart_db.Cart.ProductID) {1} --> r-tainted: read(cart_db.Cart.ProductID) {1}
[___3] (BasicObject BasicType) ProductID string
      --> w-tainted: write(cart_db.Cart.Quantity) {1}           --> w-tainted: write(cart_db.Cart.Quantity) {1} --> r-tainted: read(cart_db.Cart.Quantity) {1}
[__2] (FieldObject FieldType) Quantity int
       --> w-tainted: write(cart_db.Cart.Quantity) {1}             --> w-tainted: write(cart_db.Cart.Quantity) {1} --> r-tainted: read(cart_db.Cart.Quantity) {1}
[___3] (BasicObject BasicType) 1 int

[0] (InterfaceObject UserType) _ .error

