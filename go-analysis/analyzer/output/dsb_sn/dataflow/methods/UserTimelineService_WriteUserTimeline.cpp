[0] (PointerObject PointerType) u (*socialnetwork.UserTimelineServiceImpl struct{userTimelineCache Cache, userTimelineDB NoSQLDatabase, postStorageService socialnetwork.PostStorageService, CacheHits int64, CacheMiss int64, NumRequests int64})
[_1] (StructObject UserType) socialnetwork.UserTimelineServiceImpl struct{userTimelineCache Cache, userTimelineDB NoSQLDatabase, postStorageService socialnetwork.PostStorageService, CacheHits int64, CacheMiss int64, NumRequests int64}
[__2] (FieldObject FieldType) CacheHits int64
[___3] (BasicObject BasicType) 1 int
[__2] (FieldObject FieldType) CacheMiss int64
[___3] (BasicObject BasicType) 1 int
[__2] (FieldObject FieldType) NumRequests int64
[___3] (BasicObject BasicType) 1 int
[__2] (FieldObject FieldType) postStorageService socialnetwork.PostStorageService
[___3] (ServiceObject ServiceType) postStorageService socialnetwork.PostStorageService
[__2] (FieldObject FieldType) userTimelineCache Cache
[___3] (BlueprintBackendObject BlueprintBackendType) userTimelineCache Cache
[__2] (FieldObject FieldType) userTimelineDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) userTimelineDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) reqID int64

[0] (BasicObject BasicType) postID int64

[0] (BasicObject BasicType) userID int64

[0] (BasicObject BasicType) timestamp int64

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = usertimeline, collection = usertimeline}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ int64,  "userid" string,  int64}

[0] (BlueprintBackendObject BlueprintBackendType) results NoSQLCursor {database = usertimeline, collection = usertimeline}
[_1] (ArrayObject ArrayType) userPosts []socialnetwork.UserPosts struct{UserID int64, Posts []socialnetwork.PostInfo struct{PostID int64, Timestamp int64}}

[0] (InterfaceObject UserType) err .error

[0] (ArrayObject ArrayType) userPosts []socialnetwork.UserPosts struct{UserID int64, Posts []socialnetwork.PostInfo struct{PostID int64, Timestamp int64}}

[0] (StructObject UserType) userPosts socialnetwork.UserPosts struct{UserID int64, Posts []socialnetwork.PostInfo struct{PostID int64, Timestamp int64}}
[_1] (FieldObject FieldType) Posts []socialnetwork.PostInfo struct{PostID int64, Timestamp int64}
[__2] (ArrayObject ArrayType) []socialnetwork.PostInfo struct{PostID int64, Timestamp int64}
[___3] (StructObject UserType) socialnetwork.PostInfo struct{PostID int64, Timestamp int64}
[____4] (FieldObject FieldType) PostID int64
[_____5] (BasicObject BasicType) postID int64
[____4] (FieldObject FieldType) Timestamp int64
[_____5] (BasicObject BasicType) timestamp int64
[_1] (FieldObject FieldType) UserID int64
[__2] (BasicObject BasicType) userID int64

[0] (InterfaceObject UserType) err .error

[0] (ArrayObject ArrayType) postInfo []socialnetwork.PostInfo struct{PostID int64, Timestamp int64}
[_1] (StructObject UserType) socialnetwork.PostInfo struct{PostID int64, Timestamp int64}
[__2] (FieldObject FieldType) PostID int64
[___3] (BasicObject BasicType) postID int64
[__2] (FieldObject FieldType) Timestamp int64
[___3] (BasicObject BasicType) timestamp int64

[0] (BasicObject BasicType) userIDStr string
[_1] (BasicObject BasicType) userID int64
[_1] (BasicObject BasicType) 10 int
[_1] (BasicObject BasicType) userID int64
[_1] (BasicObject BasicType) 10 int

[0] (BasicObject BasicType) _ bool

[0] (InterfaceObject UserType) err .error

[0] (ArrayObject ArrayType) postInfo []socialnetwork.PostInfo struct{PostID int64, Timestamp int64}
[_1] (StructObject UserType) socialnetwork.PostInfo struct{PostID int64, Timestamp int64}
[__2] (FieldObject FieldType) PostID int64
[___3] (BasicObject BasicType) postID int64
[__2] (FieldObject FieldType) Timestamp int64
[___3] (BasicObject BasicType) timestamp int64

