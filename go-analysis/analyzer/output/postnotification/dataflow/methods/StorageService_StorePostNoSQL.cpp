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
[_1] (Reference UserType) ref <ctx context.Context> @ UploadService

    --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}
[0] (BasicObject BasicType) reqID int64
     --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}         --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[_1] (Reference BasicType) ref <reqID int64> @ UploadService

    --> w-tainted: write(posts_db.Post) {1}
[0] (StructObject UserType) post postnotification.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification.Creator struct{Username string}}
     --> w-tainted: write(posts_db.Post) {1}
[_1] (Reference UserType) ref <postnotification.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification.Creator struct{Username "some username" string}}> @ UploadService
      --> w-tainted: write(posts_db.Post.Creator) {1}
[__2] (FieldObject FieldType) Creator postnotification.Creator struct{Username "some username" string}
       --> w-tainted: write(posts_db.Post.Creator) {1}
[___3] (StructObject UserType) postnotification.Creator struct{Username "some username" string}
        --> w-tainted: write(posts_db.Post.Creator.Username) {1}
[____4] (FieldObject FieldType) Username "some username" string
         --> w-tainted: write(posts_db.Post.Creator.Username) {1}
[_____5] (BasicObject BasicType) "some username" string
      --> w-tainted: write(posts_db.Post.Mentions) {1}
[__2] (FieldObject FieldType) Mentions []string
       --> w-tainted: write(posts_db.Post.Mentions) {1}
[___3] (ArrayObject ArrayType) mentions []string
        --> w-tainted: write(posts_db.Post.Mentions) {1}
[____4] (BasicObject BasicType) "alice" string
        --> w-tainted: write(posts_db.Post.Mentions) {1}
[____4] (BasicObject BasicType) "bob" string
      --> w-tainted: write(posts_db.Post.ReqID) {1}
[__2] (FieldObject FieldType) ReqID int64
       --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}             --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[___3] (BasicObject BasicType) reqID int64
      --> w-tainted: write(posts_db.Post.Text) {1}
[__2] (FieldObject FieldType) Text string
       --> w-tainted: write(posts_db.Post.Text) {1}
[___3] (BasicObject BasicType) text string
      --> w-tainted: write(posts_db.Post.Timestamp) {1}
[__2] (FieldObject FieldType) Timestamp int64
       --> w-tainted: write(posts_db.Post.Timestamp) {1}
[___3] (BasicObject BasicType) timestamp int64
     --> w-tainted: write(posts_db.Post.PostID) {1}
[_1] (FieldObject FieldType) PostID int64
[__2] (BasicObject BasicType) 2 int

    --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}       --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[0] (BasicObject BasicType) postID_STORAGE_SVC int64

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = post, collection = post}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

    --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage) {1}       --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage) {1} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage) {1}
[0] (StructObject UserType) message postnotification.TriggerAnalyticsMessage struct{PostID string}
     --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}         --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[_1] (FieldObject FieldType) PostID string
      --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}           --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[__2] (BasicObject BasicType) string
       --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {3}             --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {3} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[___3] (BasicObject BasicType) int64
       --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {3}             --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {3} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[___3] (BasicObject BasicType) int64

[0] (BasicObject BasicType) _ bool

[0] (InterfaceObject UserType) err .error

