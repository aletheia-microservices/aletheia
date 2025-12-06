[0] (PointerObject PointerType) c (*config.ConfigServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) config.ConfigServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (StructObject UserType) conf config.Config struct{Name string, Value string, Description string}
[_1] (FieldObject FieldType) Name string
[__2] (BasicObject BasicType) Name string

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = config, collection = config}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "name" string,  string}

