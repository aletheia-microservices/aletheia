[0] (PointerObject PointerType) api (*mediamicroservices.APIServiceImpl struct{movieIdService mediamicroservices.MovieIdService, movieInfoService mediamicroservices.MovieInfoService})
[_1] (StructObject UserType) mediamicroservices.APIServiceImpl struct{movieIdService mediamicroservices.MovieIdService, movieInfoService mediamicroservices.MovieInfoService}
[__2] (FieldObject FieldType) movieIdService mediamicroservices.MovieIdService
[___3] (ServiceObject ServiceType) movieIdService mediamicroservices.MovieIdService
[__2] (FieldObject FieldType) movieInfoService mediamicroservices.MovieInfoService
[___3] (ServiceObject ServiceType) movieInfoService mediamicroservices.MovieInfoService

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) reqID int64

    --> w-tainted: write(movieinfo_db.MovieInfo.MovieID) {1}
[0] (BasicObject BasicType) movieID string

    --> w-tainted: write(movieinfo_db.MovieInfo.Title) {1}
[0] (BasicObject BasicType) title string

    --> w-tainted: write(movieinfo_db.MovieInfo.Casts) {1}
[0] (BasicObject BasicType) casts string

