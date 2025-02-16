[0] (PointerObject PointerType) c (*shopping_app.OrderServiceImpl struct{stock_service shopping_app.StockService, billing_service shopping_app.BillingService, product_service shopping_app.ProductService, order_db NoSQLDatabase, shipment_queue Queue, analytics_queue Queue})
[_1] (StructObject UserType) shopping_app.OrderServiceImpl struct{stock_service shopping_app.StockService, billing_service shopping_app.BillingService, product_service shopping_app.ProductService, order_db NoSQLDatabase, shipment_queue Queue, analytics_queue Queue}
[__2] (FieldObject FieldType) analytics_queue Queue
[___3] (BlueprintBackendObject BlueprintBackendType) analytics_queue Queue
[__2] (FieldObject FieldType) billing_service shopping_app.BillingService
[___3] (ServiceObject ServiceType) billing_service shopping_app.BillingService
[__2] (FieldObject FieldType) order_db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) order_db NoSQLDatabase
[__2] (FieldObject FieldType) product_service shopping_app.ProductService
[___3] (ServiceObject ServiceType) product_service shopping_app.ProductService
[__2] (FieldObject FieldType) shipment_queue Queue
[___3] (BlueprintBackendObject BlueprintBackendType) shipment_queue Queue
[__2] (FieldObject FieldType) stock_service shopping_app.StockService
[___3] (ServiceObject ServiceType) stock_service shopping_app.StockService

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> w-tainted: write(order_db.Order.OrderID, shipment_queue.ShipmentMessage.OrderID) {2}       --> w-tainted: write(order_db.Order.OrderID, shipment_queue.ShipmentMessage.OrderID) {2} --> r-tainted: read(shipment_queue.ShipmentMessage.OrderID) {1}
[0] (BasicObject BasicType) cartID string
     --> w-tainted: write(order_db.Order.OrderID, shipment_queue.ShipmentMessage.OrderID) {2}         --> w-tainted: write(order_db.Order.OrderID, shipment_queue.ShipmentMessage.OrderID) {2} --> r-tainted: read(shipment_queue.ShipmentMessage.OrderID) {1}
[_1] (Reference BasicType) ref <cartID string> @ Frontend

    --> w-tainted: write(order_db.Order.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID, billing_db.Bill.UserID) {4}       --> w-tainted: write(order_db.Order.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID, billing_db.Bill.UserID) {4} --> r-tainted: read(analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {2}
[0] (BasicObject BasicType) userID string
     --> w-tainted: write(order_db.Order.UserID, billing_db.Bill.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {4}         --> w-tainted: write(order_db.Order.UserID, billing_db.Bill.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {4} --> r-tainted: read(analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {2}
[_1] (Reference BasicType) ref <userID string> @ Frontend

    --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3}       --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[0] (BasicObject BasicType) productID string
     --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3}         --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[_1] (Reference BasicType) ref <productID string> @ Frontend

    --> w-tainted: write(order_db.Order.Quantity, stock_db.Stock.Quantity, billing_db.Bill.Quantity) {3}
[0] (BasicObject BasicType) quantity int
     --> w-tainted: write(order_db.Order.Quantity, stock_db.Stock.Quantity, billing_db.Bill.Quantity) {3}
[_1] (Reference BasicType) ref <quantity int> @ Frontend

    --> w-tainted: write(billing_db.Bill.PricePerUnit) {1}
[0] (BasicObject BasicType) price int
     --> w-tainted: write(billing_db.Bill.PricePerUnit) {1}
[_1] (Reference BasicType) ref <price int> @ Frontend

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = order_database, collection = order_collection}

[0] (InterfaceObject UserType) _ .error

    --> w-tainted: write(order_db.Order) {1}
[0] (StructObject UserType) order shopping_app.Order struct{OrderID string, UserID string, ProductID string, Quantity int, Timestamp 1 int}
     --> w-tainted: write(order_db.Order.OrderID) {1}
[_1] (FieldObject FieldType) OrderID string
      --> w-tainted: write(order_db.Order.OrderID, shipment_queue.ShipmentMessage.OrderID) {2}           --> w-tainted: write(order_db.Order.OrderID, shipment_queue.ShipmentMessage.OrderID) {2} --> r-tainted: read(shipment_queue.ShipmentMessage.OrderID) {1}
[__2] (BasicObject BasicType) cartID string
       --> w-tainted: write(order_db.Order.OrderID, shipment_queue.ShipmentMessage.OrderID) {2}             --> w-tainted: write(order_db.Order.OrderID, shipment_queue.ShipmentMessage.OrderID) {2} --> r-tainted: read(shipment_queue.ShipmentMessage.OrderID) {1}
[___3] (Reference BasicType) ref <cartID string> @ Frontend
     --> w-tainted: write(order_db.Order.ProductID) {1}
[_1] (FieldObject FieldType) ProductID string
      --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3}           --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[__2] (BasicObject BasicType) productID string
       --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3}             --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[___3] (Reference BasicType) ref <productID string> @ Frontend
     --> w-tainted: write(order_db.Order.Quantity) {1}
