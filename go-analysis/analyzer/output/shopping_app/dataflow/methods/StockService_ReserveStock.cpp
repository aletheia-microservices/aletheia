[0] (PointerObject PointerType) s (*shopping_app.StockServiceImpl struct{stock_db NoSQLDatabase})
[_1] (StructObject UserType) shopping_app.StockServiceImpl struct{stock_db NoSQLDatabase}
[__2] (FieldObject FieldType) stock_db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) stock_db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ OrderService
[__2] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> w-tainted: write(stock_db.Stock.ProductID, order_db.Order.ProductID, billing_db.Bill.ProductID) {3}       --> w-tainted: write(stock_db.Stock.ProductID, order_db.Order.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[0] (BasicObject BasicType) productID string
     --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3}         --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[_1] (Reference BasicType) ref <productID string> @ OrderService
      --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3}           --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[__2] (Reference BasicType) ref <productID string> @ Frontend

    --> w-tainted: write(stock_db.Stock.Quantity, order_db.Order.Quantity, billing_db.Bill.TotalCost) {3}
[0] (BasicObject BasicType) quantity int
     --> w-tainted: write(order_db.Order.Quantity, stock_db.Stock.Quantity, billing_db.Bill.TotalCost) {3}
[_1] (Reference BasicType) ref <quantity int> @ OrderService
      --> w-tainted: write(order_db.Order.Quantity, stock_db.Stock.Quantity, billing_db.Bill.TotalCost) {3}
[__2] (Reference BasicType) ref <quantity int> @ Frontend

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = stock_database, collection = stock_collection}

[0] (InterfaceObject UserType) _ .error

    --> w-tainted: write(stock_db.Stock) {1}
[0] (StructObject UserType) stock shopping_app.Stock struct{ProductID string, Quantity int}
     --> w-tainted: write(stock_db.Stock.ProductID) {1}
[_1] (FieldObject FieldType) ProductID string
      --> w-tainted: write(stock_db.Stock.ProductID, order_db.Order.ProductID, billing_db.Bill.ProductID) {3}           --> w-tainted: write(stock_db.Stock.ProductID, order_db.Order.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[__2] (BasicObject BasicType) productID string
       --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3}             --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[___3] (Reference BasicType) ref <productID string> @ OrderService
        --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3}               --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[____4] (Reference BasicType) ref <productID string> @ Frontend
     --> w-tainted: write(stock_db.Stock.Quantity) {1}
[_1] (FieldObject FieldType) Quantity int
      --> w-tainted: write(stock_db.Stock.Quantity, order_db.Order.Quantity, billing_db.Bill.TotalCost) {3}
[__2] (BasicObject BasicType) quantity int
       --> w-tainted: write(order_db.Order.Quantity, stock_db.Stock.Quantity, billing_db.Bill.TotalCost) {3}
[___3] (Reference BasicType) ref <quantity int> @ OrderService
        --> w-tainted: write(order_db.Order.Quantity, stock_db.Stock.Quantity, billing_db.Bill.TotalCost) {3}
[____4] (Reference BasicType) ref <quantity int> @ Frontend

