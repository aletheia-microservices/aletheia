[0] (PointerObject PointerType) c (*consignprice.ConsignPriceServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) consignprice.ConsignPriceServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ ConsignPriceService

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = consign, collection = consign}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ 0 int,  "index" string,  0 int}

    --> r-tainted: read(consignprice_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = consign, collection = consign}
     --> r-tainted: read(consignprice_db.ConsignPrice) {1}
[_1] (StructObject UserType) cp consignprice.ConsignPrice struct{ID string, Index int64, InitialWeight float64, InitialPrice float64, WithinPrice float64, BeyondPrice float64}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(consignprice_db.ConsignPrice) {1}
[0] (StructObject UserType) cp consignprice.ConsignPrice struct{ID string, Index int64, InitialWeight float64, InitialPrice float64, WithinPrice float64, BeyondPrice float64}

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

