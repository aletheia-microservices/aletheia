[0] (PointerObject PointerType) m (*mediamicroservices.MovieInfoServiceImpl struct{movieInfoDB NoSQLDatabase})
[_1] (StructObject UserType) mediamicroservices.MovieInfoServiceImpl struct{movieInfoDB NoSQLDatabase}
[__2] (FieldObject FieldType) movieInfoDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) movieIdDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ APIService

[0] (BasicObject BasicType) reqID int64
[_1] (Reference BasicType) ref <reqID int64> @ APIService

    --> r-tainted: read(movieinfo_db.MovieInfo.MovieID) {1}
[0] (BasicObject BasicType) movieID string
     --> r-tainted: read(movieid_db.MovieId.MovieID, movieinfo_db.MovieInfo.MovieID) {2}
[_1] (Reference BasicType) ref <movieId string> @ APIService

    --> r-tainted: read(movieinfo_db.MovieInfo) {1}
[0] (StructObject UserType) movieInfo mediamicroservices.MovieInfo struct{MovieID string, Title string, Casts string}

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = movie-info, collection = movie-info}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "movieid" string, Key "movieid" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "movieid" string
[___3] (BasicObject BasicType) "movieid" string
[__2] (FieldObject FieldType) Value string
       --> r-tainted: read(movieinfo_db.MovieInfo.MovieID) {1}
[___3] (BasicObject BasicType) movieID string
        --> r-tainted: read(movieid_db.MovieId.MovieID, movieinfo_db.MovieInfo.MovieID) {2}
[____4] (Reference BasicType) ref <movieId string> @ APIService

    --> r-tainted: read(movieinfo_db.MovieInfo) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = movie-info, collection = movie-info}
     --> r-tainted: read(movieinfo_db.MovieInfo) {1}
[_1] (StructObject UserType) movieInfo mediamicroservices.MovieInfo struct{MovieID string, Title string, Casts string}

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) res bool

[0] (InterfaceObject UserType) err .error

