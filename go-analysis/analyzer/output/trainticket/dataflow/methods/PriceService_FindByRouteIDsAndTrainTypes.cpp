[0] (PointerObject PointerType) p (*price.PriceServiceImpl struct{priceDB NoSQLDatabase})
[_1] (StructObject UserType) price.PriceServiceImpl struct{priceDB NoSQLDatabase}
[__2] (FieldObject FieldType) priceDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (ArrayObject ArrayType) rtsAndTypes []string

[0] (MapObject MapType) res map[string]price.PriceConfig struct{ID string, TrainType string, RouteID string, BasicPriceRate float64, FirstClassPriceRate float64}

[0] (BasicObject BasicType) rt string
[_1] (BasicObject BasicType) * string

[0] (SliceObject SliceType) pieces []string
[_1] (BasicObject BasicType) rt string
[__2] (BasicObject BasicType) * string
[_1] (BasicObject BasicType) ":" string
[_1] (BasicObject BasicType) rt string
[__2] (BasicObject BasicType) * string
[_1] (BasicObject BasicType) ":" string

    --> r-tainted: read(price_db.PriceConfig.RouteID) {1}
[0] (SliceObject SliceType) routeid []string

    --> r-tainted: read(price_db.PriceConfig.TrainType) {1}
[0] (SliceObject SliceType) trainType []string

[0] (StructObject UserType) pc price.PriceConfig struct{ID string, TrainType string, RouteID string, BasicPriceRate float64, FirstClassPriceRate float64}
     --> r-tainted: read(price_db.PriceConfig) {1}
[_1] (Reference UserType) ref <pc price.PriceConfig struct{ID string, TrainType string, RouteID string, BasicPriceRate float64, FirstClassPriceRate float64}> @ FindByRouteIDAndTrainType

[0] (InterfaceObject UserType) _ .error
[_1] (Reference BasicType) ref <nil> @ FindByRouteIDAndTrainType

