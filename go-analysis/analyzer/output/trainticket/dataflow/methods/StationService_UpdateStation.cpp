[0] (PointerObject PointerType) s (*station.StationServiceImpl struct{stationDB NoSQLDatabase})
[_1] (StructObject UserType) station.StationServiceImpl struct{stationDB NoSQLDatabase}
[__2] (FieldObject FieldType) stationDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (StructObject UserType) station station.Station struct{ID string, Name string, StayTime int64}
[_1] (FieldObject FieldType) ID string
[__2] (BasicObject BasicType) ID string

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = station, collection = station}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "id" string,  string}

