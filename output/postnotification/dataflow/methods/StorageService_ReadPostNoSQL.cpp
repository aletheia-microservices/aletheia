[0] (PointerObject PointerType) s (*postnotification.StorageServiceImpl struct{analyticsService postnotification.AnalyticsService, mediaService postnotification.MediaService, posts_cache Cache, postsDb NoSQLDatabase, analyticsQueue Queue})
[_1] (StructObject UserType) postnotification.StorageServiceImpl struct{analyticsService postnotification.AnalyticsService, mediaService postnotification.MediaService, posts_cache Cache, postsDb NoSQLDatabase, analyticsQueue Queue}
[__2] (FieldObject FieldType) analyticsQueue Queue
[___3] (BlueprintBackendObject BlueprintBackendType) analyticsQueue Queue
[__2] (FieldObject FieldType) analyticsService postnotification.AnalyticsService
[___3] (ServiceObject ServiceType) analyticsService postnotification.AnalyticsService
[__2] (FieldObject FieldType) mediaService postnotification.MediaService
[___3] (ServiceObject ServiceType) mediaService postnotification.MediaService
[__2] (FieldObject FieldType) postsDb NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) postsDb NoSQLDatabase
[__2] (FieldObject FieldType) posts_cache Cache
[___3] (BlueprintBackendObject BlueprintBackendType) posts_cache Cache

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ NotifyService
[__2] (Reference UserType) ref <ctx context.Context> @ NotifyService

[0] (BasicObject BasicType) reqID int64
[_1] (Reference BasicType) ref <reqID int64> @ NotifyService
      --> w-tainted: write(notifications_queue.Message.ReqID) {1}
[__2] (BasicObject BasicType) ReqID string
       --> w-tainted: write(notifications_queue.Message.ReqID) {1}             --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[___3] (Reference FieldType) ref <ReqID string> @ NotifyService
        --> w-tainted: write(notifications_queue.Message.ReqID) {1}               --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[____4] (BasicObject BasicType) string
         --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}                 --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[_____5] (BasicObject BasicType) reqID int64
         --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}                 --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[_____5] (BasicObject BasicType) reqID int64
      --> w-tainted: write(notifications_queue.Message.ReqID) {1}
[__2] (BasicObject BasicType) ReqID string
       --> w-tainted: write(notifications_queue.Message.ReqID) {1}             --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[___3] (Reference FieldType) ref <ReqID string> @ NotifyService
        --> w-tainted: write(notifications_queue.Message.ReqID) {1}               --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[____4] (BasicObject BasicType) string
         --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}                 --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[_____5] (BasicObject BasicType) reqID int64
         --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}                 --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[_____5] (BasicObject BasicType) reqID int64

    --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[0] (BasicObject BasicType) postID_STORAGE_SVC_READ int64
     --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[_1] (Reference BasicType) ref <postID_NOTIFY_SVC int64> @ NotifyService
      --> w-tainted: write(notifications_queue.Message.PostID) {1}           --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[__2] (BasicObject BasicType) PostID string
       --> w-tainted: write(notifications_queue.Message.PostID) {1}             --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[___3] (Reference FieldType) ref <PostID string> @ NotifyService
        --> w-tainted: write(notifications_queue.Message.PostID) {1}               --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[____4] (BasicObject BasicType) string
         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                 --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_____5] (BasicObject BasicType) postID_UPLOAD_SVC int64
          --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                   --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[______6] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                 --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_____5] (BasicObject BasicType) postID_UPLOAD_SVC int64
          --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                   --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[______6] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
      --> w-tainted: write(notifications_queue.Message.PostID) {1}           --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[__2] (BasicObject BasicType) PostID string
       --> w-tainted: write(notifications_queue.Message.PostID) {1}             --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[___3] (Reference FieldType) ref <PostID string> @ NotifyService
        --> w-tainted: write(notifications_queue.Message.PostID) {1}               --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[____4] (BasicObject BasicType) string
         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                 --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_____5] (BasicObject BasicType) postID_UPLOAD_SVC int64
          --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                   --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[______6] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                 --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_____5] (BasicObject BasicType) postID_UPLOAD_SVC int64
          --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                   --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[______6] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService

    --> r-tainted: read(posts_db.Post) {1}
[0] (StructObject UserType) post postnotification.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification.Creator struct{Username string}}

[0] (StructObject UserType) analytics postnotification.Analytics struct{PostID int64}

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = post, collection = post}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) postsQuery primitive.D
[_1] (StructObject StructType) struct{Key "postid" string, Key "postid" string, Value int64, Value int64}
[__2] (FieldObject FieldType) Key "postid" string
[___3] (BasicObject BasicType) "postid" string
[__2] (FieldObject FieldType) Value int64
       --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[___3] (BasicObject BasicType) postID_STORAGE_SVC_READ int64
        --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[____4] (Reference BasicType) ref <postID_NOTIFY_SVC int64> @ NotifyService
         --> w-tainted: write(notifications_queue.Message.PostID) {1}                 --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[_____5] (BasicObject BasicType) PostID string
          --> w-tainted: write(notifications_queue.Message.PostID) {1}                   --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[______6] (Reference FieldType) ref <PostID string> @ NotifyService
           --> w-tainted: write(notifications_queue.Message.PostID) {1}                     --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_______7] (BasicObject BasicType) string
            --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                       --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[________8] (BasicObject BasicType) postID_UPLOAD_SVC int64
             --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_________9] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
            --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                       --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[________8] (BasicObject BasicType) postID_UPLOAD_SVC int64
             --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_________9] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
         --> w-tainted: write(notifications_queue.Message.PostID) {1}                 --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[_____5] (BasicObject BasicType) PostID string
          --> w-tainted: write(notifications_queue.Message.PostID) {1}                   --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[______6] (Reference FieldType) ref <PostID string> @ NotifyService
           --> w-tainted: write(notifications_queue.Message.PostID) {1}                     --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_______7] (BasicObject BasicType) string
            --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                       --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[________8] (BasicObject BasicType) postID_UPLOAD_SVC int64
             --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_________9] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
            --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                       --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[________8] (BasicObject BasicType) postID_UPLOAD_SVC int64
             --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_________9] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService

    --> r-tainted: read(posts_db.Post) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = post, collection = post}
     --> r-tainted: read(posts_db.Post) {1}
[_1] (StructObject UserType) post postnotification.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification.Creator struct{Username string}}

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) res bool

[0] (InterfaceObject UserType) err .error

[0] (StructObject UserType) analytics postnotification.Analytics struct{PostID int64}
     --> r-tainted: read(analytics_db.Analytics) {1}
[_1] (Reference UserType) ref <analytics postnotification.Analytics struct{PostID int64}> @ AnalyticsService

[0] (InterfaceObject UserType) err .error
[_1] (Reference BasicType) ref <nil> @ AnalyticsService

