[0] (PointerObject PointerType) s (*coupons_app_sql.CouponServiceImpl struct{couponsDB RelationalDB})
[_1] (StructObject UserType) coupons_app_sql.CouponServiceImpl struct{couponsDB RelationalDB}
[__2] (FieldObject FieldType) couponsDB RelationalDB
[___3] (BlueprintBackendObject BlueprintBackendType) couponsDB RelationalDB

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> w-tainted: write(coupons_db.claimed_coupons.coupon_id) {1}       --> w-tainted: write(coupons_db.claimed_coupons.coupon_id) {1} --> r-tainted: read(coupons_db.claimed_coupons.coupon_id, coupons_db.int) {2}
[0] (BasicObject BasicType) couponID int
     --> w-tainted: write(coupons_db.claimed_coupons.coupon_id) {1}         --> w-tainted: write(coupons_db.claimed_coupons.coupon_id) {1} --> r-tainted: read(coupons_db.int) {1}
[_1] (Reference BasicType) ref <couponID int> @ Frontend

    --> w-tainted: write(coupons_db.claimed_coupons.user_id) {1}
[0] (BasicObject BasicType) userID int
     --> w-tainted: write(coupons_db.claimed_coupons.user_id) {1}         --> w-tainted: write(coupons_db.claimed_coupons.user_id) {1} --> r-tainted: read(students_db.int) {1}
[_1] (Reference BasicType) ref <studentID int> @ Frontend

    --> r-tainted: read(coupons_db.Coupon) {1}
[0] (StructObject UserType) coupon coupons_app_sql.Coupon struct{CouponID int, Category string, Value int}
     --> r-tainted: read(coupons_db.Value) {1}
[_1] (FieldObject FieldType) Value int
      --> r-tainted: read(coupons_db.int) {1}
[__2] (BasicObject BasicType) Value int

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) _ sql.Result interface{ interface{LastInsertId() (int64, error); RowsAffected() (int64, error)} }

[0] (InterfaceObject UserType) err .error

