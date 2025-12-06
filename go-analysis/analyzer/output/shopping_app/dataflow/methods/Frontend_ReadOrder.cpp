[0] (PointerObject PointerType) f (*shopping_app.FrontendImpl struct{cart_service shopping_app.CartService, order_service shopping_app.OrderService})
[_1] (StructObject UserType) shopping_app.FrontendImpl struct{cart_service shopping_app.CartService, order_service shopping_app.OrderService}
[__2] (FieldObject FieldType) cart_service shopping_app.CartService
[___3] (ServiceObject ServiceType) cart_service shopping_app.CartService
[__2] (FieldObject FieldType) order_service shopping_app.OrderService
[___3] (ServiceObject ServiceType) order_service shopping_app.OrderService

[0] (InterfaceObject UserType) ctx context.Context

    --> r-tainted: read(order_db._.orderID, order_db.Order.OrderID) {2}
[0] (BasicObject BasicType) orderID string

