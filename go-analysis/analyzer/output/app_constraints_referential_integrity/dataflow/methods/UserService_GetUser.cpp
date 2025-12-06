[0] (PointerObject PointerType) s (*app_constraints_referential_integrity.UserServiceImpl struct{usersDB NoSQLDatabase})
[_1] (StructObject UserType) app_constraints_referential_integrity.UserServiceImpl struct{usersDB NoSQLDatabase}
[__2] (FieldObject FieldType) usersDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) usersDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ AccountService
[__2] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> r-tainted: read(users_db.User.Username) {1}
[0] (BasicObject BasicType) username string
     --> r-tainted: read(accounts_db.Account.Username, users_db.User.Username) {2}
[_1] (Reference BasicType) ref <Username string> @ AccountService

    --> r-tainted: read(users_db.User) {1}
[0] (StructObject UserType) user app_constraints_referential_integrity.User struct{ReqID string, UserID string, Username string, Timestamp int64}

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = users, collection = users}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "username" string, Key "username" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "username" string
[___3] (BasicObject BasicType) "username" string
[__2] (FieldObject FieldType) Value string
       --> r-tainted: read(users_db.User.Username) {1}
[___3] (BasicObject BasicType) username string
        --> r-tainted: read(accounts_db.Account.Username, users_db.User.Username) {2}
[____4] (Reference BasicType) ref <Username string> @ AccountService

    --> r-tainted: read(users_db.User) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = users, collection = users}
     --> r-tainted: read(users_db.User) {1}
[_1] (StructObject UserType) user app_constraints_referential_integrity.User struct{ReqID string, UserID string, Username string, Timestamp int64}

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) res bool

[0] (InterfaceObject UserType) err .error

