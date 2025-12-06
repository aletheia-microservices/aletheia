[0] (PointerObject PointerType) s (*catalogue.catalogueImpl struct{catalogue_db NoSQLDatabase})
[_1] (StructObject UserType) catalogue.catalogueImpl struct{catalogue_db NoSQLDatabase}
[__2] (FieldObject FieldType) catalogue_db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) catalogue_db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ FrontendService

[0] (ArrayObject ArrayType) tags []string
[_1] (Reference ArrayType) ref <tags []string> @ FrontendService

[0] (BasicObject BasicType) order string
[_1] (Reference BasicType) ref <order string> @ FrontendService

[0] (BasicObject BasicType) pageNum int
[_1] (Reference BasicType) ref <pageNum int> @ FrontendService

[0] (BasicObject BasicType) pageSize int
[_1] (Reference BasicType) ref <pageSize int> @ FrontendService

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = catalogue, collection = catalogue}

[0] (InterfaceObject UserType) _ .error

    --> r-tainted: read(catalogue_db.[]catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}) {1}
[0] (ArrayObject ArrayType) allSocks []catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}

[0] (ArrayObject ArrayType) socks []catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}
[_1] (StructObject UserType) sock catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}
      --> r-tainted: read(catalogue_db.[]catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}) {1}
[__2] (ArrayObject ArrayType) allSocks []catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}
[__2] (FieldObject FieldType) Tags []string
[___3] (SliceObject SliceType) Tags []string

[0] (SliceObject UserType) filter primitive.D

    --> r-tainted: read(catalogue_db.Sock.ID) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = catalogue, collection = catalogue}
     --> r-tainted: read(catalogue_db.[]catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}) {1}
[_1] (ArrayObject ArrayType) allSocks []catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}

[0] (InterfaceObject UserType) _ .error

[0] (StructObject UserType) sock catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}
     --> r-tainted: read(catalogue_db.[]catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}) {1}
[_1] (ArrayObject ArrayType) allSocks []catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}
[_1] (FieldObject FieldType) Tags []string
[__2] (SliceObject SliceType) Tags []string

[0] (BasicObject BasicType) tag string
[_1] (ArrayObject ArrayType) tags []string
[__2] (Reference ArrayType) ref <tags []string> @ FrontendService

