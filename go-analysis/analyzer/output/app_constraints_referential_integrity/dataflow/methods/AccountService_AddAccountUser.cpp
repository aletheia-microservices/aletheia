[0] (PointerObject PointerType) s (*app_constraints_referential_integrity.AccountServiceImpl struct{userService app_constraints_referential_integrity.UserService, accountsDB NoSQLDatabase})
[_1] (StructObject UserType) app_constraints_referential_integrity.AccountServiceImpl struct{userService app_constraints_referential_integrity.UserService, accountsDB NoSQLDatabase}
[__2] (FieldObject FieldType) accountsDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) accountsDB NoSQLDatabase
[__2] (FieldObject FieldType) userService app_constraints_referential_integrity.UserService
[___3] (ServiceObject ServiceType) userService app_constraints_referential_integrity.UserService

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) accountID string

[0] (BasicObject BasicType) username string

[0] (StructObject UserType) accountUsers app_constraints_referential_integrity.AccountUsers struct{AccountID string, Usernames []string}
[_1] (FieldObject FieldType) Usernames []string
[__2] (SliceObject SliceType) Usernames []string
[___3] (BasicObject BasicType) username string

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = accounts, collection = accounts_users}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "accountID" string, Key "accountID" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "accountID" string
[___3] (BasicObject BasicType) "accountID" string
[__2] (FieldObject FieldType) Value string
[___3] (BasicObject BasicType) accountID string

[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = accounts, collection = accounts_users}
[_1] (StructObject UserType) accountUsers app_constraints_referential_integrity.AccountUsers struct{AccountID string, Usernames []string}
[__2] (FieldObject FieldType) Usernames []string
[___3] (SliceObject SliceType) Usernames []string
[____4] (BasicObject BasicType) username string

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) res bool

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

