[0] (PointerObject PointerType) api (*mediamicroservices.APIServiceImpl struct{movieIdService mediamicroservices.MovieIdService, movieInfoService mediamicroservices.MovieInfoService})
[_1] (StructObject UserType) mediamicroservices.APIServiceImpl struct{movieIdService mediamicroservices.MovieIdService, movieInfoService mediamicroservices.MovieInfoService}
[__2] (FieldObject FieldType) movieIdService mediamicroservices.MovieIdService
[___3] (ServiceObject ServiceType) movieIdService mediamicroservices.MovieIdService
[__2] (FieldObject FieldType) movieInfoService mediamicroservices.MovieInfoService
[___3] (ServiceObject ServiceType) movieInfoService mediamicroservices.MovieInfoService

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) reqID int64

    --> r-tainted: read(movieid_db.MovieId.MovieID, movieinfo_db.MovieInfo.MovieID) {2}
[0] (BasicObject BasicType) movieId string

[0] (StructObject UserType) movie1 mediamicroservices.MovieId struct{MovieID string, Title string}
     --> r-tainted: read(movieid_db.MovieId) {1}
[_1] (Reference UserType) ref <movieId mediamicroservices.MovieId struct{MovieID string, Title string}> @ MovieIdService

[0] (InterfaceObject UserType) err1 .error
[_1] (Reference UserType) ref <err .error> @ MovieIdService

[0] (StructObject UserType) movie2 mediamicroservices.MovieInfo struct{MovieID string, Title string, Casts string}
     --> r-tainted: read(movieinfo_db.MovieInfo) {1}
[_1] (Reference UserType) ref <movieInfo mediamicroservices.MovieInfo struct{MovieID string, Title string, Casts string}> @ MovieInfoService

[0] (InterfaceObject UserType) err2 .error
[_1] (Reference UserType) ref <err .error> @ MovieInfoService

