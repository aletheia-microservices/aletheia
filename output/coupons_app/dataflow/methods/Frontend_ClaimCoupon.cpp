[0] (PointerObject PointerType) u (*coupons_app.FrontendImpl struct{StudentService coupons_app.StudentService, CouponService coupons_app.CouponService})
[_1] (StructObject UserType) coupons_app.FrontendImpl struct{StudentService coupons_app.StudentService, CouponService coupons_app.CouponService}
[__2] (FieldObject FieldType) CouponService coupons_app.CouponService
[___3] (ServiceObject ServiceType) CouponService coupons_app.CouponService
[__2] (FieldObject FieldType) StudentService coupons_app.StudentService
[___3] (ServiceObject ServiceType) StudentService coupons_app.StudentService

[0] (InterfaceObject UserType) ctx context.Context

    --> w-tainted: write(coupons_db.ClaimedCoupon.CouponID) {1}
[0] (BasicObject BasicType) couponID int

    --> r-tainted: read(coupons_db.Coupon.Category) {1}
[0] (BasicObject BasicType) category string

    --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1}       --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1} --> r-tainted: read(coupons_db.Coupon.userID) {1}
[0] (BasicObject BasicType) studentID int

    --> w-tainted: write(students_db.Student.Balance) {1}       --> w-tainted: write(students_db.Student.Balance) {1} --> r-tainted: read(students_db.Student.StudentID) {1}
[0] (BasicObject BasicType) value int

[0] (InterfaceObject UserType) err .error
[_1] (Reference UserType) ref <err .error> @ CouponService

    --> w-tainted: write(students_db.Student) {1}
[0] (StructObject UserType) student coupons_app.Student struct{StudentID int, Name string, Balance int}
     --> w-tainted: write(students_db.Student) {1}         --> w-tainted: write(students_db.Student) {1} --> r-tainted: read(students_db.Student, students_db.Student.StudentID) {2}
[_1] (Reference UserType) ref <student coupons_app.Student struct{StudentID int, Name string, Balance int}> @ StudentService
      --> w-tainted: write(students_db.Student.Balance) {1}           --> w-tainted: write(students_db.Student.Balance) {1} --> r-tainted: read(students_db.Student.StudentID) {1}
[__2] (FieldObject FieldType) Balance int
       --> w-tainted: write(students_db.Student.Balance) {1}             --> w-tainted: write(students_db.Student.Balance) {1} --> r-tainted: read(students_db.Student.StudentID) {1}
[___3] (BasicObject BasicType) value int
        --> w-tainted: write(students_db.Student.Balance) {1}               --> w-tainted: write(students_db.Student.Balance) {1} --> r-tainted: read(students_db.Student.StudentID) {1}
[____4] (Reference BasicType) ref <value int> @ Frontend

[0] (InterfaceObject UserType) err .error
[_1] (Reference UserType) ref <err .error> @ StudentService

