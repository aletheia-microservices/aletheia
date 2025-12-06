[0] (PointerObject PointerType) p (*price.PriceServiceImpl struct{priceDB NoSQLDatabase})
[_1] (StructObject UserType) price.PriceServiceImpl struct{priceDB NoSQLDatabase}
[__2] (FieldObject FieldType) priceDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = priceConfig, collection = priceConfig}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(price_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = priceConfig, collection = priceConfig}
     --> r-tainted: read(price_db.[]price.PriceConfig struct{ID string, TrainType string, RouteID string, BasicPriceRate float64, FirstClassPriceRate float64}) {1}
[_1] (ArrayObject ArrayType) pcs []price.PriceConfig struct{ID string, TrainType string, RouteID string, BasicPriceRate float64, FirstClassPriceRate float64}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(price_db.[]price.PriceConfig struct{ID string, TrainType string, RouteID string, BasicPriceRate float64, FirstClassPriceRate float64}) {1}
[0] (ArrayObject ArrayType) pcs []price.PriceConfig struct{ID string, TrainType string, RouteID string, BasicPriceRate float64, FirstClassPriceRate float64}

[0] (InterfaceObject UserType) err .error

