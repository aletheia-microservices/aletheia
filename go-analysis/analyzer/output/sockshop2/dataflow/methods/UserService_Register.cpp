[0] (PointerObject PointerType) s (*user.userServiceImpl struct{UserService user.UserService, db NoSQLDatabase})
[_1] (StructObject UserType) user.userServiceImpl struct{UserService user.UserService, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ FrontendService

    --> w-tainted: write(user_db.User.Username) {1}
[0] (BasicObject BasicType) username string
     --> w-tainted: write(user_db.User.Username) {1}
[_1] (Reference BasicType) ref <username string> @ FrontendService

    --> w-tainted: write(user_db.User.Password) {1}
[0] (BasicObject BasicType) password string
     --> w-tainted: write(user_db.User.Password) {1}
[_1] (Reference BasicType) ref <password string> @ FrontendService

    --> w-tainted: write(user_db.User.Email) {1}
[0] (BasicObject BasicType) email string
     --> w-tainted: write(user_db.User.Email) {1}
[_1] (Reference BasicType) ref <email string> @ FrontendService

    --> w-tainted: write(user_db.User.FirstName) {1}
[0] (BasicObject BasicType) first string
     --> w-tainted: write(user_db.User.FirstName) {1}
[_1] (Reference BasicType) ref <first string> @ FrontendService

    --> w-tainted: write(user_db.User.LastName) {1}
[0] (BasicObject BasicType) last string
     --> w-tainted: write(user_db.User.LastName) {1}
[_1] (Reference BasicType) ref <last string> @ FrontendService

    --> w-tainted: write(user_db.User) {1}
[0] (StructObject UserType) u user.User struct{FirstName string, LastName string, Email string, Username string, Password string, Addresses user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}, Cards user.Card struct{LongNum string, Expires string, CCV string, ID string}, UserID string, Salt string}
     --> w-tainted: write(user_db.User.Addresses) {1}
[_1] (FieldObject FieldType) Addresses user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}
      --> w-tainted: write(user_db.User.Addresses) {1}
[__2] (StructObject UserType) user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}
     --> w-tainted: write(user_db.User.Cards) {1}
[_1] (FieldObject FieldType) Cards user.Card struct{LongNum string, Expires string, CCV string, ID string}
      --> w-tainted: write(user_db.User.Cards) {1}
[__2] (StructObject UserType) user.Card struct{LongNum string, Expires string, CCV string, ID string}
     --> w-tainted: write(user_db.User.Email) {1}
[_1] (FieldObject FieldType) Email string
      --> w-tainted: write(user_db.User.Email) {1}
[__2] (BasicObject BasicType) email string
       --> w-tainted: write(user_db.User.Email) {1}
[___3] (Reference BasicType) ref <email string> @ FrontendService
     --> w-tainted: write(user_db.User.FirstName) {1}
[_1] (FieldObject FieldType) FirstName string
      --> w-tainted: write(user_db.User.FirstName) {1}
[__2] (BasicObject BasicType) first string
       --> w-tainted: write(user_db.User.FirstName) {1}
[___3] (Reference BasicType) ref <first string> @ FrontendService
     --> w-tainted: write(user_db.User.LastName) {1}
[_1] (FieldObject FieldType) LastName string
      --> w-tainted: write(user_db.User.LastName) {1}
[__2] (BasicObject BasicType) last string
       --> w-tainted: write(user_db.User.LastName) {1}
[___3] (Reference BasicType) ref <last string> @ FrontendService
     --> w-tainted: write(user_db.User.Password) {1}
[_1] (FieldObject FieldType) Password string
      --> w-tainted: write(user_db.User.Password) {1}
[__2] (BasicObject BasicType) password string
       --> w-tainted: write(user_db.User.Password) {1}
[___3] (Reference BasicType) ref <password string> @ FrontendService
[_1] (FieldObject FieldType) UserID string
[__2] (BasicObject BasicType) UserID string
     --> w-tainted: write(user_db.User.Username) {1}
[_1] (FieldObject FieldType) Username string
      --> w-tainted: write(user_db.User.Username) {1}
[__2] (BasicObject BasicType) username string
       --> w-tainted: write(user_db.User.Username) {1}
[___3] (Reference BasicType) ref <username string> @ FrontendService

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = users, collection = users}

[0] (InterfaceObject UserType) _ .error

