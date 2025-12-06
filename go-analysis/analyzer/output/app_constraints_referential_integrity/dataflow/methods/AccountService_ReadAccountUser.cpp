[0] (PointerObject PointerType) s (*app_constraints_referential_integrity.AccountServiceImpl struct{userService app_constraints_referential_integrity.UserService, accountsDB NoSQLDatabase})
[_1] (StructObject UserType) app_constraints_referential_integrity.AccountServiceImpl struct{userService app_constraints_referential_integrity.UserService, accountsDB NoSQLDatabase}
[__2] (FieldObject FieldType) accountsDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) accountsDB NoSQLDatabase
[__2] (FieldObject FieldType) userService app_constraints_referential_integrity.UserService
[___3] (ServiceObject ServiceType) userService app_constraints_referential_integrity.UserService

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> r-tainted: read(accounts_db.Account.AccountID) {1}
[0] (BasicObject BasicType) accountID string
     --> r-tainted: read(accounts_db.Account.AccountID) {1}
[_1] (Reference BasicType) ref <accountID string> @ Frontend

    --> r-tainted: read(accounts_db.Account, accounts_db.Account.Username) {2}
[0] (StructObject UserType) account app_constraints_referential_integrity.Account struct{ReqID string, AccountID string, Username string, Timestamp int64}
     --> r-tainted: read(accounts_db.Account.Username) {1}
[_1] (FieldObject FieldType) Username string
      --> r-tainted: read(accounts_db.Account.Username, users_db.User.Username) {2}
[__2] (BasicObject BasicType) Username string

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = accounts, collection = accounts}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "accountID" string, Key "accountID" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "accountID" string
[___3] (BasicObject BasicType) "accountID" string
[__2] (FieldObject FieldType) Value string
       --> r-tainted: read(accounts_db.Account.AccountID) {1}
[___3] (BasicObject BasicType) accountID string
        --> r-tainted: read(accounts_db.Account.AccountID) {1}
[____4] (Reference BasicType) ref <accountID string> @ Frontend

    --> r-tainted: read(accounts_db.Account) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = accounts, collection = accounts}
     --> r-tainted: read(accounts_db.Account, accounts_db.Account.Username) {2}
[_1] (StructObject UserType) account app_constraints_referential_integrity.Account struct{ReqID string, AccountID string, Username string, Timestamp int64}
      --> r-tainted: read(accounts_db.Account.Username) {1}
[__2] (FieldObject FieldType) Username string
       --> r-tainted: read(accounts_db.Account.Username, users_db.User.Username) {2}
[___3] (BasicObject BasicType) Username string

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) res bool

[0] (InterfaceObject UserType) err .error

