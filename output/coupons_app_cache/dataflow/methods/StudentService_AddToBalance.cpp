[0] (PointerObject PointerType) s (*coupons_app_cache.StudentServiceImpl struct{studentsDB RelationalDB})
[_1] (StructObject UserType) coupons_app_cache.StudentServiceImpl struct{studentsDB RelationalDB}
[__2] (FieldObject FieldType) studentsDB RelationalDB
[___3] (BlueprintBackendObject BlueprintBackendType) studentsDB RelationalDB

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> r-tainted: read(students_db.students.student_id, students_db.int) {2}
[0] (BasicObject BasicType) studentID int
     --> w-tainted: write(coupons_db.claimed_coupons.user_id) {1}         --> w-tainted: write(coupons_db.claimed_coupons.user_id) {1} --> r-tainted: read(students_db.int) {1}
[_1] (Reference BasicType) ref <studentID int> @ Frontend

[0] (BasicObject BasicType) couponValue int
[_1] (Reference BasicType) ref <value int> @ Frontend
      --> r-tainted: read(coupons_db.int) {1}
[__2] (Reference BasicType) ref <Value int> @ CouponService

    --> r-tainted: read(students_db.Student) {1}
[0] (StructObject UserType) student coupons_app_cache.Student struct{StudentID int, Name string, Balance int}

[0] (InterfaceObject UserType) err .error

    --> w-tainted: write(students_db.students.balance) {1}
[0] (BasicObject BasicType) newBalance int

[0] (InterfaceObject UserType) _ sql.Result interface{ interface{LastInsertId() (int64, error); RowsAffected() (int64, error)} }

[0] (InterfaceObject UserType) err .error

