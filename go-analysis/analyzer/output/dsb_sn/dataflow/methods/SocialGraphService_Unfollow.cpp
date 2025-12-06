[0] (PointerObject PointerType) s (*socialnetwork.SocialGraphServiceImpl struct{socialGraphCache Cache, socialGraphDB NoSQLDatabase, userIDService socialnetwork.UserIDService})
[_1] (StructObject UserType) socialnetwork.SocialGraphServiceImpl struct{socialGraphCache Cache, socialGraphDB NoSQLDatabase, userIDService socialnetwork.UserIDService}
[__2] (FieldObject FieldType) socialGraphCache Cache
[___3] (BlueprintBackendObject BlueprintBackendType) socialGraphCache Cache
[__2] (FieldObject FieldType) socialGraphDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) socialGraphDB NoSQLDatabase
[__2] (FieldObject FieldType) userIDService socialnetwork.UserIDService
[___3] (ServiceObject ServiceType) userIDService socialnetwork.UserIDService

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) reqID int64

[0] (BasicObject BasicType) userID int64

[0] (BasicObject BasicType) followeeID int64

[0] (BasicObject BasicType) userIDstr string
[_1] (BasicObject BasicType) userID int64
[_1] (BasicObject BasicType) 10 int
[_1] (BasicObject BasicType) userID int64
[_1] (BasicObject BasicType) 10 int

[0] (BasicObject BasicType) followeeIDstr string
[_1] (BasicObject BasicType) followeeID int64
[_1] (BasicObject BasicType) 10 int
[_1] (BasicObject BasicType) followeeID int64
[_1] (BasicObject BasicType) 10 int

[0] (BasicObject BasicType) err1 error

[0] (BasicObject BasicType) err2 error

[0] (BasicObject BasicType) err3 error

[0] (StructObject UserType) wg sync.WaitGroup

[0] (ArrayObject ArrayType) followees []socialnetwork.FolloweeInfo struct{FolloweeID int64, Timestamp int64}

[0] (ArrayObject ArrayType) followers []socialnetwork.FollowerInfo struct{FollowerID int64, Timestamp int64}

