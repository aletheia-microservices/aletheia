[0] (PointerObject PointerType) s (*coupons_app_sql.StudentServiceImpl struct{studentsDB RelationalDB})
[_1] (StructObject UserType) coupons_app_sql.StudentServiceImpl struct{studentsDB RelationalDB}
[__2] (FieldObject FieldType) studentsDB RelationalDB
[___3] (BlueprintBackendObject BlueprintBackendType) studentsDB RelationalDB

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

[0] (BasicObject BasicType) studentID int
[_1] (Reference BasicType) ref <studentID int> @ Frontend

[0] (BasicObject BasicType) couponValue int
[_1] (Reference BasicType) ref <value int> @ Frontend
[__2] (Reference BasicType) ref <Value int> @ CouponService

[0] (StructObject UserType) student coupons_app_sql.Student struct{StudentID int, Name string, Balance int}

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) newBalance int

[0] (InterfaceObject UserType) _ sql.Result interface{ interface{LastInsertId() (int64, error); RowsAffected() (int64, error)} }

[0] (InterfaceObject UserType) err .error

