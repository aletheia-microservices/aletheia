[0] (PointerObject PointerType) s (*coupons_app.CouponServiceImpl struct{couponsDB NoSQLDatabase})
[_1] (StructObject UserType) coupons_app.CouponServiceImpl struct{couponsDB NoSQLDatabase}
[__2] (FieldObject FieldType) couponsDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) couponsDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> w-tainted: write(coupons_db.ClaimedCoupon.CouponID) {1}
[0] (BasicObject BasicType) couponID int
     --> w-tainted: write(coupons_db.ClaimedCoupon.CouponID) {1}
[_1] (Reference BasicType) ref <couponID int> @ Frontend

[0] (BasicObject BasicType) category string
[_1] (Reference BasicType) ref <category string> @ Frontend

    --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1}
[0] (BasicObject BasicType) userID int
     --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1}
[_1] (Reference BasicType) ref <studentID int> @ Frontend

[0] (StructObject UserType) claimedCoupon coupons_app.ClaimedCoupon struct{CouponID int, UserID int}

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = coupons, collection = Coupon}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) filter primitive.D
[_1] (StructObject StructType) struct{Key "CouponID" string, Key "CouponID" string, Value int, Value int}
[__2] (FieldObject FieldType) Key "CouponID" string
[___3] (BasicObject BasicType) "CouponID" string
[__2] (FieldObject FieldType) Value int
       --> w-tainted: write(coupons_db.ClaimedCoupon.CouponID) {1}
[___3] (BasicObject BasicType) couponID int
        --> w-tainted: write(coupons_db.ClaimedCoupon.CouponID) {1}
[____4] (Reference BasicType) ref <couponID int> @ Frontend

[0] (SliceObject UserType) update primitive.D
[_1] (StructObject StructType) struct{Key "$inc" string, Key "$inc" string, Value primitive.D, Value primitive.D}
[__2] (FieldObject FieldType) Key "$inc" string
[___3] (BasicObject BasicType) "$inc" string
[__2] (FieldObject FieldType) Value primitive.D
[___3] (SliceObject UserType) primitive.D
[____4] (StructObject StructType) struct{Key "NumClaims" string, Key "NumClaims" string, Value 1 int, Value 1 int}
[_____5] (FieldObject FieldType) Key "NumClaims" string
[______6] (BasicObject BasicType) "NumClaims" string
[_____5] (FieldObject FieldType) Value 1 int
[______6] (BasicObject BasicType) 1 int

[0] (BasicObject BasicType) res int

[0] (InterfaceObject UserType) err .error

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = coupons, collection = ClaimedCoupon}

[0] (InterfaceObject UserType) err .error

    --> w-tainted: write(coupons_db.ClaimedCoupon) {1}
[0] (StructObject UserType) claimedCoupon coupons_app.ClaimedCoupon struct{CouponID int, UserID int}
     --> w-tainted: write(coupons_db.ClaimedCoupon.CouponID) {1}
[_1] (FieldObject FieldType) CouponID int
      --> w-tainted: write(coupons_db.ClaimedCoupon.CouponID) {1}
[__2] (BasicObject BasicType) couponID int
       --> w-tainted: write(coupons_db.ClaimedCoupon.CouponID) {1}
[___3] (Reference BasicType) ref <couponID int> @ Frontend
     --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1}
[_1] (FieldObject FieldType) UserID int
      --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1}
[__2] (BasicObject BasicType) userID int
       --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1}
[___3] (Reference BasicType) ref <studentID int> @ Frontend

[0] (InterfaceObject UserType) err .error

