[0] (PointerObject PointerType) s (*user.userServiceImpl struct{UserService user.UserService, db NoSQLDatabase})
[_1] (StructObject UserType) user.userServiceImpl struct{UserService user.UserService, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ FrontendService

    --> r-tainted: read(user_db._.username) {1}
[0] (BasicObject BasicType) username string
     --> r-tainted: read(user_db._.username) {1}
[_1] (Reference BasicType) ref <username string> @ FrontendService

[0] (BasicObject BasicType) password string
[_1] (Reference BasicType) ref <password string> @ FrontendService

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = users, collection = users}

[0] (InterfaceObject UserType) _ .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "username" string, Key "username" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "username" string
[___3] (BasicObject BasicType) "username" string
[__2] (FieldObject FieldType) Value string
       --> r-tainted: read(user_db._.username) {1}
[___3] (BasicObject BasicType) username string
        --> r-tainted: read(user_db._.username) {1}
[____4] (Reference BasicType) ref <username string> @ FrontendService

[0] (SliceObject UserType) projection primitive.D
[_1] (StructObject StructType) struct{Key "password" string, Key "password" string, Value true bool, Value true bool}
[__2] (FieldObject FieldType) Key "password" string
[___3] (BasicObject BasicType) "password" string
[__2] (FieldObject FieldType) Value true bool
[___3] (BasicObject BasicType) true bool

    --> r-tainted: read(user_db.User) {1}
[0] (StructObject UserType) user user.User struct{FirstName string, LastName string, Email string, Username string, Password string, Addresses user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}, Cards user.Card struct{LongNum string, Expires string, CCV string, ID string}, UserID string, Salt string}

    --> r-tainted: read(user_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = users, collection = users}
     --> r-tainted: read(user_db.User) {1}
[_1] (StructObject UserType) user user.User struct{FirstName string, LastName string, Email string, Username string, Password string, Addresses user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}, Cards user.Card struct{LongNum string, Expires string, CCV string, ID string}, UserID string, Salt string}

[0] (InterfaceObject UserType) _ .error

