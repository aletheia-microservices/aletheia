[0] (PointerObject PointerType) p (*payment.PaymentServiceImpl struct{paymentDB NoSQLDatabase, moneyDB NoSQLDatabase})
[_1] (StructObject UserType) payment.PaymentServiceImpl struct{paymentDB NoSQLDatabase, moneyDB NoSQLDatabase}
[__2] (FieldObject FieldType) moneyDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) moneyDB NoSQLDatabase
[__2] (FieldObject FieldType) paymentDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) payDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

    --> r-tainted: read(payments_db.[]payment.Payment struct{ID string, OrderID string, UserID string, Price string}) {1}
[0] (ArrayObject ArrayType) payments []payment.Payment struct{ID string, OrderID string, UserID string, Price string}

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = payment, collection = payment}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(payments_db.Payment) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = payment, collection = payment}
     --> r-tainted: read(payments_db.[]payment.Payment struct{ID string, OrderID string, UserID string, Price string}) {1}
[_1] (ArrayObject ArrayType) payments []payment.Payment struct{ID string, OrderID string, UserID string, Price string}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

