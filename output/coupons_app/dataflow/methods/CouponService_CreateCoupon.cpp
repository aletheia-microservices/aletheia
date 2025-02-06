[0] (PointerObject PointerType) s (*coupons_app.CouponServiceImpl struct{couponsDB NoSQLDatabase})
[_1] (StructObject UserType) coupons_app.CouponServiceImpl struct{couponsDB NoSQLDatabase}
[__2] (FieldObject FieldType) couponsDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) couponsDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> w-tainted: write(coupons_db.Coupon.CouponID) {1}
[0] (BasicObject BasicType) couponID int
     --> w-tainted: write(coupons_db.Coupon.CouponID) {1}
[_1] (Reference BasicType) ref <couponID int> @ Frontend

    --> w-tainted: write(coupons_db.Coupon.Category) {1}
[0] (BasicObject BasicType) category string
     --> w-tainted: write(coupons_db.Coupon.Category) {1}
[_1] (Reference BasicType) ref <category string> @ Frontend

    --> w-tainted: write(coupons_db.Coupon) {1}
[0] (StructObject UserType) coupon coupons_app.Coupon struct{CouponID int, Category string}
     --> w-tainted: write(coupons_db.Coupon.Category) {1}
[_1] (FieldObject FieldType) Category string
      --> w-tainted: write(coupons_db.Coupon.Category) {1}
[__2] (BasicObject BasicType) category string
       --> w-tainted: write(coupons_db.Coupon.Category) {1}
[___3] (Reference BasicType) ref <category string> @ Frontend
     --> w-tainted: write(coupons_db.Coupon.CouponID) {1}
[_1] (FieldObject FieldType) CouponID int
      --> w-tainted: write(coupons_db.Coupon.CouponID) {1}
[__2] (BasicObject BasicType) couponID int
       --> w-tainted: write(coupons_db.Coupon.CouponID) {1}
[___3] (Reference BasicType) ref <couponID int> @ Frontend

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = coupons, collection = coupons}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

