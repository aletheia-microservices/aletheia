[0] (PointerObject PointerType) p (*payment.PaymentServiceImpl struct{paymentDB NoSQLDatabase, moneyDB NoSQLDatabase})
[_1] (StructObject UserType) payment.PaymentServiceImpl struct{paymentDB NoSQLDatabase, moneyDB NoSQLDatabase}
[__2] (FieldObject FieldType) moneyDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) moneyDB NoSQLDatabase
[__2] (FieldObject FieldType) paymentDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) payDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (StructObject UserType) payment payment.Payment struct{ID string, OrderID string, UserID string, Price string}
[_1] (FieldObject FieldType) Price string
      --> w-tainted: write(money_db.Money.Price) {1}
[__2] (BasicObject BasicType) Price string
[_1] (FieldObject FieldType) UserID string
      --> w-tainted: write(money_db.Money.UserID) {1}
[__2] (BasicObject BasicType) UserID string

    --> w-tainted: write(money_db.Money) {1}
[0] (StructObject UserType) m payment.Money struct{ID string, UserID string, Price string}
     --> w-tainted: write(money_db.Money.ID) {1}
[_1] (FieldObject FieldType) ID string
      --> w-tainted: write(money_db.Money.ID) {1}
[__2] (BasicObject BasicType) string
     --> w-tainted: write(money_db.Money.Price) {1}
[_1] (FieldObject FieldType) Price string
      --> w-tainted: write(money_db.Money.Price) {1}
[__2] (BasicObject BasicType) Price string
     --> w-tainted: write(money_db.Money.UserID) {1}
[_1] (FieldObject FieldType) UserID string
      --> w-tainted: write(money_db.Money.UserID) {1}
[__2] (BasicObject BasicType) UserID string

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = payment, collection = money}

[0] (InterfaceObject UserType) err .error

