[0] (PointerObject PointerType) s (*employee_app.FulltimeServiceImpl struct{fulltimesDB NoSQLDatabase})
[_1] (StructObject UserType) employee_app.FulltimeServiceImpl struct{fulltimesDB NoSQLDatabase}
[__2] (FieldObject FieldType) fulltimesDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) fulltimesDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) fulltimeID string

[0] (StructObject UserType) fulltime employee_app.Fulltime struct{EmployeeID string, FulltimeID string, Salary string, Position string}

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = fulltimes, collection = fulltimes}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "fulltimeID" string, Key "fulltimeID" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "fulltimeID" string
[___3] (BasicObject BasicType) "fulltimeID" string
[__2] (FieldObject FieldType) Value string
[___3] (BasicObject BasicType) fulltimeID string

[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = fulltimes, collection = fulltimes}
[_1] (StructObject UserType) fulltime employee_app.Fulltime struct{EmployeeID string, FulltimeID string, Salary string, Position string}

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) res bool

[0] (InterfaceObject UserType) err .error

