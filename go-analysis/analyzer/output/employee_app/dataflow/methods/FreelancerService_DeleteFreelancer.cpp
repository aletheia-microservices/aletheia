[0] (PointerObject PointerType) s (*employee_app.FreelancerServiceImpl struct{freelancersDB NoSQLDatabase})
[_1] (StructObject UserType) employee_app.FreelancerServiceImpl struct{freelancersDB NoSQLDatabase}
[__2] (FieldObject FieldType) freelancersDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) freelancersDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ EmployeeService
[__2] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID) {1}
[0] (BasicObject BasicType) employeeID string
     --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID) {1}
[_1] (Reference BasicType) ref <employeeID string> @ EmployeeService
      --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID) {1}
[__2] (Reference BasicType) ref <employeeID string> @ Frontend

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = freelancers, collection = freelancers}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) filter primitive.D
[_1] (StructObject StructType) struct{Key "employeeID" string, Key "employeeID" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "employeeID" string
[___3] (BasicObject BasicType) "employeeID" string
[__2] (FieldObject FieldType) Value string
       --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID) {1}
[___3] (BasicObject BasicType) employeeID string
        --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID) {1}
[____4] (Reference BasicType) ref <employeeID string> @ EmployeeService
         --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID) {1}
[_____5] (Reference BasicType) ref <employeeID string> @ Frontend

[0] (InterfaceObject UserType) err .error

