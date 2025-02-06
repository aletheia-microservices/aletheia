[0] (PointerObject PointerType) service (*shipping.shippingImpl struct{q Queue, db NoSQLDatabase})
[_1] (StructObject UserType) shipping.shippingImpl struct{q Queue, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase
[__2] (FieldObject FieldType) q Queue
[___3] (BlueprintBackendObject BlueprintBackendType) queue Queue

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ OrderService
[__2] (Reference UserType) ref <ctx context.Context> @ FrontendService

    --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3}       --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment) {1}
[0] (StructObject UserType) shipment shipping.Shipment struct{ID string, Name string, Status string}
     --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3}         --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment) {1}
[_1] (Reference UserType) ref <shipment shipping.Shipment struct{ID string, Name string, Status "awaiting shipment" string}> @ OrderService
      --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.Shipment.ID) {3}           --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.Shipment.ID) {3} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[__2] (FieldObject FieldType) ID string
       --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.Shipment.ID) {3}             --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.Shipment.ID) {3} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[___3] (BasicObject BasicType) string
      --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3}           --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Name) {1}
[__2] (FieldObject FieldType) Name string
       --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4}             --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4} --> r-tainted: read(shipqueue.Shipment.Name, user_db._.userid) {2}
[___3] (BasicObject BasicType) customerID string
        --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3}               --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3} --> r-tainted: read(user_db._.userid, shipqueue.Shipment.Name) {2}
[____4] (Reference BasicType) ref <userID string> @ FrontendService
      --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3}           --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Status) {1}
[__2] (FieldObject FieldType) Status "awaiting shipment" string
       --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3}             --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Status) {1}
[___3] (BasicObject BasicType) "awaiting shipment" string

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = shipping_service, collection = shipments}

[0] (InterfaceObject UserType) _ .error

