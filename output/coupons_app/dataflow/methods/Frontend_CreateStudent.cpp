[0] (PointerObject PointerType) u (*coupons_app.FrontendImpl struct{StudentService coupons_app.StudentService, CouponService coupons_app.CouponService})
[_1] (StructObject UserType) coupons_app.FrontendImpl struct{StudentService coupons_app.StudentService, CouponService coupons_app.CouponService}
[__2] (FieldObject FieldType) CouponService coupons_app.CouponService
[___3] (ServiceObject ServiceType) CouponService coupons_app.CouponService
[__2] (FieldObject FieldType) StudentService coupons_app.StudentService
[___3] (ServiceObject ServiceType) StudentService coupons_app.StudentService

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) studentID int

    --> w-tainted: write(students_db.Student.Name) {1}
[0] (BasicObject BasicType) name string

    --> w-tainted: write(students_db.Student) {1}
[0] (StructObject UserType) student coupons_app.Student struct{StudentID int, Name string, Balance int, NumClaimedCoupons int}
     --> w-tainted: write(students_db.Student) {1}
[_1] (Reference UserType) ref <coupon coupons_app.Student struct{StudentID int, Name string, Balance int, NumClaimedCoupons 0 int}> @ StudentService
      --> w-tainted: write(students_db.Student.Name) {1}
[__2] (FieldObject FieldType) Name string
       --> w-tainted: write(students_db.Student.Name) {1}
[___3] (BasicObject BasicType) name string
        --> w-tainted: write(students_db.Student.Name) {1}
[____4] (Reference BasicType) ref <name string> @ Frontend
      --> w-tainted: write(students_db.Student.NumClaimedCoupons) {1}
[__2] (FieldObject FieldType) NumClaimedCoupons 0 int
       --> w-tainted: write(students_db.Student.NumClaimedCoupons) {1}
[___3] (BasicObject BasicType) 0 int

[0] (InterfaceObject UserType) err .error
[_1] (Reference UserType) ref <err .error> @ StudentService

