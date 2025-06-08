[0] (PointerObject PointerType) m (*mediamicroservices.MovieInfoServiceImpl struct{movieInfoDB NoSQLDatabase})
[_1] (StructObject UserType) mediamicroservices.MovieInfoServiceImpl struct{movieInfoDB NoSQLDatabase}
[__2] (FieldObject FieldType) movieInfoDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) movieIdDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ APIService

[0] (BasicObject BasicType) reqID int64
[_1] (Reference BasicType) ref <reqID int64> @ APIService

    --> w-tainted: write(movieinfo_db.MovieInfo.MovieID) {1}
[0] (BasicObject BasicType) movieID string
     --> w-tainted: write(movieid_db.MovieId.MovieID, movieinfo_db.MovieInfo.MovieID) {2}
[_1] (Reference BasicType) ref <movieID string> @ APIService

    --> w-tainted: write(movieinfo_db.MovieInfo.Title) {1}
[0] (BasicObject BasicType) title string
     --> w-tainted: write(movieid_db.MovieId.Title, movieinfo_db.MovieInfo.Title) {2}
[_1] (Reference BasicType) ref <title string> @ APIService

    --> w-tainted: write(movieinfo_db.MovieInfo.Casts) {1}
[0] (BasicObject BasicType) casts string
     --> w-tainted: write(movieinfo_db.MovieInfo.Casts) {1}
[_1] (Reference BasicType) ref <casts string> @ APIService

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = movie-info, collection = movie-info}

[0] (InterfaceObject UserType) err .error

    --> w-tainted: write(movieinfo_db.MovieInfo) {1}
[0] (StructObject UserType) movieInfo mediamicroservices.MovieInfo struct{MovieID string, Title string, Casts string}
     --> w-tainted: write(movieinfo_db.MovieInfo.Casts) {1}
[_1] (FieldObject FieldType) Casts string
      --> w-tainted: write(movieinfo_db.MovieInfo.Casts) {1}
[__2] (BasicObject BasicType) casts string
       --> w-tainted: write(movieinfo_db.MovieInfo.Casts) {1}
[___3] (Reference BasicType) ref <casts string> @ APIService
     --> w-tainted: write(movieinfo_db.MovieInfo.MovieID) {1}
[_1] (FieldObject FieldType) MovieID string
      --> w-tainted: write(movieinfo_db.MovieInfo.MovieID) {1}
[__2] (BasicObject BasicType) movieID string
       --> w-tainted: write(movieid_db.MovieId.MovieID, movieinfo_db.MovieInfo.MovieID) {2}
[___3] (Reference BasicType) ref <movieID string> @ APIService
     --> w-tainted: write(movieinfo_db.MovieInfo.Title) {1}
[_1] (FieldObject FieldType) Title string
      --> w-tainted: write(movieinfo_db.MovieInfo.Title) {1}
[__2] (BasicObject BasicType) title string
       --> w-tainted: write(movieid_db.MovieId.Title, movieinfo_db.MovieInfo.Title) {2}
[___3] (Reference BasicType) ref <title string> @ APIService

