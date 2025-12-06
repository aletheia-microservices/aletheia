[0] (PointerObject PointerType) u (*app_constraints_referential_integrity.FrontendImpl struct{accountService app_constraints_referential_integrity.AccountService, userService app_constraints_referential_integrity.UserService})
[_1] (StructObject UserType) app_constraints_referential_integrity.FrontendImpl struct{accountService app_constraints_referential_integrity.AccountService, userService app_constraints_referential_integrity.UserService}
[__2] (FieldObject FieldType) accountService app_constraints_referential_integrity.AccountService
[___3] (ServiceObject ServiceType) accountService app_constraints_referential_integrity.AccountService
[__2] (FieldObject FieldType) userService app_constraints_referential_integrity.UserService
[___3] (ServiceObject ServiceType) userService app_constraints_referential_integrity.UserService

[0] (InterfaceObject UserType) ctx context.Context

    --> w-tainted: write(users_db.User.ReqID, accounts_db.Account.ReqID) {2}
[0] (BasicObject BasicType) reqID string

    --> w-tainted: write(accounts_db.Account.AccountID) {1}
[0] (BasicObject BasicType) accountID string

    --> w-tainted: write(users_db.User.Username, accounts_db.Account.Username) {2}
[0] (BasicObject BasicType) username string

    --> w-tainted: write(accounts_db.Account) {1}
[0] (StructObject UserType) account app_constraints_referential_integrity.Account struct{ReqID string, AccountID string, Username string, Timestamp int64}
     --> w-tainted: write(accounts_db.Account) {1}
[_1] (Reference UserType) ref <account app_constraints_referential_integrity.Account struct{ReqID string, AccountID string, Username string, Timestamp int64}> @ AccountService
      --> w-tainted: write(accounts_db.Account.AccountID) {1}
[__2] (FieldObject FieldType) AccountID string
       --> w-tainted: write(accounts_db.Account.AccountID) {1}
[___3] (BasicObject BasicType) accountID string
        --> w-tainted: write(accounts_db.Account.AccountID) {1}
[____4] (Reference BasicType) ref <accountID string> @ Frontend
      --> w-tainted: write(accounts_db.Account.ReqID) {1}
[__2] (FieldObject FieldType) ReqID string
       --> w-tainted: write(accounts_db.Account.ReqID, users_db.User.ReqID) {2}
[___3] (BasicObject BasicType) reqID string
        --> w-tainted: write(users_db.User.ReqID, accounts_db.Account.ReqID) {2}
[____4] (Reference BasicType) ref <reqID string> @ Frontend
      --> w-tainted: write(accounts_db.Account.Timestamp) {1}
[__2] (FieldObject FieldType) Timestamp int64
       --> w-tainted: write(accounts_db.Account.Timestamp) {1}
[___3] (BasicObject BasicType) timestamp int64
      --> w-tainted: write(accounts_db.Account.Username) {1}
[__2] (FieldObject FieldType) Username string
       --> w-tainted: write(accounts_db.Account.Username, users_db.User.Username) {2}
[___3] (BasicObject BasicType) username string
        --> w-tainted: write(users_db.User.Username, accounts_db.Account.Username) {2}
[____4] (Reference BasicType) ref <username string> @ Frontend

[0] (InterfaceObject UserType) err .error
[_1] (Reference UserType) ref <err .error> @ AccountService

