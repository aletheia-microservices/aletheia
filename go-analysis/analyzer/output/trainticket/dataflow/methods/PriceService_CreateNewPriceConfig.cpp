[0] (PointerObject PointerType) p (*price.PriceServiceImpl struct{priceDB NoSQLDatabase})
[_1] (StructObject UserType) price.PriceServiceImpl struct{priceDB NoSQLDatabase}
[__2] (FieldObject FieldType) priceDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

    --> w-tainted: write(price_db.PriceConfig) {1}
[0] (StructObject UserType) pc price.PriceConfig struct{ID string, TrainType string, RouteID string, BasicPriceRate float64, FirstClassPriceRate float64}
     --> w-tainted: write(price_db.PriceConfig.ID) {1}
[_1] (FieldObject FieldType) ID string
      --> w-tainted: write(price_db.PriceConfig.ID) {1}           --> w-tainted: write(price_db.PriceConfig.ID) {1} --> r-tainted: read(price_db._.id) {1}
[__2] (BasicObject BasicType) ID string

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = priceConfig, collection = priceConfig}

[0] (InterfaceObject UserType) err .error

[0] (StructObject UserType) _ price.PriceConfig struct{ID string, TrainType string, RouteID string, BasicPriceRate float64, FirstClassPriceRate float64}
     --> r-tainted: read(price_db.PriceConfig) {1}
[_1] (Reference UserType) ref <pc price.PriceConfig struct{ID string, TrainType string, RouteID string, BasicPriceRate float64, FirstClassPriceRate float64}> @ FindByID

[0] (InterfaceObject UserType) err .error
[_1] (Reference BasicType) ref <nil> @ FindByID

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "id" string,  string}

[0] (BasicObject BasicType) ok bool

[0] (InterfaceObject UserType) err .error

