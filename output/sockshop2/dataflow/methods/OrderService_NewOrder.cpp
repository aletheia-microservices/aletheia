[0] (PointerObject PointerType) s (*order.orderImpl struct{users user.UserService, carts cart.CartService, payments payment.PaymentService, shipping shipping.ShippingService, db NoSQLDatabase})
[_1] (StructObject UserType) order.orderImpl struct{users user.UserService, carts cart.CartService, payments payment.PaymentService, shipping shipping.ShippingService, db NoSQLDatabase}
[__2] (FieldObject FieldType) carts cart.CartService
[___3] (ServiceObject ServiceType) cartService cart.CartService
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) orderDB NoSQLDatabase
[__2] (FieldObject FieldType) payments payment.PaymentService
[___3] (ServiceObject ServiceType) payments payment.PaymentService
[__2] (FieldObject FieldType) shipping shipping.ShippingService
[___3] (ServiceObject ServiceType) shipping shipping.ShippingService
[__2] (FieldObject FieldType) users user.UserService
[___3] (ServiceObject ServiceType) userService user.UserService

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ FrontendService

    --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4}       --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4} --> r-tainted: read(shipqueue.Shipment.Name, user_db._.userid) {2}
[0] (BasicObject BasicType) customerID string
     --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID) {3}         --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID) {3} --> r-tainted: read(user_db._.userid, shipqueue.Shipment.Name) {2}
[_1] (Reference BasicType) ref <userID string> @ FrontendService

    --> r-tainted: read(user_db._.addresses) {1}
[0] (BasicObject BasicType) addressID string
     --> r-tainted: read(user_db._.addresses) {1}
[_1] (Reference BasicType) ref <addressID string> @ FrontendService

    --> r-tainted: read(user_db._.cards) {1}
[0] (BasicObject BasicType) cardID string
     --> r-tainted: read(user_db._.cards) {1}
[_1] (Reference BasicType) ref <cardID string> @ FrontendService

[0] (BasicObject BasicType) cartID string
[_1] (Reference BasicType) ref <cartID string> @ FrontendService

[0] (StructObject UserType) wg sync.WaitGroup

    --> w-tainted: write(order_db.Order.Items) {1}
[0] (ArrayObject ArrayType) items []cart.Item struct{ID string, Quantity int, UnitPrice float32}

[0] (ArrayObject ArrayType) addresses []user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}

[0] (ArrayObject ArrayType) cards []user.Card struct{LongNum string, Expires string, CCV string, ID string}

    --> w-tainted: write(order_db.Order.Total) {1}
[0] (BasicObject BasicType) amount 10 float32

    --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3}       --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment) {1}
[0] (StructObject UserType) shipment shipping.Shipment struct{ID string, Name string, Status "awaiting shipment" string}
     --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID, order_db.Order.Shipment.ID) {4}         --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID, order_db.Order.Shipment.ID) {4} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[_1] (FieldObject FieldType) ID string
      --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID) {3}           --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID) {3} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[__2] (BasicObject BasicType) string
     --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3}         --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Name) {1}
[_1] (FieldObject FieldType) Name string
      --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4}           --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4} --> r-tainted: read(shipqueue.Shipment.Name, user_db._.userid) {2}
[__2] (BasicObject BasicType) customerID string
       --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID) {3}             --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID) {3} --> r-tainted: read(user_db._.userid, shipqueue.Shipment.Name) {2}
[___3] (Reference BasicType) ref <userID string> @ FrontendService
     --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3}         --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Status) {1}
[_1] (FieldObject FieldType) Status "awaiting shipment" string
      --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3}           --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Status) {1}
[__2] (BasicObject BasicType) "awaiting shipment" string

    --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3}       --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment) {1}
[0] (StructObject UserType) shipment shipping.Shipment struct{ID string, Name string, Status "awaiting shipment" string}
     --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3}         --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment) {1}
[_1] (Reference UserType) ref <shipment shipping.Shipment struct{ID string, Name string, Status string}> @ ShippingService
      --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3}           --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment) {1}
[__2] (Reference UserType) ref <shipment shipping.Shipment struct{ID string, Name string, Status "awaiting shipment" string}> @ OrderService
       --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID, order_db.Order.Shipment.ID) {4}             --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID, order_db.Order.Shipment.ID) {4} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[___3] (FieldObject FieldType) ID string
        --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID) {3}               --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID) {3} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[____4] (BasicObject BasicType) string
       --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3}             --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Name) {1}
[___3] (FieldObject FieldType) Name string
        --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4}               --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4} --> r-tainted: read(shipqueue.Shipment.Name, user_db._.userid) {2}
[____4] (BasicObject BasicType) customerID string
         --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID) {3}                 --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID) {3} --> r-tainted: read(user_db._.userid, shipqueue.Shipment.Name) {2}
[_____5] (Reference BasicType) ref <userID string> @ FrontendService
       --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3}             --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Status) {1}
