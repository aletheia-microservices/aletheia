[0] (PointerObject PointerType) s (*catalogue.catalogueImpl struct{catalogue_db NoSQLDatabase})
[_1] (StructObject UserType) catalogue.catalogueImpl struct{catalogue_db NoSQLDatabase}
[__2] (FieldObject FieldType) catalogue_db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) catalogue_db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (ArrayObject ArrayType) tags []string

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = catalogue, collection = catalogue}

[0] (InterfaceObject UserType) _ .error

[0] (ArrayObject ArrayType) socks []catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}

[0] (SliceObject UserType) filter primitive.D

[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = catalogue, collection = catalogue}
[_1] (ArrayObject ArrayType) socks []catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}

[0] (InterfaceObject UserType) _ .error

