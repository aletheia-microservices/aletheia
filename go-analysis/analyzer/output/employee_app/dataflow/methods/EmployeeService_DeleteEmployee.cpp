[0] (PointerObject PointerType) s (*employee_app.EmployeeServiceImpl struct{freelancerService employee_app.FreelancerService, fulltimeService employee_app.FulltimeService, internService employee_app.InternService, employeesDB NoSQLDatabase})
[_1] (StructObject UserType) employee_app.EmployeeServiceImpl struct{freelancerService employee_app.FreelancerService, fulltimeService employee_app.FulltimeService, internService employee_app.InternService, employeesDB NoSQLDatabase}
[__2] (FieldObject FieldType) employeesDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) employeesDB NoSQLDatabase
[__2] (FieldObject FieldType) freelancerService employee_app.FreelancerService
[___3] (ServiceObject ServiceType) freelancerService employee_app.FreelancerService
[__2] (FieldObject FieldType) fulltimeService employee_app.FulltimeService
[___3] (ServiceObject ServiceType) fulltimeService employee_app.FulltimeService
[__2] (FieldObject FieldType) internService employee_app.InternService
[___3] (ServiceObject ServiceType) internService employee_app.InternService

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

[0] (BasicObject BasicType) employeeID string
[_1] (Reference BasicType) ref <employeeID string> @ Frontend

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = employees, collection = employees}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) filter primitive.D
[_1] (StructObject StructType) struct{Key "freelancerID" string, Key "freelancerID" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "freelancerID" string
[___3] (BasicObject BasicType) "freelancerID" string
[__2] (FieldObject FieldType) Value string
[___3] (BasicObject BasicType) employeeID string
[____4] (Reference BasicType) ref <employeeID string> @ Frontend

[0] (InterfaceObject UserType) err .error

