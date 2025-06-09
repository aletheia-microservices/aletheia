[0] (PointerObject PointerType) s (*cart.cartImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) cart.cartImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ FrontendService

[0] (BasicObject BasicType) customerID string
[_1] (Reference BasicType) ref <sessionID string> @ FrontendService

[0] (StructObject UserType) item cart.Item struct{ID string, Quantity int, UnitPrice float32}
[_1] (Reference UserType) ref <cart.Item struct{ID string, Quantity 1 int, UnitPrice float32}> @ FrontendService
[__2] (FieldObject FieldType) ID string
[___3] (BasicObject BasicType) ID string
[__2] (FieldObject FieldType) Quantity int
[___3] (BasicObject BasicType) quantity int
[__2] (FieldObject FieldType) UnitPrice float32
[___3] (BasicObject BasicType) Price float32

[0] (PointerObject PointerType) cart (*cart.cart struct{ID string, Items []cart.Item struct{ID string, Quantity int, UnitPrice float32}})
[_1] (StructObject UserType) cart.cart struct{ID string, Items []cart.Item struct{ID string, Quantity int, UnitPrice float32}}
[__2] (FieldObject FieldType) ID string
[___3] (BasicObject BasicType) customerID string
[____4] (Reference BasicType) ref <sessionID string> @ FrontendService
[__2] (FieldObject FieldType) Items []cart.Item struct{ID string, Quantity int, UnitPrice float32}
[___3] (ArrayObject ArrayType) []cart.Item struct{ID string, Quantity int, UnitPrice float32}
[____4] (StructObject UserType) item cart.Item struct{ID string, Quantity int, UnitPrice float32}
[_____5] (Reference UserType) ref <cart.Item struct{ID string, Quantity 1 int, UnitPrice float32}> @ FrontendService
[______6] (FieldObject FieldType) ID string
[_______7] (BasicObject BasicType) ID string
[______6] (FieldObject FieldType) Quantity int
[_______7] (BasicObject BasicType) quantity int
[______6] (FieldObject FieldType) UnitPrice float32
[_______7] (BasicObject BasicType) Price float32

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = cart, collection = carts}

[0] (InterfaceObject UserType) _ .error

