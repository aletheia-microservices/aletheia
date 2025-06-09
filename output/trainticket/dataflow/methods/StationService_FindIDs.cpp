[0] (PointerObject PointerType) s (*station.StationServiceImpl struct{stationDB NoSQLDatabase})
[_1] (StructObject UserType) station.StationServiceImpl struct{stationDB NoSQLDatabase}
[__2] (FieldObject FieldType) stationDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (ArrayObject ArrayType) names []string

[0] (ArrayObject ArrayType) ids []string
[_1] (BasicObject BasicType) id string
      --> r-tainted: read(station_db.Station.ID) {1}
[__2] (Reference BasicType) ref <ID string> @ FindID

[0] (BasicObject BasicType) id string
     --> r-tainted: read(station_db.Station.ID) {1}
[_1] (Reference BasicType) ref <ID string> @ FindID

[0] (InterfaceObject UserType) _ .error
[_1] (Reference BasicType) ref <nil> @ FindID

[0] (ArrayObject ArrayType) ids []string
[_1] (BasicObject BasicType) id string
      --> r-tainted: read(station_db.Station.ID) {1}
[__2] (Reference BasicType) ref <ID string> @ FindID

