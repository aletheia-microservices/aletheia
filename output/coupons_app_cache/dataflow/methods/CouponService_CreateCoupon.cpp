[0] (PointerObject PointerType) s (*coupons_app_cache.CouponServiceImpl struct{couponsDB RelationalDB})
[_1] (StructObject UserType) coupons_app_cache.CouponServiceImpl struct{couponsDB RelationalDB}
[__2] (FieldObject FieldType) couponsDB RelationalDB
[___3] (BlueprintBackendObject BlueprintBackendType) couponsDB RelationalDB

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> w-tainted: write(coupons_db.coupons.category) {1}
[0] (BasicObject BasicType) category string
     --> w-tainted: write(coupons_db.coupons.category) {1}
[_1] (Reference BasicType) ref <category string> @ Frontend

    --> w-tainted: write(coupons_db.coupons.value) {1}
[0] (BasicObject BasicType) value int
     --> w-tainted: write(coupons_db.coupons.value) {1}
[_1] (Reference BasicType) ref <value int> @ Frontend

[0] (InterfaceObject UserType) _ sql.Result interface{ interface{LastInsertId() (int64, error); RowsAffected() (int64, error)} }

[0] (InterfaceObject UserType) err .error

