[0] (PointerObject PointerType) m (*mediamicroservices.MovieIdServiceImpl struct{movieIdDB NoSQLDatabase})
[_1] (StructObject UserType) mediamicroservices.MovieIdServiceImpl struct{movieIdDB NoSQLDatabase}
[__2] (FieldObject FieldType) movieIdDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) movieIdDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ APIService

[0] (BasicObject BasicType) reqID int64
[_1] (Reference BasicType) ref <reqID int64> @ APIService

    --> r-tainted: read(movieid_db.MovieId.MovieID) {1}
[0] (BasicObject BasicType) movieID string
     --> r-tainted: read(movieid_db.MovieId.MovieID, movieinfo_db.MovieInfo.MovieID) {2}
[_1] (Reference BasicType) ref <movieId string> @ APIService

    --> r-tainted: read(movieid_db.MovieId) {1}
[0] (StructObject UserType) movieId mediamicroservices.MovieId struct{MovieID string, Title string}

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = movie-id, collection = movie-id}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "movieid" string, Key "movieid" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "movieid" string
[___3] (BasicObject BasicType) "movieid" string
[__2] (FieldObject FieldType) Value string
       --> r-tainted: read(movieid_db.MovieId.MovieID) {1}
[___3] (BasicObject BasicType) movieID string
        --> r-tainted: read(movieid_db.MovieId.MovieID, movieinfo_db.MovieInfo.MovieID) {2}
[____4] (Reference BasicType) ref <movieId string> @ APIService

    --> r-tainted: read(movieid_db.MovieId) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = movie-id, collection = movie-id}
     --> r-tainted: read(movieid_db.MovieId) {1}
[_1] (StructObject UserType) movieId mediamicroservices.MovieId struct{MovieID string, Title string}

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) res bool

[0] (InterfaceObject UserType) err .error

