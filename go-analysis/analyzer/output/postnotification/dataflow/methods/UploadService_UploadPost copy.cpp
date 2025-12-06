[0] (PointerObject PointerType) u (*postnotification.UploadServiceImpl struct{storageService postnotification.StorageService, mediaService postnotification.MediaService, notificationsQueue Queue, timelineCache Cache})
[_1] (StructObject UserType) postnotification.UploadServiceImpl struct{storageService postnotification.StorageService, mediaService postnotification.MediaService, notificationsQueue Queue, timelineCache Cache}
[__2] (FieldObject FieldType) mediaService postnotification.MediaService
[___3] (ServiceObject ServiceType) mediaService postnotification.MediaService
[__2] (FieldObject FieldType) notificationsQueue Queue
[___3] (BlueprintBackendObject BlueprintBackendType) notificationsQueue Queue
[__2] (FieldObject FieldType) storageService postnotification.StorageService
[___3] (ServiceObject ServiceType) storageService postnotification.StorageService
[__2] (FieldObject FieldType) timelineCache Cache
[___3] (BlueprintBackendObject BlueprintBackendType) timelineCache Cache

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) username string

    --> w-tainted: write(posts_db.Post.Text) {1}
[0] (BasicObject BasicType) text string

    --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}       --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[0] (BasicObject BasicType) reqID int64
     --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}         --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[_1] (Reference InterfaceType) ref <ReqID interface{}> @ NotifyService

    --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}       --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[0] (BasicObject BasicType) postIDDDDD int64

    --> w-tainted: write(media_db.Media) {1}
[0] (StructObject UserType) media postnotification.Media struct{MediaID int64, Content "HELLO WORLD!" untyped string}
     --> w-tainted: write(media_db.Media.Content) {1}
[_1] (FieldObject FieldType) Content "HELLO WORLD!" untyped string
      --> w-tainted: write(media_db.Media.Content) {1}
[__2] (BasicObject BasicType) HELLO_WORLD_CONST "HELLO WORLD!" untyped string

    --> w-tainted: write(posts_db.Post.Timestamp) {1}
[0] (BasicObject BasicType) timestamp int64

    --> w-tainted: write(posts_db.Post.Mentions) {1}
[0] (ArrayObject ArrayType) mentions []string
     --> w-tainted: write(posts_db.Post.Mentions) {1}
[_1] (BasicObject BasicType) "alice" string
     --> w-tainted: write(posts_db.Post.Mentions) {1}
[_1] (BasicObject BasicType) "bob" string

    --> w-tainted: write(posts_db.Post) {1}
[0] (StructObject UserType) post postnotification.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification.Creator struct{Username "some username" string}}
     --> w-tainted: write(posts_db.Post.Creator) {1}
[_1] (FieldObject FieldType) Creator postnotification.Creator struct{Username "some username" string}
      --> w-tainted: write(posts_db.Post.Creator) {1}
[__2] (StructObject UserType) postnotification.Creator struct{Username "some username" string}
       --> w-tainted: write(posts_db.Post.Creator.Username) {1}
[___3] (FieldObject FieldType) Username "some username" string
        --> w-tainted: write(posts_db.Post.Creator.Username) {1}
[____4] (BasicObject BasicType) "some username" string
     --> w-tainted: write(posts_db.Post.Mentions) {1}
[_1] (FieldObject FieldType) Mentions []string
      --> w-tainted: write(posts_db.Post.Mentions) {1}
[__2] (ArrayObject ArrayType) mentions []string
       --> w-tainted: write(posts_db.Post.Mentions) {1}
[___3] (BasicObject BasicType) "alice" string
       --> w-tainted: write(posts_db.Post.Mentions) {1}
[___3] (BasicObject BasicType) "bob" string
     --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}         --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[_1] (FieldObject FieldType) PostID int64
      --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}           --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[__2] (BasicObject BasicType) postID_UPLOAD_SVC int64
       --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}             --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[___3] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
        --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}               --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[____4] (Reference FieldType) ref <PostID int64> @ UploadService
         --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}                 --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[_____5] (BasicObject BasicType) postIDDDDD int64
       --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {3}             --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {3} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[___3] (Reference InterfaceType) ref <PostID interface{}> @ NotifyService
     --> w-tainted: write(posts_db.Post.ReqID) {1}
[_1] (FieldObject FieldType) ReqID int64
      --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}           --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[__2] (BasicObject BasicType) reqID int64
       --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}             --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[___3] (Reference InterfaceType) ref <ReqID interface{}> @ NotifyService
     --> w-tainted: write(posts_db.Post.Text) {1}
