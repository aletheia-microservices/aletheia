[0] (PointerObject PointerType) f (*shopping_app.FrontendImpl struct{cart_service shopping_app.CartService, order_service shopping_app.OrderService})
[_1] (StructObject UserType) shopping_app.FrontendImpl struct{cart_service shopping_app.CartService, order_service shopping_app.OrderService}
[__2] (FieldObject FieldType) cart_service shopping_app.CartService
[___3] (ServiceObject ServiceType) cart_service shopping_app.CartService
[__2] (FieldObject FieldType) order_service shopping_app.OrderService
[___3] (ServiceObject ServiceType) order_service shopping_app.OrderService

[0] (InterfaceObject UserType) ctx context.Context

    --> r-tainted: read(cart_db._.cartID) {1}
[0] (BasicObject BasicType) cartID string

    --> r-tainted: read(product_db._.productID) {1}
[0] (BasicObject BasicType) productID string

