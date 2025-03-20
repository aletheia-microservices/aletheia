[0] (PointerObject PointerType) s (*employee_app.FulltimeServiceImpl struct{fulltimesDB NoSQLDatabase})
[_1] (StructObject UserType) employee_app.FulltimeServiceImpl struct{fulltimesDB NoSQLDatabase}
[__2] (FieldObject FieldType) fulltimesDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) fulltimesDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ EmployeeService
[__2] (Reference UserType) ref <ctx context.Context> @ Frontend
[_1] (Reference UserType) ref <ctx context.Context> @ EmployeeService
[__2] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID, employees_db.Employee.EmployeeID) {2}
[0] (BasicObject BasicType) employeeID string
     --> w-tainted: write(employees_db.Employee.EmployeeID, fulltimes_db.Fulltime.EmployeeID) {2}
[_1] (Reference BasicType) ref <employeeID string> @ EmployeeService
      --> w-tainted: write(employees_db.Employee.EmployeeID, fulltimes_db.Fulltime.EmployeeID) {2}
[__2] (Reference BasicType) ref <employeeID string> @ Frontend
     --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID) {1}
[_1] (Reference BasicType) ref <employeeID string> @ EmployeeService
      --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID) {1}
[__2] (Reference BasicType) ref <employeeID string> @ Frontend

    --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[0] (BasicObject BasicType) fulltimeID string
     --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[_1] (Reference BasicType) ref <fulltimeID string> @ EmployeeService
      --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[__2] (Reference BasicType) ref <fulltimeID string> @ Frontend
     --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[_1] (Reference BasicType) ref <fulltimeID string> @ EmployeeService
      --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[__2] (Reference BasicType) ref <fulltimeID string> @ Frontend

    --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[0] (BasicObject BasicType) salary string
     --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[_1] (Reference BasicType) ref <salary string> @ EmployeeService
      --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[__2] (Reference BasicType) ref <salary string> @ Frontend
     --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[_1] (Reference BasicType) ref <salary string> @ EmployeeService
      --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[__2] (Reference BasicType) ref <salary string> @ Frontend

    --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[0] (BasicObject BasicType) position string
     --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[_1] (Reference BasicType) ref <position string> @ EmployeeService
      --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[__2] (Reference BasicType) ref <position string> @ Frontend
     --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[_1] (Reference BasicType) ref <position string> @ EmployeeService
      --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[__2] (Reference BasicType) ref <position string> @ Frontend

    --> w-tainted: write(fulltimes_db.Fulltime) {1}
[0] (StructObject UserType) employee employee_app.Fulltime struct{EmployeeID string, FulltimeID string, Salary string, Position string}
     --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID) {1}
[_1] (FieldObject FieldType) EmployeeID string
      --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID, employees_db.Employee.EmployeeID) {2}
[__2] (BasicObject BasicType) employeeID string
       --> w-tainted: write(employees_db.Employee.EmployeeID, fulltimes_db.Fulltime.EmployeeID) {2}
[___3] (Reference BasicType) ref <employeeID string> @ EmployeeService
        --> w-tainted: write(employees_db.Employee.EmployeeID, fulltimes_db.Fulltime.EmployeeID) {2}
[____4] (Reference BasicType) ref <employeeID string> @ Frontend
       --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID) {1}
[___3] (Reference BasicType) ref <employeeID string> @ EmployeeService
        --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID) {1}
[____4] (Reference BasicType) ref <employeeID string> @ Frontend
     --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[_1] (FieldObject FieldType) FulltimeID string
      --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[__2] (BasicObject BasicType) fulltimeID string
       --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[___3] (Reference BasicType) ref <fulltimeID string> @ EmployeeService
        --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[____4] (Reference BasicType) ref <fulltimeID string> @ Frontend
       --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[___3] (Reference BasicType) ref <fulltimeID string> @ EmployeeService
        --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[____4] (Reference BasicType) ref <fulltimeID string> @ Frontend
     --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[_1] (FieldObject FieldType) Position string
      --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[__2] (BasicObject BasicType) position string
       --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[___3] (Reference BasicType) ref <position string> @ EmployeeService
        --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[____4] (Reference BasicType) ref <position string> @ Frontend
       --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[___3] (Reference BasicType) ref <position string> @ EmployeeService
        --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[____4] (Reference BasicType) ref <position string> @ Frontend
     --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[_1] (FieldObject FieldType) Salary string
      --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[__2] (BasicObject BasicType) salary string
       --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[___3] (Reference BasicType) ref <salary string> @ EmployeeService
        --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[____4] (Reference BasicType) ref <salary string> @ Frontend
       --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[___3] (Reference BasicType) ref <salary string> @ EmployeeService
        --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[____4] (Reference BasicType) ref <salary string> @ Frontend

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = fulltimes, collection = fulltimes}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

