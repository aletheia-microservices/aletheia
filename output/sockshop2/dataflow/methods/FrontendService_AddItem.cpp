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

    --> r-tainted: read(catalogue_db._.id) {1}
[0] (BasicObject BasicType) itemID string

[0] (BasicObject BasicType) sessionID string

[0] (StructObject UserType) sock catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}
     --> r-tainted: read(catalogue_db.Sock, catalogue_db.Sock.ImageURL, catalogue_db.Sock.string, catalogue_db.Sock.ImageURL_1, catalogue_db.Sock.ImageURL_2, catalogue_db.Sock.TagString, catalogue_db.Sock.Tags) {7}
[_1] (Reference UserType) ref <sock catalogue.Sock struct{ID string, Name string, Description string, ImageURL []string, ImageURL_1 string, ImageURL_2 string, Price float32, Quantity int, Tags []string, TagString string}> @ CatalogueService
      --> r-tainted: read(catalogue_db.Sock.ImageURL) {1}
[__2] (FieldObject FieldType) ImageURL []string
       --> r-tainted: read(catalogue_db.Sock.ImageURL) {1}
[___3] (ArrayObject SliceType) []string
        --> r-tainted: read(catalogue_db.Sock.string, catalogue_db.Sock.ImageURL_1) {2}
[____4] (BasicObject BasicType) ImageURL_1 string
        --> r-tainted: read(catalogue_db.Sock.string, catalogue_db.Sock.ImageURL_2) {2}
[____4] (BasicObject BasicType) ImageURL_2 string
      --> r-tainted: read(catalogue_db.Sock.ImageURL_1) {1}
[__2] (FieldObject FieldType) ImageURL_1 string
       --> r-tainted: read(catalogue_db.Sock.string, catalogue_db.Sock.ImageURL_1) {2}
[___3] (BasicObject BasicType) ImageURL_1 string
      --> r-tainted: read(catalogue_db.Sock.ImageURL_2) {1}
[__2] (FieldObject FieldType) ImageURL_2 string
       --> r-tainted: read(catalogue_db.Sock.string, catalogue_db.Sock.ImageURL_2) {2}
[___3] (BasicObject BasicType) ImageURL_2 string
      --> r-tainted: read(catalogue_db.Sock.TagString) {1}
[__2] (FieldObject FieldType) TagString string
       --> r-tainted: read(catalogue_db.Sock.TagString) {1}
[___3] (BasicObject BasicType) TagString string
      --> r-tainted: read(catalogue_db.Sock.Tags) {1}
[__2] (FieldObject FieldType) Tags []string
       --> r-tainted: read(catalogue_db.Sock.Tags) {1}
[___3] (SliceObject SliceType) []string
        --> r-tainted: read(catalogue_db.Sock.TagString) {1}
[____4] (BasicObject BasicType) TagString string
        --> r-tainted: read(catalogue_db.Sock.string) {1}
[____4] (BasicObject BasicType) "," string
        --> r-tainted: read(catalogue_db.Sock.TagString) {1}
[____4] (BasicObject BasicType) TagString string
        --> r-tainted: read(catalogue_db.Sock.string) {1}
[____4] (BasicObject BasicType) "," string
[_1] (FieldObject FieldType) ID string
[__2] (BasicObject BasicType) ID string
[_1] (FieldObject FieldType) Price float32
[__2] (BasicObject BasicType) Price float32

[0] (InterfaceObject UserType) err .error
[_1] (Reference BasicType) ref <nil> @ CatalogueService

[0] (StructObject UserType) _ cart.Item struct{ID string, Quantity 1 int, UnitPrice float32}
[_1] (Reference UserType) ref <item cart.Item struct{ID string, Quantity int, UnitPrice float32}> @ CartService
[__2] (Reference UserType) ref <cart.Item struct{ID string, Quantity 1 int, UnitPrice float32}> @ FrontendService
[___3] (FieldObject FieldType) ID string
[____4] (BasicObject BasicType) ID string
[___3] (FieldObject FieldType) Quantity 1 int
[____4] (BasicObject BasicType) 1 int
[___3] (FieldObject FieldType) UnitPrice float32
[____4] (BasicObject BasicType) Price float32
[_1] (Reference UserType) ref <cart.Item struct{ID string, Quantity 1 int, UnitPrice float32}> @ FrontendService
[__2] (FieldObject FieldType) ID string
[___3] (BasicObject BasicType) ID string
[__2] (FieldObject FieldType) Quantity 1 int
[___3] (BasicObject BasicType) 1 int
[__2] (FieldObject FieldType) UnitPrice float32
[___3] (BasicObject BasicType) Price float32

[0] (InterfaceObject UserType) err .error
[_1] (Reference BasicType) ref <nil> @ CartService

