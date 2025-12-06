[0] (PointerObject PointerType) s (*catalogue.catalogueImpl struct{catalogue_db NoSQLDatabase})
[_1] (StructObject UserType) catalogue.catalogueImpl struct{catalogue_db NoSQLDatabase}
[__2] (FieldObject FieldType) catalogue_db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) catalogue_db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ FrontendService

[0] (BasicObject BasicType) id string
[_1] (Reference BasicType) ref <id string> @ FrontendService

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = catalogue, collection = catalogue}

[0] (InterfaceObject UserType) _ .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "id" string, Key "id" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "id" string
[___3] (BasicObject BasicType) "id" string
[__2] (FieldObject FieldType) Value string
[___3] (BasicObject BasicType) id string
[____4] (Reference BasicType) ref <id string> @ FrontendService

