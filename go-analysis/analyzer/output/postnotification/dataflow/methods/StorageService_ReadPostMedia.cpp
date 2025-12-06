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

[0] (BasicObject BasicType) reqID int64
[_1] (Reference BasicType) ref <reqID int64> @ UploadService

    --> r-tainted: read(posts_cache.key) {1}
[0] (BasicObject BasicType) postID int64
     --> r-tainted: read(posts_cache.key) {1}
[_1] (Reference BasicType) ref <postID int64> @ UploadService

    --> r-tainted: read(posts_cache.value) {1}
[0] (StructObject UserType) post postnotification.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification.Creator struct{Username string}}
     --> r-tainted: read(posts_cache.value, media_db._.mediaid) {2}
[_1] (FieldObject FieldType) PostID int64
      --> r-tainted: read(posts_cache.value, media_db._.mediaid) {2}
[__2] (BasicObject BasicType) PostID int64

    --> r-tainted: read(posts_cache.key) {1}
[0] (BasicObject BasicType) postIDStr string
     --> r-tainted: read(posts_cache.key) {1}
[_1] (BasicObject BasicType) postID int64
      --> r-tainted: read(posts_cache.key) {1}
[__2] (Reference BasicType) ref <postID int64> @ UploadService
     --> r-tainted: read(posts_cache.key) {1}
[_1] (BasicObject BasicType) 10 int
     --> r-tainted: read(posts_cache.key) {1}
[_1] (BasicObject BasicType) postID int64
      --> r-tainted: read(posts_cache.key) {1}
[__2] (Reference BasicType) ref <postID int64> @ UploadService
     --> r-tainted: read(posts_cache.key) {1}
[_1] (BasicObject BasicType) 10 int

[0] (StructObject UserType) media postnotification.Media struct{MediaID int64, Content string}

    --> r-tainted: read(media_db._.mediaid) {1}
[0] (FieldObject FieldType) mediaID int64
     --> r-tainted: read(posts_cache.value, media_db._.mediaid) {2}
[_1] (FieldObject FieldType) PostID int64
      --> r-tainted: read(posts_cache.value, media_db._.mediaid) {2}
[__2] (BasicObject BasicType) PostID int64
     --> r-tainted: read(media_db._.mediaid) {1}
[_1] (BasicObject BasicType) int64
      --> r-tainted: read(posts_cache.value, media_db._.mediaid) {2}
[__2] (BasicObject BasicType) PostID int64

[0] (StructObject UserType) media postnotification.Media struct{MediaID int64, Content string}
     --> r-tainted: read(media_db.Media) {1}
[_1] (Reference UserType) ref <media postnotification.Media struct{MediaID int64, Content string}> @ MediaService

[0] (InterfaceObject UserType) _ .error
[_1] (Reference UserType) ref <err .error> @ MediaService