[_1] (FieldObject FieldType) Quantity int
      --> w-tainted: write(order_db.Order.Quantity, stock_db.Stock.Quantity, billing_db.Bill.Quantity) {3}
[__2] (BasicObject BasicType) quantity int
       --> w-tainted: write(order_db.Order.Quantity, stock_db.Stock.Quantity, billing_db.Bill.Quantity) {3}
[___3] (Reference BasicType) ref <quantity int> @ Frontend
     --> w-tainted: write(order_db.Order.Timestamp) {1}
[_1] (FieldObject FieldType) Timestamp 1 int
      --> w-tainted: write(order_db.Order.Timestamp) {1}
[__2] (BasicObject BasicType) 1 int
     --> w-tainted: write(order_db.Order.UserID) {1}
[_1] (FieldObject FieldType) UserID string
      --> w-tainted: write(order_db.Order.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID, billing_db.Bill.UserID) {4}           --> w-tainted: write(order_db.Order.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID, billing_db.Bill.UserID) {4} --> r-tainted: read(analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {2}
[__2] (BasicObject BasicType) userID string
       --> w-tainted: write(order_db.Order.UserID, billing_db.Bill.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {4}             --> w-tainted: write(order_db.Order.UserID, billing_db.Bill.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {4} --> r-tainted: read(analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {2}
[___3] (Reference BasicType) ref <userID string> @ Frontend

[0] (StructObject UserType) product shopping_app.Product struct{ProductID string, Description string, PricePerUnit int, Category string}

[0] (StructObject UserType) product shopping_app.Product struct{ProductID string, Description string, PricePerUnit int, Category string}
     --> r-tainted: read(product_db.Product) {1}
[_1] (Reference UserType) ref <product shopping_app.Product struct{ProductID string, Description string, PricePerUnit int, Category string}> @ ProductService
[_1] (FieldObject FieldType) Category string
      --> w-tainted: write(analytics_queue.AnalyticsMessage.ProductCategory) {1}           --> w-tainted: write(analytics_queue.AnalyticsMessage.ProductCategory) {1} --> r-tainted: read(analytics_queue.AnalyticsMessage.ProductCategory) {1}
[__2] (BasicObject BasicType) Category string

[0] (InterfaceObject UserType) _ .error
[_1] (Reference UserType) ref <err .error> @ ProductService

    --> w-tainted: write(analytics_queue.AnalyticsMessage) {1}       --> w-tainted: write(analytics_queue.AnalyticsMessage) {1} --> r-tainted: read(analytics_queue.AnalyticsMessage) {1}
[0] (StructObject UserType) analyticsMessage shopping_app.AnalyticsMessage struct{UserID string, ProductCategory string}
     --> w-tainted: write(analytics_queue.AnalyticsMessage.ProductCategory) {1}         --> w-tainted: write(analytics_queue.AnalyticsMessage.ProductCategory) {1} --> r-tainted: read(analytics_queue.AnalyticsMessage.ProductCategory) {1}
[_1] (FieldObject FieldType) ProductCategory string
      --> w-tainted: write(analytics_queue.AnalyticsMessage.ProductCategory) {1}           --> w-tainted: write(analytics_queue.AnalyticsMessage.ProductCategory) {1} --> r-tainted: read(analytics_queue.AnalyticsMessage.ProductCategory) {1}
[__2] (BasicObject BasicType) Category string
     --> w-tainted: write(analytics_queue.AnalyticsMessage.UserID) {1}         --> w-tainted: write(analytics_queue.AnalyticsMessage.UserID) {1} --> r-tainted: read(analytics_queue.AnalyticsMessage.UserID) {1}
[_1] (FieldObject FieldType) UserID string
      --> w-tainted: write(order_db.Order.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID, billing_db.Bill.UserID) {4}           --> w-tainted: write(order_db.Order.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID, billing_db.Bill.UserID) {4} --> r-tainted: read(analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {2}
[__2] (BasicObject BasicType) userID string
       --> w-tainted: write(order_db.Order.UserID, billing_db.Bill.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {4}             --> w-tainted: write(order_db.Order.UserID, billing_db.Bill.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {4} --> r-tainted: read(analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {2}
[___3] (Reference BasicType) ref <userID string> @ Frontend

    --> w-tainted: write(shipment_queue.ShipmentMessage) {1}       --> w-tainted: write(shipment_queue.ShipmentMessage) {1} --> r-tainted: read(shipment_queue.ShipmentMessage) {1}
[0] (StructObject UserType) shipmentMessage shopping_app.ShipmentMessage struct{OrderID string, UserID string}
     --> w-tainted: write(shipment_queue.ShipmentMessage.OrderID) {1}         --> w-tainted: write(shipment_queue.ShipmentMessage.OrderID) {1} --> r-tainted: read(shipment_queue.ShipmentMessage.OrderID) {1}
[_1] (FieldObject FieldType) OrderID string
      --> w-tainted: write(order_db.Order.OrderID, shipment_queue.ShipmentMessage.OrderID) {2}           --> w-tainted: write(order_db.Order.OrderID, shipment_queue.ShipmentMessage.OrderID) {2} --> r-tainted: read(shipment_queue.ShipmentMessage.OrderID) {1}
[__2] (BasicObject BasicType) cartID string
       --> w-tainted: write(order_db.Order.OrderID, shipment_queue.ShipmentMessage.OrderID) {2}             --> w-tainted: write(order_db.Order.OrderID, shipment_queue.ShipmentMessage.OrderID) {2} --> r-tainted: read(shipment_queue.ShipmentMessage.OrderID) {1}
[___3] (Reference BasicType) ref <cartID string> @ Frontend
     --> w-tainted: write(shipment_queue.ShipmentMessage.UserID) {1}         --> w-tainted: write(shipment_queue.ShipmentMessage.UserID) {1} --> r-tainted: read(shipment_queue.ShipmentMessage.UserID) {1}
[_1] (FieldObject FieldType) UserID string
      --> w-tainted: write(order_db.Order.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID, billing_db.Bill.UserID) {4}           --> w-tainted: write(order_db.Order.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID, billing_db.Bill.UserID) {4} --> r-tainted: read(analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {2}
[__2] (BasicObject BasicType) userID string
       --> w-tainted: write(order_db.Order.UserID, billing_db.Bill.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {4}             --> w-tainted: write(order_db.Order.UserID, billing_db.Bill.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {4} --> r-tainted: read(analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {2}
[___3] (Reference BasicType) ref <userID string> @ Frontend