[_1] (FieldObject FieldType) Text string
      --> w-tainted: write(posts_db.Post.Text) {1}
[__2] (BasicObject BasicType) text string
     --> w-tainted: write(posts_db.Post.Timestamp) {1}
[_1] (FieldObject FieldType) Timestamp int64
      --> w-tainted: write(posts_db.Post.Timestamp) {1}
[__2] (BasicObject BasicType) timestamp int64

    --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}       --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[0] (BasicObject BasicType) postID_UPLOAD_SVC int64
     --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}         --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[_1] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
      --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}           --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[__2] (Reference FieldType) ref <PostID int64> @ UploadService
       --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}             --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[___3] (BasicObject BasicType) postIDDDDD int64
     --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {3}         --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {3} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[_1] (Reference InterfaceType) ref <PostID interface{}> @ NotifyService

[0] (InterfaceObject UserType) _ .error
[_1] (Reference UserType) ref <err .error> @ StorageService

    --> w-tainted: write(notifications_queue.Message) {1}       --> w-tainted: write(notifications_queue.Message) {1} --> r-tainted: read(notifications_queue.Message) {1}
[0] (StructObject UserType) message postnotification.Message struct{ReqID string, PostID string, Timestamp string}
     --> w-tainted: write(notifications_queue.Message.PostID) {1}         --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID) {1}
[_1] (FieldObject FieldType) PostID string
      --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {3}           --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {3} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[__2] (Reference InterfaceType) ref <PostID interface{}> @ NotifyService
      --> w-tainted: write(notifications_queue.Message.PostID) {1}           --> w-tainted: write(notifications_queue.Message.PostID) {1} --> r-tainted: read(notifications_queue.Message.PostID) {1}
[__2] (BasicObject BasicType) string
       --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {3}             --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {3} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[___3] (Reference InterfaceType) ref <PostID interface{}> @ NotifyService
       --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}             --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[___3] (BasicObject BasicType) postID_UPLOAD_SVC int64
        --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}               --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[____4] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
         --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}                 --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[_____5] (Reference FieldType) ref <PostID int64> @ UploadService
          --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}                   --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[______6] (BasicObject BasicType) postIDDDDD int64
        --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {3}               --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {3} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[____4] (Reference InterfaceType) ref <PostID interface{}> @ NotifyService
     --> w-tainted: write(notifications_queue.Message.ReqID) {1}         --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[_1] (FieldObject FieldType) ReqID string
      --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}           --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[__2] (Reference InterfaceType) ref <ReqID interface{}> @ NotifyService
      --> w-tainted: write(notifications_queue.Message.ReqID) {1}           --> w-tainted: write(notifications_queue.Message.ReqID) {1} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[__2] (BasicObject BasicType) string
       --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}             --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[___3] (Reference InterfaceType) ref <ReqID interface{}> @ NotifyService
       --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}             --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[___3] (BasicObject BasicType) reqID int64
        --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}               --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[____4] (Reference InterfaceType) ref <ReqID interface{}> @ NotifyService

[0] (BasicObject BasicType) _ bool

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) reqIDStr string
     --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}         --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[_1] (BasicObject BasicType) reqID int64
      --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}           --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[__2] (Reference InterfaceType) ref <ReqID interface{}> @ NotifyService
[_1] (BasicObject BasicType) 10 int

    --> w-tainted: write(timeline_cache.key, timeline_cache.value) {2}
[0] (StructObject UserType) timeline postnotification.Timeline struct{ReqID int64, PostID int64}
     --> w-tainted: write(timeline_cache.key, timeline_cache.value) {2}
[_1] (FieldObject FieldType) PostID int64
      --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}           --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[__2] (BasicObject BasicType) postID_UPLOAD_SVC int64
       --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}             --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[___3] (Reference BasicType) ref <postID_STORAGE_SVC int64> @ StorageService
        --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}               --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[____4] (Reference FieldType) ref <PostID int64> @ UploadService
         --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6}                 --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[_____5] (BasicObject BasicType) postIDDDDD int64
       --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {3}             --> w-tainted: write(notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {3} --> r-tainted: read(notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {3}
[___3] (Reference InterfaceType) ref <PostID interface{}> @ NotifyService
     --> w-tainted: write(timeline_cache.key, timeline_cache.value) {2}
[_1] (FieldObject FieldType) ReqID int64
      --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}           --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[__2] (BasicObject BasicType) reqID int64
       --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4}             --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {4} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[___3] (Reference InterfaceType) ref <ReqID interface{}> @ NotifyService