[___3] (FieldObject FieldType) Status "awaiting shipment" string
        --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3}               --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Status) {1}
[____4] (BasicObject BasicType) "awaiting shipment" string
     --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3}         --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment) {1}
[_1] (Reference UserType) ref <shipment shipping.Shipment struct{ID string, Name string, Status "awaiting shipment" string}> @ OrderService
      --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID, order_db.Order.Shipment.ID) {4}           --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID, order_db.Order.Shipment.ID) {4} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[__2] (FieldObject FieldType) ID string
       --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID) {3}             --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID) {3} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[___3] (BasicObject BasicType) string
      --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3}           --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Name) {1}
[__2] (FieldObject FieldType) Name string
       --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4}             --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4} --> r-tainted: read(shipqueue.Shipment.Name, user_db._.userid) {2}
[___3] (BasicObject BasicType) customerID string
        --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID) {3}               --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID) {3} --> r-tainted: read(user_db._.userid, shipqueue.Shipment.Name) {2}
[____4] (Reference BasicType) ref <userID string> @ FrontendService
      --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3}           --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Status) {1}
[__2] (FieldObject FieldType) Status "awaiting shipment" string
       --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3}             --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Status) {1}
[___3] (BasicObject BasicType) "awaiting shipment" string
     --> w-tainted: write(order_db.Order.Shipment.ID, shipqueue.Shipment.ID, shipdb.Shipment.ID) {3}
[_1] (FieldObject FieldType) ID string
      --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID, order_db.Order.Shipment.ID) {4}           --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID, order_db.Order.Shipment.ID) {4} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[__2] (Reference FieldType) ref <ID string> @ OrderService
       --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID) {3}             --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID) {3} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[___3] (BasicObject BasicType) string
      --> w-tainted: write(order_db.Order.ID, order_db.Order.Shipment.ID, shipqueue.Shipment.ID, shipdb.Shipment.ID) {4}
[__2] (BasicObject BasicType) ID string
       --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID, order_db.Order.Shipment.ID) {4}             --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID, order_db.Order.Shipment.ID) {4} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[___3] (Reference FieldType) ref <ID string> @ OrderService
        --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID) {3}               --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID) {3} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[____4] (BasicObject BasicType) string

[0] (InterfaceObject UserType) err .error
[_1] (Reference UserType) ref <.error> @ ShippingService

    --> w-tainted: write(order_db.Order) {1}
[0] (StructObject UserType) order order.Order struct{ID string, CustomerID string, Customer user.User struct{FirstName string, LastName string, Email string, Username string, Password string, Addresses user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}, Cards user.Card struct{LongNum string, Expires string, CCV string, ID string}, UserID string, Salt string}, Address user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}, Card user.Card struct{LongNum string, Expires string, CCV string, ID string}, Items []cart.Item struct{ID string, Quantity int, UnitPrice float32}, Shipment shipping.Shipment struct{ID string, Name string, Status "awaiting shipment" string}, Date string, Total 10 float32}
     --> w-tainted: write(order_db.Order.Address) {1}
[_1] (FieldObject FieldType) Address user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}
      --> w-tainted: write(order_db.Order.Address) {1}
[__2] (StructObject UserType) * user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}
     --> w-tainted: write(order_db.Order.Card) {1}
[_1] (FieldObject FieldType) Card user.Card struct{LongNum string, Expires string, CCV string, ID string}
      --> w-tainted: write(order_db.Order.Card) {1}
[__2] (StructObject UserType) * user.Card struct{LongNum string, Expires string, CCV string, ID string}
     --> w-tainted: write(order_db.Order.CustomerID) {1}
[_1] (FieldObject FieldType) CustomerID string
      --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4}           --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4} --> r-tainted: read(shipqueue.Shipment.Name, user_db._.userid) {2}
[__2] (BasicObject BasicType) customerID string
       --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID) {3}             --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID) {3} --> r-tainted: read(user_db._.userid, shipqueue.Shipment.Name) {2}
[___3] (Reference BasicType) ref <userID string> @ FrontendService
     --> w-tainted: write(order_db.Order.Date) {1}
[_1] (FieldObject FieldType) Date string
      --> w-tainted: write(order_db.Order.Date) {1}
[__2] (BasicObject BasicType) string
     --> w-tainted: write(order_db.Order.ID) {1}
[_1] (FieldObject FieldType) ID string
      --> w-tainted: write(order_db.Order.ID, order_db.Order.Shipment.ID, shipqueue.Shipment.ID, shipdb.Shipment.ID) {4}
[__2] (BasicObject BasicType) ID string
       --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID, order_db.Order.Shipment.ID) {4}             --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID, order_db.Order.Shipment.ID) {4} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[___3] (Reference FieldType) ref <ID string> @ OrderService
        --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID) {3}               --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID) {3} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[____4] (BasicObject BasicType) string
     --> w-tainted: write(order_db.Order.Items) {1}
