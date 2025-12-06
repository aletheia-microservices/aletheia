[0] (PointerObject PointerType) u (*coupons_app_sql.FrontendImpl struct{StudentService coupons_app_sql.StudentService, CouponService coupons_app_sql.CouponService})
[_1] (StructObject UserType) coupons_app_sql.FrontendImpl struct{StudentService coupons_app_sql.StudentService, CouponService coupons_app_sql.CouponService}
[__2] (FieldObject FieldType) CouponService coupons_app_sql.CouponService
[___3] (ServiceObject ServiceType) CouponService coupons_app_sql.CouponService
[__2] (FieldObject FieldType) StudentService coupons_app_sql.StudentService
[___3] (ServiceObject ServiceType) StudentService coupons_app_sql.StudentService

[0] (InterfaceObject UserType) ctx context.Context

    --> w-tainted: write(coupons_db.claimed_coupons.coupon_id) {1}       --> w-tainted: write(coupons_db.claimed_coupons.coupon_id) {1} --> r-tainted: read(coupons_db.int) {1}
[0] (BasicObject BasicType) couponID int

    --> w-tainted: write(coupons_db.claimed_coupons.user_id) {1}       --> w-tainted: write(coupons_db.claimed_coupons.user_id) {1} --> r-tainted: read(students_db.int) {1}
[0] (BasicObject BasicType) studentID int

[0] (BasicObject BasicType) value int
     --> r-tainted: read(coupons_db.int) {1}
[_1] (Reference BasicType) ref <Value int> @ CouponService

[0] (InterfaceObject UserType) err .error
[_1] (Reference UserType) ref <err .error> @ CouponService

[0] (InterfaceObject UserType) err .error
[_1] (Reference UserType) ref <err .error> @ StudentService

