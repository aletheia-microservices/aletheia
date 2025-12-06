[0] (PointerObject PointerType) m (*mediamicroservices_sql.MovieInfoServiceImpl struct{movieInfoDB RelationalDB})
[_1] (StructObject UserType) mediamicroservices_sql.MovieInfoServiceImpl struct{movieInfoDB RelationalDB}
[__2] (FieldObject FieldType) movieInfoDB RelationalDB
[___3] (BlueprintBackendObject BlueprintBackendType) movieIdDB RelationalDB

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ APIService

[0] (BasicObject BasicType) reqID int64
[_1] (Reference BasicType) ref <reqID int64> @ APIService

    --> r-tainted: read(movieinfo_db.movieinfo.movieid) {1}
[0] (BasicObject BasicType) movieID string
     --> r-tainted: read(movieid_db.movieid.movieid, movieinfo_db.movieinfo.movieid) {2}
[_1] (Reference BasicType) ref <movieId string> @ APIService

    --> r-tainted: read(movieinfo_db.movieinfo.*) {1}
[0] (StructObject UserType) movieInfo mediamicroservices_sql.MovieInfo struct{MovieID string, Title string, Casts string}

[0] (InterfaceObject UserType) err .error

