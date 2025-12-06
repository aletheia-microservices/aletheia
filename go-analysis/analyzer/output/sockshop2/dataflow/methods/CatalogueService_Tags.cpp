[0] (PointerObject PointerType) s (*catalogue.catalogueImpl struct{catalogue_db NoSQLDatabase})
[_1] (StructObject UserType) catalogue.catalogueImpl struct{catalogue_db NoSQLDatabase}
[__2] (FieldObject FieldType) catalogue_db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) catalogue_db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ FrontendService

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = catalogue, collection = catalogue}

[0] (InterfaceObject UserType) _ .error

    --> r-tainted: read(catalogue_db.[]catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}) {1}
[0] (ArrayObject ArrayType) socks []catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}

[0] (SliceObject UserType) projection primitive.D
[_1] (StructObject StructType) struct{Key "tags" string, Key "tags" string, Value true bool, Value true bool}
[__2] (FieldObject FieldType) Key "tags" string
[___3] (BasicObject BasicType) "tags" string
[__2] (FieldObject FieldType) Value true bool
[___3] (BasicObject BasicType) true bool

[0] (SliceObject UserType) filter primitive.D

    --> r-tainted: read(catalogue_db.Sock.ID) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = catalogue, collection = catalogue}
     --> r-tainted: read(catalogue_db.[]catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}) {1}
[_1] (ArrayObject ArrayType) socks []catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}

[0] (InterfaceObject UserType) _ .error

[0] (ArrayObject ArrayType) tags []string

