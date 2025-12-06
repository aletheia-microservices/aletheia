[0] (PointerObject PointerType) s (*station.StationServiceImpl struct{stationDB NoSQLDatabase})
[_1] (StructObject UserType) station.StationServiceImpl struct{stationDB NoSQLDatabase}
[__2] (FieldObject FieldType) stationDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (ArrayObject ArrayType) ids []string

[0] (ArrayObject ArrayType) stations []station.Station struct{ID string, Name string, StayTime int64}
[_1] (StructObject UserType) st station.Station struct{ID string, Name string, StayTime int64}
      --> r-tainted: read(station_db.Station) {1}
[__2] (Reference UserType) ref <st station.Station struct{ID string, Name string, StayTime int64}> @ FindByID

[0] (StructObject UserType) st station.Station struct{ID string, Name string, StayTime int64}
     --> r-tainted: read(station_db.Station) {1}
[_1] (Reference UserType) ref <st station.Station struct{ID string, Name string, StayTime int64}> @ FindByID

[0] (InterfaceObject UserType) _ .error
[_1] (Reference BasicType) ref <nil> @ FindByID

[0] (ArrayObject ArrayType) stations []station.Station struct{ID string, Name string, StayTime int64}
[_1] (StructObject UserType) st station.Station struct{ID string, Name string, StayTime int64}
      --> r-tainted: read(station_db.Station) {1}
[__2] (Reference UserType) ref <st station.Station struct{ID string, Name string, StayTime int64}> @ FindByID

