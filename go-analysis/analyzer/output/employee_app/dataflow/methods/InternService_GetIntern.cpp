[0] (PointerObject PointerType) s (*employee_app.InternServiceImpl struct{internsDB NoSQLDatabase})
[_1] (StructObject UserType) employee_app.InternServiceImpl struct{internsDB NoSQLDatabase}
[__2] (FieldObject FieldType) internsDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) internsDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) internID string

[0] (StructObject UserType) intern employee_app.Intern struct{EmployeeID string, InternID string, Stipend string, Duration string, Mentor string}

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = interns, collection = interns}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "internID" string, Key "internID" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "internID" string
[___3] (BasicObject BasicType) "internID" string
[__2] (FieldObject FieldType) Value string
[___3] (BasicObject BasicType) internID string

[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = interns, collection = interns}
[_1] (StructObject UserType) intern employee_app.Intern struct{EmployeeID string, InternID string, Stipend string, Duration string, Mentor string}

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) res bool

[0] (InterfaceObject UserType) err .error

