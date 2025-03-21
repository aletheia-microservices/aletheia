[0] (PointerObject PointerType) s (*postnotification.MediaServiceImpl struct{mediaDb NoSQLDatabase})
[_1] (StructObject UserType) postnotification.MediaServiceImpl struct{mediaDb NoSQLDatabase}
[__2] (FieldObject FieldType) mediaDb NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) mediaDb NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ StorageService
[__2] (Reference UserType) ref <ctx context.Context> @ UploadService

    --> r-tainted: read(media_db._.mediaid) {1}
[0] (BasicObject BasicType) mediaID int64
     --> r-tainted: read(media_db._.mediaid) {1}
[_1] (Reference FieldType) ref <mediaID int64> @ StorageService
      --> r-tainted: read(posts_cache.PostID, media_db._.mediaid) {2}
[__2] (FieldObject FieldType) PostID int64
       --> r-tainted: read(posts_cache.int64, media_db._.mediaid) {2}
[___3] (BasicObject BasicType) PostID int64
      --> r-tainted: read(media_db._.mediaid) {1}
[__2] (BasicObject BasicType) int64
       --> r-tainted: read(posts_cache.int64, media_db._.mediaid) {2}
[___3] (BasicObject BasicType) PostID int64

    --> r-tainted: read(media_db.Media) {1}
[0] (StructObject UserType) media postnotification.Media struct{MediaID int64, Content string}

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = media, collection = media}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) mediaQuery primitive.D
[_1] (StructObject StructType) struct{Key "mediaid" string, Key "mediaid" string, Value int64, Value int64}
[__2] (FieldObject FieldType) Key "mediaid" string
[___3] (BasicObject BasicType) "mediaid" string
[__2] (FieldObject FieldType) Value int64
       --> r-tainted: read(media_db._.mediaid) {1}
[___3] (BasicObject BasicType) mediaID int64
        --> r-tainted: read(media_db._.mediaid) {1}
[____4] (Reference FieldType) ref <mediaID int64> @ StorageService
         --> r-tainted: read(posts_cache.PostID, media_db._.mediaid) {2}
[_____5] (FieldObject FieldType) PostID int64
          --> r-tainted: read(posts_cache.int64, media_db._.mediaid) {2}
[______6] (BasicObject BasicType) PostID int64
         --> r-tainted: read(media_db._.mediaid) {1}
[_____5] (BasicObject BasicType) int64
          --> r-tainted: read(posts_cache.int64, media_db._.mediaid) {2}
[______6] (BasicObject BasicType) PostID int64

    --> r-tainted: read(media_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = media, collection = media}
     --> r-tainted: read(media_db.Media) {1}
[_1] (StructObject UserType) media postnotification.Media struct{MediaID int64, Content string}

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) _ bool

[0] (InterfaceObject UserType) err .error

