[0] (PointerObject PointerType) s (*postnotification_simple.StorageServiceImpl struct{postsDb NoSQLDatabase, analyticsQueue Queue})
[_1] (StructObject UserType) postnotification_simple.StorageServiceImpl struct{postsDb NoSQLDatabase, analyticsQueue Queue}
[__2] (FieldObject FieldType) analyticsQueue Queue
[___3] (BlueprintBackendObject BlueprintBackendType) analyticsQueue Queue
[__2] (FieldObject FieldType) postsDb NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) postsDb NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ UploadService
[_1] (Reference UserType) ref <ctx context.Context> @ NotifyService

    --> w-tainted: write(notifications_queue.Message.ReqID) {1}       --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[0] (BasicObject BasicType) reqID int64
[_1] (Reference BasicType) ref <reqID int64> @ UploadService
     --> w-tainted: write(notifications_queue.Message.ReqID) {1}         --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[_1] (Reference BasicType) ref <ReqID int64> @ NotifyService
      --> w-tainted: write(notifications_queue.Message.ReqID) {1}           --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[__2] (Reference FieldType) ref <ReqID int64> @ NotifyService
       --> w-tainted: write(notifications_queue.Message.ReqID, posts_db.Post.ReqID) {2}             --> w-tainted: write(notifications_queue.Message.ReqID, posts_db.Post.ReqID) {2} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[___3] (BasicObject BasicType) reqID int64

    --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE) {1}       --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE) {1} --> r-tainted: read(posts_db._.postid, posts_db.Post.PostID, notifications_queue.Message.PostID_MESSAGE) {3}
[0] (BasicObject BasicType) postID int64
     --> r-tainted: read(posts_db._.postid, analytics_db._.postid, posts_db.Post.PostID) {3}
[_1] (Reference BasicType) ref <postID int64> @ UploadService
     --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE) {1}         --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE) {1} --> r-tainted: read(notifications_queue.Message.PostID_MESSAGE, posts_db._.postid, posts_db.Post.PostID) {3}
[_1] (Reference BasicType) ref <PostID_MESSAGE int64> @ NotifyService
      --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE) {1}           --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE) {1} --> r-tainted: read(notifications_queue.Message.PostID_MESSAGE, posts_db._.postid, posts_db.Post.PostID) {3}
[__2] (Reference FieldType) ref <PostID_MESSAGE int64> @ NotifyService
       --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID) {3}             --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID) {3} --> r-tainted: read(notifications_queue.Message.PostID_MESSAGE, posts_db._.postid, posts_db.Post.PostID) {3}
[___3] (BasicObject BasicType) postID_UploadSVC int64
        --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}               --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[____4] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService

    --> r-tainted: read(posts_db.Post) {1}
[0] (StructObject UserType) post postnotification_simple.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification_simple.Creator struct{Username string}}

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = post, collection = post}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "postid" string, Key "postid" string, Value int64, Value int64}
[__2] (FieldObject FieldType) Key "postid" string
[___3] (BasicObject BasicType) "postid" string
[__2] (FieldObject FieldType) Value int64
       --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE) {1}             --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE) {1} --> r-tainted: read(posts_db._.postid, posts_db.Post.PostID, notifications_queue.Message.PostID_MESSAGE) {3}
[___3] (BasicObject BasicType) postID int64
        --> r-tainted: read(posts_db._.postid, analytics_db._.postid, posts_db.Post.PostID) {3}
[____4] (Reference BasicType) ref <postID int64> @ UploadService
        --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE) {1}               --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE) {1} --> r-tainted: read(notifications_queue.Message.PostID_MESSAGE, posts_db._.postid, posts_db.Post.PostID) {3}
[____4] (Reference BasicType) ref <PostID_MESSAGE int64> @ NotifyService
         --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE) {1}                 --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE) {1} --> r-tainted: read(notifications_queue.Message.PostID_MESSAGE, posts_db._.postid, posts_db.Post.PostID) {3}
[_____5] (Reference FieldType) ref <PostID_MESSAGE int64> @ NotifyService
          --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID) {3}                   --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID) {3} --> r-tainted: read(notifications_queue.Message.PostID_MESSAGE, posts_db._.postid, posts_db.Post.PostID) {3}
[______6] (BasicObject BasicType) postID_UploadSVC int64
           --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}                     --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[_______7] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService

    --> r-tainted: read(posts_db._, posts_db.Post) {2}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = post, collection = post}
     --> r-tainted: read(posts_db.Post) {1}
[_1] (StructObject UserType) post postnotification_simple.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification_simple.Creator struct{Username string}}

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) res bool

[0] (InterfaceObject UserType) err .error

