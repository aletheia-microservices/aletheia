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

    --> w-tainted: write(posts_cache.key, posts_cache.value, posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {6}       --> w-tainted: write(posts_cache.key, posts_cache.value, posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[] (BasicObject BasicType) reqID int64
     --> w-tainted: write(posts_db.Post.ReqID, posts_cache.key, posts_cache.value, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {6}         --> w-tainted: write(posts_db.Post.ReqID, posts_cache.key, posts_cache.value, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[_] (Reference BasicType) ref <reqID int64> @ UploadService
      --> w-tainted: write(posts_cache.key, posts_cache.value, posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {6}           --> w-tainted: write(posts_cache.key, posts_cache.value, posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[__] (Reference InterfaceType) ref <ReqID interface{}> @ NotifyService

    --> w-tainted: write(posts_cache.key, posts_cache.value) {2}
[] (StructObject UserType) post postnotification.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification.Creator struct{Username string}}
     --> w-tainted: write(posts_cache.key, posts_cache.value, posts_db.Post) {3}
[_] (Reference UserType) ref <post postnotification.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification.Creator struct{Username "some username" string}}> @ UploadService
      --> w-tainted: write(posts_db.Post.Creator, posts_cache.key, posts_cache.value) {3}
[__] (FieldObject FieldType) Creator postnotification.Creator struct{Username "some username" string}
       --> w-tainted: write(posts_cache.key, posts_cache.value, posts_db.Post.Creator) {3}
[___] (StructObject UserType) postnotification.Creator struct{Username "some username" string}
        --> w-tainted: write(posts_db.Post.Creator.Username, posts_cache.key, posts_cache.value) {3}
[____] (FieldObject FieldType) Username "some username" string
         --> w-tainted: write(posts_db.Post.Creator.Username, posts_cache.key, posts_cache.value) {3}
[_____] (BasicObject BasicType) "some username" string
      --> w-tainted: write(posts_db.Post.Mentions, posts_cache.key, posts_cache.value) {3}
[__] (FieldObject FieldType) Mentions []string
       --> w-tainted: write(posts_db.Post.Mentions, posts_cache.key, posts_cache.value) {3}
[___] (ArrayObject ArrayType) mentions []string
        --> w-tainted: write(posts_cache.key, posts_cache.value, posts_db.Post.Mentions) {3}
[____] (BasicObject BasicType) "alice" string
        --> w-tainted: write(posts_cache.key, posts_cache.value, posts_db.Post.Mentions) {3}
[____] (BasicObject BasicType) "bob" string
      --> w-tainted: write(posts_db.Post.PostID, posts_cache.key, posts_cache.value, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {5}           --> w-tainted: write(posts_db.Post.PostID, posts_cache.key, posts_cache.value, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {5} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[__] (FieldObject FieldType) PostID int64
       --> w-tainted: write(media_db.Media.PostID, posts_db.Post.PostID, posts_cache.key, posts_cache.value, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {9}             --> w-tainted: write(media_db.Media.PostID, posts_db.Post.PostID, posts_cache.key, posts_cache.value, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {9} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[___] (BasicObject BasicType) postID int64
        --> w-tainted: write(media_db.Media.PostID, posts_cache.key, posts_cache.value, posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {9}               --> w-tainted: write(media_db.Media.PostID, posts_cache.key, posts_cache.value, posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {9} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {4}
[____] (Reference InterfaceType) ref <PostID interface{}> @ NotifyService
      --> w-tainted: write(posts_db.Post.ReqID, posts_cache.key, posts_cache.value) {3}
[__] (FieldObject FieldType) ReqID int64
       --> w-tainted: write(posts_db.Post.ReqID, posts_cache.key, posts_cache.value, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {6}             --> w-tainted: write(posts_db.Post.ReqID, posts_cache.key, posts_cache.value, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[___] (BasicObject BasicType) reqID int64
        --> w-tainted: write(posts_cache.key, posts_cache.value, posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {6}               --> w-tainted: write(posts_cache.key, posts_cache.value, posts_db.Post.ReqID, notifications_queue.Message.ReqID, timeline_cache.key, timeline_cache.value) {6} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[____] (Reference InterfaceType) ref <ReqID interface{}> @ NotifyService
      --> w-tainted: write(posts_db.Post.Text, posts_cache.key, posts_cache.value) {3}
[__] (FieldObject FieldType) Text string
       --> w-tainted: write(posts_db.Post.Text, posts_cache.key, posts_cache.value) {3}
[___] (BasicObject BasicType) text string
      --> w-tainted: write(posts_db.Post.Timestamp, posts_cache.key, posts_cache.value) {3}
[__] (FieldObject FieldType) Timestamp int64
       --> w-tainted: write(posts_db.Post.Timestamp, posts_cache.key, posts_cache.value) {3}
[___] (BasicObject BasicType) timestamp int64
     --> w-tainted: write(posts_cache.key, posts_cache.value, posts_db.Post.PostID) {3}
[_] (FieldObject FieldType) PostID int64
      --> w-tainted: write(posts_db.Post.PostID, posts_cache.key, posts_cache.value, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {5}           --> w-tainted: write(posts_db.Post.PostID, posts_cache.key, posts_cache.value, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {5} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[__] (Reference FieldType) ref <PostID int64> @ UploadService
       --> w-tainted: write(media_db.Media.PostID, posts_db.Post.PostID, posts_cache.key, posts_cache.value, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {9}             --> w-tainted: write(media_db.Media.PostID, posts_db.Post.PostID, posts_cache.key, posts_cache.value, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {9} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[___] (BasicObject BasicType) postID int64
        --> w-tainted: write(media_db.Media.PostID, posts_cache.key, posts_cache.value, posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {9}               --> w-tainted: write(media_db.Media.PostID, posts_cache.key, posts_cache.value, posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {9} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {4}
[____] (Reference InterfaceType) ref <PostID interface{}> @ NotifyService
      --> w-tainted: write(posts_cache.key, posts_cache.value, posts_db.Post.PostID) {3}
[__] (BasicObject BasicType) PostID int64
       --> w-tainted: write(posts_db.Post.PostID, posts_cache.key, posts_cache.value, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {5}             --> w-tainted: write(posts_db.Post.PostID, posts_cache.key, posts_cache.value, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {5} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[___] (Reference FieldType) ref <PostID int64> @ UploadService
        --> w-tainted: write(media_db.Media.PostID, posts_db.Post.PostID, posts_cache.key, posts_cache.value, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {9}               --> w-tainted: write(media_db.Media.PostID, posts_db.Post.PostID, posts_cache.key, posts_cache.value, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {9} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[____] (BasicObject BasicType) postID int64
         --> w-tainted: write(media_db.Media.PostID, posts_cache.key, posts_cache.value, posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {9}                 --> w-tainted: write(media_db.Media.PostID, posts_cache.key, posts_cache.value, posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {9} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {4}
[_____] (Reference InterfaceType) ref <PostID interface{}> @ NotifyService

[] (BasicObject BasicType) postIDStr string
     --> w-tainted: write(posts_cache.key, posts_cache.value, posts_db.Post.PostID) {3}
[_] (BasicObject BasicType) PostID int64
      --> w-tainted: write(posts_db.Post.PostID, posts_cache.key, posts_cache.value, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {5}           --> w-tainted: write(posts_db.Post.PostID, posts_cache.key, posts_cache.value, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {5} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[__] (Reference FieldType) ref <PostID int64> @ UploadService
       --> w-tainted: write(media_db.Media.PostID, posts_db.Post.PostID, posts_cache.key, posts_cache.value, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {9}             --> w-tainted: write(media_db.Media.PostID, posts_db.Post.PostID, posts_cache.key, posts_cache.value, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {9} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID) {2}
[___] (BasicObject BasicType) postID int64
        --> w-tainted: write(media_db.Media.PostID, posts_cache.key, posts_cache.value, posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {9}               --> w-tainted: write(media_db.Media.PostID, posts_cache.key, posts_cache.value, posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID, timeline_cache.key, timeline_cache.value) {9} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID, posts_db.Post.PostID, analytics_db.Analytics.PostID) {4}
[____] (Reference InterfaceType) ref <PostID interface{}> @ NotifyService
[_] (BasicObject BasicType) 10 int

