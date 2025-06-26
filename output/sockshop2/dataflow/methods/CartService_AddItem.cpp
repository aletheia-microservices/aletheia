[0] (PointerObject PointerType) s (*cart.cartImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) cart.cartImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ FrontendService

    --> w-tainted: write(cart_db.cart.ID) {1}
[0] (BasicObject BasicType) customerID string
     --> w-tainted: write(cart_db.cart.ID) {1}
[_1] (Reference BasicType) ref <sessionID string> @ FrontendService

    --> w-tainted: write(cart_db.cart.Items) {1}
[0] (StructObject UserType) item cart.Item struct{ID string, Quantity int, UnitPrice float32}
     --> w-tainted: write(cart_db.cart.Items) {1}
[_1] (Reference UserType) ref <cart.Item struct{ID string, Quantity 1 int, UnitPrice float32}> @ FrontendService
      --> w-tainted: write(cart_db.cart.Items) {1}
[__2] (FieldObject FieldType) ID string
       --> w-tainted: write(cart_db.cart.Items) {1}             --> w-tainted: write(cart_db.cart.Items) {1} --> r-tainted: read(catalogue_db.Sock.ID.id) {1}
[___3] (BasicObject BasicType) itemID string
      --> w-tainted: write(cart_db.cart.Items) {1}
[__2] (FieldObject FieldType) Quantity 1 int
       --> w-tainted: write(cart_db.cart.Items) {1}
[___3] (BasicObject BasicType) 1 int
      --> w-tainted: write(cart_db.cart.Items) {1}
[__2] (FieldObject FieldType) UnitPrice float32
       --> w-tainted: write(cart_db.cart.Items) {1}
[___3] (BasicObject BasicType) Price float32

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = cart, collection = carts}

[0] (InterfaceObject UserType) _ .error

    --> w-tainted: write(cart_db.cart) {1}
[0] (StructObject UserType) cart cart.cart struct{ID string, Items []cart.Item struct{ID string, Quantity int, UnitPrice float32}}
     --> w-tainted: write(cart_db.cart.ID) {1}
[_1] (FieldObject FieldType) ID string
      --> w-tainted: write(cart_db.cart.ID) {1}
[__2] (BasicObject BasicType) customerID string
       --> w-tainted: write(cart_db.cart.ID) {1}
[___3] (Reference BasicType) ref <sessionID string> @ FrontendService
     --> w-tainted: write(cart_db.cart.Items) {1}
[_1] (FieldObject FieldType) Items []cart.Item struct{ID string, Quantity int, UnitPrice float32}
      --> w-tainted: write(cart_db.cart.Items) {1}
[__2] (ArrayObject ArrayType) []cart.Item struct{ID string, Quantity int, UnitPrice float32}
       --> w-tainted: write(cart_db.cart.Items) {1}
[___3] (StructObject UserType) item cart.Item struct{ID string, Quantity int, UnitPrice float32}
        --> w-tainted: write(cart_db.cart.Items) {1}
[____4] (Reference UserType) ref <cart.Item struct{ID string, Quantity 1 int, UnitPrice float32}> @ FrontendService
         --> w-tainted: write(cart_db.cart.Items) {1}
[_____5] (FieldObject FieldType) ID string
          --> w-tainted: write(cart_db.cart.Items) {1}                   --> w-tainted: write(cart_db.cart.Items) {1} --> r-tainted: read(catalogue_db.Sock.ID.id) {1}
[______6] (BasicObject BasicType) itemID string
         --> w-tainted: write(cart_db.cart.Items) {1}
[_____5] (FieldObject FieldType) Quantity 1 int
          --> w-tainted: write(cart_db.cart.Items) {1}
[______6] (BasicObject BasicType) 1 int
         --> w-tainted: write(cart_db.cart.Items) {1}
[_____5] (FieldObject FieldType) UnitPrice float32
          --> w-tainted: write(cart_db.cart.Items) {1}
[______6] (BasicObject BasicType) Price float32

