[0] (PointerObject PointerType) r (*route.RouteServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) route.RouteServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

    --> r-tainted: read(route_db._.startstation) {1}
[0] (BasicObject BasicType) start string

    --> r-tainted: read(route_db._.endstation) {1}
[0] (BasicObject BasicType) end string

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = route, collection = route}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ primitive.A,  "$and" string,  primitive.A}

    --> r-tainted: read(route_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = route, collection = route}
     --> r-tainted: read(route_db.Route) {1}
[_1] (StructObject UserType) route route.Route struct{ID string, Stations []string, Distances []int64, StartStation string, EndStation string}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(route_db.Route) {1}
[0] (StructObject UserType) route route.Route struct{ID string, Stations []string, Distances []int64, StartStation string, EndStation string}

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

