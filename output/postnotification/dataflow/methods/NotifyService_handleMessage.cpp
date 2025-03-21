[0] (PointerObject PointerType) n (*postnotification.NotifyServiceImpl struct{storageService postnotification.StorageService, notificationsQueue Queue, num_workers int})
[_1] (StructObject UserType) postnotification.NotifyServiceImpl struct{storageService postnotification.StorageService, notificationsQueue Queue, num_workers 4 int}
[__2] (FieldObject FieldType) notificationsQueue Queue
[___3] (BlueprintBackendObject BlueprintBackendType) notificationsQueue Queue
[__2] (FieldObject FieldType) num_workers 4 int
[___3] (BasicObject BasicType) 4 int
[__2] (FieldObject FieldType) storageService postnotification.StorageService
[___3] (ServiceObject ServiceType) storageService postnotification.StorageService

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ NotifyService

    --> w-tainted: write(notifications_queue.Message) {1}       --> w-tainted: write(notifications_queue.Message) {1} --> r-tainted: read(notifications_queue.Message) {1}
[0] (StructObject UserType) message postnotification.Message struct{ReqID string, PostID string, Timestamp string}
     --> w-tainted: write(notifications_queue.Message) {1}         --> w-tainted: write(notifications_queue.Message) {1} --> r-tainted: read(notifications_queue.Message) {1}
[_1] (Reference UserType) ref <notifyEvent postnotification.Message struct{ReqID string, PostID string, Timestamp string}> @ NotifyService
      --> w-tainted: write(notifications_queue.Message) {1}           --> w-tainted: write(notifications_queue.Message) {1} --> r-tainted: read(notifications_queue.Message) {1}
[__2] (Reference UserType) ref <message postnotification.Message struct{ReqID string, PostID string, Timestamp string}> @ NotifyService
       --> w-tainted: write(notifications_queue.Message.PostID) {1}             --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[___3] (FieldObject FieldType) PostID string
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
       --> w-tainted: write(notifications_queue.Message.ReqID) {1}             --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[___3] (FieldObject FieldType) ReqID string
        --> w-tainted: write(notifications_queue.Message.ReqID) {1}               --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[____4] (BasicObject BasicType) string
         --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}                 --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[_____5] (BasicObject BasicType) reqID int64
         --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}                 --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[_____5] (BasicObject BasicType) reqID int64
     --> w-tainted: write(notifications_queue.Message) {1}         --> w-tainted: write(notifications_queue.Message) {1} --> r-tainted: read(notifications_queue.Message) {1}
[_1] (Reference UserType) ref <message postnotification.Message struct{ReqID string, PostID string, Timestamp string}> @ NotifyService
      --> w-tainted: write(notifications_queue.Message.PostID) {1}           --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[__2] (FieldObject FieldType) PostID string
       --> w-tainted: write(notifications_queue.Message.PostID) {1}             --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[___3] (BasicObject BasicType) string
        --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}               --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[____4] (BasicObject BasicType) postID_UPLOAD_SVC int64
         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                 --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_____5] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
        --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}               --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[____4] (BasicObject BasicType) postID_UPLOAD_SVC int64
         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                 --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_____5] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
      --> w-tainted: write(notifications_queue.Message.ReqID) {1}           --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[__2] (FieldObject FieldType) ReqID string
       --> w-tainted: write(notifications_queue.Message.ReqID) {1}             --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[___3] (BasicObject BasicType) string
        --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}               --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[____4] (BasicObject BasicType) reqID int64
        --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}               --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[____4] (BasicObject BasicType) reqID int64
     --> w-tainted: write(notifications_queue.Message.PostID) {1}
[_1] (FieldObject FieldType) PostID string
      --> w-tainted: write(notifications_queue.Message.PostID) {1}           --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[__2] (Reference FieldType) ref <PostID string> @ NotifyService
       --> w-tainted: write(notifications_queue.Message.PostID) {1}             --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[___3] (BasicObject BasicType) string
        --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}               --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[____4] (BasicObject BasicType) postID_UPLOAD_SVC int64
         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                 --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_____5] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
        --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}               --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[____4] (BasicObject BasicType) postID_UPLOAD_SVC int64
         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                 --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_____5] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
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
     --> w-tainted: write(notifications_queue.Message.ReqID) {1}
[_1] (FieldObject FieldType) ReqID string
      --> w-tainted: write(notifications_queue.Message.ReqID) {1}           --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[__2] (Reference FieldType) ref <ReqID string> @ NotifyService
       --> w-tainted: write(notifications_queue.Message.ReqID) {1}             --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[___3] (BasicObject BasicType) string
        --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}               --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[____4] (BasicObject BasicType) reqID int64
        --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}               --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[____4] (BasicObject BasicType) reqID int64
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

