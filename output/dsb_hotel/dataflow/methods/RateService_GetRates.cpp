[0] (PointerObject PointerType) r (*hotelreservation.RateServiceImpl struct{rateCache Cache, rateDB NoSQLDatabase})
[_1] (StructObject UserType) hotelreservation.RateServiceImpl struct{rateCache Cache, rateDB NoSQLDatabase}
[__2] (FieldObject FieldType) rateCache Cache
[___3] (BlueprintBackendObject BlueprintBackendType) rateCache Cache
[__2] (FieldObject FieldType) rateDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) rateDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (ArrayObject ArrayType) hotelIDs []string

[0] (BasicObject BasicType) inDate string

[0] (BasicObject BasicType) outDate string

[0] (ArrayObject ArrayType) rate_plans []hotelreservation.RatePlan struct{HotelID string, Code string, InDate string, OutDate string, RType hotelreservation.RoomType struct{BookableRate float64, Code string, RoomDescription string, TotalRate float64, TotalRateInclusive float64}}

[0] (BasicObject BasicType) hotel_id string
[_1] (ArrayObject ArrayType) hotelIDs []string

[0] (ArrayObject ArrayType) hotel_rate_plans []hotelreservation.RatePlan struct{HotelID string, Code string, InDate string, OutDate string, RType hotelreservation.RoomType struct{BookableRate float64, Code string, RoomDescription string, TotalRate float64, TotalRateInclusive float64}}

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = rate-db, collection = inventory}

[0] (InterfaceObject UserType) err2 .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "hotelid" string,  string}

[0] (BlueprintBackendObject BlueprintBackendType) rs NoSQLCursor {database = rate-db, collection = inventory}
[_1] (ArrayObject ArrayType) hotel_rate_plans []hotelreservation.RatePlan struct{HotelID string, Code string, InDate string, OutDate string, RType hotelreservation.RoomType struct{BookableRate float64, Code string, RoomDescription string, TotalRate float64, TotalRateInclusive float64}}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

