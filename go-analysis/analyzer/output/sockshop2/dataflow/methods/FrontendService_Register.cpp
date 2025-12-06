[0] (PointerObject PointerType) f (*frontend.FrontendImpl struct{user user.UserService, catalogue catalogue.CatalogueService, cart cart.CartService, order order.OrderService})
[_1] (StructObject UserType) frontend.FrontendImpl struct{user user.UserService, catalogue catalogue.CatalogueService, cart cart.CartService, order order.OrderService}
[__2] (FieldObject FieldType) cart cart.CartService
[___3] (ServiceObject ServiceType) cart cart.CartService
[__2] (FieldObject FieldType) catalogue catalogue.CatalogueService
[___3] (ServiceObject ServiceType) catalogue catalogue.CatalogueService
[__2] (FieldObject FieldType) order order.OrderService
[___3] (ServiceObject ServiceType) order order.OrderService
[__2] (FieldObject FieldType) user user.UserService
[___3] (ServiceObject ServiceType) user user.UserService

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) sessionID string

    --> w-tainted: write(user_db.User.Username) {1}
[0] (BasicObject BasicType) username string

    --> w-tainted: write(user_db.User.Password) {1}
[0] (BasicObject BasicType) password string

    --> w-tainted: write(user_db.User.Email) {1}
[0] (BasicObject BasicType) email string

    --> w-tainted: write(user_db.User.FirstName) {1}
[0] (BasicObject BasicType) first string

    --> w-tainted: write(user_db.User.LastName) {1}
[0] (BasicObject BasicType) last string

[0] (BasicObject BasicType) userID string
[_1] (Reference BasicType) ref <UserID string> @ UserService

[0] (InterfaceObject UserType) err .error
[_1] (Reference BasicType) ref <nil> @ UserService

