[0] (PointerObject PointerType) s (*cart.cartImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) cart.cartImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ OrderService
[__2] (Reference UserType) ref <ctx context.Context> @ FrontendService
[_1] (Reference UserType) ref <ctx context.Context> @ FrontendService

    --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3}       --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3} --> r-tainted: read(user_db._.userid, shipqueue.Shipment.Name) {2}
[0] (BasicObject BasicType) customerID string
     --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4}         --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4} --> r-tainted: read(shipqueue.Shipment.Name, user_db._.userid) {2}
[_1] (Reference BasicType) ref <customerID string> @ OrderService
      --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3}           --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3} --> r-tainted: read(user_db._.userid, shipqueue.Shipment.Name) {2}
[__2] (Reference BasicType) ref <userID string> @ FrontendService
[_1] (Reference BasicType) ref <sessionID string> @ FrontendService

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = cart, collection = carts}

[0] (InterfaceObject UserType) _ .error

