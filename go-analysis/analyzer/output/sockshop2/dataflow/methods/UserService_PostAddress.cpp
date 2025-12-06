[0] (PointerObject PointerType) s (*user.userServiceImpl struct{UserService user.UserService, db NoSQLDatabase})
[_1] (StructObject UserType) user.userServiceImpl struct{UserService user.UserService, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ FrontendService

[0] (BasicObject BasicType) userid string
[_1] (Reference BasicType) ref <userID string> @ FrontendService

[0] (StructObject UserType) address user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}
[_1] (Reference UserType) ref <address user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}> @ FrontendService

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = users, collection = users}

[0] (InterfaceObject UserType) _ .error

[0] (SliceObject UserType) filter primitive.D
[_1] (StructObject StructType) struct{Key "userid" string, Key "userid" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "userid" string
[___3] (BasicObject BasicType) "userid" string
[__2] (FieldObject FieldType) Value string
[___3] (BasicObject BasicType) userid string
[____4] (Reference BasicType) ref <userID string> @ FrontendService

[0] (SliceObject UserType) update primitive.D
[_1] (StructObject StructType) struct{Key "address" string, Key "address" string, Value user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}, Value user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}}
[__2] (FieldObject FieldType) Key "address" string
[___3] (BasicObject BasicType) "address" string
[__2] (FieldObject FieldType) Value user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}
[___3] (StructObject UserType) address user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}
[____4] (Reference UserType) ref <address user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}> @ FrontendService

