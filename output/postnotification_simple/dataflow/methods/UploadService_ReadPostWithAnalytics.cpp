[0] (PointerObject PointerType) u (*postnotification_simple.UploadServiceImpl struct{storageService postnotification_simple.StorageService, analyticsService postnotification_simple.AnalyticsService, notificationsQueue Queue})
[_1] (StructObject UserType) postnotification_simple.UploadServiceImpl struct{storageService postnotification_simple.StorageService, analyticsService postnotification_simple.AnalyticsService, notificationsQueue Queue}
[__2] (FieldObject FieldType) analyticsService postnotification_simple.AnalyticsService
[___3] (ServiceObject ServiceType) analyticsService postnotification_simple.AnalyticsService
[__2] (FieldObject FieldType) notificationsQueue Queue
[___3] (BlueprintBackendObject BlueprintBackendType) notificationsQueue Queue
[__2] (FieldObject FieldType) storageService postnotification_simple.StorageService
[___3] (ServiceObject ServiceType) storageService postnotification_simple.StorageService

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) reqID int64

    --> r-tainted: read(posts_db._.postid, analytics_db._.postid, posts_db.Post.PostID) {3}
[0] (BasicObject BasicType) postID int64

[0] (StructObject UserType) post postnotification_simple.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification_simple.Creator struct{Username string}}
     --> r-tainted: read(posts_db.Post) {1}
[_1] (Reference UserType) ref <post postnotification_simple.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification_simple.Creator struct{Username string}}> @ StorageService

[0] (InterfaceObject UserType) err .error
[_1] (Reference UserType) ref <err .error> @ StorageService

[0] (StructObject UserType) analytics postnotification_simple.Analytics struct{PostID int64}
     --> r-tainted: read(analytics_db.Analytics) {1}
[_1] (Reference UserType) ref <analytics postnotification_simple.Analytics struct{PostID int64}> @ AnalyticsService

[0] (InterfaceObject UserType) err .error
[_1] (Reference BasicType) ref <nil> @ AnalyticsService

