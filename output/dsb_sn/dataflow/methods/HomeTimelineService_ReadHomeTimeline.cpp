[0] (PointerObject PointerType) h (*socialnetwork.HomeTimelineServiceImpl struct{homeTimelineCache Cache, postStorageService socialnetwork.PostStorageService, socialGraphService socialnetwork.SocialGraphService})
[_1] (StructObject UserType) socialnetwork.HomeTimelineServiceImpl struct{homeTimelineCache Cache, postStorageService socialnetwork.PostStorageService, socialGraphService socialnetwork.SocialGraphService}
[__2] (FieldObject FieldType) homeTimelineCache Cache
[___3] (BlueprintBackendObject BlueprintBackendType) homeTimelineCache Cache
[__2] (FieldObject FieldType) postStorageService socialnetwork.PostStorageService
[___3] (ServiceObject ServiceType) postStorageService socialnetwork.PostStorageService
[__2] (FieldObject FieldType) socialGraphService socialnetwork.SocialGraphService
[___3] (ServiceObject ServiceType) socialGraphService socialnetwork.SocialGraphService

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Wrk2APIService

[0] (BasicObject BasicType) reqID int64
[_1] (Reference BasicType) ref <reqID int64> @ Wrk2APIService

    --> r-tainted: read(hometimeline_cache.key) {1}
[0] (BasicObject BasicType) userID int64
     --> r-tainted: read(hometimeline_cache.key) {1}
[_1] (Reference BasicType) ref <userId int64> @ Wrk2APIService

[0] (BasicObject BasicType) start int64
[_1] (Reference BasicType) ref <start int64> @ Wrk2APIService

[0] (BasicObject BasicType) stop int64
[_1] (Reference BasicType) ref <stop int64> @ Wrk2APIService

    --> r-tainted: read(hometimeline_cache.key) {1}
[0] (BasicObject BasicType) userIDStr string
     --> r-tainted: read(hometimeline_cache.key) {1}
[_1] (BasicObject BasicType) userID int64
      --> r-tainted: read(hometimeline_cache.key) {1}
[__2] (Reference BasicType) ref <userId int64> @ Wrk2APIService
     --> r-tainted: read(hometimeline_cache.key) {1}
[_1] (BasicObject BasicType) 10 int
     --> r-tainted: read(hometimeline_cache.key) {1}
[_1] (BasicObject BasicType) userID int64
      --> r-tainted: read(hometimeline_cache.key) {1}
[__2] (Reference BasicType) ref <userId int64> @ Wrk2APIService
     --> r-tainted: read(hometimeline_cache.key) {1}
[_1] (BasicObject BasicType) 10 int

    --> r-tainted: read(hometimeline_cache.value) {1}
[0] (ArrayObject ArrayType) postIDs []int64

[0] (BasicObject BasicType) _ bool

[0] (InterfaceObject UserType) err .error

