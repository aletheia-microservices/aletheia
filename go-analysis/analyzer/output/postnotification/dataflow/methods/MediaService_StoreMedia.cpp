[0] (PointerObject PointerType) s (*postnotification.MediaServiceImpl struct{mediaDb NoSQLDatabase})
[_1] (StructObject UserType) postnotification.MediaServiceImpl struct{mediaDb NoSQLDatabase}
[__2] (FieldObject FieldType) mediaDb NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) mediaDb NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ UploadService

    --> w-tainted: write(media_db.Media) {1}
[0] (StructObject UserType) media postnotification.Media struct{MediaID int64, Content string}
     --> w-tainted: write(media_db.Media) {1}
[_1] (Reference UserType) ref <media postnotification.Media struct{MediaID int64, Content "HELLO WORLD!" untyped string}> @ UploadService
      --> w-tainted: write(media_db.Media.Content) {1}
[__2] (FieldObject FieldType) Content "HELLO WORLD!" untyped string
       --> w-tainted: write(media_db.Media.Content) {1}
[___3] (BasicObject BasicType) HELLO_WORLD_CONST "HELLO WORLD!" untyped string
     --> w-tainted: write(media_db.Media.MediaID) {1}
[_1] (FieldObject FieldType) MediaID int64
      --> w-tainted: write(media_db.Media.MediaID) {1}
[__2] (BasicObject BasicType) mediaID int64

    --> w-tainted: write(media_db.Media.MediaID) {1}
[0] (BasicObject BasicType) mediaID int64

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = media, collection = media}

[0] (InterfaceObject UserType) err .error

