[0] (PointerObject PointerType) s (*coupons_app_sql.StudentServiceImpl struct{studentsDB RelationalDB})
[_1] (StructObject UserType) coupons_app_sql.StudentServiceImpl struct{studentsDB RelationalDB}
[__2] (FieldObject FieldType) studentsDB RelationalDB
[___3] (BlueprintBackendObject BlueprintBackendType) studentsDB RelationalDB

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

[0] (BasicObject BasicType) studentID int
[_1] (Reference BasicType) ref <studentID int> @ Frontend

[0] (BasicObject BasicType) name string
[_1] (Reference BasicType) ref <name string> @ Frontend

[0] (InterfaceObject UserType) _ sql.Result interface{ interface{LastInsertId() (int64, error); RowsAffected() (int64, error)} }

[0] (InterfaceObject UserType) err .error

