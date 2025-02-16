[0] (PointerObject PointerType) f (*shopping_app.FrontendImpl struct{cart_service shopping_app.CartService, order_service shopping_app.OrderService})
[_1] (StructObject UserType) shopping_app.FrontendImpl struct{cart_service shopping_app.CartService, order_service shopping_app.OrderService}
[__2] (FieldObject FieldType) cart_service shopping_app.CartService
[___3] (ServiceObject ServiceType) cart_service shopping_app.CartService
[__2] (FieldObject FieldType) order_service shopping_app.OrderService
[___3] (ServiceObject ServiceType) order_service shopping_app.OrderService

[0] (InterfaceObject UserType) ctx context.Context

    --> w-tainted: write(order_db.Order.OrderID, shipment_queue.ShipmentMessage.OrderID) {2}       --> w-tainted: write(order_db.Order.OrderID, shipment_queue.ShipmentMessage.OrderID) {2} --> r-tainted: read(shipment_queue.ShipmentMessage.OrderID) {1}
[0] (BasicObject BasicType) cartID string

    --> w-tainted: write(order_db.Order.UserID, billing_db.Bill.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {4}       --> w-tainted: write(order_db.Order.UserID, billing_db.Bill.UserID, analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {4} --> r-tainted: read(analytics_queue.AnalyticsMessage.UserID, shipment_queue.ShipmentMessage.UserID) {2}
[0] (BasicObject BasicType) userID string

    --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3}       --> w-tainted: write(order_db.Order.ProductID, stock_db.Stock.ProductID, billing_db.Bill.ProductID) {3} --> r-tainted: read(product_db._.productID) {1}
[0] (BasicObject BasicType) productID string

    --> w-tainted: write(order_db.Order.Quantity, stock_db.Stock.Quantity, billing_db.Bill.Quantity) {3}
[0] (BasicObject BasicType) quantity int

    --> w-tainted: write(billing_db.Bill.PricePerUnit) {1}
[0] (BasicObject BasicType) price int

