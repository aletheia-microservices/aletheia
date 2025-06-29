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
[_1] (Reference UserType) ref <ctx context.Context> @ ShipmentService
[__2] (Reference UserType) ref <ctx context.Context> @ ShipmentService
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> w-tainted: write(shipment_db.Shipment.OrderID) {1}       --> w-tainted: write(shipment_db.Shipment.OrderID) {1} --> r-tainted: read(order_db._.orderID, order_db.Order.OrderID) {2}
[0] (BasicObject BasicType) orderID string
     --> w-tainted: write(shipment_db.Shipment.OrderID) {1}         --> w-tainted: write(shipment_db.Shipment.OrderID) {1} --> r-tainted: read(order_db._.orderID, order_db.Order.OrderID) {2}
[_1] (Reference BasicType) ref <OrderID string> @ ShipmentService
      --> w-tainted: write(shipment_db.Shipment.OrderID) {1}           --> w-tainted: write(shipment_db.Shipment.OrderID) {1} --> r-tainted: read(order_db._.orderID, order_db.Order.OrderID) {2}
[__2] (Reference FieldType) ref <OrderID string> @ ShipmentService
       --> w-tainted: write(shipment_db.Shipment.OrderID) {1}             --> w-tainted: write(shipment_db.Shipment.OrderID) {1} --> r-tainted: read(order_db._.orderID, order_db.Order.OrderID) {2}
[___3] (BasicObject BasicType) * string
        --> w-tainted: write(shipment_db.Shipment.OrderID) {1}               --> w-tainted: write(shipment_db.Shipment.OrderID) {1} --> r-tainted: read(order_db._.orderID, order_db.Order.OrderID) {2}
[____4] (InterfaceObject InterfaceType) * interface{}
     --> r-tainted: read(order_db._.orderID, order_db.Order.OrderID) {2}
[_1] (Reference BasicType) ref <orderID string> @ Frontend

    --> r-tainted: read(order_db.Order) {1}
[0] (StructObject UserType) order shopping_app.Order struct{OrderID string, UserID string, ProductID string, Quantity int, Timestamp int64}

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = order_database, collection = order_collection}

[0] (InterfaceObject UserType) _ .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "orderID" string, Key "orderID" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "orderID" string
[___3] (BasicObject BasicType) "orderID" string
[__2] (FieldObject FieldType) Value string
       --> w-tainted: write(shipment_db.Shipment.OrderID) {1}             --> w-tainted: write(shipment_db.Shipment.OrderID) {1} --> r-tainted: read(order_db._.orderID, order_db.Order.OrderID) {2}
[___3] (BasicObject BasicType) orderID string
        --> w-tainted: write(shipment_db.Shipment.OrderID) {1}               --> w-tainted: write(shipment_db.Shipment.OrderID) {1} --> r-tainted: read(order_db._.orderID, order_db.Order.OrderID) {2}
[____4] (Reference BasicType) ref <OrderID string> @ ShipmentService
         --> w-tainted: write(shipment_db.Shipment.OrderID) {1}                 --> w-tainted: write(shipment_db.Shipment.OrderID) {1} --> r-tainted: read(order_db._.orderID, order_db.Order.OrderID) {2}
[_____5] (Reference FieldType) ref <OrderID string> @ ShipmentService
          --> w-tainted: write(shipment_db.Shipment.OrderID) {1}                   --> w-tainted: write(shipment_db.Shipment.OrderID) {1} --> r-tainted: read(order_db._.orderID, order_db.Order.OrderID) {2}
[______6] (BasicObject BasicType) * string
           --> w-tainted: write(shipment_db.Shipment.OrderID) {1}                     --> w-tainted: write(shipment_db.Shipment.OrderID) {1} --> r-tainted: read(order_db._.orderID, order_db.Order.OrderID) {2}
[_______7] (InterfaceObject InterfaceType) * interface{}
        --> r-tainted: read(order_db._.orderID, order_db.Order.OrderID) {2}
[____4] (Reference BasicType) ref <orderID string> @ Frontend

    --> r-tainted: read(order_db._, order_db.Order) {2}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = order_database, collection = order_collection}
     --> r-tainted: read(order_db.Order) {1}
[_1] (StructObject UserType) order shopping_app.Order struct{OrderID string, UserID string, ProductID string, Quantity int, Timestamp int64}

[0] (InterfaceObject UserType) _ .error

