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

    --> r-tainted: read(coupons_db.Coupon.Category) {1}
[0] (BasicObject BasicType) category string
     --> r-tainted: read(coupons_db.Coupon.Category) {1}
[_1] (Reference BasicType) ref <category string> @ Frontend

    --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1}       --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1} --> r-tainted: read(coupons_db.Coupon.userID) {1}
[0] (BasicObject BasicType) userID int
     --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1}         --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1} --> r-tainted: read(coupons_db.Coupon.userID) {1}
[_1] (Reference BasicType) ref <studentID int> @ Frontend

    --> r-tainted: read(coupons_db.ClaimedCoupon) {1}
[0] (StructObject UserType) claimedCoupon coupons_app.ClaimedCoupon struct{CouponID int, UserID int}
[_1] (FieldObject FieldType) CouponID int
[__2] (BasicObject BasicType) CouponID int
[_1] (FieldObject FieldType) UserID int
[__2] (BasicObject BasicType) UserID int

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = coupons, collection = claimed_coupons}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "category" string, Key "category" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "category" string
[___3] (BasicObject BasicType) "category" string
[__2] (FieldObject FieldType) Value string
       --> r-tainted: read(coupons_db.Coupon.Category) {1}
[___3] (BasicObject BasicType) category string
        --> r-tainted: read(coupons_db.Coupon.Category) {1}
[____4] (Reference BasicType) ref <category string> @ Frontend
[_1] (StructObject StructType) struct{Key "userID" string, Key "userID" string, Value int, Value int}
[__2] (FieldObject FieldType) Key "userID" string
[___3] (BasicObject BasicType) "userID" string
[__2] (FieldObject FieldType) Value int
       --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1}             --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1} --> r-tainted: read(coupons_db.Coupon.userID) {1}
[___3] (BasicObject BasicType) userID int
        --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1}               --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1} --> r-tainted: read(coupons_db.Coupon.userID) {1}
[____4] (Reference BasicType) ref <studentID int> @ Frontend

    --> r-tainted: read(coupons_db.Coupon) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = coupons, collection = claimed_coupons}

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) found bool

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
      --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1}           --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1} --> r-tainted: read(coupons_db.Coupon.userID) {1}
[__2] (BasicObject BasicType) userID int
       --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1}             --> w-tainted: write(coupons_db.ClaimedCoupon.UserID) {1} --> r-tainted: read(coupons_db.Coupon.userID) {1}
[___3] (Reference BasicType) ref <studentID int> @ Frontend

[0] (InterfaceObject UserType) err .error

