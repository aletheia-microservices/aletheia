[0] (PointerObject PointerType) r (*hotelreservation.RecommendationServiceImpl struct{recommendDB NoSQLDatabase, hotels map[string]hotelreservation.Hotel struct{HId string, HLat float64, HLon float64, HRate float64, HPrice float64}})
[_1] (StructObject UserType) hotelreservation.RecommendationServiceImpl struct{recommendDB NoSQLDatabase, hotels map[string]hotelreservation.Hotel struct{HId string, HLat float64, HLon float64, HRate float64, HPrice float64}}
[__2] (FieldObject FieldType) hotels map[string]hotelreservation.Hotel struct{HId string, HLat float64, HLon float64, HRate float64, HPrice float64}
[___3] (MapObject MapType) map[string]hotelreservation.Hotel struct{HId string, HLat float64, HLon float64, HRate float64, HPrice float64}
[__2] (FieldObject FieldType) recommendDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) recommendDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = recommendation-db, collection = recommendation}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) filter primitive.D

[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = recommendation-db, collection = recommendation}
[_1] (ArrayObject ArrayType) hotels []hotelreservation.Hotel struct{HId string, HLat float64, HLon float64, HRate float64, HPrice float64}

[0] (InterfaceObject UserType) err .error

[0] (ArrayObject ArrayType) hotels []hotelreservation.Hotel struct{HId string, HLat float64, HLon float64, HRate float64, HPrice float64}

