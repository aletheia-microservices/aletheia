[0] (PointerObject PointerType) u (*app_constraints_referential_integrity.FrontendImpl struct{accountService app_constraints_referential_integrity.AccountService, userService app_constraints_referential_integrity.UserService})
[_1] (StructObject UserType) app_constraints_referential_integrity.FrontendImpl struct{accountService app_constraints_referential_integrity.AccountService, userService app_constraints_referential_integrity.UserService}
[__2] (FieldObject FieldType) accountService app_constraints_referential_integrity.AccountService
[___3] (ServiceObject ServiceType) accountService app_constraints_referential_integrity.AccountService
[__2] (FieldObject FieldType) userService app_constraints_referential_integrity.UserService
[___3] (ServiceObject ServiceType) userService app_constraints_referential_integrity.UserService

[0] (InterfaceObject UserType) ctx context.Context

    --> r-tainted: read(accounts_db.Account.AccountID) {1}
[0] (BasicObject BasicType) accountID string

