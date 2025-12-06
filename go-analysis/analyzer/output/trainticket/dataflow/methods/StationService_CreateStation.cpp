[0] (PointerObject PointerType) s (*station.StationServiceImpl struct{stationDB NoSQLDatabase})
[_1] (StructObject UserType) station.StationServiceImpl struct{stationDB NoSQLDatabase}
[__2] (FieldObject FieldType) stationDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

    --> w-tainted: write(station_db.Station) {1}
[0] (StructObject UserType) station station.Station struct{ID string, Name string, StayTime int64}
     --> w-tainted: write(station_db.Station.ID) {1}
[_1] (FieldObject FieldType) ID string
      --> w-tainted: write(station_db.Station.ID) {1}           --> w-tainted: write(station_db.Station.ID) {1} --> r-tainted: read(station_db._.id) {1}
[__2] (BasicObject BasicType) ID string

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = station, collection = station}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "id" string,  string}

    --> r-tainted: read(station_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = station, collection = station}
     --> r-tainted: read(station_db.Station) {1}
[_1] (StructObject UserType) st station.Station struct{ID string, Name string, StayTime int64}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(station_db.Station) {1}
[0] (StructObject UserType) st station.Station struct{ID string, Name string, StayTime int64}

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

