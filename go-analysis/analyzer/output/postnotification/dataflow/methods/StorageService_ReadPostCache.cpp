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

[] (BasicObject BasicType) reqID int64

[] (BasicObject BasicType) postID int64

[] (StructObject UserType) post postnotification.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification.Creator struct{Username string}}

[] (BasicObject BasicType) postIDStr string
[_] (BasicObject BasicType) postID int64
[_] (BasicObject BasicType) 10 int

[] (BasicObject BasicType) _ bool

[] (InterfaceObject UserType) err .error

