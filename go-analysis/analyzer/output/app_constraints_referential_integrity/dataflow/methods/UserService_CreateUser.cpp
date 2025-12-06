[0] (PointerObject PointerType) s (*app_constraints_referential_integrity.UserServiceImpl struct{usersDB NoSQLDatabase})
[_1] (StructObject UserType) app_constraints_referential_integrity.UserServiceImpl struct{usersDB NoSQLDatabase}
[__2] (FieldObject FieldType) usersDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) usersDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ AccountService
[__2] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> w-tainted: write(users_db.User.ReqID, accounts_db.Account.ReqID) {2}
[0] (BasicObject BasicType) reqID string
     --> w-tainted: write(accounts_db.Account.ReqID, users_db.User.ReqID) {2}
[_1] (Reference BasicType) ref <reqID string> @ AccountService
      --> w-tainted: write(users_db.User.ReqID, accounts_db.Account.ReqID) {2}
[__2] (Reference BasicType) ref <reqID string> @ Frontend

    --> w-tainted: write(users_db.User.Username, accounts_db.Account.Username) {2}
[0] (BasicObject BasicType) username string
     --> w-tainted: write(accounts_db.Account.Username, users_db.User.Username) {2}
[_1] (Reference BasicType) ref <username string> @ AccountService
      --> w-tainted: write(users_db.User.Username, accounts_db.Account.Username) {2}
[__2] (Reference BasicType) ref <username string> @ Frontend

    --> w-tainted: write(users_db.User.Timestamp) {1}
[0] (BasicObject BasicType) timestamp int64

    --> w-tainted: write(users_db.User) {1}
[0] (StructObject UserType) user app_constraints_referential_integrity.User struct{ReqID string, UserID string, Username string, Timestamp int64}
     --> w-tainted: write(users_db.User.ReqID) {1}
[_1] (FieldObject FieldType) ReqID string
      --> w-tainted: write(users_db.User.ReqID, accounts_db.Account.ReqID) {2}
[__2] (BasicObject BasicType) reqID string
       --> w-tainted: write(accounts_db.Account.ReqID, users_db.User.ReqID) {2}
[___3] (Reference BasicType) ref <reqID string> @ AccountService
        --> w-tainted: write(users_db.User.ReqID, accounts_db.Account.ReqID) {2}
[____4] (Reference BasicType) ref <reqID string> @ Frontend
     --> w-tainted: write(users_db.User.Timestamp) {1}
[_1] (FieldObject FieldType) Timestamp int64
      --> w-tainted: write(users_db.User.Timestamp) {1}
[__2] (BasicObject BasicType) timestamp int64
     --> w-tainted: write(users_db.User.Username) {1}
[_1] (FieldObject FieldType) Username string
      --> w-tainted: write(users_db.User.Username, accounts_db.Account.Username) {2}
[__2] (BasicObject BasicType) username string
       --> w-tainted: write(accounts_db.Account.Username, users_db.User.Username) {2}
[___3] (Reference BasicType) ref <username string> @ AccountService
        --> w-tainted: write(users_db.User.Username, accounts_db.Account.Username) {2}
[____4] (Reference BasicType) ref <username string> @ Frontend

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = users, collection = users}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

