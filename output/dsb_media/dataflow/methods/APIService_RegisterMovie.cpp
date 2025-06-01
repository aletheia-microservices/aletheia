[0] (PointerObject PointerType) api (*mediamicroservices.APIServiceImpl struct{movieIdService mediamicroservices.MovieIdService, movieInfoService mediamicroservices.MovieInfoService})
[_1] (StructObject UserType) mediamicroservices.APIServiceImpl struct{movieIdService mediamicroservices.MovieIdService, movieInfoService mediamicroservices.MovieInfoService}
[__2] (FieldObject FieldType) movieIdService mediamicroservices.MovieIdService
[___3] (ServiceObject ServiceType) movieIdService mediamicroservices.MovieIdService
[__2] (FieldObject FieldType) movieInfoService mediamicroservices.MovieInfoService
[___3] (ServiceObject ServiceType) movieInfoService mediamicroservices.MovieInfoService

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) reqID int64

    --> w-tainted: write(movieid_db.MovieId.MovieID, movieinfo_db.MovieInfo.MovieID) {2}
[0] (BasicObject BasicType) movieID string

    --> w-tainted: write(movieid_db.MovieId.Title, movieinfo_db.MovieInfo.Title) {2}
[0] (BasicObject BasicType) title string

    --> w-tainted: write(movieinfo_db.MovieInfo.Casts) {1}
[0] (BasicObject BasicType) casts string

    --> w-tainted: write(movieid_db.MovieId) {1}
[0] (StructObject UserType) movieId mediamicroservices.MovieId struct{MovieID string, Title string}
     --> w-tainted: write(movieid_db.MovieId) {1}
[_1] (Reference UserType) ref <movieId mediamicroservices.MovieId struct{MovieID string, Title string}> @ MovieIdService
      --> w-tainted: write(movieid_db.MovieId.MovieID) {1}
[__2] (FieldObject FieldType) MovieID string
       --> w-tainted: write(movieid_db.MovieId.MovieID) {1}
[___3] (BasicObject BasicType) movieID string
        --> w-tainted: write(movieid_db.MovieId.MovieID, movieinfo_db.MovieInfo.MovieID) {2}
[____4] (Reference BasicType) ref <movieID string> @ APIService
      --> w-tainted: write(movieid_db.MovieId.Title) {1}
[__2] (FieldObject FieldType) Title string
       --> w-tainted: write(movieid_db.MovieId.Title) {1}
[___3] (BasicObject BasicType) title string
        --> w-tainted: write(movieid_db.MovieId.Title, movieinfo_db.MovieInfo.Title) {2}
[____4] (Reference BasicType) ref <title string> @ APIService

[0] (InterfaceObject UserType) err .error
[_1] (Reference UserType) ref <.error> @ MovieIdService

    --> w-tainted: write(movieinfo_db.MovieInfo) {1}
[0] (StructObject UserType) movieInfo mediamicroservices.MovieInfo struct{MovieID string, Title string, Casts string}
     --> w-tainted: write(movieinfo_db.MovieInfo) {1}
[_1] (Reference UserType) ref <movieInfo mediamicroservices.MovieInfo struct{MovieID string, Title string, Casts string}> @ MovieInfoService
      --> w-tainted: write(movieinfo_db.MovieInfo.Casts) {1}
[__2] (FieldObject FieldType) Casts string
       --> w-tainted: write(movieinfo_db.MovieInfo.Casts) {1}
[___3] (BasicObject BasicType) casts string
        --> w-tainted: write(movieinfo_db.MovieInfo.Casts) {1}
[____4] (Reference BasicType) ref <casts string> @ APIService
      --> w-tainted: write(movieinfo_db.MovieInfo.MovieID) {1}
[__2] (FieldObject FieldType) MovieID string
       --> w-tainted: write(movieinfo_db.MovieInfo.MovieID) {1}
[___3] (BasicObject BasicType) movieID string
        --> w-tainted: write(movieid_db.MovieId.MovieID, movieinfo_db.MovieInfo.MovieID) {2}
[____4] (Reference BasicType) ref <movieID string> @ APIService
      --> w-tainted: write(movieinfo_db.MovieInfo.Title) {1}
[__2] (FieldObject FieldType) Title string
       --> w-tainted: write(movieinfo_db.MovieInfo.Title) {1}
[___3] (BasicObject BasicType) title string
        --> w-tainted: write(movieid_db.MovieId.Title, movieinfo_db.MovieInfo.Title) {2}
[____4] (Reference BasicType) ref <title string> @ APIService

[0] (InterfaceObject UserType) err .error
[_1] (Reference UserType) ref <.error> @ MovieInfoService