[_1] (FieldObject FieldType) Items []cart.Item struct{ID string, Quantity int, UnitPrice float32}
      --> w-tainted: write(order_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items []cart.Item struct{ID string, Quantity int, UnitPrice float32}
     --> w-tainted: write(order_db.Order.Shipment) {1}
[_1] (FieldObject FieldType) Shipment shipping.Shipment struct{ID string, Name string, Status "awaiting shipment" string}
      --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3}           --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment) {1}
[__2] (StructObject UserType) shipment shipping.Shipment struct{ID string, Name string, Status "awaiting shipment" string}
       --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3}             --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment) {1}
[___3] (Reference UserType) ref <shipment shipping.Shipment struct{ID string, Name string, Status string}> @ ShippingService
        --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3}               --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment) {1}
[____4] (Reference UserType) ref <shipment shipping.Shipment struct{ID string, Name string, Status "awaiting shipment" string}> @ OrderService
         --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID, order_db.Order.Shipment.ID) {4}                 --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID, order_db.Order.Shipment.ID) {4} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[_____5] (FieldObject FieldType) ID string
          --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID) {3}                   --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID) {3} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[______6] (BasicObject BasicType) string
         --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3}                 --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Name) {1}
[_____5] (FieldObject FieldType) Name string
          --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4}                   --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4} --> r-tainted: read(shipqueue.Shipment.Name, user_db._.userid) {2}
[______6] (BasicObject BasicType) customerID string
           --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID) {3}                     --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID) {3} --> r-tainted: read(user_db._.userid, shipqueue.Shipment.Name) {2}
[_______7] (Reference BasicType) ref <userID string> @ FrontendService
         --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3}                 --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Status) {1}
[_____5] (FieldObject FieldType) Status "awaiting shipment" string
          --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3}                   --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Status) {1}
[______6] (BasicObject BasicType) "awaiting shipment" string
       --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3}             --> w-tainted: write(shipqueue.Shipment, shipdb.Shipment, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment) {1}
[___3] (Reference UserType) ref <shipment shipping.Shipment struct{ID string, Name string, Status "awaiting shipment" string}> @ OrderService
        --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID, order_db.Order.Shipment.ID) {4}               --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID, order_db.Order.Shipment.ID) {4} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[____4] (FieldObject FieldType) ID string
         --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID) {3}                 --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID) {3} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[_____5] (BasicObject BasicType) string
        --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3}               --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Name) {1}
[____4] (FieldObject FieldType) Name string
         --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4}                 --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4} --> r-tainted: read(shipqueue.Shipment.Name, user_db._.userid) {2}
[_____5] (BasicObject BasicType) customerID string
          --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID) {3}                   --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID) {3} --> r-tainted: read(user_db._.userid, shipqueue.Shipment.Name) {2}
[______6] (Reference BasicType) ref <userID string> @ FrontendService
        --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3}               --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Status) {1}
[____4] (FieldObject FieldType) Status "awaiting shipment" string
         --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3}                 --> w-tainted: write(shipqueue.Shipment.Status, shipdb.Shipment.Status, order_db.Order.Shipment) {3} --> r-tainted: read(shipqueue.Shipment.Status) {1}
[_____5] (BasicObject BasicType) "awaiting shipment" string
       --> w-tainted: write(order_db.Order.Shipment.ID, shipqueue.Shipment.ID, shipdb.Shipment.ID) {3}
[___3] (FieldObject FieldType) ID string
        --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID, order_db.Order.Shipment.ID) {4}               --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID, order_db.Order.Shipment.ID) {4} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[____4] (Reference FieldType) ref <ID string> @ OrderService
         --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID) {3}                 --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID) {3} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[_____5] (BasicObject BasicType) string
        --> w-tainted: write(order_db.Order.ID, order_db.Order.Shipment.ID, shipqueue.Shipment.ID, shipdb.Shipment.ID) {4}
[____4] (BasicObject BasicType) ID string
         --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID, order_db.Order.Shipment.ID) {4}                 --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID, order_db.Order.Shipment.ID) {4} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[_____5] (Reference FieldType) ref <ID string> @ OrderService
          --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID) {3}                   --> w-tainted: write(shipqueue.Shipment.ID, shipdb.Shipment.ID, order_db.Order.ID) {3} --> r-tainted: read(shipqueue.Shipment.ID) {1}
[______6] (BasicObject BasicType) string
     --> w-tainted: write(order_db.Order.Total) {1}
[_1] (FieldObject FieldType) Total 10 float32
      --> w-tainted: write(order_db.Order.Total) {1}
[__2] (BasicObject BasicType) amount 10 float32

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = order_service, collection = orders}

[0] (InterfaceObject UserType) _ .error

[0] (InterfaceObject UserType) err .error

