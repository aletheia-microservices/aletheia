[0] (PointerObject PointerType) c (*config.ConfigServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) config.ConfigServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

    --> r-tainted: read(config_db.Config.Name) {1}
[0] (BasicObject BasicType) name string

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = config, collection = config}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "name" string,  string}

    --> r-tainted: read(config_db.Config) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = config, collection = config}
     --> r-tainted: read(config_db.Config) {1}
[_1] (StructObject UserType) conf config.Config struct{Name string, Value string, Description string}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(config_db.Config) {1}
[0] (StructObject UserType) conf config.Config struct{Name string, Value string, Description string}

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

