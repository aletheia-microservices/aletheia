[0] (PointerObject PointerType) u (*user.UserServiceImpl struct{userDB NoSQLDatabase})
[_1] (StructObject UserType) user.UserServiceImpl struct{userDB NoSQLDatabase}
[__2] (FieldObject FieldType) userDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (StructObject UserType) user user.User struct{UserID string, Username string, Password string, Gender int64, DocumentType int64, DocumentNum string, Email string}
[_1] (FieldObject FieldType) UserID string
[__2] (BasicObject BasicType) UserID string

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = user, collection = user}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "userid" string,  string}

