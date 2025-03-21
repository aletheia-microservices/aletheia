[0] (PointerObject PointerType) f (*shopping_simple.FrontendImpl struct{product_service shopping_simple.ProductService, cart_service shopping_simple.CartService})
[_1] (StructObject UserType) shopping_simple.FrontendImpl struct{product_service shopping_simple.ProductService, cart_service shopping_simple.CartService}
[__2] (FieldObject FieldType) cart_service shopping_simple.CartService
[___3] (ServiceObject ServiceType) cart_service shopping_simple.CartService
[__2] (FieldObject FieldType) product_service shopping_simple.ProductService
[___3] (ServiceObject ServiceType) product_service shopping_simple.ProductService

[0] (InterfaceObject UserType) ctx context.Context

    --> w-tainted: write(cart_db.Cart.CartID) {1}       --> w-tainted: write(cart_db.Cart.CartID) {1} --> r-tainted: read(cart_db.Cart.CartID) {1}
[0] (BasicObject BasicType) cartID string

    --> w-tainted: write(cart_db.Cart.Products) {1}       --> w-tainted: write(cart_db.Cart.Products) {1} --> r-tainted: read(product_db.Product.ProductID) {1}
[0] (BasicObject BasicType) productID string

