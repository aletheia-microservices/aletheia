[0] (PointerObject PointerType) s (*shopping_app.BillingServiceImpl struct{bill_db NoSQLDatabase})
[_1] (StructObject UserType) shopping_app.BillingServiceImpl struct{bill_db NoSQLDatabase}
[__2] (FieldObject FieldType) bill_db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) bill_db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ OrderService
[__2] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> w-tainted: write(billing_db.Bill.UserID, order_db.Order.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {4}       --> w-tainted: write(billing_db.Bill.UserID, order_db.Order.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {4} --> r-tainted: read(analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {2}
[0] (BasicObject BasicType) userID string
     --> w-tainted: write(order_db.Order.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID, billing_db.Bill.UserID) {4}         --> w-tainted: write(order_db.Order.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID, billing_db.Bill.UserID) {4} --> r-tainted: read(analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {2}
[_1] (Reference BasicType) ref <userID string> @ OrderService
      --> w-tainted: write(order_db.Order.UserID, billing_db.Bill.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {4}           --> w-tainted: write(order_db.Order.UserID, billing_db.Bill.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {4} --> r-tainted: read(analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {2}
[__2] (Reference BasicType) ref <userID string> @ Frontend

    --> w-tainted: write(billing_db.Bill.ProductID, order_db.Order.ProductID, stock_db.Stock.ProductID) {3}       --> w-tainted: write(billing_db.Bill.ProductID, order_db.Order.ProductID, stock_db.Stock.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[0] (BasicObject BasicType) productID string
     --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3}         --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[_1] (Reference BasicType) ref <productID string> @ OrderService
      --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3}           --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[__2] (Reference BasicType) ref <productID string> @ Frontend

    --> w-tainted: write(billing_db.Bill.Quantity, order_db.Order.Quantity, stock_db.Stock.Quantity, billing_db.Bill.TotalCost) {4}
[0] (BasicObject BasicType) quantity int
     --> w-tainted: write(order_db.Order.Quantity, stock_db.Stock.Quantity, billing_db.Bill.Quantity) {3}
[_1] (Reference BasicType) ref <quantity int> @ OrderService
      --> w-tainted: write(order_db.Order.Quantity, stock_db.Stock.Quantity, billing_db.Bill.Quantity) {3}
[__2] (Reference BasicType) ref <quantity int> @ Frontend

    --> w-tainted: write(billing_db.Bill.PricePerUnit, billing_db.Bill.TotalCost) {2}
[0] (BasicObject BasicType) pricePerUnit int
     --> w-tainted: write(billing_db.Bill.PricePerUnit) {1}
[_1] (Reference BasicType) ref <price int> @ OrderService
      --> w-tainted: write(billing_db.Bill.PricePerUnit) {1}
[__2] (Reference BasicType) ref <price int> @ Frontend

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = bill_database, collection = bill_collection}

[0] (InterfaceObject UserType) _ .error

    --> w-tainted: write(billing_db.Bill) {1}
[0] (StructObject UserType) bill shopping_app.Bill struct{UserID string, ProductID string, Quantity int, PricePerUnit int, TotalCost }
     --> w-tainted: write(billing_db.Bill.PricePerUnit) {1}
[_1] (FieldObject FieldType) PricePerUnit int
      --> w-tainted: write(billing_db.Bill.PricePerUnit, billing_db.Bill.TotalCost) {2}
[__2] (BasicObject BasicType) pricePerUnit int
       --> w-tainted: write(billing_db.Bill.PricePerUnit) {1}
[___3] (Reference BasicType) ref <price int> @ OrderService
        --> w-tainted: write(billing_db.Bill.PricePerUnit) {1}
[____4] (Reference BasicType) ref <price int> @ Frontend
     --> w-tainted: write(billing_db.Bill.ProductID) {1}
[_1] (FieldObject FieldType) ProductID string
      --> w-tainted: write(billing_db.Bill.ProductID, order_db.Order.ProductID, stock_db.Stock.ProductID) {3}           --> w-tainted: write(billing_db.Bill.ProductID, order_db.Order.ProductID, stock_db.Stock.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[__2] (BasicObject BasicType) productID string
       --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3}             --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[___3] (Reference BasicType) ref <productID string> @ OrderService
        --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3}               --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[____4] (Reference BasicType) ref <productID string> @ Frontend
     --> w-tainted: write(billing_db.Bill.Quantity) {1}
[_1] (FieldObject FieldType) Quantity int
      --> w-tainted: write(billing_db.Bill.Quantity, order_db.Order.Quantity, stock_db.Stock.Quantity, billing_db.Bill.TotalCost) {4}
[__2] (BasicObject BasicType) quantity int
       --> w-tainted: write(order_db.Order.Quantity, stock_db.Stock.Quantity, billing_db.Bill.Quantity) {3}
[___3] (Reference BasicType) ref <quantity int> @ OrderService
        --> w-tainted: write(order_db.Order.Quantity, stock_db.Stock.Quantity, billing_db.Bill.Quantity) {3}
[____4] (Reference BasicType) ref <quantity int> @ Frontend
     --> w-tainted: write(billing_db.Bill.TotalCost) {1}
[_1] (FieldObject FieldType) TotalCost 
      --> w-tainted: write(billing_db.Bill.TotalCost) {1}
[__2] (BasicObject BasicType) 
       --> w-tainted: write(billing_db.Bill.Quantity, order_db.Order.Quantity, stock_db.Stock.Quantity, billing_db.Bill.TotalCost) {4}
[___3] (BasicObject BasicType) quantity int
        --> w-tainted: write(order_db.Order.Quantity, stock_db.Stock.Quantity, billing_db.Bill.Quantity) {3}
[____4] (Reference BasicType) ref <quantity int> @ OrderService
         --> w-tainted: write(order_db.Order.Quantity, stock_db.Stock.Quantity, billing_db.Bill.Quantity) {3}
[_____5] (Reference BasicType) ref <quantity int> @ Frontend
       --> w-tainted: write(billing_db.Bill.PricePerUnit, billing_db.Bill.TotalCost) {2}
[___3] (BasicObject BasicType) pricePerUnit int
        --> w-tainted: write(billing_db.Bill.PricePerUnit) {1}
[____4] (Reference BasicType) ref <price int> @ OrderService
         --> w-tainted: write(billing_db.Bill.PricePerUnit) {1}
[_____5] (Reference BasicType) ref <price int> @ Frontend
     --> w-tainted: write(billing_db.Bill.UserID) {1}
[_1] (FieldObject FieldType) UserID string
      --> w-tainted: write(billing_db.Bill.UserID, order_db.Order.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {4}           --> w-tainted: write(billing_db.Bill.UserID, order_db.Order.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {4} --> r-tainted: read(analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {2}
[__2] (BasicObject BasicType) userID string
       --> w-tainted: write(order_db.Order.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID, billing_db.Bill.UserID) {4}             --> w-tainted: write(order_db.Order.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID, billing_db.Bill.UserID) {4} --> r-tainted: read(analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {2}
[___3] (Reference BasicType) ref <userID string> @ OrderService
        --> w-tainted: write(order_db.Order.UserID, billing_db.Bill.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {4}               --> w-tainted: write(order_db.Order.UserID, billing_db.Bill.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {4} --> r-tainted: read(analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {2}
[____4] (Reference BasicType) ref <userID string> @ Frontend

