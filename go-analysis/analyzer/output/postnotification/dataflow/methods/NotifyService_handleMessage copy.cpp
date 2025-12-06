[] (PointerObject PointerType) n (*postnotification.NotifyServiceImpl struct{storageService postnotification.StorageService, notificationsQueue Queue, num_workers int})
[_] (StructObject UserType) postnotification.NotifyServiceImpl struct{storageService postnotification.StorageService, notificationsQueue Queue, num_workers 4 int}
[__] (FieldObject FieldType) notificationsQueue Queue
[___] (BlueprintBackendObject BlueprintBackendType) notificationsQueue Queue
[__] (FieldObject FieldType) num_workers 4 int
[___] (BasicObject BasicType) 4 int
[__] (FieldObject FieldType) storageService postnotification.StorageService
[___] (ServiceObject ServiceType) storageService postnotification.StorageService

[] (InterfaceObject UserType) ctx context.Context
[_] (Reference UserType) ref <ctx context.Context> @ NotifyService

[] (StructObject UserType) message postnotification.Message struct{ReqID string, PostID string, Timestamp string}
[_] (Reference UserType) ref <workerMessage postnotification.Message struct{ReqID string, PostID string, Timestamp string}> @ NotifyService
      --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[__] (FieldObject FieldType) PostID string
       --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[___] (BasicObject BasicType) PostID string
        --> w-tainted: write(notifications_queue.Message.PostID) {1}               --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[____] (InterfaceObject InterfaceType) PostID interface{}
[__] (FieldObject FieldType) ReqID string
[___] (BasicObject BasicType) ReqID string
        --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}               --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[____] (InterfaceObject InterfaceType) ReqID interface{}
[__] (FieldObject FieldType) Timestamp string
[___] (BasicObject BasicType) Timestamp string
        --> r-tainted: read(notifications_queue.Message) {1}
[____] (InterfaceObject InterfaceType) Timestamp interface{}
[_] (FieldObject FieldType) PostID string
      --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[__] (Reference FieldType) ref <PostID string> @ NotifyService
       --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[___] (BasicObject BasicType) PostID string
        --> w-tainted: write(notifications_queue.Message.PostID) {1}               --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[____] (InterfaceObject InterfaceType) PostID interface{}
      --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[__] (BasicObject BasicType) PostID string
       --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[___] (Reference FieldType) ref <PostID string> @ NotifyService
        --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[____] (BasicObject BasicType) PostID string
         --> w-tainted: write(notifications_queue.Message.PostID) {1}                 --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_____] (InterfaceObject InterfaceType) PostID interface{}
[_] (FieldObject FieldType) ReqID string
[__] (Reference FieldType) ref <ReqID string> @ NotifyService
[___] (BasicObject BasicType) ReqID string
        --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}               --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[____] (InterfaceObject InterfaceType) ReqID interface{}
[__] (BasicObject BasicType) ReqID string
[___] (Reference FieldType) ref <ReqID string> @ NotifyService
[____] (BasicObject BasicType) ReqID string
         --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}                 --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[_____] (InterfaceObject InterfaceType) ReqID interface{}

[] (BasicObject BasicType) reqID int64
[_] (BasicObject BasicType) ReqID string
[__] (Reference FieldType) ref <ReqID string> @ NotifyService
[___] (BasicObject BasicType) ReqID string
        --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}               --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[____] (InterfaceObject InterfaceType) ReqID interface{}

[] (InterfaceObject UserType) err .error
[_] (BasicObject BasicType) ReqID string
[__] (Reference FieldType) ref <ReqID string> @ NotifyService
[___] (BasicObject BasicType) ReqID string
        --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}               --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[____] (InterfaceObject InterfaceType) ReqID interface{}

    --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[] (BasicObject BasicType) postID_NOTIFY_SVC int64
     --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[_] (BasicObject BasicType) PostID string
      --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[__] (Reference FieldType) ref <PostID string> @ NotifyService
       --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[___] (BasicObject BasicType) PostID string
        --> w-tainted: write(notifications_queue.Message.PostID) {1}               --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[____] (InterfaceObject InterfaceType) PostID interface{}

[] (InterfaceObject UserType) err .error
     --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[_] (BasicObject BasicType) PostID string
      --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[__] (Reference FieldType) ref <PostID string> @ NotifyService
       --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[___] (BasicObject BasicType) PostID string
        --> w-tainted: write(notifications_queue.Message.PostID) {1}               --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[____] (InterfaceObject InterfaceType) PostID interface{}

[] (StructObject UserType) _ postnotification.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification.Creator struct{Username string}}
     --> r-tainted: read(posts_db.Post) {1}
[_] (Reference UserType) ref <post postnotification.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification.Creator struct{Username string}}> @ StorageService

[] (StructObject UserType) _ postnotification.Analytics struct{PostID int64}
[_] (Reference UserType) ref <analytics postnotification.Analytics struct{PostID int64}> @ StorageService
      --> r-tainted: read(analytics_db.Analytics) {1}
[__] (Reference UserType) ref <analytics postnotification.Analytics struct{PostID int64}> @ AnalyticsService

[] (InterfaceObject UserType) err .error
[_] (Reference UserType) ref <err .error> @ StorageService
[__] (Reference BasicType) ref <nil> @ AnalyticsService

