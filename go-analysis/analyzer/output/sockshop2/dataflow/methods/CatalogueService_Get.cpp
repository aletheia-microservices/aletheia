[0] (PointerObject PointerType) s (*catalogue.catalogueImpl struct{catalogue_db NoSQLDatabase})
[_1] (StructObject UserType) catalogue.catalogueImpl struct{catalogue_db NoSQLDatabase}
[__2] (FieldObject FieldType) catalogue_db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) catalogue_db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ FrontendService
[_1] (Reference UserType) ref <ctx context.Context> @ FrontendService
[_1] (Reference UserType) ref <ctx context.Context> @ FrontendService

    --> r-tainted: read(catalogue_db.Sock.ID.id) {1}
[0] (BasicObject BasicType) id string
     --> r-tainted: read(catalogue_db.Sock.ID.id) {1}
[_1] (Reference BasicType) ref <itemID string> @ FrontendService
     --> r-tainted: read(catalogue_db.Sock.ID.id) {1}
[_1] (Reference BasicType) ref <itemID string> @ FrontendService
     --> w-tainted: write(cart_db.cart.Items) {1}         --> w-tainted: write(cart_db.cart.Items) {1} --> r-tainted: read(catalogue_db.Sock.ID.id) {1}
[_1] (Reference BasicType) ref <itemID string> @ FrontendService

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = catalogue, collection = catalogue}

[0] (InterfaceObject UserType) _ .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "id" string, Key "id" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "id" string
[___3] (BasicObject BasicType) "id" string
[__2] (FieldObject FieldType) Value string
       --> r-tainted: read(catalogue_db.Sock.ID.id) {1}
[___3] (BasicObject BasicType) id string
        --> r-tainted: read(catalogue_db.Sock.ID.id) {1}
[____4] (Reference BasicType) ref <itemID string> @ FrontendService
        --> r-tainted: read(catalogue_db.Sock.ID.id) {1}
[____4] (Reference BasicType) ref <itemID string> @ FrontendService
        --> w-tainted: write(cart_db.cart.Items) {1}               --> w-tainted: write(cart_db.cart.Items) {1} --> r-tainted: read(catalogue_db.Sock.ID.id) {1}
[____4] (Reference BasicType) ref <itemID string> @ FrontendService

    --> r-tainted: read(catalogue_db.Sock, catalogue_db.Sock.ImageURL, catalogue_db.Sock.string, catalogue_db.Sock.ImageURL_1, catalogue_db.Sock.ImageURL_2, catalogue_db.Sock.TagString, catalogue_db.Sock.Tags) {7}
[0] (StructObject UserType) sock catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}
     --> r-tainted: read(catalogue_db.Sock.ImageURL) {1}
[_1] (FieldObject FieldType) ImageURL []string
      --> r-tainted: read(catalogue_db.Sock.ImageURL) {1}
[__2] (ArrayObject SliceType) []string
       --> r-tainted: read(catalogue_db.Sock.string, catalogue_db.Sock.ImageURL_1) {2}
[___3] (BasicObject BasicType) ImageURL_1 string
       --> r-tainted: read(catalogue_db.Sock.string, catalogue_db.Sock.ImageURL_2) {2}
[___3] (BasicObject BasicType) ImageURL_2 string
     --> r-tainted: read(catalogue_db.Sock.ImageURL_1) {1}
[_1] (FieldObject FieldType) ImageURL_1 string
      --> r-tainted: read(catalogue_db.Sock.string, catalogue_db.Sock.ImageURL_1) {2}
[__2] (BasicObject BasicType) ImageURL_1 string
     --> r-tainted: read(catalogue_db.Sock.ImageURL_2) {1}
[_1] (FieldObject FieldType) ImageURL_2 string
      --> r-tainted: read(catalogue_db.Sock.string, catalogue_db.Sock.ImageURL_2) {2}
[__2] (BasicObject BasicType) ImageURL_2 string
     --> r-tainted: read(catalogue_db.Sock.TagString) {1}
[_1] (FieldObject FieldType) TagString string
      --> r-tainted: read(catalogue_db.Sock.TagString) {1}
[__2] (BasicObject BasicType) TagString string
     --> r-tainted: read(catalogue_db.Sock.Tags) {1}
[_1] (FieldObject FieldType) Tags []string
      --> r-tainted: read(catalogue_db.Sock.Tags) {1}
[__2] (SliceObject SliceType) []string
       --> r-tainted: read(catalogue_db.Sock.TagString) {1}
[___3] (BasicObject BasicType) TagString string
       --> r-tainted: read(catalogue_db.Sock.string) {1}
[___3] (BasicObject BasicType) "," string
       --> r-tainted: read(catalogue_db.Sock.TagString) {1}
[___3] (BasicObject BasicType) TagString string
       --> r-tainted: read(catalogue_db.Sock.string) {1}
[___3] (BasicObject BasicType) "," string

    --> r-tainted: read(catalogue_db.Sock.ID) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = catalogue, collection = catalogue}
     --> r-tainted: read(catalogue_db.Sock, catalogue_db.Sock.ImageURL, catalogue_db.Sock.string, catalogue_db.Sock.ImageURL_1, catalogue_db.Sock.ImageURL_2, catalogue_db.Sock.TagString, catalogue_db.Sock.Tags) {7}
[_1] (StructObject UserType) sock catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}
      --> r-tainted: read(catalogue_db.Sock.ImageURL) {1}
[__2] (FieldObject FieldType) ImageURL []string
       --> r-tainted: read(catalogue_db.Sock.ImageURL) {1}
[___3] (ArrayObject SliceType) []string
        --> r-tainted: read(catalogue_db.Sock.string, catalogue_db.Sock.ImageURL_1) {2}
[____4] (BasicObject BasicType) ImageURL_1 string
        --> r-tainted: read(catalogue_db.Sock.string, catalogue_db.Sock.ImageURL_2) {2}
[____4] (BasicObject BasicType) ImageURL_2 string
      --> r-tainted: read(catalogue_db.Sock.ImageURL_1) {1}
[__2] (FieldObject FieldType) ImageURL_1 string
       --> r-tainted: read(catalogue_db.Sock.string, catalogue_db.Sock.ImageURL_1) {2}
[___3] (BasicObject BasicType) ImageURL_1 string
      --> r-tainted: read(catalogue_db.Sock.ImageURL_2) {1}
[__2] (FieldObject FieldType) ImageURL_2 string
       --> r-tainted: read(catalogue_db.Sock.string, catalogue_db.Sock.ImageURL_2) {2}
[___3] (BasicObject BasicType) ImageURL_2 string
      --> r-tainted: read(catalogue_db.Sock.TagString) {1}
[__2] (FieldObject FieldType) TagString string
       --> r-tainted: read(catalogue_db.Sock.TagString) {1}
[___3] (BasicObject BasicType) TagString string
      --> r-tainted: read(catalogue_db.Sock.Tags) {1}
[__2] (FieldObject FieldType) Tags []string
       --> r-tainted: read(catalogue_db.Sock.Tags) {1}
[___3] (SliceObject SliceType) []string
        --> r-tainted: read(catalogue_db.Sock.TagString) {1}
[____4] (BasicObject BasicType) TagString string
        --> r-tainted: read(catalogue_db.Sock.string) {1}
[____4] (BasicObject BasicType) "," string
        --> r-tainted: read(catalogue_db.Sock.TagString) {1}
[____4] (BasicObject BasicType) TagString string
        --> r-tainted: read(catalogue_db.Sock.string) {1}
[____4] (BasicObject BasicType) "," string

[0] (InterfaceObject UserType) _ .error

