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

[0] (BasicObject BasicType) reqID int64

    --> r-tainted: read(posts_cache.key) {1}
[0] (BasicObject BasicType) postID int64

[0] (StructObject UserType) media postnotification.Media struct{MediaID int64, Content string}

[0] (StructObject UserType) media postnotification.Media struct{MediaID int64, Content string}
[_1] (Reference UserType) ref <media postnotification.Media struct{MediaID int64, Content string}> @ StorageService
      --> r-tainted: read(media_db.Media) {1}
[__2] (Reference UserType) ref <media postnotification.Media struct{MediaID int64, Content string}> @ MediaService

[0] (InterfaceObject UserType) _ .error
[_1] (Reference BasicType) ref <nil> @ StorageService

