[0] (PointerObject PointerType) r (*route.RouteServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) route.RouteServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (StructObject UserType) info route.RouteInfo struct{ID string, StartStation string, EndStation string, StationList string, DistanceList string}
[_1] (FieldObject FieldType) EndStation string
      --> w-tainted: write(route_db.Route.EndStation) {1}
[__2] (BasicObject BasicType) EndStation string
[_1] (FieldObject FieldType) ID string
      --> r-tainted: read(route_db.Route.ID) {1}
[__2] (BasicObject BasicType) ID string
[_1] (FieldObject FieldType) StartStation string
      --> w-tainted: write(route_db.Route.StartStation) {1}
[__2] (BasicObject BasicType) StartStation string
[_1] (FieldObject FieldType) StationList string
      --> w-tainted: write(route_db.Route.Stations) {1}
[__2] (BasicObject BasicType) StationList string

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = route, collection = route}

[0] (InterfaceObject UserType) err .error

    --> w-tainted: write(route_db.Route.Distances) {1}
[0] (ArrayObject SliceType) distances []int64

    --> w-tainted: write(route_db.Route.Stations) {1}
[0] (SliceObject SliceType) stations []string
     --> w-tainted: write(route_db.Route.Stations) {1}
[_1] (BasicObject BasicType) StationList string
     --> w-tainted: write(route_db.Route.Stations) {1}
[_1] (BasicObject BasicType) "," string
     --> w-tainted: write(route_db.Route.Stations) {1}
[_1] (BasicObject BasicType) StationList string
     --> w-tainted: write(route_db.Route.Stations) {1}
[_1] (BasicObject BasicType) "," string

    --> w-tainted: write(route_db.Route) {1}
[0] (StructObject UserType) route route.Route struct{ID string, Stations []string, Distances []int64, StartStation string, EndStation string}
     --> w-tainted: write(route_db.Route.Distances) {1}
[_1] (FieldObject FieldType) Distances []int64
      --> w-tainted: write(route_db.Route.Distances) {1}
[__2] (ArrayObject SliceType) distances []int64
     --> w-tainted: write(route_db.Route.EndStation) {1}
[_1] (FieldObject FieldType) EndStation string
      --> w-tainted: write(route_db.Route.EndStation) {1}
[__2] (BasicObject BasicType) EndStation string
     --> w-tainted: write(route_db.Route.ID) {1}
[_1] (FieldObject FieldType) ID string
      --> r-tainted: read(route_db.Route.ID) {1}
[__2] (BasicObject BasicType) ID string
     --> w-tainted: write(route_db.Route.StartStation) {1}
[_1] (FieldObject FieldType) StartStation string
      --> w-tainted: write(route_db.Route.StartStation) {1}
[__2] (BasicObject BasicType) StartStation string
     --> w-tainted: write(route_db.Route.Stations) {1}
[_1] (FieldObject FieldType) Stations []string
      --> w-tainted: write(route_db.Route.Stations) {1}
[__2] (SliceObject SliceType) stations []string
       --> w-tainted: write(route_db.Route.Stations) {1}
[___3] (BasicObject BasicType) StationList string
       --> w-tainted: write(route_db.Route.Stations) {1}
[___3] (BasicObject BasicType) "," string
       --> w-tainted: write(route_db.Route.Stations) {1}
[___3] (BasicObject BasicType) StationList string
       --> w-tainted: write(route_db.Route.Stations) {1}
[___3] (BasicObject BasicType) "," string

[0] (BasicObject BasicType) MAXIDLEN 32 int

[0] (BasicObject BasicType) old_exists false bool

    --> r-tainted: read(route_db.Route) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = route, collection = route}
     --> r-tainted: read(route_db.Route, route_db.Route.ID) {2}
[_1] (StructObject UserType) old_route route.Route struct{ID string, Stations []string, Distances []int64, StartStation string, EndStation string}
      --> r-tainted: read(route_db.Route.ID) {1}
[__2] (FieldObject FieldType) ID string
       --> r-tainted: read(route_db.Route.ID) {1}
[___3] (BasicObject BasicType) ID string

[0] (InterfaceObject UserType) err .error

