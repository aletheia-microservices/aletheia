[0] (PointerObject PointerType) p (*price.PriceServiceImpl struct{priceDB NoSQLDatabase})
[_1] (StructObject UserType) price.PriceServiceImpl struct{priceDB NoSQLDatabase}
[__2] (FieldObject FieldType) priceDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ PriceService

    --> r-tainted: read(price_db.PriceConfig.RouteID) {1}
[0] (BasicObject BasicType) routeID string
     --> r-tainted: read(price_db.PriceConfig.RouteID) {1}
[_1] (Reference SliceType) ref <routeid []string> @ PriceService

    --> r-tainted: read(price_db.PriceConfig.TrainType) {1}
[0] (BasicObject BasicType) trainType string
     --> r-tainted: read(price_db.PriceConfig.TrainType) {1}
[_1] (Reference SliceType) ref <trainType []string> @ PriceService

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = priceConfig, collection = priceConfig}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ primitive.A,  "$and" string,  primitive.A}

    --> r-tainted: read(price_db.PriceConfig) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = priceConfig, collection = priceConfig}
     --> r-tainted: read(price_db.PriceConfig) {1}
[_1] (StructObject UserType) pc price.PriceConfig struct{ID string, TrainType string, RouteID string, BasicPriceRate float64, FirstClassPriceRate float64}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(price_db.PriceConfig) {1}
[0] (StructObject UserType) pc price.PriceConfig struct{ID string, TrainType string, RouteID string, BasicPriceRate float64, FirstClassPriceRate float64}

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

