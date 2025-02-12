[0] (PointerObject PointerType) s (*coupons_app_sql.CouponServiceImpl struct{couponsDB RelationalDB})
[_1] (StructObject UserType) coupons_app_sql.CouponServiceImpl struct{couponsDB RelationalDB}
[__2] (FieldObject FieldType) couponsDB RelationalDB
[___3] (BlueprintBackendObject BlueprintBackendType) couponsDB RelationalDB

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

[0] (BasicObject BasicType) couponID int
[_1] (Reference BasicType) ref <couponID int> @ Frontend

[0] (BasicObject BasicType) userID int
[_1] (Reference BasicType) ref <studentID int> @ Frontend

[0] (StructObject UserType) coupon coupons_app_sql.Coupon struct{CouponID int, Category string, Value int}
[_1] (FieldObject FieldType) Value int
[__2] (BasicObject BasicType) Value int

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) _ sql.Result interface{ interface{LastInsertId() (int64, error); RowsAffected() (int64, error)} }

[0] (InterfaceObject UserType) err .error

