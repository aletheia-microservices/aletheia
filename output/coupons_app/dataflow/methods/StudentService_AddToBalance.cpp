[0] (PointerObject PointerType) s (*coupons_app.StudentServiceImpl struct{studentsDB NoSQLDatabase})
[_1] (StructObject UserType) coupons_app.StudentServiceImpl struct{studentsDB NoSQLDatabase}
[__2] (FieldObject FieldType) studentsDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) studentsDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

[0] (BasicObject BasicType) studentID int
     --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1}         --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1} --> r-tainted: read(coupons_db.Coupon.userID) {1}
[_1] (Reference BasicType) ref <studentID int> @ Frontend

[0] (BasicObject BasicType) value int
[_1] (Reference BasicType) ref <value int> @ Frontend

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = students, collection = students}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) filter primitive.D
[_1] (StructObject StructType) struct{Key "studentID" string, Key "studentID" string, Value int, Value int}
[__2] (FieldObject FieldType) Key "studentID" string
[___3] (BasicObject BasicType) "studentID" string
[__2] (FieldObject FieldType) Value int
[___3] (BasicObject BasicType) studentID int
        --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1}               --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1} --> r-tainted: read(coupons_db.Coupon.userID) {1}
[____4] (Reference BasicType) ref <studentID int> @ Frontend

[0] (SliceObject UserType) update primitive.D
[_1] (StructObject StructType) struct{Key "$inc" string, Key "$inc" string, Value primitive.D, Value primitive.D}
[__2] (FieldObject FieldType) Key "$inc" string
[___3] (BasicObject BasicType) "$inc" string
[__2] (FieldObject FieldType) Value primitive.D
[___3] (SliceObject UserType) primitive.D
[____4] (StructObject StructType) struct{Key "balance" string, Key "balance" string, Value 1 int, Value 1 int}
[_____5] (FieldObject FieldType) Key "balance" string
[______6] (BasicObject BasicType) "balance" string
[_____5] (FieldObject FieldType) Value 1 int
[______6] (BasicObject BasicType) 1 int

[0] (BasicObject BasicType) res int

[0] (InterfaceObject UserType) err .error

