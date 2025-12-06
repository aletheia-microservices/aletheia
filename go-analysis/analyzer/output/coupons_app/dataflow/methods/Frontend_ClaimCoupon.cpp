[0] (PointerObject PointerType) u (*coupons_app.FrontendImpl struct{StudentService coupons_app.StudentService, CouponService coupons_app.CouponService})
[_1] (StructObject UserType) coupons_app.FrontendImpl struct{StudentService coupons_app.StudentService, CouponService coupons_app.CouponService}
[__2] (FieldObject FieldType) CouponService coupons_app.CouponService
[___3] (ServiceObject ServiceType) CouponService coupons_app.CouponService
[__2] (FieldObject FieldType) StudentService coupons_app.StudentService
[___3] (ServiceObject ServiceType) StudentService coupons_app.StudentService

[0] (InterfaceObject UserType) ctx context.Context

    --> w-tainted: write(coupons_db.ClaimedCoupon.CouponID) {1}
[0] (BasicObject BasicType) couponID int

[0] (BasicObject BasicType) category string

    --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1}
[0] (BasicObject BasicType) studentID int

[0] (BasicObject BasicType) value int

[0] (InterfaceObject UserType) err .error
[_1] (Reference UserType) ref <err .error> @ CouponService

[0] (InterfaceObject UserType) err .error
[_1] (Reference BasicType) ref <nil> @ StudentService

