[0] (PointerObject PointerType) api (*mediamicroservices_sql.APIServiceImpl struct{movieIdService mediamicroservices_sql.MovieIdService, movieInfoService mediamicroservices_sql.MovieInfoService})
[_1] (StructObject UserType) mediamicroservices_sql.APIServiceImpl struct{movieIdService mediamicroservices_sql.MovieIdService, movieInfoService mediamicroservices_sql.MovieInfoService}
[__2] (FieldObject FieldType) movieIdService mediamicroservices_sql.MovieIdService
[___3] (ServiceObject ServiceType) movieIdService mediamicroservices_sql.MovieIdService
[__2] (FieldObject FieldType) movieInfoService mediamicroservices_sql.MovieInfoService
[___3] (ServiceObject ServiceType) movieInfoService mediamicroservices_sql.MovieInfoService

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) reqID int64

    --> w-tainted: write(movieid_db.movieid.movieid, movieinfo_db.movieinfo.movieid) {2}
[0] (BasicObject BasicType) movieID string

    --> w-tainted: write(movieid_db.movieid.title, movieinfo_db.movieinfo.title) {2}
[0] (BasicObject BasicType) title string

[0] (BasicObject BasicType) casts string

[0] (StructObject UserType) movieId mediamicroservices_sql.MovieId struct{MovieID string, Title string}
[_1] (Reference UserType) ref <movieId mediamicroservices_sql.MovieId struct{MovieID string, Title string}> @ MovieIdService
[__2] (FieldObject FieldType) MovieID string
       --> w-tainted: write(movieid_db.movieid.movieid) {1}
[___3] (BasicObject BasicType) movieID string
        --> w-tainted: write(movieid_db.movieid.movieid, movieinfo_db.movieinfo.movieid) {2}
[____4] (Reference BasicType) ref <movieID string> @ APIService
[__2] (FieldObject FieldType) Title string
       --> w-tainted: write(movieid_db.movieid.title) {1}
[___3] (BasicObject BasicType) title string
        --> w-tainted: write(movieid_db.movieid.title, movieinfo_db.movieinfo.title) {2}
[____4] (Reference BasicType) ref <title string> @ APIService

[0] (InterfaceObject UserType) err .error
[_1] (Reference UserType) ref <err .error> @ MovieIdService

[0] (StructObject UserType) movieInfo mediamicroservices_sql.MovieInfo struct{MovieID string, Title string, Casts string}
[_1] (Reference UserType) ref <movieInfo mediamicroservices_sql.MovieInfo struct{MovieID string, Title string, Casts string}> @ MovieInfoService
[__2] (FieldObject FieldType) Casts string
[___3] (BasicObject BasicType) casts string
[____4] (Reference BasicType) ref <casts string> @ APIService
[__2] (FieldObject FieldType) MovieID string
       --> w-tainted: write(movieinfo_db.movieinfo.movieid) {1}
[___3] (BasicObject BasicType) movieID string
        --> w-tainted: write(movieid_db.movieid.movieid, movieinfo_db.movieinfo.movieid) {2}
[____4] (Reference BasicType) ref <movieID string> @ APIService
[__2] (FieldObject FieldType) Title string
       --> w-tainted: write(movieinfo_db.movieinfo.title) {1}
[___3] (BasicObject BasicType) title string
        --> w-tainted: write(movieid_db.movieid.title, movieinfo_db.movieinfo.title) {2}
[____4] (Reference BasicType) ref <title string> @ APIService

[0] (InterfaceObject UserType) err .error
[_1] (Reference UserType) ref <err .error> @ MovieInfoService

