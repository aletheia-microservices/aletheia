[0] (PointerObject PointerType) s (*shipping.shippingImpl struct{q Queue, db NoSQLDatabase})
[_1] (StructObject UserType) shipping.shippingImpl struct{q Queue, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase
[__2] (FieldObject FieldType) q Queue
[___3] (BlueprintBackendObject BlueprintBackendType) queue Queue

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) id string

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = shipping_service, collection = shipments}

[0] (InterfaceObject UserType) _ .error

[0] (BlueprintBackendObject BlueprintBackendType) cursor NoSQLCursor {database = shipping_service, collection = shipments}
[_1] (StructObject UserType) shipment shipping.Shipment struct{ID string, Name string, Status string}

[0] (InterfaceObject UserType) err .error

[0] (StructObject UserType) shipment shipping.Shipment struct{ID string, Name string, Status string}

