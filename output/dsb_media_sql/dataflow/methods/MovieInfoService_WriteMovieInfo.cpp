[0] (PointerObject PointerType) m (*mediamicroservices_sql.MovieInfoServiceImpl struct{movieInfoDB RelationalDB})
[_1] (StructObject UserType) mediamicroservices_sql.MovieInfoServiceImpl struct{movieInfoDB RelationalDB}
[__2] (FieldObject FieldType) movieInfoDB RelationalDB
[___3] (BlueprintBackendObject BlueprintBackendType) movieIdDB RelationalDB

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ APIService

[0] (BasicObject BasicType) reqID int64
[_1] (Reference BasicType) ref <reqID int64> @ APIService

    --> w-tainted: write(movieinfo_db.movieinfo.movieid) {1}
[0] (BasicObject BasicType) movieID string
     --> w-tainted: write(movieid_db.movieid.movieid, movieinfo_db.movieinfo.movieid) {2}
[_1] (Reference BasicType) ref <movieID string> @ APIService

    --> w-tainted: write(movieinfo_db.movieinfo.title) {1}
[0] (BasicObject BasicType) title string
     --> w-tainted: write(movieid_db.movieid.title, movieinfo_db.movieinfo.title) {2}
[_1] (Reference BasicType) ref <title string> @ APIService

[0] (BasicObject BasicType) casts string
[_1] (Reference BasicType) ref <casts string> @ APIService

[0] (StructObject UserType) movieInfo mediamicroservices_sql.MovieInfo struct{MovieID string, Title string, Casts string}
[_1] (FieldObject FieldType) Casts string
[__2] (BasicObject BasicType) casts string
[___3] (Reference BasicType) ref <casts string> @ APIService
[_1] (FieldObject FieldType) MovieID string
      --> w-tainted: write(movieinfo_db.movieinfo.movieid) {1}
[__2] (BasicObject BasicType) movieID string
       --> w-tainted: write(movieid_db.movieid.movieid, movieinfo_db.movieinfo.movieid) {2}
[___3] (Reference BasicType) ref <movieID string> @ APIService
[_1] (FieldObject FieldType) Title string
      --> w-tainted: write(movieinfo_db.movieinfo.title) {1}
[__2] (BasicObject BasicType) title string
       --> w-tainted: write(movieid_db.movieid.title, movieinfo_db.movieinfo.title) {2}
[___3] (Reference BasicType) ref <title string> @ APIService

[0] (InterfaceObject UserType) _ sql.Result interface{ interface{LastInsertId() (int64, error); RowsAffected() (int64, error)} }

[0] (InterfaceObject UserType) err .error

