[0] (PointerObject PointerType) u (*hotelreservation.UserServiceImpl struct{users map[string]string, userDB NoSQLDatabase})
[_1] (StructObject UserType) hotelreservation.UserServiceImpl struct{users map[string]string, userDB NoSQLDatabase}
[__2] (FieldObject FieldType) userDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) userDB NoSQLDatabase
[__2] (FieldObject FieldType) users map[string]string
[___3] (MapObject MapType) map[string]string

[0] (InterfaceObject UserType) ctx context.Context

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = user-db, collection = user}

[0] (InterfaceObject UserType) err .error

[0] (ArrayObject ArrayType) users []hotelreservation.User struct{Username string, Password string}

[0] (SliceObject UserType) filter primitive.D

[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = user-db, collection = user}
[_1] (ArrayObject ArrayType) users []hotelreservation.User struct{Username string, Password string}

[0] (InterfaceObject UserType) err .error

