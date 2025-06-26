[0] (PointerObject PointerType) a (*postnotification_simple.AnalyticsServiceImpl struct{analyticsQueue Queue, analyticsDb NoSQLDatabase, numWorkers int})
[_1] (StructObject UserType) postnotification_simple.AnalyticsServiceImpl struct{analyticsQueue Queue, analyticsDb NoSQLDatabase, numWorkers 4 int}
[__2] (FieldObject FieldType) analyticsDb NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) analyticsDb NoSQLDatabase
[__2] (FieldObject FieldType) analyticsQueue Queue
[___3] (BlueprintBackendObject BlueprintBackendType) analyticsQueue Queue
[__2] (FieldObject FieldType) numWorkers 4 int
[___3] (BasicObject BasicType) 4 int

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ UploadService

    --> r-tainted: read(analytics_db._.postid) {1}
[0] (BasicObject BasicType) postID int64
     --> r-tainted: read(posts_db._.postid, analytics_db._.postid, posts_db.Post.PostID) {3}
[_1] (Reference BasicType) ref <postID int64> @ UploadService

    --> r-tainted: read(analytics_db.Analytics) {1}
[0] (StructObject UserType) analytics postnotification_simple.Analytics struct{PostID int64}

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = analyticsDb, collection = analytics_collection}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) analyticsQuery primitive.D
[_1] (StructObject StructType) struct{Key "postid" string, Key "postid" string, Value int64, Value int64}
[__2] (FieldObject FieldType) Key "postid" string
[___3] (BasicObject BasicType) "postid" string
[__2] (FieldObject FieldType) Value int64
       --> r-tainted: read(analytics_db._.postid) {1}
[___3] (BasicObject BasicType) postID int64
        --> r-tainted: read(posts_db._.postid, analytics_db._.postid, posts_db.Post.PostID) {3}
[____4] (Reference BasicType) ref <postID int64> @ UploadService

    --> r-tainted: read(analytics_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = analyticsDb, collection = analytics_collection}
     --> r-tainted: read(analytics_db.Analytics) {1}
[_1] (StructObject UserType) analytics postnotification_simple.Analytics struct{PostID int64}

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) res bool

[0] (InterfaceObject UserType) err .error

