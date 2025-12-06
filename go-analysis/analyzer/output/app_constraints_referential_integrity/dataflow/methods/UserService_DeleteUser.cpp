[0] (PointerObject PointerType) s (*app_constraints_referential_integrity.UserServiceImpl struct{usersDB NoSQLDatabase})
[_1] (StructObject UserType) app_constraints_referential_integrity.UserServiceImpl struct{usersDB NoSQLDatabase}
[__2] (FieldObject FieldType) usersDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) usersDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

[0] (BasicObject BasicType) username string
[_1] (Reference BasicType) ref <username string> @ Frontend

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = users, collection = users}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) filter primitive.D
[_1] (StructObject StructType) struct{Key "username" string, Key "username" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "username" string
[___3] (BasicObject BasicType) "username" string
[__2] (FieldObject FieldType) Value string
[___3] (BasicObject BasicType) username string
[____4] (Reference BasicType) ref <username string> @ Frontend

[0] (InterfaceObject UserType) err .error

