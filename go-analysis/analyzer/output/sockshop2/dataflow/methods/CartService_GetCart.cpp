[0] (PointerObject PointerType) s (*cart.cartImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) cart.cartImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ OrderService
[__2] (Reference UserType) ref <ctx context.Context> @ FrontendService

[0] (BasicObject BasicType) customerID string
[_1] (Reference BasicType) ref <cartID string> @ OrderService
[__2] (Reference BasicType) ref <cartID string> @ FrontendService

[0] (PointerObject PointerType) cart (*cart.cart struct{ID string, Items []cart.Item struct{ID string, Quantity int, UnitPrice float32}})
[_1] (StructObject UserType) cart.cart struct{ID string, Items []cart.Item struct{ID string, Quantity int, UnitPrice float32}}
[__2] (FieldObject FieldType) ID string
[___3] (BasicObject BasicType) customerID string
[____4] (Reference BasicType) ref <cartID string> @ OrderService
[_____5] (Reference BasicType) ref <cartID string> @ FrontendService
[__2] (FieldObject FieldType) Items []cart.Item struct{ID string, Quantity int, UnitPrice float32}
[___3] (SliceObject SliceType) Items []cart.Item struct{ID string, Quantity int, UnitPrice float32}

