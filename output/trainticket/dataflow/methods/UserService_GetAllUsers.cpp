[0] (PointerObject PointerType) u (*user.UserServiceImpl struct{userDB NoSQLDatabase})
[_1] (StructObject UserType) user.UserServiceImpl struct{userDB NoSQLDatabase}
[__2] (FieldObject FieldType) userDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = user, collection = user}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(user_db.User.UserID) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = user, collection = user}
     --> r-tainted: read(user_db.[]user.User struct{UserID string, Username string, Password string, Gender int64, DocumentType int64, DocumentNum string, Email string}) {1}
[_1] (ArrayObject ArrayType) users []user.User struct{UserID string, Username string, Password string, Gender int64, DocumentType int64, DocumentNum string, Email string}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(user_db.[]user.User struct{UserID string, Username string, Password string, Gender int64, DocumentType int64, DocumentNum string, Email string}) {1}
[0] (ArrayObject ArrayType) users []user.User struct{UserID string, Username string, Password string, Gender int64, DocumentType int64, DocumentNum string, Email string}

[0] (InterfaceObject UserType) err .error

