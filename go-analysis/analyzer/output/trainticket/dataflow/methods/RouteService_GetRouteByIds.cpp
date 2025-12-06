[0] (PointerObject PointerType) r (*route.RouteServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) route.RouteServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (ArrayObject ArrayType) ids []string

[0] (ArrayObject ArrayType) routes []route.Route struct{ID string, Stations []string, Distances []int64, StartStation string, EndStation string}
[_1] (StructObject UserType) route route.Route struct{ID string, Stations []string, Distances []int64, StartStation string, EndStation string}
      --> r-tainted: read(route_db.Route) {1}
[__2] (Reference UserType) ref <route route.Route struct{ID string, Stations []string, Distances []int64, StartStation string, EndStation string}> @ GetRouteById

[0] (StructObject UserType) route route.Route struct{ID string, Stations []string, Distances []int64, StartStation string, EndStation string}
     --> r-tainted: read(route_db.Route) {1}
[_1] (Reference UserType) ref <route route.Route struct{ID string, Stations []string, Distances []int64, StartStation string, EndStation string}> @ GetRouteById

[0] (InterfaceObject UserType) _ .error
[_1] (Reference BasicType) ref <nil> @ GetRouteById

[0] (ArrayObject ArrayType) routes []route.Route struct{ID string, Stations []string, Distances []int64, StartStation string, EndStation string}
[_1] (StructObject UserType) route route.Route struct{ID string, Stations []string, Distances []int64, StartStation string, EndStation string}
      --> r-tainted: read(route_db.Route) {1}
[__2] (Reference UserType) ref <route route.Route struct{ID string, Stations []string, Distances []int64, StartStation string, EndStation string}> @ GetRouteById

