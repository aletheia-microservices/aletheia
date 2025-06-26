[0] (PointerObject PointerType) s (*postnotification_simple.StorageServiceImpl struct{postsDb NoSQLDatabase, analyticsQueue Queue})
[_1] (StructObject UserType) postnotification_simple.StorageServiceImpl struct{postsDb NoSQLDatabase, analyticsQueue Queue}
[__2] (FieldObject FieldType) analyticsQueue Queue
[___3] (BlueprintBackendObject BlueprintBackendType) analyticsQueue Queue
[__2] (FieldObject FieldType) postsDb NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) postsDb NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ UploadService

[0] (BasicObject BasicType) postID int64
[_1] (Reference BasicType) ref <postID int64> @ UploadService

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = post, collection = post}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) filter primitive.D
[_1] (StructObject StructType) struct{Key "postid" string, Key "postid" string, Value int64, Value int64}
[__2] (FieldObject FieldType) Key "postid" string
[___3] (BasicObject BasicType) "postid" string
[__2] (FieldObject FieldType) Value int64
[___3] (BasicObject BasicType) postID int64
[____4] (Reference BasicType) ref <postID int64> @ UploadService

[0] (InterfaceObject UserType) err .error

