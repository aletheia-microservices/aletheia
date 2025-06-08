[0] (PointerObject PointerType) api (*mediamicroservices_sql.APIServiceImpl struct{movieIdService mediamicroservices_sql.MovieIdService, movieInfoService mediamicroservices_sql.MovieInfoService})
[_1] (StructObject UserType) mediamicroservices_sql.APIServiceImpl struct{movieIdService mediamicroservices_sql.MovieIdService, movieInfoService mediamicroservices_sql.MovieInfoService}
[__2] (FieldObject FieldType) movieIdService mediamicroservices_sql.MovieIdService
[___3] (ServiceObject ServiceType) movieIdService mediamicroservices_sql.MovieIdService
[__2] (FieldObject FieldType) movieInfoService mediamicroservices_sql.MovieInfoService
[___3] (ServiceObject ServiceType) movieInfoService mediamicroservices_sql.MovieInfoService

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) reqID int64

    --> r-tainted: read(movieid_db.movieid.movieid, movieinfo_db.movieinfo.movieid) {2}
[0] (BasicObject BasicType) movieId string

[0] (StructObject UserType) movie1 mediamicroservices_sql.MovieId struct{MovieID string, Title string}
     --> r-tainted: read(movieid_db.movieid.*) {1}
[_1] (Reference UserType) ref <movieId mediamicroservices_sql.MovieId struct{MovieID string, Title string}> @ MovieIdService

[0] (InterfaceObject UserType) err1 .error
[_1] (Reference UserType) ref <err .error> @ MovieIdService

[0] (StructObject UserType) movie2 mediamicroservices_sql.MovieInfo struct{MovieID string, Title string, Casts string}
     --> r-tainted: read(movieinfo_db.movieinfo.*) {1}
[_1] (Reference UserType) ref <movieInfo mediamicroservices_sql.MovieInfo struct{MovieID string, Title string, Casts string}> @ MovieInfoService

[0] (InterfaceObject UserType) err2 .error
[_1] (Reference UserType) ref <err .error> @ MovieInfoService

