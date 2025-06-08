[0] (PointerObject PointerType) m (*mediamicroservices_sql.MovieIdServiceImpl struct{movieIdDB RelationalDB})
[_1] (StructObject UserType) mediamicroservices_sql.MovieIdServiceImpl struct{movieIdDB RelationalDB}
[__2] (FieldObject FieldType) movieIdDB RelationalDB
[___3] (BlueprintBackendObject BlueprintBackendType) movieIdDB RelationalDB

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ APIService

[0] (BasicObject BasicType) reqID int64
[_1] (Reference BasicType) ref <reqID int64> @ APIService

    --> r-tainted: read(movieid_db.movieid.movieid) {1}
[0] (BasicObject BasicType) movieID string
     --> r-tainted: read(movieid_db.movieid.movieid, movieinfo_db.movieinfo.movieid) {2}
[_1] (Reference BasicType) ref <movieId string> @ APIService

    --> r-tainted: read(movieid_db.movieid.*) {1}
[0] (StructObject UserType) movieId mediamicroservices_sql.MovieId struct{MovieID string, Title string}

[0] (InterfaceObject UserType) err .error