[0] (BasicObject BasicType) reqID int64
     --> w-tainted: write(notifications_queue.Message.ReqID) {1}
[_1] (BasicObject BasicType) ReqID string
      --> w-tainted: write(notifications_queue.Message.ReqID) {1}           --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[__2] (Reference FieldType) ref <ReqID string> @ NotifyService
       --> w-tainted: write(notifications_queue.Message.ReqID) {1}             --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[___3] (BasicObject BasicType) string
        --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}               --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[____4] (BasicObject BasicType) reqID int64
        --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}               --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[____4] (BasicObject BasicType) reqID int64
     --> w-tainted: write(notifications_queue.Message.ReqID) {1}
[_1] (BasicObject BasicType) ReqID string
      --> w-tainted: write(notifications_queue.Message.ReqID) {1}           --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[__2] (Reference FieldType) ref <ReqID string> @ NotifyService
       --> w-tainted: write(notifications_queue.Message.ReqID) {1}             --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[___3] (BasicObject BasicType) string
        --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}               --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[____4] (BasicObject BasicType) reqID int64
        --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}               --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[____4] (BasicObject BasicType) reqID int64

[0] (InterfaceObject UserType) err .error
     --> w-tainted: write(notifications_queue.Message.ReqID) {1}
[_1] (BasicObject BasicType) ReqID string
      --> w-tainted: write(notifications_queue.Message.ReqID) {1}           --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[__2] (Reference FieldType) ref <ReqID string> @ NotifyService
       --> w-tainted: write(notifications_queue.Message.ReqID) {1}             --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[___3] (BasicObject BasicType) string
        --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}               --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[____4] (BasicObject BasicType) reqID int64
        --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}               --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[____4] (BasicObject BasicType) reqID int64

    --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[0] (BasicObject BasicType) postID_NOTIFY_SVC int64
     --> w-tainted: write(notifications_queue.Message.PostID) {1}         --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[_1] (BasicObject BasicType) PostID string
      --> w-tainted: write(notifications_queue.Message.PostID) {1}           --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[__2] (Reference FieldType) ref <PostID string> @ NotifyService
       --> w-tainted: write(notifications_queue.Message.PostID) {1}             --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[___3] (BasicObject BasicType) string
        --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}               --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[____4] (BasicObject BasicType) postID_UPLOAD_SVC int64
         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                 --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_____5] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
        --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}               --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[____4] (BasicObject BasicType) postID_UPLOAD_SVC int64
         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                 --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_____5] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
     --> w-tainted: write(notifications_queue.Message.PostID) {1}         --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[_1] (BasicObject BasicType) PostID string
      --> w-tainted: write(notifications_queue.Message.PostID) {1}           --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[__2] (Reference FieldType) ref <PostID string> @ NotifyService
       --> w-tainted: write(notifications_queue.Message.PostID) {1}             --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[___3] (BasicObject BasicType) string
        --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}               --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[____4] (BasicObject BasicType) postID_UPLOAD_SVC int64
         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                 --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_____5] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
        --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}               --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[____4] (BasicObject BasicType) postID_UPLOAD_SVC int64
         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                 --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_____5] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService

[0] (InterfaceObject UserType) err .error
     --> w-tainted: write(notifications_queue.Message.PostID) {1}         --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(posts_db.Post.PostID, analytics_db.Analytics.PostID) {2}
[_1] (BasicObject BasicType) PostID string
      --> w-tainted: write(notifications_queue.Message.PostID) {1}           --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[__2] (Reference FieldType) ref <PostID string> @ NotifyService
       --> w-tainted: write(notifications_queue.Message.PostID) {1}             --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[___3] (BasicObject BasicType) string
        --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}               --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[____4] (BasicObject BasicType) postID_UPLOAD_SVC int64
         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                 --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_____5] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
        --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}               --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[____4] (BasicObject BasicType) postID_UPLOAD_SVC int64
         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2}                 --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.value) {2} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_____5] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService

[0] (StructObject UserType) _ postnotification.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification.Creator struct{Username string}}
     --> r-tainted: read(posts_db.Post) {1}
[_1] (Reference UserType) ref <post postnotification.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification.Creator struct{Username string}}> @ StorageService

[0] (StructObject UserType) _ postnotification.Analytics struct{PostID int64}
[_1] (Reference UserType) ref <analytics postnotification.Analytics struct{PostID int64}> @ StorageService
      --> r-tainted: read(analytics_db.Analytics) {1}
[__2] (Reference UserType) ref <analytics postnotification.Analytics struct{PostID int64}> @ AnalyticsService

[0] (InterfaceObject UserType) err .error
[_1] (Reference UserType) ref <err .error> @ StorageService
[__2] (Reference BasicType) ref <nil> @ AnalyticsService

