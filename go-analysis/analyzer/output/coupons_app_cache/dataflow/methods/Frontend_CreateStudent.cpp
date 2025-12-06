[0] (PointerObject PointerType) u (*coupons_app_cache.FrontendImpl struct{StudentService coupons_app_cache.StudentService, CouponService coupons_app_cache.CouponService})
[_1] (StructObject UserType) coupons_app_cache.FrontendImpl struct{StudentService coupons_app_cache.StudentService, CouponService coupons_app_cache.CouponService}
[__2] (FieldObject FieldType) CouponService coupons_app_cache.CouponService
[___3] (ServiceObject ServiceType) CouponService coupons_app_cache.CouponService
[__2] (FieldObject FieldType) StudentService coupons_app_cache.StudentService
[___3] (ServiceObject ServiceType) StudentService coupons_app_cache.StudentService

[0] (InterfaceObject UserType) ctx context.Context

    --> w-tainted: write(students_db.students.name) {1}
[0] (BasicObject BasicType) name string

[0] (InterfaceObject UserType) err .error
[_1] (Reference UserType) ref <err .error> @ StudentService

