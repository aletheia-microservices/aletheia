[0] (PointerObject PointerType) p (*payment.PaymentServiceImpl struct{paymentDB NoSQLDatabase, moneyDB NoSQLDatabase})
[_1] (StructObject UserType) payment.PaymentServiceImpl struct{paymentDB NoSQLDatabase, moneyDB NoSQLDatabase}
[__2] (FieldObject FieldType) moneyDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) moneyDB NoSQLDatabase
[__2] (FieldObject FieldType) paymentDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) payDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

    --> w-tainted: write(payments_db.Payment) {1}
[0] (StructObject UserType) payment payment.Payment struct{ID string, OrderID string, UserID string, Price string}
     --> w-tainted: write(payments_db.Payment.ID) {1}
[_1] (FieldObject FieldType) ID string
      --> w-tainted: write(payments_db.Payment.ID) {1}           --> w-tainted: write(payments_db.Payment.ID) {1} --> r-tainted: read(payments_db._.id) {1}
[__2] (BasicObject BasicType) ID string

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = payment, collection = payment}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(payments_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = payment, collection = payment}
     --> r-tainted: read(payments_db.Payment) {1}
[_1] (StructObject UserType) stored payment.Payment struct{ID string, OrderID string, UserID string, Price string}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(payments_db.Payment) {1}
[0] (StructObject UserType) stored payment.Payment struct{ID string, OrderID string, UserID string, Price string}

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

