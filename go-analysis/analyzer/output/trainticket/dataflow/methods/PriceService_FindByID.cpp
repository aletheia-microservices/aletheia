[0] (PointerObject PointerType) p (*price.PriceServiceImpl struct{priceDB NoSQLDatabase})
[_1] (StructObject UserType) price.PriceServiceImpl struct{priceDB NoSQLDatabase}
[__2] (FieldObject FieldType) priceDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ PriceService

    --> w-tainted: write(price_db.PriceConfig.ID) {1}       --> w-tainted: write(price_db.PriceConfig.ID) {1} --> r-tainted: read(price_db._.id) {1}
[0] (BasicObject BasicType) id string
     --> w-tainted: write(price_db.PriceConfig.ID) {1}         --> w-tainted: write(price_db.PriceConfig.ID) {1} --> r-tainted: read(price_db._.id) {1}
[_1] (Reference BasicType) ref <ID string> @ PriceService

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = priceConfig, collection = priceConfig}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "id" string,  string}

    --> r-tainted: read(price_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = priceConfig, collection = priceConfig}
     --> r-tainted: read(price_db.PriceConfig) {1}
[_1] (StructObject UserType) pc price.PriceConfig struct{ID string, TrainType string, RouteID string, BasicPriceRate float64, FirstClassPriceRate float64}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(price_db.PriceConfig) {1}
[0] (StructObject UserType) pc price.PriceConfig struct{ID string, TrainType string, RouteID string, BasicPriceRate float64, FirstClassPriceRate float64}

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

