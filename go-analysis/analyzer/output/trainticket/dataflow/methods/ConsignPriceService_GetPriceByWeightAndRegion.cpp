[0] (PointerObject PointerType) c (*consignprice.ConsignPriceServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) consignprice.ConsignPriceServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) weight float64

[0] (BasicObject BasicType) isWithinRegion bool

[0] (StructObject UserType) cp consignprice.ConsignPrice struct{ID string, Index int64, InitialWeight float64, InitialPrice float64, WithinPrice float64, BeyondPrice float64}
     --> r-tainted: read(consignprice_db.ConsignPrice) {1}
[_1] (Reference UserType) ref <cp consignprice.ConsignPrice struct{ID string, Index int64, InitialWeight float64, InitialPrice float64, WithinPrice float64, BeyondPrice float64}> @ GetPriceConfig
[_1] (FieldObject FieldType) InitialPrice float64
[__2] (BasicObject BasicType) InitialPrice float64

[0] (InterfaceObject UserType) err .error
[_1] (Reference BasicType) ref <nil> @ GetPriceConfig

[0] (FieldObject FieldType) price float64
[_1] (FieldObject FieldType) InitialPrice float64
[__2] (BasicObject BasicType) InitialPrice float64
[_1] (BasicObject BasicType) float64
[__2] (BasicObject BasicType) InitialPrice float64

