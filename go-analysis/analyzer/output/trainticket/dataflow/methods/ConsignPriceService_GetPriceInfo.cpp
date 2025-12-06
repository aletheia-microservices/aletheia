[0] (PointerObject PointerType) c (*consignprice.ConsignPriceServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) consignprice.ConsignPriceServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = consign, collection = consign}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ 0 int,  "index" string,  0 int}

    --> r-tainted: read(consignprice_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = consign, collection = consign}
     --> r-tainted: read(consignprice_db.ConsignPrice, consignprice_db.ConsignPrice.BeyondPrice, consignprice_db.ConsignPrice.InitialPrice, consignprice_db.ConsignPrice.InitialWeight, consignprice_db.ConsignPrice.WithinPrice) {5}
[_1] (StructObject UserType) cp consignprice.ConsignPrice struct{ID string, Index int64, InitialWeight float64, InitialPrice float64, WithinPrice float64, BeyondPrice float64}
      --> r-tainted: read(consignprice_db.ConsignPrice.BeyondPrice) {1}
[__2] (FieldObject FieldType) BeyondPrice float64
       --> r-tainted: read(consignprice_db.ConsignPrice.BeyondPrice) {1}
[___3] (BasicObject BasicType) BeyondPrice float64
      --> r-tainted: read(consignprice_db.ConsignPrice.InitialPrice) {1}
[__2] (FieldObject FieldType) InitialPrice float64
       --> r-tainted: read(consignprice_db.ConsignPrice.InitialPrice) {1}
[___3] (BasicObject BasicType) InitialPrice float64
      --> r-tainted: read(consignprice_db.ConsignPrice.InitialWeight) {1}
[__2] (FieldObject FieldType) InitialWeight float64
       --> r-tainted: read(consignprice_db.ConsignPrice.InitialWeight) {1}
[___3] (BasicObject BasicType) InitialWeight float64
      --> r-tainted: read(consignprice_db.ConsignPrice.WithinPrice) {1}
[__2] (FieldObject FieldType) WithinPrice float64
       --> r-tainted: read(consignprice_db.ConsignPrice.WithinPrice) {1}
[___3] (BasicObject BasicType) WithinPrice float64

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(consignprice_db.ConsignPrice, consignprice_db.ConsignPrice.BeyondPrice, consignprice_db.ConsignPrice.InitialPrice, consignprice_db.ConsignPrice.InitialWeight, consignprice_db.ConsignPrice.WithinPrice) {5}
[0] (StructObject UserType) cp consignprice.ConsignPrice struct{ID string, Index int64, InitialWeight float64, InitialPrice float64, WithinPrice float64, BeyondPrice float64}
     --> r-tainted: read(consignprice_db.ConsignPrice.BeyondPrice) {1}
[_1] (FieldObject FieldType) BeyondPrice float64
      --> r-tainted: read(consignprice_db.ConsignPrice.BeyondPrice) {1}
[__2] (BasicObject BasicType) BeyondPrice float64
     --> r-tainted: read(consignprice_db.ConsignPrice.InitialPrice) {1}
[_1] (FieldObject FieldType) InitialPrice float64
      --> r-tainted: read(consignprice_db.ConsignPrice.InitialPrice) {1}
[__2] (BasicObject BasicType) InitialPrice float64
     --> r-tainted: read(consignprice_db.ConsignPrice.InitialWeight) {1}
[_1] (FieldObject FieldType) InitialWeight float64
      --> r-tainted: read(consignprice_db.ConsignPrice.InitialWeight) {1}
[__2] (BasicObject BasicType) InitialWeight float64
     --> r-tainted: read(consignprice_db.ConsignPrice.WithinPrice) {1}
[_1] (FieldObject FieldType) WithinPrice float64
      --> r-tainted: read(consignprice_db.ConsignPrice.WithinPrice) {1}
[__2] (BasicObject BasicType) WithinPrice float64

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) info string
[_1] (BasicObject BasicType) "The price of weight within %.2f is %.2f. The price of extra weight within the region is %.2f and beyond the region is %.2f" string
     --> r-tainted: read(consignprice_db.ConsignPrice.InitialWeight) {1}
[_1] (BasicObject BasicType) InitialWeight float64
     --> r-tainted: read(consignprice_db.ConsignPrice.InitialPrice) {1}
[_1] (BasicObject BasicType) InitialPrice float64
     --> r-tainted: read(consignprice_db.ConsignPrice.WithinPrice) {1}
[_1] (BasicObject BasicType) WithinPrice float64
     --> r-tainted: read(consignprice_db.ConsignPrice.BeyondPrice) {1}
[_1] (BasicObject BasicType) BeyondPrice float64
[_1] (BasicObject BasicType) "The price of weight within %.2f is %.2f. The price of extra weight within the region is %.2f and beyond the region is %.2f" string
     --> r-tainted: read(consignprice_db.ConsignPrice.InitialWeight) {1}
[_1] (BasicObject BasicType) InitialWeight float64
     --> r-tainted: read(consignprice_db.ConsignPrice.InitialPrice) {1}
[_1] (BasicObject BasicType) InitialPrice float64
     --> r-tainted: read(consignprice_db.ConsignPrice.WithinPrice) {1}
[_1] (BasicObject BasicType) WithinPrice float64
     --> r-tainted: read(consignprice_db.ConsignPrice.BeyondPrice) {1}
[_1] (BasicObject BasicType) BeyondPrice float64

