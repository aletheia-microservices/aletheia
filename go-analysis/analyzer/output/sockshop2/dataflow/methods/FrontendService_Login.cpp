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

    --> r-tainted: read(user_db._.username) {1}
[0] (BasicObject BasicType) username string

[0] (BasicObject BasicType) password string

[0] (StructObject UserType) u user.User struct{FirstName string, LastName string, Email string, Username string, Password string, Addresses user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}, Cards user.Card struct{LongNum string, Expires string, CCV string, ID string}, UserID string, Salt string}
     --> r-tainted: read(user_db.User) {1}
[_1] (Reference UserType) ref <user user.User struct{FirstName string, LastName string, Email string, Username string, Password string, Addresses user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}, Cards user.Card struct{LongNum string, Expires string, CCV string, ID string}, UserID string, Salt string}> @ UserService
[_1] (FieldObject FieldType) UserID string
[__2] (BasicObject BasicType) UserID string

[0] (InterfaceObject UserType) err .error
[_1] (Reference BasicType) ref <nil> @ UserService

[0] (InterfaceObject UserType) err .error
[_1] (Reference UserType) ref <.error> @ CartService

