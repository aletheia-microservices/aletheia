[0] (PointerObject PointerType) a (*postnotification.AnalyticsServiceImpl struct{analyticsQueue Queue, analyticsDb NoSQLDatabase, numWorkers int})
[_1] (StructObject UserType) postnotification.AnalyticsServiceImpl struct{analyticsQueue Queue, analyticsDb NoSQLDatabase, numWorkers 4 int}
[__2] (FieldObject FieldType) analyticsDb NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) analyticsDb NoSQLDatabase
[__2] (FieldObject FieldType) analyticsQueue Queue
[___3] (BlueprintBackendObject BlueprintBackendType) analyticsQueue Queue
[__2] (FieldObject FieldType) numWorkers 4 int
[___3] (BasicObject BasicType) 4 int

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ StorageService
[__2] (Reference UserType) ref <ctx context.Context> @ NotifyService
[___3] (Reference UserType) ref <ctx context.Context> @ NotifyService

    --> r-tainted: read(analytics_db.Analytics.PostID, posts_db.Post.PostID) {2}
[0] (BasicObject BasicType) postID int64
     --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[_1] (Reference BasicType) ref <postID_STORAGE_SVC_READ int64> @ StorageService
      --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[__2] (Reference BasicType) ref <postID_NOTIFY_SVC int64> @ NotifyService
       --> w-tainted: write(notifications_queue.Message.PostID) {1}             --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[___3] (BasicObject BasicType) PostID string
        --> w-tainted: write(notifications_queue.Message.PostID) {1}               --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[____4] (Reference FieldType) ref <PostID string> @ NotifyService
         --> w-tainted: write(notifications_queue.Message.PostID) {1}                 --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_____5] (BasicObject BasicType) string
          --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                   --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[______6] (BasicObject BasicType) postID_UPLOAD_SVC int64
           --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                     --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_______7] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
          --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                   --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[______6] (BasicObject BasicType) postID_UPLOAD_SVC int64
           --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                     --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_______7] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
       --> w-tainted: write(notifications_queue.Message.PostID) {1}             --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[___3] (BasicObject BasicType) PostID string
        --> w-tainted: write(notifications_queue.Message.PostID) {1}               --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[____4] (Reference FieldType) ref <PostID string> @ NotifyService
         --> w-tainted: write(notifications_queue.Message.PostID) {1}                 --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_____5] (BasicObject BasicType) string
          --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                   --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[______6] (BasicObject BasicType) postID_UPLOAD_SVC int64
           --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                     --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_______7] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
          --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                   --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[______6] (BasicObject BasicType) postID_UPLOAD_SVC int64
           --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                     --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_______7] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService

    --> r-tainted: read(analytics_db.Analytics) {1}
[0] (StructObject UserType) analytics postnotification.Analytics struct{PostID int64}

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = analyticsDb, collection = analytics_collection}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) analyticsQuery primitive.D
[_1] (StructObject StructType) struct{Key "postid" string, Key "postid" string, Value int64, Value int64}
[__2] (FieldObject FieldType) Key "postid" string
[___3] (BasicObject BasicType) "postid" string
[__2] (FieldObject FieldType) Value int64
       --> r-tainted: read(analytics_db.Analytics.PostID, posts_db.Post.PostID) {2}
[___3] (BasicObject BasicType) postID int64
        --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[____4] (Reference BasicType) ref <postID_STORAGE_SVC_READ int64> @ StorageService
         --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[_____5] (Reference BasicType) ref <postID_NOTIFY_SVC int64> @ NotifyService
          --> w-tainted: write(notifications_queue.Message.PostID) {1}                   --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[______6] (BasicObject BasicType) PostID string
           --> w-tainted: write(notifications_queue.Message.PostID) {1}                     --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_______7] (Reference FieldType) ref <PostID string> @ NotifyService
            --> w-tainted: write(notifications_queue.Message.PostID) {1}                       --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[________8] (BasicObject BasicType) string
             --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_________9] (BasicObject BasicType) postID_UPLOAD_SVC int64
              --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                           --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[__________10] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
             --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_________9] (BasicObject BasicType) postID_UPLOAD_SVC int64
              --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                           --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[__________10] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
          --> w-tainted: write(notifications_queue.Message.PostID) {1}                   --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[______6] (BasicObject BasicType) PostID string
           --> w-tainted: write(notifications_queue.Message.PostID) {1}                     --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_______7] (Reference FieldType) ref <PostID string> @ NotifyService
            --> w-tainted: write(notifications_queue.Message.PostID) {1}                       --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[________8] (BasicObject BasicType) string
             --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_________9] (BasicObject BasicType) postID_UPLOAD_SVC int64
              --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                           --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[__________10] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
             --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_________9] (BasicObject BasicType) postID_UPLOAD_SVC int64
              --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                           --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[__________10] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService

    --> r-tainted: read(analytics_db.Analytics) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = analyticsDb, collection = analytics_collection}
     --> r-tainted: read(analytics_db.Analytics) {1}
[_1] (StructObject UserType) analytics postnotification.Analytics struct{PostID int64}

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) res bool

[0] (InterfaceObject UserType) err .error

