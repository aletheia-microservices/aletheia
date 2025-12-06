[0] (PointerObject PointerType) s (*coupons_app.StudentServiceImpl struct{studentsDB NoSQLDatabase})
[_1] (StructObject UserType) coupons_app.StudentServiceImpl struct{studentsDB NoSQLDatabase}
[__2] (FieldObject FieldType) studentsDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) studentsDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

[0] (BasicObject BasicType) studentID int
[_1] (Reference BasicType) ref <studentID int> @ Frontend

    --> w-tainted: write(students_db.Student.Name) {1}
[0] (BasicObject BasicType) name string
     --> w-tainted: write(students_db.Student.Name) {1}
[_1] (Reference BasicType) ref <name string> @ Frontend

    --> w-tainted: write(students_db.Student) {1}
[0] (StructObject UserType) coupon coupons_app.Student struct{StudentID int, Name string, Balance int, NumClaimedCoupons 0 int}
     --> w-tainted: write(students_db.Student.Name) {1}
[_1] (FieldObject FieldType) Name string
      --> w-tainted: write(students_db.Student.Name) {1}
[__2] (BasicObject BasicType) name string
       --> w-tainted: write(students_db.Student.Name) {1}
[___3] (Reference BasicType) ref <name string> @ Frontend
     --> w-tainted: write(students_db.Student.NumClaimedCoupons) {1}
[_1] (FieldObject FieldType) NumClaimedCoupons 0 int
      --> w-tainted: write(students_db.Student.NumClaimedCoupons) {1}
[__2] (BasicObject BasicType) 0 int

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = students, collection = Student}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

