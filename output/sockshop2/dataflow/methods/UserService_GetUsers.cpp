[0] (PointerObject PointerType) s (*user.userServiceImpl struct{UserService user.UserService, db NoSQLDatabase})
[_1] (StructObject UserType) user.userServiceImpl struct{UserService user.UserService, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ FrontendService
[_1] (Reference UserType) ref <ctx context.Context> @ OrderService
[__2] (Reference UserType) ref <ctx context.Context> @ FrontendService

    --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3}       --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3} --> r-tainted: read(user_db._.userid, shipqueue.Shipment.Name) {2}
[0] (BasicObject BasicType) userid string
     --> r-tainted: read(user_db._.userid) {1}
[_1] (Reference BasicType) ref <userID string> @ FrontendService
     --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4}         --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4} --> r-tainted: read(shipqueue.Shipment.Name, user_db._.userid) {2}
[_1] (Reference BasicType) ref <customerID string> @ OrderService
      --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3}           --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3} --> r-tainted: read(user_db._.userid, shipqueue.Shipment.Name) {2}
[__2] (Reference BasicType) ref <userID string> @ FrontendService

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = users, collection = users}

[0] (InterfaceObject UserType) _ .error

    --> r-tainted: read(user_db.[]user.User struct{FirstName string, LastName string, Email string, Username string, Password string, Addresses user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}, Cards user.Card struct{LongNum string, Expires string, CCV string, ID string}, UserID string, Salt string}) {1}
[0] (ArrayObject ArrayType) users []user.User struct{FirstName string, LastName string, Email string, Username string, Password string, Addresses user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}, Cards user.Card struct{LongNum string, Expires string, CCV string, ID string}, UserID string, Salt string}

[0] (SliceObject UserType) filter primitive.D
[_1] (StructObject StructType) struct{Key "userid" string, Key "userid" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "userid" string
[___3] (BasicObject BasicType) "userid" string
[__2] (FieldObject FieldType) Value string
       --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3}             --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3} --> r-tainted: read(user_db._.userid, shipqueue.Shipment.Name) {2}
[___3] (BasicObject BasicType) userid string
        --> r-tainted: read(user_db._.userid) {1}
[____4] (Reference BasicType) ref <userID string> @ FrontendService
        --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4}               --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.CustomerID, order_db.Order.Shipment) {4} --> r-tainted: read(shipqueue.Shipment.Name, user_db._.userid) {2}
[____4] (Reference BasicType) ref <customerID string> @ OrderService
         --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3}                 --> w-tainted: write(shipqueue.Shipment.Name, shipdb.Shipment.Name, order_db.Order.Shipment) {3} --> r-tainted: read(user_db._.userid, shipqueue.Shipment.Name) {2}
[_____5] (Reference BasicType) ref <userID string> @ FrontendService

    --> r-tainted: read(user_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = users, collection = users}
     --> r-tainted: read(user_db.[]user.User struct{FirstName string, LastName string, Email string, Username string, Password string, Addresses user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}, Cards user.Card struct{LongNum string, Expires string, CCV string, ID string}, UserID string, Salt string}) {1}
[_1] (ArrayObject ArrayType) users []user.User struct{FirstName string, LastName string, Email string, Username string, Password string, Addresses user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}, Cards user.Card struct{LongNum string, Expires string, CCV string, ID string}, UserID string, Salt string}

[0] (InterfaceObject UserType) _ .error

