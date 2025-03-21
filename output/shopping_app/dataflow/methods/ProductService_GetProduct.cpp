[0] (PointerObject PointerType) s (*shopping_app.ProductServiceImpl struct{product_db NoSQLDatabase})
[_1] (StructObject UserType) shopping_app.ProductServiceImpl struct{product_db NoSQLDatabase}
[__2] (FieldObject FieldType) product_db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) product_db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ CartService
[__2] (Reference UserType) ref <ctx context.Context> @ Frontend
[_1] (Reference UserType) ref <ctx context.Context> @ OrderService
[__2] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3}       --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[0] (BasicObject BasicType) productID string
     --> r-tainted: read(product_db._.productID) {1}
[_1] (Reference BasicType) ref <productID string> @ CartService
      --> r-tainted: read(product_db._.productID) {1}
[__2] (Reference BasicType) ref <productID string> @ Frontend
     --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3}         --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[_1] (Reference BasicType) ref <productID string> @ OrderService
      --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3}           --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[__2] (Reference BasicType) ref <productID string> @ Frontend

    --> r-tainted: read(product_db.Product) {1}
[0] (StructObject UserType) product shopping_app.Product struct{ProductID string, Description string, PricePerUnit int, Category string}

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = product_database, collection = product_collection}

[0] (InterfaceObject UserType) _ .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "productID" string, Key "productID" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "productID" string
[___3] (BasicObject BasicType) "productID" string
[__2] (FieldObject FieldType) Value string
       --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3}             --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[___3] (BasicObject BasicType) productID string
        --> r-tainted: read(product_db._.productID) {1}
[____4] (Reference BasicType) ref <productID string> @ CartService
         --> r-tainted: read(product_db._.productID) {1}
[_____5] (Reference BasicType) ref <productID string> @ Frontend
        --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3}               --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[____4] (Reference BasicType) ref <productID string> @ OrderService
         --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3}                 --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[_____5] (Reference BasicType) ref <productID string> @ Frontend

    --> r-tainted: read(product_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = product_database, collection = product_collection}
     --> r-tainted: read(product_db.Product) {1}
[_1] (StructObject UserType) product shopping_app.Product struct{ProductID string, Description string, PricePerUnit int, Category string}

[0] (InterfaceObject UserType) _ .error

[0] (BasicObject BasicType) _ bool

[0] (InterfaceObject UserType) err .error

