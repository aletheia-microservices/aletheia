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

[] (BasicObject BasicType) reqID int64
[_] (Reference BasicType) ref <reqID int64> @ UploadService

    --> r-tainted: read(posts_cache.int64) {1}
[] (BasicObject BasicType) postID int64
     --> r-tainted: read(posts_cache.int64) {1}
[_] (Reference BasicType) ref <postID int64> @ UploadService

    --> r-tainted: read(posts_cache.Post) {1}
[] (StructObject UserType) post postnotification.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification.Creator struct{Username string}}
     --> r-tainted: read(posts_cache.PostID) {1}
[_] (FieldObject FieldType) PostID int64
      --> r-tainted: read(posts_cache.int64) {1}
[__] (BasicObject BasicType) PostID int64

    --> r-tainted: read(posts_cache.key, posts_cache.string, posts_cache.int64, posts_cache.int) {4}
[] (BasicObject BasicType) postIDStr string
     --> r-tainted: read(posts_cache.int64) {1}
[_] (BasicObject BasicType) postID int64
      --> r-tainted: read(posts_cache.int64) {1}
[__] (Reference BasicType) ref <postID int64> @ UploadService
     --> r-tainted: read(posts_cache.int) {1}
[_] (BasicObject BasicType) 10 int

[] (StructObject UserType) media postnotification.Media struct{MediaID int64, PostID int64, Content string}

    --> r-tainted: read(media_db._.mediaid) {1}
[] (FieldObject FieldType) mediaID int64
     --> r-tainted: read(media_db._.mediaid) {1}
[_] (BasicObject BasicType) PostID int64

[] (StructObject UserType) media postnotification.Media struct{MediaID int64, PostID int64, Content string}
     --> r-tainted: read(media_db.Media) {1}
[_] (Reference UserType) ref <media postnotification.Media struct{MediaID int64, PostID int64, Content string}> @ MediaService

[] (InterfaceObject UserType) _ .error
[_] (Reference UserType) ref <err .error> @ MediaService

