[0] (PointerObject PointerType) p (*payment.PaymentServiceImpl struct{paymentDB NoSQLDatabase, moneyDB NoSQLDatabase})
[_1] (StructObject UserType) payment.PaymentServiceImpl struct{paymentDB NoSQLDatabase, moneyDB NoSQLDatabase}
[__2] (FieldObject FieldType) moneyDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) moneyDB NoSQLDatabase
[__2] (FieldObject FieldType) paymentDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) payDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (StructObject UserType) payment payment.Payment struct{ID string, OrderID string, UserID string, Price string}
[_1] (FieldObject FieldType) OrderID string
[__2] (BasicObject BasicType) OrderID string

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = payment, collection = payment}

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) ok bool

[0] (InterfaceObject UserType) err .error

