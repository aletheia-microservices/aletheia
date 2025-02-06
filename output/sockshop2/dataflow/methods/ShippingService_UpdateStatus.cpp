[0] (PointerObject PointerType) s (*shipping.shippingImpl struct{q Queue, db NoSQLDatabase})
[_1] (StructObject UserType) shipping.shippingImpl struct{q Queue, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase
[__2] (FieldObject FieldType) q Queue
[___3] (BlueprintBackendObject BlueprintBackendType) queue Queue

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ QueueMaster

    --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID) {2}       --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID) {2} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[0] (BasicObject BasicType) id string
     --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID) {2}         --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID) {2} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[_1] (Reference BasicType) ref <ID string> @ QueueMaster
      --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.Shipment.ID) {3}           --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.Shipment.ID) {3} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[__2] (Reference FieldType) ref <ID string> @ OrderService
       --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.Shipment.ID) {3}             --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.Shipment.ID) {3} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[___3] (BasicObject BasicType) string

[0] (BasicObject BasicType) status string
[_1] (Reference BasicType) ref <"shipped" string> @ QueueMaster

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = shipping_service, collection = shipments}

[0] (InterfaceObject UserType) _ .error

