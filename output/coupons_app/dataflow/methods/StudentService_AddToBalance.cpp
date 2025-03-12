[0] (PointerObject PointerType) s (*coupons_app.StudentServiceImpl struct{studentsDB NoSQLDatabase})
[_1] (StructObject UserType) coupons_app.StudentServiceImpl struct{studentsDB NoSQLDatabase}
[__2] (FieldObject FieldType) studentsDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) studentsDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

[0] (BasicObject BasicType) studentID int
     --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1}
[_1] (Reference BasicType) ref <studentID int> @ Frontend

[0] (BasicObject BasicType) value int
[_1] (Reference BasicType) ref <value int> @ Frontend

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = students, collection = Student}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) filter primitive.D
[_1] (StructObject StructType) struct{Key "StudentID" string, Key "StudentID" string, Value int, Value int}
[__2] (FieldObject FieldType) Key "StudentID" string
[___3] (BasicObject BasicType) "StudentID" string
[__2] (FieldObject FieldType) Value int
[___3] (BasicObject BasicType) studentID int
        --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1}
[____4] (Reference BasicType) ref <studentID int> @ Frontend

[0] (SliceObject UserType) update primitive.D
[_1] (StructObject StructType) struct{Key "$inc" string, Key "$inc" string, Value primitive.D, Value primitive.D}
[__2] (FieldObject FieldType) Key "$inc" string
[___3] (BasicObject BasicType) "$inc" string
[__2] (FieldObject FieldType) Value primitive.D
[___3] (SliceObject UserType) primitive.D
[____4] (StructObject StructType) struct{Key "Balance" string, Key "Balance" string, Value 1 int, Value 1 int}
[_____5] (FieldObject FieldType) Key "Balance" string
[______6] (BasicObject BasicType) "Balance" string
[_____5] (FieldObject FieldType) Value 1 int
[______6] (BasicObject BasicType) 1 int
[_1] (StructObject StructType) struct{Key "$inc" string, Key "$inc" string, Value primitive.D, Value primitive.D}
[__2] (FieldObject FieldType) Key "$inc" string
[___3] (BasicObject BasicType) "$inc" string
[__2] (FieldObject FieldType) Value primitive.D
[___3] (SliceObject UserType) primitive.D
[____4] (StructObject StructType) struct{Key "ClaimedCoupons" string, Key "ClaimedCoupons" string, Value 1 int, Value 1 int}
[_____5] (FieldObject FieldType) Key "ClaimedCoupons" string
[______6] (BasicObject BasicType) "ClaimedCoupons" string
[_____5] (FieldObject FieldType) Value 1 int
[______6] (BasicObject BasicType) 1 int

[0] (BasicObject BasicType) res int

[0] (InterfaceObject UserType) err .error

