[0] (PointerObject PointerType) api (*mediamicroservices_sql.APIServiceImpl struct{movieIdService mediamicroservices_sql.MovieIdService, movieInfoService mediamicroservices_sql.MovieInfoService})
[_1] (StructObject UserType) mediamicroservices_sql.APIServiceImpl struct{movieIdService mediamicroservices_sql.MovieIdService, movieInfoService mediamicroservices_sql.MovieInfoService}
[__2] (FieldObject FieldType) movieIdService mediamicroservices_sql.MovieIdService
[___3] (ServiceObject ServiceType) movieIdService mediamicroservices_sql.MovieIdService
[__2] (FieldObject FieldType) movieInfoService mediamicroservices_sql.MovieInfoService
[___3] (ServiceObject ServiceType) movieInfoService mediamicroservices_sql.MovieInfoService

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) reqID int64

    --> w-tainted: write(movieid_db.movieid.movieid) {1}
[0] (BasicObject BasicType) movieID string

    --> w-tainted: write(movieid_db.movieid.title) {1}
[0] (BasicObject BasicType) title string

