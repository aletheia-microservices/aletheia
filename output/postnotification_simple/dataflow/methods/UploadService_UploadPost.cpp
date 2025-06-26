[0] (PointerObject PointerType) u (*postnotification_simple.UploadServiceImpl struct{storageService postnotification_simple.StorageService, analyticsService postnotification_simple.AnalyticsService, notificationsQueue Queue})
[_1] (StructObject UserType) postnotification_simple.UploadServiceImpl struct{storageService postnotification_simple.StorageService, analyticsService postnotification_simple.AnalyticsService, notificationsQueue Queue}
[__2] (FieldObject FieldType) analyticsService postnotification_simple.AnalyticsService
[___3] (ServiceObject ServiceType) analyticsService postnotification_simple.AnalyticsService
[__2] (FieldObject FieldType) notificationsQueue Queue
[___3] (BlueprintBackendObject BlueprintBackendType) notificationsQueue Queue
[__2] (FieldObject FieldType) storageService postnotification_simple.StorageService
[___3] (ServiceObject ServiceType) storageService postnotification_simple.StorageService

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) username string

    --> w-tainted: write(posts_db.Post.Text) {1}
[0] (BasicObject BasicType) text string

    --> w-tainted: write(notifications_queue.Message.ReqID, posts_db.Post.ReqID) {2}       --> w-tainted: write(notifications_queue.Message.ReqID, posts_db.Post.ReqID) {2} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[0] (BasicObject BasicType) reqID int64

    --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID) {3}       --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID) {3} --> r-tainted: read(notifications_queue.Message.PostID_MESSAGE, posts_db._.postid, posts_db.Post.PostID) {3}
[0] (BasicObject BasicType) postID_UploadSVC int64
     --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}         --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[_1] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService

[0] (InterfaceObject UserType) _ .error
[_1] (Reference UserType) ref <err .error> @ StorageService

    --> w-tainted: write(notifications_queue.Message) {1}       --> w-tainted: write(notifications_queue.Message) {1} --> r-tainted: read(notifications_queue.Message) {1}
[0] (StructObject UserType) message postnotification_simple.Message struct{ReqID int64, PostID_MESSAGE int64, Timestamp int64}
     --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE) {1}         --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE) {1} --> r-tainted: read(notifications_queue.Message.PostID_MESSAGE, posts_db._.postid, posts_db.Post.PostID) {3}
[_1] (FieldObject FieldType) PostID_MESSAGE int64
      --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID) {3}           --> w-tainted: write(notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID) {3} --> r-tainted: read(notifications_queue.Message.PostID_MESSAGE, posts_db._.postid, posts_db.Post.PostID) {3}
[__2] (BasicObject BasicType) postID_UploadSVC int64
       --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}             --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[___3] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
     --> w-tainted: write(notifications_queue.Message.ReqID) {1}         --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[_1] (FieldObject FieldType) ReqID int64
      --> w-tainted: write(notifications_queue.Message.ReqID, posts_db.Post.ReqID) {2}           --> w-tainted: write(notifications_queue.Message.ReqID, posts_db.Post.ReqID) {2} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[__2] (BasicObject BasicType) reqID int64

[0] (BasicObject BasicType) _ bool

[0] (InterfaceObject UserType) err .error

