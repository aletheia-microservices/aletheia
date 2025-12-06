[0] (PointerObject PointerType) s (*app_constraints_referential_integrity.AccountServiceImpl struct{userService app_constraints_referential_integrity.UserService, accountsDB NoSQLDatabase})
[_1] (StructObject UserType) app_constraints_referential_integrity.AccountServiceImpl struct{userService app_constraints_referential_integrity.UserService, accountsDB NoSQLDatabase}
[__2] (FieldObject FieldType) accountsDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) accountsDB NoSQLDatabase
[__2] (FieldObject FieldType) userService app_constraints_referential_integrity.UserService
[___3] (ServiceObject ServiceType) userService app_constraints_referential_integrity.UserService

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

[0] (BasicObject BasicType) accountID string
[_1] (Reference BasicType) ref <accountID string> @ Frontend

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = accounts, collection = accounts}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) filter primitive.D
[_1] (StructObject StructType) struct{Key "accountID" string, Key "accountID" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "accountID" string
[___3] (BasicObject BasicType) "accountID" string
[__2] (FieldObject FieldType) Value string
[___3] (BasicObject BasicType) accountID string
[____4] (Reference BasicType) ref <accountID string> @ Frontend

[0] (InterfaceObject UserType) err .error

