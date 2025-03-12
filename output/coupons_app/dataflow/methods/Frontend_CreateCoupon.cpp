[0] (PointerObject PointerType) u (*coupons_app.FrontendImpl struct{StudentService coupons_app.StudentService, CouponService coupons_app.CouponService})
[_1] (StructObject UserType) coupons_app.FrontendImpl struct{StudentService coupons_app.StudentService, CouponService coupons_app.CouponService}
[__2] (FieldObject FieldType) CouponService coupons_app.CouponService
[___3] (ServiceObject ServiceType) CouponService coupons_app.CouponService
[__2] (FieldObject FieldType) StudentService coupons_app.StudentService
[___3] (ServiceObject ServiceType) StudentService coupons_app.StudentService

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) couponID int

    --> w-tainted: write(coupons_db.Coupon.Category) {1}
[0] (BasicObject BasicType) category string

    --> w-tainted: write(coupons_db.Coupon) {1}
[0] (StructObject UserType) coupon coupons_app.Coupon struct{CouponID int, Category string, NumClaims int}
     --> w-tainted: write(coupons_db.Coupon) {1}
[_1] (Reference UserType) ref <coupon coupons_app.Coupon struct{CouponID int, Category string, NumClaims int}> @ CouponService
      --> w-tainted: write(coupons_db.Coupon.Category) {1}
[__2] (FieldObject FieldType) Category string
       --> w-tainted: write(coupons_db.Coupon.Category) {1}
[___3] (BasicObject BasicType) category string
        --> w-tainted: write(coupons_db.Coupon.Category) {1}
[____4] (Reference BasicType) ref <category string> @ Frontend

[0] (InterfaceObject UserType) err .error
[_1] (Reference UserType) ref <err .error> @ CouponService

