[0] (PointerObject PointerType) q (*queuemaster.queueMasterImpl struct{q Queue, shipping shipping.ShippingService, exitOnError bool, processed int32})
[_1] (StructObject UserType) queuemaster.queueMasterImpl struct{q Queue, shipping shipping.ShippingService, exitOnError false bool, processed 0 int}
[__2] (FieldObject FieldType) exitOnError false bool
[___3] (BasicObject BasicType) false bool
[__2] (FieldObject FieldType) processed 0 int
[___3] (BasicObject BasicType) 0 int
[__2] (FieldObject FieldType) q Queue
[___3] (BlueprintBackendObject BlueprintBackendType) queue Queue
[__2] (FieldObject FieldType) shipping shipping.ShippingService
[___3] (ServiceObject ServiceType) shipping shipping.ShippingService

[0] (InterfaceObject UserType) ctx context.Context

    --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3}       --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment) {1}
[0] (StructObject UserType) shipment shipping.Shipment struct{ID string, Name string, Status string}
     --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3}         --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment) {1}
[_1] (Reference UserType) ref <shipment shipping.Shipment struct{ID string, Name string, Status string}> @ QueueMaster
      --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3}           --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment) {1}
[__2] (Reference UserType) ref <shipment shipping.Shipment struct{ID string, Name string, Status "awaiting shipment" string}> @ OrderService
       --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.Shipment.ID) {3}             --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.Shipment.ID) {3} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[___3] (FieldObject FieldType) ID string
        --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.Shipment.ID) {3}               --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.Shipment.ID) {3} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[____4] (BasicObject BasicType) string
       --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3}             --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Name) {1}
[___3] (FieldObject FieldType) Name string
        --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4}               --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4} --> r-tainted: read(shipqueue.Shipment.Name, user_db._.userid) {2}
[____4] (BasicObject BasicType) customerID string
         --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3}                 --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3} --> r-tainted: read(user_db._.userid, shipqueue.Shipment.Name) {2}
[_____5] (Reference BasicType) ref <userID string> @ FrontendService
       --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3}             --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Status) {1}
[___3] (FieldObject FieldType) Status "awaiting shipment" string
        --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3}               --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Status) {1}
[____4] (BasicObject BasicType) "awaiting shipment" string
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
     --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID) {2}         --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID) {2} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[_1] (FieldObject FieldType) ID string
      --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.Shipment.ID) {3}           --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.Shipment.ID) {3} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[__2] (Reference FieldType) ref <ID string> @ OrderService
       --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.Shipment.ID) {3}             --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.Shipment.ID) {3} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[___3] (BasicObject BasicType) string
      --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID) {2}           --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID) {2} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[__2] (BasicObject BasicType) ID string
       --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.Shipment.ID) {3}             --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.Shipment.ID) {3} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[___3] (Reference FieldType) ref <ID string> @ OrderService
        --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.Shipment.ID) {3}               --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.Shipment.ID) {3} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[____4] (BasicObject BasicType) string

