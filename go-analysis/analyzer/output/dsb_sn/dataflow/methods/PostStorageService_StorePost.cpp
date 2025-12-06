[0] (PointerObject PointerType) p (*socialnetwork.PostStorageServiceImpl struct{postStorageCache Cache, postStorageDB NoSQLDatabase, CacheHits int64, NumReqs int64, CacheMiss int64})
[_1] (StructObject UserType) socialnetwork.PostStorageServiceImpl struct{postStorageCache Cache, postStorageDB NoSQLDatabase, CacheHits int64, NumReqs int64, CacheMiss int64}
[__2] (FieldObject FieldType) CacheMiss int64
[___3] (BasicObject BasicType) 1 int
[__2] (FieldObject FieldType) NumReqs int64
[___3] (BasicObject BasicType) 1 int
[__2] (FieldObject FieldType) postStorageCache Cache
[___3] (BlueprintBackendObject BlueprintBackendType) postStorageCache Cache
[__2] (FieldObject FieldType) postStorageDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) postStorageDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) reqID int64

[0] (StructObject UserType) post socialnetwork.Post struct{PostID int64, Creator socialnetwork.Creator struct{UserID int64, Username string}, ReqID int64, Text string, UserMentions []socialnetwork.UserMention struct{UserID int64, Username string}, Medias []socialnetwork.Media struct{MediaID int64, MediaType string}, Urls []socialnetwork.URL struct{ShortenedUrl string, ExpandedUrl string}, Timestamp int64, PostType int64}

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = post, collection = post}

[0] (InterfaceObject UserType) err .error

