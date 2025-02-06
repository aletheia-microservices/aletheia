[0] (PointerObject PointerType) s (*coupons_app.StudentServiceImpl struct{studentsDB NoSQLDatabase})
[_1] (StructObject UserType) coupons_app.StudentServiceImpl struct{studentsDB NoSQLDatabase}
[__2] (FieldObject FieldType) studentsDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) studentsDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

[0] (BasicObject BasicType) studentID int
     --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1}         --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1} --> r-tainted: read(coupons_db.Coupon.userID) {1}
[_1] (Reference BasicType) ref <studentID int> @ Frontend

    --> w-tainted: write(students_db.Student.Balance) {1}       --> w-tainted: write(students_db.Student.Balance) {1} --> r-tainted: read(students_db.Student.StudentID) {1}
[0] (BasicObject BasicType) value int
     --> w-tainted: write(students_db.Student.Balance) {1}         --> w-tainted: write(students_db.Student.Balance) {1} --> r-tainted: read(students_db.Student.StudentID) {1}
[_1] (Reference BasicType) ref <value int> @ Frontend

    --> w-tainted: write(students_db.Student) {1}       --> w-tainted: write(students_db.Student) {1} --> r-tainted: read(students_db.Student, students_db.Student.StudentID) {2}
[0] (StructObject UserType) student coupons_app.Student struct{StudentID int, Name string, Balance int}
     --> w-tainted: write(students_db.Student.Balance) {1}         --> w-tainted: write(students_db.Student.Balance) {1} --> r-tainted: read(students_db.Student.StudentID) {1}
[_1] (FieldObject FieldType) Balance int
      --> w-tainted: write(students_db.Student.Balance) {1}           --> w-tainted: write(students_db.Student.Balance) {1} --> r-tainted: read(students_db.Student.StudentID) {1}
[__2] (BasicObject BasicType) value int
       --> w-tainted: write(students_db.Student.Balance) {1}             --> w-tainted: write(students_db.Student.Balance) {1} --> r-tainted: read(students_db.Student.StudentID) {1}
[___3] (Reference BasicType) ref <value int> @ Frontend

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = students, collection = students}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "studentID" string, Key "studentID" string, Value coupons_app.Student struct{StudentID int, Name string, Balance int}, Value coupons_app.Student struct{StudentID int, Name string, Balance int}}
[__2] (FieldObject FieldType) Key "studentID" string
[___3] (BasicObject BasicType) "studentID" string
[__2] (FieldObject FieldType) Value coupons_app.Student struct{StudentID int, Name string, Balance int}
       --> w-tainted: write(students_db.Student) {1}             --> w-tainted: write(students_db.Student) {1} --> r-tainted: read(students_db.Student, students_db.Student.StudentID) {2}
[___3] (StructObject UserType) student coupons_app.Student struct{StudentID int, Name string, Balance int}
        --> w-tainted: write(students_db.Student.Balance) {1}               --> w-tainted: write(students_db.Student.Balance) {1} --> r-tainted: read(students_db.Student.StudentID) {1}
[____4] (FieldObject FieldType) Balance int
         --> w-tainted: write(students_db.Student.Balance) {1}                 --> w-tainted: write(students_db.Student.Balance) {1} --> r-tainted: read(students_db.Student.StudentID) {1}
[_____5] (BasicObject BasicType) value int
          --> w-tainted: write(students_db.Student.Balance) {1}                   --> w-tainted: write(students_db.Student.Balance) {1} --> r-tainted: read(students_db.Student.StudentID) {1}
[______6] (Reference BasicType) ref <value int> @ Frontend

    --> r-tainted: read(students_db.Student) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = students, collection = students}

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) found bool

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

