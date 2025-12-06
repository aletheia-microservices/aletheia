[0] (PointerObject PointerType) g (*hotelreservation.GeoServiceImpl struct{geoDB NoSQLDatabase})
[_1] (StructObject UserType) hotelreservation.GeoServiceImpl struct{geoDB NoSQLDatabase}
[__2] (FieldObject FieldType) geoDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) geoDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = geo-db, collection = geo}

[0] (InterfaceObject UserType) err .error

[0] (ArrayObject ArrayType) points []hotelreservation.Point struct{Pid string, Plat float64, Plon float64}

[0] (SliceObject UserType) filter primitive.D

[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = geo-db, collection = geo}
[_1] (ArrayObject ArrayType) points []hotelreservation.Point struct{Pid string, Plat float64, Plon float64}

[0] (InterfaceObject UserType) err .error

