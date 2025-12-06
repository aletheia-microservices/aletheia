[0] (PointerObject PointerType) u (*user.UserServiceImpl struct{userDB NoSQLDatabase})
[_1] (StructObject UserType) user.UserServiceImpl struct{userDB NoSQLDatabase}
[__2] (FieldObject FieldType) userDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

    --> r-tainted: read(user_db.User.UserID.username) {1}
[0] (BasicObject BasicType) username string

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = user, collection = user}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "username" string,  string}

    --> r-tainted: read(user_db.User.UserID) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = user, collection = user}
     --> r-tainted: read(user_db.User) {1}
[_1] (StructObject UserType) user user.User struct{UserID string, Username string, Password string, Gender int64, DocumentType int64, DocumentNum string, Email string}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(user_db.User) {1}
[0] (StructObject UserType) user user.User struct{UserID string, Username string, Password string, Gender int64, DocumentType int64, DocumentNum string, Email string}

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

