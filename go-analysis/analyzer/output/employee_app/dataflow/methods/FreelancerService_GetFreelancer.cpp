[0] (PointerObject PointerType) s (*employee_app.FreelancerServiceImpl struct{freelancersDB NoSQLDatabase})
[_1] (StructObject UserType) employee_app.FreelancerServiceImpl struct{freelancersDB NoSQLDatabase}
[__2] (FieldObject FieldType) freelancersDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) freelancersDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) freelancerID string

[0] (StructObject UserType) freelancer employee_app.Freelancer struct{EmployeeID string, FreelancerID string, Rate string, Terms string}

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = freelancers, collection = freelancers}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "freelancerID" string, Key "freelancerID" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "freelancerID" string
[___3] (BasicObject BasicType) "freelancerID" string
[__2] (FieldObject FieldType) Value string
[___3] (BasicObject BasicType) freelancerID string

[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = freelancers, collection = freelancers}
[_1] (StructObject UserType) freelancer employee_app.Freelancer struct{EmployeeID string, FreelancerID string, Rate string, Terms string}

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) res bool

[0] (InterfaceObject UserType) err .error

