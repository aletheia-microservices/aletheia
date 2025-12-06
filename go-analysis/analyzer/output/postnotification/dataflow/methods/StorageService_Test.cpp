[] (PointerObject PointerType) s (*postnotification.StorageServiceImpl struct{analyticsService postnotification.AnalyticsService, mediaService postnotification.MediaService, posts_cache Cache, postsDb NoSQLDatabase, analyticsQueue Queue})
[_] (StructObject UserType) postnotification.StorageServiceImpl struct{analyticsService postnotification.AnalyticsService, mediaService postnotification.MediaService, posts_cache Cache, postsDb NoSQLDatabase, analyticsQueue Queue}
[__] (FieldObject FieldType) analyticsQueue Queue
[___] (BlueprintBackendObject BlueprintBackendType) analyticsQueue Queue
[__] (FieldObject FieldType) analyticsService postnotification.AnalyticsService
[___] (ServiceObject ServiceType) analyticsService postnotification.AnalyticsService
[__] (FieldObject FieldType) mediaService postnotification.MediaService
[___] (ServiceObject ServiceType) mediaService postnotification.MediaService
[__] (FieldObject FieldType) postsDb NoSQLDatabase
[___] (BlueprintBackendObject BlueprintBackendType) postsDb NoSQLDatabase
[__] (FieldObject FieldType) posts_cache Cache
[___] (BlueprintBackendObject BlueprintBackendType) posts_cache Cache

[] (InterfaceObject UserType) ctx context.Context
[_] (Reference UserType) ref <ctx context.Context> @ UploadService

    --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}       --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {3}
[] (BasicObject BasicType) postID_STORAGE_SVC_READ_TEST int64
     --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}         --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, posts_db.Post.PostID, notifications_queue.Message.PostID) {3}
[_] (Reference BasicType) ref <postID_UPLOAD_SVC int64> @ UploadService
      --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}           --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, posts_db.Post.PostID, notifications_queue.Message.PostID) {3}
[__] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
       --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}             --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, posts_db.Post.PostID, notifications_queue.Message.PostID) {3}
[___] (Reference FieldType) ref <PostID int64> @ UploadService
        --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}               --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, posts_db.Post.PostID, notifications_queue.Message.PostID) {3}
[____] (BasicObject BasicType) postID int64
      --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {3}           --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {3} --> r-tainted: read(posts_db.Post.PostID, notifications_queue.Message.PostID, analytics_db.Analytics.PostID) {3}
[__] (Reference InterfaceType) ref <PostID interface{}> @ NotifyService

    --> r-tainted: read(posts_db.Post) {1}
[] (StructObject UserType) post postnotification.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification.Creator struct{Username string}}

[] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = post, collection = post}

[] (InterfaceObject UserType) _ .error

[] (SliceObject UserType) postsQuery primitive.D
[_] (StructObject StructType) struct{Key "postid" string, Key "postid" string, Value int64, Value int64}
[__] (FieldObject FieldType) Key "postid" string
[___] (BasicObject BasicType) "postid" string
[__] (FieldObject FieldType) Value int64
       --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}             --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {3}
[___] (BasicObject BasicType) postID_STORAGE_SVC_READ_TEST int64
        --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}               --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, posts_db.Post.PostID, notifications_queue.Message.PostID) {3}
[____] (Reference BasicType) ref <postID_UPLOAD_SVC int64> @ UploadService
         --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}                 --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, posts_db.Post.PostID, notifications_queue.Message.PostID) {3}
[_____] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
          --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}                   --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, posts_db.Post.PostID, notifications_queue.Message.PostID) {3}
[______] (Reference FieldType) ref <PostID int64> @ UploadService
           --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}                     --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, posts_db.Post.PostID, notifications_queue.Message.PostID) {3}
[_______] (BasicObject BasicType) postID int64
         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {3}                 --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {3} --> r-tainted: read(posts_db.Post.PostID, notifications_queue.Message.PostID, analytics_db.Analytics.PostID) {3}
[_____] (Reference InterfaceType) ref <PostID interface{}> @ NotifyService

    --> r-tainted: read(posts_db.Post) {1}
[] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = post, collection = post}
     --> r-tainted: read(posts_db.Post) {1}
[_] (StructObject UserType) post postnotification.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification.Creator struct{Username string}}

[] (InterfaceObject UserType) _ .error

